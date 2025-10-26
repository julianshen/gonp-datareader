package datareader

import (
	"context"
	"fmt"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
	"github.com/julianshen/gonp-datareader/sources/alphavantage"
	"github.com/julianshen/gonp-datareader/sources/fred"
	"github.com/julianshen/gonp-datareader/sources/iex"
	"github.com/julianshen/gonp-datareader/sources/stooq"
	"github.com/julianshen/gonp-datareader/sources/worldbank"
	"github.com/julianshen/gonp-datareader/sources/yahoo"
)

var (
	// ErrUnknownSource is returned when an unknown data source is requested
	ErrUnknownSource = fmt.Errorf("unknown data source")
)

// DataReader creates a new reader for the specified source.
// Supported sources: "yahoo", "fred"
func DataReader(source string, opts *Options) (sources.Reader, error) {
	if source == "" {
		return nil, fmt.Errorf("%w: source cannot be empty", ErrUnknownSource)
	}

	// Convert Options to ClientOptions
	var clientOpts *internalhttp.ClientOptions
	var apiKey string
	if opts != nil {
		clientOpts = &internalhttp.ClientOptions{
			Timeout:    opts.Timeout,
			UserAgent:  opts.UserAgent,
			MaxRetries: opts.MaxRetries,
			RetryDelay: opts.RetryDelay,
			RateLimit:  opts.RateLimit,
			CacheDir:   opts.CacheDir,
			CacheTTL:   opts.CacheTTL,
		}
		apiKey = opts.APIKey
	}

	switch source {
	case "yahoo":
		return yahoo.NewYahooReader(clientOpts), nil
	case "fred":
		if apiKey != "" {
			return fred.NewFREDReaderWithAPIKey(clientOpts, apiKey), nil
		}
		return fred.NewFREDReader(clientOpts), nil
	case "worldbank":
		return worldbank.NewWorldBankReader(clientOpts), nil
	case "alphavantage":
		return alphavantage.NewAlphaVantageReader(clientOpts, apiKey), nil
	case "stooq":
		return stooq.NewStooqReader(clientOpts), nil
	case "iex":
		return iex.NewIEXReader(clientOpts, apiKey), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSource, source)
	}
}

// Read is a convenience function that creates a reader and fetches data for a single symbol.
func Read(ctx context.Context, symbol string, source string, start, end time.Time, opts *Options) (interface{}, error) {
	reader, err := DataReader(source, opts)
	if err != nil {
		return nil, err
	}

	return reader.ReadSingle(ctx, symbol, start, end)
}

// ListSources returns a list of available data source names.
func ListSources() []string {
	return []string{
		"yahoo",
		"fred",
		"worldbank",
		"alphavantage",
		"stooq",
		"iex",
	}
}
