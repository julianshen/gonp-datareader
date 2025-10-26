package alphavantage_test

import (
	"strings"
	"testing"

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
