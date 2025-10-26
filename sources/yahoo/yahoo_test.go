package yahoo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

func TestNewYahooReader(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	if reader == nil {
		t.Fatal("NewYahooReader() returned nil")
	}

	if reader.Name() != "Yahoo Finance" {
		t.Errorf("Expected name 'Yahoo Finance', got '%s'", reader.Name())
	}

	if reader.Source() != "yahoo" {
		t.Errorf("Expected source 'yahoo', got '%s'", reader.Source())
	}
}

func TestYahooReader_ValidateSymbol(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid US stock symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "valid symbol with dot",
			symbol:  "BRK.B",
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

func TestYahooReader_BuildURL(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("AAPL", start, end)

	if url == "" {
		t.Error("BuildURL() returned empty string")
	}

	// URL should contain the symbol
	if !contains(url, "AAPL") {
		t.Errorf("URL should contain symbol 'AAPL': %s", url)
	}

	// URL should contain Yahoo Finance domain
	if !contains(url, "query1.finance.yahoo.com") && !contains(url, "query2.finance.yahoo.com") {
		t.Errorf("URL should contain Yahoo Finance domain: %s", url)
	}

	// URL should contain timestamp parameters
	if !contains(url, "period1=") || !contains(url, "period2=") {
		t.Errorf("URL should contain period parameters: %s", url)
	}
}

func TestYahooReader_BuildURL_Timestamps(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)

	url := reader.BuildURL("MSFT", start, end)

	// Verify Unix timestamps are in the URL
	startTimestamp := start.Unix()
	endTimestamp := end.Unix()

	startStr := "period1=" + itoa(startTimestamp)
	endStr := "period2=" + itoa(endTimestamp)

	if !contains(url, startStr) {
		t.Errorf("URL should contain start timestamp %s: %s", startStr, url)
	}

	if !contains(url, endStr) {
		t.Errorf("URL should contain end timestamp %s: %s", endStr, url)
	}
}

func TestYahooReader_Read(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// This test verifies the interface, actual fetching will be tested separately
	_, err := reader.Read(ctx, []string{"AAPL"}, start, end)

	// We expect this to work (or fail gracefully) - we're just testing the interface exists
	if err != nil {
		// Error is acceptable for now (network issues, etc.)
		t.Logf("Read() returned error (expected in unit test): %v", err)
	}
}

func TestYahooReader_ReadSingle(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// This test verifies the interface exists
	_, err := reader.ReadSingle(ctx, "AAPL", start, end)

	// Error is acceptable for now
	if err != nil {
		t.Logf("ReadSingle() returned error (expected in unit test): %v", err)
	}
}

func TestYahooReader_ReadSingle_WithMockServer(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.239990,300.600006,295.190002,300.350006,297.450287,33911900
2020-01-03,297.149994,300.579987,296.500000,297.429993,294.558075,36607600
2020-01-06,293.790009,299.959991,292.750000,299.799988,296.906128,29596800`

	// Create mock server
	server := createMockYahooServer(csvData)
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	// Verify we got parsed data
	data, ok := result.(*yahoo.ParsedData)
	if !ok {
		t.Fatalf("Expected *yahoo.ParsedData, got %T", result)
	}

	if len(data.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(data.Rows))
	}

	// Verify column names
	expectedCols := []string{"Date", "Open", "High", "Low", "Close", "Adj Close", "Volume"}
	if len(data.Columns) != len(expectedCols) {
		t.Errorf("Expected %d columns, got %d", len(expectedCols), len(data.Columns))
	}
}

func TestYahooReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

func TestYahooReader_ReadSingle_InvalidDateRange(t *testing.T) {
	reader := yahoo.NewYahooReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for invalid date range")
	}
}

func TestYahooReader_Read_MultipleSymbols(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	server := createMockYahooServer(csvData)
	defer server.Close()

	reader := yahoo.NewYahooReaderWithBaseURL(nil, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	results, err := reader.Read(ctx, []string{"AAPL", "MSFT"}, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if results == nil {
		t.Fatal("Read() returned nil")
	}

	// Should return a map of symbol to data
	dataMap, ok := results.(map[string]*yahoo.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*yahoo.ParsedData, got %T", results)
	}

	if len(dataMap) != 2 {
		t.Errorf("Expected 2 results, got %d", len(dataMap))
	}

	if _, ok := dataMap["AAPL"]; !ok {
		t.Error("Missing AAPL data")
	}

	if _, ok := dataMap["MSFT"]; !ok {
		t.Error("Missing MSFT data")
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var buf [20]byte
	i := len(buf) - 1

	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}

	if negative {
		buf[i] = '-'
		i--
	}

	return string(buf[i+1:])
}

func createMockYahooServer(csvData string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
}
