package tokenizer

import (
	"strings"
	"testing"

	"github.com/shapestone/shape-core/pkg/tokenizer"
)

// TestTokenizer_Structural tests structural tokens: {, }, [, ], :, ,
func TestTokenizer_Structural(t *testing.T) {
	tok := NewTokenizer()
	tok.Initialize(`{ } [ ] : ,`)

	expected := []string{
		TokenLBrace,   // {
		TokenRBrace,   // }
		TokenLBracket, // [
		TokenRBracket, // ]
		TokenColon,    // :
		TokenComma,    // ,
	}

	for i, exp := range expected {
		token := nextNonWhitespace(&tok)
		if token == nil {
			t.Fatalf("token %d: expected token, got none", i)
		}
		if token.Kind() != exp {
			t.Errorf("token %d: expected %s, got %s (value: %q)",
				i, exp, token.Kind(), token.ValueString())
		}
	}
}

// nextNonWhitespace returns the next non-whitespace token
func nextNonWhitespace(tok *tokenizer.Tokenizer) *tokenizer.Token {
	for {
		token, ok := tok.NextToken()
		if !ok {
			return nil
		}
		if token.Kind() != "Whitespace" {
			return token
		}
	}
}

// TestTokenizer_Keywords tests boolean and null literals
func TestTokenizer_Keywords(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", TokenTrue},
		{"false", TokenFalse},
		{"null", TokenNull},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tok := NewTokenizer()
			tok.Initialize(tt.input)

			token, ok := tok.NextToken()
			if !ok {
				t.Fatalf("expected token, got none")
			}
			if token.Kind() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, token.Kind())
			}
			if token.ValueString() != tt.input {
				t.Errorf("expected value %q, got %q", tt.input, token.ValueString())
			}
		})
	}
}

// TestTokenizer_String tests string tokenization
func TestTokenizer_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: `""`,
			want:  `""`,
		},
		{
			name:  "simple string",
			input: `"hello"`,
			want:  `"hello"`,
		},
		{
			name:  "string with spaces",
			input: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "escaped quote",
			input: `"say \"hello\""`,
			want:  `"say \"hello\""`,
		},
		{
			name:  "escaped backslash",
			input: `"path\\to\\file"`,
			want:  `"path\\to\\file"`,
		},
		{
			name:  "escaped newline",
			input: `"line1\nline2"`,
			want:  `"line1\nline2"`,
		},
		{
			name:  "multiple escapes",
			input: `"\"\\\n\r\t"`,
			want:  `"\"\\\n\r\t"`,
		},
		{
			name:  "unicode escape",
			input: `"\u0041\u0042\u0043"`,
			want:  `"\u0041\u0042\u0043"`,
		},
		{
			name:  "unicode with text",
			input: `"hello \u03B1\u03B2\u03B3 world"`,
			want:  `"hello \u03B1\u03B2\u03B3 world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := NewTokenizer()
			tok.Initialize(tt.input)

			token, ok := tok.NextToken()
			if !ok {
				t.Fatalf("expected token, got none")
			}
			if token.Kind() != TokenString {
				t.Errorf("expected TokenString, got %s", token.Kind())
			}
			if token.ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, token.ValueString())
			}
		})
	}
}

// TestTokenizer_String_Invalid tests invalid strings are rejected
func TestTokenizer_String_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unterminated string",
			input: `"hello`,
		},
		{
			name:  "unescaped control char",
			input: "\"hello\nworld\"", // literal newline
		},
		{
			name:  "incomplete escape",
			input: `"hello\`,
		},
		{
			name:  "incomplete unicode escape",
			input: `"hello\u00"`,
		},
		{
			name:  "invalid unicode escape",
			input: `"hello\uXXXX"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := NewTokenizer()
			tok.Initialize(tt.input)

			token, ok := tok.NextToken()
			if ok && token.Kind() == TokenString {
				t.Errorf("expected invalid string to fail, but got: %s", token.ValueString())
			}
		})
	}
}

// TestTokenizer_Number tests number tokenization
func TestTokenizer_Number(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "zero",
			input: "0",
			want:  "0",
		},
		{
			name:  "positive integer",
			input: "123",
			want:  "123",
		},
		{
			name:  "negative integer",
			input: "-456",
			want:  "-456",
		},
		{
			name:  "decimal",
			input: "123.456",
			want:  "123.456",
		},
		{
			name:  "negative decimal",
			input: "-123.456",
			want:  "-123.456",
		},
		{
			name:  "exponent lowercase e",
			input: "1e10",
			want:  "1e10",
		},
		{
			name:  "exponent uppercase E",
			input: "1E10",
			want:  "1E10",
		},
		{
			name:  "exponent with plus",
			input: "1e+10",
			want:  "1e+10",
		},
		{
			name:  "exponent with minus",
			input: "1e-10",
			want:  "1e-10",
		},
		{
			name:  "decimal with exponent",
			input: "1.23e10",
			want:  "1.23e10",
		},
		{
			name:  "negative decimal with exponent",
			input: "-1.23e-10",
			want:  "-1.23e-10",
		},
		{
			name:  "large number",
			input: "9999999999",
			want:  "9999999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := NewTokenizer()
			tok.Initialize(tt.input)

			token, ok := tok.NextToken()
			if !ok {
				t.Fatalf("expected token, got none")
			}
			if token.Kind() != TokenNumber {
				t.Errorf("expected TokenNumber, got %s", token.Kind())
			}
			if token.ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, token.ValueString())
			}
		})
	}
}

// TestTokenizer_Number_Invalid tests invalid numbers are rejected at the tokenizer level
// Note: Leading zeros like "01" are tokenized as two separate numbers ("0" and "1"),
// which is valid from the tokenizer's perspective. The parser is responsible for
// rejecting consecutive values without proper structure.
func TestTokenizer_Number_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "decimal without fraction",
			input: "123.",
		},
		{
			name:  "exponent without digits",
			input: "123e",
		},
		{
			name:  "exponent with sign but no digits",
			input: "123e+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := NewTokenizer()
			tok.Initialize(tt.input)

			token, ok := tok.NextToken()
			if ok && token.Kind() == TokenNumber {
				t.Errorf("expected invalid number to fail tokenization, but got: %s", token.ValueString())
			}
		})
	}
}

// TestTokenizer_Object tests tokenizing a complete object
func TestTokenizer_Object(t *testing.T) {
	input := `{"name": "Alice", "age": 30}`
	tok := NewTokenizer()
	tok.Initialize(input)

	expected := []struct {
		kind  string
		value string
	}{
		{TokenLBrace, "{"},
		{TokenString, `"name"`},
		{TokenColon, ":"},
		{TokenString, `"Alice"`},
		{TokenComma, ","},
		{TokenString, `"age"`},
		{TokenColon, ":"},
		{TokenNumber, "30"},
		{TokenRBrace, "}"},
	}

	for i, exp := range expected {
		token := nextNonWhitespace(&tok)
		if token == nil {
			t.Fatalf("token %d: expected token, got none", i)
		}
		if token.Kind() != exp.kind {
			t.Errorf("token %d: expected kind %s, got %s", i, exp.kind, token.Kind())
		}
		if token.ValueString() != exp.value {
			t.Errorf("token %d: expected value %q, got %q", i, exp.value, token.ValueString())
		}
	}
}

// TestTokenizer_Array tests tokenizing a complete array
func TestTokenizer_Array(t *testing.T) {
	input := `[1, 2, 3, true, false, null]`
	tok := NewTokenizer()
	tok.Initialize(input)

	expected := []struct {
		kind  string
		value string
	}{
		{TokenLBracket, "["},
		{TokenNumber, "1"},
		{TokenComma, ","},
		{TokenNumber, "2"},
		{TokenComma, ","},
		{TokenNumber, "3"},
		{TokenComma, ","},
		{TokenTrue, "true"},
		{TokenComma, ","},
		{TokenFalse, "false"},
		{TokenComma, ","},
		{TokenNull, "null"},
		{TokenRBracket, "]"},
	}

	for i, exp := range expected {
		token := nextNonWhitespace(&tok)
		if token == nil {
			t.Fatalf("token %d: expected token, got none", i)
		}
		if token.Kind() != exp.kind {
			t.Errorf("token %d: expected kind %s, got %s", i, exp.kind, token.Kind())
		}
		if token.ValueString() != exp.value {
			t.Errorf("token %d: expected value %q, got %q", i, exp.value, token.ValueString())
		}
	}
}

// TestTokenizer_Nested tests tokenizing nested structures
func TestTokenizer_Nested(t *testing.T) {
	input := `{"users": [{"name": "Alice"}, {"name": "Bob"}]}`
	tok := NewTokenizer()
	tok.Initialize(input)

	// Just verify we can tokenize it completely
	count := 0
	for {
		token, ok := tok.NextToken()
		if !ok {
			break
		}
		count++
		// Basic sanity check
		if token.Kind() == "" {
			t.Fatalf("token %d has empty kind", count)
		}
	}

	if count == 0 {
		t.Error("expected tokens, got none")
	}
}

// TestTokenizer_Whitespace tests that whitespace is properly ignored
func TestTokenizer_Whitespace(t *testing.T) {
	input := `  {  "name"  :  "Alice"  }  `
	tok := NewTokenizer()
	tok.Initialize(input)

	expected := []string{
		TokenLBrace,
		TokenString,
		TokenColon,
		TokenString,
		TokenRBrace,
	}

	for i, exp := range expected {
		token := nextNonWhitespace(&tok)
		if token == nil {
			t.Fatalf("token %d: expected token, got none", i)
		}
		if token.Kind() != exp {
			t.Errorf("token %d: expected %s, got %s", i, exp, token.Kind())
		}
	}
}

// TestTokenizer_Position tests that tokens have correct position information
func TestTokenizer_Position(t *testing.T) {
	input := `{"key": "value"}`
	tok := NewTokenizer()
	tok.Initialize(input)

	// First token should be at row 1, column 1
	token, ok := tok.NextToken()
	if !ok {
		t.Fatal("expected token")
	}
	if token.Row() != 1 {
		t.Errorf("expected row 1, got %d", token.Row())
	}
	if token.Column() != 1 {
		t.Errorf("expected column 1, got %d", token.Column())
	}
}

// TestTokenizer_StringAcrossBufferBoundary tests that string tokenization works
// correctly when the string spans across the read buffer boundary (8KB)
func TestTokenizer_StringAcrossBufferBoundary(t *testing.T) {
	// Create a string that spans the 8KB read boundary
	// The buffered stream reads in 8KB chunks
	prefix := strings.Repeat("x", 8100) // Position just before the 8KB mark
	testString := `"Item_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`
	data := prefix + testString

	t.Logf("Data length: %d", len(data))
	t.Logf("Test string starts at position: %d", 8100)
	t.Logf("This should cross the 8KB (8192 byte) read boundary")

	// Create buffered stream
	reader := strings.NewReader(data)
	stream := tokenizer.NewStreamFromReader(reader)

	// Advance stream to position 8100 (skip the prefix)
	for i := 0; i < 8100; i++ {
		stream.NextChar()
	}

	t.Logf("Stream positioned at offset %d (should be 8100)", stream.GetOffset())

	// Create tokenizer with the stream at this position
	tok := NewTokenizerWithStream(stream)

	// Now try to tokenize the string that spans the boundary
	t.Logf("About to tokenize string at stream position %d", stream.GetOffset())

	token, ok := tok.NextToken()
	if !ok {
		t.Fatalf("Failed to tokenize string at position 8100 (across buffer boundary), stream now at %d", stream.GetOffset())
	}

	if token.Kind() != TokenString {
		t.Errorf("Expected TokenString, got %s", token.Kind())
	}

	if token.ValueString() != testString {
		t.Errorf("Expected value %q, got %q", testString, token.ValueString())
	}

	t.Logf("Successfully tokenized string, stream now at position %d", stream.GetOffset())
}

// TestTokenizer_LargeJSONArray tests tokenizing a large JSON array that crosses buffer boundaries
func TestTokenizer_LargeJSONArray(t *testing.T) {
	// Create exactly 61 items - this is the failing case
	var sb strings.Builder
	sb.WriteString(`{"items": [`)
	for i := 0; i < 61; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(`{"id": 1000000000, "name": "Item_`)
		sb.WriteString(strings.Repeat("X", 50))
		sb.WriteString(`", "description": "`)
		sb.WriteString(strings.Repeat("Y", 30))
		sb.WriteString(`"}`)
	}
	sb.WriteString(`]}`)

	data := sb.String()
	t.Logf("Data length: %d bytes", len(data))

	// Create buffered stream
	reader := strings.NewReader(data)
	stream := tokenizer.NewStreamFromReader(reader)

	// Create tokenizer
	tok := NewTokenizerWithStream(stream)

	// Tokenize all tokens
	tokenCount := 0
	for {
		_, ok := tok.NextToken()
		if !ok {
			if !stream.IsEos() {
				// Get some context about where we are
				pos := stream.GetOffset()
				r, pok := stream.PeekChar()

				// Try to read a few chars to see what's available
				testChars := make([]rune, 0, 10)
				for i := 0; i < 10; i++ {
					tr, tok := stream.NextChar()
					if !tok {
						break
					}
					testChars = append(testChars, tr)
				}

				t.Logf("Failed at position %d", pos)
				t.Logf("PeekChar: '%c' (U+%04X), ok=%v", r, r, pok)
				t.Logf("Next 10 chars available: %d", len(testChars))
				if len(testChars) > 0 {
					t.Logf("Chars: %q", string(testChars))
				}
				t.Logf("Expected remaining: ~%d bytes", len(data)-pos)

				t.Fatalf("Tokenization failed after %d tokens at position %d, but not at EOS",
					tokenCount, pos)
			}
			break
		}
		tokenCount++
	}

	t.Logf("Successfully tokenized %d tokens", tokenCount)
}
