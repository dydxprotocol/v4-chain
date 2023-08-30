package prices_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := constants.Prices_DefaultGenesisState

	ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)
	prices.InitGenesis(ctx, *k, genesisState)
	got := prices.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.MarketParams, got.MarketParams)
}

func TestInitGenesisEmitsMarketUpdates(t *testing.T) {
	ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)

	prices.InitGenesis(ctx, *k, constants.Prices_DefaultGenesisState)

	// Verify expected market update events.
	for _, marketPrice := range constants.Prices_DefaultGenesisState.MarketPrices {
		keepertest.AssertMarketPriceUpdateEventInIndexerBlock(t, k, ctx, marketPrice)
	}
}
