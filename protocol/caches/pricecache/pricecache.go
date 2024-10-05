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
	GetHeight() int64
	GetRound() int32
	HasValidValues(currBlock int64, round int32) bool
}

// Ensure PriceUpdatesCacheImpl implements PriceUpdatesCache
var _ PriceUpdatesCache = (*PriceUpdatesCacheImpl)(nil)

type PriceUpdatesCacheImpl struct {
	priceUpdates PriceUpdates
	height       int64
	round        int32
	mu           sync.RWMutex
}

type PriceUpdate struct {
	MarketId uint32
	Price    *big.Int
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
