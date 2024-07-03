package keeper

import (
	"errors"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetSDAIToTradingDAI(ctx sdk.Context, sDAI *big.Int) (*big.Int, error) {
	// Get the current sDAI price
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return nil, errors.New("sDAI price not found")
	}

	/*
		// implementing the maker code
		// https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea
		function deposit(uint256 assets, address receiver) external returns (uint256 shares) {
			uint256 chi = (block.timestamp > pot.rho()) ? pot.drip() : pot.chi();
			shares = assets * RAY / chi;
			_mint(assets, shares, receiver);
		}
	*/

	// Calculate shares = assets * RAY / chi
	shares := new(big.Int).Mul(sDAI, new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil))
	shares = shares.Div(shares, sDAIPrice)

	return shares, nil
}

func (k Keeper) GetTradingDAIToSDAI(ctx sdk.Context, tradingDAI *big.Int) (*big.Int, error) {
	// Get the current sDAI price
	sDAIPrice, found := k.GetSDAIPrice(ctx)
	if !found {
		return nil, errors.New("sDAI price not found")
	}

	/*
		// implementing the maker code
		// https://etherscan.deth.net/address/0x83f20f44975d03b1b09e64809b757c47f942beea
		function redeem(uint256 shares, address receiver, address owner) external returns (uint256 assets) {
			uint256 chi = (block.timestamp > pot.rho()) ? pot.drip() : pot.chi();
			assets = shares * chi / RAY;
			_burn(assets, shares, receiver, owner);
		}
	*/

	// Calculate shares = tradingDAI * RAY / sDAIPrice
	shares := new(big.Int).Mul(tradingDAI, sDAIPrice)
	shares = shares.Div(shares, new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil))
	return shares, nil
}
