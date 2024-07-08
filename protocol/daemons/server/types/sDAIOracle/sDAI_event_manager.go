package types

import (
	"sync"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
)

// sDAIEventManager maintains an array of ethereum block height
// to sDAI conversion rate. Methods are goroutine safe.
type SDAIEventManager struct {
	// Exclusive mutex taken when reading or writing
	sync.Mutex

	// Array to store the last 10 Ethereum block heights and conversion rates
	lastTenEvents [10]api.AddsDAIEventsRequest

	// Index of the array where we should store the next api.AddsDAIEventsRequest
	nextIndexInArray int
}

// NewsDAIEventManager creates a new sDAIEventManager.
func NewsDAIEventManager() *SDAIEventManager {
	return &SDAIEventManager{
		lastTenEvents:    [10]api.AddsDAIEventsRequest{},
		nextIndexInArray: 0,
	}
}

func (s *SDAIEventManager) AddsDAIEvent(event *api.AddsDAIEventsRequest) error {
	s.Lock()
	defer s.Unlock()

	// Update the array with the new event
	s.lastTenEvents[s.nextIndexInArray] = *event

	// Move to the next index, wrapping around if necessary
	s.nextIndexInArray = (s.nextIndexInArray + 1) % 10

	return nil
}

// GetLastTensDAIEvents returns the last ten sDAI events.
func (s *SDAIEventManager) GetLastTensDAIEvents() [10]api.AddsDAIEventsRequest {
	s.Lock()
	defer s.Unlock()

	return s.lastTenEvents
}

// GetLatestsDAIEvent returns the most recent sDAI event.
func (s *SDAIEventManager) GetLatestsDAIEvent() (api.AddsDAIEventsRequest, bool) {
	s.Lock()
	defer s.Unlock()

	latestIndex := (s.nextIndexInArray - 1 + 10) % 10
	if s.lastTenEvents[latestIndex].EthereumBlockNumber != "" {
		return s.lastTenEvents[latestIndex], true
	}

	return api.AddsDAIEventsRequest{}, false
}

// GetNextIndexInArray returns the next index in the array.
func (s *SDAIEventManager) GetNextIndexInArray() int {
	s.Lock()
	defer s.Unlock()

	return s.nextIndexInArray
}
