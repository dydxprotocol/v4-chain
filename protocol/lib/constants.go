package lib

import (
	"math"
	"math/big"

	sdkmath "cosmossdk.io/math"
)

const (
	OneMillion         = uint32(1_000_000)
	OneHundredThousand = uint32(100_000)
	TenThousand        = uint32(10_000)
	OneHundred         = uint32(100)
	MaxPriceChangePpm  = uint32(10_000)
	// 10^6 quantums == 1 USD.
	QuoteCurrencyAtomicResolution = int32(-6)

	ZeroUint64 = uint64(0)

	// 10^BaseDenomExponent denotes how much full coin is represented by 1 base denom.
	BaseDenomExponent = -18
	DefaultBaseDenom  = "adv4tnt"
)

// PowerReduction defines the default power reduction value for staking.
// Use 1e18, since default stake denom is assumed to be 1e-18 of a full coin.
var PowerReduction = sdkmath.NewIntFromBigInt(
	new(big.Int).SetUint64(1_000_000_000_000_000_000),
)

// BigNegMaxUint64 returns a `big.Int` that is set to -math.MaxUint64.
func BigNegMaxUint64() *big.Int {
	return new(big.Int).Neg(
		new(big.Int).SetUint64(math.MaxUint64),
	)
}

// BigMaxInt32 returns a `big.Int` that represents `MaxInt32`.
func BigMaxInt32() *big.Int {
	return big.NewInt(math.MaxInt32)
}

// BigFloatMaxUint64 returns a `big.Float` that is set to MaxUint64.
func BigFloatMaxUint64() *big.Float {
	return new(big.Float).SetUint64(math.MaxUint64)
}

// BigIntOneMillion returns a `big.Int` that is set to 1_000_000.
func BigIntOneMillion() *big.Int {
	return big.NewInt(1_000_000)
}

// BigIntOneTrillion returns a `big.Int` that is set to 1_000_000_000_000.
func BigIntOneTrillion() *big.Int {
	return big.NewInt(1_000_000_000_000)
}

// BigRatOneMillion returns a `big.Rat` that is set to 1_000_000.
func BigRatOneMillion() *big.Rat {
	return big.NewRat(1_000_000, 1)
}

// BigRat0 returns a `big.Rat` that is set to 0.
func BigRat0() *big.Rat {
	return big.NewRat(0, 1)
}

// BigRat1 returns a `big.Rat` that is set to 1.
func BigRat1() *big.Rat {
	return big.NewRat(1, 1)
}
