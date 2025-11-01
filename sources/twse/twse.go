// Package twse provides data access to Taiwan Stock Exchange (TWSE).
//
// The TWSE reader fetches stock market data from the Taiwan Stock Exchange
// using the official TWSE Open API at https://openapi.twse.com.tw/v1/.
//
// This data source supports Taiwan stock symbols (typically 4-6 digit numeric codes)
// and provides daily trading data including OHLC prices, volume, and transaction counts.
//
// Note: TWSE uses the ROC (Republic of China) calendar system where dates are
// represented as ROC Year + Month + Day. For example, "1141031" represents
// October 31, 2025 (ROC Year 114 = Gregorian Year 2025 = 114 + 1911).
//
// Example usage:
//
//	reader := twse.NewTWSEReader(nil)
//	data, err := reader.ReadSingle(ctx, "2330", startDate, endDate)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Popular Taiwan stock symbols:
//   - 2330: TSMC (Taiwan Semiconductor Manufacturing Company)
//   - 2317: Hon Hai Precision Industry (Foxconn)
//   - 2454: MediaTek Inc.
//   - 2412: Chunghwa Telecom
//   - 0050: Yuanta/P-shares Taiwan Top 50 ETF
package twse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/internal/utils"
	"github.com/julianshen/gonp-datareader/sources"
)

const (
	// twseBaseURL is the base URL for the TWSE Open API
	twseBaseURL = "https://openapi.twse.com.tw/v1"

	// dailyStocksEndpoint provides all stocks' daily trading data
	dailyStocksEndpoint = "/exchangeReport/STOCK_DAY_ALL"

	// indexEndpoint provides market indices data
	indexEndpoint = "/exchangeReport/MI_INDEX"
)

// TWSEReader fetches data from Taiwan Stock Exchange (TWSE).
type TWSEReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewTWSEReader creates a new TWSE data reader.
//
// The reader uses default client options if opts is nil.
// No API key is required for TWSE as it's a public API.
func NewTWSEReader(opts *internalhttp.ClientOptions) *TWSEReader {
	return NewTWSEReaderWithBaseURL(opts, twseBaseURL)
}

// NewTWSEReaderWithBaseURL creates a new TWSE reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewTWSEReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *TWSEReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &TWSEReader{
		BaseSource: sources.NewBaseSource("twse"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name of the data source.
func (t *TWSEReader) Name() string {
	return "Taiwan Stock Exchange"
}

// ValidateSymbol checks if a symbol is valid for TWSE.
//
// Taiwan stock symbols are typically 4-6 digit numeric codes:
//   - Regular stocks: 4 digits (e.g., "2330" for TSMC)
//   - ETFs: 4 digits starting with 00 (e.g., "0050")
//   - Warrants: 6 digits
//
// This implementation delegates to the base symbol validation which checks
// for empty strings and invalid characters. Additional TWSE-specific
// validation will be added as needed.
func (t *TWSEReader) ValidateSymbol(symbol string) error {
	return t.BaseSource.ValidateSymbol(symbol)
}

// BuildURL constructs the TWSE API URL for fetching daily stock data.
//
// This returns the URL for the STOCK_DAY_ALL endpoint which provides
// all stocks' daily trading data for the latest trading day.
func (t *TWSEReader) BuildURL() string {
	return buildDailyURL(t.baseURL)
}

// buildDailyURL constructs the URL for the daily stocks endpoint.
//
// The endpoint returns all stocks' daily trading data including:
//   - Stock code and name
//   - OHLC prices (Open, High, Low, Close)
//   - Trading volume and transaction count
//   - Price change
//
// Example: https://openapi.twse.com.tw/v1/exchangeReport/STOCK_DAY_ALL
func buildDailyURL(baseURL string) string {
	// Remove trailing slash if present to avoid double slashes
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL + dailyStocksEndpoint
}

// buildIndexURL constructs the URL for the market indices endpoint.
//
// The endpoint returns market indices data including:
//   - Index name (in Traditional Chinese)
//   - Closing index value
//   - Change direction and points
//   - Percentage change
//
// Example: https://openapi.twse.com.tw/v1/exchangeReport/MI_INDEX
func buildIndexURL(baseURL string) string {
	// Remove trailing slash if present to avoid double slashes
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL + indexEndpoint
}

// ReadSingle fetches data for a single symbol from TWSE.
//
// Note: The TWSE API currently returns the latest trading day's data.
// The start and end parameters are validated but may not affect the returned
// data range depending on API capabilities.
func (t *TWSEReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := t.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Build URL
	urlStr := t.BuildURL()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse JSON response
	allStocks, err := parseDailyStockJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// Filter for the requested symbol
	stockData, err := filterBySymbol(allStocks, symbol)
	if err != nil {
		return nil, fmt.Errorf("filter symbol: %w", err)
	}

	// Parse the stock data into ParsedData structure
	data, err := parseStockData(stockData)
	if err != nil {
		return nil, fmt.Errorf("parse stock data: %w", err)
	}

	// Filter by date range
	filteredData := filterByDateRange(data, start, end)

	return filteredData, nil
}

// Read fetches data for multiple symbols from TWSE.
//
// Symbols are fetched in parallel for better performance.
func (t *TWSEReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return t.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (t *TWSEReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
	type result struct {
		symbol string
		data   *ParsedData
		err    error
	}

	// Create channels for work distribution and results
	results := make(chan result, len(symbols))

	// Create worker pool - limit concurrency to avoid overwhelming the server
	maxWorkers := 10
	if len(symbols) < maxWorkers {
		maxWorkers = len(symbols)
	}

	// Use a semaphore pattern to limit concurrent workers
	semaphore := make(chan struct{}, maxWorkers)

	// Launch goroutines for each symbol
	for _, symbol := range symbols {
		// Capture symbol in loop variable
		sym := symbol

		go func() {
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Fetch data
			data, err := t.ReadSingle(ctx, sym, start, end)

			// Send result
			res := result{symbol: sym, err: err}
			if err == nil {
				if parsedData, ok := data.(*ParsedData); ok {
					res.data = parsedData
				}
			}
			results <- res
		}()
	}

	// Collect results
	dataMap := make(map[string]*ParsedData, len(symbols))
	for i := 0; i < len(symbols); i++ {
		res := <-results
		if res.err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", res.symbol, res.err)
		}
		dataMap[res.symbol] = res.data
	}

	return dataMap, nil
}
