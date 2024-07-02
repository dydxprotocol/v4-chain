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

func EncodeDydxCustomWasmMessage(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	var customMessage bindings.DydxCustomWasmMessage
	if err := json.Unmarshal(msg, &customMessage); err != nil {
		return []sdk.Msg{}, wasmvmtypes.InvalidRequest{Err: "Error parsing DydxCustomWasmMessage"}
	}
	switch {
	case customMessage.DepositToSubaccountV1 != nil:
		return EncodeDepositToSubaccountV1(sender, customMessage.DepositToSubaccountV1)
	case customMessage.WithdrawFromSubaccountV1 != nil:
		return EncodeWithdrawFromSubaccountV1(sender, customMessage.WithdrawFromSubaccountV1)
	case customMessage.PlaceOrderV1 != nil:
		return EncodePlaceOrderV1(sender, customMessage.PlaceOrderV1)
	case customMessage.CancelOrderV1 != nil:
		return EncodeCancelOrderV1(sender, customMessage.CancelOrderV1)
	default:
		return nil, wasmvmtypes.InvalidRequest{Err: "Unknown Dydx Wasm Message"}
	}
}

func EncodeDepositToSubaccountV1(sender sdk.AccAddress, depositToSubaccount *bindings.DepositToSubaccountV1) ([]sdk.Msg, error) {
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

// This function is called from https://github.com/CosmWasm/wasmd/blob/main/x/wasm/keeper/handler_plugin_encoders.go#L96
// which enforces sender to be the contract address
func EncodeWithdrawFromSubaccountV1(sender sdk.AccAddress, withdrawFromSubaccount *bindings.WithdrawFromSubaccountV1) ([]sdk.Msg, error) {
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

func EncodePlaceOrderV1(sender sdk.AccAddress, placeOrder *bindings.PlaceOrderV1) ([]sdk.Msg, error) {
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

func EncodeCancelOrderV1(sender sdk.AccAddress, cancelOrder *bindings.CancelOrderV1) ([]sdk.Msg, error) {
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
