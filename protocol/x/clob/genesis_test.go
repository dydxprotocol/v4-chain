package clob_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/dydxprotocol/v4/x/perpetuals"
	"github.com/dydxprotocol/v4/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	tests := map[string]struct {
		// Genesis state.
		genesis types.GenesisState

		// Expectations.
		expectedErr     string
		expectedErrType error
	}{
		"Genesis state is valid": {
			genesis: types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:                   uint32(0),
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:                   uint32(1),
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
		},
		"Genesis state is valid when bankruptcy adjustment ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion * 10,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
		},
		"Genesis state is valid when min position liquidated is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   0,
						MaxPositionPortionLiquidatedPpm: lib.OneMillion,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
		},
		"Genesis state is invalid when there is no metadata on a CLOB pair": {
			genesis: types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id:                   uint32(0),
						Metadata:             nil,
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:                   uint32(1),
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "Asset orders are not implemented",
			expectedErrType: types.ErrInvalidClobPairParameter,
		},
		"Genesis state is invalid when there is a spot metadata on a CLOB pair": {
			genesis: types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Metadata: &types.ClobPair_SpotClobMetadata{
							SpotClobMetadata: &types.SpotClobMetadata{
								BaseAssetId:  0,
								QuoteAssetId: 1,
							},
						},
						Id:                   uint32(0),
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:                   uint32(1),
						StepBaseQuantums:     5,
						SubticksPerTick:      5,
						MinOrderBaseQuantums: 10,
						Status:               types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "Asset orders are not implemented",
			expectedErrType: types.ErrInvalidClobPairParameter,
		},
		"Genesis state is invalid when spread to maintenance margin ratio ppm is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 0,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "0 is not a valid SpreadToMaintenanceMarginRatioPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when spread to maintenance margin ratio ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion + 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "1000001 is not a valid SpreadToMaintenanceMarginRatioPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when bankruptcy adjustment ppm is less than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion - 1,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "999999 is not a valid BankruptcyAdjustmentPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max liquidation fee ppm is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  0,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "0 is not a valid MaxLiquidationFeePpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max liquidation fee ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  lib.OneMillion + 1,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "1000001 is not a valid MaxLiquidationFeePpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max position portion liquidated ppm is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   10,
						MaxPositionPortionLiquidatedPpm: 0,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "0 is not a valid MaxPositionPortionLiquidatedPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max position portion liquidated ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   10,
						MaxPositionPortionLiquidatedPpm: lib.OneMillion + 1,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedErr:     "1000001 is not a valid MaxPositionPortionLiquidatedPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max notional liquidated is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    0,
						MaxQuantumsInsuranceLost: 100_000_000_000_000,
					},
				},
			},
			expectedErr:     "0 is not a valid MaxNotionalLiquidated",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max quantums insurance lost is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    100_000_000_000_000,
						MaxQuantumsInsuranceLost: 0,
					},
				},
			},
			expectedErr:     "0 is not a valid MaxQuantumsInsuranceLost",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx, k, priceKeeper, _, perpetualsKeeper, _, _, _ :=
				keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
			ctx = ctx.WithBlockTime(constants.TimeT)

			prices.InitGenesis(ctx, *priceKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// If we expect a panic, verify that initializing the genesis state causes a panic and
			// end the test.
			if tc.expectedErr != "" {
				require.PanicsWithError(
					t,
					sdkerrors.Wrap(
						tc.expectedErrType,
						tc.expectedErr,
					).Error(),
					func() { clob.InitGenesis(ctx, *k, tc.genesis) },
				)
				return
			}

			// Initialize the CLOB genesis state.
			clob.InitGenesis(ctx, *k, tc.genesis)

			require.True(
				t,
				constants.TimeT.Equal(
					k.MustGetBlockTimeForLastCommittedBlock(ctx),
				),
			)

			// Export the CLOB genesis state and verify expectations.
			got := clob.ExportGenesis(ctx, *k)
			require.NotNil(t, got)
			require.Equal(t, tc.genesis.ClobPairs, got.ClobPairs)
			require.Equal(t, tc.genesis.LiquidationsConfig, got.LiquidationsConfig)

			// The number of CLOB pairs in the store should match the amount created thus far.
			numClobPairs := k.GetNumClobPairs(ctx)
			require.Equal(t, uint32(len(got.ClobPairs)), numClobPairs)
		})
	}
}
