import { logger } from '@dydxprotocol-indexer/base';
import {
  Transaction,
  OraclePriceFromDatabase,
  OraclePriceTable,
  PriceMap,
  MarketTable,
  MarketFromDatabase,
} from '@dydxprotocol-indexer/postgres';

let priceMap: PriceMap = {};

// PriceCache should be started after BlockCache is started.
export async function startPriceCache(blockHeight: string, txId?: number): Promise<void> {
  const markets: MarketFromDatabase[] = await MarketTable.findAll({}, [], { txId });
  priceMap = await OraclePriceTable.findLatestPrices(
    blockHeight,
    Transaction.get(txId),
  );
  markets.forEach((market: MarketFromDatabase) => {
    if (
      priceMap[market.id] === undefined &&
      market.oraclePrice !== null &&
      market.oraclePrice !== undefined
    ) {
      priceMap[market.id] = market.oraclePrice;
    }
  });
}

export function getPrice(
  marketId: number,
): string {
  if (marketId in priceMap) {
    return priceMap[marketId];
  } else {
    logger.crit({
      at: 'PriceCache.getPrice',
      message: 'price not found',
      marketId,
    });
    throw Error(`price not found for marketId ${marketId} in price cache`);
  }
}

export function updatePriceCacheWithPrice(oraclePrice: OraclePriceFromDatabase): void {
  priceMap[oraclePrice.marketId] = oraclePrice.price;
}

export function getPriceMap(): PriceMap {
  return priceMap;
}

export function clearPriceMap(): void {
  if (process.env.NODE_ENV !== 'test') {
    throw Error('cannot clear price map outside of test environment');
  }

  priceMap = {};
}
