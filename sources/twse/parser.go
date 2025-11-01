package twse

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	// rocEpochYear is the offset between ROC and Gregorian calendars
	// ROC Year 1 = Gregorian Year 1912
	rocEpochYear = 1911

	// Expected ROC date format: YYYMMDD (7 digits)
	// YYY = ROC year (3 digits)
	// MM = month (2 digits)
	// DD = day (2 digits)
	rocDateLength = 7
)

// rocToGregorian converts a ROC (Republic of China) date string to a Gregorian time.Time.
//
// ROC dates are formatted as "YYYMMDD" where:
//   - YYY is the ROC year (ROC Year = Gregorian Year - 1911)
//   - MM is the month (01-12)
//   - DD is the day (01-31)
//
// Examples:
//   - "1141031" -> October 31, 2025 (ROC 114 + 1911 = 2025)
//   - "1130101" -> January 1, 2024 (ROC 113 + 1911 = 2024)
//
// The function validates the date and returns an error if:
//   - The format is invalid (not 7 digits)
//   - The date components are invalid (e.g., month 13, day 32)
//   - The date doesn't exist (e.g., Feb 29 in non-leap year)
func rocToGregorian(rocDate string) (time.Time, error) {
	if len(rocDate) != rocDateLength {
		return time.Time{}, fmt.Errorf("invalid ROC date format: expected 7 digits (YYYMMDD), got %d digits", len(rocDate))
	}

	// Parse ROC year (first 3 digits)
	rocYearStr := rocDate[0:3]
	rocYear, err := strconv.Atoi(rocYearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ROC year: %w", err)
	}

	// Parse month (next 2 digits)
	monthStr := rocDate[3:5]
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month: %w", err)
	}

	// Parse day (last 2 digits)
	dayStr := rocDate[5:7]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %w", err)
	}

	// Convert ROC year to Gregorian year
	gregorianYear := rocYear + rocEpochYear

	// Create time.Time and let it validate the date
	// This will catch invalid dates like Feb 30, April 31, etc.
	date := time.Date(gregorianYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Verify the date is valid by checking if components match
	// If date is invalid, time.Date normalizes it (e.g., Feb 30 -> March 2)
	if date.Year() != gregorianYear || date.Month() != time.Month(month) || date.Day() != day {
		return time.Time{}, fmt.Errorf("invalid date: ROC %s (Gregorian %d-%02d-%02d does not exist)",
			rocDate, gregorianYear, month, day)
	}

	return date, nil
}

// gregorianToROC converts a Gregorian time.Time to a ROC date string.
//
// The returned string is in "YYYMMDD" format where:
//   - YYY is the ROC year (Gregorian Year - 1911)
//   - MM is the month (01-12)
//   - DD is the day (01-31)
//
// Examples:
//   - October 31, 2025 -> "1141031" (2025 - 1911 = 114)
//   - January 1, 2024 -> "1130101" (2024 - 1911 = 113)
func gregorianToROC(date time.Time) string {
	rocYear := date.Year() - rocEpochYear
	return fmt.Sprintf("%03d%02d%02d", rocYear, date.Month(), date.Day())
}

// parseROCDate parses a ROC date string into a time.Time.
//
// This is an alias for rocToGregorian for consistency with common naming patterns.
func parseROCDate(rocDate string) (time.Time, error) {
	return rocToGregorian(rocDate)
}

// formatROCDate formats a time.Time into a ROC date string.
//
// This is an alias for gregorianToROC for consistency with common naming patterns.
func formatROCDate(date time.Time) string {
	return gregorianToROC(date)
}

// TWSEStockData represents a single stock's data in the TWSE API response.
//
// All numeric fields are returned as strings by the API and need to be
// parsed to appropriate numeric types.
type TWSEStockData struct {
	Date         string `json:"Date"`         // ROC date format "YYYMMDD"
	Code         string `json:"Code"`         // Stock symbol (e.g., "2330")
	Name         string `json:"Name"`         // Company name in Traditional Chinese
	TradeVolume  string `json:"TradeVolume"`  // Number of shares traded
	TradeValue   string `json:"TradeValue"`   // Total trade value
	OpeningPrice string `json:"OpeningPrice"` // Opening price
	HighestPrice string `json:"HighestPrice"` // Daily high
	LowestPrice  string `json:"LowestPrice"`  // Daily low
	ClosingPrice string `json:"ClosingPrice"` // Closing price
	Change       string `json:"Change"`       // Price change
	Transaction  string `json:"Transaction"`  // Number of transactions
}

// ParsedData represents parsed stock data ready for use.
//
// This structure contains typed data with time.Time dates and numeric values
// converted from the API's string format.
type ParsedData struct {
	Symbol       string      // Stock symbol
	Name         string      // Company name
	Date         []time.Time // Trading dates
	Open         []float64   // Opening prices
	High         []float64   // Highest prices
	Low          []float64   // Lowest prices
	Close        []float64   // Closing prices
	Volume       []int64     // Trading volumes
	Transactions []int64     // Transaction counts
	Change       []float64   // Price changes
}

// parseDailyStockJSON parses the TWSE daily stock data JSON response.
//
// The TWSE API returns an array of stock data objects where all numeric
// values are represented as strings. This function:
//   - Parses the JSON array
//   - Converts ROC dates to time.Time
//   - Converts string numbers to appropriate numeric types
//   - Handles missing/empty values
//
// Example input:
//
//	[{
//	  "Date": "1141031",
//	  "Code": "2330",
//	  "Name": "台積電",
//	  "TradeVolume": "55956524",
//	  "OpeningPrice": "64.60",
//	  "HighestPrice": "64.80",
//	  "LowestPrice": "64.40",
//	  "ClosingPrice": "64.75",
//	  "Change": "0.3500",
//	  "Transaction": "44302"
//	}]
func parseDailyStockJSON(data []byte) ([]TWSEStockData, error) {
	var stocks []TWSEStockData
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}
	return stocks, nil
}

// parseStockData converts a single TWSEStockData to ParsedData.
//
// This function handles:
//   - ROC date to time.Time conversion
//   - String to float64 conversion for prices
//   - String to int64 conversion for volumes
//   - Empty/missing value handling
func parseStockData(stock TWSEStockData) (*ParsedData, error) {
	// Parse date
	date, err := parseROCDate(stock.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", stock.Date, err)
	}

	// Parse prices
	open, err := parseFloat(stock.OpeningPrice)
	if err != nil {
		return nil, fmt.Errorf("parse opening price %q: %w", stock.OpeningPrice, err)
	}

	high, err := parseFloat(stock.HighestPrice)
	if err != nil {
		return nil, fmt.Errorf("parse highest price %q: %w", stock.HighestPrice, err)
	}

	low, err := parseFloat(stock.LowestPrice)
	if err != nil {
		return nil, fmt.Errorf("parse lowest price %q: %w", stock.LowestPrice, err)
	}

	close, err := parseFloat(stock.ClosingPrice)
	if err != nil {
		return nil, fmt.Errorf("parse closing price %q: %w", stock.ClosingPrice, err)
	}

	change, err := parseFloat(stock.Change)
	if err != nil {
		return nil, fmt.Errorf("parse change %q: %w", stock.Change, err)
	}

	// Parse volumes
	volume, err := parseInt(stock.TradeVolume)
	if err != nil {
		return nil, fmt.Errorf("parse trade volume %q: %w", stock.TradeVolume, err)
	}

	transactions, err := parseInt(stock.Transaction)
	if err != nil {
		return nil, fmt.Errorf("parse transactions %q: %w", stock.Transaction, err)
	}

	return &ParsedData{
		Symbol:       stock.Code,
		Name:         stock.Name,
		Date:         []time.Time{date},
		Open:         []float64{open},
		High:         []float64{high},
		Low:          []float64{low},
		Close:        []float64{close},
		Volume:       []int64{volume},
		Transactions: []int64{transactions},
		Change:       []float64{change},
	}, nil
}

// parseFloat converts a string to float64, handling empty strings.
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %w", err)
	}
	return f, nil
}

// parseInt converts a string to int64, handling empty strings.
func parseInt(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %w", err)
	}
	return i, nil
}

// filterBySymbol finds a specific stock symbol in the array of stocks.
//
// Returns the matching TWSEStockData or an error if the symbol is not found.
// This is used to extract data for a single symbol from the API response
// which returns all stocks.
func filterBySymbol(stocks []TWSEStockData, symbol string) (TWSEStockData, error) {
	if symbol == "" {
		return TWSEStockData{}, fmt.Errorf("symbol cannot be empty")
	}

	for _, stock := range stocks {
		if stock.Code == symbol {
			return stock, nil
		}
	}

	return TWSEStockData{}, fmt.Errorf("symbol %q not found in response", symbol)
}

// filterByDateRange filters ParsedData to include only dates within the specified range.
//
// The filtering is inclusive: both start and end dates are included if present.
// Returns a new ParsedData with filtered data, preserving all slices in sync.
func filterByDateRange(data *ParsedData, start, end time.Time) *ParsedData {
	if data == nil || len(data.Date) == 0 {
		return &ParsedData{
			Symbol: data.Symbol,
			Name:   data.Name,
		}
	}

	// Pre-allocate slices for efficiency
	filtered := &ParsedData{
		Symbol:       data.Symbol,
		Name:         data.Name,
		Date:         make([]time.Time, 0, len(data.Date)),
		Open:         make([]float64, 0, len(data.Date)),
		High:         make([]float64, 0, len(data.Date)),
		Low:          make([]float64, 0, len(data.Date)),
		Close:        make([]float64, 0, len(data.Date)),
		Volume:       make([]int64, 0, len(data.Date)),
		Transactions: make([]int64, 0, len(data.Date)),
		Change:       make([]float64, 0, len(data.Date)),
	}

	// Filter data within date range (inclusive)
	for i, date := range data.Date {
		// Compare dates (ignore time component)
		dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		startOnly := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		endOnly := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

		if (dateOnly.Equal(startOnly) || dateOnly.After(startOnly)) &&
			(dateOnly.Equal(endOnly) || dateOnly.Before(endOnly)) {
			// Date is within range, include all data for this index
			filtered.Date = append(filtered.Date, data.Date[i])
			filtered.Open = append(filtered.Open, data.Open[i])
			filtered.High = append(filtered.High, data.High[i])
			filtered.Low = append(filtered.Low, data.Low[i])
			filtered.Close = append(filtered.Close, data.Close[i])
			filtered.Volume = append(filtered.Volume, data.Volume[i])
			filtered.Transactions = append(filtered.Transactions, data.Transactions[i])
			filtered.Change = append(filtered.Change, data.Change[i])
		}
	}

	return filtered
}
