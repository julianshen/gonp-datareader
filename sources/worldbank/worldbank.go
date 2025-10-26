// Package worldbank provides a World Bank data source reader.
package worldbank

import (
	"context"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

// WorldBankReader fetches data from the World Bank API.
type WorldBankReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
}

// NewWorldBankReader creates a new World Bank data reader.
func NewWorldBankReader(opts *internalhttp.ClientOptions) *WorldBankReader {
	return &WorldBankReader{
		BaseSource: sources.NewBaseSource("worldbank"),
		client:     internalhttp.NewRetryableClient(opts),
	}
}

// ReadSingle fetches data for a single indicator and country.
// The symbol parameter should be in the format "country/indicator", e.g., "USA/NY.GDP.MKTP.CD"
func (w *WorldBankReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// Read fetches data for multiple indicators and countries.
func (w *WorldBankReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// ValidateSymbol checks if a symbol is valid for World Bank.
func (w *WorldBankReader) ValidateSymbol(symbol string) error {
	return w.BaseSource.ValidateSymbol(symbol)
}
