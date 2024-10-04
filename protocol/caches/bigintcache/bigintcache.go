package bigintcache

import (
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
//
//go:generate mockery --name BigIntCache --filename mock_bigint_cache.go
type BigIntCache interface {
	SetValue(ctx sdk.Context, value *big.Int, round int32)
	GetValue() *big.Int
	GetHeight() int64
	GetRound() int32
	HasValidValue(currBlock int64, round int32) bool
}

// Ensure BigIntCacheImpl implements BigIntCache
var _ BigIntCache = (*BigIntCacheImpl)(nil)

type BigIntCacheImpl struct {
	value  *big.Int
	height int64
	round  int32
	mu     sync.RWMutex
}

func (veCache *BigIntCacheImpl) SetValue(
	ctx sdk.Context,
	value *big.Int,
	round int32,
) {
	veCache.mu.Lock()
	defer veCache.mu.Unlock()
	veCache.value = value
	veCache.height = ctx.BlockHeight()
	veCache.round = round
}

func (veCache *BigIntCacheImpl) GetValue() *big.Int {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.value
}

func (veCache *BigIntCacheImpl) GetHeight() int64 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.height
}

func (veCache *BigIntCacheImpl) GetRound() int32 {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return veCache.round
}

func (veCache *BigIntCacheImpl) HasValidValue(currBlock int64, round int32) bool {
	veCache.mu.RLock()
	defer veCache.mu.RUnlock()
	return (veCache.height == currBlock && veCache.round == round)
}
