package types

import (
	"sync"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/contract"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// events, err := getInitialEvents(3)
	// if err != nil {
	// 	log.Fatalf("Failed to get initial events: %v", err)
	// }

	return &SDAIEventManager{
		lastTenEvents:    [10]api.AddsDAIEventsRequest{},
		nextIndexInArray: 0,
	}
}

func getInitialEvents(numOfEvents int) ([10]api.AddsDAIEventsRequest, error) {

	// Initialize an Ethereum client from an RPC endpoint.
	time.Sleep(1 * time.Second)
	ethClient, err := ethclient.Dial(types.ETHRPC)
	if err != nil {
		return [10]api.AddsDAIEventsRequest{}, err
	}

	rates, blockNumbers, err := store.QueryDaiConversionRateForPastBlocks(ethClient, int64(numOfEvents), 3)
	if err != nil {
		return [10]api.AddsDAIEventsRequest{}, err
	}

	events := [10]api.AddsDAIEventsRequest{}

	for i := 0; i < numOfEvents; i++ {
		events[i] = api.AddsDAIEventsRequest{
			EthereumBlockNumber: blockNumbers[i],
			ConversionRate:      rates[i],
		}
	}

	ethClient.Close()

	return events, nil
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
