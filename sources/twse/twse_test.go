package twse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

// TestTWSEReader_Struct tests that TWSEReader struct exists and has correct structure
func TestTWSEReader_Struct(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := NewTWSEReader(opts)

	if reader == nil {
		t.Fatal("NewTWSEReader returned nil")
	}

	// Test that it's actually a TWSEReader
	_, ok := interface{}(reader).(*TWSEReader)
	if !ok {
		t.Error("NewTWSEReader did not return *TWSEReader")
	}
}

// TestTWSEReader_EmbedsBaseSource tests that TWSEReader embeds BaseSource
func TestTWSEReader_EmbedsBaseSource(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := NewTWSEReader(opts)

	// Check if BaseSource is embedded
	if reader.BaseSource == nil {
		t.Error("TWSEReader does not embed BaseSource")
	}
}

// TestNewTWSEReader_NonNil tests that constructor returns non-nil reader
func TestNewTWSEReader_NonNil(t *testing.T) {
	tests := []struct {
		name string
		opts *internalhttp.ClientOptions
	}{
		{
			name: "with nil options",
			opts: nil,
		},
		{
			name: "with default options",
			opts: internalhttp.DefaultClientOptions(),
		},
		{
			name: "with custom options",
			opts: &internalhttp.ClientOptions{
				Timeout:    60,
				MaxRetries: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewTWSEReader(tt.opts)
			if reader == nil {
				t.Error("NewTWSEReader returned nil")
			}
		})
	}
}

// TestTWSEReader_ImplementsReader tests that TWSEReader implements Reader interface
func TestTWSEReader_ImplementsReader(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()
	reader := NewTWSEReader(opts)

	// Try to assign to sources.Reader interface - will fail at compile time if not implemented
	var _ sources.Reader = reader

	// Check that all required methods exist
	if reader.Name() == "" {
		t.Error("Name() method returned empty string")
	}

	// ValidateSymbol should exist and be callable
	err := reader.ValidateSymbol("2330")
	// We expect either nil or a validation error, not a panic
	_ = err
}

// TestTWSEReader_Name tests the Name method
func TestTWSEReader_Name(t *testing.T) {
	reader := NewTWSEReader(nil)
	name := reader.Name()

	expectedName := "Taiwan Stock Exchange"
	if name != expectedName {
		t.Errorf("Name() = %q, want %q", name, expectedName)
	}
}

// TestTWSEReader_ValidateSymbol tests basic symbol validation
func TestTWSEReader_ValidateSymbol(t *testing.T) {
	reader := NewTWSEReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid 4-digit code",
			symbol:  "2330",
			wantErr: false,
		},
		{
			name:    "valid 4-digit code with leading zero",
			symbol:  "0050",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "23 30",
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

// TestBuildDailyURL tests the buildDailyURL function
func TestBuildDailyURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{
			name:    "default base URL",
			baseURL: twseBaseURL,
			want:    "https://openapi.twse.com.tw/v1/exchangeReport/STOCK_DAY_ALL",
		},
		{
			name:    "custom base URL",
			baseURL: "https://example.com/api",
			want:    "https://example.com/api/exchangeReport/STOCK_DAY_ALL",
		},
		{
			name:    "base URL with trailing slash",
			baseURL: "https://openapi.twse.com.tw/v1/",
			want:    "https://openapi.twse.com.tw/v1/exchangeReport/STOCK_DAY_ALL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDailyURL(tt.baseURL)
			if got != tt.want {
				t.Errorf("buildDailyURL(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

// TestBuildDailyURL_ValidURL tests that buildDailyURL returns a valid URL
func TestBuildDailyURL_ValidURL(t *testing.T) {
	url := buildDailyURL(twseBaseURL)

	// Check that URL starts with https
	if !strings.HasPrefix(url, "https://") {
		t.Errorf("buildDailyURL() URL should start with https://, got: %s", url)
	}

	// Check that URL contains the endpoint
	if !strings.Contains(url, dailyStocksEndpoint) {
		t.Errorf("buildDailyURL() URL should contain %s, got: %s", dailyStocksEndpoint, url)
	}

	// Check that URL is properly formed (no double slashes except after https://)
	if strings.Contains(url[8:], "//") {
		t.Errorf("buildDailyURL() URL contains double slashes: %s", url)
	}
}

// TestBuildIndexURL tests the buildIndexURL function
func TestBuildIndexURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{
			name:    "default base URL",
			baseURL: twseBaseURL,
			want:    "https://openapi.twse.com.tw/v1/exchangeReport/MI_INDEX",
		},
		{
			name:    "custom base URL",
			baseURL: "https://example.com/api",
			want:    "https://example.com/api/exchangeReport/MI_INDEX",
		},
		{
			name:    "base URL with trailing slash",
			baseURL: "https://openapi.twse.com.tw/v1/",
			want:    "https://openapi.twse.com.tw/v1/exchangeReport/MI_INDEX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildIndexURL(tt.baseURL)
			if got != tt.want {
				t.Errorf("buildIndexURL(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

// TestBuildIndexURL_ValidURL tests that buildIndexURL returns a valid URL
func TestBuildIndexURL_ValidURL(t *testing.T) {
	url := buildIndexURL(twseBaseURL)

	// Check that URL starts with https
	if !strings.HasPrefix(url, "https://") {
		t.Errorf("buildIndexURL() URL should start with https://, got: %s", url)
	}

	// Check that URL contains the endpoint
	if !strings.Contains(url, indexEndpoint) {
		t.Errorf("buildIndexURL() URL should contain %s, got: %s", indexEndpoint, url)
	}

	// Check that URL is properly formed (no double slashes except after https://)
	if strings.Contains(url[8:], "//") {
		t.Errorf("buildIndexURL() URL contains double slashes: %s", url)
	}
}

// TestTWSEReader_BuildURL tests the BuildURL method
func TestTWSEReader_BuildURL(t *testing.T) {
	reader := NewTWSEReader(nil)

	url := reader.BuildURL()

	// Check that URL is not empty
	if url == "" {
		t.Error("BuildURL() returned empty string")
	}

	// Check that URL starts with base URL
	if !strings.HasPrefix(url, twseBaseURL) {
		t.Errorf("BuildURL() should start with %s, got: %s", twseBaseURL, url)
	}

	// Check that URL contains daily endpoint
	if !strings.Contains(url, dailyStocksEndpoint) {
		t.Errorf("BuildURL() should contain %s, got: %s", dailyStocksEndpoint, url)
	}
}

// TestTWSEReader_BuildURL_CustomBaseURL tests BuildURL with custom base URL
func TestTWSEReader_BuildURL_CustomBaseURL(t *testing.T) {
	customURL := "https://test.example.com/v1"
	reader := NewTWSEReaderWithBaseURL(nil, customURL)

	url := reader.BuildURL()

	// Check that URL starts with custom base URL
	if !strings.HasPrefix(url, customURL) {
		t.Errorf("BuildURL() should start with %s, got: %s", customURL, url)
	}
}

// TestTWSEReader_ReadSingle tests the ReadSingle method with mock data
func TestTWSEReader_ReadSingle(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock TWSE API response with sample data
		mockData := []TWSEStockData{
			{
				Date:         "1141028",
				Code:         "2330",
				Name:         "台積電",
				TradeVolume:  "25000000",
				TradeValue:   "23750000000",
				OpeningPrice: "950.00",
				HighestPrice: "960.00",
				LowestPrice:  "945.00",
				ClosingPrice: "955.00",
				Change:       "+5.00",
				Transaction:  "12500",
			},
			{
				Date:         "1141028",
				Code:         "2317",
				Name:         "鴻海",
				TradeVolume:  "15000000",
				TradeValue:   "1575000000",
				OpeningPrice: "105.00",
				HighestPrice: "106.50",
				LowestPrice:  "104.50",
				ClosingPrice: "105.50",
				Change:       "+0.50",
				Transaction:  "8500",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// Create reader with mock server URL
	reader := NewTWSEReaderWithBaseURL(nil, server.URL)

	// Test parameters
	ctx := context.Background()
	symbol := "2330"
	start := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)

	// Execute ReadSingle
	result, err := reader.ReadSingle(ctx, symbol, start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	// Verify result is not nil
	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	// Verify result is ParsedData
	data, ok := result.(*ParsedData)
	if !ok {
		t.Fatalf("ReadSingle() returned %T, want *ParsedData", result)
	}

	// Verify data contains expected symbol
	if data.Symbol != symbol {
		t.Errorf("Symbol = %q, want %q", data.Symbol, symbol)
	}

	// Verify data has at least one entry
	if len(data.Date) == 0 {
		t.Error("ParsedData has no dates")
	}

	if len(data.Open) == 0 {
		t.Error("ParsedData has no opening prices")
	}
}

// TestTWSEReader_ReadSingle_ValidatesSymbol tests that ReadSingle validates symbols
func TestTWSEReader_ReadSingle_ValidatesSymbol(t *testing.T) {
	reader := NewTWSEReader(nil)
	ctx := context.Background()
	start := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "valid symbol",
			symbol:  "2330",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := reader.ReadSingle(ctx, tt.symbol, start, end)

			if tt.wantErr && err == nil {
				t.Error("ReadSingle() expected error for invalid symbol, got nil")
			}

			if !tt.wantErr && err != nil {
				// Note: valid symbols may still fail due to network issues in this test
				// So we only check that error message mentions symbol validation if it fails
				if strings.Contains(err.Error(), "invalid symbol") {
					t.Errorf("ReadSingle() unexpected symbol validation error: %v", err)
				}
			}
		})
	}
}

// TestTWSEReader_ReadSingle_ValidatesDateRange tests date range validation
func TestTWSEReader_ReadSingle_ValidatesDateRange(t *testing.T) {
	reader := NewTWSEReader(nil)
	ctx := context.Background()
	symbol := "2330"

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr bool
	}{
		{
			name:    "end before start",
			start:   time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			end:     time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "valid range",
			start:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			end:     time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := reader.ReadSingle(ctx, symbol, tt.start, tt.end)

			if tt.wantErr && err == nil {
				t.Error("ReadSingle() expected error for invalid date range, got nil")
			}

			if !tt.wantErr && err != nil {
				// Note: valid ranges may still fail due to network issues
				// So we only check that error message mentions date range if it fails
				if strings.Contains(err.Error(), "invalid date range") {
					t.Errorf("ReadSingle() unexpected date range error: %v", err)
				}
			}
		})
	}
}

// TestTWSEReader_Read tests the Read method with multiple symbols
func TestTWSEReader_Read(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock TWSE API response with sample data
		mockData := []TWSEStockData{
			{
				Date:         "1141028",
				Code:         "2330",
				Name:         "台積電",
				TradeVolume:  "25000000",
				TradeValue:   "23750000000",
				OpeningPrice: "950.00",
				HighestPrice: "960.00",
				LowestPrice:  "945.00",
				ClosingPrice: "955.00",
				Change:       "+5.00",
				Transaction:  "12500",
			},
			{
				Date:         "1141028",
				Code:         "2317",
				Name:         "鴻海",
				TradeVolume:  "15000000",
				TradeValue:   "1575000000",
				OpeningPrice: "105.00",
				HighestPrice: "106.50",
				LowestPrice:  "104.50",
				ClosingPrice: "105.50",
				Change:       "+0.50",
				Transaction:  "8500",
			},
			{
				Date:         "1141028",
				Code:         "2454",
				Name:         "聯發科",
				TradeVolume:  "8000000",
				TradeValue:   "7600000000",
				OpeningPrice: "950.00",
				HighestPrice: "955.00",
				LowestPrice:  "945.00",
				ClosingPrice: "950.00",
				Change:       "0.00",
				Transaction:  "5500",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// Create reader with mock server URL
	reader := NewTWSEReaderWithBaseURL(nil, server.URL)

	// Test parameters
	ctx := context.Background()
	symbols := []string{"2330", "2317", "2454"}
	start := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)

	// Execute Read
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	// Verify result is not nil
	if result == nil {
		t.Fatal("Read() returned nil result")
	}

	// Verify result is map[string]*ParsedData
	dataMap, ok := result.(map[string]*ParsedData)
	if !ok {
		t.Fatalf("Read() returned %T, want map[string]*ParsedData", result)
	}

	// Verify all symbols are present
	if len(dataMap) != len(symbols) {
		t.Errorf("Read() returned %d symbols, want %d", len(dataMap), len(symbols))
	}

	for _, symbol := range symbols {
		data, found := dataMap[symbol]
		if !found {
			t.Errorf("Read() missing data for symbol %q", symbol)
			continue
		}

		if data.Symbol != symbol {
			t.Errorf("Symbol = %q, want %q", data.Symbol, symbol)
		}

		if len(data.Date) == 0 {
			t.Errorf("Symbol %q has no dates", symbol)
		}
	}
}

// TestTWSEReader_Read_ValidatesSymbols tests that Read validates symbol list
func TestTWSEReader_Read_ValidatesSymbols(t *testing.T) {
	reader := NewTWSEReader(nil)
	ctx := context.Background()
	start := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		symbols []string
		wantErr bool
	}{
		{
			name:    "empty symbol list",
			symbols: []string{},
			wantErr: true,
		},
		{
			name:    "nil symbol list",
			symbols: nil,
			wantErr: true,
		},
		{
			name:    "valid symbols",
			symbols: []string{"2330", "2317"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := reader.Read(ctx, tt.symbols, start, end)

			if tt.wantErr && err == nil {
				t.Error("Read() expected error for invalid symbols, got nil")
			}

			if !tt.wantErr && err != nil {
				// Note: valid symbols may still fail due to network issues
				// So we only check that error message mentions symbol validation if it fails
				if strings.Contains(err.Error(), "invalid symbols") {
					t.Errorf("Read() unexpected symbol validation error: %v", err)
				}
			}
		})
	}
}
