package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// TODO: Store these variables in state.
const (
	// Number of asks / bids to place from each vault.
	NUM_LAYERS = uint8(1)
	SPREAD     = 2
)

func (k Keeper) RefreshAllVaultOrders(ctx sdk.Context) {
	clobPairs := k.clobKeeper.GetAllClobPairs(ctx)
	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)
	perpetuals := k.perpetualsKeeper.GetAllPerpetuals(ctx)

	clobPairIdToClobPair := make(map[uint32]clobtypes.ClobPair)
	for _, clobPair := range clobPairs {
		clobPairIdToClobPair[clobPair.Id] = clobPair
	}
	marketIdToPrice := make(map[uint32]pricestypes.MarketPrice)
	for _, marketPrice := range marketPrices {
		marketIdToPrice[marketPrice.Id] = marketPrice
	}
	perpIdToPerp := make(map[uint32]perptypes.Perpetual)
	for _, perpetual := range perpetuals {
		perpIdToPerp[perpetual.GetId()] = perpetual
	}

	err := k.SetTotalShares(
		ctx,
		types.VaultId{
			Type:   types.VaultType_VAULT_TYPE_CLOB,
			Number: 0,
		},
		types.NumShares{
			NumShares: dtypes.NewInt(2),
		},
	)
	if err != nil {
		panic(err)
	}
	err = k.SetTotalShares(
		ctx,
		types.VaultId{
			Type:   types.VaultType_VAULT_TYPE_CLOB,
			Number: 1,
		},
		types.NumShares{
			NumShares: dtypes.NewInt(2),
		},
	)
	if err != nil {
		panic(err)
	}

	// Iterate through all vaults.
	totalSharesIterator := k.getTotalSharesIterator(ctx)
	defer totalSharesIterator.Close()
	for ; totalSharesIterator.Valid(); totalSharesIterator.Next() {
		var vaultId types.VaultId
		k.cdc.MustUnmarshal(totalSharesIterator.Key(), &vaultId)
		var totalShares types.NumShares
		k.cdc.MustUnmarshal(totalSharesIterator.Value(), &totalShares)

		// Skip if TotalShares is non-positive.
		if totalShares.NumShares.Cmp(dtypes.NewInt(0)) <= 0 {
			continue
		}

		fmt.Println("vaultId", vaultId)

		switch vaultId.Type {
		case types.VaultType_VAULT_TYPE_CLOB:
			// Get corresponding clob pair, perpetual, and market price.
			clobPair := clobPairIdToClobPair[vaultId.Number]
			perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
			perpetual := perpIdToPerp[perpId]
			marketPrice := marketIdToPrice[perpetual.Params.MarketId]

			ordersToCancel := k.ConstructVaultClobOrders(
				ctx.WithBlockHeight(ctx.BlockHeight()-1),
				vaultId,
				clobPair,
				perpetual,
				marketPrice,
			)
			fmt.Println("height", ctx.BlockHeight(), "ordersToCancel", ordersToCancel)
			// Cancel existing orders for this vault.
			// TODO (TRA-127): store existing orders in state and cancel them (as order IDs can change).
			for _, order := range ordersToCancel {
				if _, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId); exists {
					err := k.clobKeeper.HandleMsgCancelOrder(ctx, clobtypes.NewMsgCancelOrderStateful(
						order.OrderId,
						uint32(ctx.BlockTime().Unix())+50,
					))
					if err != nil {
						fmt.Println("failed to cancel order", "order", order, "error", err.Error())
					} else {
						fmt.Println("height", ctx.BlockHeight(), "cancelled order", "order", order)
					}
				}
			}

			// can't place orders if their IDs are the same as the ones cancelled above. two approaches:
			// 0. switch off a bit at different block height (0 for even 1 for odd)

			// Construct orders for this vault.
			orders := k.ConstructVaultClobOrders(ctx, vaultId, clobPair, perpetual, marketPrice)
			fmt.Println("height", ctx.BlockHeight(), "orders", orders)
			// Place orders for this vault.
			for _, order := range orders {
				err := k.clobKeeper.HandleMsgPlaceOrder(ctx, clobtypes.NewMsgPlaceOrder(order))
				if err != nil {
					fmt.Println("failed to place order", "order", order, "error", err.Error())
				} else {
					fmt.Println("height", ctx.BlockHeight(), "placed order", "order", order)
				}
			}
		}
	}
}

// ConstructVaultClobOrders constructs a list of orders for a given vault, with its
// corresponding clob pair, perpetual, and market price.
func (k Keeper) ConstructVaultClobOrders(
	ctx sdk.Context,
	vaultId types.VaultId,
	clobPair clobtypes.ClobPair,
	perpetual perptypes.Perpetual,
	marketPrice pricestypes.MarketPrice,
) (orders []clobtypes.Order) {
	// Get vault (subaccount 0 of the associated module account).
	// vault := k.subaccountsKeeper.GetSubaccount(ctx, satypes.SubaccountId{
	// 	Owner:  vaultId.ToModuleAccountAddress(),
	// 	Number: 0,
	// })

	// Construct one ask and one bid for each layer.
	subticks := clobtypes.PriceToSubticks(
		marketPrice,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	for layer := uint8(0); layer < NUM_LAYERS; layer++ {
		// Construct ask at this layer.
		ask := clobtypes.Order{
			OrderId: clobtypes.OrderId{
				// SubaccountId: vault,
				SubaccountId: satypes.SubaccountId{
					Owner:  "dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv",
					Number: 0,
				},
				ClientId:   k.getClientId(ctx, clobtypes.Order_SIDE_SELL, uint8(layer)),
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
				ClobPairId: clobPair.Id,
			},
			Side:     clobtypes.Order_SIDE_SELL,
			Quantums: clobPair.StepBaseQuantums,
			Subticks: lib.BigRatRoundToNearestMultiple(
				new(big.Rat).Mul(subticks, big.NewRat(100+SPREAD, 100)),
				clobPair.SubticksPerTick,
				false,
			),
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + 50},
		}

		// Construct bid at this layer.
		bid := clobtypes.Order{
			OrderId: clobtypes.OrderId{
				// SubaccountId: vault,
				SubaccountId: satypes.SubaccountId{
					Owner:  "dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv",
					Number: 0,
				},
				ClientId:   k.getClientId(ctx, clobtypes.Order_SIDE_BUY, uint8(layer)),
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
				ClobPairId: clobPair.Id,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: clobPair.StepBaseQuantums,
			Subticks: lib.BigRatRoundToNearestMultiple(
				new(big.Rat).Mul(subticks, big.NewRat(100-SPREAD, 100)),
				clobPair.SubticksPerTick,
				false,
			),
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + 50},
		}

		orders = append(orders, ask, bid)
	}

	return orders
}

// getClientId returns the client ID for a given side and layer where
// first bit is the side and the next 8 bits are the layer.
func (k Keeper) getClientId(
	ctx sdk.Context,
	side clobtypes.Order_Side,
	layer uint8,
) uint32 {
	clientIdBit := uint32(side)
	clientIdBit <<= 31

	blockHeightBit := uint32(ctx.BlockHeight() % 2)
	blockHeightBit <<= 30

	layerBits := uint32(layer) << 22

	return clientIdBit | blockHeightBit | layerBits
}
