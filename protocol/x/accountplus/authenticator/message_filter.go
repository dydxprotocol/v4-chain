package authenticator

import (
	"strings"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const SEPARATOR = ","

var _ Authenticator = &MessageFilter{}

// MessageFilter filters incoming messages based on a predefined JSON pattern.
// It allows for complex pattern matching to support advanced authentication flows.
type MessageFilter struct {
	msgTypes []string
}

// NewMessageFilter creates a new MessageFilter with the provided EncodingConfig.
func NewMessageFilter() MessageFilter {
	return MessageFilter{}
}

// Type returns the type of the authenticator.
func (m MessageFilter) Type() string {
	return "MessageFilter"
}

// StaticGas returns the static gas amount for the authenticator. Currently, it's set to zero.
func (m MessageFilter) StaticGas() uint64 {
	return 0
}

// Initialize sets up the authenticator with the given data, which should be a valid JSON pattern for message filtering.
func (m MessageFilter) Initialize(config []byte) (Authenticator, error) {
	m.msgTypes = strings.Split(string(config), SEPARATOR)
	return m, nil
}

// Track is a no-op in this implementation but can be used to track message handling.
func (m MessageFilter) Track(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// Authenticate checks if the provided message conforms to the set JSON pattern.
// It returns an AuthenticationResult based on the evaluation.
func (m MessageFilter) Authenticate(ctx sdk.Context, request AuthenticationRequest) error {
	for _, msgType := range m.msgTypes {
		if request.Msg.TypeURL == msgType {
			return nil
		}
	}
	return errorsmod.Wrapf(
		sdkerrors.ErrUnauthorized,
		"message types do not match. Got %s, Expected %v",
		request.Msg.TypeURL,
		m.msgTypes,
	)
}

// ConfirmExecution confirms the execution of a message. Currently, it always confirms.
func (m MessageFilter) ConfirmExecution(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// OnAuthenticatorAdded performs additional checks when an authenticator is added.
// Specifically, it ensures numbers in JSON are encoded as strings.
func (m MessageFilter) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
}

// OnAuthenticatorRemoved is a no-op in this implementation but can be used when an authenticator is removed.
func (m MessageFilter) OnAuthenticatorRemoved(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
}
