// Package http provides HTTP client utilities for gonp-datareader.
package http

import (
	"net/http"
	"time"
)

// ClientOptions configures the HTTP client behavior.
type ClientOptions struct {
	// Timeout specifies the HTTP request timeout
	Timeout time.Duration

	// UserAgent specifies the User-Agent header
	UserAgent string

	// MaxRetries specifies the maximum number of retry attempts
	MaxRetries int

	// RetryDelay specifies the delay between retry attempts
	RetryDelay time.Duration

	// RateLimit specifies requests per second limit (0 = unlimited)
	RateLimit float64
}

// DefaultClientOptions returns default HTTP client options.
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout:    30 * time.Second,
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// NewHTTPClient creates a new HTTP client with the specified options.
// If opts is nil, default options are used.
func NewHTTPClient(opts *ClientOptions) *http.Client {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	client := &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return client
}
