import { NodeEnv, logger, stats } from '@dydxprotocol-indexer/base';
import {
  CandleFromDatabase,
  CandleResolution,
  CandlesMap,
  CandleTable,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';

import config from '../config';

let candlesMap: CandlesMap = {};

export async function startCandleCache(txId?: number): Promise<void> {
  const startTime: number = Date.now();
  logger.info({
    at: 'candle-cache#startCandleCache',
    message: 'Starting candle cache',
  });

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {}, [], { txId },
  );
  logger.info({
    at: 'candle-cache#startCandleCache',
    message: 'Found all perpetual markets',
  });

  const tickers: string[] = _.map(
    perpetualMarkets,
    PerpetualMarketColumns.ticker,
  );
  candlesMap = await CandleTable.findCandlesMap(
    tickers,
    Object.values(CandleResolution),
    { txId },
  );
  logger.info({
    at: 'candle-cache#startCandleCache',
    message: 'Found candles map',
  });

  stats.timing(`${config.SERVICE_NAME}.start_candle_cache`, Date.now() - startTime);
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
