package keeper

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	"github.com/cometbft/cometbft/libs/log"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	flags "github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		MemClob                      types.MemClob
		UntriggeredConditionalOrders map[types.ClobPairId]*UntriggeredConditionalOrders
		PerpetualIdToClobPairId      map[uint32][]types.ClobPairId

		subaccountsKeeper   types.SubaccountsKeeper
		assetsKeeper        types.AssetsKeeper
		bankKeeper          types.BankKeeper
		blockTimeKeeper     types.BlockTimeKeeper
		feeTiersKeeper      types.FeeTiersKeeper
		perpetualsKeeper    types.PerpetualsKeeper
		statsKeeper         types.StatsKeeper
		rewardsKeeper       types.RewardsKeeper
		indexerEventManager indexer_manager.IndexerEventManager

		memStoreInitialized *atomic.Bool

		MaxLiquidationOrdersPerBlock uint32

		// mev telemetry config
		mevTelemetryHost       string
		mevTelemetryIdentifier string

		// txValidation decoder and antehandler
		txDecoder sdk.TxDecoder
		// Note that the antehandler is not set until after the BaseApp antehandler is also set.
		antehandler sdk.AnteHandler

		placeOrderRateLimiter  rate_limit.RateLimiter[*types.MsgPlaceOrder]
		cancelOrderRateLimiter rate_limit.RateLimiter[*types.MsgCancelOrder]

		// The address capable of updating ClobPairs.
		// Typically, this should be the x/gov module account.
		govAuthority string
	}
)

var _ types.ClobKeeper = &Keeper{}
var _ types.MemClobKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	liquidationsStoreKey storetypes.StoreKey,
	memClob types.MemClob,
	subaccountsKeeper types.SubaccountsKeeper,
	assetsKeeper types.AssetsKeeper,
	blockTimeKeeper types.BlockTimeKeeper,
	bankKeeper types.BankKeeper,
	feeTiersKeeper types.FeeTiersKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	statsKeeper types.StatsKeeper,
	rewardsKeeper types.RewardsKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
	txDecoder sdk.TxDecoder,
	clobFlags flags.ClobFlags,
	placeOrderRateLimiter rate_limit.RateLimiter[*types.MsgPlaceOrder],
	cancelOrderRateLimiter rate_limit.RateLimiter[*types.MsgCancelOrder],
	govAuthority string,
) *Keeper {
	keeper := &Keeper{
		cdc:                          cdc,
		storeKey:                     storeKey,
		memKey:                       memKey,
		transientStoreKey:            liquidationsStoreKey,
		MemClob:                      memClob,
		UntriggeredConditionalOrders: make(map[types.ClobPairId]*UntriggeredConditionalOrders),
		PerpetualIdToClobPairId:      make(map[uint32][]types.ClobPairId),
		subaccountsKeeper:            subaccountsKeeper,
		assetsKeeper:                 assetsKeeper,
		blockTimeKeeper:              blockTimeKeeper,
		bankKeeper:                   bankKeeper,
		feeTiersKeeper:               feeTiersKeeper,
		perpetualsKeeper:             perpetualsKeeper,
		statsKeeper:                  statsKeeper,
		rewardsKeeper:                rewardsKeeper,
		indexerEventManager:          indexerEventManager,
		memStoreInitialized:          &atomic.Bool{},
		txDecoder:                    txDecoder,
		mevTelemetryHost:             clobFlags.MevTelemetryHost,
		mevTelemetryIdentifier:       clobFlags.MevTelemetryIdentifier,
		MaxLiquidationOrdersPerBlock: clobFlags.MaxLiquidationOrdersPerBlock,
		placeOrderRateLimiter:        placeOrderRateLimiter,
		cancelOrderRateLimiter:       cancelOrderRateLimiter,
	}

	// Provide the keeper to the MemClob.
	// The MemClob utilizes the keeper to read state fill amounts.
	memClob.SetClobKeeper(keeper)

	return keeper
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

// GetGovAuthority returns the x/clob module's reference to the govAuthority address used
// to ensure only the gov module can send messages related to updating ClobPairs.
func (k Keeper) GetGovAuthority() string {
	return k.govAuthority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, "x/clob")
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
			types.KeyPrefix(keyPrefix),
		)

		// Retrieve an instance of the store.
		store := prefix.NewStore(
			ctx.KVStore(k.storeKey),
			types.KeyPrefix(keyPrefix),
		)

		// Copy over all keys and values with the current key prefix to the `MemStore`.
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			memPrefixStore.Set(iterator.Key(), iterator.Value())
		}
	}
}

// Sets the ante handler after it has been constructed. This breaks a cycle between
// when the ante handler is constructed and when the clob keeper is constructed.
func (k *Keeper) SetAnteHandler(anteHandler sdk.AnteHandler) {
	k.antehandler = anteHandler
}
