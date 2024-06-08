import { Data } from '../types';
import RestClient from './rest';
/**
 * @description REST endpoints for data unrelated to a particular address.
 */
export default class MarketsClient extends RestClient {
    getPerpetualMarkets(market?: string): Promise<Data>;
    getPerpetualMarketOrderbook(market: string): Promise<Data>;
    getPerpetualMarketTrades(market: string, startingBeforeOrAtHeight?: number | null, limit?: number | null): Promise<Data>;
    getPerpetualMarketCandles(market: string, resolution: string, fromISO?: string | null, toISO?: string | null, limit?: number | null): Promise<Data>;
    getPerpetualMarketHistoricalFunding(market: string, effectiveBeforeOrAt?: string | null, effectiveBeforeOrAtHeight?: number | null, limit?: number | null): Promise<Data>;
    getPerpetualMarketSparklines(period?: string): Promise<Data>;
}
