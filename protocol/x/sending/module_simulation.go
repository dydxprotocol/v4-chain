package sending

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sendingsimulation "github.com/dydxprotocol/v4-chain/protocol/x/sending/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
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
	sendingGenesis := types.GenesisState{
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&sendingGenesis)
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return sendingsimulation.WeightedOperations(
		simState.AppParams,
		am.keeper,
		am.accountKeeper,
		am.bankKeeper,
		am.subaccountsKeeper,
	)
}

// TODO(DEC-906): implement simulated gov proposal.
// ProposalMsgs doesn't return any content functions for governance proposals
func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}
