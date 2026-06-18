package tokenizer_test

import (
	"testing"

	shapetokenizer "github.com/shapestone/shape-core/pkg/tokenizer"
	"github.com/shapestone/shape-json/internal/tokenizer"
)

func lenientTokenize(input string) []*shapetokenizer.Token {
	tok := tokenizer.NewLenientTokenizer()
	tok.InitializeFromStream(shapetokenizer.NewStream(input))
	var tokens []*shapetokenizer.Token
	for {
		t, ok := tok.NextToken()
		if !ok {
			break
		}
		if t.Kind() == "Whitespace" {
			continue
		}
		tokens = append(tokens, t)
	}
	return tokens
}

func TestCommentMatcher_LineComment(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "// hello", "// hello"},
		{"with content after newline", "// hello\n42", "// hello"},
		{"empty comment", "//\n42", "//"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) == 0 {
				t.Fatal("expected at least one token")
			}
			if tokens[0].Kind() != tokenizer.TokenComment {
				t.Errorf("expected Comment token, got %s", tokens[0].Kind())
			}
			if tokens[0].ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, tokens[0].ValueString())
			}
		})
	}
}

func TestCommentMatcher_BlockComment(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "/* hello */", "/* hello */"},
		{"multiline", "/* line1\nline2 */", "/* line1\nline2 */"},
		{"with content after", "/* hi */ 42", "/* hi */"},
		{"empty", "/**/", "/**/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) == 0 {
				t.Fatal("expected at least one token")
			}
			if tokens[0].Kind() != tokenizer.TokenComment {
				t.Errorf("expected Comment token, got %s", tokens[0].Kind())
			}
			if tokens[0].ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, tokens[0].ValueString())
			}
		})
	}
}

func TestLenientStringMatcher_SingleQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "'hello'", "'hello'"},
		{"with internal double quotes", `'say "hi"'`, `'say "hi"'`},
		{"escaped single quote", `'it\'s'`, `'it\'s'`},
		{"empty", "''", "''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) == 0 {
				t.Fatal("expected at least one token")
			}
			if tokens[0].Kind() != tokenizer.TokenSingleString {
				t.Errorf("expected SingleString token, got %s", tokens[0].Kind())
			}
			if tokens[0].ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, tokens[0].ValueString())
			}
		})
	}
}

func TestLenientStringMatcher_DoubleQuotes(t *testing.T) {
	tokens := lenientTokenize(`"hello"`)
	if len(tokens) == 0 {
		t.Fatal("expected at least one token")
	}
	if tokens[0].Kind() != tokenizer.TokenString {
		t.Errorf("expected String token, got %s", tokens[0].Kind())
	}
	if tokens[0].ValueString() != `"hello"` {
		t.Errorf("expected %q, got %q", `"hello"`, tokens[0].ValueString())
	}
}

func TestUnquotedKeyMatcher(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "name", "name"},
		{"with underscore", "_id", "_id"},
		{"with dollar", "$ref", "$ref"},
		{"alphanumeric", "key123", "key123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) == 0 {
				t.Fatal("expected at least one token")
			}
			if tokens[0].Kind() != tokenizer.TokenIdentifier {
				t.Errorf("expected Identifier token, got %s", tokens[0].Kind())
			}
			if tokens[0].ValueString() != tt.want {
				t.Errorf("expected value %q, got %q", tt.want, tokens[0].ValueString())
			}
		})
	}
}

func TestUnquotedKeyMatcher_NotMatched(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"digit start", "123abc"},
		{"punctuation", ".abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) > 0 && tokens[0].Kind() == tokenizer.TokenIdentifier {
				t.Errorf("should not match as Identifier, got %q", tokens[0].ValueString())
			}
		})
	}
}

func TestLenientTokenizer_KeywordsBeforeIdentifiers(t *testing.T) {
	// true/false/null should match as keywords, not identifiers
	tests := []struct {
		input string
		kind  string
	}{
		{"true", tokenizer.TokenTrue},
		{"false", tokenizer.TokenFalse},
		{"null", tokenizer.TokenNull},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := lenientTokenize(tt.input)
			if len(tokens) == 0 {
				t.Fatal("expected at least one token")
			}
			if tokens[0].Kind() != tt.kind {
				t.Errorf("expected %s token, got %s", tt.kind, tokens[0].Kind())
			}
		})
	}
}

func TestLenientTokenizer_FullObject(t *testing.T) {
	input := `{name: 'Alice', age: 30}`
	tokens := lenientTokenize(input)

	expected := []struct {
		kind  string
		value string
	}{
		{tokenizer.TokenLBrace, "{"},
		{tokenizer.TokenIdentifier, "name"},
		{tokenizer.TokenColon, ":"},
		{tokenizer.TokenSingleString, `'Alice'`},
		{tokenizer.TokenComma, ","},
		{tokenizer.TokenIdentifier, "age"},
		{tokenizer.TokenColon, ":"},
		{tokenizer.TokenNumber, "30"},
		{tokenizer.TokenRBrace, "}"},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Kind() != exp.kind {
			t.Errorf("token[%d]: expected kind %s, got %s", i, exp.kind, tokens[i].Kind())
		}
		if tokens[i].ValueString() != exp.value {
			t.Errorf("token[%d]: expected value %q, got %q", i, exp.value, tokens[i].ValueString())
		}
	}
}
