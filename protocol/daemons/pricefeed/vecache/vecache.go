package vecache

import (
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VeUpdatesCache is an interface that defines the methods for caching vote extension updates
//
//go:generate mockery --name VeUpdatesCache --filename mock_ve_updates_cache.go
type VeUpdatesCache interface {
	SetPriceUpdates(ctx sdk.Context, updates PriceUpdates, round int32)
	GetPriceUpdates() PriceUpdates
	SetSDaiConversionRateAndBlockHeight(ctx sdk.Context, sDaiConversionRate *big.Int, blockHeight *big.Int, round int32)
	GetConversionRateUpdateAndBlockHeight() (*big.Int, *big.Int)
	GetHeight() int64
	GetRound() int32
	HasValidValues(currBlock int64, round int32) bool
}

// Ensure VeUpdatesCacheImpl implements VeUpdatesCache
var _ VeUpdatesCache = (*VeUpdatesCacheImpl)(nil)

// this cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs.
// sDaiConversionRate is set to nil when no sDaiUpdateShould be performed.
type VeUpdatesCacheImpl struct {
	priceUpdates         PriceUpdates
	sDaiConversionRate   *big.Int
	sDAILastUpdatedBlock *big.Int
	height               int64
	round                int32
	mu                   sync.RWMutex
}

type PriceUpdate struct {
	MarketId  uint32
	SpotPrice *big.Int
	PnlPrice  *big.Int
}

type PriceUpdates []PriceUpdate

func (veCache *VeUpdatesCacheImpl) SetPriceUpdates(
	ctx sdk.Context,
	updates PriceUpdates,
	round int32,
) {
	veCache.mu.Lock()
	defer veCache.mu.Unlock()
	veCache.priceUpdates = updates
	veCache.height = ctx.BlockHeight()
	veCache.round = round
}

func (veCache *VeUpdatesCacheImpl) GetPriceUpdates() PriceUpdates {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.priceUpdates
}

// TODO: Look into potential issues with setting the round here
func (veCache *VeUpdatesCacheImpl) SetSDaiConversionRateAndBlockHeight(
	ctx sdk.Context,
	sDaiConversionRate *big.Int,
	blockHeight *big.Int,
	round int32,
) {
	veCache.mu.Lock()
	defer veCache.mu.Unlock()
	veCache.sDaiConversionRate = sDaiConversionRate
	veCache.sDAILastUpdatedBlock = blockHeight
	veCache.height = ctx.BlockHeight()
	veCache.round = round
}

func (veCache *VeUpdatesCacheImpl) GetConversionRateUpdateAndBlockHeight() (*big.Int, *big.Int) {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.sDaiConversionRate, veCache.sDAILastUpdatedBlock
}

func (veCache *VeUpdatesCacheImpl) GetHeight() int64 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.height
}

func (veCache *VeUpdatesCacheImpl) GetRound() int32 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.round
}

func (veCache *VeUpdatesCacheImpl) HasValidValues(currBlock int64, round int32) bool {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return (veCache.height == currBlock && veCache.round == round)
}
