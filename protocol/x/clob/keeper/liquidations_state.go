package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetSubaccountLiquidationInfo returns liquidation information of the given subaccount in the
// current block. If it hasn't been liquidated in the last block, it returns the default
// `SubaccountLiquidationInfo`.
func (k Keeper) GetSubaccountLiquidationInfo(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	liquidationInfo types.SubaccountLiquidationInfo,
) {
	store := k.getSubaccountLiquidationInfoStore(ctx)

	b := store.Get(subaccountId.ToStateKey())
	if b == nil {
		return liquidationInfo
	}

	k.cdc.MustUnmarshal(b, &liquidationInfo)
	return liquidationInfo
}

// MustUpdateSubaccountPerpetualLiquidated adds a perpetual id to the subaccount liquidation info
// for the current block.
func (k Keeper) MustUpdateSubaccountPerpetualLiquidated(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) {
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	// Panic if this perpetual has already been liquidated within this block.
	if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(perpetualId) {
		panic(
			fmt.Sprintf(
				"UpdateSubaccountLiquidationInfo: subaccount %v and perpetual %d have already"+
					" been liquidated within the last block",
				subaccountId,
				perpetualId,
			),
		)
	}
	subaccountLiquidationInfo.PerpetualsLiquidated = append(
		subaccountLiquidationInfo.PerpetualsLiquidated,
		perpetualId,
	)

	store := k.getSubaccountLiquidationInfoStore(ctx)
	b := k.cdc.MustMarshal(&subaccountLiquidationInfo)
	store.Set(subaccountId.ToStateKey(), b)
}

// UpdateSubaccountLiquidationInfo updates the total notional liquidated and total insurance lost
// of the given subaccount for the current block.
func (k Keeper) UpdateSubaccountLiquidationInfo(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	notionalLiquidatedQuoteQuantums *big.Int,
	insuranceFundDeltaQuoteQuantums *big.Int,
) {
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	updatedNotionalLiquidatedQuoteQuantums := new(big.Int).Add(
		new(big.Int).Abs(notionalLiquidatedQuoteQuantums),
		new(big.Int).SetUint64(subaccountLiquidationInfo.NotionalLiquidated),
	)
	if !updatedNotionalLiquidatedQuoteQuantums.IsUint64() {
		// This should never happen, since the total notional liquidated for any subaccount should
		// never exceed the value of maximum notional liquidated (uint64) in the liquidation config.
		panic(
			errorsmod.Wrapf(
				satypes.ErrIntegerOverflow,
				"Notional liquidated update for subaccount %v overflows uint64",
				subaccountId,
			),
		)
	}

	subaccountLiquidationInfo.NotionalLiquidated = updatedNotionalLiquidatedQuoteQuantums.Uint64()

	// Update the total insurance funds lost for this subaccount if the insurance fund delta is
	// negative.
	if insuranceFundDeltaQuoteQuantums.Sign() == -1 {
		updatedQuantumsInsuranceLost := new(big.Int).Add(
			new(big.Int).Abs(insuranceFundDeltaQuoteQuantums),
			new(big.Int).SetUint64(subaccountLiquidationInfo.QuantumsInsuranceLost),
		)
		if !updatedQuantumsInsuranceLost.IsUint64() {
			// This should never happen, since the total insurance lost for any subaccount should never
			// exceed the value of maximum insurance lost (uint64) in the liquidation config.
			// This should also never exceed the maximum possible insurance fund balance.
			panic(
				errorsmod.Wrapf(
					satypes.ErrIntegerOverflow,
					"Quantums insurance lost update for subaccount %v overflows uint64",
					subaccountId,
				),
			)
		}

		subaccountLiquidationInfo.QuantumsInsuranceLost = updatedQuantumsInsuranceLost.Uint64()
	}

	store := k.getSubaccountLiquidationInfoStore(ctx)
	b := k.cdc.MustMarshal(&subaccountLiquidationInfo)
	store.Set(subaccountId.ToStateKey(), b)
}

// getSubaccountLiquidationInfoStore is an internal helper function for fetching the store
// used for reading and writing `SubaccountLiquidationInfo` messages to the transient store.
func (k Keeper) getSubaccountLiquidationInfoStore(
	ctx sdk.Context,
) prefix.Store {
	store := prefix.NewStore(
		ctx.TransientStore(k.transientStoreKey),
		[]byte(types.SubaccountLiquidationInfoKeyPrefix),
	)

	return store
}
