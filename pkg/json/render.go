// Package json provides AST rendering to JSON bytes.
//
// This file implements the core JSON rendering functionality, converting
// Shape AST nodes back into JSON byte representations.
package json

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shapestone/shape-core/pkg/ast"
)

// Render converts an AST node to compact JSON bytes.
//
// The node should be the result of Parse() or ParseReader().
// Returns JSON bytes with no unnecessary whitespace.
//
// Example:
//
//	node, _ := json.Parse(`{"name": "Alice", "age": 30}`)
//	bytes, _ := json.Render(node)
//	// bytes: {"age":30,"name":"Alice"}
func Render(node ast.SchemaNode) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	if err := renderNode(node, buf, false, "", ""); err != nil {
		return nil, err
	}

	// Must copy since buffer will be returned to pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// RenderIndent converts an AST node to pretty-printed JSON bytes with indentation.
//
// The prefix is added to the beginning of each line, and indent specifies
// the indentation string (typically spaces or tabs).
//
// Common usage:
//   - RenderIndent(node, "", "  ") - 2-space indentation
//   - RenderIndent(node, "", "\t") - tab indentation
//   - RenderIndent(node, ">>", "  ") - prefix each line with ">>"
//
// Example:
//
//	node, _ := json.Parse(`{"name":"Alice","age":30}`)
//	bytes, _ := json.RenderIndent(node, "", "  ")
//	// Output:
//	// {
//	//   "age": 30,
//	//   "name": "Alice"
//	// }
func RenderIndent(node ast.SchemaNode, prefix, indent string) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	if err := renderNode(node, buf, true, prefix, indent); err != nil {
		return nil, err
	}

	// Must copy since buffer will be returned to pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// renderNode recursively renders an AST node to the buffer.
//
// Parameters:
//   - node: The AST node to render
//   - buf: The output buffer
//   - prettyPrint: Whether to add whitespace for readability
//   - prefix: String to add at the start of each line
//   - indent: Indentation string (spaces or tabs)
func renderNode(node ast.SchemaNode, buf *bytes.Buffer, prettyPrint bool, prefix, indent string) error {
	return renderNodeWithDepth(node, buf, prettyPrint, prefix, indent, 0)
}

// renderNodeWithDepth renders a node with tracking of indentation depth.
func renderNodeWithDepth(node ast.SchemaNode, buf *bytes.Buffer, prettyPrint bool, prefix, indent string, depth int) error {
	if node == nil {
		buf.WriteString("null")
		return nil
	}

	switch n := node.(type) {
	case *ast.ObjectNode:
		return renderObject(n, buf, prettyPrint, prefix, indent, depth)
	case *ast.ArrayDataNode:
		return renderArrayData(n, buf, prettyPrint, prefix, indent, depth)
	case *ast.LiteralNode:
		return renderLiteral(n, buf)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

// renderObject renders an ObjectNode as either a JSON object or array.
//
// Arrays are detected by checking if all keys are sequential numeric strings ("0", "1", "2", ...).
func renderObject(node *ast.ObjectNode, buf *bytes.Buffer, prettyPrint bool, prefix, indent string, depth int) error {
	props := node.Properties()

	// Empty object
	if len(props) == 0 {
		buf.WriteString("{}")
		return nil
	}

	// Check if this is an array (all keys are sequential numbers)
	if isArray(props) {
		return renderArray(node, buf, prettyPrint, prefix, indent, depth)
	}

	// Render as object
	buf.WriteString("{")

	// Sort keys for consistent output
	keys := make([]string, 0, len(props))
	for key := range props {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	first := true
	for _, key := range keys {
		if !first {
			buf.WriteString(",")
		}
		first = false

		if prettyPrint {
			buf.WriteString("\n")
			buf.WriteString(prefix)
			buf.WriteString(strings.Repeat(indent, depth+1))
		}

		// Write key
		buf.WriteString(`"`)
		buf.WriteString(escapeString(key))
		buf.WriteString(`"`)
		buf.WriteString(":")
		if prettyPrint {
			buf.WriteString(" ")
		}

		// Write value
		if err := renderNodeWithDepth(props[key], buf, prettyPrint, prefix, indent, depth+1); err != nil {
			return err
		}
	}

	if prettyPrint {
		buf.WriteString("\n")
		buf.WriteString(prefix)
		buf.WriteString(strings.Repeat(indent, depth))
	}
	buf.WriteString("}")
	return nil
}

// renderArray renders an ObjectNode with numeric keys as a JSON array.
func renderArray(node *ast.ObjectNode, buf *bytes.Buffer, prettyPrint bool, prefix, indent string, depth int) error {
	props := node.Properties()

	buf.WriteString("[")

	// Get array length (highest numeric key + 1)
	length := len(props)
	if length == 0 {
		buf.WriteString("]")
		return nil
	}

	// Render elements in order
	for i := 0; i < length; i++ {
		if i > 0 {
			buf.WriteString(",")
		}

		if prettyPrint {
			buf.WriteString("\n")
			buf.WriteString(prefix)
			buf.WriteString(strings.Repeat(indent, depth+1))
		}

		key := strconv.Itoa(i)
		value, ok := props[key]
		if !ok {
			// Missing element - should not happen in valid arrays
			buf.WriteString("null")
		} else {
			if err := renderNodeWithDepth(value, buf, prettyPrint, prefix, indent, depth+1); err != nil {
				return err
			}
		}
	}

	if prettyPrint {
		buf.WriteString("\n")
		buf.WriteString(prefix)
		buf.WriteString(strings.Repeat(indent, depth))
	}
	buf.WriteString("]")
	return nil
}

// renderArrayData renders an ArrayDataNode as a JSON array.
func renderArrayData(node *ast.ArrayDataNode, buf *bytes.Buffer, prettyPrint bool, prefix, indent string, depth int) error {
	elements := node.Elements()

	buf.WriteString("[")

	if len(elements) == 0 {
		buf.WriteString("]")
		return nil
	}

	// Render each element
	for i, elem := range elements {
		if i > 0 {
			buf.WriteString(",")
		}

		if prettyPrint {
			buf.WriteString("\n")
			buf.WriteString(prefix)
			buf.WriteString(strings.Repeat(indent, depth+1))
		}

		if err := renderNodeWithDepth(elem, buf, prettyPrint, prefix, indent, depth+1); err != nil {
			return err
		}
	}

	if prettyPrint {
		buf.WriteString("\n")
		buf.WriteString(prefix)
		buf.WriteString(strings.Repeat(indent, depth))
	}
	buf.WriteString("]")
	return nil
}

// renderLiteral renders a LiteralNode as a JSON primitive.
func renderLiteral(node *ast.LiteralNode, buf *bytes.Buffer) error {
	value := node.Value()

	if value == nil {
		buf.WriteString("null")
		return nil
	}

	switch v := value.(type) {
	case string:
		buf.WriteString(`"`)
		buf.WriteString(escapeString(v))
		buf.WriteString(`"`)
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		// Format float, preserving precision
		s := strconv.FormatFloat(v, 'f', -1, 64)
		// Avoid scientific notation for readability
		if !strings.Contains(s, ".") && !strings.Contains(s, "e") && !strings.Contains(s, "E") {
			// It's a whole number stored as float - keep it that way
			s = strconv.FormatFloat(v, 'f', 1, 64)
		}
		buf.WriteString(s)
	case bool:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	default:
		return fmt.Errorf("unknown literal type: %T", value)
	}

	return nil
}

// Note: isArray function is defined in unmarshal.go and shared across the package

// escapeString escapes special characters in a string for JSON.
//
// Handles:
//   - Quotation mark (") → \"
//   - Backslash (\) → \\
//   - Forward slash (/) → \/
//   - Backspace (\b) → \b
//   - Form feed (\f) → \f
//   - Newline (\n) → \n
//   - Carriage return (\r) → \r
//   - Tab (\t) → \t
//   - Control characters (U+0000 to U+001F) → \uXXXX
//   - Unicode characters beyond ASCII (optional, for safety)
func escapeString(s string) string {
	// Quick check: if string has no special chars, return as-is
	if !needsEscaping(s) {
		return s
	}

	var buf bytes.Buffer
	buf.Grow(len(s) + 10) // Pre-allocate with some extra space

	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '/':
			// Forward slash can optionally be escaped
			buf.WriteString(`\/`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				// Control characters (U+0000 to U+001F)
				buf.WriteString(fmt.Sprintf(`\u%04x`, r))
			} else {
				// Regular character (including Unicode)
				buf.WriteRune(r)
			}
		}
	}

	return buf.String()
}

// needsEscaping checks if a string contains characters that need escaping.
func needsEscaping(s string) bool {
	for _, r := range s {
		if r == '"' || r == '\\' || r == '/' || r == '\b' || r == '\f' || r == '\n' || r == '\r' || r == '\t' || r < 0x20 {
			return true
		}
	}
	return false
}
