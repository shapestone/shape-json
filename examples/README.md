# shape-json Examples

Comprehensive examples demonstrating all features of shape-json.

## Quick Start Guide

**New to shape-json?** Start here based on your use case:

| Use Case | Example | Description |
|----------|---------|-------------|
| **🌟 Building/manipulating JSON** | [dom_builder/](dom_builder/) | **RECOMMENDED** - Fluent API, no type assertions |
| **🔄 Drop-in encoding/json replacement** | [encoding_json_api/](encoding_json_api/) | Marshal/Unmarshal with structs |
| **✅ Validating JSON** | [format_detection/](format_detection/) | Detect and validate JSON format |
| **📁 Parsing large files** | [parse_reader/](parse_reader/) | Streaming file parsing |
| **🔧 Advanced AST manipulation** | [main.go](main.go) | Low-level AST API |

## Examples Overview

### 1. DOM Builder API (Recommended) ⭐

**Location:** [dom_builder/](dom_builder/)

Build and manipulate JSON with a clean, type-safe API - no type assertions needed!

```go
// Build JSON fluently
doc := json.NewDocument().
    SetString("name", "Alice").
    SetInt("age", 30).
    SetObject("address", json.NewDocument().
        SetString("city", "NYC"))

// Access with type-safe getters
doc, _ := json.ParseDocument(`{"user":{"name":"Bob"}}`)
user, _ := doc.GetObject("user")
name, _ := user.GetString("name")  // "Bob" - clean!
```

**Why use this:**
- ✅ No type assertions (`(*ast.ObjectNode)`)
- ✅ Fluent builder pattern (method chaining)
- ✅ Type-safe getters (`GetString`, `GetInt`, `GetBool`)
- ✅ Clean array semantics
- ✅ IDE autocomplete support

**Run:**
```bash
go run dom_builder/main.go
```

**Learn more:** [dom_builder/README.md](dom_builder/README.md)

---

### 2. encoding/json Compatible API

**Location:** [encoding_json_api/](encoding_json_api/)

Drop-in replacement for Go's standard `encoding/json` package.

```go
// Marshal Go structs to JSON
person := Person{Name: "Alice", Age: 30}
data, _ := json.Marshal(person)

// Unmarshal JSON to Go structs
var person Person
json.Unmarshal(data, &person)

// Stream with Encoder/Decoder
encoder := json.NewEncoder(writer)
encoder.Encode(person)
```

**Why use this:**
- ✅ 100% `encoding/json` API compatible
- ✅ Works with existing code
- ✅ Full struct tag support (`omitempty`, `string`, `-`)
- ✅ Streaming encoder/decoder

**Run:**
```bash
go run encoding_json_api/main.go
```

**Learn more:** [encoding_json_api/README.md](encoding_json_api/README.md)

---

### 3. Format Detection

**Location:** [format_detection/](format_detection/)

Validate and detect JSON format before parsing.

```go
// Validate JSON string
format, err := json.DetectFormat(`{"name": "Alice"}`)
if err != nil {
    fmt.Println("Not valid JSON:", err)
} else {
    fmt.Println("Format:", format)  // "JSON"
}

// Validate from file
file, _ := os.Open("data.json")
format, err := json.DetectFormatFromReader(file)
```

**Why use this:**
- ✅ Pre-validation before parsing
- ✅ Descriptive error messages
- ✅ Works with strings and io.Reader
- ✅ Useful for API validation

**Run:**
```bash
go run format_detection/main.go
```

**Learn more:** [format_detection/README.md](format_detection/README.md)

---

### 4. Streaming File Parsing

**Location:** [parse_reader/](parse_reader/)

Parse large JSON files with constant memory usage.

```go
// Parse from file
file, _ := os.Open("large-data.json")
defer file.Close()

node, err := json.ParseReader(file)
// Memory-efficient streaming parsing

// Parse from any io.Reader
reader := strings.NewReader(jsonString)
node, err := json.ParseReader(reader)
```

**Why use this:**
- ✅ Constant memory usage
- ✅ Works with large files
- ✅ Supports any io.Reader
- ✅ Network stream compatible

**Run:**
```bash
go run parse_reader/main.go
```

**Learn more:** [parse_reader/README.md](parse_reader/README.md)

---

### 5. Basic AST API

**Location:** [main.go](main.go)

Low-level AST parsing for advanced use cases.

```go
// Parse to AST
node, _ := json.Parse(`{"name": "Alice", "age": 30}`)

// Access AST nodes (requires type assertions)
obj := node.(*ast.ObjectNode)
nameNode, _ := obj.GetProperty("name")
name := nameNode.(*ast.LiteralNode).Value().(string)
```

**Why use this:**
- ✅ Fine-grained control
- ✅ AST manipulation
- ✅ Shape ecosystem integration
- ⚠️ Requires type assertions
- ⚠️ More verbose

**Run:**
```bash
go run main.go
```

**Note:** For typical JSON work, use the [DOM API](dom_builder/) instead.

---

## Comparison: Which API to Use?

| Feature | DOM API | encoding/json API | AST API |
|---------|---------|-------------------|---------|
| **Ease of use** | ⭐⭐⭐ Best | ⭐⭐⭐ Best | ⭐ Verbose |
| **Type assertions** | ✅ None | ✅ None | ❌ Many |
| **Method chaining** | ✅ Yes | ❌ No | ❌ No |
| **Struct marshaling** | ❌ No | ✅ Yes | ❌ No |
| **Dynamic JSON** | ✅ Best | ❌ Limited | ✅ Advanced |
| **IDE support** | ✅ Excellent | ✅ Excellent | ⭐ Basic |
| **Use case** | Build/manipulate JSON | Go struct mapping | Advanced parsing |

**Recommendation:**
- **Most users:** Use [DOM API](dom_builder/) for building/manipulating JSON
- **Struct mapping:** Use [encoding/json API](encoding_json_api/) for Go structs
- **Advanced:** Use AST API only for parser development

---

## Running All Examples

```bash
# DOM Builder (recommended)
go run dom_builder/main.go

# encoding/json compatibility
go run encoding_json_api/main.go

# Format detection
go run format_detection/main.go

# File parsing
go run parse_reader/main.go

# Basic AST parsing
go run main.go
```

## Additional Features

### JSONPath Queries

shape-json includes a complete JSONPath implementation (RFC 9535):

```go
import "github.com/shapestone/shape-json/pkg/jsonpath"

// Parse query
expr, _ := jsonpath.ParseString("$.store.book[*].author")

// Execute against data
results := expr.Get(data)
```

See [pkg/jsonpath/README.md](../pkg/jsonpath/README.md) for details.

### Format Validation

Pre-validate JSON before processing:

```go
format, err := json.DetectFormat(input)
if err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}
```

See [format_detection/](format_detection/) example.

## Documentation

- **Main README:** [../README.md](../README.md)
- **API Documentation:** [pkg/json](../pkg/json)
- **Testing Guide:** [docs/TESTING.md](../docs/TESTING.md)
- **Grammar Spec:** [docs/grammar/json.ebnf](../docs/grammar/json.ebnf)
- **JSONPath:** [pkg/jsonpath/README.md](../pkg/jsonpath/README.md)

## Getting Help

- **Examples not working?** Check Go version (1.23+ required)
- **Type assertion errors?** Use [DOM API](dom_builder/) instead of AST API
- **Need struct marshaling?** See [encoding/json_api/](encoding_json_api/)
- **Questions?** See [main README](../README.md) or file an issue

## Example Code Quality

All examples:
- ✅ Compile and run successfully
- ✅ Include comprehensive comments
- ✅ Demonstrate best practices
- ✅ Have accompanying README documentation
- ✅ Show error handling patterns

---

**Ready to start?** → Begin with the [DOM Builder API](dom_builder/) example! 🚀
