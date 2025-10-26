// Package iex provides an IEX Cloud data source reader.
package iex

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

// IEXReader fetches data from IEX Cloud API.
type IEXReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
	apiKey string
}

// NewIEXReader creates a new IEX Cloud data reader.
// An API token is required to use the IEX Cloud API.
func NewIEXReader(opts *internalhttp.ClientOptions, apiKey string) *IEXReader {
	return &IEXReader{
		BaseSource: sources.NewBaseSource("iex"),
		client:     internalhttp.NewRetryableClient(opts),
		apiKey:     apiKey,
	}
}

// BuildURL constructs the IEX Cloud API URL for fetching historical chart data.
// The IEX Cloud format is:
// https://cloud.iexapis.com/stable/stock/{symbol}/chart/{range}?token={token}
func BuildURL(symbol, dateRange, apiKey string) string {
	return fmt.Sprintf(
		"https://cloud.iexapis.com/stable/stock/%s/chart/%s?token=%s",
		symbol,
		dateRange,
		apiKey,
	)
}

// CalculateDateRange converts start/end dates to IEX Cloud date range format.
// IEX Cloud uses ranges like: 1m, 3m, 6m, 1y, 2y, 5y
func CalculateDateRange(start, end time.Time) string {
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)

	switch {
	case days < 45: // Less than 1.5 months
		return "1m"
	case days < 135: // Less than 4.5 months
		return "3m"
	case days < 270: // Less than 9 months
		return "6m"
	case days < 548: // Less than 1.5 years
		return "1y"
	case days < 1095: // Less than 3 years
		return "2y"
	default:
		return "5y" // Max range
	}
}

// ReadSingle fetches data for a single stock symbol.
func (i *IEXReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	if err := i.ValidateSymbol(symbol); err != nil {
		return nil, err
	}

	if i.apiKey == "" {
		return nil, fmt.Errorf("API key is required for IEX Cloud")
	}

	// Calculate date range in IEX Cloud format
	dateRange := CalculateDateRange(start, end)

	// Build request URL
	url := BuildURL(symbol, dateRange, i.apiKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch IEX Cloud data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IEX Cloud returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse response
	data, err := ParseResponse(body)
	if err != nil {
		return nil, fmt.Errorf("parse IEX Cloud response: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple stock symbols from IEX Cloud.
// Symbols are fetched in parallel for better performance.
func (i *IEXReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return i.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (i *IEXReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := i.ReadSingle(ctx, sym, start, end)

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

// ValidateSymbol checks if a symbol is valid for IEX Cloud.
func (i *IEXReader) ValidateSymbol(symbol string) error {
	return i.BaseSource.ValidateSymbol(symbol)
}
