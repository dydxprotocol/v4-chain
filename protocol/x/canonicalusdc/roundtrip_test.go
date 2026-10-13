package canonicalusdc

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ics20Stub mirrors ibc-go's sendTransfer and refundPacketToken: vouchers are
// resolved to their full denom trace, burned or escrowed, and the packet
// carries the full trace. Packets go through the canonical ICS4 gate.
type ics20Stub struct {
	bank    *receiveBank
	ics4    types.ICS4Wrapper
	traces  map[string]string
	packets []channeltypes.Packet
}

func (s *ics20Stub) Transfer(
	goCtx context.Context,
	msg *transfertypes.MsgTransfer,
) (*transfertypes.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fullDenomPath := msg.Token.Denom
	if strings.HasPrefix(msg.Token.Denom, "ibc/") {
		path, ok := s.traces[msg.Token.Denom]
		if !ok {
			return nil, fmt.Errorf("unknown denom trace %s", msg.Token.Denom)
		}
		fullDenomPath = path
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if transfertypes.SenderChainIsSource(msg.SourcePort, msg.SourceChannel, fullDenomPath) {
		escrow := transfertypes.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
		if err := s.bank.move(sender, escrow, sdk.NewCoins(msg.Token)); err != nil {
			return nil, err
		}
	} else {
		balance := s.bank.amount(sender, msg.Token.Denom)
		if balance.LT(msg.Token.Amount) {
			return nil, fmt.Errorf("insufficient funds")
		}
		s.bank.set(sender, msg.Token.Denom, balance.Sub(msg.Token.Amount))
	}
	data := transfertypes.NewFungibleTokenPacketData(
		fullDenomPath, msg.Token.Amount.String(), sender.String(), msg.Receiver, msg.Memo,
	)
	sequence, err := s.ics4.SendPacket(
		ctx, nil, msg.SourcePort, msg.SourceChannel, clienttypes.ZeroHeight(), msg.TimeoutTimestamp, data.GetBytes(),
	)
	if err != nil {
		return nil, err
	}
	s.packets = append(s.packets, channeltypes.Packet{
		Sequence:      sequence,
		SourcePort:    msg.SourcePort,
		SourceChannel: msg.SourceChannel,
		Data:          data.GetBytes(),
	})
	return &transfertypes.MsgTransferResponse{Sequence: sequence}, nil
}

func (*ics20Stub) UpdateParams(
	context.Context,
	*transfertypes.MsgUpdateParams,
) (*transfertypes.MsgUpdateParamsResponse, error) {
	return &transfertypes.MsgUpdateParamsResponse{}, nil
}

// refund applies ICS20's refund for a failed or timed-out packet.
func (s *ics20Stub) refund(packet channeltypes.Packet) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err != nil {
		return err
	}
	amount, ok := math.NewIntFromString(data.Amount)
	if !ok {
		return fmt.Errorf("invalid amount")
	}
	token := sdk.NewCoin(transfertypes.ParseDenomTrace(data.Denom).IBCDenom(), amount)
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}
	if transfertypes.SenderChainIsSource(packet.SourcePort, packet.SourceChannel, data.Denom) {
		escrow := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
		return s.bank.move(escrow, sender, sdk.NewCoins(token))
	}
	s.bank.set(sender, token.Denom, s.bank.amount(sender, token.Denom).Add(token.Amount))
	return nil
}

// ics20App is the inner IBC application below the canonical middleware.
type ics20App struct {
	*receiveApp
	stub *ics20Stub
}

func (a ics20App) OnAcknowledgementPacket(
	_ sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	_ sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return err
	}
	if ack.Success() {
		return nil
	}
	return a.stub.refund(packet)
}

func (a ics20App) OnTimeoutPacket(_ sdk.Context, packet channeltypes.Packet, _ sdk.AccAddress) error {
	return a.stub.refund(packet)
}

// sequencedICS4 is the channel layer below the canonical ICS4 gate.
type sequencedICS4 struct {
	next uint64
}

func (s *sequencedICS4) SendPacket(
	sdk.Context,
	*capabilitytypes.Capability,
	string,
	string,
	clienttypes.Height,
	uint64,
	[]byte,
) (uint64, error) {
	s.next++
	return s.next, nil
}

func (*sequencedICS4) WriteAcknowledgement(
	sdk.Context,
	*capabilitytypes.Capability,
	ibcexported.PacketI,
	ibcexported.Acknowledgement,
) error {
	return nil
}

func (*sequencedICS4) GetAppVersion(sdk.Context, string, string) (string, bool) {
	return transfertypes.Version, true
}

type roundTrip struct {
	ctx        sdk.Context
	keeper     *keeper.Keeper
	bank       *receiveBank
	ratelimit  *undoRatelimit
	stub       *ics20Stub
	middleware IBCMiddleware
	controls   types.Controls
	holder     sdk.AccAddress
}

func setupRoundTrip(t *testing.T, ledger types.Ledger) roundTrip {
	t.Helper()
	ctx, k, bank, ratelimit := setupKeeperWithICS4(t, &sequencedICS4{})
	controls := receiveControls()
	stub := &ics20Stub{
		bank: bank,
		ics4: k,
		traces: map[string]string{
			controls.LogicalDenom:   controls.NobleTrace(),
			controls.InjectiveDenom: controls.InjectiveTrace(),
		},
	}
	k.SetTransferMsgServer(stub)
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, ledger))
	return roundTrip{
		ctx:        ctx,
		keeper:     k,
		bank:       bank,
		ratelimit:  ratelimit,
		stub:       stub,
		middleware: NewIBCMiddleware(k, ics20App{receiveApp: &receiveApp{bank: bank}, stub: stub}),
		controls:   controls,
		holder:     sdk.AccAddress("holder"),
	}
}

func (r roundTrip) send(t *testing.T, channel string, amount int64) channeltypes.Packet {
	t.Helper()
	_, err := keeper.NewTransferDecorator(r.keeper).Transfer(r.ctx, &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    channel,
		Token:            sdk.NewInt64Coin(r.controls.LogicalDenom, amount),
		Sender:           r.holder.String(),
		Receiver:         "remote1receiver",
		TimeoutTimestamp: 1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, r.stub.packets)
	return r.stub.packets[len(r.stub.packets)-1]
}

func (r roundTrip) requireInvariant(t *testing.T) {
	t.Helper()
	reason, broken := keeper.BackingInvariant(*r.keeper)(r.ctx)
	require.False(t, broken, reason)
}

func errorAck() []byte {
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("remote error")).Acknowledgement()
}

func successAck() []byte {
	return channeltypes.NewResultAcknowledgement([]byte{1}).Acknowledgement()
}

func TestRoundTripNobleWithdrawal(t *testing.T) {
	cases := []struct {
		name      string
		settle    func(roundTrip, channeltypes.Packet) error
		wantNoble string
		wantUser  string
	}{
		{
			name: "success",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, successAck(), nil)
			},
			wantNoble: "6",
			wantUser:  "6",
		},
		{
			name: "error ack",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, errorAck(), nil)
			},
			wantNoble: "10",
			wantUser:  "10",
		},
		{
			name: "timeout",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnTimeoutPacket(r.ctx, packet, nil)
			},
			wantNoble: "10",
			wantUser:  "10",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ledger := types.DefaultLedger()
			ledger.NobleBacking = dtypes.NewInt(10)
			r := setupRoundTrip(t, ledger)
			r.bank.set(r.holder, r.controls.LogicalDenom, math.NewInt(10))
			r.requireInvariant(t)

			packet := r.send(t, r.controls.NobleChannel, 4)
			var data transfertypes.FungibleTokenPacketData
			require.NoError(t, transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
			require.Equal(t, r.controls.NobleTrace(), data.Denom)
			require.Equal(t, "6", r.bank.amount(r.holder, r.controls.LogicalDenom).String())
			require.Equal(t, "4", r.keeper.GetLedger(r.ctx).PendingNoble.String())
			r.requireInvariant(t)

			require.NoError(t, tc.settle(r, packet))
			require.Zero(t, r.keeper.PendingCount(r.ctx))
			require.Equal(t, tc.wantNoble, r.keeper.GetLedger(r.ctx).NobleBacking.String())
			require.Equal(t, "0", r.keeper.GetLedger(r.ctx).PendingNoble.String())
			require.Equal(t, tc.wantUser, r.bank.amount(r.holder, r.controls.LogicalDenom).String())
			r.requireInvariant(t)
		})
	}
}

func TestRoundTripInjectiveWithdrawal(t *testing.T) {
	cases := []struct {
		name          string
		settle        func(roundTrip, channeltypes.Packet) error
		wantInjective string
		wantUser      string
		wantUndone    int64
	}{
		{
			name: "success",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, successAck(), nil)
			},
			wantInjective: "6",
			wantUser:      "6",
		},
		{
			name: "error ack",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, errorAck(), nil)
			},
			wantInjective: "10",
			wantUser:      "10",
			wantUndone:    4,
		},
		{
			name: "timeout",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnTimeoutPacket(r.ctx, packet, nil)
			},
			wantInjective: "10",
			wantUser:      "10",
			wantUndone:    4,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ledger := types.DefaultLedger()
			ledger.InjectiveBacking = dtypes.NewInt(10)
			r := setupRoundTrip(t, ledger)
			r.bank.set(r.holder, r.controls.LogicalDenom, math.NewInt(10))
			r.bank.set(types.ModuleAddress, r.controls.InjectiveDenom, math.NewInt(10))
			r.requireInvariant(t)

			packet := r.send(t, r.controls.InjectiveChannel, 4)
			var data transfertypes.FungibleTokenPacketData
			require.NoError(t, transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
			require.Equal(t, r.controls.InjectiveTrace(), data.Denom)
			require.Equal(t, types.ModuleAddress.String(), data.Sender)
			require.Equal(t, "4", r.bank.amount(types.ModuleAddress, r.controls.LogicalDenom).String())
			require.Equal(t, "6", r.bank.amount(types.ModuleAddress, r.controls.InjectiveDenom).String())
			r.requireInvariant(t)

			require.NoError(t, tc.settle(r, packet))
			require.Zero(t, r.keeper.PendingCount(r.ctx))
			require.Equal(t, tc.wantInjective, r.keeper.GetLedger(r.ctx).InjectiveBacking.String())
			require.Equal(t, "0", r.keeper.GetLedger(r.ctx).PendingInjective.String())
			require.Equal(t, tc.wantUser, r.bank.amount(r.holder, r.controls.LogicalDenom).String())
			require.Equal(t, "0", r.bank.amount(types.ModuleAddress, r.controls.LogicalDenom).String())
			require.Equal(t, tc.wantUndone, r.ratelimit.undone)
			r.requireInvariant(t)
		})
	}
}

func TestRoundTripDirectLogicalSendIsGated(t *testing.T) {
	ledger := types.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(10)
	r := setupRoundTrip(t, ledger)
	r.bank.set(r.holder, r.controls.LogicalDenom, math.NewInt(10))
	for _, channel := range []string{r.controls.NobleChannel, "channel-9"} {
		_, err := r.stub.Transfer(r.ctx, &transfertypes.MsgTransfer{
			SourcePort:       transfertypes.PortID,
			SourceChannel:    channel,
			Token:            sdk.NewInt64Coin(r.controls.LogicalDenom, 1),
			Sender:           r.holder.String(),
			Receiver:         "remote1receiver",
			TimeoutTimestamp: 1,
		})
		require.ErrorIs(t, err, types.ErrBypass)
	}
}

func TestRoundTripExistingInjectiveHolderCanUnwind(t *testing.T) {
	settlements := []struct {
		name     string
		settle   func(roundTrip, channeltypes.Packet) error
		wantUser string
	}{
		{
			name: "success",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, successAck(), nil)
			},
			wantUser: "2",
		},
		{
			name: "error ack",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnAcknowledgementPacket(r.ctx, packet, errorAck(), nil)
			},
			wantUser: "5",
		},
		{
			name: "timeout",
			settle: func(r roundTrip, packet channeltypes.Packet) error {
				return r.middleware.OnTimeoutPacket(r.ctx, packet, nil)
			},
			wantUser: "5",
		},
	}
	for _, mode := range []types.Mode{types.Mode_MODE_GRADUAL, types.Mode_MODE_PAUSED} {
		for _, tc := range settlements {
			t.Run(mode.String()+"/"+tc.name, func(t *testing.T) {
				ledger := types.DefaultLedger()
				ledger.InjectiveBacking = dtypes.NewInt(10)
				r := setupRoundTrip(t, ledger)
				r.controls.Mode = mode
				require.NoError(t, r.keeper.SetControls(r.ctx, r.controls))
				r.bank.set(r.holder, r.controls.LogicalDenom, math.NewInt(10))
				r.bank.set(r.holder, r.controls.InjectiveDenom, math.NewInt(5))
				r.bank.set(types.ModuleAddress, r.controls.InjectiveDenom, math.NewInt(10))
				r.requireInvariant(t)

				_, err := keeper.NewTransferDecorator(r.keeper).Transfer(r.ctx, &transfertypes.MsgTransfer{
					SourcePort:       transfertypes.PortID,
					SourceChannel:    r.controls.InjectiveChannel,
					Token:            sdk.NewInt64Coin(r.controls.InjectiveDenom, 3),
					Sender:           r.holder.String(),
					Receiver:         "inj1receiver",
					TimeoutTimestamp: 1,
				})
				require.NoError(t, err)
				packet := r.stub.packets[len(r.stub.packets)-1]
				var data transfertypes.FungibleTokenPacketData
				require.NoError(t, transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
				require.Equal(t, r.controls.InjectiveTrace(), data.Denom)
				require.Equal(t, r.holder.String(), data.Sender)
				require.Zero(t, r.keeper.PendingCount(r.ctx))
				require.Equal(t, "2", r.bank.amount(r.holder, r.controls.InjectiveDenom).String())
				r.requireInvariant(t)

				require.NoError(t, tc.settle(r, packet))
				require.Equal(t, tc.wantUser, r.bank.amount(r.holder, r.controls.InjectiveDenom).String())
				require.Equal(t, "10", r.bank.amount(r.holder, r.controls.LogicalDenom).String())
				require.Equal(t, "10", r.bank.amount(types.ModuleAddress, r.controls.InjectiveDenom).String())
				require.Equal(t, "10", r.keeper.GetLedger(r.ctx).InjectiveBacking.String())
				require.Equal(t, int64(0), r.ratelimit.charged)
				require.Equal(t, int64(0), r.ratelimit.undone)
				r.requireInvariant(t)
			})
		}
	}
}
