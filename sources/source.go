// Package sources defines the base interface for all data sources.
package sources

import (
	"context"
	"time"

	"github.com/julianshen/gonp-datareader/internal/utils"
)

// Reader is the main interface for all data sources.
// Implementations must be safe for concurrent use.
type Reader interface {
	// Read fetches data for the given symbols within the date range.
	// It returns an error if any symbol is invalid or if the request fails.
	// The return type is interface{} to allow flexibility for different data sources
	// until we integrate with gonp DataFrames.
	Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error)

	// ReadSingle fetches data for a single symbol within the date range.
	// This is a convenience method that may be more efficient than Read for single symbols.
	ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error)

	// ValidateSymbol checks if a symbol is valid for this data source.
	ValidateSymbol(symbol string) error

	// Name returns the display name of the data source.
	Name() string

	// Source returns the identifier of the data source (e.g., "yahoo", "fred").
	Source() string
}

// BaseSource provides common functionality for data source implementations.
type BaseSource struct {
	source string
}

// NewBaseSource creates a new BaseSource.
func NewBaseSource(source string) *BaseSource {
	return &BaseSource{
		source: source,
	}
}

// Name returns the display name of the data source.
// By default, it returns the source identifier.
func (b *BaseSource) Name() string {
	return b.source
}

// Source returns the identifier of the data source.
func (b *BaseSource) Source() string {
	return b.source
}

// ValidateSymbol validates a symbol using the common validation rules.
// Data sources can override this method for source-specific validation.
func (b *BaseSource) ValidateSymbol(symbol string) error {
	return utils.ValidateSymbol(symbol)
}
