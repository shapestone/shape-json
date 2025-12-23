// Package main demonstrates the fluent DOM API for JSON manipulation.
//
// This example shows how to build and manipulate JSON documents without
// type assertions, using a clean, chainable API.
package main

import (
	"fmt"
	"log"

	"github.com/shapestone/shape-json/pkg/json"
)

func main() {
	fmt.Println("=== JSON DOM Builder Examples ===")

	// Example 1: Simple Document Building
	example1SimpleBuilder()

	// Example 2: Reading and Type-Safe Access
	example2TypeSafeAccess()

	// Example 3: Complex Nested Structures
	example3NestedStructures()

	// Example 4: Working with Arrays
	example4Arrays()

	// Example 5: Round-Trip (Build, Marshal, Parse)
	example5RoundTrip()

	// Example 6: Comparison with Old AST API
	example6Comparison()
}

// Example 1: Simple document building with method chaining
func example1SimpleBuilder() {
	fmt.Println("=== Example 1: Simple Document Building ===")

	// Build a document using fluent API (method chaining)
	doc := json.NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetBool("active", true).
		SetFloat("score", 95.5).
		SetNull("middle_name")

	jsonStr, _ := doc.JSON()
	fmt.Println("Built document:")
	fmt.Println(jsonStr)
	fmt.Println()
}

// Example 2: Reading with type-safe access (no type assertions!)
func example2TypeSafeAccess() {
	fmt.Println("=== Example 2: Type-Safe Access ===")

	jsonStr := `{"name":"Bob","age":25,"active":true,"score":88.5}`

	// Parse into DOM
	doc, err := json.ParseDocument(jsonStr)
	if err != nil {
		log.Fatal(err)
	}

	// Access values with type-safe getters (no type assertions!)
	name, _ := doc.GetString("name")
	age, _ := doc.GetInt("age")
	active, _ := doc.GetBool("active")
	score, _ := doc.GetFloat("score")

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Active: %v\n", active)
	fmt.Printf("Score: %.1f\n", score)

	// Check for missing keys safely
	if _, ok := doc.GetString("email"); !ok {
		fmt.Println("Email: (not provided)")
	}

	fmt.Println()
}

// Example 3: Complex nested structures
func example3NestedStructures() {
	fmt.Println("=== Example 3: Complex Nested Structures ===")

	// Build a complex nested document
	doc := json.NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetObject("address", json.NewDocument().
			SetString("street", "123 Main St").
			SetString("city", "NYC").
			SetString("zip", "10001").
			SetObject("coordinates", json.NewDocument().
				SetFloat("lat", 40.7128).
				SetFloat("lng", -74.0060))).
		SetArray("tags", json.NewArray().
			AddString("engineer").
			AddString("golang").
			AddString("opensource"))

	jsonStr, _ := doc.JSON()
	fmt.Println("Complex document:")
	fmt.Println(jsonStr)

	// Access nested values
	if addr, ok := doc.GetObject("address"); ok {
		city, _ := addr.GetString("city")
		fmt.Printf("\nCity: %s\n", city)

		if coords, ok := addr.GetObject("coordinates"); ok {
			lat, _ := coords.GetFloat("lat")
			lng, _ := coords.GetFloat("lng")
			fmt.Printf("Coordinates: %.4f, %.4f\n", lat, lng)
		}
	}

	fmt.Println()
}

// Example 4: Working with arrays
func example4Arrays() {
	fmt.Println("=== Example 4: Working with Arrays ===")

	// Build an array
	arr := json.NewArray().
		AddString("apple").
		AddString("banana").
		AddString("cherry")

	jsonStr, _ := arr.JSON()
	fmt.Println("Simple array:")
	fmt.Println(jsonStr)

	// Array of objects
	users := json.NewArray().
		AddObject(json.NewDocument().
			SetString("name", "Alice").
			SetInt("age", 30)).
		AddObject(json.NewDocument().
			SetString("name", "Bob").
			SetInt("age", 25)).
		AddObject(json.NewDocument().
			SetString("name", "Charlie").
			SetInt("age", 35))

	// Access array elements
	fmt.Println("\nUsers:")
	for i := 0; i < users.Len(); i++ {
		if user, ok := users.GetObject(i); ok {
			name, _ := user.GetString("name")
			age, _ := user.GetInt("age")
			fmt.Printf("  %d. %s (age %d)\n", i+1, name, age)
		}
	}

	// Mixed type array
	mixed := json.NewArray().
		AddString("text").
		AddInt(42).
		AddBool(true).
		AddFloat(3.14).
		AddNull().
		AddObject(json.NewDocument().SetString("key", "value")).
		AddArray(json.NewArray().AddInt(1).AddInt(2))

	jsonStr, _ = mixed.JSON()
	fmt.Println("\nMixed type array:")
	fmt.Println(jsonStr)
	fmt.Println()
}

// Example 5: Round-trip (build, marshal, parse, access)
func example5RoundTrip() {
	fmt.Println("=== Example 5: Round-Trip Example ===")

	// Build a user profile
	original := json.NewDocument().
		SetString("username", "alice123").
		SetString("email", "alice@example.com").
		SetObject("preferences", json.NewDocument().
			SetBool("notifications", true).
			SetString("theme", "dark").
			SetInt("fontSize", 14)).
		SetArray("roles", json.NewArray().
			AddString("admin").
			AddString("editor"))

	// Marshal to JSON string
	jsonStr, err := original.JSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original document:")
	fmt.Println(jsonStr)

	// Parse back from JSON
	parsed, err := json.ParseDocument(jsonStr)
	if err != nil {
		log.Fatal(err)
	}

	// Access parsed values
	username, _ := parsed.GetString("username")
	fmt.Printf("\nUsername: %s\n", username)

	if prefs, ok := parsed.GetObject("preferences"); ok {
		theme, _ := prefs.GetString("theme")
		fontSize, _ := prefs.GetInt("fontSize")
		fmt.Printf("Theme: %s (font size: %d)\n", theme, fontSize)
	}

	if roles, ok := parsed.GetArray("roles"); ok {
		fmt.Print("Roles: ")
		for i := 0; i < roles.Len(); i++ {
			role, _ := roles.GetString(i)
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(role)
		}
		fmt.Println()
	}

	fmt.Println()
}

// Example 6: Comparison with old AST API
func example6Comparison() {
	fmt.Println("=== Example 6: Old AST API vs New DOM API ===")

	jsonStr := `{"user":{"name":"Alice","tags":["go","json"]}}`

	fmt.Println("Input JSON:", jsonStr)
	fmt.Println()

	// OLD WAY: Using AST (confusing, type assertions everywhere)
	fmt.Println("OLD AST API (confusing):")
	fmt.Println("  node, _ := json.Parse(input)")
	fmt.Println("  obj := node.(*ast.ObjectNode)              // Type assertion #1")
	fmt.Println("  userNode, _ := obj.GetProperty(\"user\")")
	fmt.Println("  userObj := userNode.(*ast.ObjectNode)      // Type assertion #2")
	fmt.Println("  nameNode, _ := userObj.GetProperty(\"name\")")
	fmt.Println("  name := nameNode.(*ast.LiteralNode).Value().(string) // Type assertions #3 and #4")
	fmt.Println("  // Result: name = \"Alice\"")
	fmt.Println()

	// NEW WAY: Using DOM API (clean, no type assertions)
	fmt.Println("NEW DOM API (clean):")
	fmt.Println("  doc, _ := json.ParseDocument(input)")
	fmt.Println("  user, _ := doc.GetObject(\"user\")")
	fmt.Println("  name, _ := user.GetString(\"name\")")
	fmt.Println("  // Result: name = \"Alice\"")
	fmt.Println()

	// Actually demonstrate it
	doc, _ := json.ParseDocument(jsonStr)
	user, _ := doc.GetObject("user")
	name, _ := user.GetString("name")
	tags, _ := user.GetArray("tags")
	firstTag, _ := tags.GetString(0)

	fmt.Printf("Actual result: name=%s, first tag=%s\n", name, firstTag)
	fmt.Println()

	// Building is also cleaner
	fmt.Println("Building JSON:")
	fmt.Println()

	fmt.Println("OLD WAY: Complex AST manipulation")
	fmt.Println("  (requires creating AST nodes manually)")
	fmt.Println()

	fmt.Println("NEW WAY: Fluent builder pattern")
	built := json.NewDocument().
		SetString("name", "Alice").
		SetInt("age", 30).
		SetObject("address", json.NewDocument().
			SetString("city", "NYC"))

	jsonStr, _ = built.JSON()
	fmt.Printf("  Result: %s\n", jsonStr)
}
