package ve_utils_test

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/math"

	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	valutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/staking"
	cometabcitypes "github.com/cometbft/cometbft/abci/types"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAreVEEnabled(t *testing.T) {
	tests := []struct {
		name              string
		consensusParams   *cometbftproto.ConsensusParams
		blockHeight       int64
		expectedVEEnabled bool
	}{
		{
			name: "VE disabled: nil ABCI",
			consensusParams: &cometbftproto.ConsensusParams{
				Abci: nil,
			},
			blockHeight:       10,
			expectedVEEnabled: false,
		},
		{
			name: "VE disabled: VoteExtensionsEnableHeight is 0",
			consensusParams: &cometbftproto.ConsensusParams{
				Abci: &cometbftproto.ABCIParams{
					VoteExtensionsEnableHeight: 0,
				},
			},
			blockHeight:       10,
			expectedVEEnabled: false,
		},
		{
			name: "VE disabled: BlockHeight <= 1",
			consensusParams: &cometbftproto.ConsensusParams{
				Abci: &cometbftproto.ABCIParams{
					VoteExtensionsEnableHeight: 5,
				},
			},
			blockHeight:       1,
			expectedVEEnabled: false,
		},
		{
			name: "VE enabled: BlockHeight > VoteExtensionsEnableHeight",
			consensusParams: &cometbftproto.ConsensusParams{
				Abci: &cometbftproto.ABCIParams{
					VoteExtensionsEnableHeight: 5,
				},
			},
			blockHeight:       10,
			expectedVEEnabled: true,
		},
		{
			name: "VE disabled: BlockHeight <= VoteExtensionsEnableHeight",
			consensusParams: &cometbftproto.ConsensusParams{
				Abci: &cometbftproto.ABCIParams{
					VoteExtensionsEnableHeight: 10,
				},
			},
			blockHeight:       5,
			expectedVEEnabled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := sdktypes.Context{}.
				WithConsensusParams(*tc.consensusParams).
				WithBlockHeight(tc.blockHeight)

			result := veutils.AreVEEnabled(ctx)
			assert.Equal(t, tc.expectedVEEnabled, result)
		})
	}
}

func TestGetPriceFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		id      uint32
		input   []byte
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "Valid positive price",
			id:      1,
			input:   mustEncodePrice(t, big.NewInt(100)),
			want:    big.NewInt(100),
			wantErr: false,
		},
		{
			name:    "Valid zero price",
			id:      2,
			input:   mustEncodePrice(t, big.NewInt(0)),
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name:    "Invalid negative price",
			id:      3,
			input:   mustEncodePrice(t, big.NewInt(-100)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid input",
			id:      4,
			input:   []byte("invalid"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := veutils.GetPriceFromBytes(tc.id, tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetVEDecodedPrice(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "Valid positive price",
			input:   mustEncodePrice(t, big.NewInt(100)),
			want:    big.NewInt(100),
			wantErr: false,
		},
		{
			name:    "Valid zero price",
			input:   mustEncodePrice(t, big.NewInt(0)),
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name:    "Invalid negative price",
			input:   mustEncodePrice(t, big.NewInt(-100)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid input",
			input:   []byte("invalid"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := veutils.GetVEDecodedPrice(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetValCmtPubKeyFromVote(t *testing.T) {
	ctx := sdktypes.Context{}
	mockValidatorStore := &mocks.ValidatorStore{}

	testValPower := int64(1000)
	testVal := valutils.BuildTestValidator("alice", math.NewInt(testValPower))
	testValPubKey, err := testVal.ConsPubKey()
	require.NoError(t, err)
	protoPubKey, err := cryptocodec.ToCmtProtoPublicKey(testValPubKey)
	require.NoError(t, err)

	testValConsAddr := constants.AliceConsAddress

	mockVote := cometabcitypes.ExtendedVoteInfo{
		Validator: cometabcitypes.Validator{
			Address: testValConsAddr,
			Power:   testValPower,
		},
	}

	t.Run("Successful public key retrieval", func(t *testing.T) {

		mockValidatorStore.On("GetPubKeyByConsAddr", mock.Anything, testValConsAddr).Return(protoPubKey, nil).Once()

		pubKey, err := veutils.GetValPubKeyFromVote(ctx, mockVote, mockValidatorStore)
		require.NoError(t, err)
		require.NotNil(t, pubKey)

		expectedPubKey := constants.AlicePubKey
		require.Equal(t, expectedPubKey.Bytes(), pubKey.Bytes())
	})

	t.Run("Validator not found", func(t *testing.T) {
		unknownConsAddr := constants.BobConsAddress
		unknownVote := cometabcitypes.ExtendedVoteInfo{
			Validator: cometabcitypes.Validator{
				Address: unknownConsAddr,
				Power:   100,
			},
		}

		mockValidatorStore.On("GetPubKeyByConsAddr", mock.Anything, mock.Anything).Return(cmtprotocrypto.PublicKey{}, fmt.Errorf("error")).Once()

		_, err := veutils.GetValPubKeyFromVote(ctx, unknownVote, mockValidatorStore)
		require.Error(t, err)
		require.IsType(t, &veutils.ValidatorNotFoundError{}, err)
	})

	t.Run("Invalid public key: validator not found error", func(t *testing.T) {
		invalidCmtPubKey := cmtprotocrypto.PublicKey{
			Sum: &cmtprotocrypto.PublicKey_Ed25519{
				Ed25519: []byte("invalid"),
			},
		}

		mockValidatorStore.On("GetPubKeyByConsAddr", mock.Anything, mock.Anything).Return(invalidCmtPubKey, nil).Once()

		_, err := veutils.GetValPubKeyFromVote(ctx, mockVote, mockValidatorStore)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to convert validator")
	})
}

func TestGetVEEncodedPrice(t *testing.T) {
	tests := map[string]struct {
		price           *big.Int
		expectedVEBytes []byte
		expectedError   bool
	}{
		"Positive price": {
			price:           big.NewInt(100),
			expectedVEBytes: mustEncodePrice(t, big.NewInt(100)),
			expectedError:   false,
		},
		"Zero price": {
			price:           big.NewInt(0),
			expectedVEBytes: mustEncodePrice(t, big.NewInt(0)),
			expectedError:   false,
		},
		"Negative price": {
			price:           big.NewInt(-100),
			expectedVEBytes: nil,
			expectedError:   true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			veBytes, err := veutils.GetVEEncodedPrice(tc.price)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedVEBytes, veBytes)
			}
		})
	}
}

func TestMarshalDelimited(t *testing.T) {
	tests := map[string]struct {
		input          proto.Message
		expectedOutput []byte
		expectedError  bool
	}{
		"Valid message": {
			input: &cometbftproto.BlockID{
				Hash: []byte("testhash"),
			},
			expectedOutput: mustEncodeDelimited(t, &cometbftproto.BlockID{
				Hash: []byte("testhash"),
			}),
			expectedError: false,
		},
		"valid canconical vote": {
			input: &cometbftproto.CanonicalVoteExtension{
				Extension: []byte("test"),
				Height:    1,
				Round:     1,
				ChainId:   "test",
			},
			expectedOutput: mustEncodeDelimited(t, &cometbftproto.CanonicalVoteExtension{
				Extension: []byte("test"),
				Height:    1,
				Round:     1,
				ChainId:   "test",
			}),
			expectedError: false,
		},
		"Nil message": {
			input:          nil,
			expectedOutput: nil,
			expectedError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := veutils.MarshalDelimited(tc.input)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func mustEncodeDelimited(t *testing.T, msg proto.Message) []byte {
	var buf bytes.Buffer
	err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg)
	require.NoError(t, err)
	return buf.Bytes()
}

func mustEncodePrice(t *testing.T, price *big.Int) []byte {
	encoded, err := price.GobEncode()
	require.NoError(t, err)
	return encoded
}
