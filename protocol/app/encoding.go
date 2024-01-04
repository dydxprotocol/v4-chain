package app

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

var encodingConfig = initEncodingConfig()

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// GetEncodingConfig returns the EncodingConfig.
func GetEncodingConfig() EncodingConfig {
	return encodingConfig
}

// makeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := module.InterfaceRegistry
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// initEncodingConfig initializes an EncodingConfig.
func initEncodingConfig() EncodingConfig {
	encConfig := makeEncodingConfig()

	// This is currently required in order to support various CLI commands such as the `dydxprotocold status` command.
	std.RegisterLegacyAminoCodec(encConfig.Amino)
	std.RegisterInterfaces(encConfig.InterfaceRegistry)

	// Skipping `ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)` because it's not needed.
	basic_manager.ModuleBasics.RegisterInterfaces(encConfig.InterfaceRegistry)

	return encConfig
}
