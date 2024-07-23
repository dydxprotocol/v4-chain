package types_test

import (
	"fmt"
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
		ConversionRate: sdaitypes.TestSDAIEventRequests[0].ConversionRate,
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
			ConversionRate: sdaitypes.TestSDAIEventRequests[i].ConversionRate,
		}
		require.NoError(t, sdaiEventManager.AddsDAIEvent(event))
	}

	// Check if the events were added correctly
	lastEvents = sdaiEventManager.GetLastTensDAIEventsUnordered()
	// Check if the lastEvents array is a rotated version of TestSDAIEventRequests
	offset := -1
	for i := 0; i < 10; i++ {
		if lastEvents[0] == sdaitypes.TestSDAIEventRequests[i] {
			offset = i
			break
		}
	}
	require.NotEqual(t, -1, offset, "lastEvents is not a rotated version of TestSDAIEventRequests")

	for i := 0; i < 10; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[(i+offset)%10], lastEvents[i])
	}

	expectedLatestEvent := sdaitypes.TestSDAIEventRequests[len(sdaitypes.TestSDAIEventRequests)-1]
	latest, ok := sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, expectedLatestEvent, latest)

	// Add one more event to test the circular buffer wrap-around
	event = &api.AddsDAIEventsRequest{
		ConversionRate: "1106681181716810314385961731",
	}
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the event were added correctly
	latest, ok = sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, *event, latest)

	lastEvents = sdaiEventManager.GetLastTensDAIEventsUnordered()

	fmt.Println(lastEvents)
	fmt.Println(sdaitypes.TestSDAIEventRequests)
	// Check if the lastEvents array is a rotated version of TestSDAIEventRequests
	offset = -1
	for i := 0; i < 10; i++ {
		if lastEvents[0] == sdaitypes.TestSDAIEventRequests[i] {
			offset = i
			break
		}
	}
	require.NotEqual(t, -1, offset, "lastEvents is not a rotated version of TestSDAIEventRequests")

	for i := 0; i < 10; i++ {
		if i != (offset+10-2)%10 {
			require.EqualValues(t, sdaitypes.TestSDAIEventRequests[(i+offset)%10], lastEvents[i])
		} else {
			require.EqualValues(t, *event, lastEvents[i])
		}
	}

	require.EqualValues(t, 5, sdaiEventManager.GetNextIndexInArray())
}
