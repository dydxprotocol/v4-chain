package util_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	"github.com/stretchr/testify/require"
)

const (
	transferPort  = "transfer"
	channelOnHost = "channel-1"
)

func hashDenomTrace(denomTrace string) string {
	trace32byte := sha256.Sum256([]byte(denomTrace))
	var traceTmByte tmbytes.HexBytes = trace32byte[:]
	return fmt.Sprintf("ibc/%s", traceTmByte)
}

func TestParsePacketInfo(t *testing.T) {
	sourceChannel := "channel-100"
	destinationChannel := "channel-200"
	denom := "denom"
	amountString := "100"
	amountInt := big.NewInt(100)
	sender := "cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg"
	receiver := "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5"

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

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	require.NoError(t, err)

	receiverAddr, err := sdk.AccAddressFromBech32(receiver)
	require.NoError(t, err)

	// Send 'denom' from channel-100 -> channel-200
	expectedSendPacketInfo := types.IBCTransferPacketInfo{
		ChannelID: sourceChannel,
		Denom:     denom,
		Amount:    amountInt,
		Sender:    senderAddr,
		Receiver:  receiverAddr,
	}

	actualSendPacketInfo, err := util.ParsePacketInfo(packet, types.PACKET_SEND)
	require.NoError(t, err, "no error expected when parsing send packet")
	require.Equal(t, expectedSendPacketInfo, actualSendPacketInfo, "send packet")

	// Receive 'denom' from channel-100 -> channel-200
	expectedRecvPacketInfo := types.IBCTransferPacketInfo{
		ChannelID: destinationChannel,
		Denom:     hashDenomTrace(fmt.Sprintf("transfer/%s/%s", destinationChannel, denom)),
		Amount:    amountInt,
		Sender:    senderAddr,
		Receiver:  receiverAddr,
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

func TestGetValidatedFungibleTokenPacketData(t *testing.T) {
	testCases := []struct {
		name           string
		packetData     ibctransfertypes.FungibleTokenPacketData
		expectedAmount *big.Int
		expectedSender sdk.AccAddress
		expectedRecv   sdk.AccAddress
		expectedError  string
	}{
		{
			name: "valid_packet_data",
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    "atom",
				Amount:   "100",
				Sender:   "cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg",
				Receiver: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			},
			expectedAmount: big.NewInt(100),
			expectedSender: sdk.MustAccAddressFromBech32("cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg"),
			expectedRecv:   sdk.MustAccAddressFromBech32("cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5"),
		},
		{
			name: "invalid_amount",
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    "atom",
				Amount:   "invalid",
				Sender:   "cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg",
				Receiver: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			},
			expectedError: "Unable to cast packet amount 'invalid' to big.Int",
		},
		{
			name: "invalid_sender",
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    "atom",
				Amount:   "100",
				Sender:   "invalid",
				Receiver: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			},
			expectedError: "Unable to convert sender address",
		},
		{
			name: "invalid_receiver",
			packetData: ibctransfertypes.FungibleTokenPacketData{
				Denom:    "atom",
				Amount:   "100",
				Sender:   "cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg",
				Receiver: "invalid",
			},
			expectedError: "Unable to convert receiver address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			amount, sender, receiver, err := util.GetValidatedFungibleTokenPacketData(tc.packetData)

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedAmount, amount)
				require.Equal(t, tc.expectedSender, sender)
				require.Equal(t, tc.expectedRecv, receiver)
			}
		})
	}
}
