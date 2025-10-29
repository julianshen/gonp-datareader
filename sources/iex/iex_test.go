package iex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/iex"
)

func TestNewIEXReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}
	apiKey := "test_api_key"

	reader := iex.NewIEXReader(opts, apiKey)

	if reader == nil {
		t.Fatal("NewIEXReader returned nil")
	}

	if reader.Name() != "iex" {
		t.Errorf("Expected name 'iex', got %q", reader.Name())
	}
}

func TestNewIEXReader_RequiresAPIKey(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	// Empty API key should still create reader (validation happens at request time)
	reader := iex.NewIEXReader(opts, "")

	if reader == nil {
		t.Fatal("NewIEXReader returned nil")
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		dateRange string
		apiKey    string
		wantParts []string
	}{
		{
			name:      "6m range",
			symbol:    "AAPL",
			dateRange: "6m",
			apiKey:    "test_token",
			wantParts: []string{
				"cloud.iexapis.com",
				"/stable/stock/AAPL/chart/6m",
				"token=test_token",
			},
		},
		{
			name:      "1y range",
			symbol:    "MSFT",
			dateRange: "1y",
			apiKey:    "demo",
			wantParts: []string{
				"/stable/stock/MSFT/chart/1y",
				"token=demo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := iex.BuildURL(tt.symbol, tt.dateRange, tt.apiKey)

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

func TestCalculateDateRange(t *testing.T) {
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		want  string
	}{
		{
			name:  "less than 1 month",
			start: time.Now().AddDate(0, 0, -20),
			end:   time.Now(),
			want:  "1m",
		},
		{
			name:  "3 months",
			start: time.Now().AddDate(0, -3, 0),
			end:   time.Now(),
			want:  "3m",
		},
		{
			name:  "1 year",
			start: time.Now().AddDate(-1, 0, 0),
			end:   time.Now(),
			want:  "1y",
		},
		{
			name:  "5 years (max)",
			start: time.Now().AddDate(-5, 0, 0),
			end:   time.Now(),
			want:  "5y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iex.CalculateDateRange(tt.start, tt.end)
			if got != tt.want {
				t.Errorf("CalculateDateRange() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestIEXReader_ValidateSymbol tests symbol validation
func TestIEXReader_ValidateSymbol(t *testing.T) {
	reader := iex.NewIEXReader(nil, "test_key")

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

// TestIEXReader_ReadSingle_WithMockServer tests ReadSingle with mock HTTP server
func TestIEXReader_ReadSingle_WithMockServer(t *testing.T) {
	// Sample IEX Cloud response
	jsonResponse := `[
		{
			"date": "2023-01-05",
			"open": 130.00,
			"high": 135.00,
			"low": 129.00,
			"close": 134.50,
			"volume": 75000000
		},
		{
			"date": "2023-01-04",
			"open": 128.00,
			"high": 132.00,
			"low": 127.50,
			"close": 130.00,
			"volume": 70000000
		}
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	// Create reader with custom base URL
	reader := iex.NewIEXReaderWithBaseURL(nil, "test_key", server.URL+"?symbol=%s&range=%s&token=%s")

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

	data, ok := result.(*iex.ParsedData)
	if !ok {
		t.Fatalf("Expected *iex.ParsedData, got %T", result)
	}

	if len(data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(data.Rows))
	}
}

// TestIEXReader_ReadSingle_InvalidSymbol tests error handling for invalid symbols
func TestIEXReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := iex.NewIEXReader(nil, "test_key")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

// TestIEXReader_ReadSingle_NoAPIKey tests error when API key is missing
func TestIEXReader_ReadSingle_NoAPIKey(t *testing.T) {
	reader := iex.NewIEXReader(nil, "")

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

// TestIEXReader_Read_MultipleSymbols tests fetching multiple symbols
func TestIEXReader_Read_MultipleSymbols(t *testing.T) {
	jsonResponse := `[
		{
			"date": "2023-01-05",
			"open": 100.00,
			"high": 105.00,
			"low": 99.00,
			"close": 104.00,
			"volume": 10000000
		}
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	reader := iex.NewIEXReaderWithBaseURL(nil, "test_key", server.URL+"?symbol=%s&range=%s&token=%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*iex.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*iex.ParsedData, got %T", result)
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

// TestIEXReader_Read_InvalidDateRange tests error handling for invalid date ranges
func TestIEXReader_Read_InvalidDateRange(t *testing.T) {
	reader := iex.NewIEXReader(nil, "test_key")

	ctx := context.Background()
	start := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"AAPL"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

// TestIEXReader_HTTPError tests handling of HTTP errors
func TestIEXReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := iex.NewIEXReaderWithBaseURL(nil, "test_key", server.URL+"?symbol=%s&range=%s&token=%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}
}
