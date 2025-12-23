# JSONPath Query Engine

A Go implementation of a JSONPath query engine for querying JSON-like data structures. This package supports a subset of RFC 9535 JSONPath syntax.

## Features

- **Root selector**: `$`
- **Child selector**: `.property` or `['property']`
- **Wildcard**: `*` or `[*]`
- **Array index**: `[0]`, `[1]`, etc.
- **Array slice**: `[0:5]`, `[:5]`, `[2:]`
- **Recursive descent**: `..property`
- **Multiple selectors**: `$.a.b.c`
- **Filter expressions**: `[?(@.field operator value)]`
  - Comparison operators: `<`, `>`, `<=`, `>=`, `==`, `!=`
  - Logical operators: `&&`, `||`
  - Regex matching: `=~`
  - Field existence: `[?(@.field)]`

## Installation

```go
import "github.com/shapestone/shape-json/pkg/jsonpath"
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/shapestone/shape-json/pkg/jsonpath"
)

func main() {
    // Sample data
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

    // Parse a JSONPath query
    expr, err := jsonpath.ParseString("$.store.book[*].author")
    if err != nil {
        log.Fatal(err)
    }

    // Execute the query
    results := expr.Get(data)

    // Print results
    for _, author := range results {
        fmt.Println(author)
    }
    // Output:
    // Nigel Rees
    // Evelyn Waugh
}
```

### Supported Query Examples

#### Child Selector

```go
// Get a nested property
expr, _ := jsonpath.ParseString("$.user.name")
results := expr.Get(data)
```

#### Array Index

```go
// Get the first book
expr, _ := jsonpath.ParseString("$.store.book[0]")
results := expr.Get(data)
```

#### Array Wildcard

```go
// Get all books
expr, _ := jsonpath.ParseString("$.store.book[*]")
results := expr.Get(data)
```

#### Array Slice

```go
// Get books 0-2 (exclusive of 2)
expr, _ := jsonpath.ParseString("$.store.book[0:2]")
results := expr.Get(data)

// Get first 3 books
expr, _ := jsonpath.ParseString("$.store.book[:3]")
results := expr.Get(data)

// Get all books from index 1 onwards
expr, _ := jsonpath.ParseString("$.store.book[1:]")
results := expr.Get(data)
```

#### Recursive Descent

```go
// Find all "price" properties at any level
expr, _ := jsonpath.ParseString("$..price")
results := expr.Get(data)
```

#### Wildcard

```go
// Get all properties of an object
expr, _ := jsonpath.ParseString("$.store.*")
results := expr.Get(data)
```

#### Bracket Notation

```go
// Access property with bracket notation
expr, _ := jsonpath.ParseString("$['user']['name']")
results := expr.Get(data)
```

#### Filter Expressions

Filter expressions allow you to select array elements based on conditions.

```go
data := map[string]interface{}{
    "books": []interface{}{
        map[string]interface{}{"title": "Book A", "price": 8.99},
        map[string]interface{}{"title": "Book B", "price": 12.99},
        map[string]interface{}{"title": "Book C", "price": 5.50},
    },
}

// Get books with price less than 10
expr, _ := jsonpath.ParseString("$.books[?(@.price < 10)]")
results := expr.Get(data)
// Results: [{"title":"Book A","price":8.99}, {"title":"Book C","price":5.50}]

// Get just the titles of books under $10
expr, _ := jsonpath.ParseString("$.books[?(@.price < 10)].title")
results := expr.Get(data)
// Results: ["Book A", "Book C"]
```

**Comparison Operators:**
- `<` - Less than
- `>` - Greater than
- `<=` - Less than or equal
- `>=` - Greater than or equal
- `==` - Equal
- `!=` - Not equal

**Logical Operators:**
- `&&` - AND
- `||` - OR

**Special Operators:**
- `=~` - Regex match (supports patterns like `/pattern/` or `/pattern/i` for case-insensitive)

**Filter Examples:**

```go
// String comparison
"$.users[?(@.role == 'admin')]"

// Boolean comparison
"$.users[?(@.active == true)]"

// Numeric comparison
"$.products[?(@.stock > 20)]"

// Field existence
"$.records[?(@.email)]"

// Regex matching
"$.items[?(@.name =~ /Apple/)]"

// Case-insensitive regex
"$.items[?(@.name =~ /(?i)apple/)]"

// Multiple conditions with AND
"$.products[?(@.price < 15 && @.inStock == true)]"

// Multiple conditions with OR
"$.items[?(@.status == 'active' || @.status == 'pending')]"

// Nested field references
"$.items[?(@.details.category == 'electronics')]"
```

### Complex Example

```go
data := map[string]interface{}{
    "users": []interface{}{
        map[string]interface{}{
            "id":   1,
            "name": "John",
            "profile": map[string]interface{}{
                "age":   30,
                "email": "john@example.com",
            },
        },
        map[string]interface{}{
            "id":   2,
            "name": "Jane",
            "profile": map[string]interface{}{
                "age":   25,
                "email": "jane@example.com",
            },
        },
    },
}

// Get all user names
expr, _ := jsonpath.ParseString("$.users[*].name")
names := expr.Get(data)
// Results: ["John", "Jane"]

// Get all email addresses recursively
expr, _ := jsonpath.ParseString("$..email")
emails := expr.Get(data)
// Results: ["john@example.com", "jane@example.com"]

// Get the first user's profile
expr, _ := jsonpath.ParseString("$.users[0].profile")
profile := expr.Get(data)
// Results: [map[age:30 email:john@example.com]]
```

## API Reference

### ParseString

```go
func ParseString(query string) (Expr, error)
```

Parses a JSONPath query string into a compiled expression.

**Parameters:**
- `query`: A JSONPath query string

**Returns:**
- `Expr`: A compiled expression that can be executed
- `error`: An error if the query is invalid

### Expr Interface

```go
type Expr interface {
    Get(data interface{}) []interface{}
}
```

A compiled JSONPath expression.

#### Get Method

```go
func (e Expr) Get(data interface{}) []interface{}
```

Executes the query against the provided data.

**Parameters:**
- `data`: The data to query (typically `map[string]interface{}` or `[]interface{}`)

**Returns:**
- `[]interface{}`: A slice of all values that match the query

## Error Handling

The package returns errors for invalid queries:

```go
// Empty query
expr, err := jsonpath.ParseString("")
// Error: "query string cannot be empty"

// Invalid syntax
expr, err := jsonpath.ParseString("name")
// Error: "query must start with '$' (root selector)"

// Unclosed bracket
expr, err := jsonpath.ParseString("$[0")
// Error: "expected ']' after number"
```

## Limitations

- Negative array indices are supported (e.g., `$[-1]` gets the last element)
- Script expressions are not supported
- Step values in slices (e.g., `$[0:10:2]`) are not supported
- Filter expressions support most common use cases, but some advanced features like:
  - Comparing two field values (e.g., `[?(@.price < @.maxPrice)]`) is not yet supported
  - Filter expressions with parentheses for grouping are not supported

## Performance Considerations

- Recursive descent operations (`..`) may be slower on deeply nested structures due to the depth-first search
- The implementation uses depth limiting (max 100 levels) to prevent infinite loops
- Query parsing is fast, so you can compile queries on-the-fly or cache compiled expressions for repeated use

## Testing

Run the test suite:

```bash
go test ./pkg/jsonpath/...
```

Run with coverage:

```bash
go test ./pkg/jsonpath/... -cover
```

Current test coverage: **85.7%**

## License

This package is part of the Shape library and follows the same license terms.
