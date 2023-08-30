package types

type EpochInfoName string

const (
	FundingSampleEpochInfoName EpochInfoName = "funding-sample"
	FundingTickEpochInfoName   EpochInfoName = "funding-tick"
	StatsEpochInfoName         EpochInfoName = "stats-epoch"

	FundingTickEpochDuration   uint32 = 3600
	FundingSampleEpochDuration uint32 = 60
	StatsEpochDuration         uint32 = 3600
)

var GenesisEpochs = []EpochInfo{
	// Ticks every hour on the hour.
	{
		Name:                   string(FundingTickEpochInfoName),
		Duration:               FundingTickEpochDuration,
		NextTick:               0,
		CurrentEpoch:           0,
		CurrentEpochStartBlock: 0,
		FastForwardNextTick:    true,
	},
	// Ticks every minute at 30-seconds-past-the-minute.
	{
		Name:                   string(FundingSampleEpochInfoName),
		Duration:               FundingSampleEpochDuration,
		NextTick:               30,
		CurrentEpoch:           0,
		CurrentEpochStartBlock: 0,
		FastForwardNextTick:    true,
	},
	{
		Name:                   string(StatsEpochInfoName),
		Duration:               StatsEpochDuration,
		NextTick:               0,
		CurrentEpoch:           0,
		CurrentEpochStartBlock: 0,
		FastForwardNextTick:    true,
	},
}
