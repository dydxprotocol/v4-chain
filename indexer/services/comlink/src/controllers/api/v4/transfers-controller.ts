import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetTable,
  DEFAULT_POSTGRES_OPTIONS,
  IsoString,
  Ordering,
  QueryableField,
  SubaccountColumns,
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
import { complianceAndGeoCheck } from '../../../lib/compliance-and-geo-check';
import { NotFoundError } from '../../../lib/errors';
import { getChildSubaccountNums, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import {
  CheckLimitAndCreatedBeforeOrAtSchema,
  CheckParentSubaccountSchema,
  CheckSubaccountSchema,
} from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  transferToParentSubaccountResponseObject,
  transferToResponseObject,
} from '../../../request-helpers/request-transformer';
import {
  AssetById,
  SubaccountById,
  TransferRequest,
  TransferResponse,
  ParentSubaccountTransferRequest,
  ParentSubaccountTransferResponse,
  ParentSubaccountTransferResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'transfers-controller';

@Route('transfers')
class TransfersController extends Controller {
  @Get('/')
  async getTransfers(
    @Query() address: string,
      @Query() subaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<TransferResponse> {
    const subaccountId: string = SubaccountTable.uuid(address, subaccountNumber);

    // TODO(DEC-656): Change to a cache in Redis similar to Librarian instead of querying DB.
    const [subaccount, transfers, assets]: [
      SubaccountFromDatabase | undefined,
      TransferFromDatabase[],
      AssetFromDatabase[]
    ] = await
    Promise.all([
      SubaccountTable.findById(
        subaccountId,
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
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]],
        },
      ),
      AssetTable.findAll(
        {},
        [],
      ),
    ]);
    if (subaccount === undefined) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and subaccountNumber ${subaccountNumber}`,
      );
    }
    const recipientSubaccountIds: string[] = _
      .map(transfers, TransferColumns.recipientSubaccountId)
      .filter(
        (recipientSubaccountId: string | undefined) => recipientSubaccountId !== undefined,
      ) as string[];
    const senderSubaccountIds: string[] = _
      .map(transfers, TransferColumns.senderSubaccountId)
      .filter(
        (senderSubaccountId: string | undefined) => senderSubaccountId !== undefined,
      ) as string[];

    const subaccountIds: string[] = _.uniq([
      ...recipientSubaccountIds,
      ...senderSubaccountIds,
    ]);
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        id: subaccountIds,
      },
      [],
    );
    const idToSubaccount: SubaccountById = _.keyBy(
      subaccounts,
      SubaccountColumns.id,
    );

    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    return {
      transfers: transfers.map((transfer: TransferFromDatabase) => {
        return transferToResponseObject(transfer, idToAsset, idToSubaccount, subaccountId);
      }),
    };
  }

  @Get('/parentSubaccountNumber')
  async getTransfersForParentSubaccount(
    @Query() address: string,
      @Query() parentSubaccountNumber: number,
      @Query() limit?: number,
      @Query() createdBeforeOrAtHeight?: number,
      @Query() createdBeforeOrAt?: IsoString,
  ): Promise<ParentSubaccountTransferResponse> {

    // get all child subaccountIds for the parent subaccount number
    const subaccountIds: string[] = getChildSubaccountNums(parentSubaccountNumber).map(
      (childSubaccountNumber: number) => SubaccountTable.uuid(address, childSubaccountNumber),
    );

    // TODO(DEC-656): Change to a cache in Redis similar to Librarian instead of querying DB.
    const [subaccounts, transfers, assets]: [
      SubaccountFromDatabase[] | undefined,
      TransferFromDatabase[],
      AssetFromDatabase[]
    ] = await
    Promise.all([
      SubaccountTable.findAll(
        { id: subaccountIds },
        [],
      ),
      TransferTable.findAllToOrFromSubaccountId(
        {
          subaccountId: subaccountIds,
          limit,
          createdBeforeOrAtHeight: createdBeforeOrAtHeight
            ? createdBeforeOrAtHeight.toString()
            : undefined,
          createdBeforeOrAt,
        },
        [QueryableField.LIMIT],
        {
          ...DEFAULT_POSTGRES_OPTIONS,
          orderBy: [[TransferColumns.createdAtHeight, Ordering.DESC]],
        },
      ),
      AssetTable.findAll(
        {},
        [],
      ),
    ]);
    if (subaccounts === undefined || subaccounts.length === 0) {
      throw new NotFoundError(
        `No subaccount found with address ${address} and parentSubaccountNumber ${parentSubaccountNumber}`,
      );
    }
    const recipientSubaccountIds: string[] = _
      .map(transfers, TransferColumns.recipientSubaccountId)
      .filter(
        (recipientSubaccountId: string | undefined) => recipientSubaccountId !== undefined,
      ) as string[];
    const senderSubaccountIds: string[] = _
      .map(transfers, TransferColumns.senderSubaccountId)
      .filter(
        (senderSubaccountId: string | undefined) => senderSubaccountId !== undefined,
      ) as string[];

    const allSubaccountIds: string[] = _.uniq([
      ...recipientSubaccountIds,
      ...senderSubaccountIds,
    ]);
    const allSubaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        id: allSubaccountIds,
      },
      [],
    );
    const idToSubaccount: SubaccountById = _.keyBy(
      allSubaccounts,
      SubaccountColumns.id,
    );

    const idToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    const transfersWithParentSubaccount: ParentSubaccountTransferResponseObject[] = transfers.map(
      (transfer: TransferFromDatabase) => {
        return transferToParentSubaccountResponseObject(
          transfer,
          idToAsset,
          idToSubaccount,
          parentSubaccountNumber);
      });

    // Filter out transfers where the sender and recipient parent subaccount numbers are the same
    const transfersFiltered:
    ParentSubaccountTransferResponseObject[] = transfersWithParentSubaccount.filter(
      (transfer) => {
        return transfer.sender.parentSubaccountNumber !== transfer.recipient.parentSubaccountNumber;
      });

    return { transfers: transfersFiltered };
  }
}

router.get(
  '/',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
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
        req,
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

router.get(
  '/parentSubaccountNumber',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckParentSubaccountSchema,
  ...CheckLimitAndCreatedBeforeOrAtSchema,
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
    }: ParentSubaccountTransferRequest = matchedData(req) as ParentSubaccountTransferRequest;

    // The schema checks allow subaccountNumber to be a string, but we know it's a number here.
    const parentSubaccountNum: number = +parentSubaccountNumber;

    try {
      const controllers: TransfersController = new TransfersController();
      const response: TransferResponse = await controllers.getTransfersForParentSubaccount(
        address,
        parentSubaccountNum,
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
        req,
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
