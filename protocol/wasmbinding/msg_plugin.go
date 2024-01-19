package wasmbinding

import (
	"encoding/json"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	bindings "github.com/dydxprotocol/v4-chain/protocol/wasmbinding/bindings"

	sendingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(
	sending *sendingkeeper.Keeper,
	clob *clobkeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
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
func (m *CustomMessenger) DispatchMsg(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	contractIBCPortID string,
	msg wasmvmtypes.CosmosMsg,
) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.DydxMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "dydx msg")
		}
		if contractMsg.CreateTransfer != nil {
			return m.createTransfer(ctx, contractAddr, contractMsg.CreateTransfer)
		}
		if contractMsg.DepositToSubaccount != nil {
			return m.depositToSubaccount(ctx, contractAddr, contractMsg.DepositToSubaccount)
		}
		if contractMsg.PlaceOrder != nil {
			return m.placeOrder(ctx, contractAddr, contractMsg.PlaceOrder)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) createTransfer(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	createTransfer *bindings.CreateTransfer,
) ([]sdk.Event, [][]byte, error) {
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

func (m *CustomMessenger) depositToSubaccount(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	depositToSubaccount *bindings.DepositToSubaccount,
) ([]sdk.Event, [][]byte, error) {
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

func (m *CustomMessenger) placeOrder(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	placeOrder *bindings.PlaceOrder,
) ([]sdk.Event, [][]byte, error) {
	if placeOrder == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create transfer null transfer"}
	}

	address, err := parseAddress(placeOrder.Order.OrderId.SubaccountId.Owner)
	if err != nil {
		return nil, nil, err
	}

	number := placeOrder.Order.OrderId.SubaccountId.Number

	order := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: satypes.SubaccountId{
				Owner:  address.String(),
				Number: number,
			},
			ClientId:   placeOrder.Order.OrderId.ClientId,
			OrderFlags: placeOrder.Order.OrderId.OrderFlags,
			ClobPairId: placeOrder.Order.OrderId.ClobPairId,
		},
		Side:                            clobtypes.Order_Side(placeOrder.Order.Side),
		Quantums:                        placeOrder.Order.Quantums,
		Subticks:                        placeOrder.Order.Subticks,
		TimeInForce:                     clobtypes.Order_TimeInForce(placeOrder.Order.TimeInForce),
		ReduceOnly:                      placeOrder.Order.ReduceOnly,
		ClientMetadata:                  placeOrder.Order.ClientMetadata,
		ConditionType:                   clobtypes.Order_ConditionType(placeOrder.Order.ConditionType),
		ConditionalOrderTriggerSubticks: placeOrder.Order.ConditionalOrderTriggerSubticks,
	}

	if placeOrder.Order.OrderId.OrderFlags == clobtypes.OrderIdFlags_ShortTerm {
		order.GoodTilOneof = &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: placeOrder.Order.GoodTilBlock,
		}
		// Only process short term orders in CheckTx because short term order placements
		// are never on chain.
		if ctx.IsCheckTx() {
			fmt.Println("Placing short term order")
			_, _, err = m.clob.PlaceShortTermOrder(ctx, &clobtypes.MsgPlaceOrder{Order: order})
		}
	} else {
		order.GoodTilOneof = &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: lib.MustConvertIntegerToUint32(time.Now().Unix() + 60*60),
		}

		if ctx.IsCheckTx() || ctx.IsReCheckTx() {
			// We don't process stateful orders in CheckTx, so do nothing.
			return nil, nil, nil
		}

		fmt.Println("Placing stateful order")
		processProposerMatchesEvents := m.clob.GetProcessProposerMatchesEvents(ctx)

		if err := m.clob.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: order}); err != nil {
			return nil, nil, err
		}

		if order.IsConditionalOrder() {
			m.clob.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewConditionalOrderPlacementEvent(
						order,
					),
				),
			)
			processProposerMatchesEvents.PlacedConditionalOrderIds = append(
				processProposerMatchesEvents.PlacedConditionalOrderIds,
				order.OrderId,
			)
		} else {
			m.clob.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewLongTermOrderPlacementEvent(
						order,
					),
				),
			)
			processProposerMatchesEvents.PlacedLongTermOrderIds = append(
				processProposerMatchesEvents.PlacedLongTermOrderIds,
				order.OrderId,
			)
		}
		m.clob.MustSetProcessProposerMatchesEvents(
			ctx,
			processProposerMatchesEvents,
		)
	}

	return nil, nil, err
}
