package json

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
)

func TestParse_Valid(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "simple object",
			json: `{"name": "Alice", "age": 30}`,
		},
		{
			name: "simple array",
			json: `[1, 2, 3, 4, 5]`,
		},
		{
			name: "nested object",
			json: `{"user": {"name": "Bob", "scores": [10, 20, 30]}}`,
		},
		{
			name: "empty object",
			json: `{}`,
		},
		{
			name: "empty array",
			json: `[]`,
		},
		{
			name: "null",
			json: `null`,
		},
		{
			name: "boolean",
			json: `true`,
		},
		{
			name: "number",
			json: `42`,
		},
		{
			name: "string",
			json: `"hello"`,
		},
		{
			name: "with whitespace",
			json: `{
  "name": "Alice",
  "age": 30
}`,
		},
		{
			name: "complex nested",
			json: `{
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25}
  ],
  "count": 2,
  "active": true
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.json)
			if err != nil {
				t.Errorf("Parse() failed for valid JSON: %v\nJSON:\n%s", err, tt.json)
			}
			if node == nil {
				t.Errorf("Parse() returned nil node for valid JSON:\n%s", tt.json)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "unclosed object",
			json: `{"name": "Alice"`,
		},
		{
			name: "unclosed array",
			json: `[1, 2, 3`,
		},
		{
			name: "trailing comma in object",
			json: `{"name": "Alice",}`,
		},
		{
			name: "trailing comma in array",
			json: `[1, 2, 3,]`,
		},
		{
			name: "missing quotes on key",
			json: `{name: "Alice"}`,
		},
		{
			name: "single quotes instead of double",
			json: `{'name': 'Alice'}`,
		},
		{
			name: "missing comma",
			json: `{"name": "Alice" "age": 30}`,
		},
		{
			name: "unclosed string",
			json: `{"name": "Alice}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.json)
			if err == nil {
				t.Errorf("expected error for invalid JSON:\n%s", tt.json)
			}
		})
	}
}

func TestParse_EmptyString(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestParse_Unicode(t *testing.T) {
	json := `{"message": "Hello 世界 🌍"}`
	node, err := Parse(json)
	if err != nil {
		t.Errorf("Parse() failed for JSON with unicode: %v", err)
	}
	if node == nil {
		t.Error("Parse() returned nil node for valid JSON with unicode")
	}
}

func TestParse_EscapedCharacters(t *testing.T) {
	json := `{"path": "C:\\Users\\Alice\\file.txt", "quote": "He said \"Hello\""}`
	node, err := Parse(json)
	if err != nil {
		t.Errorf("Parse() failed for JSON with escaped characters: %v", err)
	}
	if node == nil {
		t.Error("Parse() returned nil node for valid JSON with escaped characters")
	}
}

func TestFormat(t *testing.T) {
	format := Format()
	if format != "JSON" {
		t.Errorf("expected Format() to return 'JSON', got %q", format)
	}
}

// TestDetectFormat tests JSON format detection from strings
func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFormat  string
		expectError bool
	}{
		// Valid JSON - Objects
		{
			name:        "valid object",
			input:       `{"name": "Alice", "age": 30}`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "empty object",
			input:       `{}`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "nested object",
			input:       `{"user": {"name": "Alice", "settings": {"theme": "dark"}}}`,
			wantFormat:  "JSON",
			expectError: false,
		},

		// Valid JSON - Arrays
		{
			name:        "valid array",
			input:       `[1, 2, 3, 4, 5]`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "empty array",
			input:       `[]`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "array of objects",
			input:       `[{"id": 1}, {"id": 2}]`,
			wantFormat:  "JSON",
			expectError: false,
		},

		// Valid JSON - Primitives
		{
			name:        "string value",
			input:       `"hello world"`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "integer",
			input:       `42`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "negative integer",
			input:       `-123`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "float",
			input:       `3.14159`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "scientific notation",
			input:       `1.23e10`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "true boolean",
			input:       `true`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "false boolean",
			input:       `false`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "null value",
			input:       `null`,
			wantFormat:  "JSON",
			expectError: false,
		},

		// Valid JSON - With whitespace
		{
			name:        "object with whitespace",
			input:       "  \n  {  \"key\"  :  \"value\"  }  \n  ",
			wantFormat:  "JSON",
			expectError: false,
		},

		// Invalid JSON
		{
			name:        "empty string",
			input:       "",
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   \n\t  ",
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "invalid syntax - missing quotes",
			input:       `{key: value}`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "invalid syntax - trailing comma",
			input:       `{"key": "value",}`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "invalid syntax - unclosed brace",
			input:       `{"key": "value"`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "invalid syntax - unclosed bracket",
			input:       `[1, 2, 3`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "invalid syntax - single quotes",
			input:       `{'key': 'value'}`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "plain text",
			input:       `this is not json`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "xml content",
			input:       `<root><item>value</item></root>`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "partial json",
			input:       `{"incomplete":`,
			wantFormat:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := DetectFormat(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("DetectFormat() expected error but got none")
				}
				if format != "" {
					t.Errorf("DetectFormat() format = %q, want empty string for invalid input", format)
				}
			} else {
				if err != nil {
					t.Errorf("DetectFormat() unexpected error: %v", err)
				}
				if format != tt.wantFormat {
					t.Errorf("DetectFormat() format = %q, want %q", format, tt.wantFormat)
				}
			}
		})
	}
}

// TestDetectFormatFromReader tests JSON format detection from io.Reader
func TestDetectFormatFromReader(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFormat  string
		expectError bool
	}{
		{
			name:        "valid JSON object from reader",
			input:       `{"name": "Alice", "age": 30}`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "valid JSON array from reader",
			input:       `[1, 2, 3, 4, 5]`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "valid JSON primitives from reader",
			input:       `"hello"`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "large valid JSON from reader",
			input:       `{"users": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}, {"id": 3, "name": "Charlie"}], "total": 3}`,
			wantFormat:  "JSON",
			expectError: false,
		},
		{
			name:        "invalid JSON from reader",
			input:       `{invalid json}`,
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "empty reader",
			input:       "",
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "plain text from reader",
			input:       "this is not json",
			wantFormat:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			format, err := DetectFormatFromReader(reader)

			if tt.expectError {
				if err == nil {
					t.Errorf("DetectFormatFromReader() expected error but got none")
				}
				if format != "" {
					t.Errorf("DetectFormatFromReader() format = %q, want empty string for invalid input", format)
				}
			} else {
				if err != nil {
					t.Errorf("DetectFormatFromReader() unexpected error: %v", err)
				}
				if format != tt.wantFormat {
					t.Errorf("DetectFormatFromReader() format = %q, want %q", format, tt.wantFormat)
				}
			}
		})
	}
}

// Validate Tests

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		// Valid JSON - should return nil
		{name: "valid object", input: `{"name": "Alice"}`, expectErr: false},
		{name: "valid array", input: `[1, 2, 3]`, expectErr: false},
		{name: "valid primitive", input: `42`, expectErr: false},

		// Invalid JSON - should return error
		{name: "invalid empty", input: ``, expectErr: true},
		{name: "invalid syntax", input: `{invalid}`, expectErr: true},
		{name: "invalid unclosed", input: `{"name": "Alice"`, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, expectErr %v", err, tt.expectErr)
			}
			if tt.expectErr && err == nil {
				t.Error("Validate() expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

// ValidateReader Tests

func TestValidateReader(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		// Valid JSON - should return nil
		{name: "valid object from reader", input: `{"name": "Alice"}`, expectErr: false},
		{name: "valid array from reader", input: `[1, 2, 3]`, expectErr: false},
		{name: "valid primitive from reader", input: `42`, expectErr: false},

		// Invalid JSON - should return error
		{name: "invalid empty reader", input: ``, expectErr: true},
		{name: "invalid syntax from reader", input: `{invalid}`, expectErr: true},
		{name: "invalid unclosed from reader", input: `{"name": "Alice"`, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			err := ValidateReader(reader)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateReader() error = %v, expectErr %v", err, tt.expectErr)
			}
			if tt.expectErr && err == nil {
				t.Error("ValidateReader() expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("ValidateReader() unexpected error: %v", err)
			}
		})
	}
}

// ParseReader Tests

func TestParseReader_StringsReader(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "simple object",
			json: `{"name": "Alice", "age": 30}`,
		},
		{
			name: "simple array",
			json: `[1, 2, 3, 4, 5]`,
		},
		{
			name: "nested object",
			json: `{"user": {"name": "Bob", "scores": [10, 20, 30]}}`,
		},
		{
			name: "empty object",
			json: `{}`,
		},
		{
			name: "empty array",
			json: `[]`,
		},
		{
			name: "null",
			json: `null`,
		},
		{
			name: "boolean",
			json: `true`,
		},
		{
			name: "number",
			json: `42`,
		},
		{
			name: "string",
			json: `"hello"`,
		},
		{
			name: "complex nested",
			json: `{
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25}
  ],
  "count": 2,
  "active": true
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			node, err := ParseReader(reader)
			if err != nil {
				t.Errorf("ParseReader() failed for valid JSON: %v\nJSON:\n%s", err, tt.json)
			}
			if node == nil {
				t.Errorf("ParseReader() returned nil node for valid JSON:\n%s", tt.json)
			}
		})
	}
}

func TestParseReader_BytesBuffer(t *testing.T) {
	jsonData := `{"message": "Hello from buffer", "count": 42}`
	buffer := bytes.NewBufferString(jsonData)

	node, err := ParseReader(buffer)
	if err != nil {
		t.Errorf("ParseReader() failed with bytes.Buffer: %v", err)
	}
	if node == nil {
		t.Error("ParseReader() returned nil node for valid JSON from buffer")
	}

	// Verify the parsed structure
	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatal("expected *ast.ObjectNode")
	}

	messageNode, ok := obj.GetProperty("message")
	if !ok {
		t.Error("expected 'message' property")
	}
	messageLit, ok := messageNode.(*ast.LiteralNode)
	if !ok {
		t.Fatal("expected message to be *ast.LiteralNode")
	}
	if messageLit.Value().(string) != "Hello from buffer" {
		t.Errorf("expected message 'Hello from buffer', got %q", messageLit.Value())
	}
}

func TestParseReader_File(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := filepath.Join("testdata", "parse_reader")
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases with different file sizes
	tests := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "small file",
			filename: "small.json",
			content:  `{"name": "Alice", "age": 30}`,
		},
		{
			name:     "medium file with arrays",
			filename: "medium.json",
			content: `{
  "users": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"},
    {"id": 3, "name": "Charlie", "email": "charlie@example.com"}
  ],
  "metadata": {
    "count": 3,
    "version": "1.0"
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write test file
			filePath := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Parse from file
			file, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer file.Close()

			node, err := ParseReader(file)
			if err != nil {
				t.Errorf("ParseReader() failed for file: %v", err)
			}
			if node == nil {
				t.Error("ParseReader() returned nil node for valid JSON file")
			}
		})
	}
}

func TestParseReader_LargeJSON(t *testing.T) {
	// Test parsing large JSON that exceeds buffer size to verify buffered stream works correctly

	// Create a large JSON object that exceeds the buffer size (64KB)
	// Each entry is about 100 bytes, so 1000 entries = ~100KB
	var sb strings.Builder
	sb.WriteString(`{"items": [`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		// Build valid JSON - use fmt to format proper numbers
		sb.WriteString(`{"id": 1000000000, "name": "Item_`)
		sb.WriteString(strings.Repeat("X", 50)) // Padding
		sb.WriteString(`", "description": "`)
		sb.WriteString(strings.Repeat("Y", 30)) // More padding
		sb.WriteString(`"}`)
	}
	sb.WriteString(`]}`)

	largeJSON := sb.String()
	t.Logf("Generated JSON size: %d bytes", len(largeJSON))

	// First verify with Parse() that the JSON is valid
	_, err := Parse(largeJSON)
	if err != nil {
		t.Fatalf("Parse() failed - JSON is invalid: %v", err)
	}

	// Now test with ParseReader()
	reader := strings.NewReader(largeJSON)
	node, err := ParseReader(reader)
	if err != nil {
		t.Errorf("ParseReader() failed for large JSON: %v", err)
	}
	if node == nil {
		t.Error("ParseReader() returned nil node for large JSON")
		return
	}

	// Verify it's an object with "items" property
	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatal("expected *ast.ObjectNode")
	}

	itemsNode, ok := obj.GetProperty("items")
	if !ok {
		t.Error("expected 'items' property in parsed object")
	}
	if itemsNode == nil {
		t.Error("'items' property is nil")
	}
}

func TestParseReader_IdenticalToParse(t *testing.T) {
	// Verify that ParseReader produces the same AST as Parse for the same input
	tests := []string{
		`{"name": "Alice", "age": 30}`,
		`[1, 2, 3, 4, 5]`,
		`{"nested": {"object": {"with": "values"}}}`,
		`true`,
		`null`,
		`42`,
		`"string value"`,
		`{"array": [1, 2, 3], "object": {"a": 1}}`,
	}

	for _, jsonStr := range tests {
		t.Run(jsonStr, func(t *testing.T) {
			// Parse with Parse()
			node1, err1 := Parse(jsonStr)
			if err1 != nil {
				t.Fatalf("Parse() failed: %v", err1)
			}

			// Parse with ParseReader()
			reader := strings.NewReader(jsonStr)
			node2, err2 := ParseReader(reader)
			if err2 != nil {
				t.Fatalf("ParseReader() failed: %v", err2)
			}

			// Compare results (basic check - both should be non-nil and same type)
			if node1 == nil || node2 == nil {
				t.Error("one or both nodes are nil")
				return
			}

			// Type assertion to check they're the same type
			switch n1 := node1.(type) {
			case *ast.ObjectNode:
				n2, ok := node2.(*ast.ObjectNode)
				if !ok {
					t.Error("ParseReader returned different type than Parse")
				}
				// Compare property count
				if len(n1.Properties()) != len(n2.Properties()) {
					t.Errorf("different property counts: Parse=%d, ParseReader=%d",
						len(n1.Properties()), len(n2.Properties()))
				}
			case *ast.LiteralNode:
				n2, ok := node2.(*ast.LiteralNode)
				if !ok {
					t.Error("ParseReader returned different type than Parse")
				}
				// Compare values
				if n1.Value() != n2.Value() {
					t.Errorf("different values: Parse=%v, ParseReader=%v",
						n1.Value(), n2.Value())
				}
			}
		})
	}
}

func TestParseReader_Invalid(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "unclosed object",
			json: `{"name": "Alice"`,
		},
		{
			name: "unclosed array",
			json: `[1, 2, 3`,
		},
		{
			name: "trailing comma in object",
			json: `{"name": "Alice",}`,
		},
		{
			name: "trailing comma in array",
			json: `[1, 2, 3,]`,
		},
		{
			name: "missing quotes on key",
			json: `{name: "Alice"}`,
		},
		{
			name: "unclosed string",
			json: `{"name": "Alice}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			_, err := ParseReader(reader)
			if err == nil {
				t.Errorf("expected error for invalid JSON:\n%s", tt.json)
			}
		})
	}
}

func TestParseReader_EmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	_, err := ParseReader(reader)
	if err == nil {
		t.Error("expected error for empty reader")
	}
}

func TestParseReader_Unicode(t *testing.T) {
	json := `{"message": "Hello 世界 🌍", "emoji": "😀"}`
	reader := strings.NewReader(json)
	node, err := ParseReader(reader)
	if err != nil {
		t.Errorf("ParseReader() failed for JSON with unicode: %v", err)
	}
	if node == nil {
		t.Error("ParseReader() returned nil node for valid JSON with unicode")
	}

	// Verify unicode was preserved
	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatal("expected *ast.ObjectNode")
	}

	messageNode, ok := obj.GetProperty("message")
	if !ok {
		t.Error("expected 'message' property")
	}
	messageLit, ok := messageNode.(*ast.LiteralNode)
	if !ok {
		t.Fatal("expected message to be *ast.LiteralNode")
	}
	messageStr := messageLit.Value().(string)
	if !strings.Contains(messageStr, "世界") || !strings.Contains(messageStr, "🌍") {
		t.Errorf("unicode characters not preserved: got %q", messageStr)
	}
}

func TestParseReader_PositionTracking(t *testing.T) {
	// Test that position information is tracked correctly
	json := `{
  "name": "Alice",
  "age": 30
}`
	reader := strings.NewReader(json)
	node, err := ParseReader(reader)
	if err != nil {
		t.Fatalf("ParseReader() failed: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatal("expected *ast.ObjectNode")
	}

	// Check that the object has position information
	pos := obj.Position()
	if pos.Offset < 0 || pos.Line < 1 || pos.Column < 1 {
		t.Errorf("invalid position tracking: offset=%d, line=%d, column=%d",
			pos.Offset, pos.Line, pos.Column)
	}
}
