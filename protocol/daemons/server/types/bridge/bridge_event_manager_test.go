package types_test

import (
	"testing"

	bdtypes "github.com/dydxprotocol/v4/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func setupEventManager() *bdtypes.BridgeEventManager {
	timeProvider := mocks.TimeProvider{}
	timeProvider.On("Now").Return(constants.TimeT)
	bem := bdtypes.NewBridgeEventManager(&timeProvider)

	return bem
}

func TestNewBridgeEventManager(t *testing.T) {
	bem := setupEventManager()

	require.EqualValues(t, 0, bem.GetNextRecognizedEventId())
}

func TestBridgeEventManager_SetNextRecognizedEventId(t *testing.T) {
	bem := setupEventManager()

	// Check default value.
	require.EqualValues(t, 0, bem.GetNextRecognizedEventId())

	// Increase by 1.
	require.NoError(t, bem.SetNextRecognizedEventId(1))
	require.EqualValues(t, 1, bem.GetNextRecognizedEventId())

	// Increase by more than 1.
	require.NoError(t, bem.SetNextRecognizedEventId(3))
	require.EqualValues(t, 3, bem.GetNextRecognizedEventId())

	// Keep the same
	require.NoError(t, bem.SetNextRecognizedEventId(3))
	require.EqualValues(t, 3, bem.GetNextRecognizedEventId())

	// Cannot decrease.
	require.ErrorContains(t, bem.SetNextRecognizedEventId(2), "nextRecognizedEventId cannot be set to a lower value")
}

func TestBridgeEventManager_AddBridgeEvents(t *testing.T) {
	tests := map[string]struct {
		initialNREI  uint32
		expectedNREI uint32
		events       []types.BridgeEvent
		errorMsg     string
	}{
		"Empty": {
			initialNREI:  0,
			expectedNREI: 0,
			events:       []types.BridgeEvent{},
		},
		"Single": {
			initialNREI:  constants.BridgeEvent_55.Id,
			expectedNREI: constants.BridgeEvent_55.Id + 1,
			events: []types.BridgeEvent{
				constants.BridgeEvent_55,
			},
		},
		"Multiple": {
			initialNREI:  0,
			expectedNREI: 2,
			events: []types.BridgeEvent{
				constants.BridgeEvent_0,
				constants.BridgeEvent_1,
			},
		},
		"Previous": {
			initialNREI:  1,
			expectedNREI: 2,
			events: []types.BridgeEvent{
				constants.BridgeEvent_0,
				constants.BridgeEvent_1,
			},
		},
		"Error Not Contiguous": {
			initialNREI: 0,
			events: []types.BridgeEvent{
				constants.BridgeEvent_0,
				constants.BridgeEvent_1,
				constants.BridgeEvent_55,
			},
			errorMsg: "contiguous",
		},
		"Error Skip": {
			initialNREI: 0,
			events: []types.BridgeEvent{
				constants.BridgeEvent_1,
			},
			errorMsg: "is greater than the Next Id",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			bem := setupEventManager()
			err := bem.SetNextRecognizedEventId(tc.initialNREI)
			require.NoError(t, err)

			// add the events
			err = bem.AddBridgeEvents(tc.events)

			// check for expected errors
			if tc.errorMsg != "" {
				require.ErrorContains(t, err, tc.errorMsg)
				return // no other checks needed
			}

			// ensure result is correct
			require.EqualValues(t, tc.expectedNREI, bem.GetNextRecognizedEventId())
			for _, event := range tc.events {
				_, _, found := bem.GetBridgeEventById(event.Id)
				if event.Id >= tc.initialNREI {
					require.True(t, found)
				} else {
					require.False(t, found)
				}
			}
		})
	}
}

func TestBridgeEventManager_GetBridgeEventById_Empty(t *testing.T) {
	bem := setupEventManager()

	_, _, found := bem.GetBridgeEventById(0)
	require.Equal(t, false, found)
}

func TestBridgeEventManager_GetBridgeEventById_Success(t *testing.T) {
	bem := setupEventManager()

	err := bem.AddBridgeEvents([]types.BridgeEvent{
		constants.BridgeEvent_0,
	})
	require.NoError(t, err)

	result, timestamp, found := bem.GetBridgeEventById(constants.BridgeEvent_0.Id)
	require.True(t, found)
	require.Equal(t, constants.BridgeEvent_0, result)
	require.Equal(t, constants.TimeT, timestamp)
}
