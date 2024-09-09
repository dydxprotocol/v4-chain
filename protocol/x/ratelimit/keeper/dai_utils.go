package keeper

import (
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
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

func ConvertStringToBigInt(str string) (*big.Int, error) {

	bigint, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return nil, errorsmod.Wrap(
			types.ErrUnableToDecodeBigInt,
			"Unable to convert the sDAI conversion rate to a big int",
		)
	}

	return bigint, nil
}

func ConvertStringToBigIntWithPanicOnErr(str string) *big.Int {
	bigint, err := ConvertStringToBigInt(str)

	if err != nil {
		panic(fmt.Sprintf("Could not convert string to big.Int with err %v", err))
	}

	return bigint
}

func ConvertStringToBigRatWithPanicOnErr(str string) *big.Rat {
	bigrat, ok := new(big.Rat).SetString(str)

	if !ok {
		panic("Could not convert string to big.Rat")
	}

	return bigrat
}
