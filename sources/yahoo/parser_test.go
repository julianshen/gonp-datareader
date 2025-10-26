package yahoo_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

func TestParseCSV(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.239990,300.600006,295.190002,300.350006,297.450287,33911900
2020-01-03,297.149994,300.579987,296.500000,297.429993,294.558075,36607600
2020-01-06,293.790009,299.959991,292.750000,299.799988,296.906128,29596800`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseCSV() returned nil result")
	}

	// Should have 3 rows (excluding header)
	rows := result.Rows
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(rows))
	}

	// Check columns exist
	cols := result.Columns
	expectedCols := []string{"Date", "Open", "High", "Low", "Close", "Adj Close", "Volume"}
	if len(cols) != len(expectedCols) {
		t.Errorf("Expected %d columns, got %d", len(expectedCols), len(cols))
	}

	for i, expected := range expectedCols {
		if i >= len(cols) || cols[i] != expected {
			t.Errorf("Expected column %d to be '%s', got '%s'", i, expected, cols[i])
		}
	}
}

func TestParseCSV_EmptyData(t *testing.T) {
	csvData := ``

	_, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err == nil {
		t.Error("ParseCSV() should return error for empty data")
	}
}

func TestParseCSV_HeaderOnly(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	// Should have 0 rows (only header)
	if len(result.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(result.Rows))
	}
}

func TestParseCSV_ParsesNumbers(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(result.Rows))
	}

	row := result.Rows[0]

	// Verify data values
	if row["Date"] != "2020-01-02" {
		t.Errorf("Expected Date '2020-01-02', got '%v'", row["Date"])
	}

	if row["Open"] != "296.24" {
		t.Errorf("Expected Open '296.24', got '%v'", row["Open"])
	}

	if row["Volume"] != "33911900" {
		t.Errorf("Expected Volume '33911900', got '%v'", row["Volume"])
	}
}

func TestParseCSV_HandlesNullValues(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,null,300.60,295.19,null,297.45,33911900
2020-01-03,296.24,null,295.19,300.35,null,33911900`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	if len(result.Rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(result.Rows))
	}

	// Null values should be handled (kept as "null" or empty)
	row1 := result.Rows[0]
	if row1["Open"] != "null" && row1["Open"] != "" {
		t.Logf("Null value handling: Open = %v", row1["Open"])
	}
}

func TestParsedData_GetColumn(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900
2020-01-03,297.15,300.58,296.50,297.43,294.56,36607600`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	// Test getting a column
	dates := result.GetColumn("Date")
	if len(dates) != 2 {
		t.Errorf("Expected 2 dates, got %d", len(dates))
	}

	if dates[0] != "2020-01-02" {
		t.Errorf("Expected first date '2020-01-02', got '%v'", dates[0])
	}

	if dates[1] != "2020-01-03" {
		t.Errorf("Expected second date '2020-01-03', got '%v'", dates[1])
	}
}

func TestParsedData_GetColumn_NonExistent(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.24,300.60,295.19,300.35,297.45,33911900`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	// Getting non-existent column should return empty slice or nil
	values := result.GetColumn("NonExistent")
	if values != nil && len(values) != 0 {
		t.Errorf("Expected empty/nil for non-existent column, got %v", values)
	}
}

// Benchmark tests

func BenchmarkParseCSV(b *testing.B) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.239990,300.600006,295.190002,300.350006,297.450287,33911900
2020-01-03,297.149994,300.579987,296.500000,297.429993,294.558075,36607600
2020-01-06,293.790009,299.959991,292.750000,299.799988,296.906128,29596800
2020-01-07,302.799988,305.130005,300.690002,304.940002,302.001221,33125300
2020-01-08,303.190002,303.239990,299.429993,300.790009,297.888916,28239600`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := yahoo.ParseCSV(strings.NewReader(csvData))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseCSV_LargeDataset(b *testing.B) {
	// Generate a larger CSV dataset (100 rows)
	var sb strings.Builder
	sb.WriteString("Date,Open,High,Low,Close,Adj Close,Volume\n")
	for i := 0; i < 100; i++ {
		sb.WriteString("2020-01-02,296.239990,300.600006,295.190002,300.350006,297.450287,33911900\n")
	}
	csvData := sb.String()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := yahoo.ParseCSV(strings.NewReader(csvData))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetColumn(b *testing.B) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2020-01-02,296.239990,300.600006,295.190002,300.350006,297.450287,33911900
2020-01-03,297.149994,300.579987,296.500000,297.429993,294.558075,36607600
2020-01-06,293.790009,299.959991,292.750000,299.799988,296.906128,29596800`

	result, err := yahoo.ParseCSV(strings.NewReader(csvData))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = result.GetColumn("Close")
	}
}
