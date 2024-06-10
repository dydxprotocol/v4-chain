package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

func TestQueryPMLEnabledStatus(t *testing.T) {
	tests := map[string]struct {
		pmlEnabled bool
	}{
		"PML enabled true": {
			pmlEnabled: true,
		},
		"PML enabled false": {
			pmlEnabled: false,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.ListingKeeper

				// set permissionless listing to true for test
				err := k.SetPermissionlessListingEnable(ctx, tc.pmlEnabled)
				require.NoError(t, err)

				// query permissionless market listing status
				resp, err := k.PermissionlessMarketListingStatus(ctx, &types.QueryPermissionlessMarketListingStatus{})
				require.NoError(t, err)
				require.Equal(t, resp.Enabled, tc.pmlEnabled)
			},
		)
	}
}
