package ante

import (
	errorsmod "cosmossdk.io/errors"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

const (
	// WasmExecMaxGasLimit is the maximum gas limit for a cosmwasm execution transaction.
	WasmExecMaxGasLimit = 2_000_000
)

var _ sdktypes.AnteDecorator = (*WasmExecDecorator)(nil)

type WasmExecDecorator struct {
	ratelimitKeeper ratelimittypes.RatelimitKeeper
}

func NewWasmExecDecorator(
	ratelimitKeeper ratelimittypes.RatelimitKeeper,
) WasmExecDecorator {
	return WasmExecDecorator{
		ratelimitKeeper: ratelimitKeeper,
	}
}

// IsSingleWasmExecTx returns `true` if the supplied `tx` consist of a single Cosmwasm `MsgExecuteContract` message.
// If `msgs` consist of multiple `MsgExecuteContract` messages, or a mix of it and other messages, an error is returned.
func IsSingleWasmExecTx(tx sdktypes.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	var hasMessage = false

	for _, msg := range msgs {
		switch msg.(type) {
		case *wasmtypes.MsgExecuteContract:
			hasMessage = true
		}

		if hasMessage {
			break
		}
	}

	if !hasMessage {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing `MsgExecuteContract` may not contain more than one message",
		)
	}

	return true, nil
}

// CheckGasLimit checks the specified gas limit of the transaction is within the allowed range.
func (c WasmExecDecorator) CheckGasLimit(ctx sdktypes.Context, tx sdktypes.Tx) error {
	feeTx, ok := tx.(sdktypes.FeeTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	gasLimit := feeTx.GetGas()
	if gasLimit > WasmExecMaxGasLimit {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidGasLimit,
			"CosmWasm execution specified gas limit (%v) exceeds `WasmExecMaxGasLimit` (%v)",
			gasLimit,
			WasmExecMaxGasLimit,
		)
	}

	return nil
}

// RateLimitWasmExec rate limits the execution of Cosmwasm contracts.
func (c WasmExecDecorator) RateLimitWasmExec(ctx sdktypes.Context, tx sdktypes.Tx) error {
	// Do not rate-limit in DeliverTx mode.
	if lib.IsDeliverTxMode(ctx) {
		return nil
	}

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *wasmtypes.MsgExecuteContract:
			if err := c.ratelimitKeeper.RatelimitWasmExute(ctx, msg.Sender); err != nil {
				return err
			}

		default:
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidType,
				"unexpected message type %T in WasmExecDecorator",
				msg,
			)
		}
	}

	return nil
}

func (c WasmExecDecorator) AnteHandle(
	ctx sdktypes.Context,
	tx sdktypes.Tx, simulate bool,
	next sdktypes.AnteHandler,
) (newCtx sdktypes.Context, err error) {
	isSingleWasmEx, err := IsSingleWasmExecTx(tx)
	if err != nil {
		return ctx, err
	}
	// No-op if transaction isn't a single Cosmwasm `MsgExecuteContract` message.
	if !isSingleWasmEx {
		return next(ctx, tx, simulate)
	}

	// Rate limit Cosmwasm execution.
	if err := c.RateLimitWasmExec(ctx, tx); err != nil {
		return ctx, err
	}

	// Check specified gas limit is within the allowed range.
	if err := c.CheckGasLimit(ctx, tx); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}
