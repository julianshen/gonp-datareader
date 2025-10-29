// Package main demonstrates basic usage of gonp-datareader with Yahoo Finance.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
	fmt.Println("gonp-datareader - Basic Example")
	fmt.Println("================================")

	ctx := context.Background()

	// Define date range - last month
	end := time.Now()
	start := end.AddDate(0, -1, 0)

	fmt.Printf("\nFetching data from %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	// Method 1: Using the convenience function
	fmt.Println("\n--- Method 1: Convenience Function ---")
	data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	parsedData := data.(*yahoo.ParsedData)
	fmt.Printf("✓ Fetched %d days of data for AAPL\n", len(parsedData.Rows))

	if len(parsedData.Rows) > 0 {
		firstRow := parsedData.Rows[0]
		fmt.Printf("  First day: %s - Close: %s, Volume: %s\n",
			firstRow["Date"], firstRow["Close"], firstRow["Volume"])
	}

	// Method 2: Using the factory
	fmt.Println("\n--- Method 2: Using Factory ---")
	reader, err := datareader.DataReader("yahoo", nil)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	data2, err := reader.ReadSingle(ctx, "MSFT", start, end)
	if err != nil {
		log.Fatalf("Failed to fetch MSFT data: %v", err)
	}

	parsedData2 := data2.(*yahoo.ParsedData)
	fmt.Printf("✓ Fetched %d days of data for MSFT\n", len(parsedData2.Rows))

	if len(parsedData2.Rows) > 0 {
		lastRow := parsedData2.Rows[len(parsedData2.Rows)-1]
		fmt.Printf("  Last day: %s - Close: %s\n",
			lastRow["Date"], lastRow["Close"])
	}

	// Method 3: Fetch multiple symbols
	fmt.Println("\n--- Method 3: Multiple Symbols ---")
	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	results, err := reader.Read(ctx, symbols, start, end)
	if err != nil {
		log.Fatalf("Failed to fetch multiple symbols: %v", err)
	}

	dataMap := results.(map[string]*yahoo.ParsedData)
	for symbol, data := range dataMap {
		fmt.Printf("✓ %s: %d days of data\n", symbol, len(data.Rows))
	}

	// Extract closing prices
	fmt.Println("\n--- Extracting Specific Data ---")
	appleData := dataMap["AAPL"]
	closePrices := appleData.GetColumn("Close")
	fmt.Printf("AAPL closing prices (first 5): %v\n", closePrices[:min(5, len(closePrices))])

	fmt.Println("\n✓ Example completed successfully!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
