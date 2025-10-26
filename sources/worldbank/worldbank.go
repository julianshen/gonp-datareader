// Package worldbank provides a World Bank data source reader.
package worldbank

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

// WorldBankReader fetches data from the World Bank API.
type WorldBankReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
}

// NewWorldBankReader creates a new World Bank data reader.
func NewWorldBankReader(opts *internalhttp.ClientOptions) *WorldBankReader {
	return &WorldBankReader{
		BaseSource: sources.NewBaseSource("worldbank"),
		client:     internalhttp.NewRetryableClient(opts),
	}
}

// BuildURL constructs the World Bank API URL for fetching indicator data.
// The World Bank API format is:
// https://api.worldbank.org/v2/country/{countries}/indicator/{indicator}?date={start}:{end}&format=json
func BuildURL(country, indicator string, start, end time.Time) string {
	startYear := start.Year()
	endYear := end.Year()

	return fmt.Sprintf(
		"https://api.worldbank.org/v2/country/%s/indicator/%s?date=%d:%d&format=json&per_page=1000",
		country,
		indicator,
		startYear,
		endYear,
	)
}

// ReadSingle fetches data for a single indicator and country.
// The symbol parameter should be in the format "country/indicator", e.g., "USA/NY.GDP.MKTP.CD"
func (w *WorldBankReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate symbol
	if err := w.ValidateSymbol(symbol); err != nil {
		return nil, err
	}

	// Parse symbol into country and indicator
	// For World Bank, symbol format is "country/indicator"
	// Example: "USA/NY.GDP.MKTP.CD"
	parts := splitSymbol(symbol)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid symbol format: expected 'country/indicator', got %q", symbol)
	}

	country := parts[0]
	indicator := parts[1]

	// Build URL
	url := BuildURL(country, indicator, start, end)

	// Create HTTP request
	req, err := newRequest(ctx, "GET", url)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := readAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse response
	data, err := ParseResponse(body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple indicators and countries from World Bank.
// Symbols are fetched in parallel for better performance.
func (w *WorldBankReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return w.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple indicators in parallel using a worker pool.
func (w *WorldBankReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := w.ReadSingle(ctx, sym, start, end)

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

// ValidateSymbol checks if a symbol is valid for World Bank.
func (w *WorldBankReader) ValidateSymbol(symbol string) error {
	return w.BaseSource.ValidateSymbol(symbol)
}

// splitSymbol splits a World Bank symbol into country and indicator.
// Expected format: "country/indicator" or "country;country2/indicator"
func splitSymbol(symbol string) []string {
	return strings.Split(symbol, "/")
}

// newRequest creates a new HTTP request with context.
func newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, nil)
}

// readAll reads all data from a reader.
func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
