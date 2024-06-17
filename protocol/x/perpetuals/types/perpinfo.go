package types

import (
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PerpInfo contains all information needed to calculate margin requirements for a perpetual.
type PerpInfo struct {
	Perpetual     Perpetual
	Price         pricestypes.MarketPrice
	LiquidityTier LiquidityTier
}
