// Package tokenizer provides JSON tokenization using Shape's tokenizer framework.
package tokenizer

// Token type constants for JSON format.
// These correspond to the terminals in the JSON grammar (docs/grammar/json.ebnf).
const (
	// Structural tokens
	TokenLBrace   = "LBrace"   // {
	TokenRBrace   = "RBrace"   // }
	TokenLBracket = "LBracket" // [
	TokenRBracket = "RBracket" // ]
	TokenColon    = "Colon"    // :
	TokenComma    = "Comma"    // ,

	// Literal value tokens
	TokenString = "String" // "..." (quoted string with possible escapes)
	TokenNumber = "Number" // 123, -45.67, 1.23e10
	TokenTrue   = "True"   // true
	TokenFalse  = "False"  // false
	TokenNull   = "Null"   // null

	// Special token
	TokenEOF = "EOF" // End of file
)
