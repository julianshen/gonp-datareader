package tiingo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/tiingo"
)

func TestNewTiingoReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{}
	reader := tiingo.NewTiingoReader(opts)

	if reader == nil {
		t.Fatal("NewTiingoReader() returned nil")
	}

	if reader.Name() != "Tiingo" {
		t.Errorf("Expected name 'Tiingo', got '%s'", reader.Name())
	}

	if reader.Source() != "tiingo" {
		t.Errorf("Expected source 'tiingo', got '%s'", reader.Source())
	}
}

func TestTiingoReader_ValidateSymbol(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid stock symbol",
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

func TestTiingoReader_BuildURL(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("AAPL", start, end, "test-api-key")

	if url == "" {
		t.Error("BuildURL() returned empty string")
	}

	// URL should contain the symbol
	if !contains(url, "AAPL") {
		t.Errorf("URL should contain symbol 'AAPL': %s", url)
	}

	// URL should contain Tiingo domain
	if !contains(url, "api.tiingo.com") {
		t.Errorf("URL should contain Tiingo domain: %s", url)
	}

	// URL should contain date parameters
	if !contains(url, "startDate=") || !contains(url, "endDate=") {
		t.Errorf("URL should contain date parameters: %s", url)
	}

	// URL should contain API token
	if !contains(url, "token=test-api-key") {
		t.Errorf("URL should contain API token: %s", url)
	}
}

func TestTiingoReader_ReadSingle_WithMockServer(t *testing.T) {
	jsonData := `[
		{
			"date": "2020-01-02T00:00:00.000Z",
			"close": 300.35,
			"high": 300.60,
			"low": 295.19,
			"open": 296.24,
			"volume": 33911900,
			"adjClose": 297.45,
			"adjHigh": 300.60,
			"adjLow": 295.19,
			"adjOpen": 296.24,
			"adjVolume": 33911900,
			"divCash": 0.0,
			"splitFactor": 1.0
		},
		{
			"date": "2020-01-03T00:00:00.000Z",
			"close": 297.43,
			"high": 300.58,
			"low": 296.50,
			"open": 297.15,
			"volume": 36607600,
			"adjClose": 294.56,
			"adjHigh": 300.58,
			"adjLow": 296.50,
			"adjOpen": 297.15,
			"adjVolume": 36607600,
			"divCash": 0.0,
			"splitFactor": 1.0
		}
	]`

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := tiingo.NewTiingoReaderWithBaseURL(nil, server.URL+"/tiingo/daily/%s/prices")
	reader.SetAPIKey("test-api-key") // Set API key for testing

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
	data, ok := result.(*tiingo.ParsedData)
	if !ok {
		t.Fatalf("Expected *tiingo.ParsedData, got %T", result)
	}

	if len(data.Dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(data.Dates))
	}

	if len(data.Prices) != 2 {
		t.Errorf("Expected 2 prices, got %d", len(data.Prices))
	}
}

func TestTiingoReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

func TestTiingoReader_ReadSingle_InvalidDateRange(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for invalid date range")
	}
}

func TestTiingoReader_RequiresAPIKey(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("Expected error when API key is missing")
	}

	// Error should mention API key
	if err != nil && !contains(err.Error(), "API key") {
		t.Errorf("Expected error to mention API key, got: %v", err)
	}
}

// TestTiingoReader_Read_MultipleSymbols tests parallel fetching of multiple symbols
func TestTiingoReader_Read_MultipleSymbols(t *testing.T) {
	jsonData := `[
		{
			"date": "2020-01-02T00:00:00.000Z",
			"close": 300.35,
			"high": 300.60,
			"low": 295.19,
			"open": 296.24,
			"volume": 33911900,
			"adjClose": 297.45,
			"adjHigh": 300.60,
			"adjLow": 295.19,
			"adjOpen": 296.24,
			"adjVolume": 33911900,
			"divCash": 0.0,
			"splitFactor": 1.0
		}
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := tiingo.NewTiingoReaderWithBaseURL(nil, server.URL+"/tiingo/daily/%s/prices")
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*tiingo.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*tiingo.ParsedData, got %T", result)
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

// TestTiingoReader_Read_InvalidDateRange tests error handling for invalid date ranges
func TestTiingoReader_Read_InvalidDateRange(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"AAPL"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

// TestTiingoReader_HTTPError tests handling of HTTP errors
func TestTiingoReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := tiingo.NewTiingoReaderWithBaseURL(nil, server.URL+"/tiingo/daily/%s/prices")
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}

	if !contains(err.Error(), "status 500") {
		t.Errorf("Expected error to mention status 500, got: %v", err)
	}
}

// TestTiingoReader_SetAPIKey tests the API key setter
func TestTiingoReader_SetAPIKey(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)

	// Initially no API key
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err == nil {
		t.Error("Expected error when API key is not set")
	}

	// Set API key and test it's used
	reader.SetAPIKey("test-key")

	// Create a mock server to verify the key is passed
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != "test-key" {
			t.Errorf("Expected token 'test-key', got '%s'", token)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"date": "2020-01-02T00:00:00.000Z", "close": 300.35, "high": 300.60, "low": 295.19, "open": 296.24, "volume": 33911900}]`))
	}))
	defer server.Close()

	reader2 := tiingo.NewTiingoReaderWithBaseURL(nil, server.URL+"/tiingo/daily/%s/prices")
	reader2.SetAPIKey("test-key")

	_, err = reader2.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Errorf("ReadSingle() should succeed with valid API key: %v", err)
	}
}

// TestTiingoReader_APIKeyFromContext tests getting API key from context
func TestTiingoReader_APIKeyFromContext(t *testing.T) {
	jsonData := `[{"date": "2020-01-02T00:00:00.000Z", "close": 300.35, "high": 300.60, "low": 295.19, "open": 296.24, "volume": 33911900}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != "context-key" {
			t.Errorf("Expected token 'context-key' from context, got '%s'", token)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := tiingo.NewTiingoReaderWithBaseURL(nil, server.URL+"/tiingo/daily/%s/prices")

	// Set API key via context
	ctx := context.WithValue(context.Background(), "apiKey", "context-key")
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Errorf("ReadSingle() should succeed with API key from context: %v", err)
	}
}

// TestTiingoReader_Read_EmptySymbols tests error handling for empty symbol list
func TestTiingoReader_Read_EmptySymbols(t *testing.T) {
	reader := tiingo.NewTiingoReader(nil)
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.Read(ctx, []string{}, start, end)
	if err == nil {
		t.Error("Read() should return error for empty symbols list")
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
