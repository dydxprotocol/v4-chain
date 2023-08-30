package process

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

var (
	msgAddPremiumVotesType = reflect.TypeOf(types.MsgAddPremiumVotes{})
)

// AddPremiumVotesTx represents `MsgAddPremiumVotes` tx that can be validated.
type AddPremiumVotesTx struct {
	msg *types.MsgAddPremiumVotes
}

// DecodeAddPremiumVotesTx returns a new `AddPremiumVotesTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the num of msgs in the tx matches expectations
//   - checks the msg is of expected type
//
// If error occurs during any of the checks, returns error.
func DecodeAddPremiumVotesTx(decoder sdk.TxDecoder, txBytes []byte) (*AddPremiumVotesTx, error) {
	// Decode.
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, getDecodingError(msgAddPremiumVotesType, err)
	}

	// Check msg length.
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, getUnexpectedNumMsgsError(msgAddPremiumVotesType, 1, len(msgs))
	}

	// Check msg type.
	addPremiumVotes, ok := msgs[0].(*types.MsgAddPremiumVotes)
	if !ok {
		return nil, getUnexpectedMsgTypeError(msgAddPremiumVotesType, msgs[0])
	}

	return &AddPremiumVotesTx{msg: addPremiumVotes}, nil
}

// Validate returns an error if the underlying msg fails `ValidateBasic`.
func (afst *AddPremiumVotesTx) Validate() error {
	if err := afst.msg.ValidateBasic(); err != nil {
		return getValidateBasicError(afst.msg, err)
	}
	return nil
}

// GetMsg returns the underlying `MsgAddPremiumVotes`.
func (afst *AddPremiumVotesTx) GetMsg() sdk.Msg {
	return afst.msg
}
