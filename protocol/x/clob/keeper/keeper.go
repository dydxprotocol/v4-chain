package keeper

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/dydxprotocol/v4/indexer/indexer_manager"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		transientStoreKey storetypes.StoreKey

		MemClob             types.MemClob
		subaccountsKeeper   types.SubaccountsKeeper
		assetsKeeper        types.AssetsKeeper
		bankKeeper          types.BankKeeper
		perpetualsKeeper    types.PerpetualsKeeper
		indexerEventManager indexer_manager.IndexerEventManager

		// in-memory stores
		seenPlaceOrderIds map[types.OrderId]types.Order_GoodTilBlock
		seenGoodTilBlocks map[types.Order_GoodTilBlock]map[types.OrderId]bool

		memStoreInitialized *atomic.Bool
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
	bankKeeper types.BankKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	keeper := &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		transientStoreKey:   liquidationsStoreKey,
		MemClob:             memClob,
		subaccountsKeeper:   subaccountsKeeper,
		assetsKeeper:        assetsKeeper,
		bankKeeper:          bankKeeper,
		perpetualsKeeper:    perpetualsKeeper,
		indexerEventManager: indexerEventManager,
		seenPlaceOrderIds:   make(map[types.OrderId]types.Order_GoodTilBlock),
		seenGoodTilBlocks:   make(map[types.Order_GoodTilBlock]map[types.OrderId]bool),
		memStoreInitialized: &atomic.Bool{},
	}

	// Provide the keeper to the MemClob.
	// The MemClob utilizes the keeper to read state fill amounts.
	memClob.SetClobKeeper(keeper)

	return keeper
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.setNumClobPairs(ctx, uint32(0))
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
		types.StatefulOrderPlacementKeyPrefix,
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
