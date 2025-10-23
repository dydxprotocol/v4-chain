package keeper

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
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
	indexerEventManager indexer_manager.IndexerEventManager,
	storeKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)

		accountsKeeper, _ := createAccountKeeper(
			stateStore,
			db,
			cdc,
			registry)
		bankKeeper, _ := createBankKeeper(stateStore, db, cdc, accountsKeeper)
		stakingKeeper, _ := createStakingKeeper(
			stateStore,
			db,
			cdc,
			accountsKeeper,
			bankKeeper,
		)
		statsKeeper, _ := createStatsKeeper(
			stateStore,
			epochsKeeper,
			db,
			cdc,
			stakingKeeper,
		)
		affiliatesKeeper, _ := createAffiliatesKeeper(stateStore, db, cdc, statsKeeper, transientStoreKey, true)
		vaultKeeper, _ := createVaultKeeper(
			stateStore,
			db,
			cdc,
			transientStoreKey,
		)
		feetiersKeeper, _ = createFeeTiersKeeper(stateStore, statsKeeper, vaultKeeper, affiliatesKeeper, db, cdc)
		revShareKeeper, _, _ := createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feetiersKeeper, statsKeeper)
		marketMapKeeper, _ := createMarketMapKeeper(stateStore, db, cdc)
		pricesKeeper, _, _, _ = createPricesKeeper(stateStore, db, cdc, transientStoreKey, revShareKeeper, marketMapKeeper)
		// Mock time provider response for market creation.
		assetsKeeper, _ = createAssetsKeeper(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			transientStoreKey,
			true,
		)
		rewardsKeeper, storeKey = createRewardsKeeper(
			stateStore,
			assetsKeeper,
			bankKeeper,
			feetiersKeeper,
			pricesKeeper,
			indexerEventManager,
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
	return ctx, rewardsKeeper, feetiersKeeper, bankKeeper, assetsKeeper, pricesKeeper, indexerEventManager, storeKey
}

func createRewardsKeeper(
	stateStore storetypes.CommitMultiStore,
	assetsKeeper *assetskeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	feeTiersKeeper *feetierskeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (*rewardskeeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	transientStoreKey := storetypes.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	authorities := []string{
		delaymsgtypes.ModuleAddress.String(),
		lib.GovModuleAddress.String(),
	}
	k := rewardskeeper.NewKeeper(
		cdc,
		storeKey,
		transientStoreKey,
		assetsKeeper,
		bankKeeper,
		feeTiersKeeper,
		pricesKeeper,
		indexerEventManager,
		authorities,
	)

	return k, storeKey
}

func GetTradingRewardEventsFromIndexerTendermintBlock(
	block indexer_manager.IndexerTendermintBlock,
) []*indexerevents.TradingRewardsEventV1 {
	var rewardEvents []*indexerevents.TradingRewardsEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeTradingReward {
			continue
		}
		var rewardEvent indexerevents.TradingRewardsEventV1
		err := proto.Unmarshal(event.DataBytes, &rewardEvent)
		if err != nil {
			panic(err)
		}
		rewardEvents = append(rewardEvents, &rewardEvent)
	}
	return rewardEvents
}
