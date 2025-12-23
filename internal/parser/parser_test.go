package parser

import (
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
)

// TestParse_EmptyObject tests parsing an empty object
func TestParse_EmptyObject(t *testing.T) {
	input := `{}`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 0 {
		t.Errorf("expected empty object, got %d properties", len(obj.Properties()))
	}
}

// TestParse_EmptyArray tests parsing an empty array
func TestParse_EmptyArray(t *testing.T) {
	input := `[]`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, ok := node.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected *ast.ArrayDataNode, got %T", node)
	}

	if arr.Len() != 0 {
		t.Errorf("expected empty array, got %d elements", arr.Len())
	}
}

// TestParse_String tests parsing string literals
func TestParse_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
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
			name:  "escaped quote",
			input: `"say \"hello\""`,
			want:  `say "hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			node, err := p.Parse()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lit, ok := node.(*ast.LiteralNode)
			if !ok {
				t.Fatalf("expected *ast.LiteralNode, got %T", node)
			}

			str, ok := lit.Value().(string)
			if !ok {
				t.Fatalf("expected string value, got %T", lit.Value())
			}

			if str != tt.want {
				t.Errorf("expected %q, got %q", tt.want, str)
			}
		})
	}
}

// TestParse_Number tests parsing number literals
func TestParse_Number(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  string // "int" or "float"
		wantInt   int64
		wantFloat float64
	}{
		{
			name:     "zero",
			input:    "0",
			wantType: "int",
			wantInt:  0,
		},
		{
			name:     "positive integer",
			input:    "123",
			wantType: "int",
			wantInt:  123,
		},
		{
			name:     "negative integer",
			input:    "-456",
			wantType: "int",
			wantInt:  -456,
		},
		{
			name:      "decimal",
			input:     "123.456",
			wantType:  "float",
			wantFloat: 123.456,
		},
		{
			name:      "negative decimal",
			input:     "-123.456",
			wantType:  "float",
			wantFloat: -123.456,
		},
		{
			name:      "exponent",
			input:     "1e10",
			wantType:  "float",
			wantFloat: 1e10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			node, err := p.Parse()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lit, ok := node.(*ast.LiteralNode)
			if !ok {
				t.Fatalf("expected *ast.LiteralNode, got %T", node)
			}

			if tt.wantType == "int" {
				i, ok := lit.Value().(int64)
				if !ok {
					t.Fatalf("expected int64 value, got %T", lit.Value())
				}
				if i != tt.wantInt {
					t.Errorf("expected %d, got %d", tt.wantInt, i)
				}
			} else {
				f, ok := lit.Value().(float64)
				if !ok {
					t.Fatalf("expected float64 value, got %T", lit.Value())
				}
				if f != tt.wantFloat {
					t.Errorf("expected %g, got %g", tt.wantFloat, f)
				}
			}
		})
	}
}

// TestParse_Boolean tests parsing boolean literals
func TestParse_Boolean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := NewParser(tt.input)
			node, err := p.Parse()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lit, ok := node.(*ast.LiteralNode)
			if !ok {
				t.Fatalf("expected *ast.LiteralNode, got %T", node)
			}

			b, ok := lit.Value().(bool)
			if !ok {
				t.Fatalf("expected bool value, got %T", lit.Value())
			}

			if b != tt.want {
				t.Errorf("expected %v, got %v", tt.want, b)
			}
		})
	}
}

// TestParse_Null tests parsing null literal
func TestParse_Null(t *testing.T) {
	input := "null"
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lit, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("expected *ast.LiteralNode, got %T", node)
	}

	if lit.Value() != nil {
		t.Errorf("expected nil value, got %v", lit.Value())
	}
}

// TestParse_ObjectSingleProperty tests parsing object with one property
func TestParse_ObjectSingleProperty(t *testing.T) {
	input := `{"name": "Alice"}`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 1 {
		t.Fatalf("expected 1 property, got %d", len(obj.Properties()))
	}

	nameNode, ok := obj.GetProperty("name")
	if !ok {
		t.Fatal("property 'name' not found")
	}

	nameLit, ok := nameNode.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("expected *ast.LiteralNode for 'name', got %T", nameNode)
	}

	if nameLit.Value() != "Alice" {
		t.Errorf("expected 'Alice', got %v", nameLit.Value())
	}
}

// TestParse_ObjectMultipleProperties tests parsing object with multiple properties
func TestParse_ObjectMultipleProperties(t *testing.T) {
	input := `{
		"name": "Alice",
		"age": 30,
		"active": true
	}`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 3 {
		t.Fatalf("expected 3 properties, got %d", len(obj.Properties()))
	}

	// Check name
	nameNode, _ := obj.GetProperty("name")
	if nameNode.(*ast.LiteralNode).Value() != "Alice" {
		t.Errorf("expected name='Alice', got %v", nameNode.(*ast.LiteralNode).Value())
	}

	// Check age
	ageNode, _ := obj.GetProperty("age")
	if ageNode.(*ast.LiteralNode).Value() != int64(30) {
		t.Errorf("expected age=30, got %v", ageNode.(*ast.LiteralNode).Value())
	}

	// Check active
	activeNode, _ := obj.GetProperty("active")
	if activeNode.(*ast.LiteralNode).Value() != true {
		t.Errorf("expected active=true, got %v", activeNode.(*ast.LiteralNode).Value())
	}
}

// TestParse_Array tests parsing arrays
func TestParse_Array(t *testing.T) {
	input := `[1, 2, 3]`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, ok := node.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected *ast.ArrayDataNode, got %T", node)
	}

	if arr.Len() != 3 {
		t.Fatalf("expected 3 elements, got %d", arr.Len())
	}

	// Check elements by index
	for i := 0; i < 3; i++ {
		elem := arr.Get(i)
		if elem == nil {
			t.Fatalf("element %d not found", i)
		}

		lit, ok := elem.(*ast.LiteralNode)
		if !ok {
			t.Fatalf("expected *ast.LiteralNode for element %d, got %T", i, elem)
		}

		val := lit.Value().(int64)
		if val != int64(i+1) {
			t.Errorf("element %d: expected %d, got %d", i, i+1, val)
		}
	}
}

// TestParse_NestedStructures tests parsing nested objects and arrays
func TestParse_NestedStructures(t *testing.T) {
	input := `{
		"user": {
			"name": "Alice",
			"age": 30
		},
		"tags": ["go", "json"]
	}`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	// Check user object
	userNode, ok := obj.GetProperty("user")
	if !ok {
		t.Fatal("property 'user' not found")
	}

	userObj, ok := userNode.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode for 'user', got %T", userNode)
	}

	userName, _ := userObj.GetProperty("name")
	if userName.(*ast.LiteralNode).Value() != "Alice" {
		t.Errorf("expected user.name='Alice', got %v", userName.(*ast.LiteralNode).Value())
	}

	// Check tags array
	tagsNode, ok := obj.GetProperty("tags")
	if !ok {
		t.Fatal("property 'tags' not found")
	}

	tagsArr, ok := tagsNode.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected *ast.ArrayDataNode for 'tags', got %T", tagsNode)
	}

	if tagsArr.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", tagsArr.Len())
	}
}

// TestParse_Errors tests error cases
func TestParse_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed object",
			input: `{"name": "Alice"`,
		},
		{
			name:  "unclosed array",
			input: `[1, 2, 3`,
		},
		{
			name:  "missing colon in object",
			input: `{"name" "Alice"}`,
		},
		{
			name:  "trailing comma in object",
			input: `{"name": "Alice",}`,
		},
		{
			name:  "trailing comma in array",
			input: `[1, 2, 3,]`,
		},
		{
			name:  "empty input",
			input: ``,
		},
		{
			name:  "unexpected content after value",
			input: `{} {}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			_, err := p.Parse()

			if err == nil {
				t.Errorf("expected error for invalid JSON: %s", tt.input)
			}
		})
	}
}

// TestParse_Whitespace tests that whitespace is properly handled
func TestParse_Whitespace(t *testing.T) {
	input := `
		{
			"name": "Alice",
			"age": 30
		}
	`
	p := NewParser(input)
	node, err := p.Parse()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 2 {
		t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
	}
}
