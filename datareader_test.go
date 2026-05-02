package datareader_test

import (
	"context"
	"testing"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources"
)

func TestDataReader(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantErr    bool
		wantName   string
		wantSource string
	}{
		{
			name:       "yahoo source",
			source:     "yahoo",
			wantErr:    false,
			wantName:   "Yahoo Finance",
			wantSource: "yahoo",
		},
		{
			name:       "fred source",
			source:     "fred",
			wantErr:    false,
			wantName:   "FRED",
			wantSource: "fred",
		},
		{
			name:       "tpex source",
			source:     "tpex",
			wantErr:    false,
			wantName:   "Taipei Exchange",
			wantSource: "tpex",
		},
		{
			name:    "unknown source",
			source:  "unknown",
			wantErr: true,
		},
		{
			name:    "empty source",
			source:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := datareader.DataReader(tt.source, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("DataReader(%q) error = %v, wantErr %v", tt.source, err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if reader == nil {
				t.Fatal("DataReader() returned nil reader")
			}

			if reader.Name() != tt.wantName {
				t.Errorf("Expected name %q, got %q", tt.wantName, reader.Name())
			}

			if reader.Source() != tt.wantSource {
				t.Errorf("Expected source %q, got %q", tt.wantSource, reader.Source())
			}
		})
	}
}

func TestDataReader_WithOptions(t *testing.T) {
	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 5,
	}

	reader, err := datareader.DataReader("yahoo", opts)
	if err != nil {
		t.Fatalf("DataReader() error = %v", err)
	}

	if reader == nil {
		t.Fatal("DataReader() returned nil")
	}
}

func TestRead_ConvenienceFunction(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	// This will likely fail without network, but we're testing the interface
	_, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)

	// Error is acceptable (network issues, rate limiting, etc.)
	if err != nil {
		t.Logf("Read() returned error (expected in unit test): %v", err)
	}
}

func TestRead_InvalidSource(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := datareader.Read(ctx, "AAPL", "unknown", start, end, nil)
	if err == nil {
		t.Error("Read() should return error for unknown source")
	}
}

func TestRead_InvalidSymbol(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := datareader.Read(ctx, "", "yahoo", start, end, nil)
	if err == nil {
		t.Error("Read() should return error for empty symbol")
	}
}

func TestRead_InvalidDateRange(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
	if err == nil {
		t.Error("Read() should return error for invalid date range")
	}
}

func TestListSources(t *testing.T) {
	sources := datareader.ListSources()

	if len(sources) < 2 {
		t.Errorf("ListSources() should return at least 2 sources, got %d", len(sources))
	}

	// Check for expected sources
	expectedSources := map[string]bool{
		"yahoo": false,
		"fred":  false,
		"tpex":  false,
	}

	for _, source := range sources {
		if _, ok := expectedSources[source]; ok {
			expectedSources[source] = true
		}
	}

	for source, found := range expectedSources {
		if !found {
			t.Errorf("ListSources() should include '%s'", source)
		}
	}
}

func TestDataReader_ImplementsInterface(t *testing.T) {
	reader, err := datareader.DataReader("yahoo", nil)
	if err != nil {
		t.Fatalf("DataReader() error = %v", err)
	}

	// Verify it implements the Reader interface
	var _ sources.Reader = reader

	// Test interface methods exist
	if reader.Name() == "" {
		t.Error("Name() should return non-empty string")
	}

	if reader.Source() == "" {
		t.Error("Source() should return non-empty string")
	}

	// Test ValidateSymbol
	err = reader.ValidateSymbol("AAPL")
	if err != nil {
		t.Errorf("ValidateSymbol() should not error for valid symbol: %v", err)
	}

	err = reader.ValidateSymbol("")
	if err == nil {
		t.Error("ValidateSymbol() should error for empty symbol")
	}
}

// TestDataReader_TWSE tests TWSE factory registration
func TestDataReader_TWSE(t *testing.T) {
	reader, err := datareader.DataReader("twse", nil)
	if err != nil {
		t.Fatalf("DataReader('twse') error = %v", err)
	}

	if reader == nil {
		t.Fatal("DataReader('twse') returned nil reader")
	}

	// Check name
	expectedName := "Taiwan Stock Exchange"
	if reader.Name() != expectedName {
		t.Errorf("Expected name %q, got %q", expectedName, reader.Name())
	}

	// Check source
	expectedSource := "twse"
	if reader.Source() != expectedSource {
		t.Errorf("Expected source %q, got %q", expectedSource, reader.Source())
	}

	// Test symbol validation
	err = reader.ValidateSymbol("2330")
	if err != nil {
		t.Errorf("ValidateSymbol('2330') should not error: %v", err)
	}

	err = reader.ValidateSymbol("")
	if err == nil {
		t.Error("ValidateSymbol('') should error for empty symbol")
	}
}

// TestListSources_IncludesTWSE tests that TWSE is in the sources list
func TestListSources_IncludesTWSE(t *testing.T) {
	sources := datareader.ListSources()

	found := false
	for _, source := range sources {
		if source == "twse" {
			found = true
			break
		}
	}

	if !found {
		t.Error("ListSources() should include 'twse'")
	}
}

// TestRead_TWSE tests the convenience function with TWSE
func TestRead_TWSE(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	// This will likely fail without network or if the symbol doesn't exist
	// but we're testing that the factory registration works
	_, err := datareader.Read(ctx, "2330", "twse", start, end, nil)

	// Error is acceptable (network issues, rate limiting, etc.)
	// We just want to ensure no "unknown source" error
	if err != nil {
		t.Logf("Read() returned error (may be expected): %v", err)
		// Make sure it's not an "unknown source" error
		if err.Error() == "unknown data source: twse" {
			t.Errorf("Read() failed with unknown source error: %v", err)
		}
	}
}

// TestDataReader_FinMind tests FinMind factory registration
func TestDataReader_FinMind(t *testing.T) {
	reader, err := datareader.DataReader("finmind", nil)
	if err != nil {
		t.Fatalf("DataReader('finmind') error = %v", err)
	}

	if reader == nil {
		t.Fatal("DataReader('finmind') returned nil reader")
	}

	// Check name
	expectedName := "FinMind"
	if reader.Name() != expectedName {
		t.Errorf("Expected name %q, got %q", expectedName, reader.Name())
	}

	// Check source
	expectedSource := "finmind"
	if reader.Source() != expectedSource {
		t.Errorf("Expected source %q, got %q", expectedSource, reader.Source())
	}

	// Test symbol validation
	err = reader.ValidateSymbol("2330")
	if err != nil {
		t.Errorf("ValidateSymbol('2330') should not error: %v", err)
	}

	err = reader.ValidateSymbol("")
	if err == nil {
		t.Error("ValidateSymbol('') should error for empty symbol")
	}
}

// TestDataReader_FinMind_WithAPIKey tests FinMind with API key
func TestDataReader_FinMind_WithAPIKey(t *testing.T) {
	opts := &datareader.Options{
		APIKey: "test-token-123",
	}

	reader, err := datareader.DataReader("finmind", opts)
	if err != nil {
		t.Fatalf("DataReader('finmind') with API key error = %v", err)
	}

	if reader == nil {
		t.Fatal("DataReader('finmind') returned nil reader")
	}

	if reader.Name() != "FinMind" {
		t.Errorf("Expected name 'FinMind', got %q", reader.Name())
	}
}

// TestListSources_IncludesFinMind tests that FinMind is in the sources list
func TestListSources_IncludesFinMind(t *testing.T) {
	sources := datareader.ListSources()

	found := false
	for _, source := range sources {
		if source == "finmind" {
			found = true
			break
		}
	}

	if !found {
		t.Error("ListSources() should include 'finmind'")
	}
}

// TestRead_FinMind tests the convenience function with FinMind
func TestRead_FinMind(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 30, 0, 0, 0, 0, time.UTC)

	// This will likely fail without network or if the symbol doesn't exist
	// but we're testing that the factory registration works
	_, err := datareader.Read(ctx, "2330", "finmind", start, end, nil)

	// Error is acceptable (network issues, rate limiting, etc.)
	// We just want to ensure no "unknown source" error
	if err != nil {
		t.Logf("Read() returned error (may be expected): %v", err)
		// Make sure it's not an "unknown source" error
		if err.Error() == "unknown data source: finmind" {
			t.Errorf("Read() failed with unknown source error: %v", err)
		}
	}
}

// TestDataReader_AllSources tests all data source factory registrations
func TestDataReader_AllSources(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantName   string
		wantSource string
		apiKey     string
	}{
		{
			name:       "yahoo",
			source:     "yahoo",
			wantName:   "Yahoo Finance",
			wantSource: "yahoo",
		},
		{
			name:       "fred without api key",
			source:     "fred",
			wantName:   "FRED",
			wantSource: "fred",
		},
		{
			name:       "fred with api key",
			source:     "fred",
			wantName:   "FRED",
			wantSource: "fred",
			apiKey:     "test-api-key",
		},
		{
			name:       "worldbank",
			source:     "worldbank",
			wantName:   "worldbank",
			wantSource: "worldbank",
		},
		{
			name:       "alphavantage",
			source:     "alphavantage",
			wantName:   "alphavantage",
			wantSource: "alphavantage",
			apiKey:     "test-api-key",
		},
		{
			name:       "stooq",
			source:     "stooq",
			wantName:   "stooq",
			wantSource: "stooq",
		},
		{
			name:       "iex",
			source:     "iex",
			wantName:   "iex",
			wantSource: "iex",
			apiKey:     "test-api-key",
		},
		{
			name:       "tiingo without api key",
			source:     "tiingo",
			wantName:   "Tiingo",
			wantSource: "tiingo",
		},
		{
			name:       "tiingo with api key",
			source:     "tiingo",
			wantName:   "Tiingo",
			wantSource: "tiingo",
			apiKey:     "test-api-key",
		},
		{
			name:       "oecd",
			source:     "oecd",
			wantName:   "OECD",
			wantSource: "oecd",
		},
		{
			name:       "eurostat",
			source:     "eurostat",
			wantName:   "Eurostat",
			wantSource: "eurostat",
		},
		{
			name:       "twse",
			source:     "twse",
			wantName:   "Taiwan Stock Exchange",
			wantSource: "twse",
		},
		{
			name:       "finmind without api key",
			source:     "finmind",
			wantName:   "FinMind",
			wantSource: "finmind",
		},
		{
			name:       "finmind with api key",
			source:     "finmind",
			wantName:   "FinMind",
			wantSource: "finmind",
			apiKey:     "test-api-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &datareader.Options{
				APIKey: tt.apiKey,
			}
			reader, err := datareader.DataReader(tt.source, opts)

			if err != nil {
				t.Fatalf("DataReader(%q) error = %v", tt.source, err)
			}

			if reader == nil {
				t.Fatal("DataReader() returned nil reader")
			}

			if reader.Name() != tt.wantName {
				t.Errorf("Expected name %q, got %q", tt.wantName, reader.Name())
			}

			if reader.Source() != tt.wantSource {
				t.Errorf("Expected source %q, got %q", tt.wantSource, reader.Source())
			}
		})
	}
}

// TestDataReader_NilOptions tests DataReader with nil options
func TestDataReader_NilOptions(t *testing.T) {
	sources := []string{"yahoo", "fred", "worldbank", "stooq", "oecd", "eurostat", "twse", "finmind"}

	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			reader, err := datareader.DataReader(source, nil)
			if err != nil {
				t.Fatalf("DataReader(%q, nil) error = %v", source, err)
			}
			if reader == nil {
				t.Error("DataReader() returned nil reader")
			}
		})
	}
}
