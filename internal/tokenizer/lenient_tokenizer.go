package tokenizer

import (
	"github.com/shapestone/shape-core/pkg/tokenizer"
)

// NewLenientTokenizer creates a tokenizer that handles lenient JSON.
// It extends the strict tokenizer with support for comments, single-quoted
// strings, and unquoted identifier keys.
//
// Matcher ordering:
// 1. Comments (before anything else to strip them early)
// 2. Keywords (true, false, null — before identifiers to avoid mismatching)
// 3. Structural tokens
// 4. String literals (double-quoted and single-quoted)
// 5. Numbers
// 6. Unquoted identifiers (last, as a catch-all for bare keys)
func NewLenientTokenizer() tokenizer.Tokenizer {
	return tokenizer.NewTokenizer(
		// Comments first
		CommentMatcher(),

		// Keywords (before identifier matching)
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

		// String literals (handles both " and ')
		LenientStringMatcher(),

		// Numbers
		NumberMatcher(),

		// Unquoted identifiers (catch-all, must be last)
		UnquotedKeyMatcher(),
	)
}

// NewLenientTokenizerWithStream creates a lenient tokenizer using a pre-configured stream.
func NewLenientTokenizerWithStream(stream tokenizer.Stream) tokenizer.Tokenizer {
	tok := NewLenientTokenizer()
	tok.InitializeFromStream(stream)
	return tok
}
