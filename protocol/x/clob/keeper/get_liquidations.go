package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assetstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	subaccountskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	abcicomet "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) FetchInformationForLiquidations(
	ctx sdk.Context,
	extendedCommitInfo *abcicomet.ExtendedCommitInfo,
) (
	subaccounts []satypes.Subaccount,
	marketPricesMap map[uint32]pricestypes.MarketPrice,
	perpetualsMap map[uint32]perptypes.Perpetual,
	liquidityTiersMap map[uint32]perptypes.LiquidityTier,
) {

	subaccounts = k.subaccountsKeeper.GetAllSubaccount(ctx)

	perpetuals := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	perpetualsMap = lib.UniqueSliceToMap(perpetuals, func(p perptypes.Perpetual) uint32 {
		return p.Params.Id
	})

	liquidityTiers := k.perpetualsKeeper.GetAllLiquidityTiers(ctx)
	liquidityTiersMap = lib.UniqueSliceToMap(liquidityTiers, func(l perptypes.LiquidityTier) uint32 {
		return l.Id
	})

	marketPrices := k.GetNextBlocksPricesFromExtendedCommitInfo(ctx, extendedCommitInfo)
	marketPricesMap = lib.UniqueSliceToMap(marketPrices, func(m pricestypes.MarketPrice) uint32 {
		return m.Id
	})

	return subaccounts, marketPricesMap, perpetualsMap, liquidityTiersMap
}

func (k Keeper) GetLiquidatableAndTNCSubaccountIds(
	ctx sdk.Context,
	extendedCommitInfo *abcicomet.ExtendedCommitInfo,
) (
	liquidatableSubaccountIds []satypes.SubaccountId,
	negativeTncSubaccountIds []satypes.SubaccountId,
	err error,
) {

	subaccounts, marketPrices, perpetuals, liquidityTiers := k.FetchInformationForLiquidations(ctx, extendedCommitInfo)

	liquidatableSubaccountIds = make([]satypes.SubaccountId, 0)
	negativeTncSubaccountIds = make([]satypes.SubaccountId, 0)
	for _, subaccount := range subaccounts {
		// Skip subaccounts with no open positions.
		if len(subaccount.PerpetualPositions) == 0 {
			continue
		}

		// Check if the subaccount is liquidatable.
		isLiquidatable, hasNegativeTnc, err := k.CheckSubaccountCollateralization(
			subaccount,
			marketPrices,
			perpetuals,
			liquidityTiers,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "Error checking collateralization status")
		}

		if isLiquidatable {
			liquidatableSubaccountIds = append(liquidatableSubaccountIds, *subaccount.Id)
		}
		if hasNegativeTnc {
			negativeTncSubaccountIds = append(negativeTncSubaccountIds, *subaccount.Id)
		}
	}

	return liquidatableSubaccountIds, negativeTncSubaccountIds, nil
}

func (k Keeper) CheckSubaccountCollateralization(
	unsettledSubaccount satypes.Subaccount,
	marketPrices map[uint32]pricestypes.MarketPrice,
	perpetuals map[uint32]perptypes.Perpetual,
	liquidityTiers map[uint32]perptypes.LiquidityTier,
) (
	isLiquidatable bool,
	hasNegativeTnc bool,
	err error,
) {

	// Funding payments are lazily settled, so get the settled subaccount
	// to ensure that the funding payments are included in the net collateral calculation.
	settledSubaccount, _, err := subaccountskeeper.GetSettledSubaccountWithPerpetuals(
		unsettledSubaccount,
		perpetuals,
	)
	if err != nil {
		return false, false, err
	}

	bigTotalNetCollateral := big.NewInt(0)
	bigTotalMaintenanceMargin := big.NewInt(0)

	// Calculate the net collateral and maintenance margin for each of the asset positions.
	// Note that we only expect USDC before multi-collateral support is added.
	for _, assetPosition := range settledSubaccount.AssetPositions {
		if assetPosition.AssetId != assetstypes.AssetUsdc.Id {
			return false, false, errorsmod.Wrapf(
				assetstypes.ErrNotImplementedMulticollateral,
				"Asset %d is not supported",
				assetPosition.AssetId,
			)
		}
		// Net collateral for USDC is the quantums of the position.
		// Margin requirements for USDC are zero.
		bigTotalNetCollateral.Add(bigTotalNetCollateral, assetPosition.GetBigQuantums())
	}

	// Calculate the net collateral and maintenance margin for each of the perpetual positions.
	for _, perpetualPosition := range settledSubaccount.PerpetualPositions {
		perpetual, ok := perpetuals[perpetualPosition.PerpetualId]
		if !ok {
			return false, false, errorsmod.Wrapf(
				perptypes.ErrPerpetualDoesNotExist,
				"Perpetual not found for perpetual id %d",
				perpetualPosition.PerpetualId,
			)
		}

		marketPrice, ok := marketPrices[perpetual.Params.MarketId]
		if !ok {
			return false, false, errorsmod.Wrapf(
				pricestypes.ErrMarketPriceDoesNotExist,
				"MarketPrice not found for perpetual %+v",
				perpetual,
			)
		}

		bigQuantums := perpetualPosition.GetBigQuantums()

		// Get the net collateral for the position.
		bigNetCollateralQuoteQuantums := perpkeeper.GetNetNotionalInQuoteQuantums(perpetual, marketPrice, bigQuantums)
		bigTotalNetCollateral.Add(bigTotalNetCollateral, bigNetCollateralQuoteQuantums)

		liquidityTier, ok := liquidityTiers[perpetual.Params.LiquidityTier]
		if !ok {
			return false, false, errorsmod.Wrapf(
				perptypes.ErrLiquidityTierDoesNotExist,
				"LiquidityTier not found for perpetual %+v",
				perpetual,
			)
		}

		// Get the maintenance margin requirement for the position.
		_, bigMaintenanceMarginQuoteQuantums := perpkeeper.GetMarginRequirementsInQuoteQuantums(
			perpetual,
			marketPrice,
			liquidityTier,
			bigQuantums,
		)
		bigTotalMaintenanceMargin.Add(bigTotalMaintenanceMargin, bigMaintenanceMarginQuoteQuantums)
	}

	return CanLiquidateSubaccount(bigTotalNetCollateral, bigTotalMaintenanceMargin),
		bigTotalNetCollateral.Sign() == -1,
		nil
}

func (k Keeper) GetNextBlocksPricesFromExtendedCommitInfo(ctx sdk.Context, extendedCommitInfo *abcicomet.ExtendedCommitInfo) (marketPrices []pricestypes.MarketPrice) {

	branchedCtx, ctxErr := ctx.CacheContext()
	marketPrices = k.pricesKeeper.GetAllMarketPrices(ctx)

	// from cometbft so is either nil or is valid and > 2/3
	if (extendedCommitInfo != &abcicomet.ExtendedCommitInfo{}) && ctxErr == nil {
		veCodec := vecodec.NewDefaultVoteExtensionCodec()
		votes, err := veaggregator.FetchVotesFromExtCommitInfo(*extendedCommitInfo, veCodec)
		if err == nil {
			prices, err := k.PriceApplier.VoteAggregator().AggregateDaemonVEIntoFinalPrices(ctx, votes)
			if err == nil {
				err = k.PriceApplier.WritePricesToStoreAndMaybeCache(branchedCtx, prices, 0, false)
				if err == nil {
					marketPrices = k.pricesKeeper.GetAllMarketPrices(branchedCtx)
				}
			}
		}
	}

	return marketPrices
}
