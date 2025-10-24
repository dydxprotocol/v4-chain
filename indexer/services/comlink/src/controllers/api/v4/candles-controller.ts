import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  CandleFromDatabase, CandleResolution, CandleTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { candlesRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckLimitSchema, CheckTickerParamSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import { candleToResponseObject } from '../../../request-helpers/request-transformer';
import { CandleRequest, CandleResponse } from '../../../types';

const router = express.Router();
const controllerName: string = 'candles-controller';
const candlesCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_CANDLES,
);

@Route('candles')
class CandleController extends Controller {
  @Get('/perpetualMarkets/:ticker')
  async getCandles(
    @Path() ticker: string,
      @Query() resolution: CandleResolution,
      @Query() limit?: number,
      @Query() fromISO?: string,
      @Query() toISO?: string,
  ): Promise<CandleResponse> {
    const candles: CandleFromDatabase[] = await CandleTable.findAll(
      {
        ticker: [ticker],
        resolution,
        fromISO,
        toISO,
        limit,
      },
      [],
    );

    return {
      candles: candles.map(candleToResponseObject),
    };
  }
}

router.get(
  '/perpetualMarkets/:ticker',
  rateLimiterMiddleware(candlesRateLimiter),
  candlesCacheControlMiddleware,
  ...CheckLimitSchema,
  ...CheckTickerParamSchema,
  ...checkSchema({
    resolution: {
      in: 'query',
      isString: true,
      isIn: {
        options: [Object.values(CandleResolution)],
      },
      errorMessage: `resolution must be a valid Candle Resolution, one of ${Object.values(CandleResolution)}`,
    },
    fromISO: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
    toISO: {
      in: 'query',
      optional: true,
      isISO8601: true,
    },
  }),
  handleValidationErrors,
  async (req: express.Request, res: express.Response) => {
    const {
      ticker,
      resolution,
      fromISO,
      toISO,
      limit,
    }: CandleRequest = matchedData(req) as CandleRequest;

    const start: number = Date.now();
    try {
      const controller: CandleController = new CandleController();
      const response: CandleResponse = await controller.getCandles(
        ticker,
        resolution,
        limit,
        fromISO,
        toISO,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'CandlesController GET /perpetualMarkets/:ticker',
        'Candles error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_candles_market.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
