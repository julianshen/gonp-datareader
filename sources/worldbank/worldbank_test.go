package worldbank_test

import (
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
