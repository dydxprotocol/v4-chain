package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestQueryMarketMapperRevShareDetails(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper

	setDetails := types.MarketMapperRevShareDetails{
		ExpirationTs: 1735707600,
	}
	marketId := uint32(42)
	k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)

	resp, err := k.MarketMapperRevShareDetails(
		ctx, &types.QueryMarketMapperRevShareDetails{
			MarketId: marketId,
		},
	)
	require.NoError(t, err)
	require.Equal(t, resp.Details.ExpirationTs, setDetails.ExpirationTs)
}

func TestQueryMarketMapperRevShareDetailsFailure(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RevShareKeeper

	// Query for revshare details of non-existent market
	_, err := k.MarketMapperRevShareDetails(
		ctx, &types.QueryMarketMapperRevShareDetails{
			MarketId: 42,
		},
	)
	require.ErrorIs(t, err, types.ErrMarketMapperRevShareDetailsNotFound)
}
