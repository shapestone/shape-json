package tokenizer

import (
	"github.com/shapestone/shape-core/pkg/tokenizer"
)

// NewTokenizer creates a tokenizer for JSON format.
// The tokenizer matches JSON tokens in order of specificity.
//
// Ordering is critical:
// 1. Whitespace (automatically handled by Shape's framework, now SWAR-optimized)
// 2. Keywords (true, false, null) before general strings
// 3. Structural tokens ({, }, [, ], :, ,)
// 4. String literals (quoted with escapes)
// 5. Numbers (integers and floats with optional exponents)
func NewTokenizer() tokenizer.Tokenizer {
	return tokenizer.NewTokenizer(
		// Keywords (before string matching)
		tokenizer.StringMatcherFunc(TokenTrue, "true"),
		tokenizer.StringMatcherFunc(TokenFalse, "false"),
		tokenizer.StringMatcherFunc(TokenNull, "null"),

		// Structural tokens
		tokenizer.StringMatcherFunc(TokenLBrace, "{"),
		tokenizer.StringMatcherFunc(TokenRBrace, "}"),
		tokenizer.StringMatcherFunc(TokenLBracket, "["),
		tokenizer.StringMatcherFunc(TokenRBracket, "]"),
		tokenizer.StringMatcherFunc(TokenColon, ":"),
		tokenizer.StringMatcherFunc(TokenComma, ","),

		// String literals (complex pattern)
		StringMatcher(),

		// Numbers (complex pattern)
		NumberMatcher(),
	)
}

// NewTokenizerWithStream creates a tokenizer for JSON format using a pre-configured stream.
// This is used internally to support streaming from io.Reader.
func NewTokenizerWithStream(stream tokenizer.Stream) tokenizer.Tokenizer {
	tok := NewTokenizer()
	tok.InitializeFromStream(stream)
	return tok
}

// StringMatcher creates a matcher for JSON string literals.
// Matches: "..." with escape sequences \", \\, \/, \b, \f, \n, \r, \t, \uXXXX
//
// Grammar:
//
//	String = '"' { Character } '"' ;
//	Character = UnescapedChar | EscapeSequence ;
//	EscapeSequence = "\\" ( '"' | "\\" | "/" | "b" | "f" | "n" | "r" | "t" | UnicodeEscape ) ;
//	UnicodeEscape = "u" HexDigit HexDigit HexDigit HexDigit ;
//
// Performance: Uses ByteStream for fast ASCII scanning with SWAR acceleration.
func StringMatcher() tokenizer.Matcher {
	return func(stream tokenizer.Stream) *tokenizer.Token {
		// Try ByteStream fast path for ASCII strings
		if byteStream, ok := stream.(tokenizer.ByteStream); ok {
			return stringMatcherByte(byteStream)
		}

		// Fallback to rune-based matcher for non-ByteStream
		return stringMatcherRune(stream)
	}
}

// stringMatcherByte uses ByteStream + SWAR for optimal performance.
func stringMatcherByte(stream tokenizer.ByteStream) *tokenizer.Token {
	// Opening quote
	b, ok := stream.PeekByte()
	if !ok {
		return nil
	}
	if b != '"' {
		return nil
	}

	startPos := stream.BytePosition()
	stream.NextByte() // consume opening quote

	// Use SWAR to find closing quote or escape quickly
	for {
		// Find next quote or backslash using SWAR
		offset := tokenizer.FindEscapeOrQuote(stream.RemainingBytes())

		if offset == -1 {
			// No quote or escape found - unterminated string
			return nil
		}

		// Advance to the found position
		for i := 0; i < offset; i++ {
			b, ok := stream.NextByte()
			if !ok {
				return nil
			}
			// Check for control characters
			if b < 0x20 {
				return nil
			}
		}

		// Now at quote or backslash
		b, ok := stream.NextByte()
		if !ok {
			return nil
		}

		if b == '"' {
			// Found closing quote - extract string
			value := stream.SliceFrom(startPos)

			// Convert to runes for token (maintains compatibility)
			return tokenizer.NewToken(TokenString, []rune(string(value)))
		}

		if b == '\\' {
			// Escape sequence - consume next character
			escaped, ok := stream.NextByte()
			if !ok {
				return nil
			}

			// If \u, consume 4 hex digits
			if escaped == 'u' {
				for i := 0; i < 4; i++ {
					hex, ok := stream.NextByte()
					if !ok {
						return nil
					}
					if !isHexDigitByte(hex) {
						return nil
					}
				}
			}
		}
	}
}

// stringMatcherRune is the fallback rune-based implementation.
func stringMatcherRune(stream tokenizer.Stream) *tokenizer.Token {
	var value []rune

	// Opening quote
	r, ok := stream.NextChar()
	if !ok || r != '"' {
		return nil
	}
	value = append(value, r)

	// Characters until closing quote
	for {
		r, ok := stream.NextChar()
		if !ok {
			return nil
		}

		value = append(value, r)

		if r == '"' {
			return tokenizer.NewToken(TokenString, value)
		}

		if r == '\\' {
			r, ok := stream.NextChar()
			if !ok {
				return nil
			}
			value = append(value, r)

			if r == 'u' {
				for i := 0; i < 4; i++ {
					r, ok := stream.NextChar()
					if !ok {
						return nil
					}
					if !isHexDigit(r) {
						return nil
					}
					value = append(value, r)
				}
			}
		} else if r < 0x20 {
			return nil
		}
	}
}

// isHexDigitByte checks if a byte is a hex digit (0-9, a-f, A-F).
func isHexDigitByte(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

// NumberMatcher creates a matcher for JSON number literals.
// Matches: integers and floats with optional sign and exponent
//
// Grammar:
//
//	Number = [ "-" ] Integer [ Fraction ] [ Exponent ] ;
//	Integer = "0" | ( [1-9] { Digit } ) ;
//	Fraction = "." Digit+ ;
//	Exponent = ( "e" | "E" ) [ "+" | "-" ] Digit+ ;
//
// Examples: 0, -123, 123.456, 1e10, 1.5e-3
// Performance: Uses ByteStream for fast ASCII number scanning.
func NumberMatcher() tokenizer.Matcher {
	return func(stream tokenizer.Stream) *tokenizer.Token {
		// Try ByteStream fast path for ASCII numbers
		if byteStream, ok := stream.(tokenizer.ByteStream); ok {
			return numberMatcherByte(byteStream)
		}

		// Fallback to rune-based matcher
		return numberMatcherRune(stream)
	}
}

// numberMatcherByte uses ByteStream for optimal number parsing.
func numberMatcherByte(stream tokenizer.ByteStream) *tokenizer.Token {
	startPos := stream.BytePosition()

	// Optional minus sign
	b, ok := stream.PeekByte()
	if ok && b == '-' {
		stream.NextByte()
	}

	// Integer part (required)
	b, ok = stream.PeekByte()
	if !ok || !isDigitByte(b) {
		return nil
	}

	// Special case: if it's '0', it must be alone (no leading zeros)
	if b == '0' {
		stream.NextByte()
	} else {
		// Digits 1-9 followed by more digits
		for {
			b, ok := stream.PeekByte()
			if !ok || !isDigitByte(b) {
				break
			}
			stream.NextByte()
		}
	}

	// Optional fraction
	b, ok = stream.PeekByte()
	if ok && b == '.' {
		stream.NextByte()

		// Must have at least one digit after decimal
		b, ok = stream.PeekByte()
		if !ok || !isDigitByte(b) {
			return nil
		}

		// Consume digits
		for {
			b, ok := stream.PeekByte()
			if !ok || !isDigitByte(b) {
				break
			}
			stream.NextByte()
		}
	}

	// Optional exponent
	b, ok = stream.PeekByte()
	if ok && (b == 'e' || b == 'E') {
		stream.NextByte()

		// Optional sign
		b, ok = stream.PeekByte()
		if ok && (b == '+' || b == '-') {
			stream.NextByte()
		}

		// Must have at least one digit
		b, ok = stream.PeekByte()
		if !ok || !isDigitByte(b) {
			return nil
		}

		// Consume digits
		for {
			b, ok := stream.PeekByte()
			if !ok || !isDigitByte(b) {
				break
			}
			stream.NextByte()
		}
	}

	// Extract the number as bytes and convert to runes
	value := stream.SliceFrom(startPos)
	return tokenizer.NewToken(TokenNumber, []rune(string(value)))
}

// numberMatcherRune is the fallback rune-based number matcher.
func numberMatcherRune(stream tokenizer.Stream) *tokenizer.Token {
	var value []rune

	// Optional minus sign
	r, ok := stream.PeekChar()
	if ok && r == '-' {
		stream.NextChar()
		value = append(value, r)
	}

	// Integer part (required)
	r, ok = stream.PeekChar()
	if !ok || !isDigit(r) {
		return nil
	}

	// Special case: if it's '0', it must be alone (no leading zeros)
	if r == '0' {
		stream.NextChar()
		value = append(value, r)
	} else {
		// Digits 1-9 followed by more digits
		for {
			r, ok := stream.PeekChar()
			if !ok || !isDigit(r) {
				break
			}
			stream.NextChar()
			value = append(value, r)
		}
	}

	// Optional fraction
	r, ok = stream.PeekChar()
	if ok && r == '.' {
		stream.NextChar()
		value = append(value, r)

		// Must have at least one digit after decimal
		r, ok = stream.PeekChar()
		if !ok || !isDigit(r) {
			return nil
		}

		// Consume digits
		for {
			r, ok := stream.PeekChar()
			if !ok || !isDigit(r) {
				break
			}
			stream.NextChar()
			value = append(value, r)
		}
	}

	// Optional exponent
	r, ok = stream.PeekChar()
	if ok && (r == 'e' || r == 'E') {
		stream.NextChar()
		value = append(value, r)

		// Optional sign
		r, ok = stream.PeekChar()
		if ok && (r == '+' || r == '-') {
			stream.NextChar()
			value = append(value, r)
		}

		// Must have at least one digit
		r, ok = stream.PeekChar()
		if !ok || !isDigit(r) {
			return nil
		}

		// Consume digits
		for {
			r, ok := stream.PeekChar()
			if !ok || !isDigit(r) {
				break
			}
			stream.NextChar()
			value = append(value, r)
		}
	}

	return tokenizer.NewToken(TokenNumber, value)
}

// isDigitByte checks if a byte is a decimal digit (0-9).
func isDigitByte(b byte) bool {
	return b >= '0' && b <= '9'
}

// isDigit returns true if r is a decimal digit (0-9).
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isHexDigit returns true if r is a hexadecimal digit (0-9, a-f, A-F).
func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
