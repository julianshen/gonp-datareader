# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-01-26

### Added

#### Core Framework
- HTTP client with automatic retries and exponential backoff
- Token bucket rate limiting for API compliance
- File-based response caching with TTL support
- Context support for cancellation and timeouts
- Comprehensive error handling with wrapped errors
- Parallel symbol fetching with worker pools (4.5x speedup)

#### Data Sources (9 Total)
- **Yahoo Finance**: Free stock market OHLCV data
- **FRED**: Federal Reserve Economic Data with optional API key
- **World Bank**: International development indicators
- **Alpha Vantage**: Real-time & historical stock data (requires API key)
- **Stooq**: Free international stock market data
- **IEX Cloud**: Professional stock market data (requires API token)
- **Tiingo**: High-quality stock market data & fundamentals (requires API token)
- **OECD**: Economic indicators and statistics (SDMX-JSON format)
- **Eurostat**: European Union statistics (JSON-stat format)

#### Features
- Unified API across all data sources
- Factory pattern for reader creation
- Type-safe interfaces
- Automatic retries with exponential backoff
- Configurable rate limiting (token bucket algorithm)
- Response caching with file-based storage
- Parallel fetching for multiple symbols
- Context cancellation support
- CSV parsing for Yahoo Finance and Stooq
- JSON parsing for FRED, World Bank, Alpha Vantage, IEX, Tiingo
- SDMX-JSON parsing for OECD
- JSON-stat parsing for Eurostat

#### Performance Optimizations
- CSV parser optimization: 10% speedup
- Buffer pooling: 140x faster allocations (0 allocations/op)
- Parallel symbol fetching: 4.5x speedup for multiple symbols
- Worker pool pattern with semaphore (max 10 concurrent requests)

#### Testing
- Comprehensive test suite with >70% coverage
- Main package: 71.1% coverage
- Infrastructure packages: 89.2%-100% coverage
- Mock server testing for all data sources
- Benchmark tests for parsers and critical paths
- Edge case testing
- Concurrent request testing

#### Documentation
- Complete godoc documentation for all exported types and functions
- README with quick start guide
- API key configuration instructions
- 10 working examples (one for each source + basic + advanced)
- Development guidelines (CLAUDE.md)
- Implementation roadmap (plan.md)
- Performance documentation (PERFORMANCE.md)

#### Examples
- Basic usage example
- Advanced usage with caching and rate limiting
- FRED example
- World Bank example
- Alpha Vantage example
- Stooq example
- IEX Cloud example
- Tiingo example
- OECD example
- Eurostat example

### Technical Details

#### Supported Go Versions
- Go 1.21 or later

#### Dependencies
- Zero external dependencies for core functionality
- All dependencies are standard library

#### API Compatibility
- Stable API for v0.1.x series
- Breaking changes will increment major version

### Performance Metrics
- CSV parsing: 641K ops/sec, 1,902 ns/op
- Buffer pool: 168M ops/sec, 7.2 ns/op, 0 allocs/op
- Parallel fetching: 4.5x speedup (5 symbols @ 10ms each: 11ms vs 50ms)

### Test Coverage
- Main package: 71.1%
- internal/cache: 89.2%
- internal/http: 95.3%
- internal/ratelimit: 100.0%
- internal/utils: 100.0%
- sources: 100.0%
- sources/yahoo: 90.7%
- Average: >75%

### License
- MIT License

### Contributors
- Julian Shen (@julianshen)
- Claude (AI pair programmer)

### Acknowledgments
- Inspired by Python's [pandas-datareader](https://pandas-datareader.readthedocs.io/)
- Built with Test-Driven Development (TDD) methodology
- Follows Go best practices and idioms

---

## Unreleased

### Planned Features
- Additional data sources (as community requests)
- Enhanced caching strategies
- More examples and tutorials

[0.1.0]: https://github.com/julianshen/gonp-datareader/releases/tag/v0.1.0
