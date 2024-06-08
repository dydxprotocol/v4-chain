/**
 * Simple JS example demostrating accessing market data with Indexer REST endpoints
 */

import { Network } from '../src/clients/constants';
import { IndexerClient } from '../src/clients/indexer-client';

// ------------ Markets ------------
export const MARKET_BTC_USD: string = 'BTC-USD';

async function test(): Promise<void> {
  const client = new IndexerClient(Network.testnet().indexerConfig);

  // Get perp markets
  try {
    const response = await client.markets.getPerpetualMarkets();
    console.log(response);
    console.log('markets');
    const btcMarket = response.markets['BTC-USD'];
    const btcMarketStatus = btcMarket.status;
    console.log(btcMarketStatus);
  } catch (error) {
    console.log(error.message);
  }

  try {
    const response = await client.markets.getPerpetualMarkets(MARKET_BTC_USD);
    console.log(response);
    console.log('markets');
    const btcMarket = response.markets['BTC-USD'];
    const btcMarketStatus = btcMarket.status;
    console.log(btcMarketStatus);
  } catch (error) {
    console.log(error.message);
  }

  // Get sparklines
  try {
    const response = await client.markets.getPerpetualMarketSparklines();
    console.log(response);
    console.log('sparklines');
    const btcSparklines = response['BTC-USD'];
    console.log(btcSparklines);
  } catch (error) {
    console.log(error.message);
  }

  // Get perp market trades
  try {
    const response = await client.markets.getPerpetualMarketTrades(MARKET_BTC_USD);
    console.log(response);
    console.log('trades');
    const trades = response.trades;
    console.log(trades);
  } catch (error) {
    console.log(error.message);
  }

  // Get perp market orderbook
  try {
    const response = await client.markets.getPerpetualMarketOrderbook(MARKET_BTC_USD);
    console.log(response);
    console.log('orderbook');
    const asks = response.asks;
    const bids = response.bids;
    if (asks.length > 0) {
      const asks0 = asks[0];
      const asks0Price = asks0.price;
      const asks0Size = asks0.size;
      console.log(asks0Price);
      console.log(asks0Size);
    }
    if (bids.length > 0) {
      const bids0 = bids[0];
      const bids0Price = bids0.price;
      const bids0Size = bids0.size;
      console.log(bids0Price);
      console.log(bids0Size);
    }
    const trades = response.trades;
    console.log(trades);
  } catch (error) {
    console.log(error.message);
  }

  // Get perp market candles
  try {
    const response = await client.markets.getPerpetualMarketCandles(MARKET_BTC_USD, '1MIN');
    console.log(response);
    console.log('candles');
    const candles = response.candles;
    if (candles.length > 0) {
      const candles0 = candles[0];
      const startedAt = candles0.startedAt;
      const low = candles0.low;
      const high = candles0.high;
      const open = candles0.open;
      const close = candles0.close;
      const baseTokenVolume = candles0.baseTokenVolume;
      const usdVolume = candles0.usdVolume;
      const trades = candles0.trades;
      console.log(startedAt);
      console.log(low);
      console.log(high);
      console.log(open);
      console.log(close);
      console.log(baseTokenVolume);
      console.log(usdVolume);
      console.log(trades);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get perp market historical funding rates
  try {
    const response = await client.markets.getPerpetualMarketHistoricalFunding(MARKET_BTC_USD);
    console.log(response);
    console.log('historical funding');
    const historicalFunding = response.historicalFunding;
    if (historicalFunding.length > 0) {
      const historicalFunding0 = historicalFunding[0];
      console.log(historicalFunding0);
    }
  } catch (error) {
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
