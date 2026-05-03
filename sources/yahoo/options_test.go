package yahoo_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/julianshen/gonp-datareader/sources/yahoo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseOptionsJSON_ValidResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/options_aapl.json")
	require.NoError(t, err)

	chain, err := yahoo.ParseOptionsJSON(bytes.NewReader(data))
	require.NoError(t, err)
	assert.NotNil(t, chain)
	assert.NotEmpty(t, chain.Calls)
	assert.NotEmpty(t, chain.Puts)

	call := chain.Calls[0]
	assert.NotEmpty(t, call.ContractSymbol)
	assert.Greater(t, call.Strike, 0.0)
}

func TestParseOptionsJSON_EmptyResult(t *testing.T) {
	jsonData := `{"optionChain":{"result":[],"error":null}}`
	_, err := yahoo.ParseOptionsJSON(strings.NewReader(jsonData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no options data")
}

func TestParseOptionsJSON_YahooError(t *testing.T) {
	jsonData := `{"optionChain":{"result":[],"error":{"code":"Not Found","description":"No data found"}}}`
	_, err := yahoo.ParseOptionsJSON(strings.NewReader(jsonData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yahoo finance error")
}

func TestOptionsReader_Struct(t *testing.T) {
	reader := yahoo.NewOptionsReader(nil)
	assert.NotNil(t, reader)
	assert.NotNil(t, reader.Client())
}

func TestOptionsReader_Name(t *testing.T) {
	reader := yahoo.NewOptionsReader(nil)
	assert.Equal(t, "Yahoo Finance Options", reader.Name())
}
