package pricecache

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
// Make sure to use this in thread-safe scenarios (e.g., when flow is holding lock on
// the entire app).
//
//go:generate mockery --name PriceUpdatesCache --filename mock_price_updates_cache.go
type PriceUpdatesCache interface {
	SetPriceUpdates(ctx sdk.Context, updates PriceUpdates, txHash []byte)
	GetPriceUpdates() PriceUpdates
	HasValidValues(currTxHash []byte) bool
}

// Ensure PriceUpdatesCacheImpl implements PriceUpdatesCache
var _ PriceUpdatesCache = (*PriceUpdatesCacheImpl)(nil)

type PriceUpdatesCacheImpl struct {
	priceUpdates PriceUpdates
	txHash       []byte
}

type PriceUpdate struct {
	MarketId uint32
	Price    *big.Int
}

type PriceUpdates []PriceUpdate

func (veCache *PriceUpdatesCacheImpl) SetPriceUpdates(
	ctx sdk.Context,
	updates PriceUpdates,
	txHash []byte,
) {
	veCache.priceUpdates = updates
	veCache.txHash = txHash
}

func (veCache *PriceUpdatesCacheImpl) GetPriceUpdates() PriceUpdates {
	return veCache.priceUpdates
}

func (veCache *PriceUpdatesCacheImpl) HasValidValues(currTxHash []byte) bool {
	return bytes.Equal(veCache.txHash, currTxHash)
}
