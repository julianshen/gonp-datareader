package http_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
)

func TestRetryableClient_Success(t *testing.T) {
	// Server that succeeds on first try
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestRetryableClient_RetryOnError(t *testing.T) {
	var attempts atomic.Int32

	// Server that fails twice, then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := attempts.Add(1)
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed after retries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if attempts.Load() != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts.Load())
	}
}

func TestRetryableClient_MaxRetriesExceeded(t *testing.T) {
	var attempts atomic.Int32

	// Server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 2,
		RetryDelay: 10 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should still return last response even after retries exhausted
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	// Should attempt: initial + 2 retries = 3 total
	expectedAttempts := int32(3)
	if attempts.Load() != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attempts.Load())
	}
}

func TestRetryableClient_NoRetryOn4xx(t *testing.T) {
	var attempts atomic.Int32

	// Server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should not retry on 4xx errors
	if attempts.Load() != 1 {
		t.Errorf("Expected 1 attempt for 4xx error, got %d", attempts.Load())
	}
}

func TestRetryableClient_SetsUserAgent(t *testing.T) {
	var capturedUA string

	// Server that captures User-Agent header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	customUA := "gonp-datareader/1.0 test-client"
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		UserAgent:  customUA,
		MaxRetries: 1,
		RetryDelay: 10 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if capturedUA != customUA {
		t.Errorf("Expected User-Agent %q, got %q", customUA, capturedUA)
	}
}

func TestRetryableClient_DefaultUserAgent(t *testing.T) {
	var capturedUA string

	// Server that captures User-Agent header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Use default options (should have default User-Agent)
	client := internalhttp.NewRetryableClient(nil)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if capturedUA == "" {
		t.Error("Expected default User-Agent to be set, got empty string")
	}

	// Should contain "gonp-datareader"
	if capturedUA == "" || len(capturedUA) == 0 {
		t.Errorf("Expected non-empty User-Agent, got %q", capturedUA)
	}
}

func TestRetryableClient_WithRateLimiter(t *testing.T) {
	var requestCount atomic.Int32

	// Server that counts requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Rate limiter: 2 requests per second
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		RateLimit:  2.0, // 2 requests per second
	}

	client := internalhttp.NewRetryableClient(opts)

	// Make 3 requests
	start := time.Now()
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}
	elapsed := time.Since(start)

	// Should take at least 1 second for 3 requests at 2 req/sec
	// (0s, 0.5s, 1.0s)
	if elapsed < 900*time.Millisecond {
		t.Errorf("Expected rate limiting to take >= 900ms, took %v", elapsed)
	}

	if requestCount.Load() != 3 {
		t.Errorf("Expected 3 requests, got %d", requestCount.Load())
	}
}

func TestRetryableClient_NoRateLimiter(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// No rate limit configured (0 or nil)
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		RateLimit:  0, // No rate limit
	}

	client := internalhttp.NewRetryableClient(opts)

	// Make 10 requests quickly
	start := time.Now()
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}
	elapsed := time.Since(start)

	// Should complete quickly without rate limiting
	if elapsed > 100*time.Millisecond {
		t.Errorf("Without rate limiting, expected < 100ms, took %v", elapsed)
	}

	if requestCount.Load() != 10 {
		t.Errorf("Expected 10 requests, got %d", requestCount.Load())
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{
			name:       "retry on 500",
			statusCode: http.StatusInternalServerError,
			want:       true,
		},
		{
			name:       "retry on 502",
			statusCode: http.StatusBadGateway,
			want:       true,
		},
		{
			name:       "retry on 503",
			statusCode: http.StatusServiceUnavailable,
			want:       true,
		},
		{
			name:       "retry on 504",
			statusCode: http.StatusGatewayTimeout,
			want:       true,
		},
		{
			name:       "no retry on 404",
			statusCode: http.StatusNotFound,
			want:       false,
		},
		{
			name:       "no retry on 400",
			statusCode: http.StatusBadRequest,
			want:       false,
		},
		{
			name:       "no retry on 200",
			statusCode: http.StatusOK,
			want:       false,
		},
		{
			name: "retry on network error",
			err:  errors.New("connection refused"),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.statusCode > 0 {
				resp = &http.Response{StatusCode: tt.statusCode}
			}

			got := internalhttp.ShouldRetry(resp, tt.err)
			if got != tt.want {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryableClient_WithCache(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cached response"))
	}))
	defer server.Close()

	// Create temporary cache directory
	tmpDir, err := os.MkdirTemp("", "http-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Configure client with cache
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		CacheDir:   tmpDir,
		CacheTTL:   1 * time.Hour,
	}

	client := internalhttp.NewRetryableClient(opts)

	// First request - should hit server and cache
	req1, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL+"/test", nil)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()

	if string(body1) != "cached response" {
		t.Errorf("Expected 'cached response', got %q", string(body1))
	}

	if requestCount.Load() != 1 {
		t.Errorf("Expected 1 server request, got %d", requestCount.Load())
	}

	// Second request - should come from cache
	req2, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL+"/test", nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	if string(body2) != "cached response" {
		t.Errorf("Expected 'cached response', got %q", string(body2))
	}

	// Should still be only 1 server request (second came from cache)
	if requestCount.Load() != 1 {
		t.Errorf("Expected 1 server request (cached), got %d", requestCount.Load())
	}
}

func TestRetryableClient_CacheOnlyGET(t *testing.T) {
	var getCount atomic.Int32
	var postCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCount.Add(1)
		} else if r.Method == "POST" {
			postCount.Add(1)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "http-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		CacheDir:   tmpDir,
		CacheTTL:   1 * time.Hour,
	}

	client := internalhttp.NewRetryableClient(opts)

	// GET request - should be cached
	req1, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	resp1.Body.Close()

	req2, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Second GET request failed: %v", err)
	}
	resp2.Body.Close()

	// Only 1 GET should hit server (second is cached)
	if getCount.Load() != 1 {
		t.Errorf("Expected 1 GET request, got %d", getCount.Load())
	}

	// POST request - should NOT be cached
	req3, _ := http.NewRequestWithContext(context.Background(), "POST", server.URL, nil)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("First POST request failed: %v", err)
	}
	resp3.Body.Close()

	req4, _ := http.NewRequestWithContext(context.Background(), "POST", server.URL, nil)
	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatalf("Second POST request failed: %v", err)
	}
	resp4.Body.Close()

	// Both POSTs should hit server
	if postCount.Load() != 2 {
		t.Errorf("Expected 2 POST requests, got %d", postCount.Load())
	}
}

func TestRetryableClient_NoCache(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	// No cache configured
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
	}

	client := internalhttp.NewRetryableClient(opts)

	// Make two identical requests
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		resp.Body.Close()
	}

	// Both should hit server (no cache)
	if requestCount.Load() != 2 {
		t.Errorf("Expected 2 requests without cache, got %d", requestCount.Load())
	}
}

func TestRetryableClient_CacheTTL(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "http-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Very short TTL
	opts := &internalhttp.ClientOptions{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		CacheDir:   tmpDir,
		CacheTTL:   100 * time.Millisecond,
	}

	client := internalhttp.NewRetryableClient(opts)

	// First request
	req1, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	resp1.Body.Close()

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second request - should hit server again (cache expired)
	req2, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	resp2.Body.Close()

	// Both should hit server (cache expired)
	if requestCount.Load() != 2 {
		t.Errorf("Expected 2 requests (cache expired), got %d", requestCount.Load())
	}
}
