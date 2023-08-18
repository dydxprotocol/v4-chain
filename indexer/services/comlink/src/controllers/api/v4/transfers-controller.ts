import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetTable,
  IsoString,
  Ordering,
  QueryableField,
  SubaccountFromDatabase,
  SubaccountTable,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import { matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckLimitAndCreatedBeforeOrAtSchema, CheckSubaccountSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { transferToResponseObject } from '../../../request-helpers/request-transformer';
import { AssetById, TransferRequest, TransferResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'transfers-controller';

@Route('transfers')
class TransfersController extends Controller {
  @Get('/')
  async getTransfers(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<TransferResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    // TODO(DEC-656): Change to a cache in Redis similar to Librarian instead of querying DB.
    const [subaccount, transfers, assets] : [
      SubaccountFromDatabase | undefined,
      TransferFromDatabase[],
      AssetFromDatabase[]
    ] = await
    Promise.all([
      SubaccountTable.findById(
        subaccountId,
        { readReplica: true },
      ),
      TransferTable.findAllToOrFromSubaccountId(
        {
          subaccountId: [subaccountId],
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
        },
        [QueryableField.LIMIT],
        {
          readReplica: true,
          orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]],
        },
      ),
      AssetTable.findAll(
        {},
        [],
        { readReplica: true },
      ),
    ]);
    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }
    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    return {
      transfers: transfers.map((transfer: TransferFromDatabase) => {
        return transferToResponseObject(transfer, idToAsset);
      }),
    };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      address,
      subaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    }: TransferRequest = matchedData(req) as TransferRequest;

    try {
      const controllers: TransfersController = new TransfersController();
      const response: TransferResponse = await controllers.getTransfers(
        address,
        subaccountNumber,
        limit,
        createdBeforeOrAtHeight,
        createdBeforeOrAt,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TransfersController GET /',
        'Transfers error',
        error,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_transfers.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
