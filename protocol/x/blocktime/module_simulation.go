package blocktime

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	ratelimitsimulation "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = ratelimitsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

// GenerateGenesisState creates a randomized GenState of the module.
// Note the blocktime module is intentionally initialized as empty to avoid triggering
// chain outage withdrawal gating during simulation testing for withdrawals and transfers.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	blocktimeGenesis := types.GenesisState{}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&blocktimeGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
