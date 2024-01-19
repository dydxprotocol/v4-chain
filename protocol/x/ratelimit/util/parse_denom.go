package util

import (
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// Parse the denom from the Send Packet that will be used by the rate limit module
// The denom that the rate limiter will use for a SEND packet depends on whether
// it was a NATIVE token (e.g. ustrd, stuatom, etc.) or NON-NATIVE token (e.g. ibc/...)...
//
// We can identify if the token is native or not by parsing the trace denom from the packet
// If the token is NATIVE, it will not have a prefix (e.g. ustrd),
// and if it is NON-NATIVE, it will have a prefix (e.g. transfer/channel-2/uosmo)
//
// For NATIVE denoms, return as is (e.g. ustrd)
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
// The denom that the rate limiter will use for a RECEIVE packet depends on whether it was a source or sink
//
//	Sink:   The token moves forward, to a chain different than its previous hop
//	        The new port and channel are APPENDED to the denom trace.
//	        (e.g. A -> B, B is a sink) (e.g. A -> B -> C, C is a sink)
//
//	Source: The token moves backwards (i.e. revisits the last chain it was sent from)
//	        The port and channel are REMOVED from the denom trace - undoing the last hop.
//	        (e.g. A -> B -> A, A is a source) (e.g. A -> B -> C -> B, B is a source)
//
//	If the chain is acting as a SINK: We add on the Stride port and channel and hash it
//	  Ex1: uosmo sent from Osmosis to Stride
//	       Packet Denom:   uosmo
//	        -> Add Prefix: transfer/channel-X/uosmo
//	        -> Hash:       ibc/...
//
//	  Ex2: ujuno sent from Osmosis to Stride
//	       PacketDenom:    transfer/channel-Y/ujuno  (channel-Y is the Juno <> Osmosis channel)
//	        -> Add Prefix: transfer/channel-X/transfer/channel-Y/ujuno
//	        -> Hash:       ibc/...
//
//	If the chain is acting as a SOURCE: First, remove the prefix. Then if there is still a denom trace, hash it
//	  Ex1: ustrd sent back to Stride from Osmosis
//	       Packet Denom:      transfer/channel-X/ustrd
//	        -> Remove Prefix: ustrd
//	        -> Leave as is:   ustrd
//
//	  Ex2: juno was sent to Stride, then to Osmosis, then back to Stride
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
