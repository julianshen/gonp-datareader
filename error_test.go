package datareader_test

import (
	"errors"
	"testing"

	"github.com/julianshen/gonp-datareader"
)

func TestDataReaderError_Error(t *testing.T) {
	err := &datareader.DataReaderError{
		Type:    datareader.ErrInvalidSymbol,
		Source:  "yahoo",
		Message: "symbol cannot be empty",
		Cause:   errors.New("validation failed"),
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() should return non-empty string")
	}

	// Should contain source and message
	if errMsg != "yahoo: symbol cannot be empty: validation failed" {
		t.Errorf("Expected formatted error message, got: %s", errMsg)
	}
}

func TestDataReaderError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &datareader.DataReaderError{
		Type:    datareader.ErrNetworkError,
		Source:  "fred",
		Message: "network failure",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() should return cause error, got %v", unwrapped)
	}
}

func TestDataReaderError_Is(t *testing.T) {
	err := &datareader.DataReaderError{
		Type:    datareader.ErrInvalidSymbol,
		Source:  "yahoo",
		Message: "invalid symbol",
	}

	// Test matching error
	matchingErr := &datareader.DataReaderError{
		Type:    datareader.ErrInvalidSymbol,
		Source:  "yahoo",
		Message: "invalid symbol",
	}

	if !errors.Is(err, matchingErr) {
		t.Error("errors.Is should return true for matching DataReaderError")
	}

	// Test non-matching type
	differentType := &datareader.DataReaderError{
		Type:    datareader.ErrNetworkError,
		Source:  "yahoo",
		Message: "invalid symbol",
	}

	if errors.Is(err, differentType) {
		t.Error("errors.Is should return false for different error type")
	}

	// Test non-matching source
	differentSource := &datareader.DataReaderError{
		Type:    datareader.ErrInvalidSymbol,
		Source:  "fred",
		Message: "invalid symbol",
	}

	if errors.Is(err, differentSource) {
		t.Error("errors.Is should return false for different source")
	}

	// Test non-matching message
	differentMessage := &datareader.DataReaderError{
		Type:    datareader.ErrInvalidSymbol,
		Source:  "yahoo",
		Message: "different message",
	}

	if errors.Is(err, differentMessage) {
		t.Error("errors.Is should return false for different message")
	}

	// Test non-DataReaderError
	regularErr := errors.New("regular error")
	if errors.Is(err, regularErr) {
		t.Error("errors.Is should return false for non-DataReaderError")
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that error type constants are defined
	errorTypes := []datareader.ErrorType{
		datareader.ErrInvalidSymbol,
		datareader.ErrInvalidDateRange,
		datareader.ErrNetworkError,
		datareader.ErrAPILimit,
		datareader.ErrAuthenticationFailed,
		datareader.ErrDataNotFound,
		datareader.ErrParsingFailed,
	}

	for i, et := range errorTypes {
		if et != datareader.ErrorType(i) {
			t.Errorf("ErrorType constant mismatch at index %d", i)
		}
	}
}

func TestNewDataReaderError(t *testing.T) {
	tests := []struct {
		name    string
		errType datareader.ErrorType
		source  string
		message string
		cause   error
	}{
		{
			name:    "invalid symbol error",
			errType: datareader.ErrInvalidSymbol,
			source:  "yahoo",
			message: "symbol is empty",
			cause:   nil,
		},
		{
			name:    "network error with cause",
			errType: datareader.ErrNetworkError,
			source:  "fred",
			message: "connection timeout",
			cause:   errors.New("timeout"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := datareader.NewDataReaderError(tt.errType, tt.source, tt.message, tt.cause)

			if err == nil {
				t.Fatal("NewDataReaderError returned nil")
			}

			if err.Type != tt.errType {
				t.Errorf("Expected type %v, got %v", tt.errType, err.Type)
			}

			if err.Source != tt.source {
				t.Errorf("Expected source %s, got %s", tt.source, err.Source)
			}

			if err.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, err.Message)
			}

			if err.Cause != tt.cause {
				t.Errorf("Expected cause %v, got %v", tt.cause, err.Cause)
			}
		})
	}
}
