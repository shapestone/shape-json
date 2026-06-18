package lenient

import (
	"fmt"

	"github.com/shapestone/shape-core/pkg/ast"
)

// CorrectionKind identifies the type of auto-correction applied during lenient parsing.
type CorrectionKind int

// Correction kinds for each type of JSON error that can be auto-corrected.
const (
	CorrectionTrailingComma CorrectionKind = iota
	CorrectionSingleQuote
	CorrectionUnquotedKey
	CorrectionLineComment
	CorrectionBlockComment
	CorrectionUnescapedQuote
	CorrectionDuplicateKey
)

func (k CorrectionKind) String() string {
	switch k {
	case CorrectionTrailingComma:
		return "trailing_comma"
	case CorrectionSingleQuote:
		return "single_quote"
	case CorrectionUnquotedKey:
		return "unquoted_key"
	case CorrectionLineComment:
		return "line_comment"
	case CorrectionBlockComment:
		return "block_comment"
	case CorrectionUnescapedQuote:
		return "unescaped_quote"
	case CorrectionDuplicateKey:
		return "duplicate_key"
	default:
		return fmt.Sprintf("unknown(%d)", int(k))
	}
}

// Correction records a single auto-correction applied during lenient parsing.
type Correction struct {
	Kind     CorrectionKind
	Position ast.Position
	Original string
	Fixed    string
	Message  string
}

// CorrectionCollector accumulates corrections during lenient parsing.
type CorrectionCollector struct {
	corrections []Correction
}

// NewCorrectionCollector creates a new empty collector.
func NewCorrectionCollector() *CorrectionCollector {
	return &CorrectionCollector{}
}

// Add records a correction.
func (c *CorrectionCollector) Add(kind CorrectionKind, pos ast.Position, original, fixed, message string) {
	c.corrections = append(c.corrections, Correction{
		Kind:     kind,
		Position: pos,
		Original: original,
		Fixed:    fixed,
		Message:  message,
	})
}

// Corrections returns all recorded corrections.
func (c *CorrectionCollector) Corrections() []Correction {
	return c.corrections
}
