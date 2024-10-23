package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestMarketParam_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		input     types.MarketParam
		expErrMsg string
	}{
		{
			name: "Valid MarketParam",
			input: types.MarketParam{
				Pair:              "BTC-USD",
				MinPriceChangePpm: 1_000,
			},
			expErrMsg: "",
		},
		{
			name: "Empty pair",
			input: types.MarketParam{
				Pair:              "",
				MinPriceChangePpm: 1_000,
			},
			expErrMsg: "Pair cannot be empty",
		},
		{
			name: "Invalid MinPriceChangePpm",
			input: types.MarketParam{
				Pair:              "BTC-USD",
				MinPriceChangePpm: 0,
			},
			expErrMsg: "Min price change in parts-per-million must be greater than 0",
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
