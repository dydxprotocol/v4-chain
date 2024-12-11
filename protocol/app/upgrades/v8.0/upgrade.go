package v_8_0

import (
	"bytes"
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	accountplustypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// Migrate accountplus AccountState in kvstore from non-prefixed keys to prefixed keys
func MigrateAccountplusAccountState(ctx sdk.Context, k accountpluskeeper.Keeper) {
	ctx.Logger().Info("Migrating accountplus module AccountState in kvstore from non-prefixed keys to prefixed keys")

	store := ctx.KVStore(k.GetStoreKey())

	// Iterate on unprefixed keys
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	var keysToDelete [][]byte
	var accountStatesToSet []struct {
		address      sdk.AccAddress
		accountState accountplustypes.AccountState
	}
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

		accountStatesToSet = append(accountStatesToSet, struct {
			address      sdk.AccAddress
			accountState accountplustypes.AccountState
		}{key, accountState})

		keysToDelete = append(keysToDelete, key)
	}

	// Set prefixed keys
	for _, item := range accountStatesToSet {
		k.SetAccountState(ctx, item.address, item.accountState)
	}

	// Delete unprefixed keys
	for _, key := range keysToDelete {
		store.Delete(key)
	}

	ctx.Logger().Info("Successfully migrated accountplus AccountState keys")
}

const (
	ID_NUM = 300
)

// Set market, perpetual, and clob ids to a set number
// This is done so that the ids are consistent for convenience
func setMarketListingBaseIds(
	ctx sdk.Context,
	pricesKeeper pricestypes.PricesKeeper,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
) {
	// Set all ids to a set number
	pricesKeeper.SetNextMarketID(ctx, ID_NUM)

	perpetualsKeeper.SetNextPerpetualID(ctx, ID_NUM)

	clobKeeper.SetNextClobPairID(ctx, ID_NUM)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	pricesKeeper pricestypes.PricesKeeper,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
	accountplusKeeper accountpluskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		MigrateAccountplusAccountState(sdkCtx, accountplusKeeper)

		accountplusKeeper.SetActiveState(sdkCtx, true)

		// Set market, perpetual, and clob ids to a set number
		setMarketListingBaseIds(sdkCtx, pricesKeeper, perpetualsKeeper, clobKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
