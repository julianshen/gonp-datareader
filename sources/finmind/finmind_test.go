package finmind_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/finmind"
)

func TestNewFinMindReader(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	if reader == nil {
		t.Fatal("NewFinMindReader() returned nil")
	}

	if reader.Name() != "FinMind" {
		t.Errorf("Expected name 'FinMind', got %q", reader.Name())
	}

	if reader.Source() != "finmind" {
		t.Errorf("Expected source 'finmind', got %q", reader.Source())
	}
}

func TestNewFinMindReaderWithToken(t *testing.T) {
	token := "test-token-123"
	reader := finmind.NewFinMindReaderWithToken(nil, token)

	if reader == nil {
		t.Fatal("NewFinMindReaderWithToken() returned nil")
	}

	if reader.Name() != "FinMind" {
		t.Errorf("Expected name 'FinMind', got %q", reader.Name())
	}

	if reader.Source() != "finmind" {
		t.Errorf("Expected source 'finmind', got %q", reader.Source())
	}
}

func TestFinMindReader_ValidateSymbol(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid Taiwan 4-digit stock code",
			symbol:  "2330",
			wantErr: false,
		},
		{
			name:    "valid Taiwan 6-digit warrant code",
			symbol:  "123456",
			wantErr: false,
		},
		{
			name:    "valid US stock symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			symbol:  "   ",
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

func TestFinMindReader_BuildURL(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("2330", start, end)

	// Check that URL contains required parameters
	expectedParams := []string{
		"dataset=TaiwanStockPrice",
		"data_id=2330",
		"start_date=2020-04-02",
		"end_date=2020-04-12",
	}

	for _, param := range expectedParams {
		if !contains(url, param) {
			t.Errorf("BuildURL() missing parameter: %s\nGot: %s", param, url)
		}
	}

	// Check base URL
	if !contains(url, "https://api.finmindtrade.com/api/v4/data") {
		t.Errorf("BuildURL() missing base URL\nGot: %s", url)
	}
}

func TestFinMindReader_BuildURL_CustomDataset(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)
	reader.SetDataset("TaiwanStockDividend")

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("2330", start, end)

	if !contains(url, "dataset=TaiwanStockDividend") {
		t.Errorf("BuildURL() should use custom dataset\nGot: %s", url)
	}
}

func TestFinMindReader_BuildURL_USStock(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)
	reader.SetDataset("USStockPrice")

	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("AAPL", start, end)

	expectedParams := []string{
		"dataset=USStockPrice",
		"data_id=AAPL",
		"start_date=2021-01-01",
		"end_date=2021-12-31",
	}

	for _, param := range expectedParams {
		if !contains(url, param) {
			t.Errorf("BuildURL() missing parameter: %s\nGot: %s", param, url)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestFinMindReader_ReadSingle(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		if query.Get("dataset") != "TaiwanStockPrice" {
			t.Errorf("Expected dataset=TaiwanStockPrice, got %s", query.Get("dataset"))
		}
		if query.Get("data_id") != "2330" {
			t.Errorf("Expected data_id=2330, got %s", query.Get("data_id"))
		}

		// Return mock response
		mockData := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"date":             "2020-04-06",
					"stock_id":         "2330",
					"Trading_Volume":   59712754,
					"Trading_money":    16324198154,
					"open":             273.0,
					"max":              275.5,
					"min":              270.0,
					"close":            275.5,
					"spread":           4.0,
					"Trading_turnover": 19971,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// Create reader with custom endpoint
	reader := finmind.NewFinMindReaderWithEndpoint(nil, server.URL)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "2330", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	data := result.(*finmind.ParsedData)

	if data.Symbol != "2330" {
		t.Errorf("Expected symbol '2330', got %q", data.Symbol)
	}

	if len(data.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(data.Rows))
	}

	row := data.Rows[0]
	if row["close"] != "275.5" {
		t.Errorf("Expected close '275.5', got %q", row["close"])
	}
}

func TestFinMindReader_ReadSingle_WithToken(t *testing.T) {
	// Create mock server that checks Authorization header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer test-token-123"
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuth, authHeader)
		}

		// Return mock response
		mockData := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"date":             "2020-04-06",
					"stock_id":         "2330",
					"Trading_Volume":   59712754,
					"Trading_money":    16324198154,
					"open":             273.0,
					"max":              275.5,
					"min":              270.0,
					"close":            275.5,
					"spread":           4.0,
					"Trading_turnover": 19971,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// Create reader with token and custom endpoint
	opts := &internalhttp.ClientOptions{
		Timeout: 10 * time.Second,
	}
	reader := finmind.NewFinMindReaderWithTokenAndEndpoint(opts, "test-token-123", server.URL)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "2330", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	data := result.(*finmind.ParsedData)
	if data.Symbol != "2330" {
		t.Errorf("Expected symbol '2330', got %q", data.Symbol)
	}
}

func TestFinMindReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should error on empty symbol")
	}
}

func TestFinMindReader_ReadSingle_InvalidDateRange(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "2330", start, end)
	if err == nil {
		t.Error("ReadSingle() should error on invalid date range")
	}
}

func TestFinMindReader_ReadSingle_HTTPError(t *testing.T) {
	// Create mock server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := finmind.NewFinMindReaderWithEndpoint(nil, server.URL)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "2330", start, end)
	if err == nil {
		t.Error("ReadSingle() should error on HTTP 500")
	}
}

func TestFinMindReader_Read_MultipleSymbols(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("data_id")

		mockData := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"date":             "2020-04-06",
					"stock_id":         symbol,
					"Trading_Volume":   59712754,
					"Trading_money":    16324198154,
					"open":             273.0,
					"max":              275.5,
					"min":              270.0,
					"close":            275.5,
					"spread":           4.0,
					"Trading_turnover": 19971,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	reader := finmind.NewFinMindReaderWithEndpoint(nil, server.URL)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	symbols := []string{"2330", "2317", "2454"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap := result.(map[string]*finmind.ParsedData)

	if len(dataMap) != 3 {
		t.Fatalf("Expected 3 symbols, got %d", len(dataMap))
	}

	for _, symbol := range symbols {
		data, ok := dataMap[symbol]
		if !ok {
			t.Errorf("Missing data for symbol %s", symbol)
			continue
		}
		if data.Symbol != symbol {
			t.Errorf("Expected symbol %s, got %s", symbol, data.Symbol)
		}
	}
}

func TestFinMindReader_Read_EmptySymbols(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC)

	result, err := reader.Read(ctx, []string{}, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap := result.(map[string]*finmind.ParsedData)
	if len(dataMap) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(dataMap))
	}
}
