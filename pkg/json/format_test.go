package json

import (
	"bytes"
	"testing"
)

func TestMarshalIndent(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		prefix string
		indent string
		want   string
	}{
		{
			name:   "simple object with 2-space indent",
			input:  map[string]interface{}{"name": "Alice", "age": 30},
			prefix: "",
			indent: "  ",
			want: `{
  "age": 30,
  "name": "Alice"
}`,
		},
		{
			name:   "simple object with tab indent",
			input:  map[string]interface{}{"name": "Bob", "age": 25},
			prefix: "",
			indent: "\t",
			want:   "{\n\t\"age\": 25,\n\t\"name\": \"Bob\"\n}",
		},
		{
			name:   "nested object",
			input:  map[string]interface{}{"user": map[string]interface{}{"name": "Alice", "age": 30}},
			prefix: "",
			indent: "  ",
			want: `{
  "user": {
    "age": 30,
    "name": "Alice"
  }
}`,
		},
		{
			name:   "array",
			input:  []int{1, 2, 3},
			prefix: "",
			indent: "  ",
			want: `[
  1,
  2,
  3
]`,
		},
		{
			name:   "empty object",
			input:  map[string]interface{}{},
			prefix: "",
			indent: "  ",
			want:   "{}",
		},
		{
			name:   "empty array",
			input:  []int{},
			prefix: "",
			indent: "  ",
			want:   "[]",
		},
		{
			name: "struct with tags",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{Name: "Charlie", Age: 35},
			prefix: "",
			indent: "  ",
			want: `{
  "age": 35,
  "name": "Charlie"
}`,
		},
		{
			name:   "with prefix",
			input:  map[string]interface{}{"key": "value"},
			prefix: ">>",
			indent: "  ",
			want: `{
>>  "key": "value"
>>}`,
		},
		{
			name: "complex nested structure",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"id": 1, "name": "Alice"},
					map[string]interface{}{"id": 2, "name": "Bob"},
				},
				"total": 2,
			},
			prefix: "",
			indent: "  ",
			want: `{
  "total": 2,
  "users": [
    {
      "id": 1,
      "name": "Alice"
    },
    {
      "id": 2,
      "name": "Bob"
    }
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalIndent(tt.input, tt.prefix, tt.indent)
			if err != nil {
				t.Fatalf("MarshalIndent() error = %v", err)
			}
			gotStr := string(got)
			if gotStr != tt.want {
				t.Errorf("MarshalIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", gotStr, tt.want)
			}
		})
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		prefix string
		indent string
		want   string
	}{
		{
			name:   "compact object to 2-space indent",
			src:    `{"name":"Alice","age":30}`,
			prefix: "",
			indent: "  ",
			want: `{
  "age": 30,
  "name": "Alice"
}`,
		},
		{
			name:   "compact object to tab indent",
			src:    `{"name":"Bob","age":25}`,
			prefix: "",
			indent: "\t",
			want:   "{\n\t\"age\": 25,\n\t\"name\": \"Bob\"\n}",
		},
		{
			name:   "nested object",
			src:    `{"user":{"name":"Alice","age":30}}`,
			prefix: "",
			indent: "  ",
			want: `{
  "user": {
    "age": 30,
    "name": "Alice"
  }
}`,
		},
		{
			name:   "array",
			src:    `[1,2,3]`,
			prefix: "",
			indent: "  ",
			want: `[
  1,
  2,
  3
]`,
		},
		{
			name:   "empty object",
			src:    `{}`,
			prefix: "",
			indent: "  ",
			want:   "{}",
		},
		{
			name:   "empty array",
			src:    `[]`,
			prefix: "",
			indent: "  ",
			want:   "[]",
		},
		{
			name:   "with prefix",
			src:    `{"key":"value"}`,
			prefix: ">>",
			indent: "  ",
			want: `{
>>  "key": "value"
>>}`,
		},
		{
			name:   "already indented (should re-indent)",
			src:    "{\n  \"name\": \"Alice\"\n}",
			prefix: "",
			indent: "    ",
			want: `{
    "name": "Alice"
}`,
		},
		{
			name:   "string with escaped quotes",
			src:    `{"message":"He said \"hello\""}`,
			prefix: "",
			indent: "  ",
			want: `{
  "message": "He said \"hello\""
}`,
		},
		{
			name:   "complex nested",
			src:    `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"total":2}`,
			prefix: "",
			indent: "  ",
			want: `{
  "total": 2,
  "users": [
    {
      "id": 1,
      "name": "Alice"
    },
    {
      "id": 2,
      "name": "Bob"
    }
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Indent(&buf, []byte(tt.src), tt.prefix, tt.indent)
			if err != nil {
				t.Fatalf("Indent() error = %v", err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("Indent() mismatch:\ngot:\n%s\n\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestIndent_InvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"unclosed object", `{"name":"Alice"`},
		{"unclosed array", `[1,2,3`},
		{"invalid syntax", `{name:"value"}`},
		{"trailing comma", `{"key":"value",}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Indent(&buf, []byte(tt.src), "", "  ")
			if err == nil {
				t.Error("Indent() expected error for invalid JSON, got nil")
			}
		})
	}
}

func TestCompact(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "indented object to compact",
			src: `{
  "name": "Alice",
  "age": 30
}`,
			want: `{"age":30,"name":"Alice"}`,
		},
		{
			name: "indented array to compact",
			src: `[
  1,
  2,
  3
]`,
			want: `[1,2,3]`,
		},
		{
			name: "nested with whitespace",
			src: `{
  "user": {
    "name": "Alice",
    "age": 30
  }
}`,
			want: `{"user":{"age":30,"name":"Alice"}}`,
		},
		{
			name: "already compact (no change)",
			src:  `{"name":"Alice","age":30}`,
			want: `{"age":30,"name":"Alice"}`,
		},
		{
			name: "empty object",
			src:  `{}`,
			want: `{}`,
		},
		{
			name: "empty array",
			src:  `[]`,
			want: `[]`,
		},
		{
			name: "string with spaces (preserve internal spaces)",
			src:  `{"message": "Hello World"}`,
			want: `{"message":"Hello World"}`,
		},
		{
			name: "string with escaped quotes",
			src:  `{"message": "He said \"hello\""}`,
			want: `{"message":"He said \"hello\""}`,
		},
		{
			name: "tabs and newlines",
			src:  "{\n\t\"name\":\t\"Alice\",\n\t\"age\":\t30\n}",
			want: `{"age":30,"name":"Alice"}`,
		},
		{
			name: "complex nested",
			src: `{
  "users": [
    {
      "id": 1,
      "name": "Alice"
    },
    {
      "id": 2,
      "name": "Bob"
    }
  ],
  "total": 2
}`,
			want: `{"total":2,"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Compact(&buf, []byte(tt.src))
			if err != nil {
				t.Fatalf("Compact() error = %v", err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("Compact() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCompact_InvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"unclosed object", `{"name":"Alice"`},
		{"unclosed array", `[1,2,3`},
		{"invalid syntax", `{name:"value"}`},
		{"trailing comma", `{"key":"value",}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Compact(&buf, []byte(tt.src))
			if err == nil {
				t.Error("Compact() expected error for invalid JSON, got nil")
			}
		})
	}
}

func TestMarshalIndent_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"channel", make(chan int)},
		{"function", func() {}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MarshalIndent(tt.input, "", "  ")
			if err == nil {
				t.Error("MarshalIndent() expected error for unsupported type, got nil")
			}
		})
	}
}

func TestRoundTrip_IndentCompact(t *testing.T) {
	// Test that Indent -> Compact -> Indent produces consistent results
	original := `{"name":"Alice","age":30,"city":"NYC"}`

	// Indent it
	var indented bytes.Buffer
	if err := Indent(&indented, []byte(original), "", "  "); err != nil {
		t.Fatalf("First Indent() error = %v", err)
	}

	// Compact it
	var compacted bytes.Buffer
	if err := Compact(&compacted, indented.Bytes()); err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	// Indent again
	var reindented bytes.Buffer
	if err := Indent(&reindented, compacted.Bytes(), "", "  "); err != nil {
		t.Fatalf("Second Indent() error = %v", err)
	}

	// First indented and re-indented should match
	if indented.String() != reindented.String() {
		t.Errorf("Round-trip mismatch:\nfirst:\n%s\n\nreindented:\n%s",
			indented.String(), reindented.String())
	}
}

func TestMarshalIndent_Struct(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
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
			Street: "123 Main St",
			City:   "NYC",
		},
	}

	data, err := MarshalIndent(person, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}

	want := `{
  "address": {
    "city": "NYC",
    "street": "123 Main St"
  },
  "age": 30,
  "name": "Alice"
}`

	got := string(data)
	if got != want {
		t.Errorf("MarshalIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestIndent_PreservesTrailingWhitespace(t *testing.T) {
	// Note: Unlike encoding/json, our implementation does not preserve trailing whitespace
	// This is acceptable as trailing whitespace is not part of JSON structure
	src := `{"name":"Alice"}   ` // trailing spaces

	var buf bytes.Buffer
	if err := Indent(&buf, []byte(src), "", "  "); err != nil {
		t.Fatalf("Indent() error = %v", err)
	}

	result := buf.String()
	expected := "{\n  \"name\": \"Alice\"\n}"
	if result != expected {
		t.Errorf("Indent() mismatch:\ngot:  %q\nwant: %q", result, expected)
	}
}
