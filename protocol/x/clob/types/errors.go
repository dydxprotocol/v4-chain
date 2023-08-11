package types

// DONTCOVER

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/clob module sentinel errors
// TODO(CLOB-553) Clean up sentinel errors not in use.
var (
	ErrMemClobOrderDoesNotExist = sdkerrors.Register(
		ModuleName,
		2,
		"Order does not exist in memclob",
	)
	ErrInvalidOrderSide = sdkerrors.Register(
		ModuleName,
		3,
		"Invalid order side",
	)
	ErrInvalidOrderQuantums     = sdkerrors.Register(ModuleName, 4, "Invalid order quantums")
	ErrInvalidOrderGoodTilBlock = sdkerrors.Register(ModuleName, 5, "Invalid order goodTilBlock")
	ErrInvalidOrderSubticks     = sdkerrors.Register(ModuleName, 6, "Invalid order subticks")
	ErrOrderIsCanceled          = sdkerrors.Register(
		ModuleName,
		7,
		"Attempt to place an order that is already canceled",
	)
	ErrInvalidClob = sdkerrors.Register(
		ModuleName,
		8,
		"ClobPair ID does not reference a valid CLOB",
	)
	ErrMemClobCancelAlreadyExists = sdkerrors.Register(
		ModuleName,
		9,
		"A cancel already exists in the memclob for this order with a greater than or equal GoodTilBlock",
	)
	ErrHeightExceedsGoodTilBlock = sdkerrors.Register(
		ModuleName,
		10,
		"The next block height is greater than the GoodTilBlock of the message",
	)
	ErrGoodTilBlockExceedsShortBlockWindow = sdkerrors.Register(
		ModuleName,
		11,
		"The GoodTilBlock of the message is further than ShortBlockWindow blocks into the future",
	)
	ErrInvalidPlaceOrder = sdkerrors.Register(
		ModuleName,
		12,
		"MsgPlaceOrder is invalid",
	)
	ErrInvalidMsgProposedMatchOrders = sdkerrors.Register(
		ModuleName,
		13,
		"MsgProposedMatchOrders is invalid",
	)
	ErrStateFilledAmountNoChange = sdkerrors.Register(
		ModuleName,
		14,
		"State filled amount cannot be unchanged",
	)
	ErrStateFilledAmountDecreasing = sdkerrors.Register(
		ModuleName,
		15,
		"State filled amount cannot decrease",
	)
	ErrInvalidPruneStateFilledAmount = sdkerrors.Register(
		ModuleName,
		16,
		"Cannot prune state fill amount that does not exist",
	)
	ErrOrderWouldExceedMaxOpenOrdersPerClobAndSide = sdkerrors.Register(
		ModuleName,
		17,
		fmt.Sprintf(
			"Subaccount cannot open more than %d orders on a given CLOB and side",
			MaxSubaccountOrdersPerClobAndSide,
		),
	)
	ErrFillAmountNotDivisibleByStepSize = sdkerrors.Register(
		ModuleName,
		18,
		"`FillAmount` is not divisible by `StepBaseQuantums` of the specified `ClobPairId`",
	)
	ErrNoClobPairForPerpetual = sdkerrors.Register(
		ModuleName,
		19,
		"The provided perpetual ID does not have any associated CLOB pairs",
	)
	ErrInvalidReplacement = sdkerrors.Register(
		ModuleName,
		20,
		"An order with the same `OrderId` already exists for this CLOB with a greater-than-or-equal `GoodTilBlock` "+
			"or Order Hash",
	)
	ErrClobPairAndPerpetualDoNotMatch = sdkerrors.Register(
		ModuleName,
		21,
		"Clob pair and perpetual ids do not match",
	)
	ErrMatchedOrderNegativeFee = sdkerrors.Register(
		ModuleName,
		22,
		"Matched order has negative fee",
	)
	ErrSubaccountFeeTransferFailed = sdkerrors.Register(
		ModuleName,
		23,
		"Subaccounts updated for a matched order, but fee transfer to fee-collector failed",
	)
	ErrOrderFullyFilled = sdkerrors.Register(
		ModuleName,
		24,
		"Order is fully filled",
	)
	ErrPremiumWithNonPerpetualClobPair = sdkerrors.Register(
		ModuleName,
		25,
		"Attempting to get price premium with a non-perpetual CLOB pair",
	)
	ErrZeroIndexPriceForPremiumCalculation = sdkerrors.Register(
		ModuleName,
		26,
		"Index price is zero when calculating price premium",
	)
	ErrInvalidClobPairParameter = sdkerrors.Register(
		ModuleName,
		27,
		"Invalid ClobPair parameter",
	)
	ErrZeroPriceForOracle = sdkerrors.Register(
		ModuleName,
		28,
		"Oracle price must be > 0.",
	)
	ErrInvalidStatefulOrderCancellation = sdkerrors.Register(
		ModuleName,
		29,
		"Invalid stateful order cancellation",
	)
	ErrOrderReprocessed = sdkerrors.Register(
		ModuleName,
		30,
		"An order with the same `OrderId` and `OrderHash` has already been processed for this CLOB",
	)
	ErrMissingMidPrice = sdkerrors.Register(
		ModuleName,
		31,
		"Missing mid price for ClobPair",
	)

	// Liquidations errors.
	ErrInvalidLiquidationsConfig = sdkerrors.Register(
		ModuleName,
		1000,
		"Proposed LiquidationsConfig is invalid",
	)
	ErrNoPerpetualPositionsToLiquidate = sdkerrors.Register(
		ModuleName,
		1001,
		"Subaccount has no perpetual positions to liquidate",
	)
	ErrSubaccountNotLiquidatable = sdkerrors.Register(
		ModuleName,
		1002,
		"Subaccount is not liquidatable",
	)
	ErrNoOpenPositionForPerpetual = sdkerrors.Register(
		ModuleName,
		1003,
		"Subaccount does not have an open position for perpetual",
	)
	ErrInvalidLiquidationOrderTotalSize = sdkerrors.Register(
		ModuleName,
		1004,
		"Liquidation order has invalid size",
	)
	ErrInvalidLiquidationOrderSide = sdkerrors.Register(
		ModuleName,
		1005,
		"Liquidation order is on the wrong side",
	)
	ErrTotalFillAmountExceedsOrderSize = sdkerrors.Register(
		ModuleName,
		1006,
		"Total fills amount exceeds size of liquidation order",
	)
	ErrLiquidationContainsNoFills = sdkerrors.Register(
		ModuleName,
		1007,
		"Liquidation order does not contain any fills",
	)
	ErrSubaccountHasLiquidatedPerpetual = sdkerrors.Register(
		ModuleName,
		1008,
		"Subaccount has previously liquidated this perpetual in the current block",
	)
	ErrLiquidationOrderSizeSmallerThanMin = sdkerrors.Register(
		ModuleName,
		1009,
		"Liquidation order has size smaller than min position notional specified in the liquidation config",
	)
	ErrLiquidationOrderSizeGreaterThanMax = sdkerrors.Register(
		ModuleName,
		1010,
		"Liquidation order has size greater than max position notional specified in the liquidation config",
	)
	ErrLiquidationExceedsSubaccountMaxNotionalLiquidated = sdkerrors.Register(
		ModuleName,
		1011,
		"Liquidation exceeds the maximum notional amount that a single subaccount can have liquidated per block",
	)
	ErrLiquidationExceedsSubaccountMaxInsuranceLost = sdkerrors.Register(
		ModuleName,
		1012,
		"Liquidation exceeds the maximum insurance fund payout amount for a given subaccount per block",
	)
	ErrInsuranceFundHasInsufficientFunds = sdkerrors.Register(
		ModuleName,
		1013,
		"Insurance fund does not have sufficient funds to cover liquidation losses",
	)
	ErrInvalidPerpetualPositionSizeDelta = sdkerrors.Register(
		ModuleName,
		1014,
		"Invalid perpetual position size delta",
	)
	ErrInvalidQuantumsForInsuranceFundDeltaCalculation = sdkerrors.Register(
		ModuleName,
		1015,
		"Invalid delta base and/or quote quantums for insurance fund delta calculation",
	)
	ErrEmptyDeleveragingFills = sdkerrors.Register(
		ModuleName,
		1016,
		"Deleveraging fills length must be greater than 0",
	)
	ErrDeleveragingAgainstSelf = sdkerrors.Register(
		ModuleName,
		1017,
		"Cannot deleverage subaccount against itself",
	)
	ErrDuplicateDeleveragingFillSubaccounts = sdkerrors.Register(
		ModuleName,
		1018,
		"Deleveraging match cannot have fills with same id",
	)
	ErrZeroDeleveragingFillAmount = sdkerrors.Register(
		ModuleName,
		1019,
		"Deleveraging match cannot have fills with zero amount",
	)
	ErrPositionCannotBeFullyOffset = sdkerrors.Register(
		ModuleName,
		1020,
		"Position cannot be fully offset",
	)

	// Advanced order type errors.
	ErrFokOrderCouldNotBeFullyFilled = sdkerrors.Register(
		ModuleName,
		2000,
		"FillOrKill order could not be fully filled",
	)
	ErrReduceOnlyWouldIncreasePositionSize = sdkerrors.Register(
		ModuleName,
		2001,
		"Reduce-only orders cannot increase the position size",
	)
	ErrReduceOnlyWouldChangePositionSide = sdkerrors.Register(
		ModuleName,
		2002,
		"Reduce-only orders cannot change the position side",
	)
	ErrPostOnlyWouldCrossMakerOrder = sdkerrors.Register(
		ModuleName,
		2003,
		"Post-only order would cross one or more maker orders",
	)

	// Stateful order errors.
	ErrInvalidOrderFlag = sdkerrors.Register(
		ModuleName,
		3000,
		"Invalid order flags",
	)
	ErrInvalidStatefulOrderGoodTilBlockTime = sdkerrors.Register(
		ModuleName,
		3001,
		"Invalid order goodTilBlockTime",
	)
	ErrStatefulOrdersCannotRequireImmediateExecution = sdkerrors.Register(
		ModuleName,
		3002,
		"Stateful orders cannot require immediate execution",
	)
	ErrTimeExceedsGoodTilBlockTime = sdkerrors.Register(
		ModuleName,
		3003,
		"The block time is greater than the GoodTilBlockTime of the message",
	)
	ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow = sdkerrors.Register(
		ModuleName,
		3004,
		"The GoodTilBlockTime of the message is further than StatefulOrderTimeWindow into the future",
	)
	ErrStatefulOrderAlreadyExists = sdkerrors.Register(
		ModuleName,
		3005,
		"Existing stateful order has higher-or-equal priority than the new one",
	)
	ErrStatefulOrderDoesNotExist = sdkerrors.Register(
		ModuleName,
		3006,
		"Stateful order does not exist",
	)
	ErrStatefulOrderCollateralizationCheckFailed = sdkerrors.Register(
		ModuleName,
		3007,
		"Stateful order collateralization check failed",
	)
	ErrStatefulOrderPreviouslyCancelled = sdkerrors.Register(
		ModuleName,
		3008,
		"Stateful order was previously cancelled and therefore cannot be placed",
	)

	// Operations Queue validation errors
	ErrInvalidMsgProposedOperations = sdkerrors.Register(
		ModuleName,
		4000,
		"MsgProposedOperations is invalid",
	)
	ErrInvalidMatchOrder = sdkerrors.Register(
		ModuleName,
		4001,
		"Match Order is invalid",
	)
	ErrOrderPlacementNotInOperationsQueue = sdkerrors.Register(
		ModuleName,
		4002,
		"Order was not previously placed in operations queue",
	)
	ErrFillAmountIsZero = sdkerrors.Register(
		ModuleName,
		4003,
		"Fill amount cannot be zero",
	)
	ErrInvalidDeleveragingFills = sdkerrors.Register(
		ModuleName,
		4004,
		"Generated deleveraging fills do not match operations queue deleveraging fills",
	)
	ErrDeleveragedSubaccountNotLiquidatable = sdkerrors.Register(
		ModuleName,
		4005,
		"Deleveraged subaccount in proposed match operation is not liquidatable",
	)

	// Block rate limit errors.
	ErrInvalidBlockRateLimitConfig = sdkerrors.Register(
		ModuleName,
		5000,
		"Proposed BlockRateLimitConfig is invalid",
	)
	ErrBlockRateLimitExceeded = sdkerrors.Register(
		ModuleName,
		5001,
		"Block rate limit exceeded",
	)

	// Errors for unimplemented and disabled functionality.
	ErrAssetOrdersNotImplemented = sdkerrors.Register(
		ModuleName,
		9000,
		"Asset orders are not implemented",
	)
	ErrAssetUpdateNotImplemented = sdkerrors.Register(
		ModuleName,
		9001,
		"Updates for assets other than USDC are not implemented",
	)
	ErrNotImplemented = sdkerrors.Register(
		ModuleName,
		9002,
		"This function is not implemented",
	)
	ErrReduceOnlyDisabled = sdkerrors.Register(
		ModuleName,
		9003,
		"Reduce-only is currently disabled",
	)
)
