package ve_testutils

import (
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetExtCommitInfoBz(
	consAddr sdk.ConsAddress,
	prices map[uint32][]byte,
	voteCodec vecodec.VoteExtensionCodec,
	extCodec vecodec.ExtendedCommitCodec,
) []byte {
	commitInfo, err := CreateExtendedVoteInfoWithPower(consAddr, 1, prices, voteCodec)
	if err != nil {
		return nil
	}

	_, bz, err := CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{commitInfo}, extCodec)
	if err != nil {
		return nil
	}
	return bz
}

// CreateExtendedCommitInfo creates an extended commit info with the given commit info.
func CreateExtendedCommitInfo(
	commitInfo []cometabci.ExtendedVoteInfo,
	codec vecodec.ExtendedCommitCodec,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	extendedCommitInfo := cometabci.ExtendedCommitInfo{
		Votes: commitInfo,
	}

	bz, err := codec.Encode(extendedCommitInfo)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, err
	}

	return extendedCommitInfo, bz, nil
}

// CreateExtendedVoteInfo creates an extended vote info with the given prices, timestamp and height.
func CreateExtendedVoteInfo(
	consAddr sdk.ConsAddress,
	prices map[uint32][]byte,
	codec vecodec.VoteExtensionCodec,
) (cometabci.ExtendedVoteInfo, error) {
	return CreateExtendedVoteInfoWithPower(consAddr, 1, prices, codec)
}

// CreateExtendedVoteInfoWithPower CreateExtendedVoteInfo creates an extended vote info
// with the given power, prices, timestamp and height.
func CreateExtendedVoteInfoWithPower(
	consAddr sdk.ConsAddress,
	power int64,
	prices map[uint32][]byte,
	codec vecodec.VoteExtensionCodec,
) (cometabci.ExtendedVoteInfo, error) {
	ve, err := CreateVoteExtensionBytes(prices, codec)

	if err != nil {
		return cometabci.ExtendedVoteInfo{}, err
	}
	voteInfo := cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: consAddr,
			Power:   power,
		},
		VoteExtension: ve,
		BlockIdFlag:   cometproto.BlockIDFlagCommit,
	}

	return voteInfo, nil
}

// CreateVoteExtensionBytes creates a vote extension bytes with the given prices, timestamp and height.
func CreateVoteExtensionBytes(
	prices map[uint32][]byte,
	codec vecodec.VoteExtensionCodec,
) ([]byte, error) {
	voteExtension := CreateVoteExtension(prices)
	voteExtensionBz, err := codec.Encode(voteExtension)
	if err != nil {
		return nil, err
	}

	return voteExtensionBz, nil
}

// CreateVoteExtension creates a vote extension with the given prices, timestamp and height.
func CreateVoteExtension(
	prices map[uint32][]byte,
) vetypes.DaemonVoteExtension {

	return vetypes.DaemonVoteExtension{
		Prices: prices,
	}
}
