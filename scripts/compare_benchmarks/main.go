package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// BenchmarkMetadata contains information about a benchmark run
type BenchmarkMetadata struct {
	Timestamp   string `json:"timestamp"`
	GitCommit   string `json:"commit"`
	Platform    string `json:"platform"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	GoVersion   string `json:"go_version"`
	BenchTime   string `json:"bench_time"`
	Description string `json:"description"`
}

// HistoryEntry represents a benchmark history entry
type HistoryEntry struct {
	Timestamp    string
	Path         string
	Metadata     *BenchmarkMetadata
	BenchmarkDir string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <old> <new>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  <old>, <new>: timestamp, 'latest', 'previous', or path to benchmark_output.txt\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s latest previous          # Compare latest vs previous run\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 2025-12-20_14-30-00 2025-12-21_16-45-00\n", os.Args[0])
		os.Exit(1)
	}

	oldArg := os.Args[1]
	newArg := os.Args[2]

	// Find project root
	projectRoot := findProjectRoot(".")
	if projectRoot == "" {
		fatal("Could not find project root (looking for go.mod)")
	}

	historyDir := filepath.Join(projectRoot, "benchmarks", "history")

	// Get benchmark file paths
	oldPath, err := resolveBenchmarkPath(historyDir, oldArg)
	if err != nil {
		fatal("Failed to resolve old benchmark: %v", err)
	}

	newPath, err := resolveBenchmarkPath(historyDir, newArg)
	if err != nil {
		fatal("Failed to resolve new benchmark: %v", err)
	}

	// Load metadata for both runs
	oldMeta := loadMetadata(filepath.Dir(oldPath))
	newMeta := loadMetadata(filepath.Dir(newPath))

	// Display comparison header
	fmt.Println("Benchmark Comparison")
	fmt.Println("====================")
	fmt.Println()
	displayRunInfo("Old", oldMeta, oldPath)
	fmt.Println()
	displayRunInfo("New", newMeta, newPath)
	fmt.Println()

	// Check if benchstat is available
	if !commandExists("benchstat") {
		fmt.Println("Warning: benchstat not found. Install with:")
		fmt.Println("  go install golang.org/x/perf/cmd/benchstat@latest")
		fmt.Println()
		fmt.Println("Showing simple comparison instead:")
		fmt.Println()
		showSimpleComparison(oldPath, newPath)
		return
	}

	// Run benchstat comparison
	fmt.Println("Statistical Comparison (benchstat)")
	fmt.Println("==================================")
	fmt.Println()

	cmd := exec.Command("benchstat", oldPath, newPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fatal("Failed to run benchstat: %v", err)
	}

	fmt.Println()
	fmt.Println("Interpretation:")
	fmt.Println("  ~ means no significant change")
	fmt.Println("  + means new is slower (regression)")
	fmt.Println("  - means new is faster (improvement)")
	fmt.Println()
}

// findProjectRoot walks up the directory tree to find go.mod
func findProjectRoot(startDir string) string {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "" // reached root without finding go.mod
		}
		dir = parent
	}
}

// resolveBenchmarkPath resolves a benchmark argument to a file path
func resolveBenchmarkPath(historyDir, arg string) (string, error) {
	// If it's a file path, use it directly
	if _, err := os.Stat(arg); err == nil {
		return arg, nil
	}

	// Handle special keywords
	if arg == "latest" || arg == "previous" {
		entries, err := listHistoryEntries(historyDir)
		if err != nil {
			return "", err
		}

		if len(entries) == 0 {
			return "", fmt.Errorf("no benchmark history found")
		}

		var entry *HistoryEntry
		if arg == "latest" {
			entry = entries[0]
		} else if arg == "previous" {
			if len(entries) < 2 {
				return "", fmt.Errorf("no previous benchmark found (only one entry in history)")
			}
			entry = entries[1]
		}

		return filepath.Join(entry.BenchmarkDir, "benchmark_output.txt"), nil
	}

	// Assume it's a timestamp
	benchPath := filepath.Join(historyDir, arg, "benchmark_output.txt")
	if _, err := os.Stat(benchPath); err != nil {
		return "", fmt.Errorf("benchmark not found: %s", benchPath)
	}

	return benchPath, nil
}

// listHistoryEntries returns all benchmark history entries sorted by timestamp (newest first)
func listHistoryEntries(historyDir string) ([]*HistoryEntry, error) {
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("history directory does not exist: %s", historyDir)
	}

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, err
	}

	var results []*HistoryEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		timestamp := entry.Name()
		benchDir := filepath.Join(historyDir, timestamp)

		// Check if benchmark_output.txt exists
		benchPath := filepath.Join(benchDir, "benchmark_output.txt")
		if _, err := os.Stat(benchPath); err != nil {
			continue
		}

		metadata := loadMetadata(benchDir)

		results = append(results, &HistoryEntry{
			Timestamp:    timestamp,
			Path:         benchPath,
			Metadata:     metadata,
			BenchmarkDir: benchDir,
		})
	}

	// Sort by timestamp descending (newest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp > results[j].Timestamp
	})

	return results, nil
}

// loadMetadata loads metadata from a benchmark directory
func loadMetadata(benchDir string) *BenchmarkMetadata {
	metaPath := filepath.Join(benchDir, "metadata.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil
	}

	var metadata BenchmarkMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil
	}

	return &metadata
}

// displayRunInfo displays information about a benchmark run
func displayRunInfo(label string, meta *BenchmarkMetadata, path string) {
	fmt.Printf("%s: %s\n", label, filepath.Dir(path))

	if meta != nil {
		if meta.GitCommit != "" {
			commitShort := meta.GitCommit
			if len(commitShort) > 12 {
				commitShort = commitShort[:12]
			}
			fmt.Printf("  Commit: %s\n", commitShort)
		}
		if meta.Platform != "" {
			fmt.Printf("  Platform: %s (%s/%s)\n", meta.Platform, meta.OS, meta.Arch)
		}
		if meta.GoVersion != "" {
			fmt.Printf("  Go: %s\n", meta.GoVersion)
		}
		if meta.Description != "" {
			fmt.Printf("  Description: %s\n", meta.Description)
		}
	}
}

// showSimpleComparison shows a simple file comparison
func showSimpleComparison(oldPath, newPath string) {
	fmt.Println("Old benchmark results:")
	fmt.Println("----------------------")
	showBenchmarkLines(oldPath)

	fmt.Println()
	fmt.Println("New benchmark results:")
	fmt.Println("----------------------")
	showBenchmarkLines(newPath)
}

// showBenchmarkLines displays benchmark result lines from a file
func showBenchmarkLines(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", path, err)
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Benchmark") {
			fmt.Println(line)
		}
	}
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
