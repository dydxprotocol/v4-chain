package keeper

import (
	"fmt"
	"testing"

	storetypes "cosmossdk.io/store/types"
	deleveragingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/deleveraging"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	streaming "github.com/StreamFinance-Protocol/stream-chain/protocol/streaming/grpc"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	asskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/keeper"
	blocktimekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/rate_limit"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	delaymsgmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	feetierskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/keeper"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	statskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/keeper"
	subkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
)

type ClobKeepersTestContext struct {
	Ctx               sdk.Context
	ClobKeeper        *keeper.Keeper
	PricesKeeper      *priceskeeper.Keeper
	AssetsKeeper      *asskeeper.Keeper
	BlockTimeKeeper   *blocktimekeeper.Keeper
	FeeTiersKeeper    *feetierskeeper.Keeper
	RatelimitKeeper   *ratelimitkeeper.Keeper
	PerpetualsKeeper  *perpkeeper.Keeper
	StatsKeeper       *statskeeper.Keeper
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
	var mockTimeProvider *mocks.TimeProvider
	ks.Ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		indexerEventsTransientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		ks.PricesKeeper, _, _, _, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, indexerEventsTransientStoreKey)
		// Mock time provider response for market creation.
		mockTimeProvider.On("Now").Return(constants.TimeT)
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
		ks.RatelimitKeeper, _ = createRatelimitKeeper(
			stateStore,
			db,
			cdc,
			ks.BlockTimeKeeper,
			bankKeeper,
			ks.PerpetualsKeeper,
			ks.AssetsKeeper,
			indexerEventsTransientStoreKey,
			true,
		)
		ks.SubaccountsKeeper, _ = createSubaccountsKeeper(
			stateStore,
			db,
			cdc,
			ks.AssetsKeeper,
			bankKeeper,
			ks.PerpetualsKeeper,
			ks.RatelimitKeeper,
			ks.BlockTimeKeeper,
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
			ks.PricesKeeper,
			ks.StatsKeeper,
			ks.SubaccountsKeeper,
			indexerEventManager,
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
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	memClob types.MemClob,
	aKeeper *asskeeper.Keeper,
	blockTimeKeeper types.BlockTimeKeeper,
	bankKeeper types.BankKeeper,
	feeTiersKeeper types.FeeTiersKeeper,
	perpKeeper *perpkeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	statsKeeper *statskeeper.Keeper,
	saKeeper *subkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (*keeper.Keeper, storetypes.StoreKey, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	transientStoreKey := storetypes.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		transientStoreKey,
		[]string{
			delaymsgmoduletypes.ModuleAddress.String(),
			lib.GovModuleAddress.String(),
		},
		memClob,
		saKeeper,
		aKeeper,
		blockTimeKeeper,
		bankKeeper,
		feeTiersKeeper,
		perpKeeper,
		pricesKeeper,
		statsKeeper,
		indexerEventManager,
		streaming.NewNoopGrpcStreamingManager(),
		constants.TestEncodingCfg.TxConfig.TxDecoder(),
		flags.GetDefaultClobFlags(),
		rate_limit.NewNoOpRateLimiter[sdk.Msg](),
		deleveragingtypes.NewDaemonDeleveragingInfo(),
		nil,
	)
	k.SetAnteHandler(constants.EmptyAnteHandler)

	return k, storeKey, memKey
}

func CreateTestClobPairs(
	t *testing.T,
	ctx sdk.Context,
	clobKeeper *keeper.Keeper,
	clobPairs []types.ClobPair,
) {
	for _, clobPair := range clobPairs {
		_, err := clobKeeper.CreatePerpetualClobPair(
			ctx,
			clobPair.Id,
			clobPair.MustGetPerpetualId(),
			satypes.BaseQuantums(clobPair.StepBaseQuantums),
			clobPair.QuantumConversionExponent,
			clobPair.SubticksPerTick,
			clobPair.Status,
		)
		require.NoError(t, err)
	}
}

func CreateNClobPair(
	t *testing.T,
	keeper *keeper.Keeper,
	perpKeeper *perpkeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	ctx sdk.Context,
	n int,
	mockIndexerEventManager *mocks.IndexerEventManager,
) []types.ClobPair {
	perps := CreateLiquidityTiersAndNPerpetuals(t, ctx, perpKeeper, pricesKeeper, n)

	items := make([]types.ClobPair, n)
	for i := range items {
		items[i].Id = uint32(i)
		items[i].Metadata = &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: uint32(i),
			},
		}
		items[i].SubticksPerTick = 5
		items[i].StepBaseQuantums = 5
		items[i].Status = types.ClobPair_STATUS_ACTIVE

		// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
		// the indexer event manager to expect these events.
		mockIndexerEventManager.On("AddTxnEvent",
			ctx,
			indexerevents.SubtypePerpetualMarket,
			indexerevents.PerpetualMarketEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewPerpetualMarketCreateEvent(
					clobtest.MustPerpetualId(items[i]),
					items[i].Id,
					perps[i].Params.Ticker,
					perps[i].Params.MarketId,
					items[i].Status,
					items[i].QuantumConversionExponent,
					perps[i].Params.AtomicResolution,
					items[i].SubticksPerTick,
					items[i].StepBaseQuantums,
					perps[i].Params.LiquidityTier,
					perps[i].Params.MarketType,
					perps[i].Params.DangerIndexPpm,
					fmt.Sprintf("%d", perps[i].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
				),
			),
		).Return()

		_, err := keeper.CreatePerpetualClobPair(
			ctx,
			items[i].Id,
			clobtest.MustPerpetualId(items[i]),
			satypes.BaseQuantums(items[i].StepBaseQuantums),
			items[i].QuantumConversionExponent,
			items[i].SubticksPerTick,
			items[i].Status,
		)
		if err != nil {
			panic(err)
		}
	}
	return items
}
