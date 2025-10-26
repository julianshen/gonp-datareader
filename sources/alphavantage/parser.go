package alphavantage

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// ParsedData represents parsed Alpha Vantage time series data.
type ParsedData struct {
	Columns []string
	Rows    []map[string]string
}

// alphaVantageResponse represents the Alpha Vantage API response structure.
type alphaVantageResponse struct {
	MetaData   map[string]string            `json:"Meta Data"`
	TimeSeries map[string]map[string]string `json:"Time Series (Daily)"`
	Note       string                       `json:"Note"`
	ErrorMsg   string                       `json:"Error Message"`
}

// ParseResponse parses the Alpha Vantage API JSON response.
func ParseResponse(data []byte) (*ParsedData, error) {
	var response alphaVantageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// Check for rate limit
	if response.Note != "" {
		return nil, errors.New("rate limit exceeded")
	}

	// Check for error message
	if response.ErrorMsg != "" {
		return nil, fmt.Errorf("API error: %s", response.ErrorMsg)
	}

	// Check if time series exists
	if len(response.TimeSeries) == 0 {
		return &ParsedData{
			Columns: []string{"Date", "Open", "High", "Low", "Close", "Volume"},
			Rows:    []map[string]string{},
		}, nil
	}

	// Extract dates and sort them
	var dates []string
	for date := range response.TimeSeries {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Build rows
	rows := make([]map[string]string, 0, len(dates))
	for _, date := range dates {
		values := response.TimeSeries[date]
		row := map[string]string{
			"Date":   date,
			"Open":   values["1. open"],
			"High":   values["2. high"],
			"Low":    values["3. low"],
			"Close":  values["4. close"],
			"Volume": values["5. volume"],
		}
		rows = append(rows, row)
	}

	return &ParsedData{
		Columns: []string{"Date", "Open", "High", "Low", "Close", "Volume"},
		Rows:    rows,
	}, nil
}
