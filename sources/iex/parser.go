package iex

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ParsedData represents parsed IEX Cloud chart data.
type ParsedData struct {
	Columns []string
	Rows    []map[string]string
}

// chartDataPoint represents a single day of IEX Cloud chart data.
type chartDataPoint struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// errorResponse represents an IEX Cloud API error.
type errorResponse struct {
	Error string `json:"error"`
}

// ParseResponse parses IEX Cloud JSON chart data response.
// IEX Cloud returns an array of daily chart data points.
func ParseResponse(data []byte) (*ParsedData, error) {
	// Check if response is an error
	var errResp errorResponse
	if err := json.Unmarshal(data, &errResp); err == nil && errResp.Error != "" {
		return nil, fmt.Errorf("API error: %s", errResp.Error)
	}

	// Parse as array of chart data
	var chartData []chartDataPoint
	if err := json.Unmarshal(data, &chartData); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// Sort by date ascending
	sort.SliceStable(chartData, func(i, j int) bool {
		return chartData[i].Date < chartData[j].Date
	})

	// Convert to parsed data format
	rows := make([]map[string]string, 0, len(chartData))
	for _, point := range chartData {
		row := map[string]string{
			"Date":   point.Date,
			"Open":   fmt.Sprintf("%.2f", point.Open),
			"High":   fmt.Sprintf("%.2f", point.High),
			"Low":    fmt.Sprintf("%.2f", point.Low),
			"Close":  fmt.Sprintf("%.2f", point.Close),
			"Volume": fmt.Sprintf("%d", point.Volume),
		}
		rows = append(rows, row)
	}

	return &ParsedData{
		Columns: []string{"Date", "Open", "High", "Low", "Close", "Volume"},
		Rows:    rows,
	}, nil
}
