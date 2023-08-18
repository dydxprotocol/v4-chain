package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"math/big"
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
	if !order.IsShortTermOrder() && !order.IsStatefulOrder() {
		panic(fmt.Sprintf("Unsupported order type for equity tiers. Order: %+v", order))
	}

	// Compute the net collateral for the subaccount.
	netCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: order.GetSubaccountId(),
		},
	)
	if err != nil {
		return err
	}

	// Get the equity tier limit for the subaccount.
	equityTierLimit := k.getEquityTierLimitForOrderTypeAndNetCollateral(ctx, order.IsShortTermOrder(), *netCollateral)
	if equityTierLimit == nil {
		return nil
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

	// Verify that opening this order would not exceed the maximum amount of orders for the equity tier.
	// Note that once we combine the count of orders on the memclob with how many `uncommitted` or `to be committed`
	// stateful orders on the memclob we should always have a negative number since we only count order
	// cancellations/removals for orders that exist.
	if k.getNumberOfOpenOrdersForSubaccount(ctx, order.OrderId.SubaccountId, order.IsShortTermOrder()) >= equityTierLimit.Limit {
		return sdkerrors.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order id: %+v",
			equityTierLimit.Limit,
			order.GetOrderId(),
		)
	}
	return nil
}

// getEquityTierLimitForOrderTypeAndNetCollateral returns the equity tier limit for the specified order type
// and net collateral. If `nil` is returned then there are no equity tier limits defined and no equity
// tier limit enforcement should occur for this order type.
func (k Keeper) getEquityTierLimitForOrderTypeAndNetCollateral(ctx sdk.Context, shortTermOrder bool, netCollateral big.Int) *types.EquityTierLimit {
	var equityTierLimits []types.EquityTierLimit
	if shortTermOrder {
		equityTierLimits = k.GetEquityTierLimitConfiguration(ctx).ShortTermOrderEquityTiers
	} else {
		equityTierLimits = k.GetEquityTierLimitConfiguration(ctx).StatefulOrderEquityTiers
	}
	// If there are no equity tier limits defined then we return nil representing that the equity tier limit
	// should not be enforced.
	if len(equityTierLimits) == 0 {
		return nil
	}

	for _, limit := range equityTierLimits {
		if netCollateral.Cmp(limit.UsdTncRequired.BigInt()) >= 0 {
			return &limit
		}
	}

	// If equity tier limits are defined and we couldn't find one then return a default with a limit of 0.
	return &types.EquityTierLimit{}
}

// getNumberOfOpenOrdersForSubaccount returns the number of open orders for the provided subaccount.
func (k Keeper) getNumberOfOpenOrdersForSubaccount(ctx sdk.Context, subaccountId satypes.SubaccountId, shortTermOrder bool) uint32 {
	var filter func(types.OrderId) bool
	if shortTermOrder {
		filter = func(id types.OrderId) bool {
			return id.IsShortTermOrder()
		}
	} else {
		filter = func(id types.OrderId) bool {
			return id.IsStatefulOrder()
		}
	}

	// Count all the open orders that are on the `MemClob`.
	count := int32(k.MemClob.CountSubaccountOrders(ctx, subaccountId, filter))

	// Include the number of stateful orders that exist outside of the `MemClob`. During `DeliverTx` we use
	// the count of to be committed stateful orders while in `CheckTx` we use the count of uncommitted stateful orders.
	if !shortTermOrder {
		if lib.IsDeliverTxMode(ctx) {
			count += k.GetToBeCommittedStatefulOrderCount(ctx, subaccountId)
		} else {
			count += k.GetUncommittedStatefulOrderCount(ctx, subaccountId)
		}
	}
	return lib.MustConvertIntegerToUint32(count)
}

// EnforceEquityTierLimits removes orders for each subaccount which had their net collateral decrease during
// block processing.
func (k Keeper) EnforceEquityTierLimits(ctx sdk.Context) {
	subaccountIds := k.subaccountsKeeper.GetAllSubaccountsWithDecreasedNetCollateral(ctx)

	for _, subaccountId := range subaccountIds {
		netCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
			ctx,
			satypes.Update{
				SubaccountId: subaccountId,
			},
		)
		if err != nil {
			continue
		}

		// Handle closing short term orders if necessary.
		{
			equityTierLimit := k.getEquityTierLimitForOrderTypeAndNetCollateral(ctx, true, *netCollateral)
			if equityTierLimit == nil {
				continue
			}

			numOpenOrders := k.getNumberOfOpenOrdersForSubaccount(ctx, subaccountId, true)
			numOrdersToClose := int(numOpenOrders) - int(equityTierLimit.Limit)
			if numOrdersToClose > 0 {
				// TODO: Close short term orders
			}
		}

		// Handle closing stateful orders if necessary.
		{
			equityTierLimit := k.getEquityTierLimitForOrderTypeAndNetCollateral(ctx, true, *netCollateral)
			if equityTierLimit == nil {
				continue
			}

			numOpenOrders := k.getNumberOfOpenOrdersForSubaccount(ctx, subaccountId, true)
			numOrdersToClose := int(numOpenOrders) - int(equityTierLimit.Limit)
			if numOrdersToClose > 0 {
				// TODO: Close stateful orders
			}
		}
	}
}
