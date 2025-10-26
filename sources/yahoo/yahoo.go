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
func (y *YahooReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Validate inputs
	if err := utils.ValidateSymbols(symbols); err != nil {
		return nil, fmt.Errorf("invalid symbols: %w", err)
	}

	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Fetch data for each symbol
	results := make(map[string]*ParsedData)
	for _, symbol := range symbols {
		data, err := y.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", symbol, err)
		}

		if parsedData, ok := data.(*ParsedData); ok {
			results[symbol] = parsedData
		}
	}

	return results, nil
}
