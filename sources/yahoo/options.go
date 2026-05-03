package yahoo

import "time"

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
