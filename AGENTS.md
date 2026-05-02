# Agent Instructions for gonp-datareader

This file provides operational rules for AI agents working in this repository.

---

## Branching Rules

**Never commit to main.**

Always create a branch before editing:

```bash
git checkout -b feat/<slug>    # new feature
git checkout -b fix/<slug>     # bug fix
git checkout -b chore/<slug>   # maintenance, deps, tooling
```

Examples: `feat/yahoo-options`, `fix/race-in-cache`, `chore/update-linter`

---

## Planning Rules

**Plan first, then implement.**

For any non-trivial work, invoke `superpowers:writing-plans` before writing code.

For any new feature or API surface (new endpoint, new conversion mode, new CLI flag, etc.), invoke `superpowers:brainstorming` before entering plan mode.

---

## PR Rules

Each PR should be a **small, independently reviewable chunk**.

Do not bundle unrelated changes. Split large changes into stacked PRs.

---

## TDD Rules

**Strict TDD: Red → Green → Refactor.**

- Write the failing test first.
- Write the minimum code to make it pass.
- Refactor with tests green.

**Each commit lands the test for the behaviour it introduces.**

Never an "implement X, tests later" commit.
Never a plan structured as "implement A, B, C" then "add tests for A, B, C".

---

## Verification Rules

**Verify before claiming done.**

Use `superpowers:verification-before-completion`: run the exact build/test commands and quote the output before saying something passes.

### Coverage Command

```bash
go test -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -n 1
```

Report the total percentage in your completion message.

---

## Coverage Rules

**Test coverage stays above 90%.**

- Measure with `go test -covermode=atomic -coverprofile=coverage.out ./...`
- Report `go tool cover -func=coverage.out | tail -n 1`
- cgo trampolines that cannot be unit-tested live in the upstream lok package — this repo's code should be pure Go and fully covered.

---

## Quality Rules

**No shortcuts.**

No disabled checks, relaxed lint rules, lowered thresholds, `// nolint`, skipped build tags, `|| true`, or `--no-verify`.

If something is hard, understand it.

---

## Test Integrity Rules

**Don't skip failing tests.**

No `t.Skip`, no commented assertions, no `_test_disabled` files, no never-set build tags.

Fix the code or fix the test.

---

## Completion Rules

**Complete implementation.**

Trace root causes; do not work around crashes by removing the call that triggers them, do not silently shrink scope, do not leave `// TODO` stubs in shipped code.

---

## Feature Change Rules

**Ask before cutting any feature.**

If a planned endpoint, option, format, or behaviour cannot be implemented as agreed, pause and ask — even removing a case from an integration smoke test counts.

---

## Derived from CLAUDE.md

For project overview, architecture, TDD workflow, Go conventions, and data source details, see `CLAUDE.md`.
