package types

type EpochInfoName string

const (
	FundingSampleEpochInfoName EpochInfoName = "funding-sample"
	FundingTickEpochInfoName   EpochInfoName = "funding-tick"

	FundingTickEpochDuration   uint32 = 3600
	FundingSampleEpochDuration uint32 = 60
)

type GenesisEpochParam struct {
	Duration uint32
	NextTick uint32
}

var GenesisEpochs = map[EpochInfoName]GenesisEpochParam{
	// Ticks every hour on the hour.
	FundingTickEpochInfoName: {
		Duration: FundingTickEpochDuration,
		NextTick: 0,
	},
	// Ticks every minute at 30-seconds-past-the-minute.
	FundingSampleEpochInfoName: {
		Duration: FundingSampleEpochDuration,
		NextTick: 30,
	},
}
