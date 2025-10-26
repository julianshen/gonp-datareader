// Package yahoo provides data access to Yahoo Finance.
package yahoo

import (
	"context"
	"fmt"
	"time"

	"github.com/julianshen/gonp-datareader"
	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

const (
	// yahooAPIURL is the base URL for Yahoo Finance historical data API
	yahooAPIURL = "https://query1.finance.yahoo.com/v7/finance/download/%s"
)

// YahooReader fetches data from Yahoo Finance.
type YahooReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
}

// NewYahooReader creates a new Yahoo Finance data reader.
func NewYahooReader(opts *datareader.Options) *YahooReader {
	if opts == nil {
		opts = datareader.DefaultOptions()
	}

	clientOpts := &internalhttp.ClientOptions{
		Timeout:    opts.Timeout,
		UserAgent:  opts.UserAgent,
		MaxRetries: opts.MaxRetries,
		RetryDelay: opts.RetryDelay,
	}

	return &YahooReader{
		BaseSource: sources.NewBaseSource("yahoo"),
		client:     internalhttp.NewRetryableClient(clientOpts),
	}
}

// Name returns the display name of the data source.
func (y *YahooReader) Name() string {
	return "Yahoo Finance"
}

// BuildURL constructs the Yahoo Finance API URL for the given symbol and date range.
func (y *YahooReader) BuildURL(symbol string, start, end time.Time) string {
	baseURL := fmt.Sprintf(yahooAPIURL, symbol)

	// Convert dates to Unix timestamps
	period1 := start.Unix()
	period2 := end.Unix()

	// Build query parameters
	url := fmt.Sprintf("%s?period1=%d&period2=%d&interval=1d&events=history&includeAdjustedClose=true",
		baseURL, period1, period2)

	return url
}

// Read fetches data for multiple symbols from Yahoo Finance.
func (y *YahooReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// For now, return a simple error - we'll implement actual fetching later
	return nil, fmt.Errorf("not yet implemented")
}

// ReadSingle fetches data for a single symbol from Yahoo Finance.
func (y *YahooReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// For now, return a simple error - we'll implement actual fetching later
	return nil, fmt.Errorf("not yet implemented")
}
