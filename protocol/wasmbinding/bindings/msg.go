package bindings

type DydxMsg struct {
	CreateTransfer      *CreateTransfer      `json:"create_transfer,omitempty"`
	DepositToSubaccount *DepositToSubaccount `json:"deposit_to_subaccount,omitempty"`
	PlaceOrder          *PlaceOrder          `json:"place_order,omitempty"`
}

type SubaccountId struct {
	Owner  string `json:"owner"`
	Number uint32 `json:"number"`
}

type OrderId struct {
	SubaccountId *SubaccountId `json:"subaccount_id"`
	ClientId     uint32        `json:"client_id"`
	OrderFlags   uint32        `json:"order_flags"`
	ClobPairId   uint32        `json:"clob_pair_id"`
}

type Order struct {
	OrderId                         *OrderId `json:"order_id"`
	Side                            uint32   `json:"side,omitempty"`
	Quantums                        uint64   `json:"quantums,omitempty"`
	Subticks                        uint64   `json:"subticks,omitempty"`
	GoodTilBlock                    uint32   `json:"good_til_block,omitempty"`
	GoodTilBlockTime                uint32   `json:"good_til_block_time,omitempty"`
	TimeInForce                     uint32   `json:"time_in_force,omitempty"`
	ReduceOnly                      bool     `json:"reduce_only,omitempty"`
	ClientMetadata                  uint32   `json:"client_metadata,omitempty"`
	ConditionType                   uint32   `json:"condition_type,omitempty"`
	ConditionalOrderTriggerSubticks uint64   `json:"conditional_order_trigger_subticks,omitempty"`
}

type CreateTransfer struct {
	Transfer *Transfer `json:"transfer,omitempty"`
}

type Transfer struct {
	Sender    *SubaccountId `json:"sender,omitempty"`
	Recipient *SubaccountId `json:"recipient,omitempty"`
	AssetId   uint32        `json:"asset_id,omitempty"`
	Amount    uint64        `json:"amount,omitempty"`
}

type DepositToSubaccount struct {
	Sender    string        `json:"sender"`
	Recipient *SubaccountId `json:"recipient,omitempty"`
	AssetId   uint32        `json:"asset_id,omitempty"`
	Quantums  uint64        `json:"quantums,omitempty"`
}

type PlaceOrder struct {
	Order *Order `json:"order"`
}
