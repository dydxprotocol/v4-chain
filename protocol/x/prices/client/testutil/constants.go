package testutil

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	pricefeed "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
)

const (
	responseErr = 503
	responseOk  = 200
)

var (
	binanceTicker_Btc100    = NewBinanceTicker("BTCUSDT", "99", "100", "101")
	binanceTicker_Eth9001   = NewBinanceTicker("ETHUSDT", "9000", "9001", "9002")
	binanceUSTicker_Btc101  = NewBinanceTicker("BTCUSD", "100", "101", "102")
	binanceUSTicker_Eth9000 = NewBinanceTicker("ETHUSD", "8999", "9000", "9001")

	bitfinexTicker_Btc102  = NewBitfinexTicker("tBTCUSD", 101.0, 102.0, 103.0)
	bitfinexTicker_Eth9002 = NewBitfinexTicker("tETHUSD", 9001.0, 9002.0, 9003.0)

	// Test Bitfinex Config
	BitfinexExchangeConfig = map[pricefeed.MarketId]pricefeed.MarketConfig{
		exchange_config.MARKET_BTC_USD: {
			Ticker: "tBTCUSD",
		},
		exchange_config.MARKET_ETH_USD: {
			Ticker: "tETHUSD",
		},
		exchange_config.MARKET_USDT_USD: {
			Ticker: "tUSTUSD",
		},
	}

	// Exchange responses.
	EmptyResponses_AllExchanges = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: {
			ResponseCode: responseErr,
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: {
			ResponseCode: responseErr,
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: {
			ResponseCode: responseErr,
		},
	}
	FullResponses_AllExchanges_Btc101_Eth9001 = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: {
			ResponseCode: responseOk,
			Tickers:      []JsonResponse{binanceTicker_Btc100, binanceTicker_Eth9001},
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: {
			ResponseCode: responseOk,
			Tickers:      []JsonResponse{binanceUSTicker_Btc101, binanceUSTicker_Eth9000},
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: {
			ResponseCode: responseOk,
			Tickers:      []JsonResponse{bitfinexTicker_Btc102, bitfinexTicker_Eth9002},
		},
	}
	PartialResponses_AllExchanges_Eth9001 = map[ExchangeIdAndName]Response{
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE,
			exchangeName: constants.BinanceExchangeName,
		}: {
			ResponseCode: responseOk,
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BINANCE_US,
			exchangeName: constants.BinanceUSExchangeName,
		}: {
			ResponseCode: responseOk,
			Tickers:      []JsonResponse{binanceUSTicker_Eth9000},
		},
		{
			exchangeId:   exchange_common.EXCHANGE_ID_BITFINEX,
			exchangeName: constants.BitfinexExchangeName,
		}: {
			ResponseCode: responseOk,
			Tickers:      []JsonResponse{bitfinexTicker_Eth9002},
		},
	}

	// ValidAuthority is an authority address that passes basic validation.
	ValidAuthority = authtypes.NewModuleAddress("test").String()
)
