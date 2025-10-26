package utils_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/internal/utils"
)

func TestValidateSymbol(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "valid single symbol",
			symbol:  "AAPL",
			wantErr: false,
		},
		{
			name:    "valid symbol with numbers",
			symbol:  "BRK.B",
			wantErr: false,
		},
		{
			name:    "valid symbol with dash",
			symbol:  "BRK-B",
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
		{
			name:    "symbol with leading space",
			symbol:  " AAPL",
			wantErr: true,
		},
		{
			name:    "symbol with trailing space",
			symbol:  "AAPL ",
			wantErr: true,
		},
		{
			name:    "symbol with invalid characters",
			symbol:  "AAPL@",
			wantErr: true,
		},
		{
			name:    "symbol with newline",
			symbol:  "AAPL\n",
			wantErr: true,
		},
		{
			name:    "symbol with tab",
			symbol:  "AAPL\t",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol(%q) error = %v, wantErr %v", tt.symbol, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSymbols(t *testing.T) {
	tests := []struct {
		name    string
		symbols []string
		wantErr bool
	}{
		{
			name:    "valid symbols",
			symbols: []string{"AAPL", "MSFT", "GOOGL"},
			wantErr: false,
		},
		{
			name:    "empty list",
			symbols: []string{},
			wantErr: true,
		},
		{
			name:    "nil list",
			symbols: nil,
			wantErr: true,
		},
		{
			name:    "contains invalid symbol",
			symbols: []string{"AAPL", "", "MSFT"},
			wantErr: true,
		},
		{
			name:    "contains symbol with spaces",
			symbols: []string{"AAPL", "MS FT"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateSymbols(tt.symbols)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbols(%v) error = %v, wantErr %v", tt.symbols, err, tt.wantErr)
			}
		})
	}
}
