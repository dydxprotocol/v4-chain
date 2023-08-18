package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MarkNetCollateralDecreasedForSubaccount stores whether the net collateral for a subaccount has decreased.
//
// This is used by equity tier limit validation to ensure that all subaccounts conform to their current equity
// tier limit.
func (k Keeper) MarkNetCollateralDecreasedForSubaccount(ctx sdk.Context, subaccountId types.SubaccountId) {
	store := prefix.NewStore(
		ctx.TransientStore(k.transientStoreKey),
		types.KeyPrefix(types.SubaccountWithDecreasedNetCollateralKeyPrefix),
	)
	key, err := subaccountId.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, []byte{1})
}

// GetAllSubaccountsWithDecreasedNetCollateral returns all subaccounts that have decreased their net collateral
// in the current block.
//
// This is used by equity tier limit validation to ensure that all subaccounts conform to their current equity
// tier limit.
func (k Keeper) GetAllSubaccountsWithDecreasedNetCollateral(ctx sdk.Context) []types.SubaccountId {
	store := prefix.NewStore(
		ctx.TransientStore(k.transientStoreKey),
		types.KeyPrefix(types.SubaccountWithDecreasedNetCollateralKeyPrefix),
	)

	subAccounts := make([]types.SubaccountId, 10)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		var val types.SubaccountId
		key := itr.Key()
		k.cdc.MustUnmarshal(key[:len(key)-len([]byte("/"))], &val)
		subAccounts = append(subAccounts, val)
	}
	return subAccounts
}
