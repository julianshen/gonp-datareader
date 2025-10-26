package iex_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/iex"
)

func TestParseResponse(t *testing.T) {
	// IEX Cloud returns an array of daily chart data
	jsonData := `[
		{
			"date": "2024-01-02",
			"open": 185.34,
			"high": 186.28,
			"low": 184.21,
			"close": 185.64,
			"volume": 45123000
		},
		{
			"date": "2024-01-03",
			"open": 185.90,
			"high": 187.45,
			"low": 185.32,
			"close": 186.89,
			"volume": 51234000
		}
	]`

	data, err := iex.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if data == nil {
		t.Fatal("ParseResponse() returned nil data")
	}

	expectedColumns := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	if len(data.Columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(data.Columns))
	}

	for i, col := range expectedColumns {
		if data.Columns[i] != col {
			t.Errorf("Expected column %d to be %q, got %q", i, col, data.Columns[i])
		}
	}

	if len(data.Rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(data.Rows))
	}

	// Check first row
	firstRow := data.Rows[0]
	if firstRow["Date"] != "2024-01-02" {
		t.Errorf("Expected date '2024-01-02', got %q", firstRow["Date"])
	}
	if firstRow["Open"] != "185.34" {
		t.Errorf("Expected open '185.34', got %q", firstRow["Open"])
	}
	if firstRow["Close"] != "185.64" {
		t.Errorf("Expected close '185.64', got %q", firstRow["Close"])
	}
	if firstRow["Volume"] != "45123000" {
		t.Errorf("Expected volume '45123000', got %q", firstRow["Volume"])
	}
}

func TestParseResponse_EmptyArray(t *testing.T) {
	jsonData := `[]`

	data, err := iex.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if len(data.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(data.Rows))
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	_, err := iex.ParseResponse([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestParseResponse_SortsByDate(t *testing.T) {
	// Dates in reverse order
	jsonData := `[
		{
			"date": "2024-01-05",
			"open": 190.0,
			"high": 191.0,
			"low": 189.0,
			"close": 190.5,
			"volume": 50000000
		},
		{
			"date": "2024-01-03",
			"open": 186.0,
			"high": 187.0,
			"low": 185.0,
			"close": 186.5,
			"volume": 51000000
		}
	]`

	data, err := iex.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if len(data.Rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(data.Rows))
	}

	// Should be sorted ascending by date
	if data.Rows[0]["Date"] != "2024-01-03" {
		t.Errorf("Expected first date '2024-01-03', got %q", data.Rows[0]["Date"])
	}
	if data.Rows[1]["Date"] != "2024-01-05" {
		t.Errorf("Expected second date '2024-01-05', got %q", data.Rows[1]["Date"])
	}
}

func TestParseResponse_ErrorMessage(t *testing.T) {
	// IEX Cloud might return error messages
	jsonData := `{"error": "Unknown symbol"}`

	_, err := iex.ParseResponse([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for error response")
	}

	if !strings.Contains(err.Error(), "API error") {
		t.Errorf("Expected 'API error' in error message, got %q", err.Error())
	}
}
