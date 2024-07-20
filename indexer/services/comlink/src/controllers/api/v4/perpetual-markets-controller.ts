import { stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  MarketTable,
  MarketFromDatabase,
  liquidityTierRefresher,
  LiquidityTiersMap,
  LiquidityTiersFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  matchedData,
} from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import {
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import { CheckLimitSchema, CheckTickerOptionalQuerySchema } from '../../../lib/validation/schemas';
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

@Route('perpetualMarkets')
class PerpetualMarketsController extends Controller {
  @Get('/')
  async listPerpetualMarkets(
    @Query() limit: number,
      @Query() ticker?: string,
  ): Promise<PerpetualMarketResponse> {
    const liquidityTiersMap: LiquidityTiersMap = liquidityTierRefresher.getLiquidityTiersMap();
    if (ticker !== undefined) {
      const perpetualMarket: (
        PerpetualMarketFromDatabase | undefined
      ) = await PerpetualMarketTable.findByTicker(ticker);

      if (perpetualMarket === undefined) {
        throw new NotFoundError(`${ticker} not found in markets of type ${MarketType.PERPETUAL}`);
      }

      const market: (
        MarketFromDatabase | undefined
      ) = await MarketTable.findById(perpetualMarket.marketId);

      if (market === undefined) {
        throw new NotFoundError(`Market not found for ticker ${ticker}`);
      }

      if (liquidityTiersMap[perpetualMarket.liquidityTierId] === undefined) {
        throw new NotFoundError(`Liquidity tier ${perpetualMarket.liquidityTierId} not found for ticker ${ticker}`);
      }

      return {
        markets: {
          [ticker]: perpetualMarketToResponseObject(
            perpetualMarket,
            liquidityTiersMap[perpetualMarket.liquidityTierId],
            market,
          ),
        },
      };
    }

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll({
      limit,
    }, []);

    const markets: MarketFromDatabase[] = await Promise.all(
      _.map(
        perpetualMarkets,
        async (perpetualMarket) => {
          return await MarketTable.findById(perpetualMarket.marketId) as MarketFromDatabase;
        }),
    );

    const liquidityTiers: LiquidityTiersFromDatabase[] = _.map(
      perpetualMarkets,
      (perpetualMarket) => {
        return liquidityTierRefresher.getLiquidityTierFromId(
          perpetualMarket.liquidityTierId,
        ) as LiquidityTiersFromDatabase;
      });

    const responseObjects: PerpetualMarketResponseObject[] = _.zipWith(
      perpetualMarkets, liquidityTiers, markets, perpetualMarketToResponseObject);

    return {
      markets: _.chain(responseObjects)
        .keyBy(PerpetualMarketColumns.ticker)
        .value(),
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitSchema,
  ...CheckTickerOptionalQuerySchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      limit,
      ticker,
    }: {
      limit: number,
      ticker?: string,
    } = matchedData(req) as PerpetualMarketRequest;

    try {
      const controller: PerpetualMarketsController = new PerpetualMarketsController();
      const response: PerpetualMarketResponse = await controller.listPerpetualMarkets(
        limit,
        ticker,
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
