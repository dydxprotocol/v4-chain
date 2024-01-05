package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestClobPairQuerySingle(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	msgs := keepertest.CreateNClobPair(t,
		ks.ClobKeeper,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		ks.Ctx,
		2,
		mockIndexerEventManager,
	)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetClobPairRequest
		response *types.QueryClobPairResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetClobPairRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryClobPairResponse{ClobPair: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetClobPairRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryClobPairResponse{ClobPair: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetClobPairRequest{
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
			response, err := ks.ClobKeeper.ClobPair(ks.Ctx, tc.request)
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

func TestClobPairQueryPaginated(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	msgs := keepertest.CreateNClobPair(t,
		ks.ClobKeeper,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		ks.Ctx,
		10,
		mockIndexerEventManager,
	)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllClobPairRequest {
		return &types.QueryAllClobPairRequest{
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
			resp, err := ks.ClobKeeper.ClobPairAll(ks.Ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(msgs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := ks.ClobKeeper.ClobPairAll(ks.Ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(msgs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := ks.ClobKeeper.ClobPairAll(ks.Ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),          //nolint:staticcheck
			nullify.Fill(resp.ClobPair), //nolint:staticcheck
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := ks.ClobKeeper.ClobPairAll(ks.Ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
