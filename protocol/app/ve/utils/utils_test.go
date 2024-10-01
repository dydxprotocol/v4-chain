package ve_utils_test

import (
	"math/big"
	"testing"

	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	ethosutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ethos"
	cometabcitypes "github.com/cometbft/cometbft/abci/types"

	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
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
			ctx := types.Context{}.
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
	ctx := sdk.Context{}
	mockValidatorStore := &mocks.CCValidatorStore{}

	testValPower := int64(1000)
	testVal := ethosutils.BuildCCValidator("alice", testValPower)
	testValConsAddr := sdk.ConsAddress(testVal.Address)

	mockVote := cometabcitypes.ExtendedVoteInfo{
		Validator: cometabcitypes.Validator{
			Address: testValConsAddr,
			Power:   testValPower,
		},
	}

	t.Run("Successful public key retrieval", func(t *testing.T) {
		mockValidatorStore.On("GetCCValidator", mock.Anything, mock.Anything).Return(testVal, true).Once()

		pubKey, err := veutils.GetValCmtPubKeyFromVote(ctx, mockVote, mockValidatorStore)
		require.NoError(t, err)
		require.NotNil(t, pubKey)

		expectedPubKey := constants.AliceEthosPubKey
		require.Equal(t, expectedPubKey.Bytes(), pubKey.Bytes())
	})

	t.Run("Validator not found", func(t *testing.T) {
		unknownConsAddr := constants.BobEthosConsAddress
		unknownVote := cometabcitypes.ExtendedVoteInfo{
			Validator: cometabcitypes.Validator{
				Address: unknownConsAddr,
				Power:   100,
			},
		}

		mockValidatorStore.On("GetCCValidator", mock.Anything, mock.Anything).Return(ccvtypes.CrossChainValidator{}, false).Once()

		_, err := veutils.GetValCmtPubKeyFromVote(ctx, unknownVote, mockValidatorStore)
		require.Error(t, err)
		require.IsType(t, &veutils.ValidatorNotFoundError{}, err)
	})

	t.Run("Invalid public key: validator not found error", func(t *testing.T) {
		invalidPubkey := &codectypes.Any{
			TypeUrl: "invalid/type/url",
			Value:   []byte(""),
		}

		invalidVal := ccvtypes.CrossChainValidator{
			Address: testVal.Address,
			Pubkey:  invalidPubkey,
			Power:   testVal.Power,
		}

		mockValidatorStore.On("GetCCValidator", mock.Anything, mock.Anything).Return(invalidVal, true).Once()

		_, err := veutils.GetValCmtPubKeyFromVote(ctx, mockVote, mockValidatorStore)
		require.Error(t, err)
		require.IsType(t, &veutils.ValidatorNotFoundError{}, err)
	})

	t.Run("Invalid public key: public key is nil", func(t *testing.T) {
		invalidVal := ccvtypes.CrossChainValidator{
			Address: testVal.Address,
			Power:   testVal.Power,
		}
		mockValidatorStore.On("GetCCValidator", mock.Anything, mock.Anything).Return(invalidVal, true).Once()

		_, err := veutils.GetValCmtPubKeyFromVote(ctx, mockVote, mockValidatorStore)
		require.Error(t, err)
		require.IsType(t, &veutils.ValidatorNotFoundError{}, err)
	})
}

func TestGetPubKeyByConsAddr(t *testing.T) {
	testValidator := ethosutils.BuildCCValidator("alice", 100)

	tests := []struct {
		name        string
		ccValidator ccvtypes.CrossChainValidator
		expectError bool
		errorString string
	}{
		{
			name:        "valid public key",
			ccValidator: testValidator,
			expectError: false,
		},
		{
			name: "nil public key",
			ccValidator: ccvtypes.CrossChainValidator{
				Address: testValidator.Address,
				Pubkey:  nil,
				Power:   testValidator.Power,
			},
			expectError: true,
			errorString: "public key is nil",
		},
		{
			name: "invalid public key",
			ccValidator: ccvtypes.CrossChainValidator{
				Address: testValidator.Address,
				Pubkey:  codectypes.UnsafePackAny(42),
				Power:   testValidator.Power,
			},
			expectError: true,
			errorString: "could not get pubkey for val",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := veutils.GetPubKeyByConsAddr(tc.ccValidator)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorString)
				require.Equal(t, cmtprotocrypto.PublicKey{}, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.IsType(t, cmtprotocrypto.PublicKey{}, result)
				expectedPubKeyProto := tc.ccValidator.Pubkey.GetCachedValue().(cryptotypes.PubKey)
				expectedPubKey, err := cryptocodec.ToCmtProtoPublicKey(expectedPubKeyProto)
				require.NoError(t, err)
				require.Equal(t, expectedPubKey, result)
			}
		})
	}
}

func mustEncodePrice(t *testing.T, price *big.Int) []byte {
	encoded, err := price.GobEncode()
	require.NoError(t, err)
	return encoded
}
