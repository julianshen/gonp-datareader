package worldbank

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ParsedData represents parsed World Bank indicator data.
type ParsedData struct {
	Dates  []string
	Values []string
}

// observation represents a single data point from the World Bank API.
type observation struct {
	Indicator struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"indicator"`
	Country struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"country"`
	CountryISO3Code string      `json:"countryiso3code"`
	Date            string      `json:"date"`
	Value           interface{} `json:"value"`
	Unit            string      `json:"unit"`
	ObsStatus       string      `json:"obs_status"`
	Decimal         int         `json:"decimal"`
}

// ParseResponse parses the World Bank API JSON response.
// The World Bank API returns: [metadata, [observations]]
func ParseResponse(data []byte) (*ParsedData, error) {
	// World Bank returns an array with 2 elements: [metadata, data]
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	if len(raw) < 2 {
		return nil, fmt.Errorf("unexpected response format: expected 2 elements, got %d", len(raw))
	}

	// Parse the observations (second element)
	var observations []observation
	if err := json.Unmarshal(raw[1], &observations); err != nil {
		return nil, fmt.Errorf("parse observations: %w", err)
	}

	// Extract dates and values, filtering out null values
	type dataPoint struct {
		date  string
		value string
	}
	var points []dataPoint

	for _, obs := range observations {
		// Skip null values
		if obs.Value == nil {
			continue
		}

		// Format the value properly (handle large numbers without scientific notation)
		var valueStr string
		switch v := obs.Value.(type) {
		case float64:
			// Use %.0f to avoid scientific notation for large numbers
			valueStr = fmt.Sprintf("%.0f", v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}

		points = append(points, dataPoint{
			date:  obs.Date,
			value: valueStr,
		})
	}

	// World Bank returns data in descending order (newest first)
	// Sort to get ascending order (oldest first)
	sort.SliceStable(points, func(i, j int) bool {
		return points[i].date < points[j].date
	})

	// Extract sorted dates and values
	result := &ParsedData{
		Dates:  make([]string, len(points)),
		Values: make([]string, len(points)),
	}
	for i, p := range points {
		result.Dates[i] = p.date
		result.Values[i] = p.value
	}

	return result, nil
}
