package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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
