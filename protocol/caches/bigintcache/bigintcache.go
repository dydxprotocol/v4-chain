package bigintcache

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceUpdatesCache is an interface that defines the methods for caching price updates
// Make sure to use this in thread-safe scenarios (e.g., when flow is holding lock on
// the entire app).
//
//go:generate mockery --name BigIntCache --filename mock_bigint_cache.go
type BigIntCache interface {
	SetValue(ctx sdk.Context, value *big.Int, txHash []byte)
	GetValue() *big.Int
	HasValidValue(currTxHash []byte) bool
}

// Ensure BigIntCacheImpl implements BigIntCache
var _ BigIntCache = (*BigIntCacheImpl)(nil)

type BigIntCacheImpl struct {
	value  *big.Int
	txHash []byte
}

func (veCache *BigIntCacheImpl) SetValue(
	ctx sdk.Context,
	value *big.Int,
	txHash []byte,
) {
	veCache.value = value
	veCache.txHash = txHash

}

func (veCache *BigIntCacheImpl) GetValue() *big.Int {
	return veCache.value
}

func (veCache *BigIntCacheImpl) HasValidValue(currTxHash []byte) bool {
	return bytes.Equal(veCache.txHash, currTxHash)
}
