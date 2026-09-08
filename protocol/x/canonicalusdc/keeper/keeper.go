package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	bankKeeper        types.BankKeeper
	ics4Wrapper       types.ICS4Wrapper
	channelKeeper     types.ChannelKeeper
	connectionKeeper  types.ConnectionKeeper
	transferMsgServer types.TransferMsgServer
	authorities       map[string]struct{}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	connectionKeeper types.ConnectionKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		bankKeeper:       bankKeeper,
		ics4Wrapper:      ics4Wrapper,
		channelKeeper:    channelKeeper,
		connectionKeeper: connectionKeeper,
		authorities:      lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) ValidateRouteBindings(ctx sdk.Context, controls types.Controls) error {
	if controls.Mode == types.Mode_MODE_DISABLED || k.channelKeeper == nil || k.connectionKeeper == nil {
		return nil
	}
	routes := []struct {
		name       string
		channel    string
		connection string
		client     string
	}{
		{name: "Noble", channel: controls.NobleChannel, connection: controls.NobleConnection, client: controls.NobleClient},
		{
			name:       "Injective",
			channel:    controls.InjectiveChannel,
			connection: controls.InjectiveConnection,
			client:     controls.InjectiveClient,
		},
	}
	for _, route := range routes {
		channel, found := k.channelKeeper.GetChannel(ctx, "transfer", route.channel)
		if !found || channel.State != channeltypes.OPEN || channel.Ordering != channeltypes.UNORDERED ||
			channel.Version != transfertypes.Version ||
			channel.Counterparty.PortId != transfertypes.PortID ||
			len(channel.ConnectionHops) != 1 || channel.ConnectionHops[0] != route.connection {
			return fmt.Errorf("%w: %s channel binding mismatch", types.ErrInvalidControls, route.name)
		}
		connection, found := k.connectionKeeper.GetConnection(ctx, route.connection)
		if !found || connection.State != connectiontypes.OPEN || connection.ClientId != route.client {
			return fmt.Errorf("%w: %s client binding mismatch", types.ErrInvalidControls, route.name)
		}
	}
	return nil
}

func (k *Keeper) SetTransferMsgServer(server types.TransferMsgServer) {
	if server == nil {
		panic("canonical USDC transfer message server must not be nil")
	}
	k.transferMsgServer = server
}

func (k Keeper) TransferMsgServer() types.TransferMsgServer {
	if k.transferMsgServer == nil {
		panic("canonical USDC transfer message server is not configured")
	}
	return k.transferMsgServer
}

func (k Keeper) HasAuthority(address string) bool {
	_, ok := k.authorities[address]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetControls(ctx sdk.Context) types.Controls {
	bz := ctx.KVStore(k.storeKey).Get([]byte(types.ControlsKey))
	if bz == nil {
		return types.DefaultControls()
	}
	var controls types.Controls
	k.cdc.MustUnmarshal(bz, &controls)
	return controls
}

func (k Keeper) SetControls(ctx sdk.Context, controls types.Controls) error {
	if err := controls.Validate(); err != nil {
		return err
	}
	ctx.KVStore(k.storeKey).Set([]byte(types.ControlsKey), k.cdc.MustMarshal(&controls))
	return nil
}

func (k Keeper) GetLedger(ctx sdk.Context) types.Ledger {
	bz := ctx.KVStore(k.storeKey).Get([]byte(types.LedgerKey))
	if bz == nil {
		return types.DefaultLedger()
	}
	var ledger types.Ledger
	k.cdc.MustUnmarshal(bz, &ledger)
	return ledger
}

func (k Keeper) SetLedger(ctx sdk.Context, ledger types.Ledger) error {
	if err := ledger.Validate(); err != nil {
		return err
	}
	controls := k.GetControls(ctx)
	if controls.Mode != types.Mode_MODE_DISABLED {
		ceiling, err := types.ParsePositiveAmount(controls.MigrationCeiling)
		if err != nil {
			return err
		}
		injective, _ := types.ParseAmount(ledger.InjectiveBacking)
		restricted, _ := types.ParseAmount(ledger.RestrictedFunding)
		pendingSwap, _ := types.ParseAmount(ledger.PendingBackingSwap)
		if injective.Add(restricted).Add(pendingSwap).GT(ceiling) {
			return types.ErrMigrationCeiling
		}
	}
	ctx.KVStore(k.storeKey).Set([]byte(types.LedgerKey), k.cdc.MustMarshal(&ledger))
	return nil
}

func (k Keeper) GetParticipant(ctx sdk.Context, controller string) (types.Participant, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ParticipantKeyPrefix))
	bz := store.Get(types.ParticipantKey(controller))
	if bz == nil {
		return types.Participant{}, false
	}
	var participant types.Participant
	k.cdc.MustUnmarshal(bz, &participant)
	return participant, true
}

func (k Keeper) GetAllParticipants(ctx sdk.Context) []types.Participant {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ParticipantKeyPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	participants := make([]types.Participant, 0)
	for ; iterator.Valid(); iterator.Next() {
		var participant types.Participant
		k.cdc.MustUnmarshal(iterator.Value(), &participant)
		participants = append(participants, participant)
	}
	return participants
}

func (k Keeper) ReplaceParticipants(ctx sdk.Context, participants []types.Participant) error {
	if len(participants) > types.MaxParticipants {
		return types.ErrInvalidParticipant
	}
	seen := make(map[string]struct{}, len(participants))
	for _, participant := range participants {
		if err := participant.Validate(); err != nil {
			return err
		}
		if _, exists := seen[participant.Controller]; exists {
			return fmt.Errorf("%w: duplicate controller", types.ErrInvalidParticipant)
		}
		seen[participant.Controller] = struct{}{}
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ParticipantKeyPrefix))
	iterator := store.Iterator(nil, nil)
	keys := make([][]byte, 0, types.MaxParticipants)
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, append([]byte(nil), iterator.Key()...))
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}
	for i := range participants {
		participant := participants[i]
		store.Set(types.ParticipantKey(participant.Controller), k.cdc.MustMarshal(&participant))
	}
	return nil
}

func (k Keeper) SetParticipant(ctx sdk.Context, participant types.Participant) error {
	if err := participant.Validate(); err != nil {
		return err
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ParticipantKeyPrefix))
	store.Set(types.ParticipantKey(participant.Controller), k.cdc.MustMarshal(&participant))
	return nil
}

func (k Keeper) PendingCount(ctx sdk.Context) uint32 {
	bz := ctx.KVStore(k.storeKey).Get([]byte(types.PendingCountKey))
	if len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint32(bz)
}

func (k Keeper) setPendingCount(ctx sdk.Context, count uint32) {
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, count)
	ctx.KVStore(k.storeKey).Set([]byte(types.PendingCountKey), bz)
}

func (k Keeper) SetPending(ctx sdk.Context, pending types.PendingSettlement) error {
	if err := pending.Validate(); err != nil {
		return err
	}
	count := k.PendingCount(ctx)
	controls := k.GetControls(ctx)
	if count >= controls.MaxPendingSettlements {
		return types.ErrPendingLimit
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingKeyPrefix))
	key := types.PendingKey(pending.Route, pending.SourceChannel, pending.Sequence)
	if store.Has(key) {
		return types.ErrDuplicateSettlement
	}
	store.Set(key, k.cdc.MustMarshal(&pending))
	k.setPendingCount(ctx, count+1)
	return nil
}

func (k Keeper) GetPendingForPacket(ctx sdk.Context, channel string, sequence uint64) (types.PendingSettlement, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingKeyPrefix))
	for route := types.Route_ROUTE_INJECTIVE; route <= types.Route_ROUTE_BACKING_SWAP; route++ {
		bz := store.Get(types.PendingKey(route, channel, sequence))
		if bz != nil {
			var pending types.PendingSettlement
			k.cdc.MustUnmarshal(bz, &pending)
			return pending, true
		}
	}
	return types.PendingSettlement{}, false
}

func (k Keeper) RemovePending(ctx sdk.Context, pending types.PendingSettlement) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingKeyPrefix))
	key := types.PendingKey(pending.Route, pending.SourceChannel, pending.Sequence)
	if !store.Has(key) {
		return false
	}
	store.Delete(key)
	count := k.PendingCount(ctx)
	if count == 0 {
		panic("canonical USDC pending count underflow")
	}
	k.setPendingCount(ctx, count-1)
	return true
}

func (k Keeper) GetAllPending(ctx sdk.Context) []types.PendingSettlement {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingKeyPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	pending := make([]types.PendingSettlement, 0, k.PendingCount(ctx))
	for ; iterator.Valid(); iterator.Next() {
		var settlement types.PendingSettlement
		k.cdc.MustUnmarshal(iterator.Value(), &settlement)
		pending = append(pending, settlement)
	}
	return pending
}

func (k Keeper) HasCompleted(ctx sdk.Context, channel string, sequence uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.CompletedKeyPrefix))
	return store.Has(types.CompletedKey(channel, sequence))
}

// markCompleted retains a bounded replay window for idempotent callbacks.
// Slots form a fixed-size ring; replacing a slot evicts its previous lookup key.
func (k Keeper) markCompleted(ctx sdk.Context, channel string, sequence uint64) {
	root := ctx.KVStore(k.storeKey)
	nextBytes := root.Get([]byte(types.CompletedNextKey))
	var next uint64
	if len(nextBytes) != 0 {
		next = binary.BigEndian.Uint64(nextBytes)
	}
	slot := next % types.HardMaxCompleted
	slots := prefix.NewStore(root, []byte(types.CompletedSlotPrefix))
	completed := prefix.NewStore(root, []byte(types.CompletedKeyPrefix))
	slotKey := types.CompletedSlotKey(slot)
	if previous := slots.Get(slotKey); previous != nil {
		completed.Delete(previous)
	}
	key := types.CompletedKey(channel, sequence)
	completed.Set(key, []byte{1})
	slots.Set(slotKey, key)

	nextBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(nextBytes, next+1)
	root.Set([]byte(types.CompletedNextKey), nextBytes)
}

func add(value string, amount math.Int) string {
	current, err := types.ParseAmount(value)
	if err != nil {
		panic(err)
	}
	return current.Add(amount).String()
}

func subtract(value string, amount math.Int) (string, error) {
	current, err := types.ParseAmount(value)
	if err != nil {
		return "", err
	}
	if current.LT(amount) {
		return "", types.ErrInsufficientBacking
	}
	return current.Sub(amount).String(), nil
}
