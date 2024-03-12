package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	// clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		bankKeeper       types.BankKeeper
		clobKeeper       types.ClobKeeper
		perpetualsKeeper types.PerpetualsKeeper
		pricesKeeper     types.PricesKeeper
		authorities      map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	clobKeeper types.ClobKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	pricesKeeper types.PricesKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		bankKeeper:       bankKeeper,
		clobKeeper:       clobKeeper,
		perpetualsKeeper: perpetualsKeeper,
		pricesKeeper:     pricesKeeper,
		authorities:      lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) ProvideLiquidity(ctx sdk.Context) error {
	clobPairs := k.clobKeeper.GetAllClobPairs(ctx)
	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)
	perpetuals := k.perpetualsKeeper.GetAllPerpetuals(ctx)

	marketIdToPrice := make(map[uint32]pricestypes.MarketPrice)
	for _, marketPrice := range marketPrices {
		marketIdToPrice[marketPrice.Id] = marketPrice
	}
	perpIdToPrice := make(map[uint32]perptypes.Perpetual)
	for _, perpetual := range perpetuals {
		perpIdToPrice[perpetual.GetId()] = perpetual
	}

	for _, clobPair := range clobPairs {
		if clobPair.Id != 0 {
			continue
		}

		switch clobPair.Metadata.(type) {
		case *clobtypes.ClobPair_PerpetualClobMetadata:
			perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
			perpetual := perpIdToPrice[perpId]
			marketPrice := marketIdToPrice[perpetual.Params.MarketId]

			subticks := clobtypes.PriceToSubticks(
				marketPrice,
				clobPair,
				perpetual.Params.AtomicResolution,
				lib.QuoteCurrencyAtomicResolution,
			)
			fmt.Println("marketPrice for clob pair: ", clobPair.Id, marketPrice)
			fmt.Println("subticks for clob pair: ", clobPair.Id, subticks.String())
			// fmt.Println("subticksPerTick for clob pair: ", clobPair.Id, clobPair.SubticksPerTick)

			buySubticks := new(big.Rat).Mul(subticks, big.NewRat(98, 100))
			fmt.Println("buySubticks", clobPair.Id, buySubticks.String())
			buySubticksRoundedDown := RoundToNearestMultiple(
				new(big.Rat).Mul(subticks, big.NewRat(98, 100)),
				clobPair.SubticksPerTick,
				false,
			)
			fmt.Println("buySubticksRoundedDown", clobPair.Id, buySubticksRoundedDown)
			buyOrder := clobtypes.Order{
				OrderId: clobtypes.OrderId{
					SubaccountId: satypes.SubaccountId{
						Owner:  "dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv",
						Number: 0,
					},
					ClientId:   clobPair.Id,
					OrderFlags: clobtypes.OrderIdFlags_LongTerm,
					ClobPairId: clobPair.Id,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: clobPair.StepBaseQuantums,
				// Subticks: RoundToNearestMultiple(
				// 	new(big.Rat).Mul(subticks, big.NewRat(98, 100)),
				// 	clobPair.SubticksPerTick,
				// 	false,
				// ),
				Subticks:     100_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + 100},
			}

			ppme := k.clobKeeper.GetProcessProposerMatchesEvents(ctx)
			// 1. Cancel existing orders
			if _, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, buyOrder.OrderId); exists {
				err := k.clobKeeper.CancelStatefulOrder(ctx, &clobtypes.MsgCancelOrder{
					OrderId:      buyOrder.OrderId,
					GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + 100},
				})
				if err != nil {
					fmt.Println("error cancelling order: ", err)
				}
				ppme.PlacedStatefulCancellationOrderIds = append(
					ppme.PlacedStatefulCancellationOrderIds,
					buyOrder.OrderId,
				)
				fmt.Println("cancelled existing order")
			}

			// 2. Place new orders
			err := k.clobKeeper.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{
				Order: buyOrder,
			})
			if err != nil {
				fmt.Println("error placing order: ", err)
			}
			ppme.PlacedLongTermOrderIds = append(
				ppme.PlacedLongTermOrderIds,
				buyOrder.OrderId,
			)
			fmt.Println("placed new order")

			k.clobKeeper.MustSetProcessProposerMatchesEvents(ctx, ppme)
		default:
			panic("unexpected clob pair metadata type")
		}
	}

	return nil
}

// RoundToNearestMultiple rounds `value` up/down to the nearest multiple of `base`.
func RoundToNearestMultiple(
	value *big.Rat,
	base uint32,
	up bool,
) uint64 {
	quotient := new(big.Rat).Quo(
		value,
		new(big.Rat).SetUint64(uint64(base)),
	)
	quotientFloored := new(big.Int).Div(quotient.Num(), quotient.Denom())

	if up && quotientFloored.Cmp(quotient.Num()) != 0 {
		return (quotientFloored.Uint64() + 1) * uint64(base)
	}

	return quotientFloored.Uint64() * uint64(base)
}
