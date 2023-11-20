package module

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// SlashingModuleBasic defines a custom wrapper around the default `x/slashing` module's `AppModuleBasic`
// implementation to provide custom default genesis state.
type SlashingModuleBasic struct {
	slashing.AppModuleBasic
}

// DefaultGenesis returns custom `x/slashing` module genesis state.
func (SlashingModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := slashingtypes.DefaultGenesisState()

	// No slashing for downtime and double-signing.
	genState.Params.SlashFractionDowntime = sdkmath.LegacyZeroDec()
	genState.Params.SlashFractionDoubleSign = sdkmath.LegacyZeroDec()

	// 1 minute jail duration for downtime.
	genState.Params.DowntimeJailDuration = 1 * time.Minute

	// 3000 blocks (about 4.167 hours @ 5 second/block).
	genState.Params.SignedBlocksWindow = 3000

	// Require 5% minimum liveness per signed block window.
	genState.Params.MinSignedPerWindow = sdkmath.LegacyMustNewDecFromStr("0.05")

	return cdc.MustMarshalJSON(genState)
}
