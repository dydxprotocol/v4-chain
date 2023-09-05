package types

import (
	moderrors "cosmossdk.io/errors"
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		NumMessages:     0,
		DelayedMessages: []*DelayedMessage{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	ids := make(map[uint32]struct{}, len(gs.DelayedMessages))

	for i, msg := range gs.DelayedMessages {
		if err := msg.Validate(); err != nil {
			return moderrors.Wrap(
				ErrInvalidGenesisState,
				fmt.Sprintf("invalid delayed message at index %v with id %v: %v", i, msg.Id, err),
			)
		}

		if msg.Id >= gs.NumMessages {
			return moderrors.Wrap(
				ErrInvalidGenesisState,
				"delayed message id exceeds total number of messages",
			)
		}
		if _, ok := ids[msg.Id]; ok {
			return moderrors.Wrap(
				ErrInvalidGenesisState,
				"duplicate delayed message id",
			)
		}
		ids[msg.Id] = struct{}{}
	}

	return nil
}
