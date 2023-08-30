package types

import (
	"fmt"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDelayMessage implements the UnpackInterfaces method for unpacking Msg.
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
	cached := dm.Msg.GetCachedValue()
	if cached == nil {
		return nil, fmt.Errorf("any cached value is nil, delayed messages must be correctly packed any values")
	}
	return cached.(sdk.Msg), nil
}
