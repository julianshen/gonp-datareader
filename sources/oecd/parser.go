package oecd

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// ParsedData holds parsed OECD data.
type ParsedData struct {
	Dates  []string
	Values []float64
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
		result := make([]string, len(p.Values))
		for i, value := range p.Values {
			result[i] = fmt.Sprintf("%g", value)
		}
		return result
	default:
		return nil
	}
}

// sdmxResponse represents the SDMX-JSON structure returned by OECD API.
type sdmxResponse struct {
	Header struct {
		ID       string `json:"id"`
		Prepared string `json:"prepared"`
	} `json:"header"`
	DataSets []struct {
		Observations map[string][]float64 `json:"observations"`
	} `json:"dataSets"`
	Structure struct {
		Dimensions struct {
			Observation []struct {
				ID     string `json:"id"`
				Values []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"values"`
			} `json:"observation"`
		} `json:"dimensions"`
	} `json:"structure"`
}

// ParseJSON parses OECD SDMX-JSON response data.
func ParseJSON(reader io.Reader) (*ParsedData, error) {
	var resp sdmxResponse

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Find time dimension index
	timeDimIndex := -1
	var timePeriods []string
	for i, dim := range resp.Structure.Dimensions.Observation {
		if dim.ID == "TIME_PERIOD" {
			timeDimIndex = i
			for _, val := range dim.Values {
				timePeriods = append(timePeriods, val.ID)
			}
			break
		}
	}

	if timeDimIndex == -1 {
		return nil, fmt.Errorf("TIME_PERIOD dimension not found")
	}

	// Extract observations
	// The key format is "dim1:dim2:dim3:timeIdx" where each number is an index
	// into the corresponding dimension's values array
	observations := make(map[string]float64)

	if len(resp.DataSets) > 0 {
		for key, values := range resp.DataSets[0].Observations {
			if len(values) > 0 {
				// Parse the key to get dimension indices
				indices := strings.Split(key, ":")
				if len(indices) > timeDimIndex {
					timeIdx, err := strconv.Atoi(indices[timeDimIndex])
					if err == nil && timeIdx < len(timePeriods) {
						timePeriod := timePeriods[timeIdx]
						observations[timePeriod] = values[0]
					}
				}
			}
		}
	}

	// Sort by time period
	dates := make([]string, 0, len(observations))
	for date := range observations {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Build result arrays
	values := make([]float64, len(dates))
	for i, date := range dates {
		values[i] = observations[date]
	}

	return &ParsedData{
		Dates:  dates,
		Values: values,
	}, nil
}
