package sources_test

import (
	"context"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader/sources"
)

// TestReaderInterface verifies that the Reader interface is defined correctly.
func TestReaderInterface(t *testing.T) {
	// This test ensures the Reader interface compiles
	var _ sources.Reader = (*mockReader)(nil)
}

// mockReader is a test implementation of the Reader interface
type mockReader struct{}

func (m *mockReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// Mock implementation
	return nil, nil
}

func (m *mockReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Mock implementation
	return nil, nil
}

func (m *mockReader) ValidateSymbol(symbol string) error {
	// Mock implementation
	return nil
}

func (m *mockReader) Name() string {
	return "mock"
}

func (m *mockReader) Source() string {
	return "mock"
}

func TestReaderInterfaceMethods(t *testing.T) {
	reader := &mockReader{}

	// Test Name method
	if reader.Name() == "" {
		t.Error("Name() should return non-empty string")
	}

	// Test Source method
	if reader.Source() == "" {
		t.Error("Source() should return non-empty string")
	}

	// Test ValidateSymbol
	err := reader.ValidateSymbol("AAPL")
	if err != nil {
		t.Errorf("ValidateSymbol should not error for valid symbol: %v", err)
	}

	// Test Read method signature
	ctx := context.Background()
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err = reader.Read(ctx, []string{"AAPL"}, start, end)
	if err != nil {
		t.Errorf("Read should not error in mock: %v", err)
	}

	// Test ReadSingle method signature
	_, err = reader.ReadSingle(ctx, "AAPL", start, end)
	if err != nil {
		t.Errorf("ReadSingle should not error in mock: %v", err)
	}
}

func TestBaseSource(t *testing.T) {
	base := sources.NewBaseSource("yahoo")

	if base.Source() != "yahoo" {
		t.Errorf("Expected source 'yahoo', got '%s'", base.Source())
	}

	if base.Name() != "yahoo" {
		t.Errorf("Expected name 'yahoo', got '%s'", base.Name())
	}
}

func TestBaseSource_ValidateSymbol(t *testing.T) {
	base := sources.NewBaseSource("test")

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "symbol with spaces",
			symbol:  "AA PL",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := base.ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol(%q) error = %v, wantErr %v", tt.symbol, err, tt.wantErr)
			}
		})
	}
}
