package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/clob module sentinel errors
// TODO(CLOB-553) Clean up sentinel errors not in use.
var (
	ErrMemClobOrderDoesNotExist = moderrors.Register(
		ModuleName,
		2,
		"Order does not exist in memclob",
	)
	ErrInvalidOrderSide = moderrors.Register(
		ModuleName,
		3,
		"Invalid order side",
	)
	ErrInvalidOrderQuantums     = moderrors.Register(ModuleName, 4, "Invalid order quantums")
	ErrInvalidOrderGoodTilBlock = moderrors.Register(ModuleName, 5, "Invalid order goodTilBlock")
	ErrInvalidOrderSubticks     = moderrors.Register(ModuleName, 6, "Invalid order subticks")
	ErrOrderIsCanceled          = moderrors.Register(
		ModuleName,
		7,
		"Attempt to place an order that is already canceled",
	)
	ErrInvalidClob = moderrors.Register(
		ModuleName,
		8,
		"ClobPair ID does not reference a valid CLOB",
	)
	ErrMemClobCancelAlreadyExists = moderrors.Register(
		ModuleName,
		9,
		"A cancel already exists in the memclob for this order with a greater than or equal GoodTilBlock",
	)
	ErrHeightExceedsGoodTilBlock = moderrors.Register(
		ModuleName,
		10,
		"The next block height is greater than the GoodTilBlock of the message",
	)
	ErrGoodTilBlockExceedsShortBlockWindow = moderrors.Register(
		ModuleName,
		11,
		"The GoodTilBlock of the message is further than ShortBlockWindow blocks into the future",
	)
	ErrInvalidPlaceOrder = moderrors.Register(
		ModuleName,
		12,
		"MsgPlaceOrder is invalid",
	)
	ErrInvalidMsgProposedMatchOrders = moderrors.Register(
		ModuleName,
		13,
		"MsgProposedMatchOrders is invalid",
	)
	ErrStateFilledAmountNoChange = moderrors.Register(
		ModuleName,
		14,
		"State filled amount cannot be unchanged",
	)
	ErrStateFilledAmountDecreasing = moderrors.Register(
		ModuleName,
		15,
		"State filled amount cannot decrease",
	)
	ErrInvalidPruneStateFilledAmount = moderrors.Register(
		ModuleName,
		16,
		"Cannot prune state fill amount that does not exist",
	)
	ErrOrderWouldExceedMaxOpenOrdersPerClobAndSide = moderrors.Register(
		ModuleName,
		17,
		"Subaccount cannot open more than 20 orders on a given CLOB and side",
	)
	ErrFillAmountNotDivisibleByStepSize = moderrors.Register(
		ModuleName,
		18,
		"`FillAmount` is not divisible by `StepBaseQuantums` of the specified `ClobPairId`",
	)
	ErrNoClobPairForPerpetual = moderrors.Register(
		ModuleName,
		19,
		"The provided perpetual ID does not have any associated CLOB pairs",
	)
	ErrInvalidReplacement = moderrors.Register(
		ModuleName,
		20,
		"An order with the same `OrderId` already exists for this CLOB with a greater-than-or-equal `GoodTilBlock` "+
			"or Order Hash",
	)
	ErrClobPairAndPerpetualDoNotMatch = moderrors.Register(
		ModuleName,
		21,
		"Clob pair and perpetual ids do not match",
	)
	ErrMatchedOrderNegativeFee = moderrors.Register(
		ModuleName,
		22,
		"Matched order has negative fee",
	)
	ErrSubaccountFeeTransferFailed = moderrors.Register(
		ModuleName,
		23,
		"Subaccounts updated for a matched order, but fee transfer to fee-collector failed",
	)
	ErrOrderFullyFilled = moderrors.Register(
		ModuleName,
		24,
		"Order is fully filled",
	)
	ErrPremiumWithNonPerpetualClobPair = moderrors.Register(
		ModuleName,
		25,
		"Attempting to get price premium with a non-perpetual CLOB pair",
	)
	ErrZeroIndexPriceForPremiumCalculation = moderrors.Register(
		ModuleName,
		26,
		"Index price is zero when calculating price premium",
	)
	ErrInvalidClobPairParameter = moderrors.Register(
		ModuleName,
		27,
		"Invalid ClobPair parameter",
	)
	ErrZeroPriceForOracle = moderrors.Register(
		ModuleName,
		28,
		"Oracle price must be > 0.",
	)
	ErrInvalidStatefulOrderCancellation = moderrors.Register(
		ModuleName,
		29,
		"Invalid stateful order cancellation",
	)
	ErrOrderReprocessed = moderrors.Register(
		ModuleName,
		30,
		"An order with the same `OrderId` and `OrderHash` has already been processed for this CLOB",
	)
	ErrMissingMidPrice = moderrors.Register(
		ModuleName,
		31,
		"Missing mid price for ClobPair",
	)
	ErrStatefulOrderCancellationAlreadyExists = moderrors.Register(
		ModuleName,
		32,
		"Existing stateful order cancellation has higher-or-equal priority than the new one",
	)
	ErrClobPairAlreadyExists = moderrors.Register(
		ModuleName,
		33,
		"ClobPair with id already exists",
	)
	ErrOrderConflictsWithClobPairStatus = moderrors.Register(
		ModuleName,
		34,
		"Order conflicts with ClobPair status",
	)
	ErrInvalidClobPairStatusTransition = moderrors.Register(
		ModuleName,
		35,
		"Invalid ClobPair status transition",
	)
	ErrOperationConflictsWithClobPairStatus = moderrors.Register(
		ModuleName,
		36,
		"Operation conflicts with ClobPair status",
	)
	ErrPerpetualDoesNotExist = moderrors.Register(
		ModuleName,
		37,
		"Perpetual does not exist in state",
	)
	ErrInvalidMsgUpdateClobPair = moderrors.Register(
		ModuleName,
		38,
		"MsgUpdateClobPair is invalid",
	)
	ErrInvalidClobPairUpdate = moderrors.Register(
		ModuleName,
		39,
		"ClobPair update is invalid",
	)

	// Liquidations errors.
	ErrInvalidLiquidationsConfig = moderrors.Register(
		ModuleName,
		1000,
		"Proposed LiquidationsConfig is invalid",
	)
	ErrNoPerpetualPositionsToLiquidate = moderrors.Register(
		ModuleName,
		1001,
		"Subaccount has no perpetual positions to liquidate",
	)
	ErrSubaccountNotLiquidatable = moderrors.Register(
		ModuleName,
		1002,
		"Subaccount is not liquidatable",
	)
	ErrNoOpenPositionForPerpetual = moderrors.Register(
		ModuleName,
		1003,
		"Subaccount does not have an open position for perpetual",
	)
	ErrInvalidLiquidationOrderTotalSize = moderrors.Register(
		ModuleName,
		1004,
		"Liquidation order has invalid size",
	)
	ErrInvalidLiquidationOrderSide = moderrors.Register(
		ModuleName,
		1005,
		"Liquidation order is on the wrong side",
	)
	ErrTotalFillAmountExceedsOrderSize = moderrors.Register(
		ModuleName,
		1006,
		"Total fills amount exceeds size of liquidation order",
	)
	ErrLiquidationContainsNoFills = moderrors.Register(
		ModuleName,
		1007,
		"Liquidation order does not contain any fills",
	)
	ErrSubaccountHasLiquidatedPerpetual = moderrors.Register(
		ModuleName,
		1008,
		"Subaccount has previously liquidated this perpetual in the current block",
	)
	ErrLiquidationOrderSizeSmallerThanMin = moderrors.Register(
		ModuleName,
		1009,
		"Liquidation order has size smaller than min position notional specified in the liquidation config",
	)
	ErrLiquidationOrderSizeGreaterThanMax = moderrors.Register(
		ModuleName,
		1010,
		"Liquidation order has size greater than max position notional specified in the liquidation config",
	)
	ErrLiquidationExceedsSubaccountMaxNotionalLiquidated = moderrors.Register(
		ModuleName,
		1011,
		"Liquidation exceeds the maximum notional amount that a single subaccount can have liquidated per block",
	)
	ErrLiquidationExceedsSubaccountMaxInsuranceLost = moderrors.Register(
		ModuleName,
		1012,
		"Liquidation exceeds the maximum insurance fund payout amount for a given subaccount per block",
	)
	ErrInsuranceFundHasInsufficientFunds = moderrors.Register(
		ModuleName,
		1013,
		"Insurance fund does not have sufficient funds to cover liquidation losses",
	)
	ErrInvalidPerpetualPositionSizeDelta = moderrors.Register(
		ModuleName,
		1014,
		"Invalid perpetual position size delta",
	)
	ErrInvalidQuantumsForInsuranceFundDeltaCalculation = moderrors.Register(
		ModuleName,
		1015,
		"Invalid delta base and/or quote quantums for insurance fund delta calculation",
	)
	ErrEmptyDeleveragingFills = moderrors.Register(
		ModuleName,
		1016,
		"Deleveraging fills length must be greater than 0",
	)
	ErrDeleveragingAgainstSelf = moderrors.Register(
		ModuleName,
		1017,
		"Cannot deleverage subaccount against itself",
	)
	ErrDuplicateDeleveragingFillSubaccounts = moderrors.Register(
		ModuleName,
		1018,
		"Deleveraging match cannot have fills with same id",
	)
	ErrZeroDeleveragingFillAmount = moderrors.Register(
		ModuleName,
		1019,
		"Deleveraging match cannot have fills with zero amount",
	)
	ErrPositionCannotBeFullyOffset = moderrors.Register(
		ModuleName,
		1020,
		"Position cannot be fully offset",
	)

	// Advanced order type errors.
	ErrFokOrderCouldNotBeFullyFilled = moderrors.Register(
		ModuleName,
		2000,
		"FillOrKill order could not be fully filled",
	)
	ErrReduceOnlyWouldIncreasePositionSize = moderrors.Register(
		ModuleName,
		2001,
		"Reduce-only orders cannot increase the position size",
	)
	ErrReduceOnlyWouldChangePositionSide = moderrors.Register(
		ModuleName,
		2002,
		"Reduce-only orders cannot change the position side",
	)
	ErrPostOnlyWouldCrossMakerOrder = moderrors.Register(
		ModuleName,
		2003,
		"Post-only order would cross one or more maker orders",
	)

	// Stateful order errors.
	ErrInvalidOrderFlag = moderrors.Register(
		ModuleName,
		3000,
		"Invalid order flags",
	)
	ErrInvalidStatefulOrderGoodTilBlockTime = moderrors.Register(
		ModuleName,
		3001,
		"Invalid order goodTilBlockTime",
	)
	ErrStatefulOrdersCannotRequireImmediateExecution = moderrors.Register(
		ModuleName,
		3002,
		"Stateful orders cannot require immediate execution",
	)
	ErrTimeExceedsGoodTilBlockTime = moderrors.Register(
		ModuleName,
		3003,
		"The block time is greater than the GoodTilBlockTime of the message",
	)
	ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow = moderrors.Register(
		ModuleName,
		3004,
		"The GoodTilBlockTime of the message is further than StatefulOrderTimeWindow into the future",
	)
	ErrStatefulOrderAlreadyExists = moderrors.Register(
		ModuleName,
		3005,
		"Existing stateful order has higher-or-equal priority than the new one",
	)
	ErrStatefulOrderDoesNotExist = moderrors.Register(
		ModuleName,
		3006,
		"Stateful order does not exist",
	)
	ErrStatefulOrderCollateralizationCheckFailed = moderrors.Register(
		ModuleName,
		3007,
		"Stateful order collateralization check failed",
	)
	ErrStatefulOrderPreviouslyCancelled = moderrors.Register(
		ModuleName,
		3008,
		"Stateful order was previously cancelled and therefore cannot be placed",
	)
	ErrStatefulOrderPreviouslyRemoved = moderrors.Register(
		ModuleName,
		3009,
		"Stateful order was previously removed and therefore cannot be placed",
	)

	// Operations Queue validation errors
	ErrInvalidMsgProposedOperations = moderrors.Register(
		ModuleName,
		4000,
		"MsgProposedOperations is invalid",
	)
	ErrInvalidMatchOrder = moderrors.Register(
		ModuleName,
		4001,
		"Match Order is invalid",
	)
	ErrOrderPlacementNotInOperationsQueue = moderrors.Register(
		ModuleName,
		4002,
		"Order was not previously placed in operations queue",
	)
	ErrFillAmountIsZero = moderrors.Register(
		ModuleName,
		4003,
		"Fill amount cannot be zero",
	)
	ErrInvalidDeleveragingFill = moderrors.Register(
		ModuleName,
		4004,
		"Deleveraging fill is invalid",
	)
	ErrInvalidDeleveragedSubaccount = moderrors.Register(
		ModuleName,
		4005,
		"Deleveraged subaccount in proposed deleveraged operation failed deleveraging validation",
	)

	// Block rate limit errors.
	ErrInvalidBlockRateLimitConfig = moderrors.Register(
		ModuleName,
		5000,
		"Proposed BlockRateLimitConfig is invalid",
	)
	ErrBlockRateLimitExceeded = moderrors.Register(
		ModuleName,
		5001,
		"Block rate limit exceeded",
	)

	// Conditional order errors.
	ErrInvalidConditionType = moderrors.Register(
		ModuleName,
		6000,
		"Conditional type is invalid",
	)
	ErrInvalidConditionalOrderTriggerSubticks = moderrors.Register(
		ModuleName,
		6001,
		"Conditional order trigger subticks is invalid",
	)

	// Errors for unimplemented and disabled functionality.
	ErrAssetOrdersNotImplemented = moderrors.Register(
		ModuleName,
		9000,
		"Asset orders are not implemented",
	)
	ErrAssetUpdateNotImplemented = moderrors.Register(
		ModuleName,
		9001,
		"Updates for assets other than USDC are not implemented",
	)
	ErrNotImplemented = moderrors.Register(
		ModuleName,
		9002,
		"This function is not implemented",
	)
	ErrReduceOnlyDisabled = moderrors.Register(
		ModuleName,
		9003,
		"Reduce-only is currently disabled",
	)

	// Equity tier limit errors.
	ErrInvalidEquityTierLimitConfig = moderrors.Register(
		ModuleName,
		10000,
		"Proposed EquityTierLimitConfig is invalid",
	)
	ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit = moderrors.Register(
		ModuleName,
		10001,
		"Subaccount cannot open more orders due to equity tier limit.",
	)
)
