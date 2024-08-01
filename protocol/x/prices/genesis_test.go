package prices_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

var (
	modifiedPrice        = uint64(2500000000)
	modifiedMinExchanges = uint32(2)
	addedMarketParam     = constants.TestMarketParams[2]
	addedMarketPrice     = constants.TestMarketPrices[2]
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
			ctx, k, _, _, mockTimeProvider, _, marketMapKeeper := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			marketMapKeeper.InitGenesis(ctx, constants.MarketMap_DefaultGenesisState)
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

func TestExportGenesis_WithMutation(t *testing.T) {
	ctx, k, _, _, mockTimeProvider, _, marketMapKeeper := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketMapKeeper.InitGenesis(ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ctx, *k, *types.DefaultGenesis())

	modifiedMarketParam := types.DefaultGenesis().MarketParams[0]
	modifiedMarketParam.MinExchanges = modifiedMinExchanges

	// Add a new market.
	_, err := keepertest.CreateTestMarket(t, ctx, k, addedMarketParam, addedMarketPrice)
	require.NoError(t, err)

	// Modify a param.
	_, err = k.ModifyMarketParam(ctx, modifiedMarketParam)
	require.NoError(t, err)

	// Update a market price.
	err = k.UpdateMarketPrices(ctx, []*types.MsgUpdateMarketPrices_MarketPrice{
		{
			MarketId: 1,
			Price:    modifiedPrice,
		},
	})
	require.NoError(t, err)

	expectedExportGenesis := types.DefaultGenesis()
	expectedExportGenesis.MarketParams = append(expectedExportGenesis.MarketParams, addedMarketParam)
	expectedExportGenesis.MarketPrices = append(expectedExportGenesis.MarketPrices, addedMarketPrice)
	expectedExportGenesis.MarketParams[0].MinExchanges = modifiedMinExchanges
	expectedExportGenesis.MarketPrices[1].Price = modifiedPrice

	// Verify expected exported genesis state matches input.
	exportedState := prices.ExportGenesis(ctx, *k)
	require.Equal(t, expectedExportGenesis, exportedState)
}

// invalidGenesis returns a genesis state that doesn't pass validation.
func invalidGenesis() types.GenesisState {
	genesisState := *types.DefaultGenesis()
	// Create a genesis state with a market price that doesn't match ids the market param in the same list position.
	genesisState.MarketPrices[0].Id = genesisState.MarketParams[0].Id + 1
	return genesisState
}

func TestInitGenesis_Panics(t *testing.T) {
	ctx, k, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	// Verify InitGenesis panics when given an invalid genesis state.
	require.Panics(t, func() {
		prices.InitGenesis(ctx, *k, invalidGenesis())
	})
}

func TestInitGenesisEmitsMarketUpdates(t *testing.T) {
	ctx, k, _, _, mockTimeProvider, _, marketMapKeeper := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	marketMapKeeper.InitGenesis(ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ctx, *k, constants.Prices_DefaultGenesisState)

	// Verify expected market update events.
	for _, marketPrice := range constants.Prices_DefaultGenesisState.MarketPrices {
		keepertest.AssertMarketPriceUpdateEventInIndexerBlock(t, k, ctx, marketPrice)
	}
}
