// Package main demonstrates Stooq data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/stooq"
)

func main() {
	fmt.Println("gonp-datareader - Stooq Example")
	fmt.Println("================================")

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch Stock Data (Convenience Function) ---")

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/stooq",
		CacheTTL:   24 * time.Hour,
	}

	// Stooq uses .US suffix for US stocks
	result, err := datareader.Read(ctx, "AAPL.US", "stooq", start, end, opts)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	data := result.(*stooq.ParsedData)
	fmt.Printf("✓ Fetched stock data for AAPL.US (%d observations)\n", len(data.Rows))

	if len(data.Rows) > 0 {
		fmt.Println("\nRecent closing prices:")
		// Show last 5 trading days
		startIdx := 0
		if len(data.Rows) > 5 {
			startIdx = len(data.Rows) - 5
		}
		for i := startIdx; i < len(data.Rows); i++ {
			row := data.Rows[i]
			fmt.Printf("  %s: Close=%s, Volume=%s\n",
				row["Date"], row["Close"], row["Volume"])
		}
	}

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("stooq", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: International markets
	fmt.Println("\n--- Example 3: International Markets ---")

	// Stooq supports international markets
	symbols := map[string]string{
		"MSFT.US": "Microsoft (US)",
		"^SPX":    "S&P 500 Index",
		"^DJI":    "Dow Jones Industrial Average",
		"AAPL.US": "Apple Inc",
	}

	for symbol, name := range symbols {
		fmt.Printf("\nFetching %s...\n", name)

		result, err := reader.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			fmt.Printf("✗ Error fetching %s: %v\n", symbol, err)
			continue
		}

		stockData := result.(*stooq.ParsedData)
		fmt.Printf("✓ Fetched %d days of data\n", len(stockData.Rows))

		if len(stockData.Rows) > 0 {
			lastRow := stockData.Rows[len(stockData.Rows)-1]
			fmt.Printf("  Latest (%s): Close=%s\n",
				lastRow["Date"], lastRow["Close"])
		}
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

	// Stooq information
	fmt.Println("\n--- Stooq Information ---")
	fmt.Println("Features:")
	fmt.Println("  • Free, no API key required")
	fmt.Println("  • Historical daily OHLCV data")
	fmt.Println("  • International markets")
	fmt.Println("  • Indices and commodities")

	fmt.Println("\nSymbol Format:")
	fmt.Println("  US Stocks: {SYMBOL}.US (e.g., AAPL.US, MSFT.US)")
	fmt.Println("  Indices: ^{INDEX} (e.g., ^SPX, ^DJI)")
	fmt.Println("  Other markets: Various suffixes (.UK, .DE, .JP, etc.)")

	fmt.Println("\nPopular Symbols:")
	fmt.Println("  AAPL.US - Apple Inc")
	fmt.Println("  MSFT.US - Microsoft Corporation")
	fmt.Println("  GOOGL.US - Alphabet Inc")
	fmt.Println("  ^SPX - S&P 500 Index")
	fmt.Println("  ^DJI - Dow Jones Industrial Average")
	fmt.Println("  ^IXIC - NASDAQ Composite")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • No API key required (free access)")
	fmt.Println("  • Good coverage of international markets")
	fmt.Println("  • Simple CSV format")
	fmt.Println("  • Reliable historical data")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • Quick prototyping without API setup")
	fmt.Println("  • International market data")
	fmt.Println("  • Historical analysis")
	fmt.Println("  • Educational purposes")

	fmt.Println("\n✓ Stooq example completed successfully!")
}
