package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// SetLeverage stores leverage data for a subaccount.
// Deprecated: Use subaccountsKeeper.SetLeverage instead.
// TODO: should be unused
func (k Keeper) SetLeverage(ctx sdk.Context, subaccountId *satypes.SubaccountId, leverageMap map[uint32]uint32) {
	k.subaccountsKeeper.SetLeverage(ctx, subaccountId, leverageMap)
}

// GetLeverage retrieves leverage data for a subaccount.
// Deprecated: Use subaccountsKeeper.GetLeverage instead.
// TODO: move to subaccounts keeper
func (k Keeper) GetLeverage(ctx sdk.Context, subaccountId *satypes.SubaccountId) (map[uint32]uint32, bool) {
	return k.subaccountsKeeper.GetLeverage(ctx, subaccountId)
}

// UpdateLeverage updates leverage for specific perpetuals for a subaccount.
// Deprecated: Use subaccountsKeeper.UpdateLeverage instead.
func (k Keeper) UpdateLeverage(
	ctx sdk.Context,
	subaccountId *satypes.SubaccountId,
	perpetualLeverage map[uint32]uint32,
) error {
	return k.subaccountsKeeper.UpdateLeverage(ctx, subaccountId, perpetualLeverage)
}
