package keeper

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"
	indexer_manager "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	sdktest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sdk"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type GenesisInitializer interface {
	InitializeForGenesis(ctx sdk.Context)
}

type callback func(
	db *dbm.MemDB,
	registry codectypes.InterfaceRegistry,
	cdc *codec.ProtoCodec,
	stateStore storetypes.CommitMultiStore,
	transientStoreKey storetypes.StoreKey,
) []GenesisInitializer

func initKeepers(t testing.TB, cb callback) sdk.Context {
	ctx, stateStore, db := sdktest.NewSdkContextWithMultistore()
	// Mount transient store for indexer events, shared by all keepers that emit indexer events.
	transientStoreKey := storetypes.NewTransientStoreKey(indexer_manager.IndexerEventsKey)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)

	initializers := cb(db, module.InterfaceRegistry, cdc, stateStore, transientStoreKey)

	require.NoError(t, stateStore.LoadLatestVersion())

	for _, i := range initializers {
		i.InitializeForGenesis(ctx)
	}

	return ctx
}
