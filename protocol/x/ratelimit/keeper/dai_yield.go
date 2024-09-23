package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessNewTDaiConversionRateUpdate(ctx sdk.Context) error {

	tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted, err := k.MintNewTDaiYield(ctx)
	if err != nil {
		return err
	}

	err = k.ClaimInsuranceFundYields(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)
	if err != nil {
		return err
	}

	err = k.SetNewYieldIndex(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)
	if err != nil {
		return err
	}

	err = k.perpetualsKeeper.UpdateYieldIndexToNewMint(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)
	if err != nil {
		return err
	}

	// Emit indexer event
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return errors.New("could not find sDAI price when emitting indexer event for new yield index")
	}

	assetYieldIndex, found := k.GetAssetYieldIndex(ctx)
	if !found {
		return errors.New("could not find asset yield index when emitting indexer event for new yield index")
	}

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

func (k Keeper) ClaimInsuranceFundYields(ctx sdk.Context, tradingDaiSupplyBeforeNewEpoch *big.Int, tradingDaiMinted *big.Int) error {

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

		insuranceFundBalance := k.bankKeeper.GetBalance(ctx, insuranceFund, types.TDaiDenom)
		if insuranceFundBalance.IsZero() {
			continue
		}

		insuranceFundBalanceBigInt, _, err := k.assetsKeeper.ConvertCoinToAsset(ctx, assettypes.AssetTDai.Id, insuranceFundBalance)
		if err != nil {
			return err
		}

		insuranceFundYieldToMint := new(big.Int).Mul(insuranceFundBalanceBigInt, tradingDaiMinted)
		insuranceFundYieldToMint.Div(insuranceFundYieldToMint, tradingDaiSupplyBeforeNewEpoch)

		// Ensure the insurance fund yield to mint is non-negative
		if insuranceFundYieldToMint.Sign() <= 0 {
			continue
		}

		_, coinsToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(ctx, assettypes.AssetTDai.Id, insuranceFundYieldToMint)
		if err != nil {
			return err
		}

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
	assetYieldIndex, found := k.GetAssetYieldIndex(ctx)
	if !found {
		return errors.New("could not retrieve asset yield index")
	}

	if totalTDaiMinted.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	if totalTDaiPreMint.Cmp(big.NewInt(0)) == 0 {
		return errors.New("total t-dai minted is non-zero, while total t-dai before mint is 0")
	}

	ratio := new(big.Rat).SetFrac(totalTDaiMinted, totalTDaiPreMint)
	assetYieldIndex = assetYieldIndex.Add(assetYieldIndex, ratio)

	k.SetAssetYieldIndex(ctx, assetYieldIndex)
	return nil
}

func (k Keeper) MintNewTDaiYield(ctx sdk.Context) (*big.Int, *big.Int, error) {

	sDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.SDaiDenom)
	sDAISupply := sDAISupplyCoins.Amount.BigInt()

	if sDAISupply.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), big.NewInt(0), nil
	}

	tDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.TDaiDenom)
	tDAISupply := tDAISupplyCoins.Amount.BigInt()

	tDAIAfterYield, err := k.GetTradingDAIFromSDAIAmount(ctx, sDAISupply)
	if err != nil {
		return nil, nil, err
	}

	if tDAIAfterYield.Cmp(tDAISupply) <= 0 {
		return nil, nil, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"Trading DAI supply is less than or equal to the sDAI supply",
		)
	}

	tradingDaiToMint := tDAIAfterYield.Sub(tDAIAfterYield, tDAISupply)
	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewIntFromBigInt(tradingDaiToMint)))

	err = k.bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, tradingDAICoins)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	return tDAISupply, tradingDaiToMint, nil
}
