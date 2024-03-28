package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestBigIntToNumShares(t *testing.T) {
	tests := map[string]struct {
		num               *big.Int
		expectedNumShares types.NumShares
	}{
		"Success - 1": {
			num: big.NewInt(1),
			expectedNumShares: types.NumShares{
				NumShares: dtypes.NewInt(1),
			},
		},
		"Success - num is nil": {
			num:               nil,
			expectedNumShares: types.NumShares{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			n := types.BigIntToNumShares(tc.num)
			require.Equal(t, tc.expectedNumShares, n)
		})
	}
}
