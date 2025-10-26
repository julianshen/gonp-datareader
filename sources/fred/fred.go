// Package fred provides data access to FRED (Federal Reserve Economic Data).
package fred

import (
	"context"
	"fmt"
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
	client *internalhttp.RetryableClient
	apiKey string
}

// NewFREDReader creates a new FRED data reader.
func NewFREDReader(opts *internalhttp.ClientOptions) *FREDReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}

	return &FREDReader{
		BaseSource: sources.NewBaseSource("fred"),
		client:     internalhttp.NewRetryableClient(opts),
	}
}

// NewFREDReaderWithAPIKey creates a new FRED data reader with an API key.
func NewFREDReaderWithAPIKey(opts *internalhttp.ClientOptions, apiKey string) *FREDReader {
	reader := NewFREDReader(opts)
	reader.apiKey = apiKey
	return reader
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

	url := fmt.Sprintf("%s?series_id=%s&api_key=%s&observation_start=%s&observation_end=%s&file_type=json",
		fredAPIURL, seriesID, apiKey, startStr, endStr)

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

	// TODO: Implement actual HTTP request and parsing
	return nil, fmt.Errorf("not implemented yet")
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

	// TODO: Implement fetching multiple series
	return nil, fmt.Errorf("not implemented yet")
}
