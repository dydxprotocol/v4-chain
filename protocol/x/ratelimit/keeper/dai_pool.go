package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	sDaiToTDaiDenomExponentDecimals = new(big.Int).Abs(
		big.NewInt(types.SDaiDenomExponent - assettypes.TDaiDenomExponent),
	)
	tenScaledBysDaiToTDaiDenomDecimals = new(big.Int).Exp(
		big.NewInt(10),
		sDaiToTDaiDenomExponentDecimals,
		nil,
	)
)

// MintTradingDAIToUserAccount transfers the input sDAI amount from the user's
// account to the pool account and mints the corresponding amount of trading
// DAI into the user's account.
func (k Keeper) MintTradingDAIToUserAccount(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	sDAIAmount *big.Int,
) error {

	tradingDAIAmount, err := k.GetTradingDAIFromSDAIAmount(ctx, sDAIAmount)
	if err != nil {
		return errorsmod.Wrap(err, "failed to convert sDAI to trading DAI")
	}

	sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(sDAIAmount)))
	tradingDAICoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewIntFromBigInt(tradingDAIAmount)))

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddr, types.SDaiPoolAccount, sDAICoins); err != nil {
		return errorsmod.Wrap(err, "failed to send sDAI to ratelimit module")
	}

	if err := k.bankKeeper.MintCoins(
		ctx, types.TDaiPoolAccount, tradingDAICoins,
	); err != nil {
		return errorsmod.Wrap(err, "failed to mint new trading DAI")
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.TDaiPoolAccount, userAddr, tradingDAICoins); err != nil {
		return errorsmod.Wrap(err, "failed to send trading DAI to recipient account")
	}

	return nil
}

// Withdraws the converted amount of TDai from the TDai pool and burns it.
// Sends the amount of sDai to the user account. We round up the converted
// amount of tdai on withdrawal, since we never want to leave "dangling" tdai.
func (k Keeper) WithdrawSDaiFromTDai(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	sDaiAmount *big.Int,
) error {

	tDaiDenomAmount, err := k.GetTradingDAIFromSDAIAmountAndRoundUp(ctx, sDaiAmount)
	if err != nil {
		return err
	}

	err = k.burnTDaiInUserAccount(ctx, userAddr, tDaiDenomAmount)
	if err != nil {
		return err
	}

	err = k.sendSDaiAmountToUserAccount(ctx, userAddr, sDaiAmount)

	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnTDaiInUserAccount(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	tDaiAmount *big.Int,
) error {
	tDaiCoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, sdkmath.NewIntFromBigInt(tDaiAmount)))

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddr, types.TDaiPoolAccount, tDaiCoins)
	if err != nil {
		return errorsmod.Wrap(err, "failed to send tDAI from user account to tDai pool account")
	}

	err = k.bankKeeper.BurnCoins(ctx, types.TDaiPoolAccount, tDaiCoins)

	if err != nil {
		return errorsmod.Wrap(err, "failed to burn tDai transferred to tDai pool account")
	}

	return nil
}

func (k Keeper) sendSDaiAmountToUserAccount(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	sDaiAmount *big.Int,
) error {
	sDAICoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewIntFromBigInt(sDaiAmount)))
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.SDaiPoolAccount, userAddr, sDAICoins)

	if err != nil {
		return errorsmod.Wrap(err, "failed to send sDAI to user account")
	}

	return nil
}

// Inspired by the deposit function of the Maker code at:
// https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea
func (k Keeper) GetTradingDAIFromSDAIAmount(ctx sdk.Context, sDaiAmount *big.Int) (*big.Int, error) {
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return nil, errors.New("sDAI price not found")
	}

	if sDAIPrice.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("sDAI price is zero")
	}

	tDaiAmount := k.calculateTDaiAmount(sDaiAmount, sDAIPrice)

	scaledTDaiAmount := k.scaleTDaiAmountByDenomExponent(tDaiAmount)

	return scaledTDaiAmount, nil
}

// Inspired by the withdraw function of the Maker code at:
// https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea
// NOTE: sDaiAmount should be provided as gsDai (i.e., scaled by the gsdai denom exponent)
func (k Keeper) GetTradingDAIFromSDAIAmountAndRoundUp(
	ctx sdk.Context,
	sDaiAmount *big.Int,
) (*big.Int, error) {
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return nil, errorsmod.Wrap(
			types.ErrFailedSDaiToTDaiConversion,
			"sDai price not found",
		)
	}

	if sDAIPrice.Cmp(big.NewInt(0)) == 0 {
		return nil, errorsmod.Wrap(
			types.ErrFailedSDaiToTDaiConversion,
			"sDAI price is zero",
		)
	}

	tDaiAmount, err := k.calculateTDaiAmountAndRoundUp(sDaiAmount, sDAIPrice)
	if err != nil {
		return nil, err
	}

	scaledTDaiAmount, err := k.scaleTDaiAmountByDenomExponentAndRoundUp(tDaiAmount)
	if err != nil {
		return nil, err
	}

	return scaledTDaiAmount, nil
}

func (k Keeper) calculateTDaiAmount(sDaiAmount *big.Int, sDAIPrice *big.Int) *big.Int {
	scaledSDaiAmount := new(big.Int).Mul(sDaiAmount, sDAIPrice)
	tDaiAmount := divideAmountBySDaiDecimals(scaledSDaiAmount)
	return tDaiAmount

}

func (k Keeper) calculateTDaiAmountAndRoundUp(sDaiAmount *big.Int, sDAIPrice *big.Int) (*big.Int, error) {
	scaledSDaiAmount := new(big.Int).Mul(sDaiAmount, sDAIPrice)
	tenScaledBySDaiDecimals := getTenScaledBySDaiDecimals()
	tDaiAmount, err := divideAndRoundUp(scaledSDaiAmount, tenScaledBySDaiDecimals)
	if err != nil {
		return nil, err
	}
	return tDaiAmount, nil
}

func (k Keeper) scaleTDaiAmountByDenomExponent(tDaiAmount *big.Int) *big.Int {
	scaledTDaiAmount := new(big.Int).Div(tDaiAmount, tenScaledBysDaiToTDaiDenomDecimals)
	return scaledTDaiAmount
}

func (k Keeper) scaleTDaiAmountByDenomExponentAndRoundUp(tDaiAmount *big.Int) (*big.Int, error) {
	scaledTDaiAmount, err := divideAndRoundUp(tDaiAmount, tenScaledBysDaiToTDaiDenomDecimals)
	if err != nil {
		return big.NewInt(0), err
	}
	return scaledTDaiAmount, nil
}
