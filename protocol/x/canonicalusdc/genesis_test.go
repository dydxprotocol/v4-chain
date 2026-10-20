package canonicalusdc_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func genesisKeeper(t *testing.T) (sdk.Context, *keeper.Keeper) {
	t.Helper()
	ctx, k, _ := genesisKeeperWithAccounts(t)
	return ctx, k
}

func genesisKeeperWithAccounts(t *testing.T) (sdk.Context, *keeper.Keeper, *accountKeeperFake) {
	t.Helper()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	accounts := newAccountKeeperFake()
	k := keeper.NewKeeper(
		codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()),
		key,
		accounts,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	return ctx, k, accounts
}

func TestInitExportGenesisRoundTrip(t *testing.T) {
	ctx, k := genesisKeeper(t)
	canonicalusdc.InitGenesis(ctx, *k, *types.DefaultGenesis())
	exported := canonicalusdc.ExportGenesis(ctx, *k)
	require.Equal(t, types.DefaultControls(), exported.Controls)
	require.Equal(t, types.DefaultLedger(), exported.Ledger)
	require.Empty(t, exported.PendingSettlements)
	require.NoError(t, exported.Validate())
}

func TestGenesisValidateRejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*types.GenesisState)
		err  error
	}{
		{
			name: "disabled pending",
			mut: func(genesis *types.GenesisState) {
				genesis.Controls.MaxPendingSettlements = 10
				genesis.PendingSettlements = []types.PendingSettlement{{
					Route:         types.Route_ROUTE_NOBLE,
					SourceChannel: "channel-0",
					Sequence:      1,
					Amount:        dtypes.NewInt(1),
				}}
			},
			err: types.ErrInvalidControls,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			genesis := types.DefaultGenesis()
			tc.mut(genesis)
			require.ErrorIs(t, genesis.Validate(), tc.err)
		})
	}
}

func TestInitGenesisPanicsOnInvalidState(t *testing.T) {
	ctx, k := genesisKeeper(t)
	genesis := types.DefaultGenesis()
	genesis.Controls.LogicalDenom = "uatom"
	require.Panics(t, func() {
		canonicalusdc.InitGenesis(ctx, *k, *genesis)
	})
}

func TestInitGenesisDisabledSkipsRouteBindings(t *testing.T) {
	ctx, k := genesisKeeper(t)
	require.NotPanics(t, func() {
		canonicalusdc.InitGenesis(ctx, *k, *types.DefaultGenesis())
	})
}

func TestDefaultGenesisIsDisabled(t *testing.T) {
	require.Equal(t, types.Mode_MODE_DISABLED, types.DefaultGenesis().Controls.Mode)
	require.Equal(t, types.DefaultControls().LogicalDenom, transfertypes.ParseDenomTrace(
		"transfer/channel-0/uusdc",
	).IBCDenom())
}

func TestGenesisLedgerAmountsRoundTrip(t *testing.T) {
	ctx, k := genesisKeeper(t)
	genesis := types.DefaultGenesis()
	genesis.Ledger = types.Ledger{
		NobleBacking:     dtypes.NewInt(123),
		InjectiveBacking: dtypes.NewInt(45),
		PendingInjective: dtypes.NewInt(6),
		PendingNoble:     dtypes.NewInt(7),
	}
	canonicalusdc.InitGenesis(ctx, *k, *genesis)
	exported := canonicalusdc.ExportGenesis(ctx, *k)
	require.Equal(t, genesis.Ledger, exported.Ledger)

	cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	bz, err := cdc.MarshalJSON(exported)
	require.NoError(t, err)
	require.Contains(t, string(bz), `"noble_backing":"123"`)
	var decoded types.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(bz, &decoded))
	require.Equal(t, exported.Ledger, decoded.Ledger)
}

// accountKeeperFake mirrors x/auth: GetModuleAccount creates a missing module
// account and panics when a non-module account holds the address.
type accountKeeperFake struct {
	accounts map[string]sdk.AccountI
	next     uint64
}

func newAccountKeeperFake() *accountKeeperFake {
	return &accountKeeperFake{accounts: map[string]sdk.AccountI{}}
}

func (f *accountKeeperFake) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	return f.accounts[addr.String()]
}

func (f *accountKeeperFake) SetAccount(_ context.Context, acc sdk.AccountI) {
	f.accounts[acc.GetAddress().String()] = acc
}

func (f *accountKeeperFake) GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI {
	addr, permissions := f.GetModuleAddressAndPermissions(name)
	if existing := f.GetAccount(ctx, addr); existing != nil {
		moduleAccount, ok := existing.(sdk.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return moduleAccount
	}
	moduleAccount := authtypes.NewEmptyModuleAccount(name, permissions...)
	f.next++
	if err := moduleAccount.SetAccountNumber(f.next); err != nil {
		panic(err)
	}
	f.SetAccount(ctx, moduleAccount)
	return moduleAccount
}

func (f *accountKeeperFake) GetModuleAddressAndPermissions(name string) (sdk.AccAddress, []string) {
	if name != types.ModuleName {
		return nil, nil
	}
	return types.ModuleAddress, []string{authtypes.Minter, authtypes.Burner}
}

func requireCanonicalModuleAccount(t *testing.T, ctx sdk.Context, accounts *accountKeeperFake) sdk.ModuleAccountI {
	t.Helper()
	moduleAccount, ok := accounts.GetAccount(ctx, types.ModuleAddress).(sdk.ModuleAccountI)
	require.True(t, ok)
	require.Equal(t, types.ModuleName, moduleAccount.GetName())
	require.True(t, moduleAccount.HasPermission(authtypes.Minter))
	require.True(t, moduleAccount.HasPermission(authtypes.Burner))
	return moduleAccount
}

func TestInitGenesisCreatesModuleAccount(t *testing.T) {
	ctx, k, accounts := genesisKeeperWithAccounts(t)
	canonicalusdc.InitGenesis(ctx, *k, *types.DefaultGenesis())
	requireCanonicalModuleAccount(t, ctx, accounts)
}

func TestInitGenesisConvertsPreexistingBaseAccount(t *testing.T) {
	ctx, k, accounts := genesisKeeperWithAccounts(t)
	accounts.SetAccount(ctx, authtypes.NewBaseAccount(types.ModuleAddress, nil, 42, 3))

	canonicalusdc.InitGenesis(ctx, *k, *types.DefaultGenesis())

	moduleAccount := requireCanonicalModuleAccount(t, ctx, accounts)
	require.Equal(t, uint64(42), moduleAccount.GetAccountNumber())
	require.Equal(t, moduleAccount, accounts.GetModuleAccount(ctx, types.ModuleName))
}

func TestInitGenesisKeepsExistingModuleAccount(t *testing.T) {
	ctx, k, accounts := genesisKeeperWithAccounts(t)
	existing := accounts.GetModuleAccount(ctx, types.ModuleName)
	canonicalusdc.InitGenesis(ctx, *k, *types.DefaultGenesis())
	require.Equal(t, existing, accounts.GetAccount(ctx, types.ModuleAddress))
}

func TestEnsureModuleAccountRejects(t *testing.T) {
	vesting, err := vestingtypes.NewContinuousVestingAccount(
		authtypes.NewBaseAccountWithAddress(types.ModuleAddress),
		sdk.NewCoins(sdk.NewInt64Coin("adv4tnt", 1)),
		0,
		1,
	)
	require.NoError(t, err)
	cases := map[string]sdk.AccountI{
		"account with public key": authtypes.NewBaseAccount(
			types.ModuleAddress, secp256k1.GenPrivKey().PubKey(), 1, 0,
		),
		"unexpected account type": vesting,
	}
	for name, existing := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, k, accounts := genesisKeeperWithAccounts(t)
			accounts.SetAccount(ctx, existing)
			require.ErrorIs(t, k.EnsureModuleAccount(ctx), types.ErrModuleAccount)
			require.Equal(t, existing, accounts.GetAccount(ctx, types.ModuleAddress))
		})
	}
}
