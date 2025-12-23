package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shapestone/shape-json/pkg/json"
)

func main() {
	fmt.Println("=== JSON Validation Examples ===")

	// Example 1: Simple Validation (Idiomatic Go)
	fmt.Println("Example 1: Simple Validation (Idiomatic Go)")
	if err := json.Validate(`{"name": "Alice", "age": 30}`); err != nil {
		fmt.Printf("  ❌ Invalid JSON: %v\n", err)
	} else {
		fmt.Println("  ✅ Valid JSON!")
	}

	// Example 2: Validation with Error Details
	fmt.Println("\nExample 2: Validation with Error Details")
	invalidJSON := `{"name": "Alice"`  // Missing closing brace
	if err := json.Validate(invalidJSON); err != nil {
		fmt.Printf("  ❌ Invalid JSON: %v\n", err)
	} else {
		fmt.Println("  ✅ Valid JSON!")
	}

	// Example 3: Validate Different Types
	fmt.Println("\nExample 3: Validate Different JSON Types")
	examples := []struct {
		name  string
		input string
	}{
		{"Object", `{"key": "value"}`},
		{"Array", `[1, 2, 3]`},
		{"String", `"hello"`},
		{"Number", `42`},
		{"Boolean", `true`},
		{"Null", `null`},
	}

	for _, ex := range examples {
		if err := json.Validate(ex.input); err != nil {
			fmt.Printf("  ❌ %s: Invalid - %v\n", ex.name, err)
		} else {
			fmt.Printf("  ✅ %s: Valid\n", ex.name)
		}
	}

	// Example 4: Detect Invalid JSON
	fmt.Println("\nExample 4: Common JSON Errors")
	invalidInputs := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Plain text", "this is not json"},
		{"Missing quotes on key", `{name: "value"}`},
		{"Trailing comma in object", `{"key": "value",}`},
		{"Unclosed brace", `{"key": "value"`},
		{"Single quotes", `{'key': 'value'}`},
		{"Trailing comma in array", `[1, 2, 3,]`},
	}

	for _, test := range invalidInputs {
		if err := json.Validate(test.input); err != nil {
			fmt.Printf("  ✅ %s: Correctly rejected\n", test.name)
			fmt.Printf("     Error: %v\n", err)
		} else {
			fmt.Printf("  ❌ %s: Incorrectly accepted\n", test.name)
		}
	}

	// Example 5: Validate from io.Reader
	fmt.Println("\nExample 5: Validation from io.Reader")
	jsonData := `{
		"users": [
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"}
		],
		"total": 2
	}`
	reader := strings.NewReader(jsonData)
	if err := json.ValidateReader(reader); err != nil {
		fmt.Printf("  ❌ Reader contains invalid JSON: %v\n", err)
	} else {
		fmt.Println("  ✅ Reader contains valid JSON")
	}

	// Example 6: Validate from File
	fmt.Println("\nExample 6: Validation from File")
	testFile := "/tmp/test-shape-json.json"
	content := []byte(`{"message": "Hello from file", "timestamp": 1234567890}`)
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		fmt.Printf("  ❌ Error creating test file: %v\n", err)
	} else {
		defer os.Remove(testFile)

		file, err := os.Open(testFile)
		if err != nil {
			fmt.Printf("  ❌ Error opening file: %v\n", err)
		} else {
			defer file.Close()

			if err := json.ValidateReader(file); err != nil {
				fmt.Printf("  ❌ File contains invalid JSON: %v\n", err)
			} else {
				fmt.Println("  ✅ File contains valid JSON")
			}
		}
	}

	// Example 7: Validate Before Parsing (Best Practice)
	fmt.Println("\nExample 7: Validate Before Parsing (Recommended)")
	userInput := `{"username": "alice", "email": "alice@example.com"}`

	// Validate first
	if err := json.Validate(userInput); err != nil {
		fmt.Printf("  ❌ Invalid JSON, not attempting to parse: %v\n", err)
		return
	}

	fmt.Println("  ✅ Valid JSON detected, proceeding to parse...")

	// Now safely parse using DOM API
	doc, err := json.ParseDocument(userInput)
	if err != nil {
		fmt.Printf("  ❌ Parse error: %v\n", err)
	} else {
		username, _ := doc.GetString("username")
		email, _ := doc.GetString("email")
		fmt.Printf("  ✅ Successfully parsed: username=%s, email=%s\n", username, email)
	}

	// Example 8: Integration Testing Pattern
	fmt.Println("\nExample 8: API Input Validation Pattern")

	// Simulate API inputs
	apiInputs := []struct {
		name  string
		input string
	}{
		{"Valid user request", `{"action": "create", "user": {"name": "Alice"}}`},
		{"Invalid JSON syntax", `{"action": "create"` },
		{"Valid but incomplete", `{"action": "delete"}`},
	}

	for _, test := range apiInputs {
		fmt.Printf("  Testing: %s\n", test.name)

		// Quick validation - idiomatic Go
		if err := json.Validate(test.input); err != nil {
			fmt.Printf("    ❌ Rejected: Invalid JSON syntax - %v\n", err)
			continue
		}

		fmt.Printf("    ✅ Accepted: Valid JSON syntax\n")

		// Parse and process
		doc, _ := json.ParseDocument(test.input)
		action, _ := doc.GetString("action")
		fmt.Printf("    Action: %s\n", action)
	}

	fmt.Println("\n=== All Examples Completed ===")
}
