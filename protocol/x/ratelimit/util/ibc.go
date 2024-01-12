package util

// TODO(CORE-854): Improve attribution message.
// This file re-uses similar utilities (with minor tweaking) Stride's IBC Rate Limit implementation:
// https://github.com/Stride-Labs/stride/blob/121f2ac5d2e5f8e406f89999410a49ea4277a552/x/ratelimit

import (
	"encoding/json"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
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

func ParsePacketInfo(
	packet channeltypes.Packet,
	direction types.PacketDirection,
) (types.IBCTransferPacketInfo, error) {
	var packetData ibctransfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &packetData); err != nil {
		return types.IBCTransferPacketInfo{}, err
	}

	var channelID, denom string
	if direction == types.PACKET_SEND {
		channelID = packet.GetSourceChannel()
		denom = ParseDenomFromSendPacket(packetData)
	} else {
		channelID = packet.GetDestChannel()
		denom = ParseDenomFromRecvPacket(packet, packetData)
	}

	// From `SetString` documentation:
	// For base 0, the number prefix determines the actual base:
	// A prefix of “0b” or “0B” selects base 2, “0”, “0o” or “0O” selects base 8, and “0x” or “0X” selects base 16.
	// Otherwise, the selected base is 10 and no prefix is accepted.
	amount, ok := new(big.Int).SetString(packetData.Amount, 0)
	if !ok {
		return types.IBCTransferPacketInfo{},
			errorsmod.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"Unable to cast packet amount '%s' to sdkmath.Int",
				packetData.Amount,
			)
	}

	packetInfo := types.IBCTransferPacketInfo{
		ChannelID: channelID,
		Denom:     denom,
		Amount:    amount,
	}

	return packetInfo, nil
}

// UnpackAcknowledgementResponseForTransfer unmarshals Acknowledgements for IBC transfers, determines the status of the
// acknowledgement (success or failure), and, if applicable, assembles the message responses
func UnpackAcknowledgementResponseForTransfer(
	ctx sdk.Context,
	logger log.Logger,
	ack []byte,
) (*types.AcknowledgementResponse, error) {
	// Unmarshal the raw ack response
	var acknowledgement channeltypes.Acknowledgement
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(ack, &acknowledgement); err != nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal ICS-20 transfer packet acknowledgement: %s",
			err.Error(),
		)
	}

	// The ack can come back as either AcknowledgementResult or AcknowledgementError
	// If it comes back as AcknowledgementResult, the messages are encoded differently depending on the SDK version
	switch response := acknowledgement.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		if len(response.Result) == 0 {
			return nil, errorsmod.Wrapf(
				channeltypes.ErrInvalidAcknowledgement,
				"acknowledgement result cannot be empty",
			)
		}
		logger.Info(fmt.Sprintf("IBC transfer acknowledgement success: %+v", response))
		return &types.AcknowledgementResponse{Status: types.AckResponseStatus_SUCCESS}, nil

	case *channeltypes.Acknowledgement_Error:
		logger.Error(fmt.Sprintf("acknowledgement error: %s", response.Error))
		return &types.AcknowledgementResponse{Status: types.AckResponseStatus_FAILURE, Error: response.Error}, nil

	default:
		return nil, errorsmod.Wrapf(
			channeltypes.ErrInvalidAcknowledgement,
			"unsupported acknowledgement response field type %T",
			response,
		)
	}
}
