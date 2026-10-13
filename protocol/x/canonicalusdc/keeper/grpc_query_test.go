package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func noblePending(sequence uint64, amount int64) canonicaltypes.PendingSettlement {
	return canonicaltypes.PendingSettlement{
		Route:            canonicaltypes.Route_ROUTE_NOBLE,
		SourceChannel:    nobleChannel,
		Sequence:         sequence,
		Sender:           sdk.AccAddress("sender").String(),
		LogicalReceiver:  "noble1receiver",
		PhysicalReceiver: "noble1receiver",
		LogicalDenom:     types.UusdcDenom,
		PhysicalDenom:    types.UusdcDenom,
		Amount:           dtypes.NewInt(amount),
	}
}

func TestQueryStateAndPending(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	ledger := canonicaltypes.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(3)
	require.NoError(t, k.SetLedger(ctx, ledger))
	require.NoError(t, k.SetPending(ctx, noblePending(4, 3)))

	state, err := k.State(ctx, &canonicaltypes.QueryStateRequest{})
	require.NoError(t, err)
	require.Equal(t, canonicaltypes.Mode_MODE_GRADUAL, state.Controls.Mode)
	require.Equal(t, "3", state.Ledger.NobleBacking.String())

	pending, err := k.PendingSettlements(ctx, &canonicaltypes.QueryPendingSettlementsRequest{})
	require.NoError(t, err)
	require.Len(t, pending.PendingSettlements, 1)
	require.Equal(t, uint64(4), pending.PendingSettlements[0].Sequence)

	single, err := k.PendingSettlement(ctx, &canonicaltypes.QueryPendingSettlementRequest{
		SourceChannel: nobleChannel,
		Sequence:      4,
	})
	require.NoError(t, err)
	require.Equal(t, noblePending(4, 3), single.PendingSettlement)

	_, err = k.PendingSettlement(ctx, &canonicaltypes.QueryPendingSettlementRequest{
		SourceChannel: nobleChannel,
		Sequence:      5,
	})
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestQueryPendingSettlementsPaginates(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	for sequence := uint64(1); sequence <= 5; sequence++ {
		require.NoError(t, k.SetPending(ctx, noblePending(sequence, 1)))
	}

	seen := make([]uint64, 0, 5)
	var nextKey []byte
	for page := 0; page < 3; page++ {
		response, err := k.PendingSettlements(ctx, &canonicaltypes.QueryPendingSettlementsRequest{
			Pagination: &query.PageRequest{Key: nextKey, Limit: 2, CountTotal: page == 0},
		})
		require.NoError(t, err)
		if page == 0 {
			require.Equal(t, uint64(5), response.Pagination.Total)
		}
		for _, settlement := range response.PendingSettlements {
			seen = append(seen, settlement.Sequence)
		}
		nextKey = response.Pagination.NextKey
	}
	require.Equal(t, []uint64{1, 2, 3, 4, 5}, seen)
	require.Nil(t, nextKey)
}

func TestQueryNilRequests(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	_, err := k.State(ctx, nil)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)
	_, err = k.PendingSettlements(ctx, nil)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)
	_, err = k.PendingSettlement(ctx, nil)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)
}
