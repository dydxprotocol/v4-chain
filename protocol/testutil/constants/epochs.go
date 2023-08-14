package constants

import (
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
)

var (
	Duration_OneYear = time.Hour * 24 * 365
)

// GenerateEpochGenesisStateWithoutFunding is a testutil function to generate
// genesis state for the epochs module where no new epochs will trigger for
// `funding-tick` and `funding-sample` within the next 365 days, meaning that
// there will be no funding payment for unit and end-to-end tests that use this
// genesis state.
func GenerateEpochGenesisStateWithoutFunding() types.GenesisState {
	now := time.Now()
	return *types.DefaultGenesisWithEpochs(
		types.EpochInfo{
			Name:                   string(types.FundingSampleEpochInfoName),
			NextTick:               uint32(now.Add(Duration_OneYear).Unix()),
			Duration:               uint32(Duration_OneYear.Seconds()),
			CurrentEpoch:           0,
			CurrentEpochStartBlock: 0,
			FastForwardNextTick:    false,
		},
		types.EpochInfo{
			Name:                   string(types.FundingTickEpochInfoName),
			NextTick:               uint32(now.Add(Duration_OneYear).Unix()),
			Duration:               uint32(Duration_OneYear.Seconds()),
			CurrentEpoch:           0,
			CurrentEpochStartBlock: 0,
			FastForwardNextTick:    false,
		},
	)
}
