// Package eurostat provides data access to Eurostat API.
package eurostat

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/internal/utils"
	"github.com/julianshen/gonp-datareader/sources"
)

const (
	// eurostatAPIURL is the base URL for Eurostat Statistics API (JSON-stat format)
	eurostatAPIURL = "https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data/%s"
)

// EurostatReader fetches data from Eurostat API.
type EurostatReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewEurostatReader creates a new Eurostat data reader.
func NewEurostatReader(opts *internalhttp.ClientOptions) *EurostatReader {
	return NewEurostatReaderWithBaseURL(opts, eurostatAPIURL)
}

// NewEurostatReaderWithBaseURL creates a new Eurostat reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewEurostatReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *EurostatReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &EurostatReader{
		BaseSource: sources.NewBaseSource("eurostat"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name of the data source.
func (e *EurostatReader) Name() string {
	return "Eurostat"
}

// ValidateSymbol validates a Eurostat dataset code.
// Eurostat symbols are dataset codes like "DEMO_R_D3DENS", "GDP", etc.
func (e *EurostatReader) ValidateSymbol(symbol string) error {
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}

	// Check for invalid characters (spaces)
	if strings.Contains(symbol, " ") {
		return fmt.Errorf("symbol cannot contain spaces")
	}

	return nil
}

// BuildURL constructs the Eurostat API URL for the given symbol and date range.
func (e *EurostatReader) BuildURL(symbol string, start, end time.Time) string {
	// Build URL with dataset code
	url := fmt.Sprintf(e.baseURL, symbol)

	// Add language parameter (default to English)
	url += "?lang=EN"

	// Note: Eurostat API doesn't support date filtering in the URL
	// Date filtering would need to be done post-fetch or via dimension filters
	// For now, we fetch all data and filter client-side if needed

	return url
}

// ReadSingle fetches data for a single symbol from Eurostat.
func (e *EurostatReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := e.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Build URL
	url := e.BuildURL(symbol, start, end)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Accept header for JSON
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Eurostat returned status %d (failed to read response body: %w)", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("Eurostat returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	data, err := ParseJSON(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from Eurostat.
// Symbols are fetched in parallel for better performance.
func (e *EurostatReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return e.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (e *EurostatReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := e.ReadSingle(ctx, sym, start, end)

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
