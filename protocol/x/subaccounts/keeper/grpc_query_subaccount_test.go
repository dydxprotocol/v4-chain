package keeper_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestSubaccountQuerySingle(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
	msgs := createNSubaccount(keeper, ctx, 2, big.NewInt(1_000))
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetSubaccountRequest
		response *types.QuerySubaccountResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetSubaccountRequest{
				Owner:  msgs[0].Id.Owner,
				Number: msgs[0].Id.Number,
			},
			response: &types.QuerySubaccountResponse{Subaccount: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetSubaccountRequest{
				Owner:  msgs[1].Id.Owner,
				Number: msgs[1].Id.Number,
			},
			response: &types.QuerySubaccountResponse{Subaccount: msgs[1]},
		},
		{
			desc: "KeyNotFoundOwner",
			request: &types.QueryGetSubaccountRequest{
				Owner:  "100000",
				Number: msgs[1].Id.Number,
			},
			response: &types.QuerySubaccountResponse{Subaccount: types.Subaccount{
				Id: &types.SubaccountId{
					Owner:  "100000",
					Number: msgs[1].Id.Number,
				},
			}},
		},
		{
			desc: "KeyNotFoundNumber",
			request: &types.QueryGetSubaccountRequest{
				Owner:  msgs[1].Id.Owner,
				Number: uint32(100),
			},
			response: &types.QuerySubaccountResponse{Subaccount: types.Subaccount{
				Id: &types.SubaccountId{
					Owner:  msgs[1].Id.Owner,
					Number: uint32(100),
				},
			}},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Subaccount(ctx, tc.request)
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

func TestSubaccountQueryPaginated(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
	msgs := createNSubaccount(keeper, ctx, 5, big.NewInt(1_000))

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSubaccountRequest {
		return &types.QueryAllSubaccountRequest{
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
			resp, err := keeper.SubaccountAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(msgs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SubaccountAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(msgs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.SubaccountAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),            //nolint:staticcheck
			nullify.Fill(resp.Subaccount), //nolint:staticcheck
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.SubaccountAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
