package tiingo

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// PriceData represents a single price record from Tiingo.
type PriceData struct {
	Close  float64
	Open   float64
	High   float64
	Low    float64
	Volume int64
}

// ParsedData holds parsed Tiingo data.
type ParsedData struct {
	Dates  []string
	Prices []PriceData
}

// GetColumn returns a column of data by name.
// Supported column names: "Date", "Close", "Open", "High", "Low", "Volume"
func (p *ParsedData) GetColumn(name string) []string {
	if p == nil {
		return nil
	}

	switch name {
	case "Date":
		return p.Dates
	case "Close":
		result := make([]string, len(p.Prices))
		for i, price := range p.Prices {
			result[i] = fmt.Sprintf("%g", price.Close)
		}
		return result
	case "Open":
		result := make([]string, len(p.Prices))
		for i, price := range p.Prices {
			result[i] = fmt.Sprintf("%g", price.Open)
		}
		return result
	case "High":
		result := make([]string, len(p.Prices))
		for i, price := range p.Prices {
			result[i] = fmt.Sprintf("%g", price.High)
		}
		return result
	case "Low":
		result := make([]string, len(p.Prices))
		for i, price := range p.Prices {
			result[i] = fmt.Sprintf("%g", price.Low)
		}
		return result
	case "Volume":
		result := make([]string, len(p.Prices))
		for i, price := range p.Prices {
			result[i] = fmt.Sprintf("%d", price.Volume)
		}
		return result
	default:
		return nil
	}
}

// tiingoResponse represents the JSON structure returned by Tiingo API.
type tiingoResponse struct {
	Date      string  `json:"date"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Open      float64 `json:"open"`
	Volume    int64   `json:"volume"`
	AdjClose  float64 `json:"adjClose"`
	AdjHigh   float64 `json:"adjHigh"`
	AdjLow    float64 `json:"adjLow"`
	AdjOpen   float64 `json:"adjOpen"`
	AdjVolume int64   `json:"adjVolume"`
	DivCash   float64 `json:"divCash"`
	SplitFactor float64 `json:"splitFactor"`
}

// ParseJSON parses Tiingo JSON response data.
func ParseJSON(reader io.Reader) (*ParsedData, error) {
	var resp []tiingoResponse

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Parse records
	dates := make([]string, 0, len(resp))
	prices := make([]PriceData, 0, len(resp))

	for _, record := range resp {
		// Parse date (format: "2020-01-02T00:00:00.000Z")
		date := record.Date
		if t, err := time.Parse(time.RFC3339, date); err == nil {
			date = t.Format("2006-01-02")
		} else if strings.Contains(date, "T") {
			// Fallback: just take the date part before 'T'
			date = strings.Split(date, "T")[0]
		}

		dates = append(dates, date)
		prices = append(prices, PriceData{
			Close:  record.Close,
			Open:   record.Open,
			High:   record.High,
			Low:    record.Low,
			Volume: record.Volume,
		})
	}

	return &ParsedData{
		Dates:  dates,
		Prices: prices,
	}, nil
}
