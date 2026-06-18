package json_test

import (
	"testing"

	"github.com/shapestone/shape-json/pkg/json"
)

func TestRepair_ValidJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty object", `{}`},
		{"simple object", `{"name":"Alice","age":30}`},
		{"empty array", `[]`},
		{"simple array", `[1,2,3]`},
		{"string", `"hello"`},
		{"number", `42`},
		{"boolean", `true`},
		{"null", `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Valid JSON should round-trip (may reorder keys)
			if err := json.Validate(result); err != nil {
				t.Errorf("repaired output is not valid JSON: %v", err)
			}
		})
	}
}

func TestRepair_TrailingCommas(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"object", `{"a":1,"b":2,}`, `{"a":1,"b":2}`},
		{"array", `[1,2,3,]`, `[1,2,3]`},
		{"nested", `{"items":[1,2,],"x":1,}`, `{"items":[1,2],"x":1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestRepair_SingleQuotedStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"object with single quotes", `{'name':'Alice'}`, `{"name":"Alice"}`},
		{"mixed quotes", `{'name':"Alice"}`, `{"name":"Alice"}`},
		{"single quote value with double inside", `{'msg':'say "hi"'}`, `{"msg":"say \"hi\""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestRepair_UnquotedKeys(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", `{name:"Alice"}`, `{"name":"Alice"}`},
		{"multiple", `{name:"Alice",age:30}`, `{"age":30,"name":"Alice"}`},
		{"with underscore", `{_id:1}`, `{"_id":1}`},
		{"with dollar", `{$ref:"#/foo"}`, `{"$ref":"#\/foo"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestRepair_Comments(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"line comment", "{\n// comment\n\"a\":1\n}"},
		{"block comment", "{\"a\":1 /* comment */}"},
		{"comment before value", "// header\n{\"a\":1}"},
		{"multiple comments", "{\n// one\n\"a\":1,\n/* two */\n\"b\":2\n}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := json.Validate(result); err != nil {
				t.Errorf("repaired output is not valid JSON: %v", err)
			}
		})
	}
}

func TestRepair_DuplicateKeys(t *testing.T) {
	result, err := json.Repair(`{"key":"first","key":"second"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Last value should win
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("failed to unmarshal repaired output: %v", err)
	}
	if m["key"] != "second" {
		t.Errorf("expected last value 'second', got %v", m["key"])
	}
}

func TestRepair_UnescapedQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", `{"msg":"He said "hello" to me"}`},
		{"multiple", `{"msg":"She said "hi" and "bye""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Repair(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := json.Validate(result); err != nil {
				t.Errorf("repaired output is not valid JSON: %v\nresult: %s", err, result)
			}
		})
	}
}

func TestRepair_MixedErrors(t *testing.T) {
	input := `{
		// settings file
		name: 'Alice',
		age: 30,
	}`
	result, err := json.Repair(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := json.Validate(result); err != nil {
		t.Errorf("repaired output is not valid JSON: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("cannot unmarshal repaired output: %v", err)
	}
	if m["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", m["name"])
	}
}

func TestRepairBytes(t *testing.T) {
	input := []byte(`{"a":1,}`)
	result, err := json.RepairBytes(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != `{"a":1}` {
		t.Errorf("got %q, want %q", string(result), `{"a":1}`)
	}
}

func TestRepairWithCorrections(t *testing.T) {
	input := `{name: 'Alice',}`
	result, corrections, err := json.RepairWithCorrections(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := json.Validate(result); err != nil {
		t.Errorf("repaired output is not valid JSON: %v", err)
	}

	if len(corrections) == 0 {
		t.Fatal("expected corrections, got none")
	}

	kinds := make(map[json.CorrectionKind]bool)
	for _, c := range corrections {
		kinds[c.Kind] = true
	}

	if !kinds[json.CorrectionUnquotedKey] {
		t.Error("expected UnquotedKey correction")
	}
	if !kinds[json.CorrectionSingleQuote] {
		t.Error("expected SingleQuote correction")
	}
	if !kinds[json.CorrectionTrailingComma] {
		t.Error("expected TrailingComma correction")
	}
}

func TestRepairWithCorrections_NoCorrections(t *testing.T) {
	_, corrections, err := json.RepairWithCorrections(`{"valid":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(corrections) != 0 {
		t.Errorf("expected no corrections for valid JSON, got %d", len(corrections))
	}
}

func TestRepair_RoundTrip(t *testing.T) {
	inputs := []string{
		`{"a":1,}`,
		`{'key':'value'}`,
		`{key: "value"}`,
		`{// comment
"a":1}`,
		`{"key":"first","key":"second"}`,
		`{"msg":"He said "hello" to me"}`,
		`{
			// config
			name: 'test',
			items: [1, 2, 3,],
		}`,
	}

	for _, input := range inputs {
		t.Run("", func(t *testing.T) {
			repaired, err := json.Repair(input)
			if err != nil {
				t.Fatalf("Repair failed: %v", err)
			}

			// The repaired output must be valid strict JSON
			if err := json.Validate(repaired); err != nil {
				t.Errorf("repaired output is not valid JSON: %v\ninput:    %s\nrepaired: %s", err, input, repaired)
			}

			// It must also round-trip through Parse
			_, err = json.Parse(repaired)
			if err != nil {
				t.Errorf("Parse failed on repaired output: %v\nrepaired: %s", err, repaired)
			}
		})
	}
}

func TestRepair_UnescapedQuoteNotFirstElement(t *testing.T) {
	// Regression: accommodation-fulfillment bug — unescaped quotes in a
	// non-first array element with newlines inside strings.
	input := "{\"items\":[{\"name\":\"Room Daily\",\"kind\":\"entity\",\"definition\":\"A per-day record for a physical room.\",\"evidence\":\"route.ts\",\"register\":\"business\"},{\"name\":\"Room Type\",\"kind\":\"entity\",\"definition\":\"A named category of accommodation room (e.g. \"Deluxe King\") scoped to a program, used to group inventory blocks and link room assignments to stays.\",\"evidence\":\"route.ts\",\"register\":\"business\"}]}"

	result, err := json.Repair(input)
	if err != nil {
		t.Fatalf("Repair failed: %v", err)
	}
	if err := json.Validate(result); err != nil {
		t.Errorf("repaired output is not valid JSON: %v", err)
	}
}

func TestRepair_NewlinesInStrings(t *testing.T) {
	input := "{\"msg\":\"hello\\nworld\",\"ok\":true}"
	result, err := json.Repair(input)
	if err != nil {
		t.Fatalf("Repair failed: %v", err)
	}
	if err := json.Validate(result); err != nil {
		t.Errorf("repaired output is not valid JSON: %v", err)
	}
}

func TestRepair_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"just comma", ","},
		{"unclosed object", `{"a":1`},
		{"unclosed array", `[1,2`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := json.Repair(tt.input)
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}
