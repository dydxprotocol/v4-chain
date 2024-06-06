package keeper

import (
	"encoding/json"

	"golang.org/x/xerrors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

// Function to set permissionless listing flag in module store
func (k Keeper) SetPermissionlessListingEnable(ctx sdk.Context, enable bool) error {
	store := ctx.KVStore(k.storeKey)
	b, err := json.Marshal(enable)
	if err != nil {
		return err
	}
	store.Set([]byte(types.PermissionlessListingEnableKey), b)
	return nil
}

// Function to check if permissionless listing is enabled
func (k Keeper) IsPermissionlessListingEnabled(ctx sdk.Context) (enabled bool, err error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.PermissionlessListingEnableKey))
	if b == nil {
		return false, xerrors.Errorf("permissionless listing enable key not found")
	}
	err = json.Unmarshal(b, &enabled)
	if err != nil {
		return false, err
	}
	return enabled, nil
}
