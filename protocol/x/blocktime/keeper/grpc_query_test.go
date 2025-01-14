package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

func TestSynchronyParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BlockTimeKeeper

	for name, tc := range map[string]struct {
		req *types.QuerySynchronyParamsRequest
		res *types.QuerySynchronyParamsResponse
		err error
	}{
		"Default": {
			req: &types.QuerySynchronyParamsRequest{},
			res: &types.QuerySynchronyParamsResponse{
				Params: types.DefaultSynchronyParams(),
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
			res, err := k.SynchronyParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestDowntimeParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BlockTimeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryDowntimeParamsRequest
		res *types.QueryDowntimeParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryDowntimeParamsRequest{},
			res: &types.QueryDowntimeParamsResponse{
				Params: types.DefaultGenesis().Params,
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
			res, err := k.DowntimeParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestAllDowntimeInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BlockTimeKeeper
	info := &types.AllDowntimeInfo{
		Infos: []*types.AllDowntimeInfo_DowntimeInfo{
			{
				Duration: time.Second,
				BlockInfo: types.BlockInfo{
					Height:    1,
					Timestamp: time.Now().UTC(),
				},
			},
		},
	}
	k.SetAllDowntimeInfo(ctx, info)

	for name, tc := range map[string]struct {
		req *types.QueryAllDowntimeInfoRequest
		res *types.QueryAllDowntimeInfoResponse
		err error
	}{
		"Success": {
			req: &types.QueryAllDowntimeInfoRequest{},
			res: &types.QueryAllDowntimeInfoResponse{
				Info: info,
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
			res, err := k.AllDowntimeInfo(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestPreviousBlockInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BlockTimeKeeper
	info := &types.BlockInfo{
		Height:    1,
		Timestamp: time.Now().UTC(),
	}
	k.SetPreviousBlockInfo(ctx, info)

	for name, tc := range map[string]struct {
		req *types.QueryPreviousBlockInfoRequest
		res *types.QueryPreviousBlockInfoResponse
		err error
	}{
		"Success": {
			req: &types.QueryPreviousBlockInfoRequest{},
			res: &types.QueryPreviousBlockInfoResponse{
				Info: info,
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
			res, err := k.PreviousBlockInfo(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
