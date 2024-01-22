package util_test

import (
	"fmt"
	"testing"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"
	"github.com/stretchr/testify/require"
)

func TestParseDenomFromRecvPacket(t *testing.T) {
	nobleChannelOnDydx := "channel-0"
	nobleChannelOnOsmo := "channel-200"
	osmoChannelOnDydx := "channel-5"
	dydxChannelOnNoble := "channel-100"
	dydxChannelOnOsmo := "channel-101"
	originalUsdcDenom := "uusdc"

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
		// Sink asset two hops away:
		//   uusdc sent from Noble to Osmosis to dYdX (transfer/channel-200/uusdc)
		//   -> tack on prefix (transfer/channel-0/transfer/channel-200/uusdc) and hash
		{
			name:               "sink_two_hops",
			packetDenomTrace:   fmt.Sprintf("%s/%s/%s", transferPort, nobleChannelOnOsmo, originalUsdcDenom),
			sourceChannel:      dydxChannelOnOsmo,
			destinationChannel: osmoChannelOnDydx,
			expectedDenom: hashDenomTrace(
				fmt.Sprintf(
					"%s/%s/%s/%s/%s",
					transferPort,
					osmoChannelOnDydx,
					transferPort,
					nobleChannelOnOsmo,
					originalUsdcDenom,
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
