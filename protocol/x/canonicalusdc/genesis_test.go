package canonicalusdc_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func genesisKeeper(t *testing.T) (sdk.Context, *keeper.Keeper) {
	t.Helper()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	database := db.NewMemDB()
	multiStore := store.NewCommitMultiStore(database, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	multiStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, database)
	require.NoError(t, multiStore.LoadLatestVersion())
	ctx := sdk.NewContext(multiStore, cmtproto.Header{}, false, log.NewNopLogger())
	k := keeper.NewKeeper(
		codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()),
		key,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	return ctx, k
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
					Amount:        "1",
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
