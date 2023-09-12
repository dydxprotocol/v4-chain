package types

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// MutableExchangeMarketConfig stores all mutable market configuration per exchange.
type MutableExchangeMarketConfig struct {
	Id ExchangeId
	// We use the keys of MarketToMarketConfig to infer which markets are supported
	// by the exchange.
	MarketToMarketConfig map[MarketId]MarketConfig
}

// Copy returns a copy of the MutableExchangeMarketConfig.
func (memc *MutableExchangeMarketConfig) Copy() *MutableExchangeMarketConfig {
	marketToMarketConfig := make(map[MarketId]MarketConfig, len(memc.MarketToMarketConfig))
	for market, config := range memc.MarketToMarketConfig {
		marketToMarketConfig[market] = config.Copy()
	}
	return &MutableExchangeMarketConfig{
		Id:                   memc.Id,
		MarketToMarketConfig: marketToMarketConfig,
	}
}

// GetMarketIds returns the ordered list of market ids supported by the exchange. This set is
// currently implicitly defined by the keys of the MarketToTicker map.
func (memc *MutableExchangeMarketConfig) GetMarketIds() []MarketId {
	return lib.GetSortedKeys[lib.Sortable[uint32]](memc.MarketToMarketConfig)
}

func (memc *MutableExchangeMarketConfig) Validate(marketConfigs []*MutableMarketConfig) error {
	marketToMutableConfig := make(map[MarketId]*MutableMarketConfig, len(marketConfigs))
	for _, mutableMarketConfig := range marketConfigs {
		marketToMutableConfig[mutableMarketConfig.Id] = mutableMarketConfig
	}

	for marketId, config := range memc.MarketToMarketConfig {
		mutableMarketConfig, exists := marketToMutableConfig[marketId]
		if !exists {
			return fmt.Errorf("no market config for market %v on exchange '%v'", marketId, memc.Id)
		}
		if err := mutableMarketConfig.Validate(); err != nil {
			return fmt.Errorf("invalid market config for market %v on exchange '%v': %w", marketId, memc.Id, err)
		}

		if config.AdjustByMarket != nil {
			if _, exists := marketToMutableConfig[*config.AdjustByMarket]; !exists {
				return fmt.Errorf(
					"no market config for adjust-by market %v used to convert market %v price on exchange '%v'",
					*config.AdjustByMarket,
					marketId,
					memc.Id,
				)
			}
		}
	}
	return nil
}

// Equal returns true if the two MutableExchangeMarketConfig objects are equal.
func (memc *MutableExchangeMarketConfig) Equal(other *MutableExchangeMarketConfig) bool {
	if memc.Id != other.Id {
		return false
	}
	if len(memc.MarketToMarketConfig) != len(other.MarketToMarketConfig) {
		return false
	}

	for market, config := range memc.MarketToMarketConfig {
		otherConfig, exists := other.MarketToMarketConfig[market]
		if !exists {
			return false
		}
		if !config.Equal(otherConfig) {
			return false
		}
	}

	return true
}
