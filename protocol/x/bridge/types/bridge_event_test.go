package types_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestBridgeEvent_Equal(t *testing.T) {
	tests := map[string]struct {
		a   types.BridgeEvent
		b   types.BridgeEvent
		res bool
	}{
		"Equal": {
			a: types.BridgeEvent{
				Id:             1,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(17)),
				Address:        "address",
				EthBlockHeight: 128,
			},
			b: types.BridgeEvent{
				Id:             1,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(17)),
				Address:        "address",
				EthBlockHeight: 128,
			},
			res: true,
		},
		"Id not equal": {
			a: types.BridgeEvent{
				Id:             1,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(17)),
				Address:        "address",
				EthBlockHeight: 128,
			},
			b: types.BridgeEvent{
				Id:             2,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(17)),
				Address:        "address",
				EthBlockHeight: 128,
			},
			res: false,
		},
		"Coin denom not equal": {
			a: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(171)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			b: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test2", sdkmath.NewInt(171)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			res: false,
		},
		"Coin amount not equal": {
			a: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(171)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			b: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(1711)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			res: false,
		},
		"Address not equal": {
			a: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(171)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			b: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(1711)),
				Address:        "address1",
				EthBlockHeight: 1280,
			},
			res: false,
		},
		"Eth block height not equal": {
			a: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(171)),
				Address:        "address",
				EthBlockHeight: 1280,
			},
			b: types.BridgeEvent{
				Id:             10,
				Coin:           sdk.NewCoin("test", sdkmath.NewInt(1711)),
				Address:        "address",
				EthBlockHeight: 1281,
			},
			res: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.res, tc.a.Equal(tc.b))
		})
	}
}
