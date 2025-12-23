package json

import (
	"testing"
)

// ============================================================================
// Document Tests - Builder Methods
// ============================================================================

func TestNewDocument(t *testing.T) {
	doc := NewDocument()
	if doc == nil {
		t.Fatal("NewDocument() returned nil")
	}
	if doc.Size() != 0 {
		t.Errorf("Expected empty document, got size %d", doc.Size())
	}
}

func TestDocument_SetString(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetString("city", "NYC")

	if val, ok := doc.GetString("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got '%s' (ok=%v)", val, ok)
	}
	if val, ok := doc.GetString("city"); !ok || val != "NYC" {
		t.Errorf("Expected 'NYC', got '%s' (ok=%v)", val, ok)
	}
}

func TestDocument_SetInt(t *testing.T) {
	doc := NewDocument().
		SetInt("age", 30).
		SetInt("count", 0)

	if val, ok := doc.GetInt("age"); !ok || val != 30 {
		t.Errorf("Expected 30, got %d (ok=%v)", val, ok)
	}
	if val, ok := doc.GetInt("count"); !ok || val != 0 {
		t.Errorf("Expected 0, got %d (ok=%v)", val, ok)
	}
}

func TestDocument_SetInt64(t *testing.T) {
	doc := NewDocument().
		SetInt64("big", 9223372036854775807)

	if val, ok := doc.GetInt64("big"); !ok || val != 9223372036854775807 {
		t.Errorf("Expected 9223372036854775807, got %d (ok=%v)", val, ok)
	}
}

func TestDocument_SetBool(t *testing.T) {
	doc := NewDocument().
		SetBool("active", true).
		SetBool("deleted", false)

	if val, ok := doc.GetBool("active"); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
	if val, ok := doc.GetBool("deleted"); !ok || val != false {
		t.Errorf("Expected false, got %v (ok=%v)", val, ok)
	}
}

func TestDocument_SetFloat(t *testing.T) {
	doc := NewDocument().
		SetFloat("pi", 3.14159).
		SetFloat("zero", 0.0)

	if val, ok := doc.GetFloat("pi"); !ok || val != 3.14159 {
		t.Errorf("Expected 3.14159, got %f (ok=%v)", val, ok)
	}
	if val, ok := doc.GetFloat("zero"); !ok || val != 0.0 {
		t.Errorf("Expected 0.0, got %f (ok=%v)", val, ok)
	}
}

func TestDocument_SetNull(t *testing.T) {
	doc := NewDocument().SetNull("value")

	if !doc.IsNull("value") {
		t.Error("Expected value to be null")
	}
	if !doc.Has("value") {
		t.Error("Expected value key to exist")
	}
}

func TestDocument_SetObject(t *testing.T) {
	nested := NewDocument().
		SetString("street", "123 Main St").
		SetString("city", "NYC")

	doc := NewDocument().SetObject("address", nested)

	addr, ok := doc.GetObject("address")
	if !ok {
		t.Fatal("Expected address object")
	}

	if val, ok := addr.GetString("city"); !ok || val != "NYC" {
		t.Errorf("Expected 'NYC', got '%s' (ok=%v)", val, ok)
	}
}

func TestDocument_SetArray(t *testing.T) {
	arr := NewArray().
		AddString("go").
		AddString("json")

	doc := NewDocument().SetArray("tags", arr)

	tags, ok := doc.GetArray("tags")
	if !ok {
		t.Fatal("Expected tags array")
	}

	if tags.Len() != 2 {
		t.Errorf("Expected 2 tags, got %d", tags.Len())
	}

	if val, ok := tags.GetString(0); !ok || val != "go" {
		t.Errorf("Expected 'go', got '%s' (ok=%v)", val, ok)
	}
}

// ============================================================================
// Document Tests - Getter Methods
// ============================================================================

func TestDocument_GetString_Missing(t *testing.T) {
	doc := NewDocument()
	if val, ok := doc.GetString("missing"); ok {
		t.Errorf("Expected missing key to return false, got '%s'", val)
	}
}

func TestDocument_GetString_WrongType(t *testing.T) {
	doc := NewDocument().SetInt("age", 30)
	if val, ok := doc.GetString("age"); ok {
		t.Errorf("Expected wrong type to return false, got '%s'", val)
	}
}

func TestDocument_GetInt_FromFloat(t *testing.T) {
	doc := NewDocument().SetFloat("value", 42.0)
	if val, ok := doc.GetInt("value"); !ok || val != 42 {
		t.Errorf("Expected 42, got %d (ok=%v)", val, ok)
	}
}

func TestDocument_GetFloat_FromInt(t *testing.T) {
	doc := NewDocument().SetInt("value", 42)
	if val, ok := doc.GetFloat("value"); !ok || val != 42.0 {
		t.Errorf("Expected 42.0, got %f (ok=%v)", val, ok)
	}
}

func TestDocument_Has(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetNull("value")

	if !doc.Has("name") {
		t.Error("Expected 'name' key to exist")
	}
	if !doc.Has("value") {
		t.Error("Expected 'value' key (null) to exist")
	}
	if doc.Has("missing") {
		t.Error("Expected 'missing' key to not exist")
	}
}

func TestDocument_Keys(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30)

	keys := doc.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	if !keyMap["name"] || !keyMap["age"] {
		t.Errorf("Expected keys 'name' and 'age', got %v", keys)
	}
}

func TestDocument_Remove(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		Remove("age")

	if doc.Has("age") {
		t.Error("Expected 'age' to be removed")
	}
	if !doc.Has("name") {
		t.Error("Expected 'name' to still exist")
	}
	if doc.Size() != 1 {
		t.Errorf("Expected size 1, got %d", doc.Size())
	}
}

// ============================================================================
// Document Tests - JSON Marshaling
// ============================================================================

func TestDocument_JSON(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetBool("active", true)

	jsonStr, err := doc.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Parse back to verify using our own Unmarshal
	var result map[string]interface{}
	if err := Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	// Our implementation returns int64 for whole numbers, not float64
	if result["age"] != int64(30) {
		t.Errorf("Expected age=30, got %v", result["age"])
	}
	if result["active"] != true {
		t.Errorf("Expected active=true, got %v", result["active"])
	}
}

func TestParseDocument(t *testing.T) {
	jsonStr := `{"name":"Alice","age":30,"active":true}`

	doc, err := ParseDocument(jsonStr)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if val, ok := doc.GetString("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got '%s' (ok=%v)", val, ok)
	}
	if val, ok := doc.GetInt("age"); !ok || val != 30 {
		t.Errorf("Expected 30, got %d (ok=%v)", val, ok)
	}
	if val, ok := doc.GetBool("active"); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
}

func TestParseDocument_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid json", `{invalid}`},
		{"array not object", `[1,2,3]`},
		{"primitive not object", `"string"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDocument(tt.input)
			if err == nil {
				t.Error("Expected error for invalid input")
			}
		})
	}
}

func TestDocument_MarshalJSON(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30)

	bytes, err := Marshal(doc)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}

	var result map[string]interface{}
	if err := Unmarshal(bytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got %v (JSON: %s)", result["name"], string(bytes))
	}
	// Check age too - our implementation returns int64 for whole numbers
	if result["age"] != int64(30) {
		t.Errorf("Expected age=30, got %v (JSON: %s)", result["age"], string(bytes))
	}
}

func TestDocument_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"name":"Alice","age":30}`

	var doc Document
	if err := Unmarshal([]byte(jsonStr), &doc); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}

	if val, ok := doc.GetString("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got '%s' (ok=%v)", val, ok)
	}
}

// ============================================================================
// Array Tests - Builder Methods
// ============================================================================

func TestNewArray(t *testing.T) {
	arr := NewArray()
	if arr == nil {
		t.Fatal("NewArray() returned nil")
	}
	if arr.Len() != 0 {
		t.Errorf("Expected empty array, got length %d", arr.Len())
	}
}

func TestArray_AddString(t *testing.T) {
	arr := NewArray().
		AddString("go").
		AddString("json").
		AddString("parser")

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	if val, ok := arr.GetString(0); !ok || val != "go" {
		t.Errorf("Expected 'go', got '%s' (ok=%v)", val, ok)
	}
	if val, ok := arr.GetString(2); !ok || val != "parser" {
		t.Errorf("Expected 'parser', got '%s' (ok=%v)", val, ok)
	}
}

func TestArray_AddInt(t *testing.T) {
	arr := NewArray().
		AddInt(1).
		AddInt(2).
		AddInt(3)

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	if val, ok := arr.GetInt(1); !ok || val != 2 {
		t.Errorf("Expected 2, got %d (ok=%v)", val, ok)
	}
}

func TestArray_AddBool(t *testing.T) {
	arr := NewArray().
		AddBool(true).
		AddBool(false)

	if val, ok := arr.GetBool(0); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
	if val, ok := arr.GetBool(1); !ok || val != false {
		t.Errorf("Expected false, got %v (ok=%v)", val, ok)
	}
}

func TestArray_AddFloat(t *testing.T) {
	arr := NewArray().
		AddFloat(3.14).
		AddFloat(2.71)

	if val, ok := arr.GetFloat(0); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %f (ok=%v)", val, ok)
	}
}

func TestArray_AddNull(t *testing.T) {
	arr := NewArray().
		AddString("value").
		AddNull().
		AddInt(42)

	if !arr.IsNull(1) {
		t.Error("Expected index 1 to be null")
	}
	if arr.IsNull(0) {
		t.Error("Expected index 0 to not be null")
	}
}

func TestArray_AddObject(t *testing.T) {
	obj := NewDocument().SetString("name", "Alice")
	arr := NewArray().AddObject(obj)

	result, ok := arr.GetObject(0)
	if !ok {
		t.Fatal("Expected object at index 0")
	}

	if val, ok := result.GetString("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got '%s' (ok=%v)", val, ok)
	}
}

func TestArray_AddArray(t *testing.T) {
	nested := NewArray().AddInt(1).AddInt(2)
	arr := NewArray().AddArray(nested)

	result, ok := arr.GetArray(0)
	if !ok {
		t.Fatal("Expected array at index 0")
	}

	if result.Len() != 2 {
		t.Errorf("Expected length 2, got %d", result.Len())
	}

	if val, ok := result.GetInt(0); !ok || val != 1 {
		t.Errorf("Expected 1, got %d (ok=%v)", val, ok)
	}
}

// ============================================================================
// Array Tests - Getter Methods
// ============================================================================

func TestArray_Get_OutOfBounds(t *testing.T) {
	arr := NewArray().AddString("value")

	if _, ok := arr.Get(-1); ok {
		t.Error("Expected negative index to return false")
	}
	if _, ok := arr.Get(1); ok {
		t.Error("Expected out of bounds index to return false")
	}
}

func TestArray_GetString_WrongType(t *testing.T) {
	arr := NewArray().AddInt(42)

	if val, ok := arr.GetString(0); ok {
		t.Errorf("Expected wrong type to return false, got '%s'", val)
	}
}

func TestArray_GetInt_FromFloat(t *testing.T) {
	arr := NewArray().AddFloat(42.0)

	if val, ok := arr.GetInt(0); !ok || val != 42 {
		t.Errorf("Expected 42, got %d (ok=%v)", val, ok)
	}
}

func TestArray_GetFloat_FromInt(t *testing.T) {
	arr := NewArray().AddInt(42)

	if val, ok := arr.GetFloat(0); !ok || val != 42.0 {
		t.Errorf("Expected 42.0, got %f (ok=%v)", val, ok)
	}
}

// ============================================================================
// Array Tests - JSON Marshaling
// ============================================================================

func TestArray_JSON(t *testing.T) {
	arr := NewArray().
		AddString("go").
		AddInt(42).
		AddBool(true)

	jsonStr, err := arr.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Parse back to verify using our own Unmarshal
	var result []interface{}
	if err := Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	if result[0] != "go" {
		t.Errorf("Expected 'go', got %v", result[0])
	}
	// Our implementation returns int64 for whole numbers, not float64
	if result[1] != int64(42) {
		t.Errorf("Expected 42, got %v", result[1])
	}
	if result[2] != true {
		t.Errorf("Expected true, got %v", result[2])
	}
}

func TestParseArray(t *testing.T) {
	jsonStr := `["go","json",42,true]`

	arr, err := ParseArray(jsonStr)
	if err != nil {
		t.Fatalf("ParseArray() error: %v", err)
	}

	if arr.Len() != 4 {
		t.Errorf("Expected length 4, got %d", arr.Len())
	}

	if val, ok := arr.GetString(0); !ok || val != "go" {
		t.Errorf("Expected 'go', got '%s' (ok=%v)", val, ok)
	}
	if val, ok := arr.GetInt(2); !ok || val != 42 {
		t.Errorf("Expected 42, got %d (ok=%v)", val, ok)
	}
	if val, ok := arr.GetBool(3); !ok || val != true {
		t.Errorf("Expected true, got %v (ok=%v)", val, ok)
	}
}

func TestParseArray_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid json", `[invalid]`},
		{"object not array", `{"key":"value"}`},
		{"primitive not array", `"string"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseArray(tt.input)
			if err == nil {
				t.Error("Expected error for invalid input")
			}
		})
	}
}

func TestArray_MarshalJSON(t *testing.T) {
	arr := NewArray().
		AddString("go").
		AddInt(42)

	bytes, err := Marshal(arr)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}

	var result []interface{}
	if err := Unmarshal(bytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected length 2, got %d (JSON: %s)", len(result), string(bytes))
	}
	if result[0] != "go" {
		t.Errorf("Expected 'go', got %v", result[0])
	}
	// Our implementation returns int64 for whole numbers, not float64
	if result[1] != int64(42) {
		t.Errorf("Expected 42, got %v", result[1])
	}
}

func TestArray_UnmarshalJSON(t *testing.T) {
	jsonStr := `["go",42,true]`

	var arr Array
	if err := Unmarshal([]byte(jsonStr), &arr); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}

	if val, ok := arr.GetString(0); !ok || val != "go" {
		t.Errorf("Expected 'go', got '%s' (ok=%v)", val, ok)
	}
}

// ============================================================================
// Complex Nested Structure Tests
// ============================================================================

func TestComplexNestedStructure(t *testing.T) {
	// Build a complex nested structure using fluent API
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetObject("address", NewDocument().
			SetString("street", "123 Main St").
			SetString("city", "NYC").
			SetString("zip", "10001")).
		SetArray("tags", NewArray().
			AddString("go").
			AddString("json").
			AddString("parser")).
		SetArray("history", NewArray().
			AddObject(NewDocument().
				SetString("action", "created").
				SetInt("timestamp", 1234567890)).
			AddObject(NewDocument().
				SetString("action", "updated").
				SetInt("timestamp", 1234567900)))

	// Verify nested object
	addr, ok := doc.GetObject("address")
	if !ok {
		t.Fatal("Expected address object")
	}
	if val, ok := addr.GetString("city"); !ok || val != "NYC" {
		t.Errorf("Expected 'NYC', got '%s' (ok=%v)", val, ok)
	}

	// Verify array
	tags, ok := doc.GetArray("tags")
	if !ok {
		t.Fatal("Expected tags array")
	}
	if tags.Len() != 3 {
		t.Errorf("Expected 3 tags, got %d", tags.Len())
	}
	if val, ok := tags.GetString(1); !ok || val != "json" {
		t.Errorf("Expected 'json', got '%s' (ok=%v)", val, ok)
	}

	// Verify array of objects
	history, ok := doc.GetArray("history")
	if !ok {
		t.Fatal("Expected history array")
	}
	if history.Len() != 2 {
		t.Errorf("Expected 2 history items, got %d", history.Len())
	}

	item, ok := history.GetObject(0)
	if !ok {
		t.Fatal("Expected object at history[0]")
	}
	if val, ok := item.GetString("action"); !ok || val != "created" {
		t.Errorf("Expected 'created', got '%s' (ok=%v)", val, ok)
	}

	// Marshal to JSON and verify
	jsonStr, err := doc.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Parse back and verify
	parsed, err := ParseDocument(jsonStr)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if val, ok := parsed.GetString("name"); !ok || val != "Alice" {
		t.Errorf("Expected 'Alice', got '%s' (ok=%v)", val, ok)
	}
}

func TestRoundTripComplexStructure(t *testing.T) {
	// Create complex structure
	original := NewDocument().
		SetString("type", "user").
		SetArray("permissions", NewArray().
			AddString("read").
			AddString("write")).
		SetObject("metadata", NewDocument().
			SetInt("version", 1).
			SetBool("active", true))

	// Marshal to JSON
	jsonStr, err := original.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	// Parse back
	parsed, err := ParseDocument(jsonStr)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	// Verify all fields
	if val, ok := parsed.GetString("type"); !ok || val != "user" {
		t.Errorf("Expected 'user', got '%s' (ok=%v)", val, ok)
	}

	perms, ok := parsed.GetArray("permissions")
	if !ok {
		t.Fatal("Expected permissions array")
	}
	if perms.Len() != 2 {
		t.Errorf("Expected 2 permissions, got %d", perms.Len())
	}

	metadata, ok := parsed.GetObject("metadata")
	if !ok {
		t.Fatal("Expected metadata object")
	}
	if val, ok := metadata.GetInt("version"); !ok || val != 1 {
		t.Errorf("Expected version=1, got %d (ok=%v)", val, ok)
	}
	if val, ok := metadata.GetBool("active"); !ok || val != true {
		t.Errorf("Expected active=true, got %v (ok=%v)", val, ok)
	}
}

// ============================================================================
// Document Tests - JSONIndent
// ============================================================================

func TestDocument_JSONIndent(t *testing.T) {
	doc := NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30)

	// Test with 2-space indent
	pretty, err := doc.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `{
  "age": 30,
  "name": "Alice"
}`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestDocument_JSONIndent_Nested(t *testing.T) {
	doc := NewDocument().
		SetString("type", "user").
		SetObject("details", NewDocument().
			SetString("name", "Bob").
			SetInt("age", 25))

	pretty, err := doc.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `{
  "details": {
    "age": 25,
    "name": "Bob"
  },
  "type": "user"
}`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestDocument_JSONIndent_WithArray(t *testing.T) {
	doc := NewDocument().
		SetString("type", "user").
		SetArray("tags", NewArray().AddString("go").AddString("json"))

	pretty, err := doc.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `{
  "tags": [
    "go",
    "json"
  ],
  "type": "user"
}`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestDocument_JSONIndent_TabIndent(t *testing.T) {
	doc := NewDocument().
		SetString("key", "value")

	pretty, err := doc.JSONIndent("", "\t")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := "{\n\t\"key\": \"value\"\n}"
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestDocument_JSONIndent_WithPrefix(t *testing.T) {
	doc := NewDocument().
		SetString("key", "value")

	pretty, err := doc.JSONIndent(">>", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := "{\n>>  \"key\": \"value\"\n>>}"
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestDocument_JSONIndent_EmptyDocument(t *testing.T) {
	doc := NewDocument()

	pretty, err := doc.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := "{}"
	if pretty != want {
		t.Errorf("JSONIndent() = %q, want %q", pretty, want)
	}
}

// ============================================================================
// Array Tests - JSONIndent
// ============================================================================

func TestArray_JSONIndent(t *testing.T) {
	arr := NewArray().
		AddString("apple").
		AddString("banana").
		AddInt(42)

	pretty, err := arr.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `[
  "apple",
  "banana",
  42
]`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestArray_JSONIndent_Nested(t *testing.T) {
	arr := NewArray().
		AddArray(NewArray().AddInt(1).AddInt(2)).
		AddArray(NewArray().AddInt(3).AddInt(4))

	pretty, err := arr.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `[
  [
    1,
    2
  ],
  [
    3,
    4
  ]
]`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestArray_JSONIndent_WithObjects(t *testing.T) {
	arr := NewArray().
		AddObject(NewDocument().SetString("name", "Alice")).
		AddObject(NewDocument().SetString("name", "Bob"))

	pretty, err := arr.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := `[
  {
    "name": "Alice"
  },
  {
    "name": "Bob"
  }
]`
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}

func TestArray_JSONIndent_EmptyArray(t *testing.T) {
	arr := NewArray()

	pretty, err := arr.JSONIndent("", "  ")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := "[]"
	if pretty != want {
		t.Errorf("JSONIndent() = %q, want %q", pretty, want)
	}
}

func TestArray_JSONIndent_TabIndent(t *testing.T) {
	arr := NewArray().AddInt(1).AddInt(2)

	pretty, err := arr.JSONIndent("", "\t")
	if err != nil {
		t.Fatalf("JSONIndent() error = %v", err)
	}

	want := "[\n\t1,\n\t2\n]"
	if pretty != want {
		t.Errorf("JSONIndent() mismatch:\ngot:\n%s\n\nwant:\n%s", pretty, want)
	}
}
