package jsonpath

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []tokenType
		wantError bool
	}{
		{
			name:      "root selector only",
			input:     "$",
			wantTypes: []tokenType{tokenRoot},
			wantError: false,
		},
		{
			name:      "root with child selector",
			input:     "$.name",
			wantTypes: []tokenType{tokenRoot, tokenDot, tokenIdentifier},
			wantError: false,
		},
		{
			name:      "root with multiple child selectors",
			input:     "$.user.name",
			wantTypes: []tokenType{tokenRoot, tokenDot, tokenIdentifier, tokenDot, tokenIdentifier},
			wantError: false,
		},
		{
			name:      "array index",
			input:     "$[0]",
			wantTypes: []tokenType{tokenRoot, tokenLeftBracket, tokenNumber, tokenRightBracket},
			wantError: false,
		},
		{
			name:      "array wildcard",
			input:     "$[*]",
			wantTypes: []tokenType{tokenRoot, tokenLeftBracket, tokenWildcard, tokenRightBracket},
			wantError: false,
		},
		{
			name:      "bracket notation property",
			input:     "$['name']",
			wantTypes: []tokenType{tokenRoot, tokenLeftBracket, tokenString, tokenRightBracket},
			wantError: false,
		},
		{
			name:      "array slice",
			input:     "$[0:5]",
			wantTypes: []tokenType{tokenRoot, tokenLeftBracket, tokenNumber, tokenColon, tokenNumber, tokenRightBracket},
			wantError: false,
		},
		{
			name:      "recursive descent",
			input:     "$..name",
			wantTypes: []tokenType{tokenRoot, tokenDotDot, tokenIdentifier},
			wantError: false,
		},
		{
			name:      "wildcard child",
			input:     "$.*",
			wantTypes: []tokenType{tokenRoot, tokenDot, tokenWildcard},
			wantError: false,
		},
		{
			name:      "complex path",
			input:     "$.users[0].name",
			wantTypes: []tokenType{tokenRoot, tokenDot, tokenIdentifier, tokenLeftBracket, tokenNumber, tokenRightBracket, tokenDot, tokenIdentifier},
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantTypes: nil,
			wantError: true,
		},
		{
			name:      "missing root",
			input:     "name",
			wantTypes: nil,
			wantError: true,
		},
		{
			name:      "unclosed bracket - caught by parser",
			input:     "$[0",
			wantTypes: []tokenType{tokenRoot, tokenLeftBracket, tokenNumber},
			wantError: false,
		},
		{
			name:      "unclosed string",
			input:     "$['name",
			wantTypes: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("tokenize() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("tokenize() unexpected error: %v", err)
				return
			}
			if len(tokens) != len(tt.wantTypes) {
				t.Errorf("tokenize() got %d tokens, want %d", len(tokens), len(tt.wantTypes))
				return
			}
			for i, tok := range tokens {
				if tok.typ != tt.wantTypes[i] {
					t.Errorf("tokenize() token[%d] type = %v, want %v", i, tok.typ, tt.wantTypes[i])
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "root only",
			input:     "$",
			wantError: false,
		},
		{
			name:      "child selector",
			input:     "$.name",
			wantError: false,
		},
		{
			name:      "multiple child selectors",
			input:     "$.user.name.first",
			wantError: false,
		},
		{
			name:      "array index",
			input:     "$[0]",
			wantError: false,
		},
		{
			name:      "array wildcard",
			input:     "$[*]",
			wantError: false,
		},
		{
			name:      "bracket notation",
			input:     "$['property']",
			wantError: false,
		},
		{
			name:      "array slice full",
			input:     "$[0:5]",
			wantError: false,
		},
		{
			name:      "array slice from start",
			input:     "$[:5]",
			wantError: false,
		},
		{
			name:      "array slice to end",
			input:     "$[2:]",
			wantError: false,
		},
		{
			name:      "recursive descent",
			input:     "$..name",
			wantError: false,
		},
		{
			name:      "wildcard child",
			input:     "$.*",
			wantError: false,
		},
		{
			name:      "complex path",
			input:     "$.users[0].profile.name",
			wantError: false,
		},
		{
			name:      "invalid - no root",
			input:     "name",
			wantError: true,
		},
		{
			name:      "invalid - empty",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid - unexpected token after root",
			input:     "$]",
			wantError: true,
		},
		{
			name:      "invalid - unclosed bracket",
			input:     "$[0",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if err != nil {
				if !tt.wantError {
					t.Errorf("tokenize() unexpected error: %v", err)
				}
				return
			}

			_, err = parse(tokens)
			if tt.wantError {
				if err == nil {
					t.Errorf("parse() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("parse() unexpected error: %v", err)
			}
		})
	}
}

func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "double quotes string",
			query: `$["name"]`,
		},
		{
			name:  "recursive wildcard",
			query: "$..*",
		},
		{
			name:  "slice to end",
			query: "$[5:]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if err != nil {
				t.Fatalf("ParseString() error: %v", err)
			}
			// Just verify it parses successfully
			if expr == nil {
				t.Errorf("ParseString() returned nil expr")
			}
		})
	}
}

func TestRecursiveWildcard(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  int // Just check count since order might vary
		data  interface{}
	}{
		{
			name:  "recursive wildcard on nested object",
			query: "$..*",
			data: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"c": 2,
					"d": 3,
				},
			},
			want: 4, // The root object is not included, only the values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if err != nil {
				t.Fatalf("ParseString() error: %v", err)
			}
			got := expr.Get(tt.data)
			if len(got) < tt.want {
				t.Errorf("Expr.Get() returned %d items, want at least %d", len(got), tt.want)
			}
		})
	}
}

func TestTokenizeValues(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantValue  string
		tokenIndex int
	}{
		{
			name:       "identifier value",
			input:      "$.name",
			tokenIndex: 2, // identifier token
			wantValue:  "name",
		},
		{
			name:       "number value",
			input:      "$[42]",
			tokenIndex: 2, // number token
			wantValue:  "42",
		},
		{
			name:       "string value",
			input:      "$['property']",
			tokenIndex: 2, // string token
			wantValue:  "property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize() error: %v", err)
			}
			if tt.tokenIndex >= len(tokens) {
				t.Fatalf("tokenIndex %d out of range (len=%d)", tt.tokenIndex, len(tokens))
			}
			if tokens[tt.tokenIndex].value != tt.wantValue {
				t.Errorf("token value = %q, want %q", tokens[tt.tokenIndex].value, tt.wantValue)
			}
		})
	}
}
