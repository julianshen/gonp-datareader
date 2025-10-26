package stooq_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/stooq"
)

func TestParseCSV(t *testing.T) {
	// Stooq CSV format
	csvData := `Date,Open,High,Low,Close,Volume
2024-01-15,185.89,186.95,185.30,186.51,51234567
2024-01-12,182.16,185.92,182.00,185.59,70538800
2024-01-11,180.21,182.41,179.50,182.32,64280300
`

	result, err := stooq.ParseCSV([]byte(csvData))
	if err != nil {
		t.Fatalf("ParseCSV failed: %v", err)
	}

	if len(result.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(result.Rows))
	}

	// Check columns
	expectedCols := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	if len(result.Columns) != len(expectedCols) {
		t.Errorf("Expected %d columns, got %d", len(expectedCols), len(result.Columns))
	}

	for i, expected := range expectedCols {
		if i < len(result.Columns) && result.Columns[i] != expected {
			t.Errorf("Expected column %d to be %q, got %q", i, expected, result.Columns[i])
		}
	}

	// Check first row (sorted ascending by date)
	if len(result.Rows) > 0 {
		firstRow := result.Rows[0]
		if firstRow["Date"] != "2024-01-11" {
			t.Errorf("Expected first date '2024-01-11', got %q", firstRow["Date"])
		}
		if firstRow["Close"] != "182.32" {
			t.Errorf("Expected close '182.32', got %q", firstRow["Close"])
		}
	}
}

func TestParseCSV_EmptyData(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Volume
`

	result, err := stooq.ParseCSV([]byte(csvData))
	if err != nil {
		t.Fatalf("ParseCSV failed: %v", err)
	}

	if len(result.Rows) != 0 {
		t.Errorf("Expected 0 rows for empty data, got %d", len(result.Rows))
	}
}

func TestParseCSV_InvalidCSV(t *testing.T) {
	// Empty data
	csvData := ``

	_, err := stooq.ParseCSV([]byte(csvData))
	if err == nil {
		t.Error("Expected error for empty CSV")
	}
}
