package twse

import (
	"testing"
	"time"
)

// TestRocToGregorian tests conversion from ROC date string to Gregorian time.Time
func TestRocToGregorian(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "ROC 1141031 = 2025-10-31",
			rocDate:   "1141031",
			wantYear:  2025,
			wantMonth: time.October,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "ROC 1130101 = 2024-01-01",
			rocDate:   "1130101",
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC 1131231 = 2024-12-31",
			rocDate:   "1131231",
			wantYear:  2024,
			wantMonth: time.December,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "ROC 1120229 = 2023-02-28 (not leap year)",
			rocDate:   "1120229",
			wantYear:  0, // Invalid date
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "ROC 1120228 = 2023-02-28",
			rocDate:   "1120228",
			wantYear:  2023,
			wantMonth: time.February,
			wantDay:   28,
			wantErr:   false,
		},
		{
			name:      "ROC year 100 = 2011",
			rocDate:   "1000101",
			wantYear:  2011,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "empty string",
			rocDate:   "",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - too short",
			rocDate:   "11410",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - too long",
			rocDate:   "11410311",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - non-numeric",
			rocDate:   "abc1031",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid month",
			rocDate:   "1141331",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid day",
			rocDate:   "1140132",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rocToGregorian(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("rocToGregorian(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear {
					t.Errorf("rocToGregorian(%q) year = %d, want %d", tt.rocDate, got.Year(), tt.wantYear)
				}
				if got.Month() != tt.wantMonth {
					t.Errorf("rocToGregorian(%q) month = %v, want %v", tt.rocDate, got.Month(), tt.wantMonth)
				}
				if got.Day() != tt.wantDay {
					t.Errorf("rocToGregorian(%q) day = %d, want %d", tt.rocDate, got.Day(), tt.wantDay)
				}
			}
		})
	}
}

// TestGregorianToROC tests conversion from Gregorian time.Time to ROC date string
func TestGregorianToROC(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			want: "1141031",
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1130101",
		},
		{
			name: "2024-12-31",
			date: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			want: "1131231",
		},
		{
			name: "2023-02-28",
			date: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			want: "1120228",
		},
		{
			name: "2011-01-01 (ROC 100)",
			date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1000101",
		},
		{
			name: "2020-02-29 (leap year)",
			date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			want: "1090229",
		},
		{
			name: "1912-01-01 (ROC year 1)",
			date: time.Date(1912, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "0010101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gregorianToROC(tt.date)
			if got != tt.want {
				t.Errorf("gregorianToROC(%v) = %q, want %q", tt.date, got, tt.want)
			}
		})
	}
}

// TestParseROCDate tests parsing ROC date strings to time.Time
func TestParseROCDate(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "valid date 1141031",
			rocDate:   "1141031",
			wantYear:  2025,
			wantMonth: time.October,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "valid date 1130101",
			rocDate:   "1130101",
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "invalid format",
			rocDate:   "invalid",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "empty string",
			rocDate:   "",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseROCDate(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseROCDate(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear {
					t.Errorf("parseROCDate(%q) year = %d, want %d", tt.rocDate, got.Year(), tt.wantYear)
				}
				if got.Month() != tt.wantMonth {
					t.Errorf("parseROCDate(%q) month = %v, want %v", tt.rocDate, got.Month(), tt.wantMonth)
				}
				if got.Day() != tt.wantDay {
					t.Errorf("parseROCDate(%q) day = %d, want %d", tt.rocDate, got.Day(), tt.wantDay)
				}
			}
		})
	}
}

// TestFormatROCDate tests formatting time.Time to ROC date string
func TestFormatROCDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			want: "1141031",
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1130101",
		},
		{
			name: "single digit month and day",
			date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			want: "1130105",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatROCDate(tt.date)
			if got != tt.want {
				t.Errorf("formatROCDate(%v) = %q, want %q", tt.date, got, tt.want)
			}
		})
	}
}

// TestROCDateRoundTrip tests that conversion is reversible
func TestROCDateRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "2020-02-29 (leap year)",
			date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to ROC and back
			rocStr := formatROCDate(tt.date)
			got, err := parseROCDate(rocStr)
			if err != nil {
				t.Fatalf("parseROCDate(%q) error = %v", rocStr, err)
			}

			// Compare year, month, day (ignore time components)
			if got.Year() != tt.date.Year() || got.Month() != tt.date.Month() || got.Day() != tt.date.Day() {
				t.Errorf("round trip failed: %v -> %s -> %v", tt.date, rocStr, got)
			}
		})
	}
}

// TestROCDateEdgeCases tests edge cases like leap years and year boundaries
func TestROCDateEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "leap year Feb 29, 2024",
			rocDate:   "1130229",
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
			wantErr:   false,
		},
		{
			name:      "leap year Feb 29, 2020",
			rocDate:   "1090229",
			wantYear:  2020,
			wantMonth: time.February,
			wantDay:   29,
			wantErr:   false,
		},
		{
			name:      "non-leap year Feb 29, 2023 (invalid)",
			rocDate:   "1120229",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "year boundary - Dec 31",
			rocDate:   "1131231",
			wantYear:  2024,
			wantMonth: time.December,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "year boundary - Jan 1",
			rocDate:   "1140101",
			wantYear:  2025,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC year 100 (2011)",
			rocDate:   "1000101",
			wantYear:  2011,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC year 1 (1912)",
			rocDate:   "0010101",
			wantYear:  1912,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rocToGregorian(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("rocToGregorian(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear || got.Month() != tt.wantMonth || got.Day() != tt.wantDay {
					t.Errorf("rocToGregorian(%q) = %v, want %d-%02d-%02d",
						tt.rocDate, got, tt.wantYear, tt.wantMonth, tt.wantDay)
				}
			}
		})
	}
}

// TestParseDailyStockJSON tests parsing TWSE daily stock JSON response
func TestParseDailyStockJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    int
		wantErr bool
	}{
		{
			name: "valid single stock",
			json: `[{
				"Date": "1141031",
				"Code": "2330",
				"Name": "台積電",
				"TradeVolume": "55956524",
				"TradeValue": "3616991558",
				"OpeningPrice": "64.60",
				"HighestPrice": "64.80",
				"LowestPrice": "64.40",
				"ClosingPrice": "64.75",
				"Change": "0.3500",
				"Transaction": "44302"
			}]`,
			want:    1,
			wantErr: false,
		},
		{
			name: "valid multiple stocks",
			json: `[
				{
					"Date": "1141031",
					"Code": "2330",
					"Name": "台積電",
					"TradeVolume": "55956524",
					"OpeningPrice": "64.60",
					"HighestPrice": "64.80",
					"LowestPrice": "64.40",
					"ClosingPrice": "64.75",
					"Change": "0.3500",
					"Transaction": "44302"
				},
				{
					"Date": "1141031",
					"Code": "2317",
					"Name": "鴻海",
					"TradeVolume": "12345678",
					"OpeningPrice": "100.00",
					"HighestPrice": "102.50",
					"LowestPrice": "99.00",
					"ClosingPrice": "101.50",
					"Change": "1.5000",
					"Transaction": "10000"
				}
			]`,
			want:    2,
			wantErr: false,
		},
		{
			name:    "empty array",
			json:    `[]`,
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			json:    `invalid json`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			json:    ``,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDailyStockJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDailyStockJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("parseDailyStockJSON() got %d stocks, want %d", len(got), tt.want)
			}
		})
	}
}

// TestParseDailyStockJSON_FieldExtraction tests that all fields are extracted correctly
func TestParseDailyStockJSON_FieldExtraction(t *testing.T) {
	json := `[{
		"Date": "1141031",
		"Code": "2330",
		"Name": "台積電",
		"TradeVolume": "55956524",
		"TradeValue": "3616991558",
		"OpeningPrice": "64.60",
		"HighestPrice": "64.80",
		"LowestPrice": "64.40",
		"ClosingPrice": "64.75",
		"Change": "0.3500",
		"Transaction": "44302"
	}]`

	stocks, err := parseDailyStockJSON([]byte(json))
	if err != nil {
		t.Fatalf("parseDailyStockJSON() error = %v", err)
	}

	if len(stocks) != 1 {
		t.Fatalf("parseDailyStockJSON() got %d stocks, want 1", len(stocks))
	}

	stock := stocks[0]

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Date", stock.Date, "1141031"},
		{"Code", stock.Code, "2330"},
		{"Name", stock.Name, "台積電"},
		{"TradeVolume", stock.TradeVolume, "55956524"},
		{"TradeValue", stock.TradeValue, "3616991558"},
		{"OpeningPrice", stock.OpeningPrice, "64.60"},
		{"HighestPrice", stock.HighestPrice, "64.80"},
		{"LowestPrice", stock.LowestPrice, "64.40"},
		{"ClosingPrice", stock.ClosingPrice, "64.75"},
		{"Change", stock.Change, "0.3500"},
		{"Transaction", stock.Transaction, "44302"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

// TestParseStockData tests converting TWSEStockData to ParsedData
func TestParseStockData(t *testing.T) {
	tests := []struct {
		name    string
		stock   TWSEStockData
		wantErr bool
		check   func(*testing.T, *ParsedData)
	}{
		{
			name: "valid stock data",
			stock: TWSEStockData{
				Date:         "1141031",
				Code:         "2330",
				Name:         "台積電",
				TradeVolume:  "55956524",
				OpeningPrice: "64.60",
				HighestPrice: "64.80",
				LowestPrice:  "64.40",
				ClosingPrice: "64.75",
				Change:       "0.3500",
				Transaction:  "44302",
			},
			wantErr: false,
			check: func(t *testing.T, p *ParsedData) {
				if p.Symbol != "2330" {
					t.Errorf("Symbol = %q, want %q", p.Symbol, "2330")
				}
				if p.Name != "台積電" {
					t.Errorf("Name = %q, want %q", p.Name, "台積電")
				}
				if len(p.Date) != 1 {
					t.Errorf("Date length = %d, want 1", len(p.Date))
				}
				if len(p.Open) != 1 || p.Open[0] != 64.60 {
					t.Errorf("Open = %v, want [64.60]", p.Open)
				}
				if len(p.High) != 1 || p.High[0] != 64.80 {
					t.Errorf("High = %v, want [64.80]", p.High)
				}
				if len(p.Low) != 1 || p.Low[0] != 64.40 {
					t.Errorf("Low = %v, want [64.40]", p.Low)
				}
				if len(p.Close) != 1 || p.Close[0] != 64.75 {
					t.Errorf("Close = %v, want [64.75]", p.Close)
				}
				if len(p.Volume) != 1 || p.Volume[0] != 55956524 {
					t.Errorf("Volume = %v, want [55956524]", p.Volume)
				}
				if len(p.Transactions) != 1 || p.Transactions[0] != 44302 {
					t.Errorf("Transactions = %v, want [44302]", p.Transactions)
				}
				if len(p.Change) != 1 || p.Change[0] != 0.3500 {
					t.Errorf("Change = %v, want [0.3500]", p.Change)
				}
			},
		},
		{
			name: "empty values",
			stock: TWSEStockData{
				Date:         "1141031",
				Code:         "2330",
				Name:         "台積電",
				TradeVolume:  "",
				OpeningPrice: "",
				HighestPrice: "",
				LowestPrice:  "",
				ClosingPrice: "",
				Change:       "",
				Transaction:  "",
			},
			wantErr: false,
			check: func(t *testing.T, p *ParsedData) {
				if p.Open[0] != 0 {
					t.Errorf("Empty OpeningPrice should be 0, got %v", p.Open[0])
				}
				if p.Volume[0] != 0 {
					t.Errorf("Empty TradeVolume should be 0, got %v", p.Volume[0])
				}
			},
		},
		{
			name: "invalid date",
			stock: TWSEStockData{
				Date:         "invalid",
				Code:         "2330",
				OpeningPrice: "64.60",
				HighestPrice: "64.80",
				LowestPrice:  "64.40",
				ClosingPrice: "64.75",
				TradeVolume:  "1000",
				Transaction:  "100",
				Change:       "0.5",
			},
			wantErr: true,
		},
		{
			name: "invalid price",
			stock: TWSEStockData{
				Date:         "1141031",
				Code:         "2330",
				OpeningPrice: "invalid",
				HighestPrice: "64.80",
				LowestPrice:  "64.40",
				ClosingPrice: "64.75",
				TradeVolume:  "1000",
				Transaction:  "100",
				Change:       "0.5",
			},
			wantErr: true,
		},
		{
			name: "invalid volume",
			stock: TWSEStockData{
				Date:         "1141031",
				Code:         "2330",
				OpeningPrice: "64.60",
				HighestPrice: "64.80",
				LowestPrice:  "64.40",
				ClosingPrice: "64.75",
				TradeVolume:  "invalid",
				Transaction:  "100",
				Change:       "0.5",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStockData(tt.stock)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStockData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

// TestParseFloat tests string to float conversion
func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"valid integer", "100", 100.0, false},
		{"valid decimal", "64.75", 64.75, false},
		{"valid negative", "-10.5", -10.5, false},
		{"empty string", "", 0, false},
		{"invalid string", "abc", 0, true},
		{"invalid format", "12.34.56", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseFloat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseInt tests string to int conversion
func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"valid integer", "12345", 12345, false},
		{"valid large number", "55956524", 55956524, false},
		{"empty string", "", 0, false},
		{"invalid string", "abc", 0, true},
		{"decimal number", "12.34", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInt(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseInt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
