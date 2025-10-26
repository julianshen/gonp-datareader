# gonp-datareader Project Overview

A comprehensive package for fetching financial and economic data in Go, designed to work with [gonp](https://github.com/julianshen/gonp).

---

## 📦 Generated Documentation Files

### 1. **spec.md** (17.5 KB)
Complete project specification and requirements document.

**Key Sections:**
- Project Overview & Goals
- Supported Data Sources (Yahoo, FRED, World Bank, Alpha Vantage, etc.)
- Package Structure & Architecture
- Core API Design with Examples
- Data Source Specifications
- Error Handling Strategy
- Performance Requirements
- Testing Requirements
- Security Considerations
- Release Plan (v0.1.0 → v1.0.0)

**Use for:** Understanding project scope, API design decisions, and requirements

---

### 2. **CLAUDE.md** (26 KB)
Development methodology combining TDD, Tidy First, and Effective Go.

**Key Sections:**

**Part 1:** TDD Methodology (Red → Green → Refactor)
**Part 2:** Tidy First Approach (Structural vs Behavioral Changes)
**Part 3:** Go Idioms (Naming, Errors, Interfaces, Concurrency, Documentation)
**Part 4:** Code Quality Standards
**Part 5:** Commit Discipline
**Part 6:** Testing Standards
**Part 7:** Performance & Optimization
**Part 8:** Example Workflows
**Part 9:** Tools & Automation
**Part 10:** Summary Checklist

**Use for:** Daily development guidance, code reviews, and quality standards

---

### 3. **plan.md** (23 KB)
Step-by-step implementation plan with 200+ discrete tasks.

**Phases:**
- Phase 0: Project Setup
- Phase 1: Foundation (Error Handling & HTTP Client)
- Phase 2: Base Reader Interface
- Phase 3: Yahoo Finance Reader (MVP)
- Phase 4: DataReader Factory
- Phase 5: FRED Reader
- Phase 6: Rate Limiting
- Phase 7: Response Caching
- Phase 8: World Bank Reader
- Phase 9: Alpha Vantage Reader
- Phase 10: Documentation & Examples
- Phase 11: Testing & Quality
- Phase 12: Performance Optimization
- Phase 13: Additional Data Sources
- Phase 14: Release Preparation

**Each task follows:** ☐ Test → Implement → Verify → Commit → ☑

**Use for:** Daily development workflow, progress tracking, knowing what to build next

---

### 4. **README-DOCS.md** (11 KB)
Guide to using the documentation effectively.

**Key Sections:**
- Document Overview & Purpose
- Getting Started for New Contributors
- Development Workflow Examples
- Key Principles & Quick Reference
- Progress Tracking
- Finding Information (Reference Table)
- Tools Setup
- Learning Path (Week-by-week)
- Common Pitfalls & Solutions
- Checklist Before Each Commit

**Use for:** Onboarding, understanding how documents relate, finding information quickly

---

## 🎯 Project Goals

**Seamless Integration** with gonp DataFrames  
**Multiple Data Sources** for comprehensive coverage  
**Type Safety** leveraging Go's type system  
**High Performance** with efficient parsing and minimal allocations  
**Extensibility** through plugin architecture  
**Production Ready** with comprehensive error handling and testing

---

## 🏗️ Architecture Overview

```
gonp-datareader/
├── datareader.go          # Main interface & factory
├── config.go              # Configuration
├── error.go               # Custom error types
├── sources/               # Data source implementations
│   ├── source.go         # Base interface
│   ├── yahoo/            # Yahoo Finance
│   ├── fred/             # Federal Reserve Economic Data
│   ├── worldbank/        # World Bank
│   └── alphavantage/     # Alpha Vantage
├── internal/             
│   ├── http/            # HTTP client with retry
│   ├── cache/           # Response caching
│   └── utils/           # Common utilities
├── examples/            # Usage examples
├── docs/                # Documentation
└── testdata/            # Test fixtures
```

---

## 🚀 Quick Start Guide

### 1. Read the Documentation
```
1. spec.md      → Understand WHAT we're building
2. CLAUDE.md    → Learn HOW to build it  
3. plan.md      → Follow the step-by-step PLAN
```

### 2. Setup Environment
```bash
# Initialize project
go mod init github.com/yourorg/gonp-datareader

# Install tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create Makefile (see CLAUDE.md Section 9.2)
```

### 3. Start Development
```bash
# Open plan.md
# Find first unchecked item: ☐
# Write test (RED)
go test -v ./...  # Confirm failure

# Implement (GREEN)
go test -v ./...  # Confirm pass

# Refactor if needed

# Quality checks
make check

# Mark complete: ☑
# Commit with proper message
git commit -m "feat: description"
```

---

## 📊 Feature Roadmap

### v0.1.0 - MVP
- ✅ Core DataReader interface
- ✅ Yahoo Finance support
- ✅ FRED support
- ✅ Basic error handling
- ✅ Essential documentation

### v0.2.0 - Extended Sources
- ✅ World Bank support
- ✅ Alpha Vantage support
- ✅ Response caching
- ✅ Rate limiting

### v0.3.0 - Production Ready
- ✅ Additional sources (IEX, Stooq, OECD)
- ✅ Comprehensive examples
- ✅ Performance optimizations

### v1.0.0 - Stable Release
- ✅ All planned sources
- ✅ 90%+ test coverage
- ✅ Production battle-tested
- ✅ Complete documentation

---

## 💻 Example Usage

### Basic Usage
```go
package main

import (
    "context"
    "time"
    
    dr "github.com/yourorg/gonp-datareader"
)

func main() {
    ctx := context.Background()
    start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
    
    // Fetch Yahoo Finance data
    df, err := dr.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(df.Head())
}
```

### With Custom Options
```go
opts := &dr.Options{
    Timeout:     30 * time.Second,
    MaxRetries:  3,
    EnableCache: true,
    RateLimit:   10.0, // 10 requests per second
}

reader, err := dr.DataReader("yahoo", opts)
df, err := reader.Read(ctx, []string{"AAPL", "MSFT"}, start, end)
```

### Multiple Sources
```go
// Yahoo Finance
stockData, _ := dr.Read(ctx, "AAPL", "yahoo", start, end, nil)

// FRED Economic Data
fredOpts := &dr.Options{APIKey: "your-fred-api-key"}
gdp, _ := dr.Read(ctx, "GDP", "fred", start, end, fredOpts)

// World Bank
wbData, _ := dr.Read(ctx, "NY.GDP.PCAP.CD", "worldbank", start, end, nil)
```

---

## 🧪 Testing Philosophy

From **CLAUDE.md**:

### TDD Cycle
```
RED: Write failing test
  ↓
GREEN: Make it pass (minimal code)
  ↓
REFACTOR: Improve structure
  ↓
COMMIT: With proper message
```

### Coverage Goals
- Minimum: 80% overall
- Target: 90%+ for v1.0
- Critical paths: 100% (parsers, error handling)

### Test Types
- Unit Tests: Individual functions
- Integration Tests: Real API calls (with VCR)
- Benchmark Tests: Performance profiling
- Example Tests: Documentation validation

---

## 📐 Design Principles

From **Effective Go** and **CLAUDE.md**:

1. **Simplicity First:** Clear over clever
2. **Test-Driven:** Write tests before code
3. **Small Functions:** <50 lines, single responsibility
4. **Explicit Dependencies:** No hidden globals
5. **Interface Segregation:** Small, focused interfaces
6. **Error Handling:** Always check, wrap with context
7. **Concurrency Safety:** Share memory by communicating
8. **Documentation:** Godoc for all exports

---

## 🔧 Development Tools

### Required
```bash
go 1.21+                 # Go compiler
goimports               # Import management
golangci-lint           # Comprehensive linting
gosec                   # Security scanning
```

### Makefile Targets
```bash
make test               # Run all tests
make test-coverage      # Generate coverage report
make lint               # Run all linters
make fmt                # Format code
make check              # ALL quality checks (run before commit!)
make build              # Build project
```

---

## 📈 Progress Tracking

Track your progress in **plan.md**:

```markdown
☐ - Not started
☑ - Completed

Statistics:
- Total Items: ~200
- Completed: ___
- Remaining: ___
- Percentage: ___%
```

---

## 🤝 Contributing

### Before Starting
1. Read spec.md (requirements)
2. Read CLAUDE.md (methodology)
3. Read plan.md (find next task)

### Development Process
1. Pick next ☐ from plan.md
2. Write test (RED)
3. Implement minimum (GREEN)
4. Refactor if needed
5. Run `make check`
6. Mark ☑ in plan.md
7. Commit with proper message
8. Repeat

### Commit Format
```
<type>: <description>

Types:
- feat: New feature
- fix: Bug fix
- refactor: Code restructuring
- test: Add/update tests
- docs: Documentation
- chore: Build/CI/deps
- perf: Performance
```

---

## 📚 Learning Resources

### Project Docs
- **spec.md** - Requirements & Design
- **CLAUDE.md** - Development Guide
- **plan.md** - Implementation Roadmap
- **README-DOCS.md** - Documentation Guide

### External Resources
- [Effective Go](https://go.dev/doc/effective_go)
- [gonp Documentation](https://github.com/julianshen/gonp)
- [pandas-datareader Docs](https://pandas-datareader.readthedocs.io/)

---

## ✅ Quality Checklist

Before each commit (from CLAUDE.md):
- ☐ All tests passing
- ☐ All linters passing
- ☐ Code formatted
- ☐ Single logical change
- ☐ Proper commit message
- ☐ Item marked in plan.md

---

## 🎯 Success Criteria

From **spec.md**:

**Functionality:** All data sources working reliably  
**Performance:** Sub-second typical queries  
**Reliability:** 99%+ success rate  
**Usability:** Clear API, comprehensive docs  
**Quality:** >80% test coverage, clean code

---

## 📞 Getting Help

### When Stuck
1. Review relevant CLAUDE.md section
2. Check examples in spec.md
3. Verify TDD cycle adherence
4. Simplify the test
5. Break into smaller steps

### Quick Reference
- API design? → spec.md
- How to test? → CLAUDE.md Part 1
- What's next? → plan.md
- Commit format? → CLAUDE.md Part 5
- Go idioms? → CLAUDE.md Part 3

---

## 🎉 Getting Started

**Ready to begin?**

1. Open **plan.md**
2. Go to Phase 0: Project Setup
3. Check first ☐ item
4. Follow TDD cycle from CLAUDE.md
5. Build something awesome!

The documentation is comprehensive, the plan is clear, and the methodology is proven. You have everything you need to create a high-quality Go library.

**Let's build gonp-datareader!** 🚀
