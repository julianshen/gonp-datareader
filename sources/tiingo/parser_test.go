package tiingo_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/tiingo"
)

func TestParseJSON(t *testing.T) {
	jsonData := `[
		{
			"date": "2020-01-02T00:00:00.000Z",
			"close": 300.35,
			"high": 300.60,
			"low": 295.19,
			"open": 296.24,
			"volume": 33911900,
			"adjClose": 297.45,
			"adjHigh": 300.60,
			"adjLow": 295.19,
			"adjOpen": 296.24,
			"adjVolume": 33911900,
			"divCash": 0.0,
			"splitFactor": 1.0
		},
		{
			"date": "2020-01-03T00:00:00.000Z",
			"close": 297.43,
			"high": 300.58,
			"low": 296.50,
			"open": 297.15,
			"volume": 36607600,
			"adjClose": 294.56,
			"adjHigh": 300.58,
			"adjLow": 296.50,
			"adjOpen": 297.15,
			"adjVolume": 36607600,
			"divCash": 0.0,
			"splitFactor": 1.0
		}
	]`

	data, err := tiingo.ParseJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if data == nil {
		t.Fatal("ParseJSON returned nil data")
	}

	if len(data.Dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(data.Dates))
	}

	if len(data.Prices) != 2 {
		t.Errorf("Expected 2 prices, got %d", len(data.Prices))
	}

	// Check first record
	if data.Dates[0] != "2020-01-02" {
		t.Errorf("Expected first date '2020-01-02', got '%s'", data.Dates[0])
	}

	firstPrice := data.Prices[0]
	if firstPrice.Close != 300.35 {
		t.Errorf("Expected close 300.35, got %f", firstPrice.Close)
	}

	if firstPrice.Open != 296.24 {
		t.Errorf("Expected open 296.24, got %f", firstPrice.Open)
	}

	if firstPrice.High != 300.60 {
		t.Errorf("Expected high 300.60, got %f", firstPrice.High)
	}

	if firstPrice.Low != 295.19 {
		t.Errorf("Expected low 295.19, got %f", firstPrice.Low)
	}

	if firstPrice.Volume != 33911900 {
		t.Errorf("Expected volume 33911900, got %d", firstPrice.Volume)
	}
}

func TestParseJSON_EmptyData(t *testing.T) {
	jsonData := `[]`

	data, err := tiingo.ParseJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(data.Dates) != 0 {
		t.Errorf("Expected 0 dates for empty data, got %d", len(data.Dates))
	}

	if len(data.Prices) != 0 {
		t.Errorf("Expected 0 prices for empty data, got %d", len(data.Prices))
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json}`

	_, err := tiingo.ParseJSON(strings.NewReader(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParsedData_GetColumn(t *testing.T) {
	data := &tiingo.ParsedData{
		Dates: []string{"2020-01-02", "2020-01-03"},
		Prices: []tiingo.PriceData{
			{Close: 300.35, Open: 296.24, High: 300.60, Low: 295.19, Volume: 33911900},
			{Close: 297.43, Open: 297.15, High: 300.58, Low: 296.50, Volume: 36607600},
		},
	}

	dates := data.GetColumn("Date")
	if len(dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(dates))
	}

	if dates[0] != "2020-01-02" {
		t.Errorf("Expected first date '2020-01-02', got '%s'", dates[0])
	}

	closes := data.GetColumn("Close")
	if len(closes) != 2 {
		t.Errorf("Expected 2 closes, got %d", len(closes))
	}

	if closes[0] != "300.35" {
		t.Errorf("Expected first close '300.35', got '%s'", closes[0])
	}

	// Non-existent column should return nil
	unknown := data.GetColumn("Unknown")
	if unknown != nil {
		t.Errorf("Expected nil for unknown column, got %v", unknown)
	}
}

// Benchmark tests
func BenchmarkParseJSON(b *testing.B) {
	jsonData := `[
		{"date": "2020-01-02T00:00:00.000Z", "close": 300.35, "high": 300.60, "low": 295.19, "open": 296.24, "volume": 33911900, "adjClose": 297.45},
		{"date": "2020-01-03T00:00:00.000Z", "close": 297.43, "high": 300.58, "low": 296.50, "open": 297.15, "volume": 36607600, "adjClose": 294.56}
	]`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tiingo.ParseJSON(strings.NewReader(jsonData))
		if err != nil {
			b.Fatal(err)
		}
	}
}
