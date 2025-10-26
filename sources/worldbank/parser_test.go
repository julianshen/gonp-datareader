package worldbank_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/worldbank"
)

func TestParseResponse(t *testing.T) {
	// World Bank API returns JSON in format: [metadata, [observations]]
	jsonData := `[
		{
			"page": 1,
			"pages": 1,
			"per_page": "1000",
			"total": 3
		},
		[
			{
				"indicator": {
					"id": "NY.GDP.MKTP.CD",
					"value": "GDP (current US$)"
				},
				"country": {
					"id": "US",
					"value": "United States"
				},
				"countryiso3code": "USA",
				"date": "2022",
				"value": 25462700000000,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			},
			{
				"indicator": {
					"id": "NY.GDP.MKTP.CD",
					"value": "GDP (current US$)"
				},
				"country": {
					"id": "US",
					"value": "United States"
				},
				"countryiso3code": "USA",
				"date": "2021",
				"value": 23315100000000,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			},
			{
				"indicator": {
					"id": "NY.GDP.MKTP.CD",
					"value": "GDP (current US$)"
				},
				"country": {
					"id": "US",
					"value": "United States"
				},
				"countryiso3code": "USA",
				"date": "2020",
				"value": null,
				"unit": "",
				"obs_status": "",
				"decimal": 0
			}
		]
	]`

	result, err := worldbank.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if len(result.Dates) != 2 {
		t.Errorf("Expected 2 dates (null filtered), got %d", len(result.Dates))
	}

	if len(result.Values) != 2 {
		t.Errorf("Expected 2 values (null filtered), got %d", len(result.Values))
	}

	// Check that dates are in order (World Bank returns newest first)
	if result.Dates[0] != "2021" {
		t.Errorf("Expected first date '2021', got %q", result.Dates[0])
	}

	if result.Dates[1] != "2022" {
		t.Errorf("Expected second date '2022', got %q", result.Dates[1])
	}

	// Check values
	if result.Values[0] != "23315100000000" {
		t.Errorf("Expected first value '23315100000000', got %q", result.Values[0])
	}
}

func TestParseResponse_EmptyData(t *testing.T) {
	jsonData := `[
		{
			"page": 1,
			"pages": 0,
			"per_page": "1000",
			"total": 0
		},
		[]
	]`

	result, err := worldbank.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if len(result.Dates) != 0 {
		t.Errorf("Expected 0 dates for empty data, got %d", len(result.Dates))
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	_, err := worldbank.ParseResponse([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
