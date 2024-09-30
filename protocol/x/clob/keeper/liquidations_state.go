package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/store/prefix"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	liqInfoBytes := store.Get(subaccountId.ToStateKey())
	if liqInfoBytes == nil {
		return liquidationInfo
	}

	k.cdc.MustUnmarshal(liqInfoBytes, &liquidationInfo)
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
	liqInfoBytes := k.cdc.MustMarshal(&subaccountLiquidationInfo)
	store.Set(subaccountId.ToStateKey(), liqInfoBytes)
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

func (k Keeper) getCumulativeInsuranceFundDeltaStore(
	ctx sdk.Context,
) prefix.Store {
	store := prefix.NewStore(
		ctx.TransientStore(k.transientStoreKey),
		[]byte(types.CumulativeInsuranceFundDeltaKeyPrefix),
	)

	return store
}

// GetCumulativeInsuranceFundDelta returns the cumulative insurance fund delta as a big integer.
func (k Keeper) GetCumulativeInsuranceFundDelta(
	ctx sdk.Context,
	perpetualId uint32,
) (*big.Int, error) {

	insuranceFundName, err := k.perpetualsKeeper.GetInsuranceFundName(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	store := k.getCumulativeInsuranceFundDeltaStore(ctx)
	insuranceFundDeltaBytes := store.Get([]byte(insuranceFundName))
	if insuranceFundDeltaBytes == nil {
		return big.NewInt(0), nil
	}

	var delta big.Int
	delta.UnmarshalText(insuranceFundDeltaBytes)

	return &delta, nil
}

// IncrementCumulativeInsuranceFundDelta increments the cumulative insurance fund delta by the given amount.
func (k Keeper) IncrementCumulativeInsuranceFundDelta(
	ctx sdk.Context,
	perpetualId uint32,
	insuranceFundDelta *big.Int,
) error {

	insuranceFundName, err := k.perpetualsKeeper.GetInsuranceFundName(ctx, perpetualId)
	if err != nil {
		return err
	}

	store := k.getCumulativeInsuranceFundDeltaStore(ctx)
	currentDelta, err := k.GetCumulativeInsuranceFundDelta(ctx, perpetualId)
	if err != nil {
		return err
	}

	newDelta := new(big.Int).Sub(currentDelta, insuranceFundDelta)

	insuranceFundDeltaBytes, err := newDelta.MarshalText()
	if err != nil {
		return err
	}

	store.Set([]byte(insuranceFundName), insuranceFundDeltaBytes)
	return nil
}
