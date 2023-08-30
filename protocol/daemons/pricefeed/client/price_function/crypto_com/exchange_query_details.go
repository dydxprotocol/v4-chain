package crypto_com

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	CryptoComDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_CRYPTO_COM,
		Url:           "https://api.crypto.com/v2/public/get-ticker",
		PriceFunction: CryptoComPriceFunction,
		IsMultiMarket: true,
	}
)
