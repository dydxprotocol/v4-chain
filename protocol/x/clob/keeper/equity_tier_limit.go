package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
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
	b := store.Get([]byte(types.EquityTierLimitConfigKey))

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
	store.Set([]byte(types.EquityTierLimitConfigKey), b)

	return nil
}

// getEquityTierLimitForSubaccount returns the equity tier limit for a subaccount based on its net collateral.
// The equity tier limit is the maximum amount of open orders a subaccount can have based on its net collateral
// and the equity tier limits configuration. An error is returned if calculating the net collateral returns an error or the user has zero allowed open orders based upon their net collateral.
// Also returns the net collateral of the subaccount for debug purposes.
func (k Keeper) getEquityTierLimitForSubaccount(
	ctx sdk.Context, subaccountId satypes.SubaccountId,
	equityTierLimits []types.EquityTierLimit,
) (equityTier types.EquityTierLimit, bigNetCollateral *big.Int, err error) {
	netCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: subaccountId,
		},
	)
	if err != nil {
		return types.EquityTierLimit{}, nil, err
	}

	var equityTierLimit types.EquityTierLimit
	for _, limit := range equityTierLimits {
		if netCollateral.Cmp(limit.UsdTncRequired.BigInt()) < 0 {
			break
		}
		equityTierLimit = limit
	}

	// Return immediately if the amount the subaccount can open is 0
	if equityTierLimit.Limit == 0 {
		return types.EquityTierLimit{}, nil, errorsmod.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d, for subaccount %+v with net collateral %+v",
			equityTierLimit.Limit,
			subaccountId,
			netCollateral,
		)
	}

	return equityTierLimit, nil, nil
}

// ValidateSubaccountEquityTierLimitForShortTermOrder returns an error if adding the order would exceed the equity
// tier limit on how many short term open orders a subaccount can have. Short-term fill-or-kill and immediate-or-cancel
// orders never rest on the book and will always be allowed as they do not apply to the number of open orders that
// equity tier limits enforce.
func (k Keeper) ValidateSubaccountEquityTierLimitForShortTermOrder(ctx sdk.Context, order types.Order) error {
	if order.RequiresImmediateExecution() {
		return nil
	}

	equityTierLimits := k.GetEquityTierLimitConfiguration(ctx).ShortTermOrderEquityTiers
	if len(equityTierLimits) == 0 {
		return nil
	}

	equityTierLimit, netCollateral, err := k.getEquityTierLimitForSubaccount(ctx, order.GetSubaccountId(), equityTierLimits)
	if err != nil {
		return err
	}

	// For short term orders we just count how many orders exist on the memclob.
	equityTierCount := k.MemClob.CountSubaccountShortTermOrders(ctx, order.GetSubaccountId())

	// Verify that opening this order would not exceed the maximum amount of orders for the equity tier.
	if lib.MustConvertIntegerToUint32(equityTierCount) >= equityTierLimit.Limit {
		return errorsmod.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order count: %d, total net collateral: %+v, order id: %+v",
			equityTierLimit.Limit,
			netCollateral,
			equityTierCount,
			order.GetOrderId(),
		)
	}
	return nil
}

// ValidateSubaccountEquityTierLimitForStatefulOrder returns an error if adding the order would exceed the equity
// tier limit on how many open orders a subaccount can have.
//
// Note that the method is dependent on whether we are executing on `checkState` or on `deliverState` for
// stateful orders. For `deliverState`, we sum:
//   - the number of long term orders.
//   - the number of conditional orders.
//   - the number of uncommitted stateful orders in the mempool (checkState only)
func (k Keeper) ValidateSubaccountEquityTierLimitForStatefulOrder(ctx sdk.Context, order types.Order) error {
	equityTierLimits := k.GetEquityTierLimitConfiguration(ctx).StatefulOrderEquityTiers
	if len(equityTierLimits) == 0 {
		return nil
	}

	equityTierLimit, netCollateral, err := k.getEquityTierLimitForSubaccount(ctx, order.GetSubaccountId(), equityTierLimits)
	if err != nil {
		return err
	}

	equityTierCount := k.GetStatefulOrderCount(ctx, order.OrderId.SubaccountId)
	// If this is `CheckTx` then we must also add the number of uncommitted stateful orders that this validator
	// is aware of (orders that are part of the mempool but have yet to proposed in a block).
	if !lib.IsDeliverTxMode(ctx) {
		equityTierCountMaybeNegative := k.GetUncommittedStatefulOrderCount(ctx, order.OrderId) + int32(equityTierCount)
		if equityTierCountMaybeNegative < 0 {
			panic(
				fmt.Errorf(
					"Expected ValidateSubaccountEquityTierLimitForStatefulOrder for new order %+v to be >= 0. "+
						"equityTierCount %d, statefulOrderCount %d, uncommittedStatefulOrderCount %d.",
					order,
					equityTierCountMaybeNegative,
					k.GetStatefulOrderCount(ctx, order.OrderId.SubaccountId),
					k.GetUncommittedStatefulOrderCount(ctx, order.OrderId),
				),
			)
		}
		equityTierCount = uint32(equityTierCountMaybeNegative)
	}

	// Verify that opening this order would not exceed the maximum amount of orders for the equity tier.
	if lib.MustConvertIntegerToUint32(equityTierCount) >= equityTierLimit.Limit {
		return errorsmod.Wrapf(
			types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
			"Opening order would exceed equity tier limit of %d. Order count: %d, total net collateral: %+v, order id: %+v",
			equityTierLimit.Limit,
			equityTierCount,
			netCollateral,
			order.GetOrderId(),
		)
	}
	return nil
}
