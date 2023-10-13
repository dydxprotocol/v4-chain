package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/common"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
)

// CreateUsdcAsset creates USDC in the assets module for tests.
func CreateUsdcAsset(ctx sdk.Context, assetsKeeper *keeper.Keeper) error {
	_, err := assetsKeeper.CreateAsset(
		ctx,
		constants.Usdc.Id,
		constants.Usdc.Symbol,
		constants.Usdc.Denom,
		constants.Usdc.DenomExponent,
		constants.Usdc.HasMarket,
		constants.Usdc.MarketId,
		constants.Usdc.AtomicResolution,
	)
	return err
}

func AssetsKeepers(
	t testing.TB,
	msgSenderEnabled bool,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	accountKeeper *authkeeper.AccountKeeper,
	bankKeeper *bankkeeper.BaseKeeper,
	storeKey storetypes.StoreKey,
) {
	var mockTimeProvider *mocks.TimeProvider
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		pricesKeeper, _, _, _, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		accountKeeper, _ = createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		keeper, storeKey = createAssetsKeeper(stateStore, db, cdc, pricesKeeper, transientStoreKey, msgSenderEnabled)

		return []GenesisInitializer{pricesKeeper, keeper}
	})
	// Mock time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)
	return ctx, keeper, pricesKeeper, accountKeeper, bankKeeper, storeKey
}

func createAssetsKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
	transientStoreKey storetypes.StoreKey,
	msgSenderEnabled bool,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(msgSenderEnabled)
	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		pk,
		mockIndexerEventsManager,
	)

	return k, storeKey
}

// GetAssetCreateEventsFromIndexerBlock returns the asset create events in the
// Indexer Block event Kafka message.
func GetAssetCreateEventsFromIndexerBlock(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) []*indexerevents.AssetCreateEventV1 {
	var assetEvents []*indexerevents.AssetCreateEventV1
	block := keeper.GetIndexerEventManager().ProduceBlock(ctx)
	if block == nil {
		return assetEvents
	}
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeAsset {
			continue
		}
		unmarshaler := common.UnmarshalerImpl{}
		var assetEvent indexerevents.AssetCreateEventV1
		err := unmarshaler.Unmarshal(event.DataBytes, &assetEvent)
		if err != nil {
			panic(err)
		}
		assetEvents = append(assetEvents, &assetEvent)
	}
	return assetEvents
}
