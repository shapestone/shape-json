# JSON Grammar Specification

This directory contains EBNF grammar specifications for JSON parsing.

## Files

### json.ebnf
Complete RFC 8259 compliant JSON grammar specification with full support for:
- Unicode characters and escape sequences
- All JSON value types (objects, arrays, strings, numbers, booleans, null)
- Complete number format (integers, floats, scientific notation)
- Full escape sequence support (`\"`, `\\`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`, `\uXXXX`)

This is the primary grammar specification and serves as the authoritative documentation for JSON parsing implementation.

**Note**: Uses advanced EBNF features (character classes with negation and hex escapes like `[^"\\\x00-\x1F]`) that are not yet fully supported by Shape's EBNF parser.

### json-simple.ebnf
Simplified JSON grammar for grammar verification tests (ADR 0005).

Uses only basic EBNF features supported by Shape's grammar parser:
- Simple character class ranges (`[a-zA-Z0-9]`)
- No negated character classes
- No hex escape sequences in character classes
- Simplified string content (basic ASCII characters only)

This grammar is used by `internal/parser/grammar_test.go` for automated grammar verification testing.

## Purpose

These grammars serve multiple purposes:

1. **Documentation**: Formal specification of JSON syntax for implementers
2. **Testing**: Automated verification that parser matches grammar (ADR 0005)
3. **Validation**: Reference for parser implementation correctness
4. **Coverage**: Ensures all grammar rules are exercised by tests

## Implementation Notes

The parser implementation in `internal/parser/parser.go` follows the full `json.ebnf` grammar specification using LL(1) recursive descent parsing as described in Shape ADR 0004.

Each grammar rule maps to a parse function:
- `Value` → `Parse()`
- `Object` → `parseObject()`
- `Array` → `parseArray()`
- `String` → `parseString()`
- `Number` → `parseNumber()`
- `Boolean` → `parseBoolean()`
- `Null` → `parseNull()`

## Future Work

When Shape's EBNF parser supports advanced character class syntax (negation, hex escapes), the simplified grammar can be deprecated and all grammar tests can use the full `json.ebnf` specification directly.
