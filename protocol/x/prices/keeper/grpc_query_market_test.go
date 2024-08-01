package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestMarketPriceQuerySingle(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	msgs := keepertest.CreateNMarkets(t, ctx, keeper, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryMarketPriceRequest
		response *types.QueryMarketPriceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryMarketPriceRequest{
				Id: msgs[0].Param.Id,
			},
			response: &types.QueryMarketPriceResponse{MarketPrice: msgs[0].Price},
		},
		{
			desc: "Second",
			request: &types.QueryMarketPriceRequest{
				Id: msgs[1].Param.Id,
			},
			response: &types.QueryMarketPriceResponse{MarketPrice: msgs[1].Price},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryMarketPriceRequest{
				Id: uint32(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.MarketPrice(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response), //nolint:staticcheck
					nullify.Fill(response),    //nolint:staticcheck
				)
			}
		})
	}
}

func TestMarketPriceQueryPaginated(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	msgs := keepertest.CreateNMarkets(t, ctx, keeper, 5)
	prices := make([]types.MarketPrice, len(msgs))
	for i := range msgs {
		prices[i] = msgs[i].Price
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMarketPricesRequest {
		return &types.QueryAllMarketPricesRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(prices); i += step {
			resp, err := keeper.AllMarketPrices(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MarketPrices), step)
			require.Subset(t,
				nullify.Fill(prices),            //nolint:staticcheck
				nullify.Fill(resp.MarketPrices), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(prices); i += step {
			resp, err := keeper.AllMarketPrices(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MarketPrices), step)
			require.Subset(t,
				nullify.Fill(prices),            //nolint:staticcheck
				nullify.Fill(resp.MarketPrices), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AllMarketPrices(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(prices), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(prices),            //nolint:staticcheck
			nullify.Fill(resp.MarketPrices), //nolint:staticcheck
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AllMarketPrices(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestMarketParamQuerySingle(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	msgs := keepertest.CreateNMarkets(t, ctx, keeper, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryMarketParamRequest
		response *types.QueryMarketParamResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryMarketParamRequest{
				Id: msgs[0].Param.Id,
			},
			response: &types.QueryMarketParamResponse{MarketParam: msgs[0].Param},
		},
		{
			desc: "Second",
			request: &types.QueryMarketParamRequest{
				Id: msgs[1].Param.Id,
			},
			response: &types.QueryMarketParamResponse{MarketParam: msgs[1].Param},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryMarketParamRequest{
				Id: uint32(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.MarketParam(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response), //nolint:staticcheck
					nullify.Fill(response),    //nolint:staticcheck
				)
			}
		})
	}
}

func TestMarketParamQueryPaginated(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	msgs := keepertest.CreateNMarkets(t, ctx, keeper, 5)
	params := make([]types.MarketParam, len(msgs))
	for i := range msgs {
		params[i] = msgs[i].Param
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMarketParamsRequest {
		return &types.QueryAllMarketParamsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(params); i += step {
			resp, err := keeper.AllMarketParams(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MarketParams), step)
			require.Subset(t,
				nullify.Fill(params),            //nolint:staticcheck
				nullify.Fill(resp.MarketParams), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(params); i += step {
			resp, err := keeper.AllMarketParams(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MarketParams), step)
			require.Subset(t,
				nullify.Fill(params),            //nolint:staticcheck
				nullify.Fill(resp.MarketParams), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AllMarketParams(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(params), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(params),            //nolint:staticcheck
			nullify.Fill(resp.MarketParams), //nolint:staticcheck
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AllMarketParams(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
