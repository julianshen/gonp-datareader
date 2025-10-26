package datareader

import "fmt"

// ErrorType represents the type of error that occurred.
type ErrorType int

const (
	// ErrInvalidSymbol indicates an invalid symbol was provided
	ErrInvalidSymbol ErrorType = iota
	// ErrInvalidDateRange indicates an invalid date range was provided
	ErrInvalidDateRange
	// ErrNetworkError indicates a network-related error
	ErrNetworkError
	// ErrAPILimit indicates an API rate limit was exceeded
	ErrAPILimit
	// ErrAuthenticationFailed indicates authentication failed
	ErrAuthenticationFailed
	// ErrDataNotFound indicates the requested data was not found
	ErrDataNotFound
	// ErrParsingFailed indicates data parsing failed
	ErrParsingFailed
)

// DataReaderError provides detailed error information for data reader operations.
type DataReaderError struct {
	// Type indicates the category of error
	Type ErrorType
	// Source is the name of the data source (e.g., "yahoo", "fred")
	Source string
	// Message is a human-readable description of the error
	Message string
	// Cause is the underlying error that caused this error, if any
	Cause error
}

// Error implements the error interface.
func (e *DataReaderError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Source, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Source, e.Message)
}

// Unwrap returns the underlying cause error.
func (e *DataReaderError) Unwrap() error {
	return e.Cause
}

// Is implements error matching for errors.Is.
func (e *DataReaderError) Is(target error) bool {
	t, ok := target.(*DataReaderError)
	if !ok {
		return false
	}
	return e.Type == t.Type && e.Source == t.Source && e.Message == t.Message
}

// NewDataReaderError creates a new DataReaderError.
func NewDataReaderError(errType ErrorType, source, message string, cause error) *DataReaderError {
	return &DataReaderError{
		Type:    errType,
		Source:  source,
		Message: message,
		Cause:   cause,
	}
}
