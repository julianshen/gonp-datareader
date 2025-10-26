// Package datareader provides remote data access for financial and economic data sources.
//
// This package offers a unified interface to fetch data from various sources including
// Yahoo Finance, FRED, World Bank, Alpha Vantage, Stooq, and IEX Cloud. It supports
// features like automatic retries, rate limiting, caching, and flexible configuration.
//
// # Quick Start
//
// The simplest way to fetch data is using the Read function:
//
//	import (
//		"context"
//		"time"
//		"github.com/julianshen/gonp-datareader"
//	)
//
//	ctx := context.Background()
//	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Now()
//
//	// Fetch stock data from Yahoo Finance
//	data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Using the Factory Pattern
//
// For more control, create a reader instance using the DataReader factory:
//
//	opts := &datareader.Options{
//		Timeout:    60 * time.Second,
//		MaxRetries: 3,
//		RetryDelay: 2 * time.Second,
//	}
//
//	reader, err := datareader.DataReader("yahoo", opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	data, err := reader.ReadSingle(ctx, "AAPL", start, end)
//
// # API Keys
//
// Some sources require API keys. Set them via the APIKey field in Options:
//
//	opts := &datareader.Options{
//		APIKey: "your-api-key-here",
//	}
//
//	// FRED with API key
//	reader, err := datareader.DataReader("fred", opts)
//
//	// Alpha Vantage requires API key
//	reader, err := datareader.DataReader("alphavantage", opts)
//
//	// IEX Cloud requires API token
//	reader, err := datareader.DataReader("iex", opts)
//
// # Supported Data Sources
//
// The following data sources are currently supported:
//
//   - yahoo: Yahoo Finance - Free stock market data
//   - fred: Federal Reserve Economic Data - Economic indicators (optional API key)
//   - worldbank: World Bank - International economic indicators
//   - alphavantage: Alpha Vantage - Stock market data (requires API key)
//   - stooq: Stooq - Free international stock market data
//   - iex: IEX Cloud - Professional stock market data (requires API token)
//   - tiingo: Tiingo - Stock market data with high-quality fundamentals (requires API token)
//   - oecd: OECD - Economic indicators and statistics (no API key required)
//
// Use ListSources() to get a list of all available sources at runtime.
//
// # Configuration Options
//
// The Options struct provides extensive configuration:
//
//	opts := &datareader.Options{
//		// API authentication
//		APIKey: "your-api-key",
//
//		// HTTP client settings
//		Timeout:    30 * time.Second,
//		UserAgent:  "MyApp/1.0",
//
//		// Retry configuration
//		MaxRetries: 3,
//		RetryDelay: 1 * time.Second,
//
//		// Rate limiting (requests per second)
//		RateLimit: 5.0,
//
//		// Response caching
//		CacheDir: ".cache/datareader",
//		CacheTTL: 24 * time.Hour,
//	}
//
// # Error Handling
//
// All functions return errors that can be checked:
//
//	data, err := datareader.Read(ctx, "INVALID", "yahoo", start, end, nil)
//	if err != nil {
//		if errors.Is(err, datareader.ErrUnknownSource) {
//			// Handle unknown source
//		}
//		// Handle other errors
//	}
package datareader

import (
	"context"
	"fmt"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
	"github.com/julianshen/gonp-datareader/sources/alphavantage"
	"github.com/julianshen/gonp-datareader/sources/fred"
	"github.com/julianshen/gonp-datareader/sources/iex"
	"github.com/julianshen/gonp-datareader/sources/oecd"
	"github.com/julianshen/gonp-datareader/sources/stooq"
	"github.com/julianshen/gonp-datareader/sources/tiingo"
	"github.com/julianshen/gonp-datareader/sources/worldbank"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

var (
	// ErrUnknownSource is returned when an unknown or unsupported data source is requested.
	// Use ListSources() to get a list of all available sources.
	ErrUnknownSource = fmt.Errorf("unknown data source")
)

// DataReader creates a new reader for the specified data source.
//
// The source parameter specifies which data source to use. Currently supported sources:
//   - "yahoo": Yahoo Finance - free stock market data (no API key required)
//   - "fred": Federal Reserve Economic Data - economic indicators (optional API key)
//   - "worldbank": World Bank - international economic indicators (no API key required)
//   - "alphavantage": Alpha Vantage - stock market data (API key required)
//   - "stooq": Stooq - free international stock market data (no API key required)
//   - "tiingo": Tiingo - stock market data (API key required)
//   - "iex": IEX Cloud - professional stock market data (API token required)
//   - "oecd": OECD - economic indicators and statistics (no API key required)
//
// The opts parameter provides configuration for the reader. If nil, default options are used.
// See the Options struct for available configuration settings.
//
// # Example Usage
//
//	// Create a Yahoo Finance reader with default options
//	reader, err := datareader.DataReader("yahoo", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create a FRED reader with custom options and API key
//	opts := &datareader.Options{
//		APIKey:     "your-fred-api-key",
//		Timeout:    60 * time.Second,
//		MaxRetries: 3,
//	}
//	reader, err := datareader.DataReader("fred", opts)
//
// # Error Handling
//
// Returns ErrUnknownSource if the source is not recognized.
// Use ListSources() to get a list of valid source names.
func DataReader(source string, opts *Options) (sources.Reader, error) {
	if source == "" {
		return nil, fmt.Errorf("%w: source cannot be empty", ErrUnknownSource)
	}

	// Convert Options to ClientOptions
	var clientOpts *internalhttp.ClientOptions
	var apiKey string
	if opts != nil {
		clientOpts = &internalhttp.ClientOptions{
			Timeout:    opts.Timeout,
			UserAgent:  opts.UserAgent,
			MaxRetries: opts.MaxRetries,
			RetryDelay: opts.RetryDelay,
			RateLimit:  opts.RateLimit,
			CacheDir:   opts.CacheDir,
			CacheTTL:   opts.CacheTTL,
		}
		apiKey = opts.APIKey
	}

	switch source {
	case "yahoo":
		return yahoo.NewYahooReader(clientOpts), nil
	case "fred":
		if apiKey != "" {
			return fred.NewFREDReaderWithAPIKey(clientOpts, apiKey), nil
		}
		return fred.NewFREDReader(clientOpts), nil
	case "worldbank":
		return worldbank.NewWorldBankReader(clientOpts), nil
	case "alphavantage":
		return alphavantage.NewAlphaVantageReader(clientOpts, apiKey), nil
	case "stooq":
		return stooq.NewStooqReader(clientOpts), nil
	case "iex":
		return iex.NewIEXReader(clientOpts, apiKey), nil
	case "tiingo":
		reader := tiingo.NewTiingoReader(clientOpts)
		if apiKey != "" {
			reader.SetAPIKey(apiKey)
		}
		return reader, nil
	case "oecd":
		return oecd.NewOECDReader(clientOpts), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSource, source)
	}
}

// Read is a convenience function that creates a reader and fetches data for a single symbol.
//
// This is the simplest way to fetch data. It combines DataReader() and ReadSingle()
// into a single call.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - symbol: The symbol to fetch (e.g., "AAPL" for stocks, "GDP" for FRED series)
//   - source: The data source name (use ListSources() to see available sources)
//   - start: Start date for the data range
//   - end: End date for the data range
//   - opts: Configuration options (can be nil for defaults)
//
// The return type depends on the source:
//   - Yahoo, Stooq, Alpha Vantage, IEX: Returns *ParsedData with OHLCV data
//   - FRED: Returns *ParsedData with dates and values
//   - World Bank: Returns *ParsedData with indicator values
//
// # Example Usage
//
//	ctx := context.Background()
//	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Now()
//
//	// Fetch stock data from Yahoo Finance
//	data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Fetch economic data from FRED
//	opts := &datareader.Options{
//		APIKey: "your-fred-api-key",
//	}
//	data, err := datareader.Read(ctx, "GDP", "fred", start, end, opts)
//
// # Context Cancellation
//
// The context can be used to cancel long-running requests:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
func Read(ctx context.Context, symbol string, source string, start, end time.Time, opts *Options) (interface{}, error) {
	reader, err := DataReader(source, opts)
	if err != nil {
		return nil, err
	}

	return reader.ReadSingle(ctx, symbol, start, end)
}

// ListSources returns a list of all available data source names.
//
// This function is useful for discovering which sources are supported
// and for validating user input.
//
// # Example Usage
//
//	sources := datareader.ListSources()
//	fmt.Println("Available sources:")
//	for _, source := range sources {
//		fmt.Printf("  - %s\n", source)
//	}
//
//	// Validate user input
//	userSource := "yahoo"
//	found := false
//	for _, source := range datareader.ListSources() {
//		if source == userSource {
//			found = true
//			break
//		}
//	}
//	if !found {
//		log.Fatalf("Unknown source: %s", userSource)
//	}
func ListSources() []string {
	return []string{
		"yahoo",
		"fred",
		"worldbank",
		"alphavantage",
		"stooq",
		"iex",
		"tiingo",
		"oecd",
	}
}
