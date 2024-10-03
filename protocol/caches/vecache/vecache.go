package vecache

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs
type VeCache struct {
	height        int64
	consAddresses map[string]struct{}
	mu            sync.RWMutex
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
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.height = ctx.BlockHeight()
	pc.consAddresses = consAddresses
}

func (pc *VeCache) GetSeenVotesInCache() map[string]struct{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.consAddresses
}

func (pc *VeCache) GetHeight() int64 {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.height
}
