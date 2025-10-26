# gonp-datareader

Remote data access for financial and economic data sources in Go, designed to work with [gonp](https://github.com/julianshen/gonp) DataFrames.

[![Go Reference](https://pkg.go.dev/badge/github.com/julianshen/gonp-datareader.svg)](https://pkg.go.dev/github.com/julianshen/gonp-datareader)
[![Go Report Card](https://goreportcard.com/badge/github.com/julianshen/gonp-datareader)](https://goreportcard.com/report/github.com/julianshen/gonp-datareader)

## Overview

gonp-datareader is the Go equivalent of Python's pandas-datareader, providing a simple and unified interface to fetch data from various internet sources into gonp DataFrames.

## Features

- **Multiple Data Sources**: Yahoo Finance, FRED, World Bank, Alpha Vantage, and more
- **Simple API**: Easy-to-use interface for fetching financial and economic data
- **Type Safe**: Leverages Go's type system for compile-time safety
- **High Performance**: Efficient data fetching and parsing with minimal allocations
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
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
    ctx := context.Background()
    start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

    // Method 1: Use convenience function
    data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        panic(err)
    }

    parsedData := data.(*yahoo.ParsedData)
    fmt.Printf("Fetched %d days of data\n", len(parsedData.Rows))

    // Method 2: Use factory for more control
    reader, _ := datareader.DataReader("yahoo", nil)
    results, _ := reader.Read(ctx, []string{"AAPL", "MSFT"}, start, end)

    dataMap := results.(map[string]*yahoo.ParsedData)
    for symbol, data := range dataMap {
        closes := data.GetColumn("Close")
        fmt.Printf("%s: %d prices, first close: %s\n",
            symbol, len(closes), closes[0])
    }
}
```

## Supported Data Sources

| Source | Description | API Key Required |
|--------|-------------|------------------|
| Yahoo Finance | Stock prices, dividends, splits | No |
| FRED | Federal Reserve Economic Data | Optional |
| World Bank | Development indicators | No |
| Alpha Vantage | Real-time and historical data | Yes |
| Tiingo | Financial data | Yes |
| IEX Cloud | Investors Exchange data | Yes |

## Examples

See the [examples](./examples/) directory for complete working examples:

- [Basic Usage](./examples/basic/) - Simple example showing the three main ways to fetch data
- [Advanced Usage](./examples/advanced/) - Custom options, error handling, and data analysis

Run an example:
```bash
cd examples/basic
go run main.go
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/julianshen/gonp-datareader)
- [Examples](./examples/) - Working code examples
- [Data Sources](./docs/sources.md)
- [Development Guide](./CLAUDE.md)

## Development Status

**Current Status: Production-Ready for Yahoo Finance**

Phase 1 (Yahoo Finance integration) is complete and tested:
- ✅ Core framework and interfaces
- ✅ Yahoo Finance reader with full data fetching
- ✅ Comprehensive test coverage (100%)
- ✅ Working examples (basic and advanced)
- ✅ Integration tests with mock server
- ✅ Error handling and retry logic

See [spec.md](./spec.md) for the complete specification and [plan.md](./plan.md) for the implementation roadmap.

## Contributing

Contributions are welcome! Please see [CLAUDE.md](./CLAUDE.md) for development guidelines.

This project follows Test-Driven Development (TDD) methodology. All contributions should:
- Include tests
- Follow Go conventions
- Maintain test coverage above 80%
- Pass all linters and tests

## License

MIT License - see [LICENSE](./LICENSE) file for details.

## Acknowledgments

Inspired by Python's [pandas-datareader](https://pandas-datareader.readthedocs.io/).
