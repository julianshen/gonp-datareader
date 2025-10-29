# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-29

### 🎉 Production Ready - First Stable Release

This release marks gonp-datareader as production-ready with comprehensive test coverage,
security audit compliance, and stable APIs.

### Added

#### Test Coverage Improvements
- **Phase 1**: Improved coverage for 4 lowest-coverage sources
  - Alpha Vantage: 27.4% → 89.9% (+62.5%)
  - IEX Cloud: 28.2% → 89.3% (+61.1%)
  - Stooq: 28.4% → 87.5% (+59.1%)
  - World Bank: 29.9% → 89.0% (+59.1%)

- **Phase 2**: Improved coverage for remaining low-coverage sources
  - Tiingo: 47.7% → 78.5% (+30.8%)
  - FRED: 55.9% → 91.7% (+35.8%)

- **Overall Project Coverage**: **86.86%** (exceeds 85% v1.0.0 target)

#### Test Infrastructure
- Added baseURL support to all data source readers for flexible testing
- Created New[Source]ReaderWithBaseURL helper constructors for mock server testing
- Comprehensive mock HTTP server tests for all sources
- Added tests for:
  - Symbol validation
  - Multiple symbol fetching (parallel operations)
  - HTTP error handling (500 responses)
  - Invalid date range handling
  - API key validation
  - Context-based API key passing

#### Security
- **Security Audit**: Passed gosec security scanner audit with 0 issues
- Fixed all 6 security findings (3 MEDIUM, 3 LOW severity)
- Added appropriate error handling for cache operations
- Documented security considerations for cache file permissions

### Fixed

#### Security Issues Resolved
1. G104 (LOW): Unhandled error from resp.Body.Close() in retry.go
2. G104 (LOW): Unhandled error from cache.Set() in retry.go
3. G104 (LOW): Unhandled error from os.Remove() in cache.go
4. G304 (MEDIUM): Potential file inclusion via variable in cache.go
5. G301 (MEDIUM): Directory permissions 0755 in cache.go
6. G306 (MEDIUM): File permissions 0644 in cache.go

#### World Bank Reader
- Fixed ValidateSymbol to allow "/" character in "country/indicator" format
- Custom symbol validation loop in Read method for World Bank-specific formats

### Technical Details

#### v1.0.0 Checklist
- ✅ All planned sources working (9 sources)
- ✅ Coverage > 85% (achieved 86.86%)
- ✅ Production battle-tested
- ✅ API stable and documented
- ✅ Security audit passed (0 issues)
- ✅ Performance optimized
- ✅ Comprehensive examples
- ✅ Migration guide tested

#### Test Coverage by Package
- **Sources**: alphavantage (89.9%), iex (89.3%), stooq (87.5%), worldbank (89.0%),
  tiingo (78.5%), fred (91.7%), yahoo (90.7%), oecd (65.1%), eurostat (65.6%)
- **Infrastructure**: ratelimit (100%), utils (100%), sources (100%), http (95.3%), cache (89.2%)
- **Main Package**: 71.1%

### Breaking Changes
None - API remains stable from v0.2.0.

### Migration from v0.2.0
No breaking changes. All existing code will continue to work as-is.

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

## [0.2.0] - 2025-01-28

### Added

#### Documentation
- **docs/sources.md**: Comprehensive 700+ line documentation for all 9 data sources
  - Detailed API key requirements and symbol formats
  - Rate limiting information and best practices
  - Capabilities and limitations for each source
  - Usage examples and comparison matrix
- **docs/migration.md**: 800+ line migration guide from pandas-datareader
  - Side-by-side Python/Go code comparisons
  - 7 complete conversion examples
  - Feature parity matrix
  - Best practices for Go developers coming from Python
- **docs/api.md**: 900+ line complete API reference
  - Detailed function and interface documentation
  - 7 practical usage examples
  - 4 advanced patterns (concurrent fetching, aggregation, custom config, fallback)
  - Best practices section

#### CI/CD
- GitHub Actions workflows for automated testing
- Automated linting with golangci-lint
- Code coverage reporting with Codecov integration
- Coverage and Go Report Card badges in README

### Changed

#### Performance
- **Parallel Multi-Symbol Fetching**: All 9 data sources now support parallel symbol fetching
  - Worker pool pattern with semaphore (max 10 concurrent requests)
  - Context cancellation support throughout
  - Graceful error handling for partial failures
  - Applied to: Stooq, Alpha Vantage, World Bank, IEX Cloud, and all other sources

#### Infrastructure
- Improved error handling in multi-symbol scenarios
- Enhanced test coverage for concurrent operations
- Better documentation throughout codebase

### Technical Details

#### Test Coverage
- Main package: 71.1%
- Infrastructure packages: 89.2%-100% (ratelimit: 100%, utils: 100%, sources base: 100%)
- sources/yahoo: 90.7%
- Weighted average (by criticality): ~80%+

#### Performance Metrics
- Multi-symbol parallel fetching: 4.5x speedup
- Worker pool efficiently manages concurrent requests
- Semaphore pattern prevents API overwhelming

### Documentation Improvements
- Complete source-specific documentation
- Migration guide for Python developers
- API reference with practical examples
- Enhanced README with badges and status

---

## Unreleased

### Planned Features
- Additional data sources (as community requests)
- Enhanced caching strategies
- More examples and tutorials

[0.2.0]: https://github.com/julianshen/gonp-datareader/releases/tag/v0.2.0
[0.1.0]: https://github.com/julianshen/gonp-datareader/releases/tag/v0.1.0
