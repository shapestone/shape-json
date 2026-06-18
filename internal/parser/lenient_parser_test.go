package parser_test

import (
	"strings"
	"testing"

	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-json/internal/lenient"
	"github.com/shapestone/shape-json/internal/parser"
)

func parseLenient(t *testing.T, input string) (ast.SchemaNode, []lenient.Correction) {
	t.Helper()
	p := parser.NewLenientParser(input)
	node, corrections, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return node, corrections
}

func TestLenientParser_ValidJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty object", `{}`},
		{"simple object", `{"name": "Alice", "age": 30}`},
		{"empty array", `[]`},
		{"simple array", `[1, 2, 3]`},
		{"nested", `{"items": [1, 2], "meta": {"count": 2}}`},
		{"string value", `"hello"`},
		{"number", `42`},
		{"boolean", `true`},
		{"null", `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, corrections := parseLenient(t, tt.input)
			if len(corrections) != 0 {
				t.Errorf("expected no corrections for valid JSON, got %d", len(corrections))
			}
		})
	}
}

func TestLenientParser_TrailingCommaObject(t *testing.T) {
	node, corrections := parseLenient(t, `{"name": "Alice", "age": 30,}`)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 2 {
		t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
	}

	if len(corrections) != 1 {
		t.Fatalf("expected 1 correction, got %d", len(corrections))
	}
	if corrections[0].Kind != lenient.CorrectionTrailingComma {
		t.Errorf("expected TrailingComma correction, got %s", corrections[0].Kind)
	}
}

func TestLenientParser_TrailingCommaArray(t *testing.T) {
	node, corrections := parseLenient(t, `[1, 2, 3,]`)

	arr, ok := node.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected ArrayDataNode, got %T", node)
	}

	if len(arr.Elements()) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements()))
	}

	if len(corrections) != 1 {
		t.Fatalf("expected 1 correction, got %d", len(corrections))
	}
	if corrections[0].Kind != lenient.CorrectionTrailingComma {
		t.Errorf("expected TrailingComma correction, got %s", corrections[0].Kind)
	}
}

func TestLenientParser_SingleQuotedStrings(t *testing.T) {
	node, corrections := parseLenient(t, `{'name': 'Alice'}`)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	nameNode, ok := obj.GetProperty("name")
	if !ok {
		t.Fatal("property 'name' not found")
	}

	lit, ok := nameNode.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("expected LiteralNode, got %T", nameNode)
	}

	if lit.Value().(string) != "Alice" {
		t.Errorf("expected 'Alice', got %q", lit.Value())
	}

	// Should have 2 corrections: one for key, one for value
	singleQuoteCount := 0
	for _, c := range corrections {
		if c.Kind == lenient.CorrectionSingleQuote {
			singleQuoteCount++
		}
	}
	if singleQuoteCount != 2 {
		t.Errorf("expected 2 SingleQuote corrections, got %d", singleQuoteCount)
	}
}

func TestLenientParser_UnquotedKeys(t *testing.T) {
	node, corrections := parseLenient(t, `{name: "Alice", age: 30}`)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	if _, ok := obj.GetProperty("name"); !ok {
		t.Error("property 'name' not found")
	}
	if _, ok := obj.GetProperty("age"); !ok {
		t.Error("property 'age' not found")
	}

	unquotedCount := 0
	for _, c := range corrections {
		if c.Kind == lenient.CorrectionUnquotedKey {
			unquotedCount++
		}
	}
	if unquotedCount != 2 {
		t.Errorf("expected 2 UnquotedKey corrections, got %d", unquotedCount)
	}
}

func TestLenientParser_LineComment(t *testing.T) {
	input := `{
		// this is a comment
		"name": "Alice"
	}`
	node, corrections := parseLenient(t, input)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	if _, ok := obj.GetProperty("name"); !ok {
		t.Error("property 'name' not found")
	}

	commentCount := 0
	for _, c := range corrections {
		if c.Kind == lenient.CorrectionLineComment {
			commentCount++
		}
	}
	if commentCount != 1 {
		t.Errorf("expected 1 LineComment correction, got %d", commentCount)
	}
}

func TestLenientParser_BlockComment(t *testing.T) {
	input := `{
		"name": "Alice" /* inline comment */,
		"age": 30
	}`
	node, corrections := parseLenient(t, input)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 2 {
		t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
	}

	blockCount := 0
	for _, c := range corrections {
		if c.Kind == lenient.CorrectionBlockComment {
			blockCount++
		}
	}
	if blockCount != 1 {
		t.Errorf("expected 1 BlockComment correction, got %d", blockCount)
	}
}

func TestLenientParser_DuplicateKeys(t *testing.T) {
	node, corrections := parseLenient(t, `{"key": "first", "key": "second"}`)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	keyNode, ok := obj.GetProperty("key")
	if !ok {
		t.Fatal("property 'key' not found")
	}

	lit := keyNode.(*ast.LiteralNode)
	if lit.Value().(string) != "second" {
		t.Errorf("expected last value 'second', got %q", lit.Value())
	}

	dupCount := 0
	for _, c := range corrections {
		if c.Kind == lenient.CorrectionDuplicateKey {
			dupCount++
		}
	}
	if dupCount != 1 {
		t.Errorf("expected 1 DuplicateKey correction, got %d", dupCount)
	}
}

func TestLenientParser_MixedCorrections(t *testing.T) {
	input := `{
		// settings
		name: 'Alice',
		age: 30,
	}`
	node, corrections := parseLenient(t, input)

	obj, ok := node.(*ast.ObjectNode)
	if !ok {
		t.Fatalf("expected ObjectNode, got %T", node)
	}

	if len(obj.Properties()) != 2 {
		t.Errorf("expected 2 properties, got %d", len(obj.Properties()))
	}

	kinds := make(map[lenient.CorrectionKind]int)
	for _, c := range corrections {
		kinds[c.Kind]++
	}

	if kinds[lenient.CorrectionLineComment] != 1 {
		t.Errorf("expected 1 LineComment, got %d", kinds[lenient.CorrectionLineComment])
	}
	if kinds[lenient.CorrectionUnquotedKey] != 2 {
		t.Errorf("expected 2 UnquotedKey, got %d", kinds[lenient.CorrectionUnquotedKey])
	}
	if kinds[lenient.CorrectionSingleQuote] != 1 {
		t.Errorf("expected 1 SingleQuote, got %d", kinds[lenient.CorrectionSingleQuote])
	}
	if kinds[lenient.CorrectionTrailingComma] != 1 {
		t.Errorf("expected 1 TrailingComma, got %d", kinds[lenient.CorrectionTrailingComma])
	}
}

func TestLenientParser_UnescapedQuotes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantKey string
		wantVal string
	}{
		{
			name:    "simple unescaped",
			input:   `{"msg": "He said "hello" to me"}`,
			wantKey: "msg",
			wantVal: `He said "hello" to me`,
		},
		{
			name:    "multiple unescaped",
			input:   `{"msg": "She said "hi" and he said "bye" to them"}`,
			wantKey: "msg",
			wantVal: `She said "hi" and he said "bye" to them`,
		},
		{
			name:    "unescaped at start",
			input:   `{"msg": ""hello" world"}`,
			wantKey: "msg",
			wantVal: `"hello" world`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, corrections := parseLenient(t, tt.input)

			obj, ok := node.(*ast.ObjectNode)
			if !ok {
				t.Fatalf("expected ObjectNode, got %T", node)
			}

			valNode, ok := obj.GetProperty(tt.wantKey)
			if !ok {
				t.Fatalf("property %q not found", tt.wantKey)
			}

			lit, ok := valNode.(*ast.LiteralNode)
			if !ok {
				t.Fatalf("expected LiteralNode, got %T", valNode)
			}

			if lit.Value().(string) != tt.wantVal {
				t.Errorf("expected %q, got %q", tt.wantVal, lit.Value())
			}

			hasUnescaped := false
			for _, c := range corrections {
				if c.Kind == lenient.CorrectionUnescapedQuote {
					hasUnescaped = true
				}
			}
			if !hasUnescaped {
				t.Error("expected UnescapedQuote correction")
			}
		})
	}
}

func TestLenientParser_UnescapedQuotesInArray(t *testing.T) {
	input := `["He said "hello"", "normal"]`
	node, _ := parseLenient(t, input)

	arr, ok := node.(*ast.ArrayDataNode)
	if !ok {
		t.Fatalf("expected ArrayDataNode, got %T", node)
	}

	if len(arr.Elements()) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr.Elements()))
	}
}

func TestLenientParser_UnescapedQuotesNotFirstElement(t *testing.T) {
	// Regression: unescaped quote recovery must work when the defective
	// element is not the first in an array and strings contain newlines.
	tests := []struct {
		name  string
		input string
	}{
		{
			"second array element with unescaped quotes",
			`{"items":[{"name":"Clean"},{"name":"Broken","def":"room (e.g. "Deluxe King") scoped."}]}`,
		},
		{
			"newline in defective string",
			"{\"items\":[{\"ok\":\"clean\"},{\"def\":\"room (e.g. \"Deluxe King\") scoped,\nstays.\"}]}",
		},
		{
			"newline before unescaped quote",
			"{\"items\":[{\"ok\":\"clean\"},{\"def\":\"A named category of\naccommodation room (e.g. \"Deluxe King\") scoped.\"}]}",
		},
		{
			"full accommodation-fulfillment repro",
			"{\"items\":[{\"name\":\"Room Daily\",\"kind\":\"entity\",\"definition\":\"A per-day record for a physical room.\",\"evidence\":\"route.ts\",\"register\":\"business\"},{\"name\":\"Room Type\",\"kind\":\"entity\",\"definition\":\"A named category of accommodation room (e.g. \"Deluxe King\") scoped to a program, used to group inventory blocks and link room assignments to stays.\",\"evidence\":\"route.ts\",\"register\":\"business\"}]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, corrections := parseLenient(t, tt.input)
			hasUnescaped := false
			for _, c := range corrections {
				if c.Kind == lenient.CorrectionUnescapedQuote {
					hasUnescaped = true
				}
			}
			if !hasUnescaped {
				t.Error("expected UnescapedQuote correction")
			}
		})
	}
}

func TestLenientParser_UnescapedQuotesWithMultibyteUTF8(t *testing.T) {
	// Regression: Token.Offset() returns a rune index, not a byte offset.
	// Multi-byte characters (em-dash U+2014 = 3 bytes, 1 rune) in earlier
	// strings caused the unescaped-quote recovery to misindex into the input.
	tests := []struct {
		name  string
		input string
	}{
		{
			"em-dash in first element",
			"{\"items\":[{\"name\":\"A\",\"evidence\":\"route.ts — handler\",\"register\":\"b\"},{\"name\":\"B\",\"definition\":\"a room (e.g. \"Deluxe King\") scoped.\",\"register\":\"b\"}]}",
		},
		{
			"em-dash and backtick in first element",
			"{\"items\":[{\"name\":\"A\",\"evidence\":\"route.ts — `GET`\",\"register\":\"b\"},{\"name\":\"B\",\"definition\":\"a room (e.g. \"Deluxe King\") scoped.\",\"register\":\"b\"}]}",
		},
		{
			"multiple multi-byte chars before defective string",
			"{\"items\":[{\"desc\":\"éàü — café\"},{\"def\":\"say \"hello\" please\"}]}",
		},
		{
			"emoji before defective string",
			"{\"items\":[{\"icon\":\"\U0001F680\U0001F30D\"},{\"def\":\"say \"hello\" ok\"}]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, corrections := parseLenient(t, tt.input)
			hasUnescaped := false
			for _, c := range corrections {
				if c.Kind == lenient.CorrectionUnescapedQuote {
					hasUnescaped = true
				}
			}
			if !hasUnescaped {
				t.Error("expected UnescapedQuote correction")
			}
		})
	}
}

func TestLenientParser_NewlinesInStrings(t *testing.T) {
	input := "{\"msg\":\"line one\\nline two\"}"
	node, _ := parseLenient(t, input)
	if node == nil {
		t.Fatal("expected non-nil node")
	}
}

func TestLenientParser_InvalidJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty input", "", "unexpected end of input"},
		{"just comma", ",", "expected JSON value"},
		{"unclosed object", `{"name": "Alice"`, "expected RBrace"},
		{"unclosed array", `[1, 2`, "expected RBracket"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewLenientParser(tt.input)
			_, _, err := p.Parse()
			if err == nil {
				t.Fatal("expected error but got none")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
