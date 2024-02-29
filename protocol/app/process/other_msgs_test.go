package process_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	"github.com/stretchr/testify/require"
)

func TestDecodeOtherMsgsTx(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()

	tests := map[string]struct {
		txBytes []byte

		expectedErr              error
		expectedErrTypeCheckOnly bool
		expectedMsgs             []sdk.Msg
	}{
		"Error: decode fails": {
			txBytes:                  []byte{1, 2, 3}, // invalid bytes.
			expectedErr:              errorsmod.Wrap(process.ErrDecodingTxBytes, "OtherMsgsTx Error"),
			expectedErrTypeCheckOnly: true,
		},
		"Error: empty bytes": {
			txBytes:     []byte{}, // empty returns 0 msgs.
			expectedErr: errorsmod.Wrap(process.ErrUnexpectedNumMsgs, "OtherMsgs len cannot be zero"),
		},
		"Error: app-injected msg type is not allowed": {
			txBytes: constants.ValidMsgUpdateMarketPricesTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *types.MsgUpdateMarketPrices",
			),
		},
		"Error: internal msg type is not allowed": {
			txBytes: testmsgs.MsgSoftwareUpgradeTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *types.MsgSoftwareUpgrade",
			),
		},
		"Error: unsupported msg type is not allowed": {
			txBytes: testmsgs.GovBetaMsgSubmitProposalTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *v1beta1.MsgSubmitProposal",
			),
		},
		"Error: nested msg type with unsupported inner is not allowed": {
			txBytes: testmsgs.MsgSubmitProposalWithUnsupportedInnerTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *v1.MsgSubmitProposal",
			),
		},
		"Error: nested msg type with app-injected inner is not allowed": {
			txBytes: testmsgs.MsgSubmitProposalWithAppInjectedInnerTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *v1.MsgSubmitProposal",
			),
		},
		"Error: nested msg type with double-nested inner is not allowed": {
			txBytes: testmsgs.MsgSubmitProposalWithDoubleNestedInnerTxBytes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Invalid msg type or content in OtherTxs *v1.MsgSubmitProposal",
			),
		},
		"Error: place order is not allowed": {
			txBytes: constants.Msg_PlaceOrder_TxBtyes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Msg type *types.MsgPlaceOrder is not allowed in OtherTxs",
			),
		},
		"Error: cancel order is not allowed": {
			txBytes: constants.Msg_CancelOrder_TxBtyes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Msg type *types.MsgCancelOrder is not allowed in OtherTxs",
			),
		},
		"Error: batch cancel order is not allowed": {
			txBytes: constants.Msg_BatchCancel_TxBtyes,
			expectedErr: errorsmod.Wrap(
				process.ErrUnexpectedMsgType,
				"Msg type *types.MsgBatchCancel is not allowed in OtherTxs",
			),
		},
		"Valid: single msg": {
			txBytes:      constants.Msg_Send_TxBytes,
			expectedMsgs: []sdk.Msg{constants.Msg_Send},
		},
		"Valid: mult msgs": {
			txBytes:      constants.Msg_SendAndTransfer_TxBytes,
			expectedMsgs: []sdk.Msg{constants.Msg_Send, constants.Msg_Transfer},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			omt, err := process.DecodeOtherMsgsTx(encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			if tc.expectedErr != nil {
				require.Nil(t, omt)
				if tc.expectedErrTypeCheckOnly {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.Equal(t, tc.expectedErr.Error(), err.Error())
				}
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.expectedMsgs, omt.GetMsgs())
			}
		})
	}
}

func TestOtherMsgsTx_Validate(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txBuilder := encodingCfg.TxConfig.NewTxBuilder()

	// Fails `ValidateBasic`
	failingSingleTx := constants.Msg_Transfer_Invalid_SameSenderAndRecipient_TxBytes

	_ = txBuilder.SetMsgs(constants.Msg_Send, constants.Msg_Transfer_Invalid_SameSenderAndRecipient) // invalid.
	failingMultiTx, _ := encodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())

	tests := map[string]struct {
		txBytes     []byte
		expectedErr error
	}{
		"Error Single: ValidateBasic fails": {
			txBytes:     failingSingleTx,
			expectedErr: errorsmod.Wrap(process.ErrMsgValidateBasic, "Sender is the same as recipient"),
		},
		"Error Multi: ValidateBasic fails": {
			txBytes:     failingMultiTx,
			expectedErr: errorsmod.Wrap(process.ErrMsgValidateBasic, "Sender is the same as recipient"),
		},
		"Valid Single: ValidateBasic passes": {
			txBytes: constants.Msg_Send_TxBytes,
		},
		"Valid Multi: ValidateBasic passes": {
			txBytes: constants.Msg_SendAndTransfer_TxBytes,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			omt, err := process.DecodeOtherMsgsTx(encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			require.NoError(t, err)

			err = omt.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOtherMsgsTx_GetMsgs(t *testing.T) {
	tests := map[string]struct {
		txWrapper    process.OtherMsgsTx
		txBytes      []byte
		expectedMsgs []sdk.Msg
	}{
		"Returns nil msg": {
			txWrapper: process.OtherMsgsTx{},
		},
		"Returns valid msg: single": {
			txBytes:      constants.Msg_Send_TxBytes,
			expectedMsgs: []sdk.Msg{constants.Msg_Send},
		},
		"Returns valid msg: multi": {
			txBytes:      constants.Msg_SendAndTransfer_TxBytes,
			expectedMsgs: []sdk.Msg{constants.Msg_Send, constants.Msg_Transfer},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msgs []sdk.Msg
			if tc.txBytes != nil {
				omt, err := process.DecodeOtherMsgsTx(constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
				require.NoError(t, err)
				msgs = omt.GetMsgs()
			} else {
				msgs = tc.txWrapper.GetMsgs()
			}
			require.Equal(t, tc.expectedMsgs, msgs)
		})
	}
}
