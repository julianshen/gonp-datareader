package finmind

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FinMindResponse represents the JSON response from FinMind API.
type FinMindResponse struct {
	Data []FinMindStockData `json:"data"`
}

// FinMindStockData represents a single stock data entry from FinMind.
type FinMindStockData struct {
	Date            string  `json:"date"`
	StockID         string  `json:"stock_id"`
	TradingVolume   int64   `json:"Trading_Volume"`
	TradingMoney    int64   `json:"Trading_money"`
	Open            float64 `json:"open"`
	Max             float64 `json:"max"`
	Min             float64 `json:"min"`
	Close           float64 `json:"close"`
	Spread          float64 `json:"spread"`
	TradingTurnover int64   `json:"Trading_turnover"`
}

// ParsedData represents parsed stock data in a tabular format.
//
// This structure is compatible with the existing datareader pattern
// and provides easy access to stock data for analysis.
type ParsedData struct {
	Symbol  string              // Stock symbol
	Columns []string            // Column names
	Rows    []map[string]string // Data rows (as string maps for flexibility)
}

// ParseFinMindResponse parses the JSON response from FinMind API.
//
// The response contains a "data" array with stock information. Each entry
// includes date, stock_id, OHLCV data, volume, and trading statistics.
//
// Returns ParsedData with columns and rows suitable for analysis.
// Returns an error if the JSON is malformed.
//
// Example:
//
//	data, err := ParseFinMindResponse(responseBody)
//	if err != nil {
//	    return nil, err
//	}
//	fmt.Printf("Symbol: %s, Rows: %d\n", data.Symbol, len(data.Rows))
func ParseFinMindResponse(body []byte) (*ParsedData, error) {
	var response FinMindResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	// Handle empty data
	if len(response.Data) == 0 {
		return &ParsedData{
			Symbol:  "",
			Columns: []string{},
			Rows:    []map[string]string{},
		}, nil
	}

	// Define columns (matching FinMind API response fields)
	columns := []string{
		"date",
		"stock_id",
		"Trading_Volume",
		"Trading_money",
		"open",
		"max",
		"min",
		"close",
		"spread",
		"Trading_turnover",
	}

	// Extract symbol from first row
	symbol := response.Data[0].StockID

	// Convert data to rows
	rows := make([]map[string]string, 0, len(response.Data))
	for _, entry := range response.Data {
		row := map[string]string{
			"date":             entry.Date,
			"stock_id":         entry.StockID,
			"Trading_Volume":   strconv.FormatInt(entry.TradingVolume, 10),
			"Trading_money":    strconv.FormatInt(entry.TradingMoney, 10),
			"open":             formatFloat(entry.Open),
			"max":              formatFloat(entry.Max),
			"min":              formatFloat(entry.Min),
			"close":            formatFloat(entry.Close),
			"spread":           formatFloat(entry.Spread),
			"Trading_turnover": strconv.FormatInt(entry.TradingTurnover, 10),
		}
		rows = append(rows, row)
	}

	return &ParsedData{
		Symbol:  symbol,
		Columns: columns,
		Rows:    rows,
	}, nil
}

// formatFloat converts a float64 to string, removing unnecessary decimals.
func formatFloat(f float64) string {
	// Check if the float is actually an integer
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
