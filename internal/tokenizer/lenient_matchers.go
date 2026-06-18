package tokenizer

import (
	"github.com/shapestone/shape-core/pkg/tokenizer"
)

// CommentMatcher creates a matcher for JSON comments (not part of RFC 8259).
// Matches line comments (// to end of line) and block comments (/* to */).
func CommentMatcher() tokenizer.Matcher {
	return func(stream tokenizer.Stream) *tokenizer.Token {
		r, ok := stream.PeekChar()
		if !ok || r != '/' {
			return nil
		}

		stream.NextChar() // consume first '/'

		r, ok = stream.PeekChar()
		if !ok {
			return nil
		}

		if r == '/' {
			return matchLineComment(stream)
		}
		if r == '*' {
			return matchBlockComment(stream)
		}

		return nil
	}
}

func matchLineComment(stream tokenizer.Stream) *tokenizer.Token {
	stream.NextChar() // consume second '/'
	var value []rune
	value = append(value, '/', '/')

	for {
		r, ok := stream.PeekChar()
		if !ok {
			break
		}
		if r == '\n' {
			break
		}
		stream.NextChar()
		value = append(value, r)
	}

	return tokenizer.NewToken(TokenComment, value)
}

func matchBlockComment(stream tokenizer.Stream) *tokenizer.Token {
	stream.NextChar() // consume '*'
	var value []rune
	value = append(value, '/', '*')

	for {
		r, ok := stream.NextChar()
		if !ok {
			// Unterminated block comment — return what we have
			return tokenizer.NewToken(TokenComment, value)
		}
		value = append(value, r)

		if r == '*' {
			r2, ok := stream.PeekChar()
			if ok && r2 == '/' {
				stream.NextChar()
				value = append(value, '/')
				return tokenizer.NewToken(TokenComment, value)
			}
		}
	}
}

// LenientStringMatcher creates a matcher that handles both double-quoted and
// single-quoted JSON strings. Single-quoted strings are normalized to
// double-quoted form in the token value.
func LenientStringMatcher() tokenizer.Matcher {
	return func(stream tokenizer.Stream) *tokenizer.Token {
		r, ok := stream.PeekChar()
		if !ok {
			return nil
		}

		if r == '"' {
			// Delegate to the standard string matcher for double-quoted strings
			if byteStream, ok := stream.(tokenizer.ByteStream); ok {
				return stringMatcherByte(byteStream)
			}
			return stringMatcherRune(stream)
		}

		if r == '\'' {
			return singleQuotedStringMatcher(stream)
		}

		return nil
	}
}

// singleQuotedStringMatcher parses a single-quoted string and emits a
// TokenSingleString with the original value preserved (including quotes).
// The parser is responsible for normalizing to double quotes.
func singleQuotedStringMatcher(stream tokenizer.Stream) *tokenizer.Token {
	stream.NextChar() // consume opening '

	var value []rune
	value = append(value, '\'')

	for {
		r, ok := stream.NextChar()
		if !ok {
			return nil
		}

		value = append(value, r)

		if r == '\'' {
			return tokenizer.NewToken(TokenSingleString, value)
		}

		if r == '\\' {
			escaped, ok := stream.NextChar()
			if !ok {
				return nil
			}
			value = append(value, escaped)
			continue
		}

		if r < 0x20 {
			return nil
		}
	}
}

// UnquotedKeyMatcher creates a matcher for unquoted JavaScript-style identifiers.
// Matches: [a-zA-Z_$][a-zA-Z0-9_$]*
// Emits TokenIdentifier. The parser converts these to quoted string keys.
func UnquotedKeyMatcher() tokenizer.Matcher {
	return func(stream tokenizer.Stream) *tokenizer.Token {
		r, ok := stream.PeekChar()
		if !ok {
			return nil
		}

		if !isIdentStart(r) {
			return nil
		}

		var value []rune
		stream.NextChar()
		value = append(value, r)

		for {
			r, ok := stream.PeekChar()
			if !ok {
				break
			}
			if !isIdentContinue(r) {
				break
			}
			stream.NextChar()
			value = append(value, r)
		}

		return tokenizer.NewToken(TokenIdentifier, value)
	}
}

func isIdentStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
}

func isIdentContinue(r rune) bool {
	return isIdentStart(r) || (r >= '0' && r <= '9')
}
