import { stats, cacheControlMiddleware } from '@dydxprotocol-indexer/base';
import {
  CandleFromDatabase,
  CandleResolution,
  CandleTable,
  DEFAULT_POSTGRES_OPTIONS,
  PerpetualMarketColumns,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { sparklinesRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP, SPARKLINE_TIME_PERIOD_TO_LOOKBACK_MAP } from '../../../lib/constants';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import { candlesToSparklineResponseObject } from '../../../request-helpers/request-transformer';
import { SparklineResponseObject, SparklinesRequest, SparklineTimePeriod } from '../../../types';

const router = express.Router();
const controllerName: string = 'sparklines-controller';
const sparklinesCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_SPARKLINES,
);

@Route('sparklines')
class SparklinesController extends Controller {
  @Get('/')
  async get(
    @Query() timePeriod: SparklineTimePeriod,
  ): Promise<SparklineResponseObject> {

    const tickers: string[] = _.map(
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      PerpetualMarketColumns.ticker,
    );

    const resolution: CandleResolution = SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP[timePeriod];
    const lookbackMs: number = SPARKLINE_TIME_PERIOD_TO_LOOKBACK_MAP[timePeriod];

    const ungroupedTickerCandles
    : CandleFromDatabase[] = await CandleTable.findByResAndLookbackPeriod(
      resolution,
      lookbackMs,
      DEFAULT_POSTGRES_OPTIONS,
    );

    return candlesToSparklineResponseObject(tickers, ungroupedTickerCandles);
  }
}

router.get(
  '/',
  rateLimiterMiddleware(sparklinesRateLimiter),
  sparklinesCacheControlMiddleware,
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
      const controller: SparklinesController = new SparklinesController();
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
