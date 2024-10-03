package keeper

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/dydxprotocol/v4-chain/protocol/finalizeblock"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	dydxlog "github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	streamingtypes "github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	flags "github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		transientStoreKey storetypes.StoreKey
		authorities       map[string]struct{}

		MemClob                 types.MemClob
		PerpetualIdToClobPairId map[uint32][]types.ClobPairId

		subaccountsKeeper types.SubaccountsKeeper
		assetsKeeper      types.AssetsKeeper
		bankKeeper        types.BankKeeper
		blockTimeKeeper   types.BlockTimeKeeper
		feeTiersKeeper    types.FeeTiersKeeper
		perpetualsKeeper  types.PerpetualsKeeper
		pricesKeeper      types.PricesKeeper
		statsKeeper       types.StatsKeeper
		rewardsKeeper     types.RewardsKeeper
		affiliatesKeeper  types.AffiliatesKeeper
		revshareKeeper    types.RevShareKeeper

		indexerEventManager      indexer_manager.IndexerEventManager
		streamingManager         streamingtypes.FullNodeStreamingManager
		finalizeBlockEventStager finalizeblock.EventStager[*types.ClobStagedFinalizeBlockEvent]

		initialized         *atomic.Bool
		memStoreInitialized *atomic.Bool

		Flags flags.ClobFlags

		mevTelemetryConfig MevTelemetryConfig

		// txValidation decoder and antehandler
		txDecoder sdk.TxDecoder
		// Note that the antehandler is not set until after the BaseApp antehandler is also set.
		antehandler sdk.AnteHandler

		placeCancelOrderRateLimiter rate_limit.RateLimiter[sdk.Msg]

		DaemonLiquidationInfo *liquidationtypes.DaemonLiquidationInfo
	}
)

var (
	_ types.ClobKeeper    = &Keeper{}
	_ types.MemClobKeeper = &Keeper{}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	transientStoreKey storetypes.StoreKey,
	authorities []string,
	memClob types.MemClob,
	subaccountsKeeper types.SubaccountsKeeper,
	assetsKeeper types.AssetsKeeper,
	blockTimeKeeper types.BlockTimeKeeper,
	bankKeeper types.BankKeeper,
	feeTiersKeeper types.FeeTiersKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	pricesKeeper types.PricesKeeper,
	statsKeeper types.StatsKeeper,
	rewardsKeeper types.RewardsKeeper,
	affiliatesKeeper types.AffiliatesKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
	streamingManager streamingtypes.FullNodeStreamingManager,
	txDecoder sdk.TxDecoder,
	clobFlags flags.ClobFlags,
	placeCancelOrderRateLimiter rate_limit.RateLimiter[sdk.Msg],
	daemonLiquidationInfo *liquidationtypes.DaemonLiquidationInfo,
	revshareKeeper types.RevShareKeeper,
) *Keeper {
	keeper := &Keeper{
		cdc:                     cdc,
		storeKey:                storeKey,
		memKey:                  memKey,
		transientStoreKey:       transientStoreKey,
		authorities:             lib.UniqueSliceToSet(authorities),
		MemClob:                 memClob,
		PerpetualIdToClobPairId: make(map[uint32][]types.ClobPairId),
		subaccountsKeeper:       subaccountsKeeper,
		assetsKeeper:            assetsKeeper,
		blockTimeKeeper:         blockTimeKeeper,
		bankKeeper:              bankKeeper,
		feeTiersKeeper:          feeTiersKeeper,
		perpetualsKeeper:        perpetualsKeeper,
		pricesKeeper:            pricesKeeper,
		statsKeeper:             statsKeeper,
		rewardsKeeper:           rewardsKeeper,
		affiliatesKeeper:        affiliatesKeeper,
		indexerEventManager:     indexerEventManager,
		streamingManager:        streamingManager,
		memStoreInitialized:     &atomic.Bool{}, // False by default.
		initialized:             &atomic.Bool{}, // False by default.
		txDecoder:               txDecoder,
		mevTelemetryConfig: MevTelemetryConfig{
			Enabled:    clobFlags.MevTelemetryEnabled,
			Hosts:      clobFlags.MevTelemetryHosts,
			Identifier: clobFlags.MevTelemetryIdentifier,
		},
		Flags:                       clobFlags,
		placeCancelOrderRateLimiter: placeCancelOrderRateLimiter,
		DaemonLiquidationInfo:       daemonLiquidationInfo,
		revshareKeeper:              revshareKeeper,
		finalizeBlockEventStager: finalizeblock.NewEventStager[*types.ClobStagedFinalizeBlockEvent](
			transientStoreKey,
			cdc,
			types.StagedEventsCountKey,
			types.StagedEventsKeyPrefix,
		),
	}

	// Provide the keeper to the MemClob.
	// The MemClob utilizes the keeper to read state fill amounts.
	memClob.SetClobKeeper(keeper)

	return keeper
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) GetFullNodeStreamingManager() streamingtypes.FullNodeStreamingManager {
	return k.streamingManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(
		log.ModuleKey, "x/clob",
		metrics.ProposerConsAddress, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress),
		metrics.CheckTx, ctx.IsCheckTx(),
		metrics.ReCheckTx, ctx.IsReCheckTx(),
	)
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}

// IsInitialized returns whether the clob keeper has been hydrated.
func (k Keeper) IsInitialized() bool {
	return k.initialized.Load()
}

// Initialize hydrates the clob keeper with the necessary in memory data structures.
func (k Keeper) Initialize(ctx sdk.Context) {
	alreadyInitialized := k.initialized.Swap(true)
	if alreadyInitialized {
		return
	}

	// Initialize memstore in clobKeeper with order fill amounts and stateful orders.
	k.InitMemStore(ctx)

	// Branch the context for hydration.
	// This means that new order matches from hydration will get added to the operations
	// queue but the corresponding state changes will be discarded.
	// This is needed because we are hydrating in memory structures in PreBlock
	// which operates on deliver state. Writing optimistic matches breaks consensus.
	checkCtx, _ := ctx.CacheContext()
	checkCtx = checkCtx.WithIsCheckTx(true)

	// Initialize memclob in clobKeeper with orderbooks using `ClobPairs` in state.
	k.InitMemClobOrderbooks(checkCtx)
	// Initialize memclob with all existing stateful orders.
	// TODO(DEC-1348): Emit indexer messages to indicate that application restarted.
	k.InitStatefulOrders(checkCtx)

	// Initialize the untriggered conditional orders data structure with untriggered
	// conditional orders in state.
	k.HydrateClobPairAndPerpetualMapping(checkCtx)
}

func (k Keeper) GetStagedClobFinalizeBlockEvents(ctx sdk.Context) []*types.ClobStagedFinalizeBlockEvent {
	return k.finalizeBlockEventStager.GetStagedFinalizeBlockEvents(
		ctx,
		func() *types.ClobStagedFinalizeBlockEvent {
			return &types.ClobStagedFinalizeBlockEvent{}
		},
	)
}

func (k Keeper) ProcessStagedFinalizeBlockEvents(ctx sdk.Context) {
	stagedEvents := k.GetStagedClobFinalizeBlockEvents(ctx)
	for _, stagedEvent := range stagedEvents {
		if stagedEvent == nil {
			// We don't ever expect this. However, should not panic since we are in Precommit.
			dydxlog.ErrorLog(
				ctx,
				"got nil ClobStagedFinalizeBlockEvent, skipping",
				"staged_events",
				stagedEvents,
			)
			continue
		}

		switch event := stagedEvent.Event.(type) {
		case *types.ClobStagedFinalizeBlockEvent_CreateClobPair:
			k.ApplySideEffectsForNewClobPair(ctx, *event.CreateClobPair)
		default:
			dydxlog.ErrorLog(
				ctx,
				"got unknown ClobStagedFinalizeBlockEvent",
				"event",
				event,
			)
		}
	}
}

// InitMemStore initializes the memstore of the `clob` keeper.
// This is called during app initialization in `app.go`, before any ABCI calls are received.
func (k Keeper) InitMemStore(ctx sdk.Context) {
	alreadyInitialized := k.memStoreInitialized.Swap(true)
	if alreadyInitialized {
		panic(errors.New("Memory store already initialized and is not intended to be invoked more then once."))
	}

	memStore := ctx.KVStore(k.memKey)
	memStoreType := memStore.GetStoreType()
	if memStoreType != storetypes.StoreTypeMemory {
		panic(
			fmt.Sprintf(
				"invalid memory store type; got %s, expected: %s",
				memStoreType,
				storetypes.StoreTypeMemory,
			),
		)
	}

	// Initialize all the necessary memory stores.
	for _, keyPrefix := range []string{
		types.StatefulOrderKeyPrefix,
	} {
		// Retrieve an instance of the memstore.
		memPrefixStore := prefix.NewStore(
			memStore,
			[]byte(keyPrefix),
		)

		// Retrieve an instance of the store.
		store := prefix.NewStore(
			ctx.KVStore(k.storeKey),
			[]byte(keyPrefix),
		)

		// Copy over all keys and values with the current key prefix to the `MemStore`.
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			memPrefixStore.Set(iterator.Key(), iterator.Value())
		}
	}

	// Ensure that the stateful order count is accurately represented in the memstore on restart.
	statefulOrders := k.GetAllStatefulOrders(ctx)
	for _, order := range statefulOrders {
		subaccountId := order.GetSubaccountId()
		k.SetStatefulOrderCount(
			ctx,
			subaccountId,
			k.GetStatefulOrderCount(ctx, subaccountId)+1,
		)
	}
}

// Sets the ante handler after it has been constructed. This breaks a cycle between
// when the ante handler is constructed and when the clob keeper is constructed.
func (k *Keeper) SetAnteHandler(anteHandler sdk.AnteHandler) {
	k.antehandler = anteHandler
}

func (k Keeper) GetSubaccountSnapshotsForInitStreams(
	ctx sdk.Context,
) (
	subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
) {
	lib.AssertCheckTxMode(ctx)

	return k.GetFullNodeStreamingManager().GetSubaccountSnapshotsForInitStreams(
		func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate {
			subaccountUpdate := k.subaccountsKeeper.GetStreamSubaccountUpdate(
				ctx,
				subaccountId,
				true,
			)
			return &subaccountUpdate
		},
	)
}

// InitializeNewStreams initializes new streams for all uninitialized clob pairs
// by sending the corresponding orderbook snapshots.
func (k Keeper) InitializeNewStreams(
	ctx sdk.Context,
	subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
) {
	streamingManager := k.GetFullNodeStreamingManager()

	streamingManager.InitializeNewStreams(
		func(clobPairId types.ClobPairId) *types.OffchainUpdates {
			return k.MemClob.GetOffchainUpdatesForOrderbookSnapshot(
				ctx,
				clobPairId,
			)
		},
		subaccountSnapshots,
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
}

// SendOrderbookUpdates sends the offchain updates to the Full Node streaming manager.
func (k Keeper) SendOrderbookUpdates(
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
) {
	if len(offchainUpdates.Messages) == 0 {
		return
	}

	k.GetFullNodeStreamingManager().SendOrderbookUpdates(
		offchainUpdates,
		ctx,
	)
}

// SendOrderbookFillUpdate sends the orderbook fills to the Full Node streaming manager.
func (k Keeper) SendOrderbookFillUpdate(
	ctx sdk.Context,
	orderbookFill types.StreamOrderbookFill,
) {
	k.GetFullNodeStreamingManager().SendOrderbookFillUpdate(
		orderbookFill,
		ctx,
		k.PerpetualIdToClobPairId,
	)
}

// SendTakerOrderStatus sends the taker order with its status to the Full Node streaming manager.
func (k Keeper) SendTakerOrderStatus(
	ctx sdk.Context,
	takerOrder types.StreamTakerOrder,
) {
	k.GetFullNodeStreamingManager().SendTakerOrderStatus(
		takerOrder,
		ctx,
	)
}
