package json_test

import (
	"encoding/json"
	"testing"
	"time"

	shapejson "github.com/shapestone/shape-json/pkg/json"
)

// ================================
// Marshal Benchmarks
// ================================

// benchStruct is a representative struct used for marshal benchmarks.
type benchStruct struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Active bool    `json:"active"`
	Score  float64 `json:"score"`
}

// benchNestedStruct is a more complex struct for medium-complexity marshal benchmarks.
type benchNestedStruct struct {
	ID       int               `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
	Profile  struct {
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		Age       int     `json:"age"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"profile"`
}

var (
	benchSimple = benchStruct{
		Name:   "Alice",
		Age:    30,
		Active: true,
		Score:  95.5,
	}

	benchNested = benchNestedStruct{
		ID:       1,
		Username: "alice_wonder",
		Email:    "alice@example.com",
		Tags:     []string{"admin", "user", "moderator"},
		Metadata: map[string]string{
			"theme":    "dark",
			"language": "en",
			"timezone": "PST",
		},
		Profile: struct {
			FirstName string  `json:"firstName"`
			LastName  string  `json:"lastName"`
			Age       int     `json:"age"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			FirstName: "Alice",
			LastName:  "Wonder",
			Age:       30,
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
)

// --- Shape-JSON Marshal ---

func BenchmarkShapeJSON_Marshal_Simple(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchSimple)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Marshal_Nested(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchNested)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Marshal_Map(b *testing.B) {
	m := map[string]interface{}{
		"name": "Alice", "age": 30, "active": true, "score": 95.5,
		"tags": []interface{}{"admin", "user"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(m)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Marshal_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}
	// Unmarshal large JSON into interface{} then re-marshal
	var v interface{}
	if err := shapejson.Unmarshal([]byte(largeJSON), &v); err != nil {
		b.Fatalf("Failed to prepare data: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(v)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

// --- encoding/json v1 Marshal ---

func BenchmarkEncodingJSON_Marshal_Simple(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(benchSimple)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Nested(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(benchNested)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Map(b *testing.B) {
	m := map[string]interface{}{
		"name": "Alice", "age": 30, "active": true, "score": 95.5,
		"tags": []interface{}{"admin", "user"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(m)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Large(b *testing.B) {
	if err := loadBenchmarkData(); err != nil {
		b.Fatalf("Failed to load benchmark data: %v", err)
	}
	var v interface{}
	if err := json.Unmarshal([]byte(largeJSON), &v); err != nil {
		b.Fatalf("Failed to prepare data: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(v)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

// ================================
// Type-Specific Benchmarks
// ================================

// --- Float32 ---

type benchFloat32Struct struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

var benchFloat32 = benchFloat32Struct{X: 1.1, Y: 2.2, Z: 3.3}
var benchFloat32JSON = []byte(`{"x":1.1,"y":2.2,"z":3.3}`)

func BenchmarkShapeJSON_Marshal_Float32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchFloat32)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Float32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(benchFloat32)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Unmarshal_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchFloat32Struct
		if err := shapejson.Unmarshal(benchFloat32JSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

func BenchmarkEncodingJSON_Unmarshal_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchFloat32Struct
		if err := json.Unmarshal(benchFloat32JSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// --- Uint64 ---

type benchUint64Struct struct {
	ID    uint64 `json:"id"`
	Count uint64 `json:"count"`
}

var benchUint64 = benchUint64Struct{ID: 18446744073709551615, Count: 9223372036854775808}
var benchUint64JSON = []byte(`{"id":18446744073709551615,"count":9223372036854775808}`)

func BenchmarkShapeJSON_Marshal_Uint64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchUint64)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Uint64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(benchUint64)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Unmarshal_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchUint64Struct
		if err := shapejson.Unmarshal(benchUint64JSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

func BenchmarkEncodingJSON_Unmarshal_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchUint64Struct
		if err := json.Unmarshal(benchUint64JSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// --- time.Time ---

type benchTimeStruct struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

var benchTime = benchTimeStruct{
	Created: time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC),
	Updated: time.Date(2025, 6, 15, 14, 45, 30, 123456789, time.UTC),
}
var benchTimeJSON = []byte(`{"created":"2025-06-15T10:30:00Z","updated":"2025-06-15T14:45:30.123456789Z"}`)

func BenchmarkShapeJSON_Marshal_Time(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchTime)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_Time(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(benchTime)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Unmarshal_Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchTimeStruct
		if err := shapejson.Unmarshal(benchTimeJSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

func BenchmarkEncodingJSON_Unmarshal_Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchTimeStruct
		if err := json.Unmarshal(benchTimeJSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// --- time.Duration ---

type benchDurationStruct struct {
	Timeout  time.Duration `json:"timeout"`
	Interval time.Duration `json:"interval"`
}

var benchDuration = benchDurationStruct{
	Timeout:  5*time.Second + 500*time.Millisecond,
	Interval: 1*time.Hour + 30*time.Minute,
}
var benchDurationJSON = []byte(`{"timeout":"PT5.5S","interval":"PT1H30M"}`)

func BenchmarkShapeJSON_Marshal_Duration(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(benchDuration)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Unmarshal_Duration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v benchDurationStruct
		if err := shapejson.Unmarshal(benchDurationJSON, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// --- UTF-8 String ---

func BenchmarkShapeJSON_Marshal_UTF8(b *testing.B) {
	b.ReportAllocs()
	data := map[string]string{
		"greeting": "Hello, \u4e16\u754c!",
		"emoji":    "\U0001f600\U0001f680\U0001f30d",
		"mixed":    "caf\u00e9 na\u00efve r\u00e9sum\u00e9",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := shapejson.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkEncodingJSON_Marshal_UTF8(b *testing.B) {
	b.ReportAllocs()
	data := map[string]string{
		"greeting": "Hello, \u4e16\u754c!",
		"emoji":    "\U0001f600\U0001f680\U0001f30d",
		"mixed":    "caf\u00e9 na\u00efve r\u00e9sum\u00e9",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkShapeJSON_Unmarshal_UTF8(b *testing.B) {
	data := []byte(`{"greeting":"Hello, \u4e16\u754c!","emoji":"\ud83d\ude00\ud83d\ude80\ud83c\udf0d","mixed":"caf\u00e9 na\u00efve r\u00e9sum\u00e9"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v map[string]string
		if err := shapejson.Unmarshal(data, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

func BenchmarkEncodingJSON_Unmarshal_UTF8(b *testing.B) {
	data := []byte(`{"greeting":"Hello, \u4e16\u754c!","emoji":"\ud83d\ude00\ud83d\ude80\ud83c\udf0d","mixed":"caf\u00e9 na\u00efve r\u00e9sum\u00e9"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v map[string]string
		if err := json.Unmarshal(data, &v); err != nil {
			b.Fatal(err)
		}
		_ = v
	}
}

// ================================
// Unmarshal Struct Comparison (v1 baseline)
// ================================

func BenchmarkEncodingJSON_Unmarshal_Struct(b *testing.B) {
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
		if err := json.Unmarshal(jsonData, &p); err != nil {
			b.Fatal(err)
		}
		_ = p
	}
}
