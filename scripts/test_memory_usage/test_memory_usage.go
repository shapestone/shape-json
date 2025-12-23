package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shapestone/shape-json/pkg/json"
)

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getMemStats() runtime.MemStats {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

func testParseString(filename string) error {
	fmt.Println("\n=== Testing Parse() (loads entire file into memory) ===")

	// Baseline memory
	baseline := getMemStats()
	fmt.Printf("Baseline memory: %s\n", formatBytes(baseline.Alloc))

	// Read entire file
	fmt.Println("Reading file into memory...")
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	afterRead := getMemStats()
	fmt.Printf("After reading file: %s (delta: +%s)\n",
		formatBytes(afterRead.Alloc),
		formatBytes(afterRead.Alloc-baseline.Alloc))

	// Parse
	fmt.Println("Parsing JSON...")
	start := time.Now()
	_, err = json.Parse(string(data))
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	afterParse := getMemStats()
	fmt.Printf("After parsing: %s (delta: +%s)\n",
		formatBytes(afterParse.Alloc),
		formatBytes(afterParse.Alloc-baseline.Alloc))
	fmt.Printf("Parse time: %v\n", elapsed)

	return nil
}

func testParseReader(filename string) error {
	fmt.Println("\n=== Testing ParseReader() (streaming with constant memory) ===")

	// Baseline memory
	baseline := getMemStats()
	fmt.Printf("Baseline memory: %s\n", formatBytes(baseline.Alloc))

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse with streaming
	fmt.Println("Parsing JSON with streaming...")
	start := time.Now()
	_, err = json.ParseReader(file)
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	afterParse := getMemStats()
	fmt.Printf("After parsing: %s (delta: +%s)\n",
		formatBytes(afterParse.Alloc),
		formatBytes(afterParse.Alloc-baseline.Alloc))
	fmt.Printf("Parse time: %v\n", elapsed)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_memory_usage.go <json_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Get file size
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Testing file: %s (%.2f MB)\n", filename, float64(stat.Size())/(1024*1024))

	// Test Parse() - in-memory approach
	fmt.Println("\n============================================================")
	if err := testParseString(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Parse() error: %v\n", err)
	}

	// Wait a bit and GC before next test
	time.Sleep(500 * time.Millisecond)
	runtime.GC()
	runtime.GC()

	// Test ParseReader() - streaming approach
	fmt.Println("\n============================================================")
	if err := testParseReader(filename); err != nil {
		fmt.Fprintf(os.Stderr, "ParseReader() error: %v\n", err)
	}

	fmt.Println("\n============================================================")
	fmt.Println("\n✓ Memory test complete!")
}
