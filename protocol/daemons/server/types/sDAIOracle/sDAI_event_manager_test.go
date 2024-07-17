package types_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
	"github.com/stretchr/testify/require"
)

func setupEventManager() *sdaitypes.SDAIEventManager {
	sdaitypes.SDAIEventFetcher = &sdaitypes.MockEventFetcher{}
	return sdaitypes.NewsDAIEventManager()
}

func setupEventManagerWithNoEvents() *sdaitypes.SDAIEventManager {
	sdaitypes.SDAIEventFetcher = &sdaitypes.MockEventFetcherNoEvents{}
	return sdaitypes.NewsDAIEventManager()
}

func TestDefaultNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := setupEventManager()
	require.EqualValues(t, sdaitypes.FirstNextIndexInArray, sdaiEventManager.GetNextIndexInArray())
	actualEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()
	for i := 0; i < sdaitypes.InitialNumEvents; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], actualEvents[i])
	}
}

// TODO: This implementation is not optimal. The latestIndex is still a constant.
func TestSDAIEventManager_GetLatestsDAIEvent_Failure(t *testing.T) {
	sdaiEventManager := setupEventManagerWithNoEvents()
	lastEvent, found := sdaiEventManager.GetLatestsDAIEvent()
	require.False(t, found)
	require.EqualValues(t, api.AddsDAIEventsRequest{}, lastEvent)
}

func TestSDAIEventManager_GetLatestsDAIEvent_SuccessNoMod(t *testing.T) {
	sdaiEventManager := setupEventManager()
	lastEvent, found := sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, found)
	require.EqualValues(t, sdaitypes.TestSDAIEventRequests[sdaitypes.InitialNumEvents-1], lastEvent)
}

func TestSDAIEventManager_GetLatestsDAIEvent_SuccessMod(t *testing.T) {
	sdaiEventManager := setupEventManager()

	// Fill up entire array so that we loop back around in the array
	for i := 0; i < 10; i++ {
		event := sdaitypes.TestSDAIEventRequests[i]
		sdaiEventManager.AddsDAIEvent(&event)
	}

	lastEvent, found := sdaiEventManager.GetLatestsDAIEvent()

	require.True(t, found)
	require.EqualValues(t, sdaitypes.TestSDAIEventRequests[10-1], lastEvent)
}

func TestGetLastTensDAIEvents_NoWrapAround(t *testing.T) {
	sdaiEventManager := setupEventManager()
	for i := sdaitypes.FirstNextIndexInArray; i < 10; i++ {
		sdaiEventManager.AddsDAIEvent(&sdaitypes.TestSDAIEventRequests[i])
	}
	lastTenEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()

	for i := 0; i < 10; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], lastTenEvents[i])
	}
}

func TestSDAIEventManager_AddsDAIEvent(t *testing.T) {
	sdaiEventManager := setupEventManager()
	newEventIndex := sdaiEventManager.GetNextIndexInArray()

	// Create a new event
	event := &api.AddsDAIEventsRequest{
		ConversionRate:      sdaitypes.TestSDAIEventRequests[0].ConversionRate,
		EthereumBlockNumber: sdaitypes.TestSDAIEventRequests[0].EthereumBlockNumber,
	}

	// Add the event
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the event was added correctly
	lastEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()
	require.EqualValues(t, *event, lastEvents[newEventIndex])

	for i := 0; i < sdaitypes.InitialNumEvents; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], lastEvents[i])
	}

	// // Add more events to test the circular buffer
	for i := 0; i < 10; i++ {
		event := &api.AddsDAIEventsRequest{
			ConversionRate:      sdaitypes.TestSDAIEventRequests[i].ConversionRate,
			EthereumBlockNumber: sdaitypes.TestSDAIEventRequests[i].EthereumBlockNumber,
		}
		require.NoError(t, sdaiEventManager.AddsDAIEvent(event))
	}

	// Check if the events were added correctly
	lastEvents = sdaiEventManager.GetLastTensDAIEventsUnordered()
	sort.Slice(lastEvents[:], func(i, j int) bool {
		blockNumberI, _ := strconv.ParseInt(lastEvents[i].EthereumBlockNumber, 10, 64)
		blockNumberJ, _ := strconv.ParseInt(lastEvents[j].EthereumBlockNumber, 10, 64)
		return blockNumberI < blockNumberJ
	})

	for i := 0; i < 10; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], lastEvents[i])
	}

	expectedLatestEvent := sdaitypes.TestSDAIEventRequests[len(sdaitypes.TestSDAIEventRequests)-1]
	latest, ok := sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, expectedLatestEvent, latest)

	// Add one more event to test the circular buffer wrap-around
	event = &api.AddsDAIEventsRequest{
		ConversionRate:      "1106681181716810314385961731",
		EthereumBlockNumber: "12360",
	}
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the event were added correctly
	latest, ok = sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, *event, latest)

	lastEvents = sdaiEventManager.GetLastTensDAIEventsUnordered()
	sort.Slice(lastEvents[:], func(i, j int) bool {
		blockNumberI, _ := strconv.ParseInt(lastEvents[i].EthereumBlockNumber, 10, 64)
		blockNumberJ, _ := strconv.ParseInt(lastEvents[j].EthereumBlockNumber, 10, 64)
		return blockNumberI < blockNumberJ
	})

	for i := 0; i < 10; i++ {
		var expectedEvent api.AddsDAIEventsRequest
		if i < 9 {
			expectedEvent = sdaitypes.TestSDAIEventRequests[i+1]
		} else {
			expectedEvent = *event
		}
		require.EqualValues(t, expectedEvent, lastEvents[i])
	}

	require.EqualValues(t, 5, sdaiEventManager.GetNextIndexInArray())
}
