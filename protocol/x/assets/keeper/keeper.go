package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/assets/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		pricesKeeper types.PricesKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	pricesKeeper types.PricesKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		pricesKeeper: pricesKeeper,
	}
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.setNumAssets(ctx, uint32(0))
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}
