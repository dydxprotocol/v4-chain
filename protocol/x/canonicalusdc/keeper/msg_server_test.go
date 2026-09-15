package keeper

import (
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
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
			require.Equal(t, "0", k.GetLedger(ctx).NobleBacking)
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
		ledger.PendingNoble = "1"
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
			Amount:           "1",
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
			errContains: "cannot disable with pending",
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
		injectiveBacking string
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
			name:      "enabling paused",
			mode:      canonicaltypes.Mode_MODE_PAUSED,
			uusdc:     8,
			wantMode:  canonicaltypes.Mode_MODE_PAUSED,
			wantNoble: "8",
		},
		{
			name:             "residual after classified",
			mode:             canonicaltypes.Mode_MODE_GRADUAL,
			uusdc:            10,
			physical:         4,
			injectiveBacking: "4",
			wantMode:         canonicaltypes.Mode_MODE_GRADUAL,
			wantNoble:        "6",
			wantInjective:    "4",
		},
		{
			name:             "supply below classified",
			mode:             canonicaltypes.Mode_MODE_GRADUAL,
			uusdc:            3,
			injectiveBacking: "4",
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
			if tc.injectiveBacking != "" {
				ledger := canonicaltypes.DefaultLedger()
				ledger.InjectiveBacking = tc.injectiveBacking
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
				require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking)
				require.Equal(t, tc.wantInjective, k.GetLedger(ctx).InjectiveBacking)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantMode, k.GetControls(ctx).Mode)
			require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking)
			if tc.wantInjective != "" {
				require.Equal(t, tc.wantInjective, k.GetLedger(ctx).InjectiveBacking)
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
			require.Equal(t, tc.wantNoble, k.GetLedger(ctx).NobleBacking)
		})
	}
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
