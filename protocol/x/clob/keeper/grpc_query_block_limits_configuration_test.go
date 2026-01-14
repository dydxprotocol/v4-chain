package keeper_test

import (
	"testing"

	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlockLimitsConfiguration(t *testing.T) {
	tests := map[string]struct {
		req *types.QueryBlockLimitsConfigurationRequest
		res *types.QueryBlockLimitsConfigurationResponse
		err error
	}{
		"success": {
			req: &types.QueryBlockLimitsConfigurationRequest{},
			res: &types.QueryBlockLimitsConfigurationResponse{
				BlockLimitsConfig: types.BlockLimitsConfig_Default,
			},
		},
		"failure: nil request": {
			req: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testApp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			res, err := tApp.App.ClobKeeper.BlockLimitsConfiguration(ctx, tc.req)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
