package json

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

// bufferPool is a pool of bytes.Buffer instances to reduce allocations during marshaling.
// Buffers are returned to the pool after use to minimize GC pressure.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// getBuffer retrieves a buffer from the pool and resets it for use.
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns a buffer to the pool if it's not too large.
// Buffers larger than 64KB are not pooled to avoid holding excessive memory.
func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() <= 64*1024 {
		bufferPool.Put(buf)
	}
}

// Marshal returns the JSON encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the json.Marshaler interface, Marshal calls its MarshalJSON method to produce JSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings coerced to valid UTF-8.
//
// Array and slice values encode as JSON arrays, except that []byte encodes
// as a base64-encoded string, and a nil slice encodes as the null JSON value.
//
// Struct values encode as JSON objects. Each exported struct field becomes
// a member of the object, using the field name as the object key, unless the
// field is omitted for one of the reasons given below.
//
// The encoding of each struct field can be customized by the format string
// stored under the "json" key in the struct field's tag. The format string
// gives the name of the field, possibly followed by a comma-separated list
// of options. The name may be empty in order to specify options without
// overriding the default field name.
//
// The "omitempty" option specifies that the field should be omitted from the
// encoding if the field has an empty value, defined as false, 0, a nil pointer,
// a nil interface value, and any empty array, slice, map, or string.
//
// As a special case, if the field tag is "-", the field is always omitted.
//
// Map values encode as JSON objects. The map's key type must be a string;
// the map keys are used as JSON object keys, subject to the UTF-8 coercion
// described for string values above.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as
// the null JSON value.
//
// Interface values encode as the value contained in the interface.
// A nil interface value encodes as the null JSON value.
//
// Channel, complex, and function values cannot be encoded in JSON.
// Attempting to encode such a value causes Marshal to return an error.
//
// JSON cannot represent cyclic data structures and Marshal does not handle them.
// Passing cyclic structures to Marshal will result in an error.
func Marshal(v interface{}) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	if err := marshalValue(reflect.ValueOf(v), buf, false); err != nil {
		return nil, err
	}

	// Must copy since buffer will be returned to pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// Marshaler is the interface implemented by types that can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// marshalValue marshals a reflect.Value to a buffer
func marshalValue(rv reflect.Value, buf *bytes.Buffer, asString bool) error {
	// Handle invalid values
	if !rv.IsValid() {
		buf.WriteString("null")
		return nil
	}

	// Handle nil interface
	if rv.Kind() == reflect.Interface && rv.IsNil() {
		buf.WriteString("null")
		return nil
	}

	// Check if type implements Marshaler interface
	if rv.Type().Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()) {
		marshaler := rv.Interface().(Marshaler)
		b, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(b)
		return nil
	}

	// Dereference interface
	if rv.Kind() == reflect.Interface {
		return marshalValue(rv.Elem(), buf, asString)
	}

	// Handle pointers
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return marshalValue(rv.Elem(), buf, asString)
	}

	// If asString is true, marshal numbers and bools as strings
	if asString {
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString(`"`)
			buf.WriteString(strconv.FormatInt(rv.Int(), 10))
			buf.WriteString(`"`)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			buf.WriteString(`"`)
			buf.WriteString(strconv.FormatUint(rv.Uint(), 10))
			buf.WriteString(`"`)
			return nil
		case reflect.Float32, reflect.Float64:
			buf.WriteString(`"`)
			buf.WriteString(strconv.FormatFloat(rv.Float(), 'g', -1, 64))
			buf.WriteString(`"`)
			return nil
		case reflect.Bool:
			buf.WriteString(`"`)
			buf.WriteString(strconv.FormatBool(rv.Bool()))
			buf.WriteString(`"`)
			return nil
		}
	}

	switch rv.Kind() {
	case reflect.String:
		return marshalString(rv.String(), buf)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(rv.Int(), 10))
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(rv.Uint(), 10))
		return nil

	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(rv.Float(), 'g', -1, 64))
		return nil

	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(rv.Bool()))
		return nil

	case reflect.Struct:
		return marshalStruct(rv, buf)

	case reflect.Map:
		return marshalMap(rv, buf)

	case reflect.Slice, reflect.Array:
		return marshalSlice(rv, buf)

	default:
		return fmt.Errorf("json: unsupported type %s", rv.Type())
	}
}

// marshalString marshals a string with proper JSON escaping
func marshalString(s string, buf *bytes.Buffer) error {
	buf.WriteString(`"`)
	buf.WriteString(escapeString(s))
	buf.WriteString(`"`)
	return nil
}

// marshalStruct marshals a struct to JSON
func marshalStruct(rv reflect.Value, buf *bytes.Buffer) error {
	structType := rv.Type()

	// Collect fields with their info and values
	type fieldEntry struct {
		name     string
		value    reflect.Value
		asString bool
	}

	var fields []fieldEntry

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		info := getFieldInfo(field)

		// Skip fields with "-" tag
		if info.skip {
			continue
		}

		fieldVal := rv.Field(i)

		// Handle omitempty
		if info.omitEmpty && isEmptyValue(fieldVal) {
			continue
		}

		fields = append(fields, fieldEntry{
			name:     info.name,
			value:    fieldVal,
			asString: info.asString,
		})
	}

	// Sort fields by name for deterministic output
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].name < fields[j].name
	})

	buf.WriteString("{")
	for i, field := range fields {
		if i > 0 {
			buf.WriteString(",")
		}

		// Write field name
		buf.WriteString(`"`)
		buf.WriteString(field.name)
		buf.WriteString(`":`)

		// Write field value
		if err := marshalValue(field.value, buf, field.asString); err != nil {
			return err
		}
	}

	buf.WriteString("}")
	return nil
}

// marshalMap marshals a map to JSON
func marshalMap(rv reflect.Value, buf *bytes.Buffer) error {
	if rv.IsNil() {
		buf.WriteString("null")
		return nil
	}

	mapType := rv.Type()

	// Only support string keys
	if mapType.Key().Kind() != reflect.String {
		return fmt.Errorf("json: unsupported map key type %s", mapType.Key())
	}

	buf.WriteString("{")

	// Get keys and sort them for deterministic output
	keys := rv.MapKeys()
	strKeys := make([]string, len(keys))
	for i, key := range keys {
		strKeys[i] = key.String()
	}
	sort.Strings(strKeys)

	first := true
	for _, keyStr := range strKeys {
		key := reflect.ValueOf(keyStr)
		val := rv.MapIndex(key)

		if !first {
			buf.WriteString(",")
		}
		first = false

		// Write key
		buf.WriteString(`"`)
		buf.WriteString(keyStr)
		buf.WriteString(`":`)

		// Write value
		if err := marshalValue(val, buf, false); err != nil {
			return err
		}
	}

	buf.WriteString("}")
	return nil
}

// marshalSlice marshals a slice or array to JSON
func marshalSlice(rv reflect.Value, buf *bytes.Buffer) error {
	// Nil slices encode as null
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		buf.WriteString("null")
		return nil
	}

	buf.WriteString("[")

	length := rv.Len()
	for i := 0; i < length; i++ {
		if i > 0 {
			buf.WriteString(",")
		}

		if err := marshalValue(rv.Index(i), buf, false); err != nil {
			return err
		}
	}

	buf.WriteString("]")
	return nil
}
