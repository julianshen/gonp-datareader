package twse

import (
	"strings"
	"testing"

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
