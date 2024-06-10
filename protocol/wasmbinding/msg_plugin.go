package wasmbinding

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bindings "github.com/dydxprotocol/v4-chain/protocol/wasmbinding/bindings"

	sendingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(sending *sendingkeeper.Keeper, clob *clobkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			sending: sending,
			clob:    clob,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	sending *sendingkeeper.Keeper
	clob    *clobkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.SendingMsg
		// print the custom message in string format
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "Error Unmarshalling Custom Message")
		}
		if contractMsg.CreateTransfer != nil {
			return m.createTransfer(ctx, contractAddr, contractMsg.CreateTransfer)
		}
		if contractMsg.DepositToSubaccount != nil {
			return m.depositToSubaccount(ctx, contractAddr, contractMsg.DepositToSubaccount)
		}
		if contractMsg.WithdrawFromSubaccount != nil {
			return m.withdrawFromSubaccount(ctx, contractAddr, contractMsg.WithdrawFromSubaccount)
		}
		if contractMsg.PlaceOrder != nil {
			return m.placeOrder(ctx, contractAddr, contractMsg.PlaceOrder)
		}
		if contractMsg.CancelOrder != nil {
			return m.cancelOrder(ctx, contractAddr, contractMsg.CancelOrder)
		}
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Unknown custom message"}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) createTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, createTransfer *sendingtypes.MsgCreateTransfer) ([]sdk.Event, [][]byte, error) {
	if createTransfer == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid create transfer request: No transfer data provided"}
	}
	err := m.sending.ProcessTransfer(ctx, createTransfer.Transfer)
	return nil, nil, err
}

func (m *CustomMessenger) depositToSubaccount(ctx sdk.Context, contractAddr sdk.AccAddress, depositToSubaccount *sendingtypes.MsgDepositToSubaccount) ([]sdk.Event, [][]byte, error) {
	if depositToSubaccount == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid deposit to subaccount request: No deposit data provided"}
	}

	err := m.sending.ProcessDepositToSubaccount(ctx, depositToSubaccount)
	return nil, nil, err
}

func (m *CustomMessenger) withdrawFromSubaccount(ctx sdk.Context, contractAddr sdk.AccAddress, withdrawFromSubaccount *sendingtypes.MsgWithdrawFromSubaccount) ([]sdk.Event, [][]byte, error) {
	if withdrawFromSubaccount == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid withdraw from subaccount request: No withdraw data provided"}
	}

	err := m.sending.ProcessWithdrawFromSubaccount(ctx, withdrawFromSubaccount)
	return nil, nil, err

}

func parseAddress(addr string) (sdk.AccAddress, error) {
	parsed, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}
	err = sdk.VerifyAddressFormat(parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (m *CustomMessenger) placeOrder(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	placeOrder *clobtypes.MsgPlaceOrder,
) ([]sdk.Event, [][]byte, error) {
	if placeOrder == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid place order request: No order data provided"}
	}
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		return nil, nil, nil
	}
	err := m.clob.HandleMsgPlaceOrder(ctx, placeOrder, false)
	return nil, nil, err
}

func (m *CustomMessenger) cancelOrder(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	cancelOrder *clobtypes.MsgCancelOrder,
) ([]sdk.Event, [][]byte, error) {
	if cancelOrder == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid cancel order request: No order data provided"}
	}
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		return nil, nil, nil
	}
	err := m.clob.HandleMsgCancelOrder(ctx, cancelOrder)
	return nil, nil, err
}
