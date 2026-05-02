# Yahoo Live No Skip Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the Yahoo live integration test so it never uses `t.Skip`, while keeping the real API call explicitly opt-in.

**Architecture:** Keep mock integration coverage unchanged. Add a small environment-gate helper in `integration_test.go` and make the live Yahoo test return successfully when the gate is disabled, but fail on real API errors when the gate is enabled.

**Tech Stack:** Go integration tests, standard library `os`, existing Yahoo/datareader test helpers.

---

## Chunk 1: Yahoo Live Opt-In Gate

### Task 1: Add Env-Gate Test

**Files:**
- Modify: `integration_test.go`

- [x] **Step 1: Write the failing test**

Add a test that proves the opt-in gate is disabled when `GONP_DATAREADER_LIVE_YAHOO` is unset and enabled only when it is set to `1`.

- [x] **Step 2: Run test to verify it fails**

Run: `go test -tags=integration -run TestRealYahooIntegrationEnabled ./...`
Expected: FAIL because `realYahooIntegrationEnabled` does not exist.

- [x] **Step 3: Write minimal implementation**

Add `realYahooIntegrationEnabled()` in `integration_test.go` using `os.Getenv("GONP_DATAREADER_LIVE_YAHOO") == "1"`.

- [x] **Step 4: Run test to verify it passes**

Run: `go test -tags=integration -run TestRealYahooIntegrationEnabled ./...`
Expected: PASS.

### Task 2: Refactor Live Yahoo Test

**Files:**
- Modify: `integration_test.go`

- [x] **Step 1: Remove `t.Skip` from the live Yahoo test path**

If the opt-in env var is not enabled, log the command to enable it and `return`.

- [x] **Step 2: Make opted-in real API failures fail**

Remove the skip-on-error behavior. When opted in, any API error should call `t.Fatalf`.

- [x] **Step 3: Verify targeted integration tests**

Run: `go test -tags=integration -run 'TestIntegration_RealYahooFinance|TestRealYahooIntegrationEnabled' ./...`
Expected: PASS without making a real API request unless `GONP_DATAREADER_LIVE_YAHOO=1`.

- [x] **Step 4: Verify full suite and coverage**

Run:
`go test ./...`
`go test -covermode=atomic -coverprofile=coverage.out ./...`
`go tool cover -func=coverage.out | tail -n 1`
