package keeper

import (
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// `CompleteBridge` processes a bridge event by transfer the appropriate tokens
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

	// Convert bridge address string to sdk.AccAddress.
	bridgeAccAddress, err := sdk.AccAddressFromBech32(bridge.Address)
	if err != nil {
		return err
	}

	// Send coin from bridge module account to specified account.
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		bridgeAccAddress,
		sdk.Coins{bridge.Coin},
	); err != nil {
		return err
	}

	// Emit metric on unbridged balance.
	unbridged := k.bankKeeper.GetBalance(
		ctx,
		sdk.MustAccAddressFromBech32(k.GetBridgeAuthority()),
		bridge.Coin.Denom,
	)
	telemetry.SetGaugeWithLabels(
		[]string{
			types.ModuleName,
			metrics.UnbridgedBalance,
		},
		float32(unbridged.Amount.Int64()),
		[]gometrics.Label{
			metrics.GetLabelForStringValue(metrics.BridgeTokenDenom, unbridged.Denom),
		},
	)

	// Emit metric on last completed bridge id.
	telemetry.SetGauge(
		float32(bridge.Id),
		types.ModuleName,
		metrics.LastCompletedBridgeId,
	)

	return nil
}
