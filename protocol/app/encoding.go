package app

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	"github.com/dydxprotocol/v4-chain/protocol/app/params"

	"github.com/cosmos/cosmos-sdk/std"
)

var encodingConfig params.EncodingConfig = MakeEncodingConfig()

// GetEncodingConfig returns the EncodingConfig.
func GetEncodingConfig() params.EncodingConfig {
	return encodingConfig
}

// MakeEncodingConfig creates an EncodingConfig.
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	basic_manager.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
