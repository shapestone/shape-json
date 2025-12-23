// Package jsonpath provides a JSONPath query engine implementation.
// It supports RFC 9535 JSONPath syntax for querying JSON-like data structures.
package jsonpath

import (
	"fmt"
)

// Expr is a compiled JSONPath expression that can be executed against data.
type Expr interface {
	// Get executes the query against data and returns all matched values.
	// The data parameter should be a Go value representing JSON data
	// (map[string]interface{}, []interface{}, or primitive types).
	// Returns a slice of all values that match the query path.
	Get(data interface{}) []interface{}
}

// ParseString parses a JSONPath query string into a compiled expression.
// The query must follow RFC 9535 JSONPath syntax.
//
// Supported features:
//   - Root selector: $
//   - Child selector: .property or ['property']
//   - Wildcard: * or [*]
//   - Array index: [0], [1], etc.
//   - Array slice: [0:5], [:5], [2:], etc.
//   - Recursive descent: ..property
//   - Multiple selectors: $.a.b.c
//
// Example:
//
//	expr, err := jsonpath.ParseString("$.users[0].name")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	results := expr.Get(data)
func ParseString(query string) (Expr, error) {
	if query == "" {
		return nil, fmt.Errorf("query string cannot be empty")
	}

	// Parse the query into tokens
	tokens, err := tokenize(query)
	if err != nil {
		return nil, fmt.Errorf("tokenization failed: %w", err)
	}

	// Parse tokens into an expression
	expr, err := parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	return expr, nil
}

// expr is the internal implementation of the Expr interface
type expr struct {
	selectors []selector
}

// Get implements the Expr interface
func (e *expr) Get(data interface{}) []interface{} {
	return execute(e.selectors, data)
}

// selector represents a single path segment in a JSONPath expression
type selector interface {
	// apply applies this selector to the current values and returns new matches
	apply(current []interface{}) []interface{}
}
