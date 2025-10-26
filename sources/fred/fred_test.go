package fred_test

import (
	"context"
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
