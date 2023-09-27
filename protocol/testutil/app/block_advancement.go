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
	ShortTermOrdersAndOperations []interface{}     // should hold Order and OperationRaw types. Slice to allow for ordering.
	StatefulOrders               []clobtypes.Order // should hold stateful orders to include in DeliverTx after ProposedOperationsTx
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
	advanceToBlockOptions := AdvanceToBlockOptions{
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
	}

	deliverTxsOverride := b.BlockAdvancement.getDeliverTxs(ctx, tApp.App)
	if len(deliverTxsOverride) > 0 {
		advanceToBlockOptions.DeliverTxsOverride = deliverTxsOverride
	}

	return tApp.AdvanceToBlock(blockHeight, advanceToBlockOptions)
}

// getDeliverTxs returns a slice of tx bytes to be executed in DeliverTx.
func (b BlockAdvancement) getDeliverTxs(ctx sdktypes.Context, app *app.App) [][]byte {
	deliverTxs := make([][]byte, 0)

	// operations come first in block
	if len(b.ShortTermOrdersAndOperations) > 0 {
		deliverTxs = append(deliverTxs, b.getProposedOperationsTxBytes(ctx, app))
	}

	// stateful order placements come after all app-injected messages
	deliverTxs = append(deliverTxs, b.getStatefulMsgPlaceOrderTxBytes(ctx, app)...)

	return deliverTxs
}

// getStatefulMsgPlaceOrderTxBytes iterates over StatefulOrders and returns a slice of tx bytes corresponding to the
// signed set of MsgPlaceOrder txs.
func (b BlockAdvancement) getStatefulMsgPlaceOrderTxBytes(ctx sdktypes.Context, app *app.App) [][]byte {
	txs := make([][]byte, len(b.StatefulOrders))

	for i, order := range b.StatefulOrders {
		if !order.IsStatefulOrder() {
			panic("Order should be stateful")
		}
		requestTxs := MustMakeCheckTxsWithClobMsg(
			ctx,
			app,
			*clobtypes.NewMsgPlaceOrder(order),
		)
		txs[i] = requestTxs[0].Tx
	}
	return txs
}

// getProposedOperationsTxBytes iterates through the ShortTermOrdersAndOperations slice, signing every order and appending a
// short term order placement operation to the operations queue. Other elements in the list should be of type
// OperationRaw and will be appended to the operations queue as is. Transaction bytes for tx containing
// the MsgProposedOperations msg are returned.
func (b BlockAdvancement) getProposedOperationsTxBytes(ctx sdktypes.Context, app *app.App) []byte {
	operationsQueue := make([]clobtypes.OperationRaw, len(b.ShortTermOrdersAndOperations))
	for i, orderOrOperation := range b.ShortTermOrdersAndOperations {
		switch castedValue := orderOrOperation.(type) {
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

	msgProposedOperations := &clobtypes.MsgProposedOperations{
		OperationsQueue: operationsQueue,
	}

	return testtx.MustGetTxBytes(msgProposedOperations)
}
