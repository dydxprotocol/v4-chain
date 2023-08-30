package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	feetiers "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		NumMessages: 1,
		DelayedMessages: []*DelayedMessage{
			// The first 120 days of fees are set to a promotional schedule. Those fee tier parameters are populated
			// as the default genesis state of x/feetiers. We use a delayed message here to update those parameters to
			// the standard schedule after 120 days of blocks have passed.
			{
				Id:          0,
				Msg:         getStandardPerpetualFeeParamsUpdateAsAny(),
				BlockHeight: BlockHeight120Days,
			},
		},
	}
}

// getStandardPerpetualFeeParamsUpdateAsAny returns a MsgUpdatePerpetualFeeParams message with the standard fee tier
// schedule as an Any type. This is used in the default genesis state to automate switching the fee tier schedule
// after the promotional period has passed.
func getStandardPerpetualFeeParamsUpdateAsAny() *types.Any {
	msg := &feetiers.MsgUpdatePerpetualFeeParams{
		Authority: authtypes.NewModuleAddress(ModuleName).String(),
		Params:    feetiers.StandardParams(),
	}
	any, _ := types.NewAnyWithValue(msg)
	return any
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
