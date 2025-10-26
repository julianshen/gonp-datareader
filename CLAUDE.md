# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Project Overview

**gonp-datareader** is a Go library for fetching financial and economic data from various sources (Yahoo Finance, FRED, World Bank, etc.), designed to work with [gonp](https://github.com/julianshen/gonp) DataFrames. This is the Go equivalent of Python's pandas-datareader.

**Current Status:** Initial planning phase - no implementation code yet. All planning documentation is complete.

---

## Essential Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./sources/yahoo/

# Run specific test
go test -v -run TestYahooReader_Read ./sources/yahoo/

# Run tests with race detection
go test -race ./...
```

### Code Quality
```bash
# Format code
gofmt -s -w .

# Run go vet
go vet ./...

# Run golangci-lint (if available)
golangci-lint run

# Import organization
goimports -w .
```

### Build
```bash
# Build all packages
go build ./...

# Build specific package
go build ./sources/yahoo
```

### Module Management
```bash
# Initialize module (if needed)
go mod init github.com/julianshen/gonp-datareader

# Add dependencies
go get github.com/julianshen/gonp

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

---

## Architecture

This project follows a **plugin architecture** where each data source implements a common `Reader` interface:

```
gonp-datareader/
├── datareader.go          # Main Reader interface + factory function
├── config.go              # Options configuration
├── error.go               # Custom error types
├── sources/               # Data source implementations
│   ├── source.go         # Base source interface
│   ├── yahoo/            # Yahoo Finance (no API key)
│   ├── fred/             # Federal Reserve Economic Data (API key optional)
│   ├── worldbank/        # World Bank (no API key)
│   └── alphavantage/     # Alpha Vantage (API key required)
├── internal/
│   ├── http/            # HTTP client with retry logic + rate limiting
│   ├── cache/           # Optional response caching
│   └── utils/           # Validation, date parsing
├── examples/            # Usage examples
└── docs/                # Additional documentation
```

### Key Interfaces

**Reader Interface** (sources/source.go):
```go
type Reader interface {
    Read(ctx context.Context, symbols []string, start, end time.Time) (*dataframe.DataFrame, error)
    ReadSingle(ctx context.Context, symbol string, start, end time.Time) (*series.Series, error)
    ValidateSymbol(symbol string) error
    Name() string
}
```

**Factory Pattern** (datareader.go):
```go
func DataReader(source string, opts *Options) (Reader, error)
func Read(ctx context.Context, symbol string, source string, start, end time.Time, opts *Options) (*dataframe.DataFrame, error)
```

---

## Development Methodology

This project follows **strict Test-Driven Development (TDD)** with the **Tidy First** approach:

### TDD Cycle (ALWAYS follow this)
1. **RED**: Write a failing test that defines desired functionality
2. **GREEN**: Write minimum code to make the test pass
3. **REFACTOR**: Improve code structure while keeping tests green

**Never write production code without a failing test first.**

### Commit Discipline

**Separate structural and behavioral changes:**

- **Structural** (refactor): Renaming, extracting functions, moving code - NO behavior change
- **Behavioral** (feat/fix): Adding features, fixing bugs - changes behavior

**Commit prefixes:**
- `feat:` - New feature
- `fix:` - Bug fix
- `refactor:` - Code restructuring (no behavior change)
- `test:` - Adding or updating tests only
- `docs:` - Documentation changes
- `chore:` - Build, dependencies, tooling

**Commit when:**
- All tests pass (`go test ./...`)
- Code is formatted (`gofmt -s -w .`)
- Linters pass (`go vet ./...`)
- Change is atomic (one logical unit)

---

## Testing Standards

### Coverage Requirements
- Minimum: **80%** for all packages
- Critical paths: **100%** (data parsing, error handling)
- Test all exported functions
- Test error conditions and edge cases

### Test Organization

Use **table-driven tests** with subtests:
```go
func TestYahooReader_Read(t *testing.T) {
    tests := []struct {
        name    string
        symbol  string
        wantErr bool
    }{
        {name: "valid symbol", symbol: "AAPL", wantErr: false},
        {name: "invalid symbol", symbol: "", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Data
- Use `testdata/` directories for fixtures
- Real data samples (anonymized if needed)
- Document data sources

### Mocking
- Mock external HTTP requests using `httptest`
- Inject dependencies via interfaces
- Mock time-dependent code

---

## Go Conventions

### Naming
- **Variables**: Short names in small scopes (`i`, `err`, `ctx`), descriptive in larger scopes
- **Functions**: MixedCaps, no underscores
- **Getters**: No "Get" prefix - `Symbol()` not `GetSymbol()`
- **Interfaces**: Single-method use verb+"er" - `Reader`, `Writer`, `Closer`
- **Packages**: Lowercase, no underscores, singular preferred

### Error Handling
Always check errors immediately and wrap with context:
```go
data, err := fetchData()
if err != nil {
    return nil, fmt.Errorf("fetch data: %w", err)
}
```

Use custom error types for package API:
```go
var ErrInvalidSymbol = errors.New("invalid symbol")

// Check with errors.Is
if errors.Is(err, ErrInvalidSymbol) {
    // Handle
}
```

### Interface Design
- Accept interfaces, return concrete types
- Keep interfaces small and focused
- Define interfaces at point of use (consumer), not implementation (provider)

### Documentation
Every exported type, function, constant needs a doc comment starting with its name:
```go
// Reader fetches data from remote sources.
// Implementations must be safe for concurrent use.
type Reader interface {
    // Read fetches data for the given symbols within the date range.
    Read(ctx context.Context, symbols []string, start, end time.Time) (*dataframe.DataFrame, error)
}
```

---

## Project Planning Documents

This repository has comprehensive planning documentation:

- **spec.md** - Complete project specification, API design, requirements, data source details
- **plan.md** - Step-by-step implementation plan with ~200 discrete tasks following TDD
- **PROJECT-OVERVIEW.md** - Quick reference guide to all documentation
- **README-DOCS.md** - Guide for using the planning documents effectively

### Working from the Plan

The **plan.md** file contains the implementation roadmap:
1. Find next unchecked item: `☐`
2. Follow TDD cycle (Red → Green → Refactor)
3. Mark complete: `☑`
4. Commit with proper prefix

Phases are ordered by priority:
- Phase 0: Project Setup
- Phase 1: Foundation (Error Handling & HTTP Client)
- Phase 2: Base Reader Interface
- Phase 3: Yahoo Finance Reader (MVP)
- Phase 4-14: Additional sources, features, optimization

---

## Code Quality Checklist

Before each commit:
- ☐ All tests passing (`go test ./...`)
- ☐ All linters passing (`go vet ./...`)
- ☐ Code formatted (`gofmt -s -w .`)
- ☐ Single logical change
- ☐ Proper commit message with prefix
- ☐ Test coverage maintained (>80%)
- ☐ Documentation updated if needed
- ☐ Item marked complete in plan.md

---

## Common Tasks

### Adding a New Data Source

1. Create package: `sources/newsource/`
2. Write test for Reader interface implementation
3. Implement Reader interface
4. Write parser tests with real data samples
5. Implement parser
6. Add to factory in `datareader.go`
7. Add example in `examples/`
8. Update documentation

### Adding a Feature

1. Write failing test
2. Implement minimum code to pass
3. Add more tests for edge cases
4. Refactor if needed (keep tests green)
5. Commit structural changes separately from behavioral
6. Update documentation

### Fixing a Bug

1. Write test that demonstrates the bug
2. Verify test fails
3. Fix the code
4. Verify test passes
5. Consider refactoring
6. Commit with `fix:` prefix

---

## Performance Considerations

- **HTTP**: Use connection pooling, HTTP/2, keep-alive
- **Memory**: Stream large responses, minimize allocations, reuse buffers
- **Concurrency**: Safe for concurrent use, parallel fetching for multiple symbols
- **Caching**: Optional file-based cache for responses

### Benchmarking
```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkParseResponse -benchmem ./sources/yahoo/
```

---

## Dependencies

**Required:**
- Go 1.21+
- github.com/julianshen/gonp (latest)

**Testing:**
- Standard library `testing`
- github.com/stretchr/testify (for assertions, optional)

**Philosophy:** Prefer standard library, minimize external dependencies, keep dependency tree shallow.

---

## Quick Reference

| Need | Look Here |
|------|-----------|
| What to build next | plan.md |
| API design decisions | spec.md |
| TDD workflow | This file + plan.md |
| Data source details | spec.md (Data Source Specifications) |
| Error handling patterns | This file (Error Handling section) |
| Test patterns | This file (Testing Standards) |
| Commit message format | This file (Commit Discipline) |

---

## Development Workflow Example

```bash
# 1. Check next task
# Open plan.md, find next ☐ item

# 2. Write failing test (RED)
# Add test in appropriate _test.go file
go test -v ./path/to/package/
# Confirm test fails

# 3. Implement minimum code (GREEN)
# Write just enough code to pass
go test -v ./path/to/package/
# Confirm test passes

# 4. Refactor (if needed)
# Improve structure, run tests after each change
go test -v ./path/to/package/

# 5. Quality checks
gofmt -s -w .
go vet ./...
go test ./...

# 6. Commit
git add .
git commit -m "feat: add symbol validation for yahoo reader"

# 7. Mark complete in plan.md
# Change ☐ to ☑

# 8. Repeat
```

---

## Key Principles

1. **Write test first** - No production code without a failing test
2. **Minimum implementation** - Just enough to pass the test
3. **Refactor in green** - Only refactor when tests pass
4. **Separate commits** - Structural vs behavioral changes
5. **Clear over clever** - Obvious beats elegant
6. **Small functions** - Target <50 lines, single responsibility
7. **Explicit dependencies** - Use interfaces, inject dependencies
8. **Document exports** - Every exported item needs godoc
9. **Check errors** - Always check, always wrap with context
10. **Test thoroughly** - >80% coverage, test edge cases and errors
