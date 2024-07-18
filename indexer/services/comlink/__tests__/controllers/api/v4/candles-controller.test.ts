import {
  CandleResolution,
  CandleTable,
  dbHelpers,
  helpers,
  IsoString,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import _, { max, min } from 'lodash';
import request from 'supertest';

import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';
import config from '../../../../src/config';

describe('candles-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('/v4/candles/perpetualMarkets/:ticker', () => {
    it('successfully returns no candles if none exist', async () => {
      const resolution: CandleResolution = CandleResolution.ONE_DAY;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${resolution}`,
      });

      expect(response.body.candles).toEqual([]);
    });

    it.each(
      _.map(Object.values(CandleResolution), (resolution: CandleResolution) => [resolution]),
    )('successfully returns resolution %s candles', async (resolution: CandleResolution) => {
      await Promise.all(
        // eslint-disable-next-line @typescript-eslint/require-await
        _.map(Object.values(CandleResolution), async (res: CandleResolution) => {
          return CandleTable.create({
            ...testConstants.defaultCandle,
            resolution: res,
          });
        }),
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${resolution}`,
      });

      expect(response.body.candles).toEqual([{
        ...testConstants.defaultCandle,
        resolution,
      }]);
    });

    it('successfully returns at most API_LIMIT_V4 candles', async () => {
      await Promise.all(
        // eslint-disable-next-line @typescript-eslint/require-await
        _.times(config.API_LIMIT_V4 + 1, async (count: number) => {
          return CandleTable.create({
            ...testConstants.defaultCandle,
            startedAt: testConstants.createdDateTime.minus({ minutes: count }).toISO(),
            resolution: CandleResolution.ONE_DAY,
          });
        }),
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_DAY}`,
      });

      expect(response.body.candles.length).toEqual(config.API_LIMIT_V4);
    });

    it('accepts includeOrderbook as a parameter', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=true`,
      });
      expect(response.statusCode).toEqual(200);
    });
  });

  describe('getCandles', () => {
    it('returns unaltered candles when includeOrderbook is false', async () => {
      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=false`,
      });

      expect(response.body.candles).toEqual([{
        ...testConstants.defaultCandle,
      }]);
    });

    it('returns candles using orderbookMidPriceOpen and orderbookMidPriceClose when includeOrderbook is true', async () => {
      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        trades: 0,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=true`,
      });

      const low = min([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.orderbookMidPriceOpen]);
      const high = max([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.orderbookMidPriceOpen]);

      expect(response.body.candles).toEqual([{
        startedAt: testConstants.defaultCandle.startedAt,
        ticker: testConstants.defaultCandle.ticker,
        resolution: testConstants.defaultCandle.resolution,
        low,
        high,
        open: testConstants.defaultCandle.orderbookMidPriceOpen,
        close: testConstants.defaultCandle.orderbookMidPriceClose,
        baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
        usdVolume: testConstants.defaultCandle.usdVolume,
        trades: 0,
        startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
        orderbookMidPriceOpen: testConstants.defaultCandle.orderbookMidPriceOpen,
        orderbookMidPriceClose: testConstants.defaultCandle.orderbookMidPriceClose,
      }]);
    });

    it('when orderbookMidPriceClose is null, returns the original candle close values', async () => {
      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        orderbookMidPriceClose: undefined,
        trades: 0,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=true`,
      });

      const low = min([testConstants.defaultCandle.close,
        testConstants.defaultCandle.orderbookMidPriceOpen]);
      const high = max([testConstants.defaultCandle.close,
        testConstants.defaultCandle.orderbookMidPriceOpen]);

      expect(response.body.candles).toEqual([{
        startedAt: testConstants.defaultCandle.startedAt,
        ticker: testConstants.defaultCandle.ticker,
        resolution: testConstants.defaultCandle.resolution,
        low,
        high,
        open: testConstants.defaultCandle.orderbookMidPriceOpen,
        close: testConstants.defaultCandle.close,
        baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
        usdVolume: testConstants.defaultCandle.usdVolume,
        trades: 0,
        startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
        orderbookMidPriceOpen: testConstants.defaultCandle.orderbookMidPriceOpen,
        orderbookMidPriceClose: null,
      }]);

    });

    it('when orderbookMidPriceOpen is null, returns the original candle open values', async () => {
      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        orderbookMidPriceOpen: undefined,
        trades: 0,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=true`,
      });

      const low = min([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.open]);
      const high = max([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.open]);

      expect(response.body.candles).toEqual([{
        startedAt: testConstants.defaultCandle.startedAt,
        ticker: testConstants.defaultCandle.ticker,
        resolution: testConstants.defaultCandle.resolution,
        low,
        high,
        open: testConstants.defaultCandle.open,
        close: testConstants.defaultCandle.orderbookMidPriceClose,
        baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
        usdVolume: testConstants.defaultCandle.usdVolume,
        trades: 0,
        startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
        orderbookMidPriceOpen: null,
        orderbookMidPriceClose: testConstants.defaultCandle.orderbookMidPriceClose,
      }]);
    });

    it('correctly formats multiple candles when includeOrderbook is true', async () => {
      const startedAtMinusOne: IsoString = helpers.calculateNormalizedCandleStartTime(
        testConstants.createdDateTime.minus({ minutes: 1 }),
        CandleResolution.ONE_MINUTE,
      ).toISO();

      const startedAtMinusTwo: IsoString = helpers.calculateNormalizedCandleStartTime(
        testConstants.createdDateTime.minus({ minutes: 2 }),
        CandleResolution.ONE_MINUTE,
      ).toISO();

      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        trades: 0,
      });

      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        startedAt: startedAtMinusOne,
        trades: 0,
      });

      await CandleTable.create({
        ...testConstants.defaultCandle,
        resolution: CandleResolution.ONE_MINUTE,
        startedAt: startedAtMinusTwo,
        trades: 0,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/candles/perpetualMarkets/${testConstants.defaultCandle.ticker}` +
          `?resolution=${CandleResolution.ONE_MINUTE}&includeOrderbook=true`,
      });

      const low = min([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.orderbookMidPriceOpen]);
      const high = max([testConstants.defaultCandle.orderbookMidPriceClose,
        testConstants.defaultCandle.orderbookMidPriceOpen]);

      expect(response.body.candles).toEqual([
        {
          startedAt: testConstants.defaultCandle.startedAt,
          ticker: testConstants.defaultCandle.ticker,
          resolution: testConstants.defaultCandle.resolution,
          low,
          high,
          open: testConstants.defaultCandle.orderbookMidPriceOpen,
          close: testConstants.defaultCandle.orderbookMidPriceClose,
          baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
          usdVolume: testConstants.defaultCandle.usdVolume,
          trades: 0,
          startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
          orderbookMidPriceOpen: testConstants.defaultCandle.orderbookMidPriceOpen,
          orderbookMidPriceClose: testConstants.defaultCandle.orderbookMidPriceClose,
        },
        {
          startedAt: startedAtMinusOne,
          ticker: testConstants.defaultCandle.ticker,
          resolution: testConstants.defaultCandle.resolution,
          low,
          high,
          open: testConstants.defaultCandle.orderbookMidPriceOpen,
          close: testConstants.defaultCandle.orderbookMidPriceClose,
          baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
          usdVolume: testConstants.defaultCandle.usdVolume,
          trades: 0,
          startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
          orderbookMidPriceOpen: testConstants.defaultCandle.orderbookMidPriceOpen,
          orderbookMidPriceClose: testConstants.defaultCandle.orderbookMidPriceClose,
        },
        {
          startedAt: startedAtMinusTwo,
          ticker: testConstants.defaultCandle.ticker,
          resolution: testConstants.defaultCandle.resolution,
          low,
          high,
          open: testConstants.defaultCandle.orderbookMidPriceOpen,
          close: testConstants.defaultCandle.orderbookMidPriceClose,
          baseTokenVolume: testConstants.defaultCandle.baseTokenVolume,
          usdVolume: testConstants.defaultCandle.usdVolume,
          trades: 0,
          startingOpenInterest: testConstants.defaultCandle.startingOpenInterest,
          orderbookMidPriceOpen: testConstants.defaultCandle.orderbookMidPriceOpen,
          orderbookMidPriceClose: testConstants.defaultCandle.orderbookMidPriceClose,
        },
      ]);
    });
  });
});
