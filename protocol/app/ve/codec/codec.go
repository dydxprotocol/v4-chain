package codec

import (
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
)

type VoteExtensionCodec interface {
	// Encode encodes the vote extension into a byte array.
	Encode(ve vetypes.DeamonVoteExtension) ([]byte, error)

	// Decode decodes the vote extension from a byte array.
	Decode([]byte) (vetypes.DeamonVoteExtension, error)
}

type ExtendedCommitCodec interface {
	// Encode encodes the extended commit info into a byte array.
	Encode(cometabci.ExtendedCommitInfo) ([]byte, error)

	// Decode decodes the extended commit info from a byte array.
	Decode([]byte) (cometabci.ExtendedCommitInfo, error)
}

func NewDefaultVoteExtensionCodec() *DefaultVoteExtensionCodec {
	return &DefaultVoteExtensionCodec{}
}

// DefaultVoteExtensionCodec is the default implementation of VoteExtensionCodec. It uses the
// vanilla implementations of Unmarshal / Marshal under the hood.
type DefaultVoteExtensionCodec struct{}

func (codec *DefaultVoteExtensionCodec) Encode(ve vetypes.DeamonVoteExtension) ([]byte, error) {
	return ve.Marshal()
}

func (codec *DefaultVoteExtensionCodec) Decode(bz []byte) (vetypes.DeamonVoteExtension, error) {
	var ve vetypes.DeamonVoteExtension
	return ve, ve.Unmarshal(bz)
}

// DefaultExtendedCommitCodec is the default implementation of ExtendedCommitCodec. It uses the
// vanilla implementations of Unmarshal / Marshal under the hood.
type DefaultExtendedCommitCodec struct{}

// NewDefaultExtendedCommitCodec returns a new DefaultExtendedCommitCodec.
func NewDefaultExtendedCommitCodec() *DefaultExtendedCommitCodec {
	return &DefaultExtendedCommitCodec{}
}

func (codec *DefaultExtendedCommitCodec) Encode(extendedCommitInfo cometabci.ExtendedCommitInfo) ([]byte, error) {
	return extendedCommitInfo.Marshal()
}

func (codec *DefaultExtendedCommitCodec) Decode(bz []byte) (cometabci.ExtendedCommitInfo, error) {
	if len(bz) == 0 {
		return cometabci.ExtendedCommitInfo{}, nil
	}

	var extendedCommitInfo cometabci.ExtendedCommitInfo
	return extendedCommitInfo, extendedCommitInfo.Unmarshal(bz)
}
