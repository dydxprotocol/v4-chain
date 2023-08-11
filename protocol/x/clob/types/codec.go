package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgProposedOperations{}, "clob/ProposedOperations", nil)
	cdc.RegisterConcrete(&MsgPlaceOrder{}, "clob/PlaceOrder", nil)
	cdc.RegisterConcrete(&MsgCancelOrder{}, "clob/CancelOrder", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposedOperations{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPlaceOrder{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCancelOrder{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
