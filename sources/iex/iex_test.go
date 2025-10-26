package iex_test

import (
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
