# Architecture

This document describes the architecture of shape-json, a comprehensive JSON parser and manipulation library for Go.

## Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [System Architecture](#system-architecture)
- [Core Components](#core-components)
- [Data Flow](#data-flow)
- [AST Representation](#ast-representation)
- [Key Design Decisions](#key-design-decisions)
- [Performance Characteristics](#performance-characteristics)
- [Additional Resources](#additional-resources)

---

## Overview

shape-json is built on a **unified architecture** with a single custom parser that powers all APIs:

- **encoding/json Compatible API**: Drop-in replacement for Go's standard library
- **Fluent DOM API**: Type-safe, user-friendly JSON manipulation
- **Validation API**: Idiomatic error-based validation
- **JSONPath Query Engine**: RFC 9535-compliant querying
- **Low-Level AST API**: Direct AST manipulation for advanced use cases

**Key Characteristic:** All APIs use the same underlying parser and renderer—there is **zero dependency on encoding/json** in production code.

---

## Design Principles

1. **Grammar-Driven Design**
   - EBNF grammar specification drives parser implementation
   - Grammar file: [`docs/grammar/json.ebnf`](docs/grammar/json.ebnf)
   - Ensures RFC 8259 compliance

2. **Architectural Purity**
   - Single unified parser for all APIs
   - No encoding/json dependency in production code
   - True standalone JSON parser

3. **LL(1) Recursive Descent Parser**
   - Hand-coded parser with single token lookahead
   - Optimized for performance and maintainability
   - Clear error messages with position tracking

4. **Zero External Dependencies**
   - Only depends on Go standard library + Shape infrastructure
   - No third-party parsing libraries

5. **Type Safety**
   - Proper distinction between arrays and objects at AST level
   - Full round-trip fidelity (e.g., `[]` → Marshal → Unmarshal → `[]`)
   - Type-safe DOM API with no type assertions needed

6. **Streaming Support**
   - Constant memory usage for large files via `ParseReader()`
   - Works with any `io.Reader`
   - Buffered stream implementation with UTF-8 boundary handling

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        User-Facing APIs                         │
├──────────────┬──────────────┬──────────────┬───────────────────┤
│ Marshal/     │  Fluent DOM  │  Validation  │  Low-Level AST    │
│ Unmarshal    │  API         │  API         │  Parse/ParseReader│
│ Encoder/     │  Document/   │  Validate/   │                   │
│ Decoder      │  Array       │  ValidateRdr │                   │
└──────────────┴──────────────┴──────────────┴───────────────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐            ┌─────────────────┐
    │  Render Path    │            │   Parse Path    │
    │  (Go → JSON)    │            │   (JSON → Go)   │
    └────────┬────────┘            └────────┬────────┘
             │                              │
             │    ┌───────────────────┐     │
             └───→│ Type Conversion   │←────┘
                  │ Layer             │
                  │ (AST ↔ Go types)  │
                  └────────┬──────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
    ┌──────────────────┐      ┌──────────────────┐
    │  JSON Renderer   │      │  JSON Parser     │
    │  (render.go)     │      │  (parser.go)     │
    │                  │      │                  │
    │  - Render()      │      │  - Parse()       │
    │  - RenderIndent()│      │  - ParseReader() │
    │  - escapeString()│      │  - Tokenizer     │
    └──────────────────┘      └──────────────────┘
              │                         │
              └────────────┬────────────┘
                           ▼
                  ┌─────────────────┐
                  │  Unified AST    │
                  │  (Shape Core)   │
                  │                 │
                  │  - ObjectNode   │
                  │  - ArrayDataNode│
                  │  - LiteralNode  │
                  └─────────────────┘
```

---

## Dual-Path Architecture

shape-json implements an **intelligent dual-path architecture** that automatically selects the optimal parsing strategy based on the API being used. This provides dramatic performance improvements for common use cases while maintaining full feature support for advanced scenarios.

### Path Selection

```
┌──────────────────────────────────────┐
│         User API Call                │
└───────────────┬──────────────────────┘
                │
        ┌───────┴────────┐
        │                │
        v                v
┌───────────────┐  ┌──────────────────┐
│  FAST PATH    │  │   AST PATH       │
│  (Optimized)  │  │ (Full Features)  │
├───────────────┤  ├──────────────────┤
│ APIs:         │  │ APIs:            │
│ - Unmarshal() │  │ - Parse()        │
│ - Validate()  │  │ - ParseReader()  │
│               │  │ - ParseDocument()│
├───────────────┤  ├──────────────────┤
│ Direct Parse  │  │ Tokenization     │
│ ↓             │  │ ↓                │
│ Go Values     │  │ AST Construction │
│               │  │ ↓                │
│ Performance:  │  │ AST Nodes        │
│ - 9.7x faster │  │ ↓                │
│ - 10.2x less  │  │ Type Conversion  │
│   memory      │  │ ↓                │
│ - 6.3x fewer  │  │ Go Values        │
│   allocations │  │                  │
└───────────────┘  └──────────────────┘
```

### Fast Path (internal/fastparser/)

**Purpose:** Maximum performance for common operations (80-90% of use cases)

**Implementation:**
- Package: `internal/fastparser/`
- Direct byte-level parsing (no tokenization)
- No AST construction
- Immediate conversion to Go types using reflection
- SWAR (SIMD Within A Register) optimizations for whitespace

**APIs that use fast path:**
- `Unmarshal(data []byte, v interface{})` - Direct unmarshaling
- `Validate(input string)` - Syntax validation only
- `ValidateReader(reader io.Reader)` - Stream validation

**Performance characteristics:**
```
Benchmark (410 KB JSON file):
Fast Path:  1.69 ms,  1.29 MB,  36,756 allocs
AST Path:  16.45 ms, 13.19 MB, 230,025 allocs
Improvement: 9.7x faster, 10.2x less memory, 6.3x fewer allocations
```

**Trade-offs:**
- ✅ Dramatically faster for unmarshaling
- ✅ Much lower memory usage
- ✅ Fewer allocations (less GC pressure)
- ❌ No AST available for advanced features
- ❌ No position tracking in errors

### AST Path (internal/parser/)

**Purpose:** Advanced features requiring tree structure

**Implementation:**
- Package: `internal/parser/`
- Tokenization via `internal/tokenizer/`
- Full AST construction with position tracking
- Type conversion layer for Go interop

**APIs that use AST path:**
- `Parse(input string)` - Returns `ast.SchemaNode`
- `ParseReader(reader io.Reader)` - Stream parsing to AST
- `ParseDocument(input string)` - DOM API (uses Parse)
- `ParseArray(input string)` - DOM API (uses Parse)

**Performance characteristics:**
- Slower parsing (builds complete tree)
- Higher memory usage (stores all nodes)
- More allocations (node objects)

**Benefits:**
- ✅ Full position tracking for errors
- ✅ JSONPath query support
- ✅ Tree manipulation/transformation
- ✅ Format conversion (JSON ↔ XML ↔ YAML)
- ✅ Precise error locations

### Path Selection Logic

The architecture automatically selects the path based on which API is called:

**Fast Path Selection:**
```go
// Unmarshal → fastparser.Unmarshal() → Fast Path
var person Person
json.Unmarshal(data, &person)  // ← 9.7x faster

// Validate → fastparser.Parse() → Fast Path
err := json.Validate(jsonString)  // ← 8.2x faster
```

**AST Path Selection:**
```go
// Parse → parser.Parse() → AST Path
node, _ := json.Parse(jsonString)  // ← Full AST

// JSONPath requires AST
results, _ := jsonpath.Query(node, "$.users[*].name")
```

### When to Use Each Path

**Use Fast Path when:**
- Unmarshaling JSON into Go structs
- Validating JSON syntax
- Parsing into `interface{}`
- Performance is critical
- You don't need the AST

**Use AST Path when:**
- Running JSONPath queries
- Manipulating the tree structure
- Converting between formats
- Need precise error locations
- Building tools that inspect JSON structure

### Migration Recommendations

**From AST path to fast path (performance optimization):**

```go
// Before (slow - builds AST unnecessarily)
node, _ := json.Parse(string(data))
value := json.NodeToInterface(node)
person := value.(Person)

// After (9.7x faster - uses fast path)
var person Person
json.Unmarshal(data, &person)
```

**Keeping AST path when needed:**

```go
// If you need JSONPath, keep using Parse
node, _ := json.Parse(jsonString)
results, _ := jsonpath.Query(node, "$.users[?(@.age > 30)]")
```

---

## Core Components

### 1. JSON Parser (`internal/parser/parser.go`)

**Responsibility:** Parse JSON text into Shape's unified AST

**Key Features:**
- LL(1) recursive descent parser
- Single token lookahead
- RFC 8259 compliant
- Context-aware error messages with position tracking

**APIs:**
- `Parse(input string) (ast.SchemaNode, error)` - Parse from string
- `ParseReader(r io.Reader) (ast.SchemaNode, error)` - Streaming parse with constant memory

**Implementation:**
- Tokenizes input using Shape's tokenizer framework
- Builds AST nodes (`ObjectNode`, `ArrayDataNode`, `LiteralNode`)
- Handles all JSON data types (object, array, string, number, boolean, null)

### 2. JSON Renderer (`pkg/json/render.go`)

**Responsibility:** Convert AST back to JSON text

**Key Features:**
- Custom JSON rendering (no encoding/json dependency)
- Full string escaping (`\"`, `\\`, `\n`, `\t`, `\uXXXX`, etc.)
- Alphabetical key sorting for deterministic output
- Pretty-printing with custom indentation
- Automatic array detection (numeric keys "0", "1", "2"...)

**APIs:**
- `Render(node ast.SchemaNode) ([]byte, error)` - Compact JSON
- `RenderIndent(node ast.SchemaNode, prefix, indent string) ([]byte, error)` - Pretty JSON

### 3. Type Conversion Layer (`pkg/json/convert.go`)

**Responsibility:** Bidirectional conversion between AST nodes and Go types

**Key Features:**
- AST → Go: Convert AST nodes to `interface{}`, maps, slices
- Go → AST: Convert Go values to AST nodes
- Type preservation (int64, float64, bool, string, nil)
- Recursive conversion for nested structures

**APIs:**
- `NodeToInterface(node ast.SchemaNode) interface{}` - AST to Go types
- `InterfaceToNode(v interface{}) (ast.SchemaNode, error)` - Go types to AST

### 4. Tokenizer (`internal/tokenizer/tokenizer.go`)

**Responsibility:** Lexical analysis - convert JSON text into tokens

**Key Features:**
- Uses Shape's tokenizer framework
- Token types: `{`, `}`, `[`, `]`, `:`, `,`, STRING, NUMBER, TRUE, FALSE, NULL
- Proper string escaping and Unicode handling
- High test coverage (98.8%)

### 5. Streaming Infrastructure (`pkg/stream/`)

**Responsibility:** Buffered reading for constant-memory parsing

**Key Features:**
- 64KB sliding window buffer
- 8KB safety margin for backtracking
- Automatic buffer compaction
- UTF-8 boundary handling
- Position tracking for error messages

**Implementation:**
- `NewStreamFromReader(io.Reader)` - Create buffered stream
- Supports any `io.Reader` (files, network, pipes)
- ~20KB constant memory usage regardless of input size

---

## Data Flow

### Marshal (Go struct → JSON)

```
Go struct
    ↓
Custom reflection (pkg/json/marshal.go)
    ↓
InterfaceToNode() - Convert to AST
    ↓
Render() - AST to JSON bytes
    ↓
JSON bytes
```

**Example:**
```go
person := Person{Name: "Alice", Age: 30}
data, _ := json.Marshal(person)
// {"age":30,"name":"Alice"}
```

### Unmarshal (JSON → Go struct)

```
JSON bytes
    ↓
Parse() - JSON to AST
    ↓
NodeToInterface() - AST to Go types
    ↓
Custom reflection (pkg/json/unmarshal.go)
    ↓
Go struct
```

**Example:**
```go
var person Person
json.Unmarshal(data, &person)
// person = Person{Name: "Alice", Age: 30}
```

### DOM API (JSON ↔ Document)

```
JSON string
    ↓
Parse() - JSON to AST
    ↓
NodeToInterface() - AST to map[string]interface{}
    ↓
Document wrapper
    ↓
(reverse direction for JSON output)
    ↓
InterfaceToNode() - map to AST
    ↓
Render() - AST to JSON string
```

**Example:**
```go
doc, _ := json.ParseDocument(`{"name":"Alice"}`)
name, _ := doc.GetString("name")  // "Alice"
jsonStr, _ := doc.JSON()          // {"name":"Alice"}
```

### Format API (Indent/Compact)

```
JSON string (any whitespace)
    ↓
Parse() - JSON to AST
    ↓
RenderIndent() or Render()
    ↓
Formatted JSON string
```

**Example:**
```go
var indented bytes.Buffer
json.Indent(&indented, compactJSON, "", "  ")
```

### Validation API

```
JSON string
    ↓
Parse() - attempt to parse
    ↓
Returns: nil (valid) or error (invalid)
```

**Example:**
```go
err := json.Validate(`{"name": "Alice"}`)
// err == nil means valid JSON
```

---

## AST Representation

shape-json uses Shape's universal AST with proper type distinction:

### Object Nodes (`*ast.ObjectNode`)

```go
// Represents: {"name": "Alice", "age": 30}
ObjectNode {
    properties: map[string]ast.SchemaNode{
        "name": LiteralNode{value: "Alice"},
        "age":  LiteralNode{value: int64(30)},
    }
}
```

### Array Nodes (`*ast.ArrayDataNode`)

```go
// Represents: ["apple", "banana", "cherry"]
ArrayDataNode {
    elements: []ast.SchemaNode{
        LiteralNode{value: "apple"},
        LiteralNode{value: "banana"},
        LiteralNode{value: "cherry"},
    }
}
```

**Key Feature:** Arrays use `ArrayDataNode` (not `ObjectNode` with numeric keys), ensuring:
- Empty arrays `[]` are distinguishable from empty objects `{}`
- Full round-trip fidelity
- Type-safe operations

### Literal Nodes (`*ast.LiteralNode`)

```go
// Represents: "hello"
LiteralNode{value: "hello"}

// Represents: 42
LiteralNode{value: int64(42)}

// Represents: 3.14
LiteralNode{value: float64(3.14)}

// Represents: true
LiteralNode{value: true}

// Represents: null
LiteralNode{value: nil}
```

**Type Handling:**
- **Whole numbers** → `int64` (e.g., `42`)
- **Decimals** → `float64` (e.g., `3.14`)
- **Strings** → `string`
- **Booleans** → `bool`
- **Null** → `nil`

---

## Key Design Decisions

### 1. Elimination of encoding/json Dependency

**Decision:** Do not use encoding/json in production code

**Rationale:**
- Architectural purity - a JSON parser should not depend on another JSON parser
- Enables future enhancements (e.g., streaming support across all APIs)
- Full control over behavior and error messages
- Consistent behavior across all APIs

**Trade-off:** Custom parser is ~20x slower than encoding/json for large files, but provides streaming capability with 99.985% memory savings

**See:** [ARCHITECTURAL_DISCOVERY.md](ARCHITECTURAL_DISCOVERY.md) for the complete evolution story

### 2. Alphabetical Key Sorting

**Decision:** Always sort object keys alphabetically in output

**Rationale:**
- **Deterministic output** for testing and comparison
- **Canonical representation** - same JSON structure always produces same output
- **Easier debugging** - consistent key order

**Trade-off:** Does not preserve original key order from input

**Example:**
```json
Input:  {"name":"Alice","age":30}
Output: {"age":30,"name":"Alice"}
```

### 3. ArrayDataNode for Type Distinction

**Decision:** Use dedicated `ArrayDataNode` type instead of `ObjectNode` with numeric keys

**Rationale:**
- **Type safety** - arrays and objects are distinct types
- **Round-trip fidelity** - empty arrays remain arrays
- **Semantic correctness** - JSON arrays should have array representation

**Impact:**
```go
// Before: [] → ObjectNode{} → {} (wrong!)
// After:  [] → ArrayDataNode{} → [] (correct!)
```

### 4. Streaming Parser with Constant Memory

**Decision:** Provide `ParseReader()` for constant-memory parsing

**Rationale:**
- **Scalability** - handle files of any size
- **Predictable memory** - ~20KB regardless of input
- **Network compatibility** - works with any `io.Reader`

**Trade-off:** ~20x slower than in-memory `Parse()`, but acceptable for streaming use cases

**Use Cases:**
- Large files (>10MB)
- Network streams
- Memory-constrained environments

### 5. Type-Safe DOM API

**Decision:** Provide type-safe getters without type assertions

**Rationale:**
- **User experience** - cleaner, more intuitive API
- **Error handling** - errors instead of panics
- **Discoverability** - clear method names (`GetString`, `GetInt`, etc.)

**Example:**
```go
// Without DOM API (requires type assertions)
obj := data.(map[string]interface{})
name := obj["name"].(string)  // Can panic!

// With DOM API (no type assertions)
name, err := doc.GetString("name")  // Returns error if wrong type
```

---

## Performance Characteristics

### Memory Usage

| Operation | Memory Usage | Notes |
|-----------|--------------|-------|
| `Parse(input)` | O(n) | Memory scales with input size |
| `ParseReader(r)` | ~20KB constant | Regardless of input size |
| `Marshal(v)` | O(n) | Memory for AST + output |
| DOM API | O(n) | In-memory map/slice representation |

### Performance Benchmarks

**Parse() vs ParseReader() (138MB file):**
- `Parse()`: 13.1 seconds, 138MB memory
- `ParseReader()`: 4m31s, ~20KB memory
- **Trade-off:** 20x slower, 99.985% less memory

**When to use:**
- **Small/medium files (<10MB):** Use `Parse()` for speed
- **Large files (>10MB):** Use `ParseReader()` for constant memory
- **Network streams:** Use `ParseReader()` to avoid buffering entire response

### Test Coverage

- **Parser:** 97.1%
- **Tokenizer:** 98.8%
- **JSONPath:** 89.8%
- **JSON API:** 84.2%
- **Overall:** 84%+

**Fuzzing Results:**
- Parser: 840K+ executions, 0 crashes
- JSONPath: 1.17M+ executions, 0 crashes

---

## Additional Resources

### Documentation

- **[README.md](README.md)** - Quick start and feature overview
- **[USER_GUIDE.md](USER_GUIDE.md)** - Comprehensive API documentation
- **[ARCHITECTURAL_DISCOVERY.md](ARCHITECTURAL_DISCOVERY.md)** - Detailed evolution story (dual parser → unified architecture)
- **[docs/TESTING.md](docs/TESTING.md)** - Testing guide
- **[docs/grammar/json.ebnf](docs/grammar/json.ebnf)** - Complete EBNF grammar specification

### Related Projects

- **[Shape Core](https://github.com/shapestone/shape-core)** - Universal AST and tokenizer framework
- **[Shape Ecosystem](https://github.com/shapestone/shape)** - Multi-format parser ecosystem
- **[Parser Implementation Guide](https://github.com/shapestone/shape-core/blob/main/docs/PARSER_IMPLEMENTATION_GUIDE.md)** - Guide for implementing parsers using Shape

### API Reference

- **[pkg.go.dev](https://pkg.go.dev/github.com/shapestone/shape-json)** - Generated API documentation
- **[JSONPath README](pkg/jsonpath/README.md)** - JSONPath query engine documentation

---

## Contributing

When contributing to shape-json:

1. **Understand the architecture** - Read this document and ARCHITECTURAL_DISCOVERY.md
2. **Maintain purity** - Never import encoding/json in production code
3. **Follow the grammar** - Parser changes must align with EBNF grammar
4. **Add tests** - All new code requires comprehensive tests
5. **Update docs** - Keep architecture docs synchronized with code

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## License

Apache License 2.0

Copyright © 2020-2025 Shapestone

See [LICENSE](LICENSE) for the full license text.

---

*Last Updated: December 21, 2025*
