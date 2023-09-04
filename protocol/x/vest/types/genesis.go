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
				VesterAccount:   rewardstypes.VesterAccountName,
				TreasuryAccount: rewardstypes.TreasuryAccountName,
				Denom:           "dv4tnt",
				StartTime:       time.Date(2023, 9, 13, 0, 0, 0, 0, time.UTC).In(time.UTC),
				EndTime:         time.Date(2023, 10, 13, 0, 0, 0, 0, time.UTC).In(time.UTC),
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
