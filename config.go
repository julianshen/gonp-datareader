// Package datareader provides remote data access for financial
// and economic data sources, designed to work with gonp DataFrames.
//
// Basic usage:
//
//	ctx := context.Background()
//	df, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
package datareader

import "time"

// Options configures the behavior of a Reader.
type Options struct {
	// APIKey for sources that require authentication
	APIKey string

	// Timeout specifies the HTTP request timeout.
	// Zero means no timeout.
	Timeout time.Duration

	// MaxRetries for failed requests
	MaxRetries int

	// RetryDelay between retry attempts
	RetryDelay time.Duration

	// EnableCache enables response caching
	EnableCache bool

	// CacheDir specifies the cache directory
	CacheDir string

	// RateLimit specifies requests per second limit
	RateLimit float64

	// UserAgent for HTTP requests
	UserAgent string
}

// DefaultOptions returns default configuration for data readers.
func DefaultOptions() *Options {
	return &Options{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}
