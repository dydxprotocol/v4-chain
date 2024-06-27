package bindings

import (
	subaccounttypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type DydxCustomWasmMessage struct {
	DepositToSubaccount    *DepositToSubaccount    `json:"deposit_to_subaccount,omitempty"`
	WithdrawFromSubaccount *WithdrawFromSubaccount `json:"withdraw_from_subaccount,omitempty"`
	PlaceOrder             *PlaceOrder             `json:"place_order,omitempty"`
	CancelOrder            *CancelOrder            `json:"cancel_order,omitempty"`
}

type DepositToSubaccount struct {
	Recipient subaccounttypes.SubaccountId `json:"recipient"`
	AssetId   uint32                       `json:"asset_id"`
	Quantums  uint64                       `json:"quantums"`
}

type WithdrawFromSubaccount struct {
	SubaccountNumber uint32 `json:"subaccount_number"`
	Recipient        string `json:"recipient"`
	AssetId          uint32 `json:"asset_id"`
	Quantums         uint64 `json:"quantums"`
}

type PlaceOrder struct {
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

type CancelOrder struct {
	SubaccountNumber uint32 `json:"subaccount_number"`
	ClientId         uint32 `json:"client_id"`
	OrderFLags       uint32 `json:"order_flags"`
	ClobPairId       uint32 `json:"clob_pair_id"`
	GoodTilBlockTime uint32 `json:"good_til_block_time"`
}
