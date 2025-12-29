package json

import (
	"strings"
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
)

// TestRenderArray_LegacyFormat tests rendering ObjectNode with numeric keys as array
// This covers the renderArray function which has 0% coverage
func TestRenderArray_LegacyFormat(t *testing.T) {
	// Create an ObjectNode with numeric string keys "0", "1", "2" etc.
	// This is the legacy array format that renderArray handles
	pos := ast.Position{}

	// Create properties with numeric keys
	props := make(map[string]ast.SchemaNode)
	props["0"] = ast.NewLiteralNode(int64(1), pos)
	props["1"] = ast.NewLiteralNode(int64(2), pos)
	props["2"] = ast.NewLiteralNode(int64(3), pos)

	node := ast.NewObjectNode(props, pos)

	// Render should recognize this as an array and use renderArray
	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := `[1,2,3]`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderArray_Empty tests rendering empty array in legacy format
func TestRenderArray_Empty(t *testing.T) {
	pos := ast.Position{}

	// Empty object with no properties should render as {}
	props := make(map[string]ast.SchemaNode)
	node := ast.NewObjectNode(props, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := `{}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderArray_MixedTypes tests rendering array with mixed types in legacy format
func TestRenderArray_MixedTypes(t *testing.T) {
	pos := ast.Position{}

	props := make(map[string]ast.SchemaNode)
	props["0"] = ast.NewLiteralNode("hello", pos)
	props["1"] = ast.NewLiteralNode(int64(42), pos)
	props["2"] = ast.NewLiteralNode(true, pos)
	props["3"] = ast.NewLiteralNode(nil, pos)

	node := ast.NewObjectNode(props, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should render as an array
	if !strings.HasPrefix(string(result), "[") || !strings.HasSuffix(string(result), "]") {
		t.Errorf("Expected array format, got %s", string(result))
	}

	// Check it contains all elements
	if !strings.Contains(string(result), `"hello"`) {
		t.Errorf("Missing string element in: %s", string(result))
	}
	if !strings.Contains(string(result), `42`) {
		t.Errorf("Missing int element in: %s", string(result))
	}
	if !strings.Contains(string(result), `true`) {
		t.Errorf("Missing bool element in: %s", string(result))
	}
	if !strings.Contains(string(result), `null`) {
		t.Errorf("Missing null element in: %s", string(result))
	}
}

// TestRenderArray_Nested tests rendering nested arrays in legacy format
func TestRenderArray_Nested(t *testing.T) {
	pos := ast.Position{}

	// Create nested array
	innerProps := make(map[string]ast.SchemaNode)
	innerProps["0"] = ast.NewLiteralNode(int64(1), pos)
	innerProps["1"] = ast.NewLiteralNode(int64(2), pos)
	innerNode := ast.NewObjectNode(innerProps, pos)

	outerProps := make(map[string]ast.SchemaNode)
	outerProps["0"] = innerNode
	outerProps["1"] = ast.NewLiteralNode(int64(3), pos)
	outerNode := ast.NewObjectNode(outerProps, pos)

	result, err := Render(outerNode)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should render as nested array
	expected := `[[1,2],3]`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderArray_WithIndent tests rendering array with indentation in legacy format
func TestRenderArray_WithIndent(t *testing.T) {
	pos := ast.Position{}

	props := make(map[string]ast.SchemaNode)
	props["0"] = ast.NewLiteralNode(int64(1), pos)
	props["1"] = ast.NewLiteralNode(int64(2), pos)
	props["2"] = ast.NewLiteralNode(int64(3), pos)

	node := ast.NewObjectNode(props, pos)

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

// TestRenderArray_WithPrefix tests rendering array with prefix in legacy format
func TestRenderArray_WithPrefix(t *testing.T) {
	pos := ast.Position{}

	props := make(map[string]ast.SchemaNode)
	props["0"] = ast.NewLiteralNode(int64(1), pos)
	props["1"] = ast.NewLiteralNode(int64(2), pos)

	node := ast.NewObjectNode(props, pos)

	result, err := RenderIndent(node, ">>", "  ")
	if err != nil {
		t.Fatalf("RenderIndent error: %v", err)
	}

	// Check that lines have the prefix
	lines := strings.Split(string(result), "\n")
	for i, line := range lines {
		if i == 0 {
			// First line has no prefix
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

// TestRenderArray_MissingElement tests rendering array with missing element
func TestRenderArray_MissingElement(t *testing.T) {
	pos := ast.Position{}

	// Create array with gap (missing index 1)
	// Note: Arrays with gaps are not detected as arrays by isArray() since it requires
	// sequential keys. This will render as an object instead.
	props := make(map[string]ast.SchemaNode)
	props["0"] = ast.NewLiteralNode(int64(1), pos)
	// props["1"] is missing
	props["2"] = ast.NewLiteralNode(int64(3), pos)

	node := ast.NewObjectNode(props, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Since isArray requires sequential keys, this renders as an object
	expected := `{"0":1,"2":3}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderArray_NestedObjects tests rendering array containing objects
func TestRenderArray_NestedObjects(t *testing.T) {
	pos := ast.Position{}

	// Create object to put in array
	objProps := make(map[string]ast.SchemaNode)
	objProps["name"] = ast.NewLiteralNode("Alice", pos)
	objProps["age"] = ast.NewLiteralNode(int64(30), pos)
	objNode := ast.NewObjectNode(objProps, pos)

	// Create array containing the object
	arrayProps := make(map[string]ast.SchemaNode)
	arrayProps["0"] = objNode
	arrayNode := ast.NewObjectNode(arrayProps, pos)

	result, err := Render(arrayNode)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should render as array with object inside
	if !strings.HasPrefix(string(result), "[") || !strings.HasSuffix(string(result), "]") {
		t.Errorf("Expected array format, got %s", string(result))
	}

	// Check object is rendered
	if !strings.Contains(string(result), `"name":"Alice"`) {
		t.Errorf("Missing object content in: %s", string(result))
	}
}

// TestRenderNodeWithDepth_NilNode tests rendering nil node
func TestRenderNodeWithDepth_NilNode(t *testing.T) {
	// Test that nil node renders as "null"
	result, err := Render(nil)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := `null`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderLiteral_ControlCharacters tests rendering strings with control characters
func TestRenderLiteral_ControlCharacters(t *testing.T) {
	pos := ast.Position{}

	tests := []struct {
		name     string
		value    string
		contains string
	}{
		{
			name:     "null character",
			value:    "hello\x00world",
			contains: `\u0000`,
		},
		{
			name:     "control character",
			value:    "test\x01",
			contains: `\u0001`,
		},
		{
			name:     "backspace",
			value:    "test\b",
			contains: `\b`,
		},
		{
			name:     "form feed",
			value:    "test\f",
			contains: `\f`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := ast.NewLiteralNode(tt.value, pos)
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

// TestRenderLiteral_FloatFormats tests different float number formats
func TestRenderLiteral_FloatFormats(t *testing.T) {
	pos := ast.Position{}

	tests := []struct {
		name  string
		value float64
		check func(t *testing.T, result string)
	}{
		{
			name:  "whole number as float",
			value: 42.0,
			check: func(t *testing.T, result string) {
				// Should include decimal point for whole numbers stored as float
				if !strings.Contains(result, ".") {
					t.Errorf("Expected decimal point in: %s", result)
				}
			},
		},
		{
			name:  "decimal number",
			value: 3.14,
			check: func(t *testing.T, result string) {
				if result != "3.14" {
					t.Errorf("Expected 3.14, got %s", result)
				}
			},
		},
		{
			name:  "very small number",
			value: 0.000001,
			check: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Expected non-empty result")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := ast.NewLiteralNode(tt.value, pos)
			result, err := Render(node)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			tt.check(t, string(result))
		})
	}
}

// TestEscapeString_NoEscape tests the fast path when no escaping is needed
func TestEscapeString_NoEscape(t *testing.T) {
	tests := []string{
		"simple",
		"hello world",
		"1234567890",
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := escapeString(input)
			if result != input {
				t.Errorf("Expected %q, got %q", input, result)
			}
		})
	}
}

// TestNeedsEscaping tests the needsEscaping helper function
func TestNeedsEscaping(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"simple", false},
		{"hello world", false},
		{`with "quotes"`, true},
		{`with \backslash`, true},
		{"with /slash", true},
		{"with\nlinefeed", true},
		{"with\rreturn", true},
		{"with\ttab", true},
		{"with\bbackspace", true},
		{"with\fformfeed", true},
		{"with\x00null", true},
		{"with\x1fcontrol", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := needsEscaping(tt.input)
			if result != tt.expected {
				t.Errorf("needsEscaping(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRenderObject_EmptyObject tests rendering empty object
func TestRenderObject_EmptyObject(t *testing.T) {
	pos := ast.Position{}
	props := make(map[string]ast.SchemaNode)
	node := ast.NewObjectNode(props, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := "{}"
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderObject_KeySorting tests that object keys are sorted
func TestRenderObject_KeySorting(t *testing.T) {
	pos := ast.Position{}

	props := make(map[string]ast.SchemaNode)
	props["zebra"] = ast.NewLiteralNode(int64(3), pos)
	props["apple"] = ast.NewLiteralNode(int64(1), pos)
	props["middle"] = ast.NewLiteralNode(int64(2), pos)

	node := ast.NewObjectNode(props, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Keys should be sorted alphabetically
	expected := `{"apple":1,"middle":2,"zebra":3}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// TestRenderArrayData_EmptyArray tests rendering empty ArrayDataNode
func TestRenderArrayData_EmptyArray(t *testing.T) {
	pos := ast.Position{}
	elements := []ast.SchemaNode{}
	node := ast.NewArrayDataNode(elements, pos)

	result, err := Render(node)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expected := "[]"
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}
