package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestParams(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	wctx := sdk.WrapSDKContext(pc.Ctx)

	tests := map[string]struct {
		req         *types.QueryParamsRequest
		res         *types.QueryParamsResponse
		expectedErr error
	}{
		"nil request": {
			req:         nil,
			res:         nil,
			expectedErr: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"valid request": {
			req: &types.QueryParamsRequest{},
			res: &types.QueryParamsResponse{
				Params: types.DefaultGenesis().Params,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := pc.PerpetualsKeeper.Params(wctx, tc.req)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.Equal(t, tc.res, res)
			}
		})
	}
}
