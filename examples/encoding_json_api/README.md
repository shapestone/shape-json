# encoding/json Compatible API Example

This example demonstrates shape-json's complete compatibility with Go's standard `encoding/json` package, making it a drop-in replacement.

## Overview

shape-json provides 100% API compatibility with `encoding/json`, allowing you to replace:
```go
import "encoding/json"
```

with:
```go
import "github.com/shapestone/shape-json/pkg/json"
```

All existing code continues to work without changes.

## Features Demonstrated

### 1. Marshal - Convert Go Struct to JSON
```go
person := Person{Name: "Alice", Age: 30}
data, err := json.Marshal(person)
// Result: {"name":"Alice","age":30}
```

### 2. Unmarshal - Parse JSON into Go Struct
```go
jsonData := []byte(`{"name":"Bob","age":25}`)
var person Person
err := json.Unmarshal(jsonData, &person)
// person.Name == "Bob", person.Age == 25
```

### 3. Encoder - Stream JSON to io.Writer
```go
encoder := json.NewEncoder(os.Stdout)
err := encoder.Encode(person)
// Writes JSON to stdout with newline
```

### 4. Decoder - Stream JSON from io.Reader
```go
file, _ := os.Open("data.json")
decoder := json.NewDecoder(file)
var person Person
err := decoder.Decode(&person)
```

### 5. Struct Tags
```go
type User struct {
    PublicName  string `json:"name"`              // Rename field
    Password    string `json:"-"`                 // Skip field
    Email       string `json:"email,omitempty"`   // Omit if empty
    Count       int    `json:"count,string"`      // Marshal as string
    NoTag       string                            // Uses field name
}
```

**Supported Tags:**
- `json:"name"` - Rename field in JSON
- `json:"-"` - Skip field entirely
- `json:",omitempty"` - Omit if zero value
- `json:",string"` - Marshal number/bool as string
- No tag - Use field name as-is

### 6. Round Trip
```go
// Marshal
data, _ := json.Marshal(original)

// Unmarshal
var decoded Person
json.Unmarshal(data, &decoded)

// decoded == original ✓
```

## Running the Example

```bash
cd examples/encoding_json_api
go run main.go
```

## Expected Output

```
=== shape-json: encoding/json Compatible API Demo ===

--- Marshal Demo ---
Marshaled JSON:
{"address":{"city":"Seattle","state":"WA"},"age":30,"email":"alice@example.com","name":"Alice"}

--- Unmarshal Demo ---
Unmarshaled struct:
  Name: Bob
  Age: 25
  City: Portland
  State: OR

--- Encoder Demo ---
Encoding multiple values to stream:
Encoded stream:
{"address":{"city":"Seattle","state":"WA"},"age":30,"name":"Alice"}
{"address":{"city":"Portland","state":"OR"},"age":25,"name":"Bob"}
{"address":{"city":"San Francisco","state":"CA"},"age":35,"name":"Charlie"}

--- Decoder Demo ---
Decoded from stream:
  Name: Alice
  Age: 30
  City: Seattle

--- Struct Tags Demo ---
With struct tags:
{"NoTag":"visible","count":"42","display_name":"Alice","username":"alice123"}
Note: Password is skipped, Email is omitted (empty), Count is string

--- Round Trip Demo ---
Original marshaled: {"address":{"city":"Seattle","state":"WA"},"age":30,"email":"alice@example.com","name":"Alice"}
✓ Round trip successful - all fields match!
```

## API Reference

### Marshal/Unmarshal Functions

**`Marshal(v interface{}) ([]byte, error)`**
- Converts Go value to JSON bytes
- Supports structs, maps, slices, primitives
- Respects struct tags

**`Unmarshal(data []byte, v interface{}) error`**
- Parses JSON bytes into Go value
- Handles type conversion automatically
- Validates JSON structure

### Encoder/Decoder Streaming

**`NewEncoder(w io.Writer) *Encoder`**
- Creates streaming encoder
- Writes JSON to any io.Writer
- Adds newline after each value

**`NewDecoder(r io.Reader) *Decoder`**
- Creates streaming decoder
- Reads JSON from any io.Reader
- Supports buffered reading

**Methods:**
- `encoder.Encode(v interface{}) error` - Write value as JSON
- `decoder.Decode(v interface{}) error` - Read JSON into value

## Drop-in Replacement

To replace `encoding/json` in existing code:

**Before:**
```go
import "encoding/json"

data, err := json.Marshal(value)
var result MyStruct
err = json.Unmarshal(data, &result)
```

**After:**
```go
import "github.com/shapestone/shape-json/pkg/json"

// Same code - no changes needed!
data, err := json.Marshal(value)
var result MyStruct
err = json.Unmarshal(data, &result)
```

## Why Use shape-json Instead of encoding/json?

1. **Unified Ecosystem** - Works seamlessly with other Shape parsers (YAML, TOML, etc.)
2. **AST Access** - Can access parsed AST for advanced use cases
3. **Parser Framework** - Built on Shape's robust parser infrastructure
4. **Format Detection** - Includes JSON validation and detection APIs
5. **DOM API** - User-friendly fluent API for JSON manipulation

## Compatibility Notes

shape-json provides:
- ✅ Complete API compatibility with `encoding/json`
- ✅ All struct tag options (`omitempty`, `string`, `-`)
- ✅ Streaming encoder/decoder
- ✅ Marshal/Unmarshal functions
- ✅ Type conversions and validation
- ✅ RFC 8259 compliance

## Related Examples

- [DOM Builder](../dom_builder/) - Fluent API for JSON manipulation
- [Basic Parsing](../main.go) - Low-level AST parsing
- [Parse Reader](../parse_reader/) - Streaming file parsing
- [Format Detection](../format_detection/) - JSON validation

## Documentation

See [pkg/json](../../pkg/json) for complete API documentation.
