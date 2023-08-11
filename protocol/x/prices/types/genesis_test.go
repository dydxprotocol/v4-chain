package types_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/prices/types"
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
				ExchangeFeeds: []types.ExchangeFeed{
					{
						Id:   0,
						Name: constants.CoinbaseExchangeName,
					},
				},
				Markets: []types.Market{
					{
						Id:   0,
						Pair: constants.BtcUsdPair,
					},
					{
						Id:   1,
						Pair: constants.EthUsdPair,
					},
				},
			},
			expectedError: nil,
		},
		"invalid: duplicate market ids": {
			genState: &types.GenesisState{
				Markets: []types.Market{
					{
						Id:   0,
						Pair: constants.BtcUsdPair,
					},
					{
						Id:   0, // duplicate
						Pair: constants.EthUsdPair,
					},
				},
			},
			expectedError: errors.New("duplicated market id"),
		},
		"invalid: found gap in market id": {
			genState: &types.GenesisState{
				Markets: []types.Market{
					{
						Id:   0,
						Pair: constants.BtcUsdPair,
					},
					{
						Id:   2, // gap
						Pair: constants.EthUsdPair,
					},
				},
			},
			expectedError: errors.New("found gap in market id"),
		},
		"invalid: pair not set": {
			genState: &types.GenesisState{
				Markets: []types.Market{
					{
						Id:   0,
						Pair: "",
					},
				},
			},
			expectedError: errors.New("Pair must be non-empty string"),
		},
		"invalid: duplicate exchange feed ids": {
			genState: &types.GenesisState{
				ExchangeFeeds: []types.ExchangeFeed{
					{
						Id:   0,
						Name: constants.CoinbaseExchangeName,
					},
					{
						Id:   0, // duplicate
						Name: constants.BinanceExchangeName,
					},
				},
			},
			expectedError: errors.New("duplicated exchange feed id"),
		},
		"invalid: found gap in exchange feed id": {
			genState: &types.GenesisState{
				ExchangeFeeds: []types.ExchangeFeed{
					{
						Id:   0,
						Name: constants.CoinbaseExchangeName,
					},
					{
						Id:   2, // gap
						Name: constants.BinanceExchangeName,
					},
				},
			},
			expectedError: errors.New("found gap in exchange feed id"),
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
