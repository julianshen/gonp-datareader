package iex_test

import (
	"testing"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/iex"
)

func TestNewIEXReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}
	apiKey := "test_api_key"

	reader := iex.NewIEXReader(opts, apiKey)

	if reader == nil {
		t.Fatal("NewIEXReader returned nil")
	}

	if reader.Name() != "iex" {
		t.Errorf("Expected name 'iex', got %q", reader.Name())
	}
}

func TestNewIEXReader_RequiresAPIKey(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	// Empty API key should still create reader (validation happens at request time)
	reader := iex.NewIEXReader(opts, "")

	if reader == nil {
		t.Fatal("NewIEXReader returned nil")
	}
}
