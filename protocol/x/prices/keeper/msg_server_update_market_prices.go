package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) UpdateMarketPrices(
	goCtx context.Context,
	msg *types.MsgUpdateMarketPrices,
) (*types.MsgUpdateMarketPricesResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	// Validate.
	// Note that non-deterministic validation is skipped, because the prices have been deemed
	// valid w/r/t index prices in `ProcessProposal` in order for the msg to reach this step.
	if err := k.Keeper.PerformStatefulPriceUpdateValidation(ctx, msg, false); err != nil {
		log.ErrorLogWithError(
			ctx,
			"PerformStatefulPriceUpdateValidation failed",
			err,
		)
		panic(err)
	}

	// Update state.
	if err := k.Keeper.UpdateMarketPrices(ctx, msg.MarketPriceUpdates); err != nil {
		// This should never happen, because the updates were validated above.
		log.ErrorLogWithError(
			ctx,
			"UpdateMarketPrices failed",
			err,
		)
		panic(err)
	}

	telemetry.IncrCounter(1, types.ModuleName, metrics.UpdateMarketPrices, metrics.Success)
	return &types.MsgUpdateMarketPricesResponse{}, nil
}
