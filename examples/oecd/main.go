// Package main demonstrates OECD data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/oecd"
)

func main() {
	fmt.Println("gonp-datareader - OECD Example")
	fmt.Println("===============================")

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch Economic Data (Convenience Function) ---")

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	// These would be used in a real example with actual dataset IDs
	_ = ctx
	_ = start
	_ = end

	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/oecd",
		CacheTTL:   24 * time.Hour,
	}

	// Note: OECD datasets have specific formats like "DATASET/DIMENSIONS"
	// This example uses a simplified dataset identifier
	// In practice, you'd need to know the correct OECD dataset structure
	fmt.Println("\nNote: OECD API requires specific dataset identifiers.")
	fmt.Println("Examples shown use simplified formats for demonstration.")
	fmt.Println("\nTo find real dataset IDs:")
	fmt.Println("1. Visit https://data.oecd.org")
	fmt.Println("2. Select a dataset")
	fmt.Println("3. Click 'Developer API' to see the correct format")

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("oecd", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Understanding OECD Data Structure
	fmt.Println("\n--- Example 3: OECD Data Structure ---")

	fmt.Println("\nOECD uses SDMX-JSON format with complex dataset identifiers:")
	fmt.Println("  • Format: DATASET/DIMENSIONS")
	fmt.Println("  • Example: MEI/USA (Main Economic Indicators for USA)")
	fmt.Println("  • Example: QNA/AUS.GDP (Quarterly National Accounts GDP for Australia)")

	fmt.Println("\nCommon OECD Datasets:")
	fmt.Println("  • MEI: Main Economic Indicators")
	fmt.Println("  • QNA: Quarterly National Accounts")
	fmt.Println("  • EO: Economic Outlook")
	fmt.Println("  • KEI: Key Economic Indicators")
	fmt.Println("  • REGION_ECONOM: Regional Economy Statistics")

	// Example 4: Data columns
	fmt.Println("\n--- Example 4: Available Data Columns ---")
	fmt.Println("Columns: Date, Value")
	fmt.Println("\nOECD provides time series data with:")
	fmt.Println("  • Date: Time period (various formats: YYYY-MM, YYYY-QN, YYYY)")
	fmt.Println("  • Value: Indicator value (units vary by dataset)")

	// Example 5: Mock data structure
	fmt.Println("\n--- Example 5: Data Structure Example ---")

	// Create sample data to show structure
	sampleData := &oecd.ParsedData{
		Dates:  []string{"2020-Q1", "2020-Q2", "2020-Q3", "2020-Q4"},
		Values: []float64{100.0, 98.5, 102.3, 105.1},
	}

	fmt.Println("\nSample Economic Indicator Data:")
	for i := range sampleData.Dates {
		fmt.Printf("  %s: %.1f\n", sampleData.Dates[i], sampleData.Values[i])
	}

	// OECD information
	fmt.Println("\n--- OECD Information ---")
	fmt.Println("Features:")
	fmt.Println("  • Comprehensive economic and social statistics")
	fmt.Println("  • Data from OECD member countries and partners")
	fmt.Println("  • Standardized SDMX-JSON format")
	fmt.Println("  • No API key required (free access)")
	fmt.Println("  • High-quality official statistics")

	fmt.Println("\nData Coverage:")
	fmt.Println("  • Economic Indicators: GDP, unemployment, inflation, trade")
	fmt.Println("  • Social Statistics: Education, health, inequality")
	fmt.Println("  • Environmental Data: Energy, emissions, resources")
	fmt.Println("  • Government Statistics: Revenue, spending, debt")

	fmt.Println("\nAPI Access:")
	fmt.Println("  • No authentication required")
	fmt.Println("  • Free for all users")
	fmt.Println("  • Rate limiting applies to protect service")
	fmt.Println("  • SDMX-JSON standard format")

	fmt.Println("\nDataset Identifiers:")
	fmt.Println("  Format varies by dataset - check OECD Data Explorer")
	fmt.Println("  • Simple: DATASET_CODE")
	fmt.Println("  • Filtered: DATASET/COUNTRY.INDICATOR")
	fmt.Println("  • Complex: DATASET/COUNTRY.INDICATOR.MEASURE")

	fmt.Println("\nDate Ranges:")
	fmt.Println("  • Monthly: YYYY-MM (e.g., 2020-01)")
	fmt.Println("  • Quarterly: YYYY-QN (e.g., 2020-Q1)")
	fmt.Println("  • Annual: YYYY (e.g., 2020)")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • International economic comparisons")
	fmt.Println("  • Policy analysis and research")
	fmt.Println("  • Long-term trend analysis")
	fmt.Println("  • Cross-country benchmarking")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • Official government statistics")
	fmt.Println("  • Standardized across countries")
	fmt.Println("  • Extensive historical data")
	fmt.Println("  • Free and open access")
	fmt.Println("  • Well-documented datasets")

	fmt.Println("\nLimitations:")
	fmt.Println("  • Complex dataset identifiers require research")
	fmt.Println("  • SDMX format has nested structure")
	fmt.Println("  • Data release frequency varies by indicator")
	fmt.Println("  • Some datasets have limited country coverage")

	fmt.Println("\nHow to Find Dataset IDs:")
	fmt.Println("  1. Go to https://data.oecd.org")
	fmt.Println("  2. Browse or search for your indicator")
	fmt.Println("  3. Click on the dataset")
	fmt.Println("  4. Click 'Developer API' button")
	fmt.Println("  5. Copy the dataset ID from the API query")

	fmt.Println("\nLinks:")
	fmt.Println("  • Data Explorer: https://data.oecd.org")
	fmt.Println("  • API Documentation: https://data.oecd.org/api/")
	fmt.Println("  • SDMX Guide: https://data.oecd.org/api/sdmx-json-documentation/")

	fmt.Println("\nBest Practices:")
	fmt.Println("  1. Use caching to reduce API calls")
	fmt.Println("  2. Verify dataset IDs in OECD Data Explorer first")
	fmt.Println("  3. Check data frequency (monthly, quarterly, annual)")
	fmt.Println("  4. Be aware of data revision schedules")
	fmt.Println("  5. Understand dimension filters for your dataset")

	fmt.Println("\n✓ OECD example completed!")
	fmt.Println("\nNote: This example uses mock data. For real data:")
	fmt.Println("1. Find your dataset ID at https://data.oecd.org")
	fmt.Println("2. Use the ID in reader.ReadSingle(ctx, \"DATASET_ID\", start, end)")
	fmt.Println("3. Parse the returned data structure")
}
