package process

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var (
	msgUpdateMarketPricesType = reflect.TypeOf(types.MsgUpdateMarketPrices{})
)

// UpdateMarketPricesTx represents `MsgUpdateMarketPrices` tx that can be validated.
type UpdateMarketPricesTx struct {
	ctx          sdk.Context
	pricesKeeper ProcessPricesKeeper
	msg          *types.MsgUpdateMarketPrices
}

// DecodeAddPremiumVotesTx returns a new `UpdateMarketPricesTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the num of msgs in the tx matches expectations
//   - checks the msg is of expected type
//
// If error occurs during any of the checks, returns error.
func DecodeUpdateMarketPricesTx(
	ctx sdk.Context,
	pricesKeeper ProcessPricesKeeper,
	decoder sdk.TxDecoder,
	txBytes []byte,
) (*UpdateMarketPricesTx, error) {
	// Decode.
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, getDecodingError(msgUpdateMarketPricesType, err)
	}

	// Check msg length.
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, getUnexpectedNumMsgsError(msgUpdateMarketPricesType, 1, len(msgs))
	}

	// Check msg type.
	updateMarketPrices, ok := msgs[0].(*types.MsgUpdateMarketPrices)
	if !ok {
		return nil, getUnexpectedMsgTypeError(msgUpdateMarketPricesType, msgs[0])
	}

	return &UpdateMarketPricesTx{
		ctx:          ctx,
		pricesKeeper: pricesKeeper,
		msg:          updateMarketPrices,
	}, nil
}

// Validate returns an error if:
// - the underlying msg fails `ValidateBasic`
// - the underlying msg values are not "valid" according to the index price.
func (umpt *UpdateMarketPricesTx) Validate() error {
	if err := umpt.msg.ValidateBasic(); err != nil {
		return getValidateBasicError(umpt.msg, err)
	}

	if err := umpt.pricesKeeper.PerformStatefulPriceUpdateValidation(umpt.ctx, umpt.msg, false); err != nil {
		return err
	}

	return nil
}

// GetMsg returns the underlying `MsgUpdateMarketPrices`.
func (umpt *UpdateMarketPricesTx) GetMsg() sdk.Msg {
	return umpt.msg
}
