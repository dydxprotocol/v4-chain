import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  liquidityTierRefresher,
  LiquidityTiersFromDatabase,
  LiquidityTiersMap,
  MarketFromDatabase,
  MarketTable,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  PerpetualMarketWithMarket,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  matchedData,
} from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { InvalidParamError, NotFoundError } from '../../../lib/errors';
import {
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckLimitSchema, CheckMarketOptionalQuerySchema, CheckTickerOptionalQuerySchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { perpetualMarketToResponseObject } from '../../../request-helpers/request-transformer';
import {
  MarketType,
  PerpetualMarketRequest,
  PerpetualMarketResponse,
  PerpetualMarketResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'perpetual-markets-controller';
const perpetualMarketsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_PERPETUAL_MARKETS,
);

@Route('perpetualMarkets')
class PerpetualMarketsController extends Controller {
  @Get('/')
  async listPerpetualMarkets(
    @Query() limit?: number,
      @Query() ticker?: string,
      @Query() market?: string,
  ): Promise<PerpetualMarketResponse> {
    const liquidityTiersMap: LiquidityTiersMap = liquidityTierRefresher.getLiquidityTiersMap();
    if (ticker && market) {
      throw new InvalidParamError('Only one of ticker or market may be provided');
    }

    const identifier = ticker || market;

    if (identifier !== undefined) {
      const perpetualMarket: (
        PerpetualMarketFromDatabase | undefined
      ) = await PerpetualMarketTable.findByTicker(identifier);

      if (perpetualMarket === undefined) {
        throw new NotFoundError(`${identifier} not found in markets of type ${MarketType.PERPETUAL}`);
      }

      const marketTable: (
        MarketFromDatabase | undefined
      ) = await MarketTable.findById(perpetualMarket.marketId);

      if (marketTable === undefined) {
        throw new NotFoundError(`Market not found for ticker ${identifier}`);
      }

      if (liquidityTiersMap[perpetualMarket.liquidityTierId] === undefined) {
        throw new NotFoundError(`Liquidity tier ${perpetualMarket.liquidityTierId} not found for ticker ${identifier}`);
      }

      return {
        markets: {
          [identifier]: perpetualMarketToResponseObject(
            perpetualMarket,
            liquidityTiersMap[perpetualMarket.liquidityTierId],
            marketTable,
          ),
        },
      };
    }

    const perpetualWithMarkets: PerpetualMarketWithMarket[] = await PerpetualMarketTable.findAll(
      {
        limit,
        joinWithMarkets: true,
      }, []) as PerpetualMarketWithMarket[];

    const liquidityTiers: LiquidityTiersFromDatabase[] = _.map(
      perpetualWithMarkets,
      (p) => {
        return liquidityTierRefresher.getLiquidityTierFromId(
          p.liquidityTierId,
        ) as LiquidityTiersFromDatabase;
      });

    const responseObjects: PerpetualMarketResponseObject[] = _.zipWith(
      perpetualWithMarkets,
      liquidityTiers,
      (pwm, lt) => {
        // Destructure each `perpetualWithMarket` into perpetual and market
        // to be able to use existing perpetualMarketToResponseObject function
        // for response transformation.
        const {
          pair,
          exponent,
          minPriceChangePpm,
          oraclePrice,
          ...perpetual
        } = pwm;

        return perpetualMarketToResponseObject(
          perpetual,
          lt,
          {
            id: pwm.marketId,
            pair,
            exponent,
            minPriceChangePpm,
            oraclePrice,
          },
        );
      },
    );

    return {
      markets: _.chain(responseObjects)
        .keyBy(PerpetualMarketColumns.ticker)
        .value(),
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(defaultRateLimiter),
  perpetualMarketsCacheControlMiddleware,
  ...CheckLimitSchema,
  ...CheckTickerOptionalQuerySchema,
  ...CheckMarketOptionalQuerySchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      limit,
      ticker,
      market,
    }: {
      limit: number,
      ticker?: string,
      market?: string,
    } = matchedData(req) as PerpetualMarketRequest;

    try {
      const controller: PerpetualMarketsController = new PerpetualMarketsController();
      const response: PerpetualMarketResponse = await controller.listPerpetualMarkets(
        limit,
        ticker,
        market,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PerpetualMarketController GET /',
        'PerpetualMarket error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_perpetual_markets.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
