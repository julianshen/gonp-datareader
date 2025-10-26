// +build integration

package datareader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader"
	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

// TestIntegration_EndToEnd demonstrates complete end-to-end functionality
// with a mock Yahoo Finance server.
func TestIntegration_EndToEnd(t *testing.T) {
	csvData := `Date,Open,High,Low,Close,Adj Close,Volume
2023-01-03,125.07,125.42,124.17,125.07,123.45,112117500
2023-01-04,126.89,128.66,125.08,126.36,124.72,89113600
2023-01-05,127.13,127.77,124.76,125.02,123.41,80962700
2023-01-06,126.01,130.29,124.89,129.62,127.95,87754700`

	// Create mock Yahoo Finance server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	t.Log("Mock Yahoo Finance server:", server.URL)

	// Create reader with mock URL
	clientOpts := internalhttp.DefaultClientOptions()
	reader := yahoo.NewYahooReaderWithBaseURL(clientOpts, server.URL+"/%s")

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	// Test 1: Fetch single symbol
	t.Run("fetch_single_symbol", func(t *testing.T) {
		data, err := reader.ReadSingle(ctx, "AAPL", start, end)
		if err != nil {
			t.Fatalf("Failed to fetch data: %v", err)
		}

		parsedData := data.(*yahoo.ParsedData)
		if len(parsedData.Rows) != 4 {
			t.Errorf("Expected 4 rows, got %d", len(parsedData.Rows))
		}

		t.Logf("✓ Fetched %d days of data", len(parsedData.Rows))

		// Verify columns
		expectedCols := []string{"Date", "Open", "High", "Low", "Close", "Adj Close", "Volume"}
		if len(parsedData.Columns) != len(expectedCols) {
			t.Errorf("Expected %d columns, got %d", len(expectedCols), len(parsedData.Columns))
		}

		// Verify data values
		if parsedData.Rows[0]["Date"] != "2023-01-03" {
			t.Errorf("Expected first date '2023-01-03', got '%s'", parsedData.Rows[0]["Date"])
		}

		if parsedData.Rows[0]["Close"] != "125.07" {
			t.Errorf("Expected close '125.07', got '%s'", parsedData.Rows[0]["Close"])
		}

		t.Log("✓ Data values are correct")
	})

	// Test 2: Extract column data
	t.Run("extract_columns", func(t *testing.T) {
		data, err := reader.ReadSingle(ctx, "AAPL", start, end)
		if err != nil {
			t.Fatalf("Failed to fetch data: %v", err)
		}

		parsedData := data.(*yahoo.ParsedData)

		// Extract closing prices
		closes := parsedData.GetColumn("Close")
		if len(closes) != 4 {
			t.Errorf("Expected 4 close prices, got %d", len(closes))
		}

		expectedCloses := []string{"125.07", "126.36", "125.02", "129.62"}
		for i, expected := range expectedCloses {
			if closes[i] != expected {
				t.Errorf("Close price %d: expected %s, got %s", i, expected, closes[i])
			}
		}

		t.Logf("✓ Closing prices: %v", closes)

		// Extract volumes
		volumes := parsedData.GetColumn("Volume")
		if len(volumes) != 4 {
			t.Errorf("Expected 4 volumes, got %d", len(volumes))
		}

		t.Logf("✓ Volumes: %v", volumes)
	})

	// Test 3: Multiple symbols
	t.Run("fetch_multiple_symbols", func(t *testing.T) {
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		results, err := reader.Read(ctx, symbols, start, end)
		if err != nil {
			t.Fatalf("Failed to fetch multiple symbols: %v", err)
		}

		dataMap := results.(map[string]*yahoo.ParsedData)
		if len(dataMap) != 3 {
			t.Errorf("Expected 3 results, got %d", len(dataMap))
		}

		for _, symbol := range symbols {
			if data, ok := dataMap[symbol]; !ok {
				t.Errorf("Missing data for symbol %s", symbol)
			} else {
				t.Logf("✓ %s: %d rows", symbol, len(data.Rows))
			}
		}
	})

	t.Log("✓ Integration test completed successfully!")
}

// TestIntegration_RealYahooFinance tests against the real Yahoo Finance API.
// This test may fail due to rate limiting or network issues.
// Run with: go test -tags=integration -run TestIntegration_RealYahooFinance
func TestIntegration_RealYahooFinance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real API test in short mode")
	}

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	// Add delay to avoid rate limiting
	time.Sleep(2 * time.Second)

	data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
	if err != nil {
		t.Logf("Real API test failed (expected due to rate limiting): %v", err)
		t.Skip("Skipping due to API error - this is normal")
		return
	}

	parsedData := data.(*yahoo.ParsedData)
	t.Logf("✓ Successfully fetched %d days of real data from Yahoo Finance", len(parsedData.Rows))

	if len(parsedData.Rows) > 0 {
		t.Logf("First row: Date=%s, Close=%s",
			parsedData.Rows[0]["Date"],
			parsedData.Rows[0]["Close"])
	}
}
