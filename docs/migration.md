# Migration Guide: pandas-datareader to gonp-datareader

This guide helps Python users familiar with pandas-datareader transition to gonp-datareader in Go.

---

## Table of Contents

1. [Quick Comparison](#quick-comparison)
2. [Installation](#installation)
3. [Basic Concepts](#basic-concepts)
4. [Syntax Differences](#syntax-differences)
5. [Feature Parity](#feature-parity)
6. [Code Examples](#code-examples)
7. [Common Patterns](#common-patterns)
8. [Performance Considerations](#performance-considerations)

---

## Quick Comparison

### pandas-datareader (Python)
```python
import pandas_datareader as pdr
from datetime import datetime

start = datetime(2024, 1, 1)
end = datetime.now()

df = pdr.get_data_yahoo('AAPL', start, end)
```

### gonp-datareader (Go)
```go
import (
    "context"
    "time"
    "github.com/julianshen/gonp-datareader"
)

ctx := context.Background()
start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Now()

data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
```

---

## Installation

### pandas-datareader (Python)
```bash
pip install pandas-datareader
```

### gonp-datareader (Go)
```bash
go get github.com/julianshen/gonp-datareader
```

---

## Basic Concepts

### Python → Go Equivalents

| Concept | Python (pandas-datareader) | Go (gonp-datareader) |
|---------|----------------------------|----------------------|
| **Data Source** | Module function (`pdr.get_data_yahoo()`) | Source parameter (`"yahoo"`) |
| **Date Range** | `datetime` objects | `time.Time` objects |
| **Error Handling** | Exceptions (try/except) | Explicit error returns |
| **Configuration** | Keyword arguments | `Options` struct |
| **Result Type** | `pandas.DataFrame` | Source-specific types |
| **Concurrency** | Threading/asyncio | Goroutines (built-in) |

---

## Syntax Differences

### 1. Imports

**Python:**
```python
import pandas_datareader as pdr
from datetime import datetime
```

**Go:**
```go
import (
    "context"
    "time"
    "github.com/julianshen/gonp-datareader"
)
```

### 2. Date Handling

**Python:**
```python
from datetime import datetime

start = datetime(2024, 1, 1)
end = datetime.now()
```

**Go:**
```go
import "time"

start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Now()
```

### 3. Basic Data Fetching

**Python:**
```python
# Yahoo Finance
df = pdr.get_data_yahoo('AAPL', start, end)

# FRED
df = pdr.get_data_fred('GDP', start, end)

# World Bank
df = pdr.wb.download(indicator='NY.GDP.MKTP.CD', country='US', start=start, end=end)
```

**Go:**
```go
ctx := context.Background()

// Yahoo Finance
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)

// FRED
data, err := datareader.Read(ctx, "GDP", "fred", start, end, nil)

// World Bank
data, err := datareader.Read(ctx, "USA/NY.GDP.MKTP.CD", "worldbank", start, end, nil)
```

### 4. Error Handling

**Python:**
```python
try:
    df = pdr.get_data_yahoo('AAPL', start, end)
except Exception as e:
    print(f"Error: {e}")
```

**Go:**
```go
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
if err != nil {
    log.Fatal(fmt.Errorf("failed to fetch data: %w", err))
}
```

### 5. Configuration Options

**Python:**
```python
# Implicit configuration through environment or defaults
df = pdr.get_data_alphavantage('AAPL', start, end, api_key='YOUR_KEY')
```

**Go:**
```go
opts := &datareader.Options{
    APIKey:     "YOUR_KEY",
    Timeout:    60 * time.Second,
    MaxRetries: 3,
    CacheDir:   ".cache",
    CacheTTL:   24 * time.Hour,
}

data, err := datareader.Read(ctx, "AAPL", "alphavantage", start, end, opts)
```

### 6. Multiple Symbols

**Python:**
```python
# Fetch multiple symbols
symbols = ['AAPL', 'MSFT', 'GOOGL']
df = pdr.get_data_yahoo(symbols, start, end)
```

**Go:**
```go
// Fetch multiple symbols (parallel)
symbols := []string{"AAPL", "MSFT", "GOOGL"}
dataMap, err := datareader.Read(ctx, symbols, "yahoo", start, end, nil)
if err != nil {
    log.Fatal(err)
}

// Access individual symbol data
appleData := dataMap.(map[string]*yahoo.ParsedData)["AAPL"]
```

---

## Feature Parity

### Data Sources

| Source | pandas-datareader | gonp-datareader | Notes |
|--------|-------------------|-----------------|-------|
| Yahoo Finance | ✅ | ✅ | Full parity |
| FRED | ✅ | ✅ | Full parity |
| World Bank | ✅ | ✅ | Symbol format differs |
| Alpha Vantage | ✅ | ✅ | Full parity |
| Stooq | ✅ | ✅ | Full parity |
| IEX Cloud | ✅ | ✅ | Full parity |
| Tiingo | ✅ | ✅ | Full parity |
| OECD | ✅ | ✅ | SDMX-JSON format |
| Eurostat | ✅ | ✅ | JSON-stat format |
| Quandl | ✅ | ❌ | Not implemented |
| Nasdaq | ✅ | ❌ | Not implemented |
| Morningstar | ✅ | ❌ | Not implemented |

### Features

| Feature | pandas-datareader | gonp-datareader |
|---------|-------------------|-----------------|
| **Historical Data** | ✅ | ✅ |
| **Multiple Symbols** | ✅ | ✅ (parallel) |
| **Date Range** | ✅ | ✅ |
| **API Key Support** | ✅ | ✅ |
| **Caching** | ✅ (requests-cache) | ✅ (built-in) |
| **Rate Limiting** | ❌ | ✅ (built-in) |
| **Retry Logic** | ❌ | ✅ (built-in) |
| **Timeout Control** | ✅ | ✅ |
| **Concurrency** | Threading | Goroutines (native) |
| **Type Safety** | Duck typing | Static typing |

---

## Code Examples

### Example 1: Basic Stock Data

**Python:**
```python
import pandas_datareader as pdr
from datetime import datetime

start = datetime(2024, 1, 1)
end = datetime.now()

# Fetch Apple stock data
df = pdr.get_data_yahoo('AAPL', start, end)
print(df.head())
print(f"Closing prices: {df['Close']}")
```

**Go:**
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

    // Fetch Apple stock data
    result, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
    if err != nil {
        log.Fatal(err)
    }

    data := result.(*yahoo.ParsedData)
    fmt.Printf("Fetched %d rows\n", len(data.Rows))

    // Access closing prices
    for _, row := range data.Rows {
        fmt.Printf("Date: %s, Close: %s\n", row["Date"], row["Close"])
    }
}
```

### Example 2: Economic Data from FRED

**Python:**
```python
import pandas_datareader as pdr
from datetime import datetime

start = datetime(2020, 1, 1)
end = datetime.now()

# US GDP
gdp = pdr.get_data_fred('GDP', start, end)
print(gdp)
```

**Go:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/fred"
)

func main() {
    ctx := context.Background()
    start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    // US GDP
    result, err := datareader.Read(ctx, "GDP", "fred", start, end, nil)
    if err != nil {
        log.Fatal(err)
    }

    data := result.(*fred.ParsedData)
    fmt.Printf("GDP data points: %d\n", len(data.Dates))

    // Print data
    for i, date := range data.Dates {
        fmt.Printf("Date: %s, GDP: %.2f\n", date, data.Values[i])
    }
}
```

### Example 3: World Bank Indicators

**Python:**
```python
import pandas_datareader as pdr
from datetime import datetime

start = datetime(2010, 1, 1)
end = datetime(2020, 12, 31)

# Download GDP for US and China
df = pdr.wb.download(
    indicator='NY.GDP.MKTP.CD',
    country=['US', 'CN'],
    start=start,
    end=end
)
print(df)
```

**Go:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/julianshen/gonp-datareader"
    "github.com/julianshen/gonp-datareader/sources/worldbank"
)

func main() {
    ctx := context.Background()
    start := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

    // US GDP
    result, err := datareader.Read(ctx, "USA/NY.GDP.MKTP.CD", "worldbank", start, end, nil)
    if err != nil {
        log.Fatal(err)
    }

    data := result.(*worldbank.ParsedData)
    for _, obs := range data.Observations {
        fmt.Printf("%s: %s\n", obs.Date, obs.Value)
    }

    // China GDP (separate call)
    result, err = datareader.Read(ctx, "CHN/NY.GDP.MKTP.CD", "worldbank", start, end, nil)
    // Process...
}
```

### Example 4: Multiple Symbols with Caching

**Python:**
```python
import pandas_datareader as pdr
from datetime import datetime
import requests_cache

# Enable caching
session = requests_cache.CachedSession('cache')
pdr.get_data_yahoo.session = session

symbols = ['AAPL', 'MSFT', 'GOOGL']
df = pdr.get_data_yahoo(symbols, start, end)
```

**Go:**
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

    // Configure with caching
    opts := &datareader.Options{
        CacheDir: ".cache/yahoo",
        CacheTTL: 24 * time.Hour,
    }

    // Fetch multiple symbols (parallel)
    symbols := []string{"AAPL", "MSFT", "GOOGL"}
    result, err := datareader.Read(ctx, symbols, "yahoo", start, end, opts)
    if err != nil {
        log.Fatal(err)
    }

    dataMap := result.(map[string]*yahoo.ParsedData)
    for symbol, data := range dataMap {
        fmt.Printf("%s: %d rows\n", symbol, len(data.Rows))
    }
}
```

---

## Common Patterns

### 1. Fetching and Processing Data

**Python:**
```python
df = pdr.get_data_yahoo('AAPL', start, end)

# Calculate returns
df['Returns'] = df['Close'].pct_change()

# Filter by condition
high_volume = df[df['Volume'] > 1000000]
```

**Go:**
```go
result, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
if err != nil {
    log.Fatal(err)
}

data := result.(*yahoo.ParsedData)

// Process data
for _, row := range data.Rows {
    close := parseFloat(row["Close"])
    volume := parseInt(row["Volume"])

    // Custom processing
    if volume > 1000000 {
        fmt.Printf("High volume day: %s\n", row["Date"])
    }
}
```

### 2. Error Handling

**Python:**
```python
try:
    df = pdr.get_data_yahoo('INVALID', start, end)
except Exception as e:
    print(f"Error: {e}")
    df = None
```

**Go:**
```go
data, err := datareader.Read(ctx, "INVALID", "yahoo", start, end, nil)
if err != nil {
    log.Printf("Error: %v\n", err)
    return
}

// Continue with data...
```

### 3. Configuration and Reuse

**Python:**
```python
# Configuration through environment
import os
os.environ['ALPHAVANTAGE_API_KEY'] = 'YOUR_KEY'

df = pdr.get_data_alphavantage('AAPL', start, end)
```

**Go:**
```go
// Create reusable reader
opts := &datareader.Options{
    APIKey: os.Getenv("ALPHAVANTAGE_API_KEY"),
    Timeout: 60 * time.Second,
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

## Performance Considerations

### Python (pandas-datareader)
- **Single-threaded by default**
- **Threading**: Use `concurrent.futures` for parallel requests
- **Async**: Use `asyncio` with async-compatible libraries
- **Memory**: DataFrames can be memory-intensive
- **Speed**: Interpreted Python + pandas overhead

**Example (Threading):**
```python
from concurrent.futures import ThreadPoolExecutor

symbols = ['AAPL', 'MSFT', 'GOOGL']

def fetch_symbol(symbol):
    return pdr.get_data_yahoo(symbol, start, end)

with ThreadPoolExecutor(max_workers=5) as executor:
    results = list(executor.map(fetch_symbol, symbols))
```

### Go (gonp-datareader)
- **Parallel by default** for multi-symbol requests
- **Goroutines**: Lightweight concurrency built-in
- **Memory**: Efficient with static types
- **Speed**: Compiled, faster execution
- **Worker pools**: Built-in concurrency limiting

**Example (Built-in Parallelism):**
```go
// Automatic parallel fetching
symbols := []string{"AAPL", "MSFT", "GOOGL"}
dataMap, err := datareader.Read(ctx, symbols, "yahoo", start, end, nil)
// 4-5x faster than sequential
```

### Performance Comparison

| Aspect | Python | Go |
|--------|--------|-----|
| **Startup** | Slower (interpreter) | Faster (compiled) |
| **Execution** | Slower | Faster (2-10x) |
| **Concurrency** | GIL limitations | True parallelism |
| **Memory** | Higher | Lower |
| **Type Safety** | Runtime | Compile-time |

---

## Key Differences Summary

### What's Different

1. **Type System**: Go is statically typed, Python is dynamically typed
2. **Error Handling**: Go uses explicit error returns, Python uses exceptions
3. **Concurrency**: Go has goroutines (native), Python has threading/asyncio
4. **Source Selection**: Go uses string parameter, Python uses different functions
5. **Result Types**: Go returns source-specific types, Python returns DataFrames
6. **Configuration**: Go uses Options struct, Python uses keyword arguments

### What's Similar

1. **Data Sources**: Both support major sources (Yahoo, FRED, etc.)
2. **Date Handling**: Both use standard library date types
3. **API Keys**: Both support authenticated sources
4. **Caching**: Both support response caching
5. **Purpose**: Both fetch financial/economic data

### Advantages of Go

- ✅ **Performance**: Faster execution (compiled)
- ✅ **Concurrency**: Native goroutines
- ✅ **Type Safety**: Compile-time error checking
- ✅ **Built-in Features**: Rate limiting, retries included
- ✅ **Deployment**: Single binary, no dependencies
- ✅ **Memory**: More efficient

### Advantages of Python

- ✅ **Ecosystem**: Rich data science libraries (pandas, numpy)
- ✅ **Simplicity**: More concise syntax
- ✅ **Interactive**: REPL, Jupyter notebooks
- ✅ **Flexibility**: Dynamic typing
- ✅ **Community**: Larger data science community

---

## Complete Example: Side-by-Side

### Python Version

```python
import pandas_datareader as pdr
from datetime import datetime
import requests_cache

# Setup
session = requests_cache.CachedSession('cache', expire_after=86400)
start = datetime(2024, 1, 1)
end = datetime.now()

# Fetch data
try:
    symbols = ['AAPL', 'MSFT', 'GOOGL']
    df = pdr.get_data_yahoo(symbols, start, end, session=session)

    # Process
    print(f"Fetched {len(df)} rows")
    print(df['Close'].head())

except Exception as e:
    print(f"Error: {e}")
```

### Go Version

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
    // Setup
    ctx := context.Background()
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Now()

    opts := &datareader.Options{
        CacheDir: ".cache/yahoo",
        CacheTTL: 24 * time.Hour,
    }

    // Fetch data
    symbols := []string{"AAPL", "MSFT", "GOOGL"}
    result, err := datareader.Read(ctx, symbols, "yahoo", start, end, opts)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // Process
    dataMap := result.(map[string]*yahoo.ParsedData)
    for symbol, data := range dataMap {
        fmt.Printf("%s: %d rows\n", symbol, len(data.Rows))
        if len(data.Rows) > 0 {
            fmt.Printf("Latest close: %s\n", data.Rows[len(data.Rows)-1]["Close"])
        }
    }
}
```

---

## Migration Checklist

- [ ] Install Go (1.21+)
- [ ] Set up Go project (`go mod init`)
- [ ] Install gonp-datareader (`go get github.com/julianshen/gonp-datareader`)
- [ ] Convert date handling from `datetime` to `time.Time`
- [ ] Replace function calls with `datareader.Read()`
- [ ] Add explicit error handling
- [ ] Update symbol formats (especially World Bank)
- [ ] Configure options (API keys, caching, etc.)
- [ ] Update result type handling
- [ ] Test with your data sources
- [ ] Optimize with parallel fetching
- [ ] Add retry/rate limiting if needed (built-in)

---

## Getting Help

- **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/julianshen/gonp-datareader)
- **Examples**: [GitHub examples](https://github.com/julianshen/gonp-datareader/tree/main/examples)
- **Source Documentation**: [docs/sources.md](sources.md)
- **Issues**: [GitHub Issues](https://github.com/julianshen/gonp-datareader/issues)

---

## Conclusion

While Python's pandas-datareader and Go's gonp-datareader serve the same purpose, they reflect the different philosophies and strengths of their respective languages. Python excels at rapid prototyping and interactive data analysis, while Go provides performance, type safety, and built-in concurrency.

Choose Python if you need:
- Interactive analysis (Jupyter)
- Rich pandas ecosystem
- Rapid prototyping

Choose Go if you need:
- Production services
- High performance
- Type safety
- Built-in concurrency

Both are excellent tools for their respective ecosystems!
