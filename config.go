package datareader

import "time"

// Options configures the behavior of a data reader.
//
// All fields are optional. If nil is passed to DataReader or Read,
// sensible defaults are used. Use DefaultOptions() to get a pre-configured
// Options struct with recommended defaults.
//
// # Example Usage
//
//	// Basic configuration
//	opts := &datareader.Options{
//		Timeout: 60 * time.Second,
//	}
//
//	// Full configuration
//	opts := &datareader.Options{
//		// Authentication
//		APIKey: "your-api-key",
//
//		// HTTP settings
//		Timeout:   30 * time.Second,
//		UserAgent: "MyApp/1.0",
//
//		// Retry logic
//		MaxRetries: 3,
//		RetryDelay: 2 * time.Second,
//
//		// Rate limiting
//		RateLimit: 5.0, // 5 requests per second
//
//		// Caching
//		CacheDir: ".cache/datareader",
//		CacheTTL: 24 * time.Hour,
//	}
type Options struct {
	// APIKey for sources that require authentication.
	// Required for: alphavantage, iex
	// Optional for: fred (higher rate limits with key)
	// Not used for: yahoo, worldbank, stooq
	APIKey string

	// Timeout specifies the maximum duration for HTTP requests.
	// Zero or negative values mean no timeout.
	// Default: 30 seconds
	Timeout time.Duration

	// MaxRetries specifies the maximum number of retry attempts for failed requests.
	// Retries use exponential backoff with RetryDelay as the base delay.
	// Default: 3
	MaxRetries int

	// RetryDelay specifies the initial delay between retry attempts.
	// Actual delay increases exponentially: RetryDelay * 2^attempt
	// Default: 1 second
	RetryDelay time.Duration

	// EnableCache enables response caching (deprecated, use CacheDir instead).
	// Caching is automatically enabled when CacheDir is set.
	EnableCache bool

	// CacheDir specifies the directory for cached responses.
	// If empty, caching is disabled.
	// Cached responses are stored with SHA-256 hashed filenames.
	CacheDir string

	// CacheTTL specifies how long cached responses remain valid.
	// Zero means responses are cached indefinitely.
	// Expired entries are automatically cleaned on access.
	CacheTTL time.Duration

	// RateLimit specifies the maximum number of requests per second.
	// Zero or negative values mean no rate limiting.
	// Uses token bucket algorithm for smooth rate limiting.
	RateLimit float64

	// UserAgent specifies the User-Agent header for HTTP requests.
	// Some sources (like Yahoo Finance) may require a valid browser User-Agent.
	// Default: Chrome/Safari User-Agent string
	UserAgent string
}

// DefaultOptions returns a new Options struct with recommended default values.
//
// Default values:
//   - Timeout: 30 seconds
//   - MaxRetries: 3
//   - RetryDelay: 1 second
//   - UserAgent: Chrome browser User-Agent
//
// # Example Usage
//
//	opts := datareader.DefaultOptions()
//	opts.APIKey = "your-api-key"
//	opts.CacheDir = ".cache"
//
//	reader, err := datareader.DataReader("fred", opts)
func DefaultOptions() *Options {
	return &Options{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}
