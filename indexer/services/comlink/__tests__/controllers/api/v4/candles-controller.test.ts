import {
  CandleResolution, CandleTable, dbHelpers, perpetualMarketRefresher, testConstants, testMocks,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
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
  });
});
