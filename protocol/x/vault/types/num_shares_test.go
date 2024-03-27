package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestNumShares_ToBigRat(t *testing.T) {
	tests := map[string]struct {
		n           types.NumShares
		expectedRat *big.Rat
		expectedErr error
	}{
		"Success - 1/1": {
			n: types.NumShares{
				Numerator:   dtypes.NewInt(1),
				Denominator: dtypes.NewInt(1),
			},
			expectedRat: big.NewRat(1, 1),
		},
		"Success - 1/77": {
			n: types.NumShares{
				Numerator:   dtypes.NewInt(1),
				Denominator: dtypes.NewInt(77),
			},
			expectedRat: big.NewRat(1, 77),
		},
		"Success - 99/2": {
			n: types.NumShares{
				Numerator:   dtypes.NewInt(99),
				Denominator: dtypes.NewInt(2),
			},
			expectedRat: big.NewRat(99, 2),
		},
		"Failure - numerator is nil": {
			n: types.NumShares{
				Denominator: dtypes.NewInt(77),
			},
			expectedErr: types.ErrNilFraction,
		},
		"Failure - denominator is nil": {
			n: types.NumShares{
				Numerator: dtypes.NewInt(77),
			},
			expectedErr: types.ErrNilFraction,
		},
		"Failure - denominator is 0": {
			n: types.NumShares{
				Numerator:   dtypes.NewInt(77),
				Denominator: dtypes.NewInt(0),
			},
			expectedErr: types.ErrZeroDenominator,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rat, err := tc.n.ToBigRat()
			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRat, rat)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
				require.True(t, rat == nil)
			}
		})
	}
}

func TestBigRatToNumShares(t *testing.T) {
	tests := map[string]struct {
		rat               *big.Rat
		expectedNumShares types.NumShares
	}{
		"Success - 1/1": {
			rat: big.NewRat(1, 1),
			expectedNumShares: types.NumShares{
				Numerator:   dtypes.NewInt(1),
				Denominator: dtypes.NewInt(1),
			},
		},
		"Success - 1/2": {
			rat: big.NewRat(1, 2),
			expectedNumShares: types.NumShares{
				Numerator:   dtypes.NewInt(1),
				Denominator: dtypes.NewInt(2),
			},
		},
		"Success - 5/3": {
			rat: big.NewRat(5, 3),
			expectedNumShares: types.NumShares{
				Numerator:   dtypes.NewInt(5),
				Denominator: dtypes.NewInt(3),
			},
		},
		"Success - rat is nil": {
			rat:               nil,
			expectedNumShares: types.NumShares{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			n := types.BigRatToNumShares(tc.rat)
			require.Equal(t, tc.expectedNumShares, n)
		})
	}
}
