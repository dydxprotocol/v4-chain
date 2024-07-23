package pricecache

import (
	"sync"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs
type PriceCache struct {
	priceUpdates []types.MarketPriceUpdates_MarketPriceUpdate
	height       int64
	round        int32
	mu           sync.RWMutex
}

func (pc *PriceCache) SetPriceUpdates(
	ctx sdk.Context,
	updates []types.MarketPriceUpdates_MarketPriceUpdate,
	round int32,
) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.priceUpdates = updates
	pc.height = ctx.BlockHeight()
}

func (pc *PriceCache) GetPriceUpdates() []types.MarketPriceUpdates_MarketPriceUpdate {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.priceUpdates
}

func (pc *PriceCache) HasValidPrices(currBlock int64, round int32) bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return (pc.height == currBlock && pc.round == round)
}
