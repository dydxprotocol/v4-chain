package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	// LeverageKeyPrefix is the prefix for leverage storage
	LeverageKeyPrefix = "leverage/"
)

// Leverage represents leverage data for a subaccount
type Leverage struct {
	SubaccountId      *satypes.SubaccountId `protobuf:"bytes,1,opt,name=subaccount_id,json=subaccountId,proto3" json:"subaccount_id,omitempty"`
	PerpetualLeverage map[uint32]uint32     `protobuf:"bytes,2,rep,name=perpetual_leverage,json=perpetualLeverage,proto3" json:"perpetual_leverage,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

// Reset implements the proto.Message interface
func (m *Leverage) Reset() { *m = Leverage{} }

// String implements the proto.Message interface
func (m *Leverage) String() string { return "Leverage" }

// ProtoMessage implements the proto.Message interface
func (*Leverage) ProtoMessage() {}
