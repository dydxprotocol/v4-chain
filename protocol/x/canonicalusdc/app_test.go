package canonicalusdc_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	dydxapp "github.com/dydxprotocol/v4-chain/protocol/app"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestTestAppGenesisIsDisabledNoOp(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	app := tApp.App

	genesis := canonicalusdc.ExportGenesis(ctx, *app.CanonicalUsdcKeeper)
	require.Equal(t, canonicaltypes.Mode_MODE_DISABLED, genesis.Controls.Mode)
	require.Equal(t, types.UusdcDenom, genesis.Controls.LogicalDenom)
	require.NoError(t, genesis.Validate())
	require.Equal(t, uint64(1), canonicalusdc.AppModule{}.ConsensusVersion())

	require.False(t, dydxapp.BlockedAddresses()[canonicaltypes.ModuleAddress.String()])
	require.Contains(t, dydxapp.GetMaccPerms()[canonicaltypes.ModuleName], authtypes.Minter)
	require.Contains(t, dydxapp.GetMaccPerms()[canonicaltypes.ModuleName], authtypes.Burner)

	handler := app.MsgServiceRouter().Handler(&transfertypes.MsgTransfer{})
	require.NotNil(t, handler)

	state, err := app.CanonicalUsdcKeeper.State(ctx, &canonicaltypes.QueryStateRequest{})
	require.NoError(t, err)
	require.Equal(t, canonicaltypes.Mode_MODE_DISABLED, state.Controls.Mode)

	checkTx := testapp.MustMakeCheckTx(
		ctx,
		app,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.AliceAccAddress.String(),
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		&transfertypes.MsgTransfer{
			SourcePort:       transfertypes.PortID,
			SourceChannel:    "channel-0",
			Token:            sdk.NewInt64Coin(types.UusdcDenom, 1),
			Sender:           constants.AliceAccAddress.String(),
			Receiver:         "noble1receiver",
			TimeoutTimestamp: 1,
		},
	)
	result := tApp.CheckTx(checkTx)
	if !result.IsOK() {
		require.NotContains(t, result.Log, canonicaltypes.ErrUnsupportedChannel.Error())
		require.NotContains(t, result.Log, canonicaltypes.ErrPaused.Error())
		require.NotContains(t, result.Log, canonicaltypes.ErrNobleDepositsDisabled.Error())
	}
}

func TestBackingInvariantWithRealBank(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.CanonicalUsdcKeeper
	bank := tApp.App.BankKeeper
	physical := transfertypes.ParseDenomTrace("transfer/channel-1/uusdc").IBCDenom()
	controls := canonicaltypes.Controls{
		Mode:                    canonicaltypes.Mode_MODE_GRADUAL,
		LogicalDenom:            types.UusdcDenom,
		NobleChannel:            "channel-0",
		InjectiveChannel:        "channel-1",
		NoblePacketDenom:        "uusdc",
		InjectivePacketDenom:    "uusdc",
		InjectiveDenom:          physical,
		NobleWithdrawalsEnabled: true,
		MaxTransferAmount:       "1000",
		MigrationCeiling:        "10000",
		MaxPendingSettlements:   10,
		NobleClient:             "07-tendermint-0",
		NobleConnection:         "connection-0",
		InjectiveClient:         "07-tendermint-1",
		InjectiveConnection:     "connection-1",
	}
	require.NoError(t, k.SetControls(ctx, controls))

	holder := sdk.AccAddress("backing-holder-addr00")
	supply := bank.GetSupply(ctx, types.UusdcDenom).Amount
	ledger := canonicaltypes.DefaultLedger()
	ledger.NobleBacking = supply.String()
	require.NoError(t, k.SetLedger(ctx, ledger))
	reason, broken := keeper.BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)

	require.NoError(t, k.MintLogicalAndSend(ctx, holder, sdk.NewCoins(sdk.NewInt64Coin(types.UusdcDenom, 5))))
	reason, broken = keeper.BackingInvariant(*k)(ctx)
	require.True(t, broken)
	require.Contains(t, reason, "does not equal accounted backing")

	require.NoError(t, bank.MintCoins(ctx, canonicaltypes.ModuleName, sdk.NewCoins(
		sdk.NewInt64Coin(physical, 5),
	)))
	ledger.InjectiveBacking = "5"
	require.NoError(t, k.SetLedger(ctx, ledger))
	reason, broken = keeper.BackingInvariant(*k)(ctx)
	require.False(t, broken, reason)
	require.Equal(t, supply.Add(math.NewInt(5)), bank.GetSupply(ctx, types.UusdcDenom).Amount)
}
