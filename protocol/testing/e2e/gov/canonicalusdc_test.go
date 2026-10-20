package gov_test

import (
	"encoding/json"
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibccoretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateControlsDisabledProposal(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return canonicalUsdcGenesis(false)
	}).Build()
	ctx := tApp.InitChain()

	ctx = testapp.SubmitAndTallyProposal(
		t,
		ctx,
		tApp,
		[]sdk.Msg{&canonicalusdctypes.MsgUpdateControls{
			Authority: lib.GovModuleAddress.String(),
			Controls:  canonicalusdctypes.DefaultControls(),
		}},
		testapp.TestSubmitProposalTxHeight,
		false,
		false,
		govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
	)
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, tApp.App.CanonicalUsdcKeeper.GetControls(ctx).Mode)
}

func TestUpdateControlsEnableGradualProposal(t *testing.T) {
	tests := map[string]struct {
		openRoutes     bool
		expectedStatus govtypesv1.ProposalStatus
	}{
		"routes open": {
			openRoutes:     true,
			expectedStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"routes missing": {
			openRoutes:     false,
			expectedStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
				return canonicalUsdcGenesis(tc.openRoutes)
			}).Build()
			ctx := tApp.InitChain()
			supply := tApp.App.BankKeeper.GetSupply(ctx, assettypes.UusdcDenom).Amount
			require.True(t, supply.IsPositive())

			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{&canonicalusdctypes.MsgUpdateControls{
					Authority: lib.GovModuleAddress.String(),
					Controls:  gradualControls(),
				}},
				testapp.TestSubmitProposalTxHeight,
				false,
				false,
				tc.expectedStatus,
			)

			controls := tApp.App.CanonicalUsdcKeeper.GetControls(ctx)
			ledger := tApp.App.CanonicalUsdcKeeper.GetLedger(ctx)
			if !tc.openRoutes {
				require.Equal(t, canonicalusdctypes.DefaultControls(), controls)
				require.Equal(t, "0", ledger.NobleBacking.String())
				return
			}
			require.Equal(t, gradualControls(), controls)
			require.Equal(t, supply.String(), ledger.NobleBacking.String())
			require.Equal(t, "0", ledger.InjectiveBacking.String())
		})
	}
}

func gradualControls() canonicalusdctypes.Controls {
	return canonicalusdctypes.Controls{
		Mode:                    canonicalusdctypes.Mode_MODE_GRADUAL,
		LogicalDenom:            assettypes.UusdcDenom,
		NobleChannel:            "channel-0",
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		NoblePacketDenom:        "uusdc",
		NobleWithdrawalsEnabled: true,
		InjectiveChannel:        "channel-1",
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          transfertypes.ParseDenomTrace("transfer/channel-1/uusdc").IBCDenom(),
		MaxTransferAmount:       dtypes.NewInt(1_000_000),
		MigrationCeiling:        dtypes.NewInt(10_000_000),
		MaxPendingSettlements:   100,
	}
}

// canonicalUsdcGenesis shortens the voting period and, when openRoutes is set,
// adds open transfer channels and connections for both canonical routes.
func canonicalUsdcGenesis(openRoutes bool) types.GenesisDoc {
	genesis := testapp.DefaultGenesis()
	testapp.UpdateGenesisDocWithAppStateForModule(
		&genesis,
		func(genesisState *govtypesv1.GenesisState) {
			genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
		},
	)
	if !openRoutes {
		return genesis
	}

	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		panic(err)
	}
	cdc := app.GetEncodingConfig().Codec
	var ibcGenesis ibccoretypes.GenesisState
	cdc.MustUnmarshalJSON(appState[ibcexported.ModuleName], &ibcGenesis)
	controls := gradualControls()
	addTransferRoute(&ibcGenesis, controls.NobleChannel, controls.NobleConnection, controls.NobleClient)
	addTransferRoute(&ibcGenesis, controls.InjectiveChannel, controls.InjectiveConnection, controls.InjectiveClient)
	ibcGenesis.ConnectionGenesis.NextConnectionSequence = uint64(len(ibcGenesis.ConnectionGenesis.Connections))
	ibcGenesis.ChannelGenesis.NextChannelSequence = uint64(len(ibcGenesis.ChannelGenesis.Channels))
	appState[ibcexported.ModuleName] = cdc.MustMarshalJSON(&ibcGenesis)

	appStateBytes, err := json.Marshal(appState)
	if err != nil {
		panic(err)
	}
	genesis.AppState = appStateBytes
	return genesis
}

func addTransferRoute(genesis *ibccoretypes.GenesisState, channelID, connectionID, clientID string) {
	connection := connectiontypes.NewConnectionEnd(
		connectiontypes.OPEN,
		clientID,
		connectiontypes.NewCounterparty(clientID, connectionID, commitmenttypes.NewMerklePrefix([]byte("ibc"))),
		connectiontypes.GetCompatibleVersions(),
		0,
	)
	genesis.ConnectionGenesis.Connections = append(
		genesis.ConnectionGenesis.Connections,
		connectiontypes.NewIdentifiedConnection(connectionID, connection),
	)
	genesis.ConnectionGenesis.ClientConnectionPaths = append(
		genesis.ConnectionGenesis.ClientConnectionPaths,
		connectiontypes.NewConnectionPaths(clientID, []string{connectionID}),
	)
	channel := channeltypes.NewChannel(
		channeltypes.OPEN,
		channeltypes.UNORDERED,
		channeltypes.NewCounterparty(transfertypes.PortID, channelID),
		[]string{connectionID},
		transfertypes.Version,
	)
	genesis.ChannelGenesis.Channels = append(
		genesis.ChannelGenesis.Channels,
		channeltypes.NewIdentifiedChannel(transfertypes.PortID, channelID, channel),
	)
}
