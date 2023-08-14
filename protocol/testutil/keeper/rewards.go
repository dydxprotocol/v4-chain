package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	rewardskeeper "github.com/dydxprotocol/v4-chain/protocol/x/rewards/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

func createRewardsKeeper(
	stateStore storetypes.CommitMultiStore,
	assetsKeeper *assetskeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	feeTiersKeeper *feetierskeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
) (*rewardskeeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	k := rewardskeeper.NewKeeper(
		cdc,
		storeKey,
		transientStoreKey,
		assetsKeeper,
		bankKeeper,
		feeTiersKeeper,
		pricesKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	return k, storeKey
}
