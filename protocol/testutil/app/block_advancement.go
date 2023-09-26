package app

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// BlockAdvancement holds orders and matches to be placed in a block. Using this struct and building
// the ops queue with the getOperationsQueue helper function allows us to build the operations queue
// without going through CheckTx and, therefore, not affect the local memclob state. This also allows us to propose
// an invalid set of operations that an honest validator would not generate.
type BlockAdvancement struct {
	OrdersAndOperations []interface{} // should hold Order and OperationRaw types.
}

type BlockAdvancementWithError struct {
	BlockAdvancement       BlockAdvancement
	ExpectedDeliverTxError string
}

// AdvanceToBlock advances the test app to the given block height using the operations queue
// generated from the specified BlockAdvancement. It catches errors in DeliverTx and verifies that
// the error matches the expected error.
func (b BlockAdvancementWithError) AdvanceToBlock(
	ctx sdktypes.Context,
	blockHeight uint32,
	tApp *TestApp,
	t *testing.T,
) sdktypes.Context {
	msgProposedOperations := &clobtypes.MsgProposedOperations{
		OperationsQueue: b.BlockAdvancement.getOperationsQueue(ctx, tApp.App),
	}
	return tApp.AdvanceToBlock(blockHeight, AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{testtx.MustGetTxBytes(msgProposedOperations)},
		ValidateDeliverTxs: func(
			ctx sdktypes.Context,
			request abcitypes.RequestDeliverTx,
			response abcitypes.ResponseDeliverTx,
		) (haltchain bool) {
			if b.ExpectedDeliverTxError != "" {
				require.True(t, response.IsErr())
				require.Contains(t, response.Log, b.ExpectedDeliverTxError)
			} else {
				require.True(t, response.IsOK())
			}
			return false
		},
	})
}

// getOperationsQueue iterates through the ordersAndOperations slice, signing every order and appending a
// short term order placement operation to the operations queue. Other elements in the list should be of type
// OperationRaw and will be appended to the operations queue as is.
func (b BlockAdvancement) getOperationsQueue(ctx sdktypes.Context, app *app.App) []clobtypes.OperationRaw {
	operationsQueue := make([]clobtypes.OperationRaw, len(b.OrdersAndOperations))
	for i, orderOrMatch := range b.OrdersAndOperations {
		switch castedValue := orderOrMatch.(type) {
		case clobtypes.Order:
			order := castedValue
			requestTxs := MustMakeCheckTxsWithClobMsg(
				ctx,
				app,
				*clobtypes.NewMsgPlaceOrder(order),
			)
			operationsQueue[i] = clobtypes.OperationRaw{
				Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
					ShortTermOrderPlacement: requestTxs[0].Tx,
				},
			}
		case clobtypes.OperationRaw:
			operationsQueue[i] = castedValue
		default:
			panic("invalid type")
		}
	}

	return operationsQueue
}