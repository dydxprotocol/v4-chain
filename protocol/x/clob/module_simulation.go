package clob

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	v4module "github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	clobsimulation "github.com/dydxprotocol/v4-chain/protocol/x/clob/simulation"
)

// avoid unused import issue
var (
	_                            = sample.AccAddress
	_                            = clobsimulation.FindAccount
	_                            = sims.StakePerAccount
	_                            = simulation.MsgEntryKind
	_                            = baseapp.Paramspace
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasProposalMsgs     = AppModule{}
)

const (
	opWeightMsgProposedOperations = "op_weight_msg_temp_operations"
	// TODO(DEC-571): Determine the simulation weight value
	defaultWeightMsgProposedOperations int = 100

	opWeightMsgPlaceOrder = "op_weight_msg_place_order"
	// TODO(DEC-571): Determine the simulation weight value
	defaultWeightMsgPlaceOrder int = 100

	opWeightMsgCancelOrder          = "op_weight_msg_cancel_order"
	defaultWeightMsgCancelOrder int = 25

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	clobsimulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	protoCdc := codec.NewProtoCodec(v4module.InterfaceRegistry)
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgProposedOperations int
	simState.AppParams.GetOrGenerate(opWeightMsgProposedOperations, &weightMsgProposedOperations, nil,
		func(_ *rand.Rand) {
			weightMsgProposedOperations = defaultWeightMsgProposedOperations
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgProposedOperations,
		clobsimulation.SimulateMsgProposedOperations(am.accountKeeper, am.bankKeeper, *am.keeper),
	))

	var weightMsgPlaceOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgPlaceOrder, &weightMsgPlaceOrder, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceOrder = defaultWeightMsgPlaceOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPlaceOrder,
		clobsimulation.SimulateMsgPlaceOrder(am.accountKeeper, am.bankKeeper, am.subaccountsKeeper, *am.keeper, protoCdc),
	))

	var weightMsgCancelOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelOrder, &weightMsgCancelOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCancelOrder = defaultWeightMsgCancelOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelOrder,
		clobsimulation.SimulateMsgCancelOrder(
			am.accountKeeper,
			am.bankKeeper,
			am.subaccountsKeeper,
			*am.keeper,
			am.keeper.MemClob,
			protoCdc,
		),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// TODO(DEC-906): implement simulated gov proposal.
// ProposalMsgs doesn't return any content functions for governance proposals
func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}
