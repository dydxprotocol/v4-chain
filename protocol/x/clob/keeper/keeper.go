package keeper

import (
	"errors"
	"fmt"
	"sync/atomic"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	streamingtypes "github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/types"
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

		MemClob                      types.MemClob
		UntriggeredConditionalOrders map[types.ClobPairId]*UntriggeredConditionalOrders
		PerpetualIdToClobPairId      map[uint32][]types.ClobPairId

		subaccountsKeeper types.SubaccountsKeeper
		assetsKeeper      types.AssetsKeeper
		bankKeeper        types.BankKeeper
		blockTimeKeeper   types.BlockTimeKeeper
		feeTiersKeeper    types.FeeTiersKeeper
		perpetualsKeeper  types.PerpetualsKeeper
		pricesKeeper      types.PricesKeeper
		statsKeeper       types.StatsKeeper
		rewardsKeeper     types.RewardsKeeper

		indexerEventManager indexer_manager.IndexerEventManager
		streamingManager    streamingtypes.GrpcStreamingManager

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
	liquidationsStoreKey storetypes.StoreKey,
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
	indexerEventManager indexer_manager.IndexerEventManager,
	grpcStreamingManager streamingtypes.GrpcStreamingManager,
	txDecoder sdk.TxDecoder,
	clobFlags flags.ClobFlags,
	placeCancelOrderRateLimiter rate_limit.RateLimiter[sdk.Msg],
	daemonLiquidationInfo *liquidationtypes.DaemonLiquidationInfo,
) *Keeper {
	keeper := &Keeper{
		cdc:                          cdc,
		storeKey:                     storeKey,
		memKey:                       memKey,
		transientStoreKey:            liquidationsStoreKey,
		authorities:                  lib.UniqueSliceToSet(authorities),
		MemClob:                      memClob,
		UntriggeredConditionalOrders: make(map[types.ClobPairId]*UntriggeredConditionalOrders),
		PerpetualIdToClobPairId:      make(map[uint32][]types.ClobPairId),
		subaccountsKeeper:            subaccountsKeeper,
		assetsKeeper:                 assetsKeeper,
		blockTimeKeeper:              blockTimeKeeper,
		bankKeeper:                   bankKeeper,
		feeTiersKeeper:               feeTiersKeeper,
		perpetualsKeeper:             perpetualsKeeper,
		pricesKeeper:                 pricesKeeper,
		statsKeeper:                  statsKeeper,
		rewardsKeeper:                rewardsKeeper,
		indexerEventManager:          indexerEventManager,
		streamingManager:             grpcStreamingManager,
		memStoreInitialized:          &atomic.Bool{},
		txDecoder:                    txDecoder,
		mevTelemetryConfig: MevTelemetryConfig{
			Enabled:    clobFlags.MevTelemetryEnabled,
			Hosts:      clobFlags.MevTelemetryHosts,
			Identifier: clobFlags.MevTelemetryIdentifier,
		},
		Flags:                       clobFlags,
		placeCancelOrderRateLimiter: placeCancelOrderRateLimiter,
		DaemonLiquidationInfo:       daemonLiquidationInfo,
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

func (k Keeper) GetGrpcStreamingManager() streamingtypes.GrpcStreamingManager {
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
		types.OrderAmountFilledKeyPrefix,
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

// InitializeNewGrpcStreams initializes new gRPC streams for all uninitialized clob pairs
// by sending the corresponding orderbook snapshots.
func (k Keeper) InitializeNewGrpcStreams(ctx sdk.Context) {
	streamingManager := k.GetGrpcStreamingManager()
	allUpdates := types.NewOffchainUpdates()

	uninitializedClobPairIds := streamingManager.GetUninitializedClobPairIds()
	for _, clobPairId := range uninitializedClobPairIds {
		update := k.MemClob.GetOffchainUpdatesForOrderbookSnapshot(
			ctx,
			types.ClobPairId(clobPairId),
		)

		allUpdates.Append(update)
	}

	k.SendOrderbookUpdates(ctx, allUpdates, true)
}

// SendOrderbookUpdates sends the offchain updates to the gRPC streaming manager.
func (k Keeper) SendOrderbookUpdates(
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
	snapshot bool,
) {
	if len(offchainUpdates.Messages) == 0 {
		return
	}

	k.GetGrpcStreamingManager().SendOrderbookUpdates(
		offchainUpdates,
		snapshot,
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
}
