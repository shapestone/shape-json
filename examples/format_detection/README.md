# JSON Validation Example

This example demonstrates JSON validation using clean, boolean-based APIs.

## Overview

shape-json provides simple validation functions that return boolean results or detailed error messages, making it easy to validate JSON before parsing.

## API (Idiomatic Go)

### Validation from String
```go
if err := json.Validate(`{"name": "Alice"}`); err != nil {
    fmt.Println("Invalid JSON:", err)
    // Output: Invalid JSON: expected RBrace at line 1, column 17
}
// err == nil means valid JSON
```

### Validation from io.Reader
```go
file, _ := os.Open("data.json")
defer file.Close()
if err := json.ValidateReader(file); err != nil {
    fmt.Println("Invalid JSON:", err)
}
// err == nil means valid JSON
```

## API Functions

**`Validate(input string) error`**
- Returns nil if valid JSON
- Returns error with detailed message if invalid
- Idiomatic Go approach: `if err != nil`

**`ValidateReader(reader io.Reader) error`**
- Returns nil if valid JSON
- Returns error with detailed message if invalid
- Reads entire input from reader
- Idiomatic Go approach: `if err != nil`

## Running the Example

```bash
cd examples/format_detection
go run main.go
```

## Examples Demonstrated

### 1. Simple Boolean Validation
Quick true/false check for valid JSON.

### 2. Validation with Error Details
Get specific error messages explaining why JSON is invalid.

### 3. Validate Different Types
Works with all JSON types: objects, arrays, strings, numbers, booleans, null.

### 4. Common JSON Errors
Detects common mistakes:
- Empty strings
- Plain text
- Missing quotes on keys
- Trailing commas
- Unclosed braces/brackets
- Single quotes instead of double quotes

### 5. Validation from io.Reader
Validate JSON from any io.Reader (strings.Reader, files, network streams).

### 6. Validation from File
Read and validate JSON files.

### 7. Validate Before Parsing (Best Practice)
Recommended pattern: validate first, then parse.

```go
if err := json.Validate(userInput); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}

// Now safely parse
doc, _ := json.ParseDocument(userInput)
```

### 8. API Input Validation Pattern
Use validation in API handlers to reject bad requests early.

## Comparison: Old vs New API

### Old API (Deprecated)
```go
// Returns "JSON" string on success - redundant!
format, err := json.DetectFormat(input)
if err != nil {
    fmt.Println("Not JSON:", err)
} else {
    fmt.Println("Format:", format)  // Always "JSON"
}
```

**Problems:**
- ❌ Returns meaningless "JSON" string
- ❌ Confusing API (format detection in JSON package?)
- ❌ Verbose

### New API (Idiomatic Go)
```go
// Simple and clear - just check the error
if err := json.Validate(input); err != nil {
    fmt.Println("Invalid:", err)
}
```

**Benefits:**
- ✅ Idiomatic Go (error-based)
- ✅ Single function, not two
- ✅ Clear intent
- ✅ Detailed error messages

## Use Cases

**Quick Validation:**
```go
if err := json.Validate(userInput); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}
```

**Validation with Logging:**
```go
if err := json.Validate(userInput); err != nil {
    log.Printf("Invalid JSON from user %s: %v", userID, err)
    return err
}
```

**Pre-flight Check:**
```go
// Validate before expensive operations
if err := json.Validate(largeInput); err != nil {
    return nil, fmt.Errorf("invalid JSON, skipping processing: %w", err)
}

// Now safely parse and process
doc, _ := json.ParseDocument(largeInput)
```

**API Request Validation:**
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)

    if err := json.Validate(string(body)); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Process valid JSON
    doc, _ := json.ParseDocument(string(body))
    // ...
}
```

## Related Examples

- [DOM Builder](../dom_builder/) - Fluent API for JSON manipulation
- [encoding/json API](../encoding_json_api/) - Drop-in stdlib replacement
- [Parse Reader](../parse_reader/) - Streaming JSON parsing

## Documentation

See [pkg/json/parser.go](../../pkg/json/parser.go) for complete API documentation.
