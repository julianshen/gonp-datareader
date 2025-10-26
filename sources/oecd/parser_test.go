package oecd_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/oecd"
)

func TestParseJSON(t *testing.T) {
	jsonData := `{
		"header": {
			"id": "test",
			"prepared": "2020-01-01T00:00:00Z"
		},
		"dataSets": [{
			"observations": {
				"0:0:0:0": [100.0],
				"0:0:0:1": [101.5],
				"0:0:0:2": [102.3]
			}
		}],
		"structure": {
			"dimensions": {
				"observation": [
					{
						"id": "LOCATION",
						"values": [{"id": "USA"}]
					},
					{
						"id": "INDICATOR",
						"values": [{"id": "GDP"}]
					},
					{
						"id": "MEASURE",
						"values": [{"id": "IDX"}]
					},
					{
						"id": "TIME_PERIOD",
						"values": [
							{"id": "2020-Q1"},
							{"id": "2020-Q2"},
							{"id": "2020-Q3"}
						]
					}
				]
			}
		}
	}`

	reader := strings.NewReader(jsonData)
	data, err := oecd.ParseJSON(reader)

	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if data == nil {
		t.Fatal("ParseJSON() returned nil")
	}

	// Check dates
	if len(data.Dates) != 3 {
		t.Errorf("Expected 3 dates, got %d", len(data.Dates))
	}

	expectedDates := []string{"2020-Q1", "2020-Q2", "2020-Q3"}
	for i, expected := range expectedDates {
		if i < len(data.Dates) && data.Dates[i] != expected {
			t.Errorf("Date[%d] = %s, want %s", i, data.Dates[i], expected)
		}
	}

	// Check values
	if len(data.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(data.Values))
	}

	expectedValues := []float64{100.0, 101.5, 102.3}
	for i, expected := range expectedValues {
		if i < len(data.Values) && data.Values[i] != expected {
			t.Errorf("Value[%d] = %f, want %f", i, data.Values[i], expected)
		}
	}
}

func TestParseJSON_EmptyData(t *testing.T) {
	jsonData := `{
		"header": {"id": "test"},
		"dataSets": [{"observations": {}}],
		"structure": {
			"dimensions": {
				"observation": [
					{"id": "TIME_PERIOD", "values": []}
				]
			}
		}
	}`

	reader := strings.NewReader(jsonData)
	data, err := oecd.ParseJSON(reader)

	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if len(data.Dates) != 0 {
		t.Errorf("Expected 0 dates for empty data, got %d", len(data.Dates))
	}

	if len(data.Values) != 0 {
		t.Errorf("Expected 0 values for empty data, got %d", len(data.Values))
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json`

	reader := strings.NewReader(jsonData)
	_, err := oecd.ParseJSON(reader)

	if err == nil {
		t.Error("ParseJSON() should return error for invalid JSON")
	}
}

func TestParsedData_GetColumn(t *testing.T) {
	data := &oecd.ParsedData{
		Dates:  []string{"2020-Q1", "2020-Q2"},
		Values: []float64{100.0, 101.5},
	}

	tests := []struct {
		name     string
		column   string
		expected []string
	}{
		{
			name:     "Date column",
			column:   "Date",
			expected: []string{"2020-Q1", "2020-Q2"},
		},
		{
			name:     "Value column",
			column:   "Value",
			expected: []string{"100", "101.5"},
		},
		{
			name:     "Unknown column",
			column:   "Unknown",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.GetColumn(tt.column)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil for unknown column, got %v", result)
				}
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Column[%d] = %s, want %s", i, result[i], expected)
				}
			}
		})
	}
}

func BenchmarkParseJSON(b *testing.B) {
	jsonData := `{
		"header": {"id": "test"},
		"dataSets": [{"observations": {"0:0:0:0": [100.0], "0:0:0:1": [101.5]}}],
		"structure": {
			"dimensions": {
				"observation": [
					{"id": "TIME_PERIOD", "values": [{"id": "2020-Q1"}, {"id": "2020-Q2"}]}
				]
			}
		}
	}`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(jsonData)
		_, err := oecd.ParseJSON(reader)
		if err != nil {
			b.Fatalf("ParseJSON() error = %v", err)
		}
	}
}
