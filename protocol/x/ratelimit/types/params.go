package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
)

// BigBaselineMinimum1Hr defines the minimum baseline sDai for the 1-hour rate-limit.
var BigBaselineMinimum1Hr = new(big.Int).Mul(
	big.NewInt(1_000_000), // 1m full coins
	lib.BigPow10(-SDaiDenomExponent),
)

// BigBaselineMinimum1Day defines the minimum baseline sDai for the 1-day rate-limit.
var BigBaselineMinimum1Day = new(big.Int).Mul(
	big.NewInt(10_000_000), // 10m full coins
	lib.BigPow10(-SDaiDenomExponent),
)

var DefaultSDaiHourlyLimter = Limiter{
	Period:          3600 * time.Second,
	BaselineMinimum: dtypes.NewIntFromBigInt(BigBaselineMinimum1Hr),
	BaselineTvlPpm:  10_000, // 1%
}

var DefaultSDaiDailyLimiter = Limiter{
	Period:          24 * time.Hour,
	BaselineMinimum: dtypes.NewIntFromBigInt(BigBaselineMinimum1Day),
	BaselineTvlPpm:  100_000, // 10%
}

// DefaultSDaiRateLimitParams returns default rate-limit params for sDai.
func DefaultSDaiRateLimitParams() LimitParams {
	return LimitParams{
		Denom: SDaiDenom,
		Limiters: []Limiter{
			DefaultSDaiHourlyLimter,
			DefaultSDaiDailyLimiter,
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

		if limiter.BaselineMinimum.BigInt().Sign() <= 0 {
			return ErrInvalidBaselineMinimum
		}

		if limiter.BaselineTvlPpm == 0 || limiter.BaselineTvlPpm > lib.OneMillion {
			return ErrInvalidBaselineTvlPpm
		}
	}
	return nil
}
