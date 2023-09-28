package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// IsNestedMsg returns true if the given msg is a nested msg.
func IsNestedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ------- CosmosSDK default modules
		// gov
		*gov.MsgSubmitProposal:
		return true
	}
	return false
}

// ValidateNestedMsg returns err if the given msg is an invalid nested msg.
func ValidateNestedMsg(msg sdk.Msg) error {
	if !IsNestedMsg(msg) {
		return fmt.Errorf("not a nested msg")
	}

	// Get inner msgs.
	innerMsgs, err := getInnerMsgs(msg)
	if err != nil {
		return err
	}

	// Check that the inner msgs are valid.
	if err := validateInnerMsg(innerMsgs); err != nil {
		return err
	}

	return nil // is valid nested msg.
}

// getInnerMsgs returns the inner msgs of the given msg.
func getInnerMsgs(msg sdk.Msg) ([]sdk.Msg, error) {
	switch msg := msg.(type) {
	case
		*gov.MsgSubmitProposal:
		return msg.GetMsgs()
	default:
		return nil, fmt.Errorf("unsupported msg type: %T", msg)
	}
}

// validateInnerMsg returns err if the given inner msgs contain an invalid msg.
func validateInnerMsg(innerMsgs []sdk.Msg) error {
	for _, msg := range innerMsgs {
		// 1. unsupported msgs.
		if IsUnsupportedMsg(msg) {
			return fmt.Errorf("Invalid nested msg: unsupported msg type")
		}

		// 2. app-injected msgs.
		if IsAppInjectedMsg(msg) {
			return fmt.Errorf("Invalid nested msg: app-injected msg type")
		}

		// 3. double-nested msgs.
		if IsNestedMsg(msg) {
			return fmt.Errorf("Invalid nested msg: double-nested msg type")
		}

		// For "internal msgs", we allow them, because they are designed to be nested.
	}
	return nil
}
