package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

func (k Keeper) GetBridgeEventFromServer(ctx sdk.Context, id uint32) (event types.BridgeEvent, found bool) {
	event, _, found = k.bridgeEventManager.GetBridgeEventById(id)
	return event, found
}
