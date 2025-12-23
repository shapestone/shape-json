# Testing Guide

This document describes the testing strategy and coverage for shape-json.

## Test Coverage Summary

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/parser` | 97.1% | ✅ Excellent |
| `internal/tokenizer` | 98.8% | ✅ Excellent |
| `pkg/jsonpath` | 89.8% | ✅ Excellent |
| `pkg/json` | 91.6% | ✅ Excellent |
| **Overall Library** | **91.9%** | ✅ Excellent |

### Parser Test Coverage Details

All critical parser functions have excellent coverage:

| Function | Coverage |
|----------|----------|
| `NewParser` | 100% |
| `NewParserFromStream` | 100% |
| `Parse` | 100% |
| `parseValue` | 100% |
| `parseObject` | 90.0% |
| `parseMember` | 100% |
| `parseArray` | 95.2% |
| `parseString` | 100% |
| `parseNumber` | 92.9% |
| `parseBoolean` | 100% |
| `parseNull` | 100% |
| Helper functions | 100% |

### JSON API Test Coverage Details

All critical JSON API functions have excellent coverage:

| Function | Coverage |
|----------|----------|
| `Marshal` | 100% |
| `Unmarshal` | 93.3% |
| `unmarshalLiteral` | 91.7% |
| `unmarshalArray` | 91.3% |
| `marshalValue` | 91.1% |
| `marshalString` | 80.0% |
| `Encode` | 100% |
| `NewEncoder` | 100% |
| `NewDecoder` | 100% |
| `DetectFormat` | 100% |
| `DetectFormatFromReader` | 100% |

## Test Categories

### 1. Unit Tests

Located in `*_test.go` files throughout the codebase.

**Parser Tests** (`internal/parser/`):
- `parser_test.go` - Core parsing tests for all JSON value types
- `parser_coverage_test.go` - Edge case and error handling tests
- `grammar_test.go` - Grammar verification tests (ADR 0005)

**Tokenizer Tests** (`internal/tokenizer/`):
- `tokenizer_test.go` - Token recognition and generation tests

**JSON API Tests** (`pkg/json/`):
- `marshal_test.go` - Marshal/Unmarshal tests
- `decoder_test.go` - Stream decoder tests
- `encoder_test.go` - Stream encoder tests
- `json_coverage_test.go` - Comprehensive edge case and type conversion tests

**JSONPath Tests** (`pkg/jsonpath/`):
- `jsonpath_test.go` - Core JSONPath query tests
- `jsonpath_coverage_test.go` - Edge cases and comprehensive scenarios
- `filter_test.go` - Filter expression tests
- `parser_test.go` - JSONPath parser tests
- `executor_test.go` - Query execution tests
- RFC 9535 compliance tests

### 2. Fuzzing Tests

**Parser Fuzzing** (`internal/parser/parser_fuzz_test.go`):

Fuzzing ensures the parser handles arbitrary input gracefully without panicking.

#### Running Fuzzing Tests

Run all fuzz seed tests (quick):
```bash
go test ./internal/parser -run Fuzz
```

Run specific fuzzer for extended period:
```bash
# Fuzz general JSON parsing for 30 seconds
go test ./internal/parser -fuzz=FuzzParserObjects -fuzztime=30s

# Fuzz arrays
go test ./internal/parser -fuzz=FuzzParserArrays -fuzztime=30s

# Fuzz strings
go test ./internal/parser -fuzz=FuzzParserStrings -fuzztime=30s

# Fuzz numbers
go test ./internal/parser -fuzz=FuzzParserNumbers -fuzztime=30s

# Fuzz nested structures
go test ./internal/parser -fuzz=FuzzParserNested -fuzztime=30s
```

Run all fuzzers in CI/development:
```bash
# Quick smoke test (1 second each)
for fuzzer in FuzzParserObjects FuzzParserArrays FuzzParserStrings FuzzParserNumbers FuzzParserNested; do
    go test ./internal/parser -fuzz=$fuzzer -fuzztime=1s
done
```

#### Fuzzing Results

Recent fuzzing run (5 seconds):
- **Executions**: 841,856 inputs tested
- **Rate**: ~140,000 exec/sec
- **Interesting cases found**: 268 unique variations
- **Crashes**: 0 (parser is robust)

#### Fuzz Test Types

1. **FuzzParserObjects** - Tests object parsing with random variations
2. **FuzzParserArrays** - Tests array parsing with random variations
3. **FuzzParserStrings** - Tests string parsing including escape sequences
4. **FuzzParserNumbers** - Tests number parsing (integers, floats, scientific notation)
5. **FuzzParserNested** - Tests deeply nested structures

Each fuzzer includes a seed corpus of valid JSON examples that the Go fuzzer mutates to create interesting test cases.

**JSONPath Fuzzing** (`pkg/jsonpath/jsonpath_fuzz_test.go`):

- **FuzzJSONPath** - General JSONPath fuzzing with 35+ seed examples
- **FuzzJSONPathFilters** - Filter expression fuzzing
- **FuzzJSONPathSlices** - Array slicing fuzzing
- **FuzzJSONPathRecursive** - Recursive descent fuzzing
- **FuzzJSONPathBracket** - Bracket notation fuzzing
- **FuzzJSONPathComplex** - Complex nested path fuzzing

Run JSONPath fuzzers:
```bash
go test ./pkg/jsonpath -fuzz=FuzzJSONPathFilters -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathSlices -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathRecursive -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathBracket -fuzztime=30s
go test ./pkg/jsonpath -fuzz=FuzzJSONPathComplex -fuzztime=30s
```

JSONPath fuzzing results (5 seconds):
- **Executions**: 1.17M+ inputs tested
- **Rate**: ~250,000 exec/sec
- **Interesting cases found**: 269 unique variations
- **Crashes**: 0 (JSONPath engine is robust)

### 3. Grammar Verification Tests

Located in `internal/parser/grammar_test.go`.

Implements ADR 0005: Grammar-as-Verification.

- **TestGrammarVerification** - Auto-generates test cases from EBNF grammar
- **TestGrammarCoverage** - Ensures all grammar rules are exercised (88.9% coverage)
- **TestGrammarFileExists** - Validates grammar files
- **TestGrammarDocumentation** - Verifies grammar documentation

### 4. Integration Tests

Located in `examples/main.go` and various package tests.

Tests the complete flow from parsing to AST manipulation.

## Running Tests

### Run All Tests
```bash
make test
# or
go test ./...
```

### Run Tests with Coverage
```bash
make coverage
# or
go test -v -coverprofile=coverage.out ./internal/... ./pkg/...
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Package Tests
```bash
# Parser only
go test ./internal/parser -v

# Tokenizer only
go test ./internal/tokenizer -v

# JSON API only
go test ./pkg/json -v

# JSONPath only
go test ./pkg/jsonpath -v
```

### Run Grammar Tests
```bash
# All grammar tests
make grammar-test

# Verify grammar files exist
make grammar-verify
```

### Run Tests with Race Detection
```bash
go test -race ./...
```

## Test Best Practices

### Writing New Tests

1. **Follow Table-Driven Pattern**
   ```go
   tests := []struct {
       name  string
       input string
       want  result
   }{
       {name: "case1", input: "...", want: ...},
       {name: "case2", input: "...", want: ...},
   }

   for _, tt := range tests {
       t.Run(tt.name, func(t *testing.T) {
           // Test logic
       })
   }
   ```

2. **Test Both Success and Error Cases**
   - Valid inputs should parse successfully
   - Invalid inputs should return appropriate errors
   - Error messages should be descriptive

3. **Use Descriptive Test Names**
   - Good: `TestParse_ObjectWithTrailingComma`
   - Bad: `TestParse1`

4. **Test Edge Cases**
   - Empty inputs
   - Very large inputs
   - Deeply nested structures
   - Special characters and escape sequences
   - Boundary conditions (max integers, etc.)

### Coverage Goals

- **Parser code**: 95%+ (currently 97.1% ✅)
- **Tokenizer code**: 95%+ (currently 98.8% ✅)
- **JSONPath code**: 85%+ (currently 89.8% ✅)
- **JSON API code**: 85%+ (currently 91.6% ✅)
- **Overall library**: 85%+ (currently 91.9% ✅)

## Continuous Integration

The CI workflow (`.github/workflows/ci.yml`) runs:

1. Grammar verification tests
2. Full test suite with race detection
3. Coverage report (fails if below 70%)
4. Linting
5. Build verification

Coverage reports are uploaded to Codecov for tracking over time.

## Test Data

### Valid JSON Examples

The test suite includes comprehensive examples:
- Primitives: `null`, `true`, `false`, numbers, strings
- Objects: empty, single property, multiple properties, nested
- Arrays: empty, single element, multiple elements, nested
- Mixed structures: objects containing arrays, arrays containing objects
- Edge cases: whitespace variations, escape sequences, large numbers

### Invalid JSON Examples

The test suite validates error handling for:
- Trailing commas
- Missing commas
- Unclosed brackets/braces
- Invalid tokens
- Malformed numbers
- Unterminated strings
- Non-string object keys
- Unexpected tokens

## Debugging Failed Tests

### View Detailed Test Output
```bash
go test ./internal/parser -v
```

### Run Single Test
```bash
go test ./internal/parser -run TestParse_String -v
```

### Run Tests with Coverage Visualization
```bash
go test -coverprofile=coverage.out ./internal/parser
go tool cover -html=coverage.out
```

### Check Uncovered Lines
```bash
go tool cover -func=coverage.out | grep -v "100.0%"
```

## Performance Testing

### Benchmark Tests

Run benchmarks:
```bash
go test -bench=. ./internal/parser
go test -bench=. ./pkg/json
```

### Memory Usage

Check memory allocations:
```bash
go test -bench=. -benchmem ./internal/parser
```

## Related Documentation

- [Parser Implementation Guide](https://github.com/shapestone/shape-core/blob/main/docs/PARSER_IMPLEMENTATION_GUIDE.md)
- [Grammar Specification](grammar/json.ebnf)
- [ADR 0005: Grammar-as-Verification](https://github.com/shapestone/shape-core/blob/main/docs/adr/0005-grammar-as-verification.md)
