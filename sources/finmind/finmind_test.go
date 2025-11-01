package finmind_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/finmind"
)

func TestNewFinMindReader(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	if reader == nil {
		t.Fatal("NewFinMindReader() returned nil")
	}

	if reader.Name() != "FinMind" {
		t.Errorf("Expected name 'FinMind', got %q", reader.Name())
	}

	if reader.Source() != "finmind" {
		t.Errorf("Expected source 'finmind', got %q", reader.Source())
	}
}

func TestNewFinMindReaderWithToken(t *testing.T) {
	token := "test-token-123"
	reader := finmind.NewFinMindReaderWithToken(nil, token)

	if reader == nil {
		t.Fatal("NewFinMindReaderWithToken() returned nil")
	}

	if reader.Name() != "FinMind" {
		t.Errorf("Expected name 'FinMind', got %q", reader.Name())
	}

	if reader.Source() != "finmind" {
		t.Errorf("Expected source 'finmind', got %q", reader.Source())
	}
}

func TestFinMindReader_ValidateSymbol(t *testing.T) {
	reader := finmind.NewFinMindReader(nil)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid Taiwan 4-digit stock code",
			symbol:  "2330",
			wantErr: false,
		},
		{
			name:    "valid Taiwan 6-digit warrant code",
			symbol:  "123456",
			wantErr: false,
		},
		{
			name:    "valid US stock symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			symbol:  "   ",
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
