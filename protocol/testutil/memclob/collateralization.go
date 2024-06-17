package memclob

import (
	"fmt"
	"testing"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// CollateralizationCheck is a testing utility struct used to represent a collateralization check.
type CollateralizationCheck struct {
	CollatCheck map[satypes.SubaccountId][]clobtypes.PendingOpenOrder
	Result      map[satypes.SubaccountId]satypes.UpdateResult
}

// CreateCollatCheckFunction creates a collateralization check function that can be used when calling the `PlaceOrder`
// function on the memclob. It asserts that the passed in parameters to the collateralization check function are
// correct, and returns the specified collateralization check results.
// Note that this function takes a pointer parameter representing the number of collateralization checks, and
// increments it after each performed collateralization check such that the caller can assert the expected number of
// collateralization checks are performed.
func CreateCollatCheckFunction(
	t testing.TB,
	collateralCheckCounter *int,
	expectedCollatCheckParams map[int]map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
	collatCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult,
) (
	collatCheckFn clobtypes.AddOrderToOrderbookCollateralizationCheckFn,
) {
	collatCheckFn = func(
		subaccountMatchedOrders map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	) {
		// Before returning, increment the number of collateralization checks.
		defer func() {
			*collateralCheckCounter++
		}()

		collatCheckNum := *collateralCheckCounter

		// Verify the parameters are as expected.
		expectedSubaccountMatchedOrders := expectedCollatCheckParams[collatCheckNum]
		expectedNumSubaccounts := len(expectedSubaccountMatchedOrders)

		require.Len(
			t,
			subaccountMatchedOrders,
			expectedNumSubaccounts,
			fmt.Sprintf(
				"Different number of subaccounts. Collateral check %d",
				collatCheckNum,
			),
		)
		for subaccount, pendingMatches := range subaccountMatchedOrders {
			require.ElementsMatch(
				t,
				pendingMatches,
				expectedSubaccountMatchedOrders[subaccount],
				fmt.Sprintf(
					`Elements differ. List A is actual, list B is expected.
							Collateral check number: %d, subaccount.owner: %s, subaccount.number: %d`,
					collatCheckNum,
					subaccount.Owner,
					subaccount.Number,
				),
			)
		}

		// Return the result of the collateralization check.
		subaccountCollatCheckResult := make(map[satypes.SubaccountId]satypes.UpdateResult)
		success = true
		for subaccountId := range subaccountMatchedOrders {
			expectedUpdateResult, exists := collatCheckFailures[collatCheckNum][subaccountId]
			// If `collatCheckFailures` contains a successful update result, we should throw an error since
			// `collatCheckFailures` should not have entries for successful updates.
			// Else if no update result was specified for this subaccount, then we can assume success.
			if exists && expectedUpdateResult.IsSuccess() {
				require.Fail(
					t,
					fmt.Sprintf(
						"UpdateResult for collateralization check %d, subaccount %s should not be marked as successful.",
						collatCheckNum,
						subaccountId.String(),
					),
				)
			} else if !exists {
				expectedUpdateResult = satypes.Success
			}

			if !expectedUpdateResult.IsSuccess() {
				success = false
			}

			subaccountCollatCheckResult[subaccountId] = expectedUpdateResult
		}

		return success, subaccountCollatCheckResult
	}

	return collatCheckFn
}

// AlwaysSuccessfulCollatCheckFn is a collateralization check function that always returns success.
func AlwaysSuccessfulCollatCheckFn(
	subaccountMatchedOrders map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
) (success bool, successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult) {
	return true, nil
}

// CreateSimpleCollatCheckFunction creates a collateralization check function that can be used when
// calling the `PlaceOrder` function on the memclob. It asserts that the passed in parameters to
// the collateralization check function are correct, and returns the specified collateralization check results.
// Note that this function takes a pointer parameter representing the number of collateralization checks, and
// increments it after each performed collateralization check such that the caller can assert the expected number of
// collateralization checks are performed.
func CreateSimpleCollatCheckFunction(
	t testing.TB,
	collateralCheckCounter *int,
	expectedCollatCheck map[int]CollateralizationCheck,
) (
	collatCheckFn clobtypes.AddOrderToOrderbookCollateralizationCheckFn,
) {
	collatCheckFn = func(
		subaccountOpenOrders map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	) {
		// Before returning, increment the number of collateralization checks.
		defer func() {
			*collateralCheckCounter++
		}()

		collatCheckNum := *collateralCheckCounter

		// Verify the parameters are as expected.
		collatCheck := expectedCollatCheck[collatCheckNum]
		expectedNumSubaccounts := len(collatCheck.CollatCheck)

		require.Len(
			t,
			subaccountOpenOrders,
			expectedNumSubaccounts,
			fmt.Sprintf(
				"Different number of subaccounts. Collateral check %d",
				collatCheckNum,
			),
		)
		success = true
		for subaccount, openOrders := range subaccountOpenOrders {
			require.ElementsMatch(
				t,
				openOrders,
				collatCheck.CollatCheck[subaccount],
				fmt.Sprintf(
					`Elements differ. List A is actual, list B is expected.
							Collateral check number: %d, subaccount.owner: %s, subaccount.number: %d`,
					collatCheckNum,
					subaccount.Owner,
					subaccount.Number,
				),
			)

			if !collatCheck.Result[subaccount].IsSuccess() {
				success = false
			}
		}

		// Return the result of the collateralization check.
		return success, collatCheck.Result
	}

	return collatCheckFn
}
