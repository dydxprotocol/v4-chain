package finalizeblock

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	ante_types "github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// EventStager supports staging and retrieval of events (of type T) from FinalizeBlock.
type EventStager[T proto.Message] struct {
	transientStoreKey    storetypes.StoreKey
	cdc                  codec.BinaryCodec
	stagedEventCountKey  string
	stagedEventKeyPrefix string
}

// NewEventStager creates a new EventStager.
func NewEventStager[T proto.Message](
	transientStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	stagedEventCountKey string,
	stagedEventKeyPrefix string,
) EventStager[T] {
	return EventStager[T]{
		transientStoreKey:    transientStoreKey,
		cdc:                  cdc,
		stagedEventCountKey:  stagedEventCountKey,
		stagedEventKeyPrefix: stagedEventKeyPrefix,
	}
}

// GetStagedFinalizeBlockEvents retrieves all staged events from the store.
func (s EventStager[T]) GetStagedFinalizeBlockEvents(
	ctx sdk.Context,
	newStagedEvent func() T,
) []T {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	store := noGasCtx.TransientStore(s.transientStoreKey)

	count := s.getStagedEventsCount(store)
	events := make([]T, count)
	store = prefix.NewStore(store, []byte(s.stagedEventKeyPrefix))
	for i := uint32(0); i < count; i++ {
		event := newStagedEvent()
		bytes := store.Get(lib.Uint32ToKey(i))
		s.cdc.MustUnmarshal(bytes, event)
		events[i] = event
	}
	return events
}

func (s EventStager[T]) getStagedEventsCount(
	store storetypes.KVStore,
) uint32 {
	countsBytes := store.Get([]byte(s.stagedEventCountKey))
	if countsBytes == nil {
		return 0
	}
	return binary.BigEndian.Uint32(countsBytes)
}

// StageFinalizeBlockEvent stages an event in the transient store.
func (s EventStager[T]) StageFinalizeBlockEvent(
	ctx sdk.Context,
	stagedEvent T,
) {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	store := noGasCtx.TransientStore(s.transientStoreKey)

	// Increment events count.
	count := s.getStagedEventsCount(store)
	store.Set([]byte(s.stagedEventCountKey), lib.Uint32ToKey(count+1))

	// Store events keyed by index.
	store = prefix.NewStore(store, []byte(s.stagedEventKeyPrefix))
	store.Set(lib.Uint32ToKey(count), s.cdc.MustMarshal(stagedEvent))
}
