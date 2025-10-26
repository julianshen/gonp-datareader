package worldbank_test

import (
	"testing"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/worldbank"
)

func TestNewWorldBankReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	reader := worldbank.NewWorldBankReader(opts)

	if reader == nil {
		t.Fatal("NewWorldBankReader returned nil")
	}

	if reader.Name() != "worldbank" {
		t.Errorf("Expected name 'worldbank', got %q", reader.Name())
	}
}
