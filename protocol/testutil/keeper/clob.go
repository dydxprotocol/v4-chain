package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	db "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	asskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	blocktimekeeper "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	rewardskeeper "github.com/dydxprotocol/v4-chain/protocol/x/rewards/keeper"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
	subkeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
)

type ClobKeepersTestContext struct {
	Ctx               sdk.Context
	ClobKeeper        *keeper.Keeper
	PricesKeeper      *priceskeeper.Keeper
	AssetsKeeper      *asskeeper.Keeper
	BlockTimeKeeper   *blocktimekeeper.Keeper
	FeeTiersKeeper    *feetierskeeper.Keeper
	PerpetualsKeeper  *perpkeeper.Keeper
	StatsKeeper       *statskeeper.Keeper
	RewardsKeeper     *rewardskeeper.Keeper
	SubaccountsKeeper *subkeeper.Keeper
	StoreKey          storetypes.StoreKey
	MemKey            storetypes.StoreKey
	Cdc               *codec.ProtoCodec
}

func NewClobKeepersTestContext(
	t testing.TB,
	memClob types.MemClob,
	bankKeeper bankkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (ks ClobKeepersTestContext) {
	ks = NewClobKeepersTestContextWithUninitializedMemStore(t, memClob, bankKeeper, indexerEventManager)

	// Initialize the memstore.
	ks.ClobKeeper.InitMemStore(ks.Ctx)

	return ks
}

func NewClobKeepersTestContextWithUninitializedMemStore(
	t testing.TB,
	memClob types.MemClob,
	bankKeeper bankkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (ks ClobKeepersTestContext) {
	ks.Ctx = initKeepers(t, func(
		db *db.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		indexerEventsTransientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		ks.PricesKeeper, _, _, _, _ = createPricesKeeper(stateStore, db, cdc, indexerEventsTransientStoreKey)
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		ks.PerpetualsKeeper, _ = createPerpetualsKeeper(
			stateStore,
			db,
			cdc,
			ks.PricesKeeper,
			epochsKeeper,
			indexerEventsTransientStoreKey,
		)
		ks.AssetsKeeper, _ = createAssetsKeeper(
			stateStore,
			db,
			cdc,
			ks.PricesKeeper,
			indexerEventsTransientStoreKey,
			true,
		)
		ks.BlockTimeKeeper, _ = createBlockTimeKeeper(stateStore, db, cdc)
		ks.StatsKeeper, _ = createStatsKeeper(
			stateStore,
			epochsKeeper,
			db,
			cdc,
		)
		ks.FeeTiersKeeper, _ = createFeeTiersKeeper(
			stateStore,
			ks.StatsKeeper,
			db,
			cdc,
		)
		ks.RewardsKeeper, _ = createRewardsKeeper(
			stateStore,
			ks.AssetsKeeper,
			bankKeeper,
			ks.FeeTiersKeeper,
			ks.PricesKeeper,
			db,
			cdc,
		)
		ks.SubaccountsKeeper, _ = createSubaccountsKeeper(
			stateStore,
			db,
			cdc,
			ks.AssetsKeeper,
			bankKeeper,
			ks.PerpetualsKeeper,
			indexerEventsTransientStoreKey,
			true,
		)
		ks.ClobKeeper, ks.StoreKey, ks.MemKey = createClobKeeper(
			stateStore,
			db,
			cdc,
			memClob,
			ks.AssetsKeeper,
			ks.BlockTimeKeeper,
			bankKeeper,
			ks.FeeTiersKeeper,
			ks.PerpetualsKeeper,
			ks.StatsKeeper,
			ks.RewardsKeeper,
			ks.SubaccountsKeeper,
			indexerEventManager,
			indexerEventsTransientStoreKey,
		)
		ks.Cdc = cdc

		return []GenesisInitializer{
			ks.PricesKeeper,
			ks.PerpetualsKeeper,
			ks.AssetsKeeper,
			ks.SubaccountsKeeper,
			ks.ClobKeeper,
			ks.FeeTiersKeeper,
			ks.StatsKeeper,
		}
	})

	if err := ks.ClobKeeper.InitializeEquityTierLimit(ks.Ctx, types.EquityTierLimitConfiguration{}); err != nil {
		panic(err)
	}

	return ks
}

func createClobKeeper(
	stateStore storetypes.CommitMultiStore,
	db *db.MemDB,
	cdc *codec.ProtoCodec,
	memClob types.MemClob,
	aKeeper *asskeeper.Keeper,
	blockTimeKeeper types.BlockTimeKeeper,
	bankKeeper types.BankKeeper,
	feeTiersKeeper types.FeeTiersKeeper,
	perpKeeper *perpkeeper.Keeper,
	statsKeeper *statskeeper.Keeper,
	rewardsKeeper types.RewardsKeeper,
	saKeeper *subkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
	indexerEventsTransientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(types.TransientStoreKey)
	untriggeredConditionalOrders := make(map[types.ClobPairId]*keeper.UntriggeredConditionalOrders)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		transientStoreKey,
		memClob,
		untriggeredConditionalOrders,
		saKeeper,
		aKeeper,
		blockTimeKeeper,
		bankKeeper,
		feeTiersKeeper,
		perpKeeper,
		statsKeeper,
		rewardsKeeper,
		indexerEventManager,
		constants.TestEncodingCfg.TxConfig.TxDecoder(),
		flags.GetDefaultClobFlags(),
		rate_limit.NewNoOpRateLimiter[*types.MsgPlaceOrder](),
		rate_limit.NewNoOpRateLimiter[*types.MsgCancelOrder](),
	)
	k.SetAnteHandler(constants.EmptyAnteHandler)

	return k, storeKey, memKey
}
