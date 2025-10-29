// Package stooq provides a Stooq data source reader.
package stooq

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

// StooqReader fetches data from Stooq.
type StooqReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string // For testing with mock servers
}

// NewStooqReader creates a new Stooq data reader.
func NewStooqReader(opts *internalhttp.ClientOptions) *StooqReader {
	return NewStooqReaderWithBaseURL(opts, "https://stooq.com/q/d/l/?s=%s&i=d")
}

// NewStooqReaderWithBaseURL creates a new Stooq reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewStooqReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *StooqReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &StooqReader{
		BaseSource: sources.NewBaseSource("stooq"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// BuildURL constructs the Stooq URL for fetching historical data.
// The Stooq format is:
// https://stooq.com/q/d/l/?s={symbol}&i=d
// where i=d means daily data
func BuildURL(symbol string) string {
	return fmt.Sprintf(
		"https://stooq.com/q/d/l/?s=%s&i=d",
		url.QueryEscape(symbol),
	)
}

// ReadSingle fetches data for a single symbol.
func (s *StooqReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate symbol
	if err := s.ValidateSymbol(symbol); err != nil {
		return nil, err
	}

	// Build URL - use custom baseURL if set (for testing), otherwise use standard format
	var urlStr string
	if s.baseURL != "" {
		urlStr = fmt.Sprintf(s.baseURL, url.QueryEscape(symbol))
	} else {
		urlStr = BuildURL(symbol)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	resp, err := s.client.Do(req)
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

	// Parse CSV
	data, err := ParseCSV(body)
	if err != nil {
		return nil, fmt.Errorf("parse CSV: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from Stooq.
// Symbols are fetched in parallel for better performance.
func (s *StooqReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return s.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (s *StooqReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := s.ReadSingle(ctx, sym, start, end)

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

// ValidateSymbol checks if a symbol is valid for Stooq.
func (s *StooqReader) ValidateSymbol(symbol string) error {
	return s.BaseSource.ValidateSymbol(symbol)
}
