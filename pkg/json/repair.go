package json

import (
	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-json/internal/lenient"
	"github.com/shapestone/shape-json/internal/parser"
)

// CorrectionKind identifies the type of correction applied during Repair.
type CorrectionKind = lenient.CorrectionKind

const (
	CorrectionTrailingComma  = lenient.CorrectionTrailingComma
	CorrectionSingleQuote    = lenient.CorrectionSingleQuote
	CorrectionUnquotedKey    = lenient.CorrectionUnquotedKey
	CorrectionLineComment    = lenient.CorrectionLineComment
	CorrectionBlockComment   = lenient.CorrectionBlockComment
	CorrectionUnescapedQuote = lenient.CorrectionUnescapedQuote
	CorrectionDuplicateKey   = lenient.CorrectionDuplicateKey
)

// Correction records a single auto-correction applied during Repair.
type Correction struct {
	Kind     CorrectionKind
	Position ast.Position
	Original string
	Fixed    string
	Message  string
}

// Repair takes potentially invalid JSON and returns valid RFC 8259 JSON.
//
// It auto-corrects common errors:
//   - Trailing commas in objects and arrays
//   - Single-quoted strings (converted to double-quoted)
//   - Unquoted object keys (quoted)
//   - Line comments (//) and block comments (/* */) (removed)
//   - Unescaped double quotes inside strings (escaped)
//   - Duplicate keys (last value wins)
//
// If the input is already valid JSON, it is parsed and re-rendered (keys are
// sorted alphabetically, whitespace is normalized).
//
// Returns an error only if the input is too broken to recover.
//
// Example:
//
//	// Try fast path, repair on failure
//	err := json.Unmarshal(data, &v)
//	if err != nil {
//	    fixed, repairErr := json.RepairBytes(data)
//	    if repairErr != nil {
//	        return repairErr
//	    }
//	    err = json.Unmarshal(fixed, &v)
//	}
func Repair(input string) (string, error) {
	p := parser.NewLenientParser(input)
	node, _, err := p.Parse()
	if err != nil {
		return "", err
	}
	result, err := Render(node)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// RepairBytes is like Repair but operates on byte slices.
// Convenient for composing with Unmarshal.
func RepairBytes(data []byte) ([]byte, error) {
	p := parser.NewLenientParser(string(data))
	node, _, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return Render(node)
}

// RepairWithCorrections is like Repair but also returns what was fixed.
// Each Correction describes one auto-correction: its kind, position in the
// original input, and a human-readable message.
//
// Example:
//
//	fixed, corrections, err := json.RepairWithCorrections(rawInput)
//	for _, c := range corrections {
//	    log.Printf("Fixed %s at line %d: %s", c.Kind, c.Position.Row(), c.Message)
//	}
func RepairWithCorrections(input string) (string, []Correction, error) {
	p := parser.NewLenientParser(input)
	node, internalCorrections, err := p.Parse()
	if err != nil {
		return "", nil, err
	}
	result, err := Render(node)
	if err != nil {
		return "", nil, err
	}
	return string(result), toPublicCorrections(internalCorrections), nil
}

func toPublicCorrections(internal []lenient.Correction) []Correction {
	if len(internal) == 0 {
		return nil
	}
	out := make([]Correction, len(internal))
	for i, c := range internal {
		out[i] = Correction{
			Kind:     c.Kind,
			Position: c.Position,
			Original: c.Original,
			Fixed:    c.Fixed,
			Message:  c.Message,
		}
	}
	return out
}
