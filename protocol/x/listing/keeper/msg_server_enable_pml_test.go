package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"
)

func TestMsgEnablePermissionlessMarketListing(t *testing.T) {

	tests := map[string]struct {
		// Msg.
		msg *types.MsgEnablePermissionlessMarketListing
		// Expected error
		expectedErr string
	}{
		"Success - Enabled": {
			msg: &types.MsgEnablePermissionlessMarketListing{
				Authority:                         lib.GovModuleAddress.String(),
				EnablePermissionlessMarketListing: true,
			},
		},
		"Success - Disabled": {
			msg: &types.MsgEnablePermissionlessMarketListing{
				Authority:                         lib.GovModuleAddress.String(),
				EnablePermissionlessMarketListing: false,
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgEnablePermissionlessMarketListing{
				Authority:                         constants.AliceAccAddress.String(),
				EnablePermissionlessMarketListing: true,
			},
			expectedErr: "invalid authority",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.ListingKeeper
			ms := keeper.NewMsgServerImpl(k)
			_, err := ms.EnablePermissionlessMarketListing(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				enabledFlag, err := k.IsPermissionlessListingEnabled(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.msg.EnablePermissionlessMarketListing, enabledFlag)
			}
		})
	}
}
