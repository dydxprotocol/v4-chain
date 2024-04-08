package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestSubticksToBigInt(t *testing.T) {
	num := uint64(5)
	st := types.Subticks(5)

	require.Zero(t, new(big.Int).SetUint64(num).Cmp(st.ToBigInt()))
}

func TestSubticksToBigRat(t *testing.T) {
	num := uint64(5)
	st := types.Subticks(5)

	require.Zero(t, new(big.Rat).SetUint64(num).Cmp(st.ToBigRat()))
}

func TestSubticksToUInt64(t *testing.T) {
	num := uint64(5)
	st := types.Subticks(5)

	require.Equal(t, num, st.ToUint64())
}

func TestString(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		orderStatus types.OrderStatus

		// Expectations.
		expectedString string
	}{
		"Order status is Success": {
			orderStatus: types.Success,

			expectedString: "Success",
		},
		"Order status is Undercollateralized": {
			orderStatus: types.Undercollateralized,

			expectedString: "Undercollateralized",
		},
		"Order status is InternalError": {
			orderStatus: types.InternalError,

			expectedString: "InternalError",
		},
		"Order status is ImmediateOrCancelWouldRestOnBook": {
			orderStatus: types.ImmediateOrCancelWouldRestOnBook,

			expectedString: "ImmediateOrCancelWouldRestOnBook",
		},
		"Order status is ReduceOnlyResized": {
			orderStatus: types.ReduceOnlyResized,

			expectedString: "ReduceOnlyResized",
		},
		"Order status is LiquidationRequiresDeleveraging": {
			orderStatus: types.LiquidationRequiresDeleveraging,

			expectedString: "LiquidationRequiresDeleveraging",
		},
		"Order status is LiquidationExceededSubaccountMaxNotionalLiquidated": {
			orderStatus: types.LiquidationExceededSubaccountMaxNotionalLiquidated,

			expectedString: "LiquidationExceededSubaccountMaxNotionalLiquidated",
		},
		"Order status is LiquidationExceededSubaccountMaxInsuranceLost": {
			orderStatus: types.LiquidationExceededSubaccountMaxInsuranceLost,

			expectedString: "LiquidationExceededSubaccountMaxInsuranceLost",
		},
		"Order status is ViolatesIsolatedSubaccountConstraints": {
			orderStatus: types.ViolatesIsolatedSubaccountConstraints,

			expectedString: "ViolatesIsolatedSubaccountConstraints",
		},
		"Order status is unknown enum value": {
			orderStatus: 999,

			expectedString: "Unknown",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualString := tc.orderStatus.String()
			require.Equal(t, tc.expectedString, actualString)
		})
	}
}

func TestIsSuccess(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		orderStatus types.OrderStatus

		// Expectations.
		expectedIsSuccess bool
	}{
		"Order status of Success is successful": {
			orderStatus: types.Success,

			expectedIsSuccess: true,
		},
		"Order status of Undercollateralized is not successful": {
			orderStatus: types.Undercollateralized,

			expectedIsSuccess: false,
		},
		"Order status of InternalError is not successful": {
			orderStatus: types.InternalError,

			expectedIsSuccess: false,
		},
		"Order status of unknown enum value is not successful": {
			orderStatus: 10,

			expectedIsSuccess: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualIsSuccess := tc.orderStatus.IsSuccess()
			require.Equal(t, tc.expectedIsSuccess, actualIsSuccess)
		})
	}
}
