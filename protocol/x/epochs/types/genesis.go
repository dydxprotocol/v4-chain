package types

import (
	"fmt"
	"sort"
)

// DefaultGenesis returns the default epochs genesis state.
func DefaultGenesis() *GenesisState {
	epochInfoList := []EpochInfo{}

	for EpochInfoName, params := range GenesisEpochs {
		epochInfoList = append(
			epochInfoList,
			EpochInfo{
				Name:                   string(EpochInfoName),
				Duration:               params.Duration,
				NextTick:               params.NextTick,
				CurrentEpoch:           0,
				CurrentEpochStartBlock: 0,
				FastForwardNextTick:    true,
			},
		)
	}

	// Sort the list so the order is deterministic.
	sort.SliceStable(epochInfoList, func(i, j int) bool {
		return epochInfoList[i].Name < epochInfoList[j].Name
	})

	return &GenesisState{
		EpochInfoList: epochInfoList,
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in epochInfo
	epochInfoIndexMap := make(map[string]struct{})

	for _, epochInfo := range gs.EpochInfoList {
		index := string(epochInfo.Name)
		if _, ok := epochInfoIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for epochInfo")
		}
		epochInfoIndexMap[index] = struct{}{}

		if err := epochInfo.Validate(); err != nil {
			return fmt.Errorf("failed to validate epochInfo in genesis state: %w", err)
		}
	}

	return nil
}
