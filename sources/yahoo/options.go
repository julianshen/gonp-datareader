package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

const optionsAPIURL = "https://query1.finance.yahoo.com/v7/finance/options/%s"

// OptionsReader fetches options chain data from Yahoo Finance.
type OptionsReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewOptionsReader creates a new options chain reader.
func NewOptionsReader(opts *internalhttp.ClientOptions) *OptionsReader {
	return NewOptionsReaderWithBaseURL(opts, optionsAPIURL)
}

// NewOptionsReaderWithBaseURL creates an options reader with a custom base URL.
func NewOptionsReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *OptionsReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}
	return &OptionsReader{
		BaseSource: sources.NewBaseSource("yahoo-options"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name.
func (o *OptionsReader) Name() string {
	return "Yahoo Finance Options"
}

// Client returns the underlying HTTP client.
func (o *OptionsReader) Client() *internalhttp.RetryableClient {
	return o.client
}

// GetOptionsChain fetches the options chain for a symbol.
// If expiration is nil, returns the nearest expiration.
func (o *OptionsReader) GetOptionsChain(ctx context.Context, symbol string, expiration *time.Time) (*OptionsChain, error) {
	if err := o.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}

	url := fmt.Sprintf(o.baseURL, symbol)
	if expiration != nil {
		url = fmt.Sprintf("%s?date=%d", url, expiration.Unix())
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch options: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yahoo returned %d: %s", resp.StatusCode, string(body))
	}

	return ParseOptionsJSON(resp.Body)
}

// OptionContract represents a single options contract.
type OptionContract struct {
	ContractSymbol string    `json:"contractSymbol"`
	Strike         float64   `json:"strike"`
	Expiration     time.Time `json:"expiration"`
	Type           string    `json:"type"`
	LastPrice      float64   `json:"lastPrice"`
	Bid            float64   `json:"bid"`
	Ask            float64   `json:"ask"`
	Change         float64   `json:"change"`
	PercentChange  float64   `json:"percentChange"`
	Volume         int64     `json:"volume"`
	OpenInterest   int64     `json:"openInterest"`
	ImpliedVol     float64   `json:"impliedVolatility"`
	InTheMoney     bool      `json:"inTheMoney"`
}

// OptionsChain represents all contracts for a single expiration date.
type OptionsChain struct {
	ExpirationDate time.Time        `json:"expirationDate"`
	Calls          []OptionContract `json:"calls"`
	Puts           []OptionContract `json:"puts"`
}

// rawOptionContract is used to decode JSON without the time.Time field.
type rawOptionContract struct {
	ContractSymbol string  `json:"contractSymbol"`
	Strike         float64 `json:"strike"`
	Expiration     int64   `json:"expiration"`
	Type           string  `json:"type"`
	LastPrice      float64 `json:"lastPrice"`
	Bid            float64 `json:"bid"`
	Ask            float64 `json:"ask"`
	Change         float64 `json:"change"`
	PercentChange  float64 `json:"percentChange"`
	Volume         int64   `json:"volume"`
	OpenInterest   int64   `json:"openInterest"`
	ImpliedVol     float64 `json:"impliedVolatility"`
	InTheMoney     bool    `json:"inTheMoney"`
}

// ParseOptionsJSON parses Yahoo Finance options JSON response.
func ParseOptionsJSON(r io.Reader) (*OptionsChain, error) {
	var resp struct {
		OptionChain struct {
			Result []struct {
				ExpirationDates []int64   `json:"expirationDates"`
				Strikes         []float64 `json:"strikes"`
				Options         []struct {
					ExpirationDate int64               `json:"expirationDate"`
					Calls          []rawOptionContract `json:"calls"`
					Puts           []rawOptionContract `json:"puts"`
				} `json:"options"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"optionChain"`
	}

	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode options JSON: %w", err)
	}

	if resp.OptionChain.Error != nil {
		return nil, fmt.Errorf("yahoo finance error: %v", resp.OptionChain.Error)
	}

	if len(resp.OptionChain.Result) == 0 || len(resp.OptionChain.Result[0].Options) == 0 {
		return nil, fmt.Errorf("no options data found")
	}

	result := resp.OptionChain.Result[0]
	opt := result.Options[0]

	calls := make([]OptionContract, len(opt.Calls))
	for i, rc := range opt.Calls {
		calls[i] = OptionContract{
			ContractSymbol: rc.ContractSymbol,
			Strike:         rc.Strike,
			Expiration:     time.Unix(rc.Expiration, 0),
			Type:           "CALL",
			LastPrice:      rc.LastPrice,
			Bid:            rc.Bid,
			Ask:            rc.Ask,
			Change:         rc.Change,
			PercentChange:  rc.PercentChange,
			Volume:         rc.Volume,
			OpenInterest:   rc.OpenInterest,
			ImpliedVol:     rc.ImpliedVol,
			InTheMoney:     rc.InTheMoney,
		}
	}

	puts := make([]OptionContract, len(opt.Puts))
	for i, rp := range opt.Puts {
		puts[i] = OptionContract{
			ContractSymbol: rp.ContractSymbol,
			Strike:         rp.Strike,
			Expiration:     time.Unix(rp.Expiration, 0),
			Type:           "PUT",
			LastPrice:      rp.LastPrice,
			Bid:            rp.Bid,
			Ask:            rp.Ask,
			Change:         rp.Change,
			PercentChange:  rp.PercentChange,
			Volume:         rp.Volume,
			OpenInterest:   rp.OpenInterest,
			ImpliedVol:     rp.ImpliedVol,
			InTheMoney:     rp.InTheMoney,
		}
	}

	return &OptionsChain{
		ExpirationDate: time.Unix(opt.ExpirationDate, 0),
		Calls:          calls,
		Puts:           puts,
	}, nil
}
