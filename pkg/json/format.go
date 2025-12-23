package json

import (
	"bytes"
)

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
//
// This function is compatible with encoding/json.MarshalIndent.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	person := Person{Name: "Alice", Age: 30}
//	data, _ := json.MarshalIndent(person, "", "  ")
//	// Output:
//	// {
//	//   "age": 30,
//	//   "name": "Alice"
//	// }
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	// First marshal to compact JSON
	compact, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	// Then apply indentation using pooled buffer
	buf := getBuffer()
	defer putBuffer(buf)

	if err := Indent(buf, compact, prefix, indent); err != nil {
		return nil, err
	}

	// Must copy since buffer will be returned to pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// Indent appends to dst an indented form of the JSON-encoded src.
// Each element in a JSON object or array begins on a new line
// beginning with prefix followed by one or more copies of indent
// according to the indentation nesting.
//
// The data appended to dst does not begin with the prefix nor
// any indentation, to make it easier to embed inside other formatted JSON data.
//
// Although leading space characters (space, tab, carriage return, newline)
// at the beginning of src are dropped, trailing space characters
// at the end of src are preserved and copied to dst.
//
// This function is compatible with encoding/json.Indent.
//
// Example:
//
//	compact := []byte(`{"name":"Alice","age":30}`)
//	var buf bytes.Buffer
//	json.Indent(&buf, compact, "", "  ")
//	// buf.String():
//	// {
//	//   "name": "Alice",
//	//   "age": 30
//	// }
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	// Parse the JSON to AST
	node, err := Parse(string(src))
	if err != nil {
		return err
	}

	// Render with indentation
	indented, err := RenderIndent(node, prefix, indent)
	if err != nil {
		return err
	}

	// Write to destination buffer
	dst.Write(indented)
	return nil
}

// Compact appends to dst the JSON-encoded src with
// insignificant space characters elided.
//
// This function is compatible with encoding/json.Compact.
//
// Example:
//
//	pretty := []byte(`{
//	  "name": "Alice",
//	  "age": 30
//	}`)
//	var buf bytes.Buffer
//	json.Compact(&buf, pretty)
//	// buf.String(): {"name":"Alice","age":30}
func Compact(dst *bytes.Buffer, src []byte) error {
	// Parse the JSON to AST
	node, err := Parse(string(src))
	if err != nil {
		return err
	}

	// Render compactly
	compact, err := Render(node)
	if err != nil {
		return err
	}

	// Write to destination buffer
	dst.Write(compact)
	return nil
}
