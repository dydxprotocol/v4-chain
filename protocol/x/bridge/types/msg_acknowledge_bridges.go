package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgAcknowledgeBridges) GetSigners() []sdk.AccAddress {
	// Return empty slice because app-injected msg is not expected to be signed.
	return []sdk.AccAddress{}
}

func (msg *MsgAcknowledgeBridges) ValidateBasic() error {
	// Validates that bridge event IDs are consecutive.
	for i, event := range msg.Events {
		if i > 0 && msg.Events[i-1].Id != event.Id-1 {
			return ErrBridgeIdsNotConsecutive
		}
	}
	return nil
}
