// Package stooq provides a Stooq data source reader.
package stooq

import (
	"context"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

// StooqReader fetches data from Stooq.
type StooqReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
}

// NewStooqReader creates a new Stooq data reader.
func NewStooqReader(opts *internalhttp.ClientOptions) *StooqReader {
	return &StooqReader{
		BaseSource: sources.NewBaseSource("stooq"),
		client:     internalhttp.NewRetryableClient(opts),
	}
}

// ReadSingle fetches data for a single symbol.
func (s *StooqReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// Read fetches data for multiple symbols.
func (s *StooqReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// ValidateSymbol checks if a symbol is valid for Stooq.
func (s *StooqReader) ValidateSymbol(symbol string) error {
	return s.BaseSource.ValidateSymbol(symbol)
}
