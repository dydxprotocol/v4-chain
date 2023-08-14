package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

func TestEpochInfoValidate(t *testing.T) {
	tests := map[string]struct {
		name                   string
		duration               uint32
		expectedErr            error
		currentEpoch           uint32
		currentEpochStartBlock uint32
	}{
		"validates successfully": {
			name:        "id",
			duration:    60,
			expectedErr: nil,
		},
		"Failure: invalid duration": {
			name:        "id",
			duration:    0,
			expectedErr: types.ErrDurationIsZero,
		},
		"Failure: empty id": {
			name:        "",
			duration:    60,
			expectedErr: types.ErrEmptyEpochInfoName,
		},
		"Failure: CurrentEpoch is zero, CurrentEpochStartBlock is non-zero": {
			name:                   "id",
			duration:               60,
			currentEpoch:           0,
			currentEpochStartBlock: 100,
			expectedErr:            types.ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple,
		},
		"Failure: CurrentEpoch is non-zero, CurrentEpochStartBlock is zero": {
			name:                   "id",
			duration:               60,
			currentEpoch:           10,
			currentEpochStartBlock: 0,
			expectedErr:            types.ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			epochInfo := &types.EpochInfo{
				Name:                   tc.name,
				Duration:               tc.duration,
				CurrentEpoch:           tc.currentEpoch,
				CurrentEpochStartBlock: tc.currentEpochStartBlock,
			}

			err := epochInfo.Validate()
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEpochInfoName(t *testing.T) {
	epochInfo := types.EpochInfo{
		Name: "id",
	}

	EpochInfoName := epochInfo.GetEpochInfoName()
	require.Equal(t, types.EpochInfoName("id"), EpochInfoName)
}
