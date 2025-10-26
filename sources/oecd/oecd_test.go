package oecd_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/oecd"
)

func TestNewOECDReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{}
	reader := oecd.NewOECDReader(opts)

	if reader == nil {
		t.Fatal("NewOECDReader() returned nil")
	}

	if reader.Name() != "OECD" {
		t.Errorf("Expected name 'OECD', got '%s'", reader.Name())
	}

	if reader.Source() != "oecd" {
		t.Errorf("Expected source 'oecd', got '%s'", reader.Source())
	}
}

func TestOECDReader_ValidateSymbol(t *testing.T) {
	reader := oecd.NewOECDReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid dataset with country",
			symbol:  "MEI/USA",
			wantErr: false,
		},
		{
			name:    "valid simple dataset",
			symbol:  "QNA/AUS.GDP",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "MEI / USA",
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

func TestOECDReader_BuildURL(t *testing.T) {
	reader := oecd.NewOECDReader(nil)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("MEI/USA", start, end)

	if url == "" {
		t.Error("BuildURL() returned empty string")
	}

	// URL should contain the OECD domain
	if !contains(url, "stats.oecd.org") {
		t.Errorf("URL should contain OECD domain: %s", url)
	}

	// URL should contain dataset
	if !contains(url, "MEI") {
		t.Errorf("URL should contain dataset 'MEI': %s", url)
	}

	// URL should contain date parameters
	if !contains(url, "startPeriod") {
		t.Errorf("URL should contain startPeriod parameter: %s", url)
	}
}

func TestOECDReader_ReadSingle_WithMockServer(t *testing.T) {
	// Simplified SDMX-JSON response structure
	jsonData := `{
		"header": {
			"id": "test",
			"prepared": "2020-01-01T00:00:00Z"
		},
		"dataSets": [{
			"observations": {
				"0:0:0:0": [150.5],
				"0:0:0:1": [152.3],
				"0:0:0:2": [153.1]
			}
		}],
		"structure": {
			"dimensions": {
				"observation": [
					{
						"id": "LOCATION",
						"values": [{"id": "USA"}]
					},
					{
						"id": "INDICATOR",
						"values": [{"id": "GDP"}]
					},
					{
						"id": "MEASURE",
						"values": [{"id": "IDX"}]
					},
					{
						"id": "TIME_PERIOD",
						"values": [
							{"id": "2020-01"},
							{"id": "2020-02"},
							{"id": "2020-03"}
						]
					}
				]
			}
		}
	}`

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := oecd.NewOECDReaderWithBaseURL(nil, server.URL+"/sdmx-json/data/%s/all")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "MEI/USA", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	// Verify we got parsed data
	data, ok := result.(*oecd.ParsedData)
	if !ok {
		t.Fatalf("Expected *oecd.ParsedData, got %T", result)
	}

	if len(data.Dates) != 3 {
		t.Errorf("Expected 3 dates, got %d", len(data.Dates))
	}

	if len(data.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(data.Values))
	}

	// Check first value
	if data.Values[0] != 150.5 {
		t.Errorf("Expected first value 150.5, got %f", data.Values[0])
	}
}

func TestOECDReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := oecd.NewOECDReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

func TestOECDReader_ReadSingle_InvalidDateRange(t *testing.T) {
	reader := oecd.NewOECDReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "MEI/USA", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for invalid date range")
	}
}

// Helper function
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
