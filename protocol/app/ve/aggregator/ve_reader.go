package aggregator

import (
	"fmt"

	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	abci "github.com/cometbft/cometbft/abci/types"
)

func GetDaemonVotesFromBlock(
	proposal [][]byte,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) ([]Vote, error) {
	extCommitInfo, err := FetchExtCommitInfoFromProposal(proposal, extCommitCodec)
	if err != nil {
		return nil, fmt.Errorf("error fetching extended-commit-info: %w", err)
	}

	votes, err := FetchVotesFromExtCommitInfo(extCommitInfo, veCodec)
	if err != nil {
		return nil, fmt.Errorf("error fetching votes: %w", err)
	}

	return votes, nil
}

func FetchExtCommitInfoFromProposal(
	proposal [][]byte,
	extCommitCodec codec.ExtendedCommitCodec,
) (abci.ExtendedCommitInfo, error) {
	if len(proposal) <= constants.DaemonInfoIndex {
		return abci.ExtendedCommitInfo{}, fmt.Errorf("proposal slice is too short, expected at least %d elements but got %d", constants.DaemonInfoIndex+1, len(proposal))
	}

	extCommitInfoBytes := proposal[constants.DaemonInfoIndex]

	extCommitInfo, err := extCommitCodec.Decode(extCommitInfoBytes)
	if err != nil {
		return abci.ExtendedCommitInfo{}, fmt.Errorf("error decoding extended-commit-info: %w", err)
	}

	return extCommitInfo, nil
}

func FetchVotesFromExtCommitInfo(
	extCommitInfo abci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
) ([]Vote, error) {
	votes := make([]Vote, len(extCommitInfo.Votes))
	for i, voteInfo := range extCommitInfo.Votes {
		voteExtension, err := veCodec.Decode(voteInfo.VoteExtension)
		if err != nil {
			return nil, fmt.Errorf("error decoding vote-extension: %w", err)
		}

		votes[i] = Vote{
			ConsAddress:         voteInfo.Validator.Address,
			DaemonVoteExtension: voteExtension,
		}
	}

	return votes, nil
}
