import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  IsoString,
  Ordering,
  TradingRewardAggregationColumns,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckHistoricalBlockTradingRewardsSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { tradingRewardAggregationToResponse } from '../../../request-helpers/request-transformer';
import { HistoricalTradingRewardAggregationRequest, HistoricalTradingRewardAggregationsResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-trading-reward-aggregations-controller';
const historicalTradingRewardAggregationsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_HISTORICAL_TRADING_REWARDS,
);

@Route('historicalTradingRewardAggregations')
class HistoricalTradingRewardAggregationsController extends Controller {
  @Get('/:address')
  async getAggregations(
    @Path() address: string,
      @Query() period: TradingRewardAggregationPeriod,
      @Query() limit?: number,
      @Query() startingBeforeOrAt?: IsoString,
      @Query() startingBeforeOrAtHeight?: string,
  ): Promise<HistoricalTradingRewardAggregationsResponse> {
    const tradingRewardAggregations:
    TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll({
      address,
      period,
      limit,
      startedAtBeforeOrAt: startingBeforeOrAt,
      startedAtHeightBeforeOrAt: startingBeforeOrAtHeight,
    }, [], { orderBy: [[TradingRewardAggregationColumns.startedAtHeight, Ordering.DESC]] });
    return {
      rewards: _.map(
        tradingRewardAggregations,
        tradingRewardAggregationToResponse,
      ),
    };
  }
}

router.get(
  '/:address',
  rateLimiterMiddleware(defaultRateLimiter),
  historicalTradingRewardAggregationsCacheControlMiddleware,
  ...CheckHistoricalBlockTradingRewardsSchema,
  ...checkSchema({
    period: {
      in: 'query',
      isString: true,
      isIn: {
        options: [Object.values(TradingRewardAggregationPeriod)],
      },
      errorMessage: `period must be a valid Trading Reward Aggregation Period, one of ${Object.values(TradingRewardAggregationPeriod)}`,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      limit,
      period,
      startingBeforeOrAt,
      startingBeforeOrAtHeight,
    }: HistoricalTradingRewardAggregationRequest = matchedData(
      req,
    ) as HistoricalTradingRewardAggregationRequest;

    try {
      const controller:
      HistoricalTradingRewardAggregationsController = new
      HistoricalTradingRewardAggregationsController();
      const response:
      HistoricalTradingRewardAggregationsResponse = await controller.getAggregations(
        address,
        period,
        limit,
        startingBeforeOrAt,
        startingBeforeOrAtHeight,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalTradingRewardAggregationsController GET /',
        'HistoricalTradingRewardAggregations error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_trading_reward_aggregations.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
