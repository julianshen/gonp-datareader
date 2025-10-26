// Package iex provides an IEX Cloud data source reader.
package iex

import (
	"context"
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
