import { stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  Ordering,
  PnlTicksFromDatabase,
  PnlTicksTable,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { complianceCheck } from '../../../lib/compliance-check';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { rejectRestrictedCountries } from '../../../lib/restrict-countries';
import {
  CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import { PnlTicksRequest, HistoricalPnlResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'historical-pnl-controller';

@Route('historical-pnl')
class HistoricalPnlController extends Controller {
  @Get('/')
  async getHistoricalPnl(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
      @Query() createdOnOrAfterHeight?: number,
      @Query() createdOnOrAfter?: IsoString,
  ): Promise<HistoricalPnlResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    const [subaccount, pnlTicks]: [
      SubaccountFromDatabase | undefined,
      PnlTicksFromDatabase[],
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
    };
  }
}

router.get(
  '/',
  rejectRestrictedCountries,
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtAndOnOrAfterSchema,
  handleValidationErrors,
  complianceCheck,
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

export default router;
