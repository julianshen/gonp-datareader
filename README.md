# gonp-datareader

Remote data access for financial and economic data sources in Go, designed to work with [gonp](https://github.com/julianshen/gonp) DataFrames.

[![Go Reference](https://pkg.go.dev/badge/github.com/julianshen/gonp-datareader.svg)](https://pkg.go.dev/github.com/julianshen/gonp-datareader)
[![Go Report Card](https://goreportcard.com/badge/github.com/julianshen/gonp-datareader)](https://goreportcard.com/report/github.com/julianshen/gonp-datareader)

## Overview

gonp-datareader is the Go equivalent of Python's pandas-datareader, providing a simple and unified interface to fetch financial and economic data from various internet sources. It offers built-in support for automatic retries, rate limiting, caching, and flexible configuration.

## Features

- **Multiple Data Sources**: Yahoo Finance, FRED, World Bank, Alpha Vantage, Stooq, IEX Cloud
- **Simple API**: Easy-to-use interface for fetching financial and economic data
- **Automatic Retries**: Built-in retry logic with exponential backoff
- **Rate Limiting**: Token bucket rate limiting to respect API limits
- **Response Caching**: File-based caching with TTL support
- **Type Safe**: Leverages Go's type system for compile-time safety
- **Context Support**: Full context.Context support for cancellation and timeouts
- **Concurrent Safe**: Safe for concurrent use across goroutines
- **Extensible**: Easy to add new data sources through plugin architecture

## Installation

```bash
go get github.com/julianshen/gonp-datareader
```

## Quick Start

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

    // Fetch stock data from Yahoo Finance
    data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Fetched stock data: %+v\n", data)
}
```

## Supported Data Sources

| Source | Description | API Key Required | Symbol Format |
|--------|-------------|------------------|---------------|
| **yahoo** | Yahoo Finance - Stock prices, OHLCV data | No | `AAPL`, `MSFT` |
| **fred** | Federal Reserve Economic Data | Optional* | `GDP`, `UNRATE` |
| **worldbank** | World Bank Development Indicators | No | `USA/NY.GDP.MKTP.CD` |
| **alphavantage** | Alpha Vantage - Real-time & historical data | Yes | `AAPL`, `MSFT` |
| **stooq** | Stooq - Free international stock data | No | `AAPL.US`, `^SPX` |
| **iex** | IEX Cloud - Professional stock market data | Yes | `AAPL`, `MSFT` |

*FRED works without an API key but has lower rate limits

## API Key Configuration

Some sources require API keys. Configure them via environment variables or the `Options` struct:

```go
// Using Options struct
opts := &datareader.Options{
    APIKey: "your-api-key-here",
}

reader, err := datareader.DataReader("alphavantage", opts)

// Or set environment variables (source-specific examples use this approach)
// export FRED_API_KEY=your_key
// export ALPHAVANTAGE_API_KEY=your_key
// export IEX_API_KEY=your_key
```

### Getting API Keys

- **FRED**: Free at https://fred.stlouisfed.org/docs/api/api_key.html
- **Alpha Vantage**: Free tier at https://www.alphavantage.co/support/#api-key
- **IEX Cloud**: Free tier at https://iexcloud.io/pricing/

## Advanced Usage

### Custom Configuration

```go
opts := &datareader.Options{
    // API authentication
    APIKey: "your-api-key",

    // HTTP client settings
    Timeout:    60 * time.Second,
    UserAgent:  "MyApp/1.0",

    // Retry configuration
    MaxRetries: 3,
    RetryDelay: 2 * time.Second,

    // Rate limiting (requests per second)
    RateLimit: 5.0,

    // Response caching
    CacheDir: ".cache/datareader",
    CacheTTL: 24 * time.Hour,
}

reader, err := datareader.DataReader("yahoo", opts)
```

### Using the Factory Pattern

```go
// Create a reader instance for reuse
reader, err := datareader.DataReader("yahoo", nil)
if err != nil {
    log.Fatal(err)
}

// Fetch single symbol
data, err := reader.ReadSingle(ctx, "AAPL", start, end)

// Fetch multiple symbols
data, err := reader.Read(ctx, []string{"AAPL", "MSFT", "GOOGL"}, start, end)
```

### Context and Cancellation

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)

// Manual cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel after 5 seconds
}()

data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
```

## Examples

See the [examples](./examples/) directory for complete working examples:

- **[Basic Usage](./examples/basic/)** - Simple example showing data fetching
- **[Advanced Usage](./examples/advanced/)** - Custom options, caching, rate limiting
- **[FRED](./examples/fred/)** - Federal Reserve Economic Data usage
- **[World Bank](./examples/worldbank/)** - International economic indicators
- **[Alpha Vantage](./examples/alphavantage/)** - Stock market data with API key
- **[Stooq](./examples/stooq/)** - Free international stock market data
- **[IEX Cloud](./examples/iex/)** - Professional stock market data

Run an example:
```bash
cd examples/basic
go run main.go

# For examples requiring API keys
cd examples/fred
FRED_API_KEY=your_key_here go run main.go
```

## Documentation

- **[API Reference](https://pkg.go.dev/github.com/julianshen/gonp-datareader)** - Full API documentation
- **[Examples](./examples/)** - Working code examples for all sources
- **[Development Guide](./CLAUDE.md)** - TDD methodology and development guidelines
- **[Implementation Plan](./plan.md)** - Detailed implementation roadmap

## Development Status

**Current Status: Production-Ready**

All core phases complete:
- ✅ **Core Framework**: Error handling, HTTP client, retries
- ✅ **Yahoo Finance**: Free stock market OHLCV data
- ✅ **FRED**: Federal Reserve Economic Data (with optional API key)
- ✅ **World Bank**: International economic indicators
- ✅ **Alpha Vantage**: Stock market data (requires API key)
- ✅ **Stooq**: Free international stock market data
- ✅ **IEX Cloud**: Professional stock market data (requires API token)
- ✅ **Rate Limiting**: Token bucket algorithm for API limits
- ✅ **Response Caching**: File-based caching with TTL
- ✅ **Comprehensive Tests**: >70% test coverage
- ✅ **Full Documentation**: Package docs and examples

## Testing

Run all tests:
```bash
make test

# With coverage
make test-coverage

# Run linters
make lint

# Format code
make fmt
```

## Contributing

Contributions are welcome! Please see [CLAUDE.md](./CLAUDE.md) for development guidelines.

This project follows Test-Driven Development (TDD) methodology. All contributions should:
- Include tests with >80% coverage
- Follow Go conventions and idioms
- Pass all linters (`go vet`, `golangci-lint`)
- Include documentation for exported functions
- Use the Red → Green → Refactor cycle

## Roadmap

Future enhancements:
- Additional data sources (Tiingo, OECD, Eurostat)
- Improved performance optimizations
- Enhanced documentation
- More comprehensive examples

See [plan.md](./plan.md) for the detailed implementation roadmap.

## License

MIT License - see [LICENSE](./LICENSE) file for details.

## Acknowledgments

Inspired by Python's [pandas-datareader](https://pandas-datareader.readthedocs.io/).
