package assets

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sample"
	assetssimulation "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/simulation"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = assetssimulation.FindAccount
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
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	assetsGenesis := types.GenesisState{
		Assets: []types.Asset{
			types.AssetTDai,
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&assetsGenesis)
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	// this line is used by starport scaffolding # simapp/module/operation
	return operations
}

// TODO(DEC-906): implement simulated gov proposal.
// ProposalMsgs doesn't return any content functions for governance proposals
func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}
