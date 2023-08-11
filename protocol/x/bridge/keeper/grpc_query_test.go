package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

func TestEventParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryEventParamsRequest
		res *types.QueryEventParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryEventParamsRequest{},
			res: &types.QueryEventParamsResponse{
				Params: types.DefaultGenesis().EventParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.EventParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestProposeParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryProposeParamsRequest
		res *types.QueryProposeParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryProposeParamsRequest{},
			res: &types.QueryProposeParamsResponse{
				Params: types.DefaultGenesis().ProposeParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.ProposeParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestSafetyParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QuerySafetyParamsRequest
		res *types.QuerySafetyParamsResponse
		err error
	}{
		"Success": {
			req: &types.QuerySafetyParamsRequest{},
			res: &types.QuerySafetyParamsResponse{
				Params: types.DefaultGenesis().SafetyParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.SafetyParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
