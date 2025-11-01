package finmind_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/finmind"
)

func TestParseFinMindResponse(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"date": "2020-04-06",
				"stock_id": "2330",
				"Trading_Volume": 59712754,
				"Trading_money": 16324198154,
				"open": 273,
				"max": 275.5,
				"min": 270,
				"close": 275.5,
				"spread": 4,
				"Trading_turnover": 19971
			},
			{
				"date": "2020-04-07",
				"stock_id": "2330",
				"Trading_Volume": 45123456,
				"Trading_money": 12345678901,
				"open": 275,
				"max": 278,
				"min": 274,
				"close": 277,
				"spread": 2,
				"Trading_turnover": 15432
			}
		]
	}`

	data, err := finmind.ParseFinMindResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseFinMindResponse() error = %v", err)
	}

	if data.Symbol != "2330" {
		t.Errorf("Expected symbol '2330', got %q", data.Symbol)
	}

	if len(data.Rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(data.Rows))
	}

	// Check first row
	row1 := data.Rows[0]
	if row1["date"] != "2020-04-06" {
		t.Errorf("Expected date '2020-04-06', got %q", row1["date"])
	}
	if row1["stock_id"] != "2330" {
		t.Errorf("Expected stock_id '2330', got %q", row1["stock_id"])
	}
	if row1["open"] != "273" {
		t.Errorf("Expected open '273', got %q", row1["open"])
	}
	if row1["close"] != "275.5" {
		t.Errorf("Expected close '275.5', got %q", row1["close"])
	}
	if row1["Trading_Volume"] != "59712754" {
		t.Errorf("Expected Trading_Volume '59712754', got %q", row1["Trading_Volume"])
	}

	// Check second row
	row2 := data.Rows[1]
	if row2["date"] != "2020-04-07" {
		t.Errorf("Expected date '2020-04-07', got %q", row2["date"])
	}
	if row2["close"] != "277" {
		t.Errorf("Expected close '277', got %q", row2["close"])
	}
}

func TestParseFinMindResponse_EmptyData(t *testing.T) {
	jsonData := `{"data": []}`

	data, err := finmind.ParseFinMindResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseFinMindResponse() should not error on empty data: %v", err)
	}

	if data.Symbol != "" {
		t.Errorf("Expected empty symbol, got %q", data.Symbol)
	}

	if len(data.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(data.Rows))
	}
}

func TestParseFinMindResponse_InvalidJSON(t *testing.T) {
	jsonData := `{"data": [invalid json]}`

	_, err := finmind.ParseFinMindResponse([]byte(jsonData))
	if err == nil {
		t.Error("ParseFinMindResponse() should error on invalid JSON")
	}
}

func TestParseFinMindResponse_Columns(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"date": "2020-04-06",
				"stock_id": "2330",
				"Trading_Volume": 59712754,
				"Trading_money": 16324198154,
				"open": 273,
				"max": 275.5,
				"min": 270,
				"close": 275.5,
				"spread": 4,
				"Trading_turnover": 19971
			}
		]
	}`

	data, err := finmind.ParseFinMindResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseFinMindResponse() error = %v", err)
	}

	expectedColumns := []string{
		"date", "stock_id", "Trading_Volume", "Trading_money",
		"open", "max", "min", "close", "spread", "Trading_turnover",
	}

	if len(data.Columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(data.Columns))
	}

	// Check all expected columns are present
	columnMap := make(map[string]bool)
	for _, col := range data.Columns {
		columnMap[col] = true
	}

	for _, expected := range expectedColumns {
		if !columnMap[expected] {
			t.Errorf("Missing expected column: %s", expected)
		}
	}
}
