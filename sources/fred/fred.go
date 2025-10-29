// Package fred provides data access to FRED (Federal Reserve Economic Data).
package fred

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
	// fredAPIURL is the base URL for FRED API
	fredAPIURL = "https://api.stlouisfed.org/fred/series/observations"
)

// FREDReader fetches data from FRED (Federal Reserve Economic Data).
type FREDReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	apiKey  string
	baseURL string // For testing with mock servers
}

// NewFREDReader creates a new FRED data reader.
func NewFREDReader(opts *internalhttp.ClientOptions) *FREDReader {
	return NewFREDReaderWithBaseURL(opts, fredAPIURL)
}

// NewFREDReaderWithBaseURL creates a new FRED reader with a custom base URL.
// This is primarily used for testing with mock servers.
func NewFREDReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *FREDReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &FREDReader{
		BaseSource: sources.NewBaseSource("fred"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// NewFREDReaderWithAPIKey creates a new FRED data reader with an API key.
func NewFREDReaderWithAPIKey(opts *internalhttp.ClientOptions, apiKey string) *FREDReader {
	reader := NewFREDReader(opts)
	reader.apiKey = apiKey
	return reader
}

// SetAPIKey sets the API key for FRED requests.
func (f *FREDReader) SetAPIKey(apiKey string) {
	f.apiKey = apiKey
}

// GetAPIKey returns the currently configured API key.
func (f *FREDReader) GetAPIKey() string {
	return f.apiKey
}

// Name returns the display name of the data source.
func (f *FREDReader) Name() string {
	return "FRED"
}

// BuildURL constructs the FRED API URL for the given series and date range.
func (f *FREDReader) BuildURL(seriesID string, start, end time.Time, apiKey string) string {
	// Format dates as YYYY-MM-DD
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	// Use custom baseURL if set (for testing), otherwise use standard FRED URL
	baseURL := f.baseURL
	if baseURL == "" {
		baseURL = fredAPIURL
	}

	url := fmt.Sprintf("%s?series_id=%s&api_key=%s&observation_start=%s&observation_end=%s&file_type=json",
		baseURL, seriesID, apiKey, startStr, endStr)

	return url
}

// ReadSingle fetches data for a single series from FRED.
func (f *FREDReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := f.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Check API key
	if f.apiKey == "" {
		return nil, fmt.Errorf("FRED API key is required")
	}

	// Build URL
	url := f.BuildURL(symbol, start, end, f.apiKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("FRED API returned status %d (failed to read response body: %w)", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("FRED API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	data, err := ParseJSON(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple series from FRED.
func (f *FREDReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Check API key
	if f.apiKey == "" {
		return nil, fmt.Errorf("FRED API key is required")
	}

	// Fetch data for each series
	results := make(map[string]*ParsedData)
	for _, symbol := range symbols {
		data, err := f.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", symbol, err)
		}

		if parsedData, ok := data.(*ParsedData); ok {
			results[symbol] = parsedData
		}
	}

	return results, nil
}
