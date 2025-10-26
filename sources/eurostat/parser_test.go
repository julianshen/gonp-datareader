package eurostat_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/eurostat"
)

func TestParseJSON(t *testing.T) {
	jsonData := `{
		"version": "2.0",
		"class": "dataset",
		"label": "Population density",
		"id": ["geo", "time"],
		"size": [2, 3],
		"dimension": {
			"geo": {
				"label": "Geopolitical entity",
				"category": {
					"index": {"DE": 0, "FR": 1},
					"label": {"DE": "Germany", "FR": "France"}
				}
			},
			"time": {
				"label": "Time",
				"category": {
					"index": {"2020": 0, "2021": 1, "2022": 2},
					"label": {"2020": "2020", "2021": "2021", "2022": "2022"}
				}
			}
		},
		"value": [100.0, 101.5, 102.3, 200.0, 201.5, 202.3]
	}`

	reader := strings.NewReader(jsonData)
	data, err := eurostat.ParseJSON(reader)

	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if data == nil {
		t.Fatal("ParseJSON() returned nil")
	}

	// Check dates (should have 3 unique time periods)
	if len(data.Dates) != 3 {
		t.Errorf("Expected 3 dates, got %d", len(data.Dates))
	}

	expectedDates := []string{"2020", "2021", "2022"}
	for i, expected := range expectedDates {
		if i < len(data.Dates) && data.Dates[i] != expected {
			t.Errorf("Date[%d] = %s, want %s", i, data.Dates[i], expected)
		}
	}

	// Check values (should be aggregated/averaged across geo dimension)
	if len(data.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(data.Values))
	}

	// Values should be averages: (100+200)/2=150, (101.5+201.5)/2=151.5, (102.3+202.3)/2=152.3
	expectedValues := []float64{150.0, 151.5, 152.3}
	for i, expected := range expectedValues {
		if i < len(data.Values) && data.Values[i] != expected {
			t.Errorf("Value[%d] = %f, want %f", i, data.Values[i], expected)
		}
	}
}

func TestParseJSON_EmptyData(t *testing.T) {
	jsonData := `{
		"version": "2.0",
		"class": "dataset",
		"id": ["time"],
		"size": [0],
		"dimension": {
			"time": {
				"category": {
					"index": {},
					"label": {}
				}
			}
		},
		"value": []
	}`

	reader := strings.NewReader(jsonData)
	data, err := eurostat.ParseJSON(reader)

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
	_, err := eurostat.ParseJSON(reader)

	if err == nil {
		t.Error("ParseJSON() should return error for invalid JSON")
	}
}

func TestParsedData_GetColumn(t *testing.T) {
	data := &eurostat.ParsedData{
		Dates:  []string{"2020", "2021"},
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
			expected: []string{"2020", "2021"},
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
		"version": "2.0",
		"class": "dataset",
		"id": ["geo", "time"],
		"size": [1, 2],
		"dimension": {
			"geo": {"category": {"index": {"EU": 0}, "label": {"EU": "EU"}}},
			"time": {"category": {"index": {"2020": 0, "2021": 1}, "label": {"2020": "2020", "2021": "2021"}}}
		},
		"value": [100.0, 101.5]
	}`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(jsonData)
		_, err := eurostat.ParseJSON(reader)
		if err != nil {
			b.Fatalf("ParseJSON() error = %v", err)
		}
	}
}
