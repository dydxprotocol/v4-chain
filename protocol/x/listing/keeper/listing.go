package keeper

import (
	"encoding/json"

	"golang.org/x/xerrors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

// Function to set hard cap on listed markets in module store
func (k Keeper) SetMarketsHardCap(ctx sdk.Context, hardCap uint32) error {
	store := ctx.KVStore(k.storeKey)
	b, err := json.Marshal(hardCap)
	if err != nil {
		return err
	}
	store.Set([]byte(types.HardCapForMarketsKey), b)
	return nil
}

// Function to get hard cap on listed markets from module store
func (k Keeper) GetMarketsHardCap(ctx sdk.Context) (hardCap uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.HardCapForMarketsKey))
	if b == nil {
		return 0, xerrors.Errorf("market listing hard cap not found")
	}
	err = json.Unmarshal(b, &hardCap)
	if err != nil {
		return 0, err
	}
	return hardCap, nil
}
