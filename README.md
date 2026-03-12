# shape-json

A fast, drop-in JSON library for Go — replaces encoding/json with JSONPath queries, a fluent DOM builder, and a high-performance parser.

![Build Status](https://github.com/shapestone/shape-json/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shapestone/shape-json?v=2)](https://goreportcard.com/report/github.com/shapestone/shape-json)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![codecov](https://codecov.io/gh/shapestone/shape-json/branch/main/graph/badge.svg?v=3)](https://codecov.io/gh/shapestone/shape-json)
![Go Version](https://img.shields.io/github/go-mod/go-version/shapestone/shape-json)
![Latest Release](https://img.shields.io/github/v/release/shapestone/shape-json?v=2)
[![GoDoc](https://pkg.go.dev/badge/github.com/shapestone/shape-json.svg)](https://pkg.go.dev/github.com/shapestone/shape-json)

[![CodeQL](https://github.com/shapestone/shape-json/actions/workflows/codeql.yml/badge.svg)](https://github.com/shapestone/shape-json/actions/workflows/codeql.yml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/shapestone/shape-json/badge)](https://securityscorecards.dev/viewer/?uri=github.com/shapestone/shape-json)
[![Security Policy](https://img.shields.io/badge/Security-Policy-brightgreen)](SECURITY.md)

**Repository:** github.com/shapestone/shape-json

A JSON parser for the [Shape Parser™](https://github.com/shapestone/shape) ecosystem.

Parses JSON data (RFC 8259, the JSON internet standard) into Shape Parser's™ unified AST representation.

## What it does

shape-json is a complete JSON parser and serializer for Go programs. It parses JSON (as defined by RFC 8259, the JSON internet standard) into either Go structs or a traversable Abstract Syntax Tree (AST). It is a pure Go implementation — it does not use or wrap `encoding/json`.

It provides four APIs at different levels of abstraction: a drop-in `encoding/json` replacement, a fluent DOM (Document Object Model) builder, a JSONPath query engine, and a low-level AST API for parser tooling.

## Who it's for

- Go developers who want a drop-in replacement for `encoding/json` with better performance
- Developers who need to query JSON with JSONPath (e.g., `$.users[?(@.role == "admin")].name`)
- Tooling authors who need a traversable Abstract Syntax Tree (AST) from JSON input
- Developers building data pipelines that process JSON files, APIs, or config files

## Use Cases

- **API clients**: Unmarshal JSON responses from REST APIs into Go structs
- **Configuration loading**: Parse JSON config files with validation
- **Data transformation pipelines**: Use the AST path to read, modify, and re-serialize JSON
- **JSONPath queries**: Extract values from deeply nested JSON with RFC 9535 filter expressions
- **JSON validation services**: Use `Validate()` / `ValidateReader()` for input validation without full parsing
- **JSON builder tools**: Use the fluent DOM API to construct JSON from code without string templates
- **Parser/tooling authors**: Use the AST directly for static analysis, linting, or format conversion

## Installation

```bash
go get github.com/shapestone/shape-json
```

## Quick Start

```bash
# Install
go get github.com/shapestone/shape-json

# Run tests
make test

# Run tests with coverage report
make coverage

# Build
make build

# Run all checks (grammar, tests, lint, build, coverage)
make all
```

## Usage

### Drop-in Replacement for encoding/json

shape-json provides a complete encoding/json compatible API, making it a drop-in replacement for the standard library:

```go
import "github.com/shapestone/shape-json/pkg/json"

// Marshal: Convert Go structs to JSON
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

person := Person{Name: "Alice", Age: 30}
data, err := json.Marshal(person)
// data: {"age":30,"name":"Alice"}

// Unmarshal: Parse JSON into Go structs
var decoded Person
err = json.Unmarshal(data, &decoded)

// Decoder: Stream JSON from io.Reader
file, _ := os.Open("data.json")
decoder := json.NewDecoder(file)
var result Person
err = decoder.Decode(&result)

// Encoder: Stream JSON to io.Writer
encoder := json.NewEncoder(os.Stdout)
err = encoder.Encode(person)
```

**Struct Tag Support:**
```go
type User struct {
    PublicName  string `json:"name"`              // Rename field
    Password    string `json:"-"`                 // Skip field
    Email       string `json:"email,omitempty"`   // Omit if empty
    Count       int    `json:"count,string"`      // Marshal as string
}
```

**JSON Formatting:**
```go
// Pretty printing with indentation
person := Person{Name: "Alice", Age: 30}
pretty, _ := json.MarshalIndent(person, "", "  ")
// {
//   "age": 30,
//   "name": "Alice"
// }

// Compact existing JSON (remove whitespace)
var compact bytes.Buffer
json.Compact(&compact, prettyJSON)

// Add indentation to existing JSON
var indented bytes.Buffer
json.Indent(&indented, compactJSON, "", "  ")
```

### Fluent DOM API (Recommended)

The DOM API provides a user-friendly, type-safe way to build and manipulate JSON without type assertions:

```go
import "github.com/shapestone/shape-json/pkg/json"

// Build JSON with fluent chaining (no type assertions needed!)
doc := json.NewDocument().
    SetString("name", "Alice").
    SetInt("age", 30).
    SetBool("active", true).
    SetObject("address", json.NewDocument().
        SetString("city", "NYC").
        SetString("zip", "10001")).
    SetArray("tags", json.NewArray().
        AddString("go").
        AddString("json"))

jsonStr, _ := doc.JSON()
// {"name":"Alice","age":30,"active":true,"address":{"city":"NYC","zip":"10001"},"tags":["go","json"]}

// Pretty print with indentation
pretty, _ := doc.JSONIndent("", "  ")
// {
//   "name": "Alice",
//   "age": 30,
//   ...
// }

// Parse and access with type-safe getters (no type assertions!)
doc, _ := json.ParseDocument(`{"user":{"name":"Bob","age":25}}`)
user, _ := doc.GetObject("user")
name, _ := user.GetString("name")  // "Bob" - clean and simple!
age, _ := user.GetInt("age")       // 25

// Work with arrays naturally
arr := json.NewArray().AddString("apple").AddString("banana")
fruit, _ := arr.GetString(0)  // "apple"
```

See [examples/dom_builder](examples/dom_builder) for comprehensive examples.

### JSON Validation

Validate JSON input before parsing (idiomatic Go approach):

```go
import "github.com/shapestone/shape-json/pkg/json"

// Validate from string - idiomatic Go
if err := json.Validate(`{"name": "Alice"}`); err != nil {
    fmt.Println("Invalid JSON:", err)
    // Output: Invalid JSON: unexpected end of input
}
// err == nil means valid JSON

// Validate from io.Reader - idiomatic Go
file, _ := os.Open("data.json")
defer file.Close()
if err := json.ValidateReader(file); err != nil {
    fmt.Println("Invalid JSON:", err)
}
// err == nil means valid JSON
```

### Low-Level AST API

For advanced use cases, access the underlying AST:

```go
import (
    "github.com/shapestone/shape/pkg/ast"
    "github.com/shapestone/shape-json/pkg/json"
)

// Parse JSON into AST
node, err := json.Parse(`{"name": "Alice", "age": 30}`)
if err != nil {
    log.Fatal("Parse error:", err)
}

// Access parsed data
obj := node.(*ast.ObjectNode)
nameNode, _ := obj.GetProperty("name")
name := nameNode.(*ast.LiteralNode).Value().(string)  // "Alice"
```

## Features

- **Fluent DOM (Document Object Model) API**: User-friendly JSON manipulation (Recommended)
  - Type-safe getters (`GetString`, `GetInt`, `GetBool`, etc.) - No type assertions!
  - Fluent builder pattern with method chaining
  - Clean array semantics (not objects with numeric keys)
  - `Document` and `Array` types for intuitive JSON handling
  - Pretty-printing with `JSONIndent()` method
  - See [examples/dom_builder](examples/dom_builder) for details
- **encoding/json Compatible API**: Drop-in replacement for standard library
  - `Marshal()` / `Unmarshal()` - Convert between Go structs and JSON
  - `MarshalIndent()` / `Indent()` / `Compact()` - JSON formatting and pretty-printing
  - `Encoder` / `Decoder` - Streaming JSON I/O
  - Full struct tag support (`json:"name,omitempty,string,-"`)
  - **Pure implementation**: Does NOT use encoding/json internally
- **JSON Validation**: Idiomatic error-based validation
  - `Validate()` / `ValidateReader()` - Returns nil if valid, error with details if invalid
- **Complete JSON Support**: Full RFC 8259 (the JSON internet standard) compliance
- **Proper Type Distinction**: Empty arrays `[]` and empty objects `{}` are properly distinguished with full round-trip fidelity
- **LL(1) (top-down, one-token lookahead) Recursive Descent Parser**: Hand-coded, optimized parser
- **Shape AST Integration**: Returns unified AST nodes for advanced use cases
- **JSONPath Query Engine**: RFC 9535-compliant JSONPath implementation (see [pkg/jsonpath](pkg/jsonpath/README.md))
- **Comprehensive Error Messages**: Context-aware error reporting
- **High Test Coverage**: 91.0% JSON API, 90.2% fastparser, 92.2% parser, 69.9% tokenizer, 89.8% JSONPath
- **Zero External Dependencies** (except Shape infrastructure)
- **Architectural Purity**: Single unified parser - no encoding/json dependency

## Performance

shape-json uses an **intelligent dual-path architecture** that automatically selects the optimal parsing strategy:

### ⚡ Fast Path (Default - Optimized for Speed)

The fast path bypasses AST construction for maximum performance:

- **APIs**: `Unmarshal()`, `Validate()`, `ValidateReader()`
- **Performance**:
  - **9.7x faster** than AST path
  - **10.2x less memory** usage
  - **6.3x fewer allocations**
- **Use when**: You just need to populate Go types or validate JSON syntax
- **Automatically selected** when you use `Unmarshal()` or `Validate()`

```go
// Fast path - automatically selected (9.7x faster!)
var person Person
json.Unmarshal(data, &person)

// Fast path - validation only (8.2x faster!)
if err := json.Validate(jsonString); err != nil {
    // Invalid JSON
}
```

**Benchmark results (410 KB JSON file)**:
- Fast Path: 1.69 ms, 1.29 MB, 36,756 allocs
- AST Path: 16.45 ms, 13.19 MB, 230,025 allocs

### 🌳 AST Path (Full Features)

The AST path builds a complete Abstract Syntax Tree for advanced use cases:

- **APIs**: `Parse()`, `ParseReader()`, `ParseDocument()`, `ParseArray()`
- **Performance**: Slower, more memory (enables advanced features)
- **Use when**: You need JSONPath queries, tree manipulation, or format conversion
- **Trade-off**: Richer features at the cost of performance

```go
// AST path - full tree structure for advanced features
node, _ := json.Parse(jsonString)
results := jsonpath.Query(node, "$.users[?(@.age > 30)].name")

// Or use DOM API (built on AST path)
doc, _ := json.ParseDocument(jsonString)
user, _ := doc.GetObject("user")
```

### Marshal Performance

shape-json v0.10.0 introduces a compiled encoder cache for `Marshal()`, delivering **5.4x faster struct marshaling** and **3.1x faster map marshaling** compared to the previous version, with 6-9x fewer allocations. The compiled encoder eliminates per-call reflection by caching type-level encoders with pre-encoded key bytes.

### Choosing the Right API

**Most common use cases** (80-90% of usage) → Use Fast Path:
- ✅ Unmarshal JSON into Go structs
- ✅ Validate JSON syntax
- ✅ Parse into `interface{}`

**Advanced features** → Use AST Path:
- 🔍 JSONPath queries
- 🌳 Tree manipulation/transformation
- 📊 Format conversion (JSON ↔ XML ↔ YAML)

**Performance tip**: If you only need data, use `Unmarshal()`. If you need both data AND the AST, use `Parse()` followed by `NodeToInterface()`:

```go
// Option 1: Fast path for speed (recommended for most cases)
var data MyStruct
json.Unmarshal(bytes, &data)  // 9.7x faster

// Option 2: AST path when you need the tree
node, _ := json.Parse(string(bytes))
jsonpath.Query(node, "$.items[*].price")  // JSONPath needs AST
```

## Architecture

shape-json uses a **unified architecture** with a single custom parser for all APIs:

- **Grammar-Driven**: EBNF (Extended Backus-Naur Form) grammar in `docs/grammar/json.ebnf`
- **Tokenizer**: Custom tokenizer using Shape's framework
- **Parser**: LL(1) (top-down, one-token lookahead) recursive descent with single token lookahead
- **Rendering**: Custom JSON renderer (no encoding/json dependency)
- **AST Representation**:
  - Objects → `*ast.ObjectNode` with properties map
  - Arrays → `*ast.ArrayDataNode` with elements slice
  - Primitives → `*ast.LiteralNode` (string, int64, float64, bool, nil)
- **Data Flow**:
  - Parse: `JSON string → Tokenizer → Parser → AST`
  - Render: `AST → Renderer → JSON string`
  - Marshal: `Go types → compiled encoder cache → JSON bytes`
  - Unmarshal: `JSON bytes → Parser → AST → Go types`
  - DOM: `JSON → AST → map/slice → Document/Array`

## JSONPath Query Engine

shape-json includes a complete JSONPath implementation for querying JSON data:

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
            map[string]interface{}{"author": "Nigel Rees"},
            map[string]interface{}{"author": "Evelyn Waugh"},
        },
    },
}

results := expr.Get(data)
// Results: ["Nigel Rees", "Evelyn Waugh"]
```

**Supported features:**
- Root selector (`$`)
- Child selectors (`.property`, `['property']`)
- Wildcards (`*`, `[*]`)
- Array indexing (`[0]`, `[-1]`)
- Array slicing (`[0:5]`, `[:5]`, `[2:]`)
- Recursive descent (`..property`)
- Filter expressions (`[?(@.price < 10)]`, `[?(@.role == 'admin')]`)

See [pkg/jsonpath/README.md](pkg/jsonpath/README.md) for complete documentation.

## Examples

The `examples/` directory contains comprehensive usage examples:

### DOM Builder API (Recommended)
Build and manipulate JSON with a fluent, type-safe API:
```bash
go run examples/dom_builder/main.go
```

See [examples/dom_builder/README.md](examples/dom_builder/README.md) for detailed documentation and API comparison.

### JSON Validation
Validate JSON input before parsing:
```bash
go run examples/format_detection/main.go
```

See [examples/format_detection/README.md](examples/format_detection/README.md) for detailed documentation.

### Other Examples
- [main.go](examples/main.go) - Basic parsing with AST
- [parse_reader/](examples/parse_reader/) - Streaming JSON parsing from files
- [encoding_json_api/](examples/encoding_json_api/) - encoding/json compatibility

## Grammar

See [docs/grammar/json.ebnf](docs/grammar/json.ebnf) for the complete EBNF (Extended Backus-Naur Form) specification.

Key grammar rules:
```ebnf
Value = Object | Array | String | Number | Boolean | Null ;
Object = "{" [ Member { "," Member } ] "}" ;
Array = "[" [ Value { "," Value } ] "]" ;
```

## Thread Safety

**shape-json is thread-safe.** All public APIs can be called concurrently from multiple goroutines without external synchronization.

### Safe for Concurrent Use

```go
// ✅ SAFE: Multiple goroutines can call these concurrently
go func() {
    var v1 interface{}
    json.Unmarshal(data1, &v1)
}()

go func() {
    var v2 interface{}
    json.Unmarshal(data2, &v2)
}()

// ✅ SAFE: Parse, Marshal, Validate all create new instances
go func() { json.Parse(input1) }()
go func() { json.Marshal(obj1) }()
go func() { json.Validate(input2) }()
```

### Thread Safety Guarantees

- **`Unmarshal()`, `Marshal()`** - Thread-safe, use internal buffer pools
- **`Parse()`, `Validate()`** - Thread-safe, create new parser instances
- **`NewDecoder()`, `NewEncoder()`** - Thread-safe factory functions
- **Race detector verified** - All tests pass with `go test -race`

This matches `encoding/json`'s thread safety guarantees.

## Testing

shape-json has comprehensive test coverage including unit tests, fuzzing, and grammar verification.

### Coverage Summary
- Fast Parser: 90.2% ✅
- Parser: 92.2% ✅
- Tokenizer: 69.9% ✅
- JSONPath: 89.8% ✅
- JSON API: 91.0% ✅
- **Overall Library**: **85.84%** ✅

### Quick Start

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run fuzzing tests (5 seconds each)
go test ./internal/parser -fuzz=FuzzParserObjects -fuzztime=5s

# Run grammar verification
make grammar-test

# Run example
go run examples/main.go
```

### Fuzzing

The parser includes extensive fuzzing tests to ensure robustness:

```bash
# Fuzz parser (objects, arrays, strings, numbers, nested structures)
go test ./internal/parser -fuzz=FuzzParserObjects -fuzztime=30s
go test ./internal/parser -fuzz=FuzzParserArrays -fuzztime=30s
go test ./internal/parser -fuzz=FuzzParserStrings -fuzztime=30s

# Fuzz JSONPath (filters, slices, recursive descent, brackets)
go test ./pkg/jsonpath -fuzz=FuzzJSONPathFilters -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathSlices -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathRecursive -fuzztime=30s
```

**Recent fuzzing results:**
- Parser: 840K+ executions, 268 interesting cases, 0 crashes
- JSONPath: 1.17M+ executions, 269 interesting cases, 0 crashes

See [docs/TESTING.md](docs/TESTING.md) for comprehensive testing documentation.

## Documentation

- [EBNF Grammar](docs/grammar/json.ebnf) - Complete JSON grammar specification
- [JSONPath Query Engine](pkg/jsonpath/README.md) - JSONPath query documentation
- [Parser Implementation Guide](https://github.com/shapestone/shape-core/blob/main/docs/PARSER_IMPLEMENTATION_GUIDE.md) - Guide for implementing parsers
- [Shape Core Infrastructure](https://github.com/shapestone/shape-core) - Universal AST and tokenizer framework
- [Shape Ecosystem](https://github.com/shapestone/shape) - Documentation and examples

## Development

```bash
# Run tests
make test

# Generate coverage report
make coverage

# Build
make build

# Run all checks
make all
```

### Performance Benchmarking

shape-json includes a comprehensive benchmarking system with historical tracking:

#### Quick Start

```bash
# Run benchmarks with statistical output
make bench

# Generate performance report (automatically saves to history)
make performance-report

# View benchmark history
make bench-history

# Compare latest vs previous benchmark
make bench-compare-history
```

#### Available Benchmark Targets

**`make bench`** - Run all benchmarks with standard settings
- Quick performance check
- Shows ns/op, MB/s, B/op, allocs/op

**`make bench-report`** - Run benchmarks and save output to file
- Saves results to `benchmarks/results.txt`
- Useful for manual inspection

**`make bench-compare`** - Run benchmarks 10 times for statistical analysis
- Creates `benchmarks/benchstat.txt` with multiple runs
- Analyze with: `benchstat benchmarks/benchstat.txt`
- Requires: `go install golang.org/x/perf/cmd/benchstat@latest`

**`make bench-profile`** - Run benchmarks with CPU and memory profiling
- Generates `benchmarks/cpu.prof` and `benchmarks/mem.prof`
- Analyze with: `go tool pprof benchmarks/cpu.prof`

**`make performance-report`** - Generate comprehensive performance report
- Runs all benchmarks with 3-second minimum per test
- Creates detailed `PERFORMANCE_REPORT.md` with analysis
- Automatically saves timestamped history to `benchmarks/history/`
- Includes metadata (git commit, platform, Go version)

**`make bench-history`** - List all historical benchmark runs
- Shows timestamp, git commit, and platform for each run
- Helps track performance over time

**`make bench-compare-history`** - Compare latest vs previous benchmark
- Uses benchstat for statistical comparison
- Shows performance improvements/regressions
- Requires benchstat: `go install golang.org/x/perf/cmd/benchstat@latest`

#### Tracking Performance Over Time

The benchmark system automatically saves historical data:

```bash
# Run a benchmark and save to history
make performance-report

# View all historical runs
make bench-history

# Example output:
# Available benchmark history:
#
#   2025-12-21_16-30-45
#     "commit": "ee2758d",
#     "platform": "Apple M1 Pro",
#
#   2025-12-20_14-22-10
#     "commit": "a1b2c3d",
#     "platform": "Apple M1 Pro",
```

Each benchmark run is saved to `benchmarks/history/YYYY-MM-DD_HH-MM-SS/`:
- `benchmark_output.txt` - Raw benchmark results
- `PERFORMANCE_REPORT.md` - Generated performance report
- `metadata.json` - Git commit, platform info, timestamp

#### Comparing Benchmark Runs

Compare two benchmark runs to track performance changes:

```bash
# Compare latest vs previous (most common)
make bench-compare-history

# Or use the comparison tool directly
go run scripts/compare_benchmarks/main.go latest previous

# Compare specific timestamps
go run scripts/compare_benchmarks/main.go 2025-12-20_14-22-10 2025-12-21_16-30-45
```

The comparison uses benchstat to show:
- Statistical significance of changes
- Speed improvements/regressions
- Memory usage changes
- Allocation differences

#### Custom Benchmark Runs

Add descriptions to track specific changes:

```bash
# Run with description
go run scripts/generate_benchmark_report/main.go -description "After parser optimization"

# Disable history saving
go run scripts/generate_benchmark_report/main.go -save-history=false
```

For detailed information on the benchmarking system, see:
- [PERFORMANCE_REPORT.md](PERFORMANCE_REPORT.md) - Latest benchmark results
- [benchmarks/history/README.md](benchmarks/history/README.md) - History tracking guide

## Make Targets

| Target | Description |
|--------|-------------|
| `make test` | Run all tests with race detector |
| `make lint` | Run golangci-lint static analysis |
| `make build` | Build the project |
| `make coverage` | Generate HTML coverage report in `coverage/` |
| `make all` | Run grammar-verify, test, lint, build, coverage |
| `make grammar-test` | Run grammar verification tests |
| `make grammar-verify` | Verify grammar files exist and are valid |
| `make bench` | Run all benchmarks |
| `make bench-report` | Run benchmarks and save to `benchmarks/results.txt` |
| `make bench-compare` | Run 10x for statistical analysis (benchstat) |
| `make bench-profile` | Run benchmarks with CPU and memory profiling |
| `make performance-report` | Generate full report + save to `benchmarks/history/` |
| `make bench-history` | List all historical benchmark runs |
| `make bench-compare-history` | Compare latest vs previous benchmark run |
| `make clean` | Remove `coverage/` and `benchmarks/` directories |

## Related Projects

| Repository | Description |
|-----------|-------------|
| [shapestone/shape](https://github.com/shapestone/shape) | Shape Parser™ — multi-format parser ecosystem |
| [shapestone/shape-core](https://github.com/shapestone/shape-core) | Universal AST and tokenizer framework |
| [shapestone/inkling](https://github.com/shapestone/inkling) | Example custom DSL using Shape's tokenizer |

## License

Apache License 2.0

Copyright © 2020-2025 Shapestone

See [LICENSE](LICENSE) for the full license text and [NOTICE](NOTICE) for third-party attributions.
