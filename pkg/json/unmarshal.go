package json

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-json/internal/fastparser"
)

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
//
// This function uses a high-performance fast path that bypasses AST construction for
// optimal performance. If you need the AST for advanced features (JSONPath, etc.), use
// Parse() followed by NodeToInterface() or manual AST traversal.
//
// Unmarshal uses the inverse of the encodings that Marshal uses, allocating maps, slices,
// and pointers as necessary, with the following additional rules:
//
// To unmarshal JSON into a pointer, Unmarshal first handles the case of the JSON being
// the JSON literal null. In that case, Unmarshal sets the pointer to nil. Otherwise,
// Unmarshal unmarshals the JSON into the value pointed at by the pointer. If the pointer
// is nil, Unmarshal allocates a new value for it to point to.
//
// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys
// used by Marshal (either the struct field name or its tag), preferring an exact match
// but also accepting a case-insensitive match. Unmarshal will only set exported fields.
//
// To unmarshal JSON into an interface value, Unmarshal stores one of these in the interface value:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON strings
//	[]interface{}, for JSON arrays
//	map[string]interface{}, for JSON objects
//	nil for JSON null
//
// If the JSON is not valid, Unmarshal returns a parse error.
func Unmarshal(data []byte, v interface{}) error {
	// Fast path: Direct parsing without AST construction (4-5x faster)
	return fastparser.Unmarshal(data, v)
}

// UnmarshalWithAST parses the JSON-encoded data into an AST first, then unmarshals into v.
// This is the slower path but allows access to the AST for advanced features.
// Most users should use Unmarshal() instead for better performance.
func UnmarshalWithAST(data []byte, v interface{}) error {
	// Parse JSON into AST
	node, err := Parse(string(data))
	if err != nil {
		return err
	}

	return unmarshalFromNode(node, v)
}

// Unmarshaler is the interface implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// unmarshalFromNode unmarshals an AST node into a Go value
// This is used by both Unmarshal and Decoder.Decode
func unmarshalFromNode(node ast.SchemaNode, v interface{}) error {
	// Use reflection to populate v from AST
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || v == nil {
		return errors.New("json: Unmarshal(nil)")
	}

	if rv.Kind() != reflect.Ptr {
		return errors.New("json: Unmarshal(non-pointer " + rv.Type().String() + ")")
	}

	if rv.IsNil() {
		return errors.New("json: Unmarshal(nil " + rv.Type().String() + ")")
	}

	// Check if type implements Unmarshaler interface
	if rv.Type().Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		// Render node back to JSON
		jsonBytes, err := Render(node)
		if err != nil {
			return err
		}
		unmarshaler := rv.Interface().(Unmarshaler)
		return unmarshaler.UnmarshalJSON(jsonBytes)
	}

	return unmarshalValue(node, rv.Elem())
}

// unmarshalValue unmarshals an AST node into a reflect.Value
func unmarshalValue(node ast.SchemaNode, rv reflect.Value) error {
	// Handle null
	if lit, ok := node.(*ast.LiteralNode); ok && lit.Value() == nil {
		// Set to zero value (nil for pointers, zero for values)
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}

	// Handle interface{} specially
	if rv.Kind() == reflect.Interface && rv.NumMethod() == 0 {
		val := nodeToInterface(node)
		rv.Set(reflect.ValueOf(val))
		return nil
	}

	// Handle pointers
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return unmarshalValue(node, rv.Elem())
	}

	switch node.Type() {
	case ast.NodeTypeLiteral:
		return unmarshalLiteral(node.(*ast.LiteralNode), rv)
	case ast.NodeTypeObject:
		return unmarshalObject(node.(*ast.ObjectNode), rv)
	case ast.NodeTypeArrayData:
		return unmarshalArrayData(node.(*ast.ArrayDataNode), rv)
	default:
		return fmt.Errorf("json: unsupported node type %s", node.Type())
	}
}

// unmarshalLiteral unmarshals a literal node into a reflect.Value
func unmarshalLiteral(node *ast.LiteralNode, rv reflect.Value) error {
	val := node.Value()

	switch rv.Kind() {
	case reflect.String:
		if s, ok := val.(string); ok {
			rv.SetString(s)
			return nil
		}
		return fmt.Errorf("json: cannot unmarshal %T into Go value of type string", val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := val.(type) {
		case int64:
			if rv.OverflowInt(v) {
				return fmt.Errorf("json: value %d overflows %s", v, rv.Type())
			}
			rv.SetInt(v)
			return nil
		case float64:
			// Allow conversion from float to int if it's a whole number
			if v == float64(int64(v)) {
				i := int64(v)
				if rv.OverflowInt(i) {
					return fmt.Errorf("json: value %v overflows %s", v, rv.Type())
				}
				rv.SetInt(i)
				return nil
			}
			return fmt.Errorf("json: cannot unmarshal number %v into Go value of type %s", v, rv.Type())
		}
		return fmt.Errorf("json: cannot unmarshal %T into Go value of type %s", val, rv.Type())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := val.(type) {
		case int64:
			if v < 0 || rv.OverflowUint(uint64(v)) {
				return fmt.Errorf("json: value %d overflows %s", v, rv.Type())
			}
			rv.SetUint(uint64(v))
			return nil
		case float64:
			if v < 0 || v != float64(uint64(v)) {
				return fmt.Errorf("json: cannot unmarshal number %v into Go value of type %s", v, rv.Type())
			}
			u := uint64(v)
			if rv.OverflowUint(u) {
				return fmt.Errorf("json: value %v overflows %s", v, rv.Type())
			}
			rv.SetUint(u)
			return nil
		}
		return fmt.Errorf("json: cannot unmarshal %T into Go value of type %s", val, rv.Type())

	case reflect.Float32, reflect.Float64:
		switch v := val.(type) {
		case float64:
			if rv.OverflowFloat(v) {
				return fmt.Errorf("json: value %v overflows %s", v, rv.Type())
			}
			rv.SetFloat(v)
			return nil
		case int64:
			f := float64(v)
			if rv.OverflowFloat(f) {
				return fmt.Errorf("json: value %v overflows %s", v, rv.Type())
			}
			rv.SetFloat(f)
			return nil
		}
		return fmt.Errorf("json: cannot unmarshal %T into Go value of type %s", val, rv.Type())

	case reflect.Bool:
		if b, ok := val.(bool); ok {
			rv.SetBool(b)
			return nil
		}
		return fmt.Errorf("json: cannot unmarshal %T into Go value of type bool", val)

	default:
		return fmt.Errorf("json: unsupported type %s", rv.Type())
	}
}

// unmarshalObject unmarshals an object node into a reflect.Value (struct, map, or slice)
func unmarshalObject(node *ast.ObjectNode, rv reflect.Value) error {
	props := node.Properties()

	// Check if this is an array (all keys are numeric strings "0", "1", "2", etc.)
	if isArray(props) {
		return unmarshalArray(node, rv)
	}

	switch rv.Kind() {
	case reflect.Struct:
		return unmarshalStruct(node, rv)
	case reflect.Map:
		return unmarshalMap(node, rv)
	case reflect.Slice:
		return unmarshalArray(node, rv)
	default:
		return fmt.Errorf("json: cannot unmarshal object into Go value of type %s", rv.Type())
	}
}

// isArray checks if the object node represents a JSON array (numeric string keys)
func isArray(props map[string]ast.SchemaNode) bool {
	if len(props) == 0 {
		return false
	}

	for i := 0; i < len(props); i++ {
		if _, ok := props[strconv.Itoa(i)]; !ok {
			return false
		}
	}
	return true
}

// unmarshalStruct unmarshals an object node into a struct
func unmarshalStruct(node *ast.ObjectNode, rv reflect.Value) error {
	props := node.Properties()
	structType := rv.Type()

	// Build a map of JSON field names to struct field indices
	fieldMap := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" { // Skip unexported fields
			continue
		}

		info := getFieldInfo(field)
		if info.skip {
			continue
		}

		fieldMap[info.name] = i
	}

	// Set struct fields from JSON properties
	for jsonName, propNode := range props {
		if fieldIdx, ok := fieldMap[jsonName]; ok {
			fieldVal := rv.Field(fieldIdx)
			if err := unmarshalValue(propNode, fieldVal); err != nil {
				return err
			}
		}
	}

	return nil
}

// unmarshalMap unmarshals an object node into a map
func unmarshalMap(node *ast.ObjectNode, rv reflect.Value) error {
	props := node.Properties()
	mapType := rv.Type()

	// Create the map if nil
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(mapType))
	}

	keyType := mapType.Key()
	valueType := mapType.Elem()

	// Only support string keys
	if keyType.Kind() != reflect.String {
		return fmt.Errorf("json: unsupported map key type %s", keyType)
	}

	for key, propNode := range props {
		// Create a new value of the map's value type
		elemVal := reflect.New(valueType).Elem()

		// Unmarshal the property into the value
		if err := unmarshalValue(propNode, elemVal); err != nil {
			return err
		}

		// Set the map entry
		rv.SetMapIndex(reflect.ValueOf(key), elemVal)
	}

	return nil
}

// unmarshalArray unmarshals an array (object with numeric keys) into a slice
func unmarshalArray(node *ast.ObjectNode, rv reflect.Value) error {
	props := node.Properties()

	// Determine array length
	arrayLen := len(props)

	switch rv.Kind() {
	case reflect.Slice:
		// Create a new slice of the correct length
		sliceType := rv.Type()
		slice := reflect.MakeSlice(sliceType, arrayLen, arrayLen)

		// Unmarshal each element
		for i := 0; i < arrayLen; i++ {
			key := strconv.Itoa(i)
			if propNode, ok := props[key]; ok {
				elemVal := slice.Index(i)
				if err := unmarshalValue(propNode, elemVal); err != nil {
					return err
				}
			}
		}

		rv.Set(slice)
		return nil

	case reflect.Array:
		if arrayLen > rv.Len() {
			return fmt.Errorf("json: array length %d exceeds target array length %d", arrayLen, rv.Len())
		}

		// Unmarshal each element
		for i := 0; i < arrayLen; i++ {
			key := strconv.Itoa(i)
			if propNode, ok := props[key]; ok {
				elemVal := rv.Index(i)
				if err := unmarshalValue(propNode, elemVal); err != nil {
					return err
				}
			}
		}

		return nil

	default:
		return fmt.Errorf("json: cannot unmarshal array into Go value of type %s", rv.Type())
	}
}

// unmarshalArrayData unmarshals an ArrayDataNode into a slice or array.
func unmarshalArrayData(node *ast.ArrayDataNode, rv reflect.Value) error {
	elements := node.Elements()
	arrayLen := len(elements)

	switch rv.Kind() {
	case reflect.Slice:
		// Create a new slice of the correct length
		sliceType := rv.Type()
		slice := reflect.MakeSlice(sliceType, arrayLen, arrayLen)

		// Unmarshal each element
		for i, elem := range elements {
			elemVal := slice.Index(i)
			if err := unmarshalValue(elem, elemVal); err != nil {
				return err
			}
		}

		rv.Set(slice)
		return nil

	case reflect.Array:
		if arrayLen > rv.Len() {
			return fmt.Errorf("json: array length %d exceeds target array length %d", arrayLen, rv.Len())
		}

		// Unmarshal each element
		for i, elem := range elements {
			elemVal := rv.Index(i)
			if err := unmarshalValue(elem, elemVal); err != nil {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("json: cannot unmarshal array into Go value of type %s", rv.Type())
	}
}

// nodeToInterface is a wrapper around NodeToInterface for backward compatibility.
// New code should use NodeToInterface from convert.go directly.
func nodeToInterface(node ast.SchemaNode) interface{} {
	return NodeToInterface(node)
}
