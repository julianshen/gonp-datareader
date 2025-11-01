package finmind_test

import (
	"testing"
	"time"

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
