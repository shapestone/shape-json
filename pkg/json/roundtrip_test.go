package json

import (
	"reflect"
	"testing"
)

// TestEmptyArrayRoundTrip tests that empty arrays maintain their type through marshal/unmarshal cycles.
// This test validates the fix for the architectural issue where empty arrays were indistinguishable from
// empty objects, breaking JSON's type system invariants.
func TestEmptyArrayRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name: "empty array in map",
			input: map[string]interface{}{
				"items": []interface{}{},
			},
		},
		{
			name: "multiple empty arrays",
			input: map[string]interface{}{
				"array1": []interface{}{},
				"array2": []interface{}{},
				"nested": map[string]interface{}{
					"emptyArr": []interface{}{},
				},
			},
		},
		{
			name: "empty array and empty object",
			input: map[string]interface{}{
				"emptyArray":  []interface{}{},
				"emptyObject": map[string]interface{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			jsonBytes, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal
			var result map[string]interface{}
			err = Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify type preservation
			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Round-trip failed:\ninput:  %#v\nresult: %#v\nJSON:   %s", tt.input, result, string(jsonBytes))
			}

			// Specifically check that arrays are still arrays, not maps
			checkArraysAreArrays(t, "", tt.input, result)
		})
	}
}

// checkArraysAreArrays recursively verifies that all arrays in the original remain arrays in the result
func checkArraysAreArrays(t *testing.T, path string, original, result interface{}) {
	switch orig := original.(type) {
	case []interface{}:
		// Original was an array - result MUST be an array
		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("At path %q: expected []interface{}, got %T (value: %v)", path, result, result)
			return
		}
		// Recursively check elements
		for i := range orig {
			if i < len(res) {
				checkArraysAreArrays(t, path+"["+string(rune('0'+i))+"]", orig[i], res[i])
			}
		}

	case map[string]interface{}:
		// Original was a map - result MUST be a map
		res, ok := result.(map[string]interface{})
		if !ok {
			t.Errorf("At path %q: expected map[string]interface{}, got %T (value: %v)", path, result, result)
			return
		}
		// Recursively check values
		for key, val := range orig {
			newPath := path + "." + key
			if path == "" {
				newPath = key
			}
			if resVal, exists := res[key]; exists {
				checkArraysAreArrays(t, newPath, val, resVal)
			}
		}
	}
}

// TestEmptyArrayJSONRepresentation verifies that empty arrays render as [] not {}
func TestEmptyArrayJSONRepresentation(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "standalone empty array",
			input:    []interface{}{},
			expected: `[]`,
		},
		{
			name:     "empty array in struct",
			input:    struct{ Items []int }{Items: []int{}},
			expected: `{"Items":[]}`,
		},
		{
			name: "empty array in map",
			input: map[string]interface{}{
				"data": []interface{}{},
			},
			expected: `{"data":[]}`,
		},
		{
			name: "empty object vs empty array",
			input: map[string]interface{}{
				"arr": []interface{}{},
				"obj": map[string]interface{}{},
			},
			expected: `{"arr":[],"obj":{}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Marshal() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

// TestArrayTypeFidelity tests that []interface{} and map[string]interface{} are never confused
func TestArrayTypeFidelity(t *testing.T) {
	// Create a structure with both empty arrays and empty objects
	data := struct {
		EmptyArray  []interface{}          `json:"emptyArray"`
		EmptyObject map[string]interface{} `json:"emptyObject"`
		FilledArray []int                  `json:"filledArray"`
	}{
		EmptyArray:  []interface{}{},
		EmptyObject: map[string]interface{}{},
		FilledArray: []int{1, 2, 3},
	}

	// Marshal
	jsonBytes, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify JSON representation
	expectedJSON := `{"emptyArray":[],"emptyObject":{},"filledArray":[1,2,3]}`
	if string(jsonBytes) != expectedJSON {
		t.Errorf("JSON representation:\ngot:  %s\nwant: %s", string(jsonBytes), expectedJSON)
	}

	// Unmarshal back
	var result map[string]interface{}
	err = Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Verify emptyArray is a slice
	if emptyArr, ok := result["emptyArray"]; !ok {
		t.Error("emptyArray key not found")
	} else if _, isSlice := emptyArr.([]interface{}); !isSlice {
		t.Errorf("emptyArray has wrong type: got %T, want []interface{}", emptyArr)
	}

	// Verify emptyObject is a map
	if emptyObj, ok := result["emptyObject"]; !ok {
		t.Error("emptyObject key not found")
	} else if _, isMap := emptyObj.(map[string]interface{}); !isMap {
		t.Errorf("emptyObject has wrong type: got %T, want map[string]interface{}", emptyObj)
	}

	// Verify filledArray is a slice
	if filledArr, ok := result["filledArray"]; !ok {
		t.Error("filledArray key not found")
	} else if _, isSlice := filledArr.([]interface{}); !isSlice {
		t.Errorf("filledArray has wrong type: got %T, want []interface{}", filledArr)
	}
}
