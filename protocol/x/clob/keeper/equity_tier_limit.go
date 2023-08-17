package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetEquityTierLimitConfiguration gets the equity tier limit configuration from state.
// The configuration is guaranteed to have been initialized ensuring specific field orderings.
func (k Keeper) GetEquityTierLimitConfiguration(
	ctx sdk.Context,
) (config types.EquityTierLimitConfiguration) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(
		types.KeyPrefix(
			types.EquityTierLimitConfigKey,
		),
	)

	// The equity tier limit configuration should be set in state by the genesis logic.
	// If it's not found, then that indicates it was never set in state, which is invalid.
	if b == nil {
		panic("GetEquityTierLimitConfiguration: EquityTierLimitConfiguration was never set in state")
	}

	k.cdc.MustUnmarshal(b, &config)

	return config
}

// InitializeEquityTierLimit initializes the equity tier limit configuration in state.
// This function should only be called from CLOB genesis.
func (k *Keeper) InitializeEquityTierLimit(
	ctx sdk.Context,
	config types.EquityTierLimitConfiguration,
) error {
	// Initialize the configuration, this effectively sorts the fields in a specific order
	// that the application expects.
	config.Initialize()

	// Validate the equity tier limit config before writing it to state.
	if err := config.Validate(); err != nil {
		return err
	}

	// Write the rate limit configuration to state.
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&config)
	store.Set(
		types.KeyPrefix(
			types.EquityTierLimitConfigKey,
		),
		b,
	)

	return nil
}

// ValidateSubaccountEquityTierLimitForNewOrder returns an error if adding the order would exceed the equity
// tier limit on how many open orders a subaccount can have.
//
// Note that the method is dependent on whether we are executing on `checkState` or on `deliverState` for
// stateful orders. During `checkState` we rely on the uncommitted order count to tell us how many stateful
// orders exist outside of the MemClob. For `deliverState` we use the `ProcessProposerMatchesEvents`
// to find out how many orders (minus removals) exists outside of the `MemClob`.
func (k Keeper) ValidateSubaccountEquityTierLimitForNewOrder(ctx sdk.Context, order types.Order) error {
	var equityTierLimits []types.EquityTierLimit
	var filter func(types.OrderId) bool
	if order.IsShortTermOrder() {
		equityTierLimits = k.GetEquityTierLimitConfiguration(ctx).ShortTermOrderEquityTiers
		filter = func(id types.OrderId) bool {
			return id.IsShortTermOrder()
		}
	} else if order.IsStatefulOrder() {
		equityTierLimits = k.GetEquityTierLimitConfiguration(ctx).StatefulOrderEquityTiers
		filter = func(id types.OrderId) bool {
			return id.IsStatefulOrder()
		}
	} else {
		panic(fmt.Sprintf("Unsupported order type for equity tiers. Order: %+v", order))
	}
	if len(equityTierLimits) == 0 {
		return nil
	}

	subaccountId := order.GetSubaccountId()
	netCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: subaccountId,
		},
	)
	if err != nil {
		return err
	}

	equityTierLimit := types.EquityTierLimit{}
	for _, limit := range equityTierLimits {
		if netCollateral.Cmp(limit.UsdTncRequired.BigInt()) >= 0 {
			equityTierLimit = limit
		} else {
			break
		}
	}
	// Return immediately if the amount the subaccount can open is 0.
	if equityTierLimit.Limit == 0 {
		return sdkerrors.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order id: %+v",
			equityTierLimit.Limit,
			order.GetOrderId(),
		)
	}

	// Count all the open orders that are on the `MemClob`.
	equityTierCount := int32(k.MemClob.CountSubaccountOrders(ctx, subaccountId, filter))

	// Include the number of stateful orders that exist outside of the `MemClob`. During `DeliverTx` we use
	// the `ProcessProposerMatchesEvents` to figure out the delta in how many orders will exist on
	// the `MemClob` while in `CheckTx` we use the count of uncommitted stateful orders.
	if order.IsStatefulOrder() {
		if lib.IsDeliverTxMode(ctx) {
			processProposerMatchesEvents := k.GetProcessProposerMatchesEvents(ctx)
			// Increment the count for every order that would be added to the order book for this subaccount.
			for _, newOrders := range [][]types.OrderId{
				processProposerMatchesEvents.PlacedLongTermOrderIds,
				processProposerMatchesEvents.PlacedConditionalOrderIds,
			} {
				for _, orderId := range newOrders {
					if subaccountId == orderId.SubaccountId && filter(orderId) {
						equityTierCount++
					}
				}
			}
			// Decrement the count for every order that would be removed from the order book for this subaccount
			for _, removedOrders := range [][]types.OrderId{
				processProposerMatchesEvents.PlacedStatefulCancellationOrderIds,
				processProposerMatchesEvents.ExpiredStatefulOrderIds,
				processProposerMatchesEvents.OrderIdsFilledInLastBlock,
			} {
				for _, orderId := range removedOrders {
					if subaccountId == orderId.SubaccountId && filter(orderId) {
						equityTierCount--
					}
				}
			}
		} else {
			equityTierCount += k.GetUncommittedStatefulOrderCount(ctx, order.OrderId)
		}
	}

	// Verify that opening this order would not exceed the maximum amount of orders for the equity tier.
	if lib.MustConvertIntegerToUint32(equityTierCount) >= equityTierLimit.Limit {
		return sdkerrors.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order id: %+v",
			equityTierLimit.Limit,
			order.GetOrderId(),
		)
	}
	return nil
}
