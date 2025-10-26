# gonp-datareader Project Documentation

This project provides comprehensive documentation for building a golang data reader library similar to pandas-datareader, designed to work with the [gonp](https://github.com/julianshen/gonp) library.

## 📚 Document Overview

### 1. **spec.md** - Project Specification
**Purpose:** Comprehensive project requirements and architecture

**Contents:**
- Project goals and vision
- Supported data sources (Yahoo Finance, FRED, World Bank, etc.)
- Complete package structure
- Core API design with code examples
- Data source specifications
- Error handling strategy
- Performance requirements
- Testing requirements
- Security considerations
- Release plan (v0.1.0 → v1.0.0)

**When to reference:**
- Understanding project scope and requirements
- Designing new features
- Making architectural decisions
- Planning releases
- Onboarding new contributors

---

### 2. **CLAUDE.md** - Development Guidelines
**Purpose:** Development methodology combining TDD, Tidy First, and Go best practices

**Contents:**

**Part 1: TDD Methodology**
- Red → Green → Refactor cycle
- Writing tests first
- Minimum implementation principles
- Refactoring rules

**Part 2: Tidy First Approach**
- Separating structural from behavioral changes
- When and how to refactor
- Commit discipline

**Part 3: Go Idioms**
- Code organization
- Naming conventions (from Effective Go)
- Error handling patterns
- Interface design
- Concurrency patterns
- Documentation standards

**Part 4-10:** Code quality, performance, examples, and checklists

**When to reference:**
- Before starting any development work
- When writing tests
- During code review
- When refactoring
- When committing code
- When stuck or unsure how to proceed

---

### 3. **plan.md** - Implementation Plan
**Purpose:** Step-by-step TDD implementation checklist

**Contents:**
- 200+ discrete implementation tasks
- Organized into 14 phases
- Each task follows: Test → Implement → Verify → Commit
- Checkboxes (☐/☑) for progress tracking
- Commit message guidelines for each task
- Refactoring checkpoints
- Progress tracking section

**Phases:**
0. Project Setup
1. Foundation (Error Handling & HTTP Client)
2. Base Reader Interface
3. Yahoo Finance Reader (MVP)
4. DataReader Factory
5. FRED Reader
6. Rate Limiting
7. Response Caching
8. World Bank Reader
9. Alpha Vantage Reader
10. Documentation & Examples
11. Testing & Quality
12. Performance Optimization
13. Additional Data Sources
14. Release Preparation

**When to reference:**
- Daily development work
- Deciding what to build next
- Tracking project progress
- Ensuring nothing is missed
- Following TDD methodology

---

## 🚀 Getting Started

### For New Contributors

1. **Read in this order:**
   1. `spec.md` - Understand WHAT we're building
   2. `CLAUDE.md` - Learn HOW to build it
   3. `plan.md` - See the step-by-step PLAN

2. **Set up your environment:**
   ```bash
   # Clone and setup
   mkdir gonp-datareader
   cd gonp-datareader
   go mod init github.com/yourorg/gonp-datareader
   
   # Install tools (from CLAUDE.md Section 9.1)
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Create Makefile (from CLAUDE.md Section 9.2)
   # Copy Makefile content from CLAUDE.md
   ```

3. **Start developing:**
   ```bash
   # Open plan.md
   # Find the first unchecked item (☐)
   # Follow the TDD cycle from CLAUDE.md
   # Mark complete (☑) when done
   ```

---

## 💡 Development Workflow

### Daily Workflow

```bash
# 1. Check plan.md for next task
# Example: "☐ Test: YahooReader struct exists"

# 2. Write the test (RED)
cat > sources/yahoo/yahoo_test.go << 'EOF'
package yahoo_test

import (
    "testing"
)

func TestYahooReader_StructExists(t *testing.T) {
    reader := yahoo.YahooReader{}
    if reader == nil {
        t.Fatal("YahooReader should not be nil")
    }
}
EOF

# 3. Run test - confirm it fails
go test -v ./sources/yahoo/
# FAIL: undefined: yahoo.YahooReader

# 4. Implement minimum code (GREEN)
cat > sources/yahoo/yahoo.go << 'EOF'
package yahoo

type YahooReader struct {}
EOF

# 5. Run test - confirm it passes
go test -v ./sources/yahoo/
# PASS: TestYahooReader_StructExists

# 6. Refactor if needed (keep tests GREEN)

# 7. Quality checks
make fmt    # Format code
make lint   # Run linters
make test   # All tests

# 8. Mark complete in plan.md
# Change: ☐ Test: YahooReader struct exists
# To:     ☑ Test: YahooReader struct exists

# 9. Commit
git add .
git commit -m "feat: create YahooReader structure"

# 10. Move to next item in plan.md
```

---

## 📖 Key Principles

### From CLAUDE.md

**Always:**
- ✅ Write tests BEFORE implementation
- ✅ Write minimum code to pass tests
- ✅ Keep functions small and focused
- ✅ Separate structural changes (refactor) from behavioral changes (feat/fix)
- ✅ Run all tests before committing
- ✅ Follow Go naming conventions

**Never:**
- ❌ Write code without a failing test first
- ❌ Mix refactoring with feature changes
- ❌ Commit with failing tests
- ❌ Skip refactoring phase
- ❌ Use clever solutions over clear ones

### Test-Driven Development Cycle

```
RED (Write failing test)
  ↓
GREEN (Make it pass - minimal code)
  ↓
REFACTOR (Improve structure)
  ↓
COMMIT (With appropriate message)
  ↓
REPEAT (Next test)
```

---

## 🎯 Quick Reference

### When Writing Tests
```go
// Good test name (from CLAUDE.md Section 1.2)
func TestYahooReader_ReadSingle_ValidSymbol_ReturnsDataFrame(t *testing.T)

// Table-driven pattern (from CLAUDE.md Section 1.2)
tests := []struct {
    name    string
    symbol  string
    wantErr bool
}{
    {
        name:    "valid symbol returns data",
        symbol:  "AAPL",
        wantErr: false,
    },
    {
        name:    "empty symbol returns error",
        symbol:  "",
        wantErr: true,
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### When Committing
```bash
# Format (from CLAUDE.md Section 5.2)
<type>: <description>

# Types:
feat:     # New feature
fix:      # Bug fix
refactor: # Code restructuring without behavior change
test:     # Adding or updating tests
docs:     # Documentation changes
chore:    # Build, CI, dependencies
perf:     # Performance improvement

# Examples from plan.md:
git commit -m "feat: implement Yahoo Finance URL builder"
git commit -m "refactor: extract HTTP client configuration"
git commit -m "test: add edge case tests for date parsing"
git commit -m "docs: add API reference documentation"
```

---

## 📊 Tracking Progress

### In plan.md

```markdown
## Progress Tracking

**Current Phase:** Phase 3: Yahoo Finance Reader
**Last Completed:** 3.2 Yahoo URL Building
**Next Up:** 3.3 Yahoo HTTP Request

**Statistics:**
- Total Items: ~200
- Completed: 25
- Remaining: 175
- Percentage: 12.5%
```

### Checking Coverage
```bash
# From CLAUDE.md Section 6.1
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Target: >80% coverage
```

---

## 🔍 Finding Information

### "How do I...?"

| Question | Document | Section |
|----------|----------|---------|
| What data sources are supported? | spec.md | Supported Data Sources |
| How do I write a test? | CLAUDE.md | Part 1: TDD Methodology |
| What's the package structure? | spec.md | Package Structure |
| How do I name functions? | CLAUDE.md | Part 3.2: Naming Conventions |
| How do I handle errors? | CLAUDE.md | Part 3.3: Error Handling |
| What do I build next? | plan.md | Find next ☐ item |
| How do I commit changes? | CLAUDE.md | Part 5: Commit Discipline |
| What's the API design? | spec.md | Core API Design |
| How do I optimize performance? | CLAUDE.md | Part 7: Performance |
| What examples are needed? | plan.md | Phase 10 |

---

## 🛠️ Tools Setup

### Required Tools (from CLAUDE.md)
```bash
# Go 1.21+
go version

# Formatting
go install golang.org/x/tools/cmd/goimports@latest

# Linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Security
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### Makefile Targets (from CLAUDE.md)
```bash
make test          # Run all tests
make test-coverage # Generate coverage report
make lint          # Run linters
make fmt           # Format code
make check         # Run all quality checks (ALWAYS before commit)
make build         # Build project
make clean         # Clean generated files
```

---

## 🎓 Learning Path

### Week 1: Foundation
1. Read all three documents completely
2. Set up development environment
3. Complete Phase 0: Project Setup from plan.md
4. Complete Phase 1: Foundation from plan.md
5. Practice TDD cycle with simple examples

### Week 2: Core Implementation
1. Complete Phase 2: Base Reader Interface
2. Complete Phase 3: Yahoo Finance Reader (MVP)
3. Complete Phase 4: DataReader Factory
4. Review and refactor

### Week 3-4: Additional Sources
1. Complete Phase 5: FRED Reader
2. Complete Phase 6-7: Rate Limiting & Caching
3. Complete Phase 8-9: World Bank & Alpha Vantage
4. Add comprehensive tests

### Week 5: Polish
1. Complete Phase 10: Documentation & Examples
2. Complete Phase 11: Testing & Quality
3. Complete Phase 12: Performance Optimization
4. Prepare for v0.1.0 release

---

## 🚨 Common Pitfalls & Solutions

### Pitfall 1: Writing code before tests
**Solution:** Always refer to plan.md - every item starts with "Test:"

### Pitfall 2: Making tests too complex
**Solution:** Review CLAUDE.md Part 1.3 - start with simplest possible test

### Pitfall 3: Mixing refactoring with features
**Solution:** CLAUDE.md Part 2.2 - separate commits always

### Pitfall 4: Unclear commit messages
**Solution:** Use plan.md suggested commit messages and CLAUDE.md format

### Pitfall 5: Not running tests frequently
**Solution:** Run `make test` after every implementation step

---

## 📞 Support

### When Stuck
1. Re-read relevant CLAUDE.md section
2. Review examples in spec.md
3. Check if you're following TDD cycle
4. Simplify the test
5. Break into smaller steps from plan.md

### Resources
- **Effective Go:** https://go.dev/doc/effective_go
- **gonp Documentation:** https://github.com/julianshen/gonp
- **pandas-datareader Reference:** https://pandas-datareader.readthedocs.io/

---

## 📋 Checklist Before Each Commit

From CLAUDE.md Part 5.1:
- ☐ All tests passing (`go test ./...`)
- ☐ All linters passing (`make lint`)
- ☐ Code formatted (`make fmt`)
- ☐ Change is single logical unit
- ☐ Commit message follows format
- ☐ Item marked complete in plan.md

---

## 🎉 Success Metrics

### From spec.md
- ✅ All priority data sources working reliably
- ✅ Sub-second response for typical queries
- ✅ 99%+ success rate for valid requests
- ✅ Clear API, comprehensive documentation
- ✅ High test coverage (>80%)
- ✅ Clean code, minimal bugs

---

## 📝 Summary

These three documents work together:

1. **spec.md** = WHAT to build (requirements, architecture, design)
2. **CLAUDE.md** = HOW to build it (methodology, standards, practices)
3. **plan.md** = Step-by-step ROADMAP (ordered tasks with checkboxes)

Follow the plan, use the methodology, meet the spec. The result will be a high-quality, well-tested, idiomatic Go library.

**Start here:** Open plan.md, find Phase 0, check the first ☐ item, and begin your TDD journey!
