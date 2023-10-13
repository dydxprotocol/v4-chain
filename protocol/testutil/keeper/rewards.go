package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	rewardskeeper "github.com/dydxprotocol/v4-chain/protocol/x/rewards/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

func RewardsKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	rewardsKeeper *rewardskeeper.Keeper,
	feetiersKeeper *feetierskeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	assetsKeeper *assetskeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
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
		// Mock time provider response for market creation.
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		assetsKeeper, _ = createAssetsKeeper(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			transientStoreKey,
			true,
		)
		statsKeeper, _ := createStatsKeeper(
			stateStore,
			epochsKeeper,
			db,
			cdc,
		)
		feetiersKeeper, _ = createFeeTiersKeeper(
			stateStore,
			statsKeeper,
			db,
			cdc,
		)
		rewardsKeeper, storeKey = createRewardsKeeper(
			stateStore,
			assetsKeeper,
			bankKeeper,
			feetiersKeeper,
			pricesKeeper,
			db,
			cdc,
		)

		return []GenesisInitializer{
			pricesKeeper,
			assetsKeeper,
			feetiersKeeper,
			statsKeeper,
		}
	})
	return ctx, rewardsKeeper, feetiersKeeper, bankKeeper, assetsKeeper, pricesKeeper, storeKey
}

func createRewardsKeeper(
	stateStore storetypes.CommitMultiStore,
	assetsKeeper *assetskeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	feeTiersKeeper *feetierskeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
) (*rewardskeeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	authorities := []string{
		authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	}
	k := rewardskeeper.NewKeeper(
		cdc,
		storeKey,
		transientStoreKey,
		assetsKeeper,
		bankKeeper,
		feeTiersKeeper,
		pricesKeeper,
		authorities,
	)

	return k, storeKey
}
