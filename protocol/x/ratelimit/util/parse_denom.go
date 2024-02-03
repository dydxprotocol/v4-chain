package util

// This file includes utility methods used by the IBC middleware for parsing IBC denoms.
// Re-uses Stride x/ratelimit implementation: https://github.com/Stride-Labs/stride/tree/4913e1dd1a/x/ratelimit
// See v4-chain/protocol/x/ratelimit/LICENSE and v4-chain/protocol/x/ratelimit/README.md for licensing information.
import (
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// Parse the denom from the Send Packet that will be used by the rate limit module
// The denom that the rate limiter will use for a SEND packet depends on whether
// it was a NATIVE token (e.g. adv4tnt, etc.) or NON-NATIVE token (e.g. ibc/...)...
//
// We can identify if the token is native or not by parsing the trace denom from the packet
// If the token is NATIVE, it will not have a prefix (e.g. adv4tnt),
// and if it is NON-NATIVE, it will have a prefix (e.g. transfer/channel-2/uosmo)
//
// For NATIVE denoms, return as is (e.g. adv4tnt)
// For NON-NATIVE denoms, take the ibc hash (e.g. hash "transfer/channel-2/usoms" into "ibc/...")
func ParseDenomFromSendPacket(packet ibctransfertypes.FungibleTokenPacketData) (denom string) {
	// Determine the denom by looking at the denom trace path
	denomTrace := ibctransfertypes.ParseDenomTrace(packet.Denom)

	// Native assets will have an empty trace path and can be returned as is
	if denomTrace.Path == "" {
		denom = packet.Denom
	} else {
		// Non-native assets should be hashed
		denom = denomTrace.IBCDenom()
	}

	return denom
}

// Parse the denom from the Recv Packet that will be used by the rate limit module
// The denom that the rate limiter will use for a RECEIVE packet depends on whether it was a `source` or `sink`,
// explained here: https://github.com/cosmos/ibc-go/blob/04531d83bf/modules/apps/transfer/keeper/relay.go#L23-L54
//
//	If the chain is acting as a SINK: Add on the dYdX Chain port and channel and hash it
//	  Ex1: uusdc sent from Noble to dYdX
//	       Packet Denom:   uusdc
//	        -> Add Prefix: transfer/channel-0/uusdc
//	        -> Hash:       ibc/...
//
//	  Ex2: ujuno sent from Osmosis to dYdX Chain
//	       PacketDenom:    transfer/channel-Y/ujuno  (channel-Y is the Juno <> Osmosis channel)
//	        -> Add Prefix: transfer/channel-X/transfer/channel-Y/ujuno
//	        -> Hash:       ibc/...
//
//	If the chain is acting as a SOURCE: First, remove the prefix. Then if there is still a denom trace, hash it
//	  Ex1: adv4tnt sent back to dYdX chain from Osmosis
//	       Packet Denom:      transfer/channel-X/adv4tnt
//	        -> Remove Prefix: adv4tnt
//	        -> Leave as is:   adv4tnt
//
//	  Ex2: juno was sent to dYdX Chain, then to Osmosis, then back to dYdX Chain
//	       Packet Denom:      transfer/channel-X/transfer/channel-Z/ujuno
//	        -> Remove Prefix: transfer/channel-Z/ujuno
//	        -> Hash:          ibc/...
func ParseDenomFromRecvPacket(
	packet channeltypes.Packet,
	packetData ibctransfertypes.FungibleTokenPacketData,
) (denom string) {
	// To determine the denom, first check whether Stride is acting as source
	if ibctransfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), packetData.Denom) {
		// Remove the source prefix (e.g. transfer/channel-X/transfer/channel-Z/ujuno -> transfer/channel-Z/ujuno)
		sourcePrefix := ibctransfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := packetData.Denom[len(sourcePrefix):]

		// Native assets will have an empty trace path and can be returned as is
		denomTrace := ibctransfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path == "" {
			denom = unprefixedDenom
		} else {
			// Non-native assets should be hashed
			denom = denomTrace.IBCDenom()
		}
	} else {
		// Prefix the destination channel - this will contain the trailing slash (e.g. transfer/channel-X/)
		destinationPrefix := ibctransfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		prefixedDenom := destinationPrefix + packetData.Denom

		// Hash the denom trace
		denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)
		denom = denomTrace.IBCDenom()
	}

	return denom
}
