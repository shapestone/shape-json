package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_large_json.go <output_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Generate a large JSON array with many objects
	// Target: ~100MB file
	// Each object is roughly 200 bytes, so we need ~500,000 objects

	file.WriteString("[\n")

	for i := 0; i < 500000; i++ {
		if i > 0 {
			file.WriteString(",\n")
		}

		// Write a JSON object with various fields
		fmt.Fprintf(file, `  {
    "id": %d,
    "name": "User %d",
    "email": "user%d@example.com",
    "age": %d,
    "active": %v,
    "score": %d.%d,
    "tags": ["tag1", "tag2", "tag3"],
    "metadata": {
      "created": "2024-01-01T00:00:00Z",
      "updated": "2024-12-08T00:00:00Z"
    }
  }`, i, i, i, 20+(i%60), i%2 == 0, i%100, i%100)

		// Progress indicator
		if i > 0 && i%10000 == 0 {
			fmt.Printf("Generated %d objects...\n", i)
		}
	}

	file.WriteString("\n]\n")

	// Get file size
	stat, _ := file.Stat()
	sizeMB := float64(stat.Size()) / (1024 * 1024)
	fmt.Printf("\nGenerated %s: %.2f MB\n", filename, sizeMB)
}
