// Package oecd provides data access to OECD API.
package oecd

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
	// oecdAPIURL is the base URL for OECD SDMX-JSON API
	oecdAPIURL = "https://stats.oecd.org/sdmx-json/data/%s/all"
)

// OECDReader fetches data from OECD API.
type OECDReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewOECDReader creates a new OECD data reader.
func NewOECDReader(opts *internalhttp.ClientOptions) *OECDReader {
	return NewOECDReaderWithBaseURL(opts, oecdAPIURL)
}

// NewOECDReaderWithBaseURL creates a new OECD reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewOECDReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *OECDReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &OECDReader{
		BaseSource: sources.NewBaseSource("oecd"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name of the data source.
func (o *OECDReader) Name() string {
	return "OECD"
}

// ValidateSymbol validates an OECD dataset identifier.
// OECD symbols are in the format "DATASET/DIMENSIONS" or just "DATASET".
// Examples: "MEI/USA", "QNA/AUS.GDP", "REGION_ECONOM"
func (o *OECDReader) ValidateSymbol(symbol string) error {
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}

	// Check for invalid characters (spaces)
	if strings.Contains(symbol, " ") {
		return fmt.Errorf("symbol cannot contain spaces")
	}

	return nil
}

// BuildURL constructs the OECD API URL for the given symbol and date range.
func (o *OECDReader) BuildURL(symbol string, start, end time.Time) string {
	// Format dates as YYYY-MM or YYYY-QN for quarterly data
	startDate := start.Format("2006-01")
	endDate := end.Format("2006-01")

	// Build URL with symbol and date range
	url := fmt.Sprintf(o.baseURL, symbol)
	url += fmt.Sprintf("?startPeriod=%s&endPeriod=%s&dimensionAtObservation=AllDimensions",
		startDate, endDate)

	return url
}

// ReadSingle fetches data for a single symbol from OECD.
func (o *OECDReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := o.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Build URL
	url := o.BuildURL(symbol, start, end)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Accept header for JSON
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("OECD returned status %d (failed to read response body: %w)", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("OECD returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	data, err := ParseJSON(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols from OECD.
// Symbols are fetched in parallel for better performance.
func (o *OECDReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Use parallel fetching for multiple symbols
	return o.readParallel(ctx, symbols, start, end)
}

// readParallel fetches multiple symbols in parallel using a worker pool.
func (o *OECDReader) readParallel(ctx context.Context, symbols []string, start, end time.Time) (map[string]*ParsedData, error) {
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
			data, err := o.ReadSingle(ctx, sym, start, end)

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
