package vault

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SkewAntiderivative returns the antiderivative of skew given a vault's skew
// factor and leverage.
// skew_antiderivative = skew_factor * leverage^2 + skew_factor^2 * leverage^3 / 3
func SkewAntiderivative(
	skewFactorPpm uint32,
	leverage *big.Rat,
) *big.Rat {
	bigSkewFactorPpm := new(big.Rat).SetUint64(uint64(skewFactorPpm))
	bigOneMillion := lib.BigRatOneMillion()

	// a = skew_factor * leverage^2.
	a := new(big.Rat).Mul(leverage, leverage)
	a.Mul(a, bigSkewFactorPpm)

	// b = skew_factor^2 * leverage^3 / 3.
	b := new(big.Rat).Set(a)
	b.Mul(b, leverage)
	b.Mul(b, bigSkewFactorPpm)
	b.Quo(b, big.NewRat(3, 1))

	// normalize `a` whose unit currently is ppm.
	a.Quo(a, bigOneMillion)
	// normalize `b` whose unit currently is ppm * ppm.
	b.Quo(b, bigOneMillion)
	b.Quo(b, bigOneMillion)

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
