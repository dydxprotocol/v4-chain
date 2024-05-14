package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		assetsKeeper        types.AssetsKeeper
		bankKeeper          types.BankKeeper
		perpetualsKeeper    types.PerpetualsKeeper
		blocktimeKeeper     types.BlocktimeKeeper
		indexerEventManager indexer_manager.IndexerEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	assetsKeeper types.AssetsKeeper,
	bankKeeper types.BankKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	blocktimeKeeper types.BlocktimeKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		assetsKeeper:        assetsKeeper,
		bankKeeper:          bankKeeper,
		perpetualsKeeper:    perpetualsKeeper,
		blocktimeKeeper:     blocktimeKeeper,
		indexerEventManager: indexerEventManager,
	}
}

func (k Keeper) DebugCheckOpenInterestForPerpetuals(ctx sdk.Context) {
	perpOIMap := make(map[uint32]*big.Int)

	// Iterate through all subaccounts and perp positions for each subaccount.
	// Aggregate open interest for each perpetual market.
	subaccounts := k.GetAllSubaccount(ctx)
	for _, sa := range subaccounts {
		for _, perpPosition := range sa.PerpetualPositions {
			if perpPosition.Quantums.BigInt().Sign() <= 0 {
				// Only record positive positions for total open interest.
				// Total negative position size should be equal to total positive position size.
				continue
			}
			if openInterest, exists := perpOIMap[perpPosition.PerpetualId]; exists {
				// Already seen this perpetual. Add to open interest.
				openInterest.Add(
					openInterest,
					perpPosition.Quantums.BigInt(),
				)
			} else {
				// Haven't seen this pereptual. Initialize open interest.
				perpOIMap[perpPosition.PerpetualId] = new(big.Int).Set(
					perpPosition.Quantums.BigInt(),
				)
			}
		}
	}

	perps := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	for _, perp := range perps {
		if openInterest, exists := perpOIMap[perp.Params.Id]; exists {
			fmt.Printf("!!! perp: %+v\n", perp)
			if openInterest.Cmp(perp.OpenInterest.BigInt()) != 0 {
				fmt.Printf(
					"!!! Open Interest Check Failed: Perpetual %d: calculated %s, got %s\n",
					perp.Params.Id,
					openInterest.String(),
					perp.OpenInterest.BigInt().String(),
				)
				panic("Open Interest Check Failed")

			}
		} else {
			if perp.OpenInterest.BigInt().Sign() != 0 {
				fmt.Printf(
					"!!! Open Interest Check Failed: Perpetual %d: calculated 0, got %s\n",
					perp.Params.Id,
					perp.OpenInterest.BigInt().String(),
				)
				panic("Open Interest Check Failed")
			}
		}
	}

	fmt.Printf("!! Open Interest Check Succeeds, perpOIMap = %+v\n", perpOIMap)
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
