package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shapestone/shape-core/pkg/ast"
	"github.com/shapestone/shape-json/pkg/json"
)

func main() {
	// Example 1: Parse from strings.Reader
	fmt.Println("=== Example 1: Parse from strings.Reader ===")
	jsonStr := `{
		"name": "Alice",
		"age": 30,
		"email": "alice@example.com"
	}`

	reader := strings.NewReader(jsonStr)
	node, err := json.ParseReader(reader)
	if err != nil {
		log.Fatalf("ParseReader error: %v", err)
	}

	obj := node.(*ast.ObjectNode)
	fmt.Printf("Parsed object with %d properties\n", len(obj.Properties()))

	if nameNode, ok := obj.GetProperty("name"); ok {
		name := nameNode.(*ast.LiteralNode).Value().(string)
		fmt.Printf("Name: %s\n", name)
	}

	// Example 2: Parse from bytes.Buffer
	fmt.Println("\n=== Example 2: Parse from bytes.Buffer ===")
	jsonBytes := []byte(`{
		"product": "Widget",
		"price": 29.99,
		"inStock": true
	}`)

	buffer := bytes.NewBuffer(jsonBytes)
	node, err = json.ParseReader(buffer)
	if err != nil {
		log.Fatalf("ParseReader error: %v", err)
	}

	obj = node.(*ast.ObjectNode)
	if priceNode, ok := obj.GetProperty("price"); ok {
		price := priceNode.(*ast.LiteralNode).Value().(float64)
		fmt.Printf("Product price: $%.2f\n", price)
	}

	// Example 3: Parse from file
	fmt.Println("\n=== Example 3: Parse from file ===")

	// Create a temporary JSON file
	tmpFile := "/tmp/example_data.json"
	fileContent := `{
		"users": [
			{"id": 1, "name": "Alice", "role": "admin"},
			{"id": 2, "name": "Bob", "role": "user"},
			{"id": 3, "name": "Charlie", "role": "user"}
		],
		"metadata": {
			"version": "1.0",
			"lastUpdated": "2025-12-08"
		}
	}`

	err = os.WriteFile(tmpFile, []byte(fileContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Parse the file using ParseReader
	file, err := os.Open(tmpFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	node, err = json.ParseReader(file)
	if err != nil {
		log.Fatalf("ParseReader error: %v", err)
	}

	obj = node.(*ast.ObjectNode)
	fmt.Printf("Parsed file successfully\n")

	// Access nested data
	if usersNode, ok := obj.GetProperty("users"); ok {
		usersArray := usersNode.(*ast.ObjectNode)
		fmt.Printf("Number of users: %d\n", len(usersArray.Properties()))

		// Get first user
		if user0, ok := usersArray.GetProperty("0"); ok {
			userObj := user0.(*ast.ObjectNode)
			if nameNode, ok := userObj.GetProperty("name"); ok {
				name := nameNode.(*ast.LiteralNode).Value().(string)
				fmt.Printf("First user: %s\n", name)
			}
		}
	}

	// Example 4: Comparison - Parse vs ParseReader
	fmt.Println("\n=== Example 4: Parse vs ParseReader ===")
	testJSON := `{"message": "Hello, World!"}`

	// Using Parse (for small strings in memory)
	node1, err := json.Parse(testJSON)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	obj1 := node1.(*ast.ObjectNode)
	msg1, _ := obj1.GetProperty("message")
	fmt.Printf("Parse() result: %v\n", msg1.(*ast.LiteralNode).Value())

	// Using ParseReader (for streams and large files)
	reader2 := strings.NewReader(testJSON)
	node2, err := json.ParseReader(reader2)
	if err != nil {
		log.Fatalf("ParseReader error: %v", err)
	}
	obj2 := node2.(*ast.ObjectNode)
	msg2, _ := obj2.GetProperty("message")
	fmt.Printf("ParseReader() result: %v\n", msg2.(*ast.LiteralNode).Value())

	// Both produce identical results
	fmt.Println("\nBoth methods produce the same AST representation")

	// Example 5: When to use ParseReader
	fmt.Println("\n=== Example 5: When to use ParseReader ===")
	fmt.Println("Use ParseReader when:")
	fmt.Println("  - Reading JSON from files")
	fmt.Println("  - Processing large JSON documents that don't fit in memory")
	fmt.Println("  - Streaming JSON from network connections")
	fmt.Println("  - Reading from compressed streams (gzip, etc.)")
	fmt.Println("  - Any io.Reader source")
	fmt.Println("\nUse Parse when:")
	fmt.Println("  - JSON is already in memory as a string")
	fmt.Println("  - Working with small JSON documents")
	fmt.Println("  - Simplicity is preferred over streaming")
}
