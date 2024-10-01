package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assetstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/heap"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) FetchInformationForLiquidations(
	ctx sdk.Context,
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

	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)
	marketPricesMap = lib.UniqueSliceToMap(marketPrices, func(m pricestypes.MarketPrice) uint32 {
		return m.Id
	})

	return subaccounts, marketPricesMap, perpetualsMap, liquidityTiersMap
}

func (k Keeper) GetLiquidatableAndNegativeTncSubaccountIds(
	ctx sdk.Context,
) (
	liquidatableSubaccountIds *heap.LiquidationPriorityHeap,
	negativeTncSubaccountIds []satypes.SubaccountId,
	err error,
) {

	subaccounts, marketPrices, perpetuals, liquidityTiers := k.FetchInformationForLiquidations(ctx)

	negativeTncSubaccountIds = make([]satypes.SubaccountId, 0)
	liquidatableSubaccountIds = heap.NewLiquidationPriorityHeap()
	for _, subaccount := range subaccounts {

		if len(subaccount.PerpetualPositions) == 0 {
			continue
		}

		isLiquidatable, hasNegativeTnc, liquidationPriority, err := k.GetSubaccountCollateralizationInfo(ctx, subaccount, marketPrices, perpetuals, liquidityTiers)

		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "Error checking collateralization status")
		}

		if isLiquidatable {
			liquidatableSubaccountIds.AddSubaccount(*subaccount.Id, liquidationPriority)
		}
		if hasNegativeTnc {
			negativeTncSubaccountIds = append(negativeTncSubaccountIds, *subaccount.Id)
		}
	}

	return liquidatableSubaccountIds, negativeTncSubaccountIds, nil
}

func (k Keeper) GetSubaccountCollateralizationInfo(
	ctx sdk.Context,
	unsettledSubaccount satypes.Subaccount,
	marketPrices map[uint32]pricestypes.MarketPrice,
	perpetuals map[uint32]perptypes.Perpetual,
	liquidityTiers map[uint32]perptypes.LiquidityTier,
) (
	isLiquidatable bool,
	hasNegativeTnc bool,
	liquidationPriority *big.Float,
	err error,
) {
	bigTotalNetCollateral := big.NewInt(0)
	bigTotalMaintenanceMargin := big.NewInt(0)
	bigWeightedMaintenanceMargin := big.NewInt(0)

	settledSubaccount, _, _, err := k.subaccountsKeeper.GetSettledSubaccount(ctx, unsettledSubaccount)
	if err != nil {
		return false, false, nil, err
	}

	err = updateCollateralizationInfoGivenAssets(settledSubaccount, bigTotalNetCollateral)
	if err != nil {
		return false, false, nil, err
	}

	for _, perpetualPosition := range settledSubaccount.PerpetualPositions {
		perpetual, price, liquidityTier, err := getPerpetualLiquidityTierAndPrice(perpetualPosition.PerpetualId, perpetuals, marketPrices, liquidityTiers)
		if err != nil {
			return false, false, nil, err
		}
		updateCollateralizationInfoGivenPerp(perpetual, price, liquidityTier, perpetualPosition.GetBigQuantums(), bigTotalNetCollateral, bigWeightedMaintenanceMargin, bigTotalMaintenanceMargin)
	}

	return finalizeCollateralizationInfo(bigTotalNetCollateral, bigTotalMaintenanceMargin, bigWeightedMaintenanceMargin)
}

func getPerpetualLiquidityTierAndPrice(
	perpetualId uint32,
	perpetuals map[uint32]perptypes.Perpetual,
	marketPrices map[uint32]pricestypes.MarketPrice,
	liquidityTiers map[uint32]perptypes.LiquidityTier,
) (
	perpetual perptypes.Perpetual,
	price pricestypes.MarketPrice,
	liquidityTier perptypes.LiquidityTier,
	err error,
) {
	perpetual, ok := perpetuals[perpetualId]
	if !ok {
		return perptypes.Perpetual{}, pricestypes.MarketPrice{}, perptypes.LiquidityTier{}, errorsmod.Wrapf(
			perptypes.ErrPerpetualDoesNotExist,
			"Perpetual not found for perpetual id %d",
			perpetualId,
		)
	}

	price, ok = marketPrices[perpetual.Params.MarketId]
	if !ok {
		return perptypes.Perpetual{}, pricestypes.MarketPrice{}, perptypes.LiquidityTier{}, errorsmod.Wrapf(
			pricestypes.ErrMarketPriceDoesNotExist,
			"MarketPrice not found for perpetual %+v",
			perpetual,
		)
	}

	liquidityTier, ok = liquidityTiers[perpetual.Params.LiquidityTier]
	if !ok {
		return perptypes.Perpetual{}, pricestypes.MarketPrice{}, perptypes.LiquidityTier{}, errorsmod.Wrapf(
			perptypes.ErrLiquidityTierDoesNotExist,
			"LiquidityTier not found for perpetual %+v",
			perpetual,
		)
	}

	return perpetual, price, liquidityTier, nil
}

func updateCollateralizationInfoGivenAssets(
	settledSubaccount satypes.Subaccount,
	bigTotalNetCollateral *big.Int,
) error {

	// Note that we only expect TDai before multi-collateral support is added.
	for _, assetPosition := range settledSubaccount.AssetPositions {
		if assetPosition.AssetId != assetstypes.AssetTDai.Id {
			return errorsmod.Wrapf(
				assetstypes.ErrNotImplementedMulticollateral,
				"Asset %d is not supported",
				assetPosition.AssetId,
			)
		}
		bigTotalNetCollateral.Add(bigTotalNetCollateral, assetPosition.GetBigQuantums())
	}
	return nil
}

func updateCollateralizationInfoGivenPerp(
	perpetual perptypes.Perpetual,
	price pricestypes.MarketPrice,
	liquidityTier perptypes.LiquidityTier,
	bigPositionQuantums *big.Int,
	bigTotalNetCollateral *big.Int,
	bigWeightedMaintenanceMargin *big.Int,
	bigTotalMaintenanceMargin *big.Int,
) {
	updateNetCollateral(perpetual, price, bigPositionQuantums, bigTotalNetCollateral)
	updateWeightedMaintenanceMargin(perpetual, price, bigPositionQuantums, bigWeightedMaintenanceMargin)
	updateTotalMaintenanceMargin(perpetual, price, liquidityTier, bigPositionQuantums, bigTotalMaintenanceMargin)
}

func updateNetCollateral(
	perpetual perptypes.Perpetual,
	price pricestypes.MarketPrice,
	bigPositionQuantums *big.Int,
	bigTotalNetCollateral *big.Int,
) {
	bigPositionQuoteQuantums := perpkeeper.GetNetNotionalInQuoteQuantums(perpetual, price, bigPositionQuantums)
	bigTotalNetCollateral.Add(bigTotalNetCollateral, bigPositionQuoteQuantums)
}

func updateWeightedMaintenanceMargin(
	perpetual perptypes.Perpetual,
	price pricestypes.MarketPrice,
	bigPositionQuantums *big.Int,
	bigWeightedMaintenanceMargin *big.Int,
) {
	bigPositionQuoteQuantums := perpkeeper.GetNetNotionalInQuoteQuantums(perpetual, price, bigPositionQuantums)
	weightedPositionQuoteQuantums := new(big.Int).Mul(bigPositionQuoteQuantums.Abs(bigPositionQuoteQuantums), new(big.Int).SetUint64(uint64(perpetual.Params.DangerIndexPpm)))
	bigWeightedMaintenanceMargin.Add(bigWeightedMaintenanceMargin, weightedPositionQuoteQuantums)
}

func updateTotalMaintenanceMargin(
	perpetual perptypes.Perpetual,
	price pricestypes.MarketPrice,
	liquidityTier perptypes.LiquidityTier,
	bigPositionQuantums *big.Int,
	bigTotalMaintenanceMargin *big.Int,
) {
	_, bigMaintenanceMarginQuoteQuantums := perpkeeper.GetMarginRequirementsInQuoteQuantums(perpetual, price, liquidityTier, bigPositionQuantums)
	bigTotalMaintenanceMargin.Add(bigTotalMaintenanceMargin, bigMaintenanceMarginQuoteQuantums)
}

func finalizeCollateralizationInfo(
	bigTotalNetCollateral *big.Int,
	bigTotalMaintenanceMargin *big.Int,
	bigWeightedMaintenanceMargin *big.Int,
) (
	isLiquidatable bool,
	hasNegativeTnc bool,
	liquidationPriority *big.Float,
	err error,
) {
	isLiquidatable = CanLiquidateSubaccount(bigTotalNetCollateral, bigTotalMaintenanceMargin)
	hasNegativeTnc = bigTotalNetCollateral.Sign() == -1
	liquidationPriority = CalculateLiquidationPriority(
		bigTotalNetCollateral,
		bigTotalMaintenanceMargin,
		bigWeightedMaintenanceMargin,
	)

	return isLiquidatable, hasNegativeTnc, liquidationPriority, nil
}
