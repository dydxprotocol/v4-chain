package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/vault"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetVaultWithdrawalSlippagePpm returns the slippage that should be incurred from the specified
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
) (*big.Rat, error) {
	totalShares := k.GetTotalShares(ctx).NumShares.BigInt()
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

	_, perpetual, marketParam, marketPrice, err := k.GetVaultClobPerpAndMarket(ctx, vaultId)
	if err != nil {
		return nil, err
	}

	// Get vault leverage.
	leverage, _, err := k.GetVaultLeverageAndEquity(ctx, vaultId, perpetual, marketPrice)
	if err != nil {
		return nil, err
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
		vault.SpreadPpm(&quotingParams, &marketParam),
	)

	// Return min(simple_slippage, estimated_slippage).
	return lib.BigRatMin(
		simpleSlippage,
		estimatedSlippage,
	), nil
}
