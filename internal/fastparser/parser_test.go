package fastparser

import (
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple string",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "empty string",
			input: `""`,
			want:  "",
		},
		{
			name:  "string with spaces",
			input: `"hello world"`,
			want:  "hello world",
		},
		{
			name:  "string with escape sequences",
			input: `"hello\nworld"`,
			want:  "hello\nworld",
		},
		{
			name:  "string with quote escape",
			input: `"say \"hello\""`,
			want:  `say "hello"`,
		},
		{
			name:  "string with backslash",
			input: `"path\\to\\file"`,
			want:  `path\to\file`,
		},
		{
			name:  "string with all escapes",
			input: `"tab\ttab\nline\rreturn\bbackspace\fform"`,
			want:  "tab\ttab\nline\rreturn\bbackspace\fform",
		},
		{
			name:  "string with unicode escape",
			input: `"\u0041\u0042\u0043"`,
			want:  "ABC",
		},
		{
			name:    "unclosed string",
			input:   `"hello`,
			wantErr: true,
		},
		{
			name:    "invalid escape",
			input:   `"hello\x"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:  "integer zero",
			input: "0",
			want:  int64(0),
		},
		{
			name:  "positive integer",
			input: "123",
			want:  int64(123),
		},
		{
			name:  "negative integer",
			input: "-456",
			want:  int64(-456),
		},
		{
			name:  "float with decimal",
			input: "123.456",
			want:  123.456,
		},
		{
			name:  "negative float",
			input: "-123.456",
			want:  -123.456,
		},
		{
			name:  "exponent positive",
			input: "1e10",
			want:  1e10,
		},
		{
			name:  "exponent negative",
			input: "1e-5",
			want:  1e-5,
		},
		{
			name:  "exponent uppercase",
			input: "1E10",
			want:  1E10,
		},
		{
			name:  "decimal with exponent",
			input: "1.23e10",
			want:  1.23e10,
		},
		{
			name:    "leading zero",
			input:   "01",
			wantErr: true,
		},
		{
			name:    "just minus",
			input:   "-",
			wantErr: true,
		},
		{
			name:    "decimal without digits",
			input:   "1.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestParseLiterals(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:  "true",
			input: "true",
			want:  true,
		},
		{
			name:  "false",
			input: "false",
			want:  false,
		},
		{
			name:  "null",
			input: "null",
			want:  nil,
		},
		{
			name:    "invalid true",
			input:   "tru",
			wantErr: true,
		},
		{
			name:    "invalid false",
			input:   "fals",
			wantErr: true,
		},
		{
			name:    "invalid null",
			input:   "nul",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{
			name:    "empty array",
			input:   "[]",
			wantLen: 0,
		},
		{
			name:    "simple array",
			input:   "[1, 2, 3]",
			wantLen: 3,
		},
		{
			name:    "mixed types",
			input:   `[1, "hello", true, null]`,
			wantLen: 4,
		},
		{
			name:    "nested array",
			input:   "[[1, 2], [3, 4]]",
			wantLen: 2,
		},
		{
			name:    "unclosed array",
			input:   "[1, 2",
			wantErr: true,
		},
		{
			name:    "missing comma",
			input:   "[1 2]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				arr, ok := got.([]interface{})
				if !ok {
					t.Errorf("Parse() returned %T, want []interface{}", got)
					return
				}
				if len(arr) != tt.wantLen {
					t.Errorf("Parse() array length = %d, want %d", len(arr), tt.wantLen)
				}
			}
		})
	}
}

func TestParseObject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{
			name:    "empty object",
			input:   "{}",
			wantLen: 0,
		},
		{
			name:    "simple object",
			input:   `{"name": "Alice", "age": 30}`,
			wantLen: 2,
		},
		{
			name:    "nested object",
			input:   `{"user": {"name": "Alice"}}`,
			wantLen: 1,
		},
		{
			name:    "object with array",
			input:   `{"tags": [1, 2, 3]}`,
			wantLen: 1,
		},
		{
			name:    "unclosed object",
			input:   `{"name": "Alice"`,
			wantErr: true,
		},
		{
			name:    "missing colon",
			input:   `{"name" "Alice"}`,
			wantErr: true,
		},
		{
			name:    "missing comma",
			input:   `{"name": "Alice" "age": 30}`,
			wantErr: true,
		},
		{
			name:    "non-string key",
			input:   `{123: "value"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				obj, ok := got.(map[string]interface{})
				if !ok {
					t.Errorf("Parse() returned %T, want map[string]interface{}", got)
					return
				}
				if len(obj) != tt.wantLen {
					t.Errorf("Parse() object length = %d, want %d", len(obj), tt.wantLen)
				}
			}
		})
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		{
			name:  "leading whitespace",
			input: "   true",
			want:  true,
		},
		{
			name:  "trailing whitespace",
			input: "true   ",
			want:  true,
		},
		{
			name:  "whitespace in object",
			input: `{  "name"  :  "Alice"  }`,
			want:  map[string]interface{}{"name": "Alice"},
		},
		{
			name:  "whitespace in array",
			input: `[  1  ,  2  ,  3  ]`,
			want:  []interface{}{int64(1), int64(2), int64(3)},
		},
		{
			name:  "tabs and newlines",
			input: "{\n\t\"name\": \"Alice\"\n}",
			want:  map[string]interface{}{"name": "Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			got, err := p.Parse()
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			// For complex types, just check type
			switch tt.want.(type) {
			case map[string]interface{}:
				if _, ok := got.(map[string]interface{}); !ok {
					t.Errorf("Parse() returned %T, want map[string]interface{}", got)
				}
			case []interface{}:
				if _, ok := got.([]interface{}); !ok {
					t.Errorf("Parse() returned %T, want []interface{}", got)
				}
			default:
				if got != tt.want {
					t.Errorf("Parse() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestParseComplexJSON(t *testing.T) {
	input := `{
		"name": "Alice",
		"age": 30,
		"active": true,
		"balance": 123.45,
		"tags": ["go", "json", "parser"],
		"address": {
			"city": "NYC",
			"zip": "10001"
		},
		"metadata": null
	}`

	p := NewParser([]byte(input))
	got, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	obj, ok := got.(map[string]interface{})
	if !ok {
		t.Fatalf("Parse() returned %T, want map[string]interface{}", got)
	}

	// Check fields
	if obj["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", obj["name"])
	}
	if obj["age"] != int64(30) {
		t.Errorf("age = %v, want 30", obj["age"])
	}
	if obj["active"] != true {
		t.Errorf("active = %v, want true", obj["active"])
	}
	if obj["balance"] != 123.45 {
		t.Errorf("balance = %v, want 123.45", obj["balance"])
	}
	if obj["metadata"] != nil {
		t.Errorf("metadata = %v, want nil", obj["metadata"])
	}

	// Check array
	tags, ok := obj["tags"].([]interface{})
	if !ok {
		t.Fatalf("tags is %T, want []interface{}", obj["tags"])
	}
	if len(tags) != 3 {
		t.Errorf("tags length = %d, want 3", len(tags))
	}

	// Check nested object
	address, ok := obj["address"].(map[string]interface{})
	if !ok {
		t.Fatalf("address is %T, want map[string]interface{}", obj["address"])
	}
	if address["city"] != "NYC" {
		t.Errorf("city = %v, want NYC", address["city"])
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "invalid character",
			input: "@",
		},
		{
			name:  "trailing data",
			input: "true false",
		},
		{
			name:  "incomplete object",
			input: `{"name"`,
		},
		{
			name:  "incomplete array",
			input: `[1, 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser([]byte(tt.input))
			_, err := p.Parse()
			if err == nil {
				t.Errorf("Parse() expected error for input %q", tt.input)
			}
		})
	}
}
