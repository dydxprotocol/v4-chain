package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
// This function should only be called from CLOB genesis or when an equity tier limit configuration
// change is accepted via governance.
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
// tier limit on how many open orders a subaccount can have. Short-term fill-or-kill and immediate-or-cancel orders
// never rest on the book and will always be allowed as they do not apply to the number of open orders that equity
// tier limits enforce.
//
// Note that the method is dependent on whether we are executing on `checkState` or on `deliverState` for
// stateful orders. During `checkState` we rely on the uncommitted order count to tell us how many stateful
// orders exist outside of the MemClob. For `deliverState` we use the `ProcessProposerMatchesEvents`
// to find out how many orders (minus removals) exists outside of the `MemClob`.
func (k Keeper) ValidateSubaccountEquityTierLimitForNewOrder(ctx sdk.Context, order types.Order) error {
	// Always allow short-term FoK or IoC orders as they will either fill immediately or be cancelled and won't rest on
	// the book.
	if order.IsShortTermOrder() && order.RequiresImmediateExecution() {
		return nil
	}

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
		return errorsmod.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order id: %+v",
			equityTierLimit.Limit,
			order.GetOrderId(),
		)
	}

	// Count all the open orders that are on the `MemClob`.
	equityTierCount := int32(k.MemClob.CountSubaccountOrders(ctx, subaccountId, filter))

	// Include the number of stateful orders that exist outside of the `MemClob`. This includes the number of
	// untriggered conditional orders and during `DeliverTx` we add the count of to be committed stateful orders
	// while in `CheckTx` we add the count of uncommitted stateful orders.
	if order.IsStatefulOrder() {
		equityTierCount += int32(k.CountUntriggeredSubaccountOrders(ctx, subaccountId, filter))
		if lib.IsDeliverTxMode(ctx) {
			equityTierCount += k.GetToBeCommittedStatefulOrderCount(ctx, order.OrderId)
		} else {
			equityTierCount += k.GetUncommittedStatefulOrderCount(ctx, order.OrderId)
		}
	}

	if equityTierCount < 0 {
		if lib.IsDeliverTxMode(ctx) {
			err = fmt.Errorf(
				"Expected ValidateSubaccountEquityTierLimitForNewOrder for new order %+v to be >= 0. "+
					"equityTierCount %d, memClobCount %d, (stateful order only) toBeCommittedCount %d.",
				order,
				equityTierCount,
				k.MemClob.CountSubaccountOrders(ctx, subaccountId, filter),
				k.GetToBeCommittedStatefulOrderCount(ctx, order.OrderId),
			)
		} else {
			err = fmt.Errorf(
				"Expected ValidateSubaccountEquityTierLimitForNewOrder for new order %+v to be >= 0. "+
					"equityTierCount %d, memClobCount %d, (stateful order only) uncommittedCount %d.",
				order,
				equityTierCount,
				k.MemClob.CountSubaccountOrders(ctx, subaccountId, filter),
				k.GetUncommittedStatefulOrderCount(ctx, order.OrderId),
			)
		}
		panic(err)
	}

	// Verify that opening this order would not exceed the maximum amount of orders for the equity tier.
	// Note that once we combine the count of orders on the memclob with how many `uncommitted` or `to be committed`
	// stateful orders on the memclob we should always have a negative number since we only count order
	// cancellations/removals for orders that exist.
	if lib.MustConvertIntegerToUint32(equityTierCount) >= equityTierLimit.Limit {
		return errorsmod.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order id: %+v",
			equityTierLimit.Limit,
			order.GetOrderId(),
		)
	}
	return nil
}
