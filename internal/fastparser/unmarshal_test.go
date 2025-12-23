package fastparser

import (
	"testing"
)

func TestUnmarshalStruct(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Active bool
	}

	tests := []struct {
		name    string
		input   string
		want    Person
		wantErr bool
	}{
		{
			name:  "simple struct",
			input: `{"Name": "Alice", "Age": 30, "Active": true}`,
			want:  Person{Name: "Alice", Age: 30, Active: true},
		},
		{
			name:  "partial fields",
			input: `{"Name": "Bob"}`,
			want:  Person{Name: "Bob"},
		},
		{
			name:  "extra fields ignored",
			input: `{"Name": "Charlie", "Age": 25, "Unknown": "value"}`,
			want:  Person{Name: "Charlie", Age: 25},
		},
		{
			name:    "type mismatch",
			input:   `{"Name": 123}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Person
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalMap(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{
			name:    "string map",
			input:   `{"key1": "value1", "key2": "value2"}`,
			wantLen: 2,
		},
		{
			name:    "int map",
			input:   `{"a": 1, "b": 2, "c": 3}`,
			wantLen: 3,
		},
		{
			name:    "empty map",
			input:   `{}`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]interface{}
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("Unmarshal() map length = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestUnmarshalSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{
			name:    "int slice",
			input:   `[1, 2, 3, 4, 5]`,
			wantLen: 5,
		},
		{
			name:    "string slice",
			input:   `["a", "b", "c"]`,
			wantLen: 3,
		},
		{
			name:    "empty slice",
			input:   `[]`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []interface{}
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("Unmarshal() slice length = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestUnmarshalInterface(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(interface{}) bool
	}{
		{
			name:  "string",
			input: `"hello"`,
			check: func(v interface{}) bool {
				s, ok := v.(string)
				return ok && s == "hello"
			},
		},
		{
			name:  "number",
			input: `42`,
			check: func(v interface{}) bool {
				i, ok := v.(int64)
				return ok && i == 42
			},
		},
		{
			name:  "bool",
			input: `true`,
			check: func(v interface{}) bool {
				b, ok := v.(bool)
				return ok && b == true
			},
		},
		{
			name:  "null",
			input: `null`,
			check: func(v interface{}) bool {
				return v == nil
			},
		},
		{
			name:  "object",
			input: `{"name": "Alice"}`,
			check: func(v interface{}) bool {
				m, ok := v.(map[string]interface{})
				return ok && len(m) == 1
			},
		},
		{
			name:  "array",
			input: `[1, 2, 3]`,
			check: func(v interface{}) bool {
				a, ok := v.([]interface{})
				return ok && len(a) == 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			err := Unmarshal([]byte(tt.input), &got)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			if !tt.check(got) {
				t.Errorf("Unmarshal() = %v (%T), check failed", got, got)
			}
		})
	}
}

func TestUnmarshalNestedStruct(t *testing.T) {
	type Address struct {
		City string
		Zip  string
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
	}

	input := `{
		"Name": "Alice",
		"Age": 30,
		"Address": {
			"City": "NYC",
			"Zip": "10001"
		}
	}`

	var got Person
	err := Unmarshal([]byte(input), &got)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", got.Name)
	}
	if got.Age != 30 {
		t.Errorf("Age = %v, want 30", got.Age)
	}
	if got.Address.City != "NYC" {
		t.Errorf("Address.City = %v, want NYC", got.Address.City)
	}
	if got.Address.Zip != "10001" {
		t.Errorf("Address.Zip = %v, want 10001", got.Address.Zip)
	}
}

func TestUnmarshalStructWithSlice(t *testing.T) {
	type Person struct {
		Name string
		Tags []string
	}

	input := `{"Name": "Alice", "Tags": ["go", "json", "parser"]}`

	var got Person
	err := Unmarshal([]byte(input), &got)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", got.Name)
	}
	if len(got.Tags) != 3 {
		t.Errorf("Tags length = %d, want 3", len(got.Tags))
	}
	if len(got.Tags) > 0 && got.Tags[0] != "go" {
		t.Errorf("Tags[0] = %v, want go", got.Tags[0])
	}
}

func TestUnmarshalErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		dest  interface{}
	}{
		{
			name:  "nil pointer",
			input: `{}`,
			dest:  nil,
		},
		{
			name:  "non-pointer",
			input: `{}`,
			dest:  struct{}{},
		},
		{
			name:  "nil struct pointer",
			input: `{}`,
			dest:  (*struct{})(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.dest)
			if err == nil {
				t.Errorf("Unmarshal() expected error for %s", tt.name)
			}
		})
	}
}

func TestUnmarshalComplexJSON(t *testing.T) {
	type Person struct {
		Name     string
		Age      int
		Active   bool
		Balance  float64
		Tags     []string
		Metadata map[string]interface{}
	}

	input := `{
		"Name": "Alice",
		"Age": 30,
		"Active": true,
		"Balance": 123.45,
		"Tags": ["go", "json"],
		"Metadata": {
			"level": 5,
			"verified": true
		}
	}`

	var got Person
	err := Unmarshal([]byte(input), &got)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", got.Name)
	}
	if got.Age != 30 {
		t.Errorf("Age = %v, want 30", got.Age)
	}
	if !got.Active {
		t.Errorf("Active = %v, want true", got.Active)
	}
	if got.Balance != 123.45 {
		t.Errorf("Balance = %v, want 123.45", got.Balance)
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(got.Tags))
	}
	if len(got.Metadata) != 2 {
		t.Errorf("Metadata length = %d, want 2", len(got.Metadata))
	}
}
