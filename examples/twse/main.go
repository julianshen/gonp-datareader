// Package main demonstrates Taiwan Stock Exchange (TWSE) data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/twse"
)

func main() {
	fmt.Println("gonp-datareader - TWSE (Taiwan Stock Exchange) Example")
	fmt.Println("=======================================================")

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch Taiwan Stock Data (Convenience Function) ---")

	start := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/twse",
		CacheTTL:   24 * time.Hour,
	}

	// Fetch TSMC (2330) stock data
	result, err := datareader.Read(ctx, "2330", "twse", start, end, opts)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	data := result.(*twse.ParsedData)
	fmt.Printf("✓ Fetched stock data for %s (%d observations)\n", data.Symbol, len(data.Date))

	if len(data.Date) > 0 {
		fmt.Println("\nRecent trading data:")
		// Show last 5 trading days
		startIdx := 0
		if len(data.Date) > 5 {
			startIdx = len(data.Date) - 5
		}
		for i := startIdx; i < len(data.Date); i++ {
			fmt.Printf("  %s: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%d\n",
				data.Date[i].Format("2006-01-02"),
				data.Open[i], data.High[i], data.Low[i], data.Close[i], data.Volume[i])
		}
	}

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("twse", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Popular Taiwan stocks
	fmt.Println("\n--- Example 3: Popular Taiwan Stocks ---")

	symbols := map[string]string{
		"2330": "TSMC (Taiwan Semiconductor)",
		"2317": "Hon Hai Precision (Foxconn)",
		"2454": "MediaTek",
		"2412": "Chunghwa Telecom",
		"0050": "Yuanta Taiwan 50 ETF",
	}

	for symbol, name := range symbols {
		fmt.Printf("\nFetching %s...\n", name)

		result, err := reader.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			fmt.Printf("✗ Error fetching %s: %v\n", symbol, err)
			continue
		}

		stockData := result.(*twse.ParsedData)
		fmt.Printf("✓ Fetched %d days of data\n", len(stockData.Date))

		if len(stockData.Date) > 0 {
			lastIdx := len(stockData.Date) - 1
			fmt.Printf("  Latest (%s): Close=%.2f, Change=%.2f\n",
				stockData.Date[lastIdx].Format("2006-01-02"),
				stockData.Close[lastIdx],
				stockData.Change[lastIdx])
		}
	}

	// Example 4: Data structure
	if len(data.Date) > 0 {
		fmt.Println("\n--- Example 4: Available Data Fields ---")
		lastIdx := len(data.Date) - 1
		fmt.Println("ParsedData structure contains:")
		fmt.Printf("  Symbol: %s\n", data.Symbol)
		fmt.Printf("  Date: %v\n", data.Date[lastIdx])
		fmt.Printf("  Open: %.2f\n", data.Open[lastIdx])
		fmt.Printf("  High: %.2f\n", data.High[lastIdx])
		fmt.Printf("  Low: %.2f\n", data.Low[lastIdx])
		fmt.Printf("  Close: %.2f\n", data.Close[lastIdx])
		fmt.Printf("  Volume: %d\n", data.Volume[lastIdx])
		fmt.Printf("  Transactions: %d\n", data.Transactions[lastIdx])
		fmt.Printf("  Change: %.2f\n", data.Change[lastIdx])
	}

	// Example 5: Symbol validation
	fmt.Println("\n--- Example 5: Symbol Validation ---")

	testSymbols := []string{"2330", "0050", "123456", "ABC", "233"}
	for _, sym := range testSymbols {
		err := reader.ValidateSymbol(sym)
		if err != nil {
			fmt.Printf("  %s: ✗ Invalid (%v)\n", sym, err)
		} else {
			fmt.Printf("  %s: ✓ Valid\n", sym)
		}
	}

	// TWSE information
	fmt.Println("\n--- TWSE Information ---")
	fmt.Println("Features:")
	fmt.Println("  • Free, no API key required")
	fmt.Println("  • Official Taiwan Stock Exchange data")
	fmt.Println("  • Daily trading data (OHLCV)")
	fmt.Println("  • ROC (Republic of China) calendar support")
	fmt.Println("  • Transaction counts and price changes")

	fmt.Println("\nSymbol Format:")
	fmt.Println("  Regular stocks: 4-digit codes (e.g., 2330, 2317, 2454)")
	fmt.Println("  ETFs: 4-digit codes starting with 00 (e.g., 0050, 0056)")
	fmt.Println("  Warrants: 6-digit codes")

	fmt.Println("\nPopular Symbols:")
	fmt.Println("  2330 - TSMC (Taiwan Semiconductor Manufacturing)")
	fmt.Println("  2317 - Hon Hai Precision Industry (Foxconn)")
	fmt.Println("  2454 - MediaTek")
	fmt.Println("  2412 - Chunghwa Telecom")
	fmt.Println("  2891 - CTBC Financial Holding")
	fmt.Println("  0050 - Yuanta/P-shares Taiwan Top 50 ETF")
	fmt.Println("  0056 - Yuanta FTSE Taiwan Dividend Plus ETF")

	fmt.Println("\nROC Calendar:")
	fmt.Println("  TWSE uses Taiwan's ROC (Republic of China) calendar")
	fmt.Println("  ROC Year = Gregorian Year - 1911")
	fmt.Println("  Example: 2025 = ROC Year 114 (2025 - 1911)")
	fmt.Println("  The library handles conversion automatically")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • No API key required (official public API)")
	fmt.Println("  • Comprehensive Taiwan market coverage")
	fmt.Println("  • High data quality from official source")
	fmt.Println("  • Includes transaction counts")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • Taiwan stock market analysis")
	fmt.Println("  • TSMC and semiconductor sector research")
	fmt.Println("  • Asian market diversification")
	fmt.Println("  • ETF tracking and analysis")

	fmt.Println("\n✓ TWSE example completed successfully!")
}
