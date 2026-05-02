package tpex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const rocEpochYear = 1911

// TPEXMainboardQuote represents one row from /tpex_mainboard_daily_close_quotes.
type TPEXMainboardQuote struct {
	Date                  string `json:"Date"`
	SecuritiesCompanyCode string `json:"SecuritiesCompanyCode"`
	CompanyName           string `json:"CompanyName"`
	Close                 string `json:"Close"`
	Change                string `json:"Change"`
	Open                  string `json:"Open"`
	High                  string `json:"High"`
	Low                   string `json:"Low"`
	TradingShares         string `json:"TradingShares"`
	TransactionNumber     string `json:"TransactionNumber"`
}

// TPEXEmergingQuote represents one row from /tpex_esb_latest_statistics.
type TPEXEmergingQuote struct {
	Date                  string `json:"Date"`
	SecuritiesCompanyCode string `json:"SecuritiesCompanyCode"`
	CompanyName           string `json:"CompanyName"`
	Highest               string `json:"Highest"`
	Lowest                string `json:"Lowest"`
	Average               string `json:"Average"`
	LatestPrice           string `json:"LatestPrice"`
	TransactionVolume     string `json:"TransactionVolume"`
}

// TPEXIndexData represents one row from /tpex_index.
type TPEXIndexData struct {
	Date   string `json:"Date"`
	Open   string `json:"Open"`
	High   string `json:"High"`
	Low    string `json:"Low"`
	Close  string `json:"Close"`
	Change string `json:"Change"`
}

// ParsedData contains typed TPEX data.
type ParsedData struct {
	Symbol       string
	Name         string
	Date         []time.Time
	Open         []float64
	High         []float64
	Low          []float64
	Close        []float64
	Volume       []int64
	Transactions []int64
	Change       []float64
}

func parseMainboardJSON(data []byte) ([]TPEXMainboardQuote, error) {
	var rows []TPEXMainboardQuote
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("unmarshal mainboard JSON: %w", err)
	}
	return rows, nil
}

func parseEmergingJSON(data []byte) ([]TPEXEmergingQuote, error) {
	var rows []TPEXEmergingQuote
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("unmarshal emerging JSON: %w", err)
	}
	return rows, nil
}

func parseIndexJSON(data []byte) ([]TPEXIndexData, error) {
	var rows []TPEXIndexData
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("unmarshal index JSON: %w", err)
	}
	return rows, nil
}

func parseMainboardData(row TPEXMainboardQuote) (*ParsedData, error) {
	date, err := parseTPEXDate(row.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", row.Date, err)
	}
	open, high, low, close, change, err := parseOHLCChange(row.Open, row.High, row.Low, row.Close, row.Change)
	if err != nil {
		return nil, err
	}
	volume, err := parseInt(row.TradingShares)
	if err != nil {
		return nil, fmt.Errorf("parse trading shares %q: %w", row.TradingShares, err)
	}
	transactions, err := parseInt(row.TransactionNumber)
	if err != nil {
		return nil, fmt.Errorf("parse transaction number %q: %w", row.TransactionNumber, err)
	}

	return oneRow(row.SecuritiesCompanyCode, row.CompanyName, date, open, high, low, close, volume, transactions, change), nil
}

func parseEmergingData(row TPEXEmergingQuote) (*ParsedData, error) {
	date, err := parseTPEXDate(row.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", row.Date, err)
	}
	open, high, low, close, _, err := parseOHLCChange(row.Average, row.Highest, row.Lowest, row.LatestPrice, "")
	if err != nil {
		return nil, err
	}
	volume, err := parseInt(row.TransactionVolume)
	if err != nil {
		return nil, fmt.Errorf("parse transaction volume %q: %w", row.TransactionVolume, err)
	}

	return oneRow("esb:"+row.SecuritiesCompanyCode, row.CompanyName, date, open, high, low, close, volume, 0, 0), nil
}

func parseIndexData(row TPEXIndexData) (*ParsedData, error) {
	date, err := parseTPEXDate(row.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", row.Date, err)
	}
	open, high, low, close, change, err := parseOHLCChange(row.Open, row.High, row.Low, row.Close, row.Change)
	if err != nil {
		return nil, err
	}

	return oneRow("index", "TPEx Index", date, open, high, low, close, 0, 0, change), nil
}

func parseOHLCChange(openRaw, highRaw, lowRaw, closeRaw, changeRaw string) (float64, float64, float64, float64, float64, error) {
	open, err := parseFloat(openRaw)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("parse open %q: %w", openRaw, err)
	}
	high, err := parseFloat(highRaw)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("parse high %q: %w", highRaw, err)
	}
	low, err := parseFloat(lowRaw)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("parse low %q: %w", lowRaw, err)
	}
	close, err := parseFloat(closeRaw)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("parse close %q: %w", closeRaw, err)
	}
	change, err := parseFloat(changeRaw)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("parse change %q: %w", changeRaw, err)
	}
	return open, high, low, close, change, nil
}

func oneRow(symbol, name string, date time.Time, open, high, low, close float64, volume, transactions int64, change float64) *ParsedData {
	return &ParsedData{
		Symbol:       symbol,
		Name:         name,
		Date:         []time.Time{date},
		Open:         []float64{open},
		High:         []float64{high},
		Low:          []float64{low},
		Close:        []float64{close},
		Volume:       []int64{volume},
		Transactions: []int64{transactions},
		Change:       []float64{change},
	}
}

func parseTPEXDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"2006/01/02", "2006-01-02"} {
		if date, err := time.ParseInLocation(layout, value, time.UTC); err == nil {
			return date, nil
		}
	}

	if len(value) == 7 {
		rocYear, err := strconv.Atoi(value[:3])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid ROC year: %w", err)
		}
		month, err := strconv.Atoi(value[3:5])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid month: %w", err)
		}
		day, err := strconv.Atoi(value[5:7])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day: %w", err)
		}
		year := rocYear + rocEpochYear
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if date.Year() != year || int(date.Month()) != month || date.Day() != day {
			return time.Time{}, fmt.Errorf("invalid date %q", value)
		}
		return date, nil
	}

	return time.Time{}, fmt.Errorf("unsupported date format %q", value)
}

func parseFloat(value string) (float64, error) {
	normalized := normalizeNumber(value)
	if normalized == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %w", err)
	}
	return parsed, nil
}

func parseInt(value string) (int64, error) {
	normalized := normalizeNumber(value)
	if normalized == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(normalized, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %w", err)
	}
	return parsed, nil
}

func normalizeNumber(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, ",", "")
	value = strings.TrimPrefix(value, "+")
	if value == "-" || value == "--" {
		return ""
	}
	return value
}

func filterMainboardBySymbol(rows []TPEXMainboardQuote, symbol string) (TPEXMainboardQuote, error) {
	for _, row := range rows {
		if row.SecuritiesCompanyCode == symbol {
			return row, nil
		}
	}
	return TPEXMainboardQuote{}, fmt.Errorf("symbol %q not found in response", symbol)
}

func filterEmergingBySymbol(rows []TPEXEmergingQuote, symbol string) (TPEXEmergingQuote, error) {
	code := strings.TrimPrefix(symbol, "esb:")
	for _, row := range rows {
		if row.SecuritiesCompanyCode == code {
			return row, nil
		}
	}
	return TPEXEmergingQuote{}, fmt.Errorf("symbol %q not found in response", symbol)
}

func filterByDateRange(data *ParsedData, start, end time.Time) *ParsedData {
	if data == nil {
		return &ParsedData{}
	}
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
	start = dayOnly(start)
	end = dayOnly(end)
	for i, date := range data.Date {
		current := dayOnly(date)
		if current.Before(start) || current.After(end) {
			continue
		}
		filtered.Date = append(filtered.Date, data.Date[i])
		filtered.Open = append(filtered.Open, data.Open[i])
		filtered.High = append(filtered.High, data.High[i])
		filtered.Low = append(filtered.Low, data.Low[i])
		filtered.Close = append(filtered.Close, data.Close[i])
		filtered.Volume = append(filtered.Volume, data.Volume[i])
		filtered.Transactions = append(filtered.Transactions, data.Transactions[i])
		filtered.Change = append(filtered.Change, data.Change[i])
	}
	return filtered
}

func dayOnly(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}
