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
	ret := PremiumStore{
		NumPremiums: numPremiums,
	}

	// Get a list of sorted perpetual Ids.
	perpetualIds := []uint32{}
	for perpId := range m {
		perpetualIds = append(perpetualIds, perpId)
	}
	sort.Slice(perpetualIds, func(i, j int) bool {
		return perpetualIds[i] < perpetualIds[j]
	})

	// Iterate through the sorted list of perpetual Ids and add market premiums.
	for _, perpId := range perpetualIds {
		ret.AllMarketPremiums = append(ret.AllMarketPremiums,
			m[perpId],
		)
	}
	return &ret
}
