package epochs

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	epochssimulation "github.com/dydxprotocol/v4-chain/protocol/x/epochs/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = epochssimulation.FindAccount
	_ = sims.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

var (
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasProposalMsgs     = AppModule{}
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	epochssimulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the epoch module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil // `Epochs` module does not support any operations.
}

// TODO(DEC-906): implement simulated gov proposal.
// ProposalMsgs doesn't return any content functions for governance proposals
func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}
