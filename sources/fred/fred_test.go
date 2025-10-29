package fred_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
	"github.com/julianshen/gonp-datareader/sources/fred"
)

func TestNewFREDReader(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	if reader == nil {
		t.Fatal("NewFREDReader returned nil")
	}
}

func TestFREDReader_Name(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	name := reader.Name()
	expected := "FRED"

	if name != expected {
		t.Errorf("Expected name %q, got %q", expected, name)
	}
}

func TestFREDReader_Source(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	source := reader.Source()
	expected := "fred"

	if source != expected {
		t.Errorf("Expected source %q, got %q", expected, source)
	}
}

func TestFREDReader_ImplementsInterface(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	// Verify it implements Reader interface
	var _ sources.Reader = reader
}

func TestFREDReader_ValidateSymbol(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid series ID",
			symbol:  "GDP",
			wantErr: false,
		},
		{
			name:    "valid series with numbers",
			symbol:  "DGS10",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "GD P",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := reader.ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFREDReader_BuildURL(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("GDP", start, end, "test_api_key")

	// Should contain base URL
	expectedBase := "https://api.stlouisfed.org/fred/series/observations"
	if len(url) < len(expectedBase) || url[:len(expectedBase)] != expectedBase {
		t.Errorf("URL should start with %q, got %q", expectedBase, url)
	}

	// Should contain series ID
	if !contains(url, "series_id=GDP") {
		t.Errorf("URL should contain series_id=GDP, got %q", url)
	}

	// Should contain API key
	if !contains(url, "api_key=test_api_key") {
		t.Errorf("URL should contain api_key parameter, got %q", url)
	}

	// Should contain dates
	if !contains(url, "observation_start=2020-01-01") {
		t.Errorf("URL should contain observation_start, got %q", url)
	}

	if !contains(url, "observation_end=2020-12-31") {
		t.Errorf("URL should contain observation_end, got %q", url)
	}

	// Should contain JSON format
	if !contains(url, "file_type=json") {
		t.Errorf("URL should contain file_type=json, got %q", url)
	}
}

func TestFREDReader_ReadSingle_RequiresAPIKey(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	// Should fail without API key configured
	_, err := reader.ReadSingle(ctx, "GDP", start, end)
	if err == nil {
		t.Error("Expected error when API key not configured, got nil")
	}
}

func TestFREDReader_SetAPIKey(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReader(opts)

	reader.SetAPIKey("test_key")

	if reader.GetAPIKey() != "test_key" {
		t.Errorf("Expected API key 'test_key', got '%s'", reader.GetAPIKey())
	}
}

func TestFREDReader_Read_InvalidSymbols(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReaderWithAPIKey(opts, "test_key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.Read(ctx, []string{}, start, end)
	if err == nil {
		t.Error("Expected error for empty symbols list, got nil")
	}
}

func TestFREDReader_ReadSingle_InvalidDateRange(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := fred.NewFREDReaderWithAPIKey(opts, "test_key")

	ctx := context.Background()
	start := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "GDP", start, end)
	if err == nil {
		t.Error("Expected error for invalid date range, got nil")
	}
}

func TestFREDReader_ReadSingle_WithMockServer(t *testing.T) {
	jsonData := `{
		"observations": [
			{
				"realtime_start": "2020-01-01",
				"realtime_end": "2020-01-01",
				"date": "2020-01-01",
				"value": "21427.91"
			},
			{
				"realtime_start": "2020-01-02",
				"realtime_end": "2020-01-02",
				"date": "2020-01-02",
				"value": "21433.22"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := fred.NewFREDReaderWithBaseURL(nil, server.URL)
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "GDP", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	data, ok := result.(*fred.ParsedData)
	if !ok {
		t.Fatalf("Expected *fred.ParsedData, got %T", result)
	}

	if len(data.Dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(data.Dates))
	}

	if len(data.Values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(data.Values))
	}
}

func TestFREDReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := fred.NewFREDReaderWithAPIKey(nil, "test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

func TestFREDReader_Read_MultipleSymbols(t *testing.T) {
	jsonData := `{
		"observations": [
			{
				"realtime_start": "2020-01-01",
				"realtime_end": "2020-01-01",
				"date": "2020-01-01",
				"value": "21427.91"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := fred.NewFREDReaderWithBaseURL(nil, server.URL)
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"GDP", "DGS10", "UNRATE"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*fred.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*fred.ParsedData, got %T", result)
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

func TestFREDReader_Read_InvalidDateRange(t *testing.T) {
	reader := fred.NewFREDReaderWithAPIKey(nil, "test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"GDP"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

func TestFREDReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := fred.NewFREDReaderWithBaseURL(nil, server.URL)
	reader.SetAPIKey("test-api-key")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "GDP", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}

	if !contains(err.Error(), "status 500") {
		t.Errorf("Expected error to mention status 500, got: %v", err)
	}
}

func TestFREDReader_Read_RequiresAPIKey(t *testing.T) {
	reader := fred.NewFREDReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.Read(ctx, []string{"GDP"}, start, end)
	if err == nil {
		t.Error("Read() should return error when API key not set")
	}

	if !contains(err.Error(), "API key") {
		t.Errorf("Expected error to mention API key, got: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
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
