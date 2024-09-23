package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: [YBCP-89]
func (k Keeper) ClaimYieldForSubaccountFromIdAndSetNewState(
	ctx sdk.Context,
	subaccountId *types.SubaccountId,
) (
	err error,
) {

	if subaccountId == nil {
		return types.ErrSubaccountIdIsNil
	}

	subaccount := k.GetSubaccount(ctx, *subaccountId)
	if len(subaccount.AssetPositions) == 0 && len(subaccount.PerpetualPositions) == 0 {
		return types.ErrNoYieldToClaim
	}

	perpIdToPerp, assetYieldIndex, availableYield, err := k.fetchParamsToSettleSubaccount(ctx, subaccount)
	if err != nil {
		return err
	}

	settledSubaccount, totalYieldInQuantums, err := AddYieldToSubaccount(subaccount, perpIdToPerp, assetYieldIndex, availableYield)
	if err != nil {
		return err
	}

	err = k.DepositYieldToSubaccount(ctx, *settledSubaccount.Id, totalYieldInQuantums)
	if err != nil {
		return err
	}

	k.SetSubaccount(ctx, settledSubaccount)

	return nil
}

func AddYieldToSubaccount(
	subaccount types.Subaccount,
	perpIdToPerp map[uint32]perptypes.Perpetual,
	assetYieldIndex *big.Rat,
	availableYieldInQuantums *big.Int,
) (
	settledSubaccount types.Subaccount,
	totalNewYieldInQuantums *big.Int,
	err error,
) {
	assetYield, err := getYieldFromAssetPositions(subaccount, assetYieldIndex)
	if err != nil {
		return types.Subaccount{}, nil, err
	}

	totalNewPerpYield, newPerpetualPositions, err := getYieldFromPerpPositions(subaccount, perpIdToPerp)
	if err != nil {
		return types.Subaccount{}, nil, err
	}

	totalNewYieldInQuantums = new(big.Int).Add(assetYield, totalNewPerpYield)

	totalNewYieldInQuantums = HandleInsufficientYieldDueToNegativeTNC(totalNewYieldInQuantums, availableYieldInQuantums)

	assetYieldIndexString := assetYieldIndex.String()
	newSubaccount := types.Subaccount{
		Id:                 subaccount.Id,
		AssetPositions:     subaccount.AssetPositions,
		PerpetualPositions: newPerpetualPositions,
		MarginEnabled:      subaccount.MarginEnabled,
		AssetYieldIndex:    assetYieldIndexString,
	}

	if totalNewYieldInQuantums.Cmp(big.NewInt(0)) < 0 {
		totalNewYieldInQuantums = big.NewInt(0)
	}

	newTDaiPosition := new(big.Int).Add(subaccount.GetTDaiPosition(), totalNewYieldInQuantums)

	// TODO(CLOB-993): Remove this function and use `UpdateAssetPositions` instead.
	newSubaccount.SetTDaiAssetPosition(newTDaiPosition)
	return newSubaccount, totalNewYieldInQuantums, nil
}

func HandleInsufficientYieldDueToNegativeTNC(
	totalNewYield *big.Int,
	availableYield *big.Int,
) (
	yieldToTransfer *big.Int,
) {

	yieldToTransfer = new(big.Int).Set(totalNewYield)
	if availableYield.Cmp(totalNewYield) < 0 {
		yieldToTransfer.Set(availableYield)
	}

	return yieldToTransfer
}

// -------------------ASSET YIELD --------------------------

func getYieldFromAssetPositions(
	subaccount types.Subaccount,
	assetYieldIndex *big.Rat,
) (
	newAssetYield *big.Int,
	err error,
) {
	for _, assetPosition := range subaccount.AssetPositions {
		if assetPosition.AssetId != assettypes.AssetTDai.Id {
			return nil, assettypes.ErrNotImplementedMulticollateral
		}

		newAssetYield, err := calculateAssetYieldInQuoteQuantums(subaccount, assetYieldIndex, assetPosition)
		if err != nil {
			return nil, err
		} else {
			return newAssetYield, err
		}
	}
	return big.NewInt(0), nil
}

func calculateAssetYieldInQuoteQuantums(
	subaccount types.Subaccount,
	generalYieldIndex *big.Rat,
	assetPosition *types.AssetPosition,
) (
	newYield *big.Int,
	err error,
) {

	if assetPosition == nil {
		return nil, types.ErrPositionIsNil
	}

	if generalYieldIndex == nil {
		return nil, types.ErrGlobaYieldIndexNil
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) < 0 {
		return nil, types.ErrGlobalYieldIndexNegative
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) == 0 {
		return big.NewInt(0), nil
	}

	if subaccount.AssetYieldIndex == "" {
		return nil, types.ErrYieldIndexUninitialized
	}

	currentYieldIndex, success := new(big.Rat).SetString(subaccount.AssetYieldIndex)
	if !success {
		return nil, types.ErrRatConversion
	}

	if generalYieldIndex.Cmp(currentYieldIndex) < 0 {
		return nil, types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount
	}

	assetAmount := new(big.Rat).SetInt(assetPosition.GetBigQuantums())
	currYieldIndexdivisor := currentYieldIndex
	if currYieldIndexdivisor.Cmp(big.NewRat(0, 1)) == 0 {
		currYieldIndexdivisor = big.NewRat(1, 1)
	}

	yieldIndexQuotient := new(big.Rat).Quo(generalYieldIndex, currYieldIndexdivisor)
	newAssetAmount := new(big.Rat).Mul(assetAmount, yieldIndexQuotient)
	newYieldRat := assetAmount.Sub(newAssetAmount, assetAmount)

	newYield = lib.BigRatRound(newYieldRat, false)

	return newYield, nil
}

// -------------------PERP YIELD --------------------------

func getYieldFromPerpPositions(
	subaccount types.Subaccount,
	perpIdToPerp map[uint32]perptypes.Perpetual,
) (
	totalNewPerpYield *big.Int,
	newPerpetualPositions []*types.PerpetualPosition,
	err error,
) {
	totalNewPerpYield = big.NewInt(0)
	newPerpetualPositions = []*types.PerpetualPosition{}

	for _, perpetualPosition := range subaccount.PerpetualPositions {
		perpetual, found := perpIdToPerp[perpetualPosition.PerpetualId]
		if !found {
			return nil,
				nil,
				errorsmod.Wrap(
					perptypes.ErrPerpetualDoesNotExist, lib.UintToString(perpetualPosition.PerpetualId),
				)
		}

		perpYield, perpYieldIndex, err := calculateNewPerpYield(perpetual, perpetualPosition)
		if err != nil {
			return nil, nil, err
		}
		totalNewPerpYield = new(big.Int).Add(totalNewPerpYield, perpYield)

		newPerpetualPosition := types.PerpetualPosition{
			PerpetualId:  perpetualPosition.PerpetualId,
			Quantums:     perpetualPosition.Quantums,
			FundingIndex: perpetualPosition.FundingIndex,
			YieldIndex:   perpYieldIndex.String(),
		}
		newPerpetualPositions = append(newPerpetualPositions, &newPerpetualPosition)
	}
	return totalNewPerpYield, newPerpetualPositions, nil
}

func calculateNewPerpYield(
	perpetual perptypes.Perpetual,
	perpetualPosition *types.PerpetualPosition,
) (
	newPerpYield *big.Int,
	perpYieldIndex *big.Rat,
	err error,
) {
	perpYieldIndex, err = getCurrentYieldIndexForPerp(perpetual)
	if err != nil {
		return nil, nil, err
	}

	newPerpYield, err = calculatePerpetualYieldInQuoteQuantums(perpetualPosition, perpYieldIndex)
	if err != nil {
		return nil, nil, err
	}

	return newPerpYield, perpYieldIndex, nil
}

func getCurrentYieldIndexForPerp(
	perp perptypes.Perpetual,
) (
	yieldIndex *big.Rat,
	err error,
) {
	if perp.YieldIndex == "" {
		return nil, types.ErrYieldIndexUninitialized
	}

	generalYieldIndex, success := new(big.Rat).SetString(perp.YieldIndex)
	if !success {
		return nil, types.ErrRatConversion
	}
	return generalYieldIndex, nil
}

func calculatePerpetualYieldInQuoteQuantums(
	perpPosition *types.PerpetualPosition,
	generalYieldIndex *big.Rat,
) (
	newYield *big.Int,
	err error,
) {
	if perpPosition == nil {
		return nil, types.ErrPositionIsNil
	}

	if generalYieldIndex == nil {
		return nil, types.ErrGlobaYieldIndexNil
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) < 0 {
		return nil, types.ErrGlobalYieldIndexNegative
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) == 0 {
		return big.NewInt(0), nil
	}

	if perpPosition.YieldIndex == "" {
		return nil, types.ErrYieldIndexUninitialized
	}

	currentYieldIndex, success := new(big.Rat).SetString(perpPosition.YieldIndex)
	if !success {
		return nil, types.ErrRatConversion
	}

	if generalYieldIndex.Cmp(currentYieldIndex) < 0 {
		return nil, types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount
	}

	yieldIndexDifference := new(big.Rat).Sub(generalYieldIndex, currentYieldIndex)
	perpAmount := new(big.Rat).SetInt(perpPosition.GetBigQuantums())
	newYieldRat := new(big.Rat).Mul(perpAmount, yieldIndexDifference)
	newYield = lib.BigRatRound(newYieldRat, false)

	return newYield, nil
}

// -------------------YIELD ON BANK LEVEL --------------------------

func (k Keeper) DepositYieldToSubaccount(
	ctx sdk.Context,
	subaccountId types.SubaccountId,
	totalYieldInQuantums *big.Int,
) error {
	if totalYieldInQuantums == nil {
		return nil
	}

	if totalYieldInQuantums.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	if totalYieldInQuantums.Cmp(big.NewInt(0)) == -1 {
		return types.ErrTryingToDepositNegativeYield
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetTDai.Id,
		totalYieldInQuantums,
	)
	if err != nil {
		return err
	}

	collateralPoolAddr, err := k.GetCollateralPoolForSubaccount(ctx, subaccountId)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		ratelimittypes.TDaiPoolAccount,
		collateralPoolAddr,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	return nil
}
