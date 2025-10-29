package alphavantage_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/alphavantage"
)

func TestNewAlphaVantageReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}
	apiKey := "test_api_key"

	reader := alphavantage.NewAlphaVantageReader(opts, apiKey)

	if reader == nil {
		t.Fatal("NewAlphaVantageReader returned nil")
	}

	if reader.Name() != "alphavantage" {
		t.Errorf("Expected name 'alphavantage', got %q", reader.Name())
	}
}

func TestNewAlphaVantageReader_RequiresAPIKey(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	// Empty API key should still create reader (validation happens at request time)
	reader := alphavantage.NewAlphaVantageReader(opts, "")

	if reader == nil {
		t.Fatal("NewAlphaVantageReader returned nil")
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		apiKey    string
		wantParts []string
	}{
		{
			name:   "daily time series",
			symbol: "AAPL",
			apiKey: "demo",
			wantParts: []string{
				"alphavantage.co",
				"function=TIME_SERIES_DAILY",
				"symbol=AAPL",
				"apikey=demo",
				"outputsize=full",
			},
		},
		{
			name:   "different symbol",
			symbol: "MSFT",
			apiKey: "test_key",
			wantParts: []string{
				"symbol=MSFT",
				"apikey=test_key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := alphavantage.BuildURL(tt.symbol, tt.apiKey)

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

// TestAlphaVantageReader_ValidateSymbol tests symbol validation
func TestAlphaVantageReader_ValidateSymbol(t *testing.T) {
	reader := alphavantage.NewAlphaVantageReader(nil, "test_key")

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "valid symbol with dot",
			symbol:  "BRK.B",
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

// TestAlphaVantageReader_ReadSingle_WithMockServer tests ReadSingle with mock HTTP server
func TestAlphaVantageReader_ReadSingle_WithMockServer(t *testing.T) {
	// Sample Alpha Vantage response
	jsonResponse := `{
		"Meta Data": {
			"1. Information": "Daily Prices",
			"2. Symbol": "AAPL",
			"3. Last Refreshed": "2023-01-05",
			"4. Output Size": "Full size",
			"5. Time Zone": "US/Eastern"
		},
		"Time Series (Daily)": {
			"2023-01-05": {
				"1. open": "130.00",
				"2. high": "135.00",
				"3. low": "129.00",
				"4. close": "134.50",
				"5. volume": "75000000"
			},
			"2023-01-04": {
				"1. open": "128.00",
				"2. high": "132.00",
				"3. low": "127.50",
				"4. close": "130.00",
				"5. volume": "70000000"
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	// Create reader with custom base URL
	reader := alphavantage.NewAlphaVantageReaderWithBaseURL(nil, "test_key", server.URL+"?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	data, ok := result.(*alphavantage.ParsedData)
	if !ok {
		t.Fatalf("Expected *alphavantage.ParsedData, got %T", result)
	}

	if len(data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(data.Rows))
	}
}

// TestAlphaVantageReader_ReadSingle_InvalidSymbol tests error handling for invalid symbols
func TestAlphaVantageReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := alphavantage.NewAlphaVantageReader(nil, "test_key")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

// TestAlphaVantageReader_ReadSingle_NoAPIKey tests error when API key is missing
func TestAlphaVantageReader_ReadSingle_NoAPIKey(t *testing.T) {
	reader := alphavantage.NewAlphaVantageReader(nil, "")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error when API key is missing")
	}

	if !strings.Contains(err.Error(), "API key is required") {
		t.Errorf("Expected 'API key is required' error, got: %v", err)
	}
}

// TestAlphaVantageReader_Read_MultipleSymbols tests fetching multiple symbols
func TestAlphaVantageReader_Read_MultipleSymbols(t *testing.T) {
	jsonResponse := `{
		"Meta Data": {
			"2. Symbol": "TEST"
		},
		"Time Series (Daily)": {
			"2023-01-05": {
				"1. open": "100.00",
				"2. high": "105.00",
				"3. low": "99.00",
				"4. close": "104.00",
				"5. volume": "10000000"
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	reader := alphavantage.NewAlphaVantageReaderWithBaseURL(nil, "test_key", server.URL+"?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*alphavantage.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*alphavantage.ParsedData, got %T", result)
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

// TestAlphaVantageReader_Read_InvalidDateRange tests error handling for invalid date ranges
func TestAlphaVantageReader_Read_InvalidDateRange(t *testing.T) {
	reader := alphavantage.NewAlphaVantageReader(nil, "test_key")

	ctx := context.Background()
	start := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"AAPL"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

// TestAlphaVantageReader_HTTPError tests handling of HTTP errors
func TestAlphaVantageReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := alphavantage.NewAlphaVantageReaderWithBaseURL(nil, "test_key", server.URL+"?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}
}
