package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintNewTDaiAndSetNewYieldIndex(ctx sdk.Context) error {

	tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted, err := k.MintNewTDaiYield(ctx)
	if err != nil {
		return err
	}

	err = k.SetNewYieldIndex(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)
	if err != nil {
		return err
	}

	k.perpetualsKeeper.UpdateYieldIndexToNewMint(ctx, tradingDaiSupplyBeforeNewEpoch, tradingDaiMinted)

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

	tDAISupplyCoins := k.bankKeeper.GetSupply(ctx, types.TradingDAIDenom)
	tDAISupply := tDAISupplyCoins.Amount.BigInt()

	tradingDAIAfterYield, err := k.GetTradingDAIFromSDAIAmount(ctx, sDAISupply)
	if err != nil {
		return nil, nil, err
	}

	if tradingDAIAfterYield.Cmp(tDAISupply) <= 0 {
		return nil, nil, errorsmod.Wrap(
			types.ErrInvalidSDAIConversionRate,
			"Trading DAI supply is less than or equal to the sDAI supply",
		)
	}

	tradingDaiToMint := tradingDAIAfterYield.Sub(tradingDAIAfterYield, tDAISupply)
	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewIntFromBigInt(tradingDaiToMint)))

	err = k.bankKeeper.MintCoins(ctx, types.PoolAccount, tradingDAICoins)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	return tDAISupply, tradingDaiToMint, nil
}
