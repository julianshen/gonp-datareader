package alphavantage_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/alphavantage"
)

func TestParseResponse(t *testing.T) {
	// Alpha Vantage returns JSON with "Time Series (Daily)" key
	jsonData := `{
		"Meta Data": {
			"1. Information": "Daily Prices (open, high, low, close) and Volumes",
			"2. Symbol": "AAPL",
			"3. Last Refreshed": "2024-01-15",
			"4. Output Size": "Compact",
			"5. Time Zone": "US/Eastern"
		},
		"Time Series (Daily)": {
			"2024-01-15": {
				"1. open": "185.89",
				"2. high": "186.95",
				"3. low": "185.30",
				"4. close": "186.51",
				"5. volume": "51234567"
			},
			"2024-01-12": {
				"1. open": "182.16",
				"2. high": "185.92",
				"3. low": "182.00",
				"4. close": "185.59",
				"5. volume": "70538800"
			},
			"2024-01-11": {
				"1. open": "180.21",
				"2. high": "182.41",
				"3. low": "179.50",
				"4. close": "182.32",
				"5. volume": "64280300"
			}
		}
	}`

	result, err := alphavantage.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if len(result.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(result.Rows))
	}

	// Check columns exist
	expectedCols := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	if len(result.Columns) != len(expectedCols) {
		t.Errorf("Expected %d columns, got %d", len(expectedCols), len(result.Columns))
	}

	for _, col := range expectedCols {
		found := false
		for _, c := range result.Columns {
			if c == col {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing expected column %q", col)
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

func TestParseResponse_EmptyTimeSeries(t *testing.T) {
	jsonData := `{
		"Meta Data": {
			"1. Information": "Daily Prices",
			"2. Symbol": "INVALID"
		},
		"Time Series (Daily)": {}
	}`

	result, err := alphavantage.ParseResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if len(result.Rows) != 0 {
		t.Errorf("Expected 0 rows for empty time series, got %d", len(result.Rows))
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	_, err := alphavantage.ParseResponse([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestParseResponse_RateLimitError(t *testing.T) {
	// Alpha Vantage returns this when rate limit is exceeded
	jsonData := `{
		"Note": "Thank you for using Alpha Vantage! Our standard API call frequency is 5 calls per minute and 500 calls per day."
	}`

	_, err := alphavantage.ParseResponse([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for rate limit response")
	}

	if err != nil && err.Error() != "rate limit exceeded" {
		t.Errorf("Expected 'rate limit exceeded' error, got %q", err.Error())
	}
}
