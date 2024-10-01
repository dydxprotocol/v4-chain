package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/stretchr/testify/require"
)

func TestDefaultNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()
	actualEvents := sdaiEventManager.GetSDaiPrice()
	require.EqualValues(t, sdaitypes.TestSDAIEventRequest, actualEvents)
}

func TestEmptyNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager(true)
	actualEvents := sdaiEventManager.GetSDaiPrice()
	require.EqualValues(t, api.AddsDAIEventsRequest{}, actualEvents)
}

func TestSDAIEventManager_AddsDAIEvent(t *testing.T) {
	sdaiEventManager := sdaitypes.SetupMockEventManager()

	// Create a new event
	event := &api.AddsDAIEventsRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	}

	// Add the event
	require.NoError(t, sdaiEventManager.AddsDAIEvent(event))

	// Check if the event was added correctly
	actualEvents := sdaiEventManager.GetSDaiPrice()
	require.EqualValues(t, *event, actualEvents)
}
