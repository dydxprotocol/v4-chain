package types

import (
	"sort"
)

// GetMarketPremiumsMap converts the `MarketPremiums` list stored in
// `PremiumStore` to a map form `perpetualId` to `MarketPremiums`.
func (ps *PremiumStore) GetMarketPremiumsMap() map[uint32]MarketPremiums {
	ret := make(map[uint32]MarketPremiums)
	for _, marketPremiums := range ps.AllMarketPremiums {
		ret[marketPremiums.PerpetualId] = marketPremiums
	}
	return ret
}

// NewPremiumStoreFromMarketPremiumMap returns a new `PremiumStore` struct
// from a MarketPremiumMap.
func NewPremiumStoreFromMarketPremiumMap(
	m map[uint32]MarketPremiums,
	numPremiums uint32,
) *PremiumStore {
	// Get a list of sorted MarketPremiums.
	premiumList := []MarketPremiums{}
	for _, premium := range m {
		premiumList = append(premiumList, premium)
	}
	sort.Slice(premiumList, func(i, j int) bool {
		return premiumList[i].GetPerpetualId() < premiumList[j].GetPerpetualId()
	})

	return &PremiumStore{
		NumPremiums:       numPremiums,
		AllMarketPremiums: premiumList,
	}
}
