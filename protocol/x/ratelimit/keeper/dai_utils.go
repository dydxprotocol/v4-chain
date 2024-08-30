package keeper

import (
	"errors"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
)

func divideAmountBySDaiDecimals(scaledSDaiAmount *big.Int) *big.Int {
	tenScaledBySDaiDecimals := getTenScaledBySDaiDecimals()

	dividedAmount := new(big.Int).Div(
		scaledSDaiAmount,
		tenScaledBySDaiDecimals,
	)

	return dividedAmount
}

func getTenScaledBySDaiDecimals() *big.Int {
	return new(big.Int).Exp(
		big.NewInt(types.BASE_10),
		big.NewInt(types.SDAI_DECIMALS),
		nil,
	)
}

// DivideAndRoundUp performs division with rounding up: calculates x / y and rounds up to the nearest whole number
func divideAndRoundUp(x *big.Int, y *big.Int) (*big.Int, error) {
	// Handle nil inputs
	if x == nil || y == nil {
		return nil, errors.New("input values cannot be nil")
	}

	// Handle negative inputs
	if x.Cmp(big.NewInt(0)) < 0 || y.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("input values cannot be negative")
	}

	// Handle division by zero
	if y.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("division by zero")
	}

	// Handle x being zero
	if x.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	result := new(big.Int).Sub(x, big.NewInt(1))
	result = result.Div(result, y)
	result = result.Add(result, big.NewInt(1))
	return result, nil
}
