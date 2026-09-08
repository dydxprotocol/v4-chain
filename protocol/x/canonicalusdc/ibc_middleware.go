package canonicalusdc

import (
	"fmt"

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
		return im.receiveInjective(ctx, packet, data, relayer, controls)
	}
	var legacyReturnAmount math.Int
	if physicalDenom == controls.LogicalDenom {
		amount, ok := math.NewIntFromString(data.Amount)
		if !ok || !amount.IsPositive() {
			return channeltypes.NewErrorAcknowledgement(types.ErrInvalidAmount)
		}
		legacyDownstream, err := types.ParseAmount(im.keeper.GetLedger(ctx).LegacyDownstream)
		if err != nil || legacyDownstream.LT(amount) {
			return channeltypes.NewErrorAcknowledgement(types.ErrUnclassifiedReturn)
		}
		legacyReturnAmount = amount
	}
	ack := im.app.OnRecvPacket(ctx, packet, relayer)
	if ack != nil && ack.Success() && !legacyReturnAmount.IsNil() {
		ledger := im.keeper.GetLedger(ctx)
		legacyDownstream, _ := types.ParseAmount(ledger.LegacyDownstream)
		ledger.LegacyDownstream = legacyDownstream.Sub(legacyReturnAmount).String()
		if err := im.keeper.SetLedger(ctx, ledger); err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	}
	return ack
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
	maxTransfer, err := types.ParsePositiveAmount(controls.MaxTransferAmount)
	if err != nil || amount.GT(maxTransfer) {
		return channeltypes.NewErrorAcknowledgement(types.ErrInvalidAmount)
	}
	memo, isFunding, err := types.ParseFundingMemo(data.Memo)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	if !isFunding {
		if _, err := sdk.AccAddressFromBech32(data.Receiver); err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	} else if err := im.validateFunding(ctx, data.Receiver, memo, controls); err != nil {
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
	ledger := im.keeper.GetLedger(ctx)
	if isFunding {
		participant, _ := im.keeper.GetParticipant(ctx, memo.Controller)
		participant.Funded = addLedger(participant.Funded, amount)
		ledger.RestrictedFunding = addLedger(ledger.RestrictedFunding, amount)
		if err := im.keeper.SetParticipant(ctx, participant); err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	} else {
		receiver, _ := sdk.AccAddressFromBech32(data.Receiver)
		coins := sdk.NewCoins(sdk.NewCoin(controls.LogicalDenom, amount))
		if err := im.keeper.MintLogicalAndSend(ctx, receiver, coins); err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
		ledger.InjectiveBacking = addLedger(ledger.InjectiveBacking, amount)
	}
	if err := im.keeper.SetLedger(ctx, ledger); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	event := sdk.NewEvent(
		types.EventTypeBackingFunding,
		sdk.NewAttribute(types.AttributeKeyPhysicalDenom, controls.InjectiveDenom),
		sdk.NewAttribute(types.AttributeKeyLogicalDenom, controls.LogicalDenom),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, amount.String()),
	)
	if isFunding {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyController, memo.Controller),
			sdk.NewAttribute(types.AttributeKeyNobleRecipient, memo.NobleRecipient),
		)
	}
	ctx.EventManager().EmitEvent(event)
	return ack
}

func (im IBCMiddleware) validateFunding(
	ctx sdk.Context,
	receiver string,
	memo types.FundingMemo,
	controls types.Controls,
) error {
	if receiver != types.ModuleAddress.String() ||
		memo.Version != controls.MemoVersion ||
		memo.Action != types.FundBackingSwapAction {
		return types.ErrInvalidFundingMemo
	}
	participant, ok := im.keeper.GetParticipant(ctx, memo.Controller)
	if !ok || participant.NobleRecipient != memo.NobleRecipient {
		return types.ErrUnauthorized
	}
	return nil
}

func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	pending, found := im.keeper.GetPendingForPacket(ctx, packet.SourceChannel, packet.Sequence)
	if !found {
		if im.keeper.HasCompleted(ctx, packet.SourceChannel, packet.Sequence) {
			return nil
		}
		return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}
	if err := validatePendingPacket(pending, im.keeper.GetControls(ctx), packet); err != nil {
		return err
	}
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return err
	}
	outcome := types.AttributeValueError
	if ack.Success() {
		outcome = types.AttributeValueSuccess
	}
	cacheCtx, write := ctx.CacheContext()
	if err := im.app.OnAcknowledgementPacket(cacheCtx, packet, acknowledgement, relayer); err != nil {
		return err
	}
	if err := im.keeper.Settle(cacheCtx, pending, ack.Success(), outcome); err != nil {
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
		if im.keeper.HasCompleted(ctx, packet.SourceChannel, packet.Sequence) {
			return nil
		}
		return im.app.OnTimeoutPacket(ctx, packet, relayer)
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
		return fmt.Errorf("pending settlement source does not match packet")
	}
	if data.Amount != pending.Amount || data.Receiver != pending.PhysicalReceiver {
		return fmt.Errorf("pending settlement does not match packet")
	}
	expectedDenom := controls.NoblePacketDenom
	if pending.Route == types.Route_ROUTE_INJECTIVE {
		expectedDenom = controls.InjectivePacketDenom
	}
	if data.Denom != expectedDenom {
		return fmt.Errorf("pending settlement denom does not match packet")
	}
	expectedSender := pending.Sender
	if pending.Route == types.Route_ROUTE_INJECTIVE || pending.Route == types.Route_ROUTE_BACKING_SWAP {
		expectedSender = types.ModuleAddress.String()
	}
	if data.Sender != expectedSender {
		return fmt.Errorf("pending settlement sender does not match packet")
	}
	return nil
}

func addLedger(value string, amount math.Int) string {
	current, err := types.ParseAmount(value)
	if err != nil {
		panic(err)
	}
	return current.Add(amount).String()
}
