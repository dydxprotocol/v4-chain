package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/ante"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/stretchr/testify/require"
)

func TestIsTimestampNonceTx(t *testing.T) {
	tests := map[string]struct {
		seqs           []uint64
		expectedResult bool
		expectedErr    bool
	}{
		"Returns false for non-ts nonce": {
			seqs:           []uint64{0},
			expectedResult: false,
			expectedErr:    false,
		},
		"Returns true for ts nonce": {
			seqs:           []uint64{keeper.TimestampNonceSequenceCutoff},
			expectedResult: true,
			expectedErr:    false,
		},
		"Returns false with no error if multisignature with regular seq number": {
			seqs:           []uint64{1, 1},
			expectedResult: false,
			expectedErr:    false,
		},
		"Returns error for multisignature with timestamp nonce": {
			seqs:           []uint64{keeper.TimestampNonceSequenceCutoff, keeper.TimestampNonceSequenceCutoff},
			expectedResult: false,
			expectedErr:    true,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize some test setup which builds a test transaction from a slice of messages.
			var reg codectypes.InterfaceRegistry
			protoCfg := authtx.NewTxConfig(codec.NewProtoCodec(reg), authtx.DefaultSignModes)
			builder := protoCfg.NewTxBuilder()
			err := builder.SetMsgs([]sdk.Msg{constants.Msg_Send}...)
			require.NoError(t, err)

			// Create signatures
			var signatures []signing.SignatureV2
			for _, seq := range tc.seqs {
				signatures = append(signatures, getSignature(seq))
			}
			err = builder.SetSignatures(signatures...)

			require.NoError(t, err)
			tx := builder.GetTx()
			ctx, _, _ := sdktest.NewSdkContextWithMultistore()

			// Invoke the function under test.
			result, err := ante.IsTimestampNonceTx(ctx, tx)

			// Assert the results.
			if tc.expectedErr {
				require.NotNil(t, err)
			}
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func getSignature(seq uint64) signing.SignatureV2 {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	return signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: seq,
	}
}
