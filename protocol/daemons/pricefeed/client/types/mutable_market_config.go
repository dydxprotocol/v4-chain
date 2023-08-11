package types

import "fmt"

// MutableMarketConfig stores the metadata that is common to a market across exchanges.
type MutableMarketConfig struct {
	Id           MarketId
	Pair         string
	Exponent     Exponent
	MinExchanges uint32
}

// Copy returns a copy of the MutableMarketConfig.
func (mmc *MutableMarketConfig) Copy() *MutableMarketConfig {
	return &MutableMarketConfig{
		Id:           mmc.Id,
		Pair:         mmc.Pair,
		Exponent:     mmc.Exponent,
		MinExchanges: mmc.MinExchanges,
	}
}

func (mmc *MutableMarketConfig) Validate() error {
	if mmc.Pair == "" {
		return fmt.Errorf("pair cannot be empty")
	}
	if mmc.MinExchanges == 0 {
		return fmt.Errorf("min exchanges cannot be 0")
	}

	return nil
}
