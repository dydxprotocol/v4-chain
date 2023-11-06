import { stats } from '@dydxprotocol-indexer/base';
import {
  CandleColumns,
  CandleFromDatabase,
  CandleResolution,
  CandleTable,
  DEFAULT_POSTGRES_OPTIONS,
  Ordering,
  PerpetualMarketColumns,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { SPARKLINE_TIME_PERIOD_TO_LIMIT_MAP, SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP } from '../../../lib/constants';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import { candlesToSparklineResponseObject } from '../../../request-helpers/request-transformer';
import { SparklineResponseObject, SparklinesRequest, SparklineTimePeriod } from '../../../types';

const router = express.Router();
const controllerName: string = 'sparklines-controller';

@Route('sparklines')
class FillsController extends Controller {
  @Get('/')
  async get(
    @Query() timePeriod: SparklineTimePeriod,
  ): Promise<SparklineResponseObject> {
    const tickers: string[] = _.map(
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      PerpetualMarketColumns.ticker,
    );

    const resolution: CandleResolution = SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP[timePeriod];
    const limit: number = SPARKLINE_TIME_PERIOD_TO_LIMIT_MAP[timePeriod];

    const ungroupedTickerCandles: CandleFromDatabase[] = await CandleTable.findAll(
      {
        ticker: tickers,
        resolution,
        limit: limit * tickers.length,
      },
      [],
      { ...DEFAULT_POSTGRES_OPTIONS, orderBy: [[CandleColumns.startedAt, Ordering.DESC]] },
    );

    return candlesToSparklineResponseObject(tickers, ungroupedTickerCandles, limit);
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...checkSchema({
    timePeriod: {
      in: 'query',
      isString: true,
      isIn: {
        options: [Object.values(SparklineTimePeriod)],
      },
      errorMessage: `timePeriod must be a valid Time Period, one of ${Object.values(SparklineTimePeriod)}`,
    },
  }),
  handleValidationErrors,
  async (req: express.Request, res: express.Response) => {
    const {
      timePeriod,
    } = matchedData(req) as SparklinesRequest;

    const start: number = Date.now();
    try {
      const controller: FillsController = new FillsController();
      const sparklineResponse: SparklineResponseObject = await controller.get(
        timePeriod,
      );

      return res.send(sparklineResponse);
    } catch (error) {
      return handleControllerError(
        'SparklinesController GET /',
        'Sparklines error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_sparklines.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
