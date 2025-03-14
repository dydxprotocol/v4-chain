import {
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  FillTable,
  Market24HourTradeVolumes,
  MarketOpenInterest,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketsMap,
  PerpetualMarketTable,
  PerpetualPositionTable,
  Transaction,
  OraclePriceTable,
  PriceMap,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable, LiquidityTiersMap, LiquidityTiersColumns,
} from '@dydxprotocol-indexer/postgres';
import { NextFundingCache } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';
import _ from 'lodash';

import config from '../config';
import { redisClient } from '../helpers/redis';
import { compareAndHandleMarketsWebsocketMessage } from '../helpers/websocket';

export function getPriceChange(
  marketId: number,
  latestPrices: PriceMap,
  prices24hAgo: PriceMap,
): string | undefined {
  const latestPrice: string | undefined = latestPrices[marketId];
  const price24hAgo: string | undefined = prices24hAgo[marketId];
  if (latestPrice === undefined || price24hAgo === undefined) {
    logger.info({
      at: 'market-updater#getPriceChange',
      message: 'No 24-hour old oracle price found for market',
      marketId,
    });
    return undefined;
  }
  return new Big(latestPrice)
    .minus(price24hAgo)
    .toFixed();
}

export default async function runTask(): Promise<void> {
  const start: number = Date.now();

  // Run all initial database queries in parallel
  const [
    liquidityTiers,
    perpetualMarkets,
    latestPrices,
    prices24hAgo,
  ] : [
    LiquidityTiersFromDatabase[],
    PerpetualMarketFromDatabase[],
    PriceMap,
    PriceMap,
  ] = await Promise.all([
    LiquidityTiersTable.findAll({}, []),
    PerpetualMarketTable.findAll({}, []),
    OraclePriceTable.getLatestPrices(),
    OraclePriceTable.getPricesFrom24hAgo(),
  ]);

  // Derive data from perpetual markets
  const perpetualMarketIds: string[] = _.map(perpetualMarkets, PerpetualMarketColumns.id);
  const clobPairIds: string[] = _.map(perpetualMarkets, PerpetualMarketColumns.clobPairId);
  const tickerDefaultFundingRate1HPairs: [string, string][] = _.map(
    perpetualMarkets,
    (market) => [
      market[PerpetualMarketColumns.ticker],
      // Use 0 as default for null default funding rate
      market[PerpetualMarketColumns.defaultFundingRate1H] ?? '0',
    ],
  );

  stats.timing(
    `${config.SERVICE_NAME}.market_updater_initial_queries`,
    Date.now() - start,
  );

  try {
    const [
      perpetualMarketStats,
      openInterest,
      fundingRates,
    ]: [
      _.Dictionary<Market24HourTradeVolumes>,
      _.Dictionary<MarketOpenInterest>,
      _.Dictionary<Big | undefined>,
    ] = await Promise.all([
      // TODO(DEC-1149 Add support for pulling information from candles
      FillTable.get24HourInformation(clobPairIds),
      PerpetualPositionTable.getOpenInterestLong(perpetualMarketIds),
      NextFundingCache.getNextFunding(redisClient, tickerDefaultFundingRate1HPairs),
    ]);

    stats.timing(
      `${config.SERVICE_NAME}.markets_updater_get_fills_positions_and_markets_timing`,
      Date.now() - start,
    );

    const clobPairIdToPerpetualMarketId: _.Dictionary<string> = _.chain(perpetualMarkets)
      .keyBy(PerpetualMarketColumns.clobPairId)
      .mapValues(PerpetualMarketColumns.id)
      .value();
    const perpetualMarketStatsByPerpetualMarketId: _.Dictionary<Market24HourTradeVolumes> = _.keyBy(
      perpetualMarketStats,
      (stat: Market24HourTradeVolumes) => {
        return clobPairIdToPerpetualMarketId[stat.clobPairId];
      },
    );

    const perpetualMarketMap: PerpetualMarketsMap = _.keyBy(
      perpetualMarkets,
      PerpetualMarketColumns.id,
    );
    const liquidityTiersMap: LiquidityTiersMap = _.keyBy(
      liquidityTiers,
      LiquidityTiersColumns.id,
    );
    const updatedMarketsMap: PerpetualMarketsMap = {};

    const txId: number = await Transaction.start();
    try {
      await PerpetualMarketTable.updateMarketCheckerFields(
        perpetualMarkets.map((perpetualMarket: PerpetualMarketFromDatabase) => {
          const perpetualMarketId: string = perpetualMarket.id;
          const marketId: number = perpetualMarket.marketId;
          const marketUpdates: {
            trades24H: number,
            volume24H: string,
            priceChange24H: string,
            openInterest: string,
            nextFundingRate: string,
          } = {
            trades24H: parseInt(
              perpetualMarketStatsByPerpetualMarketId[perpetualMarketId].trades24H,
              10,
            ) ?? 0,
            volume24H: perpetualMarketStatsByPerpetualMarketId[perpetualMarketId].volume24H ?? '0',
            priceChange24H: getPriceChange(marketId, latestPrices, prices24hAgo) ??
              perpetualMarket.priceChange24H,
            openInterest: openInterest[perpetualMarketId]?.openInterest ?? '0',
            nextFundingRate: fundingRates[perpetualMarket.ticker]?.toFixed() ??
              perpetualMarket.nextFundingRate,
          };

          // Keep track of updated markets to create the websocket message
          updatedMarketsMap[perpetualMarketId] = {
            ...perpetualMarketMap[perpetualMarketId],
            ...marketUpdates,
          };

          // Return field updates for bulk update
          return {
            id: perpetualMarketId,
            ...marketUpdates,
          };
        }),
        Transaction.get(txId),
      );
      await Transaction.commit(txId);
    } catch (error) {
      await Transaction.rollback(txId);
      throw error;
    }

    compareAndHandleMarketsWebsocketMessage({
      oldMarkets: perpetualMarketMap,
      newMarkets: updatedMarketsMap,
      liquidityTiers: liquidityTiersMap,
    });
  } catch (error) {
    logger.error({
      at: 'market-checker#runTask',
      message: 'Error occurred in task to update markets',
      error,
    });
  } finally {
    stats.timing(
      `${config.SERVICE_NAME}.market_updater_update_timing`,
      Date.now() - start,
    );
  }
}
