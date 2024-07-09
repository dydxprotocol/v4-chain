package bindings

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	subaccounttypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type DydxCustomWasmMessage struct {
	DepositToSubaccountV1    *DepositToSubaccountV1    `json:"deposit_to_subaccount_v1,omitempty"`
	WithdrawFromSubaccountV1 *WithdrawFromSubaccountV1 `json:"withdraw_from_subaccount_v1,omitempty"`
	PlaceOrderV1             *PlaceOrderV1             `json:"place_order_v1,omitempty"`
	CancelOrderV1            *CancelOrderV1            `json:"cancel_order_v1,omitempty"`
	BatchCancelV1            *BatchCancelV1            `json:"batch_cancel_v1,omitempty"`
}

type DepositToSubaccountV1 struct {
	Recipient subaccounttypes.SubaccountId `json:"recipient"`
	AssetId   uint32                       `json:"asset_id"`
	Quantums  uint64                       `json:"quantums"`
}

type WithdrawFromSubaccountV1 struct {
	SubaccountNumber uint32 `json:"subaccount_number"`
	Recipient        string `json:"recipient"`
	AssetId          uint32 `json:"asset_id"`
	Quantums         uint64 `json:"quantums"`
}

type PlaceOrderV1 struct {
	SubaccountNumber                uint32 `json:"subaccount_number"`
	ClientId                        uint32 `json:"client_id"`
	OrderFLags                      uint32 `json:"order_flags"`
	ClobPairId                      uint32 `json:"clob_pair_id"`
	Side                            int32  `json:"side"`
	Quantums                        uint64 `json:"quantums"`
	Subticks                        uint64 `json:"subticks"`
	GoodTilBlockTime                uint32 `json:"good_til_block_time"`
	ReduceOnly                      bool   `json:"reduce_only"`
	ClientMetadata                  uint32 `json:"client_metadata"`
	ConditionType                   int32  `json:"condition_type"`
	ConditionalOrderTriggerSubticks uint64 `json:"conditional_order_trigger_subticks"`
}

type CancelOrderV1 struct {
	SubaccountNumber uint32 `json:"subaccount_number"`
	ClientId         uint32 `json:"client_id"`
	OrderFLags       uint32 `json:"order_flags"`
	ClobPairId       uint32 `json:"clob_pair_id"`
	GoodTilBlockTime uint32 `json:"good_til_block_time"`
}

type BatchCancelV1 struct {
	SubaccountNumber uint32                 `json:"subaccount_number"`
	ShortTermCancels []clobtypes.OrderBatch `json:"short_term_cancels"`
	GoodTilBlock     uint32                 `json:"good_til_block"`
}
