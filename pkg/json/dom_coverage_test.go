package json

import (
	"testing"
)

// TestDocument_Set tests the Set method with generic interface{}
func TestDocument_Set(t *testing.T) {
	doc := NewDocument()

	// Test setting various types
	doc.Set("string", "value")
	doc.Set("int", 42)
	doc.Set("bool", true)
	doc.Set("float", 3.14)
	doc.Set("null", nil)

	// Verify values
	if val, ok := doc.Get("string"); !ok || val != "value" {
		t.Errorf("Expected 'value', got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("int"); !ok || val != 42 {
		t.Errorf("Expected 42, got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("bool"); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("float"); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("null"); !ok || val != nil {
		t.Errorf("Expected nil, got %v (ok=%v)", val, ok)
	}
}

// TestDocument_Set_Chaining tests that Set returns Document for chaining
func TestDocument_Set_Chaining(t *testing.T) {
	doc := NewDocument().
		Set("a", 1).
		Set("b", 2).
		Set("c", 3)

	if doc.Size() != 3 {
		t.Errorf("Expected size 3, got %d", doc.Size())
	}
}

// TestDocument_Get tests the Get method
func TestDocument_Get(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetNull("value")

	// Test existing keys
	if val, ok := doc.Get("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("age"); !ok || val != 30 {
		t.Errorf("Expected 30, got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.Get("value"); !ok || val != nil {
		t.Errorf("Expected nil, got %v (ok=%v)", val, ok)
	}

	// Test missing key
	if val, ok := doc.Get("missing"); ok {
		t.Errorf("Expected missing key to return false, got value: %v", val)
	}
}

// TestDocument_ToMap tests the ToMap method
func TestDocument_ToMap(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetBool("active", true)

	m := doc.ToMap()

	if m == nil {
		t.Fatal("ToMap() returned nil")
	}

	if len(m) != 3 {
		t.Errorf("Expected map size 3, got %d", len(m))
	}

	if m["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got %v", m["name"])
	}
	if m["age"] != 30 {
		t.Errorf("Expected age=30, got %v", m["age"])
	}
	if m["active"] != true {
		t.Errorf("Expected active=true, got %v", m["active"])
	}
}

// TestDocument_ToMap_Modification tests that modifying the map affects the document
func TestDocument_ToMap_Modification(t *testing.T) {
	doc := NewDocument().SetString("name", "Alice")

	m := doc.ToMap()
	m["age"] = 30

	// The map is the underlying data, so changes should be reflected
	if val, ok := doc.GetInt("age"); !ok || val != 30 {
		t.Errorf("Expected age=30 after map modification, got %v (ok=%v)", val, ok)
	}
}

// TestDocument_GetInt64_AllPaths tests all code paths in GetInt64
func TestDocument_GetInt64_AllPaths(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Document
		key      string
		expected int64
		ok       bool
	}{
		{
			name: "int64 value",
			setup: func() *Document {
				return NewDocument().SetInt64("value", 9223372036854775807)
			},
			key:      "value",
			expected: 9223372036854775807,
			ok:       true,
		},
		{
			name: "int value",
			setup: func() *Document {
				return NewDocument().SetInt("value", 42)
			},
			key:      "value",
			expected: 42,
			ok:       true,
		},
		{
			name: "float value",
			setup: func() *Document {
				return NewDocument().SetFloat("value", 42.0)
			},
			key:      "value",
			expected: 42,
			ok:       true,
		},
		{
			name: "missing key",
			setup: func() *Document {
				return NewDocument()
			},
			key:      "missing",
			expected: 0,
			ok:       false,
		},
		{
			name: "wrong type",
			setup: func() *Document {
				return NewDocument().SetString("value", "not a number")
			},
			key:      "value",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := tt.setup()
			val, ok := doc.GetInt64(tt.key)
			if ok != tt.ok {
				t.Errorf("Expected ok=%v, got %v", tt.ok, ok)
			}
			if val != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, val)
			}
		})
	}
}

// TestArray_Add tests the Add method with generic interface{}
func TestArray_Add(t *testing.T) {
	arr := NewArray()

	// Test adding various types
	arr.Add("string")
	arr.Add(42)
	arr.Add(true)
	arr.Add(3.14)
	arr.Add(nil)

	if arr.Len() != 5 {
		t.Errorf("Expected length 5, got %d", arr.Len())
	}

	// Verify values
	if val, ok := arr.Get(0); !ok || val != "string" {
		t.Errorf("Expected 'string', got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.Get(1); !ok || val != 42 {
		t.Errorf("Expected 42, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.Get(2); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.Get(3); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.Get(4); !ok || val != nil {
		t.Errorf("Expected nil, got %v (ok=%v)", val, ok)
	}
}

// TestArray_Add_Chaining tests that Add returns Array for chaining
func TestArray_Add_Chaining(t *testing.T) {
	arr := NewArray().
		Add(1).
		Add(2).
		Add(3)

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}
}

// TestArray_AddInt64 tests the AddInt64 method
func TestArray_AddInt64(t *testing.T) {
	arr := NewArray().
		AddInt64(9223372036854775807).
		AddInt64(-9223372036854775808).
		AddInt64(0)

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	// Verify values
	if val, ok := arr.GetInt64(0); !ok || val != 9223372036854775807 {
		t.Errorf("Expected max int64, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.GetInt64(1); !ok || val != -9223372036854775808 {
		t.Errorf("Expected min int64, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.GetInt64(2); !ok || val != 0 {
		t.Errorf("Expected 0, got %v (ok=%v)", val, ok)
	}
}

// TestArray_GetInt64_AllPaths tests all code paths in GetInt64
func TestArray_GetInt64_AllPaths(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Array
		index    int
		expected int64
		ok       bool
	}{
		{
			name: "int64 value",
			setup: func() *Array {
				return NewArray().AddInt64(9223372036854775807)
			},
			index:    0,
			expected: 9223372036854775807,
			ok:       true,
		},
		{
			name: "int value",
			setup: func() *Array {
				return NewArray().AddInt(42)
			},
			index:    0,
			expected: 42,
			ok:       true,
		},
		{
			name: "float value",
			setup: func() *Array {
				return NewArray().AddFloat(42.0)
			},
			index:    0,
			expected: 42,
			ok:       true,
		},
		{
			name: "out of bounds negative",
			setup: func() *Array {
				return NewArray().AddInt64(42)
			},
			index:    -1,
			expected: 0,
			ok:       false,
		},
		{
			name: "out of bounds positive",
			setup: func() *Array {
				return NewArray().AddInt64(42)
			},
			index:    1,
			expected: 0,
			ok:       false,
		},
		{
			name: "wrong type",
			setup: func() *Array {
				return NewArray().AddString("not a number")
			},
			index:    0,
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := tt.setup()
			val, ok := arr.GetInt64(tt.index)
			if ok != tt.ok {
				t.Errorf("Expected ok=%v, got %v", tt.ok, ok)
			}
			if val != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, val)
			}
		})
	}
}

// TestArray_ToSlice tests the ToSlice method
func TestArray_ToSlice(t *testing.T) {
	arr := NewArray().
		AddString("go").
		AddInt(42).
		AddBool(true)

	slice := arr.ToSlice()

	if slice == nil {
		t.Fatal("ToSlice() returned nil")
	}

	if len(slice) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(slice))
	}

	if slice[0] != "go" {
		t.Errorf("Expected slice[0]='go', got %v", slice[0])
	}
	if slice[1] != 42 {
		t.Errorf("Expected slice[1]=42, got %v", slice[1])
	}
	if slice[2] != true {
		t.Errorf("Expected slice[2]=true, got %v", slice[2])
	}
}

// TestArray_ToSlice_Modification tests that modifying the slice affects the array
func TestArray_ToSlice_Modification(t *testing.T) {
	arr := NewArray().AddString("original")

	slice := arr.ToSlice()
	slice = append(slice, "new")

	// The original array should not be affected by append
	if arr.Len() != 1 {
		t.Errorf("Expected array length 1 after slice append, got %d", arr.Len())
	}

	// But direct modification should affect it
	arr2 := NewArray().AddString("original")
	slice2 := arr2.ToSlice()
	slice2[0] = "modified"

	// The slice is the underlying data, so changes should be reflected
	if val, ok := arr2.GetString(0); !ok || val != "modified" {
		t.Errorf("Expected 'modified' after slice modification, got %v (ok=%v)", val, ok)
	}
}

// TestDocument_JSON_ErrorPath tests JSON method error handling
func TestDocument_JSON_ErrorPath(t *testing.T) {
	// Create a document with a value that can't be converted
	doc := NewDocument()
	doc.data["invalid"] = make(chan int) // channels can't be converted

	_, err := doc.JSON()
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestArray_JSON_ErrorPath tests JSON method error handling
func TestArray_JSON_ErrorPath(t *testing.T) {
	// Create an array with a value that can't be converted
	arr := NewArray()
	arr.data = append(arr.data, make(chan int)) // channels can't be converted

	_, err := arr.JSON()
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestDocument_JSONIndent_ErrorPath tests JSONIndent method error handling
func TestDocument_JSONIndent_ErrorPath(t *testing.T) {
	// Create a document with a value that can't be marshaled
	doc := NewDocument()
	doc.data["invalid"] = make(chan int)

	_, err := doc.JSONIndent("", "  ")
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestArray_JSONIndent_ErrorPath tests JSONIndent method error handling
func TestArray_JSONIndent_ErrorPath(t *testing.T) {
	// Create an array with a value that can't be marshaled
	arr := NewArray()
	arr.data = append(arr.data, make(chan int))

	_, err := arr.JSONIndent("", "  ")
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestDocument_UnmarshalJSON_ErrorPath tests UnmarshalJSON error handling
func TestDocument_UnmarshalJSON_ErrorPath(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"invalid json", `{invalid}`},
		{"array instead of object", `[1,2,3]`},
		{"primitive instead of object", `"string"`},
		{"null", `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Document
			err := doc.UnmarshalJSON([]byte(tt.data))
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

// TestArray_UnmarshalJSON_ErrorPath tests UnmarshalJSON error handling
func TestArray_UnmarshalJSON_ErrorPath(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"invalid json", `[invalid]`},
		{"object instead of array", `{"key":"value"}`},
		{"primitive instead of array", `"string"`},
		{"null", `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var arr Array
			err := arr.UnmarshalJSON([]byte(tt.data))
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

// TestDocument_MarshalJSON_ErrorPath tests MarshalJSON error handling
func TestDocument_MarshalJSON_ErrorPath(t *testing.T) {
	doc := NewDocument()
	doc.data["invalid"] = make(chan int)

	_, err := doc.MarshalJSON()
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestArray_MarshalJSON_ErrorPath tests MarshalJSON error handling
func TestArray_MarshalJSON_ErrorPath(t *testing.T) {
	arr := NewArray()
	arr.data = append(arr.data, make(chan int))

	_, err := arr.MarshalJSON()
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestArray_IsNull_EdgeCases tests IsNull with edge cases
func TestArray_IsNull_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddString("value").
		AddNull().
		AddInt(42)

	// Out of bounds negative
	if arr.IsNull(-1) {
		t.Error("IsNull(-1) should return false for out of bounds")
	}

	// Out of bounds positive
	if arr.IsNull(10) {
		t.Error("IsNull(10) should return false for out of bounds")
	}

	// Non-null value
	if arr.IsNull(0) {
		t.Error("IsNull(0) should return false for non-null value")
	}

	// Null value
	if !arr.IsNull(1) {
		t.Error("IsNull(1) should return true for null value")
	}
}

// TestArray_GetBool_EdgeCases tests GetBool with various edge cases
func TestArray_GetBool_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddBool(true).
		AddBool(false).
		AddInt(42)

	// Valid bool true
	if val, ok := arr.GetBool(0); !ok || !val {
		t.Errorf("Expected true at index 0, got %v (ok=%v)", val, ok)
	}

	// Valid bool false
	if val, ok := arr.GetBool(1); !ok || val {
		t.Errorf("Expected false at index 1, got %v (ok=%v)", val, ok)
	}

	// Wrong type
	if val, ok := arr.GetBool(2); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", val)
	}

	// Out of bounds negative
	if val, ok := arr.GetBool(-1); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}

	// Out of bounds positive
	if val, ok := arr.GetBool(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}
}

// TestArray_GetFloat_EdgeCases tests GetFloat with int64 conversion
func TestArray_GetFloat_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddFloat(3.14).
		AddInt(42).
		AddInt64(int64(100))

	// Float value
	if val, ok := arr.GetFloat(0); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %v (ok=%v)", val, ok)
	}

	// Int value (should convert)
	if val, ok := arr.GetFloat(1); !ok || val != 42.0 {
		t.Errorf("Expected 42.0, got %v (ok=%v)", val, ok)
	}

	// Int64 value (should convert)
	if val, ok := arr.GetFloat(2); !ok || val != 100.0 {
		t.Errorf("Expected 100.0, got %v (ok=%v)", val, ok)
	}
}

// TestArray_GetObject_EdgeCases tests GetObject edge cases
func TestArray_GetObject_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddObject(NewDocument().SetString("name", "Alice")).
		AddInt(42)

	// Valid object
	if obj, ok := arr.GetObject(0); !ok {
		t.Error("Expected object at index 0")
	} else if name, ok := obj.GetString("name"); !ok || name != "Alice" {
		t.Errorf("Expected name='Alice', got %v (ok=%v)", name, ok)
	}

	// Wrong type
	if obj, ok := arr.GetObject(1); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", obj)
	}

	// Out of bounds negative
	if obj, ok := arr.GetObject(-1); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", obj)
	}

	// Out of bounds positive
	if obj, ok := arr.GetObject(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", obj)
	}
}

// TestArray_GetArray_EdgeCases tests GetArray edge cases
func TestArray_GetArray_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddArray(NewArray().AddInt(1).AddInt(2)).
		AddInt(42)

	// Valid array
	if nestedArr, ok := arr.GetArray(0); !ok {
		t.Error("Expected array at index 0")
	} else if nestedArr.Len() != 2 {
		t.Errorf("Expected nested array length 2, got %d", nestedArr.Len())
	}

	// Wrong type
	if nestedArr, ok := arr.GetArray(1); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", nestedArr)
	}

	// Out of bounds negative
	if nestedArr, ok := arr.GetArray(-1); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", nestedArr)
	}

	// Out of bounds positive
	if nestedArr, ok := arr.GetArray(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", nestedArr)
	}
}

// TestDocument_GetBool_EdgeCases tests GetBool when value exists but is false
func TestDocument_GetBool_EdgeCases(t *testing.T) {
	doc := NewDocument().
		SetBool("t", true).
		SetBool("f", false).
		SetInt("notBool", 42)

	// True value
	if val, ok := doc.GetBool("t"); !ok || !val {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}

	// False value (make sure it returns ok=true)
	if val, ok := doc.GetBool("f"); !ok || val {
		t.Errorf("Expected false with ok=true, got %v (ok=%v)", val, ok)
	}

	// Wrong type
	if val, ok := doc.GetBool("notBool"); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", val)
	}

	// Missing key
	if val, ok := doc.GetBool("missing"); ok {
		t.Errorf("Expected ok=false for missing key, got %v", val)
	}
}

// TestDocument_GetFloat_Int64Conversion tests GetFloat with int64 values
func TestDocument_GetFloat_Int64Conversion(t *testing.T) {
	doc := NewDocument().
		SetFloat("f", 3.14).
		SetInt("i", 42).
		SetInt64("i64", int64(100))

	// Float value
	if val, ok := doc.GetFloat("f"); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %v (ok=%v)", val, ok)
	}

	// Int value (should convert)
	if val, ok := doc.GetFloat("i"); !ok || val != 42.0 {
		t.Errorf("Expected 42.0, got %v (ok=%v)", val, ok)
	}

	// Int64 value (should convert)
	if val, ok := doc.GetFloat("i64"); !ok || val != 100.0 {
		t.Errorf("Expected 100.0, got %v (ok=%v)", val, ok)
	}
}

// TestDocument_GetObject_EdgeCases tests GetObject when not found or wrong type
func TestDocument_GetObject_EdgeCases(t *testing.T) {
	doc := NewDocument().
		SetObject("obj", NewDocument().SetString("name", "Alice")).
		SetInt("notObj", 42)

	// Valid object
	if obj, ok := doc.GetObject("obj"); !ok {
		t.Error("Expected object")
	} else if name, ok := obj.GetString("name"); !ok || name != "Alice" {
		t.Errorf("Expected name='Alice', got %v (ok=%v)", name, ok)
	}

	// Wrong type
	if obj, ok := doc.GetObject("notObj"); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", obj)
	}

	// Missing key
	if obj, ok := doc.GetObject("missing"); ok {
		t.Errorf("Expected ok=false for missing key, got %v", obj)
	}
}

// TestDocument_GetArray_EdgeCases tests GetArray when not found or wrong type
func TestDocument_GetArray_EdgeCases(t *testing.T) {
	doc := NewDocument().
		SetArray("arr", NewArray().AddInt(1).AddInt(2)).
		SetInt("notArr", 42)

	// Valid array
	if arr, ok := doc.GetArray("arr"); !ok {
		t.Error("Expected array")
	} else if arr.Len() != 2 {
		t.Errorf("Expected array length 2, got %d", arr.Len())
	}

	// Wrong type
	if arr, ok := doc.GetArray("notArr"); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", arr)
	}

	// Missing key
	if arr, ok := doc.GetArray("missing"); ok {
		t.Errorf("Expected ok=false for missing key, got %v", arr)
	}
}

// TestArray_Get_EdgeCases tests Array Get with edge cases
func TestArray_Get_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddString("value").
		AddNull()

	// Valid index
	if val, ok := arr.Get(0); !ok || val != "value" {
		t.Errorf("Expected 'value', got %v (ok=%v)", val, ok)
	}

	// Null value (should still return ok=true)
	if val, ok := arr.Get(1); !ok || val != nil {
		t.Errorf("Expected nil with ok=true, got %v (ok=%v)", val, ok)
	}

	// Out of bounds negative
	if val, ok := arr.Get(-1); ok {
		t.Errorf("Expected ok=false for negative index, got %v", val)
	}

	// Out of bounds positive
	if val, ok := arr.Get(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}
}

// TestArray_GetString_EdgeCases tests GetString with edge cases
func TestArray_GetString_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddString("hello").
		AddInt(42)

	// Valid string
	if val, ok := arr.GetString(0); !ok || val != "hello" {
		t.Errorf("Expected 'hello', got %v (ok=%v)", val, ok)
	}

	// Wrong type
	if val, ok := arr.GetString(1); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", val)
	}

	// Out of bounds negative
	if val, ok := arr.GetString(-1); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}

	// Out of bounds positive
	if val, ok := arr.GetString(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}
}

// TestArray_GetInt_EdgeCases tests GetInt with edge cases and type conversions
func TestArray_GetInt_EdgeCases(t *testing.T) {
	arr := NewArray().
		AddInt(42).
		AddFloat(100.0).
		AddInt64(int64(200)).
		AddString("notInt")

	// Int value
	if val, ok := arr.GetInt(0); !ok || val != 42 {
		t.Errorf("Expected 42, got %v (ok=%v)", val, ok)
	}

	// Float value (should convert if whole number)
	if val, ok := arr.GetInt(1); !ok || val != 100 {
		t.Errorf("Expected 100, got %v (ok=%v)", val, ok)
	}

	// Int64 value (should convert)
	if val, ok := arr.GetInt(2); !ok || val != 200 {
		t.Errorf("Expected 200, got %v (ok=%v)", val, ok)
	}

	// Wrong type
	if val, ok := arr.GetInt(3); ok {
		t.Errorf("Expected ok=false for wrong type, got %v", val)
	}

	// Out of bounds negative
	if val, ok := arr.GetInt(-1); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}

	// Out of bounds positive
	if val, ok := arr.GetInt(10); ok {
		t.Errorf("Expected ok=false for out of bounds, got %v", val)
	}
}
