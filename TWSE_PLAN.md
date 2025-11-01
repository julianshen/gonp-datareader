# TWSE (Taiwan Stock Exchange) Reader Implementation Plan

## Overview

Add Taiwan Stock Exchange (TWSE) as a new data source to gonp-datareader, providing access to Taiwan stock market data through the official TWSE Open API.

## API Characteristics

**Base URL:** `https://openapi.twse.com.tw/v1/`

**Authentication:** None required (public API)

**Data Format:** JSON

**Date Format:** ROC (Republic of China) calendar - e.g., "1141031" represents 2025-10-31 (114th year of ROC)

**Rate Limiting:** Unknown, but should implement conservative rate limiting (1-2 requests/second)

## Key Endpoints to Support

### 1. Daily Stock Data (Primary)
**Endpoint:** `/exchangeReport/STOCK_DAY_ALL`

**Description:** All stocks daily trading data for the current/latest trading day

**Response Fields:**
- `Date`: Trading date in ROC format (e.g., "1141031")
- `Code`: Stock symbol (e.g., "2330" for TSMC)
- `Name`: Company name in Traditional Chinese
- `TradeVolume`: Number of shares traded
- `TradeValue`: Total trade value
- `OpeningPrice`: Opening price
- `HighestPrice`: Daily high
- `LowestPrice`: Daily low
- `ClosingPrice`: Closing price
- `Change`: Price change from previous day
- `Transaction`: Number of transactions

**Data Format Example:**
```json
{
  "Date": "1141031",
  "Code": "2330",
  "Name": "台積電",
  "TradeVolume": "55956524",
  "TradeValue": "3616991558",
  "OpeningPrice": "64.60",
  "HighestPrice": "64.80",
  "LowestPrice": "64.40",
  "ClosingPrice": "64.75",
  "Change": "0.3500",
  "Transaction": "44302"
}
```

### 2. Market Indices (Secondary)
**Endpoint:** `/exchangeReport/MI_INDEX`

**Description:** Major market indices including weighted index

**Response Fields:**
- `日期` (Date): Date in ROC format
- `指數` (Index Name): Index name in Chinese
- `收盤指數` (Closing Index): Closing index value
- `漲跌` (Direction): Up/Down indicator
- `漲跌點數` (Change Points): Point change
- `漲跌百分比` (Change %): Percentage change

### 3. Listed Companies (Metadata - Optional)
**Endpoint:** `/opendata/t187ap03_L`

**Description:** Comprehensive list of all TWSE listed companies with detailed information

## Implementation Tasks (Following TDD)

### Phase 1: TWSE Reader Structure

#### Task 1.1: Create Package Structure
- ☐ Create `sources/twse/` directory
- ☐ Create `sources/twse/twse.go` main file
- ☐ Create `sources/twse/twse_test.go` test file
- ☐ Create `sources/twse/parser.go` for JSON parsing
- ☐ Create `sources/twse/parser_test.go` for parser tests

**Commit:** `feat: create TWSE reader package structure`

#### Task 1.2: TWSEReader Structure
- ☐ Test: TWSEReader struct exists
- ☐ Implement: TWSEReader in sources/twse/twse.go
- ☐ Test: TWSEReader embeds BaseSource
- ☐ Implement: Embed BaseSource
- ☐ Test: NewTWSEReader returns non-nil reader
- ☐ Implement: NewTWSEReader constructor
- ☐ Test: TWSEReader implements Reader interface
- ☐ Verify: All Reader methods present

**Commit:** `feat: implement TWSE reader structure`

### Phase 2: Date Conversion Utilities

**Critical Component:** ROC ↔ Gregorian calendar conversion

#### Task 2.1: Date Conversion Functions
- ☐ Test: rocToGregorian converts "1141031" to time.Time correctly
- ☐ Implement: rocToGregorian function (add 1911 years)
- ☐ Test: gregorianToROC converts time.Time to "1141031" correctly
- ☐ Implement: gregorianToROC function (subtract 1911 years)
- ☐ Test: parseROCDate handles "1141031" format
- ☐ Implement: parseROCDate with format "YYYMMDD"
- ☐ Test: formatROCDate creates correct ROC string
- ☐ Implement: formatROCDate function
- ☐ Test: Edge cases (leap years, year boundaries)
- ☐ Implement: Comprehensive date handling

**Formula:** ROC Year = Gregorian Year - 1911

**Example:**
- ROC 114/10/31 = 2025/10/31
- ROC 113/01/01 = 2024/01/01

**Commit:** `feat: add ROC calendar conversion utilities`

### Phase 3: URL Building

#### Task 3.1: Build URL for Daily Data
- ☐ Test: buildDailyURL creates valid endpoint URL
- ☐ Implement: buildDailyURL function
- ☐ Test: buildDailyURL uses correct base URL
- ☐ Implement: URL constant and formatting
- ☐ Test: buildIndexURL creates valid index endpoint
- ☐ Implement: buildIndexURL function

**URL Format:**
```go
const (
    twseBaseURL = "https://openapi.twse.com.tw/v1"
    dailyStocksEndpoint = "/exchangeReport/STOCK_DAY_ALL"
    indexEndpoint = "/exchangeReport/MI_INDEX"
)
```

**Commit:** `feat: implement TWSE URL builders`

### Phase 4: JSON Response Parsing

#### Task 4.1: Parse Daily Stock Data
- ☐ Test: parseDailyStockJSON parses valid response
- ☐ Implement: parseDailyStockJSON function
- ☐ Test: parseDailyStockJSON extracts stock code
- ☐ Implement: Extract Code field
- ☐ Test: parseDailyStockJSON extracts OHLC data
- ☐ Implement: Parse OpeningPrice, HighestPrice, LowestPrice, ClosingPrice
- ☐ Test: parseDailyStockJSON converts string numbers to float64
- ☐ Implement: String to float conversion with error handling
- ☐ Test: parseDailyStockJSON handles missing/null values
- ☐ Implement: Null value handling
- ☐ Test: parseDailyStockJSON converts ROC date to time.Time
- ☐ Implement: Use rocToGregorian for date conversion
- ☐ Test: parseDailyStockJSON extracts volume and transaction count
- ☐ Implement: Parse TradeVolume and Transaction

**Response Structure:**
```go
type TWSEStockData struct {
    Date          string `json:"Date"`          // ROC format
    Code          string `json:"Code"`
    Name          string `json:"Name"`
    TradeVolume   string `json:"TradeVolume"`   // String numbers
    TradeValue    string `json:"TradeValue"`
    OpeningPrice  string `json:"OpeningPrice"`
    HighestPrice  string `json:"HighestPrice"`
    LowestPrice   string `json:"LowestPrice"`
    ClosingPrice  string `json:"ClosingPrice"`
    Change        string `json:"Change"`
    Transaction   string `json:"Transaction"`
}

type ParsedData struct {
    Symbol       string
    Date         []time.Time
    Open         []float64
    High         []float64
    Low          []float64
    Close        []float64
    Volume       []int64
    Transactions []int64
    Change       []float64
}
```

**Commit:** `feat: implement TWSE JSON parser for daily stock data`

#### Task 4.2: Filter Data by Symbol and Date Range
- ☐ Test: filterBySymbol extracts single symbol from all stocks
- ☐ Implement: Filter function for symbol matching
- ☐ Test: filterByDateRange filters data within start/end dates
- ☐ Implement: Date range filtering (note: API returns only latest day)
- ☐ Test: Handle case when symbol not found
- ☐ Implement: Return appropriate error

**Note:** TWSE API returns only the latest trading day data, not historical ranges. May need to call multiple times or use different endpoint for historical data.

**Commit:** `feat: add symbol and date filtering for TWSE data`

### Phase 5: Reader Integration

#### Task 5.1: Implement ReadSingle
- ☐ Test: TWSEReader.ReadSingle fetches data for "2330" (TSMC)
- ☐ Implement: ReadSingle method
- ☐ Test: TWSEReader.ReadSingle validates symbol
- ☐ Implement: Symbol validation (4-digit Taiwan stock codes)
- ☐ Test: TWSEReader.ReadSingle validates date range
- ☐ Implement: Date validation
- ☐ Test: TWSEReader.ReadSingle returns ParsedData
- ☐ Implement: Complete integration with parser
- ☐ Test: TWSEReader.ReadSingle handles HTTP errors
- ☐ Implement: Error handling for network issues
- ☐ Test: TWSEReader.ReadSingle handles invalid symbols
- ☐ Implement: Return ErrDataNotFound for missing symbols

**Commit:** `feat: implement TWSE ReadSingle method`

#### Task 5.2: Implement Read for Multiple Symbols
- ☐ Test: TWSEReader.Read handles multiple symbols
- ☐ Implement: Parallel fetching pattern (similar to Yahoo/Stooq)
- ☐ Test: TWSEReader.Read respects rate limits
- ☐ Implement: Rate limiting between requests
- ☐ Test: TWSEReader.Read combines results into map
- ☐ Implement: Result aggregation

**Commit:** `feat: implement TWSE Read method for multiple symbols`

### Phase 6: Error Handling

#### Task 6.1: TWSE-Specific Errors
- ☐ Test: Returns error for invalid Taiwan stock code format
- ☐ Implement: Stock code validation (typically 4 digits)
- ☐ Test: Returns error for non-trading days
- ☐ Implement: Handle weekends/holidays gracefully
- ☐ Test: Returns error for API timeout
- ☐ Implement: Timeout handling
- ☐ Test: Returns descriptive error messages
- ☐ Implement: Error message formatting

**Commit:** `feat: add comprehensive error handling for TWSE reader`

### Phase 7: Factory Registration

#### Task 7.1: Register with DataReader Factory
- ☐ Test: DataReader("twse") returns TWSE reader
- ☐ Implement: Register in init() function
- ☐ Test: Read with "twse" source works end-to-end
- ☐ Verify: Complete factory integration
- ☐ Test: TWSE reader available in source list
- ☐ Implement: Add to supported sources

**Registration Code:**
```go
func init() {
    RegisterSource("twse", func(opts *Options) (Reader, error) {
        return NewTWSEReader(opts.ToClientOptions()), nil
    })
}
```

**Commit:** `feat: register TWSE reader with factory`

### Phase 8: Documentation

#### Task 8.1: Code Documentation
- ☐ Add package-level godoc for twse package
- ☐ Document TWSEReader struct
- ☐ Document all exported functions
- ☐ Add usage examples in godoc
- ☐ Document ROC calendar conversion

**Commit:** `docs: add TWSE reader godoc documentation`

#### Task 8.2: Update sources.md
- ☐ Add TWSE section to docs/sources.md
- ☐ Document TWSE API capabilities
- ☐ Document symbol format (4-digit Taiwan codes)
- ☐ Document date format and limitations
- ☐ Add data limitations (latest day only)
- ☐ Add rate limit recommendations
- ☐ Add example Taiwan stock symbols

**Commit:** `docs: add TWSE to data sources documentation`

#### Task 8.3: Create Example
- ☐ Create examples/twse/main.go
- ☐ Example: Fetch TSMC (2330) data
- ☐ Example: Fetch multiple Taiwan stocks
- ☐ Example: Error handling
- ☐ Add README for TWSE example

**Commit:** `docs: add TWSE usage example`

#### Task 8.4: Update Migration Guide
- ☐ Add TWSE to docs/migration.md
- ☐ Compare with pandas-datareader (if Taiwan support exists)
- ☐ Add Taiwan-specific notes

**Commit:** `docs: add TWSE to migration guide`

### Phase 9: Testing

#### Task 9.1: Unit Tests with Mock Data
- ☐ Test: Parse valid TWSE JSON response
- ☐ Test: Handle malformed JSON
- ☐ Test: Handle empty response array
- ☐ Test: ROC date conversion edge cases
- ☐ Test: String to number conversion errors
- ☐ Test: Symbol validation (valid/invalid formats)

**Test Data:** Use actual TWSE API response samples

**Commit:** `test: add comprehensive unit tests for TWSE reader`

#### Task 9.2: Integration Tests (Optional)
- ☐ Test: Fetch real data from TWSE API
- ☐ Test: Verify data structure correctness
- ☐ Test: Test rate limiting behavior

**Commit:** `test: add TWSE integration tests`

### Phase 10: Update README and CHANGELOG

#### Task 10.1: Update README
- ☐ Add TWSE to supported sources list
- ☐ Add TWSE to quick start examples
- ☐ Update feature matrix

**Commit:** `docs: add TWSE to README`

#### Task 10.2: Update CHANGELOG
- ☐ Add TWSE feature to CHANGELOG.md
- ☐ Document any API changes
- ☐ Note version for release

**Commit:** `docs: update CHANGELOG for TWSE support`

## Symbol Format

**Taiwan Stock Exchange symbols are typically 4-6 digits:**
- Regular stocks: 4 digits (e.g., "2330" for TSMC)
- ETFs: 4 digits starting with 00 (e.g., "0050" for 元大台灣50)
- Warrants: 6 digits

**Popular Taiwan Stocks for Testing:**
- 2330 - 台積電 (TSMC)
- 2317 - 鴻海 (Hon Hai/Foxconn)
- 2454 - 聯發科 (MediaTek)
- 2412 - 中華電 (Chunghwa Telecom)
- 0050 - 元大台灣50 (Taiwan 50 ETF)

## Validation Rules

1. **Symbol Format:** 4-6 digits, all numeric
2. **Date Handling:**
   - API returns latest trading day only
   - For historical data, may need to call different endpoints
3. **ROC Calendar:** All dates must be converted to/from ROC format

## Rate Limiting Recommendations

- Conservative: 1-2 requests per second
- Add configurable rate limit in Options
- Use existing rate limiter infrastructure

## Potential Challenges

1. **Historical Data:** Current endpoint may only return latest day
   - May need to explore other TWSE endpoints for historical data
   - Could require daily polling/caching strategy

2. **Chinese Characters:** Company names are in Traditional Chinese
   - Keep as-is for authenticity
   - Could add English name mapping (future enhancement)

3. **Trading Calendar:** Taiwan market holidays differ from US markets
   - No data on non-trading days
   - Should handle gracefully

4. **Data Availability:** Real-time vs delayed data
   - Verify data freshness
   - Document update frequency

## Future Enhancements

1. **Historical Data Support:**
   - Explore TWSE historical endpoints
   - Implement data caching strategy

2. **Additional Endpoints:**
   - Company fundamentals
   - Institutional holdings
   - Revenue reports

3. **Symbol Lookup:**
   - Company name to symbol mapping
   - English/Chinese name search

4. **Market Indices:**
   - TAIEX (Taiwan Weighted Index)
   - Sector indices

## Success Criteria

- ☐ Can fetch latest trading data for any valid Taiwan stock symbol
- ☐ Proper ROC calendar conversion
- ☐ Clean JSON parsing with error handling
- ☐ Integration with existing gonp-datareader architecture
- ☐ >80% test coverage
- ☐ Comprehensive documentation
- ☐ Working examples

## Timeline Estimate

- Phase 1-3 (Structure + Date Utils): 2-3 hours
- Phase 4-5 (Parsing + Integration): 3-4 hours
- Phase 6-7 (Error Handling + Registration): 1-2 hours
- Phase 8-10 (Documentation + Testing): 2-3 hours

**Total:** 8-12 hours of focused development following TDD

## Notes

- Follow existing patterns from Stooq/Yahoo readers (CSV) and FRED reader (JSON)
- TWSE uses JSON like FRED, but simpler structure
- ROC calendar conversion is the unique challenge
- No API key required makes it accessible for all users
