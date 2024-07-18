package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

const maxTimestampNonceArrSize = 20
const maxTimeInPastMs = 30000
const maxTimeInFutureMs = 30000

var InitialTimestampNonceDetails = &types.TimestampNonceDetails{
	MaxEjectedNonce: 0,
	TimestampNonces: []uint64{},
}

func (k Keeper) DeepCopyTimestampNonceDetails(details *types.TimestampNonceDetails) *types.TimestampNonceDetails {
	if details == nil {
		return nil
	}

	copyDetails := &types.TimestampNonceDetails{
		MaxEjectedNonce: details.MaxEjectedNonce,
		TimestampNonces: make([]uint64, len(details.TimestampNonces)),
	}

	// Copy the slice elements
	copy(copyDetails.TimestampNonces, details.TimestampNonces)

	return copyDetails
}
