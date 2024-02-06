package util

// This file includes IBC utility methods used by the IBC middleware.
// Re-uses/adapts Stride x/ratelimit implementation: https://github.com/Stride-Labs/stride/tree/4913e1dd1a/x/ratelimit
// See v4-chain/protocol/x/ratelimit/LICENSE and v4-chain/protocol/x/ratelimit/README.md for licensing information.

import (
	"encoding/json"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

// Parses the following information from a packet needed for transfer rate limit:
// - channelID
// - denom
// - amount
//
// This function is similar to Stride's implementation below except it ignores the `Sender`
// and `Receiver` information.
// https://github.com/Stride-Labs/stride/blob/eb3564c7/x/ratelimit/keeper/packet.go#L127
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
				"Unable to cast packet amount '%s' to big.Int",
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
// acknowledgement (success or failure).
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
		logger.Info(
			"IBC transfer acknowledgement success",
			"response",
			response,
		)
		return &types.AcknowledgementResponse{Status: types.AckResponseStatus_SUCCESS}, nil

	case *channeltypes.Acknowledgement_Error:
		logger.Error(
			"received acknowledgement error",
			"error",
			response.Error,
		)
		return &types.AcknowledgementResponse{Status: types.AckResponseStatus_FAILURE, Error: response.Error}, nil

	default:
		return nil, errorsmod.Wrapf(
			channeltypes.ErrInvalidAcknowledgement,
			"unsupported acknowledgement response field type %T",
			response,
		)
	}
}
