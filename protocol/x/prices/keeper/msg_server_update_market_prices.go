package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) UpdateMarketPrices(
	goCtx context.Context,
	msg *types.MsgUpdateMarketPrices,
) (*types.MsgUpdateMarketPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate.
	// Note that non-deterministic validation is skipped, because the prices have been deemed
	// valid w/r/t index prices in `ProcessProposal` in order for the msg to reach this step.
	if err := k.Keeper.PerformStatefulPriceUpdateValidation(ctx, msg, false); err != nil {
		errMsg := fmt.Sprintf("PerformStatefulPriceUpdateValidation failed, err = %v", err)
		k.Keeper.Logger(ctx).Error(errMsg)
		panic(err)
	}

	// Update state.
	if err := k.Keeper.UpdateMarketPrices(ctx, msg.MarketPriceUpdates); err != nil {
		// This should never happen, because the updates were validated above.
		errMsg := fmt.Sprintf("UpdateMarketPrices failed, err = %v", err)
		k.Keeper.Logger(ctx).Error(errMsg)
		panic(err)
	}

	telemetry.IncrCounter(1, types.ModuleName, metrics.UpdateMarketPrices, metrics.Success)
	return &types.MsgUpdateMarketPricesResponse{}, nil
}
