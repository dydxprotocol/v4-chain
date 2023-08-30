package types

import (
	"fmt"
	"sort"
)

func createGenesis(epochs []EpochInfo) *GenesisState {
	genesisState := &GenesisState{}
	genesisState.EpochInfoList = append(genesisState.EpochInfoList, epochs...)

	// Sort the list so the order is deterministic.
	sort.SliceStable(genesisState.EpochInfoList, func(i, j int) bool {
		return genesisState.EpochInfoList[i].Name < genesisState.EpochInfoList[j].Name
	})

	return genesisState
}

// DefaultGenesis returns the default epochs genesis state.
func DefaultGenesis() *GenesisState {
	return createGenesis(GenesisEpochs)
}

// DefaultGenesisWithEpochs returns the default genesis state with input epochs added or overwritten.
func DefaultGenesisWithEpochs(epochs ...EpochInfo) *GenesisState {
	newEpochs := append([]EpochInfo(nil), GenesisEpochs...)
	for _, epoch := range epochs {
		for i := range newEpochs {
			if newEpochs[i].Name == epoch.Name {
				newEpochs[i] = epoch
			}
		}
	}
	return createGenesis(newEpochs)
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
