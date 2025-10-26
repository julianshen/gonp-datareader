package eurostat

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ParsedData holds parsed Eurostat data.
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

// jsonStatResponse represents the JSON-stat structure returned by Eurostat API.
type jsonStatResponse struct {
	Version   string                       `json:"version"`
	Class     string                       `json:"class"`
	Label     string                       `json:"label"`
	ID        []string                     `json:"id"`
	Size      []int                        `json:"size"`
	Dimension map[string]jsonStatDimension `json:"dimension"`
	Value     []interface{}                `json:"value"`
}

type jsonStatDimension struct {
	Label    string           `json:"label"`
	Category jsonStatCategory `json:"category"`
}

type jsonStatCategory struct {
	Index map[string]int    `json:"index"`
	Label map[string]string `json:"label"`
}

// ParseJSON parses Eurostat JSON-stat response data.
func ParseJSON(reader io.Reader) (*ParsedData, error) {
	var resp jsonStatResponse

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Find time dimension
	timeDimIndex := -1
	var timeCategories []string
	for i, dimID := range resp.ID {
		if dimID == "time" {
			timeDimIndex = i
			if dim, ok := resp.Dimension[dimID]; ok {
				// Extract time categories sorted by index
				timeMap := make(map[int]string)
				for cat, idx := range dim.Category.Index {
					timeMap[idx] = cat
				}
				// Sort by index
				indices := make([]int, 0, len(timeMap))
				for idx := range timeMap {
					indices = append(indices, idx)
				}
				sort.Ints(indices)
				for _, idx := range indices {
					timeCategories = append(timeCategories, timeMap[idx])
				}
			}
			break
		}
	}

	if timeDimIndex == -1 {
		return nil, fmt.Errorf("time dimension not found")
	}

	if len(timeCategories) == 0 {
		return &ParsedData{
			Dates:  []string{},
			Values: []float64{},
		}, nil
	}

	// Calculate values aggregated by time
	// Values in JSON-stat are in row-major order
	numTimes := len(timeCategories)
	timeValues := make(map[int][]float64)

	// Calculate stride for time dimension
	stride := 1
	for i := len(resp.Size) - 1; i > timeDimIndex; i-- {
		stride *= resp.Size[i]
	}

	// Extract values for each time period
	for i, val := range resp.Value {
		if val == nil {
			continue
		}

		// Calculate which time index this value belongs to
		timeIdx := (i / stride) % numTimes

		// Convert value to float64
		var floatVal float64
		switch v := val.(type) {
		case float64:
			floatVal = v
		case int:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
		default:
			continue
		}

		timeValues[timeIdx] = append(timeValues[timeIdx], floatVal)
	}

	// Build result arrays
	dates := make([]string, numTimes)
	values := make([]float64, numTimes)

	for i := 0; i < numTimes; i++ {
		dates[i] = timeCategories[i]

		// Average values for this time period across other dimensions
		if vals, ok := timeValues[i]; ok && len(vals) > 0 {
			sum := 0.0
			for _, v := range vals {
				sum += v
			}
			values[i] = sum / float64(len(vals))
		}
	}

	return &ParsedData{
		Dates:  dates,
		Values: values,
	}, nil
}
