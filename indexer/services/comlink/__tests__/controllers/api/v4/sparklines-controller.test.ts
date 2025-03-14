import {
  CandleFromDatabase,
  CandleResolution,
  CandleTable,
  dbHelpers,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  testConstants,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import { DateTime } from 'luxon';
import request from 'supertest';
import { SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP } from '../../../../src/lib/constants';

import { RequestMethod, SparklineTimePeriod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';
import Big from 'big.js';
import * as SubaccountTable from '@dydxprotocol-indexer/postgres/build/src/stores/subaccount-table';
import {
  defaultLiquidityTier,
  defaultLiquidityTier2,
  defaultMarket,
  defaultMarket2,
  defaultMarket3,
  defaultPerpetualMarket,
  defaultPerpetualMarket2,
  defaultPerpetualMarket3,
  defaultSubaccount,
  defaultSubaccount2,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import * as MarketTable from '@dydxprotocol-indexer/postgres/build/src/stores/market-table';
import * as LiquidityTiersTable from '@dydxprotocol-indexer/postgres/build/src/stores/liquidity-tiers-table';
import * as PerpetualMarketTable from '@dydxprotocol-indexer/postgres/build/src/stores/perpetual-market-table';

// helper function to seed data
async function seedData() {
  await Promise.all([
    SubaccountTable.create(defaultSubaccount),
    SubaccountTable.create(defaultSubaccount2),
  ]);
  await Promise.all([
    MarketTable.create(defaultMarket),
    MarketTable.create(defaultMarket2),
    MarketTable.create(defaultMarket3),
  ]);
  await Promise.all([
    LiquidityTiersTable.create(defaultLiquidityTier),
    LiquidityTiersTable.create(defaultLiquidityTier2),
  ]);
  await Promise.all([
    PerpetualMarketTable.create(defaultPerpetualMarket),
    PerpetualMarketTable.create(defaultPerpetualMarket2),
    PerpetualMarketTable.create(defaultPerpetualMarket3),
  ]);
}

describe('sparklines-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('/v4/sparklines', () => {
    it('successfully returns no sparklines if no candles exist', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/sparklines?timePeriod=ONE_DAY',
      });

      expect(response.body).toEqual({
        [testConstants.defaultPerpetualMarket.ticker]: [],
        [testConstants.defaultPerpetualMarket2.ticker]: [],
        [testConstants.defaultPerpetualMarket3.ticker]: [],
      });
    });

    it.each(
      _.map(Object.values(SparklineTimePeriod), (timePeriod: SparklineTimePeriod) => [timePeriod]),
    )('successfully returns time period %s sparklines', async (timePeriod: SparklineTimePeriod) => {
      const tickerToBasePrice: Record<string, number> = {
        [testConstants.defaultPerpetualMarket.ticker]: 20000,
        [testConstants.defaultPerpetualMarket2.ticker]: 1000,
        [testConstants.defaultPerpetualMarket3.ticker]: 0.00000062,
      };
      const tickerToCandles: Record<string, Record<CandleResolution, string>> = _.mapValues(
        tickerToBasePrice,
        (basePrice: number): Record<CandleResolution, string> => {
          return {
            [CandleResolution.ONE_DAY]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.FOUR_HOURS]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.ONE_HOUR]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.THIRTY_MINUTES]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.FIFTEEN_MINUTES]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.FIVE_MINUTES]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
            [CandleResolution.ONE_MINUTE]:
              Big(Math.random().toFixed(2)).mul(basePrice).toFixed(),
          };
        },
      );

      await Promise.all(
        // eslint-disable-next-line @typescript-eslint/require-await
        _.flatten(
          _.map(
            [
              testConstants.defaultPerpetualMarket,
              testConstants.defaultPerpetualMarket2,
              testConstants.defaultPerpetualMarket3,
            ],
            (perpetualMarket: PerpetualMarketFromDatabase): Promise<CandleFromDatabase>[] => {
              return _.map(
                Object.values(CandleResolution),
                // eslint-disable-next-line @typescript-eslint/require-await
                async (res: CandleResolution): Promise<CandleFromDatabase> => {
                  return CandleTable.create({
                    ...testConstants.defaultCandle,
                    ticker: perpetualMarket.ticker,
                    resolution: res,
                    close: tickerToCandles[perpetualMarket.ticker][res],
                  });
                },
              );
            },
          ),
        ),
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/sparklines?timePeriod=${timePeriod}`,
      });

      const resolution: CandleResolution = SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP[timePeriod];
      expect(response.body).toEqual({
        [testConstants.defaultPerpetualMarket.ticker]: [
          tickerToCandles[testConstants.defaultPerpetualMarket.ticker][resolution],
        ],
        [testConstants.defaultPerpetualMarket2.ticker]: [
          tickerToCandles[testConstants.defaultPerpetualMarket2.ticker][resolution],
        ],
        [testConstants.defaultPerpetualMarket3.ticker]: [
          tickerToCandles[testConstants.defaultPerpetualMarket3.ticker][resolution],
        ],
      });
    });

    it('successfully returns a sparkline for a time period', async () => {
      const defaultTimePeriod: SparklineTimePeriod = SparklineTimePeriod.ONE_DAY;
      const resolution:
      CandleResolution = SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP[defaultTimePeriod];
      const now = DateTime.now().startOf('hour');  // Round to current hour
      const allClosePrices: string[] = [];

      // Create 100 hourly candles from oldest to newest, aligned to hour marks
      await Promise.all(
        // eslint-disable-next-line @typescript-eslint/require-await
        _.times(100, async (i: number) => {
          const close = Math.floor(Math.random() * 20000).toString();
          // Store prices oldest to newest
          allClosePrices.push(close);
          // Create candles from 99h ago to now, aligned to hour marks, in chronological order.
          const startedAt = now.minus({ hours: 99 - i }).toISO();
          return CandleTable.create({
            ...testConstants.defaultCandle,
            resolution,
            close,
            startedAt,
          });
        }),
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/sparklines?timePeriod=${defaultTimePeriod}`,
      });

      // Should only get back the most recent 24 candles
      // Reverse since we expect response to be in reverse chronological order.
      const last24Prices = allClosePrices.slice(-24).reverse();
      expect(response.body).toEqual({
        [testConstants.defaultPerpetualMarket.ticker]: last24Prices,
        [testConstants.defaultPerpetualMarket2.ticker]: [],
        [testConstants.defaultPerpetualMarket3.ticker]: [],
      });
    });

    it('successfully returns multiple sparklines when one sparkline has less than enough candles',
      async () => {
        const timePeriod: SparklineTimePeriod = SparklineTimePeriod.ONE_DAY;
        const resolution: CandleResolution = SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP[timePeriod];
        const firstClosing: string = Math.floor(Math.random() * 20000).toString();

        const numCandles: number = 24; // enough for ONE_HOUR for a day
        await Promise.all(
          _.times(numCandles, (i: number) => {
            return CandleTable.create({
              ...testConstants.defaultCandle,
              startedAt: DateTime
                .fromISO(testConstants.defaultCandle.startedAt)
                .minus({ hour: i })
                .toISO(),
              ticker: testConstants.defaultPerpetualMarket.ticker,
              resolution,
              close: firstClosing,
            });
          }),
        );

        const secondClosing: string = Math.floor(Math.random() * 20000).toString();

        const numCandles2: number = numCandles - 10; // not enough for ONE_HOUR for a day
        await Promise.all(
          _.times(numCandles2, (i: number) => {
            return CandleTable.create({
              ...testConstants.defaultCandle,
              startedAt: DateTime
                .fromISO(testConstants.defaultCandle.startedAt)
                .minus({ hour: i })
                .toISO(),
              ticker: testConstants.defaultPerpetualMarket2.ticker,
              resolution,
              close: secondClosing,
            });
          }),
        );

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/sparklines?timePeriod=${timePeriod}`,
        });

        expect(response.body).toEqual({
          [testConstants.defaultPerpetualMarket.ticker]: _.times(numCandles, () => firstClosing),
          [testConstants.defaultPerpetualMarket2.ticker]: _.times(numCandles2, () => secondClosing),
          [testConstants.defaultPerpetualMarket3.ticker]: [],
        });
      },
    );
  });
});
