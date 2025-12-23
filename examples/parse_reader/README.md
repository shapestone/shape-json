# ParseReader Example

This example demonstrates how to use the `ParseReader()` API to parse JSON from various `io.Reader` sources.

## Running the Example

```bash
go run main.go
```

## What it Demonstrates

1. **Parsing from strings.Reader** - Converting a string to a reader and parsing
2. **Parsing from bytes.Buffer** - Parsing JSON from a byte buffer
3. **Parsing from files** - Reading and parsing JSON files using `os.Open()`
4. **Comparison** - Shows that `Parse()` and `ParseReader()` produce identical results
5. **Usage guidelines** - When to use each API

## When to Use ParseReader

Use `ParseReader()` when:

- Reading JSON from files
- Processing large JSON documents (constant memory usage)
- Streaming JSON from network connections
- Reading from compressed streams (gzip, etc.)
- Working with any `io.Reader` source

## When to Use Parse

Use `Parse()` when:

- JSON is already in memory as a string
- Working with small JSON documents
- Simplicity is preferred over streaming

## Memory Efficiency

The `ParseReader()` API uses a buffered stream implementation that reads data in chunks (64KB buffer).
This allows parsing very large JSON files without loading the entire file into memory at once,
maintaining constant memory usage regardless of file size.

## Example Output

```
=== Example 1: Parse from strings.Reader ===
Parsed object with 3 properties
Name: Alice

=== Example 2: Parse from bytes.Buffer ===
Product price: $29.99

=== Example 3: Parse from file ===
Parsed file successfully
Number of users: 3
First user: Alice

=== Example 4: Parse vs ParseReader ===
Parse() result: Hello, World!
ParseReader() result: Hello, World!

Both methods produce the same AST representation
```
