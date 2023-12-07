package keeper

import (
	"fmt"
	"math/big"
	"runtime/debug"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/mev_telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type MevTelemetryConfig struct {
	Enabled    bool
	Hosts      []string
	Identifier string
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
	ClobPair              types.ClobPair
	MidPriceSubticks      types.Subticks
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
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.MevLatency,
		time.Now(),
	)

	// Recover from any panics that occur during MEV calculation.
	defer func() {
		if r := recover(); r != nil {
			k.Logger(ctx).Error(
				"panic when recording mev metrics",
				"panic",
				r,
				"stack trace",
				string(debug.Stack()),
			)
		}
	}()

	clobMidPrices, clobPairs := k.GetClobMetadata(ctx)

	// Initialize cumulative PnL for block proposer and validator.
	blockProposerPnL, validatorPnL := k.InitializeCumulativePnLs(
		ctx,
		perpetualKeeper,
		clobMidPrices,
		clobPairs,
	)

	// Calculate the block proposer's PnL from regular and liquidation matches.
	blockProposerMevMatches, err := k.GetMEVDataFromOperations(
		ctx,
		msgProposedOperations.GetOperationsQueue(),
		clobPairs,
	)
	if err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"Failed to create MEV matches for block proposer operations: Error: %+v, Operations: %+v",
				err.Error(),
				msgProposedOperations.GetOperationsQueue(),
			),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}
	if err := k.CalculateSubaccountPnLForMevMatches(
		ctx,
		blockProposerPnL,
		blockProposerMevMatches,
	); err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"Failed to calculate match PnL for block proposer: Error: %+v, MEV matches: %+v",
				err.Error(),
				blockProposerMevMatches,
			),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}

	// Calculate the validator's PnL from regular and liquidation matches.
	validatorMevMatches, err := k.GetMEVDataFromOperations(
		ctx,
		k.GetOperations(ctx).GetOperationsQueue(),
		clobPairs,
	)
	if err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"Failed to create MEV matches for validator operations: Error: %+v, Operations: %+v",
				err.Error(),
				k.GetOperations(ctx).GetOperationsQueue(),
			),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}
	if err := k.CalculateSubaccountPnLForMevMatches(
		ctx,
		validatorPnL,
		validatorMevMatches,
	); err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"Failed to calculate match PnL for validator: Error: %+v, MEV matches: %+v",
				err.Error(),
				validatorMevMatches,
			),
		)
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}

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
		k.Logger(ctx).Error("Failed to get consensus round")
		metrics.IncrCounter(
			metrics.ClobMevErrorCount,
			1,
		)
		return
	}

	// Add label for the block proposer.
	proposerConsAddress := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	proposer, found := stakingKeeper.GetValidatorByConsAddr(ctx, proposerConsAddress)
	if !found {
		k.Logger(ctx).Error(
			"Failed to get proposer by consensus address",
			"proposer",
			proposerConsAddress.String(),
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
			validatorPnL[clobPairId].MidPriceSubticks.ToUint64(),
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
		mevClobMidPrices := make([]types.ClobMidPrice, 0, len(clobPairs))
		for _, clobPair := range clobPairs {
			mevClobMidPrices = append(
				mevClobMidPrices,
				types.ClobMidPrice{
					ClobPair: clobPair,
					Subticks: clobMidPrices[types.ClobPairId(clobPair.Id)].ToUint64(),
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
// This function falls back to use the oracle price if any of the mid prices are missing.
func (k Keeper) GetClobMetadata(
	ctx sdk.Context,
) (
	clobMidPrices map[types.ClobPairId]types.Subticks,
	clobPairs map[types.ClobPairId]types.ClobPair,
) {
	clobMidPrices = make(map[types.ClobPairId]types.Subticks)
	clobPairs = make(map[types.ClobPairId]types.ClobPair)

	for _, clobPair := range k.GetAllClobPairs(ctx) {
		clobPairId := clobPair.GetClobPairId()
		var midPriceSubticks types.Subticks

		// Get the mid price if it exists, otherwise get the oracle price.
		if midPrice, exist := k.MemClob.GetMidPrice(ctx, clobPairId); exist {
			midPriceSubticks = midPrice
		} else {
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
			midPriceSubticks = types.Subticks(oraclePriceSubticksInt.Uint64())
		}

		// Set the CLOB mid price and CLOB pair.
		clobMidPrices[clobPairId] = midPriceSubticks
		clobPairs[types.ClobPairId(clobPairId)] = clobPair
	}

	return clobMidPrices, clobPairs
}

// InitializeCumulativePnLs initializes the cumulative PnLs for the block proposer and the
// current validator.
func (k Keeper) InitializeCumulativePnLs(
	ctx sdk.Context,
	perpetualKeeper process.ProcessPerpetualKeeper,
	clobMidPrices map[types.ClobPairId]types.Subticks,
	clobPairs map[types.ClobPairId]types.ClobPair,
) (
	blockProposerPnL map[types.ClobPairId]*CumulativePnL,
	validatorPnL map[types.ClobPairId]*CumulativePnL,
) {
	blockProposerPnL = make(map[types.ClobPairId]*CumulativePnL)
	validatorPnL = make(map[types.ClobPairId]*CumulativePnL)

	if len(clobMidPrices) != len(clobPairs) {
		panic(
			fmt.Sprintf(
				"InitializeCumulativePnLs: clob mid prices %+v and clob pairs %+v have different lengths",
				clobMidPrices,
				clobPairs,
			),
		)
	}

	for clobPairId, clobPair := range clobPairs {
		var midPriceSubticks types.Subticks

		// Panic if the mid price does not exist
		midPriceSubticks, exists := clobMidPrices[clobPairId]
		if !exists {
			panic(
				fmt.Sprintf(
					"InitializeCumulativePnLs: mid price does not exist for clob pair %+v",
					clobPair,
				),
			)
		}

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
				ClobPair:                    clobPair,
				MidPriceSubticks:            midPriceSubticks,
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
	clobPairs map[types.ClobPairId]types.ClobPair,
) (
	validatorMevMatches *types.ValidatorMevMatches,
	err error,
) {
	// Collect all the short-term orders placed for subsequent lookups.
	placedShortTermOrders := make(map[types.OrderId]types.Order)

	// Populate `mevMatches` and `mevLiquidationMatches` from the local validator's match operations.
	mevMatches := make([]types.MEVMatch, 0)
	mevLiquidationMatches := make([]types.MEVLiquidationMatch, 0)
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
						),

						MakerOrderSubaccountId: &makerOrder.OrderId.SubaccountId,
						MakerOrderSubticks:     makerOrder.Subticks,
						MakerOrderIsBuy:        makerOrder.IsBuy(),
						MakerFeePpm: k.feeTiersKeeper.GetPerpetualFeePpm(
							ctx,
							makerOrder.GetSubaccountId().Owner,
							false,
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
) (
	err error,
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
			if err := cumulativePnL.AddPnLForTradeWithFilledSubticks(
				p.subaccountId,
				p.isBuy,
				types.Subticks(matchWithOrders.MakerOrderSubticks),
				satypes.BaseQuantums(matchWithOrders.FillAmount),
				p.feePpm,
			); err != nil {
				return err
			}
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
			if err := cumulativePnL.AddPnLForTradeWithFilledSubticks(
				p.subaccountId,
				p.isBuy,
				types.Subticks(mevLiquidation.MakerOrderSubticks),
				satypes.BaseQuantums(mevLiquidation.FillAmount),
				p.feePpm,
			); err != nil {
				return err
			}
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

	return nil
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
						if err := cumulativePnL.AddPnLForTradeWithFilledQuoteQuantums(
							p.subaccountId,
							p.isBuy,
							absQuoteQuantums,
							satypes.BaseQuantums(fill.FillAmount),
							p.feePpm,
						); err != nil {
							return err
						}
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
		perpetualId := cumulativePnL.ClobPair.MustGetPerpetualId()
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
			bigNetSettlementPpm, _, err := perpetualKeeper.GetSettlementPpm(
				ctx,
				perpetualId,
				deltaQuantums,
				// Use the position's old funding index to calculate the funding payment.
				fundingIndex,
			)
			if err != nil {
				return err
			}

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
) (err error) {
	// Get the fill quote quantums using the filled subticks and filled quantums.
	filledQuoteQuantums, err := getFillQuoteQuantums(c.ClobPair, filledSubticks, filledQuantums)
	if err != nil {
		return err
	}
	return c.AddPnLForTradeWithFilledQuoteQuantums(
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
) (err error) {
	// Get the fill quote quantums using the mid price subticks and filled quantums.
	filledQuoteQuantumsUsingMidPrice, err := getFillQuoteQuantums(c.ClobPair, c.MidPriceSubticks, filledQuantums)
	if err != nil {
		return err
	}

	// Calculate PnL for the given subaccount.
	var pnl *big.Int
	if isBuy {
		pnl = new(big.Int).Sub(filledQuoteQuantumsUsingMidPrice, filledQuoteQuantums)
	} else {
		pnl = new(big.Int).Sub(filledQuoteQuantums, filledQuoteQuantumsUsingMidPrice)
	}

	// Calculate fees.
	bigFeeQuoteQuantums := lib.BigIntMulSignedPpm(filledQuoteQuantums, feePpm, true)
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
	return nil
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
