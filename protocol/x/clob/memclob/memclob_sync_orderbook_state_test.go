package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// helper to create a minimal perpetual ClobPair with desired params.
func newPerpClobPair(id uint32, subticksPerTick uint32, stepBaseQuantums uint64) types.ClobPair {
	return types.ClobPair{
		Id:               id,
		SubticksPerTick:  subticksPerTick,
		StepBaseQuantums: stepBaseQuantums,
		Metadata: &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: id,
			},
		},
	}
}

func TestSyncOrderbookState_PanicsWhenOrderbookMissing(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	clobPair := newPerpClobPair(1, 100, 1000)
	require.Panics(t, func() {
		mem.SyncOrderbookState(clobPair)
	})
}

func TestSyncOrderbookState_PanicsWhenSubticksPerTickZero(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	update := newPerpClobPair(1, 0, 1000)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_PanicsWhenSubticksPerTickIncreased(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	// Increase from 100 -> 200 should panic.
	update := newPerpClobPair(1, 200, 1000)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_PanicsWhenSubticksPerTickNotDivisor(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	// 100 % 30 != 0 should panic even though decreased.
	update := newPerpClobPair(1, 30, 1000)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_PanicsWhenMinOrderBaseQuantumsZero(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	update := newPerpClobPair(1, 100, 0)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_PanicsWhenMinOrderBaseQuantumsIncreased(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	// Increase from 1000 -> 2000 should panic.
	update := newPerpClobPair(1, 100, 2000)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_PanicsWhenMinOrderBaseQuantumsNotDivisor(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	// 1000 % 300 != 0 should panic even though decreased.
	update := newPerpClobPair(1, 100, 300)
	require.Panics(t, func() {
		mem.SyncOrderbookState(update)
	})
}

func TestSyncOrderbookState_SucceedsAndUpdatesValues(t *testing.T) {
	mem := NewMemClobPriceTimePriority(false)
	initial := newPerpClobPair(1, 100, 1000)
	mem.CreateOrderbook(initial)

	// Valid decreases to positive divisors.
	update := newPerpClobPair(1, 50, 200)
	require.NotPanics(t, func() {
		mem.SyncOrderbookState(update)
	})

	ob := mem.orderbooks[update.GetClobPairId()]
	require.Equal(t, types.SubticksPerTick(50), ob.SubticksPerTick)
	require.Equal(t, ob.MinOrderBaseQuantums, update.GetClobPairMinOrderBaseQuantums())
}
