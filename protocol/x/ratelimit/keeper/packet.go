package keeper

// This file includes keeper methods used by the IBC middleware for processing IBC packets.
// Re-uses/adapts Stride x/ratelimit implementation: https://github.com/Stride-Labs/stride/tree/4913e1dd1a/x/ratelimit
// See v4-chain/protocol/x/ratelimit/LICENSE and v4-chain/protocol/x/ratelimit/README.md for licensing information.

import (
	"encoding/json"
	"math/big"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// Remove a pending packet sequence number from the store
// Used after the ack or timeout for a packet has been received
func (k Keeper) RemovePendingSendPacket(ctx sdk.Context, channelId string, sequence uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingSendPacketPrefix))
	key := types.GetPendingSendPacketKey(channelId, sequence)
	store.Delete(key)
}

// Sets a pending packet sequence number in the store
func (k Keeper) SetPendingSendPacket(ctx sdk.Context, channelId string, sequence uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingSendPacketPrefix))
	key := types.GetPendingSendPacketKey(channelId, sequence)
	store.Set(key, []byte{1}) // Use a single bit to indicate packet is pending.
}

// Checks whether the packet sequence number is in the store - indicating that it is a pending
// packet.
func (k Keeper) HasPendingSendPacket(ctx sdk.Context, channelId string, sequence uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingSendPacketPrefix))
	key := types.GetPendingSendPacketKey(channelId, sequence)
	return store.Has(key)
}

// Middleware implementation for OnAckPacket
// It is called on the sender chain when a relayer relays back the acknowledgement from the receiver chain.
// On the dYdX chain, this includes the “response” of the receiver chain for outbound transfer from dYdX.
func (k Keeper) AcknowledgeIBCTransferPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	// Check whether the ack was a success or error
	ackResponse, err := util.UnpackAcknowledgementResponseForTransfer(ctx, k.Logger(ctx), acknowledgement)
	if err != nil {
		return err
	}

	// Parse the denom, channelId, and amount from the packet
	packetInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	if err != nil {
		return err
	}

	// If the ack was successful, remove the pending packet
	if ackResponse.Status == types.AckResponseStatus_SUCCESS {
		k.RemovePendingSendPacket(ctx, packetInfo.ChannelID, packet.Sequence)
		return nil
	}

	// If the ack failed, undo the change to the capacity
	k.UndoSendPacket(ctx, packetInfo.ChannelID, packet.Sequence, packetInfo.Denom, packetInfo.Amount)
	return nil
}

// Middleware implementation for OnAckPacket
// It is called on the sender chain when a relayer relays back the acknowledgement from the receiver chain.
// On the dYdX chain, this includes the “response” of the receiver chain for outbound transfer from dYdX.
func (k Keeper) RedoMintTradingDAIIfAcknowledgeIBCTransferPacketFails(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	// Check whether the ack was a success or error
	ackResponse, err := util.UnpackAcknowledgementResponseForTransfer(ctx, k.Logger(ctx), acknowledgement)
	if err != nil {
		return err
	}

	// If the ack was successful return
	if ackResponse.Status == types.AckResponseStatus_SUCCESS {
		return nil
	}

	// Parse the denom, channelId, and amount from the packet
	packetInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	if err != nil {
		return err
	}

	// We use sDaiDenom, since denom is hashed when we parse
	if packetInfo.Denom != types.SDaiDenom {
		return nil
	}

	// Redeposit sDAI
	return k.MintTradingDAIToUserAccount(ctx, packetInfo.Sender, packetInfo.Amount)
}

// Middleware implementation for OnTimeout
// It is triggered by a relayer with MsgTimeout on the sender chain when timeoutHeight is
// reached for a sent packet but acknowledgement has not been received. It should therefore
// revert the capacity change.
func (k Keeper) TimeoutIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	packetInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	if err != nil {
		return err
	}

	k.UndoSendPacket(ctx, packetInfo.ChannelID, packet.Sequence, packetInfo.Denom, packetInfo.Amount)
	return nil
}

func (k Keeper) UndoMintTradingDAIIfAfterTimeoutIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	packetInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	if err != nil {
		return err
	}

	// We use sDaiDenom, since denom is hashed when we parse
	if packetInfo.Denom != types.SDaiDenom {
		return nil
	}

	return k.MintTradingDAIToUserAccount(ctx, packetInfo.Sender, packetInfo.Amount)
}

// If a SendPacket fails or times out, undo the capacity decrease that happened during the send
// Idempotent - has no effect on a previously removed pending packet.
func (k Keeper) UndoSendPacket(
	ctx sdk.Context,
	channelId string,
	sequence uint64,
	denom string,
	amount *big.Int,
) {
	if !k.HasPendingSendPacket(ctx, channelId, sequence) {
		return
	}
	// Undo'ing capacity change from the withdrawal.
	k.UndoWithdrawal(ctx, denom, amount)
	k.RemovePendingSendPacket(ctx, channelId, sequence)

	k.Logger(ctx).Info(
		"SendPacket timeout'ed or failed acknowledgement on the receiver chain. Reverted capacity change.",
		"channel_id",
		channelId,
		"sequence",
		sequence,
		"denom",
		denom,
		"amount",
		amount.String(),
	)
}

// SendPacket wraps IBC ChannelKeeper's SendPacket function
// If the packet does not get rate limited, it passes the packet to the IBC Channel keeper
func (k Keeper) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	// The packet must first be sent up the stack to get the sequence number from the channel keeper
	sequence, err = k.ics4Wrapper.SendPacket(
		ctx,
		channelCap,
		sourcePort,
		sourceChannel,
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	if err != nil {
		return sequence, err
	}

	err = k.TrySendRateLimitedPacket(ctx, channeltypes.Packet{
		Sequence:         sequence,
		SourceChannel:    sourceChannel,
		SourcePort:       sourcePort,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Data:             data,
	})
	if err != nil {
		k.Logger(ctx).Info(
			"ICS20 packet send was denied",
			"exec_mode",
			ctx.ExecMode(),
			"error",
			err.Error(),
		)
		return 0, err
	}
	return sequence, err
}

// Middleware implementation for SendPacket with rate limiting
// Checks whether the rate limit has been exceeded - and if it hasn't, sends the packet
func (k Keeper) TrySendRateLimitedPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	packetInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	if err != nil {
		return err
	}

	if err := k.ProcessWithdrawal(ctx, packetInfo.Denom, packetInfo.Amount); err != nil {
		// Some of the capacities were inefficient. Return error to fail the transaction.
		return err
	}

	// Store the sequence number of the packet so that if the transfer fails,
	// we can identify if it was sent during this quota and can revert the outflow
	k.SetPendingSendPacket(ctx, packetInfo.ChannelID, packet.Sequence)

	return nil
}

// WriteAcknowledgement wraps IBC ChannelKeeper's WriteAcknowledgement function
func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
	acknowledgement ibcexported.Acknowledgement,
) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, acknowledgement)
}

// GetAppVersion wraps IBC ChannelKeeper's GetAppVersion function
func (k Keeper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return k.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}

// PreprocessSendPacket implements the ICS4WrapperWithPreprocess interface
func (k Keeper) PreprocessSendPacket(ctx sdk.Context, packet []byte) error {
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet, &packetData); err != nil {
		return err
	}

	if packetData.Denom == types.SDaiBaseDenomFullPath {
		amount, senderAddress, _, err := util.GetValidatedFungibleTokenPacketData(packetData)
		if err != nil {
			return err
		}
		err = k.WithdrawSDaiFromTDai(ctx, senderAddress, amount)
		if err != nil {
			return err
		}
	}
	return nil
}
