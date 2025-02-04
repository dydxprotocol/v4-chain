interface OraclePriceCache {
  [ticker: string]: string | undefined,
}

const oraclePriceCache: OraclePriceCache = {};

/**
 * Refresh loop to cache the list of all perpetual markets from the database in-memory.
 */
export async function start(): Promise<void> {
  // await loopHelpers.startUpdateLoop(
  //   updateOrderbookMidPrices,
  //   config.ORDERBOOK_MID_PRICE_REFRESH_INTERVAL_MS,
  //   'updateOrderbookMidPrices',
  // );
}

export function setOraclePrice(ticker: string, price: string) {
  oraclePriceCache[ticker] = price;
}

export function getOracleCachePrice(ticker: string): string | undefined {
  return oraclePriceCache[ticker];
}
