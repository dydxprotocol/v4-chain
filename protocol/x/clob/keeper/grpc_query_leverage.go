package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Leverage returns the leverage for a subaccount.
func (k Keeper) Leverage(
	c context.Context,
	req *types.QueryLeverageRequest,
) (*types.QueryLeverageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	// Get leverage for the subaccount
	leverage, exists := k.GetLeverage(ctx, &satypes.SubaccountId{
		Owner:  req.Owner,
		Number: req.Number,
	})
	if !exists {
		leverage = make(map[uint32]uint32)
	}

	clobPairLeverage := make([]*types.ClobPairLeverageInfo, 0, len(leverage))
	for clobPairId, leverage := range leverage {
		clobPairLeverage = append(clobPairLeverage, &types.ClobPairLeverageInfo{
			ClobPairId: clobPairId,
			Leverage:   leverage,
		})
	}

	return &types.QueryLeverageResponse{
		ClobPairLeverage: clobPairLeverage,
	}, nil
}
