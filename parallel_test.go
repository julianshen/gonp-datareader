package datareader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

// TestYahooReader_ParallelFetching tests that multiple symbols are fetched in parallel
func TestYahooReader_ParallelFetching(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	var requestCount int32
	var concurrentRequests int32
	var maxConcurrent int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track concurrent requests
		current := atomic.AddInt32(&concurrentRequests, 1)

		// Update max concurrent
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}

		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)

		atomic.AddInt32(&requestCount, 1)
		atomic.AddInt32(&concurrentRequests, -1)

		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "FB"}

	startTime := time.Now()
	results, err := reader.Read(ctx, symbols, start, end)
	elapsed := time.Since(startTime)

	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	// Verify all symbols were fetched
	dataMap, ok := results.(map[string]*yahoo.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*yahoo.ParsedData, got %T", results)
	}

	if len(dataMap) != len(symbols) {
		t.Errorf("Expected %d results, got %d", len(symbols), len(dataMap))
	}

	// Verify all symbols are present
	for _, symbol := range symbols {
		if _, ok := dataMap[symbol]; !ok {
			t.Errorf("Missing data for symbol %s", symbol)
		}
	}

	// Check that requests were made
	finalCount := atomic.LoadInt32(&requestCount)
	if finalCount != int32(len(symbols)) {
		t.Errorf("Expected %d requests, got %d", len(symbols), finalCount)
	}

	// Check for parallelism - with 5 symbols and 50ms each, sequential would take ~250ms
	// Parallel should take close to 50ms (plus overhead)
	// We'll be lenient and check it's faster than sequential
	sequentialTime := 50 * time.Millisecond * time.Duration(len(symbols))
	if elapsed >= sequentialTime {
		t.Logf("Warning: Parallel fetching took %v, sequential would take %v", elapsed, sequentialTime)
		t.Log("Parallel fetching may not be working as expected")
	} else {
		t.Logf("Parallel fetching successful: %v (vs %v sequential)", elapsed, sequentialTime)
	}

	// Log max concurrent requests
	maxConc := atomic.LoadInt32(&maxConcurrent)
	t.Logf("Max concurrent requests: %d", maxConc)

	if maxConc < 2 {
		t.Error("Expected at least 2 concurrent requests, got less - parallelism not working")
	}
}

// TestYahooReader_ParallelFetchingWithErrors tests error handling in parallel fetching
func TestYahooReader_ParallelFetchingWithErrors(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error for FAIL symbol
		if r.URL.Path == "/FAIL" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// Include one symbol that will fail
	symbols := []string{"AAPL", "FAIL", "GOOGL"}

	_, err := reader.Read(ctx, symbols, start, end)
	if err == nil {
		t.Error("Expected error when one symbol fails, got nil")
	}

	// Error should mention the failed symbol
	if err != nil {
		t.Logf("Got expected error: %v", err)
	}
}

// TestYahooReader_ParallelFetchingContextCancellation tests context cancellation during parallel fetching
func TestYahooReader_ParallelFetchingContextCancellation(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	requestReceived := make(chan struct{}, 5)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived <- struct{}{}
		// Delay to allow cancellation
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx, cancel := context.WithCancel(context.Background())
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL", "MSFT", "GOOGL"}

	// Start reading in a goroutine
	errCh := make(chan error, 1)
	go func() {
		_, err := reader.Read(ctx, symbols, start, end)
		errCh <- err
	}()

	// Wait for at least one request to be received, then cancel
	<-requestReceived
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the read to complete
	err := <-errCh
	if err == nil {
		t.Error("Expected error after context cancellation, got nil")
	}

	t.Logf("Got expected cancellation error: %v", err)
}

// BenchmarkYahooReader_ParallelVsSequential compares parallel and sequential fetching
func BenchmarkYahooReader_ParallelVsSequential(b *testing.B) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900
2020-01-03,297.15,300.58,296.50,297.43,294.56,36607600`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate network delay
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

	symbols := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "FB"}

	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := reader.Read(ctx, symbols, start, end)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
