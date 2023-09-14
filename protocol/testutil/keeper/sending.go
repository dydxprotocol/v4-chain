package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

func SendingKeepers(t testing.TB) (
	ctx sdk.Context,
	sendingKeeper *keeper.Keeper,
	accountKeeper *authkeeper.AccountKeeper,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpkeeper.Keeper,
	assetsKeeper *assetskeeper.Keeper,
	subaccountsKeeper types.SubaccountsKeeper,
	storeKey storetypes.StoreKey,
) {
	return SendingKeepersWithSubaccountsKeeper(t, nil)
}

func SendingKeepersWithSubaccountsKeeper(t testing.TB, saKeeper types.SubaccountsKeeper) (
	ctx sdk.Context,
	sendingKeeper *keeper.Keeper,
	accountKeeper *authkeeper.AccountKeeper,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpkeeper.Keeper,
	assetsKeeper *assetskeeper.Keeper,
	subaccountsKeeper types.SubaccountsKeeper,
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
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		pricesKeeper, _, _, _, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		perpetualsKeeper, _ = createPerpetualsKeeper(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			epochsKeeper,
			transientStoreKey,
		)
		assetsKeeper, _ = createAssetsKeeper(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			transientStoreKey,
			true,
		)
		accountKeeper, _ = createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ := createBankKeeper(stateStore, db, cdc, accountKeeper)
		if saKeeper == nil {
			subaccountsKeeper, _ = createSubaccountsKeeper(
				stateStore,
				db,
				cdc,
				assetsKeeper,
				bankKeeper,
				perpetualsKeeper,
				transientStoreKey,
				true,
			)
		} else {
			subaccountsKeeper = saKeeper
		}
		sendingKeeper, storeKey = createSendingKeeper(
			stateStore,
			db,
			cdc,
			accountKeeper,
			bankKeeper,
			subaccountsKeeper,
			transientStoreKey,
		)

		return []GenesisInitializer{pricesKeeper, perpetualsKeeper, assetsKeeper, sendingKeeper}
	})

	// Mock time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)

	return ctx,
		sendingKeeper,
		accountKeeper,
		pricesKeeper,
		perpetualsKeeper,
		assetsKeeper,
		subaccountsKeeper,
		storeKey
}

func createSendingKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	accKeeper *authkeeper.AccountKeeper,
	bankKeeper types.BankKeeper,
	saKeeper types.SubaccountsKeeper,
	transientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		accKeeper,
		bankKeeper,
		saKeeper,
		mockIndexerEventsManager,
		[]string{
			authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		},
	)

	return k, storeKey
}
