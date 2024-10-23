package types_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState      *types.GenesisState
		expectedError error
	}{
		"valid: default": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid": {
			genState: &types.GenesisState{
				MarketParams: []types.MarketParam{
					{
						Id:                0,
						Pair:              constants.BtcUsdPair,
						MinPriceChangePpm: 1,
					},
					{
						Id:                1,
						Pair:              constants.EthUsdPair,
						MinPriceChangePpm: 1,
					},
				},
				MarketPrices: []types.MarketPrice{
					{
						Id:    0,
						Price: constants.FiveBillion,
					},
					{
						Id:    1,
						Price: constants.FiveBillion,
					},
				},
			},
			expectedError: nil,
		},
		"invalid: duplicate market param ids": {
			genState: &types.GenesisState{
				MarketParams: []types.MarketParam{
					{
						Id:                0,
						Pair:              constants.BtcUsdPair,
						MinPriceChangePpm: 1,
					},
					{
						Id:                0,
						Pair:              constants.EthUsdPair,
						MinPriceChangePpm: 1,
					},
				},
			},
			expectedError: errors.New("duplicated market param id"),
		},
		"invalid: market param invalid (pair unset)": {
			genState: &types.GenesisState{
				MarketParams: []types.MarketParam{
					{
						Id:   0,
						Pair: "",
					},
				},
			},
			expectedError: errorsmod.Wrap(types.ErrInvalidInput, "Pair cannot be empty"),
		},
		"invalid: mismatched number of market params and prices": {
			genState: &types.GenesisState{
				MarketParams: []types.MarketParam{
					{
						Id:                0,
						Pair:              constants.BtcUsdPair,
						MinPriceChangePpm: 1,
					},
					{
						Id:                1,
						Pair:              constants.EthUsdPair,
						MinPriceChangePpm: 1,
					},
				},
				MarketPrices: []types.MarketPrice{
					{
						Id:    0,
						Price: constants.FiveBillion,
					},
				},
			},
			expectedError: errors.New("expected the same number of market prices and market params"),
		},
		"invalid: market prices don't correspond to params": {
			genState: &types.GenesisState{
				MarketParams: []types.MarketParam{
					{
						Id:                0,
						Pair:              constants.BtcUsdPair,
						MinPriceChangePpm: 1,
					},
					{
						Id:                1,
						Pair:              constants.EthUsdPair,
						MinPriceChangePpm: 1,
					},
				},
				MarketPrices: []types.MarketPrice{
					{
						Id:    0,
						Price: constants.FiveBillion,
					},
					{
						Id:    2, // nonconsecutive id
						Price: constants.FiveBillion,
					},
				},
			},
			expectedError: errorsmod.Wrap(types.ErrInvalidInput, "market param id 1 does not match market price id 2"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}
