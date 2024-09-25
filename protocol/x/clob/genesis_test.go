package clob_test

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"

	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices"
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
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							Limit:     200,
							NumBlocks: 1,
						},
					},
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							Limit:     2,
							NumBlocks: 1,
						},
						{
							Limit:     20,
							NumBlocks: 100,
						},
					},
				},
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(0),
							Limit:          0,
						},
						{
							UsdTncRequired: dtypes.NewInt(20),
							Limit:          1,
						},
						{
							UsdTncRequired: dtypes.NewInt(100),
							Limit:          5,
						},
						{
							UsdTncRequired: dtypes.NewInt(1000),
							Limit:          10,
						},
						{
							UsdTncRequired: dtypes.NewInt(10000),
							Limit:          100,
						},
						{
							UsdTncRequired: dtypes.NewInt(100000),
							Limit:          1000,
						},
					},
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(0),
							Limit:          0,
						},
						{
							UsdTncRequired: dtypes.NewInt(20),
							Limit:          1,
						},
						{
							UsdTncRequired: dtypes.NewInt(100),
							Limit:          5,
						},
						{
							UsdTncRequired: dtypes.NewInt(1000),
							Limit:          10,
						},
						{
							UsdTncRequired: dtypes.NewInt(10000),
							Limit:          100,
						},
						{
							UsdTncRequired: dtypes.NewInt(100000),
							Limit:          200,
						},
					},
				},
				ClobPairs: []types.ClobPair{
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:               uint32(0),
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 1,
							},
						},
						Id:               uint32(1),
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
		},
		"Genesis state is valid when bankruptcy adjustment ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion * 10,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
		},
		"Genesis state is valid when min position liquidated is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm:             lib.OneMillion,
					ValidatorFeePpm:                 200_000,
					LiquidityFeePpm:                 800_000,
					FillablePriceConfig:             constants.FillablePriceConfig_Default,
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
		},
		"Genesis state is invalid when there is no metadata on a CLOB pair": {
			genesis: types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id:               uint32(0),
						Metadata:         nil,
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:               uint32(1),
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
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
						Id:               uint32(0),
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
					{
						Metadata: &types.ClobPair_PerpetualClobMetadata{
							PerpetualClobMetadata: &types.PerpetualClobMetadata{
								PerpetualId: 0,
							},
						},
						Id:               uint32(1),
						StepBaseQuantums: 5,
						SubticksPerTick:  5,
						Status:           types.ClobPair_STATUS_ACTIVE,
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr:     "Asset orders are not implemented",
			expectedErrType: types.ErrInvalidClobPairParameter,
		},
		"Genesis state is invalid when spread to maintenance margin ratio ppm is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 0,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr:     "0 is not a valid SpreadToMaintenanceMarginRatioPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is valid when spread to maintenance margin ratio ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion + 1,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
		},
		"Genesis state is invalid when bankruptcy adjustment ppm is less than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion - 1,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr:     "999999 is not a valid BankruptcyAdjustmentPpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max liquidation fee ppm is 0": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm:             0,
					ValidatorFeePpm:                 200_000,
					LiquidityFeePpm:                 800_000,
					FillablePriceConfig:             constants.FillablePriceConfig_Default,
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr:     "0 is not a valid InsuranceFundFeePpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when max liquidation fee ppm is greater than one million": {
			genesis: types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm:             lib.OneMillion + 1,
					ValidatorFeePpm:                 200_000,
					LiquidityFeePpm:                 800_000,
					FillablePriceConfig:             constants.FillablePriceConfig_Default,
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr:     "1000001 is not a valid InsuranceFundFeePpm",
			expectedErrType: types.ErrInvalidLiquidationsConfig,
		},
		"Genesis state is invalid when BlockRateLimitConfiguration is invalid": {
			genesis: types.GenesisState{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							Limit:     1,
							NumBlocks: 0,
						},
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr: "0 is not a valid NumBlocks for MaxShortTermOrdersAndCancelsPerNBlocks rate limit " +
				"{NumBlocks:0 Limit:1}",
			expectedErrType: types.ErrInvalidBlockRateLimitConfig,
		},
		"Genesis state is invalid when EquityTierLimitConfiguration is invalid": {
			genesis: types.GenesisState{
				EquityTierLimitConfig: types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							Limit:          0,
							UsdTncRequired: dtypes.NewInt(-1),
						},
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 100_000,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedErr: "-1 is not a valid UsdTncRequired for ShortTermOrderEquityTiers equity tier limit " +
				"{UsdTncRequired:-1 Limit:0}",
			expectedErrType: types.ErrInvalidEquityTierLimitConfig,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			ctx := ks.Ctx.WithBlockTime(constants.TimeT)

			prices.InitGenesis(ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			if tc.expectedErr == "" {
				for i, clobPair := range tc.genesis.ClobPairs {
					perpetualId := clobtest.MustPerpetualId(clobPair)
					perpetual := constants.Perpetuals_DefaultGenesisState.Perpetuals[perpetualId]
					mockIndexerEventManager.On("AddTxnEvent",
						ctx,
						indexerevents.SubtypePerpetualMarket,
						indexerevents.PerpetualMarketEventVersion,
						indexer_manager.GetBytes(
							indexerevents.NewPerpetualMarketCreateEvent(
								perpetualId,
								uint32(i),
								perpetual.Params.Ticker,
								perpetual.Params.MarketId,
								clobPair.Status,
								clobPair.QuantumConversionExponent,
								perpetual.Params.AtomicResolution,
								clobPair.SubticksPerTick,
								clobPair.StepBaseQuantums,
								perpetual.Params.LiquidityTier,
								perpetual.Params.MarketType,
								perpetual.Params.DangerIndexPpm,
								fmt.Sprintf("%d", perpetual.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
							),
						),
					).Once().Return()
				}
			}
			// If we expect a panic, verify that initializing the genesis state causes a panic and
			// end the test.
			if tc.expectedErr != "" {
				require.PanicsWithError(
					t,
					errorsmod.Wrap(
						tc.expectedErrType,
						tc.expectedErr,
					).Error(),
					func() { clob.InitGenesis(ctx, ks.ClobKeeper, tc.genesis) },
				)
				return
			}

			// Initialize the CLOB genesis state.
			clob.InitGenesis(ctx, ks.ClobKeeper, tc.genesis)

			// Export the CLOB genesis state and verify expectations.
			got := clob.ExportGenesis(ctx, *ks.ClobKeeper)
			require.NotNil(t, got)
			require.Equal(t, tc.genesis.ClobPairs, got.ClobPairs)
			require.Equal(t, tc.genesis.LiquidationsConfig, got.LiquidationsConfig)
			require.Equal(t, tc.genesis.BlockRateLimitConfig, got.BlockRateLimitConfig)
			require.Equal(t, tc.genesis.EquityTierLimitConfig, got.EquityTierLimitConfig)
		})
	}
}
