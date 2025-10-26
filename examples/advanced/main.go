// Package main demonstrates advanced usage with custom options and error handling.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
	fmt.Println("gonp-datareader - Advanced Example")
	fmt.Println("===================================")

	ctx := context.Background()

	// Custom options with longer timeout and more retries
	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 5,
		RetryDelay: 2 * time.Second,
		UserAgent:  "gonp-datareader-example/1.0",
	}

	fmt.Println("\nConfiguration:")
	fmt.Printf("  Timeout: %v\n", opts.Timeout)
	fmt.Printf("  Max Retries: %d\n", opts.MaxRetries)
	fmt.Printf("  Retry Delay: %v\n", opts.RetryDelay)

	// Create reader with custom options
	reader, err := datareader.DataReader("yahoo", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("\n✓ Created %s reader\n", reader.Name())

	// Demonstrate validation
	fmt.Println("\n--- Input Validation ---")

	validSymbols := []string{"AAPL", "MSFT", "GOOGL"}
	for _, symbol := range validSymbols {
		if err := reader.ValidateSymbol(symbol); err == nil {
			fmt.Printf("✓ %s is valid\n", symbol)
		}
	}

	invalidSymbols := []string{"", "AA PL", "AAPL@"}
	for _, symbol := range invalidSymbols {
		if err := reader.ValidateSymbol(symbol); err != nil {
			fmt.Printf("✗ %q is invalid: %v\n", symbol, err)
		}
	}

	// Fetch data with error handling
	fmt.Println("\n--- Fetching Data with Error Handling ---")

	end := time.Now()
	start := end.AddDate(0, 0, -7) // Last week

	result, err := reader.ReadSingle(ctx, "TSLA", start, end)
	if err != nil {
		// Handle different error types
		fmt.Printf("✗ Error fetching data: %v\n", err)

		// You could check for specific error types here
		// if errors.Is(err, datareader.ErrNetworkError) { ... }

		log.Fatalf("Failed to fetch data")
	}

	data := result.(*yahoo.ParsedData)
	fmt.Printf("✓ Successfully fetched %d days of TSLA data\n", len(data.Rows))

	// Analyze the data
	fmt.Println("\n--- Data Analysis ---")

	if len(data.Rows) > 0 {
		// Get all closing prices
		closes := data.GetColumn("Close")
		volumes := data.GetColumn("Volume")

		fmt.Printf("Data points: %d\n", len(closes))
		fmt.Printf("First close: %s\n", closes[0])
		fmt.Printf("Last close: %s\n", closes[len(closes)-1])

		if len(volumes) > 0 {
			fmt.Printf("First volume: %s\n", volumes[0])
		}

		// Show all available columns
		fmt.Printf("\nAvailable columns: %v\n", data.Columns)

		// Show first row
		fmt.Println("\nFirst row of data:")
		for _, col := range data.Columns {
			fmt.Printf("  %s: %s\n", col, data.Rows[0][col])
		}
	}

	// List all available sources
	fmt.Println("\n--- Available Data Sources ---")
	sources := datareader.ListSources()
	for _, source := range sources {
		fmt.Printf("  • %s\n", source)
	}

	fmt.Println("\n✓ Advanced example completed successfully!")
}
