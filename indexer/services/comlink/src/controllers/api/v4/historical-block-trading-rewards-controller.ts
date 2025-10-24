import { stats, cacheControlMiddleware } from '@dydxprotocol-indexer/base';
import {
  IsoString,
  Ordering,
  TradingRewardColumns,
  TradingRewardFromDatabase,
  TradingRewardTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
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
import { tradingRewardToResponse } from '../../../request-helpers/request-transformer';
import { HistoricalBlockTradingRewardRequest as HistoricalBlockTradingRewardsRequest, HistoricalBlockTradingRewardsResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-block-trading-rewards-controller';
const historicalBlockTradingRewardsCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_HISTORICAL_TRADING_REWARDS,
);

@Route('historicalBlockTradingRewards')
class HistoricalBlockTradingRewardsController extends Controller {
  @Get('/:address')
  async getTradingRewards(
    @Path() address: string,
      @Query() limit?: number,
      @Query() startingBeforeOrAt?: IsoString,
      @Query() startingBeforeOrAtHeight?: string,
  ): Promise<HistoricalBlockTradingRewardsResponse> {
    const tradingRewardAggregations:
    TradingRewardFromDatabase[] = await TradingRewardTable.findAll({
      address,
      limit,
      blockTimeBeforeOrAt: startingBeforeOrAt,
      blockHeightBeforeOrAt: startingBeforeOrAtHeight,
    }, [], { orderBy: [[TradingRewardColumns.blockHeight, Ordering.DESC]] });

    return {
      rewards: _.map(
        tradingRewardAggregations,
        tradingRewardToResponse,
      ),
    };
  }
}

router.get(
  '/:address',
  rateLimiterMiddleware(defaultRateLimiter),
  historicalBlockTradingRewardsCacheControlMiddleware,
  ...CheckHistoricalBlockTradingRewardsSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      limit,
      startingBeforeOrAt,
      startingBeforeOrAtHeight,
    }: HistoricalBlockTradingRewardsRequest = matchedData(
      req,
    ) as HistoricalBlockTradingRewardsRequest;

    try {
      const controller:
      HistoricalBlockTradingRewardsController = new HistoricalBlockTradingRewardsController();
      const response: HistoricalBlockTradingRewardsResponse = await controller.getTradingRewards(
        address,
        limit,
        startingBeforeOrAt,
        startingBeforeOrAtHeight,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalBlockTradingRewardsController GET /',
        'HistoricalBlockTradingRewards error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_block_trading_reward.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
