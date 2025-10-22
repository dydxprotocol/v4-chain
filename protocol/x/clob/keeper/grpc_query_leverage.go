package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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
	leverageMap, exists := k.subaccountsKeeper.GetLeverage(ctx, &satypes.SubaccountId{
		Owner:  req.Owner,
		Number: req.Number,
	})
	if !exists {
		leverageMap = make(map[uint32]uint32)
	}

	clobPairLeverage := make([]*types.ClobPairLeverageInfo, 0, len(leverageMap))
	for perpetualId, custom_imf_ppm := range leverageMap {
		clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
		if err != nil {
			return nil, status.Error(codes.Internal, errorsmod.Wrap(err, "failed to get clob pair id for perpetual").Error())
		}
		clobPairLeverage = append(clobPairLeverage, &types.ClobPairLeverageInfo{
			ClobPairId:   clobPairId.ToUint32(),
			CustomImfPpm: custom_imf_ppm,
		})
	}
	return &types.QueryLeverageResponse{ClobPairLeverage: clobPairLeverage}, nil
}
