package keeper

import (
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/stretchr/testify/require"
)

type GenesisInitializer interface {
	InitializeForGenesis(ctx sdk.Context)
}

type callback func(
	db *tmdb.MemDB,
	registry codectypes.InterfaceRegistry,
	cdc *codec.ProtoCodec,
	stateStore storetypes.CommitMultiStore,
	transientStoreKey storetypes.StoreKey,
) []GenesisInitializer

func initKeepers(t testing.TB, cb callback) sdk.Context {
	ctx, stateStore, db := sdktest.NewSdkContextWithMultistore()
	// Mount transient store for indexer events, shared by all keepers that emit indexer events.
	transientStoreKey := sdk.NewTransientStoreKey(indexer_manager.IndexerEventsKey)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	initializers := cb(db, registry, cdc, stateStore, transientStoreKey)

	require.NoError(t, stateStore.LoadLatestVersion())

	for _, i := range initializers {
		i.InitializeForGenesis(ctx)
	}

	return ctx
}
