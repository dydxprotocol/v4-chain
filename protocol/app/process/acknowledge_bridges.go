package process

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	gometrics "github.com/hashicorp/go-metrics"
)

var (
	msgAcknowledgeBridgesType = reflect.TypeOf(types.MsgAcknowledgeBridges{})
)

// AcknowledgeBridgesTx represents `MsgAcknowledgeBridge`s tx that can be validated.
type AcknowledgeBridgesTx struct {
	ctx          sdk.Context
	bridgeKeeper ProcessBridgeKeeper
	msg          *types.MsgAcknowledgeBridges
}

// DecodeAcknowledgeBridgesTx returns a new `AcknowledgeBridgesTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the msg is of expected type
//
// If error occurs during any of the checks, returns error.
func DecodeAcknowledgeBridgesTx(
	ctx sdk.Context,
	bridgeKeeper ProcessBridgeKeeper,
	decoder sdk.TxDecoder,
	txBytes []byte,
) (*AcknowledgeBridgesTx, error) {
	// Decode.
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, getDecodingError(msgAcknowledgeBridgesType, err)
	}

	// Check msg length.
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, getUnexpectedNumMsgsError(msgAcknowledgeBridgesType, 1, len(msgs))
	}

	// Check msg type.
	acknowledgeBridges, ok := msgs[0].(*types.MsgAcknowledgeBridges)
	if !ok {
		return nil, getUnexpectedMsgTypeError(msgAcknowledgeBridgesType, msgs[0])
	}

	return &AcknowledgeBridgesTx{
		ctx:          ctx,
		bridgeKeeper: bridgeKeeper,
		msg:          acknowledgeBridges,
	}, nil
}

// Validate returns an error if:
// - msg fails `ValidateBasic`.
// - bridge events are non empty and bridging is disabled.
// - first bridge event ID is not the one to be next acknowledged.
// - last bridge event ID has not been recognized.
// - a bridge event's content is not the same as in server state.
func (abt *AcknowledgeBridgesTx) Validate() error {
	// `ValidateBasic` validates that bridge event IDs are consecutive.
	if err := abt.msg.ValidateBasic(); err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{
				ModuleName,
				metrics.AcknowledgeBridgesTx,
				metrics.Validate,
				metrics.Error,
			},
			1,
			[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Error, metrics.ValidateBasic)},
		)
		return getValidateBasicError(abt.msg, err)
	}

	if len(abt.msg.Events) == 0 {
		// If there is no bridge event, return nil.
		return nil
	} else if abt.bridgeKeeper.GetSafetyParams(abt.ctx).IsDisabled {
		// If there is any bridge event when bridging is disabled, return error.
		return types.ErrBridgingDisabled
	}

	// Validate that first bridge event ID is the one to be next acknowledged.
	acknowledgedEventInfo := abt.bridgeKeeper.GetAcknowledgedEventInfo(abt.ctx)
	if acknowledgedEventInfo.NextId != abt.msg.Events[0].Id {
		telemetry.IncrCounterWithLabels(
			[]string{
				ModuleName,
				metrics.AcknowledgeBridgesTx,
				metrics.Validate,
				metrics.Error,
			},
			1,
			[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Error, types.ErrBridgeIdNotNextToAcknowledge.Error())},
		)
		return types.ErrBridgeIdNotNextToAcknowledge
	}

	// Validate that last bridge event ID has been recognized.
	recognizedEventInfo := abt.bridgeKeeper.GetRecognizedEventInfo(abt.ctx)
	if recognizedEventInfo.NextId <= abt.msg.Events[len(abt.msg.Events)-1].Id {
		telemetry.IncrCounterWithLabels(
			[]string{
				ModuleName,
				metrics.AcknowledgeBridgesTx,
				metrics.Validate,
				metrics.Error,
			},
			1,
			[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Error, types.ErrBridgeIdNotRecognized.Error())},
		)
		return types.ErrBridgeIdNotRecognized
	}

	// Validate that bridge events' content is the same as in server state.
	for _, event := range abt.msg.Events {
		eventInState, found := abt.bridgeKeeper.GetBridgeEventFromServer(abt.ctx, event.Id)
		if !found {
			return types.ErrBridgeEventNotFound
		}
		if !eventInState.Equal(event) {
			return types.ErrBridgeEventContentMismatch
		}
	}

	return nil
}

// GetMsg returns the underlying `MsgAcknowledgeBridges`.
func (abt *AcknowledgeBridgesTx) GetMsg() sdk.Msg {
	return abt.msg
}
