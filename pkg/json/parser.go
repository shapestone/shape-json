// Package json provides JSON format parsing and AST generation.
//
// This package implements a complete JSON parser following RFC 8259.
// It parses JSON data into Shape's unified AST representation.
//
// Grammar: See docs/grammar/json.ebnf for the complete EBNF specification.
//
// This parser uses LL(1) recursive descent parsing (see Shape ADR 0004).
// Each production rule in the grammar corresponds to a parse function in internal/parser/parser.go.
//
// # Parsing APIs
//
// The package provides two parsing functions:
//
//   - Parse(string) - Parses JSON from a string in memory
//   - ParseReader(io.Reader) - Parses JSON from any io.Reader with streaming support
//
// Use Parse() for small JSON documents that are already in memory as strings.
// Use ParseReader() for large files, network streams, or any io.Reader source.
//
// # Example usage with Parse:
//
//	jsonStr := `{"name": "Alice", "age": 30}`
//	node, err := json.Parse(jsonStr)
//	if err != nil {
//	    // handle error
//	}
//	// node is now a *ast.ObjectNode representing the JSON data
//
// # Example usage with ParseReader:
//
//	file, err := os.Open("data.json")
//	if err != nil {
//	    // handle error
//	}
//	defer file.Close()
//
//	node, err := json.ParseReader(file)
//	if err != nil {
//	    // handle error
//	}
//	// node is now a *ast.ObjectNode representing the JSON data
//
// For more examples, see the examples/parse_reader directory.
package json

import (
	"io"

	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-core/pkg/tokenizer"
	"github.com/shapestone/shape-json/internal/fastparser"
	"github.com/shapestone/shape-json/internal/parser"
)

// Parse parses JSON format into an AST from a string.
//
// The input is a complete JSON value (object, array, string, number, boolean, or null).
//
// Returns an ast.SchemaNode representing the parsed JSON:
//   - *ast.ObjectNode for objects (and arrays, represented as objects with numeric string keys)
//   - *ast.LiteralNode for primitives (string, number, boolean, null)
//
// For parsing large files or streaming data, use ParseReader instead.
//
// Example:
//
//	node, err := json.Parse(`{"name": "Alice", "age": 30}`)
//	obj := node.(*ast.ObjectNode)
//	nameNode, _ := obj.GetProperty("name")
//	name := nameNode.(*ast.LiteralNode).Value().(string) // "Alice"
func Parse(input string) (ast.SchemaNode, error) {
	p := parser.NewParser(input)
	return p.Parse()
}

// ParseReader parses JSON format into an AST from an io.Reader.
//
// This function is designed for parsing large JSON files or streaming data with
// constant memory usage. It uses a buffered stream implementation that reads data
// in chunks, making it suitable for files that don't fit entirely in memory.
//
// The reader can be any io.Reader implementation:
//   - os.File for reading from files
//   - strings.Reader for reading from strings
//   - bytes.Buffer for reading from byte slices
//   - Network streams, compressed streams, etc.
//
// Returns an ast.SchemaNode representing the parsed JSON:
//   - *ast.ObjectNode for objects (and arrays, represented as objects with numeric string keys)
//   - *ast.LiteralNode for primitives (string, number, boolean, null)
//
// Example parsing from a file:
//
//	file, err := os.Open("data.json")
//	if err != nil {
//	    // handle error
//	}
//	defer file.Close()
//
//	node, err := json.ParseReader(file)
//	if err != nil {
//	    // handle error
//	}
//	// ast is now a *ast.ObjectNode representing the JSON data
//
// Example parsing from a string:
//
//	reader := strings.NewReader(`{"name": "Alice", "age": 30}`)
//	node, err := json.ParseReader(reader)
func ParseReader(reader io.Reader) (ast.SchemaNode, error) {
	stream := tokenizer.NewStreamFromReader(reader)
	p := parser.NewParserFromStream(stream)
	return p.Parse()
}

// Format returns the format identifier for this parser.
// Returns "JSON" to identify this as the JSON data format parser.
func Format() string {
	return "JSON"
}

// Validate checks if the input string is valid JSON.
//
// This function uses a high-performance fast path that bypasses AST construction.
//
// Returns nil if the input is valid JSON.
// Returns an error with details about why the JSON is invalid.
//
// This is the idiomatic Go approach - check the error:
//
//	if err := json.Validate(input); err != nil {
//	    // Invalid JSON
//	    fmt.Println("Invalid JSON:", err)
//	}
//	// Valid JSON - err is nil
//
// Valid JSON includes:
//   - Objects: {"key": "value"}
//   - Arrays: [1, 2, 3]
//   - Strings: "hello"
//   - Numbers: 42, 3.14, 1e10
//   - Booleans: true, false
//   - Null: null
func Validate(input string) error {
	// Fast path: Just parse and discard (4-5x faster than AST construction)
	parser := fastparser.NewParser([]byte(input))
	_, err := parser.Parse()
	return err
}

// ValidateReader checks if the input from an io.Reader is valid JSON.
//
// This function uses a high-performance fast path that bypasses AST construction.
//
// Returns nil if the input is valid JSON.
// Returns an error with details about why the JSON is invalid.
// This reads the entire input from the reader.
//
// This is the idiomatic Go approach - check the error:
//
//	file, _ := os.Open("data.json")
//	defer file.Close()
//	if err := json.ValidateReader(file); err != nil {
//	    // Invalid JSON
//	    fmt.Println("Invalid JSON:", err)
//	}
//	// Valid JSON - err is nil
//
// For large streams, consider using ParseReader directly and handling errors.
func ValidateReader(reader io.Reader) error {
	// Fast path: Read all data and parse without AST
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	parser := fastparser.NewParser(data)
	_, err = parser.Parse()
	return err
}

// DetectFormat attempts to detect if the input is valid JSON.
//
// Deprecated: Use Validate() instead. This function returns "JSON" on success,
// which is redundant. The idiomatic Go approach is:
//
//	if err := json.Validate(input); err != nil {
//	    // Invalid
//	}
func DetectFormat(input string) (string, error) {
	err := Validate(input)
	if err != nil {
		return "", err
	}
	return "JSON", nil
}

// DetectFormatFromReader attempts to detect if input from io.Reader is valid JSON.
//
// Deprecated: Use ValidateReader() instead. This function returns "JSON" on success,
// which is redundant. The idiomatic Go approach is:
//
//	if err := json.ValidateReader(reader); err != nil {
//	    // Invalid
//	}
func DetectFormatFromReader(reader io.Reader) (string, error) {
	err := ValidateReader(reader)
	if err != nil {
		return "", err
	}
	return "JSON", nil
}
