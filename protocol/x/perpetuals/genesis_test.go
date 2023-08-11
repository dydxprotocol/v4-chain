package perpetuals_test

import (
	"testing"

	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/perpetuals"
	"github.com/dydxprotocol/v4/x/perpetuals/types"
	"github.com/dydxprotocol/v4/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	pricesGenesisState := constants.Prices_DefaultGenesisState
	genesisState := constants.Perpetuals_DefaultGenesisState

	ctx, k, priceKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	prices.InitGenesis(ctx, *priceKeeper, pricesGenesisState)
	perpetuals.InitGenesis(ctx, *k, genesisState)
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
		fundingRateClampFactorPpm uint32
		premiumVoteClampFactorPpm uint32
	}{
		"MarketId doesn't reference a valid Market": {
			marketId:                  999,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Ticker is empty": {
			marketId:                  0,
			ticker:                    "",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Initial Margin Ppm exceeds maximum": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          lib.OneMillion + 1,
			maintenanceFractionPpm:    0,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion + 1,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Base Position Notional is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion + 1,
			basePositionNotional:      0,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 1,
		},
		"Funding Rate Clamp Factor Ppm is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 0,
			premiumVoteClampFactorPpm: 1,
		},
		"Premium Vote Clamp Factor Ppm is zero": {
			marketId:                  0,
			ticker:                    "genesis_ticker",
			initialMarginPpm:          0,
			maintenanceFractionPpm:    lib.OneMillion,
			basePositionNotional:      1,
			fundingRateClampFactorPpm: 1,
			premiumVoteClampFactorPpm: 0,
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
					},
				},
				Params: types.Params{
					FundingRateClampFactorPpm: tc.fundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: tc.premiumVoteClampFactorPpm,
				},
				Perpetuals: []types.Perpetual{
					{
						MarketId:      tc.marketId,
						Ticker:        tc.ticker,
						FundingIndex:  dtypes.ZeroInt(),
						LiquidityTier: 0,
					},
				},
			}

			require.Panics(t, func() {
				perpetuals.InitGenesis(ctx, *k, genesisState)
			})
		})
	}
}
