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
