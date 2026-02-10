package json

import (
	"bytes"
	"reflect"
	"sync"
)

// bufPool pools []byte slices for the compiled-encoder fast path.
var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 1024)
		return &b
	},
}

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
	if v == nil {
		return []byte("null"), nil
	}

	// Fast path: try the type-switch encoder first (no reflect at all)
	bp := bufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	buf, err := appendInterface(buf, v)
	if err == nil {
		// Success — copy result so pooled buffer can be reused
		result := make([]byte, len(buf))
		copy(result, buf)
		*bp = buf
		bufPool.Put(bp)
		return result, nil
	}

	if err != errNeedReflect {
		// Real error from appendInterface (e.g. Marshaler failed)
		*bp = buf
		bufPool.Put(bp)
		return nil, err
	}

	// Fall back to the compiled encoder cache
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			*bp = buf
			bufPool.Put(bp)
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	enc := encoderForType(rv.Type())
	buf, err = enc(buf, rv)
	if err != nil {
		*bp = buf
		bufPool.Put(bp)
		return nil, err
	}

	result := make([]byte, len(buf))
	copy(result, buf)
	*bp = buf
	bufPool.Put(bp)
	return result, nil
}

// Marshaler is the interface implemented by types that can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
