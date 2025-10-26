package fred_test

import (
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/fred"
)

func TestParseJSON(t *testing.T) {
	jsonData := `{
		"realtime_start": "2024-01-01",
		"realtime_end": "2024-01-01",
		"observation_start": "2020-01-01",
		"observation_end": "2020-12-31",
		"units": "lin",
		"output_type": 1,
		"file_type": "json",
		"order_by": "observation_date",
		"sort_order": "asc",
		"count": 4,
		"offset": 0,
		"limit": 100000,
		"observations": [
			{
				"realtime_start": "2024-01-01",
				"realtime_end": "2024-01-01",
				"date": "2020-01-01",
				"value": "21734.056"
			},
			{
				"realtime_start": "2024-01-01",
				"realtime_end": "2024-01-01",
				"date": "2020-04-01",
				"value": "19520.114"
			},
			{
				"realtime_start": "2024-01-01",
				"realtime_end": "2024-01-01",
				"date": "2020-07-01",
				"value": "21170.252"
			},
			{
				"realtime_start": "2024-01-01",
				"realtime_end": "2024-01-01",
				"date": "2020-10-01",
				"value": "21494.731"
			}
		]
	}`

	data, err := fred.ParseJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if data == nil {
		t.Fatal("ParseJSON returned nil data")
	}

	if len(data.Dates) != 4 {
		t.Errorf("Expected 4 dates, got %d", len(data.Dates))
	}

	if len(data.Values) != 4 {
		t.Errorf("Expected 4 values, got %d", len(data.Values))
	}

	// Check first observation
	if data.Dates[0] != "2020-01-01" {
		t.Errorf("Expected first date '2020-01-01', got '%s'", data.Dates[0])
	}

	if data.Values[0] != "21734.056" {
		t.Errorf("Expected first value '21734.056', got '%s'", data.Values[0])
	}

	// Check last observation
	if data.Dates[3] != "2020-10-01" {
		t.Errorf("Expected last date '2020-10-01', got '%s'", data.Dates[3])
	}

	if data.Values[3] != "21494.731" {
		t.Errorf("Expected last value '21494.731', got '%s'", data.Values[3])
	}
}

func TestParseJSON_EmptyData(t *testing.T) {
	jsonData := `{
		"observations": []
	}`

	data, err := fred.ParseJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(data.Dates) != 0 {
		t.Errorf("Expected 0 dates for empty observations, got %d", len(data.Dates))
	}

	if len(data.Values) != 0 {
		t.Errorf("Expected 0 values for empty observations, got %d", len(data.Values))
	}
}

func TestParseJSON_MissingValues(t *testing.T) {
	jsonData := `{
		"observations": [
			{
				"date": "2020-01-01",
				"value": "100.5"
			},
			{
				"date": "2020-01-02",
				"value": "."
			},
			{
				"date": "2020-01-03",
				"value": "102.3"
			}
		]
	}`

	data, err := fred.ParseJSON(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	// Should skip the observation with missing value (".")
	if len(data.Dates) != 2 {
		t.Errorf("Expected 2 dates (skipping missing value), got %d", len(data.Dates))
	}

	if len(data.Values) != 2 {
		t.Errorf("Expected 2 values (skipping missing value), got %d", len(data.Values))
	}

	// Check that we got the non-missing values
	if data.Dates[0] != "2020-01-01" || data.Values[0] != "100.5" {
		t.Errorf("First observation incorrect: date=%s, value=%s", data.Dates[0], data.Values[0])
	}

	if data.Dates[1] != "2020-01-03" || data.Values[1] != "102.3" {
		t.Errorf("Second observation incorrect: date=%s, value=%s", data.Dates[1], data.Values[1])
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json}`

	_, err := fred.ParseJSON(strings.NewReader(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseJSON_APIError(t *testing.T) {
	jsonData := `{
		"error_code": 400,
		"error_message": "Bad Request. The value for variable api_key is not registered."
	}`

	_, err := fred.ParseJSON(strings.NewReader(jsonData))
	if err == nil {
		t.Error("Expected error for API error response, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "Bad Request") {
		t.Errorf("Expected error to contain 'Bad Request', got: %v", err)
	}
}

func TestParsedData_GetColumn(t *testing.T) {
	data := &fred.ParsedData{
		Dates:  []string{"2020-01-01", "2020-01-02", "2020-01-03"},
		Values: []string{"100.5", "101.2", "102.3"},
	}

	dates := data.GetColumn("Date")
	if len(dates) != 3 {
		t.Errorf("Expected 3 dates, got %d", len(dates))
	}

	if dates[0] != "2020-01-01" {
		t.Errorf("Expected first date '2020-01-01', got '%s'", dates[0])
	}

	values := data.GetColumn("Value")
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	if values[0] != "100.5" {
		t.Errorf("Expected first value '100.5', got '%s'", values[0])
	}

	// Non-existent column should return nil
	unknown := data.GetColumn("Unknown")
	if unknown != nil {
		t.Errorf("Expected nil for unknown column, got %v", unknown)
	}
}
