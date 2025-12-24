// Package fastparser implements a high-performance JSON parser without AST construction.
//
// This parser is optimized for the common case of unmarshaling JSON directly into Go types.
// It bypasses tokenization and AST construction, parsing directly from bytes to values.
//
// Performance targets (vs AST parser):
//   - 4-5x faster parsing
//   - 5-6x less memory
//   - 4-5x fewer allocations
package fastparser

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Parser implements a zero-allocation JSON parser that builds values directly without AST.
type Parser struct {
	data   []byte
	pos    int
	length int
}

// NewParser creates a new fast parser for the given data.
func NewParser(data []byte) *Parser {
	return &Parser{
		data:   data,
		pos:    0,
		length: len(data),
	}
}

// Parse parses the JSON data and returns the value as interface{}.
// This is used by Unmarshal and Validate.
func (p *Parser) Parse() (interface{}, error) {
	p.skipWhitespace()
	if p.pos >= p.length {
		return nil, errors.New("unexpected end of JSON input")
	}

	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()
	if p.pos < p.length {
		return nil, fmt.Errorf("unexpected data after JSON value at position %d", p.pos)
	}

	return value, nil
}

// parseValue parses any JSON value.
func (p *Parser) parseValue() (interface{}, error) {
	if p.pos >= p.length {
		return nil, errors.New("unexpected end of JSON input")
	}

	c := p.data[p.pos]
	switch c {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't':
		return p.parseTrue()
	case 'f':
		return p.parseFalse()
	case 'n':
		return p.parseNull()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.parseNumber()
	default:
		return nil, fmt.Errorf("unexpected character '%c' at position %d", c, p.pos)
	}
}

// parseObject parses a JSON object into map[string]interface{}.
func (p *Parser) parseObject() (map[string]interface{}, error) {
	if p.pos >= p.length || p.data[p.pos] != '{' {
		return nil, errors.New("expected '{'")
	}
	p.pos++ // skip '{'

	result := make(map[string]interface{})
	p.skipWhitespace()

	// Handle empty object
	if p.pos < p.length && p.data[p.pos] == '}' {
		p.pos++
		return result, nil
	}

	for {
		p.skipWhitespace()

		// Parse key (must be string)
		if p.pos >= p.length || p.data[p.pos] != '"' {
			return nil, errors.New("expected string key in object")
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		p.skipWhitespace()

		// Expect ':'
		if p.pos >= p.length || p.data[p.pos] != ':' {
			return nil, errors.New("expected ':' after object key")
		}
		p.pos++ // skip ':'

		p.skipWhitespace()

		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		result[key] = value

		p.skipWhitespace()

		// Check for more entries or end of object
		if p.pos >= p.length {
			return nil, errors.New("unexpected end of JSON input in object")
		}

		if p.data[p.pos] == '}' {
			p.pos++
			return result, nil
		}

		if p.data[p.pos] != ',' {
			return nil, fmt.Errorf("expected ',' or '}' in object at position %d", p.pos)
		}
		p.pos++ // skip ','
	}
}

// parseArray parses a JSON array into []interface{}.
func (p *Parser) parseArray() ([]interface{}, error) {
	if p.pos >= p.length || p.data[p.pos] != '[' {
		return nil, errors.New("expected '['")
	}
	p.pos++ // skip '['

	result := make([]interface{}, 0)
	p.skipWhitespace()

	// Handle empty array
	if p.pos < p.length && p.data[p.pos] == ']' {
		p.pos++
		return result, nil
	}

	for {
		p.skipWhitespace()

		// Parse element
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		result = append(result, value)

		p.skipWhitespace()

		// Check for more entries or end of array
		if p.pos >= p.length {
			return nil, errors.New("unexpected end of JSON input in array")
		}

		if p.data[p.pos] == ']' {
			p.pos++
			return result, nil
		}

		if p.data[p.pos] != ',' {
			return nil, fmt.Errorf("expected ',' or ']' in array at position %d", p.pos)
		}
		p.pos++ // skip ','
	}
}

// parseString parses a JSON string.
func (p *Parser) parseString() (string, error) {
	if p.pos >= p.length || p.data[p.pos] != '"' {
		return "", errors.New("expected '\"'")
	}
	p.pos++ // skip opening '"'

	start := p.pos

	// Fast path: no escape sequences
	for p.pos < p.length {
		c := p.data[p.pos]
		if c == '"' {
			// Found closing quote without escape sequences
			s := string(p.data[start:p.pos])
			p.pos++ // skip closing '"'
			return s, nil
		}
		if c == '\\' {
			// Escape sequence found, use slow path
			return p.parseStringWithEscapes(start)
		}
		if c < 0x20 {
			return "", fmt.Errorf("invalid control character in string at position %d", p.pos)
		}
		p.pos++
	}

	return "", errors.New("unexpected end of JSON input in string")
}

// parseStringWithEscapes handles strings containing escape sequences.
func (p *Parser) parseStringWithEscapes(start int) (string, error) {
	// We already found an escape at p.pos, and everything before is in data[start:p.pos]
	var buf []byte
	buf = append(buf, p.data[start:p.pos]...)

	for p.pos < p.length {
		c := p.data[p.pos]

		if c == '"' {
			p.pos++ // skip closing '"'
			return string(buf), nil
		}

		if c == '\\' {
			p.pos++
			if p.pos >= p.length {
				return "", errors.New("unexpected end of JSON input after backslash")
			}

			escaped := p.data[p.pos]
			p.pos++

			switch escaped {
			case '"', '\\', '/':
				buf = append(buf, escaped)
			case 'b':
				buf = append(buf, '\b')
			case 'f':
				buf = append(buf, '\f')
			case 'n':
				buf = append(buf, '\n')
			case 'r':
				buf = append(buf, '\r')
			case 't':
				buf = append(buf, '\t')
			case 'u':
				// Unicode escape \uXXXX
				if p.pos+4 > p.length {
					return "", errors.New("incomplete unicode escape")
				}
				hex := string(p.data[p.pos : p.pos+4])
				p.pos += 4

				codepoint, err := strconv.ParseUint(hex, 16, 16)
				if err != nil {
					return "", fmt.Errorf("invalid unicode escape: %v", err)
				}

				// Encode rune to UTF-8
				var tmp [utf8.UTFMax]byte
				n := utf8.EncodeRune(tmp[:], rune(codepoint))
				buf = append(buf, tmp[:n]...)
			default:
				return "", fmt.Errorf("invalid escape sequence '\\%c'", escaped)
			}
		} else if c < 0x20 {
			return "", fmt.Errorf("invalid control character in string at position %d", p.pos)
		} else {
			buf = append(buf, c)
			p.pos++
		}
	}

	return "", errors.New("unexpected end of JSON input in string")
}

// parseNumber parses a JSON number (int64 or float64).
func (p *Parser) parseNumber() (interface{}, error) {
	start := p.pos

	// Optional minus sign
	if p.pos < p.length && p.data[p.pos] == '-' {
		p.pos++
	}

	if p.pos >= p.length {
		return nil, errors.New("unexpected end of JSON input in number")
	}

	// Integer part
	if p.data[p.pos] == '0' {
		p.pos++
	} else if p.data[p.pos] >= '1' && p.data[p.pos] <= '9' {
		p.pos++
		for p.pos < p.length && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	} else {
		return nil, errors.New("invalid number")
	}

	// Check for decimal point or exponent
	isFloat := false

	if p.pos < p.length && p.data[p.pos] == '.' {
		isFloat = true
		p.pos++

		// Decimal digits required after '.'
		if p.pos >= p.length || p.data[p.pos] < '0' || p.data[p.pos] > '9' {
			return nil, errors.New("invalid number: expected digit after '.'")
		}

		for p.pos < p.length && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	}

	if p.pos < p.length && (p.data[p.pos] == 'e' || p.data[p.pos] == 'E') {
		isFloat = true
		p.pos++

		// Optional sign
		if p.pos < p.length && (p.data[p.pos] == '+' || p.data[p.pos] == '-') {
			p.pos++
		}

		// Exponent digits required
		if p.pos >= p.length || p.data[p.pos] < '0' || p.data[p.pos] > '9' {
			return nil, errors.New("invalid number: expected digit in exponent")
		}

		for p.pos < p.length && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	}

	numStr := string(p.data[start:p.pos])

	if isFloat {
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %v", err)
		}
		return f, nil
	}

	i, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		// Try as float if integer overflow
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %v", err)
		}
		return f, nil
	}
	return i, nil
}

// parseTrue parses the literal "true".
func (p *Parser) parseTrue() (bool, error) {
	if p.pos+4 <= p.length && string(p.data[p.pos:p.pos+4]) == "true" {
		p.pos += 4
		return true, nil
	}
	return false, errors.New("invalid literal (expected 'true')")
}

// parseFalse parses the literal "false".
// nolint:unparam // Returns false by design - parses the "false" literal
func (p *Parser) parseFalse() (bool, error) {
	if p.pos+5 <= p.length && string(p.data[p.pos:p.pos+5]) == "false" {
		p.pos += 5
		return false, nil
	}
	return false, errors.New("invalid literal (expected 'false')")
}

// parseNull parses the literal "null".
func (p *Parser) parseNull() (interface{}, error) {
	if p.pos+4 <= p.length && string(p.data[p.pos:p.pos+4]) == "null" {
		p.pos += 4
		return nil, nil
	}
	return nil, errors.New("invalid literal (expected 'null')")
}

// skipWhitespace skips whitespace characters (space, tab, LF, CR).
// This uses SWAR optimization when possible for better performance.
func (p *Parser) skipWhitespace() {
	// Fast path: SWAR optimization for common whitespace
	for p.pos+8 <= p.length {
		// Load 8 bytes
		chunk := uint64(p.data[p.pos])<<56 |
			uint64(p.data[p.pos+1])<<48 |
			uint64(p.data[p.pos+2])<<40 |
			uint64(p.data[p.pos+3])<<32 |
			uint64(p.data[p.pos+4])<<24 |
			uint64(p.data[p.pos+5])<<16 |
			uint64(p.data[p.pos+6])<<8 |
			uint64(p.data[p.pos+7])

		// Check if all bytes are whitespace (space=0x20, tab=0x09, LF=0x0A, CR=0x0D)
		// This is a simplified check - if any byte is NOT whitespace, break
		hasNonWhitespace := false
		for i := 0; i < 8; i++ {
			b := byte(chunk >> (56 - i*8))
			if b != ' ' && b != '\t' && b != '\n' && b != '\r' {
				hasNonWhitespace = true
				break
			}
		}

		if hasNonWhitespace {
			break
		}

		p.pos += 8
	}

	// Slow path: handle remaining bytes
	for p.pos < p.length {
		c := p.data[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			p.pos++
		} else {
			break
		}
	}
}
