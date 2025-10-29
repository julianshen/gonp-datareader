package worldbank_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/worldbank"
)

func TestNewWorldBankReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	reader := worldbank.NewWorldBankReader(opts)

	if reader == nil {
		t.Fatal("NewWorldBankReader returned nil")
	}

	if reader.Name() != "worldbank" {
		t.Errorf("Expected name 'worldbank', got %q", reader.Name())
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		country   string
		indicator string
		start     time.Time
		end       time.Time
		wantParts []string
	}{
		{
			name:      "single country GDP",
			country:   "USA",
			indicator: "NY.GDP.MKTP.CD",
			start:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			end:       time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			wantParts: []string{
				"api.worldbank.org",
				"/v2/country/USA/indicator/NY.GDP.MKTP.CD",
				"date=2020:2023",
				"format=json",
			},
		},
		{
			name:      "multiple countries",
			country:   "USA;CHN;GBR",
			indicator: "SP.POP.TOTL",
			start:     time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			end:       time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
			wantParts: []string{
				"/v2/country/USA;CHN;GBR/indicator/SP.POP.TOTL",
				"date=2015:2020",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := worldbank.BuildURL(tt.country, tt.indicator, tt.start, tt.end)

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

// TestWorldBankReader_ValidateSymbol tests symbol validation
func TestWorldBankReader_ValidateSymbol(t *testing.T) {
	reader := worldbank.NewWorldBankReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid country/indicator",
			symbol:  "USA/NY.GDP.MKTP.CD",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "US A/GDP",
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

// TestWorldBankReader_ReadSingle_WithMockServer tests ReadSingle with mock HTTP server
func TestWorldBankReader_ReadSingle_WithMockServer(t *testing.T) {
	// Sample World Bank JSON response
	jsonResponse := `[
		{"page":1,"pages":1,"per_page":50,"total":2},
		[
			{
				"indicator":{"id":"NY.GDP.MKTP.CD","value":"GDP (current US$)"},
				"country":{"id":"US","value":"United States"},
				"countryiso3code":"USA",
				"date":"2023",
				"value":25462700000000,
				"unit":"",
				"obs_status":"",
				"decimal":0
			},
			{
				"indicator":{"id":"NY.GDP.MKTP.CD","value":"GDP (current US$)"},
				"country":{"id":"US","value":"United States"},
				"countryiso3code":"USA",
				"date":"2022",
				"value":25035200000000,
				"unit":"",
				"obs_status":"",
				"decimal":0
			}
		]
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	// Create reader with custom base URL
	reader := worldbank.NewWorldBankReaderWithBaseURL(nil, server.URL+"?country=%s&indicator=%s&start=%d&end=%d")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "USA/NY.GDP.MKTP.CD", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	data, ok := result.(*worldbank.ParsedData)
	if !ok {
		t.Fatalf("Expected *worldbank.ParsedData, got %T", result)
	}

	if len(data.Dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(data.Dates))
	}
}

// TestWorldBankReader_ReadSingle_InvalidSymbol tests error handling for invalid symbols
func TestWorldBankReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := worldbank.NewWorldBankReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

// TestWorldBankReader_ReadSingle_InvalidFormat tests error for wrong symbol format
func TestWorldBankReader_ReadSingle_InvalidFormat(t *testing.T) {
	reader := worldbank.NewWorldBankReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "INVALID_FORMAT", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for invalid symbol format")
	}

	if !strings.Contains(err.Error(), "invalid symbol format") {
		t.Errorf("Expected 'invalid symbol format' error, got: %v", err)
	}
}

// TestWorldBankReader_Read_MultipleSymbols tests fetching multiple symbols
func TestWorldBankReader_Read_MultipleSymbols(t *testing.T) {
	jsonResponse := `[
		{"page":1,"pages":1,"per_page":50,"total":1},
		[
			{
				"indicator":{"id":"NY.GDP.MKTP.CD","value":"GDP"},
				"country":{"id":"US","value":"United States"},
				"countryiso3code":"USA",
				"date":"2023",
				"value":25462700000000,
				"unit":"",
				"obs_status":"",
				"decimal":0
			}
		]
	]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	reader := worldbank.NewWorldBankReaderWithBaseURL(nil, server.URL+"?country=%s&indicator=%s&start=%d&end=%d")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	symbols := []string{"USA/NY.GDP.MKTP.CD", "CHN/NY.GDP.MKTP.CD", "GBR/NY.GDP.MKTP.CD"}
	result, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	dataMap, ok := result.(map[string]*worldbank.ParsedData)
	if !ok {
		t.Fatalf("Expected map[string]*worldbank.ParsedData, got %T", result)
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

// TestWorldBankReader_Read_InvalidDateRange tests error handling for invalid date ranges
func TestWorldBankReader_Read_InvalidDateRange(t *testing.T) {
	reader := worldbank.NewWorldBankReader(nil)

	ctx := context.Background()
	start := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) // end before start

	_, err := reader.Read(ctx, []string{"USA/NY.GDP.MKTP.CD"}, start, end)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

// TestWorldBankReader_HTTPError tests handling of HTTP errors
func TestWorldBankReader_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	reader := worldbank.NewWorldBankReaderWithBaseURL(nil, server.URL+"?country=%s&indicator=%s&start=%d&end=%d")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "USA/NY.GDP.MKTP.CD", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for HTTP 500")
	}
}
