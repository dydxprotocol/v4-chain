package keeper

import (
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestUpdateControlsRequiresAuthorityAndBindings(t *testing.T) {
	cases := []struct {
		name string
		prep func(*Keeper)
		msg  *canonicaltypes.MsgUpdateControls
		err  error
	}{
		{name: "nil message", msg: nil, err: canonicaltypes.ErrUnauthorized},
		{
			name: "wrong authority",
			msg: &canonicaltypes.MsgUpdateControls{
				Authority: sdk.AccAddress("not-authority").String(),
				Controls:  activeControls(),
			},
			err: canonicaltypes.ErrUnauthorized,
		},
		{
			name: "valid",
			prep: func(k *Keeper) { bindOpenRoutes(k, activeControls()) },
			msg: &canonicaltypes.MsgUpdateControls{
				Authority: lib.GovModuleAddress.String(),
				Controls:  activeControls(),
			},
		},
		{
			name: "closed channel",
			prep: func(k *Keeper) {
				bindOpenRoutes(k, activeControls())
				closed := k.channelKeeper.(channelMock)[injChannel]
				closed.State = channeltypes.CLOSED
				k.channelKeeper.(channelMock)[injChannel] = closed
			},
			msg: &canonicaltypes.MsgUpdateControls{
				Authority: lib.GovModuleAddress.String(),
				Controls:  activeControls(),
			},
			err: canonicaltypes.ErrInvalidControls,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, _, _, _ := setupKeeper(t)
			if tc.prep != nil {
				tc.prep(k)
			}
			_, err := NewMsgServer(k).UpdateControls(ctx, tc.msg)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, canonicaltypes.Mode_MODE_GRADUAL, k.GetControls(ctx).Mode)
			require.Equal(t, "0", k.GetLedger(ctx).NobleBacking.String())
		})
	}
}

func TestUpdateControlsPendingFreeze(t *testing.T) {
	setupPending := func(t *testing.T) (sdk.Context, *Keeper) {
		t.Helper()
		ctx, k, _, _, _ := setupKeeper(t)
		bindOpenRoutes(k, activeControls())
		require.NoError(t, k.SetControls(ctx, activeControls()))
		ledger := canonicaltypes.DefaultLedger()
		ledger.PendingNoble = dtypes.NewInt(1)
		require.NoError(t, k.SetLedger(ctx, ledger))
		require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
			Route:            canonicaltypes.Route_ROUTE_NOBLE,
			SourceChannel:    nobleChannel,
			Sequence:         1,
			Sender:           sdk.AccAddress("pending-sender").String(),
			LogicalReceiver:  "noble1receiver",
			PhysicalReceiver: "noble1receiver",
			LogicalDenom:     types.UusdcDenom,
			PhysicalDenom:    types.UusdcDenom,
			Amount:           dtypes.NewInt(1),
		}))
		return ctx, k
	}
	cases := []struct {
		name        string
		prep        func(*Keeper, *canonicaltypes.Controls)
		mode        canonicaltypes.Mode
		errContains string
		wantMode    canonicaltypes.Mode
	}{
		{
			name:        "disable",
			mode:        canonicaltypes.Mode_MODE_DISABLED,
			errContains: "cannot disable while",
		},
		{
			name: "retarget",
			prep: func(k *Keeper, controls *canonicaltypes.Controls) {
				controls.InjectiveChannel = "channel-9"
				controls.InjectiveDenom = transfertypes.ParseDenomTrace("transfer/channel-9/uusdc").IBCDenom()
				bindOpenRoutes(k, *controls)
			},
			mode:        canonicaltypes.Mode_MODE_GRADUAL,
			errContains: "route identity",
		},
		{
			name: "rebind client",
			prep: func(k *Keeper, controls *canonicaltypes.Controls) {
				controls.NobleClient = "07-tendermint-9"
				bindOpenRoutes(k, *controls)
			},
			mode:        canonicaltypes.Mode_MODE_GRADUAL,
			errContains: "route identity",
		},
		{name: "pause", mode: canonicaltypes.Mode_MODE_PAUSED, wantMode: canonicaltypes.Mode_MODE_PAUSED},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k := setupPending(t)
			controls := activeControls()
			controls.Mode = tc.mode
			if tc.prep != nil {
				tc.prep(k, &controls)
			}
			_, err := NewMsgServer(k).UpdateControls(ctx, &canonicaltypes.MsgUpdateControls{
				Authority: lib.GovModuleAddress.String(),
				Controls:  controls,
			})
			if tc.errContains != "" {
				require.ErrorIs(t, err, canonicaltypes.ErrInvalidControls)
				require.Contains(t, err.Error(), tc.errContains)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantMode, k.GetControls(ctx).Mode)
		})
	}
}

func TestUpdateControlsSnapshotsNobleBacking(t *testing.T) {
	holder := sdk.AccAddress("holder")
	cases := []struct {
		name             string
		mode             canonicaltypes.Mode
		uusdc            int64
		physical         int64
		injectiveBacking int64
		wantErr          error
		wantMode         canonicaltypes.Mode
		wantNoble        string
		wantInjective    string
		thenPause        bool
	}{
		{
			name:      "from supply",
			mode:      canonicaltypes.Mode_MODE_GRADUAL,
			uusdc:     10,
			wantMode:  canonicaltypes.Mode_MODE_GRADUAL,
			wantNoble: "10",
			thenPause: true,
		},
		{
			name:      "pausing from disabled does not snapshot",
			mode:      canonicaltypes.Mode_MODE_PAUSED,
			uusdc:     8,
			wantMode:  canonicaltypes.Mode_MODE_PAUSED,
			wantNoble: "0",
		},
		{
			name:             "residual after classified",
			mode:             canonicaltypes.Mode_MODE_GRADUAL,
			uusdc:            10,
			physical:         4,
			injectiveBacking: 4,
			wantMode:         canonicaltypes.Mode_MODE_GRADUAL,
			wantNoble:        "6",
			wantInjective:    "4",
		},
		{
			name:             "supply below classified",
			mode:             canonicaltypes.Mode_MODE_GRADUAL,
			uusdc:            3,
			injectiveBacking: 4,
			wantErr:          canonicaltypes.ErrInvalidLedger,
			wantMode:         canonicaltypes.Mode_MODE_DISABLED,
			wantNoble:        "0",
			wantInjective:    "4",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, bank, _, _ := setupKeeper(t)
			bindOpenRoutes(k, activeControls())
			bank.set(holder, types.UusdcDenom, math.NewInt(tc.uusdc))
			if tc.physical != 0 {
				bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(tc.physical))
			}
			if tc.injectiveBacking != 0 {
				ledger := canonicaltypes.DefaultLedger()
				ledger.InjectiveBacking = dtypes.NewInt(tc.injectiveBacking)
				require.NoError(t, k.SetLedger(ctx, ledger))
			}
			controls := activeControls()
			controls.Mode = tc.mode
			_, err := NewMsgServer(k).UpdateControls(ctx, &canonicaltypes.MsgUpdateControls{
				Authority: lib.GovModuleAddress.String(),
				Controls:  controls,
			})
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Equal(t, tc.wantMode, k.GetControls(ctx).Mode)
				require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking.String())
				require.Equal(t, tc.wantInjective, k.GetLedger(ctx).InjectiveBacking.String())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantMode, k.GetControls(ctx).Mode)
			require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking.String())
			if tc.wantInjective != "" {
				require.Equal(t, tc.wantInjective, k.GetLedger(ctx).InjectiveBacking.String())
			}
			requireInvariant(t, ctx, k)
			if !tc.thenPause {
				return
			}
			paused := activeControls()
			paused.Mode = canonicaltypes.Mode_MODE_PAUSED
			_, err = NewMsgServer(k).UpdateControls(ctx, &canonicaltypes.MsgUpdateControls{
				Authority: lib.GovModuleAddress.String(),
				Controls:  paused,
			})
			require.NoError(t, err)
			require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking.String())
		})
	}
}

func updateControls(t *testing.T, ctx sdk.Context, k *Keeper, controls canonicaltypes.Controls) error {
	t.Helper()
	_, err := NewMsgServer(k).UpdateControls(ctx, &canonicaltypes.MsgUpdateControls{
		Authority: lib.GovModuleAddress.String(),
		Controls:  controls,
	})
	return err
}

func withMode(mode canonicaltypes.Mode) canonicaltypes.Controls {
	controls := activeControls()
	controls.Mode = mode
	return controls
}

func TestUpdateControlsNobleWithdrawalsRequiredWhileNobleBackingRemains(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	bank.set(sdk.AccAddress("holder"), types.UusdcDenom, math.NewInt(10))

	closed := activeControls()
	closed.NobleWithdrawalsEnabled = false
	require.ErrorIs(t, updateControls(t, ctx, k, closed), canonicaltypes.ErrNobleBackingRemains)
	require.Equal(t, canonicaltypes.Mode_MODE_DISABLED, k.GetControls(ctx).Mode)

	require.NoError(t, updateControls(t, ctx, k, activeControls()))
	require.Equal(t, "10", k.GetLedger(ctx).NobleBacking.String())
	require.ErrorIs(t, updateControls(t, ctx, k, closed), canonicaltypes.ErrNobleBackingRemains)
	require.True(t, k.GetControls(ctx).NobleWithdrawalsEnabled)

	// Pausing is never blocked by remaining Noble backing.
	pausedClosed := closed
	pausedClosed.Mode = canonicaltypes.Mode_MODE_PAUSED
	require.NoError(t, updateControls(t, ctx, k, pausedClosed))
}

func TestUpdateControlsNobleWithdrawalsBlockedByPendingNoble(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	require.NoError(t, updateControls(t, ctx, k, activeControls()))
	ledger := k.GetLedger(ctx)
	ledger.PendingNoble = dtypes.NewInt(1)
	require.NoError(t, k.SetLedger(ctx, ledger))
	require.NoError(t, k.SetPending(ctx, canonicaltypes.PendingSettlement{
		Route: canonicaltypes.Route_ROUTE_NOBLE, SourceChannel: nobleChannel, Sequence: 1,
		Sender: sdk.AccAddress("s").String(), LogicalReceiver: "r", PhysicalReceiver: "r",
		LogicalDenom: types.UusdcDenom, PhysicalDenom: types.UusdcDenom, Amount: dtypes.NewInt(1),
	}))

	closed := activeControls()
	closed.NobleWithdrawalsEnabled = false
	require.ErrorIs(t, updateControls(t, ctx, k, closed), canonicaltypes.ErrNobleBackingRemains)

	require.NoError(t, k.Settle(ctx, mustPending(t, k, ctx, nobleChannel, 1), true, canonicaltypes.AttributeValueSuccess))
	require.NoError(t, updateControls(t, ctx, k, closed))
	require.False(t, k.GetControls(ctx).NobleWithdrawalsEnabled)
}

func TestUpdateControlsPauseThenResumeResnapshotsDrift(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	bank.set(sdk.AccAddress("holder"), types.UusdcDenom, math.NewInt(10))
	require.NoError(t, updateControls(t, ctx, k, activeControls()))
	requireInvariant(t, ctx, k)

	// Logical supply rises without a ledger credit.
	bank.set(sdk.AccAddress("refunded"), types.UusdcDenom, math.NewInt(3))
	_, broken := BackingInvariant(*k)(ctx)
	require.True(t, broken)

	// Updates that stay in GRADUAL do not re-check the supply equation.
	limits := activeControls()
	limits.MaxTransferAmount = dtypes.NewInt(500)
	require.NoError(t, updateControls(t, ctx, k, limits))

	require.NoError(t, updateControls(t, ctx, k, withMode(canonicaltypes.Mode_MODE_PAUSED)))
	require.Equal(t, "10", k.GetLedger(ctx).NobleBacking.String(), "pausing does not snapshot")

	require.NoError(t, updateControls(t, ctx, k, activeControls()))
	require.Equal(t, "13", k.GetLedger(ctx).NobleBacking.String())
	requireInvariant(t, ctx, k)
}

func TestUpdateControlsEnteringGradualRequiresInvariant(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	bank.set(sdk.AccAddress("holder"), types.UusdcDenom, math.NewInt(10))
	ledger := canonicaltypes.DefaultLedger()
	ledger.InjectiveBacking = dtypes.NewInt(4)
	require.NoError(t, k.SetLedger(ctx, ledger))

	// The module does not hold the physical coins backing injective_backing.
	err := updateControls(t, ctx, k, activeControls())
	require.ErrorIs(t, err, canonicaltypes.ErrInvalidLedger)
	require.Contains(t, err.Error(), "physical reserve")
	require.Equal(t, canonicaltypes.Mode_MODE_DISABLED, k.GetControls(ctx).Mode)

	// Pausing is still allowed.
	require.NoError(t, updateControls(t, ctx, k, withMode(canonicaltypes.Mode_MODE_PAUSED)))
}

func mustPending(
	t *testing.T,
	k *Keeper,
	ctx sdk.Context,
	channel string,
	sequence uint64,
) canonicaltypes.PendingSettlement {
	t.Helper()
	pending, found := k.GetPendingForPacket(ctx, channel, sequence)
	require.True(t, found)
	return pending
}
