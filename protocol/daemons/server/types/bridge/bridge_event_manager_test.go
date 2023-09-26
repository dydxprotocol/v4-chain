package types_test

import (
	"testing"

	bdtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

var (
	DefaultBridgeEventInfo = types.BridgeEventInfo{
		NextId:         0,
		EthBlockHeight: 0,
	}
)

func setupEventManager() *bdtypes.BridgeEventManager {
	timeProvider := mocks.TimeProvider{}
	timeProvider.On("Now").Return(constants.TimeT)
	bem := bdtypes.NewBridgeEventManager(&timeProvider)

	return bem
}

func TestNewBridgeEventManager(t *testing.T) {
	bem := setupEventManager()

	require.EqualValues(t, DefaultBridgeEventInfo, bem.GetRecognizedEventInfo())
}

func TestBridgeEventManager_SetRecognizedEventInfo(t *testing.T) {
	bem := setupEventManager()

	// Check default value.
	require.EqualValues(t, DefaultBridgeEventInfo, bem.GetRecognizedEventInfo())

	// Increase `NextId` by 1.
	eventInfo := types.BridgeEventInfo{
		NextId:         1,
		EthBlockHeight: 0,
	}
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())

	// Increase `NextId` by more than 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 0,
	}
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())

	// Keep `NextId` the same
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())

	// Cannot decrease `NextId`.
	eventInfo = types.BridgeEventInfo{
		NextId:         2,
		EthBlockHeight: 0,
	}
	require.ErrorContains(t, bem.SetRecognizedEventInfo(eventInfo), "NextId cannot be set to a lower value")

	// Increase `EthBlockHeight` by 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 1,
	}
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())

	// Increase `EthBlockHeight` by more than 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 4,
	}
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())

	// Cannot decrease `EthBlockHeight`.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 3,
	}
	require.ErrorContains(t, bem.SetRecognizedEventInfo(eventInfo), "EthBlockHeight cannot be set to a lower value")

	// Increase `NextId` and `EthBlockHeight` at the same time.
	eventInfo = types.BridgeEventInfo{
		NextId:         5,
		EthBlockHeight: 5,
	}
	require.NoError(t, bem.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, bem.GetRecognizedEventInfo())
}

func TestBridgeEventManager_AddBridgeEvents(t *testing.T) {
	tests := map[string]struct {
		initialREI  types.BridgeEventInfo
		expectedREI types.BridgeEventInfo
		events      []types.BridgeEvent
		errorMsg    string
	}{
		"Empty": {
			initialREI: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
			expectedREI: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
			events: []types.BridgeEvent{},
		},
		"Single": {
			initialREI: types.BridgeEventInfo{
				NextId:         constants.BridgeEvent_Id55_Height15.Id,
				EthBlockHeight: 0,
			},
			expectedREI: types.BridgeEventInfo{
				NextId:         constants.BridgeEvent_Id55_Height15.Id + 1,
				EthBlockHeight: constants.BridgeEvent_Id55_Height15.EthBlockHeight,
			},
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id55_Height15,
			},
		},
		"Multiple": {
			initialREI: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
			expectedREI: types.BridgeEventInfo{
				NextId:         2,
				EthBlockHeight: constants.BridgeEvent_Id1_Height0.EthBlockHeight,
			},
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
		},
		"Previous": {
			initialREI: types.BridgeEventInfo{
				NextId:         1,
				EthBlockHeight: 0,
			},
			expectedREI: types.BridgeEventInfo{
				NextId:         2,
				EthBlockHeight: constants.BridgeEvent_Id1_Height0.EthBlockHeight,
			},
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
		},
		"Error Not Contiguous": {
			initialREI: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id55_Height15,
			},
			errorMsg: "contiguous",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			bem := setupEventManager()
			err := bem.SetRecognizedEventInfo(tc.initialREI)
			require.NoError(t, err)

			// add the events
			err = bem.AddBridgeEvents(tc.events)

			// check for expected errors
			if tc.errorMsg != "" {
				require.ErrorContains(t, err, tc.errorMsg)
				return // no other checks needed
			}

			// ensure result is correct
			require.EqualValues(t, tc.expectedREI, bem.GetRecognizedEventInfo())
			for _, event := range tc.events {
				_, _, found := bem.GetBridgeEventById(event.Id)
				if event.Id >= tc.initialREI.NextId {
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
		constants.BridgeEvent_Id0_Height0,
	})
	require.NoError(t, err)

	result, timestamp, found := bem.GetBridgeEventById(constants.BridgeEvent_Id0_Height0.Id)
	require.True(t, found)
	require.Equal(t, constants.BridgeEvent_Id0_Height0, result)
	require.Equal(t, constants.TimeT, timestamp)
}
