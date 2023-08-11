package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAcknowledgeBridge{}, "bridge/AcknowledgeBridge", nil)
	cdc.RegisterConcrete(&MsgCompleteBridge{}, "bridge/CompleteBridge", nil)
	cdc.RegisterConcrete(&MsgUpdateEventParams{}, "bridge/UpdateEventParams", nil)
	cdc.RegisterConcrete(&MsgUpdateProposeParams{}, "bridge/UpdateProposeParams", nil)
	cdc.RegisterConcrete(&MsgUpdateSafetyParams{}, "bridge/UpdateSafetyParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAcknowledgeBridge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCompleteBridge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateEventParams{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateProposeParams{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateSafetyParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
