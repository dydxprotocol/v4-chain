package types

import (
	"fmt"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDelayMessage implements the UnpackInterfaces method for unpacking Msg, which is encoded as an Any type.
// Implementing this interface is necessary to decode the Msg, see https://docs.cosmos.network/v0.45/core/encoding.html
var _ codec.UnpackInterfacesMessage = &DelayedMessage{}

func (dm *DelayedMessage) Validate() error {
	if dm.Msg == nil {
		return ErrMsgIsNil
	}
	return nil
}

func (dm *DelayedMessage) UnpackInterfaces(unpacker codec.AnyUnpacker) error {
	var sdkMsg sdk.Msg
	// Unpack the Any into the sdk.Msg type. This should hydrate the cached value.
	return unpacker.UnpackAny(dm.Msg, &sdkMsg)
}

func (dm *DelayedMessage) GetMessage() (sdk.Msg, error) {
	if dm.Msg == nil {
		return nil, ErrMsgIsNil
	}
	cached := dm.Msg.GetCachedValue()
	if cached == nil {
		return nil, fmt.Errorf("any cached value is nil, delayed messages must be correctly packed any values")
	}
	casted, ok := cached.(sdk.Msg)
	if !ok {
		return nil, fmt.Errorf("cached value is not a sdk.Msg")
	}
	return casted, nil
}
