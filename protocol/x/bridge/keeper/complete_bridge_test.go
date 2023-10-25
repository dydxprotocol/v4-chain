package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestCompleteBridge(t *testing.T) {
	tests := map[string]struct {
		// Initial balance of bridge module account.
		initialModAccBalance sdk.Coin
		// Bridge event to complete.
		bridgeEvent types.BridgeEvent
		// Whether bridging is disabled.
		bridgingDisabled bool

		// Expected error, if any.
		expectedError string
		// Expected balance of bridge module account after bridge completion.
		expectedModAccBalance sdk.Coin
	}{
		"Success": {
			initialModAccBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent:           constants.BridgeEvent_Id0_Height0,           // bridges 888 tokens.
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(112)), // 1000 - 888.
		},
		"Success: coin amount is 0": {
			initialModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:             7,
				Address:        constants.BobAccAddress.String(),
				Coin:           sdk.NewCoin("adv4tnt", sdkmath.ZeroInt()),
				EthBlockHeight: 3,
			},
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Success: coin amount is negative": {
			initialModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:      7,
				Address: constants.BobAccAddress.String(),
				Coin: sdk.Coin{
					Denom:  "adv4tnt",
					Amount: sdkmath.NewInt(-1),
				},
				EthBlockHeight: 3,
			},
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Failure: invalid address string": {
			initialModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:             4,
				Address:        "not an address string",
				Coin:           sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
				EthBlockHeight: 2,
			},
			expectedError:         "decoding bech32 failed",
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Failure: bridge module account has insufficient balance": {
			initialModAccBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(500)),
			bridgeEvent:           constants.BridgeEvent_Id0_Height0, // bridges 888 tokens.
			expectedError:         "insufficient funds",
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(500)),
		},
		"Failure: bridging is disabled": {
			initialModAccBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent:           constants.BridgeEvent_Id0_Height0, // bridges 888 tokens.
			bridgingDisabled:      true,
			expectedError:         types.ErrBridgingDisabled.Error(),
			expectedModAccBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)), // same as initial.
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize context and keeper.
			ctx, bridgeKeeper, _, _, _, bankKeeper, _ := keepertest.BridgeKeepers(t)
			err := bridgeKeeper.UpdateSafetyParams(ctx, types.SafetyParams{
				IsDisabled:  tc.bridgingDisabled,
				DelayBlocks: bridgeKeeper.GetSafetyParams(ctx).DelayBlocks,
			})
			require.NoError(t, err)
			// Fund bridge module account with enough balance.
			err = bankKeeper.MintCoins(
				ctx,
				types.ModuleName,
				sdk.NewCoins(tc.initialModAccBalance),
			)
			require.NoError(t, err)

			// Complete bridge.
			err = bridgeKeeper.CompleteBridge(ctx, tc.bridgeEvent)

			// Assert expectations.
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)

				// Assert that target account's balance of bridged token is as expected.
				balance := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(tc.bridgeEvent.Address),
					tc.bridgeEvent.Coin.Denom,
				)
				expectedBalance := sdk.NewCoin(tc.bridgeEvent.Coin.Denom, sdkmath.ZeroInt())
				if tc.bridgeEvent.Coin.Amount.IsPositive() {
					expectedBalance = tc.bridgeEvent.Coin
				}
				require.Equal(t, expectedBalance.Denom, balance.Denom)
				require.Equal(t, expectedBalance.Amount, balance.Amount)
			}
			// Assert that bridge module account's balance is as expected.
			modAccBalance := bankKeeper.GetBalance(
				ctx,
				types.ModuleAddress,
				tc.bridgeEvent.Coin.Denom,
			)
			require.Equal(t, tc.expectedModAccBalance.Denom, modAccBalance.Denom)
			require.Equal(t, tc.expectedModAccBalance.Amount, modAccBalance.Amount)
		})
	}
}
