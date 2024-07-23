package keeper_test

import (
	"testing"
<<<<<<< HEAD

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDeepCopyTimestampNonceDetails(t *testing.T) {
	details := keeper.InitialTimestampNonceDetails
	detailsCopy := keeper.DeepCopyTimestampNonceDetails(details)

	detailsCopy.MaxEjectedNonce = details.MaxEjectedNonce + 1
	detailsCopy.TimestampNonces = append(detailsCopy.TimestampNonces, []uint64{1, 2, 3}...)

	require.NotEqual(t, details.MaxEjectedNonce, detailsCopy.MaxEjectedNonce)
	require.False(
		t,
		cmp.Equal(details.GetTimestampNonces(), detailsCopy.GetTimestampNonces()),
		"TimestampNonces not deepcopy",
	)
}
=======
)

func Placeholder(t *testing.T) {}
>>>>>>> main
