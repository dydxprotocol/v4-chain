package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestValidateSubaccountEquityTierLimitForShortTermOrder_ResultHeldUntilCleared(t *testing.T) {
	ks, ctx := setupShortTermOrderKeeper(t, constants.Alice_Num0_10_000USD, []types.EquityTierLimit{
		{UsdTncRequired: dtypes.NewInt(0), Limit: 0},
		{UsdTncRequired: dtypes.NewInt(70_000_000_000), Limit: 1000}, // $70,000
	})
	order := restingBuyOrder(constants.Alice_Num0, 0)

	err := ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order)
	require.ErrorIs(t, err, types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit)

	ks.SubaccountsKeeper.SetSubaccount(ctx, constants.Alice_Num0_100_000USD)
	err = ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order)
	require.ErrorIs(t, err, types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit)

	ks.ClobKeeper.ClearShortTermEquityTierResults()
	require.NoError(t, ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order))

	ks.SubaccountsKeeper.SetSubaccount(ctx, constants.Alice_Num0_10_000USD)
	require.NoError(t, ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order))

	ks.ClobKeeper.ClearShortTermEquityTierResults()
	err = ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order)
	require.ErrorIs(t, err, types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit)
}
