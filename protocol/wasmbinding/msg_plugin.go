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

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(sending *sendingkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			sending: sending,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	sending *sendingkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.SendingMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "sending msg")
		}
		if contractMsg.CreateTransfer != nil {
			return m.createTransfer(ctx, contractAddr, contractMsg.CreateTransfer)
		}
		if contractMsg.DepositToSubaccount != nil {
			return m.depositToSubaccount(ctx, contractAddr, contractMsg.DepositToSubaccount)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) createTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, createTransfer *bindings.CreateTransfer) ([]sdk.Event, [][]byte, error) {
	if createTransfer == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create transfer null transfer"}
	}

	senderAddress, err := parseAddress(createTransfer.Transfer.Sender.Owner)
	if err != nil {
		return nil, nil, err
	}

	senderNumber := createTransfer.Transfer.Sender.Number

	rcptAddress, err := parseAddress(createTransfer.Transfer.Recipient.Owner)
	if err != nil {
		return nil, nil, err
	}

	rcptNumber := createTransfer.Transfer.Recipient.Number

	pendingTransfer := sendingtypes.Transfer{
		Sender: satypes.SubaccountId{
			Owner:  senderAddress.String(),
			Number: senderNumber,
		},
		Recipient: satypes.SubaccountId{
			Owner:  rcptAddress.String(),
			Number: rcptNumber,
		},
		AssetId: createTransfer.Transfer.AssetId,
		Amount:  createTransfer.Transfer.Amount,
	}

	err = m.sending.ProcessTransfer(ctx, &pendingTransfer)

	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func (m *CustomMessenger) depositToSubaccount(ctx sdk.Context, contractAddr sdk.AccAddress, depositToSubaccount *bindings.DepositToSubaccount) ([]sdk.Event, [][]byte, error) {
	if depositToSubaccount == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "deposit to subaccount null deposit"}
	}

	senderAddress, err := parseAddress(depositToSubaccount.Sender)
	if err != nil {
		return nil, nil, err
	}

	rcptAddress, err := parseAddress(depositToSubaccount.Recipient.Owner)
	if err != nil {
		return nil, nil, err
	}

	rcptNumber := depositToSubaccount.Recipient.Number

	deposit := sendingtypes.MsgDepositToSubaccount{
		Sender: senderAddress.String(),
		Recipient: satypes.SubaccountId{
			Owner:  rcptAddress.String(),
			Number: rcptNumber,
		},
		AssetId:  depositToSubaccount.AssetId,
		Quantums: depositToSubaccount.Quantums,
	}

	err = m.sending.ProcessDepositToSubaccount(ctx, &deposit)

	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
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
