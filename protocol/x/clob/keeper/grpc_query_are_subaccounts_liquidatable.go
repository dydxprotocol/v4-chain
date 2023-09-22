package keeper

import (
	"context"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AreSubaccountsLiquidatable(
	c context.Context,
	req *types.AreSubaccountsLiquidatableRequest,
) (
	*types.AreSubaccountsLiquidatableResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	negativeTncCount := 0
	results := make([]types.AreSubaccountsLiquidatableResponse_Result, len(req.SubaccountIds))
	for i, subaccountId := range req.SubaccountIds {
		isLiquidatable, isNegativeTnc, err := k.IsLiquidatable(ctx, subaccountId)

		if isNegativeTnc {
			negativeTncCount++
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		results[i] = types.AreSubaccountsLiquidatableResponse_Result{
			SubaccountId:   subaccountId,
			IsLiquidatable: isLiquidatable,
		}
	}

	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.SubaccountsNegativeTnc, metrics.Count},
		float32(negativeTncCount),
		[]gometrics.Label{
			metrics.GetLabelForStringValue(
				metrics.Callback,
				metrics.AreSubaccountsLiquidatable,
			),
		},
	)

	return &types.AreSubaccountsLiquidatableResponse{
		Results: results,
	}, nil
}
