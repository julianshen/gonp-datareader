package tpex

import (
	"math"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "gregorian slash date",
			input: "2025/10/31",
			want:  time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "gregorian dash date",
			input: "2025-10-31",
			want:  time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "roc compact date",
			input: "1141031",
			want:  time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "invalid date",
			input:   "bad-date",
			wantErr: true,
		},
		{
			name:    "invalid roc month",
			input:   "1141331",
			wantErr: true,
		},
		{
			name:    "non numeric roc day",
			input:   "11410aa",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTPEXDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseTPEXDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Fatalf("parseTPEXDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{input: "64.75", want: 64.75},
		{input: "+0.35", want: 0.35},
		{input: "1,234.50", want: 1234.50},
	}

	for _, tt := range tests {
		got, err := parseFloat(tt.input)
		if err != nil {
			t.Fatalf("parseFloat(%q) error = %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("parseFloat(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseNumberSentinelsAsNaN(t *testing.T) {
	for _, input := range []string{"", "-", "--", "---", " --- "} {
		got, err := parseFloat(input)
		if err != nil {
			t.Fatalf("parseFloat(%q) error = %v", input, err)
		}
		if !math.IsNaN(got) {
			t.Fatalf("parseFloat(%q) = %v, want NaN", input, got)
		}
	}
}

func TestParseJSONErrors(t *testing.T) {
	if _, err := parseMainboardJSON([]byte("bad")); err == nil {
		t.Fatal("parseMainboardJSON() error = nil, want error")
	}
	if _, err := parseEmergingJSON([]byte("bad")); err == nil {
		t.Fatal("parseEmergingJSON() error = nil, want error")
	}
	if _, err := parseIndexJSON([]byte("bad")); err == nil {
		t.Fatal("parseIndexJSON() error = nil, want error")
	}
}

func TestParseInvalidNumbers(t *testing.T) {
	if _, err := parseFloat("not-a-number"); err == nil {
		t.Fatal("parseFloat() error = nil, want error")
	}
	if _, err := parseInt("1.25"); err == nil {
		t.Fatal("parseInt() error = nil, want error")
	}
}

func TestParseMainboardStockData(t *testing.T) {
	stock := tpexMainboardQuote{
		Date:                  "2025/10/31",
		SecuritiesCompanyCode: "8069",
		CompanyName:           "元太",
		Open:                  "212.50",
		High:                  "218.00",
		Low:                   "210.00",
		Close:                 "216.50",
		Change:                "+4.00",
		TradingShares:         "12,345,000",
		TransactionNumber:     "8,901",
	}

	got, err := parseMainboardData(stock)
	if err != nil {
		t.Fatalf("parseMainboardData() error = %v", err)
	}

	if got.Symbol != "8069" || got.Name != "元太" {
		t.Fatalf("symbol/name = %q/%q", got.Symbol, got.Name)
	}
	if got.Open[0] != 212.50 || got.High[0] != 218.00 || got.Low[0] != 210.00 || got.Close[0] != 216.50 {
		t.Fatalf("OHLC = %v/%v/%v/%v", got.Open, got.High, got.Low, got.Close)
	}
	if got.Change[0] != 4.00 || got.Volume[0] != 12345000 || got.Transactions[0] != 8901 {
		t.Fatalf("change/volume/transactions = %v/%v/%v", got.Change, got.Volume, got.Transactions)
	}
}

func TestParseMainboardStockDataErrors(t *testing.T) {
	stock := tpexMainboardQuote{
		Date:                  "2025/10/31",
		SecuritiesCompanyCode: "8069",
		Open:                  "bad",
	}
	if _, err := parseMainboardData(stock); err == nil {
		t.Fatal("parseMainboardData() error = nil, want error")
	}
}

func TestParseEmergingStockData(t *testing.T) {
	stock := tpexEmergingQuote{
		Date:                  "2025/10/31",
		SecuritiesCompanyCode: "6871",
		CompanyName:           "訊芯",
		Average:               "95.10",
		Highest:               "98.20",
		Lowest:                "93.00",
		LatestPrice:           "97.50",
		TransactionVolume:     "1,250",
	}

	got, err := parseEmergingData(stock)
	if err != nil {
		t.Fatalf("parseEmergingData() error = %v", err)
	}

	if got.Symbol != "esb:6871" || got.Name != "訊芯" {
		t.Fatalf("symbol/name = %q/%q", got.Symbol, got.Name)
	}
	if !math.IsNaN(got.Open[0]) || got.High[0] != 98.20 || got.Low[0] != 93.00 || got.Close[0] != 97.50 {
		t.Fatalf("OHLC = %v/%v/%v/%v", got.Open, got.High, got.Low, got.Close)
	}
	if got.Volume[0] != 1250 {
		t.Fatalf("volume = %v, want 1250", got.Volume)
	}
}

func TestParseEmergingStockDataErrors(t *testing.T) {
	tests := []struct {
		name  string
		stock tpexEmergingQuote
	}{
		{
			name: "bad high",
			stock: tpexEmergingQuote{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "6871",
				Highest:               "bad",
			},
		},
		{
			name: "bad low",
			stock: tpexEmergingQuote{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "6871",
				Highest:               "2",
				Lowest:                "bad",
			},
		},
		{
			name: "bad latest price",
			stock: tpexEmergingQuote{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "6871",
				Highest:               "2",
				Lowest:                "1",
				LatestPrice:           "bad",
			},
		},
		{
			name: "bad transaction volume",
			stock: tpexEmergingQuote{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "6871",
				Highest:               "2",
				Lowest:                "1",
				LatestPrice:           "1.5",
				TransactionVolume:     "bad",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseEmergingData(tt.stock); err == nil {
				t.Fatal("parseEmergingData() error = nil, want error")
			}
		})
	}
}

func TestParseIndexData(t *testing.T) {
	row := tpexIndexData{
		Date:   "2025/10/31",
		Open:   "250.10",
		High:   "252.20",
		Low:    "249.30",
		Close:  "251.80",
		Change: "+1.20",
	}

	got, err := parseIndexData(row)
	if err != nil {
		t.Fatalf("parseIndexData() error = %v", err)
	}

	if got.Symbol != "index" || got.Name != "TPEx Index" {
		t.Fatalf("symbol/name = %q/%q", got.Symbol, got.Name)
	}
	if got.Open[0] != 250.10 || got.High[0] != 252.20 || got.Low[0] != 249.30 || got.Close[0] != 251.80 {
		t.Fatalf("OHLC = %v/%v/%v/%v", got.Open, got.High, got.Low, got.Close)
	}
	if got.Change[0] != 1.20 {
		t.Fatalf("change = %v, want 1.20", got.Change)
	}
}

func TestParseIndexDataErrors(t *testing.T) {
	if _, err := parseIndexData(tpexIndexData{Date: "bad"}); err == nil {
		t.Fatal("parseIndexData() error = nil, want error")
	}
}

func TestFilterEmergingBySymbolNotFound(t *testing.T) {
	_, err := filterEmergingBySymbol([]tpexEmergingQuote{{SecuritiesCompanyCode: "6871"}}, "esb:9999")
	if err == nil {
		t.Fatal("filterEmergingBySymbol() error = nil, want error")
	}
}

func TestFilterByDateRange(t *testing.T) {
	data := &ParsedData{
		Symbol:       "index",
		Name:         "TPEx Index",
		Date:         []time.Time{time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)},
		Open:         []float64{1, 2},
		High:         []float64{3, 4},
		Low:          []float64{5, 6},
		Close:        []float64{7, 8},
		Volume:       []int64{9, 10},
		Transactions: []int64{11, 12},
		Change:       []float64{13, 14},
	}

	got := filterByDateRange(data, time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC))
	if len(got.Date) != 1 || got.Close[0] != 8 {
		t.Fatalf("filtered data = %+v, want one row with close 8", got)
	}
}
