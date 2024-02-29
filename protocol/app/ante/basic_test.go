package ante_test

import (
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/stretchr/testify/require"
)

// Forked from:
// https://github.com/cosmos/cosmos-sdk/blob/1e8e923d3174cdfdb42454a96c27251ad72b6504/x/auth/ante/basic.go#L18
func TestValidateBasic_AppInjectedMsgWrapper(t *testing.T) {
	// The Ante suite and test setup from Cosmos SDK has "cosmos" hardcoded in it as the address codec which
	// is why we are required to use it below.
	tests := map[string]struct {
		msgOne         sdk.Msg
		msgTwo         sdk.Msg
		isRecheck      bool
		txHasSignature bool

		expectedErr error
	}{
		"fails ValidateBasic: no msg": {
			txHasSignature: false, // this should cause ValidateBasic to fail.

			expectedErr: sdkerrors.ErrNoSignatures,
		},
		"valid ValidateBasic: single msg": {
			msgOne:         &testdata.TestMsg{Signers: []string{constants.AliceAccAddress.String()}},
			txHasSignature: true, // this should allow ValidateBasic to pass.

			expectedErr: nil,
		},
		"fails ValidateBasic: mult msgs, AppInjected msg": {
			msgOne:         &pricestypes.MsgUpdateMarketPrices{}, // AppInjected.
			msgTwo:         &testdata.TestMsg{Signers: []string{constants.AliceAccAddress.String()}},
			txHasSignature: true,

			expectedErr: nil,
		},
		"valid: mult msgs, NO AppInjected msg": {
			msgOne:         &testdata.TestMsg{Signers: []string{constants.AliceAccAddress.String()}},
			msgTwo:         &testdata.TestMsg{Signers: []string{constants.AliceAccAddress.String()}},
			txHasSignature: true, // this should allow ValidateBasic to pass.

			expectedErr: nil,
		},
		"skip ValidateBasic: recheck": {
			msgOne:         &pricestypes.MsgUpdateMarketPrices{}, // AppInjected.
			msgTwo:         &testdata.TestMsg{Signers: []string{constants.AliceAccAddress.String()}},
			isRecheck:      true,
			txHasSignature: false, // this should cause ValidateBasic to fail, but this is skipped.

			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			suite := testante.SetupTestSuite(t, true)
			suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

			vbd := ante.NewValidateBasicDecorator()
			antehandler := sdk.ChainAnteDecorators(vbd)

			msgs := make([]sdk.Msg, 0)
			if tc.msgOne != nil {
				msgs = append(msgs, tc.msgOne)
			}
			if tc.msgTwo != nil {
				msgs = append(msgs, tc.msgTwo)
			}
			require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

			var privs []cryptotypes.PrivKey
			var accNums []uint64
			var accSeqs []uint64

			if tc.txHasSignature {
				priv1, _, _ := testdata.KeyTestPubAddr()
				privs, accNums, accSeqs = []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
			} else {
				// Empty private key, so tx's signature should be empty.
				privs, accNums, accSeqs = []cryptotypes.PrivKey{}, []uint64{}, []uint64{}
			}

			tx, err := suite.CreateTestTx(
				suite.Ctx,
				privs,
				accNums,
				accSeqs,
				suite.Ctx.ChainID(),
				signing.SignMode_SIGN_MODE_DIRECT,
			)
			require.NoError(t, err)

			if tc.isRecheck {
				suite.Ctx = suite.Ctx.WithIsReCheckTx(true) // test decorator skips on recheck
			}

			_, err = antehandler(suite.Ctx, tx, false)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
