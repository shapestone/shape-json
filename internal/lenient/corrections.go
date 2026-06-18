package lenient

import (
	"fmt"

	"github.com/shapestone/shape-core/pkg/ast"
)

type CorrectionKind int

const (
	CorrectionTrailingComma  CorrectionKind = iota
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

type Correction struct {
	Kind     CorrectionKind
	Position ast.Position
	Original string
	Fixed    string
	Message  string
}

type CorrectionCollector struct {
	corrections []Correction
}

func NewCorrectionCollector() *CorrectionCollector {
	return &CorrectionCollector{}
}

func (c *CorrectionCollector) Add(kind CorrectionKind, pos ast.Position, original, fixed, message string) {
	c.corrections = append(c.corrections, Correction{
		Kind:     kind,
		Position: pos,
		Original: original,
		Fixed:    fixed,
		Message:  message,
	})
}

func (c *CorrectionCollector) Corrections() []Correction {
	return c.corrections
}
