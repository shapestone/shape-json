package json

import (
	"testing"
)

// TestUnmarshal_BasicTypes tests unmarshaling into basic Go types
func TestUnmarshal_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		expected interface{}
	}{
		{
			name:     "string",
			json:     `"hello"`,
			target:   new(string),
			expected: "hello",
		},
		{
			name:     "int",
			json:     `42`,
			target:   new(int),
			expected: 42,
		},
		{
			name:     "int64",
			json:     `42`,
			target:   new(int64),
			expected: int64(42),
		},
		{
			name:     "float64",
			json:     `3.14`,
			target:   new(float64),
			expected: 3.14,
		},
		{
			name:     "bool true",
			json:     `true`,
			target:   new(bool),
			expected: true,
		},
		{
			name:     "bool false",
			json:     `false`,
			target:   new(bool),
			expected: false,
		},
		{
			name:     "null to pointer",
			json:     `null`,
			target:   new(*string),
			expected: (*string)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Extract the actual value from pointer
			var actual interface{}
			switch v := tt.target.(type) {
			case *string:
				actual = *v
			case *int:
				actual = *v
			case *int64:
				actual = *v
			case *float64:
				actual = *v
			case *bool:
				actual = *v
			case **string:
				actual = *v
			}

			if actual != tt.expected {
				t.Errorf("Unmarshal() = %v (%T), want %v (%T)", actual, actual, tt.expected, tt.expected)
			}
		})
	}
}

// TestUnmarshal_Struct tests unmarshaling into structs
func TestUnmarshal_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name     string
		json     string
		target   interface{}
		validate func(t *testing.T, v interface{})
	}{
		{
			name:   "simple struct",
			json:   `{"Name": "Alice", "Age": 30}`,
			target: new(Person),
			validate: func(t *testing.T, v interface{}) {
				p := v.(*Person)
				if p.Name != "Alice" || p.Age != 30 {
					t.Errorf("got Name=%s Age=%d, want Name=Alice Age=30", p.Name, p.Age)
				}
			},
		},
		{
			name:   "partial struct",
			json:   `{"Name": "Bob"}`,
			target: new(Person),
			validate: func(t *testing.T, v interface{}) {
				p := v.(*Person)
				if p.Name != "Bob" || p.Age != 0 {
					t.Errorf("got Name=%s Age=%d, want Name=Bob Age=0", p.Name, p.Age)
				}
			},
		},
		{
			name:   "empty struct",
			json:   `{}`,
			target: new(Person),
			validate: func(t *testing.T, v interface{}) {
				p := v.(*Person)
				if p.Name != "" || p.Age != 0 {
					t.Errorf("got Name=%s Age=%d, want zero values", p.Name, p.Age)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, tt.target)
		})
	}
}

// TestUnmarshal_StructTags tests unmarshaling with json struct tags
func TestUnmarshal_StructTags(t *testing.T) {
	type Tagged struct {
		PublicName  string `json:"name"`
		InternalAge int    `json:"age"`
		Ignored     string `json:"-"`
		NoTag       string
	}

	tests := []struct {
		name     string
		json     string
		validate func(t *testing.T, v *Tagged)
	}{
		{
			name: "tagged fields",
			json: `{"name": "Alice", "age": 30, "NoTag": "visible"}`,
			validate: func(t *testing.T, v *Tagged) {
				if v.PublicName != "Alice" {
					t.Errorf("PublicName = %s, want Alice", v.PublicName)
				}
				if v.InternalAge != 30 {
					t.Errorf("InternalAge = %d, want 30", v.InternalAge)
				}
				if v.NoTag != "visible" {
					t.Errorf("NoTag = %s, want visible", v.NoTag)
				}
			},
		},
		{
			name: "ignored field in json",
			json: `{"name": "Bob", "age": 25, "Ignored": "should not be set"}`,
			validate: func(t *testing.T, v *Tagged) {
				if v.Ignored != "" {
					t.Errorf("Ignored = %s, want empty (should be ignored)", v.Ignored)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target Tagged
			err := Unmarshal([]byte(tt.json), &target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, &target)
		})
	}
}

// TestUnmarshal_NestedStruct tests unmarshaling nested structures
func TestUnmarshal_NestedStruct(t *testing.T) {
	type Address struct {
		City  string
		State string
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
	}

	json := `{
		"Name": "Alice",
		"Age": 30,
		"Address": {
			"City": "Seattle",
			"State": "WA"
		}
	}`

	var person Person
	err := Unmarshal([]byte(json), &person)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if person.Name != "Alice" {
		t.Errorf("Name = %s, want Alice", person.Name)
	}
	if person.Age != 30 {
		t.Errorf("Age = %d, want 30", person.Age)
	}
	if person.Address.City != "Seattle" {
		t.Errorf("Address.City = %s, want Seattle", person.Address.City)
	}
	if person.Address.State != "WA" {
		t.Errorf("Address.State = %s, want WA", person.Address.State)
	}
}

// TestUnmarshal_Slices tests unmarshaling into slices
func TestUnmarshal_Slices(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		validate func(t *testing.T, v interface{})
	}{
		{
			name:   "int slice",
			json:   `[1, 2, 3, 4, 5]`,
			target: new([]int),
			validate: func(t *testing.T, v interface{}) {
				slice := *v.(*[]int)
				if len(slice) != 5 {
					t.Errorf("len = %d, want 5", len(slice))
				}
				for i, want := range []int{1, 2, 3, 4, 5} {
					if slice[i] != want {
						t.Errorf("slice[%d] = %d, want %d", i, slice[i], want)
					}
				}
			},
		},
		{
			name:   "string slice",
			json:   `["a", "b", "c"]`,
			target: new([]string),
			validate: func(t *testing.T, v interface{}) {
				slice := *v.(*[]string)
				if len(slice) != 3 {
					t.Errorf("len = %d, want 3", len(slice))
				}
				for i, want := range []string{"a", "b", "c"} {
					if slice[i] != want {
						t.Errorf("slice[%d] = %s, want %s", i, slice[i], want)
					}
				}
			},
		},
		{
			name:   "empty slice",
			json:   `[]`,
			target: new([]int),
			validate: func(t *testing.T, v interface{}) {
				slice := *v.(*[]int)
				if len(slice) != 0 {
					t.Errorf("len = %d, want 0", len(slice))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, tt.target)
		})
	}
}

// TestUnmarshal_StructSlice tests unmarshaling slices of structs
func TestUnmarshal_StructSlice(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	json := `[
		{"Name": "Alice", "Age": 30},
		{"Name": "Bob", "Age": 25}
	]`

	var people []Person
	err := Unmarshal([]byte(json), &people)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if len(people) != 2 {
		t.Fatalf("len = %d, want 2", len(people))
	}

	if people[0].Name != "Alice" || people[0].Age != 30 {
		t.Errorf("people[0] = %+v, want {Name:Alice Age:30}", people[0])
	}

	if people[1].Name != "Bob" || people[1].Age != 25 {
		t.Errorf("people[1] = %+v, want {Name:Bob Age:25}", people[1])
	}
}

// TestUnmarshal_Maps tests unmarshaling into maps
func TestUnmarshal_Maps(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		validate func(t *testing.T, v interface{})
	}{
		{
			name:   "string to string map",
			json:   `{"key1": "value1", "key2": "value2"}`,
			target: new(map[string]string),
			validate: func(t *testing.T, v interface{}) {
				m := *v.(*map[string]string)
				if len(m) != 2 {
					t.Errorf("len = %d, want 2", len(m))
				}
				if m["key1"] != "value1" {
					t.Errorf("m[key1] = %s, want value1", m["key1"])
				}
				if m["key2"] != "value2" {
					t.Errorf("m[key2] = %s, want value2", m["key2"])
				}
			},
		},
		{
			name:   "string to int map",
			json:   `{"a": 1, "b": 2, "c": 3}`,
			target: new(map[string]int),
			validate: func(t *testing.T, v interface{}) {
				m := *v.(*map[string]int)
				if len(m) != 3 {
					t.Errorf("len = %d, want 3", len(m))
				}
				if m["a"] != 1 || m["b"] != 2 || m["c"] != 3 {
					t.Errorf("map = %v, want {a:1 b:2 c:3}", m)
				}
			},
		},
		{
			name:   "empty map",
			json:   `{}`,
			target: new(map[string]string),
			validate: func(t *testing.T, v interface{}) {
				m := *v.(*map[string]string)
				if len(m) != 0 {
					t.Errorf("len = %d, want 0", len(m))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, tt.target)
		})
	}
}

// TestUnmarshal_Pointers tests unmarshaling with pointer fields
func TestUnmarshal_Pointers(t *testing.T) {
	type Person struct {
		Name *string
		Age  *int
	}

	tests := []struct {
		name     string
		json     string
		validate func(t *testing.T, v *Person)
	}{
		{
			name: "non-null pointers",
			json: `{"Name": "Alice", "Age": 30}`,
			validate: func(t *testing.T, v *Person) {
				if v.Name == nil {
					t.Error("Name is nil, want non-nil")
				} else if *v.Name != "Alice" {
					t.Errorf("*Name = %s, want Alice", *v.Name)
				}
				if v.Age == nil {
					t.Error("Age is nil, want non-nil")
				} else if *v.Age != 30 {
					t.Errorf("*Age = %d, want 30", *v.Age)
				}
			},
		},
		{
			name: "null pointers",
			json: `{"Name": null, "Age": null}`,
			validate: func(t *testing.T, v *Person) {
				if v.Name != nil {
					t.Errorf("Name = %v, want nil", v.Name)
				}
				if v.Age != nil {
					t.Errorf("Age = %v, want nil", v.Age)
				}
			},
		},
		{
			name: "missing fields (should be nil)",
			json: `{}`,
			validate: func(t *testing.T, v *Person) {
				if v.Name != nil {
					t.Errorf("Name = %v, want nil", v.Name)
				}
				if v.Age != nil {
					t.Errorf("Age = %v, want nil", v.Age)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var person Person
			err := Unmarshal([]byte(tt.json), &person)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, &person)
		})
	}
}

// TestUnmarshal_Errors tests error cases
func TestUnmarshal_Errors(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		target      interface{}
		expectError bool
	}{
		{
			name:        "non-pointer target",
			json:        `{"name": "Alice"}`,
			target:      struct{ Name string }{},
			expectError: true,
		},
		{
			name:        "nil target",
			json:        `{"name": "Alice"}`,
			target:      nil,
			expectError: true,
		},
		{
			name:        "invalid json",
			json:        `{invalid}`,
			target:      new(map[string]string),
			expectError: true,
		},
		{
			name:        "type mismatch - string to int",
			json:        `"hello"`,
			target:      new(int),
			expectError: true,
		},
		{
			name:        "type mismatch - number to string",
			json:        `42`,
			target:      new(string),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.json), tt.target)
			if tt.expectError && err == nil {
				t.Error("Unmarshal() error = nil, want error")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unmarshal() error = %v, want nil", err)
			}
		})
	}
}

// TestUnmarshal_Interface tests unmarshaling into interface{}
func TestUnmarshal_Interface(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		validate func(t *testing.T, v interface{})
	}{
		{
			name: "object to interface",
			json: `{"name": "Alice", "age": 30}`,
			validate: func(t *testing.T, v interface{}) {
				m, ok := v.(map[string]interface{})
				if !ok {
					t.Fatalf("type = %T, want map[string]interface{}", v)
				}
				if m["name"] != "Alice" {
					t.Errorf("name = %v, want Alice", m["name"])
				}
				if m["age"] != int64(30) {
					t.Errorf("age = %v, want 30", m["age"])
				}
			},
		},
		{
			name: "array to interface",
			json: `[1, 2, 3]`,
			validate: func(t *testing.T, v interface{}) {
				arr, ok := v.([]interface{})
				if !ok {
					t.Fatalf("type = %T, want []interface{}", v)
				}
				if len(arr) != 3 {
					t.Errorf("len = %d, want 3", len(arr))
				}
			},
		},
		{
			name: "string to interface",
			json: `"hello"`,
			validate: func(t *testing.T, v interface{}) {
				s, ok := v.(string)
				if !ok {
					t.Fatalf("type = %T, want string", v)
				}
				if s != "hello" {
					t.Errorf("value = %s, want hello", s)
				}
			},
		},
		{
			name: "number to interface",
			json: `42`,
			validate: func(t *testing.T, v interface{}) {
				n, ok := v.(int64)
				if !ok {
					t.Fatalf("type = %T, want int64", v)
				}
				if n != 42 {
					t.Errorf("value = %d, want 42", n)
				}
			},
		},
		{
			name: "bool to interface",
			json: `true`,
			validate: func(t *testing.T, v interface{}) {
				b, ok := v.(bool)
				if !ok {
					t.Fatalf("type = %T, want bool", v)
				}
				if !b {
					t.Error("value = false, want true")
				}
			},
		},
		{
			name: "null to interface",
			json: `null`,
			validate: func(t *testing.T, v interface{}) {
				if v != nil {
					t.Errorf("value = %v, want nil", v)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target interface{}
			err := Unmarshal([]byte(tt.json), &target)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			tt.validate(t, target)
		})
	}
}
