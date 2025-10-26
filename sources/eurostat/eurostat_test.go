package eurostat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/eurostat"
)

func TestNewEurostatReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{}
	reader := eurostat.NewEurostatReader(opts)

	if reader == nil {
		t.Fatal("NewEurostatReader() returned nil")
	}

	if reader.Name() != "Eurostat" {
		t.Errorf("Expected name 'Eurostat', got '%s'", reader.Name())
	}

	if reader.Source() != "eurostat" {
		t.Errorf("Expected source 'eurostat', got '%s'", reader.Source())
	}
}

func TestEurostatReader_ValidateSymbol(t *testing.T) {
	reader := eurostat.NewEurostatReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid dataset code",
			symbol:  "DEMO_R_D3DENS",
			wantErr: false,
		},
		{
			name:    "valid simple code",
			symbol:  "GDP",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "DEMO R D3DENS",
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

func TestEurostatReader_BuildURL(t *testing.T) {
	reader := eurostat.NewEurostatReader(nil)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	url := reader.BuildURL("DEMO_R_D3DENS", start, end)

	if url == "" {
		t.Error("BuildURL() returned empty string")
	}

	// URL should contain the Eurostat domain
	if !contains(url, "ec.europa.eu/eurostat") {
		t.Errorf("URL should contain Eurostat domain: %s", url)
	}

	// URL should contain dataset code
	if !contains(url, "DEMO_R_D3DENS") {
		t.Errorf("URL should contain dataset 'DEMO_R_D3DENS': %s", url)
	}

	// URL should contain language parameter
	if !contains(url, "lang=") {
		t.Errorf("URL should contain language parameter: %s", url)
	}
}

func TestEurostatReader_ReadSingle_WithMockServer(t *testing.T) {
	// Simplified JSON-stat response structure
	jsonData := `{
		"version": "2.0",
		"class": "dataset",
		"label": "Test Dataset",
		"id": ["geo", "time"],
		"size": [1, 3],
		"dimension": {
			"geo": {
				"label": "Geopolitical entity",
				"category": {
					"index": {"EU27_2020": 0},
					"label": {"EU27_2020": "European Union"}
				}
			},
			"time": {
				"label": "Time",
				"category": {
					"index": {"2020": 0, "2021": 1, "2022": 2},
					"label": {"2020": "2020", "2021": "2021", "2022": "2022"}
				}
			}
		},
		"value": [100.5, 102.3, 104.1]
	}`

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	reader := eurostat.NewEurostatReaderWithBaseURL(nil, server.URL+"/statistics/1.0/data/%s")

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)

	result, err := reader.ReadSingle(ctx, "DEMO_R_D3DENS", start, end)
	if err != nil {
		t.Fatalf("ReadSingle() error = %v", err)
	}

	if result == nil {
		t.Fatal("ReadSingle() returned nil result")
	}

	// Verify we got parsed data
	data, ok := result.(*eurostat.ParsedData)
	if !ok {
		t.Fatalf("Expected *eurostat.ParsedData, got %T", result)
	}

	if len(data.Dates) != 3 {
		t.Errorf("Expected 3 dates, got %d", len(data.Dates))
	}

	if len(data.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(data.Values))
	}

	// Check first value
	if data.Values[0] != 100.5 {
		t.Errorf("Expected first value 100.5, got %f", data.Values[0])
	}
}

func TestEurostatReader_ReadSingle_InvalidSymbol(t *testing.T) {
	reader := eurostat.NewEurostatReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "", start, end)
	if err == nil {
		t.Error("ReadSingle() should return error for empty symbol")
	}
}

func TestEurostatReader_ReadSingle_InvalidDateRange(t *testing.T) {
	reader := eurostat.NewEurostatReader(nil)

	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := reader.ReadSingle(ctx, "DEMO_R_D3DENS", start, end)
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
