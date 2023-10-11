package prices_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestExportGenesis(t *testing.T) {
	tests := map[string]struct {
		genesisState *types.GenesisState
	}{
		"default genesis": {
			genesisState: types.DefaultGenesis(),
		},
		"empty genesis": {
			genesisState: &types.GenesisState{
				MarketParams: []types.MarketParam{},
				MarketPrices: []types.MarketPrice{},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, k, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			prices.InitGenesis(ctx, *k, *tc.genesisState)

			// Verify expected keeper state after InitGenesis.
			require.Equal(t, tc.genesisState.MarketParams, k.GetAllMarketParams(ctx))
			require.Equal(t, tc.genesisState.MarketPrices, k.GetAllMarketPrices(ctx))

			// Verify expected exported genesis state matches input.
			exportedState := prices.ExportGenesis(ctx, *k)
			require.Equal(t, tc.genesisState, exportedState)
		})
	}
}

// invalidGenesis returns a genesis state that doesn't pass validation.
func invalidGenesis() types.GenesisState {
	genesisState := constants.Prices_DefaultGenesisState
	// Create a genesis state with a market price that doesn't match ids the market param in the same list position.
	genesisState.MarketPrices[0].Id = genesisState.MarketParams[0].Id + 1
	return genesisState
}

func TestInitGenesis_Panics(t *testing.T) {
	ctx, k, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	// Verify InitGenesis panics when given an invalid genesis state.
	require.Panics(t, func() {
		prices.InitGenesis(ctx, *k, invalidGenesis())
	})
}

func TestInitGenesisEmitsMarketUpdates(t *testing.T) {
	ctx, k, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	prices.InitGenesis(ctx, *k, constants.Prices_DefaultGenesisState)

	// Verify expected market update events.
	for _, marketPrice := range constants.Prices_DefaultGenesisState.MarketPrices {
		keepertest.AssertMarketPriceUpdateEventInIndexerBlock(t, k, ctx, marketPrice)
	}
}
