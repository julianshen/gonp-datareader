// Package iex provides an IEX Cloud data source reader.
package iex

import (
	"context"
	"fmt"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
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
	// TODO: Implement
	return nil, nil
}

// Read fetches data for multiple stock symbols.
func (i *IEXReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// ValidateSymbol checks if a symbol is valid for IEX Cloud.
func (i *IEXReader) ValidateSymbol(symbol string) error {
	return i.BaseSource.ValidateSymbol(symbol)
}
