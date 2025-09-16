package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateLeverage handles MsgUpdateLeverage
func (k msgServer) UpdateLeverage(goCtx context.Context, msg *types.MsgUpdateLeverage) (*types.MsgUpdateLeverageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert from LeverageEntry slice to map
	perpetualLeverageMap := make(map[uint32]uint32)
	for _, entry := range msg.PerpetualLeverage {
		perpetualLeverageMap[entry.PerpetualId] = entry.Leverage
	}

	// Update leverage for the subaccount (validation happens inside UpdateLeverage)
	if err := k.Keeper.UpdateLeverage(ctx, msg.SubaccountId, perpetualLeverageMap); err != nil {
		return nil, err
	}

	return &types.MsgUpdateLeverageResponse{}, nil
}
