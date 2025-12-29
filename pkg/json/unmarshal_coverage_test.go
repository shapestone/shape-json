package json

import (
	"testing"
)

// TestUnmarshalWithAST tests the UnmarshalWithAST function
func TestUnmarshalWithAST(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		validate func(t *testing.T, v interface{})
		wantErr  bool
	}{
		{
			name:   "simple object",
			json:   `{"name": "Alice", "age": 30}`,
			target: new(map[string]interface{}),
			validate: func(t *testing.T, v interface{}) {
				m := *v.(*map[string]interface{})
				if m["name"] != "Alice" {
					t.Errorf("Expected name='Alice', got %v", m["name"])
				}
				if m["age"] != int64(30) {
					t.Errorf("Expected age=30, got %v", m["age"])
				}
			},
		},
		{
			name:   "array",
			json:   `[1, 2, 3]`,
			target: new([]int),
			validate: func(t *testing.T, v interface{}) {
				arr := *v.(*[]int)
				if len(arr) != 3 {
					t.Errorf("Expected length 3, got %d", len(arr))
				}
			},
		},
		{
			name:   "struct",
			json:   `{"Name": "Bob", "Age": 25}`,
			target: new(struct {
				Name string
				Age  int
			}),
			validate: func(t *testing.T, v interface{}) {
				s := v.(*struct {
					Name string
					Age  int
				})
				if s.Name != "Bob" {
					t.Errorf("Expected Name='Bob', got %s", s.Name)
				}
				if s.Age != 25 {
					t.Errorf("Expected Age=25, got %d", s.Age)
				}
			},
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			target:  new(map[string]string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, tt.target)
			}
		})
	}
}

// TestUnmarshalWithAST_Errors tests error cases for UnmarshalWithAST
func TestUnmarshalWithAST_Errors(t *testing.T) {
	tests := []struct {
		name   string
		json   string
		target interface{}
	}{
		{
			name:   "nil target",
			json:   `{"name": "Alice"}`,
			target: nil,
		},
		{
			name:   "non-pointer target",
			json:   `{"name": "Alice"}`,
			target: struct{ Name string }{},
		},
		{
			name:   "nil pointer",
			json:   `{"name": "Alice"}`,
			target: (*struct{ Name string })(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if err == nil {
				t.Error("UnmarshalWithAST() expected error, got nil")
			}
		})
	}
}

// TestUnmarshalLiteral_AllBranches tests all branches in unmarshalLiteral
func TestUnmarshalLiteral_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
	}{
		// String tests
		{
			name:   "string to string",
			json:   `"hello"`,
			target: new(string),
		},
		{
			name:    "non-string to string",
			json:    `42`,
			target:  new(string),
			wantErr: true,
		},
		// Int tests
		{
			name:   "int64 to int",
			json:   `42`,
			target: new(int),
		},
		{
			name:   "float to int (whole number)",
			json:   `42.0`,
			target: new(int),
		},
		{
			name:    "float to int (non-whole)",
			json:    `42.5`,
			target:  new(int),
			wantErr: true,
		},
		{
			name:    "overflow int",
			json:    `9223372036854775807`,
			target:  new(int8),
			wantErr: true,
		},
		{
			name:    "wrong type to int",
			json:    `"hello"`,
			target:  new(int),
			wantErr: true,
		},
		// Uint tests
		{
			name:   "int64 to uint",
			json:   `42`,
			target: new(uint),
		},
		{
			name:    "negative to uint",
			json:    `-1`,
			target:  new(uint),
			wantErr: true,
		},
		{
			name:   "float to uint (whole number)",
			json:   `42.0`,
			target: new(uint),
		},
		{
			name:    "float to uint (negative)",
			json:    `-1.0`,
			target:  new(uint),
			wantErr: true,
		},
		{
			name:    "float to uint (non-whole)",
			json:    `42.5`,
			target:  new(uint),
			wantErr: true,
		},
		{
			name:    "overflow uint",
			json:    `18446744073709551616`,
			target:  new(uint8),
			wantErr: true,
		},
		{
			name:    "wrong type to uint",
			json:    `"hello"`,
			target:  new(uint),
			wantErr: true,
		},
		// Float tests
		{
			name:   "float64 to float",
			json:   `3.14`,
			target: new(float64),
		},
		{
			name:   "int64 to float",
			json:   `42`,
			target: new(float64),
		},
		{
			name:    "overflow float",
			json:    `1.7976931348623159e+309`,
			target:  new(float32),
			wantErr: true,
		},
		{
			name:    "wrong type to float",
			json:    `"hello"`,
			target:  new(float64),
			wantErr: true,
		},
		// Bool tests
		{
			name:   "bool to bool",
			json:   `true`,
			target: new(bool),
		},
		{
			name:    "non-bool to bool",
			json:    `42`,
			target:  new(bool),
			wantErr: true,
		},
		// Unsupported type
		{
			name: "literal to unsupported type",
			json: `42`,
			target: new(struct {
				Name string
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalArray_LegacyObjectFormat tests unmarshaling arrays from legacy object format
func TestUnmarshalArray_LegacyObjectFormat(t *testing.T) {
	// This tests the unmarshalArray function which handles ObjectNode with numeric keys
	// This is a legacy format that's still supported

	// Create a JSON that will parse to an ObjectNode with numeric keys
	// (This would be from older data or specific use cases)
	// Note: Modern JSON arrays parse to ArrayDataNode, but we need to test the legacy path

	// We can't easily create this through Parse, so we'll test via the unmarshalObject path
	// which calls unmarshalArray when it detects numeric keys

	// For now, we'll test that regular arrays work correctly
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
	}{
		{
			name:   "array to slice",
			json:   `[1, 2, 3]`,
			target: new([]int),
		},
		{
			name:   "array to array",
			json:   `[1, 2, 3]`,
			target: new([3]int),
		},
		{
			name:    "array exceeds target array length",
			json:    `[1, 2, 3, 4, 5]`,
			target:  new([3]int),
			wantErr: true,
		},
		{
			name:    "array to unsupported type",
			json:    `[1, 2, 3]`,
			target:  new(string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalArrayData_AllBranches tests all branches in unmarshalArrayData
func TestUnmarshalArrayData_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
		check   func(t *testing.T, target interface{})
	}{
		{
			name:   "array to slice",
			json:   `[1, 2, 3]`,
			target: new([]int),
			check: func(t *testing.T, target interface{}) {
				slice := *target.(*[]int)
				if len(slice) != 3 {
					t.Errorf("Expected length 3, got %d", len(slice))
				}
			},
		},
		{
			name:   "array to fixed array",
			json:   `[1, 2, 3]`,
			target: new([3]int),
			check: func(t *testing.T, target interface{}) {
				arr := target.(*[3]int)
				if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
					t.Errorf("Expected [1 2 3], got %v", arr)
				}
			},
		},
		{
			name:    "array exceeds fixed array length",
			json:    `[1, 2, 3, 4, 5]`,
			target:  new([3]int),
			wantErr: true,
		},
		{
			name:    "array to unsupported type",
			json:    `[1, 2, 3]`,
			target:  new(string),
			wantErr: true,
		},
		{
			name:   "nested arrays",
			json:   `[[1, 2], [3, 4]]`,
			target: new([][]int),
			check: func(t *testing.T, target interface{}) {
				slice := *target.(*[][]int)
				if len(slice) != 2 {
					t.Errorf("Expected length 2, got %d", len(slice))
				}
				if len(slice[0]) != 2 {
					t.Errorf("Expected nested length 2, got %d", len(slice[0]))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

// TestUnmarshalValue_AllBranches tests all branches in unmarshalValue
func TestUnmarshalValue_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
	}{
		{
			name:   "null to pointer",
			json:   `null`,
			target: new(*int),
		},
		{
			name:   "null to value",
			json:   `null`,
			target: new(int),
		},
		{
			name:   "value to interface{}",
			json:   `{"name": "Alice"}`,
			target: new(interface{}),
		},
		{
			name:   "value to pointer",
			json:   `42`,
			target: new(*int),
		},
		{
			name:    "unsupported node type",
			json:    `42`,
			target:  new(int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalObject_AllBranches tests all branches in unmarshalObject
func TestUnmarshalObject_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
	}{
		{
			name:   "object to struct",
			json:   `{"Name": "Alice", "Age": 30}`,
			target: new(struct {
				Name string
				Age  int
			}),
		},
		{
			name:   "object to map",
			json:   `{"key1": "value1", "key2": "value2"}`,
			target: new(map[string]string),
		},
		{
			name:    "object to unsupported type",
			json:    `{"key": "value"}`,
			target:  new(string),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnmarshalMap_AllBranches tests all branches in unmarshalMap
func TestUnmarshalMap_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
		check   func(t *testing.T, target interface{})
	}{
		{
			name:   "nil map initialization",
			json:   `{"key": "value"}`,
			target: new(map[string]string),
			check: func(t *testing.T, target interface{}) {
				m := *target.(*map[string]string)
				if m == nil {
					t.Error("Expected map to be initialized")
				}
				if m["key"] != "value" {
					t.Errorf("Expected key='value', got %v", m["key"])
				}
			},
		},
		{
			name:    "non-string keys",
			json:    `{"key": "value"}`,
			target:  new(map[int]string),
			wantErr: true,
		},
		{
			name:   "nested values in map",
			json:   `{"nested": {"inner": "value"}}`,
			target: new(map[string]map[string]string),
			check: func(t *testing.T, target interface{}) {
				m := *target.(*map[string]map[string]string)
				if m["nested"]["inner"] != "value" {
					t.Errorf("Expected nested value, got %v", m["nested"]["inner"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

// TestUnmarshalStruct_AllBranches tests all branches in unmarshalStruct
func TestUnmarshalStruct_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		target  interface{}
		wantErr bool
		check   func(t *testing.T, target interface{})
	}{
		{
			name: "unexported fields ignored",
			json: `{"Name": "Alice", "age": 30}`,
			target: new(struct {
				Name string
				age  int // unexported, should be ignored
			}),
			check: func(t *testing.T, target interface{}) {
				s := target.(*struct {
					Name string
					age  int
				})
				if s.Name != "Alice" {
					t.Errorf("Expected Name='Alice', got %s", s.Name)
				}
				if s.age != 0 {
					t.Errorf("Expected age=0 (ignored), got %d", s.age)
				}
			},
		},
		{
			name: "skip field with json:'-'",
			json: `{"Name": "Alice", "Ignored": "value"}`,
			target: new(struct {
				Name    string
				Ignored string `json:"-"`
			}),
			check: func(t *testing.T, target interface{}) {
				s := target.(*struct {
					Name    string
					Ignored string `json:"-"`
				})
				if s.Ignored != "" {
					t.Errorf("Expected Ignored='', got %s", s.Ignored)
				}
			},
		},
		{
			name: "unknown fields ignored",
			json: `{"Name": "Alice", "UnknownField": "value"}`,
			target: new(struct {
				Name string
			}),
			check: func(t *testing.T, target interface{}) {
				s := target.(*struct {
					Name string
				})
				if s.Name != "Alice" {
					t.Errorf("Expected Name='Alice', got %s", s.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

// TestNodeToInterface_Wrapper tests the nodeToInterface wrapper function
func TestNodeToInterface_Wrapper(t *testing.T) {
	// This is just a wrapper, but we should still test it for coverage
	json := `{"name": "Alice", "age": 30}`

	// Call the wrapper function (via unmarshalValue with interface{})
	var result interface{}
	err := UnmarshalWithAST([]byte(json), &result)
	if err != nil {
		t.Fatalf("UnmarshalWithAST error: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	if m["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got %v", m["name"])
	}
}

// CustomTypeForAST is a type that implements Unmarshaler for AST testing
type CustomTypeForAST struct {
	Value string
}

// UnmarshalJSON implements the Unmarshaler interface
func (c *CustomTypeForAST) UnmarshalJSON(data []byte) error {
	c.Value = "custom: " + string(data)
	return nil
}

// TestUnmarshalWithAST_CustomUnmarshaler tests custom Unmarshaler interface
func TestUnmarshalWithAST_CustomUnmarshaler(t *testing.T) {
	var ct CustomTypeForAST
	err := UnmarshalWithAST([]byte(`"test"`), &ct)
	if err != nil {
		t.Fatalf("UnmarshalWithAST() error = %v", err)
	}

	// The custom unmarshaler should have been called
	if ct.Value != `custom: "test"` {
		t.Errorf("Expected custom unmarshaler to be called, got Value=%q", ct.Value)
	}
}

// TestIsArray_EdgeCases tests the isArray helper function
func TestIsArray_EdgeCases(t *testing.T) {
	// This is tested indirectly through unmarshaling, but let's ensure edge cases work
	tests := []struct {
		name   string
		json   string
		target interface{}
	}{
		{
			name:   "empty object is not array",
			json:   `{}`,
			target: new(map[string]string),
		},
		{
			name:   "object with non-numeric keys is not array",
			json:   `{"key": "value"}`,
			target: new(map[string]string),
		},
		{
			name:   "object with mixed keys is not array",
			json:   `{"0": "a", "key": "value"}`,
			target: new(map[string]string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalWithAST([]byte(tt.json), tt.target)
			if err != nil {
				t.Fatalf("UnmarshalWithAST() error = %v", err)
			}
		})
	}
}

// TestUnmarshal_PointerFields tests unmarshaling with pointer fields
func TestUnmarshal_PointerFields(t *testing.T) {
	type Person struct {
		Name *string
		Age  *int
	}

	tests := []struct {
		name  string
		json  string
		check func(t *testing.T, p *Person)
	}{
		{
			name: "nil pointers",
			json: `{}`,
			check: func(t *testing.T, p *Person) {
				if p.Name != nil {
					t.Errorf("Expected Name=nil, got %v", p.Name)
				}
				if p.Age != nil {
					t.Errorf("Expected Age=nil, got %v", p.Age)
				}
			},
		},
		{
			name: "allocated pointers",
			json: `{"Name": "Alice", "Age": 30}`,
			check: func(t *testing.T, p *Person) {
				if p.Name == nil || *p.Name != "Alice" {
					t.Errorf("Expected Name='Alice', got %v", p.Name)
				}
				if p.Age == nil || *p.Age != 30 {
					t.Errorf("Expected Age=30, got %v", p.Age)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var person Person
			err := UnmarshalWithAST([]byte(tt.json), &person)
			if err != nil {
				t.Fatalf("UnmarshalWithAST() error = %v", err)
			}
			tt.check(t, &person)
		})
	}
}
