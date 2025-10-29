// Package main demonstrates Tiingo data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/tiingo"
)

func main() {
	fmt.Println("gonp-datareader - Tiingo Example")
	fmt.Println("=================================")

	// Tiingo requires an API token
	apiKey := os.Getenv("TIINGO_API_KEY")
	if apiKey == "" {
		fmt.Println("\n⚠️  TIINGO_API_KEY environment variable not set")
		fmt.Println("\nTo use Tiingo:")
		fmt.Println("  1. Sign up at https://www.tiingo.com")
		fmt.Println("  2. Get your API token from https://www.tiingo.com/account/api/token")
		fmt.Println("  3. export TIINGO_API_KEY=your_token_here")
		fmt.Println("  4. Run this example again")
		fmt.Println("\nNote: Free tier allows unlimited API calls with reasonable rate limits")
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
		CacheDir:   ".cache/tiingo",
		CacheTTL:   24 * time.Hour, // Cache for 24 hours
	}

	fmt.Println("\nFetching AAPL data from Tiingo...")
	result, err := datareader.Read(ctx, "AAPL", "tiingo", start, end, opts)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	data := result.(*tiingo.ParsedData)
	fmt.Printf("✓ Fetched stock data for AAPL (%d observations)\n", len(data.Dates))

	if len(data.Prices) > 0 {
		fmt.Println("\nRecent closing prices:")
		// Show last 5 trading days
		startIdx := 0
		if len(data.Prices) > 5 {
			startIdx = len(data.Prices) - 5
		}
		for i := startIdx; i < len(data.Prices); i++ {
			price := data.Prices[i]
			fmt.Printf("  %s: Open=$%.2f, High=$%.2f, Low=$%.2f, Close=$%.2f, Volume=%d\n",
				data.Dates[i], price.Open, price.High, price.Low, price.Close, price.Volume)
		}
	}

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("tiingo", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Different stock symbols
	fmt.Println("\n--- Example 3: Multiple Stock Symbols ---")

	symbols := map[string]string{
		"MSFT":  "Microsoft Corporation",
		"GOOGL": "Alphabet Inc",
		"TSLA":  "Tesla Inc",
	}

	for symbol, name := range symbols {
		fmt.Printf("\nFetching %s (%s)...\n", name, symbol)

		result, err := reader.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			fmt.Printf("✗ Error fetching %s: %v\n", symbol, err)
			continue
		}

		stockData := result.(*tiingo.ParsedData)
		fmt.Printf("✓ Fetched %d days of data\n", len(stockData.Dates))

		if len(stockData.Prices) > 0 {
			lastIdx := len(stockData.Prices) - 1
			lastPrice := stockData.Prices[lastIdx]
			fmt.Printf("  Latest (%s): Close=$%.2f, Volume=%d\n",
				stockData.Dates[lastIdx], lastPrice.Close, lastPrice.Volume)
		}
	}

	// Example 4: Data columns
	if len(data.Prices) > 0 {
		fmt.Println("\n--- Example 4: Available Data Columns ---")
		fmt.Println("Columns: Date, Open, High, Low, Close, Volume")
		fmt.Println("\nTiingo provides daily OHLCV data with:")
		fmt.Println("  • Date: Trading date (YYYY-MM-DD)")
		fmt.Println("  • Open: Opening price")
		fmt.Println("  • High: Highest price")
		fmt.Println("  • Low: Lowest price")
		fmt.Println("  • Close: Closing price")
		fmt.Println("  • Volume: Trading volume")

		fmt.Println("\nSample row (most recent):")
		lastIdx := len(data.Prices) - 1
		price := data.Prices[lastIdx]
		fmt.Printf("  Date: %s\n", data.Dates[lastIdx])
		fmt.Printf("  Open: $%.2f\n", price.Open)
		fmt.Printf("  High: $%.2f\n", price.High)
		fmt.Printf("  Low: $%.2f\n", price.Low)
		fmt.Printf("  Close: $%.2f\n", price.Close)
		fmt.Printf("  Volume: %d\n", price.Volume)
	}

	// Tiingo information
	fmt.Println("\n--- Tiingo Information ---")
	fmt.Println("Features:")
	fmt.Println("  • High-quality stock market data and fundamentals")
	fmt.Println("  • Real-time and historical EOD (End-of-Day) prices")
	fmt.Println("  • Corporate actions and dividend data")
	fmt.Println("  • Comprehensive US and international stock coverage")
	fmt.Println("  • News and fundamental data available")

	fmt.Println("\nAuthentication:")
	fmt.Println("  • API token required (free tier available)")
	fmt.Println("  • Sign up: https://www.tiingo.com")
	fmt.Println("  • Get token: https://www.tiingo.com/account/api/token")
	fmt.Println("  • Use environment variable: TIINGO_API_KEY")

	fmt.Println("\nData Quality:")
	fmt.Println("  • Institutional-grade data with quality checks")
	fmt.Println("  • Corporate action adjustments available")
	fmt.Println("  • Point-in-time data for backtesting")
	fmt.Println("  • Historical data back to 1962 for many symbols")

	fmt.Println("\nAPI Rate Limits:")
	fmt.Println("  • Free tier: Unlimited requests with reasonable rate limits")
	fmt.Println("  • Premium tiers: Higher rate limits and additional features")
	fmt.Println("  • Typically 50-100 requests per hour for free tier")
	fmt.Println("  • Use caching to minimize API calls")

	fmt.Println("\nSymbol Format:")
	fmt.Println("  • US Stocks: Standard ticker (e.g., AAPL, MSFT, GOOGL)")
	fmt.Println("  • International: May require exchange suffix")

	fmt.Println("\nPricing Tiers:")
	fmt.Println("  • Free: EOD data with rate limits (this example)")
	fmt.Println("  • Starter: Higher limits, more endpoints")
	fmt.Println("  • Power: Real-time data and fundamentals")
	fmt.Println("  • Commercial: Redistribution rights")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • Quantitative research and backtesting")
	fmt.Println("  • Portfolio analysis and tracking")
	fmt.Println("  • Financial modeling and analysis")
	fmt.Println("  • Long-term historical data analysis")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • Generous free tier with unlimited requests")
	fmt.Println("  • High data quality with corporate action adjustments")
	fmt.Println("  • Extensive historical data coverage")
	fmt.Println("  • Simple, well-documented API")
	fmt.Println("  • Point-in-time data for accurate backtesting")

	fmt.Println("\nBest Practices:")
	fmt.Println("  1. Use caching to reduce API calls")
	fmt.Println("  2. Implement error handling for rate limits")
	fmt.Println("  3. Store historical data locally for frequent access")
	fmt.Println("  4. Consider premium tier for production applications")

	fmt.Println("\nLinks:")
	fmt.Println("  • Website: https://www.tiingo.com")
	fmt.Println("  • Documentation: https://api.tiingo.com/documentation/general/overview")
	fmt.Println("  • Get Token: https://www.tiingo.com/account/api/token")
	fmt.Println("  • Pricing: https://www.tiingo.com/pricing")

	fmt.Println("\n✓ Tiingo example completed successfully!")
}
