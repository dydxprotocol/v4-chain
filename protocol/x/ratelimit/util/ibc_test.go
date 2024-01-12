package util_test

// TODO(CORE-854): Improve attribution message.
// This file re-uses similar utilities (with minor tweaking) Stride's IBC Rate Limit implementation:
// https://github.com/Stride-Labs/stride/blob/121f2ac5d2e5f8e406f89999410a49ea4277a552/x/ratelimit

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	"github.com/stretchr/testify/require"
)

const (
	transferPort  = "transfer"
	uusdc         = "uuusdc"
	channelOnHost = "channel-1"
)

func hashDenomTrace(denomTrace string) string {
	trace32byte := sha256.Sum256([]byte(denomTrace))
	var traceTmByte tmbytes.HexBytes = trace32byte[:]
	return fmt.Sprintf("ibc/%s", traceTmByte)
}

func TestParseDenomFromRecvPacket(t *testing.T) {
	nobleChannelOnDydx := "channel-0"
	dydxChannelOnNoble := "channel-100"

	testCases := []struct {
		name               string
		packetDenomTrace   string
		sourceChannel      string
		destinationChannel string
		expectedDenom      string
	}{
		// Sink asset one hop away:
		//   uusdc sent from Noble to dYdX
		//   -> tack on prefix (transfer/channel-0/uusdc) and hash
		{
			name:               "sink_one_hop",
			packetDenomTrace:   assettypes.UusdcDenom,
			sourceChannel:      dydxChannelOnNoble,
			destinationChannel: nobleChannelOnDydx,
			expectedDenom: hashDenomTrace(fmt.Sprintf(
				"%s/%s/%s",
				transferPort,
				nobleChannelOnDydx,
				assettypes.UusdcDenom,
			)),
		},
		// Native source assets
		//    lib.DefaultBaseDenom sent from dYdX to Noble and then back to dYdX (transfer/channel-0/adv4tnt)
		//    -> remove prefix and leave as is (adv4tnt)
		{
			name:               lib.DefaultBaseDenom,
			packetDenomTrace:   fmt.Sprintf("%s/%s/%s", transferPort, dydxChannelOnNoble, lib.DefaultBaseDenom),
			sourceChannel:      dydxChannelOnNoble,
			destinationChannel: nobleChannelOnDydx,
			expectedDenom:      lib.DefaultBaseDenom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packet := channeltypes.Packet{
				SourcePort:         transferPort,
				DestinationPort:    transferPort,
				SourceChannel:      tc.sourceChannel,
				DestinationChannel: tc.destinationChannel,
			}
			packetData := ibctransfertypes.FungibleTokenPacketData{
				Denom: tc.packetDenomTrace,
			}

			parsedDenom := util.ParseDenomFromRecvPacket(packet, packetData)
			require.Equal(t, tc.expectedDenom, parsedDenom, tc.name)
		})
	}
}

func TestParseDenomFromSendPacket(t *testing.T) {
	testCases := []struct {
		name             string
		packetDenomTrace string
		expectedDenom    string
	}{
		// Native assets stay as is
		{
			name:             "base denom",
			packetDenomTrace: lib.DefaultBaseDenom,
			expectedDenom:    lib.DefaultBaseDenom,
		},
		// Non-native assets are hashed
		{
			name:             "uusdc on dYdX",
			packetDenomTrace: "transfer/channel-0/uusdc",
			expectedDenom:    assettypes.UusdcDenom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packet := ibctransfertypes.FungibleTokenPacketData{
				Denom: tc.packetDenomTrace,
			}

			parsedDenom := util.ParseDenomFromSendPacket(packet)
			require.Equal(t, tc.expectedDenom, parsedDenom, tc.name)
		})
	}
}

func TestParsePacketInfo(t *testing.T) {
	sourceChannel := "channel-100"
	destinationChannel := "channel-200"
	denom := "denom"
	amountString := "100"
	amountInt := big.NewInt(100)
	sender := "sender"
	receiver := "receiver"

	packetData, err := json.Marshal(ibctransfertypes.FungibleTokenPacketData{
		Denom:    denom,
		Amount:   amountString,
		Sender:   sender,
		Receiver: receiver,
	})
	require.NoError(t, err)

	packet := channeltypes.Packet{
		SourcePort:         transferPort,
		SourceChannel:      sourceChannel,
		DestinationPort:    transferPort,
		DestinationChannel: destinationChannel,
		Data:               packetData,
	}

	// Send 'denom' from channel-100 -> channel-200
	expectedSendPacketInfo := types.IBCTransferPacketInfo{
		ChannelID: sourceChannel,
		Denom:     denom,
		Amount:    amountInt,
	}
	actualSendPacketInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	require.NoError(t, err, "no error expected when parsing send packet")
	require.Equal(t, expectedSendPacketInfo, actualSendPacketInfo, "send packet")

	// Receive 'denom' from channel-100 -> channel-200
	expectedRecvPacketInfo := types.IBCTransferPacketInfo{
		ChannelID: destinationChannel,
		Denom:     hashDenomTrace(fmt.Sprintf("transfer/%s/%s", destinationChannel, denom)),
		Amount:    amountInt,
	}
	actualRecvPacketInfo, err := util.ParsePacketInfo(packet, types.PACKET_RECV)
	require.NoError(t, err, "no error expected when parsing recv packet")
	require.Equal(t, expectedRecvPacketInfo, actualRecvPacketInfo, "recv packet")
}

func TestUnpackAcknowledgementResponseForTransfer(t *testing.T) {
	exampleAckError := errors.New("ABCI code: 1: error handling packet: see events for details")

	testCases := []struct {
		name                string
		ack                 channeltypes.Acknowledgement
		expectedStatus      types.AckResponseStatus
		expectedNumMessages int
		packetError         string
		functionError       string
	}{
		{
			name:           "ibc_transfer_success",
			ack:            channeltypes.NewResultAcknowledgement([]byte{1}),
			expectedStatus: types.AckResponseStatus_SUCCESS,
		},
		{
			name:           "ibc_transfer_failure",
			ack:            channeltypes.NewErrorAcknowledgement(exampleAckError),
			expectedStatus: types.AckResponseStatus_FAILURE,
			packetError:    exampleAckError.Error(),
		},
		{
			name:          "ack_unmarshal_failure",
			ack:           channeltypes.Acknowledgement{},
			functionError: "cannot unmarshal ICS-20 transfer packet acknowledgement",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// If the ack is not empty, marshal it
			var err error
			var ackBz []byte
			if !bytes.Equal(tc.ack.Acknowledgement(), []byte("{}")) {
				ackBz, err = ibctransfertypes.ModuleCdc.MarshalJSON(&tc.ack)
				require.NoError(t, err, "no error expected when marshalling ack")
			}

			// Call unpack ack response and check error
			ackResponse, err := util.UnpackAcknowledgementResponseForTransfer(sdk.Context{}, log.NewNopLogger(), ackBz)
			if tc.functionError != "" {
				require.ErrorContains(t,
					err,
					tc.functionError,
					"unpacking acknowledgement response should have resulted in a function error",
				)
				return
			}
			require.NoError(t, err, "no error expected when unpacking ack")

			// Confirm the response and error status
			require.Equal(t, tc.expectedStatus, ackResponse.Status, "Acknowledgement response status")
			require.Equal(t, tc.packetError, ackResponse.Error, "AcknowledgementError")
		})
	}
}
