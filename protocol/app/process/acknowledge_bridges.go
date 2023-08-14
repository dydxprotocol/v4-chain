package process

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
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
// - msg fails `bridgeKeeper.CanAcknowledgeBridges`.
func (abt *AcknowledgeBridgesTx) Validate() error {
	// `ValidateBasic` validates that bridge event IDs are consecutive.
	if err := abt.msg.ValidateBasic(); err != nil {
		return getValidateBasicError(abt.msg, err)
	}

	// If there is no bridge event, return nil.
	if len(abt.msg.Events) == 0 {
		return nil
	}

	// Validate that first bridge event ID is the one to be next acknowledged.
	acknowledgedEventInfo := abt.bridgeKeeper.GetAcknowledgedEventInfo(abt.ctx)
	if acknowledgedEventInfo.NextId != abt.msg.Events[0].Id {
		return types.ErrBridgeIdNotNextToAcknowledge
	}

	// Validate that last bridge event ID has been recognized.
	recognizedEventInfo := abt.bridgeKeeper.GetRecognizedEventInfo(abt.ctx)
	if recognizedEventInfo.NextId <= abt.msg.Events[len(abt.msg.Events)-1].Id {
		return types.ErrBridgeIdNotRecognized
	}

	return nil
}

// GetMsg returns the underlying `MsgAcknowledgeBridges`.
func (abt *AcknowledgeBridgesTx) GetMsg() sdk.Msg {
	return abt.msg
}
