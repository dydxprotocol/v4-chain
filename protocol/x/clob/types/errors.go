package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/clob module sentinel errors
// TODO(CLOB-553) Clean up sentinel errors not in use.
var (
	ErrMemClobOrderDoesNotExist = errorsmod.Register(
		ModuleName,
		2,
		"Order does not exist in memclob",
	)
	ErrInvalidOrderSide = errorsmod.Register(
		ModuleName,
		3,
		"Invalid order side",
	)
	ErrInvalidOrderQuantums     = errorsmod.Register(ModuleName, 4, "Invalid order quantums")
	ErrInvalidOrderGoodTilBlock = errorsmod.Register(ModuleName, 5, "Invalid order goodTilBlock")
	ErrInvalidOrderSubticks     = errorsmod.Register(ModuleName, 6, "Invalid order subticks")
	ErrOrderIsCanceled          = errorsmod.Register(
		ModuleName,
		7,
		"Attempt to place an order that is already canceled",
	)
	ErrInvalidClob = errorsmod.Register(
		ModuleName,
		8,
		"ClobPair ID does not reference a valid CLOB",
	)
	ErrMemClobCancelAlreadyExists = errorsmod.Register(
		ModuleName,
		9,
		"A cancel already exists in the memclob for this order with a greater than or equal GoodTilBlock",
	)
	ErrHeightExceedsGoodTilBlock = errorsmod.Register(
		ModuleName,
		10,
		"The next block height is greater than the GoodTilBlock of the message",
	)
	ErrGoodTilBlockExceedsShortBlockWindow = errorsmod.Register(
		ModuleName,
		11,
		"The GoodTilBlock of the message is further than ShortBlockWindow blocks into the future",
	)
	ErrInvalidPlaceOrder = errorsmod.Register(
		ModuleName,
		12,
		"MsgPlaceOrder is invalid",
	)
	ErrInvalidMsgProposedMatchOrders = errorsmod.Register(
		ModuleName,
		13,
		"MsgProposedMatchOrders is invalid",
	)
	ErrStateFilledAmountNoChange = errorsmod.Register(
		ModuleName,
		14,
		"State filled amount cannot be unchanged",
	)
	ErrStateFilledAmountDecreasing = errorsmod.Register(
		ModuleName,
		15,
		"State filled amount cannot decrease",
	)
	ErrInvalidPruneStateFilledAmount = errorsmod.Register(
		ModuleName,
		16,
		"Cannot prune state fill amount that does not exist",
	)
	ErrOrderWouldExceedMaxOpenOrdersPerClobAndSide = errorsmod.Register(
		ModuleName,
		17,
		"Subaccount cannot open more than 20 orders on a given CLOB and side",
	)
	ErrFillAmountNotDivisibleByStepSize = errorsmod.Register(
		ModuleName,
		18,
		"`FillAmount` is not divisible by `StepBaseQuantums` of the specified `ClobPairId`",
	)
	ErrNoClobPairForPerpetual = errorsmod.Register(
		ModuleName,
		19,
		"The provided perpetual ID does not have any associated CLOB pairs",
	)
	ErrInvalidReplacement = errorsmod.Register(
		ModuleName,
		20,
		"Replacing an existing order failed",
	)
	ErrClobPairAndPerpetualDoNotMatch = errorsmod.Register(
		ModuleName,
		21,
		"Clob pair and perpetual ids do not match",
	)
	ErrMatchedOrderNegativeFee = errorsmod.Register(
		ModuleName,
		22,
		"Matched order has negative fee",
	)
	ErrSubaccountFeeTransferFailed = errorsmod.Register(
		ModuleName,
		23,
		"Subaccounts updated for a matched order, but fee transfer to fee-collector failed",
	)
	ErrOrderFullyFilled = errorsmod.Register(
		ModuleName,
		24,
		"Order is fully filled",
	)
	ErrPremiumWithNonPerpetualClobPair = errorsmod.Register(
		ModuleName,
		25,
		"Attempting to get price premium with a non-perpetual CLOB pair",
	)
	ErrZeroIndexPriceForPremiumCalculation = errorsmod.Register(
		ModuleName,
		26,
		"Index price is zero when calculating price premium",
	)
	ErrInvalidClobPairParameter = errorsmod.Register(
		ModuleName,
		27,
		"Invalid ClobPair parameter",
	)
	ErrZeroPriceForOracle = errorsmod.Register(
		ModuleName,
		28,
		"Oracle price must be > 0.",
	)
	ErrInvalidStatefulOrderCancellation = errorsmod.Register(
		ModuleName,
		29,
		"Invalid stateful order cancellation",
	)
	ErrOrderReprocessed = errorsmod.Register(
		ModuleName,
		30,
		"An order with the same `OrderId` and `OrderHash` has already been processed for this CLOB",
	)
	ErrMissingMidPrice = errorsmod.Register(
		ModuleName,
		31,
		"Missing mid price for ClobPair",
	)
	ErrStatefulOrderCancellationAlreadyExists = errorsmod.Register(
		ModuleName,
		32,
		"Existing stateful order cancellation has higher-or-equal priority than the new one",
	)
	ErrClobPairAlreadyExists = errorsmod.Register(
		ModuleName,
		33,
		"ClobPair with id already exists",
	)
	ErrOrderConflictsWithClobPairStatus = errorsmod.Register(
		ModuleName,
		34,
		"Order conflicts with ClobPair status",
	)
	ErrInvalidClobPairStatusTransition = errorsmod.Register(
		ModuleName,
		35,
		"Invalid ClobPair status transition",
	)
	ErrOperationConflictsWithClobPairStatus = errorsmod.Register(
		ModuleName,
		36,
		"Operation conflicts with ClobPair status",
	)
	ErrPerpetualDoesNotExist = errorsmod.Register(
		ModuleName,
		37,
		"Perpetual does not exist in state",
	)
	ErrInvalidClobPairUpdate = errorsmod.Register(
		ModuleName,
		39,
		"ClobPair update is invalid",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		40,
		"Authority is invalid",
	)
	ErrPerpetualAssociatedWithExistingClobPair = errorsmod.Register(
		ModuleName,
		41,
		"perpetual ID is already associated with an existing CLOB pair",
	)
	ErrUnexpectedTimeInForce = errorsmod.Register(
		ModuleName,
		42,
		"Unexpected time in force",
	)
	ErrOrderHasRemainingSize = errorsmod.Register(
		ModuleName,
		43,
		"Order has remaining size",
	)
	ErrInvalidTimeInForce = errorsmod.Register(
		ModuleName,
		44,
		"invalid time in force",
	)
	ErrInvalidBatchCancel = errorsmod.Register(
		ModuleName,
		45,
		"Invalid batch cancel message",
	)
	ErrBatchCancelFailed = errorsmod.Register(
		ModuleName,
		46,
		"Batch cancel has failed",
	)
	ErrClobNotInitialized = errorsmod.Register(
		ModuleName,
		47,
		"CLOB has not been initialized",
	)
	ErrDeprecatedField = errorsmod.Register(
		ModuleName,
		48,
		"This field has been deprecated",
	)
	ErrInvalidTwapOrderPlacement = errorsmod.Register(
		ModuleName,
		49,
		"Invalid TWAP order placement",
	)
	ErrInvalidBuilderCode = errorsmod.Register(
		ModuleName,
		50,
		"Invalid builder code",
	)

	// Liquidations errors.
	ErrInvalidLiquidationsConfig = errorsmod.Register(
		ModuleName,
		1000,
		"Proposed LiquidationsConfig is invalid",
	)
	ErrNoPerpetualPositionsToLiquidate = errorsmod.Register(
		ModuleName,
		1001,
		"Subaccount has no perpetual positions to liquidate",
	)
	ErrSubaccountNotLiquidatable = errorsmod.Register(
		ModuleName,
		1002,
		"Subaccount is not liquidatable",
	)
	ErrNoOpenPositionForPerpetual = errorsmod.Register(
		ModuleName,
		1003,
		"Subaccount does not have an open position for perpetual",
	)
	ErrInvalidLiquidationOrderTotalSize = errorsmod.Register(
		ModuleName,
		1004,
		"Liquidation order has invalid size",
	)
	ErrInvalidLiquidationOrderSide = errorsmod.Register(
		ModuleName,
		1005,
		"Liquidation order is on the wrong side",
	)
	ErrTotalFillAmountExceedsOrderSize = errorsmod.Register(
		ModuleName,
		1006,
		"Total fills amount exceeds size of liquidation order",
	)
	ErrLiquidationContainsNoFills = errorsmod.Register(
		ModuleName,
		1007,
		"Liquidation order does not contain any fills",
	)
	ErrSubaccountHasLiquidatedPerpetual = errorsmod.Register(
		ModuleName,
		1008,
		"Subaccount has previously liquidated this perpetual in the current block",
	)
	ErrLiquidationOrderSizeSmallerThanMin = errorsmod.Register(
		ModuleName,
		1009,
		"Liquidation order has size smaller than min position notional specified in the liquidation config",
	)
	ErrLiquidationOrderSizeGreaterThanMax = errorsmod.Register(
		ModuleName,
		1010,
		"Liquidation order has size greater than max position notional specified in the liquidation config",
	)
	ErrLiquidationExceedsSubaccountMaxNotionalLiquidated = errorsmod.Register(
		ModuleName,
		1011,
		"Liquidation exceeds the maximum notional amount that a single subaccount can have liquidated per block",
	)
	ErrLiquidationExceedsSubaccountMaxInsuranceLost = errorsmod.Register(
		ModuleName,
		1012,
		"Liquidation exceeds the maximum insurance fund payout amount for a given subaccount per block",
	)
	ErrInsuranceFundHasInsufficientFunds = errorsmod.Register(
		ModuleName,
		1013,
		"Insurance fund does not have sufficient funds to cover liquidation losses",
	)
	ErrInvalidPerpetualPositionSizeDelta = errorsmod.Register(
		ModuleName,
		1014,
		"Invalid perpetual position size delta",
	)
	ErrInvalidQuantumsForInsuranceFundDeltaCalculation = errorsmod.Register(
		ModuleName,
		1015,
		"Invalid delta base and/or quote quantums for insurance fund delta calculation",
	)
	// TODO: Should the error code be skipped or re-assigned?
	ErrDeleveragingAgainstSelf = errorsmod.Register(
		ModuleName,
		1017,
		"Cannot deleverage subaccount against itself",
	)
	ErrDuplicateDeleveragingFillSubaccounts = errorsmod.Register(
		ModuleName,
		1018,
		"Deleveraging match cannot have fills with same id",
	)
	ErrZeroDeleveragingFillAmount = errorsmod.Register(
		ModuleName,
		1019,
		"Deleveraging match cannot have fills with zero amount",
	)
	ErrPositionCannotBeFullyOffset = errorsmod.Register(
		ModuleName,
		1020,
		"Position cannot be fully offset",
	)
	ErrDeleveragingIsFinalSettlementFlagMismatch = errorsmod.Register(
		ModuleName,
		1021,
		"Deleveraging match has incorrect value for isFinalSettlement flag",
	)
	ErrLiquidationConflictsWithClobPairStatus = errorsmod.Register(
		ModuleName,
		1022,
		"Liquidation conflicts with ClobPair status",
	)

	// Advanced order type errors.
	ErrFokOrderCouldNotBeFullyFilled = errorsmod.Register(
		ModuleName,
		2000,
		"FillOrKill order could not be fully filled",
	)
	ErrReduceOnlyWouldIncreasePositionSize = errorsmod.Register(
		ModuleName,
		2001,
		"Reduce-only orders cannot increase the position size",
	)
	ErrReduceOnlyWouldChangePositionSide = errorsmod.Register(
		ModuleName,
		2002,
		"Reduce-only orders cannot change the position side",
	)
	ErrPostOnlyWouldCrossMakerOrder = errorsmod.Register(
		ModuleName,
		2003,
		"Post-only order would cross one or more maker orders",
	)
	ErrImmediateExecutionOrderAlreadyFilled = errorsmod.Register(
		ModuleName,
		2004,
		"IOC order is already filled, remaining size is cancelled.",
	)
	ErrWouldViolateIsolatedSubaccountConstraints = errorsmod.Register(
		ModuleName,
		2005,
		"Order would violate isolated subaccount constraints.",
	)

	// Stateful order errors.
	ErrInvalidOrderFlag = errorsmod.Register(
		ModuleName,
		3000,
		"Invalid order flags",
	)
	ErrInvalidStatefulOrderGoodTilBlockTime = errorsmod.Register(
		ModuleName,
		3001,
		"Invalid order goodTilBlockTime",
	)
	ErrLongTermOrdersCannotRequireImmediateExecution = errorsmod.Register(
		ModuleName,
		3002,
		"Stateful orders cannot require immediate execution",
	)
	ErrTimeExceedsGoodTilBlockTime = errorsmod.Register(
		ModuleName,
		3003,
		"The block time is greater than the GoodTilBlockTime of the message",
	)
	ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow = errorsmod.Register(
		ModuleName,
		3004,
		"The GoodTilBlockTime of the message is further than StatefulOrderTimeWindow into the future",
	)
	ErrStatefulOrderAlreadyExists = errorsmod.Register(
		ModuleName,
		3005,
		"Existing stateful order has higher-or-equal priority than the new one",
	)
	ErrStatefulOrderDoesNotExist = errorsmod.Register(
		ModuleName,
		3006,
		"Stateful order does not exist",
	)
	ErrStatefulOrderCollateralizationCheckFailed = errorsmod.Register(
		ModuleName,
		3007,
		"Stateful order collateralization check failed",
	)
	ErrStatefulOrderPreviouslyCancelled = errorsmod.Register(
		ModuleName,
		3008,
		"Stateful order was previously cancelled and therefore cannot be placed",
	)
	ErrStatefulOrderPreviouslyRemoved = errorsmod.Register(
		ModuleName,
		3009,
		"Stateful order was previously removed and therefore cannot be placed",
	)
	ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder = errorsmod.Register(
		ModuleName,
		3010,
		"Stateful order cancellation failed because the order was already removed from state",
	)

	// Operations Queue validation errors
	ErrInvalidMsgProposedOperations = errorsmod.Register(
		ModuleName,
		4000,
		"MsgProposedOperations is invalid",
	)
	ErrInvalidMatchOrder = errorsmod.Register(
		ModuleName,
		4001,
		"Match Order is invalid",
	)
	ErrOrderPlacementNotInOperationsQueue = errorsmod.Register(
		ModuleName,
		4002,
		"Order was not previously placed in operations queue",
	)
	ErrFillAmountIsZero = errorsmod.Register(
		ModuleName,
		4003,
		"Fill amount cannot be zero",
	)
	ErrInvalidDeleveragingFill = errorsmod.Register(
		ModuleName,
		4004,
		"Deleveraging fill is invalid",
	)
	ErrInvalidDeleveragedSubaccount = errorsmod.Register(
		ModuleName,
		4005,
		"Deleveraged subaccount in proposed deleveraged operation failed deleveraging validation",
	)
	ErrInvalidOrderRemoval = errorsmod.Register(
		ModuleName,
		4006,
		"Order Removal is invalid",
	)
	ErrInvalidOrderRemovalReason = errorsmod.Register(
		ModuleName,
		4007,
		"Order Removal reason is invalid",
	)
	ErrZeroFillDeleveragingForNonNegativeTncSubaccount = errorsmod.Register(
		ModuleName,
		4008,
		"Zero-fill deleveraging operation included in block for non-negative TNC subaccount",
	)

	// Block rate limit errors.
	ErrInvalidBlockRateLimitConfig = errorsmod.Register(
		ModuleName,
		5000,
		"Proposed BlockRateLimitConfig is invalid",
	)
	ErrBlockRateLimitExceeded = errorsmod.Register(
		ModuleName,
		5001,
		"Block rate limit exceeded",
	)

	// Conditional order errors.
	ErrInvalidConditionType = errorsmod.Register(
		ModuleName,
		6000,
		"Conditional type is invalid",
	)
	ErrInvalidConditionalOrderTriggerSubticks = errorsmod.Register(
		ModuleName,
		6001,
		"Conditional order trigger subticks is invalid",
	)
	ErrConditionalOrderUntriggered = errorsmod.Register(
		ModuleName,
		6002,
		"Conditional order is untriggered",
	)

	// Errors for unimplemented and disabled functionality.
	ErrAssetOrdersNotImplemented = errorsmod.Register(
		ModuleName,
		9000,
		"Asset orders are not implemented",
	)
	ErrAssetUpdateNotImplemented = errorsmod.Register(
		ModuleName,
		9001,
		"Updates for assets other than USDC are not implemented",
	)
	ErrNotImplemented = errorsmod.Register(
		ModuleName,
		9002,
		"This function is not implemented",
	)
	ErrReduceOnlyDisabled = errorsmod.Register(
		ModuleName,
		9003,
		"Reduce-only is currently disabled for non-IOC orders",
	)

	// Equity tier limit errors.
	ErrInvalidEquityTierLimitConfig = errorsmod.Register(
		ModuleName,
		10000,
		"Proposed EquityTierLimitConfig is invalid",
	)
	ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit = errorsmod.Register(
		ModuleName,
		10001,
		"Subaccount cannot open more orders due to equity tier limit.",
	)

	// Order router errors.
	ErrInvalidOrderRouterAddress = errorsmod.Register(
		ModuleName,
		11000,
		"Invalid order router address",
	)

	// Leverage update errors.
	ErrInvalidAddress = errorsmod.Register(
		ModuleName,
		11001,
		"Invalid address",
	)
	ErrInvalidLeverage = errorsmod.Register(
		ModuleName,
		11002,
		"Invalid update leverage",
	)
	ErrLeverageExceedsMaximum = errorsmod.Register(
		ModuleName,
		11003,
		"Leverage exceeds maximum allowed for perpetual",
	)
	ErrInitialMarginPpmIsZero = errorsmod.Register(
		ModuleName,
		11004,
		"Initial margin ppm is zero",
	)
)
