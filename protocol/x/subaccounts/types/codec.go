package types

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(module.InterfaceRegistry)
)
