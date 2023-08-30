package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SingleMsgTx represents a tx with a single msg.
type SingleMsgTx interface {
	// validate checks if the underlying msg is valid or not.
	// Returns error if invalid. Otherwise, returns nil.
	Validate() error

	// getMsg returns the underlying msg in the tx.
	GetMsg() sdk.Msg
}

// MultiMsgsTx represents a tx with multiple msgs.
type MultiMsgsTx interface {
	// Validate checks if the underlying msgs are valid or not.
	// Returns error if one of the msgs is invalid. Otherwise, returns nil.
	Validate() error

	// GetMsgs returns the underlying msgs in the tx.
	GetMsgs() []sdk.Msg
}
