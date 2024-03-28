package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
)

// BigIntToNumShares returns a NumShares given a big.Int.
func BigIntToNumShares(num *big.Int) (n NumShares) {
	if num == nil {
		return n
	}
	return NumShares{
		NumShares: dtypes.NewIntFromBigInt(num),
	}
}
