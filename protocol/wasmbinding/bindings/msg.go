package bindings

type SendingMsg struct {
	CreateTransfer      *CreateTransfer      `json:"create_transfer,omitempty"`
	DepositToSubaccount *DepositToSubaccount `json:"deposit_to_subaccount,omitempty"`
}

type SubaccountId struct {
	Owner  string `json:"owner"`
	Number uint32 `json:"number"`
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
