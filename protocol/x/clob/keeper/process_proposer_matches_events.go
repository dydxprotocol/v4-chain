package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetProcessProposerMatchesEvents gets the process proposer matches events from the latest block.
func (k Keeper) GetProcessProposerMatchesEvents(ctx sdk.Context) types.ProcessProposerMatchesEvents {
	// Retrieve an instance of the memory store.
	memStore := ctx.KVStore(k.memKey)

	// Retrieve the `processProposerMatchesEvents` bytes from the store.
	processProposerMatchesEventsBytes := memStore.Get(
		[]byte(types.ProcessProposerMatchesEventsKey),
	)

	// Unmarshal the `processProposerMatchesEvents` into a struct and return it.
	var processProposerMatchesEvents types.ProcessProposerMatchesEvents
	k.cdc.MustUnmarshal(processProposerMatchesEventsBytes, &processProposerMatchesEvents)
	return processProposerMatchesEvents
}

// MustSetProcessProposerMatchesEvents sets the process proposer matches events from the latest block.
// This function panics if:
//   - the current block height does not match the block height of the ProcessProposerMatchesEvents
//   - called outside of deliver TX mode
//   - Any of the ProcessProposerMatchesEvents fields have duplicates.
//
// TODO(DEC-1281): add parameter validation.
func (k Keeper) MustSetProcessProposerMatchesEvents(
	ctx sdk.Context,
	processProposerMatchesEvents types.ProcessProposerMatchesEvents,
) {
	lib.AssertDeliverTxMode(ctx)

	if err := processProposerMatchesEvents.ValidateProcessProposerMatchesEvents(ctx); err != nil {
		panic(err)
	}

	// Retrieve an instance of the memory store.
	memStore := ctx.KVStore(k.memKey)

	// Write `processProposerMatchesEvents` to the `memStore`.
	memStore.Set(
		[]byte(types.ProcessProposerMatchesEventsKey),
		k.cdc.MustMarshal(&processProposerMatchesEvents),
	)
}

// InitializeProcessProposerMatchesEvents initializes the process proposer matches events.
// This function should only be called from the CLOB genesis.
func (k Keeper) InitializeProcessProposerMatchesEvents(
	ctx sdk.Context,
) {
	processProposerMatchesEvents := types.ProcessProposerMatchesEvents{
		BlockHeight: 1,
	}

	memStore := ctx.KVStore(k.memKey)
	memStore.Set(
		[]byte(types.ProcessProposerMatchesEventsKey),
		k.cdc.MustMarshal(&processProposerMatchesEvents),
	)
}
