package app

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	"github.com/dydxprotocol/v4-chain/protocol/app/params"

	"github.com/cosmos/cosmos-sdk/std"
)

var encodingConfig params.EncodingConfig = initEncodingConfig()

// GetEncodingConfig returns the EncodingConfig.
func GetEncodingConfig() params.EncodingConfig {
	return encodingConfig
}

// initEncodingConfig initializes an EncodingConfig.
func initEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	// This is currently required in order to support various CLI commands such as the `dydxprotocold status` command.
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// Skipping `ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)` because it's not needed.
	basic_manager.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
