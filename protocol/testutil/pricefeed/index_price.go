package pricefeed

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	pricefeedapi "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
)

func GetTestMarketPriceUpdates(n int) (indexPrices []*pricefeedapi.MarketPriceUpdate) {
	for i := 0; i < n; i++ {
		indexPrices = append(
			indexPrices,
			&pricefeedapi.MarketPriceUpdate{
				MarketId: uint32(i),
				ExchangePrices: []*pricefeedapi.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
		)
	}
	return indexPrices
}
