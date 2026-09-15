package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestCompletedRingEvictsOldestAndAllowsSequenceReuse(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	const channel = nobleChannel
	for i := uint64(1); i <= canonicaltypes.HardMaxCompleted; i++ {
		k.markCompleted(ctx, channel, i)
	}
	require.True(t, k.HasCompleted(ctx, channel, 1))
	k.markCompleted(ctx, channel, canonicaltypes.HardMaxCompleted+1)
	require.False(t, k.HasCompleted(ctx, channel, 1))
	require.True(t, k.HasCompleted(ctx, channel, canonicaltypes.HardMaxCompleted+1))

	require.NoError(t, k.SetControls(ctx, activeControls()))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route:            canonicaltypes.Route_ROUTE_NOBLE,
		SourceChannel:    channel,
		Sequence:         1,
		Sender:           sdk.AccAddress("timeout-sender").String(),
		LogicalReceiver:  "noble1receiver",
		PhysicalReceiver: "noble1receiver",
		LogicalDenom:     types.UusdcDenom,
		PhysicalDenom:    types.UusdcDenom,
		Amount:           "1",
	}))
	pending, found := k.GetPendingForPacket(ctx, channel, 1)
	require.True(t, found, "evicted completion must not hide a new pending settlement")
	require.Equal(t, uint64(1), pending.Sequence)
}

func TestPendingCountMatchesStoredSettlements(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route: canonicaltypes.Route_ROUTE_NOBLE, SourceChannel: nobleChannel, Sequence: 1,
		Sender: sdk.AccAddress("s").String(), LogicalReceiver: "r", PhysicalReceiver: "r",
		LogicalDenom: types.UusdcDenom, PhysicalDenom: types.UusdcDenom, Amount: "1",
	}))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route: canonicaltypes.Route_ROUTE_INJECTIVE, SourceChannel: injChannel, Sequence: 1,
		Sender: sdk.AccAddress("s").String(), LogicalReceiver: "r", PhysicalReceiver: "r",
		LogicalDenom: types.UusdcDenom, PhysicalDenom: physicalDenom, Amount: "2",
	}))
	require.Equal(t, uint32(2), k.PendingCount(ctx))
	require.Len(t, k.GetAllPending(ctx), 2)

	reason, broken := BackingInvariant(*k)(ctx)
	require.True(t, broken)
	require.Contains(t, reason, "pending settlement totals do not match the backing ledger")
}
