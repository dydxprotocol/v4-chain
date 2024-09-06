package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ClaimYieldForSubaccountFromId(
	ctx sdk.Context,
	subaccountId *types.SubaccountId,
) (
	err error,
) {

	subaccount := k.GetSubaccount(ctx, *subaccountId)
	if len(subaccount.AssetPositions) == 0 && len(subaccount.PerpetualPositions) == 0 {
		return errors.New("there is no yield to claim for subaccount")
	}

	settledSubaccount, yieldEarned, err := k.settleSubaccountYield(ctx, subaccount)
	if err != nil {
		return err
	}

	k.SetSubaccount(ctx, settledSubaccount)

	k.DepositYieldToSubaccount(ctx, *subaccountId, yieldEarned)

	return nil
}

func (k Keeper) settleSubaccountYield(
	ctx sdk.Context,
	subaccount types.Subaccount,
) (
	settledSubaccount types.Subaccount,
	totalYield *big.Int,
	err error,
) {

	perpIdToPerp, assetYieldIndex, err := k.fetchParamsToSettleSubaccount(ctx, subaccount)
	if err != nil {
		return types.Subaccount{}, nil, err
	}

	isYieldAlreadyClaimed, err := IsYieldAlreadyClaimed(assetYieldIndex, subaccount.AssetYieldIndex)
	if err != nil {
		return types.Subaccount{}, nil, err
	}
	if isYieldAlreadyClaimed {
		return subaccount, big.NewInt(0), nil
	}

	settledSubaccount, totalYield, err = AddYieldToSubaccount(subaccount, perpIdToPerp, assetYieldIndex)
	if err != nil {
		return types.Subaccount{}, nil, err
	}

	return settledSubaccount, totalYield, nil
}

func IsYieldAlreadyClaimed(assetYieldIndex *big.Rat, subaccountAssetYieldIndex string) (bool, error) {

	currentYieldIndex, success := new(big.Rat).SetString(subaccountAssetYieldIndex)
	if !success {
		return false, errors.New("could not convert the subaccount yield index to big.Rat")
	}

	if assetYieldIndex.Cmp(currentYieldIndex) == 0 {
		return true, nil
	}

	if assetYieldIndex.Cmp(currentYieldIndex) == -1 {
		return false, errors.New("subaccount's asset yield index greater than overall asset yield index")
	}

	return false, nil
}

func AddYieldToSubaccount(
	subaccount types.Subaccount,
	perpIdToPerp map[uint32]perptypes.Perpetual,
	assetYieldIndex *big.Rat,
) (
	settledSubaccount types.Subaccount,
	totalNewYield *big.Int,
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

	totalNewYield = new(big.Int).Add(assetYield, totalNewPerpYield)

	newSubaccount := types.Subaccount{
		Id:                 subaccount.Id,
		AssetPositions:     subaccount.AssetPositions,
		PerpetualPositions: newPerpetualPositions,
		MarginEnabled:      subaccount.MarginEnabled,
		AssetYieldIndex:    assetYieldIndex.String(),
	}

	// TODO [YBCP-21]: Handle negative yield more gracefully
	if totalNewYield.Cmp(big.NewInt(0)) < 0 {
		panic("Total yield is less than 0. This should not be the case")
	}

	newTDaiPosition := new(big.Int).Add(subaccount.GetTDaiPosition(), totalNewYield)

	// TODO(CLOB-993): Remove this function and use `UpdateAssetPositions` instead.
	newSubaccount.SetTDaiAssetPosition(newTDaiPosition)
	return newSubaccount, totalNewYield, nil
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
			continue
		}

		newAssetYield, err := calculateAssetYieldInQuoteQuantums(subaccount, assetYieldIndex, assetPosition)
		if err != nil {
			return nil, err
		}

		return newAssetYield, err
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

	if generalYieldIndex == nil {
		return nil, errors.New("general yield index is nil")
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) < 0 {
		return nil, errors.New("general yield index is negative")
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) == 0 {
		return big.NewInt(0), nil
	}

	if subaccount.AssetYieldIndex == "" {
		return nil, errors.New("asset yield for subaccount is badly initialised 0/1")
	}

	currentYieldIndex, success := new(big.Rat).SetString(subaccount.AssetYieldIndex)
	if !success {
		return nil, errors.New("could not convert the subaccount yield index to big.Rat")
	}

	if generalYieldIndex.Cmp(currentYieldIndex) < 0 {
		return nil, errors.New("general yield index is less than the current yield index")
	}

	yieldIndexDifference := new(big.Rat).Sub(generalYieldIndex, currentYieldIndex)
	assetAmount := new(big.Rat).SetInt(assetPosition.GetBigQuantums())
	newYieldRat := new(big.Rat).Mul(assetAmount, yieldIndexDifference)
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

func calculatePerpetualYieldInQuoteQuantums(
	perpPosition *types.PerpetualPosition,
	generalYieldIndex *big.Rat,
) (
	newYield *big.Int,
	err error,
) {
	if perpPosition == nil {
		return nil, errors.New("could not calculate perpetual yield: perp position is nil")
	}

	if generalYieldIndex == nil {
		return nil, errors.New("could not calculate perpetual yield: perp yield index is nil")
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) < 0 {
		return nil, errors.New("general yield index is negative")
	}

	if generalYieldIndex.Cmp(big.NewRat(0, 1)) == 0 {
		return big.NewInt(0), nil
	}

	if perpPosition.YieldIndex == "" {
		return nil, errors.New("perp yield index for perp is empty string")
	}

	currentYieldIndex, success := new(big.Rat).SetString(perpPosition.YieldIndex)
	if !success {
		return nil, errors.New("could not convert yield index of perp position to big.Rat")
	}

	if generalYieldIndex.Cmp(currentYieldIndex) < 0 {
		return nil, errors.New("general yield index is less than the current yield index")
	}

	yieldIndexDifference := new(big.Rat).Sub(generalYieldIndex, currentYieldIndex)
	perpAmount := new(big.Rat).SetInt(perpPosition.GetBigQuantums())
	newYieldRat := new(big.Rat).Mul(perpAmount, yieldIndexDifference)
	newYield = lib.BigRatRound(newYieldRat, false)

	return newYield, nil
}

func getCurrentYieldIndexForPerp(
	perp perptypes.Perpetual,
) (
	yieldIndex *big.Rat,
	err error,
) {
	if perp.YieldIndex == "" {
		return nil, errors.New("perp yield index for perp is not initialised")
	}

	generalYieldIndex, success := new(big.Rat).SetString(perp.YieldIndex)
	if !success {
		return nil, errors.New("could not convert yield index of perp to big.Rat")
	}
	return generalYieldIndex, nil
}

// -------------------YIELD ON BANK LEVEL --------------------------

func (k Keeper) DepositYieldToSubaccount(
	ctx sdk.Context,
	subaccountId types.SubaccountId,
	amountToTransfer *big.Int,
) error {
	if amountToTransfer == nil {
		return nil
	}

	if amountToTransfer.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		assettypes.AssetTDai.Id,
		amountToTransfer,
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
