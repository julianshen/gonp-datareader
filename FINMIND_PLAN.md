# FinMind Data Source Implementation Plan

## Overview

**FinMind** is an open-source financial data platform providing over 50 datasets covering Taiwan and international markets. This plan outlines the implementation of FinMind as the 11th data source for gonp-datareader.

**Project:** gonp-datareader Phase 16 - FinMind Reader
**Status:** Planning Phase
**Priority:** High (complements existing TWSE reader with more comprehensive Taiwan market data)

---

## FinMind API Summary

### Key Information

- **API Endpoint:** `https://api.finmindtrade.com/api/v4/data`
- **Authentication:** Bearer token (optional but recommended)
- **Rate Limits:**
  - Without token: 300 requests/hour
  - With token: 600 requests/hour
- **Response Format:** JSON
- **Date Format:** YYYY-MM-DD (ISO 8601)

### Data Coverage

**Taiwan Market:**
- Stock prices (daily, minute, tick data)
- Technical indicators (P/E, P/B ratios)
- Fundamental data (financial statements, dividends, monthly revenue)
- Institutional data (foreign investors, margin trading)
- Derivatives (futures, options)
- Real-time quotes

**International Markets:**
- US stocks (daily and minute-level since 2021)
- US Treasury yields
- Commodities (gold, oil)
- Exchange rates
- G7/G8 central bank rates
- UK, Europe, Japan markets

### API Request Format

```http
GET https://api.finmindtrade.com/api/v4/data?dataset=TaiwanStockPrice&data_id=2330&start_date=2020-04-02&end_date=2020-04-12
Authorization: Bearer YOUR_TOKEN_HERE
```

### API Response Format

```json
{
  "data": [
    {
      "date": "2020-04-06",
      "stock_id": "2330",
      "Trading_Volume": 59712754,
      "Trading_money": 16324198154,
      "open": 273,
      "max": 275.5,
      "min": 270,
      "close": 275.5,
      "spread": 4,
      "Trading_turnover": 19971
    }
  ]
}
```

---

## Implementation Phases

### Phase 16.1: FinMind Reader Structure ✅

**Goal:** Create the basic FinMind reader structure with proper initialization.

**Tasks:**
- ✅ Test: FinMindReader struct exists
- ✅ Implement: FinMindReader in `sources/finmind/finmind.go`
- ✅ Test: FinMindReader embeds BaseSource
- ✅ Implement: Embed BaseSource with source name "finmind"
- ✅ Test: NewFinMindReader returns non-nil reader
- ✅ Implement: NewFinMindReader constructor with token support
- ✅ Test: NewFinMindReaderWithToken handles API token
- ✅ Implement: Token storage and handling
- ✅ Test: FinMindReader implements Reader interface
- ✅ Verify: All Reader methods present (Read, ReadSingle, ValidateSymbol, Name, Source)

**Files created:**
- `sources/finmind/finmind.go` (162 lines)
- `sources/finmind/finmind_test.go` (62 lines)

**Coverage:** 69.2%

**Commit:** `feat: implement Phase 16.1 - FinMind reader structure with token support` (12edd5d)

---

### Phase 16.2: API Client Configuration ⏳

**Goal:** Set up HTTP client with authentication and rate limiting.

**Tasks:**
- ☐ Test: Client sets Bearer token in Authorization header
- ☐ Implement: Token-based authentication
- ☐ Test: Client respects rate limits (600 req/hour with token)
- ☐ Implement: Rate limiting configuration
- ☐ Test: Client builds correct API URLs
- ☐ Implement: URL building with query parameters
- ☐ Test: Client handles missing token gracefully
- ☐ Implement: Optional token support

**API Details:**
- Base URL: `https://api.finmindtrade.com/api/v4/data`
- Authentication: Bearer token in header
- Rate limit: 600 requests/hour (with token), 300 without

**Commit:** `feat: configure FinMind API client with authentication`

---

### Phase 16.3: Dataset Parameter Handling ⏳

**Goal:** Implement flexible dataset selection and parameter building.

**Tasks:**
- ☐ Test: BuildParams creates correct query parameters
- ☐ Implement: BuildParams function for API requests
- ☐ Test: TaiwanStockPrice dataset parameters
- ☐ Implement: Taiwan stock price dataset support
- ☐ Test: Date format conversion (time.Time → YYYY-MM-DD)
- ☐ Implement: Date formatting utilities
- ☐ Test: data_id parameter for symbol
- ☐ Implement: Symbol to data_id mapping

**Dataset Structure:**
```go
type DatasetParams struct {
    Dataset   string // "TaiwanStockPrice", "USStockPrice", etc.
    DataID    string // Stock symbol/ID
    StartDate string // YYYY-MM-DD
    EndDate   string // YYYY-MM-DD
}
```

**Commit:** `feat: implement dataset parameter handling`

---

### Phase 16.4: JSON Response Parser ⏳

**Goal:** Parse FinMind JSON responses into structured data.

**Tasks:**
- ☐ Test: ParseFinMindResponse parses valid JSON
- ☐ Implement: JSON unmarshaling for FinMind format
- ☐ Test: Extract stock_id, date, OHLCV fields
- ☐ Implement: Field extraction from response data
- ☐ Test: Handle Trading_Volume, Trading_money fields
- ☐ Implement: Volume and money parsing
- ☐ Test: Convert numeric fields correctly
- ☐ Implement: Type conversion for all numeric fields
- ☐ Test: Handle empty data array
- ☐ Implement: Empty response handling

**Response Structures:**
```go
type FinMindResponse struct {
    Data []FinMindStockData `json:"data"`
}

type FinMindStockData struct {
    Date            string  `json:"date"`
    StockID         string  `json:"stock_id"`
    TradingVolume   int64   `json:"Trading_Volume"`
    TradingMoney    int64   `json:"Trading_money"`
    Open            float64 `json:"open"`
    Max             float64 `json:"max"`
    Min             float64 `json:"min"`
    Close           float64 `json:"close"`
    Spread          float64 `json:"spread"`
    TradingTurnover int64   `json:"Trading_turnover"`
}

type ParsedData struct {
    Symbol       string
    Date         []time.Time
    Open         []float64
    High         []float64
    Low          []float64
    Close        []float64
    Volume       []int64
    Amount       []int64  // Trading money
    Transactions []int64  // Trading turnover
}
```

**Commit:** `feat: implement FinMind JSON parser`

---

### Phase 16.5: Symbol Validation ⏳

**Goal:** Validate Taiwan stock symbols for FinMind API.

**Tasks:**
- ☐ Test: Valid Taiwan stock codes (4-digit)
- ☐ Implement: Taiwan symbol validation
- ☐ Test: Valid US stock symbols (letters)
- ☐ Implement: US symbol validation
- ☐ Test: Invalid symbols return error
- ☐ Implement: Symbol validation with error messages
- ☐ Test: Empty symbol returns error
- ☐ Verify: Comprehensive symbol validation

**Symbol Formats:**
- Taiwan: 4-digit codes (e.g., "2330", "0050")
- US: Letters and sometimes numbers (e.g., "AAPL", "MSFT")

**Commit:** `feat: add FinMind symbol validation`

---

### Phase 16.6: ReadSingle Implementation ⏳

**Goal:** Implement single symbol data fetching.

**Tasks:**
- ☐ Test: ReadSingle fetches Taiwan stock (2330)
- ☐ Implement: ReadSingle with Taiwan dataset
- ☐ Test: ReadSingle validates symbol
- ☐ Implement: Symbol validation in ReadSingle
- ☐ Test: ReadSingle validates date range
- ☐ Implement: Date range validation
- ☐ Test: ReadSingle returns ParsedData
- ☐ Implement: Complete ReadSingle integration
- ☐ Test: ReadSingle handles HTTP errors
- ☐ Implement: Error handling for network issues

**Implementation Flow:**
1. Validate symbol
2. Validate date range
3. Build API URL with parameters
4. Add Authorization header if token exists
5. Execute HTTP request
6. Parse JSON response
7. Convert to ParsedData
8. Return data

**Commit:** `feat: implement FinMind ReadSingle method`

---

### Phase 16.7: Read (Multiple Symbols) Implementation ⏳

**Goal:** Implement parallel fetching for multiple symbols.

**Tasks:**
- ☐ Test: Read fetches multiple Taiwan stocks
- ☐ Implement: Read with parallel fetching
- ☐ Test: Read respects rate limits
- ☐ Implement: Rate limiting in parallel requests
- ☐ Test: Read handles partial failures
- ☐ Implement: Error handling for failed symbols
- ☐ Test: Read returns map[string]*ParsedData
- ☐ Verify: Complete Read implementation

**Rate Limiting Strategy:**
- Max 10 concurrent requests (same as TWSE)
- Respect 600 req/hour limit (with token)
- Use semaphore pattern for concurrency control

**Commit:** `feat: implement FinMind Read method with rate limiting`

---

### Phase 16.8: Error Handling ⏳

**Goal:** Comprehensive error handling for all failure scenarios.

**Tasks:**
- ☐ Test: Invalid token returns descriptive error
- ☐ Implement: Authentication error handling
- ☐ Test: Rate limit exceeded returns error
- ☐ Implement: Rate limit error detection
- ☐ Test: HTTP 401 (Unauthorized) handling
- ☐ Implement: Auth error responses
- ☐ Test: HTTP 429 (Too Many Requests) handling
- ☐ Implement: Rate limit error responses
- ☐ Test: Invalid dataset name returns error
- ☐ Implement: Dataset validation
- ☐ Test: Symbol not found handling
- ☐ Implement: Empty data response handling
- ☐ Test: Network timeout errors
- ☐ Implement: Timeout error handling

**Error Types:**
- Authentication errors (401, invalid token)
- Rate limit errors (429, quota exceeded)
- Data not found (empty response)
- Network errors (timeout, connection failed)
- Parse errors (invalid JSON)

**Commit:** `feat: add comprehensive FinMind error handling`

---

### Phase 16.9: Factory Registration ⏳

**Goal:** Integrate FinMind into the datareader factory system.

**Tasks:**
- ☐ Test: DataReader("finmind") returns FinMind reader
- ☐ Implement: Factory registration in datareader.go
- ☐ Test: Factory passes API token from Options
- ☐ Implement: Token passing from Options.APIKey
- ☐ Test: Read("2330", "finmind") works end-to-end
- ☐ Implement: Complete factory integration
- ☐ Test: ListSources() includes "finmind"
- ☐ Verify: Factory registration complete

**Integration Points:**
- Add import: `"github.com/julianshen/gonp-datareader/sources/finmind"`
- Add case in DataReader switch
- Add "finmind" to ListSources()
- Pass token from opts.APIKey

**Commit:** `feat: register FinMind reader with factory`

---

### Phase 16.10: Documentation and Examples ⏳

**Goal:** Provide comprehensive documentation and usage examples.

**Tasks:**
- ☐ Add package-level godoc for finmind package
- ☐ Document FinMindReader struct and methods
- ☐ Document authentication requirements
- ☐ Create examples/finmind/main.go
- ☐ Add example with token authentication
- ☐ Add example for Taiwan stocks
- ☐ Add example for US stocks (if supported)
- ☐ Update README.md with FinMind entry
- ☐ Document rate limits and best practices

**Example Structure:**
- Example 1: Basic usage with token
- Example 2: Taiwan stock data (TSMC 2330)
- Example 3: Multiple symbols
- Example 4: Error handling
- Example 5: Rate limit management

**Documentation Topics:**
- How to get FinMind API token
- Rate limit information (300 vs 600 req/hour)
- Symbol formats (Taiwan vs US)
- Available datasets
- Best practices

**Commit:** `docs: add comprehensive FinMind documentation and examples`

---

### Phase 16.11: Testing and Verification ⏳

**Goal:** Ensure comprehensive test coverage and quality.

**Tasks:**
- ☐ Test: All unit tests passing
- ☐ Verify: >80% test coverage
- ☐ Test: Integration tests with mock API
- ☐ Test: Authentication with valid/invalid tokens
- ☐ Test: Rate limiting behavior
- ☐ Test: Multiple dataset types
- ☐ Run: golangci-lint
- ☐ Run: go vet
- ☐ Verify: All examples compile and run
- ☐ Verify: CI passes on all platforms

**Coverage Target:** >80% (minimum), aim for 90%+

**Commit:** `test: verify comprehensive FinMind test suite`

---

## Technical Specifications

### Constants

```go
const (
    finmindBaseURL      = "https://api.finmindtrade.com/api/v4/data"
    finmindRateLimit    = 600 // requests per hour with token
    finmindRateLimitStr = "600 requests/hour"
)

// Dataset names
const (
    DatasetTaiwanStockPrice = "TaiwanStockPrice"
    DatasetUSStockPrice     = "USStockPrice"
    // Add more as needed
)
```

### Reader Structure

```go
type FinMindReader struct {
    *sources.BaseSource
    client  *internalhttp.RetryableClient
    baseURL string
    token   string // API token for authentication
}

func NewFinMindReader(opts *internalhttp.ClientOptions) *FinMindReader {
    return NewFinMindReaderWithToken(opts, "")
}

func NewFinMindReaderWithToken(opts *internalhttp.ClientOptions, token string) *FinMindReader {
    if opts == nil {
        opts = internalhttp.DefaultClientOptions()
    }

    // Set rate limit based on token availability
    if token != "" && opts.RateLimit == 0 {
        opts.RateLimit = float64(finmindRateLimit) / 3600 // Convert to req/sec
    }

    return &FinMindReader{
        BaseSource: sources.NewBaseSource("finmind"),
        client:     internalhttp.NewRetryableClient(opts),
        baseURL:    finmindBaseURL,
        token:      token,
    }
}
```

---

## Key Differences from TWSE

### Similarities
- Both provide Taiwan stock market data
- Both use JSON responses
- Both return OHLCV data
- Both support 4-digit Taiwan stock codes

### Differences

| Feature | TWSE | FinMind |
|---------|------|---------|
| **Authentication** | Not required | Optional but recommended (Bearer token) |
| **Rate Limits** | None specified | 300/hour (no token), 600/hour (with token) |
| **Data Coverage** | Taiwan only | Taiwan + International markets |
| **Datasets** | Single endpoint | Multiple datasets (50+) |
| **Historical Data** | Recent data | Since 1994 (Taiwan) |
| **API Style** | Simple public API | RESTful with dataset parameter |
| **Response Format** | Array of all stocks | Filtered by symbol |
| **Additional Data** | Basic OHLCV | OHLCV + fundamentals + institutional |

---

## Integration Strategy

### Token Management

```go
// In datareader.go factory
case "finmind":
    if apiKey != "" {
        return finmind.NewFinMindReaderWithToken(clientOpts, apiKey), nil
    }
    return finmind.NewFinMindReader(clientOpts), nil
```

### Usage Examples

**With Token:**
```go
opts := &datareader.Options{
    APIKey: "your-finmind-token",
    RateLimit: 600.0 / 3600.0, // 600 req/hour
}

reader, err := datareader.DataReader("finmind", opts)
data, err := reader.ReadSingle(ctx, "2330", start, end)
```

**Without Token (lower rate limit):**
```go
reader, err := datareader.DataReader("finmind", nil)
data, err := reader.ReadSingle(ctx, "2330", start, end)
```

---

## Testing Strategy

### Unit Tests
- Reader structure and initialization
- Parameter building and validation
- JSON parsing with various responses
- Symbol validation (Taiwan and US formats)
- Error handling for all scenarios

### Integration Tests
- Mock HTTP server for API responses
- Token authentication flow
- Rate limiting behavior
- Multiple dataset types
- Error responses (401, 429, etc.)

### End-to-End Tests
- Factory registration
- Convenience function usage
- Cross-platform compatibility (Linux, macOS, Windows)
- Go version compatibility (1.21, 1.22, 1.23)

---

## Success Criteria

- ✅ All tests passing with >80% coverage
- ✅ FinMind reader fully integrated into factory
- ✅ Comprehensive documentation and examples
- ✅ CI passing on all platforms and Go versions
- ✅ Rate limiting properly implemented
- ✅ Token authentication working
- ✅ Both Taiwan and US market support (if feasible)
- ✅ Error handling for all failure scenarios
- ✅ Production-ready code quality

---

## Timeline Estimate

**Total Phases:** 11
**Estimated Commits:** 15-20
**Estimated Coverage:** 90%+ (based on TWSE experience)
**Complexity:** Medium-High (authentication + rate limiting + multiple datasets)

**Comparison to TWSE:**
- Similar complexity for basic functionality
- Additional complexity for token authentication
- Additional complexity for rate limiting
- Additional complexity for multiple dataset support
- More comprehensive data coverage

---

## Dependencies

### Required
- Go 1.21+
- github.com/julianshen/gonp-datareader (existing)
- Standard library: encoding/json, net/http, time

### Optional
- FinMind API token (for higher rate limits and full access)

---

## Risk Assessment

### Low Risk
- JSON parsing (similar to TWSE)
- Symbol validation (4-digit codes)
- Date handling (standard ISO 8601)

### Medium Risk
- Rate limiting implementation
- Token authentication
- Multiple dataset support
- Handling API quota exhaustion

### Mitigation Strategies
- Implement robust rate limiting with token bucket
- Clear error messages for authentication failures
- Start with TaiwanStockPrice dataset only
- Add retry logic with exponential backoff
- Document token acquisition process clearly

---

## Future Enhancements

### Phase 17+ (Optional)
- Support for more FinMind datasets:
  - TaiwanStockInfo (company information)
  - TaiwanStockDividend (dividend data)
  - TaiwanStockPER (P/E ratios)
  - USStockPrice (US market data)
  - Futures and options data
- Real-time data support
- Fundamental data integration
- Institutional investor data
- Technical indicator calculations

---

## Notes

- FinMind is complementary to TWSE, not a replacement
- TWSE: Simple, no-auth, recent data
- FinMind: Comprehensive, auth-optional, historical data since 1994
- Users can choose based on their needs
- Both readers can coexist in the library
- FinMind offers much broader data coverage (50+ datasets)

---

## References

- FinMind Documentation: https://finmind.github.io/
- FinMind GitHub: https://github.com/FinMind/FinMind
- API Endpoint: https://api.finmindtrade.com/api/v4/data
- API Docs: https://api.finmindtrade.com/docs

---

**Status:** Planning Complete ✓
**Ready for Implementation:** Yes
**Next Step:** Begin Phase 16.1 - FinMind Reader Structure
