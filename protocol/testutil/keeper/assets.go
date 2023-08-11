package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/assets/keeper"
	"github.com/dydxprotocol/v4/x/assets/types"
	priceskeeper "github.com/dydxprotocol/v4/x/prices/keeper"
)

// CreateUsdcAsset creates USDC in the assets module for tests.
func CreateUsdcAsset(ctx sdk.Context, assetsKeeper *keeper.Keeper) error {
	_, err := assetsKeeper.CreateAsset(
		ctx,
		constants.Usdc.Symbol,
		constants.Usdc.Denom,
		constants.Usdc.DenomExponent,
		constants.Usdc.HasMarket,
		constants.Usdc.MarketId,
		constants.Usdc.AtomicResolution,
	)
	return err
}

func AssetsKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	accountKeeper *authkeeper.AccountKeeper,
	bankKeeper *bankkeeper.BaseKeeper,
	storeKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		pricesKeeper, _, _, _, _ = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		accountKeeper, _ = createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		keeper, storeKey = createAssetsKeeper(stateStore, db, cdc, pricesKeeper)

		return []GenesisInitializer{pricesKeeper, keeper}
	})

	return ctx, keeper, pricesKeeper, accountKeeper, bankKeeper, storeKey
}

func createAssetsKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		pk,
	)

	return k, storeKey
}
