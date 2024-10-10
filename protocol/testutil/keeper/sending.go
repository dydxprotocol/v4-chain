package keeper

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"

	storetypes "cosmossdk.io/store/types"
	assetskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/keeper"
	delaymsgtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
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
	RatelimitKeeper   *ratelimitkeeper.Keeper
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
		ks.PricesKeeper, _, _, _, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
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
		ks.AccountKeeper, _ = createAccountKeeper(stateStore, db, cdc, registry)
		ks.BankKeeper, _ = createBankKeeper(stateStore, db, cdc, ks.AccountKeeper)
		ks.RatelimitKeeper, _ = createRatelimitKeeper(
			stateStore,
			db,
			cdc,
			blockTimeKeeper,
			ks.BankKeeper,
			ks.PerpetualsKeeper,
			ks.AssetsKeeper,
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
				ks.RatelimitKeeper,
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
