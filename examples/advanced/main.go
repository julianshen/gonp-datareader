// Package main demonstrates advanced usage with custom options and error handling.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/fred"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

func main() {
	fmt.Println("gonp-datareader - Advanced Example")
	fmt.Println("===================================")

	ctx := context.Background()

	// Custom options with longer timeout, more retries, and rate limiting
	// Note: UserAgent should be browser-like for Yahoo Finance compatibility
	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 5,
		RetryDelay: 2 * time.Second,
		UserAgent:  "Mozilla/5.0 (compatible; gonp-datareader/1.0)",
		RateLimit:  2.0, // Limit to 2 requests per second
	}

	fmt.Println("\nConfiguration:")
	fmt.Printf("  Timeout: %v\n", opts.Timeout)
	fmt.Printf("  Max Retries: %d\n", opts.MaxRetries)
	fmt.Printf("  Retry Delay: %v\n", opts.RetryDelay)
	fmt.Printf("  User Agent: %s\n", opts.UserAgent)
	fmt.Printf("  Rate Limit: %.1f req/sec\n", opts.RateLimit)

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

	// Example with FRED (if API key is available)
	fmt.Println("\n--- FRED Data Source (Economic Data) ---")

	fredAPIKey := os.Getenv("FRED_API_KEY")
	if fredAPIKey != "" {
		fredOpts := &datareader.Options{
			APIKey:     fredAPIKey,
			Timeout:    60 * time.Second,
			MaxRetries: 3,
			RetryDelay: 2 * time.Second,
		}

		fredReader, err := datareader.DataReader("fred", fredOpts)
		if err != nil {
			fmt.Printf("✗ Failed to create FRED reader: %v\n", err)
		} else {
			fmt.Printf("✓ Created %s reader\n", fredReader.Name())

			// Fetch GDP data
			startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

			result, err := fredReader.ReadSingle(ctx, "GDP", startDate, endDate)
			if err != nil {
				fmt.Printf("✗ Error fetching GDP data: %v\n", err)
			} else {
				fredData := result.(*fred.ParsedData)
				fmt.Printf("✓ Successfully fetched GDP data: %d observations\n", len(fredData.Dates))

				if len(fredData.Dates) > 0 {
					fmt.Printf("  Latest GDP: %s = $%s billion\n",
						fredData.Dates[len(fredData.Dates)-1],
						fredData.Values[len(fredData.Values)-1])
				}
			}
		}
	} else {
		fmt.Println("⚠️  FRED_API_KEY not set - skipping FRED example")
		fmt.Println("  Get a free API key at: https://fred.stlouisfed.org/docs/api/api_key.html")
	}

	// List all available sources
	fmt.Println("\n--- Available Data Sources ---")
	sources := datareader.ListSources()
	for _, source := range sources {
		fmt.Printf("  • %s\n", source)
	}

	fmt.Println("\n✓ Advanced example completed successfully!")
}
