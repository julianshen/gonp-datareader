package stooq_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/stooq"
)

func TestNewStooqReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	reader := stooq.NewStooqReader(opts)

	if reader == nil {
		t.Fatal("NewStooqReader returned nil")
	}

	if reader.Name() != "stooq" {
		t.Errorf("Expected name 'stooq', got %q", reader.Name())
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		wantParts []string
	}{
		{
			name:   "US stock",
			symbol: "AAPL.US",
			wantParts: []string{
				"stooq.com",
				"s=AAPL.US",
				"i=d",
			},
		},
		{
			name:   "index",
			symbol: "^SPX",
			wantParts: []string{
				"stooq.com",
				"s=%5ESPX", // URL encoded ^
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := stooq.BuildURL(tt.symbol)

			for _, part := range tt.wantParts {
				if !strings.Contains(url, part) {
					t.Errorf("BuildURL() missing part %q, got %q", part, url)
				}
			}

			if !strings.HasPrefix(url, "https://") {
				t.Errorf("BuildURL() should use HTTPS, got %q", url)
			}
		})
	}
}

// TestStooqReader_ValidateSymbol tests symbol validation
func TestStooqReader_ValidateSymbol(t *testing.T) {
	reader := stooq.NewStooqReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid US stock",
			symbol:  "AAPL.US",
			wantErr: false,
		},
		{
			name:    "valid symbol with dash",
			symbol:  "BRK-B",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "AA PL",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := reader.ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol(%q) error = %v, wantErr %v", tt.symbol, err, tt.wantErr)
			}
		})
	}
}

// TestStooqReader_ReadSingle_WithMockServer tests ReadSingle with mock HTTP server
func TestStooqReader_ReadSingle_WithMockServer(t *testing.T) {
	// Sample Stooq CSV response
	csvData := `Date,Open,High,Low,Close,Volume
2023-01-05,130.00,135.00,129.00,134.50,75000000
2023-01-04,128.00,132.00,127.50,130.00,70000000`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	// Create reader with custom base URL
	reader := stooq.NewStooqReaderWithBaseURL(nil, server.URL+"?s=%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "AAPL.US", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	data, ok := result.(*stooq.ParsedData)
	if !ok {
		t.Fatalf("Expected *stooq.ParsedData, got %T", result)
	}

	if len(data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(data.Rows))
	}
}

// TestStooqReader_ReadSingle_InvalidSymbol tests error handling for invalid symbols
func TestStooqReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := stooq.NewStooqReader(nil)

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

// TestStooqReader_Read_MultipleSymbols tests fetching multiple symbols
func TestStooqReader_Read_MultipleSymbols(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Volume
2023-01-05,100.00,105.00,99.00,104.00,10000000`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	reader := stooq.NewStooqReaderWithBaseURL(nil, server.URL+"?s=%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL.US", "MSFT.US", "GOOGL.US"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*stooq.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*stooq.ParsedData, got %T", result)
	}

	if len(dataMap) != len(symbols) {
		t.Errorf("Expected %d results, got %d", len(symbols), len(dataMap))
	}

	for _, symbol := range symbols {
		if _, exists := dataMap[symbol]; !exists {
			t.Errorf("Missing data for symbol %s", symbol)
		}
	}
}

// TestStooqReader_Read_InvalidDateRange tests error handling for invalid date ranges
func TestStooqReader_Read_InvalidDateRange(t *testing.T) {
	reader := stooq.NewStooqReader(nil)

	ctx := context.Background()
	start := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"AAPL.US"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

// TestStooqReader_HTTPError tests handling of HTTP errors
func TestStooqReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := stooq.NewStooqReaderWithBaseURL(nil, server.URL+"?s=%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL.US", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}
}
