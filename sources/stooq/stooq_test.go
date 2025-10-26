package stooq_test

import (
	"testing"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources/stooq"
)

func TestNewStooqReader(t *testing.T) {
	opts := &internalhttp.ClientOptions{
		Timeout: 30,
	}

	reader := stooq.NewStooqReader(opts)

	if reader == nil {
		t.Fatal("NewStooqReader returned nil")
	}

	if reader.Name() != "stooq" {
		t.Errorf("Expected name 'stooq', got %q", reader.Name())
	}
}
