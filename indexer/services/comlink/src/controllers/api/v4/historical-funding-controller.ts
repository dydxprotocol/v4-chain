import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  FundingIndexUpdatesColumns,
  FundingIndexUpdatesFromDatabase,
  FundingIndexUpdatesTable,
  IsoString,
  Ordering,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Path, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckEffectiveBeforeOrAtSchema, CheckLimitSchema, CheckTickerParamSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { historicalFundingToResponseObject } from '../../../request-helpers/request-transformer';
import {
  HistoricalFundingRequest,
  HistoricalFundingResponse,
  HistoricalFundingResponseObject,
  MarketType,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-funding-controller';
const historicalFundingCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_HISTORICAL_FUNDING,
);

@Route('historicalFunding')
class HistoricalFundingController extends Controller {
  @Get('/:ticker')
  async getHistoricalFunding(
    @Path() ticker: string,
      @Query() limit?: number,
      @Query() effectiveBeforeOrAtHeight?: number,
      @Query() effectiveBeforeOrAt?: IsoString,
  ): Promise<HistoricalFundingResponse> {
    const perpetualMarket: (
      PerpetualMarketFromDatabase | undefined
    ) = await PerpetualMarketTable.findByTicker(ticker);

    if (perpetualMarket === undefined) {
      throw new NotFoundError(`${ticker} not found in markets of type ${MarketType.PERPETUAL}`);
    }

    const fundingIndices: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        perpetualId: [perpetualMarket.id],
        effectiveBeforeOrAt,
        effectiveBeforeOrAtHeight: effectiveBeforeOrAtHeight
          ? effectiveBeforeOrAtHeight.toString()
          : undefined,
        limit,
      }, [],
      {
        ...DEFAULT_POSTGRES_OPTIONS,
        orderBy: [[FundingIndexUpdatesColumns.effectiveAtHeight, Ordering.DESC]],
      },
    );

    return {
      historicalFunding: fundingIndices.map(
        (fundingIndex: FundingIndexUpdatesFromDatabase): HistoricalFundingResponseObject => {
          return historicalFundingToResponseObject(fundingIndex, ticker);
        },
      ),
    };
  }
}

router.get(
  '/:ticker',
  rateLimiterMiddleware(defaultRateLimiter),
  historicalFundingCacheControlMiddleware,
  ...CheckLimitSchema,
  ...CheckTickerParamSchema,
  ...CheckEffectiveBeforeOrAtSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      ticker,
      limit,
      effectiveBeforeOrAtHeight,
      effectiveBeforeOrAt,
    }: {
      ticker: string,
      limit: number,
      effectiveBeforeOrAtHeight?: number,
      effectiveBeforeOrAt?: IsoString,
    } = matchedData(req) as HistoricalFundingRequest;

    try {
      const controller: HistoricalFundingController = new HistoricalFundingController();
      const response: HistoricalFundingResponse = await controller.getHistoricalFunding(
        ticker,
        limit,
        effectiveBeforeOrAtHeight,
        effectiveBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalFundingController GET /',
        'HistoricalFunding error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_funding.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
