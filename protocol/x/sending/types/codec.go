package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateTransfer{}, "sending/CreateTransfer", nil)
	cdc.RegisterConcrete(&MsgDepositToSubaccount{}, "sending/DepositToSubaccount", nil)
	cdc.RegisterConcrete(&MsgWithdrawFromSubaccount{}, "sending/WithdrawFromSubaccount", nil)
	cdc.RegisterConcrete(&MsgSendFromModuleToAccount{}, "sending/SendFromModuleToAccount", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateTransfer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDepositToSubaccount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawFromSubaccount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendFromModuleToAccount{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
