package http

import (
	"net/http"
	"time"

	"github.com/julianshen/gonp-datareader/internal/ratelimit"
)

// RetryableClient wraps an http.Client with retry logic.
type RetryableClient struct {
	client      *http.Client
	maxRetries  int
	retryDelay  time.Duration
	userAgent   string
	rateLimiter *ratelimit.RateLimiter
}

// NewRetryableClient creates a new HTTP client with retry logic.
func NewRetryableClient(opts *ClientOptions) *RetryableClient {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	// Create rate limiter if rate limit is configured
	var limiter *ratelimit.RateLimiter
	if opts.RateLimit > 0 {
		// Use burst of 1 for strict rate limiting
		limiter = ratelimit.NewRateLimiter(opts.RateLimit, 1)
	}

	return &RetryableClient{
		client:      NewHTTPClient(opts),
		maxRetries:  opts.MaxRetries,
		retryDelay:  opts.RetryDelay,
		userAgent:   opts.UserAgent,
		rateLimiter: limiter,
	}
}

// Do executes an HTTP request with retry logic.
func (c *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Apply rate limiting before making request
		if c.rateLimiter != nil {
			if err := c.rateLimiter.Wait(req.Context()); err != nil {
				return nil, err
			}
		}

		// Clone the request for retry attempts
		reqClone := req.Clone(req.Context())

		// Set User-Agent header if configured
		if c.userAgent != "" {
			reqClone.Header.Set("User-Agent", c.userAgent)
		}

		resp, err = c.client.Do(reqClone)

		// Check if we should retry
		if !ShouldRetry(resp, err) {
			return resp, err
		}

		// Don't sleep after the last attempt
		if attempt < c.maxRetries {
			time.Sleep(c.retryDelay * time.Duration(attempt+1))
		}
	}

	// Return the last response/error
	return resp, err
}

// ShouldRetry determines if a request should be retried based on the response or error.
func ShouldRetry(resp *http.Response, err error) bool {
	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on nil response (shouldn't happen but be defensive)
	if resp == nil {
		return true
	}

	// Retry on 5xx server errors
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return true
	}

	// Don't retry on success or client errors (4xx)
	return false
}
