package fred

import (
	"encoding/json"
	"fmt"
	"io"
)

// ParsedData holds parsed FRED data.
type ParsedData struct {
	Dates  []string
	Values []string
}

// GetColumn returns a column of data by name.
// Supported column names: "Date", "Value"
func (p *ParsedData) GetColumn(name string) []string {
	if p == nil {
		return nil
	}

	switch name {
	case "Date":
		return p.Dates
	case "Value":
		return p.Values
	default:
		return nil
	}
}

// fredResponse represents the JSON structure returned by FRED API.
type fredResponse struct {
	ErrorCode    int           `json:"error_code"`
	ErrorMessage string        `json:"error_message"`
	Observations []observation `json:"observations"`
}

// observation represents a single data point from FRED.
type observation struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// ParseJSON parses FRED JSON response data.
func ParseJSON(reader io.Reader) (*ParsedData, error) {
	var resp fredResponse

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Check for API errors
	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("FRED API error: %s", resp.ErrorMessage)
	}

	// Parse observations
	dates := make([]string, 0, len(resp.Observations))
	values := make([]string, 0, len(resp.Observations))

	for _, obs := range resp.Observations {
		// Skip missing values (represented as ".")
		if obs.Value == "." {
			continue
		}

		dates = append(dates, obs.Date)
		values = append(values, obs.Value)
	}

	return &ParsedData{
		Dates:  dates,
		Values: values,
	}, nil
}
