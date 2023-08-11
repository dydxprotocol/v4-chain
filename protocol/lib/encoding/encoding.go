package encoding

import (
	"cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// MakeEncodingConfig creates an EncodingConfig.
func MakeEncodingConfig(moduleBasics module.BasicManager) params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             amino,
	}

	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	moduleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
