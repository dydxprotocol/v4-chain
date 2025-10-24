import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  Ordering, PaginationFromDatabase,
  PnlTicksFromDatabase,
  PnlTicksTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { historicalPnlRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import { aggregateHourlyPnlTicks, getChildSubaccountIds, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  CheckPaginationSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import { HistoricalPnlResponse, ParentSubaccountPnlTicksRequest, PnlTicksRequest } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-pnl-controller';
const historicalPnlCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_HISTORICAL_PNL,
);

@Route('historical-pnl')
class HistoricalPnlController extends Controller {
  @Get('/')
  async getHistoricalPnl(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdOnOrAfterHeight?: number,
      @Query() createdOnOrAfter?: IsoString,
      @Query() page?: number,
  ): Promise<HistoricalPnlResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    const [subaccount,
      {
        results: pnlTicks,
        limit: pageSize,
        offset,
        total,
      },
    ]: [
      SubaccountFromDatabase | undefined,
      PaginationFromDatabase<PnlTicksFromDatabase>,
    ] = await Promise.all([
      SubaccountTable.findById(
        subaccountId,
      ),
      PnlTicksTable.findAll(
        {
          subaccountId: [subaccountId],
          limit,
          createdBeforeOrAtBlockHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
          createdOnOrAfterBlockHeight: createdOnOrAfterHeight
            ? createdOnOrAfterHeight.toString()
            : undefined,
          createdOnOrAfter,
          page,
        },
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.DESC]],
        },
      ),
    ]);
    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }

    return {
      historicalPnl: pnlTicks.map((pnlTick: PnlTicksFromDatabase) => {
        return pnlTicksToResponseObject(pnlTick);
      }),
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/parentSubaccountNumber')
  async getHistoricalPnlForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdOnOrAfterHeight?: number,
      @Query() createdOnOrAfter?: IsoString,
  ): Promise<HistoricalPnlResponse> {

    const childSubaccountIds: string[] = getChildSubaccountIds(address, parentSubaccountNumber);

    const [subaccounts,
      {
        results: pnlTicks,
      },
    ]: [
      SubaccountFromDatabase[],
      PaginationFromDatabase<PnlTicksFromDatabase>,
    ] = await Promise.all([
      SubaccountTable.findAll(
        {
          id: childSubaccountIds,
        },
        [QueryableField.ID],
      ),
      PnlTicksTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit,
          createdBeforeOrAtBlockHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
          createdOnOrAfterBlockHeight: createdOnOrAfterHeight
            ? createdOnOrAfterHeight.toString()
            : undefined,
          createdOnOrAfter,
        },
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.DESC]],
        },
      ),
    ]);

    if (subaccounts.length === 0) {
      throw new NotFoundError(
        `No subaccounts found with address ${address} and parentSubaccountNumber ${parentSubaccountNumber}`,
      );
    }

    // aggregate pnlTicks for all subaccounts grouped by blockHeight
    const aggregatedPnlTicks: PnlTicksFromDatabase[] = _.map(
      aggregateHourlyPnlTicks(pnlTicks),
      'pnlTick',
    );

    return {
      historicalPnl: aggregatedPnlTicks.map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        }),
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(historicalPnlRateLimiter),
  historicalPnlCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      page,
    }: PnlTicksRequest = matchedData(req) as PnlTicksRequest;

    try {
      const controllers: HistoricalPnlController = new HistoricalPnlController();
      const response: HistoricalPnlResponse = await controllers.getHistoricalPnl(
        address,
        subaccountNumber,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        createdOnOrAfterHeight,
        createdOnOrAfter,
        page,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalPnlController GET /',
        'Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(historicalPnlRateLimiter),
  historicalPnlCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  handleValidationErrors,
  complianceAndGeoCheck,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      parentSubaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
    }: ParentSubaccountPnlTicksRequest = matchedData(req) as ParentSubaccountPnlTicksRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controllers: HistoricalPnlController = new HistoricalPnlController();
      const response: HistoricalPnlResponse = await controllers.getHistoricalPnlForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        createdOnOrAfterHeight,
        createdOnOrAfter,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'HistoricalPnlController GET /parentSubaccountNumber',
        'Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_historical_pnl_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
