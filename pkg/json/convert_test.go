package json

import (
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
)

// TestNodeToInterface_LiteralNode tests conversion of literal nodes
func TestNodeToInterface_LiteralNode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		check    func(t *testing.T, result interface{})
	}{
		{
			name:  "string literal",
			input: `"hello"`,
			check: func(t *testing.T, result interface{}) {
				if s, ok := result.(string); !ok || s != "hello" {
					t.Errorf("expected string 'hello', got %T: %v", result, result)
				}
			},
		},
		{
			name:  "integer literal",
			input: `42`,
			check: func(t *testing.T, result interface{}) {
				if i, ok := result.(int64); !ok || i != 42 {
					t.Errorf("expected int64 42, got %T: %v", result, result)
				}
			},
		},
		{
			name:  "float literal that's a whole number",
			input: `42.0`,
			check: func(t *testing.T, result interface{}) {
				// Should be converted to int64
				if i, ok := result.(int64); !ok || i != 42 {
					t.Errorf("expected int64 42, got %T: %v", result, result)
				}
			},
		},
		{
			name:  "float literal with decimal",
			input: `3.14`,
			check: func(t *testing.T, result interface{}) {
				if f, ok := result.(float64); !ok || f != 3.14 {
					t.Errorf("expected float64 3.14, got %T: %v", result, result)
				}
			},
		},
		{
			name:  "bool true",
			input: `true`,
			check: func(t *testing.T, result interface{}) {
				if b, ok := result.(bool); !ok || !b {
					t.Errorf("expected bool true, got %T: %v", result, result)
				}
			},
		},
		{
			name:  "bool false",
			input: `false`,
			check: func(t *testing.T, result interface{}) {
				if b, ok := result.(bool); !ok || b {
					t.Errorf("expected bool false, got %T: %v", result, result)
				}
			},
		},
		{
			name:  "null literal",
			input: `null`,
			check: func(t *testing.T, result interface{}) {
				if result != nil {
					t.Errorf("expected nil, got %T: %v", result, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := NodeToInterface(node)
			tt.check(t, result)
		})
	}
}

// TestNodeToInterface_ArrayDataNode tests conversion of array data nodes
func TestNodeToInterface_ArrayDataNode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result interface{})
	}{
		{
			name:  "simple array",
			input: `[1, 2, 3]`,
			check: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{}, got %T", result)
				}
				if len(arr) != 3 {
					t.Errorf("expected length 3, got %d", len(arr))
				}
			},
		},
		{
			name:  "array with mixed types",
			input: `[1, "hello", true, null]`,
			check: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{}, got %T", result)
				}
				if len(arr) != 4 {
					t.Errorf("expected length 4, got %d", len(arr))
				}
				if arr[1] != "hello" {
					t.Errorf("expected arr[1] to be 'hello', got %v", arr[1])
				}
			},
		},
		{
			name:  "nested arrays",
			input: `[[1, 2], [3, 4]]`,
			check: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{}, got %T", result)
				}
				if len(arr) != 2 {
					t.Errorf("expected length 2, got %d", len(arr))
				}
				nested, ok := arr[0].([]interface{})
				if !ok {
					t.Errorf("expected nested array, got %T", arr[0])
				}
				if len(nested) != 2 {
					t.Errorf("expected nested length 2, got %d", len(nested))
				}
			},
		},
		{
			name:  "empty array",
			input: `[]`,
			check: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{}, got %T", result)
				}
				if len(arr) != 0 {
					t.Errorf("expected empty array, got length %d", len(arr))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := NodeToInterface(node)
			tt.check(t, result)
		})
	}
}

// TestNodeToInterface_ObjectNode tests conversion of object nodes
func TestNodeToInterface_ObjectNode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result interface{})
	}{
		{
			name:  "simple object",
			input: `{"name": "Alice", "age": 30}`,
			check: func(t *testing.T, result interface{}) {
				obj, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("expected map[string]interface{}, got %T", result)
				}
				if obj["name"] != "Alice" {
					t.Errorf("expected name='Alice', got %v", obj["name"])
				}
				if obj["age"] != int64(30) {
					t.Errorf("expected age=30, got %v", obj["age"])
				}
			},
		},
		{
			name:  "nested object",
			input: `{"user": {"name": "Bob", "age": 25}}`,
			check: func(t *testing.T, result interface{}) {
				obj, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("expected map[string]interface{}, got %T", result)
				}
				user, ok := obj["user"].(map[string]interface{})
				if !ok {
					t.Fatalf("expected nested map, got %T", obj["user"])
				}
				if user["name"] != "Bob" {
					t.Errorf("expected name='Bob', got %v", user["name"])
				}
			},
		},
		{
			name:  "empty object",
			input: `{}`,
			check: func(t *testing.T, result interface{}) {
				obj, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("expected map[string]interface{}, got %T", result)
				}
				if len(obj) != 0 {
					t.Errorf("expected empty object, got %d entries", len(obj))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := NodeToInterface(node)
			tt.check(t, result)
		})
	}
}

// TestReleaseTree tests releasing AST nodes back to pools
func TestReleaseTree(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"literal node", `"hello"`},
		{"number node", `42`},
		{"bool node", `true`},
		{"null node", `null`},
		{"array node", `[1, 2, 3]`},
		{"object node", `{"name": "Alice"}`},
		{"nested structure", `{"users": [{"name": "Alice"}, {"name": "Bob"}]}`},
		{"empty array", `[]`},
		{"empty object", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			// Should not panic
			ReleaseTree(node)
		})
	}
}

// TestReleaseTree_Nil tests releasing nil node
func TestReleaseTree_Nil(t *testing.T) {
	// Should not panic
	ReleaseTree(nil)
}

// TestInterfaceToNode_AllTypes tests all supported type conversions
func TestInterfaceToNode_AllTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		check func(t *testing.T, node ast.SchemaNode)
	}{
		{
			name:  "nil",
			input: nil,
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != nil {
					t.Errorf("expected nil value, got %v", lit.Value())
				}
			},
		},
		{
			name:  "string",
			input: "hello",
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != "hello" {
					t.Errorf("expected 'hello', got %v", lit.Value())
				}
			},
		},
		{
			name:  "bool",
			input: true,
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != true {
					t.Errorf("expected true, got %v", lit.Value())
				}
			},
		},
		{
			name:  "int",
			input: 42,
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "int64",
			input: int64(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "int32",
			input: int32(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "int16",
			input: int16(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "int8",
			input: int8(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "uint",
			input: uint(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "uint64",
			input: uint64(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "uint32",
			input: uint32(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "uint16",
			input: uint16(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "uint8",
			input: uint8(42),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != int64(42) {
					t.Errorf("expected int64(42), got %v", lit.Value())
				}
			},
		},
		{
			name:  "float64",
			input: 3.14,
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				if lit.Value() != 3.14 {
					t.Errorf("expected 3.14, got %v", lit.Value())
				}
			},
		},
		{
			name:  "float32",
			input: float32(3.14),
			check: func(t *testing.T, node ast.SchemaNode) {
				lit, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("expected *ast.LiteralNode, got %T", node)
				}
				// float32 is converted to float64
				val, ok := lit.Value().(float64)
				if !ok {
					t.Errorf("expected float64, got %T", lit.Value())
				}
				if val < 3.13 || val > 3.15 {
					t.Errorf("expected ~3.14, got %v", val)
				}
			},
		},
		{
			name:  "slice",
			input: []interface{}{1, "hello", true},
			check: func(t *testing.T, node ast.SchemaNode) {
				arr, ok := node.(*ast.ArrayDataNode)
				if !ok {
					t.Fatalf("expected *ast.ArrayDataNode, got %T", node)
				}
				if len(arr.Elements()) != 3 {
					t.Errorf("expected 3 elements, got %d", len(arr.Elements()))
				}
			},
		},
		{
			name:  "map",
			input: map[string]interface{}{"name": "Alice", "age": 30},
			check: func(t *testing.T, node ast.SchemaNode) {
				obj, ok := node.(*ast.ObjectNode)
				if !ok {
					t.Fatalf("expected *ast.ObjectNode, got %T", node)
				}
				if len(obj.Properties()) != 2 {
					t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := InterfaceToNode(tt.input)
			if err != nil {
				t.Fatalf("InterfaceToNode() error = %v", err)
			}
			tt.check(t, node)
		})
	}
}

// TestInterfaceToNode_Document tests Document type conversion
func TestInterfaceToNode_Document(t *testing.T) {
	doc := NewDocument().SetString("name", "Alice").SetInt("age", 30)

	node, err := InterfaceToNode(doc)
	if err != nil {
		t.Fatalf("InterfaceToNode() error = %v", err)
	}

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected *ast.ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 2 {
		t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
	}
}

// TestInterfaceToNode_Array tests Array type conversion
func TestInterfaceToNode_Array(t *testing.T) {
	arr := NewArray().AddString("go").AddInt(42).AddBool(true)

	node, err := InterfaceToNode(arr)
	if err != nil {
		t.Fatalf("InterfaceToNode() error = %v", err)
	}

	arrayNode, ok := node.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected *ast.ArrayDataNode, got %T", node)
	}

	if len(arrayNode.Elements()) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arrayNode.Elements()))
	}
}

// TestInterfaceToNode_Error tests error cases
func TestInterfaceToNode_Error(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"channel", make(chan int)},
		{"function", func() {}},
		{"complex64", complex64(1 + 2i)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InterfaceToNode(tt.input)
			if err == nil {
				t.Error("InterfaceToNode() expected error for unsupported type, got nil")
			}
		})
	}
}

// TestInterfaceToNode_NestedError tests error in nested conversion
func TestInterfaceToNode_NestedError(t *testing.T) {
	// Slice with unsupported type
	input := []interface{}{1, 2, make(chan int)}
	_, err := InterfaceToNode(input)
	if err == nil {
		t.Error("InterfaceToNode() expected error for nested unsupported type, got nil")
	}

	// Map with unsupported type
	input2 := map[string]interface{}{"key": make(chan int)}
	_, err = InterfaceToNode(input2)
	if err == nil {
		t.Error("InterfaceToNode() expected error for nested unsupported type, got nil")
	}
}

// TestRoundTrip_ConversionCycle tests Parse -> NodeToInterface -> InterfaceToNode -> Render
func TestRoundTrip_ConversionCycle(t *testing.T) {
	tests := []string{
		`{"name":"Alice","age":30}`,
		`[1,2,3,"four",true,null]`,
		`{"nested":{"array":[1,2,3]}}`,
		`null`,
		`true`,
		`false`,
		`"string"`,
		`42`,
		`3.14`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			// Parse
			node1, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			// Convert to interface
			iface := NodeToInterface(node1)

			// Convert back to node
			node2, err := InterfaceToNode(iface)
			if err != nil {
				t.Fatalf("InterfaceToNode error: %v", err)
			}

			// Render both and compare
			bytes1, err := Render(node1)
			if err != nil {
				t.Fatalf("Render node1 error: %v", err)
			}

			bytes2, err := Render(node2)
			if err != nil {
				t.Fatalf("Render node2 error: %v", err)
			}

			if string(bytes1) != string(bytes2) {
				t.Errorf("Round-trip mismatch:\n  original: %s\n  result:   %s", string(bytes1), string(bytes2))
			}
		})
	}
}
