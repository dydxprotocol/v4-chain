package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RefreshAllVaultOrders refreshes all orders for all vaults by
// TODO(TRA-134)
// 1. Cancelling all existing orders.
// 2. Placing new orders.
func (k Keeper) RefreshAllVaultOrders(ctx sdk.Context) {
}
