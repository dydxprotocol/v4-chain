package ante

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

const DYDX_MSG_PREFIX = "/" + constants.AppName
const SLINKY_MSG_PREFIX = "/slinky"

// IsNestedMsg returns true if the given msg is a nested msg.
func IsNestedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ------- CosmosSDK default modules
		// authz
		*authz.MsgExec,
		// gov
		*gov.MsgSubmitProposal:
		return true
	}
	return false
}

// IsDydxMsg returns true if the given msg is a dYdX custom msg.
func IsDydxMsg(msg sdk.Msg) bool {
	return strings.HasPrefix(sdk.MsgTypeURL(msg), DYDX_MSG_PREFIX)
}

// IsSlinkyMsg returns true if the given msg is a Slinky custom msg.
func IsSlinkyMsg(msg sdk.Msg) bool {
	return strings.HasPrefix(sdk.MsgTypeURL(msg), SLINKY_MSG_PREFIX)
}

// ValidateNestedMsg returns err if the given msg is an invalid nested msg.
func ValidateNestedMsg(msg sdk.Msg) error {
	if !IsNestedMsg(msg) {
		return fmt.Errorf("not a nested msg")
	}

	// Check that the inner msgs are valid.
	if err := validateInnerMsg(msg); err != nil {
		return err
	}

	return nil // is valid nested msg.
}

// validateInnerMsg returns err if the given inner msgs contain an invalid msg.
func validateInnerMsg(msg sdk.Msg) error {
	// Get inner msgs.
	innerMsgs, err := getInnerMsgs(msg)
	if err != nil {
		return err
	}

	for _, inner := range innerMsgs {
		// 1. unsupported msgs.
		if IsUnsupportedMsg(inner) {
			return fmt.Errorf("Invalid nested msg: unsupported msg type")
		}

		// 2. app-injected msgs.
		if IsAppInjectedMsg(inner) {
			return fmt.Errorf("Invalid nested msg: app-injected msg type")
		}

		// 3. double-nested msgs.
		if IsNestedMsg(inner) {
			return fmt.Errorf("Invalid nested msg: double-nested msg type")
		}

		// 4. Reject nested dydxprotocol messages in `MsgExec`.
		if _, ok := msg.(*authz.MsgExec); ok {
			metrics.IncrCountMetricWithLabels(
				metrics.Ante,
				metrics.MsgExec,
				metrics.GetLabelForStringValue(metrics.InnerMsg, sdk.MsgTypeURL(inner)),
			)
			if IsDydxMsg(inner) {
				return fmt.Errorf("Invalid nested msg for MsgExec: dydx msg type")
			}
			if IsSlinkyMsg(inner) {
				return fmt.Errorf("Invalid nested msg for MsgExec: Slinky msg type")
			}
		}

		// For "internal msgs", we allow them, because they are designed to be nested.
	}
	return nil
}

// getInnerMsgs returns the inner msgs of the given msg.
func getInnerMsgs(msg sdk.Msg) ([]sdk.Msg, error) {
	switch msg := msg.(type) {
	case *gov.MsgSubmitProposal:
		return msg.GetMsgs()
	case *authz.MsgExec:
		return msg.GetMessages()
	default:
		return nil, fmt.Errorf("unsupported msg type: %T", msg)
	}
}
