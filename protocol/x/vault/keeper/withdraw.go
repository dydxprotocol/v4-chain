package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/vault"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetVaultWithdrawalSlippage returns the slippage that should be incurred from the specified
// vault on withdrawing `sharesToWithdraw` shares.
// For example, if `sharesToWithdraw = 100` and `0.2` is returned, it means that withdrawing
// 100 shares has a 20% slippage for the given `vaultId`.
//
// Slippage is calculated as `min(simple_slippage, estimated_slippage)` where:
// - simple_slippage = leverage * initial_margin
// - estimated_slippage = spread * (1 + average_skew) * leverage
//   - average_skew = integral / (posterior_leverage - leverage)
//   - integral = skew_antiderivative(skew_factor, posterior_leverage) -
//     skew_antiderivative(skew_factor, leverage)
//   - posterior_leverage = leverage / (1 - withdrawal_portion)
//     = leverage / (1 - shares_to_withdraw / total_shares)
//     = leverage * total_shares / (total_shares - shares_to_withdraw)
//
// To simplify above formula, let l = leverage, n = total_shares, m = shares_to_withdraw
//
//	estimated_slippage
//	= spread * (1 + integral / (posterior_leverage - leverage)) * leverage
//	= spread * (1 + integral * (n - m) / (ln - l(n - m))) * l
//	= spread * (1 + integral * (n - m) / lm) * l
//	= spread * (l + integral * (n - m) / m)
func (k Keeper) GetVaultWithdrawalSlippage(
	ctx sdk.Context,
	vaultId types.VaultId,
	sharesToWithdraw *big.Int,
	totalShares *big.Int,
	leverage *big.Rat,
	perpetual *perptypes.Perpetual,
	marketParam *pricestypes.MarketParam,
) (*big.Rat, error) {
	if sharesToWithdraw.Sign() <= 0 || sharesToWithdraw.Cmp(totalShares) > 0 {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidSharesToWithdraw,
			"sharesToWithdraw: %s, totalShares: %s",
			sharesToWithdraw,
			totalShares,
		)
	}

	quotingParams, exists := k.GetVaultQuotingParams(ctx, vaultId)
	if !exists {
		return nil, types.ErrVaultParamsNotFound
	}

	// No leverage, no slippage.
	if leverage.Sign() == 0 {
		return lib.BigRat0(), nil
	}

	// Use absolute value of leverage.
	leverage.Abs(leverage)

	// Calculate simple_slippage = leverage * initial_margin.
	lt, err := k.perpetualsKeeper.GetLiquidityTier(ctx, perpetual.Params.LiquidityTier)
	if err != nil {
		return nil, err
	}
	simpleSlippage := lib.BigRatMulPpm(leverage, lt.InitialMarginPpm)

	// Return simple slippage if withdrawing 100%.
	if sharesToWithdraw.Cmp(totalShares) == 0 {
		return simpleSlippage, nil
	}

	// Calculate estimated_slippage.
	// 1. leverage_after_withdrawal
	//    = leverage / (1 - withdrawal_portion)
	//    = leverage * total_shares / (total_shares - shares_to_withdraw)
	remainingShares := new(big.Int).Sub(totalShares, sharesToWithdraw)
	posteriorLeverage := new(big.Rat).Mul(
		leverage,
		new(big.Rat).SetFrac(totalShares, remainingShares),
	)

	// 2. integral = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
	integral := vault.SkewAntiderivative(quotingParams.SkewFactorPpm, posteriorLeverage)
	integral.Sub(integral, vault.SkewAntiderivative(quotingParams.SkewFactorPpm, leverage))

	// 3. estimated_slippage
	//    = spread * (l + integral * (n - m) / m)
	estimatedSlippage := new(big.Rat).Mul(
		integral,
		new(big.Rat).SetFrac(remainingShares, sharesToWithdraw),
	)
	estimatedSlippage.Add(
		estimatedSlippage,
		leverage,
	)
	estimatedSlippage = lib.BigRatMulPpm(
		estimatedSlippage,
		vault.SpreadPpm(&quotingParams, marketParam),
	)

	// Return min(simple_slippage, estimated_slippage).
	return lib.BigRatMin(
		simpleSlippage,
		estimatedSlippage,
	), nil
}

// WithdrawFromMegavault withdraws from megavault to a subaccount.
func (k Keeper) WithdrawFromMegavault(
	ctx sdk.Context,
	toSubaccount satypes.SubaccountId,
	sharesToWithdraw *big.Int,
	minQuoteQuantums *big.Int,
) (redeemedQuoteQuantums *big.Int, err error) {
	// 1. Check that the owner is withdrawing less that or equal to their unlocked shares.
	ownerShares, exists := k.GetOwnerShares(ctx, toSubaccount.Owner)
	if !exists {
		return nil, types.ErrOwnerNotFound
	}
	ownerShareUnlocks, _ := k.GetOwnerShareUnlocks(ctx, toSubaccount.Owner)
	ownerSharesAfterWithdrawal := ownerShares.NumShares.BigInt()
	ownerSharesAfterWithdrawal.Sub(ownerSharesAfterWithdrawal, sharesToWithdraw)
	ownerLockedShares := ownerShareUnlocks.GetTotalLockedShares()
	if ownerSharesAfterWithdrawal.Cmp(ownerLockedShares) < 0 {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientWithdrawableShares,
			"shares to withdraw: %s, owner total shares: %s, owner locked shares: %s",
			sharesToWithdraw,
			ownerShares,
			ownerLockedShares,
		)
	}

	// 2. Redeem from main vault.
	totalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	redeemedQuoteQuantums, err = k.GetSubaccountEquity(ctx, types.MegavaultMainSubaccount)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get megavault main vault equity", err)
		return nil, err
	}
	redeemedQuoteQuantums.Mul(redeemedQuoteQuantums, sharesToWithdraw)
	redeemedQuoteQuantums.Quo(redeemedQuoteQuantums, totalShares)

	// 3. Redeem from each sub vault.
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(vaultParamsIterator.Key())
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault ID from state key", err)
			continue
		}
		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

		_, perpetual, marketParam, marketPrice, err := k.GetVaultClobPerpAndMarket(ctx, *vaultId)
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get perpetual and market", err, "Vault ID", vaultId)
			continue
		}
		leverage, equity, err := k.GetVaultLeverageAndEquity(ctx, *vaultId, &perpetual, &marketPrice)
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault leverage and equity", err, "Vault ID", vaultId)
			continue
		}

		slippage, err := k.GetVaultWithdrawalSlippage(
			ctx,
			*vaultId,
			sharesToWithdraw,
			totalShares,
			leverage,
			&perpetual,
			&marketParam,
		)
		if err != nil {
			log.ErrorLogWithError(
				ctx,
				"Failed to get vault withdrawal slippage",
				err,
				"Vault ID",
				vaultId,
			)
			continue
		}

		// Transfer `equity * shares / totalShares * (1 - slippage)` from sub vault to main vault.
		redeemedFromSubVault := new(big.Rat).SetFrac(equity, big.NewInt(1))
		redeemedFromSubVault.Mul(redeemedFromSubVault, new(big.Rat).SetFrac(sharesToWithdraw, totalShares))
		redeemedFromSubVault.Mul(redeemedFromSubVault, new(big.Rat).Sub(lib.BigRat1(), slippage))
		quantumsToTransfer := new(big.Int).Quo(redeemedFromSubVault.Num(), redeemedFromSubVault.Denom())

		err = k.subaccountsKeeper.TransferFundsFromSubaccountToSubaccount(
			ctx,
			*vaultId.ToSubaccountId(),
			types.MegavaultMainSubaccount,
			assetstypes.AssetUsdc.Id,
			quantumsToTransfer,
		)
		if err != nil {
			log.ErrorLogWithError(
				ctx,
				"Failed to transfer from sub vault to main vault",
				err,
				"Vault ID",
				vaultId,
				"Quantums",
				quantumsToTransfer,
			)
			continue
		}

		// Increment total redeemed quote quantums.
		redeemedQuoteQuantums.Add(redeemedQuoteQuantums, quantumsToTransfer)
	}

	// 4. Return error if less than min quote quantums are redeemed.
	if redeemedQuoteQuantums.Cmp(minQuoteQuantums) < 0 {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientRedeemedQuoteQuantums,
			"redeemed quote quantums: %s, min quote quantums: %s",
			redeemedQuoteQuantums,
			minQuoteQuantums,
		)
	}

	// 5. Transfer from main vault to destination subaccount.
	err = k.subaccountsKeeper.TransferFundsFromSubaccountToSubaccount(
		ctx,
		types.MegavaultMainSubaccount,
		toSubaccount,
		assetstypes.AssetUsdc.Id,
		redeemedQuoteQuantums,
	)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"Failed to transfer from main vault to subaccount",
			err,
			"Subaccount",
			toSubaccount,
			"Quantums",
			redeemedQuoteQuantums,
		)
		return nil, err
	}

	// 6. Decrement total and owner shares.
	err = k.SetTotalShares(
		ctx,
		types.BigIntToNumShares(new(big.Int).Sub(totalShares, sharesToWithdraw)),
	)
	if err != nil {
		return nil, err
	}
	if ownerSharesAfterWithdrawal.Sign() == 0 {
		store := k.getOwnerSharesStore(ctx)
		store.Delete([]byte(toSubaccount.Owner))
	} else {
		err := k.SetOwnerShares(ctx, toSubaccount.Owner, types.BigIntToNumShares(ownerSharesAfterWithdrawal))
		if err != nil {
			return nil, err
		}
	}

	return redeemedQuoteQuantums, nil
}
