package json

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// TestDecoder tests the streaming Decoder
func TestDecoder(t *testing.T) {
	t.Run("decode single value", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		jsonStr := `{"name":"Alice","age":30}`
		reader := strings.NewReader(jsonStr)
		decoder := NewDecoder(reader)

		var person Person
		err := decoder.Decode(&person)
		if err != nil {
			t.Fatalf("Decode() error = %v", err)
		}

		if person.Name != "Alice" || person.Age != 30 {
			t.Errorf("Decode() = %+v, want {Name:Alice Age:30}", person)
		}
	})

	t.Run("decode from separate decoders", func(t *testing.T) {
		// Each decoder reads one JSON value from its own reader
		values := []string{
			`{"name":"Alice"}`,
			`{"name":"Bob"}`,
			`{"name":"Charlie"}`,
		}

		expected := []string{"Alice", "Bob", "Charlie"}
		for i, jsonStr := range values {
			reader := strings.NewReader(jsonStr)
			decoder := NewDecoder(reader)

			var result map[string]string
			err := decoder.Decode(&result)
			if err != nil {
				t.Fatalf("Decode() iteration %d error = %v", i, err)
			}

			if result["name"] != expected[i] {
				t.Errorf("Decode() iteration %d = %s, want %s", i, result["name"], expected[i])
			}
		}
	})

	t.Run("decode array of objects", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		jsonStr := `[{"name":"Alice","age":30},{"name":"Bob","age":25}]`
		reader := strings.NewReader(jsonStr)
		decoder := NewDecoder(reader)

		var people []Person
		err := decoder.Decode(&people)
		if err != nil {
			t.Fatalf("Decode() error = %v", err)
		}

		if len(people) != 2 {
			t.Fatalf("len(people) = %d, want 2", len(people))
		}

		if people[0].Name != "Alice" || people[0].Age != 30 {
			t.Errorf("people[0] = %+v, want {Name:Alice Age:30}", people[0])
		}

		if people[1].Name != "Bob" || people[1].Age != 25 {
			t.Errorf("people[1] = %+v, want {Name:Bob Age:25}", people[1])
		}
	})

	t.Run("decode error - invalid JSON", func(t *testing.T) {
		jsonStr := `{invalid}`
		reader := strings.NewReader(jsonStr)
		decoder := NewDecoder(reader)

		var result map[string]string
		err := decoder.Decode(&result)
		if err == nil {
			t.Error("Decode() error = nil, want error")
		}
	})

	t.Run("decode error - non-pointer", func(t *testing.T) {
		jsonStr := `{"name":"Alice"}`
		reader := strings.NewReader(jsonStr)
		decoder := NewDecoder(reader)

		var result map[string]string
		err := decoder.Decode(result) // Not a pointer
		if err == nil {
			t.Error("Decode() error = nil, want error")
		}
	})
}

// TestEncoder tests the streaming Encoder
func TestEncoder(t *testing.T) {
	t.Run("encode single value", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		person := Person{Name: "Alice", Age: 30}

		var buf bytes.Buffer
		encoder := NewEncoder(&buf)

		err := encoder.Encode(person)
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		expected := `{"age":30,"name":"Alice"}`
		result := strings.TrimSpace(buf.String())
		if result != expected {
			t.Errorf("Encode() = %s, want %s", result, expected)
		}
	})

	t.Run("encode multiple values", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
		}

		people := []Person{
			{Name: "Alice"},
			{Name: "Bob"},
			{Name: "Charlie"},
		}

		var buf bytes.Buffer
		encoder := NewEncoder(&buf)

		for _, p := range people {
			err := encoder.Encode(p)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
		}

		expected := `{"name":"Alice"}
{"name":"Bob"}
{"name":"Charlie"}
`
		if buf.String() != expected {
			t.Errorf("Encode() = %q, want %q", buf.String(), expected)
		}
	})

	t.Run("encode array", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}

		var buf bytes.Buffer
		encoder := NewEncoder(&buf)

		err := encoder.Encode(people)
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		expected := `[{"age":30,"name":"Alice"},{"age":25,"name":"Bob"}]`
		result := strings.TrimSpace(buf.String())
		if result != expected {
			t.Errorf("Encode() = %s, want %s", result, expected)
		}
	})

	t.Run("encode map", func(t *testing.T) {
		data := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		var buf bytes.Buffer
		encoder := NewEncoder(&buf)

		err := encoder.Encode(data)
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		result := strings.TrimSpace(buf.String())
		// Check that it's a valid JSON object with the right keys
		if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
			t.Errorf("Encode() = %s, expected JSON object", result)
		}
		for _, key := range []string{`"a":1`, `"b":2`, `"c":3`} {
			if !strings.Contains(result, key) {
				t.Errorf("Encode() = %s, missing key %s", result, key)
			}
		}
	})
}

// failingWriter is an io.Writer that always returns an error
type failingWriter struct {
	failAfter int
	written   int
}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	if w.written >= w.failAfter {
		return 0, errors.New("write error")
	}
	w.written += len(p)
	return len(p), nil
}

// TestEncoder_WriteErrors tests error handling in Encode
func TestEncoder_WriteErrors(t *testing.T) {
	t.Run("write error on data", func(t *testing.T) {
		// Create a writer that fails immediately
		w := &failingWriter{failAfter: 0}
		encoder := NewEncoder(w)

		data := map[string]string{"name": "Alice"}
		err := encoder.Encode(data)
		if err == nil {
			t.Error("Encode() error = nil, want write error")
		}
		if !strings.Contains(err.Error(), "write error") {
			t.Errorf("Encode() error = %v, want write error", err)
		}
	})

	t.Run("write error on newline", func(t *testing.T) {
		// Create a writer that succeeds once (for data) but fails on second write (newline)
		w := &failingWriter{failAfter: 1}
		encoder := NewEncoder(w)

		data := map[string]string{"name": "Alice"}
		err := encoder.Encode(data)
		if err == nil {
			t.Error("Encode() error = nil, want write error on newline")
		}
		if !strings.Contains(err.Error(), "write error") {
			t.Errorf("Encode() error = %v, want write error", err)
		}
	})
}

// TestDecoder_Encoder_RoundTrip tests round-trip encoding/decoding
func TestDecoder_Encoder_RoundTrip(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	original := Person{Name: "Alice", Age: 30}

	// Encode
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)
	err := encoder.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	// Decode
	decoder := NewDecoder(&buf)
	var decoded Person
	err = decoder.Decode(&decoded)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	// Compare
	if decoded.Name != original.Name || decoded.Age != original.Age {
		t.Errorf("decoded = %+v, want %+v", decoded, original)
	}
}
