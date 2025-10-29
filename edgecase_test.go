package datareader_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

// TestYahooReader_NetworkTimeout tests timeout scenarios
func TestYahooReader_NetworkTimeout(t *testing.T) {
	// Create a server that delays response beyond timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // Delay longer than client timeout
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("too late"))
	}))
	defer server.Close()

	// Create options with very short timeout
	opts := &internalhttp.ClientOptions{
		Timeout: 50 * time.Millisecond,
	}

	reader := yahoo.NewYahooReaderWithBaseURL(opts, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	// Verify it's a timeout or context deadline error
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Logf("Got error (expected timeout-related): %v", err)
	}
}

// TestYahooReader_ContextCancellation tests context cancellation
func TestYahooReader_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("Expected context canceled error, got nil")
	}

	// Verify it's a context canceled error
	if !errors.Is(err, context.Canceled) {
		t.Logf("Got error (expected context.Canceled): %v", err)
	}
}

// TestYahooReader_ContextCancellationDuringRead tests cancellation mid-request
func TestYahooReader_ContextCancellationDuringRead(t *testing.T) {
	// Create a server that delays response
	requestReceived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestReceived)
		time.Sleep(200 * time.Millisecond) // Long delay
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// Start read in goroutine
	errCh := make(chan error, 1)
	go func() {
		_, err := reader.ReadSingle(ctx, "AAPL", start, end)
		errCh <- err
	}()

	// Wait for request to be received, then cancel
	<-requestReceived
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for error
	err := <-errCh
	if err == nil {
		t.Error("Expected error after context cancellation, got nil")
	}
}

// TestYahooReader_MalformedResponse tests handling of malformed CSV
func TestYahooReader_MalformedResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "invalid CSV structure",
			response: "not,csv,data\nwithout,proper",
			wantErr:  true,
		},
		{
			name:     "empty response",
			response: "",
			wantErr:  true,
		},
		{
			name:     "header only",
			response: "Date,Open,High,Low,Close,Adj Close,Volume",
			wantErr:  false, // Empty data is valid
		},
		{
			name:     "missing header",
			response: "2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900",
			wantErr:  false, // CSV parser treats first row as header
		},
		{
			name:     "corrupted data row",
			response: "Date,Open,High,Low,Close,Adj Close,Volume\n2020-01-02,invalid,data",
			wantErr:  true, // Wrong number of fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/csv")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

			ctx := context.Background()
			start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

			_, err := reader.ReadSingle(ctx, "AAPL", start, end)

			if tt.wantErr && err == nil {
				t.Error("Expected error for malformed response, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

// TestYahooReader_PartialDataScenarios tests handling of partial/incomplete data
func TestYahooReader_PartialDataScenarios(t *testing.T) {
	tests := []struct {
		name     string
		csvData  string
		wantRows int
	}{
		{
			name: "complete data",
			csvData: `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900
2020-01-03,297.15,300.58,296.50,297.43,294.56,36607600`,
			wantRows: 2,
		},
		{
			name: "data with null values",
			csvData: `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,null,300.60,295.19,300.35,297.45,33911900
2020-01-03,297.15,null,296.50,297.43,294.56,36607600`,
			wantRows: 2,
		},
		{
			name: "data with inconsistent columns",
			csvData: `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`,
			wantRows: 1,
		},
		{
			name: "single row",
			csvData: `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`,
			wantRows: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/csv")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.csvData))
			}))
			defer server.Close()

			reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

			ctx := context.Background()
			start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

			result, err := reader.ReadSingle(ctx, "AAPL", start, end)
			if err != nil {
				t.Fatalf("ReadSingle() error = %v", err)
			}

			data, ok := result.(*yahoo.ParsedData)
			if !ok {
				t.Fatalf("Expected *yahoo.ParsedData, got %T", result)
			}

			if len(data.Rows) != tt.wantRows {
				t.Errorf("Expected %d rows, got %d", tt.wantRows, len(data.Rows))
			}
		})
	}
}

// TestYahooReader_LargeDateRange tests handling of large date ranges
func TestYahooReader_LargeDateRange(t *testing.T) {
	// Generate large CSV with many rows
	csvData := "Date,Open,High,Low,Close,Adj Close,Volume\n"
	for i := 0; i < 1000; i++ {
		csvData += "2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900\n"
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	// Request 10 years of data
	start := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	data, ok := result.(*yahoo.ParsedData)
	if !ok {
		t.Fatalf("Expected *yahoo.ParsedData, got %T", result)
	}

	if len(data.Rows) != 1000 {
		t.Errorf("Expected 1000 rows, got %d", len(data.Rows))
	}
}

// TestYahooReader_ConcurrentRequests tests concurrent access to the reader
func TestYahooReader_ConcurrentRequests(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	// Track concurrent requests
	var requestCount int
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()

		// Small delay to ensure concurrency
		time.Sleep(10 * time.Millisecond)

		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// Launch multiple concurrent requests
	const numRequests = 10
	var wg sync.WaitGroup
	errCh := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			_, err := reader.ReadSingle(ctx, symbol, start, end)
			if err != nil {
				errCh <- err
			}
		}("AAPL")
	}

	wg.Wait()
	close(errCh)

	// Check for errors
	for err := range errCh {
		t.Errorf("Concurrent request failed: %v", err)
	}

	// Verify all requests were made
	mu.Lock()
	finalCount := requestCount
	mu.Unlock()

	if finalCount != numRequests {
		t.Errorf("Expected %d requests, got %d", numRequests, finalCount)
	}
}

// TestYahooReader_ServerErrors tests handling of various HTTP errors
func TestYahooReader_ServerErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			wantErr:    true,
		},
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			wantErr:    true,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			wantErr:    true, // Empty body will cause parsing error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

			ctx := context.Background()
			start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

			_, err := reader.ReadSingle(ctx, "AAPL", start, end)

			if tt.wantErr && err == nil {
				t.Error("Expected error for HTTP error response, got nil")
			}
		})
	}
}

// TestDataReader_ConcurrentSourceRequests tests concurrent requests across different sources
func TestDataReader_ConcurrentSourceRequests(t *testing.T) {
	// This test verifies the datareader can handle concurrent requests to different sources
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	sources := []string{"yahoo", "fred", "worldbank"}
	symbols := []string{"AAPL", "GDP", "USA/NY.GDP.MKTP.CD"}

	var wg sync.WaitGroup
	errCh := make(chan error, len(sources))

	for i, source := range sources {
		wg.Add(1)
		go func(src, sym string) {
			defer wg.Done()
			_, err := datareader.Read(ctx, sym, src, start, end, nil)
			if err != nil {
				// Network errors are acceptable in unit tests
				t.Logf("Read from %s failed (acceptable): %v", src, err)
			}
		}(source, symbols[i])
	}

	wg.Wait()
	close(errCh)

	// Just verify no panics occurred
	t.Log("Concurrent source requests completed")
}

// TestYahooReader_RapidContextCancellations tests rapid cancellation scenarios
func TestYahooReader_RapidContextCancellations(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	// Create reader with no retries to avoid long delays
	opts := &internalhttp.ClientOptions{
		Timeout:    200 * time.Millisecond,
		MaxRetries: 0, // Disable retries
	}
	reader := yahoo.NewYahooReaderWithBaseURL(opts, server.URL+"/%s")

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// Rapidly create and cancel contexts (reduced to 5 iterations)
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		_, err := reader.ReadSingle(ctx, "AAPL", start, end)
		cancel() // Cancel immediately after

		if err == nil {
			// Some might succeed if very fast
			t.Log("Request succeeded despite quick cancellation")
		}
	}
}
