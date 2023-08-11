package types

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
	numPerpetuals uint32,
	numPremiums uint32,
) *PremiumStore {
	ret := PremiumStore{
		NumPremiums: numPremiums,
	}
	for perpId := uint32(0); perpId < numPerpetuals; perpId++ {
		marketPremiums, found := m[perpId]
		if !found {
			// `PrmeiumStore` is used as a sparse matrix, so a perpetual Id not
			// being found inherently means all premiums for the market were zeros.
			continue
		}
		ret.AllMarketPremiums = append(ret.AllMarketPremiums,
			marketPremiums,
		)
	}
	return &ret
}
