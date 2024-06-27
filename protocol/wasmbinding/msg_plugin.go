package wasmbinding

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bindings "github.com/dydxprotocol/v4-chain/protocol/wasmbinding/bindings"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// DispatchMsg executes on the contractMsg.
func CustomEncoder(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	var customMessage bindings.DydxCustomWasmMessage
	if err := json.Unmarshal(msg, &customMessage); err != nil {
		return []sdk.Msg{}, wasmvmtypes.InvalidRequest{Err: "Error parsing DydxCustomWasmMessage"}
	}
	switch {
	case customMessage.DepositToSubaccount != nil:
		return EncodeDepositToSubaccount(sender, customMessage.DepositToSubaccount)
	case customMessage.WithdrawFromSubaccount != nil:
		return EncodeWithdrawFromSubaccount(sender, customMessage.WithdrawFromSubaccount)
	case customMessage.PlaceOrder != nil:
		return EncodePlaceOrder(sender, customMessage.PlaceOrder)
	case customMessage.CancelOrder != nil:
		return EncodeCancelOrder(sender, customMessage.CancelOrder)
	default:
		return nil, wasmvmtypes.InvalidRequest{Err: "Unknown Dydx Wasm Message"}
	}
}

func EncodeDepositToSubaccount(sender sdk.AccAddress, depositToSubaccount *bindings.DepositToSubaccount) ([]sdk.Msg, error) {
	if depositToSubaccount == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid deposit to subaccount request: No deposit data provided"}
	}

	depositToSubaccountMsg := &sendingtypes.MsgDepositToSubaccount{
		Sender:    sender.String(),
		Recipient: depositToSubaccount.Recipient,
		AssetId:   depositToSubaccount.AssetId,
		Quantums:  depositToSubaccount.Quantums,
	}
	return []sdk.Msg{depositToSubaccountMsg}, nil
}

func EncodeWithdrawFromSubaccount(sender sdk.AccAddress, withdrawFromSubaccount *bindings.WithdrawFromSubaccount) ([]sdk.Msg, error) {
	if withdrawFromSubaccount == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid withdraw from subaccount request: No withdraw data provided"}
	}

	withdrawFromSubaccountMsg := &sendingtypes.MsgWithdrawFromSubaccount{
		Sender: types.SubaccountId{
			Owner:  sender.String(),
			Number: withdrawFromSubaccount.SubaccountNumber,
		},
		Recipient: withdrawFromSubaccount.Recipient,
		AssetId:   withdrawFromSubaccount.AssetId,
		Quantums:  withdrawFromSubaccount.Quantums,
	}
	return []sdk.Msg{withdrawFromSubaccountMsg}, nil
}

func EncodePlaceOrder(sender sdk.AccAddress, placeOrder *bindings.PlaceOrder) ([]sdk.Msg, error) {
	if placeOrder == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid place order request: No order data provided"}
	}

	placeOrderMsg := &clobtypes.MsgPlaceOrder{
		Order: clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: types.SubaccountId{
					Owner:  sender.String(),
					Number: placeOrder.SubaccountNumber,
				},
				ClientId:   placeOrder.ClientId,
				OrderFlags: placeOrder.OrderFLags,
				ClobPairId: placeOrder.ClobPairId,
			},
			Side:                            clobtypes.Order_Side(placeOrder.Side),
			Quantums:                        placeOrder.Quantums,
			Subticks:                        placeOrder.Subticks,
			GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: placeOrder.GoodTilBlockTime},
			ReduceOnly:                      placeOrder.ReduceOnly,
			ClientMetadata:                  placeOrder.ClientMetadata,
			ConditionType:                   clobtypes.Order_ConditionType(placeOrder.ConditionType),
			ConditionalOrderTriggerSubticks: placeOrder.ConditionalOrderTriggerSubticks,
		},
	}
	return []sdk.Msg{placeOrderMsg}, nil
}

func EncodeCancelOrder(sender sdk.AccAddress, cancelOrder *bindings.CancelOrder) ([]sdk.Msg, error) {
	if cancelOrder == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid cancel order request: No order data provided"}
	}

	cancelOrderMsg := &clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: types.SubaccountId{
				Owner:  sender.String(),
				Number: cancelOrder.SubaccountNumber,
			},
			ClientId:   cancelOrder.ClientId,
			OrderFlags: cancelOrder.OrderFLags,
			ClobPairId: cancelOrder.ClobPairId,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: cancelOrder.GoodTilBlockTime},
	}
	return []sdk.Msg{cancelOrderMsg}, nil
}
