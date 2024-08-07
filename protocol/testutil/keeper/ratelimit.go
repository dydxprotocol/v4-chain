package keeper

import (
	sdaidaemontypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	blocktimekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/keeper"
	delaymsgtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	perpskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func createRatelimitKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	btk *blocktimekeeper.Keeper,
	bk bankkeeper.Keeper,
	perpk *perpskeeper.Keeper,
) (*ratelimitkeeper.Keeper, storetypes.StoreKey) {

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	authorities := []string{
		delaymsgtypes.ModuleAddress.String(),
		lib.GovModuleAddress.String(),
	}

	ics4wrapper := mocks.ICS4Wrapper{}

	sdaidaemontypes.SDAIEventFetcher = &sdaidaemontypes.MockEventFetcher{}
	sDAIEventManager := sdaidaemontypes.NewsDAIEventManager()

	k := ratelimitkeeper.NewKeeper(
		cdc,
		storeKey,
		sDAIEventManager,
		bk,
		*btk,
		*perpk,
		&ics4wrapper, // this is a pointer, since the mock has pointer receiver
		authorities,
	)

	return k, storeKey
}
