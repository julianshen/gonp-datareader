# Coverage Exclude Examples Module Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore the repository coverage command above 90% by keeping runnable examples out of the root library coverage profile.

**Architecture:** Treat `examples/` as a separate Go module that depends on the root module via a local `replace`. Root `go test ./...` will cover library packages only, while examples remain buildable from their own module.

**Tech Stack:** Go modules, existing example command packages, root coverage command from `AGENTS.md`.

---

## Chunk 1: Split Examples Coverage Boundary

### Task 1: Prove Current Coverage Failure

**Files:**
- Read: `coverage.out`

- [x] **Step 1: Run root coverage total**

Run: `go tool cover -func=coverage.out | tail -n 1`
Observed: `total: (statements) 50.1%`, with `examples/*` packages at `0.0%`.

### Task 2: Add Examples Module Boundary

**Files:**
- Create: `examples/go.mod`

- [x] **Step 1: Add nested module**

Create `examples/go.mod` requiring `github.com/julianshen/gonp-datareader` and replacing it with `..`.

- [x] **Step 2: Verify root tests and coverage**

Run:
`go test ./...`
`go test -covermode=atomic -coverprofile=coverage.out ./...`
`go tool cover -func=coverage.out | tail -n 1`

Expected: root coverage total exceeds 90%.

- [x] **Step 3: Verify examples still build from their module**

Run: `go test ./...` from `examples/`

Expected: all example command packages compile.
