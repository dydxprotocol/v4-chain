package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetLastTradePrices(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Get non-existent last trade price.
	minTradePriceSubticks, maxTradePriceSubticks, found := ks.ClobKeeper.GetTradePricesForPerpetual(ks.Ctx, 0)
	require.Equal(t, minTradePriceSubticks, types.Subticks(0))
	require.Equal(t, maxTradePriceSubticks, types.Subticks(0))
	require.False(t, found)

	// Set trade prices.
	ks.ClobKeeper.SetTradePricesForPerpetual(ks.Ctx, 0, types.Subticks(17))

	// Get the min and max trade prices, which should now exist.
	minTradePriceSubticks, maxTradePriceSubticks, found = ks.ClobKeeper.GetTradePricesForPerpetual(ks.Ctx, 0)
	require.Equal(t, minTradePriceSubticks, types.Subticks(17))
	require.Equal(t, maxTradePriceSubticks, types.Subticks(17))
	require.True(t, found)

	// Update the min price.
	ks.ClobKeeper.SetTradePricesForPerpetual(ks.Ctx, 0, types.Subticks(13))

	minTradePriceSubticks, maxTradePriceSubticks, found = ks.ClobKeeper.GetTradePricesForPerpetual(ks.Ctx, 0)
	require.Equal(t, minTradePriceSubticks, types.Subticks(13))
	require.Equal(t, maxTradePriceSubticks, types.Subticks(17))
	require.True(t, found)

	// Update the max price.
	ks.ClobKeeper.SetTradePricesForPerpetual(ks.Ctx, 0, types.Subticks(23))

	minTradePriceSubticks, maxTradePriceSubticks, found = ks.ClobKeeper.GetTradePricesForPerpetual(ks.Ctx, 0)
	require.Equal(t, minTradePriceSubticks, types.Subticks(13))
	require.Equal(t, maxTradePriceSubticks, types.Subticks(23))
	require.True(t, found)
}
