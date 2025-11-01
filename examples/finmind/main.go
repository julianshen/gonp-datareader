// Package main demonstrates usage of the FinMind data reader.
//
// FinMind (https://finmind.github.io/) provides comprehensive financial data
// for Taiwan and international markets with over 50 datasets. It offers
// historical data since 1994 for Taiwan stocks, along with fundamental data,
// institutional investor data, and more.
//
// Features:
//   - Optional Bearer token authentication for higher rate limits
//   - 300 requests/hour without token, 600 requests/hour with token
//   - 50+ datasets including stocks, futures, options, bonds, and commodities
//   - Historical data since 1994 for Taiwan stocks
//   - International market coverage (US stocks, commodities, currencies)
//
// To run this example:
//
//	go run examples/finmind/main.go
//
// To use with authentication token:
//
//	FINMIND_TOKEN=your-token go run examples/finmind/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/gonp-datareader"
)

func main() {
	fmt.Println("gonp-datareader - FinMind Example")
	fmt.Println("===================================")
	fmt.Println()

	// Get API token from environment (optional)
	token := os.Getenv("FINMIND_TOKEN")
	if token != "" {
		fmt.Println("Using FinMind API token (600 req/hour)")
	} else {
		fmt.Println("No API token - using public access (300 req/hour)")
		fmt.Println("To use a token: export FINMIND_TOKEN=your-token")
	}
	fmt.Println()

	// Create context with timeout
	ctx := context.Background()

	// Define date range
	start := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC)

	// Example 1: Using the convenience function
	example1ConvenienceFunction(ctx, token, start, end)

	// Example 2: Using the factory pattern
	example2FactoryPattern(ctx, token, start, end)

	// Example 3: Fetching multiple Taiwan stocks
	example3MultipleSymbols(ctx, token, start, end)

	// Example 4: Different datasets
	example4DifferentDatasets(ctx, token)

	// Example 5: Error handling
	example5ErrorHandling(ctx)
}

// example1ConvenienceFunction demonstrates the simplest way to fetch data.
func example1ConvenienceFunction(ctx context.Context, token string, start, end time.Time) {
	fmt.Println("Example 1: Convenience Function")
	fmt.Println("--------------------------------")

	// Configure options
	opts := &datareader.Options{
		APIKey:  token,
		Timeout: 30 * time.Second,
	}

	// Fetch TSMC (2330) stock data
	result, err := datareader.Read(ctx, "2330", "finmind", start, end, opts)
	if err != nil {
		log.Printf("Error fetching data: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Printf("Successfully fetched TSMC (2330) data for October 2024\n")
	fmt.Printf("Result type: %T\n", result)
	fmt.Println()
}

// example2FactoryPattern demonstrates using the factory pattern for more control.
func example2FactoryPattern(ctx context.Context, token string, start, end time.Time) {
	fmt.Println("Example 2: Factory Pattern")
	fmt.Println("--------------------------")

	// Create reader with factory
	opts := &datareader.Options{
		APIKey:     token,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	reader, err := datareader.DataReader("finmind", opts)
	if err != nil {
		log.Printf("Error creating reader: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Printf("Created FinMind reader: %s (source: %s)\n", reader.Name(), reader.Source())

	// Fetch data for MediaTek (2454)
	result, err := reader.ReadSingle(ctx, "2454", start, end)
	if err != nil {
		log.Printf("Error fetching MediaTek data: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Printf("Successfully fetched MediaTek (2454) data\n")
	fmt.Printf("Result type: %T\n", result)
	fmt.Println()
}

// example3MultipleSymbols demonstrates fetching multiple symbols in parallel.
func example3MultipleSymbols(ctx context.Context, token string, start, end time.Time) {
	fmt.Println("Example 3: Multiple Taiwan Stocks")
	fmt.Println("----------------------------------")

	// Popular Taiwan stocks
	symbols := map[string]string{
		"2330": "TSMC (Taiwan Semiconductor)",
		"2317": "Hon Hai Precision (Foxconn)",
		"2454": "MediaTek",
		"2412": "Chunghwa Telecom",
		"0050": "Yuanta Taiwan Top 50 ETF",
	}

	fmt.Println("Fetching data for popular Taiwan stocks:")
	for code, name := range symbols {
		fmt.Printf("  %s: %s\n", code, name)
	}
	fmt.Println()

	// Create reader
	opts := &datareader.Options{
		APIKey:  token,
		Timeout: 60 * time.Second,
	}

	reader, err := datareader.DataReader("finmind", opts)
	if err != nil {
		log.Printf("Error creating reader: %v\n", err)
		fmt.Println()
		return
	}

	// Fetch all symbols in parallel
	symbolList := []string{"2330", "2317", "2454", "2412", "0050"}
	result, err := reader.Read(ctx, symbolList, start, end)
	if err != nil {
		log.Printf("Error fetching multiple symbols: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Printf("Successfully fetched data for %d symbols\n", len(symbolList))
	fmt.Printf("Result type: %T\n", result)
	fmt.Println()
}

// example4DifferentDatasets demonstrates using different FinMind datasets.
func example4DifferentDatasets(ctx context.Context, token string) {
	fmt.Println("Example 4: Available Datasets")
	fmt.Println("-----------------------------")

	fmt.Println("FinMind supports 50+ datasets including:")
	fmt.Println()

	datasets := []struct {
		name        string
		description string
	}{
		{"TaiwanStockPrice", "Daily stock prices (default)"},
		{"TaiwanStockInfo", "Company information"},
		{"TaiwanStockDividend", "Dividend data"},
		{"TaiwanStockPER", "P/E ratio data"},
		{"TaiwanStockCapital", "Share capital data"},
		{"USStockPrice", "US stock prices"},
		{"TaiwanFuturesDaily", "Futures data"},
		{"TaiwanOptionsDaily", "Options data"},
		{"GoldPrice", "Gold prices"},
		{"CrudeOilPrices", "Crude oil prices"},
	}

	for _, ds := range datasets {
		fmt.Printf("  %-25s %s\n", ds.name, ds.description)
	}

	fmt.Println()
	fmt.Println("Note: This example uses TaiwanStockPrice (default dataset)")
	fmt.Println("To use other datasets, call reader.SetDataset(\"DatasetName\") before fetching")
	fmt.Println()
}

// example5ErrorHandling demonstrates proper error handling.
func example5ErrorHandling(ctx context.Context) {
	fmt.Println("Example 5: Error Handling")
	fmt.Println("-------------------------")

	reader, err := datareader.DataReader("finmind", nil)
	if err != nil {
		log.Printf("Error creating reader: %v\n", err)
		return
	}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	// Test 1: Empty symbol
	fmt.Println("Test 1: Empty symbol (should fail)")
	_, err = reader.ReadSingle(ctx, "", start, end)
	if err != nil {
		fmt.Printf("  ✓ Expected error: %v\n", err)
	} else {
		fmt.Println("  ✗ Should have returned an error")
	}

	// Test 2: Invalid date range
	fmt.Println("\nTest 2: Invalid date range (end before start)")
	invalidStart := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	invalidEnd := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err = reader.ReadSingle(ctx, "2330", invalidStart, invalidEnd)
	if err != nil {
		fmt.Printf("  ✓ Expected error: %v\n", err)
	} else {
		fmt.Println("  ✗ Should have returned an error")
	}

	// Test 3: Valid request (may fail due to network/rate limiting)
	fmt.Println("\nTest 3: Valid request (may fail if no network/rate limited)")
	_, err = reader.ReadSingle(ctx, "2330", start, end)
	if err != nil {
		fmt.Printf("  Network/API error (expected without token): %v\n", err)
	} else {
		fmt.Println("  ✓ Successfully fetched data")
	}

	fmt.Println()
}
