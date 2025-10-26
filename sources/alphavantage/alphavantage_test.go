package alphavantage_test

import (
	"testing"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/alphavantage"
)

func TestNewAlphaVantageReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}
	apiKey := "test_api_key"

	reader := alphavantage.NewAlphaVantageReader(opts, apiKey)

	if reader == nil {
		t.Fatal("NewAlphaVantageReader returned nil")
	}

	if reader.Name() != "alphavantage" {
		t.Errorf("Expected name 'alphavantage', got %q", reader.Name())
	}
}

func TestNewAlphaVantageReader_RequiresAPIKey(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	// Empty API key should still create reader (validation happens at request time)
	reader := alphavantage.NewAlphaVantageReader(opts, "")

	if reader == nil {
		t.Fatal("NewAlphaVantageReader returned nil")
	}
}
