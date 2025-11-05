import { cacheControlMiddleware, stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  Ordering,
  PaginationFromDatabase,
  PnlFromDatabase,
  QueryableField,
  SubaccountTable,
  PnlTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller,
  Get,
  Query,
  Route,
} from 'tsoa';

import { pnlRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import { getChildSubaccountIds, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckDailyOptionalSchema,
  CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  CheckPaginationSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlToResponseObject } from '../../../request-helpers/request-transformer';
import { PnlResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'pnl-controller';
const pnlCacheControlMiddleware = cacheControlMiddleware(
  config.CACHE_CONTROL_DIRECTIVE_PNL,
);

@Route('pnl')
class PnlController extends Controller {
  @Get('/')
  async getPnl(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdOnOrAfterHeight?: number,
      @Query() createdOnOrAfter?: IsoString,
      @Query() page?: number,
      @Query() daily?: boolean,
  ): Promise<PnlResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    // First check if the subaccount exists
    const subaccount = await SubaccountTable.findById(subaccountId);

    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }

    // Set up common query parameters
    const queryParams = {
      subaccountId: [subaccountId],
      limit,
      createdBeforeOrAtHeight:
        createdBeforeOrAtHeight != null ? String(createdBeforeOrAtHeight) : undefined,
      createdBeforeOrAt,
      createdOnOrAfterHeight:
        createdOnOrAfterHeight != null ? String(createdOnOrAfterHeight) : undefined,
      createdOnOrAfter,
      page,
    };

    let pnlData: PaginationFromDatabase<PnlFromDatabase>;

    if (daily === true) {
      pnlData = await PnlTable.findAllDailyAggregate(
        queryParams,
        [QueryableField.LIMIT],
        DEFAULT_POSTGRES_OPTIONS,
      );
    } else {
      pnlData = await PnlTable.findAllHourlyAggregate(
        queryParams,
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[QueryableField.CREATED_AT_HEIGHT, Ordering.DESC]],
        },
      );
    }

    // Extract the results and pagination info
    const {
      results: pnlRecords, limit: pageSize, offset, total,
    } = pnlData;

    // Return the response
    return {
      pnl: pnlRecords.map((pnl: PnlFromDatabase) => {
        return pnlToResponseObject(pnl);
      }),
      pageSize,
      totalResults: total,
      offset,
    };
  }

  @Get('/parentSubaccountNumber')
  async getPnlForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdOnOrAfterHeight?: number,
      @Query() createdOnOrAfter?: IsoString,
      @Query() daily?: boolean,
  ): Promise<PnlResponse> {
    const childSubaccountIds: string[] = getChildSubaccountIds(address, parentSubaccountNumber);

    const queryParams = {
      subaccountId: childSubaccountIds,
      limit,
      createdBeforeOrAtHeight:
    createdBeforeOrAtHeight != null ? String(createdBeforeOrAtHeight) : undefined,
      createdBeforeOrAt,
      createdOnOrAfterHeight:
    createdOnOrAfterHeight != null ? String(createdOnOrAfterHeight) : undefined,
      createdOnOrAfter,
    };

    let pnlData: PaginationFromDatabase<PnlFromDatabase>;

    if (daily === true) {
      pnlData = await PnlTable.findAllDailyAggregate(
        queryParams,
        [QueryableField.LIMIT],
        DEFAULT_POSTGRES_OPTIONS,
      );
    } else {
      pnlData = await PnlTable.findAllHourlyAggregate(
        queryParams,
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[QueryableField.CREATED_AT_HEIGHT, Ordering.DESC]],
        },
      );
    }

    if (pnlData.results.length === 0) {
      throw new NotFoundError(
        `No PnL data found for address ${address} and parentSubaccountNumber ${parentSubaccountNumber}`,
      );
    }

    return {
      pnl: pnlData.results.map(pnlToResponseObject),
      totalResults: pnlData.total,
      offset: pnlData.offset,
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(pnlRateLimiter),
  pnlCacheControlMiddleware,
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckPaginationSchema,
  ...CheckDailyOptionalSchema,
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
      daily,
    } = matchedData(req) as {
      address: string,
      subaccountNumber: number,
      limit?: number,
      createdBeforeOrAtHeight?: number,
      createdBeforeOrAt?: IsoString,
      createdOnOrAfterHeight?: number,
      createdOnOrAfter?: IsoString,
      page?: number,
      daily?: boolean,
    };

    try {
      const controllers: PnlController = new PnlController();
      const response: PnlResponse = await controllers.getPnl(
        address,
        subaccountNumber,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        createdOnOrAfterHeight,
        createdOnOrAfter,
        page,
        daily,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PnlController GET /',
        'Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(pnlRateLimiter),
  pnlCacheControlMiddleware,
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  ...CheckDailyOptionalSchema,
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
      daily,
    } = matchedData(req) as {
      address: string,
      parentSubaccountNumber: number,
      limit?: number,
      createdBeforeOrAtHeight?: number,
      createdBeforeOrAt?: IsoString,
      createdOnOrAfterHeight?: number,
      createdOnOrAfter?: IsoString,
      daily?: boolean,
    };

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controllers: PnlController = new PnlController();
      const response: PnlResponse = await controllers.getPnlForParentSubaccount(
        address,
        parentSubaccountNum,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
        createdOnOrAfterHeight,
        createdOnOrAfter,
        daily,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'PnlController GET /parentSubaccountNumber',
        'Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_pnl_parent_subaccount.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
