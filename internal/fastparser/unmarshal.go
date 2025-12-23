package fastparser

import (
	"errors"
	"fmt"
	"reflect"
)

// Unmarshaler is the interface implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// Unmarshal parses JSON and unmarshals it into the value pointed to by v.
// This is the fast path that bypasses AST construction.
func Unmarshal(data []byte, v interface{}) error {
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
		unmarshaler := rv.Interface().(Unmarshaler)
		return unmarshaler.UnmarshalJSON(data)
	}

	p := NewParser(data)
	return p.unmarshalValue(rv.Elem())
}

// unmarshalValue unmarshals JSON into a reflect.Value.
func (p *Parser) unmarshalValue(rv reflect.Value) error {
	p.skipWhitespace()
	if p.pos >= p.length {
		return errors.New("unexpected end of JSON input")
	}

	c := p.data[p.pos]

	// Handle null
	if c == 'n' {
		if err := p.expectLiteral("null"); err != nil {
			return err
		}
		// Set to zero value
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}

	// Handle interface{} specially - parse to native Go types
	if rv.Kind() == reflect.Interface && rv.NumMethod() == 0 {
		value, err := p.parseValue()
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(value))
		return nil
	}

	// Handle pointers
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return p.unmarshalValue(rv.Elem())
	}

	// Route based on JSON type
	switch c {
	case '{':
		return p.unmarshalObject(rv)
	case '[':
		return p.unmarshalArray(rv)
	case '"':
		return p.unmarshalString(rv)
	case 't', 'f':
		return p.unmarshalBool(rv)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.unmarshalNumber(rv)
	default:
		return fmt.Errorf("unexpected character '%c' at position %d", c, p.pos)
	}
}

// unmarshalObject unmarshals a JSON object.
func (p *Parser) unmarshalObject(rv reflect.Value) error {
	if p.pos >= p.length || p.data[p.pos] != '{' {
		return errors.New("expected '{'")
	}
	p.pos++ // skip '{'

	switch rv.Kind() {
	case reflect.Struct:
		return p.unmarshalStruct(rv)
	case reflect.Map:
		return p.unmarshalMap(rv)
	default:
		return fmt.Errorf("json: cannot unmarshal object into Go value of type %s", rv.Type())
	}
}

// unmarshalStruct unmarshals a JSON object into a struct.
func (p *Parser) unmarshalStruct(rv reflect.Value) error {
	structType := rv.Type()

	// Build field map
	fieldMap := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" { // Skip unexported fields
			continue
		}

		// Check JSON tag
		tag := field.Tag.Get("json")
		if tag == "-" {
			// Field should be ignored
			continue
		}

		// Get JSON name from tag or use field name
		jsonName := field.Name
		if tag != "" {
			jsonName = tag
			// Handle "name,omitempty" format
			for idx := 0; idx < len(tag); idx++ {
				if tag[idx] == ',' {
					jsonName = tag[:idx]
					break
				}
			}
		}

		fieldMap[jsonName] = i
	}

	p.skipWhitespace()

	// Handle empty object
	if p.pos < p.length && p.data[p.pos] == '}' {
		p.pos++
		return nil
	}

	for {
		p.skipWhitespace()

		// Parse key
		if p.pos >= p.length || p.data[p.pos] != '"' {
			return errors.New("expected string key in object")
		}

		key, err := p.parseString()
		if err != nil {
			return err
		}

		p.skipWhitespace()

		// Expect ':'
		if p.pos >= p.length || p.data[p.pos] != ':' {
			return errors.New("expected ':' after object key")
		}
		p.pos++

		p.skipWhitespace()

		// Unmarshal value into struct field if it exists
		if fieldIdx, ok := fieldMap[key]; ok {
			fieldVal := rv.Field(fieldIdx)
			if err := p.unmarshalValue(fieldVal); err != nil {
				return err
			}
		} else {
			// Skip unknown field
			if err := p.skipValue(); err != nil {
				return err
			}
		}

		p.skipWhitespace()

		// Check for more entries or end of object
		if p.pos >= p.length {
			return errors.New("unexpected end of JSON input in object")
		}

		if p.data[p.pos] == '}' {
			p.pos++
			return nil
		}

		if p.data[p.pos] != ',' {
			return fmt.Errorf("expected ',' or '}' in object at position %d", p.pos)
		}
		p.pos++
	}
}

// unmarshalMap unmarshals a JSON object into a map.
func (p *Parser) unmarshalMap(rv reflect.Value) error {
	mapType := rv.Type()

	// Only support string keys
	if mapType.Key().Kind() != reflect.String {
		return fmt.Errorf("json: unsupported map key type %s", mapType.Key())
	}

	// Create the map if nil
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(mapType))
	}

	valueType := mapType.Elem()

	p.skipWhitespace()

	// Handle empty object
	if p.pos < p.length && p.data[p.pos] == '}' {
		p.pos++
		return nil
	}

	for {
		p.skipWhitespace()

		// Parse key
		if p.pos >= p.length || p.data[p.pos] != '"' {
			return errors.New("expected string key in object")
		}

		key, err := p.parseString()
		if err != nil {
			return err
		}

		p.skipWhitespace()

		// Expect ':'
		if p.pos >= p.length || p.data[p.pos] != ':' {
			return errors.New("expected ':' after object key")
		}
		p.pos++

		p.skipWhitespace()

		// Create value and unmarshal
		elemVal := reflect.New(valueType).Elem()
		if err := p.unmarshalValue(elemVal); err != nil {
			return err
		}

		// Set map entry
		rv.SetMapIndex(reflect.ValueOf(key), elemVal)

		p.skipWhitespace()

		// Check for more entries or end of object
		if p.pos >= p.length {
			return errors.New("unexpected end of JSON input in object")
		}

		if p.data[p.pos] == '}' {
			p.pos++
			return nil
		}

		if p.data[p.pos] != ',' {
			return fmt.Errorf("expected ',' or '}' in object at position %d", p.pos)
		}
		p.pos++
	}
}

// unmarshalArray unmarshals a JSON array.
func (p *Parser) unmarshalArray(rv reflect.Value) error {
	if p.pos >= p.length || p.data[p.pos] != '[' {
		return errors.New("expected '['")
	}
	p.pos++ // skip '['

	switch rv.Kind() {
	case reflect.Slice:
		return p.unmarshalSlice(rv)
	case reflect.Array:
		return p.unmarshalFixedArray(rv)
	default:
		return fmt.Errorf("json: cannot unmarshal array into Go value of type %s", rv.Type())
	}
}

// unmarshalSlice unmarshals a JSON array into a slice.
func (p *Parser) unmarshalSlice(rv reflect.Value) error {
	sliceType := rv.Type()
	elemType := sliceType.Elem()

	// Collect elements
	var elements []reflect.Value

	p.skipWhitespace()

	// Handle empty array
	if p.pos < p.length && p.data[p.pos] == ']' {
		p.pos++
		rv.Set(reflect.MakeSlice(sliceType, 0, 0))
		return nil
	}

	for {
		p.skipWhitespace()

		// Create element and unmarshal
		elemVal := reflect.New(elemType).Elem()
		if err := p.unmarshalValue(elemVal); err != nil {
			return err
		}

		elements = append(elements, elemVal)

		p.skipWhitespace()

		// Check for more entries or end of array
		if p.pos >= p.length {
			return errors.New("unexpected end of JSON input in array")
		}

		if p.data[p.pos] == ']' {
			p.pos++
			break
		}

		if p.data[p.pos] != ',' {
			return fmt.Errorf("expected ',' or ']' in array at position %d", p.pos)
		}
		p.pos++
	}

	// Create slice and copy elements
	slice := reflect.MakeSlice(sliceType, len(elements), len(elements))
	for i, elem := range elements {
		slice.Index(i).Set(elem)
	}
	rv.Set(slice)

	return nil
}

// unmarshalFixedArray unmarshals a JSON array into a fixed-size array.
func (p *Parser) unmarshalFixedArray(rv reflect.Value) error {
	arrayLen := rv.Len()

	p.skipWhitespace()

	// Handle empty array
	if p.pos < p.length && p.data[p.pos] == ']' {
		p.pos++
		return nil
	}

	idx := 0
	for {
		p.skipWhitespace()

		if idx >= arrayLen {
			return fmt.Errorf("json: array length exceeds target array length %d", arrayLen)
		}

		// Unmarshal element
		elemVal := rv.Index(idx)
		if err := p.unmarshalValue(elemVal); err != nil {
			return err
		}

		idx++

		p.skipWhitespace()

		// Check for more entries or end of array
		if p.pos >= p.length {
			return errors.New("unexpected end of JSON input in array")
		}

		if p.data[p.pos] == ']' {
			p.pos++
			return nil
		}

		if p.data[p.pos] != ',' {
			return fmt.Errorf("expected ',' or ']' in array at position %d", p.pos)
		}
		p.pos++
	}
}

// unmarshalString unmarshals a JSON string.
func (p *Parser) unmarshalString(rv reflect.Value) error {
	s, err := p.parseString()
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.String {
		return fmt.Errorf("json: cannot unmarshal string into Go value of type %s", rv.Type())
	}

	rv.SetString(s)
	return nil
}

// unmarshalNumber unmarshals a JSON number.
func (p *Parser) unmarshalNumber(rv reflect.Value) error {
	num, err := p.parseNumber()
	if err != nil {
		return err
	}

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		switch v := num.(type) {
		case int64:
			i = v
		case float64:
			if v != float64(int64(v)) {
				return fmt.Errorf("json: cannot unmarshal number %v into Go value of type %s", v, rv.Type())
			}
			i = int64(v)
		default:
			return fmt.Errorf("json: unexpected number type %T", num)
		}

		if rv.OverflowInt(i) {
			return fmt.Errorf("json: value %d overflows %s", i, rv.Type())
		}
		rv.SetInt(i)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		switch v := num.(type) {
		case int64:
			if v < 0 {
				return fmt.Errorf("json: cannot unmarshal negative number into Go value of type %s", rv.Type())
			}
			u = uint64(v)
		case float64:
			if v < 0 || v != float64(uint64(v)) {
				return fmt.Errorf("json: cannot unmarshal number %v into Go value of type %s", v, rv.Type())
			}
			u = uint64(v)
		default:
			return fmt.Errorf("json: unexpected number type %T", num)
		}

		if rv.OverflowUint(u) {
			return fmt.Errorf("json: value %d overflows %s", u, rv.Type())
		}
		rv.SetUint(u)
		return nil

	case reflect.Float32, reflect.Float64:
		var f float64
		switch v := num.(type) {
		case int64:
			f = float64(v)
		case float64:
			f = v
		default:
			return fmt.Errorf("json: unexpected number type %T", num)
		}

		if rv.OverflowFloat(f) {
			return fmt.Errorf("json: value %v overflows %s", f, rv.Type())
		}
		rv.SetFloat(f)
		return nil

	default:
		return fmt.Errorf("json: cannot unmarshal number into Go value of type %s", rv.Type())
	}
}

// unmarshalBool unmarshals a JSON boolean.
func (p *Parser) unmarshalBool(rv reflect.Value) error {
	var b bool
	var err error

	if p.data[p.pos] == 't' {
		b, err = p.parseTrue()
	} else {
		b, err = p.parseFalse()
	}

	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Bool {
		return fmt.Errorf("json: cannot unmarshal bool into Go value of type %s", rv.Type())
	}

	rv.SetBool(b)
	return nil
}

// skipValue skips over a JSON value (used for unknown struct fields).
func (p *Parser) skipValue() error {
	_, err := p.parseValue()
	return err
}

// expectLiteral expects and consumes a specific literal string.
func (p *Parser) expectLiteral(literal string) error {
	if p.pos+len(literal) > p.length {
		return fmt.Errorf("expected '%s'", literal)
	}

	if string(p.data[p.pos:p.pos+len(literal)]) != literal {
		return fmt.Errorf("expected '%s'", literal)
	}

	p.pos += len(literal)
	return nil
}
