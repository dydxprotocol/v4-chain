package canonicalusdc_test

import (
	"encoding/json"
	"testing"

	"github.com/cometbft/cometbft/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibccoretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Counterparty-side channel IDs. Core checks that a packet's source matches
// the channel's counterparty before the app sees it, so adversarial packets
// keep these and vary only the data.
const (
	advNobleCounterparty     = "channel-100"
	advInjectiveCounterparty = "channel-200"
	advThirdChannel          = "channel-2"
	advThirdCounterparty     = "channel-300"
	advMaxTransfer           = 1_000_000
	advCeiling               = 3_000_000
)

var (
	advInjectiveDenom = transfertypes.ParseDenomTrace("transfer/channel-1/uusdc").IBCDenom()
	advReceiver       = constants.AliceAccAddress
	advRelayer        = constants.BobAccAddress
)

func advControls(mode canonicalusdctypes.Mode) canonicalusdctypes.Controls {
	return canonicalusdctypes.Controls{
		Mode:                    mode,
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
		InjectiveDenom:          advInjectiveDenom,
		MaxTransferAmount:       dtypes.NewInt(advMaxTransfer),
		MigrationCeiling:        dtypes.NewInt(advCeiling),
		MaxPendingSettlements:   100,
	}
}

func advAddRoute(genesis *ibccoretypes.GenesisState, channelID, counterpartyChannel, connectionID, clientID string) {
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
		channeltypes.NewCounterparty(transfertypes.PortID, counterpartyChannel),
		[]string{connectionID},
		transfertypes.Version,
	)
	genesis.ChannelGenesis.Channels = append(
		genesis.ChannelGenesis.Channels,
		channeltypes.NewIdentifiedChannel(transfertypes.PortID, channelID, channel),
	)
}

func advGenesis() types.GenesisDoc {
	genesis := testapp.DefaultGenesis()
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		panic(err)
	}
	cdc := app.GetEncodingConfig().Codec
	var ibcGenesis ibccoretypes.GenesisState
	cdc.MustUnmarshalJSON(appState[ibcexported.ModuleName], &ibcGenesis)
	advAddRoute(&ibcGenesis, "channel-0", advNobleCounterparty, "connection-0", "07-tendermint-0")
	advAddRoute(&ibcGenesis, "channel-1", advInjectiveCounterparty, "connection-1", "07-tendermint-1")
	advAddRoute(&ibcGenesis, advThirdChannel, advThirdCounterparty, "connection-2", "07-tendermint-2")
	ibcGenesis.ConnectionGenesis.NextConnectionSequence = uint64(len(ibcGenesis.ConnectionGenesis.Connections))
	ibcGenesis.ChannelGenesis.NextChannelSequence = uint64(len(ibcGenesis.ChannelGenesis.Channels))
	appState[ibcexported.ModuleName] = cdc.MustMarshalJSON(&ibcGenesis)
	bz, err := json.Marshal(appState)
	if err != nil {
		panic(err)
	}
	genesis.AppState = bz
	return genesis
}

type advEnv struct {
	t   *testing.T
	ctx sdk.Context
	app *app.App
	seq uint64
}

// advSetup starts the real app with three open transfer channels and moves
// canonical routing into mode through the real msg server, so the Noble
// backing snapshot is taken as in production.
func advSetup(t *testing.T, mode canonicalusdctypes.Mode) *advEnv {
	t.Helper()
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(advGenesis).Build()
	ctx := tApp.InitChain()
	env := &advEnv{t: t, ctx: ctx, app: tApp.App}
	if mode != canonicalusdctypes.Mode_MODE_DISABLED {
		_, err := keeper.NewMsgServer(tApp.App.CanonicalUsdcKeeper).UpdateControls(ctx, &canonicalusdctypes.MsgUpdateControls{
			Authority: lib.GovModuleAddress.String(),
			Controls:  advControls(canonicalusdctypes.Mode_MODE_GRADUAL),
		})
		require.NoError(t, err)
	}
	if mode == canonicalusdctypes.Mode_MODE_PAUSED {
		_, err := keeper.NewMsgServer(tApp.App.CanonicalUsdcKeeper).UpdateControls(ctx, &canonicalusdctypes.MsgUpdateControls{
			Authority: lib.GovModuleAddress.String(),
			Controls:  advControls(canonicalusdctypes.Mode_MODE_PAUSED),
		})
		require.NoError(t, err)
	}
	return env
}

func (e *advEnv) packet(destChannel, denom, amount, receiver, memo string) channeltypes.Packet {
	sourceChannel := map[string]string{
		"channel-0":     advNobleCounterparty,
		"channel-1":     advInjectiveCounterparty,
		advThirdChannel: advThirdCounterparty,
	}[destChannel]
	e.seq++
	data := transfertypes.FungibleTokenPacketData{
		Denom:    denom,
		Amount:   amount,
		Sender:   "counterparty1sender",
		Receiver: receiver,
		Memo:     memo,
	}
	return channeltypes.NewPacket(
		transfertypes.ModuleCdc.MustMarshalJSON(&data),
		e.seq,
		transfertypes.PortID,
		sourceChannel,
		transfertypes.PortID,
		destChannel,
		clienttypes.NewHeight(1, 1_000_000),
		0,
	)
}

func (e *advEnv) stack() ibcPortModule {
	module, ok := e.app.IBCKeeper.Router.GetRoute(transfertypes.ModuleName)
	require.True(e.t, ok)
	return module
}

type ibcPortModule interface {
	OnRecvPacket(sdk.Context, channeltypes.Packet, sdk.AccAddress) ibcexported.Acknowledgement
}

// recvRaw calls the transfer stack without core's cache, so any state the
// stack itself fails to roll back on an error ack stays visible.
func (e *advEnv) recvRaw(packet channeltypes.Packet) ibcexported.Acknowledgement {
	return e.stack().OnRecvPacket(e.ctx, packet, advRelayer)
}

// recv mirrors core RecvPacket: state is kept only for a successful ack.
func (e *advEnv) recv(packet channeltypes.Packet) ibcexported.Acknowledgement {
	cacheCtx, write := e.ctx.CacheContext()
	ack := e.stack().OnRecvPacket(cacheCtx, packet, advRelayer)
	if ack == nil || ack.Success() {
		write()
	}
	return ack
}

type advSnapshot struct {
	Ledger         canonicalusdctypes.Ledger
	Pending        uint32
	LogicalSupply  sdkmath.Int
	PhysicalSupply sdkmath.Int
	ModuleLogical  sdkmath.Int
	ModulePhysical sdkmath.Int
	RecvLogical    sdkmath.Int
	RecvPhysical   sdkmath.Int
	LogicalCap     ratelimittypes.DenomCapacity
	PhysicalCap    ratelimittypes.DenomCapacity
	Escrow1        sdk.Coins
}

func (e *advEnv) snapshot(receiver sdk.AccAddress) advSnapshot {
	k, bank := e.app.CanonicalUsdcKeeper, e.app.BankKeeper
	module := canonicalusdctypes.ModuleAddress
	return advSnapshot{
		Ledger:         k.GetLedger(e.ctx),
		Pending:        k.PendingCount(e.ctx),
		LogicalSupply:  bank.GetSupply(e.ctx, assettypes.UusdcDenom).Amount,
		PhysicalSupply: bank.GetSupply(e.ctx, advInjectiveDenom).Amount,
		ModuleLogical:  bank.GetBalance(e.ctx, module, assettypes.UusdcDenom).Amount,
		ModulePhysical: bank.GetBalance(e.ctx, module, advInjectiveDenom).Amount,
		RecvLogical:    bank.GetBalance(e.ctx, receiver, assettypes.UusdcDenom).Amount,
		RecvPhysical:   bank.GetBalance(e.ctx, receiver, advInjectiveDenom).Amount,
		LogicalCap:     e.app.RatelimitKeeper.GetDenomCapacity(e.ctx, assettypes.UusdcDenom),
		PhysicalCap:    e.app.RatelimitKeeper.GetDenomCapacity(e.ctx, advInjectiveDenom),
		Escrow1:        bank.GetAllBalances(e.ctx, transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-1")),
	}
}

func (e *advEnv) requireBacking() {
	reason, broken := keeper.BackingInvariant(*e.app.CanonicalUsdcKeeper)(e.ctx)
	require.False(e.t, broken, reason)
}

// requireErrorAck checks the ack is an ICS-04 error ack, which is what lets
// the counterparty's ICS20 refund the sender on OnAcknowledgementPacket.
func requireErrorAck(t *testing.T, ack ibcexported.Acknowledgement) {
	t.Helper()
	require.NotNil(t, ack)
	require.False(t, ack.Success())
	var decoded channeltypes.Acknowledgement
	require.NoError(t, transfertypes.ModuleCdc.UnmarshalJSON(ack.Acknowledgement(), &decoded))
	_, isErr := decoded.Response.(*channeltypes.Acknowledgement_Error)
	require.True(t, isErr)
}

func TestAdvMalformedInjectivePacketsLeaveNoState(t *testing.T) {
	cases := map[string]struct {
		amount   string
		receiver string
	}{
		"zero amount":          {amount: "0", receiver: advReceiver.String()},
		"negative amount":      {amount: "-5", receiver: advReceiver.String()},
		"non-numeric amount":   {amount: "abc", receiver: advReceiver.String()},
		"257-bit amount":       {amount: "231584178474632390847141970017375815706539969331281128078915168015826259279872", receiver: advReceiver.String()},
		"over max transfer":    {amount: "1000001", receiver: advReceiver.String()},
		"invalid receiver":     {amount: "10", receiver: "not-bech32"},
		"wrong hrp receiver":   {amount: "10", receiver: "cosmos1qqqsyqcyq5rqwzqfpg9scrgwpugpzysn5cyl7e"},
		"blocked receiver":     {amount: "10", receiver: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()},
		"transfer escrow mod":  {amount: "10", receiver: authtypes.NewModuleAddress(transfertypes.ModuleName).String()},
		"empty receiver":       {amount: "10", receiver: ""},
		"leading space amount": {amount: " 10", receiver: advReceiver.String()},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
			before := env.snapshot(advReceiver)
			ack := env.recvRaw(env.packet("channel-1", "uusdc", tc.amount, tc.receiver, ""))
			requireErrorAck(t, ack)
			require.Equal(t, before, env.snapshot(advReceiver))
			env.requireBacking()
		})
	}
}

func TestAdvMigrationCeilingRollsBackEverything(t *testing.T) {
	env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
	for i := 0; i < 3; i++ {
		require.True(t, env.recv(env.packet("channel-1", "uusdc", "1000000", advReceiver.String(), "")).Success())
	}
	env.requireBacking()
	before := env.snapshot(advReceiver)
	require.Equal(t, "3000000", before.Ledger.InjectiveBacking.String())

	ack := env.recvRaw(env.packet("channel-1", "uusdc", "1", advReceiver.String(), ""))
	requireErrorAck(t, ack)
	require.Equal(t, before, env.snapshot(advReceiver))
	env.requireBacking()
}

func TestAdvAlternateAmountEncodingsAreConsistent(t *testing.T) {
	// big.Int base-0 parsing is shared by ICS20, x/ratelimit and the
	// middleware, so prefixed forms must credit the same amount everywhere.
	cases := map[string]int64{"+10": 10, "0x10": 16, "1_000": 1000, "0b101": 5, "0o17": 15}
	for amount, want := range cases {
		t.Run(amount, func(t *testing.T) {
			env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
			before := env.snapshot(advReceiver)
			ack := env.recv(env.packet("channel-1", "uusdc", amount, advReceiver.String(), ""))
			require.True(t, ack.Success())
			after := env.snapshot(advReceiver)
			require.Equal(t, before.RecvLogical.AddRaw(want), after.RecvLogical)
			require.Equal(t, before.ModulePhysical.AddRaw(want), after.ModulePhysical)
			require.Equal(t, before.LogicalSupply.AddRaw(want), after.LogicalSupply)
			require.Equal(t, sdkmath.NewInt(want).String(), after.Ledger.InjectiveBacking.String())
			env.requireBacking()
		})
	}
}

func TestAdvInjectiveDepositCreditsLogicalOnly(t *testing.T) {
	env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
	before := env.snapshot(advReceiver)
	memo := `{"forward":{"receiver":"x","port":"transfer","channel":"channel-0"},"wasm":{"contract":"x"},` +
		`"src_callback":{"address":"x"},"dest_callback":{"address":"x"}}`
	ack := env.recv(env.packet("channel-1", "uusdc", "500", advReceiver.String(), memo))
	require.True(t, ack.Success())
	after := env.snapshot(advReceiver)
	require.Equal(t, before.RecvLogical.AddRaw(500), after.RecvLogical)
	require.Equal(t, before.RecvPhysical, after.RecvPhysical)
	require.Equal(t, before.ModulePhysical.AddRaw(500), after.ModulePhysical)
	require.Equal(t, "500", after.Ledger.InjectiveBacking.String())
	// Inbound Injective USDC raises capacity for the physical denom only.
	require.Equal(t, before.LogicalCap, after.LogicalCap)
	env.requireBacking()
}

func TestAdvInjectiveDepositToCanonicalModuleAddress(t *testing.T) {
	env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
	module := canonicalusdctypes.ModuleAddress
	before := env.snapshot(module)
	ack := env.recv(env.packet("channel-1", "uusdc", "500", module.String(), ""))
	// Accepted: the minted USDC is stranded in the module, but accounting holds.
	require.True(t, ack.Success())
	after := env.snapshot(module)
	require.Equal(t, before.ModuleLogical.AddRaw(500), after.ModuleLogical)
	require.Equal(t, "500", after.Ledger.InjectiveBacking.String())
	env.requireBacking()
}

func TestAdvNonCanonicalDenomsOnInjectiveChannelPassThrough(t *testing.T) {
	cases := map[string]string{
		// Noble USDC routed Noble -> Injective -> dYdX.
		"multi-hop noble via injective": "transfer/channel-77/uusdc",
		"case variant":                  "UUSDC",
		"other native":                  "inj",
	}
	for name, denom := range cases {
		t.Run(name, func(t *testing.T) {
			env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
			before := env.snapshot(advReceiver)
			ack := env.recv(env.packet("channel-1", denom, "500", advReceiver.String(), ""))
			require.True(t, ack.Success())
			after := env.snapshot(advReceiver)
			voucher := transfertypes.ParseDenomTrace("transfer/channel-1/" + denom).IBCDenom()
			require.NotEqual(t, advInjectiveDenom, voucher)
			require.NotEqual(t, assettypes.UusdcDenom, voucher)
			require.Equal(t, "500", env.app.BankKeeper.GetBalance(env.ctx, advReceiver, voucher).Amount.String())
			require.Equal(t, before.Ledger, after.Ledger)
			require.Equal(t, before.LogicalSupply, after.LogicalSupply)
			require.Equal(t, before.RecvLogical, after.RecvLogical)
			env.requireBacking()
		})
	}
}

func TestAdvSourceTracesCannotForgeCanonicalDenoms(t *testing.T) {
	// A counterparty claiming to return dYdX-origin coins makes ICS20 unescrow
	// from the destination channel's escrow, never mint.
	cases := map[string]struct {
		channel string
		denom   string
	}{
		"injective voucher unwinding over injective": {
			channel: "channel-1", denom: "transfer/" + advInjectiveCounterparty + "/transfer/channel-1/uusdc",
		},
		"logical unwinding over injective": {
			channel: "channel-1", denom: "transfer/" + advInjectiveCounterparty + "/transfer/channel-0/uusdc",
		},
		"logical unwinding over third channel": {
			channel: advThirdChannel, denom: "transfer/" + advThirdCounterparty + "/transfer/channel-0/uusdc",
		},
		"injective voucher over third channel": {
			channel: advThirdChannel, denom: "transfer/" + advThirdCounterparty + "/transfer/channel-1/uusdc",
		},
		"logical unwinding over noble": {
			channel: "channel-0", denom: "transfer/" + advNobleCounterparty + "/transfer/channel-0/uusdc",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
			before := env.snapshot(advReceiver)
			ack := env.recvRaw(env.packet(tc.channel, tc.denom, "500", advReceiver.String(), ""))
			requireErrorAck(t, ack)
			require.Equal(t, before, env.snapshot(advReceiver))
			env.requireBacking()
		})
	}
}

func TestAdvLegacyEscrowUnwindsOverInjectiveWithoutLedgerChange(t *testing.T) {
	env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
	escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-1")
	coins := sdk.NewCoins(sdk.NewInt64Coin(assettypes.UusdcDenom, 700))
	// Logical USDC escrowed on channel-1 while DISABLED.
	require.NoError(t, env.app.BankKeeper.SendCoins(env.ctx, constants.CarlAccAddress, escrow, coins))
	env.app.TransferKeeper.SetTotalEscrowForDenom(env.ctx, coins[0])
	before := env.snapshot(advReceiver)
	denom := "transfer/" + advInjectiveCounterparty + "/transfer/channel-0/uusdc"
	ack := env.recv(env.packet("channel-1", denom, "700", advReceiver.String(), ""))
	require.True(t, ack.Success())
	after := env.snapshot(advReceiver)
	require.Equal(t, before.RecvLogical.AddRaw(700), after.RecvLogical)
	require.Equal(t, before.Ledger, after.Ledger)
	require.Equal(t, before.LogicalSupply, after.LogicalSupply)
	env.requireBacking()

	// Unwinding more than was escrowed fails and changes nothing.
	before = env.snapshot(advReceiver)
	requireErrorAck(t, env.recvRaw(env.packet("channel-1", denom, "1", advReceiver.String(), "")))
	require.Equal(t, before, env.snapshot(advReceiver))
}

func TestAdvInjectiveDenomOnOtherChannels(t *testing.T) {
	for _, channel := range []string{"channel-0", advThirdChannel} {
		t.Run(channel, func(t *testing.T) {
			env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
			before := env.snapshot(advReceiver)
			ack := env.recvRaw(env.packet(channel, "uusdc", "500", advReceiver.String(), ""))
			after := env.snapshot(advReceiver)
			if channel == "channel-0" {
				requireErrorAck(t, ack)
				require.Equal(t, before, after)
			} else {
				require.True(t, ack.Success())
				require.Equal(t, before.Ledger, after.Ledger)
				require.Equal(t, before.LogicalSupply, after.LogicalSupply)
				require.Equal(t, before.RecvLogical, after.RecvLogical)
				require.Equal(t, before.ModulePhysical, after.ModulePhysical)
			}
			env.requireBacking()
		})
	}
}

func TestAdvModes(t *testing.T) {
	t.Run("paused rejects injective and noble USDC", func(t *testing.T) {
		env := advSetup(t, canonicalusdctypes.Mode_MODE_PAUSED)
		for _, channel := range []string{"channel-0", "channel-1"} {
			before := env.snapshot(advReceiver)
			requireErrorAck(t, env.recvRaw(env.packet(channel, "uusdc", "500", advReceiver.String(), "")))
			require.Equal(t, before, env.snapshot(advReceiver))
		}
		env.requireBacking()
	})
	t.Run("disabled passes injective USDC through as a plain voucher", func(t *testing.T) {
		env := advSetup(t, canonicalusdctypes.Mode_MODE_DISABLED)
		before := env.snapshot(advReceiver)
		require.True(t, env.recv(env.packet("channel-1", "uusdc", "500", advReceiver.String(), "")).Success())
		after := env.snapshot(advReceiver)
		require.Equal(t, before.RecvPhysical.AddRaw(500), after.RecvPhysical)
		require.Equal(t, before.RecvLogical, after.RecvLogical)
		require.Equal(t, before.Ledger, after.Ledger)
	})
	t.Run("disabled accepts noble USDC", func(t *testing.T) {
		env := advSetup(t, canonicalusdctypes.Mode_MODE_DISABLED)
		before := env.snapshot(advReceiver)
		require.True(t, env.recv(env.packet("channel-0", "uusdc", "500", advReceiver.String(), "")).Success())
		require.Equal(t, before.RecvLogical.AddRaw(500), env.snapshot(advReceiver).RecvLogical)
	})
}

func TestAdvMalformedPacketData(t *testing.T) {
	env := advSetup(t, canonicalusdctypes.Mode_MODE_GRADUAL)
	packet := env.packet("channel-1", "uusdc", "10", advReceiver.String(), "")
	for name, data := range map[string][]byte{
		"not json":      []byte("garbage"),
		"empty":         {},
		"amount number": []byte(`{"denom":"uusdc","amount":10,"sender":"s","receiver":"` + advReceiver.String() + `"}`),
		"unknown field": []byte(`{"denom":"uusdc","amount":"10","sender":"s","receiver":"` + advReceiver.String() + `","x":1}`),
	} {
		t.Run(name, func(t *testing.T) {
			p := packet
			p.Data = data
			before := env.snapshot(advReceiver)
			requireErrorAck(t, env.recvRaw(p))
			require.Equal(t, before, env.snapshot(advReceiver))
		})
	}
}
