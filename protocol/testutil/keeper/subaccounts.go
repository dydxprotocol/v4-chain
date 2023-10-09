package keeper

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

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
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	asskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perpskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func SubaccountsKeepers(
	t testing.TB,
	msgSenderEnabled bool,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpskeeper.Keeper,
	accountKeeper *authkeeper.AccountKeeper,
	bankKeeper *bankkeeper.BaseKeeper,
	assetsKeeper *asskeeper.Keeper,
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
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		perpetualsKeeper, _ = createPerpetualsKeeper(stateStore, db, cdc, pricesKeeper, epochsKeeper, transientStoreKey)
		assetsKeeper, _ = createAssetsKeeper(stateStore, db, cdc, pricesKeeper, transientStoreKey, msgSenderEnabled)

		accountKeeper, _ = createAccountKeeper(stateStore, db, cdc, registry)

		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		keeper, storeKey = createSubaccountsKeeper(
			stateStore,
			db,
			cdc,
			assetsKeeper,
			bankKeeper,
			perpetualsKeeper,
			transientStoreKey,
			msgSenderEnabled,
		)

		return []GenesisInitializer{pricesKeeper, perpetualsKeeper, assetsKeeper, keeper}
	})

	// Mock time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)

	return ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, storeKey
}

func createSubaccountsKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	ak *asskeeper.Keeper,
	bk types.BankKeeper,
	pk *perpskeeper.Keeper,
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
		ak,
		bk,
		pk,
		mockIndexerEventsManager,
	)

	return k, storeKey
}

func CreateUsdcAssetPosition(
	quoteBalance *big.Int,
) []*types.AssetPosition {
	return []*types.AssetPosition{
		{
			AssetId:  assettypes.AssetUsdc.Id,
			Quantums: dtypes.NewIntFromBigInt(quoteBalance),
		},
	}
}

func CreateUsdcAssetUpdate(
	deltaQuoteBalance *big.Int,
) []types.AssetUpdate {
	return []types.AssetUpdate{
		{
			AssetId:          assettypes.AssetUsdc.Id,
			BigQuantumsDelta: deltaQuoteBalance,
		},
	}
}

// GetSubaccountUpdateEventsFromIndexerBlock returns the subaccount update events in the
// Indexer Block event Kafka message.
func GetSubaccountUpdateEventsFromIndexerBlock(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) []*indexerevents.SubaccountUpdateEventV1 {
	var subaccountUpdates []*indexerevents.SubaccountUpdateEventV1
	block := keeper.GetIndexerEventManager().ProduceBlock(ctx)
	if block == nil {
		return subaccountUpdates
	}
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeSubaccountUpdate {
			continue
		}
		unmarshaler := common.UnmarshalerImpl{}
		var subaccountUpdate indexerevents.SubaccountUpdateEventV1
		err := unmarshaler.Unmarshal(event.DataBytes, &subaccountUpdate)
		if err != nil {
			panic(err)
		}
		subaccountUpdates = append(subaccountUpdates, &subaccountUpdate)
	}
	return subaccountUpdates
}
