# gonp-datareader Project Specification

## Overview

gonp-datareader is a standalone Go package that provides remote data access for financial and economic data sources, designed to work seamlessly with the [gonp](https://github.com/julianshen/gonp) library. It serves as the Go equivalent of Python's pandas-datareader, enabling users to fetch data from various internet sources into gonp DataFrames.

## Project Goals

1. **Seamless Integration**: Work naturally with gonp DataFrames and Series
2. **Multiple Data Sources**: Support major financial and economic data providers
3. **Type Safety**: Leverage Go's type system for compile-time safety
4. **High Performance**: Efficient data fetching and parsing with minimal allocations
5. **Extensibility**: Easy to add new data sources through a plugin architecture
6. **Production Ready**: Comprehensive error handling, retry logic, and rate limiting

## Supported Data Sources (Priority Order)

### Phase 1: Core Data Sources
1. **Yahoo Finance** - Historical stock prices, dividends, splits
2. **FRED (Federal Reserve Economic Data)** - Economic indicators and time series
3. **World Bank** - Development indicators and global statistics
4. **Alpha Vantage** - Real-time and historical equity data (requires API key)

### Phase 2: Additional Sources
5. **Tiingo** - Financial data with API key
6. **IEX Cloud** - Investors Exchange data (requires API key)
7. **Stooq** - International stock market data
8. **OECD** - Economic statistics
9. **Eurostat** - European economic data

### Phase 3: Extended Sources
10. **Quandl** - Financial and economic datasets
11. **NASDAQ** - Listed symbols and market data
12. **Bank of Canada** - Canadian economic data
13. **Econdb** - Economic database access

## Package Structure

```
gonp-datareader/
├── datareader.go          # Main DataReader interface and factory functions
├── config.go              # Configuration and options
├── error.go               # Custom error types and handling
├── sources/               # Data source implementations
│   ├── source.go         # Base source interface
│   ├── yahoo/            # Yahoo Finance implementation
│   │   ├── yahoo.go
│   │   ├── yahoo_test.go
│   │   └── parser.go
│   ├── fred/             # FRED implementation
│   │   ├── fred.go
│   │   ├── fred_test.go
│   │   └── parser.go
│   ├── worldbank/        # World Bank implementation
│   │   ├── worldbank.go
│   │   ├── worldbank_test.go
│   │   └── parser.go
│   └── alphavantage/     # Alpha Vantage implementation
│       ├── alphavantage.go
│       ├── alphavantage_test.go
│       └── parser.go
├── internal/             # Internal utilities
│   ├── http/            # HTTP client with retry and rate limiting
│   │   ├── client.go
│   │   └── client_test.go
│   ├── cache/           # Optional response caching
│   │   ├── cache.go
│   │   └── cache_test.go
│   └── utils/           # Common utilities
│       ├── time.go      # Time/date parsing utilities
│       └── validation.go # Input validation
├── examples/             # Usage examples
│   ├── basic_usage/
│   ├── multiple_sources/
│   └── advanced_options/
├── docs/                 # Documentation
│   ├── api.md           # API reference
│   ├── sources.md       # Data source documentation
│   └── migration.md     # Migration guide from pandas-datareader
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md            # Development guidelines
├── CHANGELOG.md
└── LICENSE
```

## Core API Design

### Main DataReader Interface

```go
package datareader

import (
    "context"
    "time"
    
    "github.com/julianshen/gonp/dataframe"
    "github.com/julianshen/gonp/series"
)

// Reader is the main interface for all data sources
type Reader interface {
    // Read fetches data for the given symbol(s) within the date range
    Read(ctx context.Context, symbols []string, start, end time.Time) (*dataframe.DataFrame, error)
    
    // ReadSingle fetches data for a single symbol
    ReadSingle(ctx context.Context, symbol string, start, end time.Time) (*series.Series, error)
    
    // ValidateSymbol checks if a symbol is valid for this data source
    ValidateSymbol(symbol string) error
    
    // Name returns the name of the data source
    Name() string
}

// Options configures the data reader behavior
type Options struct {
    // APIKey for sources that require authentication
    APIKey string
    
    // Timeout for HTTP requests
    Timeout time.Duration
    
    // MaxRetries for failed requests
    MaxRetries int
    
    // RetryDelay between retry attempts
    RetryDelay time.Duration
    
    // EnableCache enables response caching
    EnableCache bool
    
    // CacheDir specifies the cache directory
    CacheDir string
    
    // RateLimit specifies requests per second limit
    RateLimit float64
    
    // UserAgent for HTTP requests
    UserAgent string
}

// DefaultOptions returns default configuration
func DefaultOptions() *Options

// DataReader creates a new reader for the specified source
func DataReader(source string, opts *Options) (Reader, error)

// Read is a convenience function that creates a reader and fetches data
func Read(ctx context.Context, symbol string, source string, start, end time.Time, opts *Options) (*dataframe.DataFrame, error)
```

### Usage Examples

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/julianshen/gonp/dataframe"
    dr "github.com/yourorg/gonp-datareader"
)

func main() {
    ctx := context.Background()
    
    // Example 1: Simple usage with defaults
    start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
    
    df, err := dr.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(df.Head())
    
    // Example 2: Using custom options
    opts := &dr.Options{
        Timeout:     30 * time.Second,
        MaxRetries:  3,
        EnableCache: true,
        CacheDir:    ".datareader_cache",
    }
    
    reader, err := dr.DataReader("yahoo", opts)
    if err != nil {
        panic(err)
    }
    
    df, err = reader.Read(ctx, []string{"AAPL", "MSFT", "GOOGL"}, start, end)
    if err != nil {
        panic(err)
    }
    
    // Example 3: FRED data with API key
    fredOpts := &dr.Options{
        APIKey: "your-fred-api-key",
    }
    
    gdp, err := dr.Read(ctx, "GDP", "fred", start, end, fredOpts)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("US GDP: %.2f\n", gdp.Loc(gdp.Index[0]))
    
    // Example 4: World Bank indicators
    wbReader, _ := dr.DataReader("worldbank", nil)
    
    // GDP per capita for multiple countries
    indicators := []string{"NY.GDP.PCAP.CD"}
    countries := []string{"US", "CN", "JP", "DE", "GB"}
    
    wbData, err := wbReader.Read(ctx, indicators, start, end)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(wbData.Head())
}
```

## Data Source Specifications

### Yahoo Finance

**Features:**
- Historical daily stock prices (OHLCV)
- Adjusted close prices
- Dividends and stock splits
- Multiple symbols in single request
- Free access, no API key required

**Symbol Format:** Standard ticker symbols (e.g., "AAPL", "MSFT")

**Date Range:** Historical data availability varies by symbol

**Output Format:** DataFrame with columns:
- Date (index)
- Open, High, Low, Close, Volume
- AdjClose (adjusted closing price)

### FRED (Federal Reserve Economic Data)

**Features:**
- Economic indicators and time series
- 500,000+ data series
- Federal Reserve and other government data
- API key recommended for higher rate limits

**Symbol Format:** FRED series IDs (e.g., "GDP", "UNRATE", "CPIAUCSL")

**Date Range:** Varies by series

**Output Format:** Series or DataFrame with:
- Date (index)
- Value for each series

### World Bank

**Features:**
- Development indicators
- Country-level statistics
- Historical data spanning decades
- No API key required

**Symbol Format:** Indicator codes (e.g., "NY.GDP.PCAP.CD" for GDP per capita)

**Date Range:** 1960 onwards (varies by indicator)

**Output Format:** DataFrame with:
- Date (index)
- Country (column)
- Indicator values

### Alpha Vantage

**Features:**
- Real-time and historical equity data
- Intraday data with various intervals
- Technical indicators
- Foreign exchange rates
- Cryptocurrency data
- Requires free API key

**Symbol Format:** Standard ticker symbols

**Date Range:** Real-time to 20+ years historical

**Output Format:** DataFrame with time-series data

## Error Handling Strategy

```go
// Custom error types for better error handling
type ErrorType int

const (
    ErrInvalidSymbol ErrorType = iota
    ErrInvalidDateRange
    ErrNetworkError
    ErrAPILimit
    ErrAuthenticationFailed
    ErrDataNotFound
    ErrParsingFailed
)

// DataReaderError provides detailed error information
type DataReaderError struct {
    Type    ErrorType
    Source  string
    Message string
    Cause   error
}

func (e *DataReaderError) Error() string
func (e *DataReaderError) Unwrap() error
func (e *DataReaderError) Is(target error) bool
```

## Performance Requirements

1. **HTTP Requests**: 
   - Use HTTP/2 where available
   - Connection pooling and keep-alive
   - Request timeout: 30 seconds default
   - Retry with exponential backoff

2. **Memory Efficiency**:
   - Stream large responses
   - Minimize allocations during parsing
   - Reuse buffers where appropriate

3. **Concurrency**:
   - Safe for concurrent use
   - Parallel fetching for multiple symbols
   - Rate limiting per source

4. **Caching**:
   - Optional file-based cache
   - Cache key: source + symbol + date range
   - TTL configurable per source

## Testing Requirements

1. **Unit Tests**: 
   - 80%+ code coverage
   - All data parsers tested with real examples
   - Error handling paths covered

2. **Integration Tests**:
   - Test against real APIs (with rate limiting)
   - Mock server tests for reliability
   - Date range boundary testing

3. **Benchmark Tests**:
   - Parsing performance
   - Memory allocation profiling
   - HTTP client overhead

4. **Example Tests**:
   - All examples must be runnable
   - Examples serve as documentation tests

## Documentation Requirements

1. **Package Documentation**:
   - Clear package-level overview
   - Godoc comments for all exported types
   - Usage examples in documentation

2. **README**:
   - Quick start guide
   - Installation instructions
   - Basic usage examples
   - Link to full documentation

3. **API Reference**:
   - Complete API documentation
   - Parameter descriptions
   - Return value specifications
   - Error conditions

4. **Source Documentation**:
   - Each source's capabilities
   - Symbol format requirements
   - Rate limits and restrictions
   - API key requirements

5. **Migration Guide**:
   - pandas-datareader to gonp-datareader
   - Code examples side-by-side
   - Feature parity matrix

## Security Considerations

1. **API Key Management**:
   - Environment variable support
   - No hardcoded keys in examples
   - Secure key storage recommendations

2. **Input Validation**:
   - Sanitize all user inputs
   - Validate date ranges
   - Check symbol formats

3. **HTTPS Only**:
   - All requests use HTTPS
   - Certificate verification enabled
   - TLS 1.2+ minimum

## Release Plan

### v0.1.0 - MVP (Phase 1)
- Core DataReader interface
- Yahoo Finance support
- FRED support
- Basic error handling
- Essential documentation
- 70%+ test coverage

### v0.2.0 - Extended Sources (Phase 2)
- World Bank support
- Alpha Vantage support
- Tiingo support
- Response caching
- 80%+ test coverage

### v0.3.0 - Production Ready (Phase 3)
- IEX Cloud support
- Stooq support
- Rate limiting
- Retry logic with backoff
- Comprehensive examples
- 85%+ test coverage

### v1.0.0 - Stable Release
- All Phase 3 sources
- Complete documentation
- Performance optimizations
- Production battle-tested
- 90%+ test coverage

## Dependencies

**Required:**
- Go 1.21+
- github.com/julianshen/gonp (latest)

**HTTP Client:**
- Standard library net/http
- Optional: golang.org/x/net/http2

**Testing:**
- Standard library testing
- github.com/stretchr/testify (assertions)

**Minimal External Dependencies:**
- Keep dependency tree shallow
- No unnecessary transitive dependencies
- Regular dependency updates

## Contributing Guidelines

1. Follow TDD methodology (see CLAUDE.md)
2. Write tests before implementation
3. Maintain code coverage above 80%
4. Follow Effective Go guidelines
5. Add examples for new features
6. Update documentation with changes
7. Run `make check` before submitting PR

## Success Metrics

1. **Functionality**: All priority data sources working reliably
2. **Performance**: Sub-second response for typical queries
3. **Reliability**: 99%+ success rate for valid requests
4. **Usability**: Clear API, comprehensive documentation
5. **Adoption**: Integration examples with gonp
6. **Quality**: High test coverage, clean code, minimal bugs

## Future Enhancements

1. **Additional Sources**: Crypto exchanges, commodities data
2. **Advanced Features**: Automatic symbol normalization
3. **Data Transformations**: Built-in resampling, filling
4. **Streaming**: Real-time data streaming support
5. **CLI Tool**: Command-line interface for data fetching
6. **gRPC API**: Server mode for data service
