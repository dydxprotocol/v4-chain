package types

import "fmt"

// MutableMarketConfig stores the metadata that is common to a market across exchanges.
type MutableMarketConfig struct {
	Id       MarketId
	Pair     string
	Exponent Exponent
}

// Copy returns a copy of the MutableMarketConfig.
func (mmc *MutableMarketConfig) Copy() *MutableMarketConfig {
	return &MutableMarketConfig{
		Id:       mmc.Id,
		Pair:     mmc.Pair,
		Exponent: mmc.Exponent,
	}
}

func (mmc *MutableMarketConfig) Validate() error {
	if mmc.Pair == "" {
		return fmt.Errorf("pair cannot be empty")
	}

	return nil
}
