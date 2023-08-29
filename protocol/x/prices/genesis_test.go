package prices_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesis(t *testing.T) {
	genesisState := constants.Prices_DefaultGenesisState

	ctx, k, _, _, _, _ := keepertest.PricesKeepers(t)
	prices.InitGenesis(ctx, *k, genesisState)
	got := prices.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.MarketParams, got.MarketParams)
}
