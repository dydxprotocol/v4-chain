package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestCompleteBridge(t *testing.T) {
	tests := map[string]struct {
		// Bridge event to complete.
		bridgeEvent types.BridgeEvent

		// Expected error, if any.
		expectedError string
	}{
		"Success": {
			bridgeEvent: constants.BridgeEvent_Id0_Height0,
		},
		"Failure: coin amount is 0": {
			bridgeEvent: types.BridgeEvent{
				Id:             7,
				Address:        constants.BobAccAddress.String(),
				Coin:           sdk.NewCoin("dummy-coin", sdk.ZeroInt()),
				EthBlockHeight: 3,
			},
			expectedError: "invalid coin",
		},
		"Failure: invalid address string": {
			bridgeEvent: types.BridgeEvent{
				Id:             4,
				Address:        "not an address string",
				Coin:           sdk.NewCoin("dummy-coin", sdk.NewInt(1)),
				EthBlockHeight: 2,
			},
			expectedError: "decoding bech32 failed",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize context and keeper.
			ctx, bridgeKeeper, _, _, _, bankKeeper := keepertest.BridgeKeepers(t)

			// Complete bridge.
			err := bridgeKeeper.CompleteBridge(ctx, tc.bridgeEvent)

			// Assert expectations.
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)

				// Assert that account's balance of bridged token is as expected.
				balance := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(tc.bridgeEvent.Address),
					tc.bridgeEvent.Coin.Denom,
				)
				require.Equal(t, tc.bridgeEvent.Coin.Denom, balance.Denom)
				require.Equal(t, tc.bridgeEvent.Coin.Amount, balance.Amount)
			}
		})
	}
}
