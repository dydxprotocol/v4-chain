package bindings

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type MarketPriceRequestWrapper struct {
	MarketPrice pricestypes.MarketPrice `json:"market_price"`
}

type SubaccountRequestWrapper struct {
	Subaccount satypes.QueryGetSubaccountRequest `json:"subaccount"`
}

type PerpeutalClobDetailsRequestWrapper struct {
	PerpetualClobDetails clobtypes.QueryGetPerpetualClobDetailsRequest `json:"perpetual_clob_details"`
}

type LiquidityTiersRequestWrapper struct {
	LiquidityTiers perptypes.QueryAllLiquidityTiersRequest `json:"liquidity_tiers"`
}
