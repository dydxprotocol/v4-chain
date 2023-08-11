package prices_test

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/prices"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := constants.Prices_DefaultGenesisState

	ctx, k, _, _, _ := keepertest.PricesKeepers(t)
	prices.InitGenesis(ctx, *k, genesisState)
	got := prices.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.Markets, got.Markets)
	require.ElementsMatch(t, genesisState.ExchangeFeeds, got.ExchangeFeeds)
}
