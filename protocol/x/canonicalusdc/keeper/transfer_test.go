package keeper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	nobleChannel = "channel-0"
	injChannel   = "channel-1"
)

var physicalDenom = transfertypes.ParseDenomTrace("transfer/" + injChannel + "/uusdc").IBCDenom()

type bankMock struct {
	balances map[string]map[string]math.Int
}

func newBankMock() *bankMock { return &bankMock{balances: make(map[string]map[string]math.Int)} }

func (b *bankMock) amount(address sdk.AccAddress, denom string) math.Int {
	denoms, ok := b.balances[address.String()]
	if !ok {
		return math.ZeroInt()
	}
	amount, ok := denoms[denom]
	if !ok {
		return math.ZeroInt()
	}
	return amount
}

func (b *bankMock) set(address sdk.AccAddress, denom string, amount math.Int) {
	if _, ok := b.balances[address.String()]; !ok {
		b.balances[address.String()] = make(map[string]math.Int)
	}
	b.balances[address.String()][denom] = amount
}

func (b *bankMock) move(from, to sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		fromAmount := b.amount(from, coin.Denom)
		if fromAmount.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds")
		}
		b.set(from, coin.Denom, fromAmount.Sub(coin.Amount))
		b.set(to, coin.Denom, b.amount(to, coin.Denom).Add(coin.Amount))
	}
	return nil
}

func (b *bankMock) GetBalance(_ context.Context, address sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, b.amount(address, denom))
}

func (b *bankMock) GetSupply(_ context.Context, denom string) sdk.Coin {
	total := math.ZeroInt()
	for _, denoms := range b.balances {
		if amount, ok := denoms[denom]; ok {
			total = total.Add(amount)
		}
	}
	return sdk.NewCoin(denom, total)
}

func (b *bankMock) SendCoinsFromAccountToModule(
	_ context.Context,
	sender sdk.AccAddress,
	module string,
	coins sdk.Coins,
) error {
	return b.move(sender, moduleAddress(module), coins)
}

func (b *bankMock) SendCoinsFromModuleToAccount(
	_ context.Context,
	module string,
	receiver sdk.AccAddress,
	coins sdk.Coins,
) error {
	return b.move(moduleAddress(module), receiver, coins)
}

func (b *bankMock) MintCoins(_ context.Context, module string, coins sdk.Coins) error {
	address := moduleAddress(module)
	for _, coin := range coins {
		b.set(address, coin.Denom, b.amount(address, coin.Denom).Add(coin.Amount))
	}
	return nil
}

func (b *bankMock) BurnCoins(_ context.Context, module string, coins sdk.Coins) error {
	address := moduleAddress(module)
	for _, coin := range coins {
		amount := b.amount(address, coin.Denom)
		if amount.LT(coin.Amount) {
			return fmt.Errorf("insufficient funds")
		}
		b.set(address, coin.Denom, amount.Sub(coin.Amount))
	}
	return nil
}

func moduleAddress(module string) sdk.AccAddress {
	if module == canonicaltypes.ModuleName {
		return canonicaltypes.ModuleAddress
	}
	panic("unexpected module")
}

type ics4Mock struct {
	sends int
	data  []byte
}

type channelMock map[string]channeltypes.Channel

func (m channelMock) GetChannel(_ sdk.Context, _ string, channel string) (channeltypes.Channel, bool) {
	value, found := m[channel]
	return value, found
}

type connectionMock map[string]connectiontypes.ConnectionEnd

func (m connectionMock) GetConnection(_ sdk.Context, connection string) (connectiontypes.ConnectionEnd, bool) {
	value, found := m[connection]
	return value, found
}

func (m *ics4Mock) SendPacket(
	_ sdk.Context,
	_ *capabilitytypes.Capability,
	_, _ string,
	_ clienttypes.Height,
	_ uint64,
	data []byte,
) (uint64, error) {
	m.sends++
	m.data = append([]byte(nil), data...)
	return 11, nil
}

func (*ics4Mock) WriteAcknowledgement(
	sdk.Context,
	*capabilitytypes.Capability,
	ibcexported.PacketI,
	ibcexported.Acknowledgement,
) error {
	return nil
}

func (*ics4Mock) GetAppVersion(sdk.Context, string, string) (string, bool) { return "ics20-1", true }

type transferMock struct {
	bank          *bankMock
	calls         int
	seen          *transfertypes.MsgTransfer
	senderBalance math.Int
	logicalOwner  sdk.AccAddress
}

func (m *transferMock) Transfer(
	_ context.Context,
	msg *transfertypes.MsgTransfer,
) (*transfertypes.MsgTransferResponse, error) {
	m.calls++
	m.seen = msg
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	m.senderBalance = m.bank.amount(m.logicalOwner, types.UusdcDenom)
	if sender.Equals(canonicaltypes.ModuleAddress) {
		if err := m.bank.BurnCoins(context.Background(), canonicaltypes.ModuleName, sdk.NewCoins(msg.Token)); err != nil {
			return nil, err
		}
	} else {
		m.bank.set(sender, msg.Token.Denom, m.bank.amount(sender, msg.Token.Denom).Sub(msg.Token.Amount))
	}
	return &transfertypes.MsgTransferResponse{Sequence: 7}, nil
}

func (*transferMock) UpdateParams(
	context.Context,
	*transfertypes.MsgUpdateParams,
) (*transfertypes.MsgUpdateParamsResponse, error) {
	return &transfertypes.MsgUpdateParamsResponse{}, nil
}

func setupKeeper(t *testing.T) (sdk.Context, *Keeper, *bankMock, *transferMock, *ics4Mock) {
	t.Helper()
	key := storetypes.NewKVStoreKey(canonicaltypes.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger()).WithBlockTime(time.Unix(100, 0))
	registry := cdctypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	bank := newBankMock()
	ics4 := &ics4Mock{}
	k := NewKeeper(cdc, key, bank, ics4, nil, nil, nil)
	transfer := &transferMock{bank: bank}
	k.SetTransferMsgServer(transfer)
	return ctx, k, bank, transfer, ics4
}

func activeControls() canonicaltypes.Controls {
	return canonicaltypes.Controls{
		Mode:                    canonicaltypes.Mode_MODE_GRADUAL,
		LogicalDenom:            types.UusdcDenom,
		NobleChannel:            nobleChannel,
		InjectiveChannel:        injChannel,
		NoblePacketDenom:        "uusdc",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          physicalDenom,
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

func TestTransferDisabledPassesOriginalMessageUnchanged(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	sender := sdk.AccAddress("disabled-sender")
	bank.set(sender, "uatom", math.NewInt(10))
	msg := &transfertypes.MsgTransfer{
		SourcePort:    transfertypes.PortID,
		SourceChannel: "channel-7",
		Token:         sdk.NewInt64Coin("uatom", 2),
		Sender:        sender.String(),
		Receiver:      "remote-receiver",
	}
	response, err := NewTransferDecorator(k).Transfer(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.Equal(t, uint64(7), response.Sequence)
	require.Same(t, msg, transfer.seen)
}

func TestActiveNonCanonicalTransferPassesOriginalMessageUnchanged(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	sender := sdk.AccAddress("noncanonical-sender")
	bank.set(sender, "uatom", math.NewInt(10))
	msg := &transfertypes.MsgTransfer{
		SourcePort:    transfertypes.PortID,
		SourceChannel: "channel-7",
		Token:         sdk.NewInt64Coin("uatom", 2),
		Sender:        sender.String(),
		Receiver:      "remote-receiver",
	}

	response, err := NewTransferDecorator(k).Transfer(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.Equal(t, uint64(7), response.Sequence)
	require.Same(t, msg, transfer.seen)
}

func TestInjectiveTransferLocksBeforePhysicalSendAndSettlesError(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	ledger := canonicaltypes.DefaultLedger()
	ledger.InjectiveBacking = "50"
	require.NoError(t, k.SetLedger(ctx, ledger))
	sender := sdk.AccAddress("injective-sender")
	transfer.logicalOwner = sender
	bank.set(sender, types.UusdcDenom, math.NewInt(50))
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(50))
	msg := &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    injChannel,
		Token:            sdk.NewInt64Coin(types.UusdcDenom, 10),
		Sender:           sender.String(),
		Receiver:         "inj1receiver",
		TimeoutTimestamp: 1,
		Memo:             "partner-memo",
	}
	response, err := NewTransferDecorator(k).Transfer(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.Equal(t, uint64(7), response.Sequence)
	require.Equal(t, "40", transfer.senderBalance.String(), "logical USDC must be locked before the original keeper runs")
	require.Equal(t, physicalDenom, transfer.seen.Token.Denom)
	require.Equal(t, canonicaltypes.ModuleAddress.String(), transfer.seen.Sender)
	require.Equal(t, types.UusdcDenom, msg.Token.Denom, "submitted message must not be mutated")
	require.Equal(t, "40", bank.amount(sender, types.UusdcDenom).String())
	require.Equal(t, "10", bank.amount(canonicaltypes.ModuleAddress, types.UusdcDenom).String())
	require.Equal(t, "40", bank.amount(canonicaltypes.ModuleAddress, physicalDenom).String())
	require.Equal(t, uint32(1), k.PendingCount(ctx))
	require.Equal(t, "40", k.GetLedger(ctx).InjectiveBacking)
	require.Equal(t, "10", k.GetLedger(ctx).PendingInjective)
	var logicalEvent sdk.Event
	for _, event := range ctx.EventManager().Events() {
		if event.Type == transfertypes.EventTypeTransfer {
			logicalEvent = event
		}
	}
	require.Equal(t, transfertypes.EventTypeTransfer, logicalEvent.Type)
	attributes := make(map[string]string, len(logicalEvent.Attributes))
	for _, attribute := range logicalEvent.Attributes {
		attributes[attribute.Key] = attribute.Value
	}
	require.Equal(t, types.UusdcDenom, attributes[transfertypes.AttributeKeyDenom])
	require.Equal(t, "partner-memo", attributes[transfertypes.AttributeKeyMemo])

	pending, found := k.GetPendingForPacket(ctx, injChannel, 7)
	require.True(t, found)
	require.ErrorIs(
		t,
		k.Settle(ctx, pending, false, canonicaltypes.AttributeValueError),
		canonicaltypes.ErrInsufficientBacking,
	)
	require.Equal(t, "40", bank.amount(sender, types.UusdcDenom).String())
	require.Equal(t, uint32(1), k.PendingCount(ctx))

	// The transfer module must refund D_I before the outer middleware restores D_N.
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(50))
	require.NoError(t, k.Settle(ctx, pending, false, canonicaltypes.AttributeValueError))
	require.Equal(t, "50", bank.amount(sender, types.UusdcDenom).String())
	require.Equal(t, "0", bank.amount(canonicaltypes.ModuleAddress, types.UusdcDenom).String())
	require.Equal(t, "50", k.GetLedger(ctx).InjectiveBacking)
	require.Equal(t, "0", k.GetLedger(ctx).PendingInjective)
	require.Zero(t, k.PendingCount(ctx))
}

func TestCanonicalTransferRejectsUnsupportedChannelBeforeMovement(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	sender := sdk.AccAddress("unsupported-sender")
	bank.set(sender, types.UusdcDenom, math.NewInt(50))
	msg := &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    "channel-9",
		Token:            sdk.NewInt64Coin(types.UusdcDenom, 10),
		Sender:           sender.String(),
		Receiver:         "remote",
		TimeoutTimestamp: 1,
	}
	_, err := NewTransferDecorator(k).Transfer(sdk.WrapSDKContext(ctx), msg)
	require.ErrorIs(t, err, canonicaltypes.ErrUnsupportedChannel)
	require.Zero(t, transfer.calls)
	require.Equal(t, "50", bank.amount(sender, types.UusdcDenom).String())
}

func TestNobleTransferCannotUseInjectiveBacking(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	ledger := canonicaltypes.DefaultLedger()
	ledger.InjectiveBacking = "50"
	require.NoError(t, k.SetLedger(ctx, ledger))
	sender := sdk.AccAddress("noble-sender")
	bank.set(sender, types.UusdcDenom, math.NewInt(50))
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(50))
	msg := &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    nobleChannel,
		Token:            sdk.NewInt64Coin(types.UusdcDenom, 10),
		Sender:           sender.String(),
		Receiver:         "noble1receiver",
		TimeoutTimestamp: 1,
	}

	_, err := NewTransferDecorator(k).Transfer(sdk.WrapSDKContext(ctx), msg)
	require.ErrorIs(t, err, canonicaltypes.ErrInsufficientBacking)
	require.Zero(t, transfer.calls)
	require.Equal(t, "50", bank.amount(sender, types.UusdcDenom).String())
}

func TestBackingInvariantRequiresFullyClassifiedLogicalSupply(t *testing.T) {
	ctx, k, bank, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	owner := sdk.AccAddress("logical-owner")
	bank.set(owner, types.UusdcDenom, math.NewInt(10))
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(16))
	ledger := canonicaltypes.DefaultLedger()
	ledger.InjectiveBacking = "10"
	ledger.RestrictedFunding = "5"
	require.NoError(t, k.SetLedger(ctx, ledger))
	require.NoError(t, k.SetParticipant(ctx, canonicaltypes.Participant{
		Controller: sdk.AccAddress("funding-controller").String(), NobleRecipient: "noble1receiver",
		MaxRelease: "5", Released: "0", Funded: "5",
	}))

	reason, broken := BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)

	ledger.InjectiveBacking = "9"
	require.NoError(t, k.SetLedger(ctx, ledger))
	reason, broken = BackingInvariant(*k)(ctx)
	require.True(t, broken)
	require.Contains(t, reason, "does not equal accounted backing")
}

func TestValidateRouteBindingsRequiresExactOpenIcs20Channels(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	controls := activeControls()
	openChannel := func(connection string) channeltypes.Channel {
		return channeltypes.NewChannel(
			channeltypes.OPEN,
			channeltypes.UNORDERED,
			channeltypes.NewCounterparty(transfertypes.PortID, "channel-remote"),
			[]string{connection},
			transfertypes.Version,
		)
	}
	k.channelKeeper = channelMock{
		nobleChannel: openChannel(controls.NobleConnection),
		injChannel:   openChannel(controls.InjectiveConnection),
	}
	k.connectionKeeper = connectionMock{
		controls.NobleConnection: connectiontypes.NewConnectionEnd(
			connectiontypes.OPEN, controls.NobleClient, connectiontypes.Counterparty{}, nil, 0,
		),
		controls.InjectiveConnection: connectiontypes.NewConnectionEnd(
			connectiontypes.OPEN, controls.InjectiveClient, connectiontypes.Counterparty{}, nil, 0,
		),
	}
	require.NoError(t, k.ValidateRouteBindings(ctx, controls))

	channel := k.channelKeeper.(channelMock)[injChannel]
	channel.State = channeltypes.CLOSED
	k.channelKeeper.(channelMock)[injChannel] = channel
	require.ErrorIs(t, k.ValidateRouteBindings(ctx, controls), canonicaltypes.ErrInvalidControls)
}

func TestBackingSwapSettlementLifecycle(t *testing.T) {
	for _, success := range []bool{true, false} {
		t.Run(fmt.Sprintf("success=%t", success), func(t *testing.T) {
			ctx, k, bank, _, _ := setupKeeper(t)
			require.NoError(t, k.SetControls(ctx, activeControls()))
			ledger := canonicaltypes.DefaultLedger()
			ledger.NobleBacking = "100"
			ledger.RestrictedFunding = "25"
			require.NoError(t, k.SetLedger(ctx, ledger))
			owner := sdk.AccAddress("backing-owner")
			bank.set(owner, types.UusdcDenom, math.NewInt(100))
			bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(25))
			controller := sdk.AccAddress("swap-controller")
			participant := canonicaltypes.Participant{
				Controller: controller.String(), NobleRecipient: "noble1receiver",
				MaxRelease: "100", Released: "0", Funded: "25",
			}
			require.NoError(t, k.SetParticipant(ctx, participant))

			sequence, err := k.ExecuteBackingSwap(
				ctx, controller.String(), math.NewInt(10), transfertypes.PortID, uint64(ctx.BlockTime().UnixNano())+1,
			)
			require.NoError(t, err)
			pending, found := k.GetPendingForPacket(ctx, nobleChannel, sequence)
			require.True(t, found)
			require.Equal(t, canonicaltypes.Route_ROUTE_BACKING_SWAP, pending.Route)
			require.Equal(t, "90", k.GetLedger(ctx).NobleBacking)
			require.Equal(t, "15", k.GetLedger(ctx).RestrictedFunding)
			require.Equal(t, "10", k.GetLedger(ctx).PendingBackingSwap)

			if !success {
				require.NoError(t, bank.MintCoins(
					ctx, canonicaltypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.UusdcDenom, 10)),
				))
			}
			require.NoError(t, k.Settle(ctx, pending, success, canonicaltypes.AttributeValueSuccess))
			stored, found := k.GetParticipant(ctx, controller.String())
			require.True(t, found)
			if success {
				require.Equal(t, "10", k.GetLedger(ctx).InjectiveBacking)
				require.Equal(t, "15", stored.Funded)
				require.Equal(t, "10", stored.Released)
			} else {
				require.Equal(t, "100", k.GetLedger(ctx).NobleBacking)
				require.Equal(t, "25", k.GetLedger(ctx).RestrictedFunding)
				require.Equal(t, "25", stored.Funded)
				require.Equal(t, "0", stored.Released)
			}
			reason, broken := BackingInvariant(*k)(ctx)
			require.False(t, broken, reason)
		})
	}
}

func TestICS4GuardRejectsCanonicalBypass(t *testing.T) {
	ctx, k, _, _, ics4 := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	data := transfertypes.FungibleTokenPacketData{
		Denom:    "uusdc",
		Amount:   "1",
		Sender:   "sender",
		Receiver: "receiver",
	}.GetBytes()
	_, err := k.SendPacket(ctx, nil, transfertypes.PortID, nobleChannel, clienttypes.ZeroHeight(), 1, data)
	require.ErrorIs(t, err, canonicaltypes.ErrBypass)
	require.Zero(t, ics4.sends)

	authorized := ctx.WithValue(transferAuthorizationKey{}, true)
	_, err = k.SendPacket(authorized, nil, transfertypes.PortID, nobleChannel, clienttypes.ZeroHeight(), 1, data)
	require.NoError(t, err)
	require.Equal(t, 1, ics4.sends)
	require.Equal(t, data, ics4.data, "authorized physical packet bytes must be forwarded unchanged")
}
