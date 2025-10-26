// Package tiingo provides data access to Tiingo API.
package tiingo

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
	// tiingoAPIURL is the base URL for Tiingo API
	tiingoAPIURL = "https://api.tiingo.com/tiingo/daily/%s/prices"
)

// TiingoReader fetches data from Tiingo API.
type TiingoReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
	apiKey  string
}

// NewTiingoReader creates a new Tiingo data reader.
func NewTiingoReader(opts *internalhttp.ClientOptions) *TiingoReader {
	return NewTiingoReaderWithBaseURL(opts, tiingoAPIURL)
}

// NewTiingoReaderWithBaseURL creates a new Tiingo reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewTiingoReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *TiingoReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &TiingoReader{
		BaseSource: sources.NewBaseSource("tiingo"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
		apiKey:     "", // Will be set from context or options
	}
}

// Name returns the display name of the data source.
func (t *TiingoReader) Name() string {
	return "Tiingo"
}

// BuildURL constructs the Tiingo API URL for the given symbol and date range.
func (t *TiingoReader) BuildURL(symbol string, start, end time.Time, apiKey string) string {
	baseURL := fmt.Sprintf(t.baseURL, symbol)

	// Format dates as YYYY-MM-DD
	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	// Build query parameters
	url := fmt.Sprintf("%s?startDate=%s&endDate=%s&token=%s",
		baseURL, startDate, endDate, apiKey)

	return url
}

// ReadSingle fetches data for a single symbol from Tiingo.
func (t *TiingoReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := t.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Get API key from context or error
	apiKey := t.getAPIKey(ctx)
	if apiKey == "" {
		return nil, fmt.Errorf("Tiingo API key is required")
	}

	// Build URL
	url := t.BuildURL(symbol, start, end, apiKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tiingo returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	data, err := ParseJSON(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from Tiingo.
// Symbols are fetched in parallel for better performance.
func (t *TiingoReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
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
func (t *TiingoReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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

// getAPIKey retrieves the API key from context or the reader's stored key.
func (t *TiingoReader) getAPIKey(ctx context.Context) string {
	// Try to get from context first
	if key := ctx.Value("apiKey"); key != nil {
		if apiKey, ok := key.(string); ok && apiKey != "" {
			return apiKey
		}
	}

	// Fall back to stored API key
	return t.apiKey
}

// SetAPIKey sets the API key for the reader.
func (t *TiingoReader) SetAPIKey(apiKey string) {
	t.apiKey = apiKey
}
