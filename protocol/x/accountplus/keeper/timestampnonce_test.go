package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

// Unit tests for ProcessTimestampNonce are not explicitly defined. sigverify_test.go contains test cases that cover
// ProcessTimestampNonce. TODO (low priority): move unit tests from sigverify_test.go to here.
// https://linear.app/dydx/issue/OTE-634/add-explicit-unit-tests-for-accountpluskeeper-processtimestampnonce

func TestIsTimestampNonce(t *testing.T) {
	tests := map[string]struct {
		tsNonce      uint64
		expectedBool bool
	}{
		"At cutoff": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff,
			expectedBool: true,
		},
		"Above cutoff": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff + uint64(100000),
			expectedBool: true,
		},
		"Below cutoff": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff - uint64(100000),
			expectedBool: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedBool, keeper.IsTimestampNonce(tc.tsNonce))
		})
	}
}

func TestIsValidTimestampNonce(t *testing.T) {
	tests := map[string]struct {
		tsNonce      uint64
		referenceTs  uint64
		expectedBool bool
	}{
		"Valid": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff,
			referenceTs:  keeper.TimestampNonceSequenceCutoff + 10000,
			expectedBool: true,
		},
		"Too old": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff,
			referenceTs:  keeper.TimestampNonceSequenceCutoff + 100000,
			expectedBool: true,
		},
		"Below early": {
			tsNonce:      keeper.TimestampNonceSequenceCutoff + 100000,
			referenceTs:  keeper.TimestampNonceSequenceCutoff,
			expectedBool: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedBool, keeper.IsTimestampNonce(tc.tsNonce))
		})
	}
}

func TestEjectStaleTsNonces(t *testing.T) {
	startTs := keeper.TimestampNonceSequenceCutoff

	tests := map[string]struct {
		timeElapsed             uint64
		expectedMaxEjectedNonce uint64
	}{
		"Will eject stale timestamp nonces": {
			timeElapsed:             keeper.MaxTimeInPastMs + 5,
			expectedMaxEjectedNonce: startTs + 5 - 1,
		},
		"Will not eject non-stale timestamp nonces": {
			timeElapsed:             keeper.MaxTimeInPastMs,
			expectedMaxEjectedNonce: startTs,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Starting state
			tsNonces := make([]uint64, keeper.MaxTimestampNonceArrSize)
			for i := 0; i < keeper.MaxTimestampNonceArrSize; i++ {
				tsNonces[i] = startTs + uint64(i) + 1
			}
			accountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: tsNonces,
					MaxEjectedNonce: startTs,
				},
			}

			// Expected state after ejection
			referenceTs := startTs + tc.timeElapsed

			var expectedTsNonces []uint64
			for _, ts := range accountState.TimestampNonceDetails.TimestampNonces {
				if ts > tc.expectedMaxEjectedNonce {
					expectedTsNonces = append(expectedTsNonces, ts)
				}
			}
			expectedAccountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: expectedTsNonces,
					MaxEjectedNonce: tc.expectedMaxEjectedNonce,
				},
			}

			keeper.EjectStaleTimestampNonces(&accountState, referenceTs)

			require.Equal(t, expectedAccountState, accountState)
		})
	}
}

func TestAttemptTimestampNonceUpdate(t *testing.T) {
	startTs := keeper.TimestampNonceSequenceCutoff
	t.Run("Will reject ts nonce <= maxEjectedNonce", func(t *testing.T) {
		tsNonce := startTs + 10

		var tsNonces []uint64
		for i := range 5 {
			tsNonces = append(tsNonces, startTs+uint64(i)+20)
		}

		accountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: startTs + 10,
			},
		}

		expectedAccountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: startTs + 10,
			},
		}

		updated := keeper.AttemptTimestampNonceUpdate(tsNonce, &accountState)

		require.False(t, updated)
		require.Equal(t, expectedAccountState, accountState)
	})

	t.Run("Will reject duplicate ts nonce", func(t *testing.T) {
		tsNonce := startTs + 20

		var tsNonces []uint64
		for i := range 5 {
			tsNonces = append(tsNonces, startTs+uint64(i)+20)
		}

		accountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: startTs + 10,
			},
		}

		expectedAccountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: startTs + 10,
			},
		}

		updated := keeper.AttemptTimestampNonceUpdate(tsNonce, &accountState)

		require.False(t, updated)
		require.Equal(t, expectedAccountState, accountState)
	})

	t.Run("Will update if ts nonces has capacity (given ts unique and ts > maxEjectedNonce)", func(t *testing.T) {
		tsNonce := startTs + 11

		var tsNonces []uint64
		for i := range 5 {
			tsNonces = append(tsNonces, startTs+uint64(i)+20)
		}

		accountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: tsNonces,
				MaxEjectedNonce: startTs + 10,
			},
		}

		expectedAccountState := types.AccountState{
			Address: constants.AliceAccAddress.String(),
			TimestampNonceDetails: types.TimestampNonceDetails{
				TimestampNonces: append(tsNonces, tsNonce),
				MaxEjectedNonce: startTs + 10,
			},
		}

		updated := keeper.AttemptTimestampNonceUpdate(tsNonce, &accountState)

		require.True(t, updated)
		require.Equal(t, expectedAccountState, accountState)
	})

	t.Run(
		"Will not update if ts nonce <= existing ts nonces (given ts unique and ts > maxEjectedNonce)",
		func(t *testing.T) {
			tsNonce := startTs + 19

			var tsNonces []uint64
			for i := range keeper.MaxTimestampNonceArrSize {
				tsNonces = append(tsNonces, startTs+uint64(i)+20)
			}

			accountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: tsNonces,
					MaxEjectedNonce: startTs, // ensure ejected less than ts nonce
				},
			}

			expectedAccountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: tsNonces,
					MaxEjectedNonce: startTs,
				},
			}

			updated := keeper.AttemptTimestampNonceUpdate(tsNonce, &accountState)

			require.False(t, updated)
			require.Equal(t, expectedAccountState, accountState)
		})

	t.Run(
		"Will update if ts nonce larger than at least one existing ts nonce (given ts unique and ts > maxEjectedNonce)",
		func(t *testing.T) {
			var tsNonces []uint64
			for i := range keeper.MaxTimestampNonceArrSize {
				tsNonces = append(tsNonces, startTs+uint64(i)+20)
			}

			tsNonce := tsNonces[len(tsNonces)-1] + 1

			accountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: tsNonces,
					MaxEjectedNonce: startTs,
				},
			}

			updatedTsNonces := make([]uint64, len(tsNonces))
			copy(updatedTsNonces, tsNonces)
			updatedTsNonces[0] = tsNonce

			expectedAccountState := types.AccountState{
				Address: constants.AliceAccAddress.String(),
				TimestampNonceDetails: types.TimestampNonceDetails{
					TimestampNonces: updatedTsNonces,
					MaxEjectedNonce: tsNonces[0],
				},
			}

			updated := keeper.AttemptTimestampNonceUpdate(tsNonce, &accountState)

			require.True(t, updated)
			require.Equal(t, expectedAccountState, accountState)
		})
}
