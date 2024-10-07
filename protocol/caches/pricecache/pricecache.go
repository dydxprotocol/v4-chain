package pricecache

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
// Make sure to use this in thread-safe scenarios (e.g., when flow is holding lock on
// the entire app).
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
	veCache.priceUpdates = updates
	veCache.height = ctx.BlockHeight()
	veCache.round = round
}

func (veCache *PriceUpdatesCacheImpl) GetPriceUpdates() PriceUpdates {
	return veCache.priceUpdates
}

func (veCache *PriceUpdatesCacheImpl) GetHeight() int64 {
	return veCache.height
}

func (veCache *PriceUpdatesCacheImpl) GetRound() int32 {
	return veCache.round
}

func (veCache *PriceUpdatesCacheImpl) HasValidValues(currBlock int64, round int32) bool {
	return (veCache.height == currBlock && veCache.round == round)
}
