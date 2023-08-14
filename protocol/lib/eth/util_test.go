package eth_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestBridgeLogToEvent(t *testing.T) {
	tests := map[string]struct {
		inputLog   ethcoretypes.Log
		inputDenom string

		expectedEvent bridgetypes.BridgeEvent
	}{
		"Success: event ID 0": {
			inputLog:   constants.EthLog_Event0,
			inputDenom: "dv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 0,
				Coin: sdk.NewCoin(
					"dv4tnt",
					sdk.NewInt(42),
				),
				Address:        "dydx1qy352euf40x77qfrg4ncn27dauqjx3t83x4ummcpydzk0zdtehhse25p74",
				EthBlockHeight: 3872013,
			},
		},
		"Success: event ID 1": {
			inputLog:   constants.EthLog_Event1,
			inputDenom: "test-token",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 1,
				Coin: sdk.NewCoin(
					"test-token",
					sdk.NewInt(222),
				),
				Address:        "dydx1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsnpjqx",
				EthBlockHeight: 3969937,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			event := libeth.BridgeLogToEvent(tc.inputLog, tc.inputDenom)
			require.Equal(t, tc.expectedEvent, event)
		})
	}
}
