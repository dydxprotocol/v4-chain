package keeper

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

type SendingKeepersTestContext struct {
	Ctx               sdk.Context
	SendingKeeper     *keeper.Keeper
	AccountKeeper     *authkeeper.AccountKeeper
	BankKeeper        *bankkeeper.BaseKeeper
	PricesKeeper      *priceskeeper.Keeper
	PerpetualsKeeper  *perpkeeper.Keeper
	AssetsKeeper      *assetskeeper.Keeper
	SubaccountsKeeper types.SubaccountsKeeper
	StoreKey          storetypes.StoreKey
}

func SendingKeepers(t testing.TB) (
	ks SendingKeepersTestContext,
) {
	return SendingKeepersWithSubaccountsKeeper(t, nil)
}

func SendingKeepersWithSubaccountsKeeper(t testing.TB, saKeeper types.SubaccountsKeeper) (
	ks SendingKeepersTestContext,
) {
	var mockTimeProvider *mocks.TimeProvider
	ks.Ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		blockTimeKeeper, _ := createBlockTimeKeeper(stateStore, db, cdc)
		ks.AccountKeeper, _ = createAccountKeeper(
			stateStore,
			db,
			cdc,
			registry)
		ks.BankKeeper, _ = createBankKeeper(stateStore, db, cdc, ks.AccountKeeper)
		stakingKeeper, _ := createStakingKeeper(
			stateStore,
			db,
			cdc,
			ks.AccountKeeper,
			ks.BankKeeper,
		)
		statsKeeper, _ := createStatsKeeper(
			stateStore,
			epochsKeeper,
			db,
			cdc,
			stakingKeeper,
		)
		affiliatesKeeper, _ := createAffiliatesKeeper(stateStore, db, cdc, statsKeeper, transientStoreKey, true)
		vaultKeeper, _ := createVaultKeeper(stateStore, db, cdc, transientStoreKey)
		feetiersKeeper, _ := createFeeTiersKeeper(stateStore, statsKeeper, vaultKeeper, affiliatesKeeper, db, cdc)
		revShareKeeper, _, _ := createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feetiersKeeper, statsKeeper)
		marketMapKeeper, _ := createMarketMapKeeper(stateStore, db, cdc)
		ks.PricesKeeper, _, _, mockTimeProvider = createPricesKeeper(
			stateStore,
			db,
			cdc,
			transientStoreKey,
			revShareKeeper,
			marketMapKeeper,
		)
		ks.PerpetualsKeeper, _ = createPerpetualsKeeper(
			stateStore,
			db,
			cdc,
			ks.PricesKeeper,
			epochsKeeper,
			transientStoreKey,
		)
		ks.AssetsKeeper, _ = createAssetsKeeper(
			stateStore,
			db,
			cdc,
			ks.PricesKeeper,
			transientStoreKey,
			true,
		)
		if saKeeper == nil {
			ks.SubaccountsKeeper, _ = createSubaccountsKeeper(
				stateStore,
				db,
				cdc,
				ks.AssetsKeeper,
				ks.BankKeeper,
				ks.PerpetualsKeeper,
				blockTimeKeeper,
				transientStoreKey,
				true,
			)
		} else {
			ks.SubaccountsKeeper = saKeeper
		}
		ks.SendingKeeper, ks.StoreKey = createSendingKeeper(
			stateStore,
			db,
			cdc,
			ks.AccountKeeper,
			ks.BankKeeper,
			ks.SubaccountsKeeper,
			transientStoreKey,
		)

		return []GenesisInitializer{ks.PricesKeeper, ks.PerpetualsKeeper, ks.AssetsKeeper, ks.SendingKeeper}
	})

	// Mock time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)

	return ks
}

func createSendingKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	accKeeper *authkeeper.AccountKeeper,
	bankKeeper types.BankKeeper,
	saKeeper types.SubaccountsKeeper,
	transientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
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
			delaymsgtypes.ModuleAddress.String(),
			lib.GovModuleAddress.String(),
		},
	)

	return k, storeKey
}
