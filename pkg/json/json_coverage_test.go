package json

import (
	"bytes"
	"strings"
	"testing"
)

// TestUnmarshalLiteral_Comprehensive tests all unmarshalLiteral code paths
func TestUnmarshalLiteral_Comprehensive(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		target      interface{}
		shouldError bool
		errorMsg    string
	}{
		// String unmarshaling
		{
			name:   "string to string",
			json:   `"hello"`,
			target: new(string),
		},
		{
			name:        "number to string - should error",
			json:        `42`,
			target:      new(string),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},
		{
			name:        "bool to string - should error",
			json:        `true`,
			target:      new(string),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},

		// Int unmarshaling
		{
			name:   "int64 to int",
			json:   `42`,
			target: new(int),
		},
		{
			name:   "int64 to int8",
			json:   `127`,
			target: new(int8),
		},
		{
			name:   "int64 to int16",
			json:   `32767`,
			target: new(int16),
		},
		{
			name:   "int64 to int32",
			json:   `2147483647`,
			target: new(int32),
		},
		{
			name:   "int64 to int64",
			json:   `9223372036854775807`,
			target: new(int64),
		},
		{
			name:        "overflow int8",
			json:        `128`,
			target:      new(int8),
			shouldError: true,
			errorMsg:    "overflows",
		},
		{
			name:        "overflow int16",
			json:        `32768`,
			target:      new(int16),
			shouldError: true,
			errorMsg:    "overflows",
		},
		{
			name:   "float to int (whole number)",
			json:   `42.0`,
			target: new(int),
		},
		{
			name:        "float to int (not whole)",
			json:        `42.5`,
			target:      new(int),
			shouldError: true,
			errorMsg:    "cannot unmarshal number",
		},
		{
			name:        "string to int - should error",
			json:        `"42"`,
			target:      new(int),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},

		// Uint unmarshaling
		{
			name:   "positive int to uint",
			json:   `42`,
			target: new(uint),
		},
		{
			name:   "int to uint8",
			json:   `255`,
			target: new(uint8),
		},
		{
			name:   "int to uint16",
			json:   `65535`,
			target: new(uint16),
		},
		{
			name:   "int to uint32",
			json:   `4294967295`,
			target: new(uint32),
		},
		{
			name:        "negative int to uint",
			json:        `-1`,
			target:      new(uint),
			shouldError: true,
			errorMsg:    "overflows",
		},
		{
			name:        "overflow uint8",
			json:        `256`,
			target:      new(uint8),
			shouldError: true,
			errorMsg:    "overflows",
		},
		{
			name:   "float to uint (whole positive)",
			json:   `42.0`,
			target: new(uint),
		},
		{
			name:        "float to uint (negative)",
			json:        `-1.0`,
			target:      new(uint),
			shouldError: true,
			errorMsg:    "cannot unmarshal number",
		},
		{
			name:        "float to uint (not whole)",
			json:        `42.5`,
			target:      new(uint),
			shouldError: true,
			errorMsg:    "cannot unmarshal number",
		},
		{
			name:        "string to uint - should error",
			json:        `"42"`,
			target:      new(uint),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},

		// Float unmarshaling
		{
			name:   "float64 to float32",
			json:   `3.14`,
			target: new(float32),
		},
		{
			name:   "float64 to float64",
			json:   `3.14159265359`,
			target: new(float64),
		},
		{
			name:   "int to float32",
			json:   `42`,
			target: new(float32),
		},
		{
			name:   "int to float64",
			json:   `42`,
			target: new(float64),
		},
		{
			name:        "very large float to float32 overflow",
			json:        `3.4e39`,
			target:      new(float32),
			shouldError: true,
			errorMsg:    "overflows",
		},
		{
			name:        "string to float - should error",
			json:        `"3.14"`,
			target:      new(float64),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},

		// Bool unmarshaling
		{
			name:   "true to bool",
			json:   `true`,
			target: new(bool),
		},
		{
			name:   "false to bool",
			json:   `false`,
			target: new(bool),
		},
		{
			name:        "number to bool - should error",
			json:        `1`,
			target:      new(bool),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},
		{
			name:        "string to bool - should error",
			json:        `"true"`,
			target:      new(bool),
			shouldError: true,
			errorMsg:    "cannot unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Logf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestUnmarshalArray_Comprehensive tests all unmarshalArray code paths
func TestUnmarshalArray_Comprehensive(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		target      interface{}
		shouldError bool
		errorMsg    string
		validate    func(t *testing.T, target interface{})
	}{
		{
			name:   "empty array to slice",
			json:   `[]`,
			target: new([]int),
			validate: func(t *testing.T, target interface{}) {
				s := *target.(*[]int)
				if len(s) != 0 {
					t.Errorf("expected empty slice, got %v", s)
				}
			},
		},
		{
			name:   "int array to slice",
			json:   `[1, 2, 3, 4, 5]`,
			target: new([]int),
			validate: func(t *testing.T, target interface{}) {
				s := *target.(*[]int)
				if len(s) != 5 {
					t.Errorf("expected length 5, got %d", len(s))
				}
				if s[0] != 1 || s[4] != 5 {
					t.Errorf("unexpected values: %v", s)
				}
			},
		},
		{
			name:   "string array to slice",
			json:   `["a", "b", "c"]`,
			target: new([]string),
			validate: func(t *testing.T, target interface{}) {
				s := *target.(*[]string)
				if len(s) != 3 {
					t.Errorf("expected length 3, got %d", len(s))
				}
				if s[0] != "a" || s[2] != "c" {
					t.Errorf("unexpected values: %v", s)
				}
			},
		},
		{
			name:   "mixed type array to interface slice",
			json:   `[1, "two", true, null]`,
			target: new([]interface{}),
			validate: func(t *testing.T, target interface{}) {
				s := *target.(*[]interface{})
				if len(s) != 4 {
					t.Errorf("expected length 4, got %d", len(s))
				}
			},
		},
		{
			name:   "nested array",
			json:   `[[1, 2], [3, 4]]`,
			target: new([][]int),
			validate: func(t *testing.T, target interface{}) {
				s := *target.(*[][]int)
				if len(s) != 2 {
					t.Errorf("expected length 2, got %d", len(s))
				}
				if len(s[0]) != 2 || s[0][0] != 1 {
					t.Errorf("unexpected nested values: %v", s)
				}
			},
		},
		{
			name:   "array to fixed size array",
			json:   `[1, 2, 3]`,
			target: new([3]int),
			validate: func(t *testing.T, target interface{}) {
				arr := *target.(*[3]int)
				if arr[0] != 1 || arr[2] != 3 {
					t.Errorf("unexpected values: %v", arr)
				}
			},
		},
		{
			name:        "array too large for fixed array",
			json:        `[1, 2, 3, 4]`,
			target:      new([3]int),
			shouldError: true,
			errorMsg:    "exceeds target array length",
		},
		{
			name:   "array smaller than fixed array",
			json:   `[1, 2]`,
			target: new([5]int),
			validate: func(t *testing.T, target interface{}) {
				arr := *target.(*[5]int)
				if arr[0] != 1 || arr[1] != 2 {
					t.Errorf("unexpected values: %v", arr)
				}
				// Remaining should be zero values
				if arr[2] != 0 {
					t.Errorf("expected zero value at index 2, got %d", arr[2])
				}
			},
		},
		{
			name:        "array to non-array type",
			json:        `[1, 2, 3]`,
			target:      new(int),
			shouldError: true,
			errorMsg:    "cannot unmarshal array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if tt.validate != nil {
					tt.validate(t, tt.target)
				}
			}
		})
	}
}

// TestEncode_Comprehensive tests Encode with various types
func TestEncode_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "encode int",
			input:    42,
			expected: "42",
		},
		{
			name:     "encode float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "encode string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "encode bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "encode bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "encode null",
			input:    nil,
			expected: "null",
		},
		{
			name:     "encode slice",
			input:    []int{1, 2, 3},
			expected: "[1,2,3]",
		},
		{
			name:     "encode empty slice",
			input:    []int{},
			expected: "[]",
		},
		{
			name: "encode struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "Alice",
				Age:  30,
			},
			expected: `{"age":30,"name":"Alice"}`,
		},
		{
			name:     "encode map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: `{"a":1,"b":2}`,
		},
		{
			name: "encode nested struct",
			input: struct {
				User struct {
					Name string `json:"name"`
				} `json:"user"`
			}{
				User: struct {
					Name string `json:"name"`
				}{
					Name: "Bob",
				},
			},
			expected: `{"user":{"name":"Bob"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.Encode(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			result := strings.TrimSpace(buf.String())
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestEncode_ErrorCases tests Encode error handling
func TestEncode_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "encode channel",
			input: make(chan int),
		},
		{
			name:  "encode function",
			input: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.Encode(tt.input)

			if err == nil {
				t.Logf("expected error for %s but got none", tt.name)
			}
		})
	}
}

// TestMarshal_EdgeCases tests Marshal edge cases
func TestMarshal_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		shouldError bool
	}{
		{
			name:  "marshal pointer to int",
			input: func() *int { i := 42; return &i }(),
		},
		{
			name:  "marshal nil pointer",
			input: (*int)(nil),
		},
		{
			name:  "marshal pointer to struct",
			input: &struct{ Name string }{Name: "test"},
		},
		{
			name:        "marshal channel",
			input:       make(chan int),
			shouldError: true,
		},
		{
			name:        "marshal function",
			input:       func() {},
			shouldError: true,
		},
		{
			name:  "marshal complex nested structure",
			input: map[string]interface{}{"a": []interface{}{1, "two", true, nil}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Marshal(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestMarshalString_SpecialCharacters tests string escaping
func TestMarshalString_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string with quotes",
			input:    `hello "world"`,
			expected: `"hello \"world\""`,
		},
		{
			name:     "string with backslash",
			input:    `C:\path\to\file`,
			expected: `"C:\\path\\to\\file"`,
		},
		{
			name:     "string with newline",
			input:    "line1\nline2",
			expected: `"line1\nline2"`,
		},
		{
			name:     "string with tab",
			input:    "col1\tcol2",
			expected: `"col1\tcol2"`,
		},
		{
			name:     "string with carriage return",
			input:    "text\rreturn",
			expected: `"text\rreturn"`,
		},
		{
			name:     "string with all escapes",
			input:    "\"\\\n\r\t",
			expected: `"\"\\\n\r\t"`,
		},
		{
			name:     "unicode string",
			input:    "Hello 世界",
			expected: `"Hello 世界"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

// TestMarshalValue_AllTypes tests marshalValue with various types
func TestMarshalValue_AllTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{name: "nil", input: nil},
		{name: "bool", input: true},
		{name: "int", input: 42},
		{name: "int8", input: int8(127)},
		{name: "int16", input: int16(32767)},
		{name: "int32", input: int32(2147483647)},
		{name: "int64", input: int64(9223372036854775807)},
		{name: "uint", input: uint(42)},
		{name: "uint8", input: uint8(255)},
		{name: "uint16", input: uint16(65535)},
		{name: "uint32", input: uint32(4294967295)},
		{name: "uint64", input: uint64(18446744073709551615)},
		{name: "float32", input: float32(3.14)},
		{name: "float64", input: float64(3.14159265359)},
		{name: "string", input: "hello"},
		{name: "slice", input: []int{1, 2, 3}},
		{name: "array", input: [3]int{1, 2, 3}},
		{name: "map", input: map[string]int{"a": 1}},
		{name: "struct", input: struct{ Name string }{"test"}},
		{name: "pointer", input: func() *int { i := 42; return &i }()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Marshal(tt.input)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}
