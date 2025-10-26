// Package main demonstrates World Bank data reader usage.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/gonp-datareader"
	"github.com/julianshen/gonp-datareader/sources/worldbank"
)

func main() {
	fmt.Println("gonp-datareader - World Bank Example")
	fmt.Println("====================================")

	ctx := context.Background()

	// Example 1: Using the convenience function
	fmt.Println("\n--- Example 1: Fetch GDP Data (Convenience Function) ---")

	start := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	// Symbol format: "country/indicator"
	// USA = United States
	// NY.GDP.MKTP.CD = GDP (current US$)
	result, err := datareader.Read(ctx, "USA/NY.GDP.MKTP.CD", "worldbank", start, end, nil)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	data := result.(*worldbank.ParsedData)
	fmt.Printf("✓ Fetched GDP data for USA (%d observations)\n", len(data.Dates))

	if len(data.Dates) > 0 {
		fmt.Println("\nRecent GDP values:")
		// Show last 5 years
		startIdx := 0
		if len(data.Dates) > 5 {
			startIdx = len(data.Dates) - 5
		}
		for i := startIdx; i < len(data.Dates); i++ {
			fmt.Printf("  %s: $%s\n", data.Dates[i], data.Values[i])
		}
	}

	// Example 2: Using the factory with options
	fmt.Println("\n--- Example 2: Factory with Options ---")

	opts := &datareader.Options{
		Timeout:    60 * time.Second,
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		CacheDir:   ".cache/worldbank",
		CacheTTL:   24 * time.Hour,
	}

	reader, err := datareader.DataReader("worldbank", opts)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	fmt.Printf("✓ Created %s reader\n", reader.Name())

	// Example 3: Population data
	fmt.Println("\n--- Example 3: Population Data ---")

	// SP.POP.TOTL = Population, total
	popResult, err := reader.ReadSingle(ctx, "CHN/SP.POP.TOTL", start, end)
	if err != nil {
		log.Fatalf("Failed to fetch population data: %v", err)
	}

	popData := popResult.(*worldbank.ParsedData)
	fmt.Printf("✓ Fetched population data for China (%d observations)\n", len(popData.Dates))

	if len(popData.Dates) > 0 {
		fmt.Println("\nRecent population:")
		startIdx := 0
		if len(popData.Dates) > 5 {
			startIdx = len(popData.Dates) - 5
		}
		for i := startIdx; i < len(popData.Dates); i++ {
			fmt.Printf("  %s: %s people\n", popData.Dates[i], popData.Values[i])
		}
	}

	// Example 4: Multiple countries
	fmt.Println("\n--- Example 4: Multiple Countries (Inflation Rate) ---")

	// FP.CPI.TOTL.ZG = Inflation, consumer prices (annual %)
	countries := []string{"USA;CHN;GBR"}
	for _, country := range countries {
		inflationResult, err := reader.ReadSingle(ctx, country+"/FP.CPI.TOTL.ZG", start, end)
		if err != nil {
			fmt.Printf("✗ Failed to fetch inflation data for %s: %v\n", country, err)
			continue
		}

		inflationData := inflationResult.(*worldbank.ParsedData)
		fmt.Printf("✓ Fetched inflation data (%d observations)\n", len(inflationData.Dates))

		if len(inflationData.Dates) > 0 {
			lastIdx := len(inflationData.Dates) - 1
			fmt.Printf("  Latest (%s): %s%%\n", inflationData.Dates[lastIdx], inflationData.Values[lastIdx])
		}
	}

	// Popular World Bank indicators
	fmt.Println("\n--- Popular World Bank Indicators ---")
	fmt.Println("Economic:")
	fmt.Println("  NY.GDP.MKTP.CD - GDP (current US$)")
	fmt.Println("  NY.GDP.MKTP.KD.ZG - GDP growth (annual %)")
	fmt.Println("  NY.GDP.PCAP.CD - GDP per capita (current US$)")
	fmt.Println("  FP.CPI.TOTL.ZG - Inflation, consumer prices (annual %)")

	fmt.Println("\nPopulation:")
	fmt.Println("  SP.POP.TOTL - Population, total")
	fmt.Println("  SP.POP.GROW - Population growth (annual %)")
	fmt.Println("  SP.URB.TOTL.IN.ZS - Urban population (% of total)")

	fmt.Println("\nEducation:")
	fmt.Println("  SE.ADT.LITR.ZS - Literacy rate, adult total (% of people ages 15 and above)")
	fmt.Println("  SE.XPD.TOTL.GD.ZS - Government expenditure on education, total (% of GDP)")

	fmt.Println("\nHealth:")
	fmt.Println("  SH.XPD.CHEX.GD.ZS - Current health expenditure (% of GDP)")
	fmt.Println("  SP.DYN.LE00.IN - Life expectancy at birth, total (years)")

	fmt.Println("\nCountry codes:")
	fmt.Println("  USA - United States")
	fmt.Println("  CHN - China")
	fmt.Println("  GBR - United Kingdom")
	fmt.Println("  DEU - Germany")
	fmt.Println("  JPN - Japan")
	fmt.Println("  IND - India")
	fmt.Println("  Multiple: USA;CHN;GBR (semicolon separated)")

	fmt.Println("\n✓ World Bank example completed successfully!")
}
