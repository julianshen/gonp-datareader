// Package ratelimit provides rate limiting functionality for HTTP requests.
package ratelimit

import (
	"context"

	"golang.org/x/time/rate"
)

// RateLimiter controls the rate of requests.
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter with the specified rate and burst.
// The rate is in requests per second. A rate of 0 means unlimited.
// The burst is the maximum number of requests that can be made at once.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	if rps <= 0 {
		// Unlimited rate
		return &RateLimiter{
			limiter: rate.NewLimiter(rate.Inf, 0),
		}
	}

	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Wait blocks until the rate limiter allows the request to proceed.
// It returns an error if the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	// Handle nil limiter (allows unlimited requests)
	if r == nil || r.limiter == nil {
		return nil
	}

	return r.limiter.Wait(ctx)
}
