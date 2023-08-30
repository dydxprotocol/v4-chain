package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
			return sdkerrors.Wrap(
				ErrInvalidGenesisState,
				fmt.Sprintf("invalid delayed message at index %v with id %v: %v", i, msg.Id, err),
			)
		}

		if msg.Id >= gs.NumMessages {
			return sdkerrors.Wrap(
				ErrInvalidGenesisState,
				"delayed message id exceeds total number of messages",
			)
		}
		if _, ok := ids[msg.Id]; ok {
			return sdkerrors.Wrap(
				ErrInvalidGenesisState,
				"duplicate delayed message id",
			)
		}
		ids[msg.Id] = struct{}{}
	}

	return nil
}
