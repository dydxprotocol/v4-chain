package bindings

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
