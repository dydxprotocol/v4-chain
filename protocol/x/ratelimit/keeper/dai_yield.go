package keeper

import (
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessNewTDaiConversionRateUpdate(ctx sdk.Context) error {

	tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted, err := k.MintNewTDaiYield(ctx)
	if err != nil {
		return err
	}

	err = k.SetNewYieldIndex(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)
	if err != nil {
		return err
	}

	k.perpetualsKeeper.UpdateYieldIndexToNewMint(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)

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

	fmt.Println("tDAI SUPPLY BEFORE YIELD ", tDAISupply)
	fmt.Println("tDAI SUPPLY after yield ", tDAIAfterYield)

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
