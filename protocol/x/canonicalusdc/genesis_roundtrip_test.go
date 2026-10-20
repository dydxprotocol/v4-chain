package canonicalusdc_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// genBank serves supply and balances for invariant checks; transfers are not
// exercised by genesis import.
type genBank struct {
	supply   map[string]math.Int
	balances map[string]math.Int
}

func newGenBank() *genBank {
	return &genBank{supply: map[string]math.Int{}, balances: map[string]math.Int{}}
}

func (b *genBank) setBalance(addr sdk.AccAddress, denom string, amount int64) {
	b.balances[addr.String()+"/"+denom] = math.NewInt(amount)
}

func (b *genBank) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	amount, ok := b.balances[addr.String()+"/"+denom]
	if !ok {
		amount = math.ZeroInt()
	}
	return sdk.NewCoin(denom, amount)
}

func (b *genBank) GetSupply(_ context.Context, denom string) sdk.Coin {
	amount, ok := b.supply[denom]
	if !ok {
		amount = math.ZeroInt()
	}
	return sdk.NewCoin(denom, amount)
}

func (*genBank) SendCoinsFromAccountToModule(context.Context, sdk.AccAddress, string, sdk.Coins) error {
	panic("unused")
}

func (*genBank) SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error {
	panic("unused")
}

func (*genBank) MintCoins(context.Context, string, sdk.Coins) error { panic("unused") }
func (*genBank) BurnCoins(context.Context, string, sdk.Coins) error { panic("unused") }

func genKeeper(t *testing.T, bank *genBank) (sdk.Context, *keeper.Keeper) {
	t.Helper()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	// Nil channel/connection keepers skip route bindings; app ordering puts IBC
	// genesis before this module, which the app-level tests cover.
	k := keeper.NewKeeper(
		codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()),
		key, newAccountKeeperFake(), bank, nil, nil, nil, nil, nil,
	)
	return ctx, k
}

func genControls(mode types.Mode) types.Controls {
	return types.Controls{
		Mode:                    mode,
		LogicalDenom:            types.DefaultControls().LogicalDenom,
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
		MaxPendingSettlements:   10,
	}
}

func genPending(route types.Route, channel string, sequence uint64, amount int64, controls types.Controls) types.PendingSettlement {
	physical, sender := controls.LogicalDenom, "dydx1sender"
	if route == types.Route_ROUTE_INJECTIVE {
		physical = controls.InjectiveDenom
	}
	return types.PendingSettlement{
		Route:            route,
		SourceChannel:    channel,
		Sequence:         sequence,
		Sender:           sender,
		LogicalReceiver:  "receiver",
		PhysicalReceiver: "receiver",
		LogicalDenom:     controls.LogicalDenom,
		PhysicalDenom:    physical,
		Amount:           dtypes.NewInt(amount),
	}
}

// genConsistentState returns a genesis whose ledger matches its pending
// settlements and a bank whose supply and module balances back it.
func genConsistentState(mode types.Mode) (types.GenesisState, *genBank) {
	controls := genControls(mode)
	genesis := types.GenesisState{
		Controls: controls,
		Ledger: types.Ledger{
			NobleBacking:     dtypes.NewInt(1_000),
			InjectiveBacking: dtypes.NewInt(200),
			PendingInjective: dtypes.NewInt(30),
			PendingNoble:     dtypes.NewInt(40),
		},
		PendingSettlements: []types.PendingSettlement{
			genPending(types.Route_ROUTE_INJECTIVE, controls.InjectiveChannel, 7, 30, controls),
			genPending(types.Route_ROUTE_NOBLE, controls.NobleChannel, 9, 40, controls),
		},
	}
	bank := newGenBank()
	bank.supply[controls.LogicalDenom] = math.NewInt(1_000 + 200 + 30)
	bank.setBalance(types.ModuleAddress, controls.InjectiveDenom, 200)
	bank.setBalance(types.ModuleAddress, controls.LogicalDenom, 30)
	return genesis, bank
}

func TestGenesisRoundTripWithPendingSettlements(t *testing.T) {
	for _, mode := range []types.Mode{types.Mode_MODE_GRADUAL, types.Mode_MODE_PAUSED} {
		t.Run(mode.String(), func(t *testing.T) {
			genesis, bank := genConsistentState(mode)
			require.NoError(t, genesis.Validate())
			ctx, k := genKeeper(t, bank)
			canonicalusdc.InitGenesis(ctx, *k, genesis)

			require.Equal(t, uint32(2), k.PendingCount(ctx))
			_, broken := keeper.BackingInvariant(*k)(ctx)
			require.False(t, broken)
			for _, pending := range genesis.PendingSettlements {
				found, ok := k.GetPendingForPacket(ctx, pending.SourceChannel, pending.Sequence)
				require.True(t, ok)
				require.Equal(t, pending, found)
			}

			exported := canonicalusdc.ExportGenesis(ctx, *k)
			require.Equal(t, genesis.Controls, exported.Controls)
			require.Equal(t, genesis.Ledger, exported.Ledger)
			require.ElementsMatch(t, genesis.PendingSettlements, exported.PendingSettlements)

			// A second import of the export must reproduce identical bytes.
			cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
			first := cdc.MustMarshalJSON(exported)
			ctx2, k2 := genKeeper(t, bank)
			canonicalusdc.InitGenesis(ctx2, *k2, *exported)
			require.Equal(t, uint32(2), k2.PendingCount(ctx2))
			require.Equal(t, first, cdc.MustMarshalJSON(canonicalusdc.ExportGenesis(ctx2, *k2)))
		})
	}
}

func TestGenesisPendingAtLimitImportsButBlocksNewSettlements(t *testing.T) {
	genesis, bank := genConsistentState(types.Mode_MODE_GRADUAL)
	genesis.Controls.MaxPendingSettlements = 2
	require.NoError(t, genesis.Validate())
	ctx, k := genKeeper(t, bank)
	canonicalusdc.InitGenesis(ctx, *k, genesis)
	extra := genPending(types.Route_ROUTE_NOBLE, genesis.Controls.NobleChannel, 10, 1, genesis.Controls)
	require.ErrorIs(t, k.SetPending(ctx, extra), types.ErrPendingLimit)
}

// Finding: Validate does not reconcile ledger pending totals with the pending
// settlements, so such a genesis imports and immediately breaks the invariant.
func TestGenesisValidateDoesNotCheckPendingTotals(t *testing.T) {
	genesis, bank := genConsistentState(types.Mode_MODE_GRADUAL)
	genesis.Ledger.PendingInjective = dtypes.NewInt(0)
	genesis.Ledger.InjectiveBacking = dtypes.NewInt(230)
	require.NoError(t, genesis.Validate())
	ctx, k := genKeeper(t, bank)
	require.NotPanics(t, func() { canonicalusdc.InitGenesis(ctx, *k, genesis) })
	reason, broken := keeper.BackingInvariant(*k)(ctx)
	require.True(t, broken, reason)
}

// Finding: GRADUAL genesis is not checked against logical supply at import.
func TestGenesisValidateDoesNotCheckSupplyEquation(t *testing.T) {
	genesis, bank := genConsistentState(types.Mode_MODE_GRADUAL)
	genesis.Ledger.NobleBacking = dtypes.NewInt(1)
	require.NoError(t, genesis.Validate())
	ctx, k := genKeeper(t, bank)
	require.NotPanics(t, func() { canonicalusdc.InitGenesis(ctx, *k, genesis) })
	reason, broken := keeper.BackingInvariant(*k)(ctx)
	require.True(t, broken, reason)
}

// Finding: a pending settlement whose route disagrees with its channel or
// denoms passes Validate; its ack would later fail validatePendingPacket.
func TestGenesisValidateAcceptsMisroutedPending(t *testing.T) {
	genesis, _ := genConsistentState(types.Mode_MODE_GRADUAL)
	genesis.PendingSettlements[0].SourceChannel = genesis.Controls.NobleChannel
	genesis.PendingSettlements[0].Sequence = 11
	genesis.PendingSettlements[1].PhysicalDenom = "uatom"
	require.NoError(t, genesis.Validate())
}

func TestGenesisDuplicatePendingPassesValidateButPanicsOnImport(t *testing.T) {
	genesis, bank := genConsistentState(types.Mode_MODE_GRADUAL)
	genesis.PendingSettlements = append(genesis.PendingSettlements, genesis.PendingSettlements[0])
	require.NoError(t, genesis.Validate())
	ctx, k := genKeeper(t, bank)
	require.Panics(t, func() { canonicalusdc.InitGenesis(ctx, *k, genesis) })
}
