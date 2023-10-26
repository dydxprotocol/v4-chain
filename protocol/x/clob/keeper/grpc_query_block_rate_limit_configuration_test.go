package keeper_test

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGetBlockRateLimitConfiguration(t *testing.T) {
	tests := map[string]struct {
		req *types.QueryBlockRateLimitConfigurationRequest
		res *types.QueryBlockRateLimitConfigurationResponse
		err error
	}{
		"success": {
			req: &types.QueryBlockRateLimitConfigurationRequest{},
			res: &types.QueryBlockRateLimitConfigurationResponse{
				BlockRateLimitConfig: types.BlockRateLimitConfiguration{
					MaxShortTermOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     200,
						},
					},
					MaxStatefulOrdersPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     02,
						},
						{
							NumBlocks: 100,
							Limit:     20,
						},
					},
					MaxShortTermOrderCancellationsPerNBlocks: []types.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 1,
							Limit:     200,
						},
					},
				},
			},
		},
		"failure: nil request": {
			req: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
			ctx := tApp.InitChain()
			res, err := tApp.App.ClobKeeper.BlockRateLimitConfiguration(sdktypes.WrapSDKContext(ctx), tc.req)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
