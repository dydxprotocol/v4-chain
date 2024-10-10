package vecache

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs.
// Make sure to use this in thread-safe scenarios (e.g., when flow is holding lock on
// the entire app).
type VeCache struct {
	height        int64
	consAddresses map[string]struct{}
}

func NewVECache() *VeCache {
	return &VeCache{
		height:        0,
		consAddresses: make(map[string]struct{}),
	}
}

func (pc *VeCache) SetSeenVotesInCache(
	ctx sdk.Context,
	consAddresses map[string]struct{},
) {
	pc.height = ctx.BlockHeight()
	pc.consAddresses = consAddresses
}

func (pc *VeCache) GetSeenVotesInCache() map[string]struct{} {
	return pc.consAddresses
}

func (pc *VeCache) GetHeight() int64 {
	return pc.height
}
