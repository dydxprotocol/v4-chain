package ante_test

import (
	"reflect"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/stretchr/testify/require"
)

const freeInfiniteGasMeterType = "*types.freeInfiniteGasMeter"

func TestValidateMsgType_FreeInfiniteGasDecorator(t *testing.T) {
	tests := map[string]struct {
		msgOne sdk.Msg
		msgTwo sdk.Msg

		expectFreeInfiniteGasMeter bool
		expectedErr                error
	}{
		"no freeInfiniteGasMeter: no msg": {
			expectFreeInfiniteGasMeter: false,
		},
		"yes freeInfiniteGasMeter: single msg, MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: single msg, Msg_CancelOrder": {
			msgOne: constants.Msg_CancelOrder,

			expectFreeInfiniteGasMeter: true,
		},
		"yes freeInfiniteGasMeter: single msg, MsgUpdateMarketPrices": {
			msgOne: &pricestypes.MsgUpdateMarketPrices{}, // app-injected.

			expectFreeInfiniteGasMeter: true,
		},
		"no freeInfiniteGasMeter: single msg": {
			msgOne: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: multi msg, MsgUpdateMarketPrices": {
			msgOne: &pricestypes.MsgUpdateMarketPrices{}, // app-injected.
			msgTwo: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: mult msgs, NO off-chain single msg clob tx": {
			msgOne: &testdata.TestMsg{Signers: []string{"meh"}},
			msgTwo: &testdata.TestMsg{Signers: []string{"meh"}},

			expectFreeInfiniteGasMeter: false,
		},
		"no freeInfiniteGasMeter: mult msgs, MsgCancelOrder with Transfer": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_Transfer,

			expectedErr: sdkerrors.ErrInvalidRequest,
		},
		"no freeInfiniteGasMeter: mult msgs, two MsgCancelOrder": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_CancelOrder,

			expectedErr: sdkerrors.ErrInvalidRequest,
		},
		"no freeInfiniteGasMeter: mult msgs, MsgPlaceOrder with Transfer": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_Transfer,

			expectedErr: sdkerrors.ErrInvalidRequest,
		},
		"no freeInfiniteGasMeter: mult msgs, two MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_PlaceOrder,

			expectedErr: sdkerrors.ErrInvalidRequest,
		},
		"no freeInfiniteGasMeter: mult msgs, MsgPlaceOrder and MsgCancelOrder": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_CancelOrder,

			expectedErr: sdkerrors.ErrInvalidRequest,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			suite := testante.SetupTestSuite(t, true)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

			wrappedHandler := ante.NewFreeInfiniteGasDecorator()
			antehandler := sdk.ChainAnteDecorators(wrappedHandler)

			msgs := make([]sdk.Msg, 0)
			if tc.msgOne != nil {
				msgs = append(msgs, tc.msgOne)
			}
			if tc.msgTwo != nil {
				msgs = append(msgs, tc.msgTwo)
			}

			require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

			// Empty private key, so tx's signature should be empty.
			privs, accNums, accSeqs := []cryptotypes.PrivKey{}, []uint64{}, []uint64{}

			tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.Ctx.ChainID())
			require.NoError(t, err)

			resultCtx, err := antehandler(suite.Ctx, tx, false)
			require.ErrorIs(t, tc.expectedErr, err)

			meter := resultCtx.GasMeter()

			if !tc.expectFreeInfiniteGasMeter || tc.expectedErr != nil {
				require.NotEqual(t, freeInfiniteGasMeterType, reflect.TypeOf(meter).String())
				require.Equal(t, suite.Ctx, resultCtx)
			} else {
				require.Equal(t, freeInfiniteGasMeterType, reflect.TypeOf(meter).String())
			}
		})
	}
}
