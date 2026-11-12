package canonicalusdc

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	ratelimitutil "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ porttypes.Middleware = IBCMiddleware{}

type IBCMiddleware struct {
	app    porttypes.IBCModule
	keeper *keeper.Keeper
}

func NewIBCMiddleware(k *keeper.Keeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{keeper: k, app: app}
}

func (im IBCMiddleware) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (im IBCMiddleware) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return im.app.OnChanOpenTry(
		ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion,
	)
}

func (im IBCMiddleware) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID,
	counterpartyChannelID,
	counterpartyVersion string,
) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (im IBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	controls := im.keeper.GetControls(ctx)
	if controls.Mode == types.Mode_MODE_DISABLED {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}
	cacheCtx, write := ctx.CacheContext()
	ack := im.onRecvPacketActive(cacheCtx, packet, relayer, controls)
	if ack != nil && ack.Success() {
		write()
	}
	return ack
}

func (im IBCMiddleware) onRecvPacketActive(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	controls types.Controls,
) ibcexported.Acknowledgement {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	physicalDenom := ratelimitutil.ParseDenomFromRecvPacket(packet, data)
	if packet.DestinationChannel == controls.NobleChannel && physicalDenom == controls.LogicalDenom {
		return channeltypes.NewErrorAcknowledgement(types.ErrNobleDepositsDisabled)
	}
	if packet.DestinationChannel == controls.InjectiveChannel &&
		data.Denom == controls.InjectivePacketDenom &&
		physicalDenom == controls.InjectiveDenom {
		if controls.Mode == types.Mode_MODE_PAUSED {
			return channeltypes.NewErrorAcknowledgement(types.ErrPaused)
		}
		if controls.InjectiveDepositsDisabled {
			return channeltypes.NewErrorAcknowledgement(types.ErrInjectiveDepositsDisabled)
		}
		return im.receiveInjective(ctx, packet, data, relayer, controls)
	}
	// Logical USDC returning on any other channel unescrows coins that are
	// already part of supply, so the ledger does not change.
	return im.app.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) receiveInjective(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	relayer sdk.AccAddress,
	controls types.Controls,
) ibcexported.Acknowledgement {
	amount, ok := math.NewIntFromString(data.Amount)
	if !ok || !amount.IsPositive() {
		return channeltypes.NewErrorAcknowledgement(types.ErrInvalidAmount)
	}
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	// The module holds the reserve and locked liabilities, so it cannot be a
	// deposit owner; the error ack refunds the sender on the source chain.
	if receiver.Equals(types.ModuleAddress) {
		return channeltypes.NewErrorAcknowledgement(types.ErrInvalidReceiver)
	}
	if err := im.keeper.CreditInjectiveInbound(ctx, amount); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	physicalData := data
	physicalData.Receiver = types.ModuleAddress.String()
	physicalPacket := packet
	physicalPacket.Data = transfertypes.ModuleCdc.MustMarshalJSON(&physicalData)
	ack := im.app.OnRecvPacket(ctx, physicalPacket, relayer)
	if ack == nil || !ack.Success() {
		return ack
	}
	coins := sdk.NewCoins(sdk.NewCoin(controls.LogicalDenom, amount))
	if err := im.keeper.MintLogicalAndSend(ctx, receiver, coins); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeInbound,
		sdk.NewAttribute(types.AttributeKeyPhysicalDenom, controls.InjectiveDenom),
		sdk.NewAttribute(types.AttributeKeyLogicalDenom, controls.LogicalDenom),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(sdk.AttributeKeySender, data.Sender),
		sdk.NewAttribute(types.AttributeKeySourceChannel, packet.SourceChannel),
		sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprint(packet.Sequence)),
	))
	return ack
}

func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	pending, found := im.keeper.GetPendingForPacket(ctx, packet.SourceChannel, packet.Sequence)
	if !found {
		return im.forwardUnsettled(ctx, packet, isErrorAck(acknowledgement), func(ctx sdk.Context) error {
			return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
		})
	}
	if err := validatePendingPacket(pending, im.keeper.GetControls(ctx), packet); err != nil {
		return err
	}
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return err
	}
	refunded := refundsSender(ack)
	outcome := types.AttributeValueSuccess
	if refunded {
		outcome = types.AttributeValueError
	}
	cacheCtx, write := ctx.CacheContext()
	if err := im.app.OnAcknowledgementPacket(cacheCtx, packet, acknowledgement, relayer); err != nil {
		return err
	}
	if err := im.keeper.Settle(cacheCtx, pending, !refunded, outcome); err != nil {
		return err
	}
	write()
	return nil
}

func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	pending, found := im.keeper.GetPendingForPacket(ctx, packet.SourceChannel, packet.Sequence)
	if !found {
		return im.forwardUnsettled(ctx, packet, true, func(ctx sdk.Context) error {
			return im.app.OnTimeoutPacket(ctx, packet, relayer)
		})
	}
	if err := validatePendingPacket(pending, im.keeper.GetControls(ctx), packet); err != nil {
		return err
	}
	cacheCtx, write := ctx.CacheContext()
	if err := im.app.OnTimeoutPacket(cacheCtx, packet, relayer); err != nil {
		return err
	}
	if err := im.keeper.Settle(cacheCtx, pending, false, types.AttributeValueTimeout); err != nil {
		return err
	}
	write()
	return nil
}

func (im IBCMiddleware) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	return im.keeper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

func (im IBCMiddleware) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
	ack ibcexported.Acknowledgement,
) error {
	return im.keeper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (im IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return im.keeper.GetAppVersion(ctx, portID, channelID)
}

// forwardUnsettled hands a packet without a pending settlement to the inner
// app. In GRADUAL every logical send is gated and recorded, so an unsettled
// refunded Noble-trace packet was sent while DISABLED: ICS20 re-mints logical
// USDC that the GRADUAL snapshot excluded, and it is credited back to Noble
// backing so the supply equation keeps holding.
func (im IBCMiddleware) forwardUnsettled(
	ctx sdk.Context,
	packet channeltypes.Packet,
	refunded bool,
	forward func(sdk.Context) error,
) error {
	refund, ok := im.legacyNobleRefund(ctx, packet, refunded)
	if !ok {
		return forward(ctx)
	}
	cacheCtx, write := ctx.CacheContext()
	if err := forward(cacheCtx); err != nil {
		return err
	}
	if err := im.keeper.CreditLegacyNobleRefund(cacheCtx, refund); err != nil {
		return err
	}
	write()
	return nil
}

// legacyNobleRefund describes a refunded Noble-trace packet as the settlement
// it would have had if it had been sent through the module.
func (im IBCMiddleware) legacyNobleRefund(
	ctx sdk.Context,
	packet channeltypes.Packet,
	refunded bool,
) (types.PendingSettlement, bool) {
	controls := im.keeper.GetControls(ctx)
	if !refunded || controls.Mode == types.Mode_MODE_DISABLED ||
		packet.SourcePort != transfertypes.PortID || packet.SourceChannel != controls.NobleChannel {
		return types.PendingSettlement{}, false
	}
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err != nil {
		return types.PendingSettlement{}, false
	}
	amount, ok := math.NewIntFromString(data.Amount)
	if !ok || !amount.IsPositive() || data.Denom != controls.NobleTrace() {
		return types.PendingSettlement{}, false
	}
	return types.PendingSettlement{
		Route:            types.Route_ROUTE_NOBLE,
		SourceChannel:    packet.SourceChannel,
		Sequence:         packet.Sequence,
		Sender:           data.Sender,
		LogicalReceiver:  data.Receiver,
		PhysicalReceiver: data.Receiver,
		LogicalDenom:     controls.LogicalDenom,
		PhysicalDenom:    controls.LogicalDenom,
		Amount:           types.NewAmount(amount),
	}, true
}

// isErrorAck reports whether the bytes decode to an ack that ICS20 refunds.
// Bytes that do not decode are left to the inner app, which rejects them.
func isErrorAck(acknowledgement []byte) bool {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return false
	}
	return refundsSender(ack)
}

// refundsSender mirrors ICS20, which refunds only on an error response and
// treats any other decodable ack as delivered.
func refundsSender(ack channeltypes.Acknowledgement) bool {
	_, isError := ack.Response.(*channeltypes.Acknowledgement_Error)
	return isError
}

func validatePendingPacket(
	pending types.PendingSettlement,
	controls types.Controls,
	packet channeltypes.Packet,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err != nil {
		return err
	}
	if packet.SourcePort != transfertypes.PortID || packet.SourceChannel != pending.SourceChannel {
		return errorsmod.Wrap(types.ErrSettlementMismatch, "source")
	}
	pendingAmount, err := types.PositiveAmountOf(pending.Amount)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalidSettlement, "amount")
	}
	packetAmount, ok := math.NewIntFromString(data.Amount)
	if !ok || !packetAmount.Equal(pendingAmount) || data.Receiver != pending.PhysicalReceiver {
		return errorsmod.Wrap(types.ErrSettlementMismatch, "amount or receiver")
	}
	expectedDenom, expectedSender := controls.NobleTrace(), pending.Sender
	if pending.Route == types.Route_ROUTE_INJECTIVE {
		expectedDenom, expectedSender = controls.InjectiveTrace(), types.ModuleAddress.String()
	}
	if data.Denom != expectedDenom {
		return errorsmod.Wrap(types.ErrSettlementMismatch, "denom")
	}
	if data.Sender != expectedSender {
		return errorsmod.Wrap(types.ErrSettlementMismatch, "sender")
	}
	return nil
}
