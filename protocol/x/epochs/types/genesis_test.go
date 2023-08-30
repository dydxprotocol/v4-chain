package types_test

import (
	"errors"
	"sort"
	"testing"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState      *types.GenesisState
		expectedError error
	}{
		"valid: default": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid": {
			genState: &types.GenesisState{
				EpochInfoList: []types.EpochInfo{
					{
						Name:     "0",
						Duration: keepertest.TestEpochDuration,
					},
					{
						Name:     "1",
						Duration: keepertest.TestEpochDuration,
					},
				},
			},
			expectedError: nil,
		},
		"invalid: duplicated epochInfo": {
			genState: &types.GenesisState{
				EpochInfoList: []types.EpochInfo{
					{
						Name:     "0",
						Duration: keepertest.TestEpochDuration,
					},
					{
						Name:     "0",
						Duration: keepertest.TestEpochDuration,
					},
				},
			},
			expectedError: errors.New("duplicated index for epochInfo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}

func TestDefaultGenesis_DefaultValue(t *testing.T) {
	defaultGenesis := types.DefaultGenesis()
	require.Equal(t,
		3,
		len(defaultGenesis.EpochInfoList),
	)

	for _, epochInfo := range defaultGenesis.EpochInfoList {
		require.Equal(t,
			true,
			epochInfo.FastForwardNextTick,
		)

		switch epochInfo.Name {
		case "funding-tick":
			require.Equal(t,
				uint32(3600),
				epochInfo.Duration,
			)
			require.Equal(t,
				false,
				epochInfo.IsInitialized,
			)
			require.Equal(t,
				uint32(0),
				epochInfo.NextTick,
			)
		case "funding-sample":
			require.Equal(t,
				uint32(60),
				epochInfo.Duration,
			)
			require.Equal(t,
				false,
				epochInfo.IsInitialized,
			)
			require.Equal(t,
				uint32(30),
				epochInfo.NextTick,
			)
		case "stats-epoch":
			require.Equal(t,
				uint32(3600),
				epochInfo.Duration,
			)
			require.Equal(t,
				false,
				epochInfo.IsInitialized,
			)
			require.Equal(t,
				uint32(0),
				epochInfo.NextTick,
			)
		default:
			t.Errorf("Unexepcted genesis epoch name:%s", epochInfo.Name)
		}
	}
}

func TestDefaultGenesis_Determinism(t *testing.T) {
	var listFromFirstIteration []types.EpochInfo

	for i := 0; i < 100; i++ {
		defaultEpochList := types.DefaultGenesis().EpochInfoList
		if i == 0 {
			// Save the list from the first iteration.
			listFromFirstIteration = append(listFromFirstIteration, defaultEpochList...)
		} else {
			require.Equal(t,
				listFromFirstIteration,
				defaultEpochList,
			)
		}

		// Assert the default genesis list is sorted in lexicographical order.
		require.Equal(t,
			true,
			sort.SliceIsSorted(defaultEpochList, func(i, j int) bool {
				return defaultEpochList[i].Name < defaultEpochList[j].Name
			}),
		)
	}
}
