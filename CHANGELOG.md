# Changelog

All notable changes to the shape-json project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.10.0] - 2026-02-10

### Added
- **Compiled encoder cache for `Marshal()`** — type-level encoder caching with pre-computed field layouts
- **Zero-reflect fast path** for common types (`appendInterface` type-switch)
- `encoder.go` — compiled struct/map/slice encoders with pre-encoded key bytes
- `encoder_helpers.go` — `appendInterface`, `appendISO8601Duration`, `sortStrings`
- `escape.go` — zero-allocation `appendEscapedString`
- Marshal benchmarks with `b.ReportAllocs()` and stdlib comparisons

### Changed
- `Marshal()` rewritten to use two-tier fast path (5.4x faster structs, 3.1x faster maps vs previous version)

## [0.9.1] - 2025-12-24

### Added
- **Thread Safety Documentation**: Comprehensive documentation of thread safety guarantees
  - New "Thread Safety" section in README.md with concurrent usage examples
  - Detailed thread safety chapter in USER_GUIDE.md with API table and best practices
  - Package-level godoc documentation on concurrent use in pkg/json/parser.go
  - All public APIs documented as thread-safe (Unmarshal, Marshal, Parse, Validate, etc.)
  - Warning about sharing Decoder/Encoder instances (matches encoding/json behavior)

### Changed
- **Go Version**: Standardized to Go 1.23 across Shape ecosystem
  - Updated go.mod from 1.25 to 1.23 for GitHub Actions compatibility
  - Updated CI workflow to use actions/setup-go@v5 with Go 1.23
  - Aligned with [Shape Go Version Policy](https://github.com/shapestone/shape-core/blob/main/docs/policies/GO_VERSION.md)

### Fixed
- CI lint issues: Applied gofmt and added nolint annotation for parseFalse function

## [0.9.0] - 2025-12-22

### Initial Release

This is the initial public release of shape-json, a high-performance JSON parser
for the Shape Parser™ ecosystem.

### Highlights

**🚀 Performance**
- **2x faster** than encoding/json for standard unmarshaling operations
- **1.6x less memory** usage than encoding/json
- Drop-in replacement with identical API

**🎯 Dual-Path Architecture**
- Fast Path: Optimized direct unmarshaling without AST overhead
- AST Path: Full tree construction for advanced features
- Automatic path selection based on API usage

**✨ Features**
- Full encoding/json API compatibility (`Unmarshal`, `Marshal`, `Decoder`, `Encoder`)
- JSONPath query support (RFC 9535 compliant)
- Streaming parser with constant memory usage
- DOM-style API for programmatic JSON manipulation
- Unified AST representation compatible with Shape ecosystem

### Added

**Core APIs**
- `json.Unmarshal()` - Fast unmarshaling (2x faster than encoding/json)
- `json.Marshal()` - Convert Go types to JSON
- `json.Validate()` - Fast JSON syntax validation
- `json.Parse()` - Parse to AST for advanced features
- `json.NewDecoder()` / `json.NewEncoder()` - Streaming APIs

**JSONPath Support**
- `jsonpath.ParseString()` - Compile JSONPath expressions
- RFC 9535-compliant implementation
- Support for filters, recursive descent, slices
- Comprehensive operator support (comparison, logical, regex)
- Zero external dependencies

**Advanced Features**
- `ParseReader()` - Constant-memory streaming parser
- `ParseDocument()` - DOM-style API
- AST node pooling with `ReleaseTree()` for memory optimization
- Position tracking for detailed error messages

**Documentation**
- Comprehensive README with examples
- USER_GUIDE.md for detailed usage patterns
- ARCHITECTURE.md explaining dual-path design
- PERFORMANCE_REPORT.md with detailed benchmarks
- Example code for all major features

**Testing & Quality**
- 84%+ test coverage
- 100+ comprehensive tests
- Fuzz testing for parser robustness
- Benchmark suite with historical tracking
- CI/CD with GitHub Actions

### Performance Benchmarks

**vs encoding/json (Unmarshal)**
- Small JSON (60B): 1.9x faster, 1.6x less memory
- Medium JSON (4.5KB): 1.9x faster, 1.3x less memory
- Large JSON (410KB): 2.0x faster, 1.4x less memory

**Fast Path vs AST Path**
- 9.7x faster for unmarshaling
- 10.2x less memory usage
- Automatic selection - no API changes needed

### Documentation

- `README.md` - Quick start and overview
- `USER_GUIDE.md` - Comprehensive usage guide
- `ARCHITECTURE.md` - Design and implementation details
- `PERFORMANCE_REPORT.md` - Detailed benchmark analysis
- `CONTRIBUTING.md` - Contribution guidelines
- `SECURITY.md` - Security policy and vulnerability reporting

### License

Apache License 2.0
Copyright © 2020-2025 Shapestone

### Links

- Repository: https://github.com/shapestone/shape-json
- Documentation: https://pkg.go.dev/github.com/shapestone/shape-json
- Issues: https://github.com/shapestone/shape-json/issues

[0.10.0]: https://github.com/shapestone/shape-json/releases/tag/v0.10.0
[0.9.0]: https://github.com/shapestone/shape-json/releases/tag/v0.9.0
