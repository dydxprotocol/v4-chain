package keeper_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func TestPerpetualQuerySingle(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	msgs := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryPerpetualRequest
		response *types.QueryPerpetualResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryPerpetualRequest{
				Id: msgs[0].Params.Id,
			},
			response: &types.QueryPerpetualResponse{Perpetual: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryPerpetualRequest{
				Id: msgs[1].Params.Id,
			},
			response: &types.QueryPerpetualResponse{Perpetual: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryPerpetualRequest{
				Id: uint32(100000),
			},
			err: status.Error(codes.NotFound, fmt.Sprintf(
				"Perpetual id %+v not found.",
				uint32(100000),
			)),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := pc.PerpetualsKeeper.Perpetual(pc.Ctx, tc.request)
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

func TestPerpetualQueryPaginated(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	msgs := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPerpetualsRequest {
		return &types.QueryAllPerpetualsRequest{
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
		for i := 0; i < len(msgs); i += step {
			resp, err := pc.PerpetualsKeeper.AllPerpetuals(pc.Ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Perpetual), step)
			require.Subset(t,
				nullify.Fill(msgs),           //nolint:staticcheck
				nullify.Fill(resp.Perpetual), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := pc.PerpetualsKeeper.AllPerpetuals(pc.Ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Perpetual), step)
			require.Subset(t,
				nullify.Fill(msgs),           //nolint:staticcheck
				nullify.Fill(resp.Perpetual), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := pc.PerpetualsKeeper.AllPerpetuals(pc.Ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),           //nolint:staticcheck
			nullify.Fill(resp.Perpetual), //nolint:staticcheck
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := pc.PerpetualsKeeper.AllPerpetuals(pc.Ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
