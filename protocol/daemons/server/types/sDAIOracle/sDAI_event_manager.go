package types

import (
	"log"
	"sync"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/contract"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	SDAIEventFetcher      EventFetcher = &EthEventFetcher{}
	InitialNumEvents                   = 3
	FirstNextIndexInArray              = InitialNumEvents
	TestSDAIEventRequests              = []api.AddsDAIEventsRequest{
		{
			ConversionRate: "1006681181716810314385961731",
		},
		{
			ConversionRate: "1016681181716810314385961731",
		},
		{
			ConversionRate: "1026681181716810314385961731",
		},
		{
			ConversionRate: "1036681181716810314385961731",
		},
		{
			ConversionRate: "1046681181716810314385961731",
		},
		{
			ConversionRate: "1056681181716810314385961731",
		},
		{
			ConversionRate: "1066681181716810314385961731",
		},
		{
			ConversionRate: "1076681181716810314385961731",
		},
		{
			ConversionRate: "1086681181716810314385961731",
		},
		{
			ConversionRate: "1096681181716810314385961731",
		},
	}
)

type MockEventFetcher struct{}

func (m *MockEventFetcher) GetInitialEvents(numOfEvents int) ([10]api.AddsDAIEventsRequest, error) {
	events := [10]api.AddsDAIEventsRequest{}
	for i := 0; i < numOfEvents; i++ {
		events[i] = TestSDAIEventRequests[i]
	}
	return events, nil
}

type MockEventFetcherNoEvents struct{}

func (m *MockEventFetcherNoEvents) GetInitialEvents(numOfEvents int) ([10]api.AddsDAIEventsRequest, error) {
	events := [10]api.AddsDAIEventsRequest{}
	return events, nil
}

type EventFetcher interface {
	GetInitialEvents(numOfEvents int) ([10]api.AddsDAIEventsRequest, error)
}

type EthEventFetcher struct{}

func (r *EthEventFetcher) GetInitialEvents(numOfEvents int) ([10]api.AddsDAIEventsRequest, error) {
	time.Sleep(1 * time.Second)
	ethClient, err := ethclient.Dial(types.ETHRPC)
	if err != nil {
		return [10]api.AddsDAIEventsRequest{}, err
	}

	rates, err := store.QueryDaiConversionRateForPastBlocks(ethClient, int64(numOfEvents), 3)
	if err != nil {
		return [10]api.AddsDAIEventsRequest{}, err
	}

	events := [10]api.AddsDAIEventsRequest{}

	for i := 0; i < numOfEvents; i++ {
		events[i] = api.AddsDAIEventsRequest{
			ConversionRate: rates[i],
		}
	}

	ethClient.Close()

	return events, nil
}

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

	events, err := SDAIEventFetcher.GetInitialEvents(InitialNumEvents)

	if err != nil {
		log.Fatalf("Failed to get initial events: %v", err)
	}

	return &SDAIEventManager{
		lastTenEvents:    events,
		nextIndexInArray: InitialNumEvents,
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
// TODO: This does not handle the circular buffer
func (s *SDAIEventManager) GetLastTensDAIEventsUnordered() [10]api.AddsDAIEventsRequest {
	s.Lock()
	defer s.Unlock()

	return s.lastTenEvents
}

// GetLatestsDAIEvent returns the most recent sDAI event.
func (s *SDAIEventManager) GetLatestsDAIEvent() (api.AddsDAIEventsRequest, bool) {
	s.Lock()
	defer s.Unlock()

	latestIndex := (s.nextIndexInArray - 1 + 10) % 10
	if s.lastTenEvents[latestIndex].ConversionRate != "" {
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
