// Package json provides a user-friendly DOM API for JSON manipulation.
//
// The DOM API provides type-safe, fluent interfaces for building and manipulating
// JSON documents without requiring type assertions or working with raw AST nodes.
//
// # Document Type
//
// Document represents a JSON object with chainable methods:
//
//	doc := json.NewDocument().
//		String("name", "Alice").
//		Int("age", 30).
//		Bool("active", true)
//
// # Array Type
//
// Array represents a JSON array with chainable methods:
//
//	arr := json.NewArray().
//		AddString("go").
//		AddString("json").
//		AddInt(42)
//
// # Type-Safe Getters
//
// Access values without type assertions:
//
//	name, ok := doc.GetString("name")  // "Alice", true
//	age, ok := doc.GetInt("age")       // 30, true
//
// # Nested Structures
//
// Build complex nested structures fluently:
//
//	doc := json.NewDocument().
//		String("name", "Alice").
//		Object("address", json.NewDocument().
//			String("city", "NYC").
//			String("zip", "10001")).
//		Array("tags", json.NewArray().
//			AddString("go").
//			AddString("json"))
package json

import (
	"fmt"
)

// Document represents a JSON object with a fluent API for manipulation.
// All setter methods return *Document to enable method chaining.
type Document struct {
	data map[string]interface{}
}

// Array represents a JSON array with a fluent API for manipulation.
// All append methods return *Array to enable method chaining.
type Array struct {
	data []interface{}
}

// NewDocument creates a new empty Document.
func NewDocument() *Document {
	return &Document{data: make(map[string]interface{})}
}

// NewArray creates a new empty Array.
func NewArray() *Array {
	return &Array{data: make([]interface{}, 0)}
}

// ParseDocument parses JSON string into a Document with a fluent API.
// Returns an error if the input is not valid JSON or not an object.
func ParseDocument(input string) (*Document, error) {
	// Parse JSON to AST
	node, err := Parse(input)
	if err != nil {
		return nil, err
	}

	// Convert AST to map[string]interface{}
	value := NodeToInterface(node)
	data, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected JSON object, got %T", value)
	}
	return &Document{data: data}, nil
}

// ParseArray parses JSON string into an Array with a fluent API.
// Returns an error if the input is not valid JSON or not an array.
func ParseArray(input string) (*Array, error) {
	// Parse JSON to AST
	node, err := Parse(input)
	if err != nil {
		return nil, err
	}

	// Convert AST to []interface{}
	value := NodeToInterface(node)
	data, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected JSON array, got %T", value)
	}
	return &Array{data: data}, nil
}

// ============================================================================
// Document Builder Methods (fluent setters that return *Document)
// ============================================================================

// Set sets a string value and returns the Document for chaining.
func (d *Document) Set(key string, value interface{}) *Document {
	d.data[key] = value
	return d
}

// SetString sets a string value and returns the Document for chaining.
func (d *Document) SetString(key, value string) *Document {
	d.data[key] = value
	return d
}

// SetInt sets an int value and returns the Document for chaining.
func (d *Document) SetInt(key string, value int) *Document {
	d.data[key] = value
	return d
}

// SetInt64 sets an int64 value and returns the Document for chaining.
func (d *Document) SetInt64(key string, value int64) *Document {
	d.data[key] = value
	return d
}

// SetBool sets a bool value and returns the Document for chaining.
func (d *Document) SetBool(key string, value bool) *Document {
	d.data[key] = value
	return d
}

// SetFloat sets a float64 value and returns the Document for chaining.
func (d *Document) SetFloat(key string, value float64) *Document {
	d.data[key] = value
	return d
}

// SetNull sets a null value and returns the Document for chaining.
func (d *Document) SetNull(key string) *Document {
	d.data[key] = nil
	return d
}

// SetObject sets a nested Document and returns the parent Document for chaining.
func (d *Document) SetObject(key string, value *Document) *Document {
	d.data[key] = value.data
	return d
}

// SetArray sets an Array and returns the Document for chaining.
func (d *Document) SetArray(key string, value *Array) *Document {
	d.data[key] = value.data
	return d
}

// ============================================================================
// Document Getter Methods (type-safe access)
// ============================================================================

// Get gets a value as interface{}. Returns nil if not found.
func (d *Document) Get(key string) (interface{}, bool) {
	val, ok := d.data[key]
	return val, ok
}

// GetString gets a string value. Returns empty string and false if not found or wrong type.
func (d *Document) GetString(key string) (string, bool) {
	if val, ok := d.data[key]; ok {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt gets an int value. Returns 0 and false if not found or wrong type.
// Handles conversion from float64 (JSON numbers are parsed as float64).
func (d *Document) GetInt(key string) (int, bool) {
	if val, ok := d.data[key]; ok {
		switch v := val.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case int64:
			return int(v), true
		}
	}
	return 0, false
}

// GetInt64 gets an int64 value. Returns 0 and false if not found or wrong type.
// Handles conversion from float64 (JSON numbers are parsed as float64).
func (d *Document) GetInt64(key string) (int64, bool) {
	if val, ok := d.data[key]; ok {
		switch v := val.(type) {
		case int64:
			return v, true
		case int:
			return int64(v), true
		case float64:
			return int64(v), true
		}
	}
	return 0, false
}

// GetBool gets a bool value. Returns false and false if not found or wrong type.
func (d *Document) GetBool(key string) (bool, bool) {
	if val, ok := d.data[key]; ok {
		if b, ok := val.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// GetFloat gets a float64 value. Returns 0.0 and false if not found or wrong type.
// Handles conversion from int types.
func (d *Document) GetFloat(key string) (float64, bool) {
	if val, ok := d.data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case int64:
			return float64(v), true
		}
	}
	return 0.0, false
}

// GetObject gets a nested Document. Returns nil and false if not found or wrong type.
func (d *Document) GetObject(key string) (*Document, bool) {
	if val, ok := d.data[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return &Document{data: m}, true
		}
	}
	return nil, false
}

// GetArray gets an Array. Returns nil and false if not found or wrong type.
func (d *Document) GetArray(key string) (*Array, bool) {
	if val, ok := d.data[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			return &Array{data: arr}, true
		}
	}
	return nil, false
}

// IsNull checks if a key exists and has a null value.
func (d *Document) IsNull(key string) bool {
	val, ok := d.data[key]
	return ok && val == nil
}

// Has checks if a key exists (including null values).
func (d *Document) Has(key string) bool {
	_, ok := d.data[key]
	return ok
}

// Remove removes a key and returns the Document for chaining.
func (d *Document) Remove(key string) *Document {
	delete(d.data, key)
	return d
}

// Keys returns all keys in the Document.
func (d *Document) Keys() []string {
	keys := make([]string, 0, len(d.data))
	for k := range d.data {
		keys = append(keys, k)
	}
	return keys
}

// Size returns the number of properties in the Document.
func (d *Document) Size() int {
	return len(d.data)
}

// ToMap returns the underlying map[string]interface{}.
func (d *Document) ToMap() map[string]interface{} {
	return d.data
}

// JSON marshals the Document to a JSON string.
func (d *Document) JSON() (string, error) {
	// Convert map to AST
	node, err := InterfaceToNode(d.data)
	if err != nil {
		return "", err
	}

	// Render AST to JSON
	bytes, err := Render(node)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JSONIndent returns a pretty-printed JSON string representation with indentation.
// The prefix is written at the beginning of each line, and indent specifies the indentation string.
//
// Common usage:
//   - JSONIndent("", "  ") - 2-space indentation
//   - JSONIndent("", "\t") - tab indentation
//   - JSONIndent(">>", "  ") - prefix each line with ">>" and use 2-space indent
//
// Example:
//
//	doc := NewDocument().
//	    SetString("name", "Alice").
//	    SetInt("age", 30)
//	pretty, _ := doc.JSONIndent("", "  ")
//	// Output:
//	// {
//	//   "age": 30,
//	//   "name": "Alice"
//	// }
func (d *Document) JSONIndent(prefix, indent string) (string, error) {
	bytes, err := MarshalIndent(d.data, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MarshalJSON implements json.Marshaler interface.
func (d *Document) MarshalJSON() ([]byte, error) {
	// Convert map to AST
	node, err := InterfaceToNode(d.data)
	if err != nil {
		return nil, err
	}

	// Render AST to JSON
	return Render(node)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (d *Document) UnmarshalJSON(data []byte) error {
	// Parse JSON to AST
	node, err := Parse(string(data))
	if err != nil {
		return err
	}

	// Convert AST to map
	value := NodeToInterface(node)
	m, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected JSON object, got %T", value)
	}
	d.data = m
	return nil
}

// ============================================================================
// Array Builder Methods (fluent append methods that return *Array)
// ============================================================================

// Add appends an interface{} value and returns the Array for chaining.
func (a *Array) Add(value interface{}) *Array {
	a.data = append(a.data, value)
	return a
}

// AddString appends a string and returns the Array for chaining.
func (a *Array) AddString(value string) *Array {
	a.data = append(a.data, value)
	return a
}

// AddInt appends an int and returns the Array for chaining.
func (a *Array) AddInt(value int) *Array {
	a.data = append(a.data, value)
	return a
}

// AddInt64 appends an int64 and returns the Array for chaining.
func (a *Array) AddInt64(value int64) *Array {
	a.data = append(a.data, value)
	return a
}

// AddBool appends a bool and returns the Array for chaining.
func (a *Array) AddBool(value bool) *Array {
	a.data = append(a.data, value)
	return a
}

// AddFloat appends a float64 and returns the Array for chaining.
func (a *Array) AddFloat(value float64) *Array {
	a.data = append(a.data, value)
	return a
}

// AddNull appends a null and returns the Array for chaining.
func (a *Array) AddNull() *Array {
	a.data = append(a.data, nil)
	return a
}

// AddObject appends a Document and returns the Array for chaining.
func (a *Array) AddObject(value *Document) *Array {
	a.data = append(a.data, value.data)
	return a
}

// AddArray appends an Array and returns the parent Array for chaining.
func (a *Array) AddArray(value *Array) *Array {
	a.data = append(a.data, value.data)
	return a
}

// ============================================================================
// Array Getter Methods (type-safe indexed access)
// ============================================================================

// Get gets a value at index as interface{}. Returns nil if out of bounds.
func (a *Array) Get(index int) (interface{}, bool) {
	if index < 0 || index >= len(a.data) {
		return nil, false
	}
	return a.data[index], true
}

// GetString gets a string at index. Returns empty string and false if not found or wrong type.
func (a *Array) GetString(index int) (string, bool) {
	if index < 0 || index >= len(a.data) {
		return "", false
	}
	if str, ok := a.data[index].(string); ok {
		return str, true
	}
	return "", false
}

// GetInt gets an int at index. Returns 0 and false if not found or wrong type.
// Handles conversion from float64 (JSON numbers are parsed as float64).
func (a *Array) GetInt(index int) (int, bool) {
	if index < 0 || index >= len(a.data) {
		return 0, false
	}
	switch v := a.data[index].(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case int64:
		return int(v), true
	}
	return 0, false
}

// GetInt64 gets an int64 at index. Returns 0 and false if not found or wrong type.
// Handles conversion from float64 (JSON numbers are parsed as float64).
func (a *Array) GetInt64(index int) (int64, bool) {
	if index < 0 || index >= len(a.data) {
		return 0, false
	}
	switch v := a.data[index].(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	}
	return 0, false
}

// GetBool gets a bool at index. Returns false and false if not found or wrong type.
func (a *Array) GetBool(index int) (bool, bool) {
	if index < 0 || index >= len(a.data) {
		return false, false
	}
	if b, ok := a.data[index].(bool); ok {
		return b, true
	}
	return false, false
}

// GetFloat gets a float64 at index. Returns 0.0 and false if not found or wrong type.
// Handles conversion from int types.
func (a *Array) GetFloat(index int) (float64, bool) {
	if index < 0 || index >= len(a.data) {
		return 0.0, false
	}
	switch v := a.data[index].(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	}
	return 0.0, false
}

// GetObject gets a Document at index. Returns nil and false if not found or wrong type.
func (a *Array) GetObject(index int) (*Document, bool) {
	if index < 0 || index >= len(a.data) {
		return nil, false
	}
	if m, ok := a.data[index].(map[string]interface{}); ok {
		return &Document{data: m}, true
	}
	return nil, false
}

// GetArray gets an Array at index. Returns nil and false if not found or wrong type.
func (a *Array) GetArray(index int) (*Array, bool) {
	if index < 0 || index >= len(a.data) {
		return nil, false
	}
	if arr, ok := a.data[index].([]interface{}); ok {
		return &Array{data: arr}, true
	}
	return nil, false
}

// IsNull checks if the value at index is null.
func (a *Array) IsNull(index int) bool {
	if index < 0 || index >= len(a.data) {
		return false
	}
	return a.data[index] == nil
}

// Len returns the length of the Array.
func (a *Array) Len() int {
	return len(a.data)
}

// ToSlice returns the underlying []interface{}.
func (a *Array) ToSlice() []interface{} {
	return a.data
}

// JSON marshals the Array to a JSON string.
func (a *Array) JSON() (string, error) {
	// Convert slice to AST
	node, err := InterfaceToNode(a.data)
	if err != nil {
		return "", err
	}

	// Render AST to JSON
	bytes, err := Render(node)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JSONIndent returns a pretty-printed JSON string representation with indentation.
// The prefix is written at the beginning of each line, and indent specifies the indentation string.
//
// Common usage:
//   - JSONIndent("", "  ") - 2-space indentation
//   - JSONIndent("", "\t") - tab indentation
//   - JSONIndent(">>", "  ") - prefix each line with ">>" and use 2-space indent
//
// Example:
//
//	arr := NewArray().
//	    AddString("apple").
//	    AddString("banana").
//	    AddInt(42)
//	pretty, _ := arr.JSONIndent("", "  ")
//	// Output:
//	// [
//	//   "apple",
//	//   "banana",
//	//   42
//	// ]
func (a *Array) JSONIndent(prefix, indent string) (string, error) {
	bytes, err := MarshalIndent(a.data, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MarshalJSON implements json.Marshaler interface.
func (a *Array) MarshalJSON() ([]byte, error) {
	// Convert slice to AST
	node, err := InterfaceToNode(a.data)
	if err != nil {
		return nil, err
	}

	// Render AST to JSON
	return Render(node)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (a *Array) UnmarshalJSON(data []byte) error {
	// Parse JSON to AST
	node, err := Parse(string(data))
	if err != nil {
		return err
	}

	// Convert AST to slice
	value := NodeToInterface(node)
	slice, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("expected JSON array, got %T", value)
	}
	a.data = slice
	return nil
}
