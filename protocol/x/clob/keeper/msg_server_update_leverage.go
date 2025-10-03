package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateLeverage handles MsgUpdateLeverage
func (k msgServer) UpdateLeverage(
	goCtx context.Context,
	msg *types.MsgUpdateLeverage,
) (*types.MsgUpdateLeverageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	perpetualLeverageMap, err := types.ValidateAndConstructPerpetualLeverageMap(ctx, msg, k.Keeper)
	if err != nil {
		return nil, err
	}

	// Update leverage for the subaccount
	if err := k.Keeper.UpdateLeverage(ctx, msg.SubaccountId, perpetualLeverageMap); err != nil {
		return nil, err
	}

	return &types.MsgUpdateLeverageResponse{}, nil
}
