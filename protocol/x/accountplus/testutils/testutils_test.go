package testutils_test

import (
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/testutils"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

var baseTsNonce = uint64(math.Pow(2, 40))

func TestCompareTimestampNonceDetails(t *testing.T) {
	detail1 := &types.TimestampNonceDetails{
		TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
		MaxEjectedNonce: baseTsNonce,
	}

	detail2 := &types.TimestampNonceDetails{
		TimestampNonces: []uint64{baseTsNonce + 2, baseTsNonce + 3, baseTsNonce + 3},
		MaxEjectedNonce: baseTsNonce,
	}

	detail3 := &types.TimestampNonceDetails{
		TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
		MaxEjectedNonce: baseTsNonce + 1,
	}

	detail4 := &types.TimestampNonceDetails{
		TimestampNonces: []uint64{baseTsNonce + 2, baseTsNonce + 2, baseTsNonce + 3},
		MaxEjectedNonce: baseTsNonce + 1,
	}

	emptyDetail := &types.TimestampNonceDetails{}

	var isEqual bool

	// Same detail
	isEqual = testutils.CompareTimestampNonceDetails(detail1, detail1)
	require.True(t, isEqual, "Details are equal but comparison returned false")

	// Different TimestampNonces
	isEqual = testutils.CompareTimestampNonceDetails(detail1, detail2)
	require.False(t, isEqual, "TimmestampNonces are different but comparison returned true")

	// Different MaxEjectedNonce
	isEqual = testutils.CompareTimestampNonceDetails(detail1, detail3)
	require.False(t, isEqual, "MaxEjectedNonce are different but comparison returned true")

	// Different TimestampNonces and MaxEjectedNonce
	isEqual = testutils.CompareTimestampNonceDetails(detail1, detail4)
	require.False(t, isEqual, "TimestampNonces and MaxEjectedNonce are different but comparison returned true")

	// Empty detail
	isEqual = testutils.CompareTimestampNonceDetails(detail1, emptyDetail)
	require.False(t, isEqual, "TimestampNonces and MaxEjectedNonce are different but comparison returned true")
}

func TestCompareAccountStates(t *testing.T) {
	accStateReference := &types.AccountState{
		Address: constants.AliceAccAddress.String(),
		TimestampNonceDetails: &types.TimestampNonceDetails{
			TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
			MaxEjectedNonce: baseTsNonce,
		},
	}

	accStateDiffAddress := &types.AccountState{
		Address: constants.BobAccAddress.String(),
		TimestampNonceDetails: &types.TimestampNonceDetails{
			TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
			MaxEjectedNonce: baseTsNonce,
		},
	}

	accStateDiffTimestampNonceDetails := &types.AccountState{
		Address: constants.AliceAccAddress.String(),
		TimestampNonceDetails: &types.TimestampNonceDetails{
			TimestampNonces: []uint64{baseTsNonce + 3, baseTsNonce + 2, baseTsNonce + 3},
			MaxEjectedNonce: baseTsNonce,
		},
	}

	var isEqual bool

	// Same AccountState
	isEqual = testutils.CompareAccountStates(accStateReference, accStateReference)
	require.True(t, isEqual, "AccountStates are equal but comparison returned false")

	// Different address
	isEqual = testutils.CompareAccountStates(accStateReference, accStateDiffAddress)
	require.False(t, isEqual, "AccountStates have different address but comparison returned true")

	// Different TimestampNonceDetails
	isEqual = testutils.CompareAccountStates(accStateReference, accStateDiffTimestampNonceDetails)
	require.False(t, isEqual, "AccountStates have different TimestampNonceDetails but comparison returned true")
}

func TestCompareAccoutStateLists(t *testing.T) {
	accountStates := []*types.AccountState{
		{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: &types.TimestampNonceDetails{
				TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
				MaxEjectedNonce: baseTsNonce,
			},
		},
		{
			Address: constants.BobAccAddress.String(),
			TimestampNonceDetails: &types.TimestampNonceDetails{
				TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
				MaxEjectedNonce: baseTsNonce + 1,
			},
		},
	}

	accountStatesDifferent := []*types.AccountState{
		{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: &types.TimestampNonceDetails{
				TimestampNonces: []uint64{baseTsNonce + 2, baseTsNonce + 2, baseTsNonce + 3},
				MaxEjectedNonce: baseTsNonce,
			},
		},
		{
			Address: constants.BobAccAddress.String(),
			TimestampNonceDetails: &types.TimestampNonceDetails{
				TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
				MaxEjectedNonce: baseTsNonce + 1,
			},
		},
	}

	var isEqual bool

	isEqual = testutils.CompareAccountStateLists(accountStates, accountStates)
	require.True(t, isEqual, "AccountState lists are different but comparison returned false")

	isEqual = testutils.CompareAccountStateLists(accountStates, accountStatesDifferent)
	require.False(t, isEqual, "AccountState lists are different but comparison returned true")
}
