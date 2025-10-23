package ante

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cometbftlog "github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	timeoutHeightLogKey = "TimeoutHeight"
)

// ClobDecorator is an AnteDecorator which is responsible for:
//   - adding short term order placements and cancelations to the in-memory orderbook (`CheckTx` only).
//   - adding stateful order placements and cancelations to state (`CheckTx` and `RecheckTx` only).
//
// This AnteDecorator also enforces that any Transaction which contains a `MsgPlaceOrder`,
// or a `MsgCancelOrder`, must consist only of a single message.
//
// This AnteDecorator is a no-op if:
//   - No messages in the transaction are `MsgPlaceOrder` or `MsgCancelOrder`.
//   - This AnteDecorator is called during `DeliverTx`.
//
// This AnteDecorator returns an error if:
//   - The transaction contains multiple messages, and one of them is a `MsgPlaceOrder`
//     or `MsgCancelOrder` message.
//   - The underlying `PlaceStatefulOrder`, `PlaceShortTermOrder`, `CancelStatefulOrder`, or `CancelShortTermOrder`
//     methods on the keeper return errors.
type ClobDecorator struct {
	clobKeeper    types.ClobKeeper
	sendingKeeper sendingtypes.SendingKeeper
}

func NewClobDecorator(
	clobKeeper types.ClobKeeper,
	sendingKeeper sendingtypes.SendingKeeper,
) ClobDecorator {
	return ClobDecorator{
		clobKeeper,
		sendingKeeper,
	}
}

func (cd ClobDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	// No need to process during `DeliverTx` or simulation, call next `AnteHandler`.
	if lib.IsDeliverTxMode(ctx) || simulate {
		return next(ctx, tx, simulate)
	}

	// Check if the transaction is a valid clob tx
	if err := ValidateMsgsInClobTx(tx); err != nil {
		return ctx, err
	}

	// Disable order placement and cancelation processing if the clob keeper is not initialized.
	if !cd.clobKeeper.IsInMemStructuresInitialized() {
		return ctx, errorsmod.Wrap(
			types.ErrClobNotInitialized,
			"clob keeper is not initialized. Please wait for the next block.",
		)
	}

	msgs := tx.GetMsgs()

	var err error
	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *types.MsgCancelOrder:
			if msg.OrderId.IsStatefulOrder() {
				err = cd.clobKeeper.CancelStatefulOrder(ctx, msg)
			} else {
				// No need to process short term order cancelations on `ReCheckTx`.
				if ctx.IsReCheckTx() {
					return next(ctx, tx, simulate)
				}

				// Note that `msg.ValidateBasic` is called before the AnteHandlers.
				// This guarantees that `MsgCancelOrder` has undergone stateless validation.
				err = cd.clobKeeper.CancelShortTermOrder(ctx, msg)
			}

			log.DebugLog(
				ctx, "Received new order cancellation",
				log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
				log.Error, err,
			)

		case *types.MsgPlaceOrder:
			if msg.Order.OrderId.IsStatefulOrder() {
				err = cd.clobKeeper.PlaceStatefulOrder(ctx, msg, false)

				log.DebugLog(
					ctx, "Received new stateful order",
					log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
					log.OrderHash, cometbftlog.NewLazySprintf("%X", msg.Order.GetOrderHash()),
					log.Error, err,
				)
			} else {
				// No need to process short term orders on `ReCheckTx`.
				if ctx.IsReCheckTx() {
					return next(ctx, tx, simulate)
				}

				// HOTFIX: Reject any short-term place orders in a transaction with a non-zero timeout height < good til block
				if timeoutHeight := GetTimeoutHeight(tx); timeoutHeight > 0 &&
					timeoutHeight < uint64(msg.Order.GetGoodTilBlock()) && ctx.IsCheckTx() {
					log.InfoLog(
						ctx,
						"Rejected short-term place order with non-zero timeout height < goodTilBlock",
						timeoutHeightLogKey,
						timeoutHeight,
					)
					return ctx, errorsmod.Wrap(
						sdkerrors.ErrInvalidRequest,
						"timeout height (if non-zero) may not be less than `goodTilBlock` for a short-term place order",
					)
				}

				var orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums
				var status types.OrderStatus
				// Note that `msg.ValidateBasic` is called before all AnteHandlers.
				// This guarantees that `MsgPlaceOrder` has undergone stateless validation.
				orderSizeOptimisticallyFilledFromMatchingQuantums, status, err = cd.clobKeeper.PlaceShortTermOrder(
					ctx,
					msg,
				)

				log.DebugLog(
					ctx,
					"Received new short term order",
					log.Tx,
					cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
					log.OrderHash,
					cometbftlog.NewLazySprintf("%X", msg.Order.GetOrderHash()),
					log.OrderStatus,
					status,
					log.OrderSizeOptimisticallyFilledFromMatchingQuantums,
					orderSizeOptimisticallyFilledFromMatchingQuantums,
					log.Error,
					err,
				)
			}
		case *types.MsgBatchCancel:
			// MsgBatchCancel currently only processes short-term cancels right now.
			// No need to process short term orders on `ReCheckTx`.
			if ctx.IsReCheckTx() {
				return next(ctx, tx, simulate)
			}

			success, failures, err := cd.clobKeeper.BatchCancelShortTermOrder(
				ctx,
				msg,
			)
			// If there are no successful cancellations and no validation errors,
			// return an error indicating no cancels have succeeded.
			if len(success) == 0 && err == nil {
				err = errorsmod.Wrapf(
					types.ErrBatchCancelFailed,
					"No successful cancellations. Failures: %+v",
					failures,
				)
			}

			log.DebugLog(
				ctx,
				"Received new batch cancellation",
				log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
				log.Error, err,
			)

		case *sendingtypes.MsgCreateTransfer:
			// We allow a single transfer msg to be batched with a PlaceOrder msg.
			// This is primarily used to transfer collateral to an isolated subaccount for an isolated position order.
			if err := cd.sendingKeeper.ProcessTransfer(ctx, msg.Transfer); err != nil {
				log.DebugLog(
					ctx,
					"Failed to process transfer msg in clob ante handler",
					log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
					log.Error, err,
				)
				return ctx, err
			}
		case *types.MsgUpdateLeverage:
			// Process UpdateLeverage message - delegate to subaccounts keeper
			// Convert from LeverageEntry slice to map
			perpetualLeverageMap, err := types.ValidateAndConstructPerpetualLeverageMap(ctx, msg, cd.clobKeeper)
			if err != nil {
				return ctx, err
			}

			// Delegate to subaccounts keeper for leverage storage and validation
			if err := cd.clobKeeper.GetSubaccountsKeeper().UpdateLeverage(
				ctx,
				msg.SubaccountId,
				perpetualLeverageMap,
			); err != nil {
				log.DebugLog(
					ctx,
					"Failed to update leverage in ante handler",
					log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
					log.Error, err,
				)
				return ctx, err
			}

			log.DebugLog(
				ctx,
				"Received new leverage update",
				log.Tx, cometbftlog.NewLazySprintf("%X", tmhash.Sum(ctx.TxBytes())),
				"subaccount", msg.SubaccountId.String(),
			)
		}
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// IsShortTermClobMsgTx returns `true` if the supplied `tx` consist of a single clob message
// (`MsgPlaceOrder` or `MsgCancelOrder` or `MsgBatchCancel`) which references a Short-Term Order.
// If `msgs` consist of multiple clob messages, or a mix of on-chain and clob messages, an error is returned.
func IsShortTermClobMsgTx(ctx sdk.Context, tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()

	var isShortTermOrder = false

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *types.MsgCancelOrder:
			{
				if msg.OrderId.IsShortTermOrder() {
					isShortTermOrder = true
				}
			}
		case *types.MsgPlaceOrder:
			{
				if msg.Order.OrderId.IsShortTermOrder() {
					isShortTermOrder = true
				}
			}
		case *types.MsgBatchCancel:
			{
				// MsgBatchCancel processes only short term orders for now.
				isShortTermOrder = true
			}
		}

		if isShortTermOrder {
			break
		}
	}

	if !isShortTermOrder {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing short term MsgCancelOrder or MsgPlaceOrder may not contain more than one message",
		)
	}

	return true, nil
}

// HasClobMsg returns `true` if the transaction has at least one clob msg
func HasClobMsg(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()

	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgCancelOrder:
			return true
		case *types.MsgPlaceOrder:
			return true
		case *types.MsgBatchCancel:
			return true
		case *types.MsgUpdateLeverage:
			return true
		}
	}
	return false
}

// ValidateMsgsInClobTx checks if the transaction contains a valid set of clob msgs
// This function assumes that the input tx has at least one clob msg
// A transaction with a clob msg must adhere to the below conditions
//   - If the tx contains a short term order msg, the tx can only have one msg
//   - If the tx contains a stateful order msg, it can only contain other stateful order msgs
//     or a single transfer msg
func ValidateMsgsInClobTx(tx sdk.Tx) error {
	msgs := tx.GetMsgs()

	var hasShortTermOrder = false
	var numTransferMsgs = 0
	// Non CLOB msgs other than a single transfer msg are not allowed in CLOB msg transactions
	// because there is no gas fee charged for CLOB transactions
	var hasDisallowedMsg = false

	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *types.MsgCancelOrder:
			if msg.OrderId.IsShortTermOrder() {
				hasShortTermOrder = true
			}
		case *types.MsgPlaceOrder:
			if msg.Order.OrderId.IsShortTermOrder() {
				hasShortTermOrder = true
			}
		case *types.MsgBatchCancel:
			// MsgBatchCancel processes only short term orders for now.
			hasShortTermOrder = true
		case *types.MsgUpdateLeverage:
			// UpdateLeverage messages are allowed in CLOB transactions
		case *sendingtypes.MsgCreateTransfer:
			numTransferMsgs += 1
		default:
			hasDisallowedMsg = true
		}
	}

	if hasShortTermOrder && len(msgs) > 1 {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing short term order may not contain more than one message",
		)
	}

	// We only expect a single transfer msg to be batched with a PlaceOrder msg. This is primarily used
	// to transfer collateral to an isolated subaccount for an isolated position order
	if numTransferMsgs > 1 {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing stateful orders can only be accompanied by 1 transfer msg",
		)
	}

	if hasDisallowedMsg {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing stateful orders cannot be accompanied by non transfer msgs",
		)
	}

	return nil
}

// GetTimeoutHeight returns the timeout height of a transaction. If the transaction does not have
// a timeout height, return 0.
func GetTimeoutHeight(tx sdk.Tx) uint64 {
	timeoutTx, ok := tx.(sdk.TxWithTimeoutHeight)
	if !ok {
		return 0
	}

	timeoutHeight := timeoutTx.GetTimeoutHeight()
	return timeoutHeight
}
