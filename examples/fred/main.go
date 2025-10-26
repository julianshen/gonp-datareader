// Package main demonstrates FRED (Federal Reserve Economic Data) usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/fred"
)

func main() {
	fmt.Println("gonp-datareader - FRED Example")
	fmt.Println("===============================")

	// FRED requires an API key (free at https://fred.stlouisfed.org/docs/api/api_key.html)
	apiKey := os.Getenv("FRED_API_KEY")
	if apiKey == "" {
		fmt.Println("\n⚠️  FRED_API_KEY environment variable not set")
		fmt.Println("Get a free API key at: https://fred.stlouisfed.org/docs/api/api_key.html")
		fmt.Println("\nUsage: FRED_API_KEY=your_key_here go run main.go")
		return
	}

	ctx := context.Background()

	// Fetch quarterly GDP data for 2020-2023
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	fmt.Printf("\nFetching FRED data from %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	// Method 1: Using convenience function
	fmt.Println("\n--- Method 1: Convenience Function ---")

	opts := &datareader.Options{
		APIKey: apiKey,
	}

	data, err := datareader.Read(ctx, "GDP", "fred", start, end, opts)
	if err != nil {
		log.Fatalf("Failed to fetch GDP data: %v", err)
	}

	parsedData := data.(*fred.ParsedData)
	fmt.Printf("✓ Fetched GDP data: %d observations\n", len(parsedData.Dates))

	if len(parsedData.Dates) > 0 {
		fmt.Printf("  First: %s = %s\n", parsedData.Dates[0], parsedData.Values[0])
		lastIdx := len(parsedData.Dates) - 1
		fmt.Printf("  Last:  %s = %s\n", parsedData.Dates[lastIdx], parsedData.Values[lastIdx])
	}

	// Method 2: Using factory pattern
	fmt.Println("\n--- Method 2: Factory Pattern ---")

	reader, err := datareader.DataReader("fred", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	// Fetch 10-Year Treasury Constant Maturity Rate
	data2, err := reader.ReadSingle(ctx, "DGS10", start, end)
	if err != nil {
		log.Fatalf("Failed to fetch DGS10 data: %v", err)
	}

	parsedData2 := data2.(*fred.ParsedData)
	fmt.Printf("✓ Fetched DGS10 (10-Year Treasury): %d observations\n", len(parsedData2.Dates))

	// Method 3: Multiple series
	fmt.Println("\n--- Method 3: Multiple Economic Indicators ---")

	series := []string{
		"GDP",      // Gross Domestic Product
		"UNRATE",   // Unemployment Rate
		"CPIAUCSL", // Consumer Price Index
	}

	results, err := reader.Read(ctx, series, start, end)
	if err != nil {
		log.Fatalf("Failed to fetch multiple series: %v", err)
	}

	dataMap := results.(map[string]*fred.ParsedData)
	for _, seriesID := range series {
		if data, ok := dataMap[seriesID]; ok {
			fmt.Printf("✓ %s: %d observations\n", seriesID, len(data.Dates))

			// Show first value
			if len(data.Dates) > 0 {
				fmt.Printf("  %s: %s\n", data.Dates[0], data.Values[0])
			}
		}
	}

	// Method 4: Using SetAPIKey
	fmt.Println("\n--- Method 4: Direct Reader with SetAPIKey ---")

	fredReader := reader.(*fred.FREDReader)

	// You can also create without API key and set it later
	// reader2 := fred.NewFREDReader(nil)
	// reader2.SetAPIKey(apiKey)

	data3, err := fredReader.ReadSingle(ctx, "FEDFUNDS", start, end)
	if err != nil {
		log.Fatalf("Failed to fetch FEDFUNDS: %v", err)
	}

	parsedData3 := data3.(*fred.ParsedData)
	fmt.Printf("✓ FEDFUNDS (Federal Funds Rate): %d observations\n", len(parsedData3.Dates))

	// Extract and display data
	fmt.Println("\n--- Data Extraction ---")

	dates := parsedData3.GetColumn("Date")
	values := parsedData3.GetColumn("Value")

	fmt.Println("Federal Funds Rate (first 5 observations):")
	for i := 0; i < min(5, len(dates)); i++ {
		fmt.Printf("  %s: %s%%\n", dates[i], values[i])
	}

	// Popular FRED series IDs
	fmt.Println("\n--- Popular FRED Series IDs ---")
	fmt.Println("Economic Indicators:")
	fmt.Println("  GDP       - Gross Domestic Product")
	fmt.Println("  UNRATE    - Unemployment Rate")
	fmt.Println("  CPIAUCSL  - Consumer Price Index")
	fmt.Println("  FEDFUNDS  - Federal Funds Effective Rate")
	fmt.Println("  DGS10     - 10-Year Treasury Constant Maturity Rate")
	fmt.Println("  DGS2      - 2-Year Treasury Constant Maturity Rate")
	fmt.Println("  DEXUSEU   - U.S. / Euro Foreign Exchange Rate")
	fmt.Println("  DTWEXBGS  - Trade Weighted U.S. Dollar Index")
	fmt.Println("\nSearch for more at: https://fred.stlouisfed.org/")

	fmt.Println("\n✓ FRED example completed successfully!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
