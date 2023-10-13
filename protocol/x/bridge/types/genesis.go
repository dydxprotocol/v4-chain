package types

import (
	"time"
)

// DefaultGenesis returns the default bridge genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		EventParams: EventParams{
			Denom:      "bridge-token",
			EthChainId: 11155111,
			EthAddress: "0xEf01c3A30eB57c91c40C52E996d29c202ae72193",
		},
		ProposeParams: ProposeParams{
			MaxBridgesPerBlock:           10,
			ProposeDelayDuration:         60 * time.Second,
			SkipRatePpm:                  800_000, // 80%
			SkipIfBlockDelayedByDuration: 5 * time.Second,
		},
		SafetyParams: SafetyParams{
			IsDisabled:  false,
			DelayBlocks: 86_400, // Seconds in a day
		},
		AcknowledgedEventInfo: BridgeEventInfo{
			NextId:         0,
			EthBlockHeight: 0,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.EventParams.Validate(); err != nil {
		return err
	}
	if err := gs.ProposeParams.Validate(); err != nil {
		return err
	}
	if err := gs.SafetyParams.Validate(); err != nil {
		return err
	}
	if err := gs.AcknowledgedEventInfo.Validate(); err != nil {
		return err
	}

	return nil
}
