package http_test

import (
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
)

func TestNewHTTPClient(t *testing.T) {
	client := internalhttp.NewHTTPClient(nil)

	if client == nil {
		t.Fatal("NewHTTPClient() returned nil")
	}
}

func TestHTTPClient_DefaultTimeout(t *testing.T) {
	client := internalhttp.NewHTTPClient(nil)

	if client.Timeout == 0 {
		t.Error("Expected default timeout to be set, got 0")
	}

	// Default should be 30 seconds
	expectedTimeout := 30 * time.Second
	if client.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, client.Timeout)
	}
}

func TestHTTPClient_CustomTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	opts := &internalhttp.ClientOptions{
		Timeout: customTimeout,
	}

	client := internalhttp.NewHTTPClient(opts)

	if client.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.Timeout)
	}
}

func TestHTTPClient_Do(t *testing.T) {
	client := internalhttp.NewHTTPClient(nil)

	// Verify client has Do method by checking it's a valid http.Client
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// Verify Transport is configured
	if client.Transport == nil {
		t.Error("Client.Transport should be configured")
	}
}

func TestClientOptions_Defaults(t *testing.T) {
	opts := internalhttp.DefaultClientOptions()

	if opts == nil {
		t.Fatal("DefaultClientOptions() returned nil")
	}

	if opts.Timeout == 0 {
		t.Error("Default timeout should be set")
	}

	if opts.UserAgent == "" {
		t.Error("Default UserAgent should be set")
	}
}
