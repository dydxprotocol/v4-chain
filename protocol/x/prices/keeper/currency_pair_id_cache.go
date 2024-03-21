package keeper

import (
	"sync"
)

// CurrencyPairIDCache handles the caching logic of currency-pairs to their corresponding IDs. This
// data-structure is thread-safe, allowing concurrent reads + synchronized writes.
type CurrencyPairIDCache struct {
	// ID -> CurrencyPair
	idToCurrencyPair map[uint64]string
	// CurrencyPair -> ID
	currencyPairToID map[string]uint64
	// lock
	sync.RWMutex
}

// NewCurrencyPairIDCache creates a new CurrencyPairIDCache
func NewCurrencyPairIDCache() *CurrencyPairIDCache {
	return &CurrencyPairIDCache{
		idToCurrencyPair: make(map[uint64]string),
		currencyPairToID: make(map[string]uint64),
	}
}

// AddCurrencyPair adds a currency pair to the cache. This method takes out a write lock on the cache
// or blocks until one is available before updating the cache.
func (c *CurrencyPairIDCache) AddCurrencyPair(id uint64, currencyPair string) {
	// acquire write lock
	c.Lock()
	defer c.Unlock()

	// update cache
	c.idToCurrencyPair[id] = currencyPair
	c.currencyPairToID[currencyPair] = id
}

// GetCurrencyPairFromID returns the currency pair from the cache
func (c *CurrencyPairIDCache) GetCurrencyPairFromID(id uint64) (string, bool) {
	// acquire read lock
	c.RLock()
	defer c.RUnlock()

	// get currency pair from cache
	currencyPair, found := c.idToCurrencyPair[id]
	return currencyPair, found
}

// GetIDForCurrencyPair returns the ID for the currency pair from the cache
func (c *CurrencyPairIDCache) GetIDForCurrencyPair(currencyPair string) (uint64, bool) {
	// acquire read lock
	c.RLock()
	defer c.RUnlock()

	// get ID for currency pair from cache
	id, found := c.currencyPairToID[currencyPair]
	return id, found
}

// Remove removes the currency-pair (by ID) from the cache. This method takes out a write lock on the
// cache or blocks until one is available before updating the cache.
func (c *CurrencyPairIDCache) Remove(id uint64) {
	// acquire write lock
	c.Lock()
	defer c.Unlock()

	// remove currency pair from cache
	currencyPair, found := c.idToCurrencyPair[id]
	if found {
		delete(c.idToCurrencyPair, id)
		delete(c.currencyPairToID, currencyPair)
	}
}
