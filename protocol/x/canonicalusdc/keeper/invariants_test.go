package keeper

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestBackingInvariantChecksSupplyOnlyInGradual(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	bank.set(sdk.AccAddress("holder"), types.UusdcDenom, math.NewInt(10))
	ledger := canonicaltypes.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(9)
	require.NoError(t, k.SetLedger(ctx, ledger))

	require.NoError(t, k.SetControls(ctx, activeControls()))
	reason, broken := BackingInvariant(*k)(ctx)
	require.True(t, broken)
	require.Contains(t, reason, "does not equal accounted backing")

	paused := activeControls()
	paused.Mode = canonicaltypes.Mode_MODE_PAUSED
	require.NoError(t, k.SetControls(ctx, paused))
	reason, broken = BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)
}

func TestPendingCountMatchesStoredSettlements(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route: canonicaltypes.Route_ROUTE_NOBLE, SourceChannel: nobleChannel, Sequence: 1,
		Sender: sdk.AccAddress("s").String(), LogicalReceiver: "r", PhysicalReceiver: "r",
		LogicalDenom: types.UusdcDenom, PhysicalDenom: types.UusdcDenom, Amount: dtypes.NewInt(1),
	}))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route: canonicaltypes.Route_ROUTE_INJECTIVE, SourceChannel: injChannel, Sequence: 1,
		Sender: sdk.AccAddress("s").String(), LogicalReceiver: "r", PhysicalReceiver: "r",
		LogicalDenom: types.UusdcDenom, PhysicalDenom: physicalDenom, Amount: dtypes.NewInt(2),
	}))
	require.Equal(t, uint32(2), k.PendingCount(ctx))
	require.Len(t, k.GetAllPending(ctx), 2)

	reason, broken := BackingInvariant(*k)(ctx)
	require.True(t, broken)
	require.Contains(t, reason, "pending settlement totals do not match the backing ledger")
}
