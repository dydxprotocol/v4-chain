package pricecache

import (
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
//
//go:generate mockery --name PriceUpdatesCache --filename mock_price_updates_cache.go
type PriceUpdatesCache interface {
	SetPriceUpdates(ctx sdk.Context, updates PriceUpdates, round int32)
	GetPriceUpdates() PriceUpdates
	SetSDaiConversionRateAndBlockHeight(ctx sdk.Context, sDaiConversionRate *big.Int, blockHeight *big.Int, round int32)
	GetConversionRateUpdateAndBlockHeight() (*big.Int, *big.Int)
	GetHeight() int64
	GetRound() int32
	HasValidValues(currBlock int64, round int32) bool
}

// Ensure PriceUpdatesCacheImpl implements PriceUpdatesCache
var _ PriceUpdatesCache = (*PriceUpdatesCacheImpl)(nil)

// this cache is used to set prices from vote extensions in processProposal
// which are fetched in ExtendVoteHandler and PreBlocker. This is to avoid
// redundant computation on calculating stake weighthed median prices in VEs.
// sDaiConversionRate is set to nil when no sDaiUpdateShould be performed.
type PriceUpdatesCacheImpl struct {
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

func (veCache *PriceUpdatesCacheImpl) SetPriceUpdates(
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

func (veCache *PriceUpdatesCacheImpl) GetPriceUpdates() PriceUpdates {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.priceUpdates
}

// TODO: Look into potential issues with setting the round here
func (veCache *PriceUpdatesCacheImpl) SetSDaiConversionRateAndBlockHeight(
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

func (veCache *PriceUpdatesCacheImpl) GetConversionRateUpdateAndBlockHeight() (*big.Int, *big.Int) {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.sDaiConversionRate, veCache.sDAILastUpdatedBlock
}

func (veCache *PriceUpdatesCacheImpl) GetHeight() int64 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.height
}

func (veCache *PriceUpdatesCacheImpl) GetRound() int32 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.round
}

func (veCache *PriceUpdatesCacheImpl) HasValidValues(currBlock int64, round int32) bool {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return (veCache.height == currBlock && veCache.round == round)
}
