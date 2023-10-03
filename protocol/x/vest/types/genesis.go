package types

import (
	time "time"

	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// TODO(CORE-530): in genesis.sh, overwrite start and end times dynamically for testnets.
		VestEntries: []VestEntry{
			{
				VesterAccount:   CommunityVesterAccountName,
				TreasuryAccount: CommunityTreasuryAccountName,
				Denom:           "dv4tnt",
				StartTime:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC),
				EndTime:         time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC),
			},
			{
				VesterAccount:   rewardstypes.VesterAccountName,
				TreasuryAccount: rewardstypes.TreasuryAccountName,
				Denom:           "dv4tnt",
				StartTime:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC),
				EndTime:         time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC),
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
