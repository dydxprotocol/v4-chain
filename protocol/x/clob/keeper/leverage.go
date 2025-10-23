package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// UpdateLeverage updates leverage for specific perpetuals for a subaccount.
func (k Keeper) UpdateLeverage(
	ctx sdk.Context,
	subaccountId *satypes.SubaccountId,
	perpetualLeverage map[uint32]uint32,
) error {
	return k.subaccountsKeeper.UpdateLeverage(ctx, subaccountId, perpetualLeverage)
}
