package types

import (
	"fmt"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgExecuteBackingSwap{}
	_ sdk.Msg = &MsgCancelBackingSwap{}
	_ sdk.Msg = &MsgUpdateControls{}
	_ sdk.Msg = &MsgUpdateParticipants{}
)

func signer(address string) []sdk.AccAddress {
	parsed, _ := sdk.AccAddressFromBech32(address)
	return []sdk.AccAddress{parsed}
}

func validateSigner(address string) error {
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return fmt.Errorf("%w: invalid signer: %v", ErrUnauthorized, err)
	}
	return nil
}

func (msg *MsgExecuteBackingSwap) GetSigners() []sdk.AccAddress { return signer(msg.Controller) }

func (msg *MsgExecuteBackingSwap) ValidateBasic() error {
	if err := validateSigner(msg.Controller); err != nil {
		return err
	}
	if _, err := ParsePositiveAmount(msg.Amount); err != nil {
		return err
	}
	if msg.SourcePort != transfertypes.PortID {
		return fmt.Errorf("source port must be %s", transfertypes.PortID)
	}
	if msg.TimeoutTimestamp == 0 {
		return fmt.Errorf("timeout timestamp must be positive")
	}
	return nil
}

func (msg *MsgCancelBackingSwap) GetSigners() []sdk.AccAddress { return signer(msg.Controller) }

func (msg *MsgCancelBackingSwap) ValidateBasic() error {
	if err := validateSigner(msg.Controller); err != nil {
		return err
	}
	_, err := ParsePositiveAmount(msg.Amount)
	return err
}

func (msg *MsgUpdateControls) GetSigners() []sdk.AccAddress { return signer(msg.Authority) }

func (msg *MsgUpdateControls) ValidateBasic() error {
	if err := validateSigner(msg.Authority); err != nil {
		return err
	}
	return msg.Controls.Validate()
}

func (msg *MsgUpdateParticipants) GetSigners() []sdk.AccAddress { return signer(msg.Authority) }

func (msg *MsgUpdateParticipants) ValidateBasic() error {
	if err := validateSigner(msg.Authority); err != nil {
		return err
	}
	if len(msg.Participants) > MaxParticipants {
		return ErrInvalidParticipant
	}
	for _, participant := range msg.Participants {
		if participant.Funded != "" && participant.Funded != "0" ||
			participant.Released != "" && participant.Released != "0" {
			return fmt.Errorf("%w: accounting fields are read-only", ErrInvalidParticipant)
		}
		participant.Funded = "0"
		participant.Released = "0"
		if err := participant.Validate(); err != nil {
			return err
		}
	}
	return nil
}
