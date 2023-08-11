package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

var (
	// Private
	coin = sdk.Coin{
		Denom:  "test-token",
		Amount: sdk.NewIntFromUint64(888),
	}

	// Public
	BridgeEvent_0 = types.BridgeEvent{
		Id:      0,
		Address: string(AliceAccAddress),
		Coin:    coin,
	}
	BridgeEvent_1 = types.BridgeEvent{
		Id:      1,
		Address: string(BobAccAddress),
		Coin:    coin,
	}
	BridgeEvent_55 = types.BridgeEvent{
		Id:      55,
		Address: string(CarlAccAddress),
		Coin:    coin,
	}
)
