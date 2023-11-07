package process_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeAddPremiumVotesTx(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txBuilder := encodingCfg.TxConfig.NewTxBuilder()

	// Valid.
	validMsgTxBytes := constants.ValidMsgAddPremiumVotesTxBytes

	// Duplicate.
	_ = txBuilder.SetMsgs(constants.ValidMsgAddPremiumVotes, constants.ValidMsgAddPremiumVotes)
	duplicateMsgTxBytes, _ := encodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())

	// Incorrect type.
	incorrectMsgTxBytes := constants.ValidMsgUpdateMarketPricesTxBytes

	tests := map[string]struct {
		txBytes []byte

		expectedErr error
		expectedMsg *types.MsgAddPremiumVotes
	}{
		"Error: decode fails": {
			txBytes:     []byte{1, 2, 3}, // invalid bytes.
			expectedErr: errors.New("tx parse error: Decoding tx bytes failed"),
		},
		"Error: empty bytes": {
			txBytes: []byte{}, // empty returns 0 msgs.
			expectedErr: errors.New("Msg Type: types.MsgAddPremiumVotes, " +
				"Expected 1 num of msgs, but got 0: Unexpected num of msgs"),
		},
		"Error: incorrect msg len": {
			txBytes: duplicateMsgTxBytes,
			expectedErr: errors.New("Msg Type: types.MsgAddPremiumVotes, " +
				"Expected 1 num of msgs, but got 2: Unexpected num of msgs"),
		},
		"Error: incorrect msg type": {
			txBytes: incorrectMsgTxBytes,
			expectedErr: errors.New(
				"Expected MsgType types.MsgAddPremiumVotes, but " +
					"got *types.MsgUpdateMarketPrices: Unexpected msg type",
			),
		},
		"Valid": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidMsgAddPremiumVotes,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			afst, err := process.DecodeAddPremiumVotesTx(encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, afst)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMsg, afst.GetMsg())
			}
		})
	}
}

func TestAddPremiumVotesTx_Validate(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()

	// Valid.
	validMsgTxBytes := constants.ValidMsgAddPremiumVotesTxBytes

	// Invalid.
	invalidMsgTxBytes := constants.InvalidMsgAddPremiumVotesTxBytes

	tests := map[string]struct {
		txBytes     []byte
		expectedErr error
	}{
		"Error: ValidateBasic fails": {
			txBytes: invalidMsgTxBytes,
			expectedErr: errors.New(
				"premium votes must be sorted by perpetual id in ascending order and cannot contain" +
					" duplicates: MsgAddPremiumVotes is invalid: ValidateBasic failed on msg",
			),
		},
		"Valid: ValidateBasic passes": {
			txBytes: validMsgTxBytes,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			afst, err := process.DecodeAddPremiumVotesTx(encodingCfg.TxConfig.TxDecoder(), tc.txBytes)
			require.NoError(t, err)

			err = afst.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddPremiumVotesTx_GetMsg(t *testing.T) {
	validMsgTxBytes := constants.ValidMsgAddPremiumVotesTxBytes

	tests := map[string]struct {
		txWrapper   process.AddPremiumVotesTx
		txBytes     []byte
		expectedMsg *types.MsgAddPremiumVotes
	}{
		"Returns nil msg": {
			txWrapper: process.AddPremiumVotesTx{},
		},
		"Returns valid msg": {
			txBytes:     validMsgTxBytes,
			expectedMsg: constants.ValidMsgAddPremiumVotes,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msg sdk.Msg
			if tc.txBytes != nil {
				afst, err := process.DecodeAddPremiumVotesTx(constants.TestEncodingCfg.TxConfig.TxDecoder(), tc.txBytes)
				require.NoError(t, err)
				msg = afst.GetMsg()
			} else {
				msg = tc.txWrapper.GetMsg()
			}
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}
