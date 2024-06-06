package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

func TestQueryPMLEnabledStatus(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.ListingKeeper

	// set permissionless listing to true for test
	k.SetPermissionlessListingEnable(ctx, true)

	// query permissionless market listing status
	resp, err := k.PermissionlessMarketListingStatus(ctx, &types.QueryPermissionlessMarketListingStatus{})
	require.NoError(t, err)
	require.True(t, resp.Enabled)
}
