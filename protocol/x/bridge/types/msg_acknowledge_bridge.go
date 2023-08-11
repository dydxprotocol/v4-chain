package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgAcknowledgeBridge) GetSigners() []sdk.AccAddress {
	// Return empty slice because app-injected msg is not expected to be signed.
	return []sdk.AccAddress{}
}

func (msg *MsgAcknowledgeBridge) ValidateBasic() error {
	return nil
}
