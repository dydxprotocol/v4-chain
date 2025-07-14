package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestQueryOrderRouterRevShare(t *testing.T) {
	testCases := map[string]struct {
		config      types.OrderRouterRevShare
		expectedPpm uint32
		expectedErr error
	}{
		"Single recipient": {
			config: types.OrderRouterRevShare{
				Address: constants.AliceAccAddress.String(),
			},
			expectedPpm: 500_000,
		},
		"Invalid address": {
			config: types.OrderRouterRevShare{
				Address: "invalid_address",
			},
			expectedErr: types.ErrInvalidAddress,
		},
		"Not found": {
			config: types.OrderRouterRevShare{
				Address: constants.BobAccAddress.String(),
			},
			expectedErr: types.ErrOrderRouterRevShareNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RevShareKeeper

			setErr := k.SetOrderRouterRevShare(ctx, constants.AliceAccAddress.String(), 500_000)
			require.NoError(t, setErr)

			resp, err := k.OrderRouterRevShare(ctx, &types.QueryOrderRouterRevShare{
				Address: tc.config.Address,
			})
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPpm, resp.OrderRouterRevShare.SharePpm)
			}
		})
	}
}
