package keeper

import (
	"fmt"
	"math/big"
	"runtime/debug"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/mev_telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var MAX_SPREAD_BEFORE_FALLING_BACK_TO_ORACLE = new(big.Rat).SetFrac64(1, 100)

type MevTelemetryConfig struct {
	Enabled    bool
	Hosts      []string
	Identifier string
}

type ClobMetadata struct {
	ClobPair    types.ClobPair
	MidPrice    types.Subticks
	OraclePrice types.Subticks
	BestBid     types.Order
	BestAsk     types.Order
}

// CumulativePnL keeps track of the cumulative PnL for each subaccount per market.
type CumulativePnL struct {
	// PnL calculations.
	SubaccountPnL               map[satypes.SubaccountId]*big.Int
	SubaccountPositionSizeDelta map[satypes.SubaccountId]*big.Int

	// Metadata.
	NumFills            int
	VolumeQuoteQuantums *big.Int

	// Cached fields used in the calculation of PnL.
	// These should not be modified after initialization.
	Metadata              ClobMetadata
	PerpetualFundingIndex *big.Int
}

type PnLCalculationParams struct {
	subaccountId satypes.SubaccountId
	isBuy        bool
	feePpm       int32
}

// RecordMevMetricsIsEnabled returns true if the MEV telemetry config is enabled.
func (k Keeper) RecordMevMetricsIsEnabled() bool {
	return k.mevTelemetryConfig.Enabled
}

// RecordMevMetrics measures and records MEV by comparing the block proposer's list of matches
// with its own list of matches.
func (k Keeper) RecordMevMetrics(
	ctx sdk.Context,
	stakingKeeper process.ProcessStakingKeeper,
	perpetualKeeper process.ProcessPerpetualKeeper,
	msgProposedOperations *types.MsgProposedOperations,
) {
	ctx = log.AddPersistentTagsToLogger(
		ctx,
		log.Module,
		"x/clob/mev_telemetry",
	)
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.MevLatency,
		time.Now(),
	)

	// Recover from any panics that occur during MEV calculation.
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				log.ErrorLog(ctx, "panic when recording mev metrics",
					log.StackTrace,
					string(debug.Stack()),
				)
			} else {
				log.ErrorLogWithError(ctx, "panic when recording mev metrics", err,
					log.StackTrace,
					string(debug.Stack()),
				)
			}
		}
	}()

	clobMetadata := k.GetClobMetadata(ctx)

	// Initialize cumulative PnL for block proposer and validator.
	blockProposerPnL, validatorPnL := k.InitializeCumulativePnLs(
		ctx,
		perpetualKeeper,
		clobMetadata,
	)

	// Calculate the block proposer's PnL from regular and liquidation matches.
	blockProposerMevMatches, err := k.GetMEVDataFromOperations(
		ctx,
		msgProposedOperations.GetOperationsQueue(),
	)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to create MEV matches for block proposer operations", err,
			log.OperationsQueue,
			msgProposedOperations.GetOperationsQueue(),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}
	k.CalculateSubaccountPnLForMevMatches(
		ctx,
		blockProposerPnL,
		blockProposerMevMatches,
	)

	// Calculate the validator's PnL from regular and liquidation matches.
	validatorMevMatches, err := k.GetMEVDataFromOperations(
		ctx,
		k.GetOperations(ctx).GetOperationsQueue(),
	)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to create MEV matches for validator operations", err,
			log.OperationsQueue,
			k.GetOperations(ctx).GetOperationsQueue(),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}
	k.CalculateSubaccountPnLForMevMatches(
		ctx,
		validatorPnL,
		validatorMevMatches,
	)

	// TODO(CLOB-742): re-enable deleveraging and funding in MEV calculation.
	// Calculate Trading PnL for block proposer.
	// if err := k.CalculateSubaccountPnLForMatches(
	// 	ctx,
	// 	blockProposerPnL,
	// 	msgProposedOperations.GetOperationsQueue(),
	// ); err != nil {
	// 	k.Logger(ctx).Error(
	// 		fmt.Sprintf(
	// 			"Failed to calculate PnL for block proposer: Error: %+v, Operations: %+v",
	// 			err.Error(),
	// 			msgProposedOperations.GetOperationsQueue(),
	// 		),
	// 	)
	// 	telemetry.IncrCounter(1, types.ModuleName, metrics.Mev, metrics.Error, metrics.Count)
	// 	return
	// }

	// TODO(CLOB-742): re-enable deleveraging and funding in MEV calculation.
	// // Calculate Trading PnL for the current validator.
	// if err := k.CalculateSubaccountPnLForMatches(
	// 	ctx,
	// 	validatorPnL,
	// 	k.GetOperations(ctx).GetOperationsQueue(),
	// ); err != nil {
	// 	k.Logger(ctx).Error(
	// 		fmt.Sprintf(
	// 			"Failed to calculate PnL for validator: Error: %+v, Operations: %+v",
	// 			err.Error(),
	// 			k.GetOperations(ctx).GetOperationsQueue(),
	// 		),
	// 	)
	// 	telemetry.IncrCounter(1, types.ModuleName, metrics.Mev, metrics.Error, metrics.Count)
	// 	return
	// }

	// Since `MaybeProcessNewFundingTickEpoch` modifies state, operate on a cached context
	// so that state updates are not persisted.
	// Funding indices for perpetuals are updated in the EndBlocker.
	// In order to measure MEV for funding received/paid in this block, update
	// the funding indices before processing operations.

	// TODO(CLOB-742): re-enable deleveraging and funding in MEV calculation.
	// cacheCtx, _ := ctx.CacheContext()
	// perpetualKeeper.MaybeProcessNewFundingTickEpoch(cacheCtx)

	// // Calculate funding using the position size delta.
	// for _, cumulativePnL := range []map[types.ClobPairId]*CumulativePnL{blockProposerPnL, validatorPnL} {
	// 	// Funding payments are additive, so BP's funding PnL is equal to the sum of funding from initial
	// 	// position size and funding from position size delta from the proposed queue.
	// 	// Similarly, validator's funding PnL is equal to the sum of funding from initial position size
	// 	// and funding from position size delta from validator's queue.
	// 	// The funding payment from initial position size cancel out, and we only need to calculate
	// 	// the funding from position size delta.
	// 	if err := k.AddSettlementForPositionDelta(
	// 		ctx,
	// 		perpetualKeeper,
	// 		cumulativePnL,
	// 	); err != nil {
	// 		k.Logger(ctx).Error(err.Error())
	// 		telemetry.IncrCounter(1, types.ModuleName, metrics.Mev, metrics.Error, metrics.Count)
	// 		return
	// 	}
	// }

	// Add label for consensus round if available.
	consensusRound, ok := ctx.Value(process.ConsensusRound).(int64)
	if !ok {
		log.ErrorLog(ctx, "Failed to get consensus round")
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}

	// Add label for the block proposer.
	proposerConsAddress := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	proposer, err := stakingKeeper.GetValidatorByConsAddr(ctx, proposerConsAddress)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get proposer by consensus address", err,
			log.Proposer, proposerConsAddress.String(),
		)

		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}

	validatorVolumeQuoteQuantumsPerMarket := make(map[types.ClobPairId]*big.Int, 0)
	mevPerMarket := make(map[types.ClobPairId]float32, 0)

	for clobPairId, blockProposerSubaccountPnL := range blockProposerPnL {
		// Calculate MEV for the given market.
		mev, _ := blockProposerSubaccountPnL.CalculateMev(validatorPnL[clobPairId]).Float32()

		validatorVolumeQuoteQuantums := new(big.Int).Div(
			validatorPnL[clobPairId].VolumeQuoteQuantums,
			big.NewInt(2),
		)

		// Log MEV metric.
		// TODO(CLOB-1051) change to use new logger library. Be careful the values are not changed
		// because mev dashboards rely on this log.
		k.Logger(ctx).Info(
			"Measuring MEV for proposed matches",
			metrics.Mev,
			mev,
			// Common metadata.
			metrics.BlockHeight,
			ctx.BlockHeight(),
			metrics.ConsensusRound,
			consensusRound,
			metrics.Proposer,
			proposer.Description.Moniker,
			metrics.ClobPairId,
			clobPairId.ToUint32(),
			metrics.MidPrice,
			validatorPnL[clobPairId].Metadata.MidPrice.ToUint64(),
			metrics.OraclePrice,
			validatorPnL[clobPairId].Metadata.OraclePrice.ToUint64(),
			metrics.BestBid,
			fmt.Sprintf("%+v", validatorPnL[clobPairId].Metadata.BestBid),
			metrics.BestAsk,
			fmt.Sprintf("%+v", validatorPnL[clobPairId].Metadata.BestAsk),
			// Validator stats.
			metrics.ValidatorNumFills,
			validatorPnL[clobPairId].NumFills,
			metrics.ValidatorVolumeQuoteQuantums,
			validatorVolumeQuoteQuantums.String(),
			// Proposer stats.
			metrics.ProposerNumFills,
			blockProposerPnL[clobPairId].NumFills,
			metrics.ProposerVolumeQuoteQuantums,
			new(big.Int).Div(
				blockProposerPnL[clobPairId].VolumeQuoteQuantums,
				big.NewInt(2),
			).String(),
		)

		metrics.SetGaugeWithLabels(
			metrics.ClobMev,
			mev,
			metrics.GetLabelForStringValue(
				metrics.Proposer,
				proposer.Description.Moniker,
			),
			metrics.GetLabelForIntValue(
				metrics.ClobPairId,
				int(clobPairId.ToUint32()),
			),
		)

		validatorVolumeQuoteQuantumsPerMarket[clobPairId] = new(big.Int).Div(
			validatorPnL[clobPairId].VolumeQuoteQuantums,
			big.NewInt(2),
		)
		mevPerMarket[clobPairId] = mev
	}

	if len(k.mevTelemetryConfig.Hosts) != 0 {
		mevClobMidPrices := make([]types.ClobMidPrice, 0, len(clobMetadata))
		for _, metadata := range clobMetadata {
			mevClobMidPrices = append(
				mevClobMidPrices,
				types.ClobMidPrice{
					ClobPair: metadata.ClobPair,
					Subticks: metadata.MidPrice.ToUint64(),
				},
			)
		}
		go mev_telemetry.SendDatapoints(
			ctx,
			k.mevTelemetryConfig.Hosts,
			types.MevMetrics{
				MevDatapoint: types.MEVDatapoint{
					Height:              lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
					ChainID:             ctx.ChainID(),
					VolumeQuoteQuantums: validatorVolumeQuoteQuantumsPerMarket,
					MEV:                 mevPerMarket,
					Identifier:          k.mevTelemetryConfig.Identifier,
				},
				MevNodeToNode: types.MevNodeToNodeMetrics{
					ValidatorMevMatches: validatorMevMatches,
					ClobMidPrices:       mevClobMidPrices,
					BpMevMatches:        blockProposerMevMatches,
					ProposalReceiveTime: uint64(time.Now().Second()),
				},
			},
		)
	}
}

// GetClobMetadata fetches the mid prices for all CLOB pairs and the CLOB pairs themselves.
// This function falls back to use the oracle price if any of the mid prices are missing
// or if the spread is greater-than-or-equal-to the max spread.
func (k Keeper) GetClobMetadata(
	ctx sdk.Context,
) (
	clobMetadata map[types.ClobPairId]ClobMetadata,
) {
	clobMetadata = make(map[types.ClobPairId]ClobMetadata)

	for _, clobPair := range k.GetAllClobPairs(ctx) {
		clobPairId := clobPair.GetClobPairId()

		midPriceSubticks, bestBid, bestAsk, exist := k.MemClob.GetMidPrice(ctx, clobPairId)
		oraclePriceSubticksRat := k.GetOraclePriceSubticksRat(ctx, clobPair)
		// Consistently round down here.
		oraclePriceSubticksInt := lib.BigRatRound(oraclePriceSubticksRat, false)
		if !oraclePriceSubticksInt.IsUint64() {
			panic(
				fmt.Sprintf(
					"GetAllMidPrices: invalid oracle price %+v for clob pair %+v",
					oraclePriceSubticksInt,
					clobPair,
				),
			)
		}
		oraclePriceSubticks := types.Subticks(oraclePriceSubticksInt.Uint64())

		// Use the oracle price instead of the mid price if the mid price doesn't exist or
		// the spread is greater-than-or-equal-to the max spread.
		if !exist || new(big.Rat).SetFrac(
			new(big.Int).SetUint64(uint64(bestAsk.Subticks-bestBid.Subticks)),
			new(big.Int).SetUint64(uint64(bestBid.Subticks)), // Note that bestBid cannot be 0 if exist is true.
		).Cmp(MAX_SPREAD_BEFORE_FALLING_BACK_TO_ORACLE) >= 0 {
			metrics.IncrCounterWithLabels(
				metrics.MevFallbackToOracle,
				1,
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(clobPairId.ToUint32())),
			)
			midPriceSubticks = oraclePriceSubticks
		}

		// Set the CLOB metadata.
		clobMetadata[clobPairId] = ClobMetadata{
			ClobPair:    clobPair,
			MidPrice:    midPriceSubticks,
			OraclePrice: oraclePriceSubticks,
			BestBid:     bestBid,
			BestAsk:     bestAsk,
		}
	}

	return clobMetadata
}

// InitializeCumulativePnLs initializes the cumulative PnLs for the block proposer and the
// current validator.
func (k Keeper) InitializeCumulativePnLs(
	ctx sdk.Context,
	perpetualKeeper process.ProcessPerpetualKeeper,
	clobMetadata map[types.ClobPairId]ClobMetadata,
) (
	blockProposerPnL map[types.ClobPairId]*CumulativePnL,
	validatorPnL map[types.ClobPairId]*CumulativePnL,
) {
	blockProposerPnL = make(map[types.ClobPairId]*CumulativePnL)
	validatorPnL = make(map[types.ClobPairId]*CumulativePnL)

	for clobPairId, metadata := range clobMetadata {
		clobPair := metadata.ClobPair
		// Get a mapping from perpetual Id to current perpetual funding index.
		perpetual, err := perpetualKeeper.GetPerpetual(ctx, clobPair.MustGetPerpetualId())
		if err != nil {
			panic(perptypes.ErrPerpetualDoesNotExist)
		}

		for _, cumulativePnL := range []map[types.ClobPairId]*CumulativePnL{
			blockProposerPnL,
			validatorPnL,
		} {
			cumulativePnL[clobPairId] = &CumulativePnL{
				SubaccountPnL:               make(map[satypes.SubaccountId]*big.Int),
				SubaccountPositionSizeDelta: make(map[satypes.SubaccountId]*big.Int),
				NumFills:                    0,
				VolumeQuoteQuantums:         big.NewInt(0),
				Metadata:                    metadata,
				PerpetualFundingIndex:       perpetual.FundingIndex.BigInt(),
			}
		}
	}
	return blockProposerPnL, validatorPnL
}

// GetMEVDataFromOperations returns the MEV matches and MEV liquidations from the provided
// operations queue. It returns an error if a short-term order cannot be decoded. Panics if
// an order cannot be found.
func (k Keeper) GetMEVDataFromOperations(
	ctx sdk.Context,
	operations []types.OperationRaw,
) (
	validatorMevMatches *types.ValidatorMevMatches,
	err error,
) {
	// Collect all the short-term orders placed for subsequent lookups.
	placedShortTermOrders := make(map[types.OrderId]types.Order)

	// Populate `mevMatches` and `mevLiquidationMatches` from the local validator's match operations.
	mevMatches := make([]types.MEVMatch, 0)
	mevLiquidationMatches := make([]types.MEVLiquidationMatch, 0)
	affiliateParameters, err := k.affiliatesKeeper.GetAffiliateParameters(ctx)
	if err != nil {
		return nil, err
	}
	for _, operation := range operations {
		switch typedOperation := operation.Operation.(type) {
		case *types.OperationRaw_ShortTermOrderPlacement:
			// Decode the short-term order for subsequent lookups.
			// Note we don't fetch the CLOB pair since it must be included in a match,
			// so we can fetch it when encoding the match.
			bytes := typedOperation.ShortTermOrderPlacement
			tx, err := k.txDecoder(bytes)
			if err != nil {
				return nil, err
			}
			msgPlaceOrder := tx.GetMsgs()[0].(*types.MsgPlaceOrder)
			order := msgPlaceOrder.Order
			placedShortTermOrders[order.GetOrderId()] = order
		case *types.OperationRaw_Match:
			switch match := typedOperation.Match.Match.(type) {
			case *types.ClobMatch_MatchOrders:
				matchOrders := match.MatchOrders

				// Add a MEV match for each fill.
				takerOrder := k.MustFetchOrderFromOrderId(ctx, matchOrders.TakerOrderId, placedShortTermOrders)
				for _, fill := range matchOrders.Fills {
					makerOrder := k.MustFetchOrderFromOrderId(ctx, fill.MakerOrderId, placedShortTermOrders)
					mevMatch := types.MEVMatch{
						TakerOrderSubaccountId: &takerOrder.OrderId.SubaccountId,
						TakerFeePpm: k.feeTiersKeeper.GetPerpetualFeePpm(
							ctx,
							takerOrder.GetSubaccountId().Owner,
							true,
							affiliateParameters.RefereeMinimumFeeTierIdx,
							takerOrder.OrderId.ClobPairId,
						),

						MakerOrderSubaccountId: &makerOrder.OrderId.SubaccountId,
						MakerOrderSubticks:     makerOrder.Subticks,
						MakerOrderIsBuy:        makerOrder.IsBuy(),
						MakerFeePpm: k.feeTiersKeeper.GetPerpetualFeePpm(
							ctx,
							makerOrder.GetSubaccountId().Owner,
							false,
							affiliateParameters.RefereeMinimumFeeTierIdx,
							takerOrder.OrderId.ClobPairId,
						),

						ClobPairId: takerOrder.OrderId.ClobPairId,
						FillAmount: fill.FillAmount,
					}
					mevMatches = append(mevMatches, mevMatch)
				}

			case *types.ClobMatch_MatchPerpetualLiquidation:
				matchLiquidation := match.MatchPerpetualLiquidation

				for _, fill := range matchLiquidation.Fills {
					makerOrder := k.MustFetchOrderFromOrderId(ctx, fill.MakerOrderId, placedShortTermOrders)

					// Calculate the insurance fund delta for this trade.
					liquidationIsBuy := !makerOrder.IsBuy()
					insuranceFundDelta, err := k.GetLiquidationInsuranceFundDelta(
						ctx,
						matchLiquidation.Liquidated,
						matchLiquidation.PerpetualId,
						liquidationIsBuy,
						fill.FillAmount,
						makerOrder.GetOrderSubticks(),
					)
					if err != nil {
						return nil, err
					}

					// `insuranceFundDelta` is measured in int64 quote quantums.
					// It represents up to ~9 trillion USDC which should always be enough for insurance fund delta.
					// We explicitly panic if there's an int64 overflow.
					if !insuranceFundDelta.IsInt64() {
						panic(fmt.Sprintf("insurance fund delta (%v) is not an int64", insuranceFundDelta.String()))
					}

					mevLiquidationMatch := types.MEVLiquidationMatch{
						LiquidatedSubaccountId: matchLiquidation.Liquidated,
						// TODO(CLOB-957): Use `SerializableInt` for insurance fund delta
						InsuranceFundDeltaQuoteQuantums: insuranceFundDelta.Int64(),

						MakerOrderSubaccountId: makerOrder.OrderId.SubaccountId,
						MakerOrderSubticks:     makerOrder.Subticks,
						MakerOrderIsBuy:        makerOrder.IsBuy(),
						MakerFeePpm: k.feeTiersKeeper.GetPerpetualFeePpm(
							ctx,
							makerOrder.GetSubaccountId().Owner,
							false,
							affiliateParameters.RefereeMinimumFeeTierIdx,
							matchLiquidation.ClobPairId,
						),

						ClobPairId: matchLiquidation.ClobPairId,
						FillAmount: fill.FillAmount,
					}
					mevLiquidationMatches = append(mevLiquidationMatches, mevLiquidationMatch)
				}
			case *types.ClobMatch_MatchPerpetualDeleveraging:
				// TODO: Encode deleveraging matches into a separate struct.
			}
		}
	}

	validatorMevMatches = &types.ValidatorMevMatches{}
	validatorMevMatches.Matches = mevMatches
	validatorMevMatches.LiquidationMatches = mevLiquidationMatches
	return validatorMevMatches, nil
}

// CalculateSubaccountPnLForMevMatches calculates the PnL for each subaccount for the given
// matches. It returns an error if any of the match CLOB pairs do not exist or `AddPnLForTradeWithFilledSubticks`
// returns an error.
func (k Keeper) CalculateSubaccountPnLForMevMatches(
	ctx sdk.Context,
	clobPairToPnLs map[types.ClobPairId]*CumulativePnL,
	matches *types.ValidatorMevMatches,
) {
	for _, matchWithOrders := range matches.Matches {
		clobPairId := matchWithOrders.ClobPairId
		cumulativePnL, exists := clobPairToPnLs[types.ClobPairId(clobPairId)]
		if !exists {
			panic(
				fmt.Sprintf(
					"CalculateSubaccountPnLForMevMatches: CLOB pair ID %+v does not exist in clobPairToPnLs",
					clobPairId,
				),
			)
		}

		// Update the PnL for the match.
		for _, p := range []PnLCalculationParams{
			{
				*matchWithOrders.TakerOrderSubaccountId,
				!matchWithOrders.MakerOrderIsBuy,
				matchWithOrders.TakerFeePpm,
			},
			{
				*matchWithOrders.MakerOrderSubaccountId,
				matchWithOrders.MakerOrderIsBuy,
				matchWithOrders.MakerFeePpm,
			},
		} {
			cumulativePnL.AddPnLForTradeWithFilledSubticks(
				p.subaccountId,
				p.isBuy,
				types.Subticks(matchWithOrders.MakerOrderSubticks),
				satypes.BaseQuantums(matchWithOrders.FillAmount),
				p.feePpm,
			)
		}

		cumulativePnL.NumFills += 1
	}

	// Calculate MEV for liquidation matches.
	for _, mevLiquidation := range matches.LiquidationMatches {
		clobPairId := mevLiquidation.ClobPairId
		cumulativePnL, exists := clobPairToPnLs[types.ClobPairId(clobPairId)]
		if !exists {
			panic(
				fmt.Sprintf(
					"CalculateSubaccountPnLForMevMatches: CLOB pair ID %+v does not exist in clobPairToPnLs",
					clobPairId,
				),
			)
		}

		// Update PnL for the match.
		liquidationIsBuy := !mevLiquidation.MakerOrderIsBuy
		for _, p := range []PnLCalculationParams{
			{mevLiquidation.LiquidatedSubaccountId, liquidationIsBuy, 0},
			{mevLiquidation.MakerOrderSubaccountId, !liquidationIsBuy, mevLiquidation.MakerFeePpm},
		} {
			cumulativePnL.AddPnLForTradeWithFilledSubticks(
				p.subaccountId,
				p.isBuy,
				types.Subticks(mevLiquidation.MakerOrderSubticks),
				satypes.BaseQuantums(mevLiquidation.FillAmount),
				p.feePpm,
			)
		}

		// Note that negative insurance fund delta (insurance fund covers losses) will
		// improve the subaccount's PnL.
		insuranceFundDelta := big.NewInt(mevLiquidation.InsuranceFundDeltaQuoteQuantums)
		cumulativePnL.AddDeltaToSubaccount(
			mevLiquidation.LiquidatedSubaccountId,
			new(big.Int).Neg(insuranceFundDelta),
		)

		cumulativePnL.NumFills += 1
	}
}

// CalculateSubaccountPnLForMatches calculates the PnL for each subaccount for the given matches.
// TODO: Delete this function.
func (k Keeper) CalculateSubaccountPnLForMatches(
	ctx sdk.Context,
	clobPairToPnLs map[types.ClobPairId]*CumulativePnL,
	operations []types.OperationRaw,
) (
	err error,
) {
	// Collect all the short-term orders placed for subsequent lookups.
	placedShortTermOrders := make(map[types.OrderId]types.Order, 0)
	for _, operation := range operations {
		switch typedOperation := operation.Operation.(type) {
		case *types.OperationRaw_ShortTermOrderPlacement:
			// Collect all the short-term orders for subsequent lookups.
			bytes := typedOperation.ShortTermOrderPlacement
			tx, err := k.txDecoder(bytes)
			if err != nil {
				return err
			}
			msgPlaceOrder := tx.GetMsgs()[0].(*types.MsgPlaceOrder)
			order := msgPlaceOrder.Order
			placedShortTermOrders[order.GetOrderId()] = order
		case *types.OperationRaw_Match:
			switch match := typedOperation.Match.Match.(type) {
			case *types.ClobMatch_MatchOrders:
				// Temporarily no-op on calculating MEV for order matches since the calculation
				// is performed in `CalculateSubaccountPnLForMevMatches`. Note that this function
				// will eventually be deleted.
			case *types.ClobMatch_MatchPerpetualLiquidation:
				// Temporarily no-op on calculating MEV for liquidation matches since the calculation
				// is performed in `CalculateSubaccountPnLForMevMatches`. Note that this function
				// will eventually be deleted.
			case *types.ClobMatch_MatchPerpetualDeleveraging:
				// Calculate MEV for deleveraging matches.
				// TODO(CLOB-742): This whole function is currently not being called since deleveraging and funding
				// are excluded from MEV calculations. Re-enable deleveraging and funding in MEV calculation.
				matchDeleveraging := match.MatchPerpetualDeleveraging
				clobPairId, err := k.GetClobPairIdForPerpetual(ctx, matchDeleveraging.PerpetualId)
				if err != nil {
					return err
				}

				cumulativePnL, exists := clobPairToPnLs[clobPairId]
				if !exists {
					return types.ErrInvalidClob
				}

				// Get the liquidated subaccount and its position.
				liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, matchDeleveraging.Liquidated)
				position, _ := liquidatedSubaccount.GetPerpetualPositionForId(matchDeleveraging.PerpetualId)
				isBuy := !position.GetIsLong()

				for _, fill := range matchDeleveraging.Fills {
					deltaQuantums := new(big.Int).SetUint64(fill.FillAmount)
					if !isBuy {
						deltaQuantums.Neg(deltaQuantums)
					}

					// Calculate the delta quote quantums for this fill.
					deltaQuoteQuantums, err := k.GetBankruptcyPriceInQuoteQuantums(
						ctx,
						matchDeleveraging.Liquidated,
						matchDeleveraging.PerpetualId,
						deltaQuantums,
					)
					if err != nil {
						return err
					}
					absQuoteQuantums := new(big.Int).Abs(deltaQuoteQuantums)

					// Update PnL for the match.
					for _, p := range []PnLCalculationParams{
						{matchDeleveraging.Liquidated, isBuy, 0},
						{fill.OffsettingSubaccountId, !isBuy, 0},
					} {
						cumulativePnL.AddPnLForTradeWithFilledQuoteQuantums(
							p.subaccountId,
							p.isBuy,
							absQuoteQuantums,
							satypes.BaseQuantums(fill.FillAmount),
							p.feePpm,
						)
					}
				}
				cumulativePnL.NumFills += len(matchDeleveraging.Fills)
			}
		}
	}
	return nil
}

// AddSettlementForPositionDelta Calculate the total settlement for a subaccount's position delta.
// This function propagates errors from perpetualKeeper.
func (k Keeper) AddSettlementForPositionDelta(
	ctx sdk.Context,
	perpetualKeeper process.ProcessPerpetualKeeper,
	clobPairToPnLs map[types.ClobPairId]*CumulativePnL,
) (err error) {
	for _, cumulativePnL := range clobPairToPnLs {
		perpetualId := cumulativePnL.Metadata.ClobPair.MustGetPerpetualId()
		for subaccountId, deltaQuantums := range cumulativePnL.SubaccountPositionSizeDelta {
			// Get the subaccount and its perpetual positions.
			subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
			var fundingIndex *big.Int
			if position, exists := subaccount.GetPerpetualPositionForId(perpetualId); exists {
				fundingIndex = position.FundingIndex.BigInt()
			} else {
				// Use the perpetual's funding index for newly created positions.
				fundingIndex = new(big.Int).Set(cumulativePnL.PerpetualFundingIndex)
			}

			// Get the funding payment for this position delta.
			perpetual, err := perpetualKeeper.GetPerpetual(ctx, perpetualId)
			if err != nil {
				return err
			}
			bigNetSettlementPpm, _ := perplib.GetSettlementPpmWithPerpetual(
				perpetual,
				deltaQuantums,
				fundingIndex,
			)

			// Add the settlement to the subaccount.
			// Note: Funding payment is the negative of settlement, i.e. positive settlement is equivalent
			// to a negative funding payment (position received funding payment) and vice versa.
			cumulativePnL.AddDeltaToSubaccount(
				subaccountId,
				bigNetSettlementPpm.Div(bigNetSettlementPpm, lib.BigIntOneMillion()),
			)
		}
	}
	return nil
}

// AddPnLForTradeWithFilledSubticks calculates PnL for a given trade using the filled subticks.
// This method calculates the filledQuoteQuantums before calling AddPnLForTradeWithFilledQuoteQuantums.
func (c *CumulativePnL) AddPnLForTradeWithFilledSubticks(
	subaccountId satypes.SubaccountId,
	isBuy bool,
	filledSubticks types.Subticks,
	filledQuantums satypes.BaseQuantums,
	feePpm int32,
) {
	// Get the fill quote quantums using the filled subticks and filled quantums.
	filledQuoteQuantums := types.FillAmountToQuoteQuantums(
		filledSubticks,
		filledQuantums,
		c.Metadata.ClobPair.QuantumConversionExponent,
	)
	c.AddPnLForTradeWithFilledQuoteQuantums(
		subaccountId,
		isBuy,
		filledQuoteQuantums,
		filledQuantums,
		feePpm,
	)
}

// AddPnLForTradeWithFilledQuoteQuantums calculates the PnL for the given trade and adds it to the cumulative PnL.
// The PnL for a buy order is calculated as:
//
//	PnL = n(p - p_mid) - f
//
// The PnL for a sell order is calculated as:
//
//	PnL = n(p_mid - p) - f
//
// where n is the size of the trade, p is the price of the trade, p_mid is the mid price of validator's ordrbook,
// and f is the fee subaccount pays for the matched trade.
// This function returns an error if the clob pair associated with the order is a spot clob pair.
func (c *CumulativePnL) AddPnLForTradeWithFilledQuoteQuantums(
	subaccountId satypes.SubaccountId,
	isBuy bool,
	filledQuoteQuantums *big.Int,
	filledQuantums satypes.BaseQuantums,
	feePpm int32,
) {
	// Get the fill quote quantums using the mid price subticks and filled quantums.
	filledQuoteQuantumsUsingMidPrice := types.FillAmountToQuoteQuantums(
		c.Metadata.MidPrice,
		filledQuantums,
		c.Metadata.ClobPair.QuantumConversionExponent,
	)

	// Calculate PnL for the given subaccount.
	var pnl *big.Int
	if isBuy {
		pnl = new(big.Int).Sub(filledQuoteQuantumsUsingMidPrice, filledQuoteQuantums)
	} else {
		pnl = new(big.Int).Sub(filledQuoteQuantums, filledQuoteQuantumsUsingMidPrice)
	}

	// Calculate fees.
	bigFeeQuoteQuantums := lib.BigMulPpm(filledQuoteQuantums, lib.BigI(feePpm), true)
	pnl.Sub(pnl, bigFeeQuoteQuantums)

	c.AddDeltaToSubaccount(subaccountId, pnl)

	// Update the position size delta.
	deltaQuantums := new(big.Int).SetUint64(filledQuantums.ToUint64())
	if !isBuy {
		deltaQuantums.Neg(deltaQuantums)
	}

	if _, ok := c.SubaccountPositionSizeDelta[subaccountId]; !ok {
		c.SubaccountPositionSizeDelta[subaccountId] = big.NewInt(0)
	}
	c.SubaccountPositionSizeDelta[subaccountId].Add(
		c.SubaccountPositionSizeDelta[subaccountId],
		deltaQuantums,
	)

	c.VolumeQuoteQuantums.Add(c.VolumeQuoteQuantums, filledQuoteQuantums)
}

// AddDeltaToSubaccount adds the given delta to the PnL for the given subaccount.
func (c *CumulativePnL) AddDeltaToSubaccount(
	subaccountId satypes.SubaccountId,
	delta *big.Int,
) {
	if _, ok := c.SubaccountPnL[subaccountId]; !ok {
		c.SubaccountPnL[subaccountId] = big.NewInt(0)
	}
	c.SubaccountPnL[subaccountId].Add(c.SubaccountPnL[subaccountId], delta)
}

// CalculateMev calculates and returns the mev value given the block proposer PnL and the validator PnL,
// using the following formula:
//
//	MEV = 1/2 * Σ|blockProposerPnL - validatorPnL|
//
// Note that this method modifies the receiver.
func (c *CumulativePnL) CalculateMev(other *CumulativePnL) *big.Float {
	// Calculate the PnL difference for each subaccount.
	for subaccountId, pnl := range other.SubaccountPnL {
		if _, ok := c.SubaccountPnL[subaccountId]; !ok {
			c.SubaccountPnL[subaccountId] = big.NewInt(0)
		}
		c.SubaccountPnL[subaccountId].Sub(c.SubaccountPnL[subaccountId], pnl)
	}

	// Sum the absolute value of the PnL for each subaccount.
	sumAbs := big.NewInt(0)
	for _, pnl := range c.SubaccountPnL {
		sumAbs.Add(sumAbs, new(big.Int).Abs(pnl))
	}

	// Calculate mev value.
	mev := new(big.Float).Quo(
		new(big.Float).SetInt(sumAbs),
		new(big.Float).SetUint64(2),
	)
	return mev
}
