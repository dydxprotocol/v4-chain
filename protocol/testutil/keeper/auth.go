package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	perpetualstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func createAccountKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	registry codectypes.InterfaceRegistry,
) (*keeper.AccountKeeper, storetypes.StoreKey) {
	types.RegisterInterfaces(registry)

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	paramKey := storetypes.NewKVStoreKey(paramtypes.StoreKey)
	stateStore.MountStoreWithDB(paramKey, storetypes.StoreTypeIAVL, db)
	paramTKey := storetypes.NewTransientStoreKey(paramtypes.TStoreKey)
	stateStore.MountStoreWithDB(paramTKey, storetypes.StoreTypeTransient, db)

	// Create default module account permissions for test.
	maccPerms := map[string][]string{
		minttypes.ModuleName:              {types.Minter},
		types.FeeCollectorName:            nil,
		satypes.ModuleName:                nil,
		perpetualstypes.InsuranceFundName: nil,
		ratelimittypes.SDaiPoolAccount:    nil,
		ibctransfertypes.ModuleName:       {types.Minter, types.Burner},
		satypes.LiquidityFeeModuleAddress: nil,
	}

	k := keeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		types.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		lib.GovModuleAddress.String(),
	)

	return &k, storeKey
}
