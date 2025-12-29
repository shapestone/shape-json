package fastparser

import (
	"testing"
)

// TestUnmarshalFixedArray tests unmarshaling into fixed-size arrays
func TestUnmarshalFixedArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:   "unmarshal to int array",
			input:  `[1, 2, 3]`,
			target: &[3]int{},
		},
		{
			name:   "unmarshal to string array",
			input:  `["a", "b", "c"]`,
			target: &[3]string{},
		},
		{
			name:   "unmarshal empty array to fixed array",
			input:  `[]`,
			target: &[3]int{},
		},
		{
			name:    "unmarshal too many elements",
			input:   `[1, 2, 3, 4, 5]`,
			target:  &[3]int{},
			wantErr: true,
		},
		{
			name:   "unmarshal fewer elements",
			input:  `[1, 2]`,
			target: &[3]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestUnmarshalNumber_AllBranches tests all number unmarshaling branches
func TestUnmarshalNumber_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
		check   func(t *testing.T, target interface{})
	}{
		{
			name:   "int8",
			input:  `42`,
			target: new(int8),
			check: func(t *testing.T, target interface{}) {
				if *target.(*int8) != 42 {
					t.Errorf("expected 42, got %v", *target.(*int8))
				}
			},
		},
		{
			name:   "int16",
			input:  `1000`,
			target: new(int16),
			check: func(t *testing.T, target interface{}) {
				if *target.(*int16) != 1000 {
					t.Errorf("expected 1000, got %v", *target.(*int16))
				}
			},
		},
		{
			name:   "int32",
			input:  `100000`,
			target: new(int32),
			check: func(t *testing.T, target interface{}) {
				if *target.(*int32) != 100000 {
					t.Errorf("expected 100000, got %v", *target.(*int32))
				}
			},
		},
		{
			name:   "uint",
			input:  `42`,
			target: new(uint),
			check: func(t *testing.T, target interface{}) {
				if *target.(*uint) != 42 {
					t.Errorf("expected 42, got %v", *target.(*uint))
				}
			},
		},
		{
			name:   "uint8",
			input:  `255`,
			target: new(uint8),
			check: func(t *testing.T, target interface{}) {
				if *target.(*uint8) != 255 {
					t.Errorf("expected 255, got %v", *target.(*uint8))
				}
			},
		},
		{
			name:   "uint16",
			input:  `65535`,
			target: new(uint16),
			check: func(t *testing.T, target interface{}) {
				if *target.(*uint16) != 65535 {
					t.Errorf("expected 65535, got %v", *target.(*uint16))
				}
			},
		},
		{
			name:   "uint32",
			input:  `4294967295`,
			target: new(uint32),
			check: func(t *testing.T, target interface{}) {
				if *target.(*uint32) != 4294967295 {
					t.Errorf("expected 4294967295, got %v", *target.(*uint32))
				}
			},
		},
		{
			name:   "uint64",
			input:  `18446744073709551615`,
			target: new(uint64),
			check: func(t *testing.T, target interface{}) {
				if *target.(*uint64) != 18446744073709551615 {
					t.Errorf("expected 18446744073709551615, got %v", *target.(*uint64))
				}
			},
		},
		{
			name:    "negative to uint",
			input:   `-1`,
			target:  new(uint),
			wantErr: true,
		},
		{
			name:   "float32",
			input:  `3.14`,
			target: new(float32),
			check: func(t *testing.T, target interface{}) {
				if *target.(*float32) < 3.13 || *target.(*float32) > 3.15 {
					t.Errorf("expected ~3.14, got %v", *target.(*float32))
				}
			},
		},
		{
			name:   "int to float",
			input:  `42`,
			target: new(float64),
			check: func(t *testing.T, target interface{}) {
				if *target.(*float64) != 42.0 {
					t.Errorf("expected 42.0, got %v", *target.(*float64))
				}
			},
		},
		{
			name:    "overflow int8",
			input:   `128`,
			target:  new(int8),
			wantErr: true,
		},
		{
			name:    "overflow uint8",
			input:   `256`,
			target:  new(uint8),
			wantErr: true,
		},
		{
			name:    "float to int non-integer",
			input:   `3.14`,
			target:  new(int),
			wantErr: true,
		},
		{
			name:    "float to uint non-integer",
			input:   `3.14`,
			target:  new(uint),
			wantErr: true,
		},
		{
			name:    "unmarshal number to unsupported type",
			input:   `42`,
			target:  new(string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

// TestUnmarshalArray_ErrorPaths tests error paths in array unmarshaling
func TestUnmarshalArray_ErrorPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "array to non-slice/array type",
			input:   `[1, 2, 3]`,
			target:  new(string),
			wantErr: true,
		},
		{
			name:    "unclosed array",
			input:   `[1, 2, 3`,
			target:  new([]int),
			wantErr: true,
		},
		{
			name:    "array with invalid separator",
			input:   `[1; 2; 3]`,
			target:  new([]int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalValue_EdgeCases tests edge cases in value unmarshaling
func TestUnmarshalValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   ``,
			target:  new(string),
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   `   `,
			target:  new(string),
			wantErr: true,
		},
		{
			name:    "unexpected character",
			input:   `@`,
			target:  new(string),
			wantErr: true,
		},
		{
			name:   "null to pointer",
			input:  `null`,
			target: new(*int),
		},
		{
			name:   "null to interface",
			input:  `null`,
			target: new(interface{}),
		},
		{
			name:   "nested pointer",
			input:  `42`,
			target: new(*int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalObject_EdgeCases tests edge cases in object unmarshaling
func TestUnmarshalObject_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "object to non-struct/map type",
			input:   `{"key": "value"}`,
			target:  new(string),
			wantErr: true,
		},
		{
			name:   "object with JSON tag '-'",
			input:  `{"name": "Alice", "ignored": "value"}`,
			target: &struct {
				Name    string
				Ignored string `json:"-"`
			}{},
		},
		{
			name:   "object with comma in tag",
			input:  `{"custom_name": "Alice"}`,
			target: &struct {
				Name string `json:"custom_name,omitempty"`
			}{},
		},
		{
			name:    "object missing opening brace",
			input:   `"key": "value"}`,
			target:  new(map[string]string),
			wantErr: true,
		},
		{
			name:    "object with invalid key",
			input:   `{123: "value"}`,
			target:  new(map[string]string),
			wantErr: true,
		},
		{
			name:    "object missing colon",
			input:   `{"key" "value"}`,
			target:  new(map[string]string),
			wantErr: true,
		},
		{
			name:    "object with unclosed value",
			input:   `{"key": "value"`,
			target:  new(map[string]string),
			wantErr: true,
		},
		{
			name:    "map with non-string keys",
			input:   `{"key": "value"}`,
			target:  new(map[int]string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalBool_ErrorPaths tests bool unmarshaling error paths
func TestUnmarshalBool_ErrorPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "bool to non-bool type",
			input:   `true`,
			target:  new(string),
			wantErr: true,
		},
		{
			name:    "invalid bool literal",
			input:   `tru`,
			target:  new(bool),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalString_ErrorPaths tests string unmarshaling error paths
func TestUnmarshalString_ErrorPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "string to non-string type",
			input:   `"hello"`,
			target:  new(int),
			wantErr: true,
		},
		{
			name:    "unclosed string",
			input:   `"hello`,
			target:  new(string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// CustomType is a type that implements Unmarshaler for testing
type CustomType struct {
	Value string
}

// UnmarshalJSON implements the Unmarshaler interface
func (c *CustomType) UnmarshalJSON(data []byte) error {
	c.Value = "custom: " + string(data)
	return nil
}

// TestUnmarshalWithUnmarshaler tests custom Unmarshaler interface
func TestUnmarshalWithUnmarshaler(t *testing.T) {
	var ct CustomType
	err := Unmarshal([]byte(`"test"`), &ct)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// The custom unmarshaler should have been called
	if ct.Value != `custom: "test"` {
		t.Errorf("Expected custom unmarshaler to be called, got Value=%q", ct.Value)
	}
}

// TestExpectLiteral tests expectLiteral function coverage
func TestExpectLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "incomplete null",
			input:   `nul`,
			target:  new(interface{}),
			wantErr: true,
		},
		{
			name:    "incomplete true",
			input:   `tru`,
			target:  new(bool),
			wantErr: true,
		},
		{
			name:    "incomplete false",
			input:   `fals`,
			target:  new(bool),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
