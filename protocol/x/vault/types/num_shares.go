package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
)

// ToBigRat converts a NumShares to a big.Rat.
// Returns an error if the denominator is zero.
func (n NumShares) ToBigRat() (*big.Rat, error) {
	if n.Numerator.IsNil() || n.Denominator.IsNil() {
		return nil, ErrNilFraction
	}
	if n.Denominator.Cmp(dtypes.NewInt(0)) == 0 {
		return nil, ErrZeroDenominator
	}

	return new(big.Rat).SetFrac(n.Numerator.BigInt(), n.Denominator.BigInt()), nil
}

func BigRatToNumShares(rat *big.Rat) (n NumShares) {
	if rat == nil {
		return n
	}
	return NumShares{
		Numerator:   dtypes.NewIntFromBigInt(rat.Num()),
		Denominator: dtypes.NewIntFromBigInt(rat.Denom()),
	}
}
