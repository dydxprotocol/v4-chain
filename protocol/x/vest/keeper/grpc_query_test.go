package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVestEntryQuery(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	defaultEntry := types.DefaultGenesis().VestEntries[0]

	for name, tc := range map[string]struct {
		req *types.QueryVestEntryRequest
		err error
	}{
		"Success - default": {
			req: &types.QueryVestEntryRequest{
				VesterAccount: defaultEntry.VesterAccount,
			},
			err: nil,
		},
		"Failure - non-existent": {
			req: &types.QueryVestEntryRequest{
				VesterAccount: "non-existent",
			},
			err: types.ErrVestEntryNotFound,
		},
		"Nil": {
			req: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.VestEntry(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, defaultEntry, res.Entry)
			}
		})
	}
}
