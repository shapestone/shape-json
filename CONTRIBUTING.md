# Contributing to Shape

Thank you for your interest in contributing to Shape! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Format Addition Policy](#format-addition-policy)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to conduct@shapestone.dev.

## Format Addition Policy

**IMPORTANT:** Shape is intentionally scoped to large, standardized, open data formats only.

### Formats That Belong in Shape ✅

Shape should include only widely-used, standardized, open formats:

- ✅ **JSON** - RFC 8259 standard
- ✅ **XML** - W3C standard
- ✅ **YAML** - YAML.org specification
- ✅ **CSV** - RFC 4180 standard
- ✅ **Properties** - Java properties format (.properties files)
- ✅ **Text** - Common text-based key:value formats

**Criteria for format inclusion:**

1. **Widely adopted:** Used across multiple industries and applications
2. **Open standard:** Has a published specification or RFC
3. **Language agnostic:** Not specific to a single programming language or framework
4. **Stable:** Mature format with established usage patterns
5. **General purpose:** Suitable for broad use cases, not domain-specific

### Formats That Don't Belong in Shape ❌

The following should **NOT** be added to Shape:

- ❌ **Custom/Proprietary DSLs:** Domain-specific languages for individual projects
- ❌ **Application-Specific Formats:** Formats unique to a single application
- ❌ **Internal Formats:** Company/project-internal configuration languages
- ❌ **Experimental Formats:** Unstable or unproven formats
- ❌ **Niche Formats:** Used by small communities only

**Examples of formats that should NOT be added:**

- Custom diagram DSLs (e.g., Inkling)
- Proprietary config formats (e.g., company-specific configs)
- Application-specific markup (e.g., custom template languages)
- Experimental data formats (e.g., research projects)

### For Custom DSLs: Use Shape's Tokenizer

If you're building a custom DSL, **use Shape's tokenization framework as a library** in your own project instead of adding your format to Shape:

```go
import "github.com/shapestone/shape/pkg/tokenizer"

// Build your custom DSL using Shape's tokenizer
tok := tokenizer.NewTokenizer(
    yourCustomMatchers...,
)
```

**Resources for Custom DSLs:**

- [Custom DSL Guide](docs/CUSTOM_DSL_GUIDE.md)
- [Tokenizer Documentation](pkg/tokenizer/README.md)
- [Example Project: Inkling](https://github.com/shapestone/inkling)

## How to Contribute

### Types of Contributions Welcome

1. **Bug Fixes:** Fix parsing errors, incorrect behavior
2. **Performance Improvements:** Optimize tokenization, parsing
3. **Documentation:** Improve guides, examples, API docs
4. **Test Coverage:** Add tests for edge cases
5. **Tooling:** Improve CI/CD, development tools
6. **Examples:** Add usage examples

### Types of Contributions We Generally Don't Accept

1. **New Format Parsers:** Unless they meet the strict criteria above
2. **Breaking API Changes:** For v1.x releases (semver)
3. **Scope Creep:** Features outside Shape's core mission (parsing/tokenization)

## Development Setup

See [Local Setup Guide](docs/contributor/local-setup.md) for detailed instructions.

### Quick Setup

```bash
# Clone repository
git clone https://github.com/shapestone/shape.git
cd shape

# Run tests
go test ./...

# Run linter
golangci-lint run

# Check coverage
go test -cover ./...
```

## Pull Request Process

1. **Fork the repository** and create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write clean, documented code
   - Add tests for new functionality
   - Update documentation as needed

3. **Run tests and linting:**
   ```bash
   go test ./...
   golangci-lint run
   ```

4. **Commit with clear messages:**
   ```bash
   git commit -m "feat: add support for X"
   git commit -m "fix: resolve issue with Y"
   git commit -m "docs: update tokenizer guide"
   ```

   Use [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `test:` - Test additions/changes
   - `refactor:` - Code refactoring
   - `perf:` - Performance improvements
   - `chore:` - Build process, tooling

5. **Push and create PR:**
   ```bash
   git push origin feature/your-feature-name
   ```

   Then create a pull request on GitHub with:
   - Clear title and description
   - Reference any related issues
   - Explain what changed and why

6. **Code Review:**
   - Maintainers will review your PR
   - Address feedback and make requested changes
   - Once approved, maintainers will merge

## Testing Guidelines

### Test Coverage Requirements

- **New Code:** Must have tests
- **Bug Fixes:** Add test that reproduces the bug
- **Target Coverage:** Maintain 90%+ coverage

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./pkg/tokenizer/

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With verbose output
go test -v ./...

# Benchmarks
go test -bench=. ./...
```

### Writing Good Tests

```go
func TestFeatureName(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"

    // Act
    result, err := SomeFunction(input)

    // Assert
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    if result != expected {
        t.Errorf("Expected %q, got %q", expected, result)
    }
}
```

## Branching Workflow

See [Branching Workflow](docs/contributor/BRANCHING_WORKFLOW.md) for details on:

- Branch naming conventions
- Merge strategies
- Release process

## Questions?

- **Issues:** [GitHub Issues](https://github.com/shapestone/shape/issues)
- **Discussions:** [GitHub Discussions](https://github.com/shapestone/shape/discussions)
- **Documentation:** [docs/](docs/)

Thank you for contributing to Shape! 🎉
