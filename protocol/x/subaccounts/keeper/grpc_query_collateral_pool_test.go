package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestQueryCollateralPoolAddress(t *testing.T) {
	for testName, tc := range map[string]struct {
		// Parameters
		request *types.QueryCollateralPoolAddressRequest

		// Expectations
		response *types.QueryCollateralPoolAddressResponse
		err      error
	}{
		"Nil request results in error": {
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Cross perpetual": {
			request: &types.QueryCollateralPoolAddressRequest{
				PerpetualId: constants.BtcUsd_NoMarginRequirement.Params.Id,
			},
			response: &types.QueryCollateralPoolAddressResponse{
				CollateralPoolAddress: types.ModuleAddress.String(),
			},
		},
		"Isolated perpetual": {
			request: &types.QueryCollateralPoolAddressRequest{
				PerpetualId: constants.IsoUsd_IsolatedMarket.Params.Id,
			},
			response: &types.QueryCollateralPoolAddressResponse{
				CollateralPoolAddress: constants.IsoCollateralPoolAddress.String(),
			},
		},
		"Perpetual not found": {
			request: &types.QueryCollateralPoolAddressRequest{
				PerpetualId: uint32(1000),
			},
			err: status.Error(codes.NotFound, fmt.Sprintf(
				"Perpetual id %+v not found.",
				uint32(1000),
			)),
		},
	} {
		t.Run(testName, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)
			response, err := keeper.CollateralPoolAddress(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
