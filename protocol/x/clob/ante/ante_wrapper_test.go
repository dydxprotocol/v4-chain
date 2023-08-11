package ante_test

import (
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/dydxprotocol/v4/mocks"
	testante "github.com/dydxprotocol/v4/testutil/ante"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/ante"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidateMsgType_NewOffChainSingleMsgClobTx(t *testing.T) {
	tests := map[string]struct {
		msgOne sdk.Msg
		msgTwo sdk.Msg

		expectSkip  bool
		expectedErr error
	}{
		"no skip: no msg": {
			expectSkip: false,
		},
		"yes skip: single msg, MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,

			expectSkip: true,
		},
		"yes skip: single msg, Msg_CancelOrder": {
			msgOne: constants.Msg_CancelOrder,

			expectSkip: true,
		},
		"no skip: mult msgs, NO off-chain single msg clob tx": {
			msgOne: &testdata.TestMsg{Signers: []string{"meh"}},
			msgTwo: &testdata.TestMsg{Signers: []string{"meh"}},

			expectSkip: false,
		},
		"no skip: mult msgs, MsgCancelOrder with Transfer": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_Transfer,

			expectedErr: errors.ErrInvalidRequest,
		},
		"no skip: mult msgs, two MsgCancelOrder": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_CancelOrder,

			expectedErr: errors.ErrInvalidRequest,
		},
		"no skip: mult msgs, MsgPlaceOrder with Transfer": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_Transfer,

			expectedErr: errors.ErrInvalidRequest,
		},
		"no skip: mult msgs, two MsgPlaceOrder": {
			msgOne: constants.Msg_PlaceOrder,
			msgTwo: constants.Msg_PlaceOrder,

			expectedErr: errors.ErrInvalidRequest,
		},
		"no skip: mult msgs, MsgCancelOrder and MsgPlaceOrder": {
			msgOne: constants.Msg_CancelOrder,
			msgTwo: constants.Msg_PlaceOrder,

			expectedErr: errors.ErrInvalidRequest,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			suite := testante.SetupTestSuite(t, true)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

			mockAntehandler := &mocks.AnteDecorator{}
			mockAntehandler.On("AnteHandle", suite.Ctx, mock.Anything, false, mock.Anything).
				Return(suite.Ctx, nil)

			wrappedHandler := ante.NewOffChainSingleMsgClobTxAnteWrapper(mockAntehandler)
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
			require.Equal(t, suite.Ctx, resultCtx)

			if tc.expectSkip || tc.expectedErr != nil {
				mockAntehandler.AssertNotCalled(
					t,
					"AnteHandle",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				)
			} else {
				mockAntehandler.AssertCalled(
					t,
					"AnteHandle",
					suite.Ctx,
					tx,
					false,
					mock.Anything,
				)
				mockAntehandler.AssertNumberOfCalls(t, "AnteHandle", 1)
			}
		})
	}
}
