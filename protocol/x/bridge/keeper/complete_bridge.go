package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// `CompleteBridge` processes a bridge event by transferring the specified coin
// from bridge module account to the given address. The id of the bridge is not
// validated as it should have already been validated by AcknowledgeBridges.
func (k Keeper) CompleteBridge(
	ctx sdk.Context,
	bridge types.BridgeEvent,
) (err error) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.CompleteBridge,
		metrics.Latency,
	)

	// Do not complete bridge if bridging is disabled.
	safetyParams := k.GetSafetyParams(ctx)
	if safetyParams.IsDisabled {
		return types.ErrBridgingDisabled
	}

	// Convert bridge address string to sdk.AccAddress.
	bridgeAccAddress, err := sdk.AccAddressFromBech32(bridge.Address)
	if err != nil {
		return err
	}

	// If coin amount is positive, send coin from bridge module account to
	// specified account.
	if bridge.Coin.Amount.IsPositive() {
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			bridgeAccAddress,
			sdk.Coins{bridge.Coin},
		); err != nil {
			return err
		}
	}

	// Emit metric on last completed bridge id.
	telemetry.SetGauge(
		float32(bridge.Id),
		types.ModuleName,
		metrics.LastCompletedBridgeId,
	)

	return nil
}

// `GetDelayedCompleteBridgeMessages` returns all delayed complete bridge
// messages and corresponding block heights at which they'll execute.
// If `address` is given, only returns messages for that address.
func (k Keeper) GetDelayedCompleteBridgeMessages(
	ctx sdk.Context,
	address string,
) (
	messages []types.DelayedCompleteBridgeMessage,
) {
	// Get all delayed messages from `x/delaymsg`.
	allDelayedMessages := k.delayMsgKeeper.GetAllDelayedMessages(ctx)
	// Iterate through all delayed messages and find `MsgCompleteBridge`s.
	messages = make([]types.DelayedCompleteBridgeMessage, 0)
	for _, delayedMsg := range allDelayedMessages {
		sdkMsg, err := delayedMsg.GetMessage()
		if err != nil {
			continue
		}

		// If the message is a complete bridge message and its address matches `address` (if given),
		// add to the list of messages to return the message itself and the block height at which
		// it will execute.
		if completeBridgeMsg, ok := sdkMsg.(*types.MsgCompleteBridge); ok {
			if address == "" || completeBridgeMsg.Event.Address == address {
				messages = append(messages, types.DelayedCompleteBridgeMessage{
					Message:     *completeBridgeMsg,
					BlockHeight: delayedMsg.GetBlockHeight(),
				})
			}
		}
	}

	return messages
}
