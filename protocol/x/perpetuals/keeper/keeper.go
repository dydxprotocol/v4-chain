package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/perpetuals/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		pricesKeeper        types.PricesKeeper
		epochsKeeper        types.EpochsKeeper
		pricePremiumGetter  types.PricePremiumGetter
		indexerEventManager indexer_manager.IndexerEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	pricesKeeper types.PricesKeeper,
	epochsKeeper types.EpochsKeeper,
	indexerEventsManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		pricesKeeper:        pricesKeeper,
		epochsKeeper:        epochsKeeper,
		indexerEventManager: indexerEventsManager,
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

// SetPricePremiumGetter sets the `PricePremiumGetter` reference, which is a Clob Keeper,
// for this Perpetuals Keeper.
// This method is called after the Perpetuals Keeper struct is initialized.
// This reference is set with an explicit method call rather than during `NewKeeper`
// due to the bidirectional dependency between the Perpetuals Keeper and the Clob Keeper.
func (k *Keeper) SetPricePremiumGetter(getter types.PricePremiumGetter) {
	k.pricePremiumGetter = getter
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.setNumPerpetuals(ctx, uint32(0))
	k.setNumLiquidityTiers(ctx, uint32(0))
}
