package types

import (
	"fmt"
	"sort"
)

// MutableExchangeMarketConfig stores all mutable market configuration per exchange.
type MutableExchangeMarketConfig struct {
	Id ExchangeId
	// We use the keys of MarketToTicker map to infer which markets are supported
	// by the exchange.
	MarketToTicker map[MarketId]string
}

// Copy returns a copy of the MutableExchangeMarketConfig.
func (memc *MutableExchangeMarketConfig) Copy() *MutableExchangeMarketConfig {
	marketToTicker := make(map[MarketId]string, len(memc.MarketToTicker))
	for marketId, ticker := range memc.MarketToTicker {
		marketToTicker[marketId] = ticker
	}
	return &MutableExchangeMarketConfig{
		Id:             memc.Id,
		MarketToTicker: marketToTicker,
	}
}

// GetMarketIds returns the ordered list of market ids supported by the exchange. This set is
// currently implicitly defined by the keys of the MarketToTicker map.
func (memc *MutableExchangeMarketConfig) GetMarketIds() []MarketId {
	marketIds := make([]MarketId, 0, len(memc.MarketToTicker))
	for marketId := range memc.MarketToTicker {
		marketIds = append(marketIds, marketId)
	}
	sort.Slice(marketIds, func(i, j int) bool {
		return marketIds[i] < marketIds[j]
	})
	return marketIds
}

func (memc *MutableExchangeMarketConfig) Validate(marketConfigs []*MutableMarketConfig) error {
	marketIdToConfig := make(map[MarketId]*MutableMarketConfig, len(marketConfigs))
	for _, marketConfig := range marketConfigs {
		marketIdToConfig[marketConfig.Id] = marketConfig
	}

	for marketId := range memc.MarketToTicker {
		marketConfig, exists := marketIdToConfig[marketId]
		if !exists {
			return fmt.Errorf("no market config for market %v on exchange '%v'", marketId, memc.Id)
		}
		if err := marketConfig.Validate(); err != nil {
			return fmt.Errorf("invalid market config for market %v on exchange '%v': %w", marketId, memc.Id, err)
		}
	}
	return nil
}

// Equal returns true if the two MutableExchangeMarketConfig objects are equal.
func (memc *MutableExchangeMarketConfig) Equal(other *MutableExchangeMarketConfig) bool {
	if memc.Id != other.Id {
		return false
	}
	if len(memc.MarketToTicker) != len(other.MarketToTicker) {
		return false
	}
	for marketId, ticker := range memc.MarketToTicker {
		otherTicker, exists := other.MarketToTicker[marketId]
		if !exists {
			return false
		}
		if ticker != otherTicker {
			return false
		}
	}
	return true
}
