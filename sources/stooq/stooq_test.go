package stooq_test

import (
	"strings"
	"testing"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/stooq"
)

func TestNewStooqReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	reader := stooq.NewStooqReader(opts)

	if reader == nil {
		t.Fatal("NewStooqReader returned nil")
	}

	if reader.Name() != "stooq" {
		t.Errorf("Expected name 'stooq', got %q", reader.Name())
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		wantParts []string
	}{
		{
			name:   "US stock",
			symbol: "AAPL.US",
			wantParts: []string{
				"stooq.com",
				"s=AAPL.US",
				"i=d",
			},
		},
		{
			name:   "index",
			symbol: "^SPX",
			wantParts: []string{
				"stooq.com",
				"s=%5ESPX", // URL encoded ^
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := stooq.BuildURL(tt.symbol)

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
