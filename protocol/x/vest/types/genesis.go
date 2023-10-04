package types

import (
	time "time"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

var (
	DefaultVestingStartTime = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC)
	DefaultVestingEndTime   = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC)
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// TODO(CORE-530): in genesis.sh, overwrite start and end times dynamically for testnets.
		VestEntries: []VestEntry{
			{
				VesterAccount:   CommunityVesterAccountName,
				TreasuryAccount: CommunityTreasuryAccountName,
				Denom:           lib.DefaultBaseDenom,
				StartTime:       DefaultVestingStartTime,
				EndTime:         DefaultVestingEndTime,
			},
			{
				VesterAccount:   rewardstypes.VesterAccountName,
				TreasuryAccount: rewardstypes.TreasuryAccountName,
				Denom:           lib.DefaultBaseDenom,
				StartTime:       DefaultVestingStartTime,
				EndTime:         DefaultVestingEndTime,
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
