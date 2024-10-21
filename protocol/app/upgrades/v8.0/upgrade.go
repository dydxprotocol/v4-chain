package v_8_0

import (
	"bytes"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	accountplustypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// Migrate accountplus AccountState in kvstore from non-prefixed keys to prefixed keys
func migrateAccountplusAccountState(ctx sdk.Context, k accountpluskeeper.Keeper) {
	ctx.Logger().Info("Migrating accountplus module AccountState in kvstore from non-prefixed keys to prefixed keys")

	store := ctx.KVStore(k.GetStoreKey())

	// Iterate on unprefixed keys
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		// Double check that key does not have prefix
		if bytes.HasPrefix(key, []byte(accountplustypes.AccountStateKeyPrefix)) {
			panic(fmt.Sprintf("unexpected key with prefix %X found during migration", accountplustypes.AccountStateKeyPrefix))
		}

		value := iterator.Value()
		var accountState accountplustypes.AccountState
		if err := k.GetCdc().Unmarshal(value, &accountState); err != nil {
			panic(fmt.Sprintf("failed to unmarshal AccountState for key %X: %s", key, err))
		}

		// SetAccountState stores with prefix
		k.SetAccountState(ctx, key, accountState)

		// Delete unprefixed key
		store.Delete(key)
	}

	ctx.Logger().Info("Successfully migrated accountplus AccountState keys")
}

// TODO: Scaffolding for upgrade: https://linear.app/dydx/issue/OTE-886/v8-upgrade-handler-scaffold
