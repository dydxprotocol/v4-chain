package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestQueryMarketMapperRevenueParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RevShareKeeper

	params := types.MarketMapperRevenueShareParams{
		Address:         constants.AliceAccAddress.String(),
		RevenueSharePpm: 100_000,
		ValidDays:       100,
	}

	err := k.SetMarketMapperRevenueShareParams(ctx, params)
	require.NoError(t, err)

	resp, err := k.MarketMapperRevenueShareParams(ctx, &types.QueryMarketMapperRevenueShareParams{})
	require.NoError(t, err)
	require.Equal(t, resp.Params, params)
}
