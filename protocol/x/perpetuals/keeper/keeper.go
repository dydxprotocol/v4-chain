package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		pricesKeeper        types.PricesKeeper
		epochsKeeper        types.EpochsKeeper
		clobKeeper          types.PerpetualsClobKeeper
		indexerEventManager indexer_manager.IndexerEventManager
		authorities         map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	pricesKeeper types.PricesKeeper,
	epochsKeeper types.EpochsKeeper,
	indexerEventsManager indexer_manager.IndexerEventManager,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		pricesKeeper:        pricesKeeper,
		epochsKeeper:        epochsKeeper,
		indexerEventManager: indexerEventsManager,
		authorities:         lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

// SetClobKeeper sets the `PerpetualsClobKeeper` reference, which is a Clob Keeper,
// for this Perpetuals Keeper.
// This method is called after the Perpetuals Keeper struct is initialized.
// This reference is set with an explicit method call rather than during `NewKeeper`
// due to the bidirectional dependency between the Perpetuals Keeper and the Clob Keeper.
func (k *Keeper) SetClobKeeper(getter types.PerpetualsClobKeeper) {
	k.clobKeeper = getter
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
