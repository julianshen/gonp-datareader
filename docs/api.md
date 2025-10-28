# API Reference

Complete API reference for gonp-datareader with detailed examples and usage patterns.

---

## Table of Contents

1. [Core Functions](#core-functions)
2. [Reader Interface](#reader-interface)
3. [Options Structure](#options-structure)
4. [Error Types](#error-types)
5. [Data Types](#data-types)
6. [Usage Examples](#usage-examples)
7. [Advanced Patterns](#advanced-patterns)

---

## Core Functions

### Read

The `Read` function is the main convenience function for fetching data from any supported source.

```go
func Read(
    ctx context.Context,
    symbols interface{},
    source string,
    start, end time.Time,
    opts *Options,
) (interface{}, error)
```

**Parameters:**
- `ctx` - Context for cancellation and timeout control
- `symbols` - Single symbol (string) or multiple symbols ([]string)
- `source` - Data source name (e.g., "yahoo", "fred", "worldbank")
- `start` - Start date for data range
- `end` - End date for data range
- `opts` - Optional configuration (can be nil for defaults)

**Returns:**
- Source-specific data structure (see [Data Types](#data-types))
- Error if request fails

**Example:**
```go
import (
    "context"
    "time"
    "github.com/julianshen/gonp-datareader"
)

ctx := context.Background()
start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Now()

// Single symbol
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
if err != nil {
    log.Fatal(err)
}

// Multiple symbols
symbols := []string{"AAPL", "MSFT", "GOOGL"}
dataMap, err := datareader.Read(ctx, symbols, "yahoo", start, end, nil)
```

---

### DataReader

The `DataReader` function creates a reusable reader for a specific data source.

```go
func DataReader(source string, opts *Options) (Reader, error)
```

**Parameters:**
- `source` - Data source name
- `opts` - Configuration options

**Returns:**
- `Reader` interface implementation
- Error if source is unknown or configuration is invalid

**Example:**
```go
opts := &datareader.Options{
    APIKey:   "your-api-key",
    Timeout:  60 * time.Second,
    CacheDir: ".cache",
}

reader, err := datareader.DataReader("alphavantage", opts)
if err != nil {
    log.Fatal(err)
}

// Reuse reader for multiple calls
data1, err := reader.ReadSingle(ctx, "AAPL", start, end)
data2, err := reader.ReadSingle(ctx, "MSFT", start, end)
```

---

### ListSources

Returns a list of all supported data sources.

```go
func ListSources() []string
```

**Returns:**
- Array of supported source names

**Example:**
```go
sources := datareader.ListSources()
fmt.Println("Supported sources:", sources)
// Output: [yahoo fred worldbank alphavantage stooq iex tiingo oecd eurostat]
```

---

## Reader Interface

The `Reader` interface defines the contract that all data source implementations must follow.

```go
type Reader interface {
    // Name returns the name of the data source
    Name() string

    // Read fetches data for one or more symbols
    Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error)

    // ReadSingle fetches data for a single symbol
    ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error)
}
```

### Methods

#### Name()

Returns the identifier for the data source.

```go
reader, _ := datareader.DataReader("yahoo", nil)
fmt.Println(reader.Name()) // Output: "yahoo"
```

#### Read()

Fetches data for multiple symbols (parallel fetching where supported).

```go
symbols := []string{"AAPL", "MSFT", "GOOGL"}
data, err := reader.Read(ctx, symbols, start, end)
if err != nil {
    log.Fatal(err)
}

// Result type depends on source
dataMap := data.(map[string]*yahoo.ParsedData)
```

#### ReadSingle()

Fetches data for a single symbol.

```go
data, err := reader.ReadSingle(ctx, "AAPL", start, end)
if err != nil {
    log.Fatal(err)
}

parsedData := data.(*yahoo.ParsedData)
```

---

## Options Structure

The `Options` structure configures reader behavior.

```go
type Options struct {
    // APIKey for sources requiring authentication
    APIKey string

    // HTTP client configuration
    Timeout   time.Duration
    UserAgent string

    // Retry configuration
    MaxRetries int
    RetryDelay time.Duration

    // Rate limiting (requests per second)
    RateLimit float64

    // Response caching
    CacheDir string
    CacheTTL time.Duration
}
```

### Fields

#### APIKey

API key for authenticated sources (Alpha Vantage, IEX Cloud, Tiingo, etc.).

```go
opts := &datareader.Options{
    APIKey: "your-api-key",
}
```

**Environment Variable Pattern:**
```go
opts := &datareader.Options{
    APIKey: os.Getenv("ALPHAVANTAGE_API_KEY"),
}
```

#### Timeout

HTTP request timeout duration.

```go
opts := &datareader.Options{
    Timeout: 60 * time.Second, // 60 second timeout
}
```

**Default:** 30 seconds

#### UserAgent

Custom User-Agent header for HTTP requests.

```go
opts := &datareader.Options{
    UserAgent: "MyApp/1.0",
}
```

**Default:** "gonp-datareader/v0.1.0"

#### MaxRetries

Maximum number of retry attempts on transient failures.

```go
opts := &datareader.Options{
    MaxRetries: 5, // Retry up to 5 times
}
```

**Default:** 3 retries

#### RetryDelay

Initial delay between retry attempts (exponential backoff).

```go
opts := &datareader.Options{
    RetryDelay: 2 * time.Second, // Start with 2 second delay
}
```

**Default:** 1 second

#### RateLimit

Maximum requests per second (token bucket algorithm).

```go
opts := &datareader.Options{
    RateLimit: 5.0, // Max 5 requests per second
}
```

**Default:** No rate limiting (0.0)

**Use Case:** Comply with API rate limits
```go
// Alpha Vantage free tier: 5 requests/minute
opts := &datareader.Options{
    RateLimit: 5.0 / 60.0, // 5 per minute = 0.083 per second
}
```

#### CacheDir

Directory for caching API responses.

```go
opts := &datareader.Options{
    CacheDir: ".cache/datareader",
}
```

**Default:** No caching (empty string)

#### CacheTTL

Time-to-live for cached responses.

```go
opts := &datareader.Options{
    CacheDir: ".cache",
    CacheTTL: 24 * time.Hour, // Cache for 24 hours
}
```

**Default:** No expiration (0)

### Complete Example

```go
opts := &datareader.Options{
    // Authentication
    APIKey: os.Getenv("ALPHAVANTAGE_API_KEY"),

    // HTTP Configuration
    Timeout:   60 * time.Second,
    UserAgent: "MyTradingBot/1.0",

    // Retry Configuration
    MaxRetries: 5,
    RetryDelay: 2 * time.Second,

    // Rate Limiting (5 per minute for Alpha Vantage free tier)
    RateLimit: 5.0 / 60.0,

    // Caching
    CacheDir: ".cache/alphavantage",
    CacheTTL: 24 * time.Hour,
}

reader, err := datareader.DataReader("alphavantage", opts)
```

---

## Error Types

### Standard Errors

```go
var (
    // ErrInvalidSymbol is returned when a symbol is invalid or empty
    ErrInvalidSymbol = errors.New("invalid symbol")

    // ErrInvalidDateRange is returned when the date range is invalid
    ErrInvalidDateRange = errors.New("invalid date range")

    // ErrUnknownSource is returned when the data source is not supported
    ErrUnknownSource = errors.New("unknown data source")
)
```

### Error Checking

Use `errors.Is()` to check for specific error types:

```go
data, err := datareader.Read(ctx, "", "yahoo", start, end, nil)
if err != nil {
    if errors.Is(err, datareader.ErrInvalidSymbol) {
        fmt.Println("Symbol is invalid or empty")
    } else if errors.Is(err, datareader.ErrUnknownSource) {
        fmt.Println("Data source not supported")
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Wrapped Errors

Errors are wrapped with context using `fmt.Errorf("%w")`:

```go
data, err := datareader.Read(ctx, "INVALID", "yahoo", start, end, nil)
if err != nil {
    fmt.Println(err)
    // Output: failed to fetch data: yahoo finance returned status 404: Not Found

    // Unwrap to get root cause
    rootErr := errors.Unwrap(err)
    fmt.Println(rootErr)
}
```

---

## Data Types

### Yahoo Finance

```go
type ParsedData struct {
    Columns []string
    Rows    []map[string]string
}
```

**Usage:**
```go
data, _ := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
yahooData := data.(*yahoo.ParsedData)

for _, row := range yahooData.Rows {
    fmt.Printf("Date: %s, Close: %s\n", row["Date"], row["Close"])
}
```

### FRED

```go
type ParsedData struct {
    Dates  []string
    Values []float64
}
```

**Usage:**
```go
data, _ := datareader.Read(ctx, "GDP", "fred", start, end, nil)
fredData := data.(*fred.ParsedData)

for i, date := range fredData.Dates {
    fmt.Printf("Date: %s, GDP: %.2f\n", date, fredData.Values[i])
}
```

### World Bank

```go
type ParsedData struct {
    Country      string
    Indicator    string
    Observations []Observation
}

type Observation struct {
    Date  string
    Value string
}
```

**Usage:**
```go
data, _ := datareader.Read(ctx, "USA/NY.GDP.MKTP.CD", "worldbank", start, end, nil)
wbData := data.(*worldbank.ParsedData)

for _, obs := range wbData.Observations {
    fmt.Printf("%s: %s\n", obs.Date, obs.Value)
}
```

### Alpha Vantage

```go
type ParsedData struct {
    Dates  []string
    Prices []PriceData
}

type PriceData struct {
    Open   float64
    High   float64
    Low    float64
    Close  float64
    Volume int64
}
```

### Stooq

```go
type ParsedData struct {
    Columns []string
    Rows    []map[string]string
}
```

### IEX Cloud

```go
type ParsedData struct {
    Rows []map[string]string
}
```

### Tiingo

```go
type ParsedData struct {
    Dates  []string
    Prices []PriceData
}

type PriceData struct {
    Close  float64
    Open   float64
    High   float64
    Low    float64
    Volume int64
}
```

### OECD

```go
type ParsedData struct {
    Dates  []string
    Values []float64
}
```

### Eurostat

```go
type ParsedData struct {
    Dates  []string
    Values []float64
}
```

---

## Usage Examples

### Example 1: Basic Stock Data

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        log.Fatal(err)
    }

    yahooData := data.(*yahoo.ParsedData)
    fmt.Printf("Fetched %d rows\n", len(yahooData.Rows))

    // Print last 5 days
    n := len(yahooData.Rows)
    for i := max(0, n-5); i < n; i++ {
        row := yahooData.Rows[i]
        fmt.Printf("%s: Close=%s, Volume=%s\n",
            row["Date"], row["Close"], row["Volume"])
    }
}
```

### Example 2: Multiple Symbols with Caching

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    // Configure caching
    opts := &datareader.Options{
        CacheDir: ".cache/yahoo",
        CacheTTL: 24 * time.Hour,
    }

    // Fetch multiple symbols (parallel)
    symbols := []string{"AAPL", "MSFT", "GOOGL", "TSLA"}
    result, err := datareader.Read(ctx, symbols, "yahoo", start, end, opts)
    if err != nil {
        log.Fatal(err)
    }

    dataMap := result.(map[string]*yahoo.ParsedData)

    for symbol, data := range dataMap {
        if len(data.Rows) > 0 {
            lastRow := data.Rows[len(data.Rows)-1]
            fmt.Printf("%s: Latest close = %s\n", symbol, lastRow["Close"])
        }
    }
}
```

### Example 3: Economic Data with API Key

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/fred"
)

func main() {
    ctx := context.Background()
    start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    opts := &datareader.Options{
        APIKey: os.Getenv("FRED_API_KEY"),
    }

    // Fetch GDP data
    result, err := datareader.Read(ctx, "GDP", "fred", start, end, opts)
    if err != nil {
        log.Fatal(err)
    }

    fredData := result.(*fred.ParsedData)

    fmt.Println("US GDP (Billions):")
    for i, date := range fredData.Dates {
        fmt.Printf("%s: $%.2fB\n", date, fredData.Values[i])
    }
}
```

### Example 4: Reusable Reader

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
)

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    // Create reader once
    opts := &datareader.Options{
        Timeout:    60 * time.Second,
        MaxRetries: 5,
        CacheDir:   ".cache/yahoo",
    }

    reader, err := datareader.DataReader("yahoo", opts)
    if err != nil {
        log.Fatal(err)
    }

    // Reuse for multiple symbols
    symbols := []string{"AAPL", "MSFT", "GOOGL"}

    for _, symbol := range symbols {
        data, err := reader.ReadSingle(ctx, symbol, start, end)
        if err != nil {
            log.Printf("Error fetching %s: %v\n", symbol, err)
            continue
        }

        fmt.Printf("Successfully fetched %s\n", symbol)
    }
}
```

### Example 5: Timeout and Context Cancellation

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
)

func main() {
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            log.Fatal("Request timed out after 5 seconds")
        }
        log.Fatal(err)
    }

    fmt.Println("Data fetched successfully")
}
```

### Example 6: Rate Limiting

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/alphavantage"
)

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    // Configure for Alpha Vantage free tier (5 calls/minute)
    opts := &datareader.Options{
        APIKey:    os.Getenv("ALPHAVANTAGE_API_KEY"),
        RateLimit: 5.0 / 60.0, // 5 per minute
    }

    reader, err := datareader.DataReader("alphavantage", opts)
    if err != nil {
        log.Fatal(err)
    }

    // Fetch multiple symbols (rate limiting applied automatically)
    symbols := []string{"AAPL", "MSFT", "GOOGL", "TSLA", "AMZN"}

    for _, symbol := range symbols {
        data, err := reader.ReadSingle(ctx, symbol, start, end)
        if err != nil {
            log.Printf("Error: %v\n", err)
            continue
        }

        avData := data.(*alphavantage.ParsedData)
        fmt.Printf("%s: %d data points\n", symbol, len(avData.Dates))
    }
}
```

### Example 7: Error Handling Patterns

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
)

func fetchWithRetry(ctx context.Context, symbol, source string, start, end time.Time) error {
    opts := &datareader.Options{
        MaxRetries: 3,
        RetryDelay: 2 * time.Second,
    }

    data, err := datareader.Read(ctx, symbol, source, start, end, opts)
    if err != nil {
        // Check for specific errors
        if errors.Is(err, datareader.ErrInvalidSymbol) {
            return fmt.Errorf("invalid symbol %q", symbol)
        }
        if errors.Is(err, datareader.ErrUnknownSource) {
            return fmt.Errorf("unsupported source %q", source)
        }
        if errors.Is(err, context.DeadlineExceeded) {
            return fmt.Errorf("request timed out")
        }

        // Generic error
        return fmt.Errorf("failed to fetch data: %w", err)
    }

    fmt.Printf("Successfully fetched %s from %s\n", symbol, source)
    _ = data // Process data
    return nil
}

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    if err := fetchWithRetry(ctx, "AAPL", "yahoo", start, end); err != nil {
        log.Fatal(err)
    }
}
```

---

## Advanced Patterns

### Pattern 1: Concurrent Fetching from Multiple Sources

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/julianshen/gonp-datareader"
)

func fetchFromSource(ctx context.Context, symbol, source string, start, end time.Time, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()

    data, err := datareader.Read(ctx, symbol, source, start, end, nil)
    if err != nil {
        results <- fmt.Sprintf("%s: Error - %v", source, err)
        return
    }

    results <- fmt.Sprintf("%s: Success", source)
    _ = data
}

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    sources := []string{"yahoo", "stooq"}
    results := make(chan string, len(sources))
    var wg sync.WaitGroup

    for _, source := range sources {
        wg.Add(1)
        go fetchFromSource(ctx, "AAPL", source, start, end, &wg, results)
    }

    wg.Wait()
    close(results)

    for result := range results {
        fmt.Println(result)
    }
}
```

### Pattern 2: Data Aggregation

```go
package main

import (
    "context"
    "fmt"
    "strconv"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/yahoo"
)

func calculateAverageClose(data *yahoo.ParsedData) (float64, error) {
    var sum float64
    var count int

    for _, row := range data.Rows {
        closeStr := row["Close"]
        close, err := strconv.ParseFloat(closeStr, 64)
        if err != nil {
            continue
        }
        sum += close
        count++
    }

    if count == 0 {
        return 0, fmt.Errorf("no valid closing prices")
    }

    return sum / float64(count), nil
}

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        panic(err)
    }

    yahooData := data.(*yahoo.ParsedData)
    avgClose, err := calculateAverageClose(yahooData)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Average closing price: $%.2f\n", avgClose)
}
```

### Pattern 3: Custom Configuration per Source

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/julianshen/gonp-datareader"
)

type Config struct {
    Yahoo struct {
        CacheDir string
        CacheTTL time.Duration
    }
    AlphaVantage struct {
        APIKey     string
        RateLimit  float64
        MaxRetries int
    }
}

func loadConfig() *Config {
    cfg := &Config{}

    cfg.Yahoo.CacheDir = ".cache/yahoo"
    cfg.Yahoo.CacheTTL = 24 * time.Hour

    cfg.AlphaVantage.APIKey = os.Getenv("ALPHAVANTAGE_API_KEY")
    cfg.AlphaVantage.RateLimit = 5.0 / 60.0
    cfg.AlphaVantage.MaxRetries = 5

    return cfg
}

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    cfg := loadConfig()

    // Yahoo with caching
    yahooOpts := &datareader.Options{
        CacheDir: cfg.Yahoo.CacheDir,
        CacheTTL: cfg.Yahoo.CacheTTL,
    }
    yahooReader, _ := datareader.DataReader("yahoo", yahooOpts)

    // Alpha Vantage with rate limiting
    avOpts := &datareader.Options{
        APIKey:     cfg.AlphaVantage.APIKey,
        RateLimit:  cfg.AlphaVantage.RateLimit,
        MaxRetries: cfg.AlphaVantage.MaxRetries,
    }
    avReader, _ := datareader.DataReader("alphavantage", avOpts)

    // Use readers
    data1, _ := yahooReader.ReadSingle(ctx, "AAPL", start, end)
    data2, _ := avReader.ReadSingle(ctx, "MSFT", start, end)

    fmt.Printf("Yahoo: %v\n", data1 != nil)
    fmt.Printf("Alpha Vantage: %v\n", data2 != nil)
}
```

### Pattern 4: Graceful Degradation

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/julianshen/gonp-datareader"
)

func fetchWithFallback(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
    sources := []string{"yahoo", "stooq"}

    var lastErr error
    for _, source := range sources {
        data, err := datareader.Read(ctx, symbol, source, start, end, nil)
        if err == nil {
            fmt.Printf("Successfully fetched from %s\n", source)
            return data, nil
        }

        fmt.Printf("Failed to fetch from %s: %v\n", source, err)
        lastErr = err
    }

    return nil, fmt.Errorf("all sources failed, last error: %w", lastErr)
}

func main() {
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    data, err := fetchWithFallback(ctx, "AAPL", start, end)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Data fetched successfully: %v\n", data != nil)
}
```

---

## Best Practices

### 1. Always Use Context

```go
// ✅ Good: Use context for cancellation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
```

```go
// ❌ Bad: Using context.Background() without timeout for long operations
ctx := context.Background() // No timeout!
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
```

### 2. Handle Errors Explicitly

```go
// ✅ Good: Check and handle errors
data, err := datareader.Read(ctx, symbol, source, start, end, nil)
if err != nil {
    return fmt.Errorf("failed to fetch %s: %w", symbol, err)
}
```

```go
// ❌ Bad: Ignoring errors
data, _ := datareader.Read(ctx, symbol, source, start, end, nil)
```

### 3. Reuse Readers When Possible

```go
// ✅ Good: Create reader once, reuse
reader, err := datareader.DataReader("yahoo", opts)
if err != nil {
    return err
}

for _, symbol := range symbols {
    data, err := reader.ReadSingle(ctx, symbol, start, end)
    // Process data
}
```

```go
// ❌ Bad: Creating new reader for each call
for _, symbol := range symbols {
    reader, _ := datareader.DataReader("yahoo", opts)
    data, _ := reader.ReadSingle(ctx, symbol, start, end)
}
```

### 4. Enable Caching for Production

```go
// ✅ Good: Use caching to reduce API calls
opts := &datareader.Options{
    CacheDir: ".cache",
    CacheTTL: 24 * time.Hour,
}
```

### 5. Configure Rate Limiting

```go
// ✅ Good: Respect API rate limits
opts := &datareader.Options{
    RateLimit: 5.0 / 60.0, // 5 per minute
}
```

### 6. Type Assert Safely

```go
// ✅ Good: Safe type assertion
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
if err != nil {
    return err
}

yahooData, ok := data.(*yahoo.ParsedData)
if !ok {
    return fmt.Errorf("unexpected data type")
}
```

```go
// ❌ Bad: Unsafe type assertion (can panic)
data, _ := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
yahooData := data.(*yahoo.ParsedData) // May panic!
```

---

## See Also

- [Source Documentation](sources.md) - Detailed info on each data source
- [Migration Guide](migration.md) - Migrating from pandas-datareader
- [Examples](../examples/) - Complete working examples
- [godoc](https://pkg.go.dev/github.com/julianshen/gonp-datareader) - API reference on pkg.go.dev

---

For more information and updates, visit the [GitHub repository](https://github.com/julianshen/gonp-datareader).
