package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

// BigBaselineMinimum1Hr defines the minimum baseline USDC for the 1-hour rate-limit.
var BigBaselineMinimum1Hr = lib.BigIntMulPow10(
	big.NewInt(1_000_000), // 1m full coins
	-assettypes.UusdcDenomExponent,
	false,
)

// BigBaselineMinimum1Day defines the minimum baseline USDC for the 1-day rate-limit.
var BigBaselineMinimum1Day = lib.BigIntMulPow10(
	big.NewInt(10_000_000), // 10m full coins
	-assettypes.UusdcDenomExponent,
	false,
)

var DefaultUsdcHourlyLimter = Limiter{
	Period:          3600 * time.Second,
	BaselineMinimum: dtypes.NewIntFromBigInt(BigBaselineMinimum1Hr),
	BaselineTvlPpm:  10_000, // 1%
}

var DefaultUsdcDailyLimiter = Limiter{
	Period:          24 * time.Hour,
	BaselineMinimum: dtypes.NewIntFromBigInt(BigBaselineMinimum1Day),
	BaselineTvlPpm:  100_000, // 10%
}

// DefaultUsdcRateLimitParams returns default rate-limit params for USDC.
func DefaultUsdcRateLimitParams() LimitParams {
	return LimitParams{
		Denom: assettypes.UusdcDenom,
		Limiters: []Limiter{
			DefaultUsdcHourlyLimter,
			DefaultUsdcDailyLimiter,
		},
	}
}

// Validate validates the set of params
func (p *LimitParams) Validate() error {
	if err := sdk.ValidateDenom(p.Denom); err != nil {
		return err
	}

	for _, limiter := range p.Limiters {
		if limiter.Period == 0 {
			return ErrInvalidRateLimitPeriod
		}

		if limiter.BaselineMinimum.Sign() <= 0 {
			return ErrInvalidBaselineMinimum
		}

		if limiter.BaselineTvlPpm == 0 || limiter.BaselineTvlPpm > lib.OneMillion {
			return ErrInvalidBaselineTvlPpm
		}
	}
	return nil
}
