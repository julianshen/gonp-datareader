package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader/internal/ratelimit"
)

func TestNewRateLimiter(t *testing.T) {
	limiter := ratelimit.NewRateLimiter(10.0, 1)

	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
}

func TestRateLimiter_AllowsRequestsAtRate(t *testing.T) {
	// 10 requests per second, burst of 1
	limiter := ratelimit.NewRateLimiter(10.0, 1)

	ctx := context.Background()

	// First request should be immediate
	start := time.Now()
	err := limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	if elapsed > 10*time.Millisecond {
		t.Errorf("First request took too long: %v", elapsed)
	}
}

func TestRateLimiter_BlocksWhenRateExceeded(t *testing.T) {
	// 5 requests per second, burst of 1
	limiter := ratelimit.NewRateLimiter(5.0, 1)

	ctx := context.Background()

	// First request is immediate
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	// Second request should block for ~200ms (1/5 second)
	start := time.Now()
	err = limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}

	// Should wait at least 150ms (with some tolerance)
	if elapsed < 150*time.Millisecond {
		t.Errorf("Expected delay >= 150ms, got %v", elapsed)
	}
}

func TestRateLimiter_RespectsContext(t *testing.T) {
	// Very slow rate: 0.1 requests per second (1 per 10 seconds)
	limiter := ratelimit.NewRateLimiter(0.1, 1)

	// Use first token
	ctx1 := context.Background()
	err := limiter.Wait(ctx1)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	// Cancel context immediately
	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Second request should fail due to context timeout
	err = limiter.Wait(ctx2)
	if err == nil {
		t.Error("Expected error due to context timeout, got nil")
	}
}

func TestRateLimiter_AllowsBurst(t *testing.T) {
	// 10 requests per second, burst of 5
	limiter := ratelimit.NewRateLimiter(10.0, 5)

	ctx := context.Background()

	start := time.Now()

	// Should allow 5 immediate requests (burst)
	for i := 0; i < 5; i++ {
		err := limiter.Wait(ctx)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
	}

	elapsed := time.Since(start)

	// All 5 requests should complete quickly (within 50ms)
	if elapsed > 50*time.Millisecond {
		t.Errorf("Burst requests took too long: %v", elapsed)
	}

	// 6th request should block
	start = time.Now()
	err := limiter.Wait(ctx)
	elapsed = time.Since(start)

	if err != nil {
		t.Fatalf("6th request failed: %v", err)
	}

	// Should wait at least 50ms
	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected delay >= 50ms for 6th request, got %v", elapsed)
	}
}

func TestRateLimiter_ZeroRate(t *testing.T) {
	// Zero rate means unlimited
	limiter := ratelimit.NewRateLimiter(0, 1)

	ctx := context.Background()

	start := time.Now()

	// Should allow many requests immediately
	for i := 0; i < 100; i++ {
		err := limiter.Wait(ctx)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
	}

	elapsed := time.Since(start)

	// All requests should complete very quickly
	if elapsed > 50*time.Millisecond {
		t.Errorf("Unlimited rate took too long: %v", elapsed)
	}
}

func TestRateLimiter_NilLimiter(t *testing.T) {
	var limiter *ratelimit.RateLimiter

	ctx := context.Background()

	// Nil limiter should allow requests without blocking
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Nil limiter should allow requests, got error: %v", err)
	}
}
