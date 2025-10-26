// Package alphavantage provides an Alpha Vantage data source reader.
package alphavantage

import (
	"context"
	"fmt"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

// AlphaVantageReader fetches data from the Alpha Vantage API.
type AlphaVantageReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
	apiKey string
}

// NewAlphaVantageReader creates a new Alpha Vantage data reader.
// An API key is required to use the Alpha Vantage API.
func NewAlphaVantageReader(opts *internalhttp.ClientOptions, apiKey string) *AlphaVantageReader {
	return &AlphaVantageReader{
		BaseSource: sources.NewBaseSource("alphavantage"),
		client:     internalhttp.NewRetryableClient(opts),
		apiKey:     apiKey,
	}
}

// BuildURL constructs the Alpha Vantage API URL for fetching daily time series data.
// The Alpha Vantage API format is:
// https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol={symbol}&apikey={apikey}&outputsize=full
func BuildURL(symbol, apiKey string) string {
	return fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full",
		symbol,
		apiKey,
	)
}

// ReadSingle fetches data for a single stock symbol.
func (a *AlphaVantageReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// Read fetches data for multiple stock symbols.
func (a *AlphaVantageReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// ValidateSymbol checks if a symbol is valid for Alpha Vantage.
func (a *AlphaVantageReader) ValidateSymbol(symbol string) error {
	return a.BaseSource.ValidateSymbol(symbol)
}
