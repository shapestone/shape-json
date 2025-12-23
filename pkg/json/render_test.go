package json

import (
	"strings"
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
)

// Test rendering simple objects
func TestRender_SimpleObject(t *testing.T) {
	input := `{"name":"Alice","age":30}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Parse the result to verify it's valid JSON
	_, err = Parse(string(result))
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Result should contain both keys (order may vary due to sorting)
	resultStr := string(result)
	if !strings.Contains(resultStr, `"name":"Alice"`) {
		t.Errorf("Result missing name field: %s", resultStr)
	}
	if !strings.Contains(resultStr, `"age":30`) {
		t.Errorf("Result missing age field: %s", resultStr)
	}
}

// Test rendering arrays
func TestRender_Array(t *testing.T) {
	input := `[1,2,3,"four",true,null]`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := `[1,2,3,"four",true,null]`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// Test rendering nested structures
func TestRender_NestedStructure(t *testing.T) {
	input := `{"user":{"name":"Bob","tags":["admin","user"]}}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Verify it's valid JSON by parsing
	_, err = Parse(string(result))
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Check structure
	resultStr := string(result)
	if !strings.Contains(resultStr, `"user"`) {
		t.Errorf("Missing user key: %s", resultStr)
	}
	if !strings.Contains(resultStr, `"name":"Bob"`) {
		t.Errorf("Missing name field: %s", resultStr)
	}
	if !strings.Contains(resultStr, `"tags"`) {
		t.Errorf("Missing tags field: %s", resultStr)
	}
}

// Test string escaping
func TestRender_StringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "quotation marks",
			input:    `{"text":"He said \"hello\""}`,
			expected: `"text":"He said \"hello\""`,
		},
		{
			name:     "backslash",
			input:    `{"path":"C:\\dir"}`,
			expected: `"path":"C:\\dir"`,
		},
		{
			name:     "newline",
			input:    `{"text":"line1\nline2"}`,
			expected: `"text":"line1\nline2"`,
		},
		{
			name:     "tab",
			input:    `{"text":"col1\tcol2"}`,
			expected: `"text":"col1\tcol2"`,
		},
		{
			name:     "carriage return",
			input:    `{"text":"line1\rline2"}`,
			expected: `"text":"line1\rline2"`,
		},
		{
			name:     "forward slash",
			input:    `{"url":"http://example.com"}`,
			expected: `"url":"http:\/\/example.com"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result, err := Render(node)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			if !strings.Contains(string(result), tt.expected) {
				t.Errorf("Expected to contain %s, got %s", tt.expected, string(result))
			}
		})
	}
}

// Test number formatting
func TestRender_Numbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "integer",
			input:    `{"num":42}`,
			contains: `"num":42`,
		},
		{
			name:     "negative integer",
			input:    `{"num":-123}`,
			contains: `"num":-123`,
		},
		{
			name:     "float",
			input:    `{"num":3.14}`,
			contains: `"num":3.14`,
		},
		{
			name:     "zero",
			input:    `{"num":0}`,
			contains: `"num":0`,
		},
		{
			name:     "large number",
			input:    `{"num":9999999999}`,
			contains: `"num":9999999999`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result, err := Render(node)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			if !strings.Contains(string(result), tt.contains) {
				t.Errorf("Expected to contain %s, got %s", tt.contains, string(result))
			}
		})
	}
}

// Test boolean and null
func TestRender_BooleanAndNull(t *testing.T) {
	input := `{"active":true,"disabled":false,"meta":null}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, `"active":true`) {
		t.Errorf("Missing true value: %s", resultStr)
	}
	if !strings.Contains(resultStr, `"disabled":false`) {
		t.Errorf("Missing false value: %s", resultStr)
	}
	if !strings.Contains(resultStr, `"meta":null`) {
		t.Errorf("Missing null value: %s", resultStr)
	}
}

// Test empty object and array
func TestRender_Empty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		note     string
	}{
		{
			name:     "empty object",
			input:    `{}`,
			expected: `{}`,
		},
		{
			name:     "empty array",
			input:    `[]`,
			expected: `[]`,
		},
		{
			name:     "object with empty array",
			input:    `{"items":[]}`,
			expected: `{"items":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result, err := Render(node)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
				if tt.note != "" {
					t.Logf("Note: %s", tt.note)
				}
			}
		})
	}
}

// Test pretty-printing with RenderIndent
func TestRenderIndent_TwoSpaces(t *testing.T) {
	input := `{"name":"Alice","age":30}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := RenderIndent(node, "", "  ")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	expected := `{
  "age": 30,
  "name": "Alice"
}`

	if string(result) != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, string(result))
	}
}

// Test pretty-printing with tabs
func TestRenderIndent_Tabs(t *testing.T) {
	input := `{"x":1,"y":2}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := RenderIndent(node, "", "\t")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	expected := "{\n\t\"x\": 1,\n\t\"y\": 2\n}"

	if string(result) != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, string(result))
	}
}

// Test pretty-printing arrays
func TestRenderIndent_Array(t *testing.T) {
	input := `[1,2,3]`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := RenderIndent(node, "", "  ")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	expected := `[
  1,
  2,
  3
]`

	if string(result) != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, string(result))
	}
}

// Test pretty-printing nested structures
func TestRenderIndent_Nested(t *testing.T) {
	input := `{"user":{"name":"Bob","age":25},"tags":["admin","user"]}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := RenderIndent(node, "", "  ")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	expected := `{
  "tags": [
    "admin",
    "user"
  ],
  "user": {
    "age": 25,
    "name": "Bob"
  }
}`

	if string(result) != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, string(result))
	}
}

// Test round-trip: Parse → Render → Parse
func TestRender_RoundTrip(t *testing.T) {
	tests := []string{
		`{"name":"Alice","age":30,"active":true}`,
		`[1,2,3,"four",true,null]`,
		`{"user":{"name":"Bob","tags":["admin"]}}`,
		`{"empty":{}}`,
		`{"emptyArray":[]}`,
		`null`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			// First parse
			node1, err := Parse(input)
			if err != nil {
				t.Fatalf("First parse error: %v", err)
			}

			// Render
			rendered, err := Render(node1)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			// Second parse
			node2, err := Parse(string(rendered))
			if err != nil {
				t.Fatalf("Second parse error: %v (rendered: %s)", err, string(rendered))
			}

			// Both nodes should be structurally equivalent
			// (We can't compare directly because keys may be reordered)
			rendered1, _ := Render(node1)
			rendered2, _ := Render(node2)

			if string(rendered1) != string(rendered2) {
				t.Errorf("Round-trip mismatch:\nFirst:  %s\nSecond: %s", string(rendered1), string(rendered2))
			}
		})
	}
}

// Test escapeString function directly
func TestEscapeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "hello", expected: "hello"},
		{input: `say "hi"`, expected: `say \"hi\"`},
		{input: `C:\path`, expected: `C:\\path`},
		{input: "line1\nline2", expected: `line1\nline2`},
		{input: "col1\tcol2", expected: `col1\tcol2`},
		{input: "page\fbreak", expected: `page\fbreak`},
		{input: "back\bspace", expected: `back\bspace`},
		{input: "return\rcarriage", expected: `return\rcarriage`},
		{input: "http://example.com", expected: `http:\/\/example.com`},
		{input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test isArray function
func TestIsArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "object",
			input:    `{"key":"value"}`,
			expected: false,
		},
		{
			name:     "array with mixed keys",
			input:    `{"0":1,"key":"value"}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			obj, ok := node.(*ast.ObjectNode)
			if !ok {
				// Arrays are now ArrayDataNode, not ObjectNode - skip those tests
				t.Skipf("Parser now returns ArrayDataNode for arrays, not ObjectNode")
			}

			result := isArray(obj.Properties())
			if result != tt.expected {
				t.Errorf("Expected isArray=%v, got %v", tt.expected, result)
			}
		})
	}
}

// Test key sorting in objects
func TestRender_KeySorting(t *testing.T) {
	input := `{"zebra":1,"apple":2,"middle":3}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Keys should be sorted alphabetically
	expected := `{"apple":2,"middle":3,"zebra":1}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// Test rendering with prefix in RenderIndent
func TestRenderIndent_WithPrefix(t *testing.T) {
	input := `{"a":1,"b":2}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := RenderIndent(node, ">>", "  ")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	// Each line should start with the prefix
	lines := strings.Split(string(result), "\n")
	for i, line := range lines {
		if i == 0 {
			// First line has no prefix (it's just the opening brace)
			continue
		}
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, ">>") {
			t.Errorf("Line %d missing prefix: %q", i, line)
		}
	}
}

// Test Unicode handling
func TestRender_Unicode(t *testing.T) {
	input := `{"emoji":"😀","chinese":"你好","arabic":"مرحبا"}`
	node, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should preserve Unicode characters
	if !strings.Contains(string(result), "😀") {
		t.Errorf("Missing emoji in result: %s", string(result))
	}
	if !strings.Contains(string(result), "你好") {
		t.Errorf("Missing Chinese in result: %s", string(result))
	}
	if !strings.Contains(string(result), "مرحبا") {
		t.Errorf("Missing Arabic in result: %s", string(result))
	}

	// Verify it's still valid JSON
	_, err = Parse(string(result))
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}
