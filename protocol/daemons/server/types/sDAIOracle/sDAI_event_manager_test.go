package types_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
	"github.com/stretchr/testify/require"
)

var (
	DefaultNextIndexInArray = 0
)

func setupEventManager() *sdaitypes.SDAIEventManager {
	return sdaitypes.NewsDAIEventManager()
}

func TestNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := setupEventManager()

	require.EqualValues(t, DefaultNextIndexInArray, sdaiEventManager.GetNextIndexInArray())
}

func TestSDAIEventManager_AddsDAIEvent(t *testing.T) {
	sdaiEventManager := setupEventManager()

	// Create a new event
	event := &api.AddsDAIEventsRequest{
		ConversionRate:      "100.0",
		EthereumBlockNumber: "123456",
	}

	// Add the event
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the event was added correctly
	lastEvents := sdaiEventManager.GetLastTensDAIEvents()
	require.EqualValues(t, *event, lastEvents[0])

	// Add more events to test the circular buffer
	for i := 1; i < 10; i++ {
		event := &api.AddsDAIEventsRequest{
			ConversionRate:      fmt.Sprintf("%d.0", 100+i),
			EthereumBlockNumber: fmt.Sprintf("%d", 123456+i),
		}
		require.NoError(t, sdaiEventManager.AddsDAIEvent(event))
	}

	// Check if the events were added correctly
	lastEvents = sdaiEventManager.GetLastTensDAIEvents()
	for i := 0; i < 10; i++ {
		expectedEvent := api.AddsDAIEventsRequest{
			ConversionRate:      fmt.Sprintf("%d.0", 100+i),
			EthereumBlockNumber: fmt.Sprintf("%d", 123456+i),
		}
		require.EqualValues(t, expectedEvent, lastEvents[i])
	}

	expectedLatestEvent := api.AddsDAIEventsRequest{
		ConversionRate:      fmt.Sprintf("%d.0", 100+9),
		EthereumBlockNumber: fmt.Sprintf("%d", 123456+9),
	}

	latest, ok := sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, expectedLatestEvent, latest)

	// Add one more event to test the circular buffer wrap-around
	event = &api.AddsDAIEventsRequest{
		ConversionRate:      "110.0",
		EthereumBlockNumber: "123466",
	}
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the events were added correctly
	lastEvents = sdaiEventManager.GetLastTensDAIEvents()

	latest, ok = sdaiEventManager.GetLatestsDAIEvent()
	require.True(t, ok)
	require.EqualValues(t, *event, latest)

	for i := 1; i < 10; i++ {
		expectedEvent := api.AddsDAIEventsRequest{
			ConversionRate:      fmt.Sprintf("%d.0", 100+i),
			EthereumBlockNumber: fmt.Sprintf("%d", 123456+i),
		}
		require.EqualValues(t, expectedEvent, lastEvents[i])
	}
	require.EqualValues(t, *event, lastEvents[0])

	require.EqualValues(t, 1, sdaiEventManager.GetNextIndexInArray())
}
