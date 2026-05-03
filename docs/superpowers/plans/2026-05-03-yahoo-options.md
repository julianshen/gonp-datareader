# Yahoo Options Chain Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Yahoo Finance options chain data reader with full test coverage following TDD.

**Architecture:** Extend the existing Yahoo source package with an `OptionsReader` that fetches options chain data from Yahoo Finance's `/v7/finance/options/{symbol}` endpoint. The reader follows the established pattern of embedding `BaseSource` and using `internalhttp.RetryableClient` for HTTP requests.

**Tech Stack:** Go 1.21+, standard library (`encoding/json`, `net/http`, `time`), existing `internal/http` retry client, `testify` for assertions.

---

## File Structure

| File | Responsibility |
|------|--------------|
| `sources/yahoo/options.go` | Options data types (`OptionContract`, `OptionsChain`, `OptionsReader`) and reader implementation |
| `sources/yahoo/options_test.go` | Unit tests for options parsing, reader behavior, error handling |
| `sources/yahoo/testdata/options_aapl.json` | Real (anonymized) Yahoo Finance options response fixture |

---

## Chunk 1: Data Types and Parser

### Task 1: Define Options Data Types

**Files:**
- Create: `sources/yahoo/options.go`
- Test: `sources/yahoo/options_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOptionContract_StructTags(t *testing.T) {
    contract := OptionContract{
        ContractSymbol: "AAPL250516C00150000",
        Strike:         150.0,
        Type:           "CALL",
        LastPrice:      5.25,
        Bid:            5.20,
        Ask:            5.30,
        Change:         0.50,
        PercentChange:  10.5,
        Volume:         1500,
        OpenInterest:   2500,
        ImpliedVol:     0.35,
        InTheMoney:     true,
    }
    assert.Equal(t, "AAPL250516C00150000", contract.ContractSymbol)
    assert.Equal(t, 150.0, contract.Strike)
    assert.Equal(t, "CALL", contract.Type)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./sources/yahoo/ -run TestOptionContract_StructTags`
Expected: FAIL with "undefined: OptionContract"

- [ ] **Step 3: Write minimal implementation**

```go
// sources/yahoo/options.go
package yahoo

import "time"

// OptionContract represents a single options contract.
type OptionContract struct {
    ContractSymbol string  `json:"contractSymbol"`
    Strike         float64 `json:"strike"`
    Expiration     time.Time `json:"expiration"`
    Type           string  `json:"type"`
    LastPrice      float64 `json:"lastPrice"`
    Bid            float64 `json:"bid"`
    Ask            float64 `json:"ask"`
    Change         float64 `json:"change"`
    PercentChange  float64 `json:"percentChange"`
    Volume         int64   `json:"volume"`
    OpenInterest   int64   `json:"openInterest"`
    ImpliedVol     float64 `json:"impliedVolatility"`
    InTheMoney     bool    `json:"inTheMoney"`
}

// OptionsChain represents all contracts for a single expiration date.
type OptionsChain struct {
    ExpirationDate time.Time        `json:"expirationDate"`
    Calls          []OptionContract `json:"calls"`
    Puts           []OptionContract `json:"puts"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./sources/yahoo/ -run TestOptionContract_StructTags`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add sources/yahoo/options.go sources/yahoo/options_test.go
git commit -m "feat: define options chain data types"
```

---

### Task 2: Implement Yahoo Options JSON Parser

**Files:**
- Modify: `sources/yahoo/options.go`
- Test: `sources/yahoo/options_test.go`
- Create: `sources/yahoo/testdata/options_aapl.json`

- [ ] **Step 1: Write the failing test**

```go
func TestParseOptionsJSON_ValidResponse(t *testing.T) {
    data, err := os.ReadFile("testdata/options_aapl.json")
    require.NoError(t, err)

    chain, err := parseOptionsJSON(bytes.NewReader(data))
    require.NoError(t, err)
    assert.NotNil(t, chain)
    assert.NotEmpty(t, chain.Calls)
    assert.NotEmpty(t, chain.Puts)

    call := chain.Calls[0]
    assert.NotEmpty(t, call.ContractSymbol)
    assert.Greater(t, call.Strike, 0.0)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./sources/yahoo/ -run TestParseOptionsJSON_ValidResponse`
Expected: FAIL with "undefined: parseOptionsJSON"

- [ ] **Step 3: Create test fixture**

Create `sources/yahoo/testdata/options_aapl.json` with a realistic Yahoo Finance `/v7/finance/options/AAPL` response structure:
- `optionChain.result[0].expirationDates` (array of unix timestamps)
- `optionChain.result[0].strikes` (array of strike prices)
- `optionChain.result[0].options[0].calls` (array of call contracts)
- `optionChain.result[0].options[0].puts` (array of put contracts)

Each contract object includes: `contractSymbol`, `strike`, `expiration`, `lastPrice`, `bid`, `ask`, `change`, `percentChange`, `volume`, `openInterest`, `impliedVolatility`, `inTheMoney`.

- [ ] **Step 4: Write minimal implementation**

```go
// parseOptionsJSON parses Yahoo Finance options JSON response.
func parseOptionsJSON(r io.Reader) (*OptionsChain, error) {
    var resp struct {
        OptionChain struct {
            Result []struct {
                ExpirationDates []int64   `json:"expirationDates"`
                Strikes         []float64 `json:"strikes"`
                Options         []struct {
                    ExpirationDate int64            `json:"expirationDate"`
                    Calls          []OptionContract `json:"calls"`
                    Puts           []OptionContract `json:"puts"`
                } `json:"options"`
            } `json:"result"`
            Error interface{} `json:"error"`
        } `json:"optionChain"`
    }

    if err := json.NewDecoder(r).Decode(&resp); err != nil {
        return nil, fmt.Errorf("decode options JSON: %w", err)
    }

    if resp.OptionChain.Error != nil {
        return nil, fmt.Errorf("yahoo finance error: %v", resp.OptionChain.Error)
    }

    if len(resp.OptionChain.Result) == 0 || len(resp.OptionChain.Result[0].Options) == 0 {
        return nil, fmt.Errorf("no options data found")
    }

    result := resp.OptionChain.Result[0]
    opt := result.Options[0]

    // Enrich contract types
    for i := range opt.Calls {
        opt.Calls[i].Type = "CALL"
        opt.Calls[i].Expiration = time.Unix(opt.ExpirationDate, 0)
    }
    for i := range opt.Puts {
        opt.Puts[i].Type = "PUT"
        opt.Puts[i].Expiration = time.Unix(opt.ExpirationDate, 0)
    }

    return &OptionsChain{
        ExpirationDate: time.Unix(opt.ExpirationDate, 0),
        Calls:          opt.Calls,
        Puts:           opt.Puts,
    }, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test -v ./sources/yahoo/ -run TestParseOptionsJSON_ValidResponse`
Expected: PASS

- [ ] **Step 6: Add edge case tests**

```go
func TestParseOptionsJSON_EmptyResult(t *testing.T) {
    jsonData := `{"optionChain":{"result":[],"error":null}}`
    _, err := parseOptionsJSON(strings.NewReader(jsonData))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no options data")
}

func TestParseOptionsJSON_YahooError(t *testing.T) {
    jsonData := `{"optionChain":{"result":[],"error":{"code":"Not Found","description":"No data found"}}}`
    _, err := parseOptionsJSON(strings.NewReader(jsonData))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "yahoo finance error")
}
```

Run: `go test -v ./sources/yahoo/ -run TestParseOptionsJSON`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add sources/yahoo/options.go sources/yahoo/options_test.go sources/yahoo/testdata/
git commit -m "feat: implement Yahoo options JSON parser"
```

---

## Chunk 2: Options Reader

### Task 3: Implement OptionsReader Structure

**Files:**
- Modify: `sources/yahoo/options.go`
- Test: `sources/yahoo/options_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOptionsReader_Struct(t *testing.T) {
    reader := NewOptionsReader(nil)
    assert.NotNil(t, reader)
    assert.NotNil(t, reader.client)
}

func TestOptionsReader_Name(t *testing.T) {
    reader := NewOptionsReader(nil)
    assert.Equal(t, "Yahoo Finance Options", reader.Name())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader`
Expected: FAIL with "undefined: NewOptionsReader"

- [ ] **Step 3: Write minimal implementation**

```go
const optionsAPIURL = "https://query1.finance.yahoo.com/v7/finance/options/%s"

// OptionsReader fetches options chain data from Yahoo Finance.
type OptionsReader struct {
    *sources.BaseSource
    client  *internalhttp.RetryableClient
    baseURL string
}

// NewOptionsReader creates a new options chain reader.
func NewOptionsReader(opts *internalhttp.ClientOptions) *OptionsReader {
    return NewOptionsReaderWithBaseURL(opts, optionsAPIURL)
}

// NewOptionsReaderWithBaseURL creates an options reader with a custom base URL.
func NewOptionsReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *OptionsReader {
    if opts == nil {
        opts = internalhttp.DefaultClientOptions()
    }
    return &OptionsReader{
        BaseSource: sources.NewBaseSource("yahoo-options"),
        client:     internalhttp.NewRetryableClient(opts),
        baseURL:    baseURL,
    }
}

// Name returns the display name.
func (o *OptionsReader) Name() string {
    return "Yahoo Finance Options"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add sources/yahoo/options.go sources/yahoo/options_test.go
git commit -m "feat: create OptionsReader structure"
```

---

### Task 4: Implement GetOptionsChain

**Files:**
- Modify: `sources/yahoo/options.go`
- Test: `sources/yahoo/options_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOptionsReader_GetOptionsChain_MockServer(t *testing.T) {
    // Setup mock server
    mux := http.NewServeMux()
    mux.HandleFunc("/v7/finance/options/", func(w http.ResponseWriter, r *http.Request) {
        symbol := strings.TrimPrefix(r.URL.Path, "/v7/finance/options/")
        assert.Equal(t, "AAPL", symbol)

        data, _ := os.ReadFile("testdata/options_aapl.json")
        w.Header().Set("Content-Type", "application/json")
        w.Write(data)
    })
    server := httptest.NewServer(mux)
    defer server.Close()

    reader := NewOptionsReaderWithBaseURL(nil, server.URL+"/v7/finance/options/%s")
    ctx := context.Background()

    chain, err := reader.GetOptionsChain(ctx, "AAPL", nil)
    require.NoError(t, err)
    assert.NotNil(t, chain)
    assert.NotEmpty(t, chain.Calls)
    assert.NotEmpty(t, chain.Puts)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader_GetOptionsChain_MockServer`
Expected: FAIL with "undefined: OptionsReader.GetOptionsChain"

- [ ] **Step 3: Write minimal implementation**

```go
// GetOptionsChain fetches the options chain for a symbol.
// If expiration is nil, returns the nearest expiration.
func (o *OptionsReader) GetOptionsChain(ctx context.Context, symbol string, expiration *time.Time) (*OptionsChain, error) {
    if err := o.ValidateSymbol(symbol); err != nil {
        return nil, fmt.Errorf("invalid symbol: %w", err)
    }

    url := fmt.Sprintf(o.baseURL, symbol)
    if expiration != nil {
        url = fmt.Sprintf("%s?date=%d", url, expiration.Unix())
    }

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    resp, err := o.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch options: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("yahoo returned %d: %s", resp.StatusCode, string(body))
    }

    return parseOptionsJSON(resp.Body)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader_GetOptionsChain_MockServer`
Expected: PASS

- [ ] **Step 5: Add error case tests**

```go
func TestOptionsReader_GetOptionsChain_InvalidSymbol(t *testing.T) {
    reader := NewOptionsReader(nil)
    _, err := reader.GetOptionsChain(context.Background(), "", nil)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid symbol")
}

func TestOptionsReader_GetOptionsChain_HTTPError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Not Found"))
    }))
    defer server.Close()

    reader := NewOptionsReaderWithBaseURL(nil, server.URL+"/%s")
    _, err := reader.GetOptionsChain(context.Background(), "INVALID", nil)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "404")
}
```

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader_GetOptionsChain`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add sources/yahoo/options.go sources/yahoo/options_test.go
git commit -m "feat: implement GetOptionsChain with mock tests"
```

---

### Task 5: Implement GetExpirationDates

**Files:**
- Modify: `sources/yahoo/options.go`
- Test: `sources/yahoo/options_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOptionsReader_GetExpirationDates(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/v7/finance/options/", func(w http.ResponseWriter, r *http.Request) {
        resp := `{"optionChain":{"result":[{"expirationDates":[1750118400,1750723200],"strikes":[150,155],"options":[]}],"error":null}}`
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(resp))
    })
    server := httptest.NewServer(mux)
    defer server.Close()

    reader := NewOptionsReaderWithBaseURL(nil, server.URL+"/v7/finance/options/%s")
    dates, err := reader.GetExpirationDates(context.Background(), "AAPL")
    require.NoError(t, err)
    assert.Len(t, dates, 2)
    assert.Equal(t, time.Unix(1750118400, 0), dates[0])
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader_GetExpirationDates`
Expected: FAIL with "undefined: OptionsReader.GetExpirationDates"

- [ ] **Step 3: Write minimal implementation**

```go
// GetExpirationDates returns available expiration dates for a symbol.
func (o *OptionsReader) GetExpirationDates(ctx context.Context, symbol string) ([]time.Time, error) {
    if err := o.ValidateSymbol(symbol); err != nil {
        return nil, fmt.Errorf("invalid symbol: %w", err)
    }

    url := fmt.Sprintf(o.baseURL, symbol)
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    resp, err := o.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch expirations: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("yahoo returned %d: %s", resp.StatusCode, string(body))
    }

    var data struct {
        OptionChain struct {
            Result []struct {
                ExpirationDates []int64 `json:"expirationDates"`
            } `json:"result"`
            Error interface{} `json:"error"`
        } `json:"optionChain"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, fmt.Errorf("decode expirations: %w", err)
    }

    if data.OptionChain.Error != nil {
        return nil, fmt.Errorf("yahoo finance error: %v", data.OptionChain.Error)
    }

    if len(data.OptionChain.Result) == 0 {
        return nil, fmt.Errorf("no expiration dates found")
    }

    dates := make([]time.Time, len(data.OptionChain.Result[0].ExpirationDates))
    for i, ts := range data.OptionChain.Result[0].ExpirationDates {
        dates[i] = time.Unix(ts, 0)
    }

    return dates, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./sources/yahoo/ -run TestOptionsReader_GetExpirationDates`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add sources/yahoo/options.go sources/yahoo/options_test.go
git commit -m "feat: implement GetExpirationDates"
```

---

## Chunk 3: Verification and Completion

### Task 6: Final Verification

- [ ] **Step 1: Run full test suite**

```bash
go test -v ./sources/yahoo/...
go test -coverprofile=coverage.out ./sources/yahoo/
go tool cover -func=coverage.out | grep "yahoo/options"
```

Expected: All tests PASS, options package coverage > 90%

- [ ] **Step 2: Run linting and formatting**

```bash
gofmt -s -w ./sources/yahoo/options.go ./sources/yahoo/options_test.go
go vet ./sources/yahoo/...
```

Expected: No issues

- [ ] **Step 3: Update plan.md progress tracking**

Add Phase 17 entry for Options Chain and mark tasks complete.

- [ ] **Step 4: Final commit**

```bash
git add plan.md
git commit -m "docs: update plan with Phase 17 Options Chain"
```

---

## Progress Tracking

**Current Phase:** Phase 17 - Yahoo Options Chain
**Last Completed:** (to be updated during execution)
**Next Up:** Phase 18 - Yahoo Holdings Reader

**Statistics:**
- New files: 2 (`options.go`, `options_test.go`)
- Test fixtures: 1 (`testdata/options_aapl.json`)
- Estimated test functions: 8-10
- Target coverage: > 90% for options package
