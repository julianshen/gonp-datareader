// Package yahoo provides data access to Yahoo Finance.
package yahoo

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
	// yahooAPIURL is the base URL for Yahoo Finance historical data API
	yahooAPIURL = "https://query1.finance.yahoo.com/v7/finance/download/%s"
)

// YahooReader fetches data from Yahoo Finance.
type YahooReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewYahooReader creates a new Yahoo Finance data reader.
func NewYahooReader(opts *internalhttp.ClientOptions) *YahooReader {
	return NewYahooReaderWithBaseURL(opts, yahooAPIURL)
}

// NewYahooReaderWithBaseURL creates a new Yahoo Finance reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewYahooReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *YahooReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &YahooReader{
		BaseSource: sources.NewBaseSource("yahoo"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name of the data source.
func (y *YahooReader) Name() string {
	return "Yahoo Finance"
}

// BuildURL constructs the Yahoo Finance API URL for the given symbol and date range.
func (y *YahooReader) BuildURL(symbol string, start, end time.Time) string {
	baseURL := fmt.Sprintf(y.baseURL, symbol)

	// Convert dates to Unix timestamps
	period1 := start.Unix()
	period2 := end.Unix()

	// Build query parameters
	url := fmt.Sprintf("%s?period1=%d&period2=%d&interval=1d&events=history&includeAdjustedClose=true",
		baseURL, period1, period2)

	return url
}

// ReadSingle fetches data for a single symbol from Yahoo Finance.
func (y *YahooReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := y.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Build URL
	url := y.BuildURL(symbol, start, end)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := y.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yahoo finance returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse CSV response
	data, err := ParseCSV(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from Yahoo Finance.
// Symbols are fetched in parallel for better performance.
func (y *YahooReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return y.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (y *YahooReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := y.ReadSingle(ctx, sym, start, end)

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
