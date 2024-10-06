package bigintcache

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
// Make sure to use this in thread-safe scenarios (e.g., when flow is holding lock on
// the entire app).
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
}

func (veCache *BigIntCacheImpl) SetValue(
	ctx sdk.Context,
	value *big.Int,
	round int32,
) {
	veCache.value = value
	veCache.height = ctx.BlockHeight()
	veCache.round = round
}

func (veCache *BigIntCacheImpl) GetValue() *big.Int {
	return veCache.value
}

func (veCache *BigIntCacheImpl) GetHeight() int64 {
	return veCache.height
}

func (veCache *BigIntCacheImpl) GetRound() int32 {
	return veCache.round
}

func (veCache *BigIntCacheImpl) HasValidValue(currBlock int64, round int32) bool {
	return (veCache.height == currBlock && veCache.round == round)
}
