package parser

import (
	"testing"
)

// FuzzParser is a fuzzing test that feeds random inputs to the parser
// to ensure it doesn't panic or exhibit undefined behavior.
//
// Run with: go test -fuzz=FuzzParser -fuzztime=30s
func FuzzParser(f *testing.F) {
	// Seed corpus with valid JSON examples
	seeds := []string{
		`{}`,
		`[]`,
		`null`,
		`true`,
		`false`,
		`0`,
		`123`,
		`-456`,
		`123.456`,
		`1.23e10`,
		`""`,
		`"hello"`,
		`"escaped \"quote\""`,
		`{"key": "value"}`,
		`{"a": 1, "b": 2}`,
		`[1, 2, 3]`,
		`["a", "b", "c"]`,
		`{"nested": {"obj": {"value": 42}}}`,
		`[[[[[[1]]]]]]`,
		`{"array": [1, 2, {"nested": true}]}`,
		// Edge cases
		`   {}   `,
		`{"":""}`,
		`[null, null]`,
		`{"a":null,"b":false,"c":0,"d":"","e":[],"f":{}}`,
		// Complex real-world-like examples
		`{"users":[{"id":1,"name":"Alice","active":true},{"id":2,"name":"Bob","active":false}]}`,
		`{"config":{"timeout":30,"retries":3,"endpoints":["http://a.com","http://b.com"]}}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// The fuzzer will generate random variations
	f.Fuzz(func(t *testing.T, input string) {
		// The parser should never panic, regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, err := parser.Parse()

		// We don't care if parsing fails (most random inputs will be invalid)
		// We only care that it doesn't panic and returns cleanly
		_ = err
	})
}

// FuzzParserObjects fuzzes specifically object parsing
func FuzzParserObjects(f *testing.F) {
	seeds := []string{
		`{}`,
		`{"a":1}`,
		`{"a":1,"b":2}`,
		`{"nested":{"inner":"value"}}`,
		`{"array":[1,2,3]}`,
		`{"bool":true,"null":null}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on object input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, _ = parser.Parse()
	})
}

// FuzzParserArrays fuzzes specifically array parsing
func FuzzParserArrays(f *testing.F) {
	seeds := []string{
		`[]`,
		`[1]`,
		`[1,2,3]`,
		`["a","b","c"]`,
		`[true,false,null]`,
		`[[1,2],[3,4]]`,
		`[{"a":1},{"b":2}]`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on array input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, _ = parser.Parse()
	})
}

// FuzzParserStrings fuzzes specifically string parsing
func FuzzParserStrings(f *testing.F) {
	seeds := []string{
		`""`,
		`"hello"`,
		`"escaped \" quote"`,
		`"newline\n"`,
		`"tab\t"`,
		`"unicode\u0041"`,
		`"backslash\\"`,
		`"mixed\n\t\r"`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on string input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, _ = parser.Parse()
	})
}

// FuzzParserNumbers fuzzes specifically number parsing
func FuzzParserNumbers(f *testing.F) {
	seeds := []string{
		`0`,
		`123`,
		`-456`,
		`123.456`,
		`-123.456`,
		`1e10`,
		`1e-10`,
		`1.23E10`,
		`1.23E+10`,
		`-1.23e-10`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on number input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, _ = parser.Parse()
	})
}

// FuzzParserNested fuzzes deeply nested structures
func FuzzParserNested(f *testing.F) {
	seeds := []string{
		`{"a":{"b":{"c":1}}}`,
		`[[[1]]]`,
		`{"a":[{"b":[{"c":1}]}]}`,
		`[{"a":1},{"b":[2,3]}]`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on nested input %q: %v", input, r)
			}
		}()

		parser := NewParser(input)
		_, _ = parser.Parse()
	})
}
