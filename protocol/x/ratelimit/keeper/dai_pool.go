package keeper

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
// Sends the amount of sDai to the user account.
func (k Keeper) WithdrawSDaiFromTDai(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	sDaiAmount *big.Int,
) error {

	tDAIAmount, err := k.GetTradingDAIFromSDAIAmountAndRoundUp(ctx, sDaiAmount)
	if err != nil {
		return err
	}

	err = k.BurnTDaiInUserAccount(ctx, userAddr, tDAIAmount)
	if err != nil {
		return err
	}

	err = k.SendSDaiAmountToUserAccount(ctx, userAddr, sDaiAmount)

	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) BurnTDaiInUserAccount(
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

func (k Keeper) SendSDaiAmountToUserAccount(
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

/*
Converts sDAI to corresponding amount of tDAI, implementing the following maker code
https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea.
Note that shares and tDaiAmount are equivalent.

	function deposit(uint256 assets, address receiver) external returns (uint256 shares) {
		uint256 chi = (block.timestamp > pot.rho()) ? pot.drip() : pot.chi();
		shares = assets * RAY / chi;
		_mint(assets, shares, receiver);
	}
*/
func (k Keeper) GetTradingDAIFromSDAIAmount(ctx sdk.Context, sDaiAmount *big.Int) (*big.Int, error) {
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return nil, errors.New("sDAI price not found")
	}

	if sDAIPrice.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("sDAI price is zero")
	}

	return k.calculateTDaiAmount(ctx, sDaiAmount, sDAIPrice), nil
}

/*
Inspired by the following maker code.
https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea
Note that shares and tDaiAmount are equivalent.

	function withdraw(uint256 assets, address receiver, address owner) external returns (uint256 shares) {
		uint256 chi = (block.timestamp > pot.rho()) ? pot.drip() : pot.chi();
		shares = _divup(assets * RAY, chi);
		_burn(assets, shares, receiver, owner);
	}
*/
func (k Keeper) GetTradingDAIFromSDAIAmountAndRoundUp(ctx sdk.Context, sDaiAmount *big.Int) (*big.Int, error) {
	// Get the current sDAI price
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

	tDaiAmount, err := k.calculateTDaiAmountAndRoundUp(ctx, sDaiAmount, sDAIPrice)
	if err != nil {
		return nil, err
	}

	return tDaiAmount, nil
}

func (k Keeper) calculateTDaiAmount(ctx sdk.Context, sDaiAmount *big.Int, sDAIPrice *big.Int) *big.Int {
	scaledSDaiAmount := scaleAmountBySDaiDecimals(ctx, sDaiAmount)
	tDaiAmount := scaledSDaiAmount.Div(scaledSDaiAmount, sDAIPrice)
	return tDaiAmount

}

func (k Keeper) calculateTDaiAmountAndRoundUp(ctx sdk.Context, sDaiAmount *big.Int, sDAIPrice *big.Int) (*big.Int, error) {
	scaledSDaiAmount := scaleAmountBySDaiDecimals(ctx, sDaiAmount)
	tDaiAmount, err := DivideAndRoundUp(scaledSDaiAmount, sDAIPrice)
	if err != nil {
		return nil, err
	}
	return tDaiAmount, nil
}

func scaleAmountBySDaiDecimals(ctx sdk.Context, sDaiAmount *big.Int) *big.Int {
	tenScaledBySDaiDecimals := new(big.Int).Exp(
		big.NewInt(types.BASE_10),
		big.NewInt(types.SDAI_DECIMALS),
		nil)

	scaledSDaiAmount := new(big.Int).Mul(
		sDaiAmount,
		tenScaledBySDaiDecimals,
	)

	return scaledSDaiAmount
}

// DivideAndRoundUp performs division with rounding up: calculates x / y and rounds up to the nearest whole number
func DivideAndRoundUp(x *big.Int, y *big.Int) (*big.Int, error) {
	// Handle nil inputs
	if x == nil || y == nil {
		return nil, errors.New("input values cannot be nil")
	}

	// Handle negative inputs
	if x.Cmp(big.NewInt(0)) < 0 || y.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("input values cannot be negative")
	}

	// Handle division by zero
	if y.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("division by zero")
	}

	// Handle x being zero
	if x.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	result := new(big.Int).Sub(x, big.NewInt(1))
	result = result.Div(result, y)
	result = result.Add(result, big.NewInt(1))
	return result, nil
}
