package process_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeProposedOperationsTx(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txBuilder := encodingCfg.TxConfig.NewTxBuilder()

	// Valid.
	validMsgTxBytes := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Duplicate.
	_ = txBuilder.SetMsgs(constants.ValidEmptyMsgProposedOperations, constants.ValidEmptyMsgProposedOperations)
	duplicateMsgTxBytes, _ := encodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())

	// Incorrect type.
	incorrectMsgTxBytes := constants.ValidMsgUpdateMarketPricesTxBytes

	tests := map[string]struct {
		txBytes []byte

		expectedErr error
		expectedMsg *types.MsgProposedOperations
	}{
		"Error: decode fails": {
			txBytes:     []byte{1, 2, 3}, // invalid bytes.
			expectedErr: errors.New("tx parse error: Decoding tx bytes failed"),
		},
		"Error: empty bytes": {
			txBytes: []byte{}, // empty returns 0 msgs.
			expectedErr: errors.New("Msg Type: types.MsgProposedOperations, " +
				"Expected 1 num of msgs, but got 0: Unexpected num of msgs"),
		},
		"Error: incorrect msg len": {
			txBytes: duplicateMsgTxBytes,
			expectedErr: errors.New("Msg Type: types.MsgProposedOperations, " +
				"Expected 1 num of msgs, but got 2: Unexpected num of msgs"),
		},
		"Error: incorrect msg type": {
			txBytes: incorrectMsgTxBytes,
			expectedErr: errors.New(
				"Expected MsgType types.MsgProposedOperations, but " +
					"got *types.MsgUpdateMarketPrices: Unexpected msg type",
			),
		},
		"Valid": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidEmptyMsgProposedOperations,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pot, err := process.DecodeProposedOperationsTx(encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, pot)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMsg, pot.GetMsg())
			}
		})
	}
}

func TestProposedOperationsTx_Validate(t *testing.T) {
	tests := map[string]struct {
		txBytes     []byte
		expectedErr error
	}{
		"Error: ValidateBasic fails": {
			txBytes: constants.InvalidProposedOperationsUnspecifiedOrderRemovalReasonTxBytes,
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidMsgProposedOperations,
				"order removal reason must be specified: {{dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 0} 0 64 0}",
			),
		},
		"Valid: ValidateBasic passes": {
			txBytes: constants.ValidEmptyMsgProposedOperationsTxBytes,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			pot, err := process.DecodeProposedOperationsTx(constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			require.NoError(t, err)

			// Run and Validate.
			err = pot.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProposedOperationsTx_GetMsg(t *testing.T) {
	validMsgTxBytes := constants.ValidEmptyMsgProposedOperationsTxBytes

	tests := map[string]struct {
		txWrapper   process.ProposedOperationsTx
		txBytes     []byte
		expectedMsg *types.MsgProposedOperations
	}{
		"Returns nil msg": {
			txWrapper: process.ProposedOperationsTx{},
		},
		"Returns valid msg": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidEmptyMsgProposedOperations,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msg sdk.Msg
			if tc.txBytes != nil {
				pot, err := process.DecodeProposedOperationsTx(constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
				require.NoError(t, err)
				msg = pot.GetMsg()
			} else {
				msg = tc.txWrapper.GetMsg()
			}
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}
