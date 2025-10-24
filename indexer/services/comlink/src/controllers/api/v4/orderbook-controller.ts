import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import { PerpetualMarketFromDatabase, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, OrderbookLevelsCache } from '@dydxprotocol-indexer/redis';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Path, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { redisReadOnlyClient } from '../../../helpers/redis/redis-controller';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { OrderbookLevelsToResponseObject } from '../../../request-helpers/request-transformer';
import { MarketType, OrderbookRequest, OrderbookResponseObject } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'orderbook-controller';
const orderbookCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_ORDERBOOK,
);

@Route('orderbooks')
class OrderbookController extends Controller {
  @Get('/perpetualMarket/:ticker')
  async getPerpetualMarket(
    @Path() ticker: string,
  ): Promise<OrderbookResponseObject> {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromTicker(ticker);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(
        `${ticker} not found in markets of type ${MarketType.PERPETUAL}`,
      );
    }

    const orderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      ticker,
      redisReadOnlyClient,
      {
        sortSides: true,
        uncrossBook: true,
        limitPerSide: config.API_ORDERBOOK_LEVELS_PER_SIDE_LIMIT,
      },
    );

    return OrderbookLevelsToResponseObject(orderbookLevels, perpetualMarket);
  }
}

router.get(
  '/perpetualMarket/:ticker',
  rateLimiterMiddleware(defaultRateLimiter),
  orderbookCacheControlMiddleware,
  ...checkSchema({
    ticker: {
      in: ['params'],
      isString: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      ticker,
    }: {
      ticker: string,
    } = matchedData(req) as OrderbookRequest;

    try {
      const controller: OrderbookController = new OrderbookController();
      return res.send(await controller.getPerpetualMarket(ticker));
    } catch (error) {
      return handleControllerError(
        'OrderbooksController GET /perpetualMarket/:ticker',
        'Orderbooks error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_orderbooks.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
