package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shapestone/shape-core/pkg/ast"
	shapetokenizer "github.com/shapestone/shape-core/pkg/tokenizer"
	"github.com/shapestone/shape-json/internal/lenient"
	"github.com/shapestone/shape-json/internal/tokenizer"
)

// LenientParser implements lenient JSON parsing with auto-correction.
// It tolerates and fixes common JSON errors:
//   - Trailing commas in objects and arrays
//   - Single-quoted strings
//   - Unquoted object keys
//   - Line and block comments
//   - Duplicate keys (last value wins)
type LenientParser struct {
	tokenizer  *shapetokenizer.Tokenizer
	current    *shapetokenizer.Token
	hasToken   bool
	inputRunes []rune
	collector  *lenient.CorrectionCollector
}

// NewLenientParser creates a lenient JSON parser for the given input string.
func NewLenientParser(input string) *LenientParser {
	stream := shapetokenizer.NewStream(input)
	tok := tokenizer.NewLenientTokenizerWithStream(stream)
	p := &LenientParser{
		tokenizer:  &tok,
		inputRunes: []rune(input),
		collector:  lenient.NewCorrectionCollector(),
	}
	p.advance()
	return p
}

// NewLenientParserFromStream creates a lenient parser from a stream.
func NewLenientParserFromStream(stream shapetokenizer.Stream) *LenientParser {
	tok := tokenizer.NewLenientTokenizerWithStream(stream)
	p := &LenientParser{
		tokenizer: &tok,
		collector: lenient.NewCorrectionCollector(),
	}
	p.advance()
	return p
}

// Parse parses the input leniently, returning the AST and any corrections applied.
func (p *LenientParser) Parse() (ast.SchemaNode, []lenient.Correction, error) {
	node, err := p.parseValue()
	if err != nil {
		return nil, p.collector.Corrections(), err
	}

	token := p.peek()
	if token != nil && p.hasToken {
		return nil, p.collector.Corrections(), fmt.Errorf("unexpected content after JSON value at %s", p.positionStr())
	}

	return node, p.collector.Corrections(), nil
}

func (p *LenientParser) parseValue() (ast.SchemaNode, error) {
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
	case tokenizer.TokenSingleString:
		return p.parseSingleQuotedString(), nil
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

func (p *LenientParser) parseObject() (*ast.ObjectNode, error) {
	startPos := p.position()

	if err := p.expect(tokenizer.TokenLBrace); err != nil {
		return nil, err
	}

	properties := make(map[string]ast.SchemaNode, 8)

	if p.peek().Kind() != tokenizer.TokenRBrace {
		key, value, err := p.parseMember()
		if err != nil {
			return nil, err
		}
		properties[key] = value

		for p.peek().Kind() == tokenizer.TokenComma {
			p.advance() // consume ","

			// Trailing comma: comma followed by closing brace
			if p.peek().Kind() == tokenizer.TokenRBrace {
				p.collector.Add(lenient.CorrectionTrailingComma, p.position(),
					",", "", "removed trailing comma in object")
				break
			}

			key, value, err := p.parseMember()
			if err != nil {
				return nil, fmt.Errorf("in object after comma: %w", err)
			}

			if _, exists := properties[key]; exists {
				p.collector.Add(lenient.CorrectionDuplicateKey, p.position(),
					key, key, fmt.Sprintf("duplicate key %q, using last value", key))
			}
			properties[key] = value
		}
	}

	if err := p.expect(tokenizer.TokenRBrace); err != nil {
		return nil, err
	}

	return ast.NewObjectNode(properties, startPos), nil
}

func (p *LenientParser) parseMember() (string, ast.SchemaNode, error) {
	var key string

	switch p.peek().Kind() {
	case tokenizer.TokenString:
		keyToken := p.current
		p.advance()
		key = unquoteDoubleString(keyToken.ValueString())

	case tokenizer.TokenSingleString:
		pos := p.position()
		keyToken := p.current
		p.advance()
		key = unquoteSingleString(keyToken.ValueString())
		p.collector.Add(lenient.CorrectionSingleQuote, pos,
			keyToken.ValueString(), fmt.Sprintf("%q", key), "converted single-quoted key to double-quoted")

	case tokenizer.TokenIdentifier:
		pos := p.position()
		key = p.current.ValueString()
		p.advance()
		p.collector.Add(lenient.CorrectionUnquotedKey, pos,
			key, fmt.Sprintf("%q", key), fmt.Sprintf("quoted bare key %q", key))

	default:
		return "", nil, fmt.Errorf("object key must be string at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	if err := p.expect(tokenizer.TokenColon); err != nil {
		return "", nil, fmt.Errorf("expected ':' after object key %q: %w", key, err)
	}

	value, err := p.parseValue()
	if err != nil {
		return "", nil, fmt.Errorf("in value for key %q: %w", key, err)
	}

	return key, value, nil
}

func (p *LenientParser) parseArray() (ast.SchemaNode, error) {
	startPos := p.position()

	if err := p.expect(tokenizer.TokenLBracket); err != nil {
		return nil, err
	}

	elements := make([]ast.SchemaNode, 0, 16)

	if p.peek().Kind() != tokenizer.TokenRBracket {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		elements = append(elements, value)

		for p.peek().Kind() == tokenizer.TokenComma {
			p.advance() // consume ","

			// Trailing comma: comma followed by closing bracket
			if p.peek().Kind() == tokenizer.TokenRBracket {
				p.collector.Add(lenient.CorrectionTrailingComma, p.position(),
					",", "", "removed trailing comma in array")
				break
			}

			value, err := p.parseValue()
			if err != nil {
				return nil, fmt.Errorf("in array element %d: %w", len(elements), err)
			}
			elements = append(elements, value)
		}
	}

	if err := p.expect(tokenizer.TokenRBracket); err != nil {
		return nil, err
	}

	return ast.NewArrayDataNode(elements, startPos), nil
}

func (p *LenientParser) parseString() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenString {
		return nil, fmt.Errorf("expected string at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	pos := p.position()
	offset := p.current.Offset()
	tokenValue := p.current.ValueString()
	p.advance()

	// Check if the next token is structurally valid after a string value.
	// If not, this may be an unescaped quote situation.
	if len(p.inputRunes) > 0 && p.hasToken && !p.isStructuralFollower() {
		if recovered, ok := p.tryRecoverUnescapedQuote(offset); ok {
			return ast.NewLiteralNode(recovered, pos), nil
		}
	}

	unquoted := unquoteDoubleString(tokenValue)
	return ast.NewLiteralNode(unquoted, pos), nil
}

// isStructuralFollower returns true if the current token can validly follow a
// string/value in JSON (comma, closing brace/bracket, colon, or EOF).
func (p *LenientParser) isStructuralFollower() bool {
	if !p.hasToken {
		return true // EOF is valid
	}
	kind := p.peek().Kind()
	return kind == tokenizer.TokenComma ||
		kind == tokenizer.TokenRBrace ||
		kind == tokenizer.TokenRBracket ||
		kind == tokenizer.TokenColon
}

// tryRecoverUnescapedQuote scans the raw input from the opening quote position
// to find the real closing quote. It scans forward to find the first `"` that
// is followed by a structural character (`,`, `}`, `]`, `:`) or EOF.
//
// openOffset is a rune index (from Token.Offset), matching the tokenizer's
// position tracking which counts runes, not bytes.
func (p *LenientParser) tryRecoverUnescapedQuote(openOffset int) (string, bool) {
	if openOffset < 0 || openOffset >= len(p.inputRunes) || p.inputRunes[openOffset] != '"' {
		return "", false
	}

	content := p.inputRunes[openOffset+1:]

	// Skip past the first `"` (the one the tokenizer already found as the close)
	// and keep scanning for a `"` followed by a structural character.
	firstClose := indexRune(content, '"')
	if firstClose < 0 {
		return "", false
	}

	for i := firstClose + 1; i < len(content); i++ {
		if content[i] != '"' {
			continue
		}
		if isStructuralAfter(content[i+1:]) {
			raw := string(content[:i])
			p.collector.Add(lenient.CorrectionUnescapedQuote, ast.NewPosition(openOffset, 0, 0),
				`"`+raw+`"`, "", "escaped unquoted double quotes inside string")

			// Advance the tokenizer past all tokens we consumed during recovery.
			// closeOffset is a rune index matching Token.Offset().
			closeOffset := openOffset + 1 + i + 1
			for p.hasToken && p.current.Offset() < closeOffset {
				p.advance()
			}

			return raw, true
		}
	}

	return "", false
}

func indexRune(runes []rune, target rune) int {
	for i, r := range runes {
		if r == target {
			return i
		}
	}
	return -1
}

func isStructuralAfter(remaining []rune) bool {
	for _, r := range remaining {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		return r == ',' || r == '}' || r == ']' || r == ':'
	}
	return true // EOF is valid
}

func (p *LenientParser) parseSingleQuotedString() *ast.LiteralNode {
	pos := p.position()
	tokenValue := p.current.ValueString()
	p.advance()

	unquoted := unquoteSingleString(tokenValue)
	p.collector.Add(lenient.CorrectionSingleQuote, pos,
		tokenValue, fmt.Sprintf("%q", unquoted), "converted single-quoted string to double-quoted")

	return ast.NewLiteralNode(unquoted, pos)
}

func (p *LenientParser) parseNumber() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenNumber {
		return nil, fmt.Errorf("expected number at %s, got %s",
			p.positionStr(), p.peek().Kind())
	}

	pos := p.position()
	tokenValue := p.current.ValueString()
	p.advance()

	if !strings.Contains(tokenValue, ".") && !strings.ContainsAny(tokenValue, "eE") {
		i, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at %s: %w", tokenValue, pos.String(), err)
		}
		return ast.NewLiteralNode(i, pos), nil
	}

	f, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q at %s: %w", tokenValue, pos.String(), err)
	}
	return ast.NewLiteralNode(f, pos), nil
}

func (p *LenientParser) parseBoolean() (*ast.LiteralNode, error) {
	kind := p.peek().Kind()
	if kind != tokenizer.TokenTrue && kind != tokenizer.TokenFalse {
		return nil, fmt.Errorf("expected boolean at %s, got %s", p.positionStr(), kind)
	}
	pos := p.position()
	value := kind == tokenizer.TokenTrue
	p.advance()
	return ast.NewLiteralNode(value, pos), nil
}

func (p *LenientParser) parseNull() (*ast.LiteralNode, error) {
	if p.peek().Kind() != tokenizer.TokenNull {
		return nil, fmt.Errorf("expected null at %s, got %s", p.positionStr(), p.peek().Kind())
	}
	pos := p.position()
	p.advance()
	return ast.NewLiteralNode(nil, pos), nil
}

// peek returns the current token, skipping whitespace and comments.
func (p *LenientParser) peek() *shapetokenizer.Token {
	for p.hasToken {
		kind := p.current.Kind()
		if kind == "Whitespace" {
			p.advance()
			continue
		}
		if kind == tokenizer.TokenComment {
			p.collector.Add(commentCorrectionKind(p.current.ValueString()), p.position(),
				p.current.ValueString(), "", "removed comment")
			p.advance()
			continue
		}
		break
	}
	return p.current
}

func commentCorrectionKind(value string) lenient.CorrectionKind {
	if strings.HasPrefix(value, "//") {
		return lenient.CorrectionLineComment
	}
	return lenient.CorrectionBlockComment
}

func (p *LenientParser) advance() {
	token, ok := p.tokenizer.NextToken()
	if ok {
		p.current = token
		p.hasToken = true
	} else {
		p.hasToken = false
	}
}

func (p *LenientParser) expect(kind string) error {
	if p.peek().Kind() != kind {
		return fmt.Errorf("expected %s at %s, got %s",
			kind, p.positionStr(), p.peek().Kind())
	}
	p.advance()
	return nil
}

func (p *LenientParser) position() ast.Position {
	if p.hasToken {
		return ast.NewPosition(
			p.current.Offset(),
			p.current.Row(),
			p.current.Column(),
		)
	}
	return ast.ZeroPosition()
}

func (p *LenientParser) positionStr() string {
	return p.position().String()
}

// unquoteDoubleString removes surrounding double quotes and processes escape sequences.
func unquoteDoubleString(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if !strings.ContainsRune(s, '\\') {
		return s
	}

	var buf strings.Builder
	buf.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			buf.WriteByte(s[i])
			continue
		}

		i++
		if i >= len(s) {
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
			if i+4 < len(s) {
				hex := s[i+1 : i+5]
				if codepoint, err := parseHex(hex); err == nil {
					buf.WriteRune(rune(codepoint))
					i += 4
				} else {
					buf.WriteString("\\u")
				}
			} else {
				buf.WriteString("\\u")
			}
		default:
			buf.WriteByte('\\')
			buf.WriteByte(s[i])
		}
	}

	return buf.String()
}

// unquoteSingleString removes surrounding single quotes and processes escape sequences.
func unquoteSingleString(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}

	if !strings.ContainsRune(s, '\\') {
		return s
	}

	var buf strings.Builder
	buf.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			buf.WriteByte(s[i])
			continue
		}

		i++
		if i >= len(s) {
			buf.WriteByte('\\')
			break
		}

		switch s[i] {
		case '\'':
			buf.WriteByte('\'')
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
			if i+4 < len(s) {
				hex := s[i+1 : i+5]
				if codepoint, err := parseHex(hex); err == nil {
					buf.WriteRune(rune(codepoint))
					i += 4
				} else {
					buf.WriteString("\\u")
				}
			} else {
				buf.WriteString("\\u")
			}
		default:
			buf.WriteByte('\\')
			buf.WriteByte(s[i])
		}
	}

	return buf.String()
}
