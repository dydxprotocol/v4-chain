package canonicalusdc

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	propSeeds         = 40
	propStepsPerSeed  = 150
	propInitialNoble  = 1000
	propMaxStepAmount = 300
)

// propChain models one chain: every operation runs in a cache context and is
// committed only when it succeeds, as a tx or an IBC callback would be. The
// stub bank and rate limiter are not store-backed, so they are snapshotted and
// restored on failure.
type propChain struct {
	r        roundTrip
	gov      string
	inflight []channeltypes.Packet
	holders  []sdk.AccAddress
}

func propSetup(t *testing.T) *propChain {
	t.Helper()
	gov := lib.GovModuleAddress.String()
	ctx, k, bank, ratelimit := propKeeper(t, []string{gov})
	controls := receiveControls()
	stub := &ics20Stub{
		bank: bank,
		ics4: k,
		traces: map[string]string{
			controls.LogicalDenom:   controls.NobleTrace(),
			controls.InjectiveDenom: controls.InjectiveTrace(),
		},
	}
	k.SetTransferMsgServer(stub)
	ledger := types.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(propInitialNoble)
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, ledger))
	inner := ics20App{receiveApp: &receiveApp{bank: bank, physicalDenom: controls.InjectiveDenom}, stub: stub}
	holders := []sdk.AccAddress{sdk.AccAddress("holder-a"), sdk.AccAddress("holder-b")}
	bank.set(holders[0], controls.LogicalDenom, math.NewInt(propInitialNoble/2))
	bank.set(holders[1], controls.LogicalDenom, math.NewInt(propInitialNoble/2))
	return &propChain{
		r: roundTrip{
			ctx:        ctx,
			keeper:     k,
			bank:       bank,
			ratelimit:  ratelimit,
			stub:       stub,
			middleware: NewIBCMiddleware(k, inner),
			controls:   controls,
		},
		gov:     gov,
		holders: holders,
	}
}

// propKeeper is setupKeeperWithICS4 with governance authorities, so the msg
// server can change modes. Channel and connection keepers are nil, which
// skips route-binding checks.
func propKeeper(t *testing.T, authorities []string) (sdk.Context, *keeper.Keeper, *receiveBank, *undoRatelimit) {
	t.Helper()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	bank := newReceiveBank()
	ratelimit := &undoRatelimit{}
	k := keeper.NewKeeper(cdc, key, nil, bank, &sequencedICS4{}, ratelimit, nil, nil, authorities)
	return ctx, k, bank, ratelimit
}

type propBankSnapshot map[string]map[string]math.Int

func (p *propChain) snapshot() (propBankSnapshot, undoRatelimit, int) {
	copied := propBankSnapshot{}
	for address, denoms := range p.r.bank.balances {
		copied[address] = map[string]math.Int{}
		for denom, amount := range denoms {
			copied[address][denom] = amount
		}
	}
	return copied, *p.r.ratelimit, len(p.r.stub.packets)
}

func (p *propChain) restore(bank propBankSnapshot, ratelimit undoRatelimit, packets int) {
	p.r.bank.balances = bank
	*p.r.ratelimit = ratelimit
	p.r.stub.packets = p.r.stub.packets[:packets]
}

// atomic runs op against a cached context and commits it only on success.
func (p *propChain) atomic(op func(ctx sdk.Context) error) error {
	bank, ratelimit, packets := p.snapshot()
	cacheCtx, write := p.r.ctx.CacheContext()
	if err := op(cacheCtx); err != nil {
		p.restore(bank, ratelimit, packets)
		return err
	}
	write()
	return nil
}

func (p *propChain) send(holder sdk.AccAddress, channel string, amount int64) error {
	return p.atomic(func(ctx sdk.Context) error {
		_, err := keeper.NewTransferDecorator(p.r.keeper).Transfer(ctx, &transfertypes.MsgTransfer{
			SourcePort:       transfertypes.PortID,
			SourceChannel:    channel,
			Token:            sdk.NewInt64Coin(p.r.controls.LogicalDenom, amount),
			Sender:           holder.String(),
			Receiver:         "remote1receiver",
			TimeoutTimestamp: 1,
		})
		if err == nil {
			p.inflight = append(p.inflight, p.r.stub.packets[len(p.r.stub.packets)-1])
		}
		return err
	})
}

func (p *propChain) recv(channel, denom string, receiver sdk.AccAddress, amount int64) error {
	data := transfertypes.NewFungibleTokenPacketData(denom, fmt.Sprint(amount), "remote1sender", receiver.String(), "")
	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         transfertypes.PortID,
		SourceChannel:      "channel-7",
		DestinationPort:    transfertypes.PortID,
		DestinationChannel: channel,
		Data:               data.GetBytes(),
	}
	return p.atomic(func(ctx sdk.Context) error {
		ack := p.r.middleware.OnRecvPacket(ctx, packet, nil)
		if ack != nil && !ack.Success() {
			return fmt.Errorf("error ack: %s", ack.Acknowledgement())
		}
		return nil
	})
}

// settle resolves an in-flight packet with a success ack, error ack or
// timeout. Packets are removed only when the callback commits, matching IBC
// core, which keeps the commitment when the callback fails.
func (p *propChain) settle(index int, outcome int) error {
	packet := p.inflight[index]
	err := p.atomic(func(ctx sdk.Context) error {
		switch outcome {
		case 0:
			return p.r.middleware.OnAcknowledgementPacket(ctx, packet, successAck(), nil)
		case 1:
			return p.r.middleware.OnAcknowledgementPacket(ctx, packet, errorAck(), nil)
		default:
			return p.r.middleware.OnTimeoutPacket(ctx, packet, nil)
		}
	})
	if err == nil {
		p.inflight = append(p.inflight[:index], p.inflight[index+1:]...)
	}
	return err
}

func (p *propChain) setMode(mode types.Mode) error {
	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.Mode = mode
	return p.updateControls(controls)
}

func (p *propChain) updateControls(controls types.Controls) error {
	return p.atomic(func(ctx sdk.Context) error {
		_, err := keeper.NewMsgServer(p.r.keeper).UpdateControls(ctx, &types.MsgUpdateControls{
			Authority: p.gov,
			Controls:  controls,
		})
		return err
	})
}

func (p *propChain) ledger() (noble, injective, pendingInjective, pendingNoble math.Int) {
	ledger := p.r.keeper.GetLedger(p.r.ctx)
	parse := func(value dtypes.SerializableInt) math.Int {
		amount, err := types.AmountOf(value)
		if err != nil {
			panic(err)
		}
		return amount
	}
	return parse(ledger.NobleBacking), parse(ledger.InjectiveBacking),
		parse(ledger.PendingInjective), parse(ledger.PendingNoble)
}

// requireModel checks the module invariant plus exact equalities the model
// guarantees because nothing outside the module touches its balances.
func (p *propChain) requireModel(t *testing.T, step string) {
	t.Helper()
	reason, broken := keeper.BackingInvariant(*p.r.keeper)(p.r.ctx)
	require.False(t, broken, "%s: %s", step, reason)
	noble, injective, pendingInjective, pendingNoble := p.ledger()
	controls := p.r.keeper.GetControls(p.r.ctx)
	require.Equal(t, injective.String(),
		p.r.bank.amount(types.ModuleAddress, controls.InjectiveDenom).String(), "%s: physical reserve", step)
	require.Equal(t, pendingInjective.String(),
		p.r.bank.amount(types.ModuleAddress, controls.LogicalDenom).String(), "%s: locked logical", step)
	supply := p.r.bank.GetSupply(p.r.ctx, controls.LogicalDenom).Amount
	require.Equal(t, supply.String(), noble.Add(injective).Add(pendingInjective).String(), "%s: supply", step)

	inflightNoble, inflightInjective := math.ZeroInt(), math.ZeroInt()
	for _, packet := range p.inflight {
		pending, found := p.r.keeper.GetPendingForPacket(p.r.ctx, packet.SourceChannel, packet.Sequence)
		require.True(t, found, "%s: in-flight packet %d has no pending", step, packet.Sequence)
		amount, err := types.PositiveAmountOf(pending.Amount)
		require.NoError(t, err)
		if pending.Route == types.Route_ROUTE_NOBLE {
			inflightNoble = inflightNoble.Add(amount)
		} else {
			inflightInjective = inflightInjective.Add(amount)
		}
	}
	require.Equal(t, inflightNoble.String(), pendingNoble.String(), "%s: pending noble", step)
	require.Equal(t, inflightInjective.String(), pendingInjective.String(), "%s: pending injective", step)
	require.Equal(t, uint32(len(p.inflight)), p.r.keeper.PendingCount(p.r.ctx), "%s: pending count", step)
}

// TestPropLedgerRandomOperations drives random deposits, withdrawals,
// settlements and mode changes and checks the model after every step.
func TestPropLedgerRandomOperations(t *testing.T) {
	for seed := int64(0); seed < propSeeds; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			p := propSetup(t)
			p.requireModel(t, "setup")
			for step := 0; step < propStepsPerSeed; step++ {
				label := propStep(t, p, rng)
				p.requireModel(t, fmt.Sprintf("step %d (%s)", step, label))
			}
		})
	}
}

func propStep(t *testing.T, p *propChain, rng *rand.Rand) string {
	t.Helper()
	holder := p.holders[rng.Intn(len(p.holders))]
	amount := 1 + rng.Int63n(propMaxStepAmount)
	controls := p.r.keeper.GetControls(p.r.ctx)
	switch op := rng.Intn(10); {
	case op < 2:
		err := p.recv(controls.InjectiveChannel, controls.InjectivePacketDenom, holder, amount)
		return fmt.Sprintf("injective deposit %d err=%v", amount, err)
	case op < 4:
		err := p.send(holder, controls.InjectiveChannel, amount)
		return fmt.Sprintf("injective withdraw %d err=%v", amount, err)
	case op < 5:
		err := p.send(holder, controls.NobleChannel, amount)
		return fmt.Sprintf("noble withdraw %d err=%v", amount, err)
	case op < 6:
		err := p.recv(controls.NobleChannel, controls.NoblePacketDenom, holder, amount)
		require.Error(t, err, "noble deposit must be refused outside DISABLED")
		return "noble deposit refused"
	case op < 9 && len(p.inflight) > 0:
		index, outcome := rng.Intn(len(p.inflight)), rng.Intn(3)
		require.NoError(t, p.settle(index, outcome), "settlement must not fail")
		return fmt.Sprintf("settle %d outcome %d", index, outcome)
	default:
		modes := []types.Mode{types.Mode_MODE_GRADUAL, types.Mode_MODE_PAUSED}
		mode := modes[rng.Intn(len(modes))]
		err := p.setMode(mode)
		return fmt.Sprintf("mode %s err=%v", mode, err)
	}
}

// TestPropDisabledEraNobleRefundRestoresBacking: a Noble withdrawal sent while
// DISABLED has no pending record. If it fails after GRADUAL is entered, ICS20
// re-mints the logical voucher, which the snapshot excluded, so it must be
// credited back to noble_backing.
func TestPropDisabledEraNobleRefundRestoresBacking(t *testing.T) {
	p := propSetup(t)
	require.NoError(t, p.setMode(types.Mode_MODE_DISABLED))
	require.NoError(t, p.send(p.holders[0], p.r.controls.NobleChannel, 100))
	inflight := p.inflight[0]
	p.inflight = nil
	require.NoError(t, p.setMode(types.Mode_MODE_GRADUAL))

	require.NoError(t, p.atomic(func(ctx sdk.Context) error {
		return p.r.middleware.OnTimeoutPacket(ctx, inflight, nil)
	}))
	reason, broken := keeper.BackingInvariant(*p.r.keeper)(p.r.ctx)
	require.False(t, broken, "refund of a DISABLED-era Noble packet breaks backing: %s", reason)
	require.Equal(t, int64(propInitialNoble/2), p.r.bank.amount(p.holders[0], p.r.controls.LogicalDenom).Int64())

	settlements := propEvents(p.r.ctx, types.EventTypeSettlement)
	require.Len(t, settlements, 1)
	require.Equal(t, map[string]string{
		types.AttributeKeyRoute:          types.Route_ROUTE_NOBLE.String(),
		types.AttributeKeySourceChannel:  inflight.SourceChannel,
		types.AttributeKeySequence:       fmt.Sprint(inflight.Sequence),
		transfertypes.AttributeKeyAmount: "100",
		sdk.AttributeKeySender:           p.holders[0].String(),
		types.AttributeKeyOutcome:        types.AttributeValueLegacyRefund,
	}, settlements[0])
}

// TestPropGradualNobleFailureRefundsSender: a Noble withdrawal sent in GRADUAL
// that fails returns the logical USDC to the sender and restores Noble backing.
func TestPropGradualNobleFailureRefundsSender(t *testing.T) {
	tests := map[string]struct {
		outcome     int
		wantOutcome string
	}{
		"error ack": {outcome: migrationAckError, wantOutcome: types.AttributeValueError},
		"timeout":   {outcome: migrationTimeout, wantOutcome: types.AttributeValueTimeout},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := propSetup(t)
			holder := p.holders[0]
			before := p.r.bank.amount(holder, p.r.controls.LogicalDenom)
			require.NoError(t, p.send(holder, p.r.controls.NobleChannel, 100))
			packet := p.inflight[0]
			noble, _, _, pendingNoble := p.ledger()
			require.Equal(t, before.SubRaw(100), p.r.bank.amount(holder, p.r.controls.LogicalDenom))
			require.Equal(t, int64(propInitialNoble-100), noble.Int64())
			require.Equal(t, int64(100), pendingNoble.Int64())
			p.requireModel(t, "noble withdrawal in flight")

			require.NoError(t, p.settle(0, tc.outcome))
			noble, _, _, pendingNoble = p.ledger()
			require.Equal(t, before, p.r.bank.amount(holder, p.r.controls.LogicalDenom), "sender refunded")
			require.Equal(t, int64(propInitialNoble), noble.Int64())
			require.True(t, pendingNoble.IsZero())
			p.requireModel(t, "noble withdrawal refunded")

			settlements := propEvents(p.r.ctx, types.EventTypeSettlement)
			require.Len(t, settlements, 1)
			require.Equal(t, map[string]string{
				types.AttributeKeyRoute:          types.Route_ROUTE_NOBLE.String(),
				types.AttributeKeySourceChannel:  packet.SourceChannel,
				types.AttributeKeySequence:       fmt.Sprint(packet.Sequence),
				transfertypes.AttributeKeyAmount: "100",
				sdk.AttributeKeySender:           holder.String(),
				types.AttributeKeyOutcome:        tc.wantOutcome,
			}, settlements[0])
		})
	}
}

// TestPropInjectiveDepositsDisabledKeepsExitsOpen: turning Injective deposits
// off refuses new inflow but leaves withdrawals and their refunds working.
func TestPropInjectiveDepositsDisabledKeepsExitsOpen(t *testing.T) {
	p := propSetup(t)
	holder := p.holders[0]
	injectiveChannel := p.r.controls.InjectiveChannel
	require.NoError(t, p.recv(injectiveChannel, p.r.controls.InjectivePacketDenom, holder, 200))

	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.InjectiveDepositsDisabled = true
	require.NoError(t, p.updateControls(controls))
	_, before, _, _ := p.ledger()
	balance := p.r.bank.amount(holder, p.r.controls.LogicalDenom)
	require.ErrorContains(t, p.recv(injectiveChannel, p.r.controls.InjectivePacketDenom, holder, 10),
		fmt.Sprintf("ABCI code: %d:", types.ErrInjectiveDepositsDisabled.ABCICode()))
	_, injective, _, _ := p.ledger()
	require.Equal(t, before, injective)
	require.Equal(t, balance, p.r.bank.amount(holder, p.r.controls.LogicalDenom))
	p.requireModel(t, "deposit refused")

	require.NoError(t, p.send(holder, injectiveChannel, 100))
	require.NoError(t, p.settle(0, migrationAckError))
	_, injective, _, _ = p.ledger()
	require.Equal(t, before, injective, "refund of a withdrawal still restores backing")
	require.Equal(t, balance, p.r.bank.amount(holder, p.r.controls.LogicalDenom))
	p.requireModel(t, "withdrawal refunded")

	require.NoError(t, p.send(holder, injectiveChannel, 50))
	require.NoError(t, p.settle(0, migrationAckSuccess))
	_, injective, _, _ = p.ledger()
	require.Equal(t, before.SubRaw(50), injective)
	p.requireModel(t, "withdrawal delivered")

	controls.InjectiveDepositsDisabled = false
	require.NoError(t, p.updateControls(controls))
	require.NoError(t, p.recv(injectiveChannel, p.r.controls.InjectivePacketDenom, holder, 10))
	p.requireModel(t, "deposits re-enabled")
}

// TestPropAckWithoutErrorResponseIsDelivered: ICS20 refunds only on an error
// response and treats any other decodable ack, such as "{}", as delivered. The
// ledger must agree on every path or the supply equation breaks.
func TestPropAckWithoutErrorResponseIsDelivered(t *testing.T) {
	emptyAck := transfertypes.ModuleCdc.MustMarshalJSON(&channeltypes.Acknowledgement{})
	t.Run("pending noble", func(t *testing.T) {
		p := propSetup(t)
		require.NoError(t, p.send(p.holders[0], p.r.controls.NobleChannel, 100))
		require.NoError(t, p.ackInflight(0, emptyAck))
		noble, _, _, _ := p.ledger()
		require.Equal(t, int64(propInitialNoble-100), noble.Int64())
		require.Equal(t, int64(propInitialNoble/2-100), p.r.bank.amount(p.holders[0], p.r.controls.LogicalDenom).Int64())
		p.requireModel(t, "empty ack on noble")
	})
	t.Run("pending injective", func(t *testing.T) {
		p := propSetup(t)
		require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, p.holders[0], 200))
		require.NoError(t, p.send(p.holders[0], p.r.controls.InjectiveChannel, 100))
		require.NoError(t, p.ackInflight(0, emptyAck))
		_, injective, _, _ := p.ledger()
		require.Equal(t, int64(100), injective.Int64())
		require.Zero(t, p.r.ratelimit.undone)
		p.requireModel(t, "empty ack on injective")
	})
	t.Run("legacy noble", func(t *testing.T) {
		p := propSetup(t)
		require.NoError(t, p.setMode(types.Mode_MODE_DISABLED))
		require.NoError(t, p.send(p.holders[0], p.r.controls.NobleChannel, 100))
		require.NoError(t, p.setMode(types.Mode_MODE_GRADUAL))
		before, _, _, _ := p.ledger()
		require.NoError(t, p.ackInflight(0, emptyAck))
		noble, _, _, _ := p.ledger()
		require.Equal(t, before, noble, "no refund, so no credit")
		require.Empty(t, propEvents(p.r.ctx, types.EventTypeSettlement))
		p.requireModel(t, "empty ack on legacy packet")
	})
}

func (p *propChain) ackInflight(index int, acknowledgement []byte) error {
	packet := p.inflight[index]
	err := p.atomic(func(ctx sdk.Context) error {
		return p.r.middleware.OnAcknowledgementPacket(ctx, packet, acknowledgement, nil)
	})
	if err == nil {
		p.inflight = append(p.inflight[:index], p.inflight[index+1:]...)
	}
	return err
}

func propEvents(ctx sdk.Context, eventType string) []map[string]string {
	var found []map[string]string
	for _, event := range ctx.EventManager().Events() {
		if event.Type != eventType {
			continue
		}
		attributes := map[string]string{}
		for _, attribute := range event.Attributes {
			attributes[attribute.Key] = attribute.Value
		}
		found = append(found, attributes)
	}
	return found
}

// TestPropDisableRejectedWhileInjectiveBacked: DISABLED would let logical USDC
// minted against Injective backing leave over Noble, pushing supply below
// injective_backing so the snapshot rejects GRADUAL forever. Governance must use
// PAUSED instead, which keeps GRADUAL reachable.
func TestPropDisableRejectedWhileInjectiveBacked(t *testing.T) {
	p := propSetup(t)
	require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, p.holders[0], 200))
	require.ErrorIs(t, p.setMode(types.Mode_MODE_DISABLED), types.ErrInvalidControls)
	require.Equal(t, types.Mode_MODE_GRADUAL, p.r.keeper.GetControls(p.r.ctx).Mode)
	require.NoError(t, p.setMode(types.Mode_MODE_PAUSED))
	require.NoError(t, p.setMode(types.Mode_MODE_GRADUAL))
	p.requireModel(t, "after pause and resume")
}

// TestPropRouteRebindRejectedWhileInjectiveBacked: GRADUAL to GRADUAL skips the
// snapshot and invariant check, so rebinding the Injective route while
// injective_backing is held in the old physical denom must be rejected.
func TestPropRouteRebindRejectedWhileInjectiveBacked(t *testing.T) {
	p := propSetup(t)
	require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, p.holders[0], 200))
	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.InjectiveChannel = "channel-5"
	controls.InjectiveDenom = transfertypes.ParseDenomTrace("transfer/channel-5/uusdc").IBCDenom()
	err := p.updateControls(controls)
	if err == nil {
		reason, broken := keeper.BackingInvariant(*p.r.keeper)(p.r.ctx)
		require.False(t, broken, "route rebind with outstanding injective backing accepted: %s", reason)
	}
}

// TestPropMaxPendingCannotDropBelowLiveCount: a max_pending_settlements below
// the live count would make ExportGenesis produce state that
// GenesisState.Validate (and InitGenesis) reject.
func TestPropMaxPendingCannotDropBelowLiveCount(t *testing.T) {
	p := propSetup(t)
	for i := 0; i < 3; i++ {
		require.NoError(t, p.send(p.holders[0], p.r.controls.NobleChannel, 10))
	}
	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.MaxPendingSettlements = 1
	if err := p.updateControls(controls); err != nil {
		return
	}
	require.NoError(t, ExportGenesis(p.r.ctx, *p.r.keeper).Validate(), "exported genesis must stay valid")
}

// TestPropInjectiveCreditOverflowIsAnErrorAck: with 256-bit governance limits,
// crediting backing overflows math.Int. The middleware should return an error
// ack instead of panicking inside OnRecvPacket.
func TestPropInjectiveCreditOverflowIsAnErrorAck(t *testing.T) {
	p := propSetup(t)
	maxInt := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), math.MaxBitLen), big.NewInt(1))
	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.MaxTransferAmount = dtypes.NewIntFromBigInt(maxInt)
	controls.MigrationCeiling = dtypes.NewIntFromBigInt(maxInt)
	require.NoError(t, p.r.keeper.SetControls(p.r.ctx, controls))
	ledger := p.r.keeper.GetLedger(p.r.ctx)
	ledger.InjectiveBacking = dtypes.NewIntFromBigInt(new(big.Int).Sub(maxInt, big.NewInt(1)))
	require.NoError(t, p.r.keeper.SetLedger(p.r.ctx, ledger))

	data := transfertypes.NewFungibleTokenPacketData(controls.InjectivePacketDenom, "5", "remote1sender",
		p.holders[0].String(), "")
	packet := channeltypes.Packet{
		Sequence: 1, SourcePort: transfertypes.PortID, SourceChannel: "channel-7",
		DestinationPort: transfertypes.PortID, DestinationChannel: controls.InjectiveChannel, Data: data.GetBytes(),
	}
	var ack ibcexported.Acknowledgement
	require.NotPanics(t, func() { ack = p.r.middleware.OnRecvPacket(p.r.ctx, packet, nil) })
	require.False(t, ack.Success())
}
