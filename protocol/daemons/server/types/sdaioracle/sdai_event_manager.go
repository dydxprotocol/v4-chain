package types

import (
	"log"
	"sync"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/contract"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	SDAIEventFetcher EventFetcher = &EthEventFetcher{}

	TestSDAIEventRequest = api.AddsDAIEventsRequest{
		ConversionRate: "1006681181716810314385961731",
	}
)

type MockEventFetcher struct{}

func (m *MockEventFetcher) GetInitialEvent(empty bool) (api.AddsDAIEventsRequest, error) {
	if empty {
		return api.AddsDAIEventsRequest{}, nil
	}
	return TestSDAIEventRequest, nil
}

type MockEventFetcherNoEvents struct{}

func (m *MockEventFetcherNoEvents) GetInitialEvent(empty bool) (api.AddsDAIEventsRequest, error) {
	return api.AddsDAIEventsRequest{}, nil
}

type EventFetcher interface {
	GetInitialEvent(empty bool) (api.AddsDAIEventsRequest, error)
}

type EthEventFetcher struct{}

func (r *EthEventFetcher) GetInitialEvent(empty bool) (api.AddsDAIEventsRequest, error) {

	if empty {
		return api.AddsDAIEventsRequest{}, nil
	}

	ethClient, err := ethclient.Dial(types.ETHRPC)
	if err != nil {
		return api.AddsDAIEventsRequest{}, err
	}

	rate, err := store.QueryDaiConversionRateWithRetries(ethClient, 3)
	if err != nil {
		return api.AddsDAIEventsRequest{}, err
	}

	ethClient.Close()

	return api.AddsDAIEventsRequest{ConversionRate: rate}, nil
}

// sDAIEventManager maintains an array of ethereum block height
// to sDAI conversion rate. Methods are goroutine safe.
type SDAIEventManager struct {
	// Exclusive mutex taken when reading or writing
	sync.Mutex

	price api.AddsDAIEventsRequest
}

// NewsDAIEventManager creates a new sDAIEventManager.
func NewsDAIEventManager(isEmpty ...bool) *SDAIEventManager {
	empty := false
	if len(isEmpty) > 0 && isEmpty[0] {
		empty = true
	}

	event, err := SDAIEventFetcher.GetInitialEvent(empty)
	if err != nil {
		log.Fatalf("Failed to get initial events: %v", err)
	}
	return &SDAIEventManager{
		price: event,
	}
}

func (s *SDAIEventManager) AddsDAIEvent(event *api.AddsDAIEventsRequest) error {
	s.Lock()
	defer s.Unlock()

	s.price = *event
	return nil
}

func (s *SDAIEventManager) GetSDaiPrice() api.AddsDAIEventsRequest {
	s.Lock()
	defer s.Unlock()

	return s.price
}

func SetupMockEventManager(isEmpty ...bool) *SDAIEventManager {
	SDAIEventFetcher = &MockEventFetcher{}

	if len(isEmpty) > 0 && isEmpty[0] {
		return NewsDAIEventManager(true)
	}
	return NewsDAIEventManager()
}

func SetupMockEventManagerWithNoEvents() *SDAIEventManager {
	SDAIEventFetcher = &MockEventFetcherNoEvents{}
	return NewsDAIEventManager()
}
