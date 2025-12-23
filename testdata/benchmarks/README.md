# Benchmark Test Data

This directory contains JSON test files used for performance benchmarking.

## Files

### small.json (~60 bytes)
Simple JSON object with basic data types:
- Single-level object
- Mix of string, number, boolean types
- Minimal nesting

**Use case**: Testing parser overhead on small payloads

### medium.json (~4.5KB)
Realistic API response structure:
- 3-4 levels of nesting
- Arrays with 5-20 items
- Mix of objects, arrays, and primitives
- Representative of typical REST API responses

**Use case**: Realistic web application payloads

### large.json (~420KB)
Large dataset for stress testing:
- 1000+ records in arrays
- 50x50 numeric matrix (2500 values)
- Deep nesting (7 levels)
- 100 large string values
- Complex nested structures

**Use case**: Testing memory efficiency and performance at scale

## Benchmark Results (Apple M1 Max)

### Shape-JSON Parse Performance
```
BenchmarkShapeJSON_Parse_Small       139558   8490 ns/op    11.43 MB/s   17088 B/op    269 allocs/op
BenchmarkShapeJSON_Parse_Medium        3082 401698 ns/op    11.43 MB/s  779486 B/op  11894 allocs/op
BenchmarkShapeJSON_Parse_Large           36  33.9ms/op      12.39 MB/s    57.8 MB/op 879168 allocs/op
```

### Shape-JSON ParseReader Performance
```
BenchmarkShapeJSON_ParseReader_Small    6150 205110 ns/op    0.51 MB/s   1.65 MB/op    438 allocs/op
BenchmarkShapeJSON_ParseReader_Medium    133   8.9ms/op     0.51 MB/s   63.7 MB/op  19531 allocs/op
BenchmarkShapeJSON_ParseReader_Large       2 858.5ms/op     0.49 MB/s    4.6 GB/op 1431620 allocs/op
```

### encoding/json Unmarshal Performance
```
BenchmarkEncodingJSON_Parse_Small   1635847   736.2 ns/op  109.58 MB/s    640 B/op     16 allocs/op
BenchmarkEncodingJSON_Parse_Medium    28818  41888 ns/op  109.58 MB/s  25120 B/op    617 allocs/op
BenchmarkEncodingJSON_Parse_Large       354   3.4ms/op    124.48 MB/s   1.79 MB/op  42754 allocs/op
```

### encoding/json Decoder Performance
```
BenchmarkEncodingJSON_Decoder_Small   1409823   867.8 ns/op   99.62 MB/s   1312 B/op     18 allocs/op
BenchmarkEncodingJSON_Decoder_Medium    26455  46077 ns/op   99.62 MB/s  34912 B/op    625 allocs/op
BenchmarkEncodingJSON_Decoder_Large       351   3.3ms/op    127.07 MB/s   2.41 MB/op  42769 allocs/op
```

## Key Observations

1. **encoding/json is significantly faster** (~10-11x faster) than shape-json for string-based parsing
   - encoding/json: 124 MB/s for large files
   - shape-json: 12 MB/s for large files

2. **Memory usage differs significantly**:
   - encoding/json uses less memory per operation (1.79 MB for large)
   - shape-json uses more memory (57.8 MB for large) due to AST construction
   - shape-json creates more allocations due to AST node creation

3. **ParseReader has high memory overhead**:
   - ParseReader uses significantly more memory (4.6 GB for large file)
   - This is likely due to buffering in the streaming tokenizer
   - Performance is also much slower (0.49 MB/s vs 12.39 MB/s)

4. **Trade-offs**:
   - encoding/json is optimized for fast parsing and minimal allocations
   - shape-json builds a complete AST for query/manipulation capabilities
   - The AST provides value for JSONPath queries and tree manipulation
   - Consider using encoding/json for simple deserialization needs

## Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./pkg/json/

# Run specific size category
go test -bench='Parse_Medium' -benchmem ./pkg/json/

# Run with longer benchtime for more accurate results
go test -bench=. -benchmem -benchtime=10s ./pkg/json/

# Compare shape-json vs encoding/json
go test -bench='Parse_Large' -benchmem ./pkg/json/

# Run only shape-json benchmarks
go test -bench='ShapeJSON' -benchmem ./pkg/json/

# Run only encoding/json benchmarks
go test -bench='EncodingJSON' -benchmem ./pkg/json/
```

## Regenerating Test Data

To regenerate large.json with different characteristics, use the generator script:

```bash
go run scripts/generate_large_json.go > testdata/benchmarks/large.json
```

## Metrics Explanation

- **ns/op**: Nanoseconds per operation (lower is better)
- **MB/s**: Throughput in megabytes per second (higher is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

## Future Improvements

Areas for potential optimization:
1. Reduce allocations in AST node creation
2. Investigate ParseReader memory usage
3. Consider object pooling for frequently allocated types
4. Benchmark with -cpuprofile and -memprofile for deeper analysis
