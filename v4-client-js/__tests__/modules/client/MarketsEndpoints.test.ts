import { Network } from '../../../src/clients/constants';
import { IndexerClient } from '../../../src/clients/indexer-client';

// ------------ Markets ------------
export const MARKET_BTC_USD: string = 'BTC-USD';

describe('IndexerClient', () => {
  const client = new IndexerClient(Network.testnet().indexerConfig);

  describe('Market Endpoints', () => {
    it('Markets', async () => {
      const response = await client.markets.getPerpetualMarkets();
      const btc = response.markets[MARKET_BTC_USD];
      const status = btc.status;
      expect(status).toBe('ACTIVE');
    });

    it('BTC Market', async () => {
      const response = await client.markets.getPerpetualMarkets(MARKET_BTC_USD);
      const btc = response.markets[MARKET_BTC_USD];
      const status = btc.status;
      expect(status).toBe('ACTIVE');
    });

    it('BTC Trades', async () => {
      const response = await client.markets.getPerpetualMarketTrades(MARKET_BTC_USD);
      const trades = response.trades;
      expect(trades).not.toBeUndefined();
    });

    it('BTC Orderbook', async () => {
      const response = await client.markets.getPerpetualMarketOrderbook(MARKET_BTC_USD);
      const asks = response.asks;
      const bids = response.bids;
      expect(asks).not.toBeUndefined();
      expect(bids).not.toBeUndefined();
    });

    it('BTC Candles', async () => {
      const response = await client.markets.getPerpetualMarketCandles(MARKET_BTC_USD, '1MIN');
      const candles = response.candles;
      expect(candles).not.toBeUndefined();
    });

    it('BTC Historical Funding', async () => {
      const response = await client.markets.getPerpetualMarketHistoricalFunding(MARKET_BTC_USD);
      expect(response).not.toBeNull();
      const historicalFunding = response.historicalFunding;
      expect(historicalFunding).not.toBeNull();
      if (historicalFunding.length > 0) {
        const historicalFunding0 = historicalFunding[0];
        expect(historicalFunding0).not.toBeNull();
      }
    });

    it('Sparklines', async () => {
      const response = await client.markets.getPerpetualMarketSparklines();
      const btcSparklines = response[MARKET_BTC_USD];
      expect(btcSparklines).not.toBeUndefined();
    });
  });
});
