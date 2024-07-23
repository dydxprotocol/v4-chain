package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/testutils"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

func TestIsTimestampNonce(t *testing.T) {
	t.Run("IsValidTimestampNonce test", func(t *testing.T) {
		require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff))
		require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff+uint64(1)))
		require.False(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff-uint64(1)))
		require.False(t, keeper.IsTimestampNonce(0))
		require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff+uint64(100000)))
		require.False(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff-uint64(100000)))
	})

}

func TestIsValidTimestampNonce(t *testing.T) {
	t.Run("IsValidTimestampNonce test", func(t *testing.T) {
		referenceTs := keeper.TimestampNonceSequenceCutoff + 10000
		require.True(t, keeper.IsValidTimestampNonce(referenceTs, referenceTs))
		require.True(t, keeper.IsValidTimestampNonce(referenceTs+uint64(1), referenceTs))
		require.True(t, keeper.IsValidTimestampNonce(referenceTs-uint64(1), referenceTs))
		require.False(t, keeper.IsValidTimestampNonce(referenceTs+uint64(100000), referenceTs))
		require.False(t, keeper.IsValidTimestampNonce(referenceTs-uint64(100000), referenceTs))
	})
}

func TestEjectStaleTsNonces(t *testing.T) {
	t.Run("Will eject stale timestamps", func(t *testing.T) {
		baseTsNonce := keeper.TimestampNonceSequenceCutoff
		tsNonces := make([]uint64, keeper.MaxTimestampNonceArrSize)
		for i := 0; i < keeper.MaxTimestampNonceArrSize; i++ {
			tsNonces[i] = baseTsNonce + uint64(i)
		}
		accountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: baseTsNonce,
			},
		}

		shift := uint64(5)
		referenceTs := baseTsNonce + keeper.MaxTimeInPastMs + shift
		expectedMaxEjectedNonce := referenceTs - keeper.MaxTimestampNonceArrSize - 1
		var expectedTsNonces []uint64
		for _, ts := range genesisState.Accounts[0].TimestampNonceDetails.TimestampNonces {
			if ts > expectedMaxEjectedNonce {
				expectedTsNonces = append(expectedTsNonces, ts)
			}
		}
		for i := shift; i < keeper.MaxTimestampNonceArrSize-shift; i++ {
			tsNonces[i] = baseTsNonce + keeper.MaxTimeInPastMs + uint64(i)
		}
		expectedTsNonceDetails := types.TimestampNonceDetails{
			TimestampNonces: expectedTsNonces,
			MaxEjectedNonce: expectedMaxEjectedNonce,
		}

		k.EjectStaleTsNonces(ctx, constants.AliceAccAddress, referenceTs)

		actualAccountState, didUpdate := k.GetAccountState(ctx, constants.AliceAccAddress)

		require.True(t, didUpdate)

		require.True(t,
			testutils.CompareTimestampNonceDetails(actualAccountState.TimestampNonceDetails, &expectedTsNonceDetails),
		)
	})
}
