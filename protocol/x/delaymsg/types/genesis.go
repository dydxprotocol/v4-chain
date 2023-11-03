package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DelayedMessages:      []*DelayedMessage{},
		NextDelayedMessageId: 0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	ids := make(map[uint32]struct{}, len(gs.DelayedMessages))

	for i, msg := range gs.DelayedMessages {
		if err := msg.Validate(); err != nil {
			return errorsmod.Wrap(
				ErrInvalidGenesisState,
				fmt.Sprintf("invalid delayed message at index %v with id %v: %v", i, msg.Id, err),
			)
		}

		if msg.Id >= gs.NextDelayedMessageId {
			return errorsmod.Wrap(
				ErrInvalidGenesisState,
				"delayed message id cannot be greater than or equal to next id",
			)
		}
		if _, ok := ids[msg.Id]; ok {
			return errorsmod.Wrap(
				ErrInvalidGenesisState,
				"duplicate delayed message id",
			)
		}
		ids[msg.Id] = struct{}{}
	}

	return nil
}
