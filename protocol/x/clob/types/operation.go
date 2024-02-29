package types

import (
	"bytes"
	fmt "fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

// decodeOperationRawShortTermOrderPlacementBytes performs stateless validation
// on a short term order placement given its underlying raw tx bytes. It also
// runs the transaction through an antehandler. The antehandler is needed to
// do signature validation. Returns an Operation if successful.
func decodeOperationRawShortTermOrderPlacementBytes(
	ctx sdk.Context,
	bytes []byte,
	decoder sdk.TxDecoder,
	anteHandler sdk.AnteHandler,
) (*InternalOperation, error) {
	tx, err := decoder(bytes)
	if err != nil {
		return nil, err
	}

	if _, err := anteHandler(ctx, tx, false); err != nil {
		return nil, err
	}

	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, fmt.Errorf("expected 1 msg, got %d", len(msgs))
	}

	msg, ok := msgs[0].(*MsgPlaceOrder)
	if !ok {
		return nil, fmt.Errorf("expected MsgPlaceOrder, got %T", msgs[0])
	}

	return &InternalOperation{
		Operation: &InternalOperation_ShortTermOrderPlacement{
			ShortTermOrderPlacement: msg,
		},
	}, nil
}

// GetInternalOperationTextString returns the text string representation of this operation.
// TODO(DEC-1772): Add method for encoding operation protos as JSON to make debugging easier.
func (o *InternalOperation) GetInternalOperationTextString() string {
	return proto.MarshalTextString(o)
}

// GetOperationsQueueString returns a string representation of the provided operations.
func GetInternalOperationsQueueTextString(operations []InternalOperation) string {
	var buf bytes.Buffer
	for _, op := range operations {
		// Note that MarshalText only throws errors if the writer throws errors which it never does
		proto.MarshalText(&buf, &op) //nolint:errcheck
		buf.Write([]byte("\n"))
	}
	return buf.String()
}
