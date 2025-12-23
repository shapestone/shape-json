package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/shapestone/shape-json/pkg/json"
)

// Example structs
type Person struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Email   string  `json:"email,omitempty"`
	Address Address `json:"address"`
}

type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
}

func main() {
	fmt.Println("=== shape-json: encoding/json Compatible API Demo ===")
	fmt.Println()

	// Demo 1: Marshal - Convert Go struct to JSON
	demoMarshal()

	// Demo 2: Unmarshal - Parse JSON into Go struct
	demoUnmarshal()

	// Demo 3: Encoder - Stream JSON to writer
	demoEncoder()

	// Demo 4: Decoder - Stream JSON from reader
	demoDecoder()

	// Demo 5: Struct tags
	demoStructTags()

	// Demo 6: Round trip
	demoRoundTrip()
}

func demoMarshal() {
	fmt.Println("--- Marshal Demo ---")

	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
		Address: Address{
			City:  "Seattle",
			State: "WA",
		},
	}

	data, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Marshaled JSON:\n%s\n\n", string(data))
}

func demoUnmarshal() {
	fmt.Println("--- Unmarshal Demo ---")

	jsonData := []byte(`{"name":"Bob","age":25,"address":{"city":"Portland","state":"OR"}}`)

	var person Person
	err := json.Unmarshal(jsonData, &person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Unmarshaled struct:\n")
	fmt.Printf("  Name: %s\n", person.Name)
	fmt.Printf("  Age: %d\n", person.Age)
	fmt.Printf("  City: %s\n", person.Address.City)
	fmt.Printf("  State: %s\n\n", person.Address.State)
}

func demoEncoder() {
	fmt.Println("--- Encoder Demo ---")

	people := []Person{
		{Name: "Alice", Age: 30, Address: Address{City: "Seattle", State: "WA"}},
		{Name: "Bob", Age: 25, Address: Address{City: "Portland", State: "OR"}},
		{Name: "Charlie", Age: 35, Address: Address{City: "San Francisco", State: "CA"}},
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	fmt.Println("Encoding multiple values to stream:")
	for _, p := range people {
		if err := encoder.Encode(p); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Encoded stream:\n%s", buf.String())
	fmt.Println()
}

func demoDecoder() {
	fmt.Println("--- Decoder Demo ---")

	// Single JSON object
	jsonStr := `{"name":"Alice","age":30,"address":{"city":"Seattle","state":"WA"}}`
	reader := strings.NewReader(jsonStr)
	decoder := json.NewDecoder(reader)

	var person Person
	err := decoder.Decode(&person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded from stream:\n")
	fmt.Printf("  Name: %s\n", person.Name)
	fmt.Printf("  Age: %d\n", person.Age)
	fmt.Printf("  City: %s\n\n", person.Address.City)
}

func demoStructTags() {
	fmt.Println("--- Struct Tags Demo ---")

	type User struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Password    string `json:"-"`               // Skip this field
		Email       string `json:"email,omitempty"` // Omit if empty
		Count       int    `json:"count,string"`    // Marshal as string
		NoTag       string // Uses field name
	}

	user := User{
		Username:    "alice123",
		DisplayName: "Alice",
		Password:    "secret", // Won't be marshaled
		Email:       "",       // Will be omitted
		Count:       42,       // Will be marshaled as "42"
		NoTag:       "visible",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("With struct tags:\n%s\n", string(data))
	fmt.Println("Note: Password is skipped, Email is omitted (empty), Count is string")
}

func demoRoundTrip() {
	fmt.Println("--- Round Trip Demo ---")

	original := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
		Address: Address{
			City:  "Seattle",
			State: "WA",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original marshaled: %s\n", string(data))

	// Unmarshal back to struct
	var decoded Person
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		log.Fatal(err)
	}

	// Verify
	if decoded.Name == original.Name &&
		decoded.Age == original.Age &&
		decoded.Email == original.Email &&
		decoded.Address.City == original.Address.City &&
		decoded.Address.State == original.Address.State {
		fmt.Println("✓ Round trip successful - all fields match!")
	} else {
		fmt.Println("✗ Round trip failed - data mismatch")
	}
	fmt.Println()
}
