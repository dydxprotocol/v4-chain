package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

// BridgeEventManager maintains a map of "Recognized" Bridge Events.
// That is, events that have been finalized on Ethereum but are
// not yet in consensus on the V4 chain. Methods are goroutine safe.
type BridgeEventManager struct {
	// Exclusive mutex taken when reading or writing
	sync.Mutex

	// Bridge events by ID
	events map[uint32]BridgeEventWithTime

	// Next unused key in the bridges map
	nextRecognizedEventId uint32

	// Time provider than can mocked out if necessary
	timeProvider lib.TimeProvider
}

// BridgeEventWithTime is a type that wraps BridgeEvent but also
// holds an additional timestamp.
type BridgeEventWithTime struct {
	event     types.BridgeEvent
	timestamp time.Time
}

// NewBridgeEventManager creates a new BridgeEventManager.
func NewBridgeEventManager(
	timeProvider lib.TimeProvider,
) *BridgeEventManager {
	return &BridgeEventManager{
		events:                make(map[uint32]BridgeEventWithTime),
		nextRecognizedEventId: 0,
		timeProvider:          timeProvider,
	}
}

// AddBridgeEvents adds bridge events to the manager (with timestamps).
// Added events must have contiguous and in-order IDs.
// Any events with ID less than the `nextRecognizedEventId` are ignored.
func (b *BridgeEventManager) AddBridgeEvents(
	events []types.BridgeEvent,
) error {
	b.Lock()
	defer b.Unlock()

	// Ignore empty lists.
	if len(events) == 0 {
		return nil
	}

	// Validate events are contiguous and in-order.
	for i, event := range events {
		if event.Id != events[0].Id+uint32(i) {
			return fmt.Errorf("AddBridgeEvents: Events must be contiguous and in-order")
		}
	}

	// Event IDs cannot be skipped.
	if events[0].Id > b.nextRecognizedEventId {
		return fmt.Errorf(
			"AddBridgeEvents: Event ID %d is greater than the Next Id %d.",
			events[0].Id,
			b.nextRecognizedEventId,
		)
	}

	now := b.timeProvider.Now()
	for _, event := range events {
		// Ignore stale events which may be the result of a race condition.
		if event.Id < b.nextRecognizedEventId {
			continue
		}

		// Due to the above validation, the eventId should always be the next expected.
		if event.Id != b.nextRecognizedEventId {
			panic(fmt.Errorf(
				"Event ID %d does not match the Next Id %d",
				event.Id,
				b.nextRecognizedEventId,
			))
		}

		// Update the BridgeEventManager
		b.events[event.Id] = BridgeEventWithTime{
			event:     event,
			timestamp: now,
		}
		b.nextRecognizedEventId++
	}

	return nil
}

// GetBridgeEventById returns a bridge event by ID.
// Found is false if the manager does not have the event.
func (b *BridgeEventManager) GetBridgeEventById(
	id uint32,
) (
	event types.BridgeEvent,
	timestamp time.Time,
	found bool,
) {
	b.Lock()
	defer b.Unlock()

	// Find the event.
	eventWithTime, found := b.events[id]
	if !found {
		return event, timestamp, found // default values
	}

	return eventWithTime.event, eventWithTime.timestamp, true
}

// SetNextRecognizedEventId sets the nextRecognizedEventId. An error is returned
// and no update occurs if the new value is lesser than the existing value.
func (b *BridgeEventManager) SetNextRecognizedEventId(
	id uint32,
) error {
	b.Lock()
	defer b.Unlock()

	if id < b.nextRecognizedEventId {
		return fmt.Errorf("nextRecognizedEventId cannot be set to a lower value")
	}
	b.nextRecognizedEventId = id
	return nil
}

// GetNextRecognizedEventId returns the nextRecognizedEventId.
func (b *BridgeEventManager) GetNextRecognizedEventId() uint32 {
	b.Lock()
	defer b.Unlock()

	return b.nextRecognizedEventId
}
