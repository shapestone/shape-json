# JSON DOM Builder API Example

This example demonstrates the fluent DOM API for building and manipulating JSON documents without type assertions.

## Overview

The DOM API provides a user-friendly, type-safe interface for working with JSON data, eliminating the need for type assertions and complex AST manipulation.

## Features Demonstrated

### 1. Simple Document Building
```go
doc := json.NewDocument().
    SetString("name", "Alice").
    SetInt("age", 30).
    SetBool("active", true)
```

### 2. Type-Safe Access
```go
doc, _ := json.ParseDocument(`{"name":"Alice","age":30}`)
name, _ := doc.GetString("name")  // No type assertions!
age, _ := doc.GetInt("age")
```

### 3. Nested Structures
```go
doc := json.NewDocument().
    SetString("name", "Alice").
    SetObject("address", json.NewDocument().
        SetString("city", "NYC").
        SetString("zip", "10001"))
```

### 4. Arrays
```go
arr := json.NewArray().
    AddString("go").
    AddString("json").
    AddInt(42)
```

### 5. Round-Trip
```go
// Build
doc := json.NewDocument().SetString("name", "Alice")

// Marshal
jsonStr, _ := doc.JSON()

// Parse back
parsed, _ := json.ParseDocument(jsonStr)

// Access
name, _ := parsed.GetString("name")
```

## Comparison: Old vs New API

### Old AST API (Confusing)
```go
// Parse
node, _ := json.Parse(input)
obj := node.(*ast.ObjectNode)              // Type assertion

// Access nested value
userNode, _ := obj.GetProperty("user")
userObj := userNode.(*ast.ObjectNode)      // Another assertion
nameNode, _ := userObj.GetProperty("name")
name := nameNode.(*ast.LiteralNode).Value().(string) // Two more assertions!
```

**Problems:**
- ❌ Multiple type assertions required
- ❌ Verbose, error-prone code
- ❌ Need to understand AST internals
- ❌ Arrays are objects with string keys ("0", "1", ...)

### New DOM API (Clean)
```go
// Parse and access
doc, _ := json.ParseDocument(input)
user, _ := doc.GetObject("user")
name, _ := user.GetString("name")  // That's it!
```

**Benefits:**
- ✅ No type assertions
- ✅ Clean, readable code
- ✅ Type-safe getters (GetString, GetInt, GetBool, etc.)
- ✅ Arrays have proper array semantics
- ✅ Fluent builder pattern for construction
- ✅ IDE autocomplete support

## Running the Example

```bash
cd examples/dom_builder
go run main.go
```

## API Reference

### Document Methods

**Setters (return *Document for chaining):**
- `SetString(key, value string) *Document`
- `SetInt(key string, value int) *Document`
- `SetInt64(key string, value int64) *Document`
- `SetBool(key string, value bool) *Document`
- `SetFloat(key string, value float64) *Document`
- `SetNull(key string) *Document`
- `SetObject(key string, value *Document) *Document`
- `SetArray(key string, value *Array) *Document`

**Getters (type-safe, return (value, bool)):**
- `GetString(key string) (string, bool)`
- `GetInt(key string) (int, bool)`
- `GetInt64(key string) (int64, bool)`
- `GetBool(key string) (bool, bool)`
- `GetFloat(key string) (float64, bool)`
- `GetObject(key string) (*Document, bool)`
- `GetArray(key string) (*Array, bool)`

**Utilities:**
- `Has(key string) bool` - Check if key exists
- `IsNull(key string) bool` - Check if value is null
- `Remove(key string) *Document` - Remove key
- `Keys() []string` - Get all keys
- `Size() int` - Get number of properties
- `JSON() (string, error)` - Marshal to JSON string

### Array Methods

**Adders (return *Array for chaining):**
- `AddString(value string) *Array`
- `AddInt(value int) *Array`
- `AddInt64(value int64) *Array`
- `AddBool(value bool) *Array`
- `AddFloat(value float64) *Array`
- `AddNull() *Array`
- `AddObject(value *Document) *Array`
- `AddArray(value *Array) *Array`

**Getters (type-safe, return (value, bool)):**
- `GetString(index int) (string, bool)`
- `GetInt(index int) (int, bool)`
- `GetInt64(index int) (int64, bool)`
- `GetBool(index int) (bool, bool)`
- `GetFloat(index int) (float64, bool)`
- `GetObject(index int) (*Document, bool)`
- `GetArray(index int) (*Array, bool)`

**Utilities:**
- `Len() int` - Get array length
- `IsNull(index int) bool` - Check if value is null
- `JSON() (string, error)` - Marshal to JSON string

## Related Examples

- [Basic Parsing](../main.go) - Basic JSON parsing with AST
- [Format Detection](../format_detection/) - JSON validation and detection
- [Parse Reader](../parse_reader/) - Streaming JSON parsing

## Documentation

See [pkg/json/dom.go](../../pkg/json/dom.go) for complete API documentation.
