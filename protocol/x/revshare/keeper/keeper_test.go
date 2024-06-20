package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.VaultKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestGetSetMarketMapperRevShareDetails(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper
	// Set the rev share details for a market
	marketId := uint32(42)
	setDetails := types.MarketMapperRevShareDetails{
		ExpirationTs: 1735707600,
	}
	err := k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)
	require.NoError(t, err)

	// Get the rev share details for the market
	getDetails, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)

	// Set expiration ts to 0
	setDetails.ExpirationTs = 0
	err = k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)
	require.NoError(t, err)

	getDetails, err = k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)
}
