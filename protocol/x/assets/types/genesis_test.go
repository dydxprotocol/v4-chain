package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState    *types.GenesisState
		expectedErr error
	}{
		"default is valid": {
			genState: types.DefaultGenesis(),
		},
		"valid genesis state": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        false,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
					{
						Id:               1,
						Symbol:           "BTC",
						Denom:            "btc-denom",
						HasMarket:        true,
						MarketId:         0,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
		},
		"empty genesis state": {
			genState: &types.GenesisState{
				Assets: []types.Asset{},
			},
			expectedErr: types.ErrNoAssetInGenesis,
		},
		"asset[0] not tdai": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        true,
						MarketId:         0,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrTDaiMustBeAssetZero,
		},
		"asset[0] is modified tdai": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        true,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrTDaiMustBeAssetZero,
		},
		"duplicated asset id": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        false,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
					{
						Id:               0,
						Denom:            "BTC",
						HasMarket:        true,
						MarketId:         0,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrAssetIdAlreadyExists,
		},
		"duplicated denom": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        false,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
					{
						Id:               1,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        true,
						MarketId:         0,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrAssetDenomAlreadyExists,
		},
		"gaps in asset id": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        false,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
					{
						Id:               2,
						Denom:            "BTC",
						HasMarket:        true,
						MarketId:         0,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrGapFoundInAssetId,
		},
		"MarketId non-zero when HasMarket is false": {
			genState: &types.GenesisState{
				Assets: []types.Asset{
					{
						Id:               0,
						Symbol:           types.AssetTDai.Symbol,
						Denom:            types.AssetTDai.Denom,
						DenomExponent:    types.AssetTDai.DenomExponent,
						HasMarket:        false,
						AtomicResolution: lib.QuoteCurrencyAtomicResolution,
						AssetYieldIndex:  "1/1",
					},
					{
						Id:               1,
						Denom:            "USDT",
						HasMarket:        false,
						MarketId:         1,
						AtomicResolution: int32(-6),
						AssetYieldIndex:  "1/1",
					},
				},
			},
			expectedErr: types.ErrInvalidMarketId,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
