package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetLastTradePrice(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Get non-existent last trade price.
	price, found := ks.ClobKeeper.GetLastTradePriceForPerpetual(ks.Ctx, 0)
	require.Equal(t, price, types.Subticks(0))
	require.False(t, found)

	// Set last trade price.
	ks.ClobKeeper.SetLastTradePriceForPerpetual(ks.Ctx, 0, types.Subticks(17))

	// Get the last trade price, which should now exist.
	price, found = ks.ClobKeeper.GetLastTradePriceForPerpetual(ks.Ctx, 0)
	require.Equal(t, price, types.Subticks(17))
	require.True(t, found)
}
