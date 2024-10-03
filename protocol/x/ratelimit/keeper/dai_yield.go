package keeper

import (
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessNewSDaiConversionRateUpdate(ctx sdk.Context, sDaiConversionRate *big.Int, blockHeight *big.Int) error {
	fmt.Println("IN PROCESS NEW SDAI CONVERSION RATE UPDATE")

	if sDaiConversionRate == nil || blockHeight == nil {
		return errors.New("sDaiConversionRate or blockHeight cannot be nil")
	}

	currBlockHeight, found := k.GetSDAILastBlockUpdated(ctx)

	if found && blockHeight.Cmp(currBlockHeight) == -1 {
		return errors.New("new block height is less than the current block height")
	}

	currConversionRate, found := k.GetSDAIPrice(ctx)

	if found {
		if sDaiConversionRate.Cmp(currConversionRate) == -1 {
			return errors.New("new sDAI conversion is less than the current sDAI conversion rate")
		}

		if sDaiConversionRate.Cmp(currConversionRate) == 0 {
			return nil
		}
	}

	fmt.Println("ALL TESTS PASSED")

	k.SetSDAIPrice(ctx, sDaiConversionRate)
	k.SetSDAILastBlockUpdated(ctx, blockHeight)

	fmt.Println("BEFORE UPDATE MINT STATE ON SDAI CONVERSION RATE UPDATE")

	return k.UpdateMintStateOnSDaiConversionRateUpdate(ctx)
}

func (k Keeper) UpdateMintStateOnSDaiConversionRateUpdate(ctx sdk.Context) error {
	tDaiSupplyDenomAmountBeforeNewEpoch, tDaiDenomAmountMinted, err := k.MintNewTDaiYield(ctx)
	if err != nil {
		return err
	}

	fmt.Println("AFTER MINT NEW TDAI YIELD")

	err = k.ClaimInsuranceFundYields(ctx, tDaiSupplyDenomAmountBeforeNewEpoch, tDaiDenomAmountMinted)
	if err != nil {
		return err
	}

	fmt.Println("AFTER CLAIM INSURANCE FUND YIELDS")

	err = k.SetNewYieldIndex(ctx, tDaiSupplyDenomAmountBeforeNewEpoch, tDaiDenomAmountMinted)
	if err != nil {
		return err
	}

	fmt.Println("AFTER SET NEW YIELD INDEX")

	err = k.perpetualsKeeper.UpdateYieldIndexToNewMint(ctx, tDaiSupplyDenomAmountBeforeNewEpoch, tDaiDenomAmountMinted)
	if err != nil {
		return err
	}

	fmt.Println("AFTER UPDATE YIELD INDEX TO NEW MINT")

	// Emit indexer event
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return errors.New("could not find sDAI price when emitting indexer event for new yield index")
	}

	fmt.Println("BEFORE GETTING ASSET YIELD INDEX. SDAI PRICE IS ", sDAIPrice)

	assetYieldIndex, found := k.GetAssetYieldIndex(ctx)
	if !found {
		return errors.New("could not find asset yield index when emitting indexer event for new yield index")
	}

	fmt.Println("ASSET YIELD INDEX IS ", assetYieldIndex)

	indexerevents.NewUpdateYieldParamsEventV1(
		sDAIPrice.String(),
		assetYieldIndex.String(),
	)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeYieldParams,
		indexerevents.UpdateYieldParamsEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewUpdateYieldParamsEventV1(
				sDAIPrice.String(),
				assetYieldIndex.String(),
			),
		),
	)

	return nil
}

func (k Keeper) ClaimInsuranceFundYields(ctx sdk.Context, tDaiSupplyDenomAmountBeforeNewEpoch *big.Int, tDaiDenomAmountMinted *big.Int) error {
	perps := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	insuranceFundsSeen := make(map[string]bool)

	for _, perpetual := range perps {
		insuranceFund, err := k.perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, perpetual.Params.Id)
		if err != nil {
			return err
		}

		if _, ok := insuranceFundsSeen[insuranceFund.String()]; ok {
			continue
		}

		insuranceFundsSeen[insuranceFund.String()] = true

		insuranceFundDenomBalance := k.bankKeeper.GetBalance(ctx, insuranceFund, types.TDaiDenom)
		if insuranceFundDenomBalance.IsZero() {
			continue
		}
		insuranceFundDenomBalanceBigInt := insuranceFundDenomBalance.Amount.BigInt()

		insuranceFundYieldToClaim := new(big.Int).Mul(insuranceFundDenomBalanceBigInt, tDaiDenomAmountMinted)
		insuranceFundYieldToClaim.Div(insuranceFundYieldToClaim, tDaiSupplyDenomAmountBeforeNewEpoch)

		// Ensure the insurance fund yield to mint is non-negative
		if insuranceFundYieldToClaim.Sign() <= 0 {
			continue
		}

		coinsToTransfer := sdk.NewCoin(assettypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(insuranceFundYieldToClaim))

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, insuranceFund, []sdk.Coin{coinsToTransfer}); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) SetNewYieldIndex(
	ctx sdk.Context,
	totalTDaiPreMint *big.Int,
	totalTDaiMinted *big.Int,
) error {
	if totalTDaiMinted.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	if totalTDaiPreMint.Cmp(big.NewInt(0)) == 0 {
		return errors.New("total t-dai minted is non-zero, while total t-dai before mint is 0")
	}

	ratio := new(big.Rat).SetFrac(totalTDaiMinted, totalTDaiPreMint)
	additionalFactor := ratio.Add(big.NewRat(1, 1), ratio)

	assetYieldIndex, found := k.GetAssetYieldIndex(ctx)

	if !found || assetYieldIndex.Cmp(big.NewRat(0, 1)) == 0 {
		assetYieldIndex = additionalFactor
	} else {
		assetYieldIndex = assetYieldIndex.Mul(assetYieldIndex, additionalFactor)
	}

	k.SetAssetYieldIndex(ctx, assetYieldIndex)
	return nil
}

func (k Keeper) MintNewTDaiYield(ctx sdk.Context) (*big.Int, *big.Int, error) {
	sDaiSupplyCoins := k.bankKeeper.GetSupply(ctx, types.SDaiDenom)
	sDaiSupplyDenomAmount := sDaiSupplyCoins.Amount.BigInt()

	if sDaiSupplyDenomAmount.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), big.NewInt(0), nil
	}

	tDaiSupplyCoins := k.bankKeeper.GetSupply(ctx, types.TDaiDenom)
	tDaiSupplyDenomAmount := tDaiSupplyCoins.Amount.BigInt()

	tDAIAfterYield, err := k.GetTradingDAIFromSDAIAmount(ctx, sDaiSupplyDenomAmount)
	if err != nil {
		return nil, nil, err
	}

	if tDAIAfterYield.Cmp(tDaiSupplyDenomAmount) <= 0 {
		return nil, nil, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"Trading DAI supply is less than or equal to the sDAI supply",
		)
	}

	tDaiDenomAmountToMint := tDAIAfterYield.Sub(tDAIAfterYield, tDaiSupplyDenomAmount)
	tDaiToMintCoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewIntFromBigInt(tDaiDenomAmountToMint)))

	err = k.bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, tDaiToMintCoins)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	return tDaiSupplyDenomAmount, tDaiDenomAmountToMint, nil
}
