// Package parser implements LL(1) recursive descent parsing for JSON format.
// Each production rule in the grammar (docs/grammar/json.ebnf) corresponds to a parse function.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shapestone/shape-core/pkg/ast"
	shapetokenizer "github.com/shapestone/shape-core/pkg/tokenizer"
	"github.com/shapestone/shape-json/internal/tokenizer"
)

// Parser implements LL(1) recursive descent parsing for JSON.
// It maintains a single token lookahead for predictive parsing.
type Parser struct {
	tokenizer *shapetokenizer.Tokenizer
	current   *shapetokenizer.Token
	hasToken  bool
}

// NewParser creates a new JSON parser for the given input string.
// For parsing from io.Reader, use NewParserFromStream instead.
func NewParser(input string) *Parser {
	return newParserWithStream(shapetokenizer.NewStream(input))
}

// NewParserFromStream creates a new JSON parser using a pre-configured stream.
// This allows parsing from io.Reader using tokenizer.NewStreamFromReader.
func NewParserFromStream(stream shapetokenizer.Stream) *Parser {
	return newParserWithStream(stream)
}

// newParserWithStream is the internal constructor that accepts a stream.
func newParserWithStream(stream shapetokenizer.Stream) *Parser {
	tok := tokenizer.NewTokenizerWithStream(stream)

	p := &Parser{
		tokenizer: &tok,
	}
	p.advance() // Load first token
	return p
}

// Parse parses the input and returns an AST representing the JSON value.
//
// Grammar:
//
//	Value = Object | Array | String | Number | Boolean | Null ;
//
// Returns ast.SchemaNode - the root of the AST.
// For JSON data, this will be ObjectNode, ArrayNode, or LiteralNode.
func (p *Parser) Parse() (ast.SchemaNode, error) {
	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// After parsing the value, we should be at EOF
	// peek() skips whitespace, so if we have a non-nil token after peek, it's extra content
	token := p.peek()
	if token != nil && p.hasToken {
		return nil, fmt.Errorf("unexpected content after JSON value at %s", p.positionStr())
	}

	return node, nil
}

// parseValue dispatches to specific parse functions based on the current token.
//
// Grammar:
//
//	Value = Object | Array | String | Number | Boolean | Null ;
//
// Uses single token lookahead (LL(1) predictive parsing).
func (p *Parser) parseValue() (ast.SchemaNode, error) {
	token := p.peek()
	if token == nil || !p.hasToken {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch token.Kind() {
	case tokenizer.TokenLBrace:
		return p.parseObject()
	case tokenizer.TokenLBracket:
		return p.parseArray()
	case tokenizer.TokenString:
		return p.parseString()
	case tokenizer.TokenNumber:
		return p.parseNumber()
	case tokenizer.TokenTrue, tokenizer.TokenFalse:
		return p.parseBoolean()
	case tokenizer.TokenNull:
		return p.parseNull()
	default:
		return nil, fmt.Errorf("expected JSON value at %s, got %s",
			p.positionStr(), token.Kind())
	}
}

// parseObject parses a JSON object.
//
// Grammar:
//
//	Object = "{" [ Member { "," Member } ] "}" ;
//
// Returns *ast.ObjectNode with properties map.
// Example valid: {}
// Example valid: {"name": "Alice"}
// Example valid: {"id": 123, "active": true}
// Example invalid: {name: "Alice"} (keys must be quoted)
// Example invalid: {"trailing": "comma",} (trailing commas not allowed)
func (p *Parser) parseObject() (*ast.ObjectNode, error) {
	startPos := p.position()

	// "{"
	if err := p.expect(tokenizer.TokenLBrace); err != nil {
		return nil, err
	}

	// Pre-size with reasonable capacity to avoid initial resizing
	properties := make(map[string]ast.SchemaNode, 8)

	// [ Member { "," Member } ]  - Optional member list
	if p.peek().Kind() != tokenizer.TokenRBrace {
		// First member
		key, value, err := p.parseMember()
		if err != nil {
			return nil, err
		}
		properties[key] = value

		// Additional members: { "," Member }
		for p.peek().Kind() == tokenizer.TokenComma {
			p.advance() // consume ","

			key, value, err := p.parseMember()
			if err != nil {
				return nil, fmt.Errorf("in object after comma: %w", err)
			}

			if _, exists := properties[key]; exists {
				return nil, fmt.Errorf("duplicate key %q in object at %s", key, p.positionStr())
			}
			properties[key] = value
		}
	}

	// "}"
	if err := p.expect(tokenizer.TokenRBrace); err != nil {
		return nil, err
	}

	return ast.NewObjectNode(properties, startPos), nil
}

// parseMember parses an object member (key-value pair).
//
// Grammar:
//
//	Member = String ":" Value ;
//
// Returns (key string, value ast.SchemaNode).
func (p *Parser) parseMember() (string, ast.SchemaNode, error) {
	// String (key)
	if p.peek().Kind() != tokenizer.TokenString {
		return "", nil, fmt.Errorf("object key must be string at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	keyToken := p.current
	p.advance()
	key := p.unquoteString(keyToken.ValueString())

	// ":"
	if err := p.expect(tokenizer.TokenColon); err != nil {
		return "", nil, fmt.Errorf("expected ':' after object key %q: %w", key, err)
	}

	// Value
	value, err := p.parseValue()
	if err != nil {
		return "", nil, fmt.Errorf("in value for key %q: %w", key, err)
	}

	return key, value, nil
}

// parseArray parses a JSON array.
//
// Grammar:
//
//	Array = "[" [ Value { "," Value } ] "]" ;
//
// Returns *ast.ArrayDataNode.
// Arrays are represented using ArrayDataNode which stores actual elements,
// properly distinguishing them from objects (including empty arrays vs empty objects).
// Example valid: []
// Example valid: [1, 2, 3]
// Example valid: [{"id": 1}, {"id": 2}]
// Example invalid: [1, 2, 3,] (trailing commas not allowed)
func (p *Parser) parseArray() (ast.SchemaNode, error) {
	startPos := p.position()

	// "["
	if err := p.expect(tokenizer.TokenLBracket); err != nil {
		return nil, err
	}

	// Build array of elements with reasonable initial capacity
	elements := make([]ast.SchemaNode, 0, 16)

	// [ Value { "," Value } ]  - Optional value list
	if p.peek().Kind() != tokenizer.TokenRBracket {
		// First value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		elements = append(elements, value)

		// Additional values: { "," Value }
		for p.peek().Kind() == tokenizer.TokenComma {
			p.advance() // consume ","

			value, err := p.parseValue()
			if err != nil {
				return nil, fmt.Errorf("in array element %d: %w", len(elements), err)
			}
			elements = append(elements, value)
		}
	}

	// "]"
	if err := p.expect(tokenizer.TokenRBracket); err != nil {
		return nil, err
	}

	// Return ArrayDataNode with actual elements
	return ast.NewArrayDataNode(elements, startPos), nil
}

// parseString parses a JSON string literal.
//
// Grammar:
//
//	String = '"' { Character } '"' ;
//
// Returns *ast.LiteralNode with the unescaped string value.
// Handles escape sequences: \", \\, \/, \b, \f, \n, \r, \t, \uXXXX
func (p *Parser) parseString() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenString {
		return nil, fmt.Errorf("expected string at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	pos := p.position()
	tokenValue := p.current.ValueString()
	p.advance()

	// Unquote and unescape the string
	unquoted := p.unquoteString(tokenValue)

	return ast.NewLiteralNode(unquoted, pos), nil
}

// parseNumber parses a JSON number literal.
//
// Grammar:
//
//	Number = [ "-" ] Integer [ Fraction ] [ Exponent ] ;
//
// Returns *ast.LiteralNode with int64 or float64 value.
// Examples: 0, -123, 123.456, 1e10, 1.5e-3
func (p *Parser) parseNumber() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenNumber {
		return nil, fmt.Errorf("expected number at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	pos := p.position()
	tokenValue := p.current.ValueString()
	p.advance()

	// Try parsing as integer first
	if !strings.Contains(tokenValue, ".") && !strings.ContainsAny(tokenValue, "eE") {
		i, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at %s: %w", tokenValue, pos.String(), err)
		}
		return ast.NewLiteralNode(i, pos), nil
	}

	// Parse as floating point
	f, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q at %s: %w", tokenValue, pos.String(), err)
	}
	return ast.NewLiteralNode(f, pos), nil
}

// parseBoolean parses a JSON boolean literal.
//
// Grammar:
//
//	Boolean = "true" | "false" ;
//
// Returns *ast.LiteralNode with bool value.
func (p *Parser) parseBoolean() (*ast.LiteralNode, error) {
	kind := p.peek().Kind()
	if kind != tokenizer.TokenTrue && kind != tokenizer.TokenFalse {
		return nil, fmt.Errorf("expected boolean at %s, got %s",
			p.positionStr(), kind)
	}

	pos := p.position()
	value := kind == tokenizer.TokenTrue
	p.advance()

	return ast.NewLiteralNode(value, pos), nil
}

// parseNull parses a JSON null literal.
//
// Grammar:
//
//	Null = "null" ;
//
// Returns *ast.LiteralNode with nil value.
func (p *Parser) parseNull() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenNull {
		return nil, fmt.Errorf("expected null at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	pos := p.position()
	p.advance()

	return ast.NewLiteralNode(nil, pos), nil
}

// Helper methods

// peek returns current token without advancing.
// Automatically skips whitespace tokens.
func (p *Parser) peek() *shapetokenizer.Token {
	// Skip whitespace tokens (WhiteSpaceMatcher still creates them, but parser ignores them)
	for p.hasToken && p.current.Kind() == "Whitespace" {
		p.advance()
	}
	return p.current
}

// advance moves to next token.
func (p *Parser) advance() {
	token, ok := p.tokenizer.NextToken()
	if ok {
		p.current = token
		p.hasToken = true
	} else {
		p.hasToken = false
	}
}

// expect consumes token of expected kind or returns error.
func (p *Parser) expect(kind string) error {
	if p.peek().Kind() != kind {
		return fmt.Errorf("expected %s at %s, got %s",
			kind, p.positionStr(), p.peek().Kind())
	}
	p.advance()
	return nil
}

// position returns current position for AST nodes.
func (p *Parser) position() ast.Position {
	if p.hasToken {
		return ast.NewPosition(
			p.current.Offset(),
			p.current.Row(),
			p.current.Column(),
		)
	}
	return ast.ZeroPosition()
}

// positionStr returns current position as a string for error messages.
func (p *Parser) positionStr() string {
	return p.position().String()
}

// unquoteString removes quotes and unescapes a JSON string.
// Handles: \", \\, \/, \b, \f, \n, \r, \t, \uXXXX
// Uses single-pass algorithm for optimal performance (5-10x faster than multiple ReplaceAll calls).
func (p *Parser) unquoteString(s string) string {
	// Remove surrounding quotes
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Fast path: no escapes
	if !strings.ContainsRune(s, '\\') {
		return s
	}

	// Single-pass escape processing
	var buf strings.Builder
	buf.Grow(len(s)) // Pre-allocate to avoid resizing

	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			buf.WriteByte(s[i])
			continue
		}

		// Handle escape sequence
		i++ // Skip backslash
		if i >= len(s) {
			// Malformed escape at end of string
			buf.WriteByte('\\')
			break
		}

		switch s[i] {
		case '"', '\\', '/':
			buf.WriteByte(s[i])
		case 'b':
			buf.WriteByte('\b')
		case 'f':
			buf.WriteByte('\f')
		case 'n':
			buf.WriteByte('\n')
		case 'r':
			buf.WriteByte('\r')
		case 't':
			buf.WriteByte('\t')
		case 'u':
			// Handle \uXXXX unicode escape
			if i+4 < len(s) {
				// Parse 4 hex digits
				hex := s[i+1 : i+5]
				if codepoint, err := parseHex(hex); err == nil {
					buf.WriteRune(rune(codepoint))
					i += 4 // Skip the 4 hex digits
				} else {
					// Invalid hex, write as-is
					buf.WriteString("\\u")
				}
			} else {
				// Not enough characters for \uXXXX
				buf.WriteString("\\u")
			}
		default:
			// Unknown escape sequence, preserve it
			buf.WriteByte('\\')
			buf.WriteByte(s[i])
		}
	}

	return buf.String()
}

// parseHex converts a 4-character hex string to an integer.
func parseHex(s string) (int, error) {
	if len(s) != 4 {
		return 0, fmt.Errorf("hex string must be 4 characters")
	}

	var result int
	for i := 0; i < 4; i++ {
		c := s[i]
		var digit int

		switch {
		case c >= '0' && c <= '9':
			digit = int(c - '0')
		case c >= 'a' && c <= 'f':
			digit = int(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			digit = int(c - 'A' + 10)
		default:
			return 0, fmt.Errorf("invalid hex character: %c", c)
		}

		result = result*16 + digit
	}

	return result, nil
}
