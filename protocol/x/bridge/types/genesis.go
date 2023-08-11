package types

import (
	"github.com/dydxprotocol/v4/lib"
)

// DefaultGenesis returns the default bridge genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		EventParams: EventParams{
			Denom:      "bridge-token",
			EthChainId: 0,
			EthAddress: "0x0000000000000000000000000000000000000000",
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
		NextAcknowledgedEventId: 0,
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

	return nil
}
