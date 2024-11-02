package util_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestParseDenomFromRecvPacket(t *testing.T) {
	nobleChannelOnKlyra := "channel-0"
	nobleChannelOnOsmo := "channel-200"
	osmoChannelOnKlyra := "channel-5"
	klyraChannelOnNoble := "channel-100"
	klyraChannelOnOsmo := "channel-101"
	originalSDaiDenom := types.SDaiDenom

	testCases := []struct {
		name               string
		packetDenomTrace   string
		sourceChannel      string
		destinationChannel string
		expectedDenom      string
	}{
		// Sink asset one hop away:
		//   sDAI sent from Noble to klyra
		//   -> tack on prefix (transfer/channel-0/gsdai) and hash
		{
			name:               "sink_one_hop",
			packetDenomTrace:   types.SDaiDenom,
			sourceChannel:      klyraChannelOnNoble,
			destinationChannel: nobleChannelOnKlyra,
			expectedDenom: hashDenomTrace(fmt.Sprintf(
				"%s/%s/%s",
				transferPort,
				nobleChannelOnKlyra,
				types.SDaiDenom,
			)),
		},
		// Native source assets
		//    lib.DefaultBaseDenom sent from klyra to Noble and then back to klyra (transfer/channel-0/adv4tnt)
		//    -> remove prefix and leave as is (adv4tnt)
		{
			name:               lib.DefaultBaseDenom,
			packetDenomTrace:   fmt.Sprintf("%s/%s/%s", transferPort, klyraChannelOnNoble, lib.DefaultBaseDenom),
			sourceChannel:      klyraChannelOnNoble,
			destinationChannel: nobleChannelOnKlyra,
			expectedDenom:      lib.DefaultBaseDenom,
		},
		// Sink asset two hops away:
		//   gsdai sent from Noble to Osmosis to klyra (transfer/channel-200/gsdai)
		//   -> tack on prefix (transfer/channel-0/transfer/channel-200/gsdai) and hash
		{
			name:               "sink_two_hops",
			packetDenomTrace:   fmt.Sprintf("%s/%s/%s", transferPort, nobleChannelOnOsmo, originalSDaiDenom),
			sourceChannel:      klyraChannelOnOsmo,
			destinationChannel: osmoChannelOnKlyra,
			expectedDenom: hashDenomTrace(
				fmt.Sprintf(
					"%s/%s/%s/%s/%s",
					transferPort,
					osmoChannelOnKlyra,
					transferPort,
					nobleChannelOnOsmo,
					originalSDaiDenom,
				),
			),
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
			name:             "gsDAI on klyra",
			packetDenomTrace: "transfer/channel-0/gsdai",
			expectedDenom:    types.SDaiDenom,
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
