// Package utils provides internal utility functions for gonp-datareader.
package utils

import (
	"errors"
	"strings"
	"unicode"
)

var (
	// ErrEmptySymbol is returned when a symbol is empty
	ErrEmptySymbol = errors.New("symbol cannot be empty")
	// ErrInvalidSymbolFormat is returned when a symbol contains invalid characters
	ErrInvalidSymbolFormat = errors.New("symbol contains invalid characters")
	// ErrEmptySymbolList is returned when a symbol list is empty or nil
	ErrEmptySymbolList = errors.New("symbol list cannot be empty")
)

// ValidateSymbol checks if a symbol is valid.
// A valid symbol:
// - Must not be empty
// - Must not contain whitespace
// - Must contain only alphanumeric characters, dots, and hyphens
func ValidateSymbol(symbol string) error {
	if symbol == "" {
		return ErrEmptySymbol
	}

	// Check for whitespace
	if strings.ContainsAny(symbol, " \t\n\r") {
		return ErrInvalidSymbolFormat
	}

	// Check for valid characters (alphanumeric, dot, hyphen)
	for _, r := range symbol {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' {
			return ErrInvalidSymbolFormat
		}
	}

	return nil
}

// ValidateSymbols validates a list of symbols.
// Returns an error if the list is empty or any symbol is invalid.
func ValidateSymbols(symbols []string) error {
	if symbols == nil || len(symbols) == 0 {
		return ErrEmptySymbolList
	}

	for _, symbol := range symbols {
		if err := ValidateSymbol(symbol); err != nil {
			return err
		}
	}

	return nil
}
