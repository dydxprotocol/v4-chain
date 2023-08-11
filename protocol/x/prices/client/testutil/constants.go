package testutil

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/testutil/constants"
)

const (
	responseErr = 503
	responseOk  = 200
)

var (
	// Error responses.
	errResponse_Binance = Response{
		ResponseCode: responseErr,
		Response:     BinanceResponse{},
	}
	errResponse_Bitfinex = Response{
		ResponseCode: responseErr,
		Response:     BitfinexResponse{},
	}

	// OK responses.
	responseBinance_100 = Response{
		ResponseCode: responseOk,
		Response:     NewBinanceResponse("99", "100", "101"),
	}
	responseBinance_101 = Response{
		ResponseCode: responseOk,
		Response:     NewBinanceResponse("100", "101", "102"),
	}
	responseBinance_9000 = Response{
		ResponseCode: responseOk,
		Response:     NewBinanceResponse("8999", "9000", "9001"),
	}
	responseBinance_9001 = Response{
		ResponseCode: responseOk,
		Response:     NewBinanceResponse("9000", "9001", "9002"),
	}

	responseBitfinex_102 = Response{
		ResponseCode: responseOk,
		Response:     NewBitfinexResponse(101.0, 102.0, 103.0),
	}
	responseBitfinex_9002 = Response{
		ResponseCode: responseOk,
		Response:     NewBitfinexResponse(9001.0, 9002.0, 9003.0),
	}

	// Exchange responses.
	errResponses_AllExchanges = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: errResponse_Binance,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: errResponse_Binance,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: errResponse_Bitfinex,
	}
	validResponses_AllExchanges_Median101 = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: responseBinance_100,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: responseBinance_101,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: responseBitfinex_102,
	}
	validResponses_AllExchanges_Median9001 = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: responseBinance_9000,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: responseBinance_9001,
		{
			exchangeId:   exchange_common.EXCHANGE_FEED_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: responseBitfinex_9002,
	}

	// All Error responses.
	AllErrorResponses = MarketToExchangeResponse{
		marketToExchangesResponse: map[MarketIdAndName]ExchangeResponse{
			{
				marketId:   exchange_common.MARKET_BTC_USD,
				marketName: constants.BtcUsdPair,
			}: {
				exchangeToResponse: errResponses_AllExchanges,
			},
			{
				marketId:   exchange_common.MARKET_ETH_USD,
				marketName: constants.EthUsdPair,
			}: {
				exchangeToResponse: errResponses_AllExchanges,
			},
		},
	}

	// Mix of Valid and Error responses.
	MixedResponses = MarketToExchangeResponse{
		marketToExchangesResponse: map[MarketIdAndName]ExchangeResponse{
			{
				marketId:   exchange_common.MARKET_BTC_USD,
				marketName: constants.BtcUsdPair,
			}: {
				exchangeToResponse: map[ExchangeIdAndName]Response{
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE,
						exchangeName: constants.BinanceExchangeName,
					}: responseBinance_100, // Valid.
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE_US,
						exchangeName: constants.BinanceUSExchangeName,
					}: errResponse_Binance, // Error.
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BITFINEX,
						exchangeName: constants.BitfinexExchangeName,
					}: errResponse_Bitfinex, // Error.
				},
			},
			{
				marketId:   exchange_common.MARKET_ETH_USD,
				marketName: constants.EthUsdPair,
			}: {
				exchangeToResponse: map[ExchangeIdAndName]Response{
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE,
						exchangeName: constants.BinanceExchangeName,
					}: responseBinance_9000, // Valid.
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BINANCE_US,
						exchangeName: constants.BinanceUSExchangeName,
					}: errResponse_Binance, // Error.
					{
						exchangeId:   exchange_common.EXCHANGE_FEED_BITFINEX,
						exchangeName: constants.BitfinexExchangeName,
					}: responseBitfinex_9002, // Valid.
				}},
		},
	}

	// All Valid responses.
	AllValidResponses = MarketToExchangeResponse{
		marketToExchangesResponse: map[MarketIdAndName]ExchangeResponse{
			{
				marketId:   exchange_common.MARKET_BTC_USD,
				marketName: constants.BtcUsdPair,
			}: {
				exchangeToResponse: validResponses_AllExchanges_Median101,
			},
			{
				marketId:   exchange_common.MARKET_ETH_USD,
				marketName: constants.EthUsdPair,
			}: {
				exchangeToResponse: validResponses_AllExchanges_Median9001,
			},
		},
	}
)
