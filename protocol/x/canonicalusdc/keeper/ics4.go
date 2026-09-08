package keeper

import (
	"fmt"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	controls := k.GetControls(ctx)
	if controls.Mode != types.Mode_MODE_DISABLED {
		var packetData transfertypes.FungibleTokenPacketData
		if err := transfertypes.ModuleCdc.UnmarshalJSON(data, &packetData); err != nil {
			return 0, err
		}
		if canonicalPacketData(controls, sourceChannel, packetData.Denom) {
			authorized, _ := ctx.Value(transferAuthorizationKey{}).(bool)
			if !authorized {
				return 0, fmt.Errorf("%w: channel %s", types.ErrBypass, sourceChannel)
			}
		}
	}
	return k.ics4Wrapper.SendPacket(
		ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data,
	)
}

func canonicalPacketData(controls types.Controls, sourceChannel, denom string) bool {
	nobleTrace := transfertypes.GetDenomPrefix(transfertypes.PortID, controls.NobleChannel) + controls.NoblePacketDenom
	injectiveTrace := transfertypes.GetDenomPrefix(
		transfertypes.PortID,
		controls.InjectiveChannel,
	) + controls.InjectivePacketDenom
	return (sourceChannel == controls.NobleChannel && denom == controls.NoblePacketDenom) ||
		(sourceChannel == controls.InjectiveChannel && denom == controls.InjectivePacketDenom) ||
		denom == nobleTrace || denom == injectiveTrace
}

func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
	ack ibcexported.Acknowledgement,
) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (k Keeper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return k.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
