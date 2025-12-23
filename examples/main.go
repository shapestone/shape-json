// Package main demonstrates basic JSON parsing using the low-level AST API.
//
// This example shows how to parse JSON and work with the Shape AST directly.
// For a more user-friendly approach, see the DOM Builder API example at:
// examples/dom_builder/
//
// The AST API provides fine-grained control and is useful for:
//   - Advanced parser use cases
//   - Integration with other Shape parsers
//   - Custom AST traversal and manipulation
//
// For typical JSON manipulation, use the DOM API instead:
//
//	doc, _ := json.ParseDocument(`{"name": "Alice"}`)
//	name, _ := doc.GetString("name")  // No type assertions needed!
//
// Examples covered:
//  1. Simple Object - Parse and access object properties
//  2. Nested Structures - Navigate nested objects and arrays
//  3. All Value Types - Work with strings, numbers, booleans, null, arrays, objects
//  4. Format Identification - Detect JSON format
//
// Usage:
//
//	go run examples/main.go
package main

import (
	"fmt"
	"log"

	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-json/pkg/json"
)

func main() {
	fmt.Println("=== shape-json: Basic AST Parsing Examples ===")
	fmt.Println("Note: For user-friendly API, see examples/dom_builder/")
	fmt.Println()

	// Example 1: Simple object
	fmt.Println("=== Example 1: Simple Object ===")
	simpleJSON := `{
		"name": "Alice",
		"age": 30,
		"active": true
	}`

	node, err := json.Parse(simpleJSON)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	obj := node.(*ast.ObjectNode)
	fmt.Printf("Parsed object with %d properties\n", len(obj.Properties()))

	// Access properties
	if nameNode, ok := obj.GetProperty("name"); ok {
		name := nameNode.(*ast.LiteralNode).Value().(string)
		fmt.Printf("Name: %s\n", name)
	}

	if ageNode, ok := obj.GetProperty("age"); ok {
		age := ageNode.(*ast.LiteralNode).Value().(int64)
		fmt.Printf("Age: %d\n", age)
	}

	// Example 2: Nested structures
	fmt.Println("\n=== Example 2: Nested Structures ===")
	nestedJSON := `{
		"user": {
			"name": "Bob",
			"email": "bob@example.com"
		},
		"tags": ["golang", "json", "parser"]
	}`

	node, err = json.Parse(nestedJSON)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	obj = node.(*ast.ObjectNode)

	// Access nested object
	if userNode, ok := obj.GetProperty("user"); ok {
		userObj := userNode.(*ast.ObjectNode)
		if emailNode, ok := userObj.GetProperty("email"); ok {
			email := emailNode.(*ast.LiteralNode).Value().(string)
			fmt.Printf("User email: %s\n", email)
		}
	}

	// Access array
	if tagsNode, ok := obj.GetProperty("tags"); ok {
		tagsObj := tagsNode.(*ast.ObjectNode)
		fmt.Printf("Number of tags: %d\n", len(tagsObj.Properties()))

		// Arrays are represented as objects with string keys "0", "1", etc.
		if tag0, ok := tagsObj.GetProperty("0"); ok {
			tag := tag0.(*ast.LiteralNode).Value().(string)
			fmt.Printf("First tag: %s\n", tag)
		}
	}

	// Example 3: All JSON value types
	fmt.Println("\n=== Example 3: All Value Types ===")
	allTypesJSON := `{
		"string": "hello",
		"number": 42,
		"float": 3.14,
		"boolean": true,
		"null": null,
		"array": [1, 2, 3],
		"object": {"nested": "value"}
	}`

	node, err = json.Parse(allTypesJSON)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	obj = node.(*ast.ObjectNode)
	for key, valueNode := range obj.Properties() {
		fmt.Printf("%s: %s (type: %s)\n", key, valueNode.String(), valueNode.Type())
	}

	// Example 4: Format identification
	fmt.Println("\n=== Example 4: Format ===")
	fmt.Printf("Parser format: %s\n", json.Format())
}
