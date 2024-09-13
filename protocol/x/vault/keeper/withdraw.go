package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/vault"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetVaultWithdrawalSlippagePpm returns the slippage (in ppm) that should be incurred
// on withdrawing `withdrawalPortionPpm` of a vault's ownership.
// For example, if `withdrawalPortionPpm = 100_000` and `200_000` is returned,
// it means that withdrawing 10% has a 20% slippage for the specified vault.
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
func (k Keeper) GetVaultWithdrawalSlippagePpm(
	ctx sdk.Context,
	vaultId types.VaultId,
	sharesToWithdraw *big.Int,
) (*big.Int, error) {
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
	leveragePpm, _, err := k.GetVaultLeverageAndEquity(ctx, vaultId, perpetual, marketPrice)
	if err != nil {
		return nil, err
	}
	// No leverage, no slippage.
	if leveragePpm.Sign() == 0 {
		return big.NewInt(0), nil
	}

	// Use absolute value of leverage.
	leveragePpm.Abs(leveragePpm)

	// Calculate simple_slippage = leverage * initial_margin (round up if necessary).
	lt, err := k.perpetualsKeeper.GetLiquidityTier(ctx, perpetual.Params.LiquidityTier)
	if err != nil {
		return nil, err
	}
	bigOneMillion := lib.BigIntOneMillion()
	simpleSlippagePpm := new(big.Int).Mul(
		leveragePpm,
		new(big.Int).SetUint64(uint64(lt.InitialMarginPpm)),
	)
	simpleSlippagePpm = lib.BigDivCeil(simpleSlippagePpm, bigOneMillion)

	// Return simple slippage if withdrawing 100%.
	if sharesToWithdraw.Cmp(totalShares) == 0 {
		return simpleSlippagePpm, nil
	}

	// Estimate slippage.
	// 1. leverage_after_withdrawal
	//    = leverage / (1 - withdrawal_portion)
	//    = leverage * total_shares / (total_shares - shares_to_withdraw)
	posteriorLeveragePpm := new(big.Int).Mul(leveragePpm, totalShares)
	posteriorLeveragePpm = lib.BigDivCeil(
		posteriorLeveragePpm,
		new(big.Int).Sub(totalShares, sharesToWithdraw),
	)

	// 2. integral = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
	estimatedSlippagePpm := vault.SkewAntiderivativePpm(quotingParams.SkewFactorPpm, posteriorLeveragePpm)
	estimatedSlippagePpm.Sub(estimatedSlippagePpm, vault.SkewAntiderivativePpm(quotingParams.SkewFactorPpm, leveragePpm))

	// 3. average_skew = integral / (posterior_leverage - leverage)
	estimatedSlippagePpm.Mul(estimatedSlippagePpm, bigOneMillion)
	estimatedSlippagePpm = lib.BigDivCeil(
		estimatedSlippagePpm,
		posteriorLeveragePpm.Sub(posteriorLeveragePpm, leveragePpm),
	)

	// 4. slippage = spread * (1 + average_skew) * leverage
	estimatedSlippagePpm.Add(estimatedSlippagePpm, bigOneMillion)
	estimatedSlippagePpm.Mul(estimatedSlippagePpm, leveragePpm)
	estimatedSlippagePpm = lib.BigIntMulPpm(
		estimatedSlippagePpm,
		vault.SpreadPpm(&quotingParams, &marketParam),
	)
	estimatedSlippagePpm = lib.BigDivCeil(estimatedSlippagePpm, bigOneMillion)

	// Return min(simple_slippage, estimated_slippage).
	return lib.BigMin(
		simpleSlippagePpm,
		estimatedSlippagePpm,
	), nil
}
