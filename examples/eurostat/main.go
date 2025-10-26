// Package main demonstrates Eurostat data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/eurostat"
)

func main() {
	fmt.Println("gonp-datareader - Eurostat Example")
	fmt.Println("===================================")

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch European Statistics (Convenience Function) ---")

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	// These would be used in a real example with actual dataset codes
	_ = ctx
	_ = start
	_ = end

	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/eurostat",
		CacheTTL:   24 * time.Hour,
	}

	// Note: Eurostat datasets have specific codes
	// This example uses mock data for demonstration
	fmt.Println("\nNote: Eurostat API requires specific dataset codes.")
	fmt.Println("Examples shown use simplified formats for demonstration.")
	fmt.Println("\nTo find real dataset codes:")
	fmt.Println("1. Visit https://ec.europa.eu/eurostat/data/database")
	fmt.Println("2. Browse to find your indicator")
	fmt.Println("3. Check the dataset code (e.g., DEMO_R_D3DENS)")

	// Example 2: Using the factory pattern
	fmt.Println("\n--- Example 2: Factory Pattern ---")

	reader, err := datareader.DataReader("eurostat", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Understanding Eurostat Data Structure
	fmt.Println("\n--- Example 3: Eurostat Data Structure ---")

	fmt.Println("\nEurostat uses JSON-stat format:")
	fmt.Println("  • Multi-dimensional statistical data format")
	fmt.Println("  • Dimensions organized in a cube model")
	fmt.Println("  • Automatic aggregation across dimensions")

	fmt.Println("\nCommon Dataset Codes:")
	fmt.Println("  • DEMO_R_D3DENS: Population density by NUTS 3 region")
	fmt.Println("  • GDP: Gross Domestic Product")
	fmt.Println("  • UNEMPLOYMENT: Unemployment rates")
	fmt.Println("  • TRADE: International trade statistics")
	fmt.Println("  • ENERGY: Energy statistics")

	// Example 4: Data columns
	fmt.Println("\n--- Example 4: Available Data Columns ---")
	fmt.Println("Columns: Date, Value")
	fmt.Println("\nEurostat provides time series data with:")
	fmt.Println("  • Date: Time period (various formats: YYYY, YYYY-MM, etc.)")
	fmt.Println("  • Value: Aggregated indicator value")
	fmt.Println("  • Note: Values are averaged across geographic dimensions")

	// Example 5: Mock data structure
	fmt.Println("\n--- Example 5: Data Structure Example ---")

	// Create sample data to show structure
	sampleData := &eurostat.ParsedData{
		Dates:  []string{"2020", "2021", "2022", "2023"},
		Values: []float64{100.0, 102.5, 105.2, 107.8},
	}

	fmt.Println("\nSample European Indicator Data:")
	for i := range sampleData.Dates {
		fmt.Printf("  %s: %.1f\n", sampleData.Dates[i], sampleData.Values[i])
	}

	// Eurostat information
	fmt.Println("\n--- Eurostat Information ---")
	fmt.Println("Features:")
	fmt.Println("  • Official European Union statistics")
	fmt.Println("  • Data from EU member states and partners")
	fmt.Println("  • JSON-stat format (lightweight, visualization-ready)")
	fmt.Println("  • No API key required (free access)")
	fmt.Println("  • Updated twice daily (11:00 and 23:00 CET)")

	fmt.Println("\nData Coverage:")
	fmt.Println("  • Economy and Finance: GDP, inflation, employment")
	fmt.Println("  • Population and Social: Demographics, education, health")
	fmt.Println("  • Industry and Services: Production, trade, ICT")
	fmt.Println("  • Agriculture and Environment: Farming, energy, emissions")
	fmt.Println("  • Regions: NUTS (Nomenclature of Territorial Units)")

	fmt.Println("\nAPI Access:")
	fmt.Println("  • No authentication required")
	fmt.Println("  • Free for all users")
	fmt.Println("  • Statistics API in JSON-stat format")
	fmt.Println("  • Rate limiting applies to protect service")

	fmt.Println("\nDataset Code Format:")
	fmt.Println("  • Codes are uppercase with underscores")
	fmt.Println("  • Example: DEMO_R_D3DENS (Population density)")
	fmt.Println("  • Example: GDP (Gross Domestic Product)")
	fmt.Println("  • Find codes in the Eurostat Data Browser")

	fmt.Println("\nLanguages Supported:")
	fmt.Println("  • EN: English (default)")
	fmt.Println("  • FR: French")
	fmt.Println("  • DE: German")

	fmt.Println("\nBest Use Cases:")
	fmt.Println("  • European economic comparisons")
	fmt.Println("  • Regional analysis within EU")
	fmt.Println("  • Policy research and benchmarking")
	fmt.Println("  • Time series analysis of EU indicators")

	fmt.Println("\nAdvantages:")
	fmt.Println("  • Official EU statistics")
	fmt.Println("  • Harmonized across member states")
	fmt.Println("  • Comprehensive coverage of EU topics")
	fmt.Println("  • Free and open access")
	fmt.Println("  • JSON-stat format for easy visualization")

	fmt.Println("\nLimitations:")
	fmt.Println("  • Focused on European countries")
	fmt.Println("  • Dataset codes require lookup")
	fmt.Println("  • No date filtering in API (client-side filtering)")
	fmt.Println("  • Multi-dimensional data is aggregated")

	fmt.Println("\nHow to Find Dataset Codes:")
	fmt.Println("  1. Go to https://ec.europa.eu/eurostat/data/database")
	fmt.Println("  2. Browse the database tree to find your topic")
	fmt.Println("  3. Click on a dataset to view details")
	fmt.Println("  4. The dataset code is shown in the page")
	fmt.Println("  5. Use that code in your API calls")

	fmt.Println("\nJSON-stat Format:")
	fmt.Println("  • Cube model: data organized in dimensions")
	fmt.Println("  • Dimensions: geo (geography), time, indicators, etc.")
	fmt.Println("  • Values: flat array in row-major order")
	fmt.Println("  • Automatic dimension aggregation")

	fmt.Println("\nLinks:")
	fmt.Println("  • Data Browser: https://ec.europa.eu/eurostat/data/database")
	fmt.Println("  • API Documentation: https://ec.europa.eu/eurostat/web/user-guides/data-browser/api-data-access")
	fmt.Println("  • Statistics Explained: https://ec.europa.eu/eurostat/statistics-explained")

	fmt.Println("\nBest Practices:")
	fmt.Println("  1. Use caching to reduce API calls")
	fmt.Println("  2. Verify dataset codes in Data Browser first")
	fmt.Println("  3. Be aware datasets update twice daily")
	fmt.Println("  4. Understand dimension aggregation")
	fmt.Println("  5. Check data frequency (annual, monthly, etc.)")

	fmt.Println("\n✓ Eurostat example completed!")
	fmt.Println("\nNote: This example uses mock data. For real data:")
	fmt.Println("1. Find your dataset code at https://ec.europa.eu/eurostat/data/database")
	fmt.Println("2. Use the code in reader.ReadSingle(ctx, \"DATASET_CODE\", start, end)")
	fmt.Println("3. Parse the returned data structure")
}
