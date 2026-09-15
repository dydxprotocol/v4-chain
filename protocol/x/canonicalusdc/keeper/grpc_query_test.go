package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestQueryStateAndPending(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	ledger := canonicaltypes.DefaultLedger()
	ledger.NobleBacking = "3"
	require.NoError(t, k.SetLedger(ctx, ledger))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route:            canonicaltypes.Route_ROUTE_NOBLE,
		SourceChannel:    nobleChannel,
		Sequence:         4,
		Sender:           sdk.AccAddress("sender").String(),
		LogicalReceiver:  "noble1receiver",
		PhysicalReceiver: "noble1receiver",
		LogicalDenom:     types.UusdcDenom,
		PhysicalDenom:    types.UusdcDenom,
		Amount:           "3",
	}))

	state, err := k.State(ctx, &canonicaltypes.QueryStateRequest{})
	require.NoError(t, err)
	require.Equal(t, canonicaltypes.Mode_MODE_GRADUAL, state.Controls.Mode)
	require.Equal(t, "3", state.Ledger.NobleBacking)

	pending, err := k.PendingSettlements(ctx, &canonicaltypes.QueryPendingSettlementsRequest{})
	require.NoError(t, err)
	require.Len(t, pending.PendingSettlements, 1)
	require.Equal(t, uint64(4), pending.PendingSettlements[0].Sequence)
}

func TestQueryNilRequests(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	_, err := k.State(ctx, nil)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)
	_, err = k.PendingSettlements(ctx, nil)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)
}
