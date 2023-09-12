package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestMarketParam_Validate(t *testing.T) {
	validExchangeConfigJson := `{"exchanges":[{"exchangeName":"Binance","ticker":"BTCUSDT"}]}`
	testCases := []struct {
		name      string
		input     types.MarketParam
		expErrMsg string
	}{
		{
			name: "Valid MarketParam",
			input: types.MarketParam{
				Pair:               "BTC-USD",
				MinExchanges:       1,
				MinPriceChangePpm:  1_000,
				ExchangeConfigJson: validExchangeConfigJson,
			},
			expErrMsg: "",
		},
		{
			name: "Empty pair",
			input: types.MarketParam{
				Pair:               "",
				MinExchanges:       1,
				MinPriceChangePpm:  1_000,
				ExchangeConfigJson: validExchangeConfigJson,
			},
			expErrMsg: "Pair cannot be empty",
		},
		{
			name: "Invalid MinPriceChangePpm",
			input: types.MarketParam{
				Pair:               "BTC-USD",
				MinExchanges:       1,
				MinPriceChangePpm:  0,
				ExchangeConfigJson: validExchangeConfigJson,
			},
			expErrMsg: "Min price change in parts-per-million must be greater than 0",
		},
		{
			name: "Empty ExchangeConfigJson",
			input: types.MarketParam{
				Pair:               "BTC-USD",
				MinExchanges:       1,
				MinPriceChangePpm:  1_000,
				ExchangeConfigJson: "",
			},
			expErrMsg: "ExchangeConfigJson string is not valid",
		},
		{
			name: "Typo in ExchangeConfigJson",
			input: types.MarketParam{
				Pair:               "BTC-USD",
				MinExchanges:       1,
				MinPriceChangePpm:  1_000,
				ExchangeConfigJson: `{"exchanges":[]`, // missing a bracket
			},
			expErrMsg: "ExchangeConfigJson string is not valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.expErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expErrMsg)
			}
		})
	}
}
