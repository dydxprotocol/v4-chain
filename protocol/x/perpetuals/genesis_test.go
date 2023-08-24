package perpetuals_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	pricesGenesisState := constants.Prices_DefaultGenesisState
	genesisState := constants.Perpetuals_DefaultGenesisState

	ctx, k, priceKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	prices.InitGenesis(ctx, *priceKeeper, pricesGenesisState)
	perpetuals.InitGenesis(ctx, *k, genesisState)
	assertLiquidityTierUpsertEventsInIndexerBlock(t, k, ctx, len(genesisState.LiquidityTiers))
	got := perpetuals.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState) //nolint:staticcheck
	nullify.Fill(got)           //nolint:staticcheck

	require.ElementsMatch(t, genesisState.Perpetuals, got.Perpetuals)
	require.Equal(t, genesisState.Params, got.Params)
}

func TestGenesis_Failure(t *testing.T) {
	tests := map[string]struct {
		marketId                  uint32
		ticker                    string
		initialMarginPpm          uint32
		maintenanceFractionPpm    uint32
		basePositionNotional      uint64
		impactNotional            uint64
		fundingRateClampFactorPpm uint32
		premiumVoteClampFactorPpm uint32
		minNumVotesPerSample      uint32
	}{
		"MarketId doesn't reference a valid Market": {
			marketId:                  999,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Ticker is empty": {
			marketId:                  0,
			ticker:                    "",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Initial Margin Ppm exceeds maximum": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          lib.OneMillion + 1,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion + 1,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Base Position Notional is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion + 1,
			basePositionNotional:      0,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Impact Notional is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion + 1,
			basePositionNotional:      1,
			impactNotional:            0,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Funding Rate Clamp Factor Ppm is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 0,
			premiumVoteClampFactorPpm: 1,
			minNumVotesPerSample:      0,
		},
		"Premium Vote Clamp Factor Ppm is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion,
			basePositionNotional:      1,
			impactNotional:            1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 0,
			minNumVotesPerSample:      0,
		},
	}

	// Test setup.
	ctx, k, priceKeeper, _, _ := keepertest.PerpetualsKeepers(t)

	pricesGenesisState := constants.Prices_DefaultGenesisState
	prices.InitGenesis(ctx, *priceKeeper, pricesGenesisState)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			genesisState := types.GenesisState{
				LiquidityTiers: []types.LiquidityTier{
					{
						Name:                   "",
						InitialMarginPpm:       tc.initialMarginPpm,
						MaintenanceFractionPpm: tc.maintenanceFractionPpm,
						BasePositionNotional:   tc.basePositionNotional,
						ImpactNotional:         tc.impactNotional,
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: tc.fundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: tc.premiumVoteClampFactorPpm,
				},
				Perpetuals: []types.Perpetual{
					{
						Params: types.PerpetualParams{
							MarketId:      tc.marketId,
							Ticker:        tc.ticker,
							LiquidityTier: 0,
						},
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			}

			require.Panics(t, func() {
				perpetuals.InitGenesis(ctx, *k, genesisState)
			})
		})
	}
}

// assertLiquidityTierUpsertEventsInIndexerBlock checks the number of liquidity tier upsert events
// included in the Indexer block kafka message.
func assertLiquidityTierUpsertEventsInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	numLiquidityTiers int,
) {
	liquidityTierUpsertEvents := keepertest.GetLiquidityTierUpsertEventsFromIndexerBlock(ctx, k)
	require.Len(t, liquidityTierUpsertEvents, numLiquidityTiers)
}
