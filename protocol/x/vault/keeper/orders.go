package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// TODO (TRA-118): store vault strategy constants in x/vault state.
const (
	// Determines how many layers of orders a vault places.
	// E.g. if num_levels=2, a vault places 2 asks and 2 bids.
	NUM_LAYERS = uint8(2)
	// Determines minimum base spread when a vault quotes around reservation price.
	MIN_BASE_SPREAD_PPM = uint32(3_000) // 30bps
	// Determines the amount to add to min_price_change_ppm to arrive at base spread.
	BASE_SPREAD_MIN_PRICE_CHANGE_PREMIUM_PPM = uint32(1_500) // 15bps
	// Determines how aggressive a vault skews its orders.
	SKEW_FACTOR_PPM = uint32(500_000) // 0.5
	// Determines the percentage of vault equity that each order is sized at.
	ORDER_SIZE_PCT_PPM = uint32(100_000) // 10%
	// Determines how long a vault's orders are valid for.
	ORDER_EXPIRATION_SECONDS = uint32(5) // 5 seconds
)

// RefreshAllVaultOrders refreshes all orders for all vaults by
// TODO(TRA-134)
// 1. Cancelling all existing orders.
// 2. Placing new orders.
func (k Keeper) RefreshAllVaultOrders(ctx sdk.Context) {
}

// GetVaultClobOrders returns a list of long term orders for a given vault, with its corresponding
// clob pair, perpetual, market parameter, and market price.
// Let n be number of layers, then the function returns orders at [a_1, b_1, a_2, b_2, ..., a_n, b_n]
// where a_i and b_i are the ask price and bid price at i-th layer. To compute a_i and b_i:
// - a_i = oraclePrice * (1 + spread)^i
// - b_i = oraclePrice * (1 - spread)^i
// TODO (TRA-144): Implement order size
// TODO (TRA-114): Implement skew
func (k Keeper) GetVaultClobOrders(
	ctx sdk.Context,
	vaultId types.VaultId,
) (orders []*clobtypes.Order, err error) {
	// Get clob pair, perpetual, market parameter, and market price that correspond to this vault.
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrClobPairNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpId)
	if err != nil {
		return orders, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketParam, exists := k.pricesKeeper.GetMarketParam(ctx, perpetual.Params.MarketId)
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrMarketParamNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketPrice, err := k.pricesKeeper.GetMarketPrice(ctx, perpetual.Params.MarketId)
	if err != nil {
		return orders, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	// Get vault (subaccount 0 of corresponding module account).
	vault := satypes.SubaccountId{
		Owner:  vaultId.ToModuleAccountAddress(),
		Number: 0,
	}
	// Calculate spread.
	spreadPpm := lib.Max(
		MIN_BASE_SPREAD_PPM,
		BASE_SPREAD_MIN_PRICE_CHANGE_PREMIUM_PPM+marketParam.MinPriceChangePpm,
	)
	// Get market price in subticks.
	subticks := clobtypes.PriceToSubticks(
		marketPrice,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	// Get order expiration time.
	goodTilBlockTime := &clobtypes.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + ORDER_EXPIRATION_SECONDS,
	}
	// Construct one ask and one bid for each layer.
	orders = make([]*clobtypes.Order, 2*NUM_LAYERS)
	askSubticks := new(big.Rat).Set(subticks)
	bidSubticks := new(big.Rat).Set(subticks)
	for i := uint8(0); i < NUM_LAYERS; i++ {
		// Calculate ask and bid subticks for this layer.
		askSubticks = lib.BigRatMulPpm(askSubticks, lib.OneMillion+spreadPpm)
		bidSubticks = lib.BigRatMulPpm(bidSubticks, lib.OneMillion-spreadPpm)

		// Construct ask at this layer.
		ask := clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: vault,
				ClientId:     k.GetVaultClobOrderClientId(ctx, clobtypes.Order_SIDE_SELL, uint8(i+1)),
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   clobPair.Id,
			},
			Side:     clobtypes.Order_SIDE_SELL,
			Quantums: clobPair.StepBaseQuantums, // TODO (TRA-144): Implement order size
			Subticks: lib.BigRatRoundToNearestMultiple(
				askSubticks,
				clobPair.SubticksPerTick,
				true, // round up for asks
			),
			GoodTilOneof: goodTilBlockTime,
		}

		// Construct bid at this layer.
		bid := clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: vault,
				ClientId:     k.GetVaultClobOrderClientId(ctx, clobtypes.Order_SIDE_BUY, uint8(i+1)),
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   clobPair.Id,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: clobPair.StepBaseQuantums, // TODO (TRA-144): Implement order size
			Subticks: lib.BigRatRoundToNearestMultiple(
				bidSubticks,
				clobPair.SubticksPerTick,
				false, // round down for bids
			),
			GoodTilOneof: goodTilBlockTime,
		}

		orders[2*i] = &ask
		orders[2*i+1] = &bid
	}

	return orders, nil
}

// GetVaultClobOrderClientId returns the client ID for a CLOB order where
// - 1st bit is `side-1` (subtract 1 as buy_side = 1, sell_side = 2)
//
// - 2nd bit is `block height % 2`
//   - block height bit alternates between 0 and 1 to ensure that client IDs
//     are different in two consecutive blocks (otherwise, order placement would
//     fail because the same order IDs are already marked for cancellation)
//
// - next 8 bits are `layer`
func (k Keeper) GetVaultClobOrderClientId(
	ctx sdk.Context,
	side clobtypes.Order_Side,
	layer uint8,
) uint32 {
	sideBit := uint32(side - 1)
	sideBit <<= 31

	blockHeightBit := uint32(ctx.BlockHeight() % 2)
	blockHeightBit <<= 30

	layerBits := uint32(layer) << 22

	return sideBit | blockHeightBit | layerBits
}
