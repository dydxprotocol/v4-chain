package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	accountKeeper     types.AccountKeeper
	bankKeeper        types.BankKeeper
	ics4Wrapper       types.ICS4Wrapper
	ratelimitKeeper   types.RatelimitKeeper
	channelKeeper     types.ChannelKeeper
	connectionKeeper  types.ConnectionKeeper
	transferMsgServer types.TransferMsgServer
	authorities       map[string]struct{}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	ics4Wrapper types.ICS4Wrapper,
	ratelimitKeeper types.RatelimitKeeper,
	channelKeeper types.ChannelKeeper,
	connectionKeeper types.ConnectionKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		ics4Wrapper:      ics4Wrapper,
		ratelimitKeeper:  ratelimitKeeper,
		channelKeeper:    channelKeeper,
		connectionKeeper: connectionKeeper,
		authorities:      lib.UniqueSliceToSet(authorities),
	}
}

// EnsureModuleAccount makes the module address a module account. Coins can be
// sent to the address before the module exists, leaving a plain BaseAccount
// there that bank refuses to mint into or burn from, so that account is
// converted in place, keeping its number and balances.
func (k Keeper) EnsureModuleAccount(ctx sdk.Context) error {
	existing := k.accountKeeper.GetAccount(ctx, types.ModuleAddress)
	if existing == nil {
		k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
		return nil
	}
	if _, ok := existing.(sdk.ModuleAccountI); ok {
		return nil
	}
	base, ok := existing.(*authtypes.BaseAccount)
	if !ok {
		return fmt.Errorf("%w: unexpected account type %T", types.ErrModuleAccount, existing)
	}
	if base.GetPubKey() != nil {
		return fmt.Errorf("%w: module address has a public key", types.ErrModuleAccount)
	}
	_, permissions := k.accountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	if len(permissions) == 0 {
		return fmt.Errorf("%w: module permissions are not registered", types.ErrModuleAccount)
	}
	k.accountKeeper.SetAccount(ctx, authtypes.NewModuleAccount(base, types.ModuleName, permissions...))
	return nil
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
		channel, found := k.channelKeeper.GetChannel(ctx, transfertypes.PortID, route.channel)
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
	ctx.KVStore(k.storeKey).Set([]byte(types.LedgerKey), k.cdc.MustMarshal(&ledger))
	return nil
}

// CreditInjectiveInbound records Injective-origin USDC received by the module.
// The resulting injective backing must not exceed the migration ceiling.
func (k Keeper) CreditInjectiveInbound(ctx sdk.Context, controls types.Controls, amount math.Int) error {
	ceiling, err := types.PositiveAmountOf(controls.MigrationCeiling)
	if err != nil {
		return err
	}
	ledger := k.GetLedger(ctx)
	current, err := types.AmountOf(ledger.InjectiveBacking)
	if err != nil {
		return err
	}
	// Compare against the remaining room so an amount near the Int bound cannot
	// overflow the addition.
	if amount.GT(ceiling.Sub(current)) {
		return types.ErrMigrationCeiling
	}
	ledger.InjectiveBacking = types.NewAmount(current.Add(amount))
	return k.SetLedger(ctx, ledger)
}

// snapshotNobleBacking classifies residual logical supply as Noble-backed.
// residual = supply(logical) - injective - pending injective.
func (k Keeper) snapshotNobleBacking(ctx sdk.Context) error {
	ledger := k.GetLedger(ctx)
	amounts, err := parseBackingAmounts(ledger)
	if err != nil {
		return err
	}
	classified := amounts.injective.Add(amounts.pendingInjective)
	supply := k.bankKeeper.GetSupply(ctx, k.GetControls(ctx).LogicalDenom).Amount
	if supply.LT(classified) {
		return fmt.Errorf(
			"%w: logical supply %s is below classified backing %s",
			types.ErrInvalidLedger,
			supply,
			classified,
		)
	}
	ledger.NobleBacking = types.NewAmount(supply.Sub(classified))
	return k.SetLedger(ctx, ledger)
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

func (k Keeper) pendingStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PendingKeyPrefix))
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
	store := k.pendingStore(ctx)
	key := types.PendingKey(pending.Route, pending.SourceChannel, pending.Sequence)
	if store.Has(key) {
		return types.ErrDuplicateSettlement
	}
	store.Set(key, k.cdc.MustMarshal(&pending))
	k.setPendingCount(ctx, count+1)
	return nil
}

func (k Keeper) GetPendingForPacket(ctx sdk.Context, channel string, sequence uint64) (types.PendingSettlement, bool) {
	store := k.pendingStore(ctx)
	for route := types.Route_ROUTE_INJECTIVE; route <= types.Route_ROUTE_NOBLE; route++ {
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
	store := k.pendingStore(ctx)
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
	iterator := k.pendingStore(ctx).Iterator(nil, nil)
	defer iterator.Close()
	pending := make([]types.PendingSettlement, 0, k.PendingCount(ctx))
	for ; iterator.Valid(); iterator.Next() {
		var settlement types.PendingSettlement
		k.cdc.MustUnmarshal(iterator.Value(), &settlement)
		pending = append(pending, settlement)
	}
	return pending
}

// GetPendingPage returns one page of pending settlements in key order.
func (k Keeper) GetPendingPage(
	ctx sdk.Context,
	pageRequest *query.PageRequest,
) ([]types.PendingSettlement, *query.PageResponse, error) {
	var pending []types.PendingSettlement
	pageResponse, err := query.Paginate(k.pendingStore(ctx), pageRequest, func(_ []byte, value []byte) error {
		var settlement types.PendingSettlement
		if err := k.cdc.Unmarshal(value, &settlement); err != nil {
			return err
		}
		pending = append(pending, settlement)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return pending, pageResponse, nil
}

func add(value dtypes.SerializableInt, amount math.Int) (dtypes.SerializableInt, error) {
	current, err := types.AmountOf(value)
	if err != nil {
		return dtypes.SerializableInt{}, err
	}
	return types.NewAmount(current.Add(amount)), nil
}

func subtract(value dtypes.SerializableInt, amount math.Int) (dtypes.SerializableInt, error) {
	current, err := types.AmountOf(value)
	if err != nil {
		return dtypes.SerializableInt{}, err
	}
	if current.LT(amount) {
		return dtypes.SerializableInt{}, types.ErrInsufficientBacking
	}
	return types.NewAmount(current.Sub(amount)), nil
}
