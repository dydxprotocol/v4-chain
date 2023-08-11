package constants

import (
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

var (
	ExchangeFeed_Coinbase = pricestypes.ExchangeFeed{
		Id:   0,
		Name: CoinbaseExchangeName,
		Memo: "Test Coinbase memo",
	}
)
