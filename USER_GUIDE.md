# shape-json: Complete User Guide

**Version:** 0.8.0-dev
**Repository:** github.com/shapestone/shape-json
**License:** Apache License 2.0

> **New in 0.8.0:** Unified architecture with zero encoding/json dependency. All APIs now use a single custom parser and renderer for architectural purity.

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [API Documentation](#api-documentation)
   - [encoding/json Compatible API](#encodingjson-compatible-api)
   - [Fluent DOM API (Recommended)](#fluent-dom-api-recommended)
   - [JSON Validation](#json-validation)
   - [Streaming Parser](#streaming-parser)
   - [Low-Level AST API](#low-level-ast-api)
5. [JSONPath Query Engine](#jsonpath-query-engine)
6. [Performance & Benchmarks](#performance--benchmarks)
7. [Architecture](#architecture)
8. [Testing & Quality](#testing--quality)
9. [Examples](#examples)
10. [Migration Guide](#migration-guide)

---

## Overview

shape-json is a comprehensive JSON parser and manipulation library for Go that provides:

- **Three API levels** to suit different use cases
- **Full RFC 8259 JSON compliance**
- **JSONPath query engine** (RFC 9535 compliant)
- **Streaming parser** for large files with constant memory usage
- **encoding/json compatible API** for easy migration
- **Fluent DOM API** for intuitive JSON manipulation
- **Zero external dependencies** (except Shape infrastructure)
- **Pure implementation** - does NOT use encoding/json internally
- **Proper type distinction** - Empty arrays `[]` and empty objects `{}` properly distinguished
- **Full round-trip fidelity** - Type information preserved through marshal/unmarshal cycles
- **High test coverage** with comprehensive test suite

### Key Features

- **Fluent DOM API**: User-friendly JSON manipulation (Recommended)
  - Type-safe getters (`GetString`, `GetInt`, `GetBool`) - No type assertions!
  - Fluent builder pattern with method chaining
  - Clean array semantics
  - Pretty-printing support

- **encoding/json Compatible API**: Drop-in replacement for standard library
  - `Marshal()` / `Unmarshal()` - Convert between Go structs and JSON
  - `MarshalIndent()` / `Indent()` / `Compact()` - JSON formatting
  - `Encoder` / `Decoder` - Streaming JSON I/O
  - Full struct tag support
  - **Pure implementation** - uses custom parser/renderer, not encoding/json
  - **Full type fidelity** - Empty arrays stay arrays, empty objects stay objects

- **JSON Validation**: Idiomatic error-based validation
- **JSONPath Query Engine**: RFC 9535-compliant implementation
- **Streaming Parser**: Constant memory (~20KB) for files of any size
- **Shape AST Integration**: Unified AST for cross-format operations

---

## Installation

```bash
go get github.com/shapestone/shape-json
```

**Requirements:**
- Go 1.25 or later
- No external dependencies (except Shape infrastructure)

---

## Quick Start

### Parsing JSON (DOM API - Recommended)

```go
import "github.com/shapestone/shape-json/pkg/json"

// Parse JSON into a Document
doc, err := json.ParseDocument(`{"name":"Alice","age":30}`)
if err != nil {
    log.Fatal(err)
}

// Access values with type-safe getters (no type assertions!)
name, _ := doc.GetString("name")  // "Alice"
age, _ := doc.GetInt("age")       // 30
```

### Building JSON (DOM API - Recommended)

```go
// Build JSON with fluent chaining
doc := json.NewDocument().
    SetString("name", "Alice").
    SetInt("age", 30).
    SetBool("active", true).
    SetArray("tags", json.NewArray().
        AddString("go").
        AddString("json"))

jsonStr, _ := doc.JSON()
// {"name":"Alice","age":30,"active":true,"tags":["go","json"]}
```

### Using encoding/json Compatible API

```go
import "github.com/shapestone/shape-json/pkg/json"

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// Marshal
person := Person{Name: "Alice", Age: 30}
data, _ := json.Marshal(person)

// Unmarshal
var decoded Person
json.Unmarshal(data, &decoded)
```

### Empty Arrays and Objects (Full Type Fidelity)

```go
import "github.com/shapestone/shape-json/pkg/json"

// Empty arrays and objects are properly distinguished
data := map[string]interface{}{
    "emptyArray":  []interface{}{},
    "emptyObject": map[string]interface{}{},
    "items":       []int{1, 2, 3},
}

// Marshal - types preserved
jsonBytes, _ := json.Marshal(data)
// {"emptyArray":[],"emptyObject":{},"items":[1,2,3]}

// Unmarshal - types preserved
var result map[string]interface{}
json.Unmarshal(jsonBytes, &result)

// Types are correctly maintained
emptyArr := result["emptyArray"].([]interface{})   // ✅ []interface{}
emptyObj := result["emptyObject"].(map[string]interface{}) // ✅ map[string]interface{}
// Full round-trip fidelity - no type information lost!
```

---

## API Documentation

### encoding/json Compatible API

shape-json provides a complete drop-in replacement for Go's standard `encoding/json`.

**Implementation Note:** While the API is compatible, shape-json uses its own custom parser and renderer internally—it does NOT use encoding/json. This ensures architectural purity and enables future enhancements like streaming support across all APIs.

#### Marshal / Unmarshal

```go
import "github.com/shapestone/shape-json/pkg/json"

type User struct {
    PublicName  string `json:"name"`              // Rename field
    Password    string `json:"-"`                 // Skip field
    Email       string `json:"email,omitempty"`   // Omit if empty
    Count       int    `json:"count,string"`      // Marshal as string
}

user := User{PublicName: "Alice", Email: "alice@example.com", Count: 42}

// Marshal to JSON
data, err := json.Marshal(user)
// {"name":"Alice","email":"alice@example.com","count":"42"}

// Unmarshal from JSON
var decoded User
err = json.Unmarshal(data, &decoded)
```

#### Pretty Printing

```go
// MarshalIndent - pretty print with indentation
person := Person{Name: "Alice", Age: 30}
pretty, _ := json.MarshalIndent(person, "", "  ")
// {
//   "age": 30,
//   "name": "Alice"
// }
// Note: Keys are sorted alphabetically

// Compact - remove whitespace
var compact bytes.Buffer
json.Compact(&compact, prettyJSON)

// Indent - add indentation to existing JSON
var indented bytes.Buffer
json.Indent(&indented, compactJSON, "", "  ")
```

#### Streaming I/O

```go
// Decoder - read from io.Reader
file, _ := os.Open("data.json")
decoder := json.NewDecoder(file)
var result Person
err := decoder.Decode(&result)

// Encoder - write to io.Writer
encoder := json.NewEncoder(os.Stdout)
err = encoder.Encode(person)
```

### Fluent DOM API (Recommended)

The DOM API provides a user-friendly, type-safe way to build and manipulate JSON:

#### Building Documents

```go
// Create a new document
doc := json.NewDocument()

// Set values with type-safe setters
doc.SetString("name", "Alice")
doc.SetInt("age", 30)
doc.SetBool("active", true)
doc.SetFloat("score", 98.5)
doc.SetNull("metadata")

// Nest objects
doc.SetObject("address", json.NewDocument().
    SetString("city", "NYC").
    SetString("zip", "10001"))

// Add arrays
doc.SetArray("tags", json.NewArray().
    AddString("go").
    AddString("json").
    AddInt(2025))

// Method chaining
doc := json.NewDocument().
    SetString("name", "Alice").
    SetInt("age", 30).
    SetBool("active", true)
```

#### Accessing Values

```go
// Parse JSON
doc, _ := json.ParseDocument(`{
    "user": {
        "name": "Bob",
        "age": 25,
        "roles": ["admin", "user"]
    }
}`)

// Get nested object
user, _ := doc.GetObject("user")

// Type-safe getters (no type assertions!)
name, _ := user.GetString("name")  // "Bob"
age, _ := user.GetInt("age")       // 25

// Get arrays
roles, _ := user.GetArray("roles")
role0, _ := roles.GetString(0)     // "admin"
```

#### Working with Arrays

```go
// Create array
arr := json.NewArray()

// Add values
arr.AddString("apple")
arr.AddInt(42)
arr.AddBool(true)
arr.AddObject(json.NewDocument().SetString("key", "value"))

// Access values
fruit, _ := arr.GetString(0)   // "apple"
num, _ := arr.GetInt(1)        // 42
flag, _ := arr.GetBool(2)      // true
obj, _ := arr.GetObject(3)

// Get length
length := arr.Len()  // 4
```

#### JSON Output

```go
// Compact JSON
jsonStr, _ := doc.JSON()
// {"name":"Alice","age":30}

// Pretty-printed JSON
pretty, _ := doc.JSONIndent("", "  ")
// {
//   "name": "Alice",
//   "age": 30
// }

// Custom indentation
custom, _ := doc.JSONIndent(">>", "\t")
```

### JSON Validation

Validate JSON before parsing (idiomatic Go error handling):

```go
import "github.com/shapestone/shape-json/pkg/json"

// Validate from string
if err := json.Validate(`{"name": "Alice"}`); err != nil {
    fmt.Println("Invalid JSON:", err)
    return
}
// err == nil means valid JSON

// Validate from io.Reader
file, _ := os.Open("data.json")
defer file.Close()
if err := json.ValidateReader(file); err != nil {
    fmt.Println("Invalid JSON:", err)
    return
}
// err == nil means valid JSON
```

### Streaming Parser

For large files, use `ParseReader()` to maintain constant memory usage:

```go
import "github.com/shapestone/shape-json/pkg/json"

// Parse large file with constant memory (~20KB)
file, _ := os.Open("huge-dataset.json")
defer file.Close()

node, err := json.ParseReader(file)
if err != nil {
    log.Fatal(err)
}

// Works with any io.Reader
resp, _ := http.Get("https://api.example.com/data")
defer resp.Body.Close()
node, _ := json.ParseReader(resp.Body)
```

**When to use:**
- **Parse()**: Small to medium files (<10MB) - faster, uses memory equal to file size
- **ParseReader()**: Large files (>10MB) - constant memory (~20KB), slower but scalable

**Performance:**
- 138MB file: Parse() uses 138MB RAM, ParseReader() uses 20KB RAM
- Memory savings: **99.985%** for large files
- Parse time: ParseReader() is ~20x slower but acceptable for streaming use cases

### Low-Level AST API

For advanced use cases requiring direct AST manipulation:

```go
import (
    "github.com/shapestone/shape/pkg/ast"
    "github.com/shapestone/shape-json/pkg/json"
)

// Parse into AST
node, err := json.Parse(`{"name": "Alice", "age": 30}`)
if err != nil {
    log.Fatal(err)
}

// Work with AST nodes
obj := node.(*ast.ObjectNode)
nameNode, _ := obj.GetProperty("name")
name := nameNode.(*ast.LiteralNode).Value().(string)  // "Alice"
```

---

## JSONPath Query Engine

shape-json includes a complete RFC 9535-compliant JSONPath implementation for querying JSON data:

### Basic Usage

```go
import "github.com/shapestone/shape-json/pkg/jsonpath"

// Parse a JSONPath query
expr, err := jsonpath.ParseString("$.store.book[*].author")
if err != nil {
    log.Fatal(err)
}

// Execute against data
data := map[string]interface{}{
    "store": map[string]interface{}{
        "book": []interface{}{
            map[string]interface{}{"author": "Nigel Rees", "price": 8.95},
            map[string]interface{}{"author": "Evelyn Waugh", "price": 12.99},
        },
    },
}

results := expr.Get(data)
// Results: ["Nigel Rees", "Evelyn Waugh"]
```

### Supported Features

#### Selectors

```go
// Root selector
"$"                           // Root of document

// Child selectors
"$.store.book"                // Dot notation
"$['store']['book']"          // Bracket notation

// Wildcards
"$.store.*"                   // All properties
"$.store.book[*]"             // All array elements

// Array indexing
"$.store.book[0]"             // First element
"$.store.book[-1]"            // Last element

// Array slicing
"$.store.book[0:2]"           // First two elements
"$.store.book[:5]"            // First five
"$.store.book[2:]"            // From third to end

// Recursive descent
"$..author"                   // All 'author' fields anywhere
```

#### Filter Expressions

```go
// Comparison operators
"$.store.book[?(@.price < 10)]"              // Books under $10
"$.store.book[?(@.price >= 10)]"             // Books $10 or more
"$.users[?(@.age > 18)]"                     // Adults only
"$.products[?(@.stock == 0)]"                // Out of stock

// Logical operators
"$.users[?(@.age > 18 && @.role == 'admin')]"   // Adult admins
"$.items[?(@.priority == 'high' || @.urgent)]"  // High priority or urgent

// Regex matching
"$.users[?(@.email =~ /.*@example\\.com$/)]"    // Example.com emails

// Field existence
"$.users[?(@.email)]"                           // Has email field
```

### Complete Example

```go
data := map[string]interface{}{
    "store": map[string]interface{}{
        "book": []interface{}{
            map[string]interface{}{
                "category": "reference",
                "author":   "Nigel Rees",
                "title":    "Sayings of the Century",
                "price":    8.95,
            },
            map[string]interface{}{
                "category": "fiction",
                "author":   "Evelyn Waugh",
                "title":    "Sword of Honour",
                "price":    12.99,
            },
        },
    },
}

// Find all authors
expr, _ := jsonpath.ParseString("$.store.book[*].author")
authors := expr.Get(data)
// ["Nigel Rees", "Evelyn Waugh"]

// Find books under $10
expr, _ = jsonpath.ParseString("$.store.book[?(@.price < 10)]")
cheapBooks := expr.Get(data)
// [{"category":"reference", "author":"Nigel Rees", ...}]

// Find all prices recursively
expr, _ = jsonpath.ParseString("$..price")
prices := expr.Get(data)
// [8.95, 12.99]
```

For complete JSONPath documentation, see [pkg/jsonpath/README.md](pkg/jsonpath/README.md).

---

## Performance & Benchmarks

### Dual-Path Architecture

shape-json automatically selects the optimal parsing strategy based on the API you use:

#### ⚡ Fast Path (Default - Optimized for Speed)

The fast path bypasses AST construction for maximum performance.

**When it's used (automatic)**:
- `Unmarshal()` - Parse JSON into Go structs or `interface{}`
- `Validate()` / `ValidateReader()` - Syntax validation only

**Performance benefits**:
- **9.7x faster** than AST path (1.69 ms vs 16.45 ms for 410 KB file)
- **10.2x less memory** (1.29 MB vs 13.19 MB)
- **6.3x fewer allocations** (36,756 vs 230,025)

**Example - Fast path (recommended for most use cases)**:
```go
// Unmarshal into struct (fast path - 9.7x faster!)
var person Person
err := json.Unmarshal(data, &person)

// Unmarshal into interface{} (fast path)
var result interface{}
err = json.Unmarshal(data, &result)

// Validate syntax only (fast path - 8.2x faster!)
if err := json.Validate(jsonString); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}
```

#### 🌳 AST Path (Full Features)

The AST path builds a complete Abstract Syntax Tree with position tracking.

**When it's used**:
- `Parse()` / `ParseReader()` - Get AST for advanced features
- `ParseDocument()` / `ParseArray()` - DOM API (uses Parse internally)
- JSONPath queries (require AST)

**Use when you need**:
- JSONPath queries (`$.users[?(@.age > 30)]`)
- Tree manipulation/transformation
- Position tracking for error reporting
- Format conversion (JSON ↔ XML ↔ YAML)

**Example - AST path for advanced features**:
```go
// Get AST for JSONPath queries
node, err := json.Parse(jsonString)
if err != nil {
    return err
}

// Run JSONPath query (requires AST)
results, err := jsonpath.Query(node, "$.users[?(@.age > 30)].name")

// Or use DOM API
doc, err := json.ParseDocument(jsonString)
user, _ := doc.GetObject("user")
name, _ := user.GetString("name")
```

#### Choosing the Right API

**Decision guide**:

| Need | Use | Path | Performance |
|------|-----|------|-------------|
| Unmarshal into Go struct | `Unmarshal()` | Fast | 9.7x faster |
| Validate JSON syntax | `Validate()` | Fast | 8.2x faster |
| Parse into `interface{}` | `Unmarshal()` | Fast | 9.7x faster |
| JSONPath queries | `Parse()` + JSONPath | AST | Full features |
| Tree manipulation | `Parse()` | AST | Full features |
| DOM API | `ParseDocument()` | AST | User-friendly |

**Performance tip**: If you don't need the AST, always use `Unmarshal()` instead of `Parse()` followed by `NodeToInterface()`:

```go
// ❌ Slow (builds AST unnecessarily)
node, _ := json.Parse(string(data))
value := json.NodeToInterface(node)
person := value.(map[string]interface{})

// ✅ Fast (9.7x faster - uses fast path)
var person map[string]interface{}
json.Unmarshal(data, &person)
```

**When you need both speed AND the AST**:

If you need to unmarshal AND run JSONPath queries, use the fast path first, then Parse if needed:

```go
// Fast unmarshal first (9.7x faster)
var data MyStruct
json.Unmarshal(bytes, &data)

// Only parse to AST if you need JSONPath
if needsQuery {
    node, _ := json.Parse(string(bytes))
    results, _ := jsonpath.Query(node, "$.items[*].price")
}
```

### Streaming Parser Performance

**Memory Usage:**
- `Parse()`: Memory usage scales with file size (O(n))
- `ParseReader()`: Constant memory usage regardless of file size
- For large files (>10MB): Memory savings exceed 99%

**Performance Characteristics:**
- `Parse()`: Faster parsing, higher memory usage
- `ParseReader()`: Slower parsing, constant low memory usage
- Trade-off: ParseReader is approximately 20x slower but uses constant memory

**Key Features:**
- Constant memory usage for files of any size
- Works with any io.Reader (files, network streams, pipes)
- Full UTF-8 support across buffer boundaries
- Position tracking for accurate error messages

**Benchmarking:**
Run `make benchmark` to see performance on your system with your data.

### Test Coverage

Run `make coverage` to generate a current coverage report.

**Coverage includes:**
- All API surfaces (Marshal/Unmarshal, DOM, Format, Parse)
- Edge cases and error conditions
- Unicode handling and escape sequences
- Struct tag parsing and field handling
- Stream parsing and buffer management

### Fuzzing

Continuous fuzzing is performed on:
- Parser (objects, arrays, strings, numbers, nested structures)
- JSONPath (filters, slices, recursive descent, complex queries)

Run fuzzing tests:
```bash
go test ./internal/parser -fuzz=FuzzParserObjects -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathFilters -fuzztime=30s
```

---

## Architecture

### Design Principles

- **Grammar-Driven**: EBNF grammar in `docs/grammar/json.ebnf`
- **LL(1) Recursive Descent Parser**: Hand-coded, optimized parser
- **Single Token Lookahead**: Efficient parsing with minimal overhead
- **Unified Architecture**: Single parser for all APIs (no encoding/json dependency)
- **Zero External Dependencies**: Uses only Go standard library + Shape infrastructure

### AST Representation

shape-json uses Shape's universal AST with proper type distinction:

- **Objects** → `*ast.ObjectNode` with properties map
- **Arrays** → `*ast.ArrayDataNode` with elements slice
- **Primitives** → `*ast.LiteralNode` (string, int64, float64, bool, nil)

**Key Feature:** Arrays and objects are distinct types at the AST level, ensuring:
- Empty arrays `[]` are distinguishable from empty objects `{}`
- Full round-trip fidelity: `[]` → Marshal → Unmarshal → `[]`
- Type-safe operations on arrays vs objects

### Unified Architecture

All APIs use the same underlying parser and renderer:

```
┌────────────────────────────────────────────────────┐
│           User-Facing APIs                         │
├─────────────┬──────────────┬───────────────────────┤
│ Marshal/    │  Fluent DOM  │  Low-Level AST API    │
│ Unmarshal   │  API         │  Parse/ParseReader    │
└─────────────┴──────────────┴───────────────────────┘
                         │
         ┌───────────────┴───────────────┐
         ▼                               ▼
┌─────────────────┐            ┌─────────────────┐
│  Render Path    │            │   Parse Path    │
│  (Go → JSON)    │            │   (JSON → Go)   │
└────────┬────────┘            └────────┬────────┘
         │                              │
         ▼                              ▼
┌─────────────────────────────────────────────────┐
│       Type Conversion Layer                     │
│   NodeToInterface() / InterfaceToNode()         │
└────────┬────────────────────────┬───────────────┘
         │                        │
         ▼                        ▼
┌──────────────┐         ┌──────────────┐
│  Renderer    │         │   Parser     │
│ (render.go)  │         │ (parser.go)  │
└──────┬───────┘         └──────┬───────┘
       │                        │
       └───────────┬────────────┘
                   ▼
          ┌─────────────────┐
          │  Unified AST    │
          │ ast.SchemaNode  │
          └─────────────────┘
```

**Data Flow:**
- **Marshal**: `Go struct → InterfaceToNode() → AST → Render() → JSON bytes`
- **Unmarshal**: `JSON bytes → Parse() → AST → NodeToInterface() → Go struct`
- **DOM**: `JSON → Parse() → AST → NodeToInterface() → map/slice → Document/Array`
- **Format**: `JSON → Parse() → AST → RenderIndent() → Pretty JSON`

### Streaming Architecture

```
io.Reader → NewStreamFromReader()
              ↓
         bufferedStreamImpl
         (64KB sliding window)
              ↓
         json.ParseReader()
              ↓
         ast.SchemaNode
```

**Buffer Management:**
- 64KB sliding window
- 8KB safety margin for backtracking
- Automatic buffer compaction
- UTF-8 boundary handling

---

## Testing & Quality

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
make coverage

# Run fuzzing (5 seconds each)
go test ./internal/parser -fuzz=FuzzParserObjects -fuzztime=5s

# Run grammar verification
make grammar-test
```

### Test Organization

- **Unit tests**: 45+ tests covering all APIs
- **Fuzzing tests**: Parser and JSONPath fuzzers
- **Integration tests**: Real-world file processing
- **Example tests**: Executable documentation
- **Benchmark tests**: Performance measurement

### Quality Metrics

- **Test Coverage**: Comprehensive coverage across all components (run `make coverage`)
- **Fuzzing**: Continuous fuzzing with zero crashes
- **Grammar Verification**: 100% EBNF compliance
- **Static Analysis**: CodeQL, OpenSSF Scorecard
- **Security**: Regular security scanning
- **Architectural Purity**: Zero encoding/json dependency in production code

---

## Examples

The `examples/` directory contains comprehensive usage examples:

### DOM Builder API (Recommended)

Build and manipulate JSON with a fluent, type-safe API:

```bash
go run examples/dom_builder/main.go
```

See [examples/dom_builder/README.md](examples/dom_builder/README.md) for detailed documentation.

### JSON Validation

Validate JSON input before parsing:

```bash
go run examples/format_detection/main.go
```

See [examples/format_detection/README.md](examples/format_detection/README.md) for detailed documentation.

### Streaming Parser

Parse large files with constant memory:

```bash
go run examples/parse_reader/main.go
```

### Other Examples

- [examples/main.go](examples/main.go) - Basic parsing with AST
- [examples/encoding_json_api/](examples/encoding_json_api/) - encoding/json compatibility

---

## Migration Guide

### From encoding/json

shape-json is a drop-in replacement. Simply change the import:

```go
// Before
import "encoding/json"

// After
import json "github.com/shapestone/shape-json/pkg/json"
```

All `Marshal()`, `Unmarshal()`, `Encoder`, `Decoder` functionality works identically.

#### Behavioral Differences

While shape-json maintains API compatibility, there are minor behavioral differences:

**1. Key Ordering**
```go
// encoding/json: preserves insertion order
{"name":"Alice","age":30}

// shape-json: alphabetical sorting for deterministic output
{"age":30,"name":"Alice"}
```

**2. Number Types**
```go
var data map[string]interface{}
json.Unmarshal([]byte(`{"age":30}`), &data)

// encoding/json: numbers are float64
age := data["age"].(float64)  // 30.0

// shape-json: whole numbers are int64
age := data["age"].(int64)     // 30
```

These differences are intentional design decisions for consistency and do not affect struct marshaling/unmarshaling.

### From shape core (JSONPath)

If you were using JSONPath from the shape core library:

```go
// Old import
import "github.com/shapestone/shape/pkg/jsonpath"

// New import
import "github.com/shapestone/shape-json/pkg/jsonpath"
```

Then update your `go.mod`:

```bash
go get github.com/shapestone/shape-json
```

All JSONPath functionality remains unchanged.

### Upgrading to Streaming Parser

For large file processing, migrate to `ParseReader()`:

```go
// Before (loads entire file into memory)
data, _ := os.ReadFile("large.json")
node, _ := json.Parse(string(data))

// After (constant memory usage)
file, _ := os.Open("large.json")
defer file.Close()
node, _ := json.ParseReader(file)
```

**Benefits:**
- 99.985% memory reduction for large files
- Handles files of any size
- Works with network streams

---

## Additional Resources

- **Go Package Documentation**: https://pkg.go.dev/github.com/shapestone/shape-json
- **Source Code**: https://github.com/shapestone/shape-json
- **Issue Tracker**: https://github.com/shapestone/shape-json/issues
- **Shape Ecosystem**: https://github.com/shapestone/shape
- **EBNF Grammar**: [docs/grammar/json.ebnf](docs/grammar/json.ebnf)
- **Testing Guide**: [docs/TESTING.md](docs/TESTING.md)

---

## License

Apache License 2.0

Copyright © 2020-2025 Shapestone

See [LICENSE](LICENSE) for the full license text and [NOTICE](NOTICE) for third-party attributions.

---

## Support

For questions, bug reports, or feature requests:
- Open an issue at https://github.com/shapestone/shape-json/issues
- Consult the examples in `examples/`
- Review the comprehensive test suite

---

## Recent Changes (v0.8.0)

### Unified Architecture

shape-json has undergone a major architectural refactor to eliminate all encoding/json usage from production code:

- **Pure Implementation**: All APIs now use a single custom parser and renderer
- **No encoding/json**: Zero dependency on Go's standard library JSON package internally
- **Consistent Behavior**: All APIs share the same parsing and rendering logic
- **Future Ready**: Enables features like streaming support across all APIs

**Benefits:**
- Architectural purity - true standalone JSON parser
- Deterministic output with alphabetical key sorting
- Consistent number handling (int64 for whole numbers)
- Foundation for future enhancements

**Minor Behavioral Changes:**
- JSON keys sorted alphabetically (was insertion order)
- Whole numbers unmarshal as int64 (was float64)

**Major Fix (v0.8.0):**
- Empty arrays now correctly render as `[]` instead of `{}`
- Arrays and objects properly distinguished at the AST level via new `ArrayDataNode` type
- Full round-trip fidelity for empty arrays

See the [Migration Guide](#migration-guide) for details on behavioral differences.

---

*Last Updated: December 17, 2025*
