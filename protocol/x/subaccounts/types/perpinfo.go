package types

import (
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type PerpInfo struct {
	Perpetual     perptypes.Perpetual
	Price         pricestypes.MarketPrice
	LiquidityTier perptypes.LiquidityTier
}
