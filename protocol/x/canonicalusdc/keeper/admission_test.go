package keeper

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// fundRoutes gives both routes backing of 50 and the sender 50 logical USDC.
func fundRoutes(t *testing.T, ctx sdk.Context, k *Keeper, bank *bankMock, sender sdk.AccAddress) {
	t.Helper()
	ledger := canonicaltypes.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(50)
	ledger.InjectiveBacking = dtypes.NewInt(50)
	require.NoError(t, k.SetLedger(ctx, ledger))
	bank.set(sender, types.UusdcDenom, math.NewInt(50))
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(50))
}

func withdrawal(sender sdk.AccAddress, channel string, amount int64) *transfertypes.MsgTransfer {
	return &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    channel,
		Token:            sdk.NewInt64Coin(types.UusdcDenom, amount),
		Sender:           sender.String(),
		Receiver:         "remote1receiver",
		TimeoutTimestamp: 1,
	}
}

func eventAttributes(t *testing.T, ctx sdk.Context, eventType string) map[string]string {
	t.Helper()
	var found []sdk.Event
	for _, event := range ctx.EventManager().Events() {
		if event.Type == eventType {
			found = append(found, event)
		}
	}
	require.Len(t, found, 1, "events of type %s", eventType)
	attributes := make(map[string]string, len(found[0].Attributes))
	for _, attribute := range found[0].Attributes {
		attributes[attribute.Key] = attribute.Value
	}
	return attributes
}

// A Noble-route memo is forwarded unless packet-forward-middleware would read
// a forwarding instruction from it.
func TestNobleWithdrawalForwardsMemosWithoutForwardKey(t *testing.T) {
	memos := map[string]string{
		"plain text":               "wallet-memo",
		"json without forward":     `{"note":"x","tag":7}`,
		"forward in another case":  `{"Forward":{"receiver":"x"}}`,
		"forward below top level":  `{"note":{"forward":{"receiver":"x"}}}`,
		"forward as a json string": `"forward"`,
		"forward in a json array":  `[{"forward":{"receiver":"x"}}]`,
		"truncated forward object": `{"forward":{"receiver":"x"}`,
	}
	for name, memo := range memos {
		t.Run(name, func(t *testing.T) {
			ctx, k, bank, transfer, _ := setupKeeper(t)
			require.NoError(t, k.SetControls(ctx, activeControls()))
			sender := sdk.AccAddress("memo-sender")
			fundRoutes(t, ctx, k, bank, sender)
			msg := withdrawal(sender, nobleChannel, 10)
			msg.Memo = memo

			_, err := NewTransferDecorator(k).Transfer(ctx, msg)
			require.NoError(t, err)
			require.Equal(t, memo, transfer.seen.Memo)
			require.Equal(t, uint32(1), k.PendingCount(ctx))
		})
	}
}

// canonical_usdc_transfer and canonical_usdc_settlement each describe the
// withdrawal on their own: who sent how much of which denom to whom.
func TestCanonicalEventsDescribeTheWithdrawal(t *testing.T) {
	routes := map[string]struct {
		channel       string
		route         canonicaltypes.Route
		physicalDenom string
	}{
		"noble":     {channel: nobleChannel, route: canonicaltypes.Route_ROUTE_NOBLE, physicalDenom: types.UusdcDenom},
		"injective": {channel: injChannel, route: canonicaltypes.Route_ROUTE_INJECTIVE, physicalDenom: physicalDenom},
	}
	for name, tc := range routes {
		t.Run(name, func(t *testing.T) {
			ctx, k, bank, transfer, _ := setupKeeper(t)
			require.NoError(t, k.SetControls(ctx, activeControls()))
			sender := sdk.AccAddress("event-sender")
			fundRoutes(t, ctx, k, bank, sender)
			transfer.logicalOwner = sender

			response, err := NewTransferDecorator(k).Transfer(ctx, withdrawal(sender, tc.channel, 10))
			require.NoError(t, err)
			common := map[string]string{
				canonicaltypes.AttributeKeyRoute:         tc.route.String(),
				canonicaltypes.AttributeKeyLogicalDenom:  types.UusdcDenom,
				canonicaltypes.AttributeKeyPhysicalDenom: tc.physicalDenom,
				canonicaltypes.AttributeKeySourceChannel: tc.channel,
				canonicaltypes.AttributeKeySequence:      "7",
				transfertypes.AttributeKeyAmount:         "10",
				sdk.AttributeKeySender:                   sender.String(),
				transfertypes.AttributeKeyReceiver:       "remote1receiver",
			}
			want := map[string]string{canonicaltypes.AttributeKeyPending: "true"}
			for key, value := range common {
				want[key] = value
			}
			require.Equal(t, want, eventAttributes(t, ctx, canonicaltypes.EventTypeCanonicalTransfer))

			pending, found := k.GetPendingForPacket(ctx, tc.channel, response.Sequence)
			require.True(t, found)
			require.NoError(t, k.Settle(ctx, pending, true, canonicaltypes.AttributeValueSuccess))
			want = map[string]string{canonicaltypes.AttributeKeyOutcome: canonicaltypes.AttributeValueSuccess}
			for key, value := range common {
				want[key] = value
			}
			require.Equal(t, want, eventAttributes(t, ctx, canonicaltypes.EventTypeSettlement))
		})
	}
}

// Rejections surface the module codespace and code in the transaction result.
func TestRejectionsCarryModuleCodespace(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	sender := sdk.AccAddress("codespace-sender")
	fundRoutes(t, ctx, k, bank, sender)
	decorator := NewTransferDecorator(k)

	withMemo := withdrawal(sender, injChannel, 10)
	withMemo.Memo = "partner-memo"
	cases := map[string]struct {
		msg  *transfertypes.MsgTransfer
		want *errorsmod.Error
	}{
		"unsupported channel": {msg: withdrawal(sender, "channel-9", 10), want: canonicaltypes.ErrUnsupportedChannel},
		"memo":                {msg: withMemo, want: canonicaltypes.ErrMemoNotAllowed},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := decorator.Transfer(ctx, tc.msg)
			require.ErrorIs(t, err, tc.want)
			codespace, code, _ := errorsmod.ABCIInfo(err, false)
			require.Equal(t, canonicaltypes.ModuleName, codespace)
			require.Equal(t, tc.want.ABCICode(), code)
		})
	}
}
