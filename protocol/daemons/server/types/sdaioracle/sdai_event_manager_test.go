package types_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/stretchr/testify/require"
)

func TestDefaultNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	require.EqualValues(t, sdaitypes.FirstNextIndexInArray, sdaiEventManager.GetNextIndexInArray())
	actualEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()
	for i := 0; i < sdaitypes.InitialNumEvents; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], actualEvents[i])
	}
}

func TestGetNextIndexInArray_Basic_HasInitialEvents(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	require.EqualValues(t, sdaitypes.INITIAL_EVENT_NUM, sdaiEventManager.GetNextIndexInArray())
}

func TestGetNextIndexInArray_Basic_NoInitialEvents(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()
	t.Cleanup(func() {
		sdaitypes.InitialNumEvents = sdaitypes.INITIAL_EVENT_NUM
		sdaitypes.FirstNextIndexInArray = sdaitypes.INITIAL_EVENT_NUM
	})
	require.EqualValues(t, sdaitypes.ZERO_EVENT_NUM, sdaiEventManager.GetNextIndexInArray())
}

func TestGetLastTensDAIEventsUnordered_Basic_HasInitialEvents(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	events := sdaiEventManager.GetLastTensDAIEventsUnordered()
	require.Equal(t, 10, len(events))
	fmt.Println("EVENTS ARE ", events)
	for i, event := range events {
		if i < sdaitypes.INITIAL_EVENT_NUM {
			require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i], event)
		} else {
			require.Empty(t, event)
		}
	}
}

func TestGetLastTensDAIEventsUnordered_Basic_NoInitialEvents(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()
	t.Cleanup(func() {
		sdaitypes.InitialNumEvents = sdaitypes.INITIAL_EVENT_NUM
		sdaitypes.FirstNextIndexInArray = sdaitypes.INITIAL_EVENT_NUM
	})
	events := sdaiEventManager.GetLastTensDAIEventsUnordered()
	require.Equal(t, 10, len(events))
	require.Empty(t, events)
}

// TODO: This implementation is not optimal. The latestIndex is still a constant.
func TestSDAIEventManager_GetLatestsDAIEvent_Failure(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()
	t.Cleanup(func() {
		sdaitypes.InitialNumEvents = sdaitypes.INITIAL_EVENT_NUM
		sdaitypes.FirstNextIndexInArray = sdaitypes.INITIAL_EVENT_NUM
	})
	lastEvent, found := sdaiEventManager.GetLatestsDAIEvent()
	require.False(t, found)
	require.EqualValues(t, api.AddsDAIEventsRequest{}, lastEvent)
}

func TestSDAIEventManager_GetLatestsDAIEvent_SuccessNoMod(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	lastEvent, found := sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, found)
	require.EqualValues(t, sdaitypes.TestSDAIEventRequests[sdaitypes.InitialNumEvents-1], lastEvent)
}

func TestSDAIEventManager_GetLatestsDAIEvent_SuccessMod(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()

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
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	for i := sdaitypes.FirstNextIndexInArray; i < 10; i++ {
		sdaiEventManager.AddsDAIEvent(&sdaitypes.TestSDAIEventRequests[i])
	}
	lastTenEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()
	require.Len(t, lastTenEvents, 10)

	for i := 0; i < 10; i++ {
		require.EqualValues(t, sdaitypes.TestSDAIEventRequests[i].ConversionRate, lastTenEvents[i].ConversionRate)
	}
}

func TestGetLastTensDAIEvents_WrapAround(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()

	insertCount := make(map[string]int)
	for i := 0; i < 25; i++ {
		event := &sdaitypes.TestSDAIEventRequests[i%10]
		sdaiEventManager.AddsDAIEvent(event)
		if i >= 15 {
			insertCount[event.ConversionRate]++
		}
	}

	lastTenEvents := sdaiEventManager.GetLastTensDAIEventsUnordered()
	require.Len(t, lastTenEvents, 10)

	returnedCount := make(map[string]int)
	for _, event := range lastTenEvents {
		returnedCount[event.ConversionRate]++
	}

	for i := 0; i < 10; i++ {
		rate := sdaitypes.TestSDAIEventRequests[i].ConversionRate
		require.Equal(t, insertCount[rate], returnedCount[rate],
			"Mismatch for ConversionRate %s", rate)
	}

	for rate, count := range returnedCount {
		require.Equal(t, insertCount[rate], count,
			"Mismatch for ConversionRate %s", rate)
	}
}

func TestSDAIEventManager_AddsDAIEvent(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
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
