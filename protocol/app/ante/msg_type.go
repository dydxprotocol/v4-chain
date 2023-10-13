package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
)

// ValidateMsgTypeDecorator checks that the tx has the expected message types.
// Specifically, if the list of msgs in the tx contains an "app-injected message", the tx
// must only contain a single message.
// This decorator will not get exeuted on ReCheckTx since it does not depend on app state.
type ValidateMsgTypeDecorator struct{}

func NewValidateMsgTypeDecorator() ValidateMsgTypeDecorator {
	return ValidateMsgTypeDecorator{}
}

func (vbd ValidateMsgTypeDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	// no need to validate format on recheck tx, call next antehandler
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	numMsgs := len(msgs)
	if numMsgs == 0 { // invalid.
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "msgs cannot be empty")
	}

	for _, msg := range msgs {
		// 1. "App-injected message" check.
		if libante.IsAppInjectedMsg(msg) {
			// "App-injected message" must be the only msg in the tx.
			if numMsgs > 1 {
				return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "app-injected msg must be the only msg in a tx")
			}

			// "App-injected message" must only be included in DeliverTx.
			if !lib.IsDeliverTxMode(ctx) {
				return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "app-injected msg must only be included in DeliverTx")
			}
		}

		// 2. Internal-only message check.
		if libante.IsInternalMsg(msg) {
			return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "internal msg cannot be submitted externally")
		}

		// 3. Nested message check.
		if libante.IsNestedMsg(msg) {
			if err := libante.ValidateNestedMsg(msg); err != nil {
				return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		}

		// 4. Unsupported message check.
		if libante.IsUnsupportedMsg(msg) {
			return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unsupported msg")
		}
	}

	return next(ctx, tx, simulate)
}
