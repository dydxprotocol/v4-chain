package ante_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

var (
	anteHandlerRejectError = errorsmod.Wrap(
		sdkerrors.ErrInvalidRequest,
		"a short term clob message may not have a non-zero timeout height, use goodTilBlock instead",
	)
)

func TestRejectSTOrderTimeoutHeightDecorator_AnteHandle(t *testing.T) {
	tests := map[string]struct {
		msgs          []sdk.Msg
		timeoutHeight uint64
		expectedErr   error
	}{
		"do nothing for non-clob messages with zero timeout height": {
			msgs: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
			},
			timeoutHeight: 0,
		},
		"do nothing for non-clob messages with non-zero timeout height": {
			msgs: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
				&banktypes.MsgSend{
					FromAddress: constants.BobAccAddress.String(),
					ToAddress:   constants.AliceAccAddress.String(),
					Amount: []sdk.Coin{
						sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
					},
				},
			},
			timeoutHeight: 1,
		},
		"do nothing for long-term place order with non-zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder_LongTerm,
			},
			timeoutHeight: 1,
		},
		"do nothing for long-term cancel order with non-zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_CancelOrder_LongTerm,
			},
			timeoutHeight: 1,
		},
		"do nothing for short-term place order with zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
			},
			timeoutHeight: 0,
		},
		"reject short-term place order with non-zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
			},
			timeoutHeight: 1,
			expectedErr:   anteHandlerRejectError,
		},
		"do nothing for short-term cancel order with zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_CancelOrder,
			},
			timeoutHeight: 0,
		},
		"reject short-term cancel order with non-zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_CancelOrder,
			},
			timeoutHeight: 1,
			expectedErr:   anteHandlerRejectError,
		},
		"do nothing for short-term batch cancel with zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_BatchCancel,
			},
			timeoutHeight: 0,
		},
		"reject short-term batch cancel with non-zero timeout height": {
			msgs: []sdk.Msg{
				constants.Msg_BatchCancel,
			},
			timeoutHeight: 1,
			expectedErr:   anteHandlerRejectError,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			wrappedHandler := ante.NewRejectSTOrderTimeoutHeightDecorator()
			anteHandler := sdk.ChainAnteDecorators(wrappedHandler)

			// Empty private key, so tx's signature should be empty.
			var (
				privs   []cryptotypes.PrivKey
				accSeqs []uint64
				accNums []uint64
			)

			tx, err := tx.CreateTestTx(
				ctx,
				tc.msgs,
				privs,
				accNums,
				accSeqs,
				tApp.App.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
				tApp.App.TxConfig(),
				tc.timeoutHeight,
			)
			require.NoError(t, err)

			_, err = anteHandler(ctx, tx, false)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Error(t, tc.expectedErr, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
