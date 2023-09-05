package module

import (
	"cosmossdk.io/math"
	"encoding/json"
	"time"

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
// TODO(DEC-1776): Adjust below values based on final state of `Network Parameters` doc.
func (SlashingModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := slashingtypes.DefaultGenesisState()

	// No slashing for downtime and double-signing.
	genState.Params.SlashFractionDowntime = math.LegacyZeroDec()
	genState.Params.SlashFractionDoubleSign = math.LegacyZeroDec()

	// 1 minute jail duration for downtime.
	genState.Params.DowntimeJailDuration = 1 * time.Minute

	// 3000 blocks (about 4.167 hours @ 5 second/block).
	genState.Params.SignedBlocksWindow = 3000

	// Require 5% minimum liveness per signed block window.
	genState.Params.MinSignedPerWindow = math.LegacyMustNewDecFromStr("0.05")

	return cdc.MustMarshalJSON(genState)
}
