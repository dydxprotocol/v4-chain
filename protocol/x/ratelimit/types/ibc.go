package types

import "math/big"

type PacketDirection int32

const (
	PACKET_SEND PacketDirection = 0
	PACKET_RECV PacketDirection = 1
)

// IBCTransferPacketInfo contains relevant information from
// an IBC packet needed for rate-limiting.
type IBCTransferPacketInfo struct {
	ChannelID string
	Denom     string
	Amount    *big.Int
}

// AckResponseStatus represents the status of an acknowledgement of IBC transfer packet.
type AckResponseStatus int

const (
	AckResponseStatus_NOT_SPECIFIED AckResponseStatus = iota
	AckResponseStatus_SUCCESS
	AckResponseStatus_TIMEOUT
	AckResponseStatus_FAILURE
)

// AcknowledgementResponse contains information about an acknowledgement of IBC transfer packet.
type AcknowledgementResponse struct {
	Status AckResponseStatus
	Error  string
}
