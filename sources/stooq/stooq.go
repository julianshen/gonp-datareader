// Package stooq provides a Stooq data source reader.
package stooq

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

// StooqReader fetches data from Stooq.
type StooqReader struct {
	*sources.BaseSource
	client *internalhttp.RetryableClient
}

// NewStooqReader creates a new Stooq data reader.
func NewStooqReader(opts *internalhttp.ClientOptions) *StooqReader {
	return &StooqReader{
		BaseSource: sources.NewBaseSource("stooq"),
		client:     internalhttp.NewRetryableClient(opts),
	}
}

// BuildURL constructs the Stooq URL for fetching historical data.
// The Stooq format is:
// https://stooq.com/q/d/l/?s={symbol}&i=d
// where i=d means daily data
func BuildURL(symbol string) string {
	return fmt.Sprintf(
		"https://stooq.com/q/d/l/?s=%s&i=d",
		url.QueryEscape(symbol),
	)
}

// ReadSingle fetches data for a single symbol.
func (s *StooqReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	// Validate symbol
	if err := s.ValidateSymbol(symbol); err != nil {
		return nil, err
	}

	// Build URL
	urlStr := BuildURL(symbol)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse CSV
	data, err := ParseCSV(body)
	if err != nil {
		return nil, fmt.Errorf("parse CSV: %w", err)
	}

	return data, nil
}

// Read fetches data for multiple symbols.
func (s *StooqReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// ValidateSymbol checks if a symbol is valid for Stooq.
func (s *StooqReader) ValidateSymbol(symbol string) error {
	return s.BaseSource.ValidateSymbol(symbol)
}
