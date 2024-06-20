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
	err := k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)
	require.NoError(t, err)

	resp, err := k.MarketMapperRevShareDetails(
		ctx, &types.QueryMarketMapperRevShareDetails{
			MarketId: marketId,
		},
	)
	require.NoError(t, err)
	require.Equal(t, resp.Details.ExpirationTs, setDetails.ExpirationTs)
}
