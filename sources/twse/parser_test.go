package twse

import (
	"testing"
	"time"
)

// TestRocToGregorian tests conversion from ROC date string to Gregorian time.Time
func TestRocToGregorian(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "ROC 1141031 = 2025-10-31",
			rocDate:   "1141031",
			wantYear:  2025,
			wantMonth: time.October,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "ROC 1130101 = 2024-01-01",
			rocDate:   "1130101",
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC 1131231 = 2024-12-31",
			rocDate:   "1131231",
			wantYear:  2024,
			wantMonth: time.December,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "ROC 1120229 = 2023-02-28 (not leap year)",
			rocDate:   "1120229",
			wantYear:  0, // Invalid date
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "ROC 1120228 = 2023-02-28",
			rocDate:   "1120228",
			wantYear:  2023,
			wantMonth: time.February,
			wantDay:   28,
			wantErr:   false,
		},
		{
			name:      "ROC year 100 = 2011",
			rocDate:   "1000101",
			wantYear:  2011,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "empty string",
			rocDate:   "",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - too short",
			rocDate:   "11410",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - too long",
			rocDate:   "11410311",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid format - non-numeric",
			rocDate:   "abc1031",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid month",
			rocDate:   "1141331",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "invalid day",
			rocDate:   "1140132",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rocToGregorian(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("rocToGregorian(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear {
					t.Errorf("rocToGregorian(%q) year = %d, want %d", tt.rocDate, got.Year(), tt.wantYear)
				}
				if got.Month() != tt.wantMonth {
					t.Errorf("rocToGregorian(%q) month = %v, want %v", tt.rocDate, got.Month(), tt.wantMonth)
				}
				if got.Day() != tt.wantDay {
					t.Errorf("rocToGregorian(%q) day = %d, want %d", tt.rocDate, got.Day(), tt.wantDay)
				}
			}
		})
	}
}

// TestGregorianToROC tests conversion from Gregorian time.Time to ROC date string
func TestGregorianToROC(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			want: "1141031",
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1130101",
		},
		{
			name: "2024-12-31",
			date: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			want: "1131231",
		},
		{
			name: "2023-02-28",
			date: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			want: "1120228",
		},
		{
			name: "2011-01-01 (ROC 100)",
			date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1000101",
		},
		{
			name: "2020-02-29 (leap year)",
			date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			want: "1090229",
		},
		{
			name: "1912-01-01 (ROC year 1)",
			date: time.Date(1912, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "0010101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gregorianToROC(tt.date)
			if got != tt.want {
				t.Errorf("gregorianToROC(%v) = %q, want %q", tt.date, got, tt.want)
			}
		})
	}
}

// TestParseROCDate tests parsing ROC date strings to time.Time
func TestParseROCDate(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "valid date 1141031",
			rocDate:   "1141031",
			wantYear:  2025,
			wantMonth: time.October,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "valid date 1130101",
			rocDate:   "1130101",
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "invalid format",
			rocDate:   "invalid",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "empty string",
			rocDate:   "",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseROCDate(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseROCDate(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear {
					t.Errorf("parseROCDate(%q) year = %d, want %d", tt.rocDate, got.Year(), tt.wantYear)
				}
				if got.Month() != tt.wantMonth {
					t.Errorf("parseROCDate(%q) month = %v, want %v", tt.rocDate, got.Month(), tt.wantMonth)
				}
				if got.Day() != tt.wantDay {
					t.Errorf("parseROCDate(%q) day = %d, want %d", tt.rocDate, got.Day(), tt.wantDay)
				}
			}
		})
	}
}

// TestFormatROCDate tests formatting time.Time to ROC date string
func TestFormatROCDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			want: "1141031",
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "1130101",
		},
		{
			name: "single digit month and day",
			date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			want: "1130105",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatROCDate(tt.date)
			if got != tt.want {
				t.Errorf("formatROCDate(%v) = %q, want %q", tt.date, got, tt.want)
			}
		})
	}
}

// TestROCDateRoundTrip tests that conversion is reversible
func TestROCDateRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "2025-10-31",
			date: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "2024-01-01",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "2020-02-29 (leap year)",
			date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to ROC and back
			rocStr := formatROCDate(tt.date)
			got, err := parseROCDate(rocStr)
			if err != nil {
				t.Fatalf("parseROCDate(%q) error = %v", rocStr, err)
			}

			// Compare year, month, day (ignore time components)
			if got.Year() != tt.date.Year() || got.Month() != tt.date.Month() || got.Day() != tt.date.Day() {
				t.Errorf("round trip failed: %v -> %s -> %v", tt.date, rocStr, got)
			}
		})
	}
}

// TestROCDateEdgeCases tests edge cases like leap years and year boundaries
func TestROCDateEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		rocDate   string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "leap year Feb 29, 2024",
			rocDate:   "1130229",
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
			wantErr:   false,
		},
		{
			name:      "leap year Feb 29, 2020",
			rocDate:   "1090229",
			wantYear:  2020,
			wantMonth: time.February,
			wantDay:   29,
			wantErr:   false,
		},
		{
			name:      "non-leap year Feb 29, 2023 (invalid)",
			rocDate:   "1120229",
			wantYear:  0,
			wantMonth: 0,
			wantDay:   0,
			wantErr:   true,
		},
		{
			name:      "year boundary - Dec 31",
			rocDate:   "1131231",
			wantYear:  2024,
			wantMonth: time.December,
			wantDay:   31,
			wantErr:   false,
		},
		{
			name:      "year boundary - Jan 1",
			rocDate:   "1140101",
			wantYear:  2025,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC year 100 (2011)",
			rocDate:   "1000101",
			wantYear:  2011,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "ROC year 1 (1912)",
			rocDate:   "0010101",
			wantYear:  1912,
			wantMonth: time.January,
			wantDay:   1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rocToGregorian(tt.rocDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("rocToGregorian(%q) error = %v, wantErr %v", tt.rocDate, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() != tt.wantYear || got.Month() != tt.wantMonth || got.Day() != tt.wantDay {
					t.Errorf("rocToGregorian(%q) = %v, want %d-%02d-%02d",
						tt.rocDate, got, tt.wantYear, tt.wantMonth, tt.wantDay)
				}
			}
		})
	}
}
