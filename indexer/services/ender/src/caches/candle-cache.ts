import { logger, NodeEnv } from '@dydxprotocol-indexer/base';
import {
  BlockTable,
  CandleFromDatabase,
  CandleResolution,
  CandlesMap,
  CandleTable,
  IsoString,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';

let candlesMap: CandlesMap = {};

export async function startCandleCache(txId?: number): Promise<void> {
  let latestBlockTime: IsoString;
  try {
    const latestBlock = await BlockTable.getLatest({ txId });
    latestBlockTime = latestBlock.time;
  } catch (error) {
    logger.error({
      at: 'ender#startCandleCache',
      message: 'Cannot fetch latest indexed block and falling back to current timestamp',
    });
    latestBlockTime = new Date().toISOString();
  }

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {}, [], { txId },
  );
  const tickers: string[] = _.map(
    perpetualMarkets,
    PerpetualMarketColumns.ticker,
  );

  candlesMap = await CandleTable.findCandlesMap(
    tickers,
    latestBlockTime,
  );
}

export function getCandle(
  ticker: string,
  resolution: CandleResolution,
): CandleFromDatabase | undefined {
  if (ticker in candlesMap && resolution in candlesMap[ticker]) {
    return candlesMap[ticker][resolution];
  }

  return undefined;
}

export function updateCandleCacheWithCandle(candle: CandleFromDatabase): void {
  if (!(candle.ticker in candlesMap)) {
    candlesMap[candle.ticker] = {};
  }

  candlesMap[candle.ticker][candle.resolution] = candle;
}

export function getCandlesMap(): CandlesMap {
  return candlesMap;
}

export function clearCandlesMap(): void {
  if (process.env.NODE_ENV !== NodeEnv.TEST) {
    throw Error('cannot clear candles map outside of test environment');
  }

  candlesMap = {};
}
