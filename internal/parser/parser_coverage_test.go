package parser

import (
	"strings"
	"testing"

	shapetokenizer "github.com/shapestone/shape-core/pkg/tokenizer"
)

// TestNewParserFromStream tests the NewParserFromStream function
func TestNewParserFromStream(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid JSON from stream",
			input: `{"name": "Alice", "age": 30}`,
			valid: true,
		},
		{
			name:  "invalid JSON from stream",
			input: `{invalid}`,
			valid: false,
		},
		{
			name:  "empty stream",
			input: "",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			stream := shapetokenizer.NewStreamFromReader(reader)
			parser := NewParserFromStream(stream)

			_, err := parser.Parse()
			if tt.valid && err != nil {
				t.Errorf("expected success but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error but parsing succeeded")
			}
		})
	}
}

// TestParseObject_EdgeCases tests edge cases in object parsing for better coverage
func TestParseObject_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "object with non-string key",
			input:       `{123: "value"}`,
			shouldError: true,
			errorMsg:    "object key must be string",
		},
		{
			name:        "object missing colon",
			input:       `{"key" "value"}`,
			shouldError: true,
			errorMsg:    "expected Colon",
		},
		{
			name:        "object with trailing comma",
			input:       `{"key": "value",}`,
			shouldError: true,
			errorMsg:    "object key must be string",
		},
		{
			name:        "object missing closing brace",
			input:       `{"key": "value"`,
			shouldError: true,
			errorMsg:    "expected RBrace",
		},
		{
			name:        "object with invalid value after comma",
			input:       `{"key1": "value1", }`,
			shouldError: true,
			errorMsg:    "object key must be string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestParseMember_EdgeCases tests edge cases in member parsing
func TestParseMember_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "member with non-string key",
			input:       `{true: "value"}`,
			shouldError: true,
			errorMsg:    "object key must be string",
		},
		{
			name:        "member missing value",
			input:       `{"key":}`,
			shouldError: true,
			errorMsg:    "expected JSON value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseArray_EdgeCases tests edge cases in array parsing
func TestParseArray_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "array with trailing comma",
			input:       `[1, 2, 3,]`,
			shouldError: true,
			errorMsg:    "expected JSON value",
		},
		{
			name:        "array missing closing bracket",
			input:       `[1, 2, 3`,
			shouldError: true,
			errorMsg:    "unexpected end of input",
		},
		{
			name:        "array with invalid element",
			input:       `[1, invalid, 3]`,
			shouldError: true,
			errorMsg:    "expected JSON value",
		},
		{
			name:        "nested arrays with errors",
			input:       `[[1, 2], [3, 4,]]`,
			shouldError: true,
			errorMsg:    "expected JSON value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseString_EdgeCases tests edge cases in string parsing
func TestParseString_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "calling parseString when not at string token",
			input:       `123`,
			shouldError: true,
			errorMsg:    "expected string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			// Try to parse as string directly (this will fail)
			_, err := parser.parseString()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseNumber_EdgeCases tests edge cases in number parsing
func TestParseNumber_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "calling parseNumber when not at number token",
			input:       `"not a number"`,
			shouldError: true,
			errorMsg:    "expected number",
		},
		{
			name:        "invalid integer",
			input:       `999999999999999999999999999999999`,
			shouldError: true,
			errorMsg:    "invalid integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.parseNumber()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseBoolean_EdgeCases tests edge cases in boolean parsing
func TestParseBoolean_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "calling parseBoolean when not at boolean token",
			input:       `123`,
			shouldError: true,
			errorMsg:    "expected boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.parseBoolean()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseNull_EdgeCases tests edge cases in null parsing
func TestParseNull_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "calling parseNull when not at null token",
			input:       `123`,
			shouldError: true,
			errorMsg:    "expected null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.parseNull()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParseValue_InvalidToken tests error handling for invalid tokens
func TestParseValue_InvalidToken(t *testing.T) {
	// This requires injecting an invalid token which is harder
	// But we can test with unexpected tokens
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "unexpected token",
			input:       `:`,
			shouldError: true,
		},
		{
			name:        "empty input",
			input:       ``,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestComplexNestedStructures tests deeply nested JSON structures
func TestComplexNestedStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "deeply nested objects",
			input: `{
				"level1": {
					"level2": {
						"level3": {
							"level4": {
								"value": "deep"
							}
						}
					}
				}
			}`,
		},
		{
			name: "deeply nested arrays",
			input: `[
				[
					[
						[
							[1, 2, 3]
						]
					]
				]
			]`,
		},
		{
			name: "mixed nested structures",
			input: `{
				"users": [
					{
						"name": "Alice",
						"addresses": [
							{
								"street": "123 Main St",
								"coordinates": [40.7128, -74.0060]
							}
						]
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestNumberFormats tests various valid number formats
func TestNumberFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "zero", input: "0"},
		{name: "positive integer", input: "123"},
		{name: "negative integer", input: "-456"},
		{name: "decimal", input: "123.456"},
		{name: "negative decimal", input: "-123.456"},
		{name: "scientific notation positive", input: "1.23e10"},
		{name: "scientific notation negative", input: "1.23e-10"},
		{name: "scientific notation capital E", input: "1.23E10"},
		{name: "scientific notation with plus", input: "1.23E+10"},
		{name: "large number", input: "9007199254740991"}, // Max safe integer in JS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if err != nil {
				t.Errorf("unexpected error for valid number %q: %v", tt.input, err)
			}
		})
	}
}

// TestStringEscapeSequences tests various escape sequences
func TestStringEscapeSequences(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "escaped quote", input: `"He said \"Hello\""`},
		{name: "escaped backslash", input: `"C:\\path\\to\\file"`},
		{name: "escaped forward slash", input: `"http:\/\/example.com"`},
		{name: "escaped backspace", input: `"text\bbackspace"`},
		{name: "escaped form feed", input: `"text\fformfeed"`},
		{name: "escaped newline", input: `"line1\nline2"`},
		{name: "escaped carriage return", input: `"text\rreturn"`},
		{name: "escaped tab", input: `"col1\tcol2"`},
		{name: "unicode escape", input: `"Hello \u0041\u0042\u0043"`}, // ABC
		{name: "emoji unicode", input: `"Smile \uD83D\uDE00"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if err != nil {
				t.Errorf("unexpected error for valid string %q: %v", tt.input, err)
			}
		})
	}
}

// TestWhitespaceVariations tests various whitespace configurations
func TestWhitespaceVariations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "spaces", input: `  {  "key"  :  "value"  }  `},
		{name: "tabs", input: "\t{\t\"key\"\t:\t\"value\"\t}\t"},
		{name: "newlines", input: "{\n\"key\"\n:\n\"value\"\n}"},
		{name: "carriage returns", input: "{\r\"key\"\r:\r\"value\"\r}"},
		{name: "mixed whitespace", input: " \t\n\r{ \t\n\r\"key\" \t\n\r: \t\n\r\"value\" \t\n\r} \t\n\r"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if err != nil {
				t.Errorf("unexpected error with whitespace: %v", err)
			}
		})
	}
}
