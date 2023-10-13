package eth_test

import (
	"sync"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestGetBridgeEventAbi(t *testing.T) {
	results := make([]*abi.ABI, 0)
	mu := sync.Mutex{}
	var wg sync.WaitGroup

	// Get the ABI 200 times.
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := libeth.GetBridgeEventAbi()

			mu.Lock()
			defer mu.Unlock()
			results = append(results, r)
		}()
	}
	wg.Wait()

	// Call the function one more time.
	// Ensure that all the pointers are equal.
	expected := libeth.GetBridgeEventAbi()
	require.NotNil(t, expected)
	for _, r := range results {
		require.Same(t, expected, r)
	}
}

func TestPadOrTruncateAddress(t *testing.T) {
	tests := map[string]struct {
		address  []byte
		expected []byte
	}{
		"nil": {
			address:  nil,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"length 0": {
			address:  []byte{},
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"length < min": {
			address:  []byte{1, 2, 3, 4},
			expected: []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"length = min": {
			address:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		"length between min and max": {
			address:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
		},
		"length = max": {
			address: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32,
			},
			expected: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32,
			},
		},
		"length > max": {
			address: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32, 33, 34,
			},
			expected: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := libeth.PadOrTruncateAddress(tc.address)
			require.Equal(t, tc.expected, actual)
			require.GreaterOrEqual(t, len(actual), libeth.MinAddrLen)
			require.LessOrEqual(t, len(actual), libeth.MaxAddrLen)
		})
	}
}

func TestBridgeLogToEvent(t *testing.T) {
	tests := map[string]struct {
		inputLog   ethcoretypes.Log
		inputDenom string

		expectedEvent bridgetypes.BridgeEvent
	}{
		"Success: event ID 0": {
			inputLog:   constants.EthLog_Event0,
			inputDenom: "adv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 0,
				Coin: sdk.NewCoin(
					"adv4tnt",
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
			inputDenom: "adv4tnt",
			expectedEvent: bridgetypes.BridgeEvent{
				Id: 4,
				Coin: sdk.NewCoin(
					"adv4tnt",
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
