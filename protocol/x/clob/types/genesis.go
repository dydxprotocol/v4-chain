package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BlockRateLimitConfig:  BlockRateLimitConfiguration{},
		ClobPairs:             []ClobPair{},
		EquityTierLimitConfig: EquityTierLimitConfiguration{},
		LiquidationsConfig:    LiquidationsConfig_Default,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated id in clobPair
	clobPairIdMap := make(map[uint32]struct{})

	for _, clobPair := range gs.ClobPairs {
		if _, ok := clobPairIdMap[clobPair.Id]; ok {
			return fmt.Errorf("duplicated id for clobPair")
		}
		clobPairIdMap[clobPair.Id] = struct{}{}
	}

	if err := gs.BlockRateLimitConfig.Validate(); err != nil {
		return err
	}

	if err := gs.EquityTierLimitConfig.Validate(); err != nil {
		return err
	}

	if err := gs.LiquidationsConfig.Validate(); err != nil {
		return err
	}

	return nil
}
