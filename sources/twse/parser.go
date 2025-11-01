package twse

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// rocEpochYear is the offset between ROC and Gregorian calendars
	// ROC Year 1 = Gregorian Year 1912
	rocEpochYear = 1911

	// Expected ROC date format: YYYMMDD (7 digits)
	// YYY = ROC year (3 digits)
	// MM = month (2 digits)
	// DD = day (2 digits)
	rocDateLength = 7
)

// rocToGregorian converts a ROC (Republic of China) date string to a Gregorian time.Time.
//
// ROC dates are formatted as "YYYMMDD" where:
//   - YYY is the ROC year (ROC Year = Gregorian Year - 1911)
//   - MM is the month (01-12)
//   - DD is the day (01-31)
//
// Examples:
//   - "1141031" -> October 31, 2025 (ROC 114 + 1911 = 2025)
//   - "1130101" -> January 1, 2024 (ROC 113 + 1911 = 2024)
//
// The function validates the date and returns an error if:
//   - The format is invalid (not 7 digits)
//   - The date components are invalid (e.g., month 13, day 32)
//   - The date doesn't exist (e.g., Feb 29 in non-leap year)
func rocToGregorian(rocDate string) (time.Time, error) {
	if len(rocDate) != rocDateLength {
		return time.Time{}, fmt.Errorf("invalid ROC date format: expected 7 digits (YYYMMDD), got %d digits", len(rocDate))
	}

	// Parse ROC year (first 3 digits)
	rocYearStr := rocDate[0:3]
	rocYear, err := strconv.Atoi(rocYearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ROC year: %w", err)
	}

	// Parse month (next 2 digits)
	monthStr := rocDate[3:5]
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month: %w", err)
	}

	// Parse day (last 2 digits)
	dayStr := rocDate[5:7]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %w", err)
	}

	// Convert ROC year to Gregorian year
	gregorianYear := rocYear + rocEpochYear

	// Create time.Time and let it validate the date
	// This will catch invalid dates like Feb 30, April 31, etc.
	date := time.Date(gregorianYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Verify the date is valid by checking if components match
	// If date is invalid, time.Date normalizes it (e.g., Feb 30 -> March 2)
	if date.Year() != gregorianYear || date.Month() != time.Month(month) || date.Day() != day {
		return time.Time{}, fmt.Errorf("invalid date: ROC %s (Gregorian %d-%02d-%02d does not exist)",
			rocDate, gregorianYear, month, day)
	}

	return date, nil
}

// gregorianToROC converts a Gregorian time.Time to a ROC date string.
//
// The returned string is in "YYYMMDD" format where:
//   - YYY is the ROC year (Gregorian Year - 1911)
//   - MM is the month (01-12)
//   - DD is the day (01-31)
//
// Examples:
//   - October 31, 2025 -> "1141031" (2025 - 1911 = 114)
//   - January 1, 2024 -> "1130101" (2024 - 1911 = 113)
func gregorianToROC(date time.Time) string {
	rocYear := date.Year() - rocEpochYear
	return fmt.Sprintf("%03d%02d%02d", rocYear, date.Month(), date.Day())
}

// parseROCDate parses a ROC date string into a time.Time.
//
// This is an alias for rocToGregorian for consistency with common naming patterns.
func parseROCDate(rocDate string) (time.Time, error) {
	return rocToGregorian(rocDate)
}

// formatROCDate formats a time.Time into a ROC date string.
//
// This is an alias for gregorianToROC for consistency with common naming patterns.
func formatROCDate(date time.Time) string {
	return gregorianToROC(date)
}
