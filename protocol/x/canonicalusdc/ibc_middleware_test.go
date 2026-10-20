package canonicalusdc

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type receiveBank struct {
	balances map[string]map[string]math.Int
}

func newReceiveBank() *receiveBank {
	return &receiveBank{balances: make(map[string]map[string]math.Int)}
}

func (b *receiveBank) amount(address sdk.AccAddress, denom string) math.Int {
	if denoms, ok := b.balances[address.String()]; ok {
		if amount, ok := denoms[denom]; ok {
			return amount
		}
	}
	return math.ZeroInt()
}

func (b *receiveBank) set(address sdk.AccAddress, denom string, amount math.Int) {
	if _, ok := b.balances[address.String()]; !ok {
		b.balances[address.String()] = make(map[string]math.Int)
	}
	b.balances[address.String()][denom] = amount
}

func (b *receiveBank) GetBalance(_ context.Context, address sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, b.amount(address, denom))
}

func (b *receiveBank) GetSupply(_ context.Context, denom string) sdk.Coin {
	total := math.ZeroInt()
	for _, denoms := range b.balances {
		if amount, ok := denoms[denom]; ok {
			total = total.Add(amount)
		}
	}
	return sdk.NewCoin(denom, total)
}

func (b *receiveBank) SendCoinsFromAccountToModule(
	_ context.Context,
	sender sdk.AccAddress,
	module string,
	coins sdk.Coins,
) error {
	if module != types.ModuleName {
		return fmt.Errorf("unexpected module")
	}
	return b.move(sender, types.ModuleAddress, coins)
}

func (b *receiveBank) move(from, to sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		balance := b.amount(from, coin.Denom)
		if balance.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds")
		}
		b.set(from, coin.Denom, balance.Sub(coin.Amount))
		b.set(to, coin.Denom, b.amount(to, coin.Denom).Add(coin.Amount))
	}
	return nil
}

func (b *receiveBank) SendCoinsFromModuleToAccount(
	_ context.Context,
	module string,
	receiver sdk.AccAddress,
	coins sdk.Coins,
) error {
	if module != types.ModuleName {
		return fmt.Errorf("unexpected module")
	}
	for _, coin := range coins {
		balance := b.amount(types.ModuleAddress, coin.Denom)
		if balance.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds")
		}
		b.set(types.ModuleAddress, coin.Denom, balance.Sub(coin.Amount))
		b.set(receiver, coin.Denom, b.amount(receiver, coin.Denom).Add(coin.Amount))
	}
	return nil
}

func (b *receiveBank) MintCoins(_ context.Context, module string, coins sdk.Coins) error {
	if module != types.ModuleName {
		return fmt.Errorf("unexpected module")
	}
	for _, coin := range coins {
		b.set(types.ModuleAddress, coin.Denom, b.amount(types.ModuleAddress, coin.Denom).Add(coin.Amount))
	}
	return nil
}

func (b *receiveBank) BurnCoins(_ context.Context, module string, coins sdk.Coins) error {
	if module != types.ModuleName {
		return fmt.Errorf("unexpected module")
	}
	for _, coin := range coins {
		balance := b.amount(types.ModuleAddress, coin.Denom)
		if balance.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds")
		}
		b.set(types.ModuleAddress, coin.Denom, balance.Sub(coin.Amount))
	}
	return nil
}

type receiveApp struct {
	bank             *receiveBank
	physicalDenom    string
	calls            int
	ackCalls         int
	timeoutCalls     int
	data             transfertypes.FungibleTokenPacketData
	expectedReceiver string
	timeoutRefund    *sdk.Coin
	timeoutReceiver  sdk.AccAddress
	ackRefund        *sdk.Coin
	failRecv         bool
	failAck          error
	failTimeout      error
}

func (a *receiveApp) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	a.calls++
	if a.failRecv {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("inner recv failed"))
	}
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &a.data); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	expectedReceiver := a.expectedReceiver
	if expectedReceiver == "" {
		expectedReceiver = types.ModuleAddress.String()
	}
	if a.data.Receiver != expectedReceiver {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("physical receiver was not reserve"))
	}
	amount, _ := math.NewIntFromString(a.data.Amount)
	_ = a.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(a.physicalDenom, amount)))
	return channeltypes.NewResultAcknowledgement([]byte{1})
}

func (a *receiveApp) OnAcknowledgementPacket(sdk.Context, channeltypes.Packet, []byte, sdk.AccAddress) error {
	a.ackCalls++
	if a.failAck != nil {
		return a.failAck
	}
	if a.ackRefund != nil {
		a.bank.set(
			types.ModuleAddress,
			a.ackRefund.Denom,
			a.bank.amount(types.ModuleAddress, a.ackRefund.Denom).Add(a.ackRefund.Amount),
		)
	}
	return nil
}
func (a *receiveApp) OnTimeoutPacket(sdk.Context, channeltypes.Packet, sdk.AccAddress) error {
	a.timeoutCalls++
	if a.failTimeout != nil {
		return a.failTimeout
	}
	if a.timeoutRefund != nil {
		a.bank.set(
			a.timeoutReceiver,
			a.timeoutRefund.Denom,
			a.bank.amount(a.timeoutReceiver, a.timeoutRefund.Denom).Add(a.timeoutRefund.Amount),
		)
	}
	return nil
}
func (*receiveApp) OnChanOpenInit(
	sdk.Context,
	channeltypes.Order,
	[]string,
	string,
	string,
	*capabilitytypes.Capability,
	channeltypes.Counterparty,
	string,
) (string, error) {
	return transfertypes.Version, nil
}
func (*receiveApp) OnChanOpenTry(
	sdk.Context,
	channeltypes.Order,
	[]string,
	string,
	string,
	*capabilitytypes.Capability,
	channeltypes.Counterparty,
	string,
) (string, error) {
	return transfertypes.Version, nil
}
func (*receiveApp) OnChanOpenAck(sdk.Context, string, string, string, string) error { return nil }
func (*receiveApp) OnChanOpenConfirm(sdk.Context, string, string) error             { return nil }
func (*receiveApp) OnChanCloseInit(sdk.Context, string, string) error               { return nil }
func (*receiveApp) OnChanCloseConfirm(sdk.Context, string, string) error            { return nil }

func setupReceiveKeeper(t *testing.T) (sdk.Context, *keeper.Keeper, *receiveBank) {
	t.Helper()
	ctx, k, bank, _ := setupReceiveKeeperWithRatelimit(t)
	return ctx, k, bank
}

func setupReceiveKeeperWithRatelimit(t *testing.T) (sdk.Context, *keeper.Keeper, *receiveBank, *undoRatelimit) {
	t.Helper()
	return setupKeeperWithICS4(t, nil)
}

func setupKeeperWithICS4(
	t *testing.T,
	ics4 types.ICS4Wrapper,
) (sdk.Context, *keeper.Keeper, *receiveBank, *undoRatelimit) {
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
	return ctx, keeper.NewKeeper(cdc, key, nil, bank, ics4, ratelimit, nil, nil, nil), bank, ratelimit
}

// undoRatelimit records refunded capacity. Charging happens in the transfer
// decorator, which these middleware tests do not exercise.
type undoRatelimit struct {
	charged int64
	undone  int64
}

func (r *undoRatelimit) ProcessWithdrawal(_ sdk.Context, _ string, amount *big.Int) error {
	r.charged += amount.Int64()
	return nil
}

func (r *undoRatelimit) UndoWithdrawal(_ sdk.Context, _ string, amount *big.Int) {
	r.undone += amount.Int64()
}

func receiveControls() types.Controls {
	physical := transfertypes.ParseDenomTrace("transfer/channel-1/uusdc").IBCDenom()
	return types.Controls{
		Mode:                    types.Mode_MODE_GRADUAL,
		LogicalDenom:            "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
		NobleChannel:            "channel-0",
		InjectiveChannel:        "channel-1",
		NoblePacketDenom:        "uusdc",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          physical,
		NobleWithdrawalsEnabled: true,
		MaxTransferAmount:       dtypes.NewInt(1000),
		MigrationCeiling:        dtypes.NewInt(10000),
		MaxPendingSettlements:   10,
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
	}
}

func TestInboundInjectiveCanonicalizationPreservesLogicalReceiver(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	receiver := sdk.AccAddress("logical-receiver")
	data := transfertypes.FungibleTokenPacketData{
		Denom:    "uusdc",
		Amount:   "25",
		Sender:   "injective-sender",
		Receiver: receiver.String(),
		Memo:     "wallet-memo",
	}
	packet := channeltypes.Packet{
		Sequence:           9,
		SourcePort:         transfertypes.PortID,
		SourceChannel:      "channel-99",
		DestinationPort:    transfertypes.PortID,
		DestinationChannel: controls.InjectiveChannel,
		Data:               data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.InjectiveDenom}
	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.True(t, ack.Success())
	require.Equal(t, 1, app.calls)
	require.Equal(t, types.ModuleAddress.String(), app.data.Receiver)
	require.Equal(t, receiver.String(), data.Receiver, "logical packet data must remain unchanged")
	require.Equal(t, "25", bank.amount(receiver, controls.LogicalDenom).String())
	require.Equal(t, "25", bank.amount(types.ModuleAddress, controls.InjectiveDenom).String())
	require.Equal(t, "25", k.GetLedger(ctx).InjectiveBacking.String())

	var inbound []sdk.Event
	for _, event := range ctx.EventManager().Events() {
		if event.Type == types.EventTypeInbound {
			inbound = append(inbound, event)
		}
	}
	require.Len(t, inbound, 1)
	attributes := map[string]string{}
	for _, attribute := range inbound[0].Attributes {
		attributes[attribute.Key] = attribute.Value
	}
	require.Equal(t, map[string]string{
		types.AttributeKeyPhysicalDenom:    controls.InjectiveDenom,
		types.AttributeKeyLogicalDenom:     controls.LogicalDenom,
		transfertypes.AttributeKeyAmount:   "25",
		transfertypes.AttributeKeyReceiver: receiver.String(),
		sdk.AttributeKeySender:             "injective-sender",
		types.AttributeKeySourceChannel:    "channel-99",
		types.AttributeKeySequence:         "9",
	}, attributes)
}

func TestInboundRecvRejectsBeforeTokenMovement(t *testing.T) {
	cases := []struct {
		name          string
		prepControls  func(types.Controls) types.Controls
		prepLedger    func(types.Ledger) types.Ledger
		denom         string
		amount        string
		sender        string
		sourceChannel string
		destChannel   string
		sequence      uint64
		injectivePhys bool
		wantErr       error
	}{
		{
			name:          "noble deposit",
			denom:         "uusdc",
			amount:        "25",
			sender:        "noble-sender",
			sourceChannel: "channel-88",
			destChannel:   receiveControls().NobleChannel,
			sequence:      1,
			wantErr:       types.ErrNobleDepositsDisabled,
		},
		{
			name: "injective transfer limit",
			prepControls: func(controls types.Controls) types.Controls {
				controls.MaxTransferAmount = dtypes.NewInt(24)
				return controls
			},
			denom:         "uusdc",
			amount:        "25",
			sender:        "injective-sender",
			sourceChannel: "channel-99",
			destChannel:   receiveControls().InjectiveChannel,
			sequence:      2,
			injectivePhys: true,
			wantErr:       types.ErrInvalidAmount,
		},
		{
			name: "migration ceiling",
			prepLedger: func(ledger types.Ledger) types.Ledger {
				ledger.InjectiveBacking = dtypes.NewInt(9990)
				return ledger
			},
			denom:         "uusdc",
			amount:        "25",
			sender:        "injective-sender",
			sourceChannel: "channel-99",
			destChannel:   receiveControls().InjectiveChannel,
			sequence:      3,
			injectivePhys: true,
			wantErr:       types.ErrMigrationCeiling,
		},
		{
			name: "paused injective deposit",
			prepControls: func(controls types.Controls) types.Controls {
				controls.Mode = types.Mode_MODE_PAUSED
				return controls
			},
			denom:         "uusdc",
			amount:        "25",
			sender:        "injective-sender",
			sourceChannel: "channel-99",
			destChannel:   receiveControls().InjectiveChannel,
			sequence:      2,
			injectivePhys: true,
			wantErr:       types.ErrPaused,
		},
		{
			name: "injective deposits disabled",
			prepControls: func(controls types.Controls) types.Controls {
				controls.InjectiveDepositsDisabled = true
				return controls
			},
			denom:         "uusdc",
			amount:        "25",
			sender:        "injective-sender",
			sourceChannel: "channel-99",
			destChannel:   receiveControls().InjectiveChannel,
			sequence:      4,
			injectivePhys: true,
			wantErr:       types.ErrInjectiveDepositsDisabled,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, bank := setupReceiveKeeper(t)
			controls := receiveControls()
			if tc.prepControls != nil {
				controls = tc.prepControls(controls)
			}
			require.NoError(t, k.SetControls(ctx, controls))
			ledger := types.DefaultLedger()
			if tc.prepLedger != nil {
				ledger = tc.prepLedger(ledger)
			}
			require.NoError(t, k.SetLedger(ctx, ledger))
			receiver := sdk.AccAddress("receiver")
			data := transfertypes.FungibleTokenPacketData{
				Denom: tc.denom, Amount: tc.amount, Sender: tc.sender, Receiver: receiver.String(),
			}
			packet := channeltypes.Packet{
				Sequence: tc.sequence, SourcePort: transfertypes.PortID, SourceChannel: tc.sourceChannel,
				DestinationPort: transfertypes.PortID, DestinationChannel: tc.destChannel, Data: data.GetBytes(),
			}
			physical := controls.LogicalDenom
			if tc.injectivePhys {
				physical = controls.InjectiveDenom
			}
			app := &receiveApp{bank: bank, physicalDenom: physical}
			ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
			require.False(t, ack.Success())
			require.Equal(t, channeltypes.NewErrorAcknowledgement(tc.wantErr), ack)
			require.Zero(t, app.calls)
			require.Equal(t, "0", bank.amount(receiver, controls.LogicalDenom).String())
			require.Equal(t, ledger, k.GetLedger(ctx))
		})
	}
}

func TestLogicalReturnOnOtherChannelPassesThroughWithoutLedgerChange(t *testing.T) {
	for _, failRecv := range []bool{false, true} {
		t.Run(fmt.Sprintf("failRecv=%t", failRecv), func(t *testing.T) {
			ctx, k, bank := setupReceiveKeeper(t)
			controls := receiveControls()
			require.NoError(t, k.SetControls(ctx, controls))
			ledger := types.DefaultLedger()
			ledger.NobleBacking = dtypes.NewInt(7)
			require.NoError(t, k.SetLedger(ctx, ledger))
			receiver := sdk.AccAddress("receiver").String()
			data := transfertypes.FungibleTokenPacketData{
				Denom:  "transfer/channel-88/transfer/channel-0/uusdc",
				Amount: "1", Sender: "remote-sender", Receiver: receiver,
			}
			packet := channeltypes.Packet{
				Sequence: 4, SourcePort: transfertypes.PortID, SourceChannel: "channel-88",
				DestinationPort: transfertypes.PortID, DestinationChannel: "channel-9", Data: data.GetBytes(),
			}
			app := &receiveApp{
				bank: bank, physicalDenom: controls.LogicalDenom, expectedReceiver: receiver, failRecv: failRecv,
			}

			ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
			require.Equal(t, !failRecv, ack.Success())
			require.Equal(t, 1, app.calls)
			if !failRecv {
				require.Equal(t, receiver, app.data.Receiver)
			}
			require.Equal(t, ledger, k.GetLedger(ctx))
		})
	}
}

func TestNobleTimeoutRefundRestoresBacking(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	ledger := types.DefaultLedger()
	ledger.NobleBacking = dtypes.NewInt(5)
	ledger.PendingNoble = dtypes.NewInt(5)
	require.NoError(t, k.SetLedger(ctx, ledger))
	pending := types.PendingSettlement{
		Route:            types.Route_ROUTE_NOBLE,
		SourceChannel:    controls.NobleChannel,
		Sequence:         17,
		Sender:           sdk.AccAddress("timeout-sender").String(),
		LogicalReceiver:  "noble1receiver",
		PhysicalReceiver: "noble1receiver",
		LogicalDenom:     controls.LogicalDenom,
		PhysicalDenom:    controls.LogicalDenom,
		Amount:           dtypes.NewInt(5),
	}
	require.NoError(t, k.SetPending(ctx, pending))
	sender, err := sdk.AccAddressFromBech32(pending.Sender)
	require.NoError(t, err)
	bank.set(sender, controls.LogicalDenom, math.NewInt(5))
	packet := outboundPacket(pending, controls.NobleTrace(), pending.Sender)
	refund := sdk.NewInt64Coin(controls.LogicalDenom, 5)
	app := &receiveApp{bank: bank, timeoutRefund: &refund, timeoutReceiver: sender}

	require.NoError(t, NewIBCMiddleware(k, app).OnTimeoutPacket(ctx, packet, nil))
	require.Equal(t, 1, app.timeoutCalls)
	require.Zero(t, k.PendingCount(ctx))
	require.Equal(t, "10", k.GetLedger(ctx).NobleBacking.String())
	require.Equal(t, "0", k.GetLedger(ctx).PendingNoble.String())
}

func TestInjectiveAcknowledgementLifecycle(t *testing.T) {
	for _, success := range []bool{true, false} {
		t.Run(fmt.Sprintf("success=%t", success), func(t *testing.T) {
			ctx, k, bank, ratelimit := setupReceiveKeeperWithRatelimit(t)
			controls := receiveControls()
			require.NoError(t, k.SetControls(ctx, controls))
			ledger := types.DefaultLedger()
			ledger.InjectiveBacking = dtypes.NewInt(40)
			ledger.PendingInjective = dtypes.NewInt(10)
			require.NoError(t, k.SetLedger(ctx, ledger))
			sender := sdk.AccAddress("ack-sender")
			bank.set(sender, controls.LogicalDenom, math.NewInt(40))
			bank.set(types.ModuleAddress, controls.LogicalDenom, math.NewInt(10))
			bank.set(types.ModuleAddress, controls.InjectiveDenom, math.NewInt(40))
			pending := types.PendingSettlement{
				Route: types.Route_ROUTE_INJECTIVE, SourceChannel: controls.InjectiveChannel,
				Sequence: 18, Sender: sender.String(), LogicalReceiver: "inj1receiver",
				PhysicalReceiver: "inj1receiver", LogicalDenom: controls.LogicalDenom,
				PhysicalDenom: controls.InjectiveDenom, Amount: dtypes.NewInt(10),
			}
			require.NoError(t, k.SetPending(ctx, pending))
			packet := outboundPacket(pending, controls.InjectiveTrace(), types.ModuleAddress.String())
			app := &receiveApp{bank: bank}
			var acknowledgement []byte
			if success {
				acknowledgement = channeltypes.NewResultAcknowledgement([]byte{1}).Acknowledgement()
			} else {
				refund := sdk.NewInt64Coin(controls.InjectiveDenom, 10)
				app.ackRefund = &refund
				acknowledgement = channeltypes.NewErrorAcknowledgement(fmt.Errorf("remote error")).Acknowledgement()
			}

			require.NoError(t, NewIBCMiddleware(k, app).OnAcknowledgementPacket(ctx, packet, acknowledgement, nil))
			require.Equal(t, 1, app.ackCalls)
			require.Zero(t, k.PendingCount(ctx))
			if success {
				require.Equal(t, "40", k.GetLedger(ctx).InjectiveBacking.String())
				require.Equal(t, "0", bank.amount(types.ModuleAddress, controls.LogicalDenom).String())
				require.Zero(t, ratelimit.undone)
			} else {
				require.Equal(t, "50", k.GetLedger(ctx).InjectiveBacking.String())
				require.Equal(t, "50", bank.amount(sender, controls.LogicalDenom).String())
				require.Equal(t, int64(10), ratelimit.undone)
			}
		})
	}
}

func TestDisabledRecvPassesThroughUnchanged(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	receiver := sdk.AccAddress("logical-receiver").String()
	data := transfertypes.FungibleTokenPacketData{
		Denom: "uusdc", Amount: "25", Sender: "noble-sender", Receiver: receiver,
	}
	packet := channeltypes.Packet{
		Sequence: 1, SourcePort: transfertypes.PortID, SourceChannel: "channel-88",
		DestinationPort: transfertypes.PortID, DestinationChannel: "channel-0", Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: receiveControls().LogicalDenom, expectedReceiver: receiver}
	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.True(t, ack.Success())
	require.Equal(t, 1, app.calls)
	require.Equal(t, receiver, app.data.Receiver)
}

func TestFailedInnerRecvDoesNotMintOrUpdateLedger(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	receiver := sdk.AccAddress("logical-receiver")
	data := transfertypes.FungibleTokenPacketData{
		Denom: "uusdc", Amount: "25", Sender: "injective-sender", Receiver: receiver.String(),
	}
	packet := channeltypes.Packet{
		Sequence: 9, SourcePort: transfertypes.PortID, SourceChannel: "channel-99",
		DestinationPort: transfertypes.PortID, DestinationChannel: controls.InjectiveChannel,
		Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.InjectiveDenom, failRecv: true}
	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.False(t, ack.Success())
	require.Equal(t, 1, app.calls)
	require.Equal(t, "0", bank.amount(receiver, controls.LogicalDenom).String())
	require.Equal(t, "0", k.GetLedger(ctx).InjectiveBacking.String())
}

func TestPendingPacketMismatchLeavesPending(t *testing.T) {
	controls := receiveControls()
	pending := types.PendingSettlement{
		Route: types.Route_ROUTE_NOBLE, SourceChannel: controls.NobleChannel, Sequence: 17,
		Sender: sdk.AccAddress("timeout-sender").String(), LogicalReceiver: "noble1receiver",
		PhysicalReceiver: "noble1receiver", LogicalDenom: controls.LogicalDenom,
		PhysicalDenom: controls.LogicalDenom, Amount: dtypes.NewInt(5),
	}
	cases := []struct {
		name   string
		mutate func(*transfertypes.FungibleTokenPacketData)
	}{
		{name: "amount", mutate: func(data *transfertypes.FungibleTokenPacketData) { data.Amount = "4" }},
		{name: "receiver", mutate: func(data *transfertypes.FungibleTokenPacketData) { data.Receiver = "other" }},
		{name: "sender", mutate: func(data *transfertypes.FungibleTokenPacketData) { data.Sender = "other" }},
		{
			name:   "base denom",
			mutate: func(data *transfertypes.FungibleTokenPacketData) { data.Denom = controls.NoblePacketDenom },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, bank := setupReceiveKeeper(t)
			require.NoError(t, k.SetControls(ctx, controls))
			ledger := types.DefaultLedger()
			ledger.PendingNoble = dtypes.NewInt(5)
			require.NoError(t, k.SetLedger(ctx, ledger))
			require.NoError(t, k.SetPending(ctx, pending))
			data := transfertypes.FungibleTokenPacketData{
				Denom: controls.NobleTrace(), Amount: "5", Sender: pending.Sender, Receiver: pending.PhysicalReceiver,
			}
			tc.mutate(&data)
			packet := channeltypes.Packet{
				Sequence: pending.Sequence, SourcePort: transfertypes.PortID,
				SourceChannel: pending.SourceChannel, Data: data.GetBytes(),
			}
			err := NewIBCMiddleware(k, &receiveApp{bank: bank}).OnTimeoutPacket(ctx, packet, nil)
			require.Error(t, err)
			require.Equal(t, uint32(1), k.PendingCount(ctx))
			require.Equal(t, "5", k.GetLedger(ctx).PendingNoble.String())
		})
	}
}

// outboundPacket builds the packet ICS20 sends for a pending settlement. ICS20
// writes the voucher's full denom trace into outbound packet data.
func outboundPacket(pending types.PendingSettlement, trace, sender string) channeltypes.Packet {
	data := transfertypes.FungibleTokenPacketData{
		Denom: trace, Amount: pending.Amount.String(), Sender: sender, Receiver: pending.PhysicalReceiver,
	}
	return channeltypes.Packet{
		Sequence: pending.Sequence, SourcePort: transfertypes.PortID,
		SourceChannel: pending.SourceChannel, Data: data.GetBytes(),
	}
}
