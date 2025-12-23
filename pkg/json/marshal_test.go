package json

import (
	"strings"
	"testing"
)

// TestMarshal_BasicTypes tests marshaling basic Go types
func TestMarshal_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string",
			value:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "int",
			value:    42,
			expected: `42`,
		},
		{
			name:     "int64",
			value:    int64(42),
			expected: `42`,
		},
		{
			name:     "float64",
			value:    3.14,
			expected: `3.14`,
		},
		{
			name:     "bool true",
			value:    true,
			expected: `true`,
		},
		{
			name:     "bool false",
			value:    false,
			expected: `false`,
		},
		{
			name:     "nil",
			value:    nil,
			expected: `null`,
		},
		{
			name:     "string with special chars",
			value:    "hello\nworld\t\"quoted\"",
			expected: `"hello\nworld\t\"quoted\""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_Struct tests marshaling structs
func TestMarshal_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "simple struct",
			value:    Person{Name: "Alice", Age: 30},
			expected: `{"Age":30,"Name":"Alice"}`,
		},
		{
			name:     "struct pointer",
			value:    &Person{Name: "Bob", Age: 25},
			expected: `{"Age":25,"Name":"Bob"}`,
		},
		{
			name:     "empty struct",
			value:    Person{},
			expected: `{"Age":0,"Name":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_StructTags tests marshaling with json struct tags
func TestMarshal_StructTags(t *testing.T) {
	type Tagged struct {
		PublicName  string `json:"name"`
		InternalAge int    `json:"age"`
		Ignored     string `json:"-"`
		NoTag       string
		Empty       string `json:"empty,omitempty"`
		ZeroInt     int    `json:"zero,omitempty"`
	}

	tests := []struct {
		name     string
		value    Tagged
		expected string
	}{
		{
			name: "all fields",
			value: Tagged{
				PublicName:  "Alice",
				InternalAge: 30,
				Ignored:     "should not appear",
				NoTag:       "visible",
				Empty:       "not empty",
				ZeroInt:     5,
			},
			expected: `{"NoTag":"visible","age":30,"empty":"not empty","name":"Alice","zero":5}`,
		},
		{
			name: "omitempty with empty values",
			value: Tagged{
				PublicName: "Bob",
				Empty:      "",
				ZeroInt:    0,
			},
			expected: `{"NoTag":"","age":0,"name":"Bob"}`,
		},
		{
			name: "ignored field not in output",
			value: Tagged{
				PublicName: "Charlie",
				Ignored:    "this should not appear",
			},
			expected: `{"NoTag":"","age":0,"name":"Charlie"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_StringOption tests the string option in struct tags
func TestMarshal_StringOption(t *testing.T) {
	type StringOpts struct {
		NormalInt  int     `json:"normal"`
		StringInt  int     `json:"stringInt,string"`
		StringBool bool    `json:"stringBool,string"`
		StringNum  float64 `json:"stringNum,string"`
	}

	value := StringOpts{
		NormalInt:  42,
		StringInt:  42,
		StringBool: true,
		StringNum:  3.14,
	}

	result, err := Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := `{"normal":42,"stringBool":"true","stringInt":"42","stringNum":"3.14"}`
	if string(result) != expected {
		t.Errorf("Marshal() = %s, want %s", string(result), expected)
	}
}

// TestMarshal_NestedStruct tests marshaling nested structures
func TestMarshal_NestedStruct(t *testing.T) {
	type Address struct {
		City  string `json:"city"`
		State string `json:"state"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	person := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			City:  "Seattle",
			State: "WA",
		},
	}

	result, err := Marshal(person)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := `{"address":{"city":"Seattle","state":"WA"},"age":30,"name":"Alice"}`
	if string(result) != expected {
		t.Errorf("Marshal() = %s, want %s", string(result), expected)
	}
}

// TestMarshal_Slices tests marshaling slices
func TestMarshal_Slices(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "int slice",
			value:    []int{1, 2, 3, 4, 5},
			expected: `[1,2,3,4,5]`,
		},
		{
			name:     "string slice",
			value:    []string{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
		{
			name:     "empty slice",
			value:    []int{},
			expected: `[]`,
		},
		{
			name:     "nil slice",
			value:    []int(nil),
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_StructSlice tests marshaling slices of structs
func TestMarshal_StructSlice(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	result, err := Marshal(people)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := `[{"age":30,"name":"Alice"},{"age":25,"name":"Bob"}]`
	if string(result) != expected {
		t.Errorf("Marshal() = %s, want %s", string(result), expected)
	}
}

// TestMarshal_Arrays tests marshaling arrays
func TestMarshal_Arrays(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "int array",
			value:    [3]int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "string array",
			value:    [2]string{"a", "b"},
			expected: `["a","b"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_Maps tests marshaling maps
func TestMarshal_Maps(t *testing.T) {
	tests := []struct {
		name         string
		value        interface{}
		expectedKeys []string // Check keys are present
	}{
		{
			name: "string to string map",
			value: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expectedKeys: []string{`"key1":"value1"`, `"key2":"value2"`},
		},
		{
			name: "string to int map",
			value: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			expectedKeys: []string{`"a":1`, `"b":2`, `"c":3`},
		},
		{
			name:         "empty map",
			value:        map[string]string{},
			expectedKeys: []string{},
		},
		{
			name:         "nil map",
			value:        map[string]string(nil),
			expectedKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			resultStr := string(result)

			// Check for nil map
			if tt.expectedKeys == nil {
				if resultStr != "null" {
					t.Errorf("Marshal() = %s, want null", resultStr)
				}
				return
			}

			// Check for empty map
			if len(tt.expectedKeys) == 0 {
				if resultStr != "{}" {
					t.Errorf("Marshal() = %s, want {}", resultStr)
				}
				return
			}

			// Check that all expected keys are present
			if !strings.HasPrefix(resultStr, "{") || !strings.HasSuffix(resultStr, "}") {
				t.Errorf("Marshal() = %s, expected object", resultStr)
			}

			for _, key := range tt.expectedKeys {
				if !strings.Contains(resultStr, key) {
					t.Errorf("Marshal() = %s, missing key %s", resultStr, key)
				}
			}
		})
	}
}

// TestMarshal_Pointers tests marshaling with pointer fields
func TestMarshal_Pointers(t *testing.T) {
	strVal := "Alice"
	intVal := 30

	type Person struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	}

	tests := []struct {
		name     string
		value    Person
		expected string
	}{
		{
			name: "non-nil pointers",
			value: Person{
				Name: &strVal,
				Age:  &intVal,
			},
			expected: `{"age":30,"name":"Alice"}`,
		},
		{
			name: "nil pointers",
			value: Person{
				Name: nil,
				Age:  nil,
			},
			expected: `{"age":null,"name":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

// TestMarshal_Interface tests marshaling interface{} values
func TestMarshal_Interface(t *testing.T) {
	t.Run("map in interface", func(t *testing.T) {
		value := map[string]interface{}{
			"name": "Alice",
			"age":  30,
		}
		result, err := Marshal(value)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		resultStr := string(result)
		if !strings.HasPrefix(resultStr, "{") || !strings.HasSuffix(resultStr, "}") {
			t.Errorf("Marshal() = %s, expected object", resultStr)
		}
		for _, key := range []string{`"name":"Alice"`, `"age":30`} {
			if !strings.Contains(resultStr, key) {
				t.Errorf("Marshal() = %s, missing key %s", resultStr, key)
			}
		}
	})

	t.Run("slice in interface", func(t *testing.T) {
		value := []interface{}{1, "two", true}
		result, err := Marshal(value)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		expected := `[1,"two",true]`
		if string(result) != expected {
			t.Errorf("Marshal() = %s, want %s", string(result), expected)
		}
	})
}

// TestMarshal_OmitEmpty tests omitempty behavior
func TestMarshal_OmitEmpty(t *testing.T) {
	type Test struct {
		String     string         `json:"string,omitempty"`
		Int        int            `json:"int,omitempty"`
		Bool       bool           `json:"bool,omitempty"`
		Slice      []int          `json:"slice,omitempty"`
		Map        map[string]int `json:"map,omitempty"`
		Ptr        *string        `json:"ptr,omitempty"`
		AlwaysShow string         `json:"always"`
	}

	value := Test{
		String:     "",
		Int:        0,
		Bool:       false,
		Slice:      nil,
		Map:        nil,
		Ptr:        nil,
		AlwaysShow: "",
	}

	result, err := Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := `{"always":""}`
	if string(result) != expected {
		t.Errorf("Marshal() = %s, want %s", string(result), expected)
	}

	// Test with non-empty values
	strVal := "test"
	value2 := Test{
		String:     "hello",
		Int:        42,
		Bool:       true,
		Slice:      []int{1, 2, 3},
		Map:        map[string]int{"a": 1},
		Ptr:        &strVal,
		AlwaysShow: "visible",
	}

	result2, err := Marshal(value2)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	resultStr := string(result2)
	// Check that all fields are present
	requiredSubstrings := []string{`"string":"hello"`, `"int":42`, `"bool":true`, `"slice":[1,2,3]`, `"ptr":"test"`, `"always":"visible"`}
	for _, substr := range requiredSubstrings {
		if !strings.Contains(resultStr, substr) {
			t.Errorf("Marshal() = %s, missing %s", resultStr, substr)
		}
	}
}

// TestMarshal_RoundTrip tests marshal/unmarshal round trip
func TestMarshal_RoundTrip(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	original := Person{Name: "Alice", Age: 30}

	// Marshal
	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal
	var decoded Person
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Compare
	if decoded.Name != original.Name || decoded.Age != original.Age {
		t.Errorf("Round trip failed: got %+v, want %+v", decoded, original)
	}
}
