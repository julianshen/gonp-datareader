// Package finmind provides access to FinMind financial data API.
//
// FinMind (https://finmind.github.io/) offers comprehensive financial data for Taiwan
// and international markets with over 50 datasets. It provides historical data since 1994
// for Taiwan stocks, along with fundamental data, institutional investor data, and more.
//
// Key features:
//   - Optional Bearer token authentication for higher rate limits
//   - 300 requests/hour without token, 600 requests/hour with token
//   - 50+ datasets including stocks, futures, options, bonds, and commodities
//   - Historical data since 1994 for Taiwan stocks
//   - International market coverage (US stocks, commodities, currencies)
//
// Basic usage:
//
//	reader := finmind.NewFinMindReader(nil)
//	data, err := reader.ReadSingle(ctx, "2330", start, end)
//
// With authentication token:
//
//	reader := finmind.NewFinMindReaderWithToken(nil, "your-token")
//	data, err := reader.ReadSingle(ctx, "2330", start, end)
package finmind

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/internal/utils"
	"github.com/julianshen/gonp-datareader/sources"
)

const (
	// DefaultAPIEndpoint is the base URL for FinMind API v4.
	DefaultAPIEndpoint = "https://api.finmindtrade.com/api/v4/data"

	// DefaultDataset is the default dataset to fetch (Taiwan stock prices).
	DefaultDataset = "TaiwanStockPrice"

	// DefaultRateLimit is the default rate limit without token (300 requests/hour).
	DefaultRateLimit = 300.0 / 3600.0 // requests per second

	// TokenRateLimit is the rate limit with token (600 requests/hour).
	TokenRateLimit = 600.0 / 3600.0 // requests per second
)

// FinMindReader fetches financial data from FinMind API.
//
// FinMind provides comprehensive financial data for Taiwan and international markets.
// It supports optional Bearer token authentication for higher rate limits and offers
// over 50 datasets including historical stock prices, fundamental data, and institutional
// investor information.
type FinMindReader struct {
	*sources.BaseSource
	client   *internalhttp.RetryableClient
	token    string
	endpoint string
	dataset  string
}

// NewFinMindReader creates a new FinMind reader without authentication token.
//
// This reader will have a rate limit of 300 requests per hour. For higher limits
// (600 requests/hour), use NewFinMindReaderWithToken.
//
// Example:
//
//	reader := finmind.NewFinMindReader(nil)
//	data, err := reader.ReadSingle(ctx, "2330", start, end)
func NewFinMindReader(opts *internalhttp.ClientOptions) *FinMindReader {
	return NewFinMindReaderWithToken(opts, "")
}

// NewFinMindReaderWithToken creates a new FinMind reader with authentication token.
//
// Using a token increases the rate limit from 300 to 600 requests per hour.
// To get a token, register at FinMind website and retrieve it from account settings.
//
// Example:
//
//	token := "your-finmind-token"
//	reader := finmind.NewFinMindReaderWithToken(nil, token)
//	data, err := reader.ReadSingle(ctx, "2330", start, end)
func NewFinMindReaderWithToken(opts *internalhttp.ClientOptions, token string) *FinMindReader {
	return NewFinMindReaderWithTokenAndEndpoint(opts, token, DefaultAPIEndpoint)
}

// NewFinMindReaderWithEndpoint creates a new FinMind reader with custom endpoint.
// This is primarily used for testing with mock servers.
func NewFinMindReaderWithEndpoint(opts *internalhttp.ClientOptions, endpoint string) *FinMindReader {
	return NewFinMindReaderWithTokenAndEndpoint(opts, "", endpoint)
}

// NewFinMindReaderWithTokenAndEndpoint creates a new FinMind reader with both token and custom endpoint.
// This is primarily used for testing with mock servers.
func NewFinMindReaderWithTokenAndEndpoint(opts *internalhttp.ClientOptions, token, endpoint string) *FinMindReader {
	// Apply default options if not provided
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	// Set appropriate rate limit based on token presence
	if token != "" && opts.RateLimit == 0 {
		opts.RateLimit = TokenRateLimit
	} else if opts.RateLimit == 0 {
		opts.RateLimit = DefaultRateLimit
	}

	return &FinMindReader{
		BaseSource: sources.NewBaseSource("finmind"),
		client:     internalhttp.NewRetryableClient(opts),
		token:      token,
		endpoint:   endpoint,
		dataset:    DefaultDataset,
	}
}

// Name returns the display name of the data source.
func (f *FinMindReader) Name() string {
	return "FinMind"
}

// SetToken sets the authentication token for the reader.
//
// This allows updating the token after reader creation. Setting a token
// increases the rate limit from 300 to 600 requests per hour.
func (f *FinMindReader) SetToken(token string) {
	f.token = token
}

// SetDataset sets the dataset to fetch from FinMind API.
//
// FinMind supports 50+ datasets. Common datasets:
//   - TaiwanStockPrice (default): Daily stock prices
//   - TaiwanStockInfo: Company information
//   - TaiwanStockDividend: Dividend data
//   - TaiwanStockPER: P/E ratio data
//   - USStockPrice: US stock prices
//   - TaiwanFuturesDaily: Futures data
//   - TaiwanOptionsDaily: Options data
//
// Example:
//
//	reader.SetDataset("TaiwanStockDividend")
func (f *FinMindReader) SetDataset(dataset string) {
	f.dataset = dataset
}

// BuildURL constructs the API URL with query parameters for FinMind API.
//
// The URL includes the following query parameters:
//   - dataset: The dataset name (e.g., "TaiwanStockPrice")
//   - data_id: The symbol/stock code
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//
// Example output:
//
//	https://api.finmindtrade.com/api/v4/data?dataset=TaiwanStockPrice&data_id=2330&start_date=2020-04-02&end_date=2020-04-12
func (f *FinMindReader) BuildURL(symbol string, start, end time.Time) string {
	// Build query parameters
	params := url.Values{}
	params.Set("dataset", f.dataset)
	params.Set("data_id", symbol)
	params.Set("start_date", formatDate(start))
	params.Set("end_date", formatDate(end))

	// Construct full URL
	return fmt.Sprintf("%s?%s", f.endpoint, params.Encode())
}

// formatDate converts a time.Time to YYYY-MM-DD format for FinMind API.
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ReadSingle fetches data for a single symbol from FinMind.
//
// The symbol format depends on the dataset:
//   - Taiwan stocks: 4-digit code (e.g., "2330" for TSMC)
//   - Taiwan warrants: 6-digit code
//   - US stocks: Stock ticker (e.g., "AAPL")
//
// The date range is inclusive of both start and end dates.
//
// Returns ParsedData containing the fetched data with columns and rows.
// Returns an error if the symbol is invalid, the request fails, or no data is found.
func (f *FinMindReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate symbol
	if err := f.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	// Validate date range
	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Build API URL
	urlStr := f.BuildURL(symbol, start, end)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Add Authorization header if token is present
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}

	// Execute HTTP request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse JSON response
	data, err := ParseFinMindResponse(body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from FinMind in parallel.
//
// This method fetches data for all symbols concurrently with a worker pool pattern
// to respect rate limits. The maximum number of concurrent workers is 10.
//
// Returns a map of symbol to ParsedData.
// Returns an error if any symbol fails to fetch.
func (f *FinMindReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement in Phase 16.7
	return nil, nil
}
