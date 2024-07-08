package types_test

import (
	"testing"

	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

var (
	DefaultNextIndexInArray = 0

	DefaultBridgeEventInfo = types.BridgeEventInfo{
		NextId:         0,
		EthBlockHeight: 0,
	}
)

func setupEventManager() *sdaitypes.SDAIEventManager {
	return sdaitypes.NewsDAIEventManager()
}

func TestNewSDAIEventManager(t *testing.T) {
	sdaiEventManager := setupEventManager()

	require.EqualValues(t, DefaultNextIndexInArray, sdaiEventManager.GetNextIndexInArray())
}

func TestBridgeEventManager_SetRecognizedEventInfo(t *testing.T) {
	sdaiEventManager := setupEventManager()

	// Check default value.
	require.EqualValues(t, DefaultBridgeEventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Increase `NextId` by 1.
	eventInfo := types.BridgeEventInfo{
		NextId:         1,
		EthBlockHeight: 0,
	}
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Increase `NextId` by more than 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 0,
	}
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Keep `NextId` the same
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Cannot decrease `NextId`.
	eventInfo = types.BridgeEventInfo{
		NextId:         2,
		EthBlockHeight: 0,
	}
	require.ErrorContains(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo), "NextId cannot be set to a lower value")

	// Increase `EthBlockHeight` by 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 1,
	}
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Increase `EthBlockHeight` by more than 1.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 4,
	}
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())

	// Cannot decrease `EthBlockHeight`.
	eventInfo = types.BridgeEventInfo{
		NextId:         3,
		EthBlockHeight: 3,
	}
	require.ErrorContains(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo), "EthBlockHeight cannot be set to a lower value")

	// Increase `NextId` and `EthBlockHeight` at the same time.
	eventInfo = types.BridgeEventInfo{
		NextId:         5,
		EthBlockHeight: 5,
	}
	require.NoError(t, sdaiEventManager.SetRecognizedEventInfo(eventInfo))
	require.EqualValues(t, eventInfo, sdaiEventManager.GetRecognizedEventInfo())
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
			sdaiEventManager := setupEventManager()
			err := sdaiEventManager.SetRecognizedEventInfo(tc.initialREI)
			require.NoError(t, err)

			// add the events
			err = sdaiEventManager.AddBridgeEvents(tc.events)

			// check for expected errors
			if tc.errorMsg != "" {
				require.ErrorContains(t, err, tc.errorMsg)
				return // no other checks needed
			}

			// ensure result is correct
			require.EqualValues(t, tc.expectedREI, sdaiEventManager.GetRecognizedEventInfo())
			for _, event := range tc.events {
				_, _, found := sdaiEventManager.GetBridgeEventById(event.Id)
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
	sdaiEventManager := setupEventManager()

	_, _, found := sdaiEventManager.GetBridgeEventById(0)
	require.Equal(t, false, found)
}

func TestBridgeEventManager_GetBridgeEventById_Success(t *testing.T) {
	sdaiEventManager := setupEventManager()

	err := sdaiEventManager.AddBridgeEvents([]types.BridgeEvent{
		constants.BridgeEvent_Id0_Height0,
	})
	require.NoError(t, err)

	result, timestamp, found := sdaiEventManager.GetBridgeEventById(constants.BridgeEvent_Id0_Height0.Id)
	require.True(t, found)
	require.Equal(t, constants.BridgeEvent_Id0_Height0, result)
	require.Equal(t, constants.TimeT, timestamp)
}
