package twse

import (
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
