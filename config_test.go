package datareader_test

import (
	"testing"
	"time"

	datareader "github.com/julianshen/gonp-datareader"
)

func TestDefaultOptions(t *testing.T) {
	opts := datareader.DefaultOptions()

	if opts == nil {
		t.Fatal("DefaultOptions() returned nil")
	}

	// Test default timeout is set
	if opts.Timeout == 0 {
		t.Error("Expected default timeout to be set, got 0")
	}

	// Test default timeout is reasonable (30 seconds)
	expectedTimeout := 30 * time.Second
	if opts.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, opts.Timeout)
	}

	// Test default max retries
	if opts.MaxRetries < 0 {
		t.Errorf("Expected MaxRetries >= 0, got %d", opts.MaxRetries)
	}
}

func TestOptions_CustomValues(t *testing.T) {
	tests := []struct {
		name string
		opts *datareader.Options
	}{
		{
			name: "custom timeout",
			opts: &datareader.Options{
				Timeout: 60 * time.Second,
			},
		},
		{
			name: "with API key",
			opts: &datareader.Options{
				APIKey:  "test-api-key",
				Timeout: 30 * time.Second,
			},
		},
		{
			name: "with cache enabled",
			opts: &datareader.Options{
				EnableCache: true,
				CacheDir:    "/tmp/cache",
			},
		},
		{
			name: "with rate limit",
			opts: &datareader.Options{
				RateLimit: 10.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opts == nil {
				t.Error("Options should not be nil")
			}
		})
	}
}
