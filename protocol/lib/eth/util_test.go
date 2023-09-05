package eth_test

import (
	sdkmath "cosmossdk.io/math"
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
					sdkmath.NewInt(12345),
				),
				Address:        "dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextx92xhln",
				EthBlockHeight: 3872013,
			},
		},
		"Success: event ID 1 - empty address": {
			inputLog:   constants.EthLog_Event1,
			inputDenom: "test-token",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 1,
				Coin: sdk.NewCoin(
					"test-token",
					sdkmath.NewInt(55),
				),
				// address shorter than 20 bytes is padded with zeros.
				Address:        "dydx1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq66wm82",
				EthBlockHeight: 3969937,
			},
		},
		"Success: event ID 2": {
			inputLog:   constants.EthLog_Event2,
			inputDenom: "test-token",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 2,
				Coin: sdk.NewCoin(
					"test-token",
					sdkmath.NewInt(777),
				),
				// 32 bytes * 8 bits / 5 bits = 51.2 characters ~ 52 bech32 characters
				Address:        "dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextxzy3rx3z4vemc3xgq42as94fpcv",
				EthBlockHeight: 4139345,
			},
		},
		"Success: event ID 3": {
			inputLog:   constants.EthLog_Event3,
			inputDenom: "test-token-2",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 3,
				Coin: sdk.NewCoin(
					"test-token-2",
					sdkmath.NewInt(888),
				),
				// address data is 62 bytes but we take the first 32 bytes only.
				// 32 bytes * 8 bits / 5 bits ~ 52 bech32 characters
				Address:        "dydx124n92ej4ve2kv4tx24n92ej4ve2kv4tx24n92ej4ve2kv4tx24nq8exmjh",
				EthBlockHeight: 4139348,
			},
		},
		"Success: event ID 4": {
			inputLog:   constants.EthLog_Event4,
			inputDenom: "dv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 4,
				Coin: sdk.NewCoin(
					"dv4tnt",
					sdkmath.NewInt(1234123443214321),
				),
				// address shorter than 20 bytes is padded with zeros.
				Address:        "dydx1zg6pydqqqqqqqqqqqqqqqqqqqqqqqqqqm0r5ra",
				EthBlockHeight: 4139349,
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
