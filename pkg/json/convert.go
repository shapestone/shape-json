// Package json provides conversion between AST nodes and Go native types.
package json

import (
	"fmt"
	"strconv"

	"github.com/shapestone/shape-core/pkg/ast"
)

// NodeToInterface converts an AST node to native Go types.
//
// Converts:
//   - *ast.LiteralNode → primitives (string, int64, float64, bool, nil)
//   - *ast.ArrayDataNode → []interface{}
//   - *ast.ObjectNode (array - legacy) → []interface{}
//   - *ast.ObjectNode (object) → map[string]interface{}
//
// This function recursively processes nested structures.
//
// Example:
//
//	node, _ := json.Parse(`{"name":"Alice","tags":["go","json"]}`)
//	data := json.NodeToInterface(node)
//	// data is map[string]interface{}{"name":"Alice", "tags":[]interface{}{"go","json"}}
func NodeToInterface(node ast.SchemaNode) interface{} {
	switch n := node.(type) {
	case *ast.LiteralNode:
		val := n.Value()
		// Ensure numbers are returned as appropriate types
		if f, ok := val.(float64); ok {
			// Check if it's a whole number
			if f == float64(int64(f)) {
				return int64(f)
			}
		}
		return val

	case *ast.ArrayDataNode:
		// Convert ArrayDataNode to []interface{}
		elements := n.Elements()
		arr := make([]interface{}, len(elements))
		for i, elem := range elements {
			arr[i] = NodeToInterface(elem)
		}
		return arr

	case *ast.ObjectNode:
		props := n.Properties()

		// Check if this represents an array (sequential numeric keys - legacy support)
		if isArray(props) {
			arr := make([]interface{}, len(props))
			for i := 0; i < len(props); i++ {
				key := strconv.Itoa(i)
				if propNode, ok := props[key]; ok {
					arr[i] = NodeToInterface(propNode)
				}
			}
			return arr
		}

		// Otherwise it's a map/object
		m := make(map[string]interface{}, len(props))
		for key, propNode := range props {
			m[key] = NodeToInterface(propNode)
		}
		return m

	default:
		return nil
	}
}

// ReleaseTree recursively releases all nodes in an AST tree back to their pools.
// This should be called when you're completely done with an AST (after conversion,
// rendering, etc.) to enable node reuse and reduce memory pressure.
//
// Example:
//
//	node, _ := json.Parse(`{"name":"Alice"}`)
//	data := json.NodeToInterface(node)
//	json.ReleaseTree(node)  // Release nodes back to pool
func ReleaseTree(node ast.SchemaNode) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.LiteralNode:
		ast.ReleaseLiteralNode(n)

	case *ast.ArrayDataNode:
		// Release children first
		for _, elem := range n.Elements() {
			ReleaseTree(elem)
		}
		ast.ReleaseArrayDataNode(n)

	case *ast.ObjectNode:
		// Release children first
		for _, child := range n.Properties() {
			ReleaseTree(child)
		}
		ast.ReleaseObjectNode(n)
	}
}

// InterfaceToNode converts native Go types to AST nodes.
//
// Converts:
//   - string → *ast.LiteralNode
//   - int, int64, int32, etc → *ast.LiteralNode
//   - float64, float32 → *ast.LiteralNode
//   - bool → *ast.LiteralNode
//   - nil → *ast.LiteralNode
//   - []interface{} → *ast.ArrayDataNode
//   - map[string]interface{} → *ast.ObjectNode
//   - *Document → *ast.ObjectNode
//   - *Array → *ast.ArrayDataNode
//
// This function recursively processes nested structures.
//
// Example:
//
//	data := map[string]interface{}{
//	    "name": "Alice",
//	    "tags": []interface{}{"go", "json"},
//	}
//	node, _ := json.InterfaceToNode(data)
//	// node is an *ast.ObjectNode representing the data
func InterfaceToNode(v interface{}) (ast.SchemaNode, error) {
	// Use empty position since we're creating nodes programmatically
	pos := ast.Position{}

	if v == nil {
		return ast.NewLiteralNode(nil, pos), nil
	}

	switch val := v.(type) {
	// Handle strings
	case string:
		return ast.NewLiteralNode(val, pos), nil

	// Handle booleans
	case bool:
		return ast.NewLiteralNode(val, pos), nil

	// Handle integers
	case int:
		return ast.NewLiteralNode(int64(val), pos), nil
	case int64:
		return ast.NewLiteralNode(val, pos), nil
	case int32:
		return ast.NewLiteralNode(int64(val), pos), nil
	case int16:
		return ast.NewLiteralNode(int64(val), pos), nil
	case int8:
		return ast.NewLiteralNode(int64(val), pos), nil

	// Handle unsigned integers
	case uint:
		return ast.NewLiteralNode(int64(val), pos), nil
	case uint64:
		return ast.NewLiteralNode(int64(val), pos), nil
	case uint32:
		return ast.NewLiteralNode(int64(val), pos), nil
	case uint16:
		return ast.NewLiteralNode(int64(val), pos), nil
	case uint8:
		return ast.NewLiteralNode(int64(val), pos), nil

	// Handle floats
	case float64:
		return ast.NewLiteralNode(val, pos), nil
	case float32:
		return ast.NewLiteralNode(float64(val), pos), nil

	// Handle slices/arrays
	case []interface{}:
		elements := make([]ast.SchemaNode, len(val))
		for i, item := range val {
			itemNode, err := InterfaceToNode(item)
			if err != nil {
				return nil, fmt.Errorf("array element %d: %w", i, err)
			}
			elements[i] = itemNode
		}
		return ast.NewArrayDataNode(elements, pos), nil

	// Handle maps
	case map[string]interface{}:
		props := make(map[string]ast.SchemaNode)
		for key, value := range val {
			valueNode, err := InterfaceToNode(value)
			if err != nil {
				return nil, fmt.Errorf("object property %s: %w", key, err)
			}
			props[key] = valueNode
		}
		return ast.NewObjectNode(props, pos), nil

	// Handle Document type
	case *Document:
		return InterfaceToNode(val.data)

	// Handle Array type
	case *Array:
		return InterfaceToNode(val.data)

	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
