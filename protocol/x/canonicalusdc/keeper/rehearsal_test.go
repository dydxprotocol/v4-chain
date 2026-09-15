package keeper

import (
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInjectiveInboundThenNobleWithdrawRehearsal(t *testing.T) {
	ctx, k, bank, transfer, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	server := NewMsgServer(k)
	decorator := NewTransferDecorator(k)
	holder := sdk.AccAddress("holder")
	operator := sdk.AccAddress("operator")
	bank.set(holder, types.UusdcDenom, math.NewInt(100))

	_, err := server.UpdateControls(ctx, &canonicaltypes.MsgUpdateControls{
		Authority: lib.GovModuleAddress.String(),
		Controls:  activeControls(),
	})
	require.NoError(t, err)
	require.Equal(t, "100", k.GetLedger(ctx).NobleBacking)

	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(100))
	require.NoError(t, k.MintLogicalAndSend(ctx, operator, sdk.NewCoins(sdk.NewInt64Coin(types.UusdcDenom, 100))))
	ledger := k.GetLedger(ctx)
	ledger.InjectiveBacking = "100"
	require.NoError(t, k.SetLedger(ctx, ledger))
	requireInvariant(t, ctx, k)

	transfer.logicalOwner = operator
	resp, err := decorator.Transfer(sdk.WrapSDKContext(ctx), &transfertypes.MsgTransfer{
		SourcePort: transfertypes.PortID, SourceChannel: nobleChannel,
		Token: sdk.NewInt64Coin(types.UusdcDenom, 100), Sender: operator.String(),
		Receiver: "noble1receiver", TimeoutTimestamp: 1,
	})
	require.NoError(t, err)
	require.NoError(t, k.Settle(
		ctx, mustPending(t, k, ctx, nobleChannel, resp.Sequence), true, canonicaltypes.AttributeValueSuccess,
	))
	require.Equal(t, "0", k.GetLedger(ctx).NobleBacking)
	require.Equal(t, "100", k.GetLedger(ctx).InjectiveBacking)
	require.Equal(t, "0", k.GetLedger(ctx).PendingNoble)
	requireInvariant(t, ctx, k)

	_, err = decorator.Transfer(sdk.WrapSDKContext(ctx), &transfertypes.MsgTransfer{
		SourcePort: transfertypes.PortID, SourceChannel: nobleChannel,
		Token: sdk.NewInt64Coin(types.UusdcDenom, 1), Sender: operator.String(),
		Receiver: "noble1receiver", TimeoutTimestamp: 1,
	})
	require.ErrorIs(t, err, canonicaltypes.ErrInsufficientBacking)
}
