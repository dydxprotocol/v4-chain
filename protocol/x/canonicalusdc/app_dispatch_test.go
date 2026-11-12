package canonicalusdc_test

import (
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/gogoproto/proto"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icagenesistypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/genesis/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibccoretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	dispatchNobleChannel      = "channel-0"
	dispatchInjectiveChannel  = "channel-1"
	dispatchICAChannel        = "channel-2"
	dispatchICAControllerPort = icatypes.ControllerPortPrefix + "owner"
	dispatchRemoteChannel     = "channel-9"
	dispatchReserve           = 1_000
	dispatchBalance           = 500
	dispatchAmount            = 40
)

var (
	// dispatchOwner signs the MsgTransfer: the authz granter, or the interchain
	// account that the ICA host executes for.
	dispatchOwner   = sdk.AccAddress("canonical-usdc-owner")
	dispatchGrantee = sdk.AccAddress("canonical-usdc-grant")
	dispatchTimeout = uint64(time.Unix(1_000_000, 0).UnixNano())
)

func dispatchControls() types.Controls {
	return types.Controls{
		Mode:                    types.Mode_MODE_GRADUAL,
		LogicalDenom:            assettypes.UusdcDenom,
		NobleChannel:            dispatchNobleChannel,
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		NoblePacketDenom:        "uusdc",
		NobleWithdrawalsEnabled: true,
		InjectiveChannel:        dispatchInjectiveChannel,
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          transfertypes.ParseDenomTrace("transfer/" + dispatchInjectiveChannel + "/uusdc").IBCDenom(),
		MaxPendingSettlements:   10,
	}
}

// dispatchGenesis starts the chain in GRADUAL with both routes open and able
// to send: active light clients, channel capabilities, registered denom traces,
// an Injective reserve, and an interchain account owned by dispatchOwner.
func dispatchGenesis(t *testing.T) cmttypes.GenesisDoc {
	t.Helper()
	genesis := testapp.DefaultGenesis()
	controls := dispatchControls()

	supply := math.ZeroInt()
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *banktypes.GenesisState) {
		state.Balances = append(state.Balances,
			banktypes.Balance{
				Address: dispatchOwner.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(controls.LogicalDenom, dispatchBalance)),
			},
			banktypes.Balance{
				Address: types.ModuleAddress.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(controls.InjectiveDenom, dispatchReserve)),
			},
		)
		for _, balance := range state.Balances {
			supply = supply.Add(balance.Coins.AmountOf(controls.LogicalDenom))
		}
	})
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *types.GenesisState) {
		state.Controls = controls
		state.Ledger = types.DefaultLedger()
		state.Ledger.InjectiveBacking = dtypes.NewInt(dispatchReserve)
		state.Ledger.NobleBacking = types.NewAmount(supply.SubRaw(dispatchReserve))
	})

	var appState map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(genesis.AppState, &appState))
	cdc := app.GetEncodingConfig().Codec
	appState[ibcexported.ModuleName] = dispatchIBCGenesis(cdc, appState[ibcexported.ModuleName], controls)
	appState[capabilitytypes.ModuleName] = cdc.MustMarshalJSON(&capabilitytypes.GenesisState{
		Index: 3,
		Owners: []capabilitytypes.GenesisOwners{
			channelCapability(1, controls.NobleChannel),
			channelCapability(2, controls.InjectiveChannel),
		},
	})
	appState[transfertypes.ModuleName] = cdc.MustMarshalJSON(&transfertypes.GenesisState{
		PortId: transfertypes.PortID,
		DenomTraces: transfertypes.Traces{
			transfertypes.ParseDenomTrace(controls.NobleTrace()),
			transfertypes.ParseDenomTrace(controls.InjectiveTrace()),
		}.Sort(),
		Params: transfertypes.DefaultParams(),
	})
	icaGenesis := icagenesistypes.DefaultGenesis()
	icaGenesis.HostGenesisState.Params = icahosttypes.NewParams(true, []string{icahosttypes.AllowAllHostMsgs})
	icaGenesis.HostGenesisState.ActiveChannels = []icagenesistypes.ActiveChannel{{
		ConnectionId: controls.NobleConnection,
		PortId:       dispatchICAControllerPort,
		ChannelId:    dispatchICAChannel,
	}}
	icaGenesis.HostGenesisState.InterchainAccounts = []icagenesistypes.RegisteredInterchainAccount{{
		ConnectionId:   controls.NobleConnection,
		PortId:         dispatchICAControllerPort,
		AccountAddress: dispatchOwner.String(),
	}}
	appState[icatypes.ModuleName] = cdc.MustMarshalJSON(icaGenesis)

	appStateBytes, err := json.Marshal(appState)
	require.NoError(t, err)
	genesis.AppState = appStateBytes
	return genesis
}

func channelCapability(index uint64, channel string) capabilitytypes.GenesisOwners {
	name := host.ChannelCapabilityPath(transfertypes.PortID, channel)
	return capabilitytypes.GenesisOwners{
		Index: index,
		IndexOwners: capabilitytypes.CapabilityOwners{Owners: []capabilitytypes.Owner{
			{Module: ibcexported.ModuleName, Name: name},
			{Module: transfertypes.ModuleName, Name: name},
		}},
	}
}

func dispatchIBCGenesis(cdc codec.Codec, raw json.RawMessage, controls types.Controls) json.RawMessage {
	var genesis ibccoretypes.GenesisState
	cdc.MustUnmarshalJSON(raw, &genesis)

	routes := []struct{ chainID, client, connection, channel string }{
		{"noble-1", controls.NobleClient, controls.NobleConnection, controls.NobleChannel},
		{"injective-1", controls.InjectiveClient, controls.InjectiveConnection, controls.InjectiveChannel},
	}
	height := clienttypes.NewHeight(1, 10)
	for _, route := range routes {
		clientState := ibctm.NewClientState(
			route.chainID, ibctm.DefaultTrustLevel, 24*time.Hour, 48*time.Hour, 10*time.Second,
			height, commitmenttypes.GetSDKSpecs(), []string{"upgrade", "upgradedIBCState"},
		)
		consensusState := ibctm.NewConsensusState(
			time.Unix(1, 0), commitmenttypes.NewMerkleRoot([]byte("root")), tmhash.Sum([]byte("validators")),
		)
		genesis.ClientGenesis.Clients = append(
			genesis.ClientGenesis.Clients, clienttypes.NewIdentifiedClientState(route.client, clientState),
		)
		genesis.ClientGenesis.ClientsConsensus = append(
			genesis.ClientGenesis.ClientsConsensus,
			clienttypes.NewClientConsensusStates(route.client, []clienttypes.ConsensusStateWithHeight{
				clienttypes.NewConsensusStateWithHeight(height, consensusState),
			}),
		)
		genesis.ConnectionGenesis.Connections = append(
			genesis.ConnectionGenesis.Connections,
			connectiontypes.NewIdentifiedConnection(route.connection, connectiontypes.NewConnectionEnd(
				connectiontypes.OPEN,
				route.client,
				connectiontypes.NewCounterparty(
					route.client, route.connection, commitmenttypes.NewMerklePrefix([]byte("ibc")),
				),
				connectiontypes.GetCompatibleVersions(),
				0,
			)),
		)
		genesis.ConnectionGenesis.ClientConnectionPaths = append(
			genesis.ConnectionGenesis.ClientConnectionPaths,
			connectiontypes.NewConnectionPaths(route.client, []string{route.connection}),
		)
		genesis.ChannelGenesis.Channels = append(
			genesis.ChannelGenesis.Channels,
			channeltypes.NewIdentifiedChannel(transfertypes.PortID, route.channel, channeltypes.NewChannel(
				channeltypes.OPEN,
				channeltypes.UNORDERED,
				channeltypes.NewCounterparty(transfertypes.PortID, dispatchRemoteChannel),
				[]string{route.connection},
				transfertypes.Version,
			)),
		)
		genesis.ChannelGenesis.SendSequences = append(
			genesis.ChannelGenesis.SendSequences,
			channeltypes.NewPacketSequence(transfertypes.PortID, route.channel, 1),
		)
	}

	metadata := icatypes.NewMetadata(
		icatypes.Version, "connection-7", controls.NobleConnection, dispatchOwner.String(),
		icatypes.EncodingProtobuf, icatypes.TxTypeSDKMultiMsg,
	)
	genesis.ChannelGenesis.Channels = append(
		genesis.ChannelGenesis.Channels,
		channeltypes.NewIdentifiedChannel(icatypes.HostPortID, dispatchICAChannel, channeltypes.NewChannel(
			channeltypes.OPEN,
			channeltypes.ORDERED,
			channeltypes.NewCounterparty(dispatchICAControllerPort, dispatchRemoteChannel),
			[]string{controls.NobleConnection},
			string(icatypes.ModuleCdc.MustMarshalJSON(&metadata)),
		)),
	)
	genesis.ClientGenesis.NextClientSequence = uint64(len(genesis.ClientGenesis.Clients))
	genesis.ConnectionGenesis.NextConnectionSequence = uint64(len(genesis.ConnectionGenesis.Connections))
	genesis.ChannelGenesis.NextChannelSequence = uint64(len(genesis.ChannelGenesis.Channels))
	return cdc.MustMarshalJSON(&genesis)
}

type dispatchEnv struct {
	t        *testing.T
	app      *app.App
	ctx      sdk.Context
	controls types.Controls
}

func newDispatchEnv(t *testing.T) dispatchEnv {
	t.Helper()
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() cmttypes.GenesisDoc {
		return dispatchGenesis(t)
	}).Build()
	ctx := tApp.InitChain().WithIsCheckTx(false)
	env := dispatchEnv{t: t, app: tApp.App, ctx: ctx, controls: dispatchControls()}
	require.Equal(t, types.Mode_MODE_GRADUAL, env.app.CanonicalUsdcKeeper.GetControls(ctx).Mode)
	env.requireInvariant()
	return env
}

func (e dispatchEnv) requireInvariant() {
	e.t.Helper()
	reason, broken := keeper.BackingInvariant(*e.app.CanonicalUsdcKeeper)(e.ctx)
	require.False(e.t, broken, reason)
}

func (e dispatchEnv) balance(address sdk.AccAddress, denom string) int64 {
	return e.app.BankKeeper.GetBalance(e.ctx, address, denom).Amount.Int64()
}

func (e dispatchEnv) supply(denom string) math.Int {
	return e.app.BankKeeper.GetSupply(e.ctx, denom).Amount
}

func (e dispatchEnv) capacity() []dtypes.SerializableInt {
	return e.app.RatelimitKeeper.GetDenomCapacity(e.ctx, e.controls.LogicalDenom).CapacityList
}

func (e dispatchEnv) module(port string) porttypes.IBCModule {
	e.t.Helper()
	module, found := e.app.IBCKeeper.Router.GetRoute(port)
	require.True(e.t, found, "IBC route %s", port)
	return module
}

func (e dispatchEnv) withdrawal(channel string) *transfertypes.MsgTransfer {
	return &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    channel,
		Token:            sdk.NewInt64Coin(e.controls.LogicalDenom, dispatchAmount),
		Sender:           dispatchOwner.String(),
		Receiver:         "remote1receiver",
		TimeoutTimestamp: dispatchTimeout,
	}
}

// viaAuthz executes the withdrawal as a grantee of dispatchOwner.
func (e dispatchEnv) viaAuthz(msg *transfertypes.MsgTransfer) {
	e.t.Helper()
	authorization := authz.NewGenericAuthorization(sdk.MsgTypeURL(msg))
	require.NoError(e.t, e.app.AuthzKeeper.SaveGrant(e.ctx, dispatchGrantee, dispatchOwner, authorization, nil))
	exec := authz.NewMsgExec(dispatchGrantee, []sdk.Msg{msg})
	handler := e.app.MsgServiceRouter().Handler(&exec)
	require.NotNil(e.t, handler)
	_, err := handler(e.ctx, &exec)
	require.NoError(e.t, err)
}

// viaICAHost executes the withdrawal as the interchain account, delivered as
// an ICS27 packet on the host channel.
func (e dispatchEnv) viaICAHost(msg *transfertypes.MsgTransfer) {
	e.t.Helper()
	data, err := icatypes.SerializeCosmosTx(e.app.AppCodec(), []proto.Message{msg}, icatypes.EncodingProtobuf)
	require.NoError(e.t, err)
	packetData := icatypes.InterchainAccountPacketData{Type: icatypes.EXECUTE_TX, Data: data}
	packet := channeltypes.NewPacket(
		packetData.GetBytes(), 1,
		dispatchICAControllerPort, dispatchRemoteChannel,
		icatypes.HostPortID, dispatchICAChannel,
		clienttypes.NewHeight(1, 100), 0,
	)
	_, err = e.app.ICAHostKeeper.OnRecvPacket(e.ctx, packet)
	require.NoError(e.t, err)
}

// sentPacket rebuilds the packet ICS20 committed for the pending settlement
// and checks it against the stored packet commitment.
func (e dispatchEnv) sentPacket(pending types.PendingSettlement) channeltypes.Packet {
	e.t.Helper()
	trace, sender := e.controls.NobleTrace(), pending.Sender
	if pending.Route == types.Route_ROUTE_INJECTIVE {
		trace, sender = e.controls.InjectiveTrace(), types.ModuleAddress.String()
	}
	data := transfertypes.NewFungibleTokenPacketData(
		trace, pending.Amount.String(), sender, pending.PhysicalReceiver, "",
	)
	packet := channeltypes.NewPacket(
		data.GetBytes(), pending.Sequence,
		transfertypes.PortID, pending.SourceChannel,
		transfertypes.PortID, dispatchRemoteChannel,
		clienttypes.ZeroHeight(), dispatchTimeout,
	)
	commitment := e.app.IBCKeeper.ChannelKeeper.GetPacketCommitment(
		e.ctx, transfertypes.PortID, pending.SourceChannel, pending.Sequence,
	)
	require.Equal(e.t, channeltypes.CommitPacket(e.app.AppCodec(), packet), commitment)
	return packet
}

// A canonical withdrawal dispatched through authz MsgExec or the ICA host
// reaches the transfer decorator like a directly signed one: the pending
// settlement names the account the message was executed for, and every
// settlement outcome returns to that account.
func TestCanonicalWithdrawalDispatchedThroughAuthzAndICAHost(t *testing.T) {
	dispatchers := map[string]func(dispatchEnv, *transfertypes.MsgTransfer){
		"authz MsgExec": dispatchEnv.viaAuthz,
		"ICA host":      dispatchEnv.viaICAHost,
	}
	routes := map[string]types.Route{
		"noble":     types.Route_ROUTE_NOBLE,
		"injective": types.Route_ROUTE_INJECTIVE,
	}
	outcomes := map[string]func(porttypes.IBCModule, sdk.Context, channeltypes.Packet) error{
		"success": func(module porttypes.IBCModule, ctx sdk.Context, packet channeltypes.Packet) error {
			ack := channeltypes.NewResultAcknowledgement([]byte{1})
			return module.OnAcknowledgementPacket(ctx, packet, ack.Acknowledgement(), nil)
		},
		"error ack": func(module porttypes.IBCModule, ctx sdk.Context, packet channeltypes.Packet) error {
			ack := channeltypes.NewErrorAcknowledgement(types.ErrPaused)
			return module.OnAcknowledgementPacket(ctx, packet, ack.Acknowledgement(), nil)
		},
		"timeout": func(module porttypes.IBCModule, ctx sdk.Context, packet channeltypes.Packet) error {
			return module.OnTimeoutPacket(ctx, packet, nil)
		},
	}
	for dispatcherName, dispatch := range dispatchers {
		for routeName, route := range routes {
			for outcomeName, settle := range outcomes {
				t.Run(dispatcherName+"/"+routeName+"/"+outcomeName, func(t *testing.T) {
					env := newDispatchEnv(t)
					channel := env.controls.NobleChannel
					if route == types.Route_ROUTE_INJECTIVE {
						channel = env.controls.InjectiveChannel
					}
					before := snapshot(env)

					dispatch(env, env.withdrawal(channel))

					pending := env.app.CanonicalUsdcKeeper.GetAllPending(env.ctx)
					require.Len(t, pending, 1)
					require.Equal(t, route, pending[0].Route)
					require.Equal(t, dispatchOwner.String(), pending[0].Sender)
					require.Equal(t, "remote1receiver", pending[0].LogicalReceiver)
					require.Equal(t, int64(dispatchAmount), pending[0].Amount.BigInt().Int64())
					inFlight := snapshot(env)
					require.Equal(t, before.owner-dispatchAmount, inFlight.owner)
					requireInFlight(t, route, before, inFlight)
					env.requireInvariant()

					packet := env.sentPacket(pending[0])
					require.NoError(t, settle(env.module(transfertypes.ModuleName), env.ctx, packet))

					require.Zero(t, env.app.CanonicalUsdcKeeper.PendingCount(env.ctx))
					settled := snapshot(env)
					if outcomeName == "success" {
						requireDelivered(t, route, before, settled)
					} else {
						require.Equal(t, before, settled, "a failed withdrawal restores every balance and counter")
					}
					env.requireInvariant()
				})
			}
		}
	}
}

// dispatchState is everything a canonical withdrawal can move.
type dispatchState struct {
	owner            int64
	lockedLogical    int64
	reserve          int64
	logicalSupply    string
	nobleBacking     string
	injectiveBacking string
	pendingNoble     string
	pendingInjective string
	hourlyCapacity   string
}

func snapshot(env dispatchEnv) dispatchState {
	ledger := env.app.CanonicalUsdcKeeper.GetLedger(env.ctx)
	capacity := env.capacity()
	require.NotEmpty(env.t, capacity)
	return dispatchState{
		owner:            env.balance(dispatchOwner, env.controls.LogicalDenom),
		lockedLogical:    env.balance(types.ModuleAddress, env.controls.LogicalDenom),
		reserve:          env.balance(types.ModuleAddress, env.controls.InjectiveDenom),
		logicalSupply:    env.supply(env.controls.LogicalDenom).String(),
		nobleBacking:     ledger.NobleBacking.String(),
		injectiveBacking: ledger.InjectiveBacking.String(),
		pendingNoble:     ledger.PendingNoble.String(),
		pendingInjective: ledger.PendingInjective.String(),
		hourlyCapacity:   capacity[0].String(),
	}
}

func minus(value string, amount int64) string {
	parsed, ok := math.NewIntFromString(value)
	if !ok {
		panic("not an integer: " + value)
	}
	return parsed.SubRaw(amount).String()
}

// requireInFlight checks the state between the send and its settlement.
func requireInFlight(t *testing.T, route types.Route, before, inFlight dispatchState) {
	t.Helper()
	want := before
	want.owner -= dispatchAmount
	want.hourlyCapacity = minus(before.hourlyCapacity, dispatchAmount)
	if route == types.Route_ROUTE_NOBLE {
		want.logicalSupply = minus(before.logicalSupply, dispatchAmount)
		want.nobleBacking = minus(before.nobleBacking, dispatchAmount)
		want.pendingNoble = "40"
	} else {
		want.lockedLogical += dispatchAmount
		want.reserve -= dispatchAmount
		want.injectiveBacking = minus(before.injectiveBacking, dispatchAmount)
		want.pendingInjective = "40"
	}
	require.Equal(t, want, inFlight)
}

// requireDelivered checks the state after a successful acknowledgement.
func requireDelivered(t *testing.T, route types.Route, before, settled dispatchState) {
	t.Helper()
	want := before
	want.owner -= dispatchAmount
	want.hourlyCapacity = minus(before.hourlyCapacity, dispatchAmount)
	want.logicalSupply = minus(before.logicalSupply, dispatchAmount)
	if route == types.Route_ROUTE_NOBLE {
		want.nobleBacking = minus(before.nobleBacking, dispatchAmount)
	} else {
		want.reserve -= dispatchAmount
		want.injectiveBacking = minus(before.injectiveBacking, dispatchAmount)
	}
	require.Equal(t, want, settled)
}
