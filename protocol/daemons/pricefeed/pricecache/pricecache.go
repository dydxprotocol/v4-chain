package pricecache

import (
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs.
// sDaiConversionRate is set to nil when no sDaiUpdateShould be performed.
type VeUpdatesCache struct {
	priceUpdates       PriceUpdates
	sDaiConversionRate *big.Int
	height             int64
	round              int32
	mu                 sync.RWMutex
}

type PriceUpdate struct {
	MarketId  uint32
	SpotPrice *big.Int
	PnlPrice  *big.Int
}

type PriceUpdates []PriceUpdate

func (pc *VeUpdatesCache) SetPriceUpdates(
	ctx sdk.Context,
	updates PriceUpdates,
	round int32,
) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.priceUpdates = updates
	pc.height = ctx.BlockHeight()
	pc.round = round
}

func (veCache *VeUpdatesCache) GetPriceUpdates() PriceUpdates {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.priceUpdates
}

// TODO: Look into potential issues with setting the round here
func (pc *VeUpdatesCache) SetSDaiConversionRate(
	ctx sdk.Context,
	sDaiConversionRate *big.Int,
	round int32,
) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.sDaiConversionRate = sDaiConversionRate
	pc.height = ctx.BlockHeight()
	pc.round = round
}

func (veCache *VeUpdatesCache) GetConversionRateUpdate() *big.Int {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.sDaiConversionRate
}

func (veCache *VeUpdatesCache) GetHeight() int64 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.height
}

func (veCache *VeUpdatesCache) GetRound() int32 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.round
}

func (veCache *VeUpdatesCache) HasValidValues(currBlock int64, round int32) bool {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return (veCache.height == currBlock && veCache.round == round)
}
