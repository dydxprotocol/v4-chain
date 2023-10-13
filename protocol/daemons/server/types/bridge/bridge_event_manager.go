package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

type EventId = uint32

// BridgeEventManager maintains a map of "Recognized" Bridge Events.
// That is, events that have been finalized on Ethereum but are
// not yet in consensus on the V4 chain. Methods are goroutine safe.
type BridgeEventManager struct {
	// Exclusive mutex taken when reading or writing
	sync.Mutex

	// Bridge events by ID
	events map[EventId]BridgeEventWithTime

	// Stores:
	// - The next unused key in the bridges map (`NextId`)
	// - The block height of the last recognized event (`EthBlockHeight`)
	recognizedEventInfo types.BridgeEventInfo

	// Time provider than can mocked out if necessary
	timeProvider libtime.TimeProvider
}

// BridgeEventWithTime is a type that wraps BridgeEvent but also
// holds an additional timestamp.
type BridgeEventWithTime struct {
	event     types.BridgeEvent
	timestamp time.Time
}

// NewBridgeEventManager creates a new BridgeEventManager.
func NewBridgeEventManager(
	timeProvider libtime.TimeProvider,
) *BridgeEventManager {
	return &BridgeEventManager{
		events: make(map[uint32]BridgeEventWithTime),
		recognizedEventInfo: types.BridgeEventInfo{
			NextId:         0,
			EthBlockHeight: 0,
		},
		timeProvider: timeProvider,
	}
}

// AddBridgeEvents adds bridge events to the manager (with timestamps).
// Added events must have contiguous and in-order IDs.
// Any events with ID less than the `recognizedEventInfo.NextId` are ignored.
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
			telemetry.IncrCounter(1, metrics.BridgeServer, metrics.AddBridgeEvents, metrics.EventIdNotSequential)
			return fmt.Errorf("AddBridgeEvents: Events must be contiguous and in-order")
		}
	}

	now := b.timeProvider.Now()
	for _, event := range events {
		// Ignore stale events which may be the result of a race condition.
		if event.Id < b.recognizedEventInfo.NextId {
			telemetry.IncrCounter(1, metrics.BridgeServer, metrics.AddBridgeEvents, metrics.EventIdAlreadyRecognized)
			continue
		}

		// Update BridgeEventManager with the new event.
		b.events[event.Id] = BridgeEventWithTime{
			event:     event,
			timestamp: now,
		}
		// Update recognized event info of BridgeEventManager.
		b.recognizedEventInfo = types.BridgeEventInfo{
			NextId:         event.Id + 1,
			EthBlockHeight: event.EthBlockHeight,
		}
	}

	// Emit metrics on updated recognized event info.
	telemetry.SetGauge(
		float32(b.recognizedEventInfo.NextId),
		metrics.BridgeServer,
		metrics.RecognizedEventInfo,
		metrics.NextId,
	)
	telemetry.SetGauge(
		float32(b.recognizedEventInfo.EthBlockHeight),
		metrics.BridgeServer,
		metrics.RecognizedEventInfo,
		metrics.EthBlockHeight,
	)

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

// GetRecognizedEventInfo returns `recognizedEventInfo`.
func (b *BridgeEventManager) GetRecognizedEventInfo() types.BridgeEventInfo {
	b.Lock()
	defer b.Unlock()

	return b.recognizedEventInfo
}

// SetRecognizedEventInfo sets `recognizedEventInfo`. An error is returned
// and no update occurs if `NextId` or `EthBlockHeight` is lesser than its
// existing value.
func (b *BridgeEventManager) SetRecognizedEventInfo(
	eventInfo types.BridgeEventInfo,
) error {
	b.Lock()
	defer b.Unlock()

	if eventInfo.NextId < b.recognizedEventInfo.NextId {
		return fmt.Errorf("NextId cannot be set to a lower value")
	} else if eventInfo.EthBlockHeight < b.recognizedEventInfo.EthBlockHeight {
		return fmt.Errorf("EthBlockHeight cannot be set to a lower value")
	}

	b.recognizedEventInfo = eventInfo
	return nil
}

func (b *BridgeEventManager) GetNow() time.Time {
	return b.timeProvider.Now()
}
