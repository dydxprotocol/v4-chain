package types

import (
	time "time"

	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		VestEntries: []VestEntry{
			{
				VesterAccount:   rewardstypes.VesterAccountName,
				TreasuryAccount: rewardstypes.TreasuryAccountName,
				Denom:           "testnet_reward_token",
				StartTime:       time.Date(2023, 8, 2, 0, 0, 0, 0, time.UTC).In(time.UTC),
				EndTime:         time.Date(2023, 8, 23, 0, 0, 0, 0, time.UTC).In(time.UTC),
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, vs := range gs.VestEntries {
		if err := vs.Validate(); err != nil {
			return err
		}
	}
	return nil
}
