package vault

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SkewAntiderivativePpm returns the antiderivative of skew given a vault's skew
// factor and leverage.
// skew_antiderivative_ppm = skew_factor * leverage^2 + skew_factor^2 * leverage^3 / 3
func SkewAntiderivativePpm(
	skewFactorPpm uint32,
	leveragePpm *big.Int,
) *big.Int {
	bigSkewFactorPpm := new(big.Int).SetUint64(uint64(skewFactorPpm))
	bigOneTrillion := lib.BigIntOneTrillion()

	// a = skew_factor * leverage^2.
	a := new(big.Int).Mul(leveragePpm, leveragePpm)
	a.Mul(a, bigSkewFactorPpm)

	// b = skew_factor^2 * leverage^3 / 3.
	b := new(big.Int).Set(a)
	b.Mul(b, leveragePpm)
	b.Mul(b, bigSkewFactorPpm)
	b = lib.BigDivCeil(b, big.NewInt(3))

	// normalize `a` whose unit currently is ppm * ppm.
	a = lib.BigDivCeil(a, bigOneTrillion)
	// normalize `b` whose unit currently is ppm * ppm * ppm.
	b = lib.BigDivCeil(b, bigOneTrillion)
	b = lib.BigDivCeil(b, bigOneTrillion)

	// return a + b.
	return a.Add(a, b)
}

// SpreadPpm returns the spread that a vault should quote at given its
// quoting params and corresponding market param.
// spread_ppm = max(spread_min_ppm, spread_buffer_ppm + min_price_change_ppm)
func SpreadPpm(
	quotingParams *types.QuotingParams,
	marketParam *pricestypes.MarketParam,
) uint32 {
	return lib.Max(
		quotingParams.SpreadMinPpm,
		quotingParams.SpreadBufferPpm+marketParam.MinPriceChangePpm,
	)
}
