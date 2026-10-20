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

func TestEnforceBackingPausesOnBreach(t *testing.T) {
	cases := map[string]struct {
		noble     int64
		injective int64
		reserve   int64
		mode      canonicaltypes.Mode
		wantMode  canonicaltypes.Mode
		reason    string
	}{
		"healthy gradual": {
			noble: 6, injective: 4, reserve: 4,
			mode: canonicaltypes.Mode_MODE_GRADUAL, wantMode: canonicaltypes.Mode_MODE_GRADUAL,
		},
		"supply drift": {
			noble: 9, mode: canonicaltypes.Mode_MODE_GRADUAL, wantMode: canonicaltypes.Mode_MODE_PAUSED,
			reason: "does not equal accounted backing",
		},
		"missing physical reserve": {
			noble: 6, injective: 4, reserve: 0, mode: canonicaltypes.Mode_MODE_GRADUAL, wantMode: canonicaltypes.Mode_MODE_PAUSED,
			reason: "physical reserve",
		},
		"paused is left alone": {
			noble: 9, mode: canonicaltypes.Mode_MODE_PAUSED, wantMode: canonicaltypes.Mode_MODE_PAUSED,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, k, bank, _, _ := setupKeeper(t)
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			bank.set(sdk.AccAddress("holder"), types.UusdcDenom, math.NewInt(10))
			bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(tc.reserve))
			ledger := canonicaltypes.DefaultLedger()
			ledger.NobleBacking = dtypes.NewInt(tc.noble)
			ledger.InjectiveBacking = dtypes.NewInt(tc.injective)
			require.NoError(t, k.SetLedger(ctx, ledger))
			controls := activeControls()
			controls.Mode = tc.mode
			require.NoError(t, k.SetControls(ctx, controls))

			k.EnforceBacking(ctx)

			require.Equal(t, tc.wantMode, k.GetControls(ctx).Mode)
			breaches := 0
			for _, event := range ctx.EventManager().Events() {
				if event.Type != canonicaltypes.EventTypeBackingBreach {
					continue
				}
				breaches++
				reason, found := event.GetAttribute(canonicaltypes.AttributeKeyReason)
				require.True(t, found)
				require.Contains(t, reason.Value, tc.reason)
			}
			if tc.reason == "" {
				require.Zero(t, breaches)
			} else {
				require.Equal(t, 1, breaches)
			}
		})
	}
}
