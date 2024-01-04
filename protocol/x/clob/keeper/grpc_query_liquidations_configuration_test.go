package keeper_test

import (
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLiquidationsConfiguration(t *testing.T) {
	tests := map[string]struct {
		req *types.QueryLiquidationsConfigurationRequest
		res *types.QueryLiquidationsConfigurationResponse
		err error
	}{
		"success": {
			req: &types.QueryLiquidationsConfigurationRequest{},
			res: &types.QueryLiquidationsConfigurationResponse{
				LiquidationsConfig: types.LiquidationsConfig_Default,
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
			res, err := tApp.App.ClobKeeper.LiquidationsConfiguration(ctx, tc.req)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
