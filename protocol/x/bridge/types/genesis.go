package types

import (
	"github.com/dydxprotocol/v4/lib"
)

// DefaultGenesis returns the default bridge genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		EventParams: EventParams{
			Denom:      "bridge-token",
			EthChainId: 11155111,
			EthAddress: "0x40ad69F5d9f7F9EA2Fc5C2009C7335F10593C935",
		},
		ProposeParams: ProposeParams{
			MaxBridgesPerBlock:           10,
			ProposeDelayDuration:         lib.MustParseDuration("60s"),
			SkipRatePpm:                  800_000, // 80%
			SkipIfBlockDelayedByDuration: lib.MustParseDuration("5s"),
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
