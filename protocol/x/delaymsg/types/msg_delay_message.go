package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDelayMessage implements the UnpackInterfaces method for unpacking Msg, which is encoded as an Any type.
// Implementing this interface is necessary to decode the Msg, see https://docs.cosmos.network/v0.45/core/encoding.html
var _ codec.UnpackInterfacesMessage = &MsgDelayMessage{}

// ValidateBasic performs basic validation on the message.
func (msg *MsgDelayMessage) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(
			ErrInvalidAuthority,
			fmt.Sprintf(
				"authority '%s' must be a valid bech32 address, but got error '%v'",
				msg.Authority,
				err.Error(),
			),
		)
	}

	// Perform basic checks that the encoded message was set.
	if msg.Msg == nil {
		return ErrMsgIsNil
	}

	return nil
}

func (msg *MsgDelayMessage) UnpackInterfaces(unpacker codec.AnyUnpacker) error {
	var sdkMsg sdk.Msg
	// Unpack the Any into the sdk.Msg type. This should hydrate the cached value.
	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

func (msg *MsgDelayMessage) GetMessage() (sdk.Msg, error) {
	if msg.Msg == nil {
		return nil, ErrMsgIsNil
	}
	cached := msg.Msg.GetCachedValue()
	if cached == nil {
		return nil, fmt.Errorf("any cached value is nil, delayed messages must be correctly packed any values")
	}
	casted, ok := cached.(sdk.Msg)
	if !ok {
		return nil, fmt.Errorf("cached value is not a sdk.Msg")
	}
	return casted, nil
}
