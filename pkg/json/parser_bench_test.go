package json_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	shapejson "github.com/shapestone/shape-json/pkg/json"
)

// Benchmark test data is loaded once and reused across all benchmarks
var (
	smallJSON  string
	mediumJSON string
	largeJSON  string
)

// loadBenchmarkData loads test data files once during test initialization
func loadBenchmarkData() error {
	if smallJSON != "" {
		return nil // already loaded
	}

	testdataDir := filepath.Join("..", "..", "testdata", "benchmarks")

	// Load small.json
	smallBytes, err := os.ReadFile(filepath.Join(testdataDir, "small.json"))
	if err != nil {
		return err
	}
	smallJSON = string(smallBytes)

	// Load medium.json
	mediumBytes, err := os.ReadFile(filepath.Join(testdataDir, "medium.json"))
	if err != nil {
		return err
	}
	mediumJSON = string(mediumBytes)

	// Load large.json
	largeBytes, err := os.ReadFile(filepath.Join(testdataDir, "large.json"))
	if err != nil {
		return err
	}
	largeJSON = string(largeBytes)

	return nil
}

// ================================
// Shape-JSON Benchmarks
// ================================

// BenchmarkShapeJSON_Parse_Small benchmarks parsing of small JSON (~100 bytes)
// using the shape-json parser.
func BenchmarkShapeJSON_Parse_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, err := shapejson.Parse(smallJSON)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = node
	}
}

// BenchmarkShapeJSON_Parse_Medium benchmarks parsing of medium JSON (~10KB)
// using the shape-json parser.
func BenchmarkShapeJSON_Parse_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(mediumJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, err := shapejson.Parse(mediumJSON)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = node
	}
}

// BenchmarkShapeJSON_Parse_Large benchmarks parsing of large JSON (~1MB)
// using the shape-json parser.
func BenchmarkShapeJSON_Parse_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, err := shapejson.Parse(largeJSON)
		if err != nil {
			b.Fatal(err)
		}
		// Release nodes back to pool for reuse (enables AST node pooling)
		shapejson.ReleaseTree(node)
	}
}

// BenchmarkShapeJSON_ParseReader_Small benchmarks parsing of small JSON (~100 bytes)
// using the shape-json ParseReader streaming parser.
func BenchmarkShapeJSON_ParseReader_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(smallJSON)
		node, err := shapejson.ParseReader(reader)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = node
	}
}

// BenchmarkShapeJSON_ParseReader_Medium benchmarks parsing of medium JSON (~10KB)
// using the shape-json ParseReader streaming parser.
func BenchmarkShapeJSON_ParseReader_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(mediumJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(mediumJSON)
		node, err := shapejson.ParseReader(reader)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = node
	}
}

// BenchmarkShapeJSON_ParseReader_Large benchmarks parsing of large JSON (~1MB)
// using the shape-json ParseReader streaming parser.
func BenchmarkShapeJSON_ParseReader_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(largeJSON)
		node, err := shapejson.ParseReader(reader)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = node
	}
}

// ================================
// encoding/json Benchmarks (Comparison)
// ================================

// BenchmarkEncodingJSON_Parse_Small benchmarks parsing of small JSON (~100 bytes)
// using the standard library encoding/json package.
// Unmarshals into interface{} to match shape-json's AST output behavior.
func BenchmarkEncodingJSON_Parse_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := json.Unmarshal([]byte(smallJSON), &v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkEncodingJSON_Parse_Medium benchmarks parsing of medium JSON (~10KB)
// using the standard library encoding/json package.
// Unmarshals into interface{} to match shape-json's AST output behavior.
func BenchmarkEncodingJSON_Parse_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(mediumJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := json.Unmarshal([]byte(mediumJSON), &v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkEncodingJSON_Parse_Large benchmarks parsing of large JSON (~1MB)
// using the standard library encoding/json package.
// Unmarshals into interface{} to match shape-json's AST output behavior.
func BenchmarkEncodingJSON_Parse_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := json.Unmarshal([]byte(largeJSON), &v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkEncodingJSON_Decoder_Small benchmarks parsing of small JSON (~100 bytes)
// using encoding/json's streaming Decoder (similar to shape-json's ParseReader).
func BenchmarkEncodingJSON_Decoder_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(smallJSON)
		decoder := json.NewDecoder(reader)
		var v interface{}
		err := decoder.Decode(&v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkEncodingJSON_Decoder_Medium benchmarks parsing of medium JSON (~10KB)
// using encoding/json's streaming Decoder (similar to shape-json's ParseReader).
func BenchmarkEncodingJSON_Decoder_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(mediumJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(mediumJSON)
		decoder := json.NewDecoder(reader)
		var v interface{}
		err := decoder.Decode(&v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkEncodingJSON_Decoder_Large benchmarks parsing of large JSON (~1MB)
// using encoding/json's streaming Decoder (similar to shape-json's ParseReader).
func BenchmarkEncodingJSON_Decoder_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(largeJSON)
		decoder := json.NewDecoder(reader)
		var v interface{}
		err := decoder.Decode(&v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// ================================
// Dual-Path Comparison Benchmarks
// ================================

// BenchmarkUnmarshal_FastPath_Large benchmarks unmarshaling large JSON (~1MB)
// using the new fast path (no AST construction).
func BenchmarkUnmarshal_FastPath_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(largeJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.Unmarshal(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkUnmarshal_ASTPath_Large benchmarks unmarshaling large JSON (~1MB)
// using the AST path (for comparison with fast path).
func BenchmarkUnmarshal_ASTPath_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(largeJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.UnmarshalWithAST(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = v
	}
}

// BenchmarkUnmarshal_FastPath_Struct benchmarks unmarshaling into a struct
// using the fast path.
func BenchmarkUnmarshal_FastPath_Struct(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	type Person struct {
		Name   string
		Age    int
		Active bool
	}

	jsonData := []byte(`{"Name": "Alice", "Age": 30, "Active": true}`)
	b.SetBytes(int64(len(jsonData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Person
		err := shapejson.Unmarshal(jsonData, &p)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = p
	}
}

// BenchmarkUnmarshal_ASTPath_Struct benchmarks unmarshaling into a struct
// using the AST path (for comparison).
func BenchmarkUnmarshal_ASTPath_Struct(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	type Person struct {
		Name   string
		Age    int
		Active bool
	}

	jsonData := []byte(`{"Name": "Alice", "Age": 30, "Active": true}`)
	b.SetBytes(int64(len(jsonData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Person
		err := shapejson.UnmarshalWithAST(jsonData, &p)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		_ = p
	}
}

// BenchmarkValidate_FastPath benchmarks validation using fast path.
func BenchmarkValidate_FastPath(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := shapejson.Validate(largeJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValidate_ASTPath benchmarks validation using AST path (via Parse).
func BenchmarkValidate_ASTPath(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	b.SetBytes(int64(len(largeJSON)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shapejson.Parse(largeJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshal_FastPath_Small benchmarks unmarshaling small JSON using fast path.
func BenchmarkUnmarshal_FastPath_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(smallJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.Unmarshal(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// BenchmarkUnmarshal_ASTPath_Small benchmarks unmarshaling small JSON using AST path.
func BenchmarkUnmarshal_ASTPath_Small(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(smallJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.UnmarshalWithAST(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// BenchmarkUnmarshal_FastPath_Medium benchmarks unmarshaling medium JSON using fast path.
func BenchmarkUnmarshal_FastPath_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(mediumJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.Unmarshal(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// BenchmarkUnmarshal_ASTPath_Medium benchmarks unmarshaling medium JSON using AST path.
func BenchmarkUnmarshal_ASTPath_Medium(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}

	data := []byte(mediumJSON)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := shapejson.UnmarshalWithAST(data, &v)
		if err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}
