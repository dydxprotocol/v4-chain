"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MARKET_BTC_USD = void 0;
const constants_1 = require("../../../src/clients/constants");
const indexer_client_1 = require("../../../src/clients/indexer-client");
// ------------ Markets ------------
exports.MARKET_BTC_USD = 'BTC-USD';
describe('IndexerClient', () => {
    const client = new indexer_client_1.IndexerClient(constants_1.Network.testnet().indexerConfig);
    describe('Market Endpoints', () => {
        it('Markets', async () => {
            const response = await client.markets.getPerpetualMarkets();
            const btc = response.markets[exports.MARKET_BTC_USD];
            const status = btc.status;
            expect(status).toBe('ACTIVE');
        });
        it('BTC Market', async () => {
            const response = await client.markets.getPerpetualMarkets(exports.MARKET_BTC_USD);
            const btc = response.markets[exports.MARKET_BTC_USD];
            const status = btc.status;
            expect(status).toBe('ACTIVE');
        });
        it('BTC Trades', async () => {
            const response = await client.markets.getPerpetualMarketTrades(exports.MARKET_BTC_USD);
            const trades = response.trades;
            expect(trades).not.toBeUndefined();
        });
        it('BTC Orderbook', async () => {
            const response = await client.markets.getPerpetualMarketOrderbook(exports.MARKET_BTC_USD);
            const asks = response.asks;
            const bids = response.bids;
            expect(asks).not.toBeUndefined();
            expect(bids).not.toBeUndefined();
        });
        it('BTC Candles', async () => {
            const response = await client.markets.getPerpetualMarketCandles(exports.MARKET_BTC_USD, '1MIN');
            const candles = response.candles;
            expect(candles).not.toBeUndefined();
        });
        it('BTC Historical Funding', async () => {
            const response = await client.markets.getPerpetualMarketHistoricalFunding(exports.MARKET_BTC_USD);
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
            const btcSparklines = response[exports.MARKET_BTC_USD];
            expect(btcSparklines).not.toBeUndefined();
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiTWFya2V0c0VuZHBvaW50cy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvY2xpZW50L01hcmtldHNFbmRwb2ludHMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4REFBeUQ7QUFDekQsd0VBQW9FO0FBRXBFLG9DQUFvQztBQUN2QixRQUFBLGNBQWMsR0FBVyxTQUFTLENBQUM7QUFFaEQsUUFBUSxDQUFDLGVBQWUsRUFBRSxHQUFHLEVBQUU7SUFDN0IsTUFBTSxNQUFNLEdBQUcsSUFBSSw4QkFBYSxDQUFDLG1CQUFPLENBQUMsT0FBTyxFQUFFLENBQUMsYUFBYSxDQUFDLENBQUM7SUFFbEUsUUFBUSxDQUFDLGtCQUFrQixFQUFFLEdBQUcsRUFBRTtRQUNoQyxFQUFFLENBQUMsU0FBUyxFQUFFLEtBQUssSUFBSSxFQUFFO1lBQ3ZCLE1BQU0sUUFBUSxHQUFHLE1BQU0sTUFBTSxDQUFDLE9BQU8sQ0FBQyxtQkFBbUIsRUFBRSxDQUFDO1lBQzVELE1BQU0sR0FBRyxHQUFHLFFBQVEsQ0FBQyxPQUFPLENBQUMsc0JBQWMsQ0FBQyxDQUFDO1lBQzdDLE1BQU0sTUFBTSxHQUFHLEdBQUcsQ0FBQyxNQUFNLENBQUM7WUFDMUIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQztRQUNoQyxDQUFDLENBQUMsQ0FBQztRQUVILEVBQUUsQ0FBQyxZQUFZLEVBQUUsS0FBSyxJQUFJLEVBQUU7WUFDMUIsTUFBTSxRQUFRLEdBQUcsTUFBTSxNQUFNLENBQUMsT0FBTyxDQUFDLG1CQUFtQixDQUFDLHNCQUFjLENBQUMsQ0FBQztZQUMxRSxNQUFNLEdBQUcsR0FBRyxRQUFRLENBQUMsT0FBTyxDQUFDLHNCQUFjLENBQUMsQ0FBQztZQUM3QyxNQUFNLE1BQU0sR0FBRyxHQUFHLENBQUMsTUFBTSxDQUFDO1lBQzFCLE1BQU0sQ0FBQyxNQUFNLENB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