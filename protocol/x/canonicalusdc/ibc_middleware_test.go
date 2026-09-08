package canonicalusdc

import (
	"context"
	"fmt"
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
	context.Context,
	sdk.AccAddress,
	string,
	sdk.Coins,
) error {
	return fmt.Errorf("unexpected account-to-module send")
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
}

func (a *receiveApp) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	a.calls++
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
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	bank := newReceiveBank()
	return ctx, keeper.NewKeeper(cdc, key, bank, nil, nil, nil, nil), bank
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
		MaxTransferAmount:       "1000",
		MigrationCeiling:        "10000",
		MaxPendingSettlements:   10,
		MemoVersion:             1,
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
	require.Equal(t, "25", k.GetLedger(ctx).InjectiveBacking)
}

func TestInboundNobleDepositRejectedBeforeTokenMovement(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	data := transfertypes.FungibleTokenPacketData{
		Denom: "uusdc", Amount: "25", Sender: "noble-sender", Receiver: sdk.AccAddress("receiver").String(),
	}
	packet := channeltypes.Packet{
		Sequence: 1, SourcePort: transfertypes.PortID, SourceChannel: "channel-88",
		DestinationPort: transfertypes.PortID, DestinationChannel: controls.NobleChannel, Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.LogicalDenom}
	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.False(t, ack.Success())
	require.Zero(t, app.calls)
}

func TestInboundInjectiveTransferLimitRejectsBeforeTokenMovement(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	controls.MaxTransferAmount = "24"
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	data := transfertypes.FungibleTokenPacketData{
		Denom: "uusdc", Amount: "25", Sender: "injective-sender", Receiver: sdk.AccAddress("receiver").String(),
	}
	packet := channeltypes.Packet{
		Sequence: 2, SourcePort: transfertypes.PortID, SourceChannel: "channel-99",
		DestinationPort: transfertypes.PortID, DestinationChannel: controls.InjectiveChannel, Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.InjectiveDenom}

	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.False(t, ack.Success())
	require.Zero(t, app.calls)
}

func TestUnclassifiedLegacyReturnRejectsBeforeTokenMovement(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	data := transfertypes.FungibleTokenPacketData{
		Denom:  "transfer/channel-88/transfer/channel-0/uusdc",
		Amount: "1", Sender: "legacy-sender", Receiver: sdk.AccAddress("receiver").String(),
	}
	packet := channeltypes.Packet{
		Sequence: 3, SourcePort: transfertypes.PortID, SourceChannel: "channel-88",
		DestinationPort: transfertypes.PortID, DestinationChannel: "channel-9", Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.LogicalDenom}

	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.False(t, ack.Success())
	require.Zero(t, app.calls)
}

func TestClassifiedLegacyReturnDecrementsOutstandingAmount(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	ledger := types.DefaultLedger()
	ledger.LegacyDownstream = "2"
	require.NoError(t, k.SetLedger(ctx, ledger))
	receiver := sdk.AccAddress("receiver").String()
	data := transfertypes.FungibleTokenPacketData{
		Denom:  "transfer/channel-88/transfer/channel-0/uusdc",
		Amount: "1", Sender: "legacy-sender", Receiver: receiver,
	}
	packet := channeltypes.Packet{
		Sequence: 4, SourcePort: transfertypes.PortID, SourceChannel: "channel-88",
		DestinationPort: transfertypes.PortID, DestinationChannel: "channel-9", Data: data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.LogicalDenom, expectedReceiver: receiver}

	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.True(t, ack.Success())
	require.Equal(t, 1, app.calls)
	require.Equal(t, "1", k.GetLedger(ctx).LegacyDownstream)
}

func TestInboundBackingFundingDoesNotMintLogicalUsdc(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	require.NoError(t, k.SetLedger(ctx, types.DefaultLedger()))
	controller := sdk.AccAddress("swap-controller")
	participant := types.Participant{
		Controller:     controller.String(),
		NobleRecipient: "noble1receiver",
		MaxRelease:     "100",
		Released:       "0",
		Funded:         "0",
	}
	require.NoError(t, k.SetParticipant(ctx, participant))
	data := transfertypes.FungibleTokenPacketData{
		Denom:    "uusdc",
		Amount:   "25",
		Sender:   "injective-sender",
		Receiver: types.ModuleAddress.String(),
		Memo: fmt.Sprintf(
			`{"canonical_usdc":{"version":1,"action":"fund_backing_swap","controller":%q,"noble_recipient":"noble1receiver"}}`,
			controller.String(),
		),
	}
	packet := channeltypes.Packet{
		Sequence:           10,
		SourcePort:         transfertypes.PortID,
		SourceChannel:      "channel-99",
		DestinationPort:    transfertypes.PortID,
		DestinationChannel: controls.InjectiveChannel,
		Data:               data.GetBytes(),
	}
	app := &receiveApp{bank: bank, physicalDenom: controls.InjectiveDenom}

	ack := NewIBCMiddleware(k, app).OnRecvPacket(ctx, packet, nil)
	require.True(t, ack.Success())
	require.Equal(t, "0", bank.GetSupply(ctx, controls.LogicalDenom).Amount.String())
	require.Equal(t, "25", bank.amount(types.ModuleAddress, controls.InjectiveDenom).String())
	require.Equal(t, "25", k.GetLedger(ctx).RestrictedFunding)
	stored, found := k.GetParticipant(ctx, controller.String())
	require.True(t, found)
	require.Equal(t, "25", stored.Funded)
}

func TestTimeoutCallbackIsIdempotent(t *testing.T) {
	ctx, k, bank := setupReceiveKeeper(t)
	controls := receiveControls()
	require.NoError(t, k.SetControls(ctx, controls))
	ledger := types.DefaultLedger()
	ledger.NobleBacking = "5"
	ledger.PendingNoble = "5"
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
		Amount:           "5",
	}
	require.NoError(t, k.SetPending(ctx, pending))
	sender, err := sdk.AccAddressFromBech32(pending.Sender)
	require.NoError(t, err)
	bank.set(sender, controls.LogicalDenom, math.NewInt(5))
	data := transfertypes.FungibleTokenPacketData{
		Denom: controls.NoblePacketDenom, Amount: pending.Amount, Sender: pending.Sender, Receiver: pending.PhysicalReceiver,
	}
	packet := channeltypes.Packet{
		Sequence:      pending.Sequence,
		SourcePort:    transfertypes.PortID,
		SourceChannel: pending.SourceChannel,
		Data:          data.GetBytes(),
	}
	refund := sdk.NewInt64Coin(controls.LogicalDenom, 5)
	app := &receiveApp{bank: bank, timeoutRefund: &refund, timeoutReceiver: sender}
	middleware := NewIBCMiddleware(k, app)

	require.NoError(t, middleware.OnTimeoutPacket(ctx, packet, nil))
	require.NoError(t, middleware.OnTimeoutPacket(ctx, packet, nil))
	require.Equal(t, 1, app.timeoutCalls)
	require.True(t, k.HasCompleted(ctx, pending.SourceChannel, pending.Sequence))
	require.Zero(t, k.PendingCount(ctx))
	require.Equal(t, "10", k.GetLedger(ctx).NobleBacking)
}

func TestInjectiveAcknowledgementLifecycleIsIdempotent(t *testing.T) {
	for _, success := range []bool{true, false} {
		t.Run(fmt.Sprintf("success=%t", success), func(t *testing.T) {
			ctx, k, bank := setupReceiveKeeper(t)
			controls := receiveControls()
			require.NoError(t, k.SetControls(ctx, controls))
			ledger := types.DefaultLedger()
			ledger.InjectiveBacking = "40"
			ledger.PendingInjective = "10"
			require.NoError(t, k.SetLedger(ctx, ledger))
			sender := sdk.AccAddress("ack-sender")
			bank.set(sender, controls.LogicalDenom, math.NewInt(40))
			bank.set(types.ModuleAddress, controls.LogicalDenom, math.NewInt(10))
			bank.set(types.ModuleAddress, controls.InjectiveDenom, math.NewInt(40))
			pending := types.PendingSettlement{
				Route: types.Route_ROUTE_INJECTIVE, SourceChannel: controls.InjectiveChannel,
				Sequence: 18, Sender: sender.String(), LogicalReceiver: "inj1receiver",
				PhysicalReceiver: "inj1receiver", LogicalDenom: controls.LogicalDenom,
				PhysicalDenom: controls.InjectiveDenom, Amount: "10",
			}
			require.NoError(t, k.SetPending(ctx, pending))
			data := transfertypes.FungibleTokenPacketData{
				Denom: controls.InjectivePacketDenom, Amount: pending.Amount,
				Sender: types.ModuleAddress.String(), Receiver: pending.PhysicalReceiver,
			}
			packet := channeltypes.Packet{
				Sequence: pending.Sequence, SourcePort: transfertypes.PortID,
				SourceChannel: pending.SourceChannel, Data: data.GetBytes(),
			}
			app := &receiveApp{bank: bank}
			var acknowledgement []byte
			if success {
				acknowledgement = channeltypes.NewResultAcknowledgement([]byte{1}).Acknowledgement()
			} else {
				refund := sdk.NewInt64Coin(controls.InjectiveDenom, 10)
				app.ackRefund = &refund
				acknowledgement = channeltypes.NewErrorAcknowledgement(fmt.Errorf("remote error")).Acknowledgement()
			}
			middleware := NewIBCMiddleware(k, app)

			require.NoError(t, middleware.OnAcknowledgementPacket(ctx, packet, acknowledgement, nil))
			require.NoError(t, middleware.OnAcknowledgementPacket(ctx, packet, acknowledgement, nil))
			require.Equal(t, 1, app.ackCalls)
			require.Zero(t, k.PendingCount(ctx))
			if success {
				require.Equal(t, "40", k.GetLedger(ctx).InjectiveBacking)
				require.Equal(t, "0", bank.amount(types.ModuleAddress, controls.LogicalDenom).String())
			} else {
				require.Equal(t, "50", k.GetLedger(ctx).InjectiveBacking)
				require.Equal(t, "50", bank.amount(sender, controls.LogicalDenom).String())
			}
		})
	}
}
