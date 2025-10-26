// Package main demonstrates IEX Cloud data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/iex"
)

func main() {
	fmt.Println("gonp-datareader - IEX Cloud Example")
	fmt.Println("====================================")

	// IEX Cloud requires an API token
	apiKey := os.Getenv("IEX_API_KEY")
	if apiKey == "" {
		fmt.Println("\n⚠️  IEX_API_KEY environment variable not set")
		fmt.Println("\nTo use IEX Cloud:")
		fmt.Println("  1. Sign up at https://iexcloud.io")
		fmt.Println("  2. Get your API token (publishable key)")
		fmt.Println("  3. export IEX_API_KEY=your_token_here")
		fmt.Println("  4. Run this example again")
		fmt.Println("\nContinuing with example structure (will fail without valid token)...")
		apiKey = "demo_token"
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
		CacheDir:   ".cache/iex",
		CacheTTL:   24 * time.Hour,
	}

	fmt.Println("\nFetching AAPL data from IEX Cloud...")
	result, err := datareader.Read(ctx, "AAPL", "iex", start, end, opts)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		fmt.Println("\n⚠️  This is expected if you don't have a valid API token")
	} else {
		data := result.(*iex.ParsedData)
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
				fmt.Printf("  %s: Open=%s, High=%s, Low=%s, Close=%s, Volume=%s\n",
					row["Date"], row["Open"], row["High"], row["Low"], row["Close"], row["Volume"])
			}
		}
	}

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("iex", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Multiple symbols
	fmt.Println("\n--- Example 3: Multiple Stock Symbols ---")

	symbols := map[string]string{
		"AAPL": "Apple Inc",
		"MSFT": "Microsoft Corporation",
		"GOOGL": "Alphabet Inc",
		"TSLA": "Tesla Inc",
	}

	for symbol, name := range symbols {
		fmt.Printf("\nFetching %s (%s)...\n", name, symbol)

		result, err := reader.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			fmt.Printf("✗ Error fetching %s: %v\n", symbol, err)
			if apiKey == "demo_token" {
				fmt.Println("  (This is expected without a valid API token)")
			}
			continue
		}

		stockData := result.(*iex.ParsedData)
		fmt.Printf("✓ Fetched %d days of data\n", len(stockData.Rows))

		if len(stockData.Rows) > 0 {
			lastRow := stockData.Rows[len(stockData.Rows)-1]
			fmt.Printf("  Latest (%s): Open=%s, High=%s, Low=%s, Close=%s\n",
				lastRow["Date"], lastRow["Open"], lastRow["High"], lastRow["Low"], lastRow["Close"])
		}
	}

	// Example 4: Data columns
	fmt.Println("\n--- Example 4: Available Data Columns ---")
	fmt.Println("Columns: Date, Open, High, Low, Close, Volume")
	fmt.Println("\nIEX Cloud provides daily OHLCV data with:")
	fmt.Println("  • Date: Trading date (YYYY-MM-DD)")
	fmt.Println("  • Open: Opening price")
	fmt.Println("  • High: Highest price")
	fmt.Println("  • Low: Lowest price")
	fmt.Println("  • Close: Closing price")
	fmt.Println("  • Volume: Trading volume")

	// IEX Cloud information
	fmt.Println("\n--- IEX Cloud Information ---")
	fmt.Println("Features:")
	fmt.Println("  • Professional-grade stock market data")
	fmt.Println("  • Real-time and historical data")
	fmt.Println("  • High data quality and reliability")
	fmt.Println("  • Comprehensive US stock coverage")

	fmt.Println("\nAuthentication:")
	fmt.Println("  • API token required (sign up at https://iexcloud.io)")
	fmt.Println("  • Free tier available with usage limits")
	fmt.Println("  • Use environment variable: IEX_API_KEY")

	fmt.Println("\nDate Ranges:")
	fmt.Println("  The reader automatically selects the appropriate range:")
	fmt.Println("  • < 45 days: 1m (1 month)")
	fmt.Println("  • < 135 days: 3m (3 months)")
	fmt.Println("  • < 270 days: 6m (6 months)")
	fmt.Println("  • < 548 days: 1y (1 year)")
	fmt.Println("  • < 1095 days: 2y (2 years)")
	fmt.Println("  • ≥ 1095 days: 5y (5 years, maximum)")

	fmt.Println("\nPricing Tiers:")
	fmt.Println("  • Free: Limited API calls per month")
	fmt.Println("  • Launch: Basic usage for individual developers")
	fmt.Println("  • Grow: For growing applications")
	fmt.Println("  • Enterprise: Unlimited usage")

	fmt.Println("\nSymbol Format:")
	fmt.Println("  • US Stocks: Standard ticker (e.g., AAPL, MSFT, GOOGL)")
	fmt.Println("  • Case-insensitive")

	fmt.Println("\nAPI Rate Limits:")
	fmt.Println("  • Free tier: 50,000 messages/month")
	fmt.Println("  • Each API call consumes messages based on data returned")
	fmt.Println("  • Monitor your usage at https://iexcloud.io/console")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • Production applications requiring reliable data")
	fmt.Println("  • Financial analysis and research")
	fmt.Println("  • Algorithmic trading development")
	fmt.Println("  • Portfolio tracking applications")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • High-quality, exchange-grade data")
	fmt.Println("  • Well-documented API")
	fmt.Println("  • Reliable uptime and performance")
	fmt.Println("  • Comprehensive data beyond just prices")

	fmt.Println("\nLinks:")
	fmt.Println("  • Website: https://iexcloud.io")
	fmt.Println("  • Documentation: https://iexcloud.io/docs/api/")
	fmt.Println("  • Console: https://iexcloud.io/console")
	fmt.Println("  • Pricing: https://iexcloud.io/pricing/")

	fmt.Println("\n✓ IEX Cloud example completed!")
}
