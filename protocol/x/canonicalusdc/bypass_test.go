package canonicalusdc

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// A Noble withdrawal sent while DISABLED burns logical USDC. If it fails after
// governance enters GRADUAL, ICS20 re-mints it without a pending settlement;
// the middleware credits it back to Noble backing so the supply equation holds.
func TestBypassLegacyNobleRefundAfterGradualRestoresBacking(t *testing.T) {
	r := setupRoundTrip(t, types.DefaultLedger())
	disabled := r.controls
	disabled.Mode = types.Mode_MODE_DISABLED
	require.NoError(t, r.keeper.SetControls(r.ctx, disabled))
	r.bank.set(r.holder, r.controls.LogicalDenom, math.NewInt(10))

	packet := r.send(t, r.controls.NobleChannel, 4)
	require.Zero(t, r.keeper.PendingCount(r.ctx))

	// Entering GRADUAL snapshots noble_backing to the current supply (6).
	require.NoError(t, r.keeper.SetControls(r.ctx, r.controls))
	ledger := types.DefaultLedger()
	ledger.NobleBacking = types.NewAmount(r.bank.GetSupply(r.ctx, r.controls.LogicalDenom).Amount)
	require.NoError(t, r.keeper.SetLedger(r.ctx, ledger))
	r.requireInvariant(t)

	require.NoError(t, r.middleware.OnTimeoutPacket(r.ctx, packet, nil))
	require.Equal(t, "10", r.bank.amount(r.holder, r.controls.LogicalDenom).String())
	require.Equal(t, "10", r.keeper.GetLedger(r.ctx).NobleBacking.String())
	r.requireInvariant(t)
}

// Pending settlements share one global cap across both routes and have no
// minimum amount, so dust withdrawals block every canonical withdrawal while
// packets are unrelayed.
func TestBypassDustWithdrawalsExhaustSharedPendingCap(t *testing.T) {
	ledger := types.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(10)
	ledger.InjectiveBacking = dtypes.NewInt(5)
	r := setupRoundTrip(t, ledger)
	controls := r.controls
	controls.MaxPendingSettlements = 2
	require.NoError(t, r.keeper.SetControls(r.ctx, controls))
	r.bank.set(r.holder, controls.LogicalDenom, math.NewInt(15))
	r.bank.set(types.ModuleAddress, controls.InjectiveDenom, math.NewInt(5))

	r.send(t, controls.NobleChannel, 1)
	r.send(t, controls.NobleChannel, 1)

	_, err := keeper.NewTransferDecorator(r.keeper).Transfer(r.ctx, &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    controls.InjectiveChannel,
		Token:            sdk.NewInt64Coin(controls.LogicalDenom, 5),
		Sender:           r.holder.String(),
		Receiver:         "remote1receiver",
		TimeoutTimestamp: 1,
	})
	require.ErrorIs(t, err, types.ErrPendingLimit)
}

// The packet-level gate catches logical USDC on any channel and module-sent
// Injective USDC whenever the caller did not go through the decorator, e.g.
// another module calling the transfer keeper directly.
func TestBypassGateRequiresAuthorization(t *testing.T) {
	ctx, k, _, _ := setupKeeperWithICS4(t, &sequencedICS4{})
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	user := sdk.AccAddress("holder").String()
	cases := []struct {
		name    string
		denom   string
		sender  string
		channel string
		wantErr bool
	}{
		{"logical on noble", controls.NobleTrace(), user, controls.NobleChannel, true},
		{"logical on other channel", controls.NobleTrace(), user, "channel-9", true},
		{"logical on injective", controls.NobleTrace(), user, controls.InjectiveChannel, true},
		{"module injective", controls.InjectiveTrace(), types.ModuleAddress.String(), controls.InjectiveChannel, true},
		{"user unwinds injective", controls.InjectiveTrace(), user, controls.InjectiveChannel, false},
		{"other denom", "uatom", user, "channel-9", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := transfertypes.NewFungibleTokenPacketData(tc.denom, "1", tc.sender, "remote", "").GetBytes()
			_, err := k.SendPacket(ctx, nil, transfertypes.PortID, tc.channel, clienttypes.ZeroHeight(), 1, data)
			if tc.wantErr {
				require.ErrorIs(t, err, types.ErrBypass)
				return
			}
			require.NoError(t, err)
		})
	}
}

// While Injective backing is outstanding, governance can neither re-bind the
// Injective route (stranding the reserve in the old physical denom) nor disable
// routing (letting Injective-backed supply leave over Noble). PAUSED stays open.
func TestBypassGovernanceRejectsTransitionsWithInjectiveBacking(t *testing.T) {
	ctx, k, bank := bypassGovKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	ledger := types.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(10)
	ledger.InjectiveBacking = dtypes.NewInt(5)
	require.NoError(t, k.SetLedger(ctx, ledger))
	bank.set(sdk.AccAddress("holder"), controls.LogicalDenom, math.NewInt(15))
	bank.set(types.ModuleAddress, controls.InjectiveDenom, math.NewInt(5))
	reason, broken := keeper.BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)
	server := keeper.NewMsgServer(k)

	rebound := controls
	rebound.InjectiveChannel = "channel-7"
	rebound.InjectiveDenom = transfertypes.ParseDenomTrace("transfer/channel-7/uusdc").IBCDenom()
	_, err := server.UpdateControls(ctx, &types.MsgUpdateControls{
		Authority: lib.GovModuleAddress.String(), Controls: rebound,
	})
	require.ErrorIs(t, err, types.ErrInvalidControls)
	require.Contains(t, err.Error(), "route identity")
	reason, broken = keeper.BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)

	disabled := controls
	disabled.Mode = types.Mode_MODE_DISABLED
	_, err = server.UpdateControls(ctx, &types.MsgUpdateControls{
		Authority: lib.GovModuleAddress.String(), Controls: disabled,
	})
	require.ErrorIs(t, err, types.ErrInvalidControls)
	require.Contains(t, err.Error(), "cannot disable")

	paused := controls
	paused.Mode = types.Mode_MODE_PAUSED
	_, err = server.UpdateControls(ctx, &types.MsgUpdateControls{
		Authority: lib.GovModuleAddress.String(), Controls: paused,
	})
	require.NoError(t, err)
	require.Equal(t, types.Mode_MODE_PAUSED, k.GetControls(ctx).Mode)
}

func bypassGovKeeper(t *testing.T) (sdk.Context, *keeper.Keeper, *receiveBank) {
	t.Helper()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	bank := newReceiveBank()
	k := keeper.NewKeeper(
		codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()), key, nil, bank, &sequencedICS4{},
		&undoRatelimit{}, nil, nil, []string{lib.GovModuleAddress.String()},
	)
	return ctx, k, bank
}
