package testutils

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/google/go-cmp/cmp"
)

func CompareTimestampNonceDetails(actualDetails, expectedDetails *types.TimestampNonceDetails) bool {
	if !cmp.Equal(
		actualDetails.GetTimestampNonces(),
		expectedDetails.GetTimestampNonces(),
	) {
		return false
	}

	if actualDetails.GetMaxEjectedNonce() != expectedDetails.GetMaxEjectedNonce() {
		return false
	}

	return true
}

func CompareAccountStates(actualAccountstate, expectedAccountState *types.AccountState) bool {
	if actualAccountstate.GetAddress() != expectedAccountState.GetAddress() {
		return false
	}

	if tsNonceDetailsEqual := CompareTimestampNonceDetails(
		actualAccountstate.GetTimestampNonceDetails(),
		expectedAccountState.GetTimestampNonceDetails(),
	); !tsNonceDetailsEqual {
		return false
	}

	return true
}

func CompareAccountStateLists(actualAccountStates, expectedAccountStates []*types.AccountState) bool {
	if len(actualAccountStates) != len(expectedAccountStates) {
		return false
	}

	// We require that the ordering of accountState be deterministic (no sorting) so that should more
	// complicated logic be introduced in the future, this test can catch any unintended effects.
	for i := range actualAccountStates {
		if !CompareAccountStates(actualAccountStates[i], expectedAccountStates[i]) {
			return false
		}
	}

	return true
}
