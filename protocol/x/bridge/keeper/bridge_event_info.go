package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

const (
	acknowledgedEventInfoKey = "AcknowledgedEventInfo"
)

// GetAcknowledgedEventInfo returns `AcknowledgedEventInfo` from state.
func (k Keeper) GetAcknowledgedEventInfo(
	ctx sdk.Context,
) (acknowledgedEventInfo types.BridgeEventInfo) {
	store := ctx.KVStore(k.storeKey)
	var rawBytes []byte = store.Get([]byte(acknowledgedEventInfoKey))

	k.cdc.MustUnmarshal(rawBytes, &acknowledgedEventInfo)
	return acknowledgedEventInfo
}

// SetAcknowledgedEventInfo sets `AcknowledgedEventInfo` in state.
func (k Keeper) SetAcknowledgedEventInfo(
	ctx sdk.Context,
	acknowledgedEventInfo types.BridgeEventInfo,
) error {
	if err := acknowledgedEventInfo.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&acknowledgedEventInfo)
	store.Set([]byte(acknowledgedEventInfoKey), b)

	return nil
}

// GetRecognizedEventInfo returns `RecognizedEventInfo` from `BridgeEventManager`.
// This has the next event id that has not yet been recognized by this node’s daemon.
// This also has the height of the highest Ethereum block from which a bridge event
// was recognized. These values are not in-consensus.
func (k Keeper) GetRecognizedEventInfo(
	ctx sdk.Context,
) (recognizedEventInfo types.BridgeEventInfo) {
	return k.bridgeEventManager.GetRecognizedEventInfo()
}
