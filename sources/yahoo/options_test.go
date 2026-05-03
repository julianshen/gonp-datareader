package yahoo_test

import (
	"testing"

	"github.com/julianshen/gonp-datareader/sources/yahoo"
	"github.com/stretchr/testify/assert"
)

func TestOptionContract_StructTags(t *testing.T) {
	contract := yahoo.OptionContract{
		ContractSymbol: "AAPL250516C00150000",
		Strike:         150.0,
		Type:           "CALL",
		LastPrice:      5.25,
		Bid:            5.20,
		Ask:            5.30,
		Change:         0.50,
		PercentChange:  10.5,
		Volume:         1500,
		OpenInterest:   2500,
		ImpliedVol:     0.35,
		InTheMoney:     true,
	}
	assert.Equal(t, "AAPL250516C00150000", contract.ContractSymbol)
	assert.Equal(t, 150.0, contract.Strike)
	assert.Equal(t, "CALL", contract.Type)
}
