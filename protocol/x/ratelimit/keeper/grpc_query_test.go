package keeper_test

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	bank_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	ratelimitutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	cometbfttypes "github.com/cometbft/cometbft/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListLimiterParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	for name, tc := range map[string]struct {
		req *types.ListLimitParamsRequest
		res *types.ListLimitParamsResponse
		err error
	}{
		"Success": {
			req: &types.ListLimitParamsRequest{},
			res: &types.ListLimitParamsResponse{
				LimitParamsList: types.DefaultGenesis().LimitParamsList,
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
			res, err := k.ListLimitParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestCapacityByDenom(t *testing.T) {
	for name, tc := range map[string]struct {
		req *types.QueryCapacityByDenomRequest
		res *types.QueryCapacityByDenomResponse
		err error
	}{
		"Success, returns default limiter and baseline capacity": {
			req: &types.QueryCapacityByDenomRequest{
				Denom: types.SDaiDenom,
			},
			res: &types.QueryCapacityByDenomResponse{
				LimiterCapacityList: []types.LimiterCapacity{
					{
						Limiter: types.DefaultSDaiHourlyLimter,
						Capacity: dtypes.NewIntFromBigInt(
							ratelimitutil.GetBaseline(
								big.NewInt(0),
								types.DefaultSDaiHourlyLimter,
							),
						),
					},
					{
						Limiter: types.DefaultSDaiDailyLimiter,
						Capacity: dtypes.NewIntFromBigInt(
							ratelimitutil.GetBaseline(
								big.NewInt(0),
								types.DefaultSDaiDailyLimiter,
							),
						),
					},
				},
			},
			err: nil,
		},
		"Success, non-existing denom": {
			req: &types.QueryCapacityByDenomRequest{
				Denom: "foo",
			},
			res: &types.QueryCapacityByDenomResponse{
				LimiterCapacityList: []types.LimiterCapacity{},
			},
			err: nil,
		},
		"Error: invalid denom": {
			req: &types.QueryCapacityByDenomRequest{
				Denom: "@@@???",
			},
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid denom: @@@???"),
		},
		"Error: nil request": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up treasury account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						// Remove any SDai balance from genesis.
						// Without additional SDai balance, this means all capacities are
						// initialized with minimum baseline.
						genesisState.Balances = bank_testutil.FilterDenomFromBalances(
							genesisState.Balances,
							types.SDaiDenom,
						)
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			res, err := k.CapacityByDenom(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestGetAllPendingSendPacket(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	channels := []string{"channel-0", "channel-1"}
	sequences := []uint64{20, 22}

	for i := range channels {
		k.SetPendingSendPacket(ctx, channels[i], sequences[i])
	}

	req := &types.QueryAllPendingSendPacketsRequest{}
	res, err := k.AllPendingSendPackets(ctx, req)
	require.NoError(t, err)
	require.Equal(t, &types.QueryAllPendingSendPacketsResponse{
		PendingSendPackets: []types.PendingSendPacket{
			{
				ChannelId: channels[0],
				Sequence:  sequences[0],
			},
			{
				ChannelId: channels[1],
				Sequence:  sequences[1],
			},
		},
	}, res)
}

func TestGetSDAIPriceQuery(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	for name, tc := range map[string]struct {
		req *types.GetSDAIPriceQueryRequest
		res *types.GetSDAIPriceQueryResponse
		err error
	}{
		"Success": {
			req: &types.GetSDAIPriceQueryRequest{},
			res: &types.GetSDAIPriceQueryResponse{
				Price: "1",
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
			k.SetSDAIPrice(ctx, big.NewInt(1))
			res, err := k.GetSDAIPriceQuery(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestGetAssetYieldIndexQuery(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	for name, tc := range map[string]struct {
		req *types.GetAssetYieldIndexQueryRequest
		res *types.GetAssetYieldIndexQueryResponse
		err error
	}{
		"Success": {
			req: &types.GetAssetYieldIndexQueryRequest{},
			res: &types.GetAssetYieldIndexQueryResponse{
				AssetYieldIndex: "1/1",
			},
			err: nil,
		},
		"Invalid Request": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"AssetYieldIndex not found": {
			req: &types.GetAssetYieldIndexQueryRequest{},
			res: nil,
			err: status.Error(codes.NotFound, "asset yield index not found"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			if !(tc.req != nil && tc.res == nil) {
				k.SetAssetYieldIndex(ctx, keeper.ConvertStringToBigRatWithPanicOnErr("1/1"))
			} else {
				store := ctx.KVStore(k.GetStoreKeyForTestingOnly())
				store.Delete([]byte(types.AssetYieldIndexPrefix))
			}
			res, err := k.GetAssetYieldIndexQuery(ctx, tc.req)
			if tc.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
