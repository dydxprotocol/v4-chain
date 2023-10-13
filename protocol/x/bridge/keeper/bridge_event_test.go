package keeper_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGetBridgeEventFromServer(t *testing.T) {
	tests := map[string]struct {
		// Bridge event to add to server.
		bridgeEvent types.BridgeEvent
		// Bridge event ID to query.
		bridgeEventId uint32

		// Expected response.
		expectedEvent types.BridgeEvent
		expectedFound bool
	}{
		"Event found": {
			bridgeEvent:   constants.BridgeEvent_Id0_Height0,
			bridgeEventId: 0,
			expectedEvent: constants.BridgeEvent_Id0_Height0,
			expectedFound: true,
		},
		"Event not found": {
			bridgeEvent:   constants.BridgeEvent_Id0_Height0,
			bridgeEventId: 1,
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize context, keeper, and bridgeEventManager.
			ctx, bridgeKeeper, _, mockTimeProvider, bridgeEventManager, _, _ := keepertest.BridgeKeepers(t)
			mockTimeProvider.On("Now").Return(time.Now())
			err := bridgeEventManager.AddBridgeEvents([]types.BridgeEvent{tc.bridgeEvent})
			require.NoError(t, err)

			// Complete bridge.
			event, found := bridgeKeeper.GetBridgeEventFromServer(ctx, tc.bridgeEventId)

			// Assert expectations.
			require.Equal(t, tc.expectedEvent, event)
			require.Equal(t, tc.expectedFound, found)
		})
	}
}
