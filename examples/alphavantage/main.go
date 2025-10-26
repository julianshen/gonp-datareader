// Package main demonstrates Alpha Vantage data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/alphavantage"
)

func main() {
	fmt.Println("gonp-datareader - Alpha Vantage Example")
	fmt.Println("========================================")

	// Alpha Vantage requires an API key
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		fmt.Println("\n⚠️  ALPHA_VANTAGE_API_KEY environment variable not set")
		fmt.Println("\nTo use Alpha Vantage:")
		fmt.Println("1. Get a free API key at: https://www.alphavantage.co/support/#api-key")
		fmt.Println("2. Set the environment variable:")
		fmt.Println("   export ALPHA_VANTAGE_API_KEY=your_api_key_here")
		fmt.Println("\nNote: Free tier allows 5 API calls per minute and 500 calls per day")
		os.Exit(1)
	}

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch Stock Data (Convenience Function) ---")

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	opts := &datareader.Options{
		APIKey:     apiKey,
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/alphavantage",
		CacheTTL:   24 * time.Hour, // Cache for 24 hours to avoid rate limits
	}

	result, err := datareader.Read(ctx, "AAPL", "alphavantage", start, end, opts)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	data := result.(*alphavantage.ParsedData)
	fmt.Printf("✓ Fetched stock data for AAPL (%d observations)\n", len(data.Rows))

	if len(data.Rows) > 0 {
		fmt.Println("\nRecent closing prices:")
		// Show last 5 trading days
		startIdx := 0
		if len(data.Rows) > 5 {
			startIdx = len(data.Rows) - 5
		}
		for i := startIdx; i < len(data.Rows); i++ {
			row := data.Rows[i]
			fmt.Printf("  %s: Close=$%s, Volume=%s\n",
				row["Date"], row["Close"], row["Volume"])
		}
	}

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern with Options ---")

	reader, err := datareader.DataReader("alphavantage", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Different stock symbols
	fmt.Println("\n--- Example 3: Multiple Stock Symbols ---")

	symbols := []string{"MSFT", "GOOGL", "TSLA"}
	for _, symbol := range symbols {
		fmt.Printf("\nFetching %s...\n", symbol)

		result, err := reader.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			fmt.Printf("✗ Error fetching %s: %v\n", symbol, err)

			// Check for rate limit
			if err.Error() == "parse response: rate limit exceeded" {
				fmt.Println("\n⚠️  Rate limit reached!")
				fmt.Println("Free tier: 5 calls/minute, 500 calls/day")
				fmt.Println("Consider using caching to reduce API calls")
				break
			}
			continue
		}

		stockData := result.(*alphavantage.ParsedData)
		fmt.Printf("✓ Fetched %d days of data\n", len(stockData.Rows))

		if len(stockData.Rows) > 0 {
			lastRow := stockData.Rows[len(stockData.Rows)-1]
			fmt.Printf("  Latest (%s): Close=$%s\n",
				lastRow["Date"], lastRow["Close"])
		}

		// Sleep to avoid rate limit (5 calls per minute = 12 seconds between calls)
		fmt.Println("  Waiting 12 seconds to avoid rate limit...")
		time.Sleep(12 * time.Second)
	}

	// Example 4: Data columns
	if len(data.Rows) > 0 {
		fmt.Println("\n--- Example 4: Available Data Columns ---")
		fmt.Printf("Columns: %v\n", data.Columns)
		fmt.Println("\nSample row (OHLCV data):")
		row := data.Rows[len(data.Rows)-1]
		for _, col := range data.Columns {
			fmt.Printf("  %s: %s\n", col, row[col])
		}
	}

	// Alpha Vantage information
	fmt.Println("\n--- Alpha Vantage Information ---")
	fmt.Println("API Limits:")
	fmt.Println("  Free tier: 5 API calls per minute, 500 calls per day")
	fmt.Println("  Premium tiers available with higher limits")

	fmt.Println("\nData Available:")
	fmt.Println("  TIME_SERIES_DAILY: Daily OHLCV data (this example)")
	fmt.Println("  Also available: Intraday, Weekly, Monthly, Adjusted, etc.")

	fmt.Println("\nAPI Key:")
	fmt.Println("  Get free key: https://www.alphavantage.co/support/#api-key")
	fmt.Println("  Set env var: export ALPHA_VANTAGE_API_KEY=your_key")

	fmt.Println("\nBest Practices:")
	fmt.Println("  1. Use caching to reduce API calls")
	fmt.Println("  2. Implement delays between requests (12 seconds for free tier)")
	fmt.Println("  3. Monitor rate limit errors")
	fmt.Println("  4. Consider premium tier for production use")

	fmt.Println("\n✓ Alpha Vantage example completed successfully!")
}
